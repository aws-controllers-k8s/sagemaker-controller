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
"""Integration tests for the SageMaker MonitoringSchedule API.
"""

import botocore
import time
import pytest
import logging

from e2e import (
    service_marker,
    create_sagemaker_resource,
    delete_custom_resource,
    wait_for_status,
    assert_tags_in_sync,
)
from e2e.replacement_values import REPLACEMENT_VALUES
from e2e.common import config as cfg
from e2e.common.fixtures import (
    xgboost_churn_data_quality_job_definition,
    xgboost_churn_endpoint,
)
from acktest.k8s import resource as k8s
from acktest.resources import random_suffix_name

RESOURCE_PLURAL = "monitoringschedules"

# Access variable so it is loaded as a fixture
_accessed = xgboost_churn_data_quality_job_definition, xgboost_churn_endpoint


@pytest.fixture(scope="module")
def xgboost_churn_data_quality_monitoring_schedule(
    xgboost_churn_data_quality_job_definition,
):
    (_, job_definition_resource) = xgboost_churn_data_quality_job_definition

    job_definition_name = job_definition_resource["spec"].get("jobDefinitionName")

    monitoring_schedule_name = random_suffix_name("monitoring-schedule", 32)

    replacements = REPLACEMENT_VALUES.copy()
    replacements["MONITORING_SCHEDULE_NAME"] = monitoring_schedule_name
    replacements["JOB_DEFINITION_NAME"] = job_definition_name
    replacements["MONITORING_TYPE"] = "DataQuality"

    reference, spec, resource = create_sagemaker_resource(
        resource_plural=RESOURCE_PLURAL,
        resource_name=monitoring_schedule_name,
        spec_file="monitoring_schedule_base",
        replacements=replacements,
    )
    assert resource is not None

    yield (reference, resource, spec)

    assert delete_custom_resource(
        reference, cfg.DELETE_WAIT_PERIOD, cfg.DELETE_WAIT_LENGTH
    )


def get_sagemaker_monitoring_schedule(sagemaker_client, monitoring_schedule_name):
    try:
        response = sagemaker_client.describe_monitoring_schedule(
            MonitoringScheduleName=monitoring_schedule_name
        )
        return response
    except botocore.exceptions.ClientError as error:
        logging.error(
            f"Could not find Monitoring Schedule with name {monitoring_schedule_name}. Error {error}"
        )
        return None


def get_monitoring_schedule_sagemaker_status(
    sagemaker_client, monitoring_schedule_name
):
    return sagemaker_client.describe_monitoring_schedule(
        MonitoringScheduleName=monitoring_schedule_name
    )["MonitoringScheduleStatus"]


def get_monitoring_schedule_resource_status(reference: k8s.CustomResourceReference):
    resource = k8s.get_resource(reference)
    assert "monitoringScheduleStatus" in resource["status"]
    return resource["status"]["monitoringScheduleStatus"]


def wait_sagemaker_monitoring_schedule_status(
    sagemaker_client,
    monitoring_schedule_name,
    expected_status: str,
    wait_periods: int = 6,
    period_length: int = 30,
):
    return wait_for_status(
        expected_status,
        wait_periods,
        period_length,
        get_monitoring_schedule_sagemaker_status,
        sagemaker_client,
        monitoring_schedule_name,
    )


def wait_resource_monitoring_schedule_status(
    reference: k8s.CustomResourceReference,
    expected_status: str,
    wait_periods: int = 6,
    period_length: int = 30,
):
    return wait_for_status(
        expected_status,
        wait_periods,
        period_length,
        get_monitoring_schedule_resource_status,
        reference,
    )


@service_marker
@pytest.mark.canary
class TestMonitoringSchedule:
    STATUS_PENDING: str = "Pending"
    STATUS_SCHEDULED: str = "Scheduled"

    def _assert_monitoring_schedule_status_in_sync(
        self,
        sagemaker_client,
        schedule_name,
        reference,
        expected_status,
        wait_periods: int = 7,
        period_length: int = 30,
    ):
        assert (
            wait_sagemaker_monitoring_schedule_status(
                sagemaker_client, schedule_name, expected_status
            )
            == wait_resource_monitoring_schedule_status(reference, expected_status, 2)
            == expected_status
        )

    def test_smoke(
        self, sagemaker_client, xgboost_churn_data_quality_monitoring_schedule
    ):
        (reference, resource, spec) = xgboost_churn_data_quality_monitoring_schedule
        assert k8s.get_resource_exists(reference)

        monitoring_schedule_name = resource["spec"].get("monitoringScheduleName")
        monitoring_schedule_desc = get_sagemaker_monitoring_schedule(
            sagemaker_client, monitoring_schedule_name
        )
        monitoring_schedule_arn = monitoring_schedule_desc["MonitoringScheduleArn"]
        assert k8s.get_resource_arn(resource) == monitoring_schedule_arn


        self._assert_monitoring_schedule_status_in_sync(
            sagemaker_client, monitoring_schedule_name, reference, self.STATUS_SCHEDULED
        )
        assert k8s.wait_on_condition(reference, "ACK.ResourceSynced", "True")

        resource_tags = resource["spec"].get("tags", None)
        assert_tags_in_sync(monitoring_schedule_arn, resource_tags)

        # Update the resource
        new_cron_expression = "cron(0 * * * ? *)"
        spec["spec"]["monitoringScheduleConfig"]["scheduleConfig"][
            "scheduleExpression"
        ] = new_cron_expression
        resource = k8s.patch_custom_resource(reference, spec)
        resource = k8s.wait_resource_consumed_by_controller(reference)
        assert resource is not None

        self._assert_monitoring_schedule_status_in_sync(
            sagemaker_client, monitoring_schedule_name, reference, self.STATUS_SCHEDULED
        )
        assert k8s.wait_on_condition(reference, "ACK.ResourceSynced", "True")

        latest_schedule = get_sagemaker_monitoring_schedule(
            sagemaker_client, monitoring_schedule_name
        )
        assert (
            latest_schedule["MonitoringScheduleConfig"]["ScheduleConfig"][
                "ScheduleExpression"
            ]
            == new_cron_expression
        )

        # Delete the k8s resource.
        assert delete_custom_resource(
            reference, cfg.DELETE_WAIT_PERIOD, cfg.DELETE_WAIT_LENGTH
        )
        assert (
            get_sagemaker_monitoring_schedule(
                sagemaker_client, monitoring_schedule_name
            )
            is None
        )
