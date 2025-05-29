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
"""Integration tests for the SageMaker Studio.
"""

import botocore
import pytest
import logging
import boto3
from typing import Dict

from acktest.resources import random_suffix_name
from acktest.k8s import resource as k8s
from e2e import (
    service_marker,
    create_sagemaker_resource,
    delete_custom_resource,
    wait_for_status,
    sagemaker_client,
    assert_tags_in_sync,
)
from e2e.replacement_values import REPLACEMENT_VALUES
from e2e.bootstrap_resources import get_bootstrap_resources
from e2e.common import config as cfg


def get_default_vpc():
    return boto3.client("ec2").describe_vpcs(
        Filters=[
            {"Name": "isDefault", "Values": ["true"]},
        ]
    )[
        "Vpcs"
    ][0]["VpcId"]


def get_subnet(vpc_id):
    return boto3.client("ec2").describe_subnets(Filters=[{"Name": "vpcId", "Values": [vpc_id]}])[
        "Subnets"
    ][0]["SubnetId"]


def get_domain_sagemaker_status(domain_id):
    response = boto3.client("sagemaker").describe_domain(DomainId=domain_id)
    return response["Status"]


def get_user_profile_sagemaker_status(domain_id, user_profile_name):
    response = boto3.client("sagemaker").describe_user_profile(
        DomainId=domain_id, UserProfileName=user_profile_name
    )
    return response["Status"]


def get_app_sagemaker_status(domain_id, user_profile_name, app_type, app_name):
    response = boto3.client("sagemaker").describe_app(
        DomainId=domain_id,
        UserProfileName=user_profile_name,
        AppType=app_type,
        AppName=app_name,
    )
    return response["Status"]


def get_k8s_resource_status(reference: k8s.CustomResourceReference):
    resource = k8s.get_resource(reference)
    assert "status" in resource["status"]
    return resource["status"]["status"]


def apply_domain_yaml(resource_name):
    replacements = REPLACEMENT_VALUES.copy()
    replacements["DOMAIN_NAME"] = resource_name
    replacements["VPC_ID"] = get_default_vpc()
    replacements["SUBNET_ID"] = get_subnet(replacements["VPC_ID"])

    reference, spec, resource = create_sagemaker_resource(
        resource_plural="domains",
        resource_name=resource_name,
        spec_file="domain",
        replacements=replacements,
    )
    return reference, resource, spec


def apply_user_profile_yaml(resource_name, domain_id):
    replacements = REPLACEMENT_VALUES.copy()
    replacements["USER_PROFILE_NAME"] = resource_name
    replacements["DOMAIN_ID"] = domain_id
    reference, spec, resource = create_sagemaker_resource(
        resource_plural="userprofiles",
        resource_name=resource_name,
        spec_file="user_profile",
        replacements=replacements,
    )
    return reference, resource, spec


def apply_app_yaml(domain_id, user_profile_name):
    replacements = REPLACEMENT_VALUES.copy()
    replacements["DOMAIN_ID"] = domain_id
    replacements["USER_PROFILE_NAME"] = user_profile_name
    reference, spec, resource = create_sagemaker_resource(
        resource_plural="apps",
        resource_name="default",
        spec_file="app",
        replacements=replacements,
    )
    return reference, resource, spec


def assert_domain_status_in_sync(domain_id, reference, expected_status):
    sm_status = wait_for_status(expected_status, 10, 30, get_domain_sagemaker_status, domain_id)
    k8s_status = wait_for_status(expected_status, 10, 30, get_k8s_resource_status, reference)
    assert sm_status == k8s_status == expected_status


def assert_user_profile_status_in_sync(domain_id, user_profile_name, reference, expected_status):
    sm_status = wait_for_status(
        expected_status,
        10,
        30,
        get_user_profile_sagemaker_status,
        domain_id,
        user_profile_name,
    )
    k8s_status = wait_for_status(expected_status, 10, 30, get_k8s_resource_status, reference)
    assert sm_status == k8s_status == expected_status


def assert_app_status_in_sync(
    domain_id, user_profile_name, app_type, app_name, reference, expected_status
):
    sm_status = wait_for_status(
        expected_status,
        10,
        30,
        get_app_sagemaker_status,
        domain_id,
        user_profile_name,
        app_type,
        app_name,
    )
    k8s_status = wait_for_status(expected_status, 10, 30, get_k8s_resource_status, reference)
    assert sm_status == k8s_status == expected_status


def patch_domain_kernel_instance(reference, spec, instance_type):
    spec["spec"]["defaultUserSettings"]["kernelGatewayAppSettings"]["defaultResourceSpec"][
        "instanceType"
    ] = instance_type
    resource = k8s.patch_custom_resource(reference, spec)
    assert resource is not None
    return resource


def patch_user_profile_kernel_instance(reference, spec, instance_type):
    spec["spec"]["userSettings"]["kernelGatewayAppSettings"]["defaultResourceSpec"][
        "instanceType"
    ] = instance_type
    resource = k8s.patch_custom_resource(reference, spec)
    assert resource is not None
    return resource


def get_domain_kernel_instance(domain_id):
    response = boto3.client("sagemaker").describe_domain(DomainId=domain_id)
    return response["DefaultUserSettings"]["KernelGatewayAppSettings"]["DefaultResourceSpec"][
        "InstanceType"
    ]


def get_user_profile_kernel_instance(domain_id, user_profile_name):
    response = boto3.client("sagemaker").describe_user_profile(
        DomainId=domain_id, UserProfileName=user_profile_name
    )
    return response["UserSettings"]["KernelGatewayAppSettings"]["DefaultResourceSpec"][
        "InstanceType"
    ]


@pytest.fixture(scope="function")
def domain_fixture():
    resource_name = random_suffix_name("sm-domain", 15)
    reference, resource, spec = apply_domain_yaml(resource_name)

    assert resource is not None
    if k8s.get_resource_arn(resource) is None:
        logging.error(f"ARN for this resource is None, resource status is: {resource['status']}")
    assert k8s.get_resource_arn(resource) is not None

    yield (reference, resource, spec)

    assert delete_custom_resource(
        reference, cfg.JOB_DELETE_WAIT_PERIODS, cfg.JOB_DELETE_WAIT_LENGTH
    )


@pytest.fixture(scope="function")
def user_profile_fixture(domain_fixture):
    (domain_reference, domain_resource, domain_spec) = domain_fixture
    assert k8s.get_resource_exists(domain_reference)

    domain_id = domain_resource["status"].get("domainID", None)
    assert domain_id is not None

    assert_domain_status_in_sync(domain_id, domain_reference, "InService")

    domain_resource = patch_domain_kernel_instance(domain_reference, domain_spec, "ml.t3.large")
    wait_for_status("ml.t3.large", 10, 30, get_domain_kernel_instance, domain_id)
    assert_domain_status_in_sync(domain_id, domain_reference, "InService")

    resource_name = random_suffix_name("profile", 15)
    (
        user_profile_reference,
        user_profile_resource,
        user_profile_spec,
    ) = apply_user_profile_yaml(resource_name, domain_id)

    assert user_profile_resource is not None
    if k8s.get_resource_arn(user_profile_resource) is None:
        logging.error(
            f"ARN for this resource is None, resource status is: {user_profile_resource['status']}"
        )
    assert k8s.get_resource_arn(user_profile_resource) is not None

    yield (
        domain_reference,
        domain_resource,
        domain_spec,
        user_profile_reference,
        user_profile_resource,
        user_profile_spec,
    )

    assert delete_custom_resource(
        user_profile_reference,
        cfg.JOB_DELETE_WAIT_PERIODS,
        cfg.JOB_DELETE_WAIT_LENGTH,
    )


@pytest.fixture(scope="function")
def app_fixture(user_profile_fixture):
    (
        domain_reference,
        domain_resource,
        domain_spec,
        user_profile_reference,
        user_profile_resource,
        user_profile_spec,
    ) = user_profile_fixture
    assert k8s.get_resource_exists(domain_reference)
    assert k8s.get_resource_exists(user_profile_reference)

    domain_id = domain_resource["status"].get("domainID", None)
    user_profile_name = user_profile_resource["spec"]["userProfileName"]
    assert_user_profile_status_in_sync(
        domain_id, user_profile_name, user_profile_reference, "InService"
    )

    user_profile_resource = patch_user_profile_kernel_instance(
        user_profile_reference, user_profile_spec, "ml.t3.large"
    )
    wait_for_status(
        "ml.t3.large",
        10,
        30,
        get_user_profile_kernel_instance,
        domain_id,
        user_profile_name,
    )
    assert_user_profile_status_in_sync(
        domain_id, user_profile_name, user_profile_reference, "InService"
    )

    (app_reference, app_resource, app_spec) = apply_app_yaml(domain_id, user_profile_name)

    assert app_resource is not None
    if k8s.get_resource_arn(app_resource) is None:
        logging.error(
            f"ARN for this resource is None, resource status is: {app_resource['status']}"
        )
    assert k8s.get_resource_arn(app_resource) is not None

    yield (
        domain_reference,
        domain_resource,
        domain_spec,
        user_profile_reference,
        user_profile_resource,
        user_profile_spec,
        app_reference,
        app_resource,
        app_spec,
    )

    assert delete_custom_resource(
        app_reference, cfg.JOB_DELETE_WAIT_PERIODS, cfg.JOB_DELETE_WAIT_LENGTH
    )


class TestDomain:
    def test_studio(self, app_fixture):
        (
            domain_reference,
            domain_resource,
            domain_spec,
            user_profile_reference,
            user_profile_resource,
            user_profile_spec,
            app_reference,
            app_resource,
            app_spec,
        ) = app_fixture

        assert k8s.get_resource_exists(domain_reference)
        assert k8s.get_resource_exists(user_profile_reference)
        assert k8s.get_resource_exists(app_reference)

        domain_id = domain_resource["status"].get("domainID", None)
        user_profile_name = user_profile_resource["spec"]["userProfileName"]
        app_type = app_resource["spec"]["appType"]
        app_name = app_resource["spec"]["appName"]

        assert_app_status_in_sync(
            domain_id,
            user_profile_name,
            app_type,
            app_name,
            app_reference,
            "InService",
        )
