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
"""Common SageMaker test fixtures.
"""

import pytest

from e2e import (
    create_sagemaker_resource,
    wait_sagemaker_endpoint_status,
)

from e2e.replacement_values import REPLACEMENT_VALUES
from acktest.resources import random_suffix_name
from acktest.k8s import resource as k8s
from e2e.common import config as cfg


@pytest.fixture(scope="module")
def xgboost_churn_endpoint(sagemaker_client):
    """Creates a SageMaker endpoint with the XGBoost churn single-variant model
    and data capture enabled.
    """
    endpoint_resource_name = random_suffix_name("xgboost-churn", 32)
    endpoint_config_resource_name = endpoint_resource_name + "-config"
    model_resource_name = endpoint_config_resource_name + "-model"

    replacements = REPLACEMENT_VALUES.copy()
    replacements["ENDPOINT_NAME"] = endpoint_resource_name
    replacements["ENDPOINT_CONFIG_NAME"] = endpoint_config_resource_name
    replacements["MODEL_NAME"] = model_resource_name
    data_bucket = replacements["SAGEMAKER_DATA_BUCKET"]
    replacements[
        "MODEL_LOCATION"
    ] = f"s3://{data_bucket}/sagemaker/model/xgb-churn-prediction-model.tar.gz"

    model_reference, model_spec, model_resource = create_sagemaker_resource(
        resource_plural=cfg.MODEL_RESOURCE_PLURAL,
        resource_name=model_resource_name,
        spec_file="xgboost_model_with_model_location",
        replacements=replacements,
    )
    assert model_resource is not None
    assert k8s.get_resource_arn(model_resource) is not None

    (
        endpoint_config_reference,
        endpoint_config_spec,
        endpoint_config_resource,
    ) = create_sagemaker_resource(
        resource_plural=cfg.ENDPOINT_CONFIG_RESOURCE_PLURAL,
        resource_name=endpoint_config_resource_name,
        spec_file="endpoint_config_data_capture_single_variant",
        replacements=replacements,
    )
    assert endpoint_config_resource is not None
    assert k8s.get_resource_arn(endpoint_config_resource) is not None

    endpoint_reference, endpoint_spec, endpoint_resource = create_sagemaker_resource(
        resource_plural=cfg.ENDPOINT_RESOURCE_PLURAL,
        resource_name=endpoint_resource_name,
        spec_file="endpoint_base",
        replacements=replacements,
    )
    assert endpoint_resource is not None
    assert k8s.get_resource_arn(endpoint_resource) is not None
    wait_sagemaker_endpoint_status(
        replacements["ENDPOINT_NAME"], "InService"
    )

    yield endpoint_spec

    for cr in (model_reference, endpoint_config_reference, endpoint_reference):
        _, deleted = k8s.delete_custom_resource(cr, 3, 10)
        assert deleted


@pytest.fixture(scope="module")
def xgboost_churn_data_quality_job_definition(xgboost_churn_endpoint):
    endpoint_spec = xgboost_churn_endpoint
    endpoint_name = endpoint_spec["spec"].get("endpointName")

    resource_name = random_suffix_name("data-quality-job-definition", 32)

    replacements = REPLACEMENT_VALUES.copy()
    replacements["JOB_DEFINITION_NAME"] = resource_name
    replacements["ENDPOINT_NAME"] = endpoint_name

    job_definition_reference, _, resource = create_sagemaker_resource(
        resource_plural=cfg.DATA_QUALITY_JOB_DEFINITION_RESOURCE_PLURAL,
        resource_name=resource_name,
        spec_file="data_quality_job_definition_xgboost_churn",
        replacements=replacements,
    )
    assert resource is not None

    job_definition_name = resource["spec"].get("jobDefinitionName")

    yield (job_definition_reference, resource)

    if k8s.get_resource_exists(job_definition_reference):
        _, deleted = k8s.delete_custom_resource(job_definition_reference, 3, 10)
        assert deleted
