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
"""Integration tests for the SageMaker Endpoint API.
"""

import boto3
import pytest
import logging
import time
from typing import Dict

from acktest.aws import s3
from acktest.resources import random_suffix_name
from acktest.k8s import resource as k8s

from e2e import (
    service_marker,
    CRD_GROUP,
    CRD_VERSION,
    create_adopted_resource,
    wait_sagemaker_endpoint_status,
    assert_endpoint_status_in_sync,
    sagemaker_client,
)
from e2e.replacement_values import REPLACEMENT_VALUES
from e2e.common import config as cfg


@pytest.fixture(scope="module")
def name_suffix():
    return random_suffix_name("sdk-endpoint", 32)


def sdk_make_model(model_name):
    data_bucket = REPLACEMENT_VALUES["SAGEMAKER_DATA_BUCKET"]
    model_input = {
        "ModelName": model_name,
        "Containers": [
            {
                "Image": REPLACEMENT_VALUES["XGBOOST_IMAGE_URI"],
                "ModelDataUrl": f"s3://{data_bucket}/sagemaker/model/xgboost-mnist-model.tar.gz",
            }
        ],
        "ExecutionRoleArn": REPLACEMENT_VALUES["SAGEMAKER_EXECUTION_ROLE_ARN"],
    }

    model_response = sagemaker_client().create_model(**model_input)
    assert model_response.get("ModelArn", None) is not None
    return model_input, model_response


def sdk_make_endpoint_config(model_name, endpoint_config_name):
    endpoint_config_input = {
        "EndpointConfigName": endpoint_config_name,
        "ProductionVariants": [
            {
                "VariantName": "variant-1",
                "ModelName": model_name,
                "InitialInstanceCount": 1,
                "InstanceType": "ml.c5.large",
            }
        ],
    }

    endpoint_config_response = sagemaker_client().create_endpoint_config(
        **endpoint_config_input
    )
    assert endpoint_config_response.get("EndpointConfigArn", None) is not None
    return endpoint_config_input, endpoint_config_response


def sdk_make_endpoint(endpoint_name, endpoint_config_name):
    endpoint_input = {
        "EndpointName": endpoint_name,
        "EndpointConfigName": endpoint_config_name,
    }
    endpoint_response = sagemaker_client().create_endpoint(**endpoint_input)
    assert endpoint_response.get("EndpointArn", None) is not None

    return endpoint_input, endpoint_response


@pytest.fixture(scope="module")
def sdk_endpoint(name_suffix):
    model_name = name_suffix + "-model"
    endpoint_config_name = name_suffix + "-config"
    endpoint_name = name_suffix

    model_input, model_response = sdk_make_model(model_name)
    endpoint_config_input, endpoint_config_response = sdk_make_endpoint_config(
        model_name, endpoint_config_name
    )
    endpoint_input, endpoint_response = sdk_make_endpoint(
        endpoint_name, endpoint_config_name
    )

    yield (
        model_input,
        model_response,
        endpoint_config_input,
        endpoint_config_response,
        endpoint_input,
        endpoint_response,
    )
    wait_sagemaker_endpoint_status(endpoint_name, cfg.ENDPOINT_STATUS_INSERVICE)
    sagemaker_client().delete_endpoint(EndpointName=endpoint_name)
    sagemaker_client().delete_endpoint_config(EndpointConfigName=endpoint_config_name)
    sagemaker_client().delete_model(ModelName=model_name)


@pytest.fixture(scope="module")
def adopted_endpoint(sdk_endpoint):
    (model_input, _, endpoint_config_input, _, endpoint_input, _) = sdk_endpoint

    replacements = REPLACEMENT_VALUES.copy()
    # adopt model
    replacements["ADOPTED_RESOURCE_NAME"] = "adopt-" + model_input["ModelName"]
    replacements["TARGET_RESOURCE_AWS"] = replacements[
        "TARGET_RESOURCE_K8S"
    ] = model_input["ModelName"]
    replacements["RESOURCE_KIND"] = "Model"

    adopt_model_reference, _, adopt_model_resource = create_adopted_resource(
        replacements=replacements,
    )
    assert adopt_model_resource is not None

    # adopt endpoint config
    replacements["ADOPTED_RESOURCE_NAME"] = (
        "adopt-" + endpoint_config_input["EndpointConfigName"]
    )
    replacements["TARGET_RESOURCE_AWS"] = replacements[
        "TARGET_RESOURCE_K8S"
    ] = endpoint_config_input["EndpointConfigName"]
    replacements["RESOURCE_KIND"] = "EndpointConfig"

    adopt_config_reference, _, adopt_config_resource = create_adopted_resource(
        replacements=replacements,
    )
    assert adopt_config_resource is not None

    # adopt endpoint
    replacements["ADOPTED_RESOURCE_NAME"] = "adopt-" + endpoint_input["EndpointName"]
    replacements["TARGET_RESOURCE_AWS"] = replacements[
        "TARGET_RESOURCE_K8S"
    ] = endpoint_input["EndpointName"]
    replacements["RESOURCE_KIND"] = "Endpoint"

    adopt_endpoint_reference, _, adopt_endpoint_resource = create_adopted_resource(
        replacements=replacements,
    )
    assert adopt_endpoint_resource is not None

    yield (adopt_model_reference, adopt_config_reference, adopt_endpoint_reference)

    for cr in (adopt_model_reference, adopt_config_reference, adopt_endpoint_reference):
        if k8s.get_resource_exists(cr):
            _, deleted = k8s.delete_custom_resource(cr, 3, 10)
            assert deleted


@service_marker
class TestAdoptedEndpoint:
    def test_smoke(self, sdk_endpoint, adopted_endpoint):
        (
            adopt_model_reference,
            adopt_config_reference,
            adopt_endpoint_reference,
        ) = adopted_endpoint

        (
            model_input,
            model_response,
            endpoint_config_input,
            endpoint_config_response,
            endpoint_input,
            endpoint_response,
        ) = sdk_endpoint

        namespace = "default"
        model_name = k8s.get_resource(adopt_model_reference)["spec"]["aws"]["nameOrID"]
        endpoint_config_name = k8s.get_resource(adopt_config_reference)["spec"]["aws"][
            "nameOrID"
        ]
        endpoint_name = k8s.get_resource(adopt_endpoint_reference)["spec"]["aws"][
            "nameOrID"
        ]

        for reference in (
            adopt_model_reference,
            adopt_config_reference,
            adopt_endpoint_reference,
        ):
            assert k8s.wait_on_condition(reference, "ACK.Adopted", "True")

        model_reference = k8s.create_reference(
            CRD_GROUP, CRD_VERSION, cfg.MODEL_RESOURCE_PLURAL, model_name, namespace
        )
        model_resource = k8s.wait_resource_consumed_by_controller(model_reference)
        assert model_resource is not None

        assert model_resource["spec"].get("modelName", None) == model_name
        assert model_resource["spec"].get("containers", None) is not None
        assert (
            model_resource["spec"].get("executionRoleARN", None)
            == model_input["ExecutionRoleArn"]
        )
        assert k8s.get_resource_arn(model_resource) == model_response.get(
            "ModelArn", None
        )

        config_reference = k8s.create_reference(
            CRD_GROUP,
            CRD_VERSION,
            cfg.ENDPOINT_CONFIG_RESOURCE_PLURAL,
            endpoint_config_name,
            namespace,
        )
        config_resource = k8s.wait_resource_consumed_by_controller(config_reference)
        assert config_resource is not None

        assert (
            config_resource["spec"].get("endpointConfigName", None)
            == endpoint_config_name
        )
        assert config_resource["spec"].get("productionVariants", None) is not None
        assert k8s.get_resource_arn(config_resource) == endpoint_config_response.get(
            "EndpointConfigArn", None
        )

        endpoint_reference = k8s.create_reference(
            CRD_GROUP,
            CRD_VERSION,
            cfg.ENDPOINT_RESOURCE_PLURAL,
            endpoint_name,
            namespace,
        )
        endpoint_resource = k8s.wait_resource_consumed_by_controller(endpoint_reference)
        assert endpoint_resource is not None

        assert endpoint_resource["spec"].get("endpointName", None) == endpoint_name
        assert (
            endpoint_resource["spec"].get("endpointConfigName", None)
            == endpoint_config_name
        )
        assert k8s.get_resource_arn(endpoint_resource) == endpoint_response.get(
            "EndpointArn", None
        )

        assert_endpoint_status_in_sync(
            endpoint_name, endpoint_reference, cfg.ENDPOINT_STATUS_INSERVICE,
        )
        assert k8s.wait_on_condition(endpoint_reference, "ACK.ResourceSynced", "True")
