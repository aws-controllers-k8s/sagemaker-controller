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

from acktest.aws import s3
from acktest.resources import random_suffix_name
from acktest.k8s import resource as k8s
from acktest.k8s import condition as ack_condition

from e2e import (
    service_marker,
    create_sagemaker_resource,
    delete_custom_resource,
    assert_endpoint_status_in_sync,
    assert_tags_in_sync,
    get_sagemaker_endpoint,
)
from e2e.replacement_values import REPLACEMENT_VALUES
from e2e.common import config as cfg

FAIL_UPDATE_ERROR_MESSAGE = "api error EndpointUpdateError: unable to update endpoint. check FailureReason. latest EndpointConfigName is "
# annontation key for last endpoint config name used for update
LAST_ENDPOINTCONFIG_UPDATE_ANNOTATION = "sagemaker.services.k8s.aws/last-endpoint-config-for-update"


@pytest.fixture(scope="module")
def name_suffix():
    return random_suffix_name("xgboost-endpoint", 32)


@pytest.fixture(scope="module")
def single_container_model(name_suffix):
    model_resource_name = name_suffix + "-model"
    replacements = REPLACEMENT_VALUES.copy()
    replacements["MODEL_NAME"] = model_resource_name

    model_reference, model_spec, model_resource = create_sagemaker_resource(
        resource_plural=cfg.MODEL_RESOURCE_PLURAL,
        resource_name=model_resource_name,
        spec_file="xgboost_model",
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
def multi_variant_config(name_suffix, single_container_model):
    config_resource_name = name_suffix + "-multi-variant-config"
    (_, model_resource) = single_container_model
    model_resource_name = model_resource["spec"].get("modelName", None)

    replacements = REPLACEMENT_VALUES.copy()
    replacements["ENDPOINT_CONFIG_NAME"] = config_resource_name
    replacements["MODEL_NAME"] = model_resource_name

    config_reference, config_spec, config_resource = create_sagemaker_resource(
        resource_plural=cfg.ENDPOINT_CONFIG_RESOURCE_PLURAL,
        resource_name=config_resource_name,
        spec_file="endpoint_config_multi_variant",
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
def single_variant_config(name_suffix, single_container_model):
    config_resource_name = name_suffix + "-single-variant-config"
    (_, model_resource) = single_container_model
    model_resource_name = model_resource["spec"].get("modelName", None)

    replacements = REPLACEMENT_VALUES.copy()
    replacements["ENDPOINT_CONFIG_NAME"] = config_resource_name
    replacements["MODEL_NAME"] = model_resource_name

    config_reference, config_spec, config_resource = create_sagemaker_resource(
        resource_plural=cfg.ENDPOINT_CONFIG_RESOURCE_PLURAL,
        resource_name=config_resource_name,
        spec_file="endpoint_config_single_variant",
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
def xgboost_endpoint(name_suffix, single_variant_config):
    endpoint_resource_name = name_suffix
    (_, config_resource) = single_variant_config
    config_resource_name = config_resource["spec"].get("endpointConfigName", None)

    replacements = REPLACEMENT_VALUES.copy()
    replacements["ENDPOINT_NAME"] = endpoint_resource_name
    replacements["ENDPOINT_CONFIG_NAME"] = config_resource_name

    reference, spec, resource = create_sagemaker_resource(
        resource_plural=cfg.ENDPOINT_RESOURCE_PLURAL,
        resource_name=endpoint_resource_name,
        spec_file="endpoint_base",
        replacements=replacements,
    )

    assert resource is not None

    yield (reference, resource, spec)

    # Delete the k8s resource if not already deleted by tests
    assert delete_custom_resource(reference, 40, cfg.DELETE_WAIT_LENGTH)


@pytest.fixture(scope="module")
def faulty_config(name_suffix, single_container_model):
    replacements = REPLACEMENT_VALUES.copy()

    # copy model data to a temp S3 location and delete it after model is created on SageMaker
    model_bucket = replacements["SAGEMAKER_DATA_BUCKET"]
    copy_source = {
        "Bucket": model_bucket,
        "Key": "sagemaker/model/xgboost-mnist-model.tar.gz",
    }
    model_destination_key = "sagemaker/model/delete/xgboost-mnist-model.tar.gz"
    s3.copy_object(model_bucket, copy_source, model_destination_key)

    model_resource_name = name_suffix + "faulty-model"
    replacements["MODEL_NAME"] = model_resource_name
    replacements["MODEL_LOCATION"] = f"s3://{model_bucket}/{model_destination_key}"
    model_reference, model_spec, model_resource = create_sagemaker_resource(
        resource_plural=cfg.MODEL_RESOURCE_PLURAL,
        resource_name=model_resource_name,
        spec_file="xgboost_model_with_model_location",
        replacements=replacements,
    )
    assert model_resource is not None
    if k8s.get_resource_arn(model_resource) is None:
        logging.error(
            f"ARN for this resource is None, resource status is: {model_resource['status']}"
        )
    assert k8s.get_resource_arn(model_resource) is not None
    s3.delete_object(model_bucket, model_destination_key)

    config_resource_name = name_suffix + "-faulty-config"
    (_, model_resource) = single_container_model
    model_resource_name = model_resource["spec"].get("modelName", None)

    replacements["ENDPOINT_CONFIG_NAME"] = config_resource_name

    config_reference, config_spec, config_resource = create_sagemaker_resource(
        resource_plural=cfg.ENDPOINT_CONFIG_RESOURCE_PLURAL,
        resource_name=config_resource_name,
        spec_file="endpoint_config_multi_variant",
        replacements=replacements,
    )
    assert config_resource is not None
    if k8s.get_resource_arn(config_resource) is None:
        logging.error(
            f"ARN for this resource is None, resource status is: {config_resource['status']}"
        )
    assert k8s.get_resource_arn(config_resource) is not None

    yield (config_reference, config_resource)

    for cr in (model_reference, config_reference):
        _, deleted = k8s.delete_custom_resource(cr, cfg.DELETE_WAIT_PERIOD, cfg.DELETE_WAIT_LENGTH)
        assert deleted


@service_marker
@pytest.mark.shallow_canary
@pytest.mark.canary
@pytest.skip("temp")
class TestEndpoint:
    def create_endpoint_test(self, xgboost_endpoint):
        (reference, resource, _) = xgboost_endpoint
        assert k8s.get_resource_exists(reference)

        # endpoint has correct arn and status
        endpoint_name = resource["spec"].get("endpointName", None)
        assert endpoint_name is not None

        endpoint_desc = get_sagemaker_endpoint(endpoint_name)
        endpoint_arn = endpoint_desc["EndpointArn"]
        assert k8s.get_resource_arn(resource) == endpoint_arn

        # endpoint transitions Creating -> InService state
        assert_endpoint_status_in_sync(endpoint_name, reference, cfg.ENDPOINT_STATUS_CREATING)
        assert k8s.wait_on_condition(
            reference, ack_condition.CONDITION_TYPE_RESOURCE_SYNCED, "False"
        )

        assert_endpoint_status_in_sync(endpoint_name, reference, cfg.ENDPOINT_STATUS_INSERVICE)
        assert k8s.wait_on_condition(
            reference, ack_condition.CONDITION_TYPE_RESOURCE_SYNCED, "True"
        )

        resource_tags = resource["spec"].get("tags", None)
        assert_tags_in_sync(endpoint_arn, resource_tags)

    def update_endpoint_failed_test(self, single_variant_config, faulty_config, xgboost_endpoint):
        (endpoint_reference, _, endpoint_spec) = xgboost_endpoint
        (_, faulty_config_resource) = faulty_config
        faulty_config_name = faulty_config_resource["spec"].get("endpointConfigName", None)
        endpoint_spec["spec"]["endpointConfigName"] = faulty_config_name
        endpoint_resource = k8s.patch_custom_resource(endpoint_reference, endpoint_spec)
        endpoint_resource = k8s.wait_resource_consumed_by_controller(endpoint_reference)
        assert endpoint_resource is not None

        # endpoint transitions Updating -> InService state
        assert_endpoint_status_in_sync(
            endpoint_reference.name,
            endpoint_reference,
            cfg.ENDPOINT_STATUS_UPDATING,
        )
        assert k8s.wait_on_condition(
            endpoint_reference, ack_condition.CONDITION_TYPE_RESOURCE_SYNCED, "False"
        )
        endpoint_resource = k8s.get_resource(endpoint_reference)
        annotations = endpoint_resource["metadata"].get("annotations", None)
        assert annotations is not None
        assert annotations[LAST_ENDPOINTCONFIG_UPDATE_ANNOTATION] == faulty_config_name

        assert_endpoint_status_in_sync(
            endpoint_reference.name,
            endpoint_reference,
            cfg.ENDPOINT_STATUS_INSERVICE,
        )

        assert k8s.wait_on_condition(
            endpoint_reference, ack_condition.CONDITION_TYPE_RESOURCE_SYNCED, "False"
        )

        (_, old_config_resource) = single_variant_config
        current_config_name = old_config_resource["spec"].get("endpointConfigName", None)
        assert k8s.assert_condition_state_message(
            endpoint_reference,
            ack_condition.CONDITION_TYPE_TERMINAL,
            "True",
            FAIL_UPDATE_ERROR_MESSAGE + current_config_name,
        )

        endpoint_resource = k8s.get_resource(endpoint_reference)
        assert endpoint_resource["status"].get("failureReason", None) is not None

    def update_endpoint_successful_test(self, multi_variant_config, xgboost_endpoint):
        (endpoint_reference, endpoint_resource, endpoint_spec) = xgboost_endpoint

        endpoint_name = endpoint_resource["spec"].get("endpointName", None)
        production_variants = get_sagemaker_endpoint(endpoint_name)["ProductionVariants"]
        old_variant_instance_count = production_variants[0]["CurrentInstanceCount"]
        old_variant_name = production_variants[0]["VariantName"]

        (_, new_config_resource) = multi_variant_config
        new_config_name = new_config_resource["spec"].get("endpointConfigName", None)
        endpoint_spec["spec"]["endpointConfigName"] = new_config_name
        endpoint_resource = k8s.patch_custom_resource(endpoint_reference, endpoint_spec)
        endpoint_resource = k8s.wait_resource_consumed_by_controller(endpoint_reference)
        assert endpoint_resource is not None

        # endpoint transitions Updating -> InService state
        assert_endpoint_status_in_sync(
            endpoint_reference.name,
            endpoint_reference,
            cfg.ENDPOINT_STATUS_UPDATING,
        )

        assert k8s.wait_on_condition(
            endpoint_reference, ack_condition.CONDITION_TYPE_RESOURCE_SYNCED, "False"
        )
        assert (
            k8s.get_resource_condition(endpoint_reference, ack_condition.CONDITION_TYPE_TERMINAL)
            is None
        )
        endpoint_resource = k8s.get_resource(endpoint_reference)
        annotations = endpoint_resource["metadata"].get("annotations", None)
        assert annotations is not None
        assert annotations[LAST_ENDPOINTCONFIG_UPDATE_ANNOTATION] == new_config_name

        assert_endpoint_status_in_sync(
            endpoint_reference.name,
            endpoint_reference,
            cfg.ENDPOINT_STATUS_INSERVICE,
        )
        assert k8s.wait_on_condition(
            endpoint_reference, ack_condition.CONDITION_TYPE_RESOURCE_SYNCED, "True"
        )
        assert (
            k8s.get_resource_condition(endpoint_reference, ack_condition.CONDITION_TYPE_TERMINAL)
            is None
        )
        endpoint_resource = k8s.get_resource(endpoint_reference)
        assert endpoint_resource["status"].get("failureReason", None) is None

        # RetainAllVariantProperties - variant properties were retained + is a multi-variant endpoint
        new_production_variants = get_sagemaker_endpoint(endpoint_name)["ProductionVariants"]
        assert len(new_production_variants) > 1
        new_variant_instance_count = None
        for variant in new_production_variants:
            if variant["VariantName"] == old_variant_name:
                new_variant_instance_count = variant["CurrentInstanceCount"]

        assert new_variant_instance_count == old_variant_instance_count

    def delete_endpoint_test(self, xgboost_endpoint):
        (reference, resource, _) = xgboost_endpoint
        endpoint_name = resource["spec"].get("endpointName", None)

        assert delete_custom_resource(reference, cfg.DELETE_WAIT_PERIOD, cfg.DELETE_WAIT_LENGTH)

        assert get_sagemaker_endpoint(endpoint_name) is None

    def test_driver(
        self,
        sagemaker_client,
        single_variant_config,
        faulty_config,
        multi_variant_config,
        xgboost_endpoint,
    ):
        self.create_endpoint_test(xgboost_endpoint)
        self.update_endpoint_failed_test(single_variant_config, faulty_config, xgboost_endpoint)
        # Note: the test has been intentionally ordered to run a successful update after a failed update
        # check that controller updates the endpoint, removes the terminal condition and clears the failure reason
        self.update_endpoint_successful_test(multi_variant_config, xgboost_endpoint)
        self.delete_endpoint_test(xgboost_endpoint)
