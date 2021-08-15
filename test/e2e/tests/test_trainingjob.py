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
"""Integration tests for the SageMaker TrainingJob API.
"""

import pytest
import logging

from acktest.resources import random_suffix_name
from acktest.k8s import resource as k8s
from e2e import (
    service_marker,
    create_sagemaker_resource,
    get_sagemaker_training_job,
    assert_training_status_in_sync,
    assert_tags_in_sync,
)
from e2e.replacement_values import REPLACEMENT_VALUES
from e2e.bootstrap_resources import get_bootstrap_resources
from e2e.common import config as cfg

RESOURCE_PLURAL = "trainingjobs"


@pytest.fixture(scope="function")
def xgboost_training_job():
    resource_name = random_suffix_name("xgboost-trainingjob", 32)
    replacements = REPLACEMENT_VALUES.copy()
    replacements["TRAINING_JOB_NAME"] = resource_name
    reference, _, resource = create_sagemaker_resource(
        resource_plural=RESOURCE_PLURAL,
        resource_name=resource_name,
        spec_file="xgboost_trainingjob",
        replacements=replacements,
    )

    assert resource is not None
    if k8s.get_resource_arn(resource) is None:
        logging.error(
            f"ARN for this resource is None, resource status is: {resource['status']}"
        )
    assert k8s.get_resource_arn(resource) is not None

    yield (reference, resource)

    if k8s.get_resource_exists(reference):
        _, deleted = k8s.delete_custom_resource(reference, 3, 10)
        assert deleted

@pytest.mark.canary
@service_marker
class TestTrainingJob:
    def test_stopped(self, xgboost_training_job):
        (reference, resource) = xgboost_training_job
        assert k8s.get_resource_exists(reference)

        training_job_name = resource["spec"].get("trainingJobName", None)
        assert training_job_name is not None

        training_job_desc = get_sagemaker_training_job(training_job_name)

        assert k8s.get_resource_arn(resource) == training_job_desc["TrainingJobArn"]
        assert training_job_desc["TrainingJobStatus"] == cfg.JOB_STATUS_INPROGRESS
        assert k8s.wait_on_condition(reference, "ACK.ResourceSynced", "False")

        assert_training_status_in_sync(
            training_job_name, reference, cfg.JOB_STATUS_INPROGRESS
        )

        # Delete the k8s resource.
        _, deleted = k8s.delete_custom_resource(reference, 3, 10)
        assert deleted is True

        training_job_desc = get_sagemaker_training_job(training_job_name)
        assert training_job_desc["TrainingJobStatus"] in cfg.LIST_JOB_STATUS_STOPPED

    def test_completed(self, xgboost_training_job):
        (reference, resource) = xgboost_training_job
        assert k8s.get_resource_exists(reference)

        training_job_name = resource["spec"].get("trainingJobName", None)
        assert training_job_name is not None

        training_job_desc = get_sagemaker_training_job(training_job_name)
        training_job_arn = training_job_desc["TrainingJobArn"]
        assert k8s.get_resource_arn(resource) == training_job_arn

        assert training_job_desc["TrainingJobStatus"] == cfg.JOB_STATUS_INPROGRESS
        assert k8s.wait_on_condition(reference, "ACK.ResourceSynced", "False")

        assert_training_status_in_sync(
            training_job_name, reference, cfg.JOB_STATUS_COMPLETED
        )
        assert k8s.wait_on_condition(reference, "ACK.ResourceSynced", "True")

        # model artifact URL is populated
        resource = k8s.get_resource(reference)
        resource["status"]["modelArtifacts"]["s3ModelArtifacts"] is not None

        resource_tags = resource["spec"].get("tags", None)
        assert_tags_in_sync(training_job_arn, resource_tags)

        # Check that you can delete a completed resource from k8s
        _, deleted = k8s.delete_custom_resource(reference, 3, 10)
        assert deleted is True
