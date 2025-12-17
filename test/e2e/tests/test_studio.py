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
"""Integration tests for the SageMaker Studio."""

import logging

import boto3
import pytest
from acktest.k8s import resource as k8s
from acktest.resources import random_suffix_name

from e2e import (
    create_sagemaker_resource,
    delete_custom_resource,
    service_marker,
    wait_for_status,
)
from e2e.replacement_values import REPLACEMENT_VALUES

STUDIO_WAIT_PERIOD = 120
STUDIO_WAIT_LENGTH = 30
STUDIO_STATUS_INSERVICE = "InService"


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


def get_space_sagemaker_status(domain_id, space_name):
    response = boto3.client("sagemaker").describe_space(DomainId=domain_id, SpaceName=space_name)
    return response["Status"]


def get_app_sagemaker_status(domain_id, app_association, resource_name, app_type, app_name):
    app_association = app_association.lower()
    assert app_association in ["user_profile", "space"]

    response = {}
    if app_association == "user_profile":
        response = boto3.client("sagemaker").describe_app(
            DomainId=domain_id,
            UserProfileName=resource_name,
            AppType=app_type,
            AppName=app_name,
        )
    elif app_association == "space":
        response = boto3.client("sagemaker").describe_app(
            DomainId=domain_id,
            SpaceName=resource_name,
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


def apply_user_profile_yaml(domain_id, resource_name):
    replacements = REPLACEMENT_VALUES.copy()
    replacements["DOMAIN_ID"] = domain_id
    replacements["USER_PROFILE_NAME"] = resource_name
    reference, spec, resource = create_sagemaker_resource(
        resource_plural="userprofiles",
        resource_name=resource_name,
        spec_file="user_profile",
        replacements=replacements,
    )
    return reference, resource, spec


def apply_space_yaml(domain_id, user_profile_name, resource_name, share_type):
    share_type = share_type.lower()
    assert share_type in ["shared", "private"]

    replacements = REPLACEMENT_VALUES.copy()
    replacements["DOMAIN_ID"] = domain_id
    replacements["SPACE_NAME"] = resource_name
    replacements["USER_PROFILE_NAME"] = user_profile_name
    reference, spec, resource = create_sagemaker_resource(
        resource_plural="spaces",
        resource_name=resource_name,
        spec_file=f"{share_type}_space",
        replacements=replacements,
    )
    return reference, resource, spec


def apply_app_yaml(domain_id, app_association, app_association_resource_name, resource_name):
    app_association = app_association.lower()
    assert app_association in ["user_profile", "space"]

    replacements = REPLACEMENT_VALUES.copy()
    replacements["DOMAIN_ID"] = domain_id
    replacements["APP_NAME"] = resource_name
    if app_association == "user_profile":
        replacements["USER_PROFILE_NAME"] = app_association_resource_name

    elif app_association == "space":
        replacements["SPACE_NAME"] = app_association_resource_name

    reference, spec, resource = create_sagemaker_resource(
        resource_plural="apps",
        resource_name=resource_name,
        spec_file=f"app_{app_association}",
        replacements=replacements,
    )
    return reference, resource, spec


def assert_domain_status_in_sync(domain_id, reference, expected_status):
    sm_status = wait_for_status(
        expected_status,
        STUDIO_WAIT_PERIOD,
        STUDIO_WAIT_LENGTH,
        get_domain_sagemaker_status,
        domain_id,
    )
    k8s_status = wait_for_status(
        expected_status, STUDIO_WAIT_PERIOD, STUDIO_WAIT_LENGTH, get_k8s_resource_status, reference
    )
    assert sm_status == k8s_status == expected_status


def assert_user_profile_status_in_sync(domain_id, user_profile_name, reference, expected_status):
    sm_status = wait_for_status(
        expected_status,
        STUDIO_WAIT_PERIOD,
        STUDIO_WAIT_LENGTH,
        get_user_profile_sagemaker_status,
        domain_id,
        user_profile_name,
    )
    k8s_status = wait_for_status(
        expected_status, STUDIO_WAIT_PERIOD, STUDIO_WAIT_LENGTH, get_k8s_resource_status, reference
    )
    assert sm_status == k8s_status == expected_status


def assert_space_status_in_sync(domain_id, space_name, reference, expected_status):
    sm_status = wait_for_status(
        expected_status,
        STUDIO_WAIT_PERIOD,
        STUDIO_WAIT_LENGTH,
        get_space_sagemaker_status,
        domain_id,
        space_name,
    )
    k8s_status = wait_for_status(
        expected_status, STUDIO_WAIT_PERIOD, STUDIO_WAIT_LENGTH, get_k8s_resource_status, reference
    )
    assert sm_status == k8s_status == expected_status


def assert_app_status_in_sync(
    domain_id, app_association, resource_name, app_type, app_name, reference, expected_status
):
    sm_status = wait_for_status(
        expected_status,
        STUDIO_WAIT_PERIOD,
        STUDIO_WAIT_LENGTH,
        get_app_sagemaker_status,
        domain_id,
        app_association,
        resource_name,
        app_type,
        app_name,
    )
    k8s_status = wait_for_status(
        expected_status, STUDIO_WAIT_PERIOD, STUDIO_WAIT_LENGTH, get_k8s_resource_status, reference
    )
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


def patch_space_jupyter_lab_instance(reference, spec, instance_type):
    spec["spec"]["spaceSettings"]["jupyterLabAppSettings"]["defaultResourceSpec"][
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


def get_space_jupyter_lab_instance(domain_id, space_name):
    response = boto3.client("sagemaker").describe_space(DomainId=domain_id, SpaceName=space_name)
    return response["SpaceSettings"]["JupyterLabAppSettings"]["DefaultResourceSpec"]["InstanceType"]


@pytest.fixture(scope="module")
def domain_fixture():
    resource_name = random_suffix_name("sm-domain", 20)
    reference, resource, spec = apply_domain_yaml(resource_name)

    assert resource is not None
    if k8s.get_resource_arn(resource) is None:
        logging.error(f"ARN for this resource is None, resource status is: {resource['status']}")
    assert k8s.get_resource_arn(resource) is not None

    yield (reference, resource, spec)

    assert delete_custom_resource(
        reference,
        STUDIO_WAIT_PERIOD,
        STUDIO_WAIT_PERIOD,
    )


@pytest.fixture(scope="module")
def user_profile_fixture(domain_fixture):
    (domain_reference, domain_resource, domain_spec) = domain_fixture
    assert k8s.get_resource_exists(domain_reference)

    domain_id = domain_resource["status"].get("domainID", None)
    assert domain_id is not None
    assert_domain_status_in_sync(domain_id, domain_reference, STUDIO_STATUS_INSERVICE)

    domain_resource = patch_domain_kernel_instance(domain_reference, domain_spec, "ml.t3.large")
    wait_for_status(
        "ml.t3.large", STUDIO_WAIT_PERIOD, STUDIO_WAIT_LENGTH, get_domain_kernel_instance, domain_id
    )
    assert_domain_status_in_sync(domain_id, domain_reference, STUDIO_STATUS_INSERVICE)

    resource_name = random_suffix_name("profile", 20)
    (
        user_profile_reference,
        user_profile_resource,
        user_profile_spec,
    ) = apply_user_profile_yaml(domain_id, resource_name)

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
        STUDIO_WAIT_PERIOD,
        STUDIO_WAIT_LENGTH,
    )


@pytest.fixture(scope="module")
def private_space_fixture(user_profile_fixture):
    (
        domain_reference,
        domain_resource,
        domain_spec,
        user_profile_reference,
        user_profile_resource,
        _,
    ) = user_profile_fixture
    assert k8s.get_resource_exists(domain_reference)
    assert k8s.get_resource_exists(user_profile_reference)

    domain_id = domain_resource["status"].get("domainID", None)
    user_profile_name = user_profile_resource["spec"]["userProfileName"]
    assert_user_profile_status_in_sync(
        domain_id, user_profile_name, user_profile_reference, STUDIO_STATUS_INSERVICE
    )

    resource_name = random_suffix_name("private-space", 20)
    (
        space_reference,
        space_resource,
        space_spec,
    ) = apply_space_yaml(domain_id, user_profile_name, resource_name, share_type="private")

    assert space_resource is not None
    if k8s.get_resource_arn(space_resource) is None:
        logging.error(
            f"ARN for this resource is None, resource status is: {space_resource['status']}"
        )
    assert k8s.get_resource_arn(space_resource) is not None

    space_name = space_resource["spec"]["spaceName"]
    space_resource = patch_space_jupyter_lab_instance(space_reference, space_spec, "ml.t3.large")
    wait_for_status(
        "ml.t3.large",
        STUDIO_WAIT_PERIOD,
        STUDIO_WAIT_LENGTH,
        get_space_jupyter_lab_instance,
        domain_id,
        space_name,
    )
    assert_space_status_in_sync(domain_id, space_name, space_reference, STUDIO_STATUS_INSERVICE)

    yield (
        domain_reference,
        domain_resource,
        domain_spec,
        space_reference,
        space_resource,
        space_spec,
    )

    assert delete_custom_resource(
        space_reference,
        STUDIO_WAIT_PERIOD,
        STUDIO_WAIT_LENGTH,
    )


@pytest.fixture(scope="module")
def shared_space_fixture(user_profile_fixture):
    (
        domain_reference,
        domain_resource,
        domain_spec,
        user_profile_reference,
        user_profile_resource,
        _,
    ) = user_profile_fixture
    assert k8s.get_resource_exists(domain_reference)
    assert k8s.get_resource_exists(user_profile_reference)

    domain_id = domain_resource["status"].get("domainID", None)
    user_profile_name = user_profile_resource["spec"]["userProfileName"]
    assert_user_profile_status_in_sync(
        domain_id, user_profile_name, user_profile_reference, STUDIO_STATUS_INSERVICE
    )

    resource_name = random_suffix_name("shared-space", 20)
    (
        space_reference,
        space_resource,
        space_spec,
    ) = apply_space_yaml(domain_id, user_profile_name, resource_name, share_type="shared")

    assert space_resource is not None
    if k8s.get_resource_arn(space_resource) is None:
        logging.error(
            f"ARN for this resource is None, resource status is: {space_resource['status']}"
        )
    assert k8s.get_resource_arn(space_resource) is not None

    yield (
        domain_reference,
        domain_resource,
        domain_spec,
        space_reference,
        space_resource,
        space_spec,
    )

    assert delete_custom_resource(
        space_reference,
        STUDIO_WAIT_PERIOD,
        STUDIO_WAIT_LENGTH,
    )


@pytest.fixture(scope="module")
def app_user_profile_fixture(user_profile_fixture):
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
        domain_id, user_profile_name, user_profile_reference, STUDIO_STATUS_INSERVICE
    )

    user_profile_resource = patch_user_profile_kernel_instance(
        user_profile_reference, user_profile_spec, "ml.t3.large"
    )
    wait_for_status(
        "ml.t3.large",
        STUDIO_WAIT_PERIOD,
        STUDIO_WAIT_LENGTH,
        get_user_profile_kernel_instance,
        domain_id,
        user_profile_name,
    )
    assert_user_profile_status_in_sync(
        domain_id, user_profile_name, user_profile_reference, STUDIO_STATUS_INSERVICE
    )

    app_association = "user_profile"
    resource_name = random_suffix_name("app-user-profile", 20)
    (app_reference, app_resource, app_spec) = apply_app_yaml(
        domain_id, app_association, user_profile_name, resource_name
    )

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
        app_reference,
        STUDIO_WAIT_PERIOD,
        STUDIO_WAIT_LENGTH,
    )


@pytest.fixture(scope="module")
def app_space_fixture(shared_space_fixture):
    (
        domain_reference,
        domain_resource,
        domain_spec,
        space_reference,
        space_resource,
        space_spec,
    ) = shared_space_fixture
    assert k8s.get_resource_exists(domain_reference)
    assert k8s.get_resource_exists(space_reference)

    domain_id = domain_resource["status"].get("domainID", None)
    space_name = space_resource["spec"]["spaceName"]
    assert_space_status_in_sync(domain_id, space_name, space_reference, STUDIO_STATUS_INSERVICE)

    app_association = "space"
    resource_name = random_suffix_name("app-space", 20)
    (app_reference, app_resource, app_spec) = apply_app_yaml(
        domain_id, app_association, space_name, resource_name
    )

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
        space_reference,
        space_resource,
        space_spec,
        app_reference,
        app_resource,
        app_spec,
    )

    assert delete_custom_resource(
        app_reference,
        STUDIO_WAIT_PERIOD,
        STUDIO_WAIT_LENGTH,
    )


@service_marker
@pytest.skip("temp")
class TestDomain:
    def create_private_space(self, private_space_fixture):
        (
            domain_reference,
            domain_resource,
            _,
            space_reference,
            space_resource,
            _,
        ) = private_space_fixture

        assert k8s.get_resource_exists(domain_reference)
        assert k8s.get_resource_exists(space_reference)

        domain_id = domain_resource["status"].get("domainID", None)
        space_name = space_resource["spec"]["spaceName"]

        assert_space_status_in_sync(
            domain_id,
            space_name,
            space_reference,
            STUDIO_STATUS_INSERVICE,
        )

    def create_app_user_profile(self, app_user_profile_fixture):
        (
            domain_reference,
            domain_resource,
            _,
            user_profile_reference,
            user_profile_resource,
            _,
            app_reference,
            app_resource,
            _,
        ) = app_user_profile_fixture

        assert k8s.get_resource_exists(domain_reference)
        assert k8s.get_resource_exists(user_profile_reference)
        assert k8s.get_resource_exists(app_reference)

        domain_id = domain_resource["status"].get("domainID", None)
        user_profile_name = user_profile_resource["spec"]["userProfileName"]
        app_type = app_resource["spec"]["appType"]
        app_name = app_resource["spec"]["appName"]
        app_association = "user_profile"

        assert_app_status_in_sync(
            domain_id,
            app_association,
            user_profile_name,
            app_type,
            app_name,
            app_reference,
            STUDIO_STATUS_INSERVICE,
        )

    def create_app_space(self, app_space_fixture):
        (
            domain_reference,
            domain_resource,
            _,
            space_reference,
            space_resource,
            _,
            app_reference,
            app_resource,
            _,
        ) = app_space_fixture

        assert k8s.get_resource_exists(domain_reference)
        assert k8s.get_resource_exists(space_reference)
        assert k8s.get_resource_exists(app_reference)

        domain_id = domain_resource["status"].get("domainID", None)
        space_name = space_resource["spec"]["spaceName"]
        app_type = app_resource["spec"]["appType"]
        app_name = app_resource["spec"]["appName"]
        app_association = "space"

        assert_app_status_in_sync(
            domain_id,
            app_association,
            space_name,
            app_type,
            app_name,
            app_reference,
            STUDIO_STATUS_INSERVICE,
        )

    def test_studio(
        self, private_space_fixture, app_user_profile_fixture, app_space_fixture
    ):
        self.create_private_space(private_space_fixture)
        self.create_app_user_profile(app_user_profile_fixture)
        self.create_app_space(app_space_fixture)
