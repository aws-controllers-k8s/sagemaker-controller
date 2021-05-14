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
"""Integration tests for the SageMaker ModelExplainabilityJobDefinition API.
"""

import botocore
import pytest
import logging

from e2e import (
    service_marker,
    create_sagemaker_resource,
)
from e2e.replacement_values import REPLACEMENT_VALUES
from e2e.common.fixtures import xgboost_churn_endpoint
from acktest.resources import random_suffix_name
from acktest.k8s import resource as k8s

RESOURCE_PLURAL = "modelexplainabilityjobdefinitions"

# Access variable so it is loaded as a fixture
_accessed = xgboost_churn_endpoint


@pytest.fixture(scope="module")
def xgboost_churn_model_explainability_job_definition(xgboost_churn_endpoint):
    endpoint_spec = xgboost_churn_endpoint
    endpoint_name = endpoint_spec["spec"].get("endpointName")

    job_definition_name = random_suffix_name("model-explain-job-definition", 32)
    replacements = REPLACEMENT_VALUES.copy()
    replacements["JOB_DEFINITION_NAME"] = job_definition_name
    replacements["ENDPOINT_NAME"] = endpoint_name

    reference, _, resource = create_sagemaker_resource(
        resource_plural=RESOURCE_PLURAL,
        resource_name=job_definition_name,
        spec_file="model_explainability_job_definition_xgboost_churn",
        replacements=replacements,
    )
    assert resource is not None

    yield (reference, resource)

    if k8s.get_resource_exists(reference):
        _, deleted = k8s.delete_custom_resource(job_definition_reference, 3, 10)
        assert deleted


def describe_sagemaker_model_explainability_job_definition(
    sagemaker_client, job_definition_name
):
    try:
        return sagemaker_client.describe_model_explainability_job_definition(
            JobDefinitionName=job_definition_name
        )
    except botocore.exceptions.ClientError as error:
        logging.error(
            f"Could not find Model Explainability Job Definition with name {job_definition_name}. Error {error}"
        )
        return None


@service_marker
@pytest.mark.canary
class TestModelExplainabilityJobDefinition:
    def test_smoke(
        self, sagemaker_client, xgboost_churn_model_explainability_job_definition
    ):
        (reference, resource) = xgboost_churn_model_explainability_job_definition
        assert k8s.get_resource_exists(reference)

        job_definition_name = resource["spec"].get("jobDefinitionName")
        assert (
            k8s.get_resource_arn(resource)
            == describe_sagemaker_model_explainability_job_definition(
                sagemaker_client, job_definition_name
            )["JobDefinitionArn"]
        )

        # Delete the k8s resource.
        _, deleted = k8s.delete_custom_resource(reference, 3, 10)
        assert deleted
        assert (
            describe_sagemaker_model_explainability_job_definition(
                sagemaker_client, job_definition_name
            )
            is None
        )
