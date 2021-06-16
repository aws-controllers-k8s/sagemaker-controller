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
"""Integration tests for the SageMaker TransformJob API.
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

RESOURCE_PLURAL = "transformjobs"


@pytest.fixture(scope="function")
def generate_job_names():
    transform_resource_name = random_suffix_name("xgboost-transformjob", 27)
    model_resource_name = "model" + transform_resource_name

    yield (transform_resource_name, model_resource_name)


@pytest.fixture(scope="function")
def xgboost_model_for_transform(generate_job_names):
    (transform_resource_name, model_resource_name) = generate_job_names
    replacements = REPLACEMENT_VALUES.copy()
    replacements["MODEL_NAME"] = model_resource_name

    reference, _, resource = create_sagemaker_resource(
        resource_plural=cfg.MODEL_RESOURCE_PLURAL,
        resource_name=model_resource_name,
        spec_file="xgboost_model",
        replacements=replacements,
    )
    assert resource is not None
    assert k8s.get_resource_arn(resource) is not None

    yield (transform_resource_name, model_resource_name)

    if k8s.get_resource_exists(reference):
        _, deleted = k8s.delete_custom_resource(reference, 3, 10)
        assert deleted


# TODO: This method can also move to a common file.
@pytest.fixture(scope="function")
def xgboost_transformjob(xgboost_model_for_transform):
    (transform_resource_name, model_resource_name) = xgboost_model_for_transform
    replacements = REPLACEMENT_VALUES.copy()
    replacements["MODEL_NAME"] = model_resource_name
    replacements["TRANSFORM_JOB_NAME"] = transform_resource_name

    reference, _, resource = create_sagemaker_resource(
        resource_plural=RESOURCE_PLURAL,
        resource_name=transform_resource_name,
        spec_file="xgboost_transformjob",
        replacements=replacements,
    )

    assert resource is not None
    assert k8s.get_resource_arn(resource) is not None

    yield (reference, resource)

    if k8s.get_resource_exists(reference):
        _, deleted = k8s.delete_custom_resource(reference, 3, 10)
        assert deleted


def get_sagemaker_transform_job(transform_job_name: str):
    try:
        transform_desc = sagemaker_client().describe_transform_job(
            TransformJobName=transform_job_name
        )
        return transform_desc
    except botocore.exceptions.ClientError as error:
        logging.error(
            f"SageMaker could not find a transform job with the name {transform_job_name}. Error {error}"
        )
        return None


def get_transform_sagemaker_status(transform_job_name: str):
    transform_sm_desc = get_sagemaker_transform_job(transform_job_name)
    return transform_sm_desc["TransformJobStatus"]


def get_transform_resource_status(reference: k8s.CustomResourceReference):
    resource = k8s.get_resource(reference)
    assert "transformJobStatus" in resource["status"]
    return resource["status"]["transformJobStatus"]


@service_marker
@pytest.mark.canary
class TestTransformJob:
    def _wait_resource_transform_status(
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
            get_transform_resource_status,
            reference,
        )

    def _wait_sagemaker_transform_status(
        self,
        transform_job_name,
        expected_status: str,
        wait_periods: int = 30,
        period_length: int = 30,
    ):
        return wait_for_status(
            expected_status,
            wait_periods,
            period_length,
            get_transform_sagemaker_status,
            transform_job_name,
        )

    def _assert_transform_status_in_sync(
        self, transform_job_name, reference, expected_status
    ):
        assert (
            self._wait_sagemaker_transform_status(transform_job_name, expected_status)
            == self._wait_resource_transform_status(reference, expected_status)
            == expected_status
        )

    def test_stopped(self, xgboost_transformjob):
        (reference, resource) = xgboost_transformjob
        assert k8s.get_resource_exists(reference)

        transform_job_name = resource["spec"].get("transformJobName", None)
        assert transform_job_name is not None

        transform_sm_desc = get_sagemaker_transform_job(transform_job_name)
        assert k8s.get_resource_arn(resource) == transform_sm_desc["TransformJobArn"]
        assert transform_sm_desc["TransformJobStatus"] == cfg.JOB_STATUS_INPROGRESS
        assert k8s.wait_on_condition(reference, "ACK.ResourceSynced", "False")

        self._assert_transform_status_in_sync(
            transform_job_name, reference, cfg.JOB_STATUS_INPROGRESS
        )

        # Delete the k8s resource.
        _, deleted = k8s.delete_custom_resource(reference, 3, 10)
        assert deleted is True

        transform_sm_desc = get_sagemaker_transform_job(transform_job_name)
        assert transform_sm_desc["TransformJobStatus"] in cfg.LIST_JOB_STATUS_STOPPED

    def test_completed(self, xgboost_transformjob):
        (reference, resource) = xgboost_transformjob
        assert k8s.get_resource_exists(reference)

        transform_job_name = resource["spec"].get("transformJobName", None)
        assert transform_job_name is not None

        transform_sm_desc = get_sagemaker_transform_job(transform_job_name)
        assert k8s.get_resource_arn(resource) == transform_sm_desc["TransformJobArn"]
        assert transform_sm_desc["TransformJobStatus"] == cfg.JOB_STATUS_INPROGRESS
        assert k8s.wait_on_condition(reference, "ACK.ResourceSynced", "False")

        self._assert_transform_status_in_sync(
            transform_job_name, reference, cfg.JOB_STATUS_COMPLETED
        )
        assert k8s.wait_on_condition(reference, "ACK.ResourceSynced", "True")

        # Check that you can delete a completed resource from k8s
        _, deleted = k8s.delete_custom_resource(reference, 3, 10)
        assert deleted is True
