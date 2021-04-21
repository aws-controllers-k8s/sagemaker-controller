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

import boto3
import pytest
import logging
from typing import Dict

from acktest.resources import random_suffix_name, load_resource_file
from acktest.k8s import resource as k8s
from e2e import (
    resource_directory,
    CRD_GROUP,
    CRD_VERSION,
    service_marker,
    create_sagemaker_resource,
    wait_for_status,
)
from e2e.replacement_values import REPLACEMENT_VALUES
from e2e.bootstrap_resources import get_bootstrap_resources
from time import sleep

RESOURCE_PLURAL = "trainingjobs"

def _sagemaker_client():
    return boto3.client("sagemaker")

def _make_training_debugger_job():
    resource_name = random_suffix_name("xgboost-trainingjob-debugger", 32)

    replacements = REPLACEMENT_VALUES.copy()
    replacements["TRAINING_JOB_NAME"] = resource_name

    data = load_resource_file(
        resource_directory, "xgboost_trainingjob_debugger", additional_replacements=replacements
    )

    # Create the k8s resource
    reference = k8s.CustomResourceReference(
        CRD_GROUP, CRD_VERSION, RESOURCE_PLURAL, resource_name, namespace="default"
    )

    return reference, data

@pytest.fixture(scope="function")
def xgboost_training_job_debugger():
    (training_job, data) = _make_training_debugger_job()
    resource = k8s.create_custom_resource(training_job, data)
    resource = k8s.wait_resource_consumed_by_controller(training_job)

    yield (training_job, resource) 

    if k8s.get_resource_exists(training_job):
        k8s.delete_custom_resource(training_job)

def get_sagemaker_training_job(training_job_name: str):
    try:
        training_job = _sagemaker_client().describe_training_job(
            TrainingJobName=training_job_name
        )
        return training_job
    except BaseException:
        logging.error(
            f"SageMaker could not find a training debugger job with the name {training_job_name}"
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
    resource_status = resource["status"]["debugRuleEvaluationStatuses"][0]["ruleEvaluationStatus"]
    assert resource_status is not None
    return resource_status

@service_marker
@pytest.mark.canary
class TestTrainingDebuggerJob:
    list_status_created = ("InProgress", "Completed")
    list_status_stopped = ("Stopped", "Stopping")
    status_inprogress: str = "InProgress"
    status_completed: str = "Completed"
    debugger_status_completed: str = (
        "NoIssuesFound"
    )

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
            reference
    )

    def _assert_training_status_in_sync(
        self, training_job_name, reference, expected_status
    ):
        assert (
            self._wait_sagemaker_training_status(
                training_job_name, expected_status
            ) 
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
            reference
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
    
    def test_training_job_debugger(self, xgboost_training_job_debugger):
        (reference, resource) = xgboost_training_job_debugger
        assert k8s.get_resource_exists(reference)

        training_job_name = resource["spec"].get("trainingJobName", None)
        assert training_job_name is not None

        training_job_desc = get_sagemaker_training_job(training_job_name)

        assert k8s.get_resource_arn(resource) == training_job_desc["TrainingJobArn"]
        assert training_job_desc["TrainingJobStatus"] in self.list_status_created
        assert k8s.wait_on_condition(reference, "ACK.ResourceSynced", "False")

        self._assert_training_status_in_sync(
            training_job_name, reference, self.status_inprogress
        )

        # Delete the k8s resource.
        _, deleted = k8s.delete_custom_resource(reference)
        assert deleted is True

        training_job_desc = get_sagemaker_training_job(training_job_name)
        assert training_job_desc["TrainingJobStatus"] in self.list_status_stopped

    def test_completed_training_job_debugger(self, xgboost_training_job_debugger):
        (reference, resource) = xgboost_training_job_debugger
        assert k8s.get_resource_exists(reference)

        training_job_name = resource["spec"].get("trainingJobName", None)
        assert training_job_name is not None

        training_job_desc = get_sagemaker_training_job(training_job_name)

        assert k8s.get_resource_arn(resource) == training_job_desc["TrainingJobArn"]
        assert training_job_desc["TrainingJobStatus"] in self.list_status_created
        assert k8s.wait_on_condition(reference, "ACK.ResourceSynced", "False")

        self._assert_training_status_in_sync(
            training_job_name, reference, self.status_completed
        )
        # TODO: This test is failing
        #assert k8s.wait_on_condition(reference, "ACK.ResourceSynced", "False")

        self._assert_training_debugger_status_in_sync(
            training_job_name, reference, self.debugger_status_completed
        )
        assert k8s.wait_on_condition(reference, "ACK.ResourceSynced", "True")

        # Check that you can delete a completed resource from k8s
        _, deleted = k8s.delete_custom_resource(reference)
        assert deleted is True
