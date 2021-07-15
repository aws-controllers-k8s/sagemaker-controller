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

from acktest.resources import random_suffix_name
from acktest.k8s import resource as k8s
from e2e import (
    service_marker,
    wait_for_status,
    create_sagemaker_resource,
    sagemaker_client,
)
from e2e.replacement_values import REPLACEMENT_VALUES
from e2e.bootstrap_resources import get_bootstrap_resources
from e2e.common import config as cfg
from time import sleep
import random

RESOURCE_PLURAL = "notebookinstances"
RESOURCE_PREFIX = "nb"
RESOURCE_SPEC_FILE = "notebook_instance"


@pytest.fixture(scope="function")
def notebook_instance():
    resource_name = RESOURCE_PREFIX + str(random.randint(0, 10000))
    replacements = REPLACEMENT_VALUES.copy()
    replacements["NOTEBOOK_INSTANCE_NAME"] = resource_name
    reference, spec, resource = create_sagemaker_resource(
        resource_plural=RESOURCE_PLURAL,
        resource_name=resource_name,
        spec_file=RESOURCE_SPEC_FILE,
        replacements=replacements,
    )

    assert resource is not None
    assert k8s.get_resource_arn(resource) is not None

    yield (reference, resource, spec)


def get_notebook_instance(notebook_instance_name: str):
    try:
        desired_notebook_instance = sagemaker_client().describe_notebook_instance(NotebookInstanceName=notebook_instance_name)
        return desired_notebook_instance
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
    def _wait_resource_notebook_status(self,reference: k8s.CustomResourceReference,
        expected_status: str,
        wait_periods: int = 30,
        period_length: int = 30):
        return wait_for_status(expected_status,
            wait_periods,
            period_length,
            get_notebook_instance_resource_status,
            reference
        )
    def _wait_sagemaker_notebook_status(self,notebook_instance_name: str,
        expected_status: str,
        wait_periods: int = 30,
        period_length: int = 30):
        return wait_for_status(expected_status,
            wait_periods,
            period_length,
            get_notebook_instance_sagemaker_status,
            notebook_instance_name
        )
    def _assert_notebook_status_in_sync(self, notebook_instance_name, reference, expected_status):
        assert(
            self._wait_sagemaker_notebook_status(notebook_instance_name,expected_status,wait_periods=60,period_length=15)
        == self._wait_resource_notebook_status(reference,expected_status,wait_periods=60,period_length=15)
        == expected_status
        )
    def testCreateUpdateAndDelete(self,notebook_instance):
        (reference, resource, spec) = notebook_instance
        assert k8s.get_resource_exists(reference)

        notebook_instance_name = resource["spec"].get("notebookInstanceName",None)
        assert notebook_instance_name is not None

        notebook_description = get_notebook_instance(notebook_instance_name)
        assert notebook_description["NotebookInstanceStatus"] == "Pending"

        #wait for the resource to go to the InService state and make sure the operator is synced with sagemaker.
        self._assert_notebook_status_in_sync(notebook_instance_name,reference,"InService")

        # Delete the k8s resource.
        _, deleted = k8s.delete_custom_resource(reference, 11, 30)
        assert deleted is True