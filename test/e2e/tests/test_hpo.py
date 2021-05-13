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
"""Integration tests for the SageMaker HyperParameterTuning API.
"""

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
    sagemaker_client,
)
from e2e.replacement_values import REPLACEMENT_VALUES
from e2e.bootstrap_resources import get_bootstrap_resources
from e2e.common import config as cfg
from time import sleep

RESOURCE_PLURAL = "hyperparametertuningjobs"


@pytest.fixture(scope="function")
def xgboost_hpojob():
    resource_name = random_suffix_name("xgboost-hpojob", 32)
    replacements = REPLACEMENT_VALUES.copy()
    replacements["HPO_JOB_NAME"] = resource_name
    reference, _, resource = create_sagemaker_resource(
        resource_plural=RESOURCE_PLURAL,
        resource_name=resource_name,
        spec_file="xgboost_hpojob",
        replacements=replacements,
    )

    assert resource is not None
    assert k8s.get_resource_arn(resource) is not None

    yield (reference, resource)

    if k8s.get_resource_exists(reference):
        _, deleted = k8s.delete_custom_resource(reference, 3, 10)
        assert deleted


def get_sagemaker_hpo_job(hpo_job_name: str):
    try:
        hpo_desc = sagemaker_client().describe_hyper_parameter_tuning_job(
            HyperParameterTuningJobName=hpo_job_name
        )
        return hpo_desc
    except botocore.exceptions.ClientError as error:
        logging.error(
            f"SageMaker could not find an hpo job with the name {hpo_job_name}. Error {error}"
        )
        return None


def get_hpo_sagemaker_status(hpo_job_name: str):
    hpo_sm_desc = get_sagemaker_hpo_job(hpo_job_name)
    return hpo_sm_desc["HyperParameterTuningJobStatus"]


def get_hpo_resource_status(reference: k8s.CustomResourceReference):
    resource = k8s.get_resource(reference)
    assert "hyperParameterTuningJobStatus" in resource["status"]
    return resource["status"]["hyperParameterTuningJobStatus"]


@service_marker
@pytest.mark.canary
class TestHPO:
    def _wait_resource_hpo_status(
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
            get_hpo_resource_status,
            reference,
        )

    def _wait_sagemaker_hpo_status(
        self,
        hpo_job_name,
        expected_status: str,
        wait_periods: int = 30,
        period_length: int = 30,
    ):
        return wait_for_status(
            expected_status,
            wait_periods,
            period_length,
            get_hpo_sagemaker_status,
            hpo_job_name,
        )

    def _assert_hpo_status_in_sync(self, hpo_job_name, reference, expected_status):
        assert (
            self._wait_sagemaker_hpo_status(hpo_job_name, expected_status)
            == self._wait_resource_hpo_status(reference, expected_status)
            == expected_status
        )

    def test_stopped(self, xgboost_hpojob):
        (reference, resource) = xgboost_hpojob
        assert k8s.get_resource_exists(reference)

        hpo_job_name = resource["spec"].get("hyperParameterTuningJobName", None)
        assert hpo_job_name is not None

        hpo_sm_desc = get_sagemaker_hpo_job(hpo_job_name)
        assert (
            k8s.get_resource_arn(resource) == hpo_sm_desc["HyperParameterTuningJobArn"]
        )
        assert hpo_sm_desc["HyperParameterTuningJobStatus"] == cfg.JOB_STATUS_INPROGRESS
        assert k8s.wait_on_condition(reference, "ACK.ResourceSynced", "False")

        self._assert_hpo_status_in_sync(
            hpo_job_name, reference, cfg.JOB_STATUS_INPROGRESS
        )

        # Delete the k8s resource.
        _, deleted = k8s.delete_custom_resource(reference, 3, 10)
        assert deleted is True

        hpo_sm_desc = get_sagemaker_hpo_job(hpo_job_name)
        assert (
            hpo_sm_desc["HyperParameterTuningJobStatus"] in cfg.LIST_JOB_STATUS_STOPPED
        )

    def test_completed(self, xgboost_hpojob):
        (reference, resource) = xgboost_hpojob
        assert k8s.get_resource_exists(reference)

        hpo_job_name = resource["spec"].get("hyperParameterTuningJobName", None)
        assert hpo_job_name is not None

        hpo_sm_desc = get_sagemaker_hpo_job(hpo_job_name)
        assert (
            k8s.get_resource_arn(resource) == hpo_sm_desc["HyperParameterTuningJobArn"]
        )
        assert hpo_sm_desc["HyperParameterTuningJobStatus"] == cfg.JOB_STATUS_INPROGRESS
        assert k8s.wait_on_condition(reference, "ACK.ResourceSynced", "False")

        self._assert_hpo_status_in_sync(
            hpo_job_name, reference, cfg.JOB_STATUS_COMPLETED
        )
        assert k8s.wait_on_condition(reference, "ACK.ResourceSynced", "True")

        # Check that you can delete a completed resource from k8s
        _, deleted = k8s.delete_custom_resource(reference, 3, 10)
        assert deleted is True
