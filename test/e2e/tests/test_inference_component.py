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

import pytest
import logging

from acktest.resources import random_suffix_name
from acktest.k8s import resource as k8s

from e2e import (
    service_marker,
    create_sagemaker_resource,
    delete_custom_resource,
    assert_inference_component_status_in_sync,
    assert_endpoint_status_in_sync,
    assert_tags_in_sync,
    get_sagemaker_inference_component,
    get_sagemaker_endpoint
)
from e2e.replacement_values import REPLACEMENT_VALUES
from e2e.common import config as cfg


@pytest.fixture(scope="module")
def name_suffix():
    return random_suffix_name("ic-xgboost", 32)


@pytest.fixture(scope="module")
def xgboost_model(name_suffix):
    model_resource_name = name_suffix + "-model"
    replacements = REPLACEMENT_VALUES.copy()
    replacements["MODEL_NAME"] = model_resource_name

    model_reference, model_spec, model_resource = create_sagemaker_resource(
        resource_plural=cfg.MODEL_RESOURCE_PLURAL,
        resource_name=model_resource_name,
        spec_file="xgboost_model_inference_component",
        replacements=replacements,
    )
    assert model_resource is not None
    if k8s.get_resource_arn(model_resource) is None:
        logging.error(
            f"ARN for this resource is None, resource status is: {model_resource['status']}"
        )
    assert k8s.get_resource_arn(model_resource) is not None

    yield (model_reference, model_resource)

    _, deleted = k8s.delete_custom_resource(
        model_reference, cfg.DELETE_WAIT_PERIOD, cfg.DELETE_WAIT_LENGTH
    )
    assert deleted


@pytest.fixture(scope="module")
def endpoint_config(name_suffix):
    config_resource_name = name_suffix + "-endpoint-config"
    replacements = REPLACEMENT_VALUES.copy()
    replacements["ENDPOINT_CONFIG_NAME"] = config_resource_name

    config_reference, config_spec, config_resource = create_sagemaker_resource(
        resource_plural=cfg.ENDPOINT_CONFIG_RESOURCE_PLURAL,
        resource_name=config_resource_name,
        spec_file="endpoint_config_inference_component",
        replacements=replacements,
    )
    assert config_resource is not None
    if k8s.get_resource_arn(config_resource) is None:
        logging.error(
            f"ARN for this resource is None, resource status is: {config_resource['status']}"
        )
    assert k8s.get_resource_arn(config_resource) is not None

    yield (config_reference, config_resource)

    _, deleted = k8s.delete_custom_resource(
        config_reference, cfg.DELETE_WAIT_PERIOD, cfg.DELETE_WAIT_LENGTH
    )
    assert deleted


@pytest.fixture(scope="module")
def endpoint(name_suffix, endpoint_config):
    endpoint_resource_name = name_suffix + "-endpoint"
    (_, config_resource) = endpoint_config
    config_resource_name = config_resource["spec"].get("endpointConfigName", None)

    replacements = REPLACEMENT_VALUES.copy()
    replacements["ENDPOINT_NAME"] = endpoint_resource_name
    replacements["ENDPOINT_CONFIG_NAME"] = config_resource_name

    endpoint_reference, endpoint_spec, endpoint_resource = create_sagemaker_resource(
        resource_plural=cfg.ENDPOINT_RESOURCE_PLURAL,
        resource_name=endpoint_resource_name,
        spec_file="endpoint_base",
        replacements=replacements,
    )

    assert endpoint_resource is not None

    # endpoint has correct arn and status
    endpoint_name = endpoint_resource["spec"].get("endpointName", None)
    assert endpoint_name is not None
    assert endpoint_name == endpoint_resource_name

    endpoint_desc = get_sagemaker_endpoint(endpoint_name)
    endpoint_arn = endpoint_desc["EndpointArn"]
    assert k8s.get_resource_arn(endpoint_resource) == endpoint_arn

    # endpoint transitions Creating -> InService state
    assert_endpoint_status_in_sync(
        endpoint_name, endpoint_reference, cfg.ENDPOINT_STATUS_CREATING
    )
    assert k8s.wait_on_condition(endpoint_reference, "ACK.ResourceSynced", "False")

    assert_endpoint_status_in_sync(
        endpoint_name, endpoint_reference, cfg.ENDPOINT_STATUS_INSERVICE
    )
    assert k8s.wait_on_condition(endpoint_reference, "ACK.ResourceSynced", "True")

    yield (endpoint_reference, endpoint_resource)

    # Delete the k8s resource if not already deleted by tests
    assert delete_custom_resource(endpoint_reference, 40, cfg.DELETE_WAIT_LENGTH)


@pytest.fixture(scope="module")
def inference_component(name_suffix, endpoint, xgboost_model):
    inference_component_resource_name = name_suffix + "-inference-component"
    (_, endpoint_resource) = endpoint
    (_, model_resource) = xgboost_model
    endpoint_resource_name = endpoint_resource["spec"].get("endpointName", None)
    model_resource_name = model_resource["spec"].get("modelName", None)

    replacements = REPLACEMENT_VALUES.copy()
    replacements["INFERENCE_COMPONENT_NAME"] = inference_component_resource_name
    replacements["ENDPOINT_NAME"] = endpoint_resource_name
    replacements["MODEL_NAME"] = model_resource_name

    reference, spec, resource = create_sagemaker_resource(
        resource_plural=cfg.INFERENCE_COMPONENT_RESOURCE_PLURAL,
        resource_name=inference_component_resource_name,
        spec_file="inference_component",
        replacements=replacements,
    )

    assert resource is not None

    yield (reference, resource, spec)

    # Delete the k8s resource if not already deleted by tests
    assert delete_custom_resource(reference, 40, cfg.DELETE_WAIT_LENGTH)


@service_marker
@pytest.mark.shallow_canary
@pytest.mark.canary
class TestInferenceComponent:
    def create_inference_component_test(self, inference_component):
        (reference, resource, _) = inference_component
        assert k8s.get_resource_exists(reference)

        # inference component has correct arn and status
        inference_component_name = resource["spec"].get("inferenceComponentName", None)
        assert inference_component_name is not None

        inference_component_desc = get_sagemaker_inference_component(inference_component_name)
        inference_component_arn = inference_component_desc["InferenceComponentArn"]
        assert k8s.get_resource_arn(resource) == inference_component_arn

        # inference_component transitions Creating -> InService state
        assert_inference_component_status_in_sync(
            inference_component_name, reference, cfg.INFERENCE_COMPONENT_STATUS_CREATING
        )
        assert k8s.wait_on_condition(reference, "ACK.ResourceSynced", "False")

        assert_inference_component_status_in_sync(
            inference_component_name, reference, cfg.INFERENCE_COMPONENT_STATUS_INSERVICE
        )
        assert k8s.wait_on_condition(reference, "ACK.ResourceSynced", "True")

        resource_tags = resource["spec"].get("tags", None)
        assert_tags_in_sync(inference_component_arn, resource_tags)

    def update_inference_component_successful_test(self, inference_component):
        (reference, resource, spec) = inference_component

        inference_component_name = resource["spec"].get("inferenceComponentName", None)

        desired_memory_required = 2024

        spec["spec"]["specification"]["computeResourceRequirements"]["minMemoryRequiredInMb"] = desired_memory_required
        resource = k8s.patch_custom_resource(reference, spec)
        resource = k8s.wait_resource_consumed_by_controller(reference)
        assert resource is not None

        # inference component transitions Updating -> InService state
        assert_inference_component_status_in_sync(
            reference.name,
            reference,
            cfg.INFERENCE_COMPONENT_STATUS_UPDATING,
        )

        assert k8s.wait_on_condition(reference, "ACK.ResourceSynced", "False")
        assert k8s.get_resource_condition(reference, "ACK.Terminal") is None
        resource = k8s.get_resource(reference)

        assert_inference_component_status_in_sync(
            reference.name,
            reference,
            cfg.INFERENCE_COMPONENT_STATUS_INSERVICE,
        )
        assert k8s.wait_on_condition(reference, "ACK.ResourceSynced", "True")
        assert k8s.get_resource_condition(reference, "ACK.Terminal") is None
        resource = k8s.get_resource(reference)
        assert resource["status"].get("failureReason", None) is None
        new_memory_required = get_sagemaker_inference_component(inference_component_name)[
            "Specification"]["ComputeResourceRequirements"]["MinMemoryRequiredInMb"]

        assert desired_memory_required == new_memory_required

    def delete_inference_component_test(self, inference_component):
        (reference, resource, _) = inference_component
        inference_component_name = resource["spec"].get("inferenceComponentName", None)

        assert delete_custom_resource(
            reference, cfg.INFERENCE_COMPONENT_DELETE_WAIT_PERIODS, cfg.INFERENCE_COMPONENT_DELETE_WAIT_LENGTH
        )

        assert get_sagemaker_inference_component(inference_component_name) is None

    def test_driver(
            self,
            inference_component
    ):
        self.create_inference_component_test(inference_component)
        self.update_inference_component_successful_test(inference_component)
        self.delete_inference_component_test(inference_component)
