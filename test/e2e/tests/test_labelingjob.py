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
"""Integration tests for the SageMaker LabelingJob API.
"""

import botocore
import pytest
import logging

from acktest.resources import random_suffix_name
from acktest.k8s import resource as k8s
from acktest.k8s import condition as ack_condition

from e2e import (
    service_marker,
    create_sagemaker_resource,
    delete_custom_resource,
    wait_for_status,
    sagemaker_client,
)
from e2e.replacement_values import REPLACEMENT_VALUES
from e2e.common import config as cfg

RESOURCE_PLURAL = "labelingjobs"


@pytest.fixture(scope="function")
def image_labeling_job():
    resource_name = random_suffix_name("image-labelingjob", 32)
    replacements = REPLACEMENT_VALUES.copy()
    replacements["LABELING_JOB_NAME"] = resource_name
    reference, _, resource = create_sagemaker_resource(
        resource_plural=RESOURCE_PLURAL,
        resource_name=resource_name,
        spec_file="image_labelingjob",
        replacements=replacements,
    )

    assert resource is not None
    if k8s.get_resource_arn(resource) is None:
        logging.error(f"ARN for this resource is None, resource status is: {resource['status']}")
    assert k8s.get_resource_arn(resource) is not None

    yield (reference, resource)

    assert delete_custom_resource(
        reference, cfg.JOB_DELETE_WAIT_PERIODS, cfg.JOB_DELETE_WAIT_LENGTH
    )


def get_sagemaker_labeling_job(labeling_job_name: str):
    try:
        labeling_job = sagemaker_client().describe_labeling_job(LabelingJobName=labeling_job_name)
        return labeling_job
    except botocore.exceptions.ClientError as error:
        logging.error(
            f"SageMaker could not find a labeling job with the name {labeling_job_name}. Error {error}"
        )
        return None


def get_labeling_sagemaker_status(labeling_job_name: str):
    labeling_sm_desc = get_sagemaker_labeling_job(labeling_job_name)
    return labeling_sm_desc["LabelingJobStatus"]


def get_labeling_resource_status(reference: k8s.CustomResourceReference):
    resource = k8s.get_resource(reference)
    assert "labelingJobStatus" in resource["status"]
    return resource["status"]["labelingJobStatus"]


@service_marker
@pytest.skip
class TestLabelingJob:
    def _wait_resource_labeling_status(
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
            get_labeling_resource_status,
            reference,
        )

    def _wait_sagemaker_labeling_status(
        self,
        labeling_job_name,
        expected_status: str,
        wait_periods: int = 30,
        period_length: int = 30,
    ):
        return wait_for_status(
            expected_status,
            wait_periods,
            period_length,
            get_labeling_sagemaker_status,
            labeling_job_name,
        )

    def _assert_labeling_status_in_sync(self, labeling_job_name, reference, expected_status):
        assert (
            self._wait_sagemaker_labeling_status(labeling_job_name, expected_status)
            == self._wait_resource_labeling_status(reference, expected_status)
            == expected_status
        )

    # NOTE: We are testing only test_stopped operation without test_completed
    #       This is due to the nature of LabelingJob where human intervention
    #       is required to bring job to the completion state.
    def test_stopped(self, image_labeling_job):
        (reference, resource) = image_labeling_job
        assert k8s.get_resource_exists(reference)

        labeling_job_name = resource["spec"].get("labelingJobName", None)
        assert labeling_job_name is not None

        labeling_job_desc = get_sagemaker_labeling_job(labeling_job_name)

        assert k8s.get_resource_arn(resource) == labeling_job_desc["LabelingJobArn"]
        assert labeling_job_desc["LabelingJobStatus"] == cfg.JOB_STATUS_INPROGRESS
        assert k8s.wait_on_condition(
            reference, ack_condition.CONDITION_TYPE_RESOURCE_SYNCED, "False"
        )

        self._assert_labeling_status_in_sync(
            labeling_job_name, reference, cfg.JOB_STATUS_INPROGRESS
        )

        # Ensure resource is not in terminal state
        assert k8s.get_resource_condition(reference, ack_condition.CONDITION_TYPE_TERMINAL) is None
        # Delete the k8s resource.
        assert delete_custom_resource(
            reference, cfg.JOB_DELETE_WAIT_PERIODS, cfg.JOB_DELETE_WAIT_LENGTH
        )

        labeling_job_desc = get_sagemaker_labeling_job(labeling_job_name)
        assert labeling_job_desc["LabelingJobStatus"] in cfg.LIST_JOB_STATUS_STOPPED
