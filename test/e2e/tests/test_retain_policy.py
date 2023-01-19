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
"""Integration tests for the Resource retain policy
"""

import pytest
import logging
import botocore

from acktest.k8s import resource as k8s
from acktest.resources import random_suffix_name
from e2e import (
    service_marker,
    create_sagemaker_resource,
    sagemaker_client,
)


from e2e.replacement_values import REPLACEMENT_VALUES


@pytest.fixture(scope="module")
def notebook_instance_lifecycle_config():
    notebook_instance_lfc_name = random_suffix_name("notebookinstancelfc", 40)
    replacements = REPLACEMENT_VALUES.copy()
    replacements["NOTEBOOK_INSTANCE_LFC_NAME"] = notebook_instance_lfc_name
    reference, spec, resource = create_sagemaker_resource(
        resource_plural="notebookinstancelifecycleconfigs",
        resource_name=notebook_instance_lfc_name,
        spec_file="notebook_instance_lifecycle_retain",
        replacements=replacements,
    )
    assert resource is not None
    yield (reference, resource, spec)

    delete_notebook_instance_lifecycle_config(notebook_instance_lfc_name)


def get_notebook_instance_lifecycle_config(notebook_instance_lfc_name: str):
    try:
        resp = sagemaker_client().describe_notebook_instance_lifecycle_config(
            NotebookInstanceLifecycleConfigName=notebook_instance_lfc_name
        )
        return resp
    except botocore.exceptions.ClientError as error:
        logging.error(
            f"SageMaker could not find a Notebook Instance Lifecycle Configuration with the name {notebook_instance_lfc_name}. Error {error}"
        )
        return None


def delete_notebook_instance_lifecycle_config(notebook_instance_lfc_name: str):
    try:
        sagemaker_client().delete_notebook_instance_lifecycle_config(
            NotebookInstanceLifecycleConfigName=notebook_instance_lfc_name
        )
    except botocore.exceptions.ClientError as error:
        logging.error(
            f"SageMaker could not find a Notebook Instance Lifecycle Configuration with the name {notebook_instance_lfc_name}. Error {error}"
        )


@service_marker
class TestNotebookInstanceLifecycleConfig:
    def test_retain_resource(self, notebook_instance_lifecycle_config):
        (reference, resource, _) = notebook_instance_lifecycle_config
        assert k8s.get_resource_exists(reference)

        # Getting the resource name
        notebook_instance_lfc_name = resource["spec"].get(
            "notebookInstanceLifecycleConfigName", None
        )
        assert notebook_instance_lfc_name is not None
        notebook_instance_lfc_desc = get_notebook_instance_lifecycle_config(
            notebook_instance_lfc_name
        )
        assert (
            k8s.get_resource_arn(resource)
            == notebook_instance_lfc_desc["NotebookInstanceLifecycleConfigArn"]
        )

        # Delete the CR in Kubernetes
        _, deleted = k8s.delete_custom_resource(reference)

        assert deleted

        # Verify that the resource was not deleted on SageMaker
        assert (
            get_notebook_instance_lifecycle_config(notebook_instance_lfc_name)
            is not None
        )