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
"""Integration tests for the Notebook Lifecycle configuration
"""

import pytest
import logging
import botocore

from acktest.k8s import resource as k8s
from acktest.resources import random_suffix_name
from e2e import (
    service_marker,
    wait_for_status,
    create_sagemaker_resource,
    sagemaker_client,
)

from e2e.bootstrap_resources import get_bootstrap_resources
import random

from e2e.replacement_values import REPLACEMENT_VALUES
from time import sleep

RESOURCE_PLURAL = "notebookinstancelifecycleconfigs"
RESOURCE_NAME_PREFIX = "nblf"
RESOURCE_SPEC_FILE = "notebook_instance_lifecycle_config"

@pytest.fixture(scope="function")
def notebook_instance_lifecycleConfig():
    notebook_instance_lfc_name = random_suffix_name(RESOURCE_NAME_PREFIX, 32)
    replacements = REPLACEMENT_VALUES.copy()
    replacements["NOTEBOOK_INSTANCE_LFC_NAME"] = notebook_instance_lfc_name
    reference, spec, resource = create_sagemaker_resource(
        resource_plural=RESOURCE_PLURAL,
        resource_name=notebook_instance_lfc_name,
        spec_file=RESOURCE_SPEC_FILE,
        replacements=replacements,
    )
    yield (reference, resource,spec)
    if k8s.get_resource_exists(reference):
        _, deleted = k8s.delete_custom_resource(reference, 10, 5)
        assert deleted

def get_notebook_instance_lifecycle_config(notebook_instance_lfc_name: str):
    try:
        desired_notebook_instance_lfc = sagemaker_client().describe_notebook_instance_lifecycle_config(
            NotebookInstanceLifecycleConfigName=notebook_instance_lfc_name)
        return desired_notebook_instance_lfc
    except botocore.exceptions.ClientError as error:
        logging.error(
            f"SageMaker could not find a Notebook Instance with the name {notebook_instance_lfc_name}. Error {error}"
        )
        return None


@service_marker
@pytest.mark.canary
class TestNotebookInstanceLifecycleConfig:
    def test_CreateUpdateDeleteNotebookLifecycleConfig(self,notebook_instance_lifecycleConfig):
        (reference, resource,spec) = notebook_instance_lifecycleConfig
        assert k8s.get_resource_exists(reference)

        #Getting the resource name
        notebook_instance_lfc_name = resource["spec"].get("notebookInstanceLifecycleConfigName",None)
        assert notebook_instance_lfc_name is not None

        #Verifying that its set correctly
        spec["spec"]["onStart"] = [{"content":"cGlwIGluc3RhbGwgc2l4"}]
        k8s.patch_custom_resource(reference,spec)

        resource = k8s.wait_resource_consumed_by_controller(reference)
        assert resource is not None
        sleep(3) #Done to avoid flakiness

        #Verifying that an update was successful
        latest_notebook_lf = get_notebook_instance_lifecycle_config(notebook_instance_lfc_name)
        assert(latest_notebook_lf["OnStart"][0]["Content"] == "cGlwIGluc3RhbGwgc2l4")

        #Deleting the resource
        _, deleted = k8s.delete_custom_resource(reference, 10, 5)
        assert deleted is True
        assert get_notebook_instance_lifecycle_config(notebook_instance_lfc_name) is None



