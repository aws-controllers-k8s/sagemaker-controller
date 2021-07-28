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
"""Integration tests for the Notebook Lifecycle configuration
"""

import pytest
import logging
import botocore

from acktest.k8s import resource as k8s
from e2e import (
    service_marker,
    wait_for_status,
    create_sagemaker_resource,
    sagemaker_client,
)

from e2e.bootstrap_resources import get_bootstrap_resources
import random

RESOURCE_PLURAL = "notebookinstancelifecycleconfigs"
RESOURCE_PREFIX = "nblf"
RESOURCE_SPEC_FILE = "notebook_instance"