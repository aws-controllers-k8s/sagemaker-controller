# Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"). You may
# not use this file except in compliance with the License. A copy of the
# License is located at
#
#	 http://aws.amazon.com/apache2.0/
#
# or in the "license" file accompanying this file. This file is distributed
# on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
# express or implied. See the License for the specific language governing
# permissions and limitations under the License.
"""Integration tests for the SageMaker TransformJob API.
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
    MODEL_RESOURCE_PLURAL,
)
from e2e.replacement_values import REPLACEMENT_VALUES
from e2e.bootstrap_resources import get_bootstrap_resources
from time import sleep

RESOURCE_PLURAL = 'transformjobs'

def _sagemaker_client():
    return boto3.client('sagemaker')

def _make_model():
    model_resource_name = random_suffix_name("xgboost-model", 32)

    replacements = REPLACEMENT_VALUES.copy()
    replacements["MODEL_NAME"] = model_resource_name

    _, _, model_resource = create_sagemaker_resource(
        resource_plural=MODEL_RESOURCE_PLURAL,
        resource_name=model_resource_name,
        spec_file="xgboost_model",
        replacements=replacements,
    )
    assert model_resource is not None
    assert k8s.get_resource_arn(model_resource) is not None

    return model_resource_name

def _make_transformjob():
    model_name = _make_model()
    resource_name = random_suffix_name("xgboost-transformjob", 32)

    replacements = REPLACEMENT_VALUES.copy()
    replacements["MODEL_NAME"] = model_name
    replacements["TRANSFORM_JOB_NAME"] = resource_name

    data = load_resource_file(
        resource_directory, "xgboost_transformjob", additional_replacements=replacements
    )

    reference = k8s.CustomResourceReference(
        CRD_GROUP, CRD_VERSION, RESOURCE_PLURAL, resource_name, namespace="default"
    )

    return reference, data

@pytest.fixture(scope="function")
def xgboost_transformjob():
    transform_job, data = _make_transformjob()
    resource = k8s.create_custom_resource(transform_job, data)
    resource = k8s.wait_resource_consumed_by_controller(transform_job)

    yield (transform_job, resource) 


def get_sagemaker_transform_job(transform_job_name: str):
    try:
        transform_desc = boto3.client('sagemaker').describe_transform_job(
            TransformJobName=transform_job_name
        )
        return transform_desc
    except BaseException:
        logging.error(
            f"SageMaker could not find a transform job with the name {transform_job_name}"
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
    list_status_created = ("InProgress", "Completed")
    list_status_stopped = ("Stopped", "Stopping")
    status_inprogress: str = "InProgress"
    status_completed: str = "Completed"

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
            reference
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
            self._wait_sagemaker_transform_status(
                transform_job_name, expected_status
            )
            == self._wait_resource_transform_status(reference, expected_status)
            == expected_status
        )

    def test_transform(self, xgboost_transformjob):
        (reference, resource) = xgboost_transformjob
        assert k8s.get_resource_exists(reference)
    
        transform_job_name = resource["spec"].get("transformJobName", None)
        assert transform_job_name is not None

        transform_sm_desc = get_sagemaker_transform_job(transform_job_name)
        assert k8s.get_resource_arn(resource) == transform_sm_desc["TransformJobArn"]
        assert transform_sm_desc["TransformJobStatus"] in self.list_status_created
        assert k8s.wait_on_condition(reference, "ACK.ResourceSynced", "False")

        self._assert_transform_status_in_sync(
            transform_job_name, reference, self.status_inprogress
        )

        # Delete the k8s resource.
        _, deleted = k8s.delete_custom_resource(reference)
        assert deleted is True

        transform_sm_desc = get_sagemaker_transform_job(transform_job_name)
        assert transform_sm_desc["TransformJobStatus"] in self.list_status_stopped

    def test_completed_transform(self, xgboost_transformjob):
        (reference, resource) = xgboost_transformjob
        assert k8s.get_resource_exists(reference)

        transform_job_name = resource["spec"].get("transformJobName", None)
        assert transform_job_name is not None

        transform_sm_desc = get_sagemaker_transform_job(transform_job_name)
        assert (
            k8s.get_resource_arn(resource) == transform_sm_desc["TransformJobArn"]
        )
        assert transform_sm_desc["TransformJobStatus"] in self.list_status_created
        assert k8s.wait_on_condition(reference, "ACK.ResourceSynced", "False")
        
        self._assert_transform_status_in_sync(
            transform_job_name, reference, self.status_completed
        )
        assert k8s.wait_on_condition(reference, "ACK.ResourceSynced", "True")

        #Check that you can delete a completed resource from k8s
        _, deleted = k8s.delete_custom_resource(reference)
        assert deleted is True
