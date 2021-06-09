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
"""Declares the structure of the bootstrapped resources and provides a loader
for them.
"""

from dataclasses import dataclass
from acktest.resources import read_bootstrap_config
from e2e import bootstrap_directory

SAGEMAKER_SOURCE_DATA_BUCKET = "source-data-bucket-592697580195-us-west-2"


@dataclass
class TestBootstrapResources:
    DataBucketName: str
    ExecutionRoleARN: str


_bootstrap_resources = None


def get_bootstrap_resources(bootstrap_file_name: str = "bootstrap.yaml"):
    global _bootstrap_resources
    if _bootstrap_resources is None:
        _bootstrap_resources = TestBootstrapResources(
            **read_bootstrap_config(
                bootstrap_directory, bootstrap_file_name=bootstrap_file_name
            ),
        )
    return _bootstrap_resources
