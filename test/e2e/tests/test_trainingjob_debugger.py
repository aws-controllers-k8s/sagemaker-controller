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
    assert_tags_in_sync,
)
from e2e.replacement_values import REPLACEMENT_VALUES
from e2e.bootstrap_resources import get_bootstrap_resources
from e2e.common import config as cfg

RESOURCE_PLURAL = "trainingjobs"


@pytest.fixture(scope="function")
def xgboost_training_job_debugger():
    resource_name = random_suffix_name("xgboost-trainingjob-debugger", 32)
    replacements = REPLACEMENT_VALUES.copy()
    replacements["TRAINING_JOB_NAME"] = resource_name
    reference, _, resource = create_sagemaker_resource(
        resource_plural=RESOURCE_PLURAL,
        resource_name=resource_name,
        spec_file="xgboost_trainingjob_debugger",
        replacements=replacements,
    )
    if k8s.get_resource_arn(resource) is None:
        logging.error(
            f"ARN for this resource is None, resource status is: {resource['status']}"
        )
    assert resource is not None
    assert k8s.get_resource_arn(resource) is not None

    yield (reference, resource)

    if k8s.get_resource_exists(reference):
        _, deleted = k8s.delete_custom_resource(reference, 3, 10)
        assert deleted


def get_sagemaker_training_job(training_job_name: str):
    try:
        training_job = sagemaker_client().describe_training_job(
            TrainingJobName=training_job_name
        )
        return training_job
    except botocore.exceptions.ClientError as error:
        logging.error(
            f"SageMaker could not find a training debugger job with the name {training_job_name}. Error {error}"
        )
        return None


# TODO: Move to __init__.py
def get_training_sagemaker_status(training_job_name: str):
    training_sm_desc = get_sagemaker_training_job(training_job_name)
    return training_sm_desc["TrainingJobStatus"]


def get_training_resource_status(reference: k8s.CustomResourceReference):
    resource = k8s.get_resource(reference)
    assert "trainingJobStatus" in resource["status"]
    return resource["status"]["trainingJobStatus"]


def get_training_debugger_sagemaker_status(training_job_name: str):
    training_sm_desc = get_sagemaker_training_job(training_job_name)
    return training_sm_desc["DebugRuleEvaluationStatuses"][0]["RuleEvaluationStatus"]


def get_training_debugger_resource_status(reference: k8s.CustomResourceReference):
    resource = k8s.get_resource(reference)
    resource_status = resource["status"]["debugRuleEvaluationStatuses"][0][
        "ruleEvaluationStatus"
    ]
    assert resource_status is not None
    return resource_status


@service_marker
class TestTrainingDebuggerJob:
    def _wait_sagemaker_training_status(
        self,
        training_job_name,
        expected_status: str,
        wait_periods: int = 30,
        period_length: int = 30,
    ):
        return wait_for_status(
            expected_status,
            wait_periods,
            period_length,
            get_training_sagemaker_status,
            training_job_name,
        )

    def _wait_resource_training_status(
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
            get_training_resource_status,
            reference,
        )

    def _assert_training_status_in_sync(
        self, training_job_name, reference, expected_status
    ):
        assert (
            self._wait_sagemaker_training_status(training_job_name, expected_status)
            == self._wait_resource_training_status(reference, expected_status)
            == expected_status
        )

    def _wait_sagemaker_training_debugger_status(
        self,
        training_job_name,
        expected_status: str,
        wait_periods: int = 30,
        period_length: int = 30,
    ):
        return wait_for_status(
            expected_status,
            wait_periods,
            period_length,
            get_training_debugger_sagemaker_status,
            training_job_name,
        )

    def _wait_resource_training_debugger_status(
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
            get_training_debugger_resource_status,
            reference,
        )

    def _assert_training_debugger_status_in_sync(
        self, training_job_name, reference, expected_status
    ):
        assert (
            self._wait_sagemaker_training_debugger_status(
                training_job_name, expected_status
            )
            == self._wait_resource_training_debugger_status(reference, expected_status)
            == expected_status
        )

    def test_completed(self, xgboost_training_job_debugger):
        (reference, resource) = xgboost_training_job_debugger
        assert k8s.get_resource_exists(reference)

        training_job_name = resource["spec"].get("trainingJobName", None)
        assert training_job_name is not None

        training_job_desc = get_sagemaker_training_job(training_job_name)
        training_job_arn = training_job_desc["TrainingJobArn"]
        assert k8s.get_resource_arn(resource) == training_job_arn

        assert training_job_desc["TrainingJobStatus"] == cfg.JOB_STATUS_INPROGRESS
        assert k8s.wait_on_condition(reference, "ACK.ResourceSynced", "False")

        self._assert_training_status_in_sync(
            training_job_name, reference, cfg.JOB_STATUS_COMPLETED
        )
        # TODO: This test is failing
        assert k8s.wait_on_condition(reference, "ACK.ResourceSynced", "False")

        self._assert_training_debugger_status_in_sync(
            training_job_name, reference, cfg.DEBUGGERJOB_STATUS_COMPLETED
        )
        assert k8s.wait_on_condition(reference, "ACK.ResourceSynced", "True")

        resource_tags = resource["spec"].get("tags", None)
        assert_tags_in_sync(training_job_arn, resource_tags)

        # Check that you can delete a completed resource from k8s
        _, deleted = k8s.delete_custom_resource(reference, 3, 10)
        assert deleted is True
