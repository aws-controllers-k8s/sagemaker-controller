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
"""Integration tests for the SageMaker DataQualityJobDefinition API.
"""

import pytest
import logging
import botocore

from e2e import service_marker
from e2e.common.fixtures import (
    xgboost_churn_data_quality_job_definition,
    xgboost_churn_endpoint,
)
from acktest.k8s import resource as k8s

# Access variable so it is loaded as a fixture
_accessed = xgboost_churn_data_quality_job_definition, xgboost_churn_endpoint


def describe_sagemaker_data_quality_job_definition(
    sagemaker_client, job_definition_name
):
    try:
        return sagemaker_client.describe_data_quality_job_definition(
            JobDefinitionName=job_definition_name
        )
    except botocore.exceptions.ClientError as error:
        logging.error(
            f"Could not find Data Quality Job Definition with name {job_definition_name}. Error {error}"
        )
        return None


@service_marker
@pytest.mark.canary
class TestDataQualityJobDefinition:
    def test_smoke(self, sagemaker_client, xgboost_churn_data_quality_job_definition):
        (reference, resource) = xgboost_churn_data_quality_job_definition
        assert k8s.get_resource_exists(reference)

        job_definition_name = resource["spec"].get("jobDefinitionName")
        assert (
            k8s.get_resource_arn(resource)
            == describe_sagemaker_data_quality_job_definition(
                sagemaker_client, job_definition_name
            )["JobDefinitionArn"]
        )

        # Delete the k8s resource.
        _, deleted = k8s.delete_custom_resource(reference, 3, 10)
        assert deleted
        assert (
            describe_sagemaker_data_quality_job_definition(
                sagemaker_client, job_definition_name
            )
            is None
        )
