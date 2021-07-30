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
"""Integration tests for the SageMaker modelPackage API.
"""

import botocore
import pytest
import logging
from typing import Dict

from acktest.resources import random_suffix_name
from acktest.k8s import resource as k8s
from e2e import (
    service_marker,
    wait_for_status,
    create_sagemaker_resource,
    sagemaker_client,
    assert_tags_in_sync,
)
from e2e.replacement_values import REPLACEMENT_VALUES
from e2e.bootstrap_resources import get_bootstrap_resources
from e2e.common import config as cfg

RESOURCE_PLURAL = "modelpackages"

DELETE_WAIT_PERIOD = 20
DELETE_WAIT_LENGTH = 30

# Disabled since versioned model package is not supported
# Issue with multiple creates.

# @pytest.fixture(scope="function")
# def xgboost_model_package_group():
#     resource_name = random_suffix_name("xgboost-model-package-group", 32)

#     replacements = REPLACEMENT_VALUES.copy()
#     replacements["MODEL_PACKAGE_GROUP_NAME"] = resource_name

#     (
#         model_package_group_reference,
#         model_package_group_spec,
#         model_package_group_resource,
#     ) = create_sagemaker_resource(
#         resource_plural=cfg.MODEL_PACKAGE_GROUP_RESOURCE_PLURAL,
#         resource_name=resource_name,
#         spec_file="xgboost_model_package_group",
#         replacements=replacements,
#     )
#     assert model_package_group_resource is not None
#     if k8s.get_resource_arn(model_package_group_resource) is None:
#         logging.error(
#             f"ARN for this resource is None, resource status is: {model_package_group_resource['status']}"
#         )
#     assert k8s.get_resource_arn(model_package_group_resource) is not None

#     yield (model_package_group_reference, model_package_group_resource)

#     # Delete the k8s resource if not already deleted by tests
#     if k8s.get_resource_exists(model_package_group_reference):
#         _, deleted = k8s.delete_custom_resource(model_package_group_reference, 3, 10)
#         assert deleted


# @pytest.fixture(scope="function")
# def xgboost_versioned_model_package(xgboost_model_package_group):
#     resource_name = random_suffix_name("xgboost-versioned-model-package", 38)
#     (_, model_package_group_resource) = xgboost_model_package_group
#     model_package_group_resource_name = model_package_group_resource["spec"].get(
#         "modelPackageGroupName", None
#     )

#     replacements = REPLACEMENT_VALUES.copy()
#     replacements["MODEL_PACKAGE_GROUP_NAME"] = model_package_group_resource_name
#     replacements["MODEL_PACKAGE_RESOURCE_NAME"] = resource_name
#     reference, spec, resource = create_sagemaker_resource(
#         resource_plural=RESOURCE_PLURAL,
#         resource_name=resource_name,
#         spec_file="xgboost_versioned_model_package",
#         replacements=replacements,
#     )
#     assert resource is not None

#     yield (reference, spec, resource)
#     # Delete the k8s resource if not already deleted by tests
#     if k8s.get_resource_exists(reference):
#         _, deleted = k8s.delete_custom_resource(reference, 3, 10)
#         assert deleted


@pytest.fixture(scope="function")
def xgboost_unversioned_model_package():
    resource_name = random_suffix_name("xgboost-unversioned-model-package", 38)
    replacements = REPLACEMENT_VALUES.copy()
    replacements["MODEL_PACKAGE_NAME"] = resource_name
    reference, _, resource = create_sagemaker_resource(
        resource_plural=RESOURCE_PLURAL,
        resource_name=resource_name,
        spec_file="xgboost_unversioned_model_package",
        replacements=replacements,
    )

    assert resource is not None

    yield (reference, resource)

    if k8s.get_resource_exists(reference):
        _, deleted = k8s.delete_custom_resource(
            reference, DELETE_WAIT_PERIOD, DELETE_WAIT_LENGTH
        )
        assert deleted


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


def get_model_package_sagemaker_status(model_package_name: str):
    model_package_sm_desc = get_sagemaker_model_package(model_package_name)
    return model_package_sm_desc["ModelPackageStatus"]


def get_model_package_resource_status(reference: k8s.CustomResourceReference):
    resource = k8s.get_resource(reference)
    assert "modelPackageStatus" in resource["status"]
    return resource["status"]["modelPackageStatus"]


@pytest.mark.canary
@service_marker
class TestmodelPackage:
    def _wait_resource_model_package_status(
        self,
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

    def _wait_sagemaker_model_package_status(
        self,
        model_package_name,
        expected_status: str,
        wait_periods: int = 30,
        period_length: int = 30,
    ):
        return wait_for_status(
            expected_status,
            wait_periods,
            period_length,
            get_model_package_sagemaker_status,
            model_package_name,
        )

    def _assert_model_package_status_in_sync(
        self, model_package_name, reference, expected_status
    ):
        assert (
            self._wait_sagemaker_model_package_status(
                model_package_name, expected_status
            )
            == self._wait_resource_model_package_status(reference, expected_status)
            == expected_status
        )

    def test_unversioned_model_package_completed(
        self, xgboost_unversioned_model_package
    ):
        (reference, resource) = xgboost_unversioned_model_package
        assert k8s.get_resource_exists(reference)

        model_package_name = resource["spec"].get("modelPackageName", None)
        assert model_package_name is not None

        model_package_desc = get_sagemaker_model_package(model_package_name)
        model_package_arn = model_package_desc["ModelPackageArn"]

        if k8s.get_resource_arn(resource) is None:
            logging.error(
                f"ARN for this resource is None, resource status is: {resource['status']}"
            )

        assert k8s.get_resource_arn(resource) == model_package_arn

        self._assert_model_package_status_in_sync(
            model_package_name, reference, cfg.JOB_STATUS_INPROGRESS
        )
        assert k8s.wait_on_condition(reference, "ACK.ResourceSynced", "False")

        self._assert_model_package_status_in_sync(
            model_package_name, reference, cfg.JOB_STATUS_COMPLETED
        )
        assert k8s.wait_on_condition(reference, "ACK.ResourceSynced", "True")

        resource_tags = resource["spec"].get("tags", None)
        assert_tags_in_sync(model_package_arn, resource_tags)

        # Check that you can delete a completed resource from k8s
        _, deleted = k8s.delete_custom_resource(
            reference, DELETE_WAIT_PERIOD, DELETE_WAIT_LENGTH
        )
        assert deleted is True
        assert get_sagemaker_model_package(model_package_name) is None

    # def test_versioned_model_package_completed(self, xgboost_versioned_model_package):
    #     (reference, spec, resource) = xgboost_versioned_model_package
    #     assert k8s.get_resource_exists(reference)

    #     model_package_group_name = resource["spec"].get("modelPackageGroupName")
    #     # Model package name for Versioned Model packages is the ARN of the resource
    #     model_package_name = sagemaker_client().list_model_packages(
    #         ModelPackageGroupName=model_package_group_name
    #     )["ModelPackageSummaryList"][0]["ModelPackageArn"]

    #     model_package_desc = get_sagemaker_model_package(model_package_name)
    #     if k8s.get_resource_arn(resource) is None:
    #         logging.error(
    #             f"ARN for this resource is None, resource status is: {resource['status']}"
    #         )

    #     assert k8s.get_resource_arn(resource) == model_package_name
    #     assert model_package_desc["ModelPackageStatus"] == cfg.JOB_STATUS_INPROGRESS
    #     self._assert_model_package_status_in_sync(
    #         model_package_name, reference, cfg.JOB_STATUS_INPROGRESS
    #     )
    #     assert k8s.wait_on_condition(reference, "ACK.ResourceSynced", "False")

    #     self._assert_model_package_status_in_sync(
    #         model_package_name, reference, cfg.JOB_STATUS_COMPLETED
    #     )
    #     assert k8s.wait_on_condition(reference, "ACK.ResourceSynced", "True")

    #     # Update the resource
    #     new_model_approval_status = "Approved"
    #     approval_description = "Approved modelpackage"
    #     spec["spec"]["modelApprovalStatus"] = new_model_approval_status
    #     spec["spec"]["approvalDescription"] = approval_description
    #     resource = k8s.patch_custom_resource(reference, spec)
    #     resource = k8s.wait_resource_consumed_by_controller(reference)
    #     assert resource is not None

    #     self._assert_model_package_status_in_sync(
    #         model_package_name, reference, cfg.JOB_STATUS_COMPLETED
    #     )
    #     assert k8s.wait_on_condition(reference, "ACK.ResourceSynced", "True")

    #     model_package_desc = get_sagemaker_model_package(model_package_name)
    #     assert model_package_desc["ModelApprovalStatus"] == new_model_approval_status
    #     assert model_package_desc["ApprovalDescription"] == approval_description

    #     assert (
    #         resource["spec"].get("modelApprovalStatus", None)
    #         == new_model_approval_status
    #     )
    #     assert resource["spec"].get("approvalDescription", None) == approval_description
    #     # Check that you can delete a completed resource from k8s
    #     _, deleted = k8s.delete_custom_resource(reference, 3, 10)
    #     assert deleted is True
    #     assert get_sagemaker_model_package(model_package_name) is None
