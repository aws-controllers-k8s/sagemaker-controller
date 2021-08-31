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
"""Integration tests for the SageMaker NotebookInstance API.
"""

import botocore
import pytest
import logging
from typing import Dict

from acktest.k8s import resource as k8s
from acktest.resources import random_suffix_name
from e2e import (
    service_marker,
    wait_for_status,
    create_sagemaker_resource,
    sagemaker_client,
)
from e2e.replacement_values import REPLACEMENT_VALUES
import random

DELETE_WAIT_PERIOD = 16
DELETE_WAIT_LENGTH = 30


@pytest.fixture(scope="module")
def notebook_instance():
    default_code_repository = "https://github.com/aws-controllers-k8s/community"
    resource_name = random_suffix_name("nb", 32)
    replacements = REPLACEMENT_VALUES.copy()
    replacements["NOTEBOOK_INSTANCE_NAME"] = resource_name
    replacements["DEFAULT_CODE_REPOSITORY"] = default_code_repository
    reference, spec, resource = create_sagemaker_resource(
        resource_plural="notebookinstances",
        resource_name=resource_name,
        spec_file="notebook_instance",
        replacements=replacements,
    )

    assert resource is not None
    yield (reference, resource, spec)

    # Delete the k8s resource if not already deleted by tests
    if k8s.get_resource_exists(reference):
        _, deleted = k8s.delete_custom_resource(
            reference, DELETE_WAIT_PERIOD, DELETE_WAIT_LENGTH
        )
        assert deleted


def get_notebook_instance(notebook_instance_name: str):
    try:
        resp = sagemaker_client().describe_notebook_instance(
            NotebookInstanceName=notebook_instance_name
        )
        return resp
    except botocore.exceptions.ClientError as error:
        logging.error(
            f"SageMaker could not find a Notebook Instance with the name {notebook_instance_name}. Error {error}"
        )
        return None


def get_notebook_instance_sagemaker_status(notebook_instance_name: str):
    notebook_instance = get_notebook_instance(notebook_instance_name)
    assert notebook_instance is not None
    return notebook_instance["NotebookInstanceStatus"]


def get_notebook_instance_resource_status(reference: k8s.CustomResourceReference):
    resource = k8s.get_resource(reference)
    assert resource is not None
    assert "notebookInstanceStatus" in resource["status"]
    return resource["status"]["notebookInstanceStatus"]


@pytest.mark.canary
@service_marker
class TestNotebookInstance:
    def _wait_resource_notebook_status(
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
            get_notebook_instance_resource_status,
            reference,
        )

    def _wait_sagemaker_notebook_status(
        self,
        notebook_instance_name: str,
        expected_status: str,
        wait_periods: int = 30,
        period_length: int = 30,
    ):
        return wait_for_status(
            expected_status,
            wait_periods,
            period_length,
            get_notebook_instance_sagemaker_status,
            notebook_instance_name,
        )

    def _assert_notebook_status_in_sync(
        self,
        notebook_instance_name,
        reference,
        expected_status,
        wait_periods=30,
        period_length=30,
    ):
        assert (
            self._wait_sagemaker_notebook_status(
                notebook_instance_name, expected_status, wait_periods, period_length
            )
            == self._wait_resource_notebook_status(
                reference, expected_status, wait_periods, period_length
            )
            == expected_status
        )

    def create_notebook_test(self, notebook_instance):
        (reference, resource, _) = notebook_instance
        assert k8s.get_resource_exists(reference)
        assert k8s.get_resource_arn(resource) is not None

        # Create the resource and verify that its Pending
        notebook_instance_name = resource["spec"].get("notebookInstanceName", None)
        assert notebook_instance_name is not None

        notebook_description = get_notebook_instance(notebook_instance_name)
        assert notebook_description["NotebookInstanceStatus"] == "Pending"

        assert k8s.wait_on_condition(reference, "ACK.ResourceSynced", "False")
        self._assert_notebook_status_in_sync(
            notebook_instance_name, reference, "Pending"
        )

        # wait for the resource to go to the InService state and make sure the operator is synced with sagemaker.
        self._assert_notebook_status_in_sync(
            notebook_instance_name, reference, "InService"
        )
        assert k8s.wait_on_condition(reference, "ACK.ResourceSynced", "True")

    def update_notebook_test(self, notebook_instance):
        (reference, resource, spec) = notebook_instance
        notebook_instance_name = resource["spec"].get("notebookInstanceName", None)
        volumeSizeInGB = 7
        additionalCodeRepositories = ["https://github.com/aws-controllers-k8s/runtime"]

        spec["spec"]["volumeSizeInGB"] = volumeSizeInGB
        # TODO: Use del spec["spec"]["defaultCodeRepository"] instead once ack.testinfra supports replacement.
        # Patch only supports updating spec fields instead of fully getting rid of them.
        spec["spec"]["defaultCodeRepository"] = None
        spec["spec"]["additionalCodeRepositories"] = additionalCodeRepositories
        k8s.patch_custom_resource(reference, spec)

        self._assert_notebook_status_in_sync(
            notebook_instance_name, reference, "Stopping"
        )
        # TODO: Replace with annotations once runtime can update annotations in readOne.
        resource = k8s.get_resource(reference)
        # Test is flakey as this field can get changed before we get resource
        # UpdateTriggered can only be in the status if beforehand it was UpdatePending
        # TODO: See if update code can be restructured to avoid this
        assert resource["status"]["stoppedByControllerMetadata"] == ( "UpdatePending" or "UpdateTriggered")

        # wait for the resource to go to the InService state and make sure the operator is synced with sagemaker.
        self._assert_notebook_status_in_sync(
            notebook_instance_name, reference, "InService"
        )
        assert k8s.wait_on_condition(reference, "ACK.ResourceSynced", "True")

        notebook_instance_desc = get_notebook_instance(notebook_instance_name)
        assert notebook_instance_desc["VolumeSizeInGB"] == volumeSizeInGB

        resource = k8s.get_resource(reference)
        assert resource["spec"]["volumeSizeInGB"] == volumeSizeInGB

        assert "DefaultCodeRepository" not in notebook_instance_desc
        assert "defaultCodeRepository" not in resource["spec"]

        assert resource["spec"]["additionalCodeRepositories"] == additionalCodeRepositories
        assert notebook_instance_desc["AdditionalCodeRepositories"] == additionalCodeRepositories

        assert "stoppedByControllerMetadata" not in resource["status"]

    def delete_notebook_test(self, notebook_instance):
        # Delete the k8s resource.
        (reference, resource, _) = notebook_instance
        notebook_instance_name = resource["spec"].get("notebookInstanceName", None)
        _, deleted = k8s.delete_custom_resource(
            reference, DELETE_WAIT_PERIOD, DELETE_WAIT_LENGTH
        )
        assert deleted is True
        assert get_notebook_instance(notebook_instance_name) is None

    def test_driver(self, notebook_instance):
        self.create_notebook_test(notebook_instance)
        self.update_notebook_test(notebook_instance)
        self.delete_notebook_test(notebook_instance)
