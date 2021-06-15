
# Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"). You may
# not use this file except in compliance with the License. A copy of the
# License is located at
#
#  http://aws.amazon.com/apache2.0/
#
# or in the "license" file accompanying this file. This file is distributed
# on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
# express or implied. See the License for the specific language governing
# permissions and limitations under the License.
"""Integration tests for the SageMaker Endpoint API.
"""

import boto3
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
#from e2e.common import config as cfg

RESOURCE_NAME_BASE = "feature-group"
RESOURCE_PLURAL = "featuregroups"
SPEC_FILE = "feature_group"

@pytest.fixture(scope="module")
def feature_group():
    feature_group_name = random_suffix_name(RESOURCE_NAME_BASE, 32)
    replacements = REPLACEMENT_VALUES.copy()
    #print("About to create a sagemaker resource!\n")
    #print("Input resource_plural: " + RESOURCE_PLURAL + "\n")
    #print("Input resource_name: " + resource_name + "\n")
    #print("Input spec_file: " + SPEC_FILE + "\n")
    #print("Input replacements: " + replacements + "\n")
    replacements["FEATURE_GROUP_NAME"] = feature_group_name
    reference, spec, resource = create_sagemaker_resource(
        resource_plural=RESOURCE_PLURAL,
        resource_name=feature_group_name,
        spec_file=SPEC_FILE,
        replacements=replacements,
    )
    #print("Created reference, spec, resource!\n")
    #print("Reference: " + reference + "\n")
    #print("Spec: " + spec + "\n")
    #print("Resource: " + resource + "\n")
    #print("About to check that resource is not none!\n")
    assert resource is not None
    #print("Resource is NOT NONE!\n")
    yield (reference, resource)
    
    # Delete the k8s resource if not already deleted by tests
    # At 10 a second wait period instead of 15, we sometimes see time out errors.
    if k8s.get_resource_exists(reference):
        _, deleted = k8s.delete_custom_resource(reference, 3, 15)
        assert deleted

def get_sagemaker_feature_group(feature_group_name: str):
    try:
        return sagemaker_client().describe_feature_group(FeatureGroupName=feature_group_name)
    except botocore.exceptions.ClientError as error:
        logging.error(
            f"SageMaker could not find a feature group with the name {feature_group_name}. Error {error}"
        )
        return None

        
@service_marker
@pytest.mark.canary
class TestFeatureGroup:
    def test_create_feature_group(self, feature_group):
        (reference, resource) = feature_group
        assert k8s.get_resource_exists(reference)
        
        feature_group_name = resource["spec"].get("featureGroupName", None)
        
        assert (
            k8s.get_resource_arn(resource)
            == get_sagemaker_feature_group(feature_group_name)["FeatureGroupArn"]
        )
        
        # Delete the k8s resource.
        # At 10 a second wait period instead of 15, we sometimes see time out errors.
        _, deleted = k8s.delete_custom_resource(reference, 3, 15)
        assert deleted
        
        assert get_sagemaker_feature_group(feature_group_name) is None
        
