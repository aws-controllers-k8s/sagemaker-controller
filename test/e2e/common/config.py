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

# Job Type Resource Statuses 
LIST_JOB_STATUS_STOPPED = (
    "Stopped", 
    "Stopping", 
    "Completed"
    )
JOB_STATUS_INPROGRESS: str = "InProgress"
JOB_STATUS_COMPLETED: str = "Completed"
DEBUGGERJOB_STATUS_COMPLETED: str = (
        "NoIssuesFound", 
        "Completed",
        "CompletedWithIssues"
    )