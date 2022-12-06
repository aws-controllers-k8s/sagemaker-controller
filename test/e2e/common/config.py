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
"""String Constants
"""

ENDPOINT_CONFIG_RESOURCE_PLURAL = "endpointconfigs"
MODEL_RESOURCE_PLURAL = "models"
ENDPOINT_RESOURCE_PLURAL = "endpoints"
DATA_QUALITY_JOB_DEFINITION_RESOURCE_PLURAL = "dataqualityjobdefinitions"
MODEL_PACKAGE_GROUP_RESOURCE_PLURAL = "modelpackagegroups"
MODEL_PACKAGE_RESOURCE_PLURAL = "modelpackages"

# Job Type Resource Statuses
LIST_JOB_STATUS_STOPPED = ("Stopped", "Stopping", "Completed")
JOB_STATUS_INPROGRESS: str = "InProgress"
JOB_STATUS_COMPLETED: str = "Completed"
JOB_STATUS_EXECUTING: str = "Executing"
JOB_STATUS_SUCCEEDED: str = "Succeeded"
RULE_STATUS_COMPLETED: str = "NoIssuesFound"

ENDPOINT_STATUS_INSERVICE = "InService"
ENDPOINT_STATUS_CREATING = "Creating"
ENDPOINT_STATUS_UPDATING = "Updating"

DELETE_WAIT_PERIOD = 4
DELETE_WAIT_LENGTH = 30

JOB_DELETE_WAIT_PERIODS = 12
JOB_DELETE_WAIT_LENGTH = 30

TAG_DELAY_SLEEP = 20
