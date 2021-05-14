# Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"). You may
# not use this file except in compliance with the License. A copy of the
# License is located at
#
# 	 http://aws.amazon.com/apache2.0/
#
# or in the "license" file accompanying this file. This file is distributed
# on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
# express or implied. See the License for the specific language governing
# permissions and limitations under the License.
"""Integration tests for the SageMaker ProcessingJob API.
"""

import botocore
import pytest
import logging
from typing import Dict

from acktest.resources import random_suffix_name
from acktest.k8s import resource as k8s
from e2e import (
    service_marker,
    create_sagemaker_resource,
    wait_for_status,
    sagemaker_client,
)
from e2e.replacement_values import REPLACEMENT_VALUES
from e2e.bootstrap_resources import get_bootstrap_resources
from e2e.common import config as cfg
from time import sleep

RESOURCE_PLURAL = "processingjobs"


@pytest.fixture(scope="function")
def kmeans_processing_job():
    resource_name = random_suffix_name("kmeans-processingjob", 32)
    replacements = REPLACEMENT_VALUES.copy()
    replacements["PROCESSING_JOB_NAME"] = resource_name
    reference, _, resource = create_sagemaker_resource(
        resource_plural=RESOURCE_PLURAL,
        resource_name=resource_name,
        spec_file="kmeans_processingjob",
        replacements=replacements,
    )

    assert resource is not None
    assert k8s.get_resource_arn(resource) is not None

    yield (reference, resource)

    if k8s.get_resource_exists(reference):
        _, deleted = k8s.delete_custom_resource(reference, 3, 10)
        assert deleted


def get_sagemaker_processing_job(processing_job_name: str):
    try:
        processing_job = sagemaker_client().describe_processing_job(
            ProcessingJobName=processing_job_name
        )
        return processing_job
    except botocore.exceptions.ClientError as error:
        logging.error(
            f"SageMaker could not find a processing job with the name {processing_job_name}. Error {error}"
        )
        return None


def get_processing_sagemaker_status(processing_job_name: str):
    processing_sm_desc = get_sagemaker_processing_job(processing_job_name)
    return processing_sm_desc["ProcessingJobStatus"]


def get_processing_resource_status(reference: k8s.CustomResourceReference):
    resource = k8s.get_resource(reference)
    assert "processingJobStatus" in resource["status"]
    return resource["status"]["processingJobStatus"]


@service_marker
@pytest.mark.canary
class TestProcessingJob:
    def _wait_resource_processing_status(
        self,
        reference: k8s.CustomResourceReference,
        expected_status: str,
        wait_periods: int = 30,
        period_length: int = 30,
    ):
        return wait_for_status(
            expected_status,
            wait_periods,
            period_length,
            get_processing_resource_status,
            reference,
        )

    def _wait_sagemaker_processing_status(
        self,
        processing_job_name,
        expected_status: str,
        wait_periods: int = 30,
        period_length: int = 30,
    ):
        return wait_for_status(
            expected_status,
            wait_periods,
            period_length,
            get_processing_sagemaker_status,
            processing_job_name,
        )

    def _assert_processing_status_in_sync(
        self, processing_job_name, reference, expected_status
    ):
        assert (
            self._wait_sagemaker_processing_status(processing_job_name, expected_status)
            == self._wait_resource_processing_status(reference, expected_status)
            == expected_status
        )

    def test_stopped(self, kmeans_processing_job):
        (reference, resource) = kmeans_processing_job
        assert k8s.get_resource_exists(reference)

        processing_job_name = resource["spec"].get("processingJobName", None)
        assert processing_job_name is not None

        processing_job_desc = get_sagemaker_processing_job(processing_job_name)

        assert k8s.get_resource_arn(resource) == processing_job_desc["ProcessingJobArn"]
        assert processing_job_desc["ProcessingJobStatus"] == cfg.JOB_STATUS_INPROGRESS
        assert k8s.wait_on_condition(reference, "ACK.ResourceSynced", "False")

        self._assert_processing_status_in_sync(
            processing_job_name, reference, cfg.JOB_STATUS_INPROGRESS
        )

        # Delete the k8s resource.
        _, deleted = k8s.delete_custom_resource(reference, 3, 10)
        assert deleted is True

        processing_job_desc = get_sagemaker_processing_job(processing_job_name)
        assert processing_job_desc["ProcessingJobStatus"] in cfg.LIST_JOB_STATUS_STOPPED

    def test_completed(self, kmeans_processing_job):
        (reference, resource) = kmeans_processing_job
        assert k8s.get_resource_exists(reference)

        processing_job_name = resource["spec"].get("processingJobName", None)
        assert processing_job_name is not None

        processing_job_desc = get_sagemaker_processing_job(processing_job_name)

        assert k8s.get_resource_arn(resource) == processing_job_desc["ProcessingJobArn"]
        assert processing_job_desc["ProcessingJobStatus"] == cfg.JOB_STATUS_INPROGRESS
        assert k8s.wait_on_condition(reference, "ACK.ResourceSynced", "False")

        self._assert_processing_status_in_sync(
            processing_job_name, reference, cfg.JOB_STATUS_COMPLETED
        )
        assert k8s.wait_on_condition(reference, "ACK.ResourceSynced", "True")

        # Check that you can delete a completed resource from k8s
        _, deleted = k8s.delete_custom_resource(reference, 3, 10)
        assert deleted is True
