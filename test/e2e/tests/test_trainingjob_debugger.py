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
    wait_for_status,
    get_sagemaker_training_job,
    assert_training_status_in_sync,
    assert_tags_in_sync,
)
from e2e.replacement_values import REPLACEMENT_VALUES
from e2e.common import config as cfg

RESOURCE_PLURAL = "trainingjobs"


@pytest.fixture(scope="function")
def xgboost_training_job_debugger():
    resource_name = random_suffix_name("xgboost-trainingjob-debugger", 50)
    replacements = REPLACEMENT_VALUES.copy()
    replacements["TRAINING_JOB_NAME"] = resource_name
    reference, _, resource = create_sagemaker_resource(
        resource_plural=RESOURCE_PLURAL,
        resource_name=resource_name,
        spec_file="xgboost_trainingjob_debugger",
        replacements=replacements,
    )
    assert resource is not None

    yield (reference, resource)

    if k8s.get_resource_exists(reference):
        _, deleted = k8s.delete_custom_resource(reference, 3, 10)
        assert deleted


def get_training_rule_eval_sagemaker_status(training_job_name: str, rule_type: str):
    training_sm_desc = get_sagemaker_training_job(training_job_name)
    return training_sm_desc[rule_type+"EvaluationStatuses"][0]["RuleEvaluationStatus"]


def get_training_rule_eval_resource_status(reference: k8s.CustomResourceReference, rule_type: str):
    resource = k8s.get_resource(reference)
    resource_status = resource["status"][rule_type+"EvaluationStatuses"][0][
        "ruleEvaluationStatus"
    ]
    assert resource_status is not None
    return resource_status

@service_marker
class TestTrainingDebuggerJob:
    def _wait_sagemaker_training_rule_eval_status(
        self,
        training_job_name,
        rule_type: str,
        expected_status: str,
        wait_periods: int = 30,
        period_length: int = 30,
    ):
        return wait_for_status(
            expected_status,
            wait_periods,
            period_length,
            get_training_rule_eval_sagemaker_status,
            training_job_name,
            rule_type,
        )

    def _wait_resource_training_rule_eval_status(
        self,
        reference: k8s.CustomResourceReference,
        rule_type: str,
        expected_status: str,
        wait_periods: int = 30,
        period_length: int = 30,
    ):
        return wait_for_status(
            expected_status,
            wait_periods,
            period_length,
            get_training_rule_eval_resource_status,
            reference,
            rule_type,
        )

    def _assert_training_rule_eval_status_in_sync(
        self, training_job_name, sagemaker_rule_type, reference, expected_status
    ):
        resource_rule_type = sagemaker_rule_type[0].lower() + sagemaker_rule_type[1:]
        assert (
            self._wait_sagemaker_training_rule_eval_status(
                training_job_name, sagemaker_rule_type, expected_status, 
            )
            == self._wait_resource_training_rule_eval_status(reference, resource_rule_type, expected_status)
            == expected_status
        )

    def test_completed(self, xgboost_training_job_debugger):
        (reference, resource) = xgboost_training_job_debugger
        assert k8s.get_resource_exists(reference)

        training_job_name = resource["spec"].get("trainingJobName", None)
        assert training_job_name is not None

        training_job_desc = get_sagemaker_training_job(training_job_name)
        training_job_arn = training_job_desc["TrainingJobArn"]
        
        resource_arn = k8s.get_resource_arn(resource)
        if resource_arn is None:
            logging.error(
                f"ARN for this resource is None, resource status is: {resource['status']}"
            )
        assert resource_arn == training_job_arn

        assert training_job_desc["TrainingJobStatus"] == cfg.JOB_STATUS_INPROGRESS
        assert k8s.wait_on_condition(reference, "ACK.ResourceSynced", "False")

        assert_training_status_in_sync(
            training_job_name, reference, cfg.JOB_STATUS_COMPLETED
        )
        assert k8s.wait_on_condition(reference, "ACK.ResourceSynced", "False")

        # Assert debugger rule evaluation completed
        self._assert_training_rule_eval_status_in_sync(
            training_job_name, "DebugRule", reference, cfg.RULE_STATUS_COMPLETED
        )
        
        # Assert profiler rule evaluation completed
        self._assert_training_rule_eval_status_in_sync(
            training_job_name, "ProfilerRule", reference, cfg.RULE_STATUS_COMPLETED
        )
        assert k8s.wait_on_condition(reference, "ACK.ResourceSynced", "True")

        resource_tags = resource["spec"].get("tags", None)
        assert_tags_in_sync(training_job_arn, resource_tags)

        # Check that you can delete a completed resource from k8s
        _, deleted = k8s.delete_custom_resource(reference, 3, 10)
        assert deleted is True
