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
"""Integration test for ACKs Cross Region Support.
"""

import logging
from acktest.aws.identity import get_region
import boto3
import botocore
import pytest
import time
from typing import Dict

from acktest.resources import random_suffix_name
from acktest.k8s import resource as k8s
from e2e import (
    service_marker,
    create_sagemaker_resource,
    assert_tags_in_sync,
)
from e2e.replacement_values import REPLACEMENT_VALUES, XGBOOST_V1_IMAGE_URIS
from e2e.common import config as cfg


def get_cross_region():
    region = get_region()
    if region != "us-west-2":
        region = "us-west-2"
    else:
        region = "us-east-1"
    return region


@pytest.fixture(scope="module")
def cross_region_model():
    resource_name = random_suffix_name("cross-region-model", 32)
    region = get_cross_region()

    replacements = REPLACEMENT_VALUES.copy()
    replacements["MODEL_NAME"] = resource_name
    replacements["REGION"] = region
    replacements[
        "XGBOOST_V1_IMAGE_URI"
    ] = f"{XGBOOST_V1_IMAGE_URIS[region]}/sagemaker-xgboost:1.0-1-cpu-py3"

    reference, spec, resource = create_sagemaker_resource(
        resource_plural=cfg.MODEL_RESOURCE_PLURAL,
        resource_name=resource_name,
        spec_file="cross_region_model",
        replacements=replacements,
    )
    assert resource is not None

    yield (reference, resource)

    # Delete the k8s resource if not already deleted by tests
    if k8s.get_resource_exists(reference):
        _, deleted = k8s.delete_custom_resource(reference, 3, 10)
        assert deleted


def get_sagemaker_cross_region_model(model_name, sagemaker_client):
    try:
        return sagemaker_client.describe_model(ModelName=model_name)
    except botocore.exceptions.ClientError as error:
        logging.error(
            f"SageMaker could not find a model with the name {model_name}. Error {error}"
        )
        return None


@service_marker
class TestCrossRegionModel:
    def test_create_cross_region_model(self, cross_region_model):
        (reference, resource) = cross_region_model
        assert k8s.get_resource_exists(reference)

        sagemaker_client = boto3.client("sagemaker", region=get_cross_region())
        model_name = resource["spec"].get("modelName", None)
        model_desc = get_sagemaker_cross_region_model(
            model_name, sagemaker_client)
        cross_region_model_arn = model_desc["ModelArn"]
        assert k8s.get_resource_arn(resource) == cross_region_model_arn

        # random sleep before we check for tags to reduce test flakyness
        time.sleep(cfg.TAG_DELAY_SLEEP)
        resource_tags = resource["spec"].get("tags", None)
        assert_tags_in_sync(cross_region_model_arn, resource_tags)

        # Delete the k8s resource.
        _, deleted = k8s.delete_custom_resource(reference, 3, 10)
        assert deleted

        assert get_sagemaker_cross_region_model(
            model_name, sagemaker_client) is None
