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

from acktest.aws.identity import get_region
from . import tags
import pytest
import logging
import botocore
import time
import boto3
from pathlib import Path
from unittest import TestCase
import random

from acktest.k8s import resource as k8s
from e2e.common import config as cfg


SERVICE_NAME = "sagemaker"
CRD_GROUP = "sagemaker.services.k8s.aws"
ADOPTED_RESOURCE_CRD_GROUP = "services.k8s.aws"
CRD_VERSION = "v1alpha1"

# PyTest marker for the current service
service_marker = pytest.mark.service(arg=SERVICE_NAME)

bootstrap_directory = Path(__file__).parent
resource_directory = Path(__file__).parent / "resources"


def sagemaker_client(region: str = None):
    region = region or get_region()
    return boto3.client("sagemaker", region_name=region)


def create_sagemaker_resource(
    resource_plural,
    resource_name,
    spec_file,
    replacements,
    namespace="default",
    wait_period=3,
    period_length=10,
):
    """
    Wrapper around k8s.load_and_create_resource to create a SageMaker resource
    """
    # Add a random sleep to prevent throttling exception before the call to load and create
    # This is because there may be many of the same resource ex: Multiple Models being created at the same time
    # If this occurs then a throttling exception may occur and cause the tests to fail, this sleep prevents this from occurring.
    rand = random.randrange(1, 4)
    time.sleep(rand)
    reference, spec, resource = k8s.load_and_create_resource(
        resource_directory,
        CRD_GROUP,
        CRD_VERSION,
        resource_plural,
        resource_name,
        spec_file,
        replacements,
        namespace,
        wait_period,
        period_length,
    )

    return reference, spec, resource


def create_adopted_resource(
    replacements, namespace="default", spec_file="adopted_resource_base"
):
    """
    Wrapper around k8s.load_and_create_resource to create a Adopoted resource
    """

    reference, spec, resource = k8s.load_and_create_resource(
        resource_directory,
        ADOPTED_RESOURCE_CRD_GROUP,
        CRD_VERSION,
        "adoptedresources",
        replacements["ADOPTED_RESOURCE_NAME"],
        spec_file,
        replacements,
        namespace,
    )

    return reference, spec, resource


def wait_for_status(
    expected_status: str,
    wait_periods: int,
    period_length: int,
    get_status_method,
    *method_args,
):
    actual_status = None
    for _ in range(wait_periods):
        time.sleep(period_length)
        actual_status = get_status_method(*method_args)
        if actual_status == expected_status:
            break
    else:
        logging.error(
            f"Wait for status: {expected_status} timed out. Actual status: {actual_status}"
        )

    return actual_status


def get_endpoint_sagemaker_status(endpoint_name):
    response = sagemaker_client().describe_endpoint(EndpointName=endpoint_name)
    return response["EndpointStatus"]


def get_endpoint_resource_status(reference: k8s.CustomResourceReference):
    resource = k8s.get_resource(reference)
    assert "endpointStatus" in resource["status"]
    return resource["status"]["endpointStatus"]


def wait_sagemaker_endpoint_status(
    endpoint_name,
    expected_status: str,
    wait_periods: int = 60,
    period_length: int = 30,
):
    return wait_for_status(
        expected_status,
        wait_periods,
        period_length,
        get_endpoint_sagemaker_status,
        endpoint_name,
    )


def wait_resource_endpoint_status(
    reference: k8s.CustomResourceReference,
    expected_status: str,
    wait_periods: int = 30,
    period_length: int = 30,
):
    return wait_for_status(
        expected_status,
        wait_periods,
        period_length,
        get_endpoint_resource_status,
        reference,
    )


def assert_endpoint_status_in_sync(endpoint_name, reference, expected_status):
    assert (
        wait_sagemaker_endpoint_status(endpoint_name, expected_status)
        == wait_resource_endpoint_status(reference, expected_status, 2)
        == expected_status
    )


def get_model_package_sagemaker_status(model_package_arn):
    response = sagemaker_client().describe_model_package(
        ModelPackageName=model_package_arn
    )
    return response["ModelPackageStatus"]


def get_model_package_resource_status(reference: k8s.CustomResourceReference):
    resource = k8s.get_resource(reference)
    assert "modelPackageStatus" in resource["status"]
    return resource["status"]["modelPackageStatus"]


def wait_sagemaker_model_package_status(
    model_package_arn,
    expected_status: str,
    wait_periods: int = 60,
    period_length: int = 30,
):
    return wait_for_status(
        expected_status,
        wait_periods,
        period_length,
        get_model_package_sagemaker_status,
        model_package_arn,
    )


def wait_resource_model_package_status(
    reference: k8s.CustomResourceReference,
    expected_status: str,
    wait_periods: int = 30,
    period_length: int = 30,
):
    return wait_for_status(
        expected_status,
        wait_periods,
        period_length,
        get_model_package_resource_status,
        reference,
    )


def assert_model_package_status_in_sync(model_package_arn, reference, expected_status):
    assert (
        wait_sagemaker_model_package_status(model_package_arn, expected_status)
        == wait_resource_model_package_status(reference, expected_status, 2)
        == expected_status
    )


def get_sagemaker_model_package_group(model_package_group_name: str):
    try:
        return sagemaker_client().describe_model_package_group(
            ModelPackageGroupName=model_package_group_name
        )
    except botocore.exceptions.ClientError as error:
        logging.error(
            f"SageMaker could not find a model package group with the name {model_package_group_name}. Error {error}"
        )
        return None


def get_sagemaker_model_package(model_package_name: str):
    try:
        model_package = sagemaker_client().describe_model_package(
            ModelPackageName=model_package_name
        )
        return model_package
    except botocore.exceptions.ClientError as error:
        logging.error(
            f"SageMaker could not find a model package with the name {model_package_name}. Error {error}"
        )
        return None


def assert_tags_in_sync(resource_arn, resource_tags):
    response = sagemaker_client().list_tags(ResourceArn=resource_arn)
    response_tags = response["Tags"]

    while "NextToken" in response:
        response = sagemaker_client().list_tags(
            ResourceArn=resource_arn, NextToken=response["NextToken"]
        )
        response_tags.extend(response["Tags"])

    # SageMaker returns tags with Key,Value while ACK returns key,value
    response_tags = [{"key": d["Key"], "value": d["Value"]} for d in response_tags]
    TestCase().assertCountEqual(tags.clean(response_tags), resource_tags)


def get_sagemaker_training_job(training_job_name: str):
    try:
        training_job = sagemaker_client().describe_training_job(
            TrainingJobName=training_job_name
        )
        return training_job
    except botocore.exceptions.ClientError as error:
        logging.error(
            f"SageMaker could not find a training job with the name {training_job_name}. Error {error}"
        )
        return None


def get_training_sagemaker_status(training_job_name: str):
    training_sm_desc = get_sagemaker_training_job(training_job_name)
    return training_sm_desc["TrainingJobStatus"]


def get_training_resource_status(reference: k8s.CustomResourceReference):
    resource = k8s.get_resource(reference)
    assert "trainingJobStatus" in resource["status"]
    return resource["status"]["trainingJobStatus"]


def wait_resource_training_status(
    reference: k8s.CustomResourceReference,
    expected_status: str,
    wait_periods: int = 60,
    period_length: int = 30,
):
    return wait_for_status(
        expected_status,
        wait_periods,
        period_length,
        get_training_resource_status,
        reference,
    )


def wait_sagemaker_training_status(
    training_job_name,
    expected_status: str,
    wait_periods: int = 60,
    period_length: int = 30,
):
    return wait_for_status(
        expected_status,
        wait_periods,
        period_length,
        get_training_sagemaker_status,
        training_job_name,
    )


def assert_training_status_in_sync(training_job_name, reference, expected_status):
    assert (
        wait_sagemaker_training_status(training_job_name, expected_status)
        == wait_resource_training_status(reference, expected_status)
        == expected_status
    )


def get_sagemaker_endpoint(endpoint_name: str):
    try:
        return sagemaker_client().describe_endpoint(EndpointName=endpoint_name)
    except botocore.exceptions.ClientError as error:
        logging.error(
            f"SageMaker could not find a endpoint with the name {endpoint_name}. Error {error}"
        )
        return None


def get_sagemaker_model(model_name: str, sm_client=None):
    sm_client = sm_client or sagemaker_client()
    try:
        return sm_client.describe_model(ModelName=model_name)
    except botocore.exceptions.ClientError as error:
        logging.error(
            f"SageMaker could not find a model with the name {model_name}. Error {error}"
        )
        return None


def get_sagemaker_endpoint_config(config_name: str):
    try:
        return sagemaker_client().describe_endpoint_config(
            EndpointConfigName=config_name
        )
    except botocore.exceptions.ClientError as error:
        logging.error(
            f"SageMaker could not find an endpoint config with the name {config_name}. Error {error}"
        )
        return None


def get_sagemaker_pipeline_execution(pipeline_execution_arn: str):
    try:
        pipeline_execution = sagemaker_client().describe_pipeline_execution(
            PipelineExecutionArn=pipeline_execution_arn
        )
        return pipeline_execution
    except botocore.exceptions.ClientError as error:
        logging.error(
            f"SageMaker could not find a pipeline execution with the arn {pipeline_execution_arn}. Error {error}"
        )
        return None


def get_sagemaker_pipeline(pipeline_name: str):
    try:
        return sagemaker_client().describe_pipeline(PipelineName=pipeline_name)
    except botocore.exceptions.ClientError as error:
        logging.error(
            f"SageMaker could not find a pipeline with the name {pipeline_name}. Error {error}"
        )
        return None


def delete_custom_resource(
    reference,
    wait_period=cfg.JOB_DELETE_WAIT_PERIODS,
    wait_length=cfg.JOB_DELETE_WAIT_LENGTH,
):
    deleted = True

    if k8s.get_resource_exists(reference): 
        _, deleted = k8s.delete_custom_resource(reference, wait_period, wait_length)
    return deleted
