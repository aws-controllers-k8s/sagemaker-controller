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

import botocore
import pytest
import logging
import time
from typing import Dict

from acktest.aws import s3
from acktest.resources import random_suffix_name
from acktest.k8s import resource as k8s

from e2e import (
    service_marker,
    create_sagemaker_resource,
    assert_endpoint_status_in_sync,
)
from e2e.replacement_values import REPLACEMENT_VALUES
from e2e.common import config as cfg

FAIL_UPDATE_ERROR_MESSAGE = "unable to update endpoint. check FailureReason"


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
    assert k8s.get_resource_arn(model_resource) is not None

    yield (model_reference, model_resource)

    _, deleted = k8s.delete_custom_resource(model_reference, 3, 10)
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
    assert k8s.get_resource_arn(config_resource) is not None

    yield (config_reference, config_resource)

    _, deleted = k8s.delete_custom_resource(config_reference, 3, 10)
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
    assert k8s.get_resource_arn(config_resource) is not None

    yield (config_reference, config_resource)

    _, deleted = k8s.delete_custom_resource(config_reference, 3, 10)
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
    if k8s.get_resource_exists(reference):
        _, deleted = k8s.delete_custom_resource(reference, 3, 10)
        assert deleted


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
    model_resource = k8s.get_resource(model_reference)
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
    assert k8s.get_resource_arn(config_resource) is not None

    yield (config_reference, config_resource)

    for cr in (model_reference, config_reference):
        _, deleted = k8s.delete_custom_resource(cr, 3, 10)
        assert deleted


@service_marker
@pytest.mark.canary
class TestEndpoint:
    def _get_resource_endpoint_arn(self, resource: Dict):
        assert (
            "ackResourceMetadata" in resource["status"]
            and "arn" in resource["status"]["ackResourceMetadata"]
        )
        return resource["status"]["ackResourceMetadata"]["arn"]

    def _describe_sagemaker_endpoint(self, sagemaker_client, endpoint_name: str):
        try:
            return sagemaker_client.describe_endpoint(EndpointName=endpoint_name)
        except botocore.exceptions.ClientError as error:
            logging.error(
                f"SageMaker could not find a endpoint with the name {endpoint_name}. Error {error}"
            )
            return None

    def create_endpoint_test(self, sagemaker_client, xgboost_endpoint):
        (reference, resource, _) = xgboost_endpoint
        assert k8s.get_resource_exists(reference)

        # endpoint has correct arn and status
        endpoint_name = resource["spec"].get("endpointName", None)
        assert endpoint_name is not None

        assert (
            self._get_resource_endpoint_arn(resource)
            == self._describe_sagemaker_endpoint(sagemaker_client, endpoint_name)[
                "EndpointArn"
            ]
        )

        # endpoint transitions Creating -> InService state
        assert_endpoint_status_in_sync(
            endpoint_name, reference, cfg.ENDPOINT_STATUS_CREATING
        )
        assert k8s.wait_on_condition(reference, "ACK.ResourceSynced", "False")

        assert_endpoint_status_in_sync(
            endpoint_name, reference, cfg.ENDPOINT_STATUS_INSERVICE
        )
        assert k8s.wait_on_condition(reference, "ACK.ResourceSynced", "True")

    def update_endpoint_failed_test(
        self, sagemaker_client, single_variant_config, faulty_config, xgboost_endpoint
    ):
        (endpoint_reference, _, endpoint_spec) = xgboost_endpoint
        (_, faulty_config_resource) = faulty_config
        faulty_config_name = faulty_config_resource["spec"].get(
            "endpointConfigName", None
        )
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
        assert k8s.wait_on_condition(endpoint_reference, "ACK.ResourceSynced", "False")
        endpoint_resource = k8s.get_resource(endpoint_reference)
        assert (
            endpoint_resource["status"].get("lastEndpointConfigNameForUpdate", None)
            == faulty_config_name
        )

        assert_endpoint_status_in_sync(
            endpoint_reference.name,
            endpoint_reference,
            cfg.ENDPOINT_STATUS_INSERVICE,
        )

        assert k8s.wait_on_condition(endpoint_reference, "ACK.ResourceSynced", "True")
        assert k8s.assert_condition_state_message(
            endpoint_reference,
            "ACK.Terminal",
            "True",
            FAIL_UPDATE_ERROR_MESSAGE,
        )

        endpoint_resource = k8s.get_resource(endpoint_reference)
        assert endpoint_resource["status"].get("failureReason", None) is not None

        # additional check: endpoint using old endpoint config
        (_, old_config_resource) = single_variant_config
        current_config_name = endpoint_resource["status"].get(
            "latestEndpointConfigName"
        )
        assert (
            current_config_name is not None
            and current_config_name
            == old_config_resource["spec"].get("endpointConfigName", None)
        )

    def update_endpoint_successful_test(
        self, sagemaker_client, multi_variant_config, xgboost_endpoint
    ):
        (endpoint_reference, endpoint_resource, endpoint_spec) = xgboost_endpoint

        endpoint_name = endpoint_resource["spec"].get("endpointName", None)
        production_variants = self._describe_sagemaker_endpoint(
            sagemaker_client, endpoint_name
        )["ProductionVariants"]
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

        assert k8s.wait_on_condition(endpoint_reference, "ACK.ResourceSynced", "False")
        assert k8s.assert_condition_state_message(
            endpoint_reference, "ACK.Terminal", "False", None
        )
        endpoint_resource = k8s.get_resource(endpoint_reference)
        assert (
            endpoint_resource["status"].get("lastEndpointConfigNameForUpdate", None)
            == new_config_name
        )

        assert_endpoint_status_in_sync(
            endpoint_reference.name,
            endpoint_reference,
            cfg.ENDPOINT_STATUS_INSERVICE,
        )
        assert k8s.wait_on_condition(endpoint_reference, "ACK.ResourceSynced", "True")
        assert k8s.assert_condition_state_message(
            endpoint_reference, "ACK.Terminal", "False", None
        )
        endpoint_resource = k8s.get_resource(endpoint_reference)
        assert endpoint_resource["status"].get("failureReason", None) is None

        # RetainAllVariantProperties - variant properties were retained + is a multi-variant endpoint
        new_production_variants = self._describe_sagemaker_endpoint(
            sagemaker_client, endpoint_name
        )["ProductionVariants"]
        assert len(new_production_variants) > 1
        new_variant_instance_count = None
        for variant in new_production_variants:
            if variant["VariantName"] == old_variant_name:
                new_variant_instance_count = variant["CurrentInstanceCount"]

        assert new_variant_instance_count == old_variant_instance_count

    def delete_endpoint_test(self, sagemaker_client, xgboost_endpoint):
        (reference, resource, _) = xgboost_endpoint
        endpoint_name = resource["spec"].get("endpointName", None)

        _, deleted = k8s.delete_custom_resource(reference, 3, 10)
        assert deleted

        # resource is removed from management from controller side if call to deleteEndpoint succeeds.
        # Sagemaker also removes a 'Deleting' endpoint pretty quickly, but there might be a lag
        # We wait maximum of 30 seconds for this clean up to happen
        endpoint_deleted = False
        for _ in range(3):
            time.sleep(10)
            if (
                self._describe_sagemaker_endpoint(sagemaker_client, endpoint_name)
                is None
            ):
                endpoint_deleted = True
                break
        assert endpoint_deleted

    def test_driver(
        self,
        sagemaker_client,
        single_variant_config,
        faulty_config,
        multi_variant_config,
        xgboost_endpoint,
    ):
        self.create_endpoint_test(sagemaker_client, xgboost_endpoint)
        self.update_endpoint_failed_test(
            sagemaker_client, single_variant_config, faulty_config, xgboost_endpoint
        )
        # Note: the test has been intentionally ordered to run a successful update after a failed update
        # check that controller updates the endpoint, removes the terminal condition and clears the failure reason
        self.update_endpoint_successful_test(
            sagemaker_client, multi_variant_config, xgboost_endpoint
        )
        self.delete_endpoint_test(sagemaker_client, xgboost_endpoint)
