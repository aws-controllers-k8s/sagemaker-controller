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
"""Integration tests for the SageMaker pipeline API.
"""

import pytest
import logging

from acktest.resources import random_suffix_name
from acktest.k8s import resource as k8s
from acktest.k8s import condition as ack_condition

from e2e import (
    service_marker,
    wait_for_status,
    create_sagemaker_resource,
    delete_custom_resource,
    get_sagemaker_pipeline,
    assert_tags_in_sync,
)
from e2e.replacement_values import REPLACEMENT_VALUES
from e2e.common import config as cfg

RESOURCE_PLURAL = "pipelines"

DELETE_WAIT_PERIOD = 20
DELETE_WAIT_LENGTH = 30


@pytest.fixture(scope="function")
def pipeline():
    resource_name = random_suffix_name("pipeline", 28)
    replacements = REPLACEMENT_VALUES.copy()
    replacements["PIPELINE_NAME"] = resource_name
    (
        pipeline_reference,
        pipeline_spec,
        pipeline_resource,
    ) = create_sagemaker_resource(
        resource_plural=RESOURCE_PLURAL,
        resource_name=resource_name,
        spec_file="pipeline",
        replacements=replacements,
    )
    assert pipeline_resource is not None
    if k8s.get_resource_arn(pipeline_resource) is None:
        logging.error(
            f"ARN for this resource is None, resource status is: {pipeline_resource['status']}"
        )
    assert k8s.get_resource_arn(pipeline_resource) is not None

    yield (
        pipeline_reference,
        pipeline_spec,
        pipeline_resource,
    )

    # Delete the k8s resource if not already deleted by tests
    assert delete_custom_resource(
        pipeline_reference, cfg.JOB_DELETE_WAIT_PERIODS, cfg.JOB_DELETE_WAIT_LENGTH
    )


def get_sagemaker_pipeline_status(pipeline_arn: str):
    sm_pipeline_desc = get_sagemaker_pipeline(pipeline_arn)
    return sm_pipeline_desc["PipelineStatus"]


def get_pipeline_resource_status(reference: k8s.CustomResourceReference):
    resource = k8s.get_resource(reference)
    assert "pipelineStatus" in resource["status"]
    return resource["status"]["pipelineStatus"]


@pytest.mark.canary
@service_marker
class TestPipeline:
    def _wait_resource_pipeline_status(
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
            get_pipeline_resource_status,
            reference,
        )

    def _wait_sagemaker_pipeline_status(
        self,
        pipeline_arn,
        expected_status: str,
        wait_periods: int = 30,
        period_length: int = 30,
    ):
        return wait_for_status(
            expected_status,
            wait_periods,
            period_length,
            get_sagemaker_pipeline_status,
            pipeline_arn,
        )

    def _assert_pipeline_status_in_sync(self, pipeline_arn, reference, expected_status):
        assert (
            self._wait_sagemaker_pipeline_status(pipeline_arn, expected_status)
            == self._wait_resource_pipeline_status(reference, expected_status)
            == expected_status
        )

    def test_pipeline_succeeded(self, pipeline):
        (reference, spec, resource) = pipeline
        assert k8s.get_resource_exists(reference)

        pipeline_name = resource["spec"].get("pipelineName")
        # Need PipelineArn to reference the resource

        pipeline_desc = get_sagemaker_pipeline(pipeline_name=pipeline_name)
        if k8s.get_resource_arn(resource) is None:
            logging.error(
                f"ARN for this resource is None, resource status is: {resource['status']}"
            )

        pipeline_arn = pipeline_desc["PipelineArn"]
        old_pipeline_last_modified_time = pipeline_desc["LastModifiedTime"]
        assert k8s.get_resource_arn(resource) == pipeline_arn

        self._assert_pipeline_status_in_sync(pipeline_arn, reference, "Active")
        assert k8s.wait_on_condition(reference, ack_condition.CONDITION_TYPE_RESOURCE_SYNCED, "True")

        # Update the resource
        new_pipeline_display_name = random_suffix_name("updated-display-name", 38)
        spec["spec"]["pipelineDisplayName"] = new_pipeline_display_name
        resource = k8s.patch_custom_resource(reference, spec)
        resource = k8s.wait_resource_consumed_by_controller(reference)
        assert resource is not None

        self._assert_pipeline_status_in_sync(pipeline_arn, reference, "Active")
        assert k8s.wait_on_condition(reference, ack_condition.CONDITION_TYPE_RESOURCE_SYNCED, "True")

        pipeline_desc = get_sagemaker_pipeline(pipeline_name)

        assert pipeline_desc["PipelineDisplayName"] == new_pipeline_display_name
        assert resource["spec"].get("pipelineDisplayName", None) == new_pipeline_display_name
        assert old_pipeline_last_modified_time != pipeline_desc["LastModifiedTime"]
        assert resource["status"].get("lastModifiedTime") != old_pipeline_last_modified_time

        resource_tags = resource["spec"].get("tags", None)
        assert_tags_in_sync(pipeline_arn, resource_tags)

        # Check that you can delete a completed resource from k8s
        assert delete_custom_resource(reference, DELETE_WAIT_PERIOD, DELETE_WAIT_LENGTH)
