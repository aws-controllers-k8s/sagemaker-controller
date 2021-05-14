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
"""Integration tests for the SageMaker EndpointConfig API.
"""

from _pytest import config
import botocore
import pytest
import logging
from typing import Dict

from acktest.resources import random_suffix_name
from acktest.k8s import resource as k8s

from e2e import (
    service_marker,
    create_sagemaker_resource,
    sagemaker_client
)
from e2e.replacement_values import REPLACEMENT_VALUES
from e2e.common import config as cfg


@pytest.fixture(scope="module")
def single_variant_config():
    config_resource_name = random_suffix_name("single-variant-config", 32)
    model_resource_name = config_resource_name + "-model"

    replacements = REPLACEMENT_VALUES.copy()
    replacements["ENDPOINT_CONFIG_NAME"] = config_resource_name
    replacements["MODEL_NAME"] = model_resource_name

    model_reference, model_spec, model_resource = create_sagemaker_resource(
        resource_plural=cfg.MODEL_RESOURCE_PLURAL,
        resource_name=model_resource_name,
        spec_file="xgboost_model",
        replacements=replacements,
    )
    assert model_resource is not None
    assert k8s.get_resource_arn(model_resource) is not None

    config_reference, config_spec, config_resource = create_sagemaker_resource(
        resource_plural=cfg.ENDPOINT_CONFIG_RESOURCE_PLURAL,
        resource_name=config_resource_name,
        spec_file="endpoint_config_single_variant",
        replacements=replacements,
    )
    assert config_resource is not None

    yield (config_reference, config_resource)

    k8s.delete_custom_resource(model_reference, 3, 10)
    # Delete the k8s resource if not already deleted by tests
    if k8s.get_resource_exists(config_reference):
        _, deleted = k8s.delete_custom_resource(config_reference, 3, 10)
        assert deleted

def get_sagemaker_endpoint_config(config_name: str):
    try:
        return sagemaker_client().describe_endpoint_config(
            EndpointConfigName=config_name
        )
    except botocore.exceptions.ClientError as error:
        logging.error(
            f"SageMaker could not find a config with the name {config_name}. Error {error}"
        )
        return None


@service_marker
@pytest.mark.canary
class TestEndpointConfig:
    def test_create_endpoint_config(self, single_variant_config):
        (reference, resource) = single_variant_config
        assert k8s.get_resource_exists(reference)

        config_name = resource["spec"].get("endpointConfigName", None)

        assert k8s.get_resource_arn(resource) == get_sagemaker_endpoint_config(config_name)["EndpointConfigArn"]

        # Delete the k8s resource.
        _, deleted = k8s.delete_custom_resource(reference, 3, 10)
        assert deleted

        assert get_sagemaker_endpoint_config(config_name) is None

