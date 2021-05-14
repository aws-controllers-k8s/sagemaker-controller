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
"""Integration tests for the SageMaker Model API.
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
    sagemaker_client,
)
from e2e.replacement_values import REPLACEMENT_VALUES
from e2e.common import config as cfg


@pytest.fixture(scope="module")
def xgboost_model():
    resource_name = random_suffix_name("xgboost-model", 32)

    replacements = REPLACEMENT_VALUES.copy()
    replacements["MODEL_NAME"] = resource_name

    reference, spec, resource = create_sagemaker_resource(
        resource_plural=cfg.MODEL_RESOURCE_PLURAL,
        resource_name=resource_name,
        spec_file="xgboost_model",
        replacements=replacements,
    )
    assert resource is not None

    yield (reference, resource)

    # Delete the k8s resource if not already deleted by tests
    if k8s.get_resource_exists(reference):
        _, deleted = k8s.delete_custom_resource(reference, 3, 10)
        assert deleted

def get_sagemaker_model(model_name: str):
    try:
        return sagemaker_client().describe_model(ModelName=model_name)
    except botocore.exceptions.ClientError as error:
        logging.error(
            f"SageMaker could not find a model with the name {model_name}. Error {error}"
        )
        return None
@service_marker
@pytest.mark.canary
class TestModel:
    def test_create_model(self, xgboost_model):
        (reference, resource) = xgboost_model
        assert k8s.get_resource_exists(reference)

        model_name = resource["spec"].get("modelName", None)

        assert k8s.get_resource_arn(resource) == get_sagemaker_model(model_name)["ModelArn"]

        # Delete the k8s resource.
        _, deleted = k8s.delete_custom_resource(reference, 3, 10)
        assert deleted

        assert get_sagemaker_model(model_name) is None

