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
"""Integration tests for the SageMaker ModelPackage API.
"""

import pytest

from acktest.resources import random_suffix_name
from acktest.k8s import resource as k8s

from e2e import (
    service_marker,
    CRD_GROUP,
    CRD_VERSION,
    create_adopted_resource,
    wait_sagemaker_model_package_status,
    assert_model_package_status_in_sync,
    delete_custom_resource,
    get_sagemaker_model_package_group,
    get_sagemaker_model_package,
    sagemaker_client,
)
from e2e.replacement_values import REPLACEMENT_VALUES
from e2e.common import config as cfg


@pytest.fixture(scope="module")
def name_suffix():
    return random_suffix_name("sdk-model-package", 32)


def sdk_make_model_package_group(model_package_group_name):
    model_package_group_input = {
        "ModelPackageGroupName": model_package_group_name,
        "ModelPackageGroupDescription": "ModelPackageGroup for adoption",
    }

    model_package_group_response = sagemaker_client().create_model_package_group(
        **model_package_group_input
    )
    assert model_package_group_response.get("ModelPackageGroupArn", None) is not None
    return model_package_group_input, model_package_group_response


def sdk_make_model_package(model_package_group_name):
    data_bucket = REPLACEMENT_VALUES["SAGEMAKER_DATA_BUCKET"]
    model_package_input = {
        "ModelPackageGroupName": model_package_group_name,
        "InferenceSpecification": {
            "Containers": [
                {
                    "Image": REPLACEMENT_VALUES["XGBOOST_IMAGE_URI"],
                    "ModelDataUrl": f"s3://{data_bucket}/sagemaker/model/xgboost-mnist-model.tar.gz",
                }
            ],
            "SupportedContentTypes": [
                "text/csv",
            ],
            "SupportedResponseMIMETypes": [
                "text/csv",
            ],
        },
    }

    model_package_response = sagemaker_client().create_model_package(**model_package_input)
    assert model_package_response.get("ModelPackageArn", None) is not None

    return model_package_input, model_package_response


@pytest.fixture(scope="module")
def sdk_model_package(name_suffix):
    model_package_group_name = name_suffix + "-group"

    (
        model_package_group_input,
        model_package_group_response,
    ) = sdk_make_model_package_group(model_package_group_name)
    model_package_input, model_package_response = sdk_make_model_package(model_package_group_name)

    yield (
        model_package_group_input,
        model_package_group_response,
        model_package_input,
        model_package_response,
    )
    model_package_arn = model_package_response.get("ModelPackageArn")
    if get_sagemaker_model_package(model_package_arn) is not None:
        wait_sagemaker_model_package_status(model_package_arn, cfg.JOB_STATUS_COMPLETED)
        sagemaker_client().delete_model_package(ModelPackageName=model_package_arn)
    if get_sagemaker_model_package_group(model_package_group_name) is not None:
        sagemaker_client().delete_model_package_group(
            ModelPackageGroupName=model_package_group_name
        )


@pytest.fixture(scope="module")
def adopted_model_package(sdk_model_package):
    (
        model_package_group_input,
        _,
        model_package_input,
        model_package_response,
    ) = sdk_model_package

    replacements = REPLACEMENT_VALUES.copy()
    # adopt model package group
    replacements["ADOPTED_RESOURCE_NAME"] = (
        "adopt-" + model_package_group_input["ModelPackageGroupName"]
    )
    replacements["TARGET_RESOURCE_AWS"] = replacements[
        "TARGET_RESOURCE_K8S"
    ] = model_package_group_input["ModelPackageGroupName"]
    replacements["RESOURCE_KIND"] = "ModelPackageGroup"

    (
        adopt_model_package_group_reference,
        _,
        adopt_model_package_group_resource,
    ) = create_adopted_resource(
        replacements=replacements,
    )
    assert adopt_model_package_group_resource is not None

    # adopt model package
    replacements["ADOPTED_RESOURCE_NAME"] = (
        "adopt-" + model_package_input["ModelPackageGroupName"] + "-child"
    )
    replacements["TARGET_RESOURCE_AWS"] = model_package_response.get("ModelPackageArn")
    replacements["TARGET_RESOURCE_K8S"] = model_package_input["ModelPackageGroupName"] + "-child"
    replacements["RESOURCE_KIND"] = "ModelPackage"

    (
        adopt_model_package_reference,
        _,
        adopt_model_package_resource,
    ) = create_adopted_resource(
        replacements=replacements,
        spec_file="adopted_resource_base_arn",
    )
    assert adopt_model_package_resource is not None

    yield (adopt_model_package_group_reference, adopt_model_package_reference)

    for cr in (adopt_model_package_group_reference, adopt_model_package_reference):
        assert delete_custom_resource(cr, cfg.DELETE_WAIT_PERIOD, cfg.DELETE_WAIT_LENGTH)


@service_marker
class TestAdoptedModelPackage:
    def test_smoke(self, sdk_model_package, adopted_model_package):
        (
            adopt_model_package_group_reference,
            adopt_model_package_reference,
        ) = adopted_model_package

        (
            model_package_group_input,
            model_package_group_response,
            model_package_input,
            model_package_response,
        ) = sdk_model_package

        namespace = "default"
        model_package_group_name = k8s.get_resource(adopt_model_package_group_reference)["spec"][
            "aws"
        ]["nameOrID"]
        model_package_arn = k8s.get_resource(adopt_model_package_reference)["spec"]["aws"]["arn"]

        for reference in (
            adopt_model_package_group_reference,
            adopt_model_package_reference,
        ):
            assert k8s.wait_on_condition(reference, k8s.CONDITION_TYPE_ADOPTED, "True")

        model_package_group_reference = k8s.create_reference(
            CRD_GROUP,
            CRD_VERSION,
            cfg.MODEL_PACKAGE_GROUP_RESOURCE_PLURAL,
            model_package_group_name,
            namespace,
        )
        model_package_group_resource = k8s.wait_resource_consumed_by_controller(
            model_package_group_reference
        )
        assert model_package_group_resource is not None

        assert (
            model_package_group_resource["spec"].get("modelPackageGroupName", None)
            == model_package_group_name
        )
        assert (
            model_package_group_resource["spec"].get("modelPackageGroupDescription", None)
            is not None
        )
        assert k8s.get_resource_arn(
            model_package_group_resource
        ) == model_package_group_response.get("ModelPackageGroupArn", None)

        model_package_reference = k8s.create_reference(
            CRD_GROUP,
            CRD_VERSION,
            cfg.MODEL_PACKAGE_RESOURCE_PLURAL,
            model_package_input["ModelPackageGroupName"] + "-child",
            namespace,
        )
        model_package_resource = k8s.wait_resource_consumed_by_controller(model_package_reference)
        assert model_package_resource is not None

        assert (
            model_package_resource["spec"].get("modelPackageGroupName", None)
            == model_package_group_name
        )
        assert model_package_resource["spec"].get("inferenceSpecification", None) is not None
        assert k8s.get_resource_arn(model_package_resource) == model_package_response.get(
            "ModelPackageArn", None
        )

        assert_model_package_status_in_sync(
            model_package_arn,
            model_package_reference,
            cfg.JOB_STATUS_COMPLETED,
        )
        assert k8s.wait_on_condition(
            model_package_reference, k8s.CONDITION_TYPE_RESOURCE_SYNCED, "True"
        )

        for cr in (model_package_reference, model_package_group_reference):
            assert delete_custom_resource(cr, cfg.DELETE_WAIT_PERIOD, cfg.DELETE_WAIT_LENGTH)
