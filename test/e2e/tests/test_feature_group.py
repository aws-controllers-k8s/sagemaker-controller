# Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"). You may
# not use this file except in compliance with the License. A copy of the
# License is located at
#
# http://aws.amazon.com/apache2.0/
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
        wait_for_status,
        sagemaker_client
)
from e2e.replacement_values import REPLACEMENT_VALUES

RESOURCE_NAME_PREFIX = "feature-group"
RESOURCE_PLURAL = "featuregroups"
SPEC_FILE = "feature_group"
FEATURE_GROUP_STATUS_CREATED = "CREATED"
WAIT_PERIOD_COUNT = 3
# A 15 second wait period is used because we sometimes see
# time out errors at a 10 second wait period.
WAIT_PERIOD_LENGTH = 15

@pytest.fixture(scope="module")
def feature_group():
    """Creates a feature group from a SPEC_FILE."""
    feature_group_name = random_suffix_name(RESOURCE_NAME_PREFIX, 32)
    replacements = REPLACEMENT_VALUES.copy()
    replacements["FEATURE_GROUP_NAME"] = feature_group_name
    reference, spec, resource = create_sagemaker_resource(
        resource_plural=RESOURCE_PLURAL,
        resource_name=feature_group_name,
        spec_file=SPEC_FILE,
        replacements=replacements,
    )
    assert resource is not None
    yield (reference, resource)
    
    # Delete the k8s resource if not already deleted by tests
    if k8s.get_resource_exists(reference):
        _, deleted = k8s.delete_custom_resource(reference, WAIT_PERIOD_COUNT, WAIT_PERIOD_LENGTH)
        assert deleted
        
def get_sagemaker_feature_group(feature_group_name: str):
    """Used to check if there is an existing feature group with a given feature_group_name."""
    try:
        return sagemaker_client().describe_feature_group(FeatureGroupName=feature_group_name)
    except botocore.exceptions.ClientError as error:
        logging.error(
            f"SageMaker could not find a feature group with the name {feature_group_name}. Error {error}"
        )
        return None
    
def get_feature_group_status(feature_group_name: str):
    feature_group_describe_response = get_sagemaker_feature_group(feature_group_name)
    return feature_group_describe_response["FeatureGroupStatus"]

@service_marker
@pytest.mark.canary
class TestFeatureGroup:
    def _wait_feature_group_status(
            self,
            feature_group_name,
            expected_status: str,
            wait_periods: int = WAIT_PERIOD_COUNT,
            period_length: int = WAIT_PERIOD_LENGTH,
    ):
        return wait_for_status(
            expected_status,
            wait_periods,
            period_length,
            get_feature_group_status,
            feature_group_name,
        )
    
    def test_create_feature_group(self, feature_group):
        """Tests that a feature group can be created and deleted
        using the Feature Group Controller.
        """
        (reference, resource) = feature_group
        assert k8s.get_resource_exists(reference)
        
        feature_group_name = resource["spec"].get("featureGroupName", None)
        
        feature_group_describe_response = get_sagemaker_feature_group(feature_group_name)
        
        assert (
            k8s.get_resource_arn(resource)
            == feature_group_describe_response["FeatureGroupArn"]
        )
        
        assert (
            self._wait_feature_group_status(
                feature_group_name, FEATURE_GROUP_STATUS_CREATED
            )
        )
        # TODO: add resource side checks.
        
        # Delete the k8s resource.
        _, deleted = k8s.delete_custom_resource(reference, WAIT_PERIOD_COUNT, WAIT_PERIOD_LENGTH)
        assert deleted
        # TODO: Implement logic to requeueOnDelete.
        # TODO: Once the delete requeue PR is merged,
        # verify that it works for DeleteFailed state.
        
        assert get_sagemaker_feature_group(feature_group_name) is None
