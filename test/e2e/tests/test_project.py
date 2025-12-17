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
"""Integration tests for the SageMaker Project."""

import logging

import boto3
import pytest
from acktest.k8s import resource as k8s
from acktest.resources import random_suffix_name

from e2e import (
    create_sagemaker_resource,
    delete_custom_resource,
    service_marker,
    wait_for_status,
)
from e2e.replacement_values import REPLACEMENT_VALUES

PROJECT_WAIT_PERIOD = 120
PROJECT_WAIT_LENGTH = 30
PROJECT_STATUS_CREATE_FAILED = "CreateFailed"
FAIL_CREATE_ERROR_MESSAGE = "Codebuild to checkin seedcode has status FAILED"


def get_project_sagemaker_status(project_name):
    response = boto3.client("sagemaker").describe_project(ProjectName=project_name)
    return response["ProjectStatus"]


def get_k8s_resource_status(reference: k8s.CustomResourceReference):
    resource = k8s.get_resource(reference)
    assert "projectStatus" in resource["status"]
    return resource["status"]["projectStatus"]


def get_k8s_resource_message(reference: k8s.CustomResourceReference):
    resource = k8s.get_resource(reference)
    assert "serviceCatalogProvisionedProductDetails" in resource["status"]
    assert (
        "provisionedProductStatusMessage"
        in resource["status"]["serviceCatalogProvisionedProductDetails"]
    )
    return resource["status"]["serviceCatalogProvisionedProductDetails"][
        "provisionedProductStatusMessage"
    ]


def apply_project_yaml(resource_name):
    replacements = REPLACEMENT_VALUES.copy()
    replacements["PROJECT_NAME"] = resource_name

    reference, spec, resource = create_sagemaker_resource(
        resource_plural="projects",
        resource_name=resource_name,
        spec_file="project",
        replacements=replacements,
    )
    return reference, resource, spec


def assert_project_status_in_sync(project_name, reference, expected_status):
    sm_status = wait_for_status(
        expected_status,
        PROJECT_WAIT_PERIOD,
        PROJECT_WAIT_LENGTH,
        get_project_sagemaker_status,
        project_name,
    )
    k8s_status = wait_for_status(
        expected_status,
        PROJECT_WAIT_PERIOD,
        PROJECT_WAIT_LENGTH,
        get_k8s_resource_status,
        reference,
    )
    assert sm_status == k8s_status == expected_status


def assert_project_status_message(reference, expected_message):
    k8s_message = wait_for_status(
        expected_message,
        PROJECT_WAIT_PERIOD,
        PROJECT_WAIT_LENGTH,
        get_k8s_resource_message,
        reference,
    )
    assert expected_message in k8s_message


@pytest.fixture(scope="module")
def project_fixture():
    resource_name = random_suffix_name("sm-project", 20)
    reference, resource, spec = apply_project_yaml(resource_name)

    assert resource is not None
    if k8s.get_resource_arn(resource) is None:
        logging.error(
            f"ARN for this resource is None, resource status is: {resource['status']}"
        )
    assert k8s.get_resource_arn(resource) is not None

    yield (reference, resource)

    assert delete_custom_resource(
        reference,
        PROJECT_WAIT_PERIOD,
        PROJECT_WAIT_PERIOD,
    )


@service_marker
class TestProject:
    def create_failed_project(self, project_fixture):
        (reference, resource) = project_fixture

        assert k8s.get_resource_exists(reference)

        project_name = resource["spec"].get("projectName", None)
        assert project_name is not None

        assert_project_status_in_sync(
            reference.name,
            reference,
            PROJECT_STATUS_CREATE_FAILED,
        )

        assert_project_status_message(reference, FAIL_CREATE_ERROR_MESSAGE)

        resource = k8s.get_resource(reference)
        assert (
            resource["status"].get("serviceCatalogProvisionedProductDetails", None)
            is not None
        )
        assert (
            resource["status"]["serviceCatalogProvisionedProductDetails"].get(
                "provisionedProductStatusMessage", None
            )
            is not None
        )

    def test_driver(self, project_fixture):
        self.create_failed_project(project_fixture)
