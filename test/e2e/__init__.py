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

import pytest
import logging
import time
import boto3
from pathlib import Path
import random

from acktest.k8s import resource as k8s

SERVICE_NAME = "sagemaker"
CRD_GROUP = "sagemaker.services.k8s.aws"
ADOPTED_RESOURCE_CRD_GROUP = "services.k8s.aws"
CRD_VERSION = "v1alpha1"

# PyTest marker for the current service
service_marker = pytest.mark.service(arg=SERVICE_NAME)

bootstrap_directory = Path(__file__).parent
resource_directory = Path(__file__).parent / "resources"


def sagemaker_client():
    return boto3.client("sagemaker")


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


def create_adopted_resource(replacements, namespace="default"):
    """
    Wrapper around k8s.load_and_create_resource to create a Adopoted resource
    """

    reference, spec, resource = k8s.load_and_create_resource(
        resource_directory,
        ADOPTED_RESOURCE_CRD_GROUP,
        CRD_VERSION,
        "adoptedresources",
        replacements["ADOPTED_RESOURCE_NAME"],
        "adopted_resource_base",
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
