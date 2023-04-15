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
"""Bootstraps the resources required to run the SageMaker integration tests.
"""

import boto3
import json
import logging
import time
import subprocess
import random

from acktest import resources
from acktest.bootstrapping import Resources, BootstrapFailureException
from acktest.bootstrapping.iam import Role
from acktest.bootstrapping.s3 import Bucket
from acktest.aws.identity import get_region, get_account_id
from acktest.aws.s3 import duplicate_bucket_contents
from e2e import bootstrap_directory
from e2e.bootstrap_resources import TestBootstrapResources, SAGEMAKER_SOURCE_DATA_BUCKET


def sync_data_bucket(bucket) -> str:
    bucket_name = bucket.name
    region = get_region()
    s3_resource = boto3.resource("s3", region_name=region)

    source_bucket = s3_resource.Bucket(SAGEMAKER_SOURCE_DATA_BUCKET)
    destination_bucket = s3_resource.Bucket(bucket_name)
    temp_dir = "/tmp/ack_s3_data"
    # awscli is not installed in test-infra container hence use boto3 to copy in us-west-2
    if region == "us-west-2":
        duplicate_bucket_contents(source_bucket, destination_bucket)
        # above method does an async copy
        # TODO: find a way to remove random wait
        time.sleep(180)
    else:
        # workaround to copy if buckets are across regions
        # TODO: check if there is a better way and merge to test-infra
        subprocess.call(["mkdir", f"{temp_dir}"])
        subprocess.call(
            [
                "aws",
                "s3",
                "sync",
                f"s3://{SAGEMAKER_SOURCE_DATA_BUCKET}",
                f"./{temp_dir}/",
                "--region",
                "us-west-2"
                "--quiet",
            ]
        )
        subprocess.call(
            ["aws", "s3", "sync", f"./{temp_dir}/", f"s3://{bucket_name}", "--quiet"]
        )

    logging.info(f"Synced data bucket")

    return bucket


def service_bootstrap() -> Resources:
    logging.getLogger().setLevel(logging.INFO)
    region = get_region()
    account_id = get_account_id()
    bucket_name = f"ack-data-bucket-{region}-{account_id}"

    resources = TestBootstrapResources(
        DataBucket=Bucket(bucket_name),
        ExecutionRole=Role(
            "ack-sagemaker-execution-role",
            "sagemaker.amazonaws.com",
            managed_policies=[
                "arn:aws:iam::aws:policy/AmazonSageMakerFullAccess",
                "arn:aws:iam::aws:policy/AmazonS3FullAccess",
            ],
        ),
    )
    try:
        resources.bootstrap()
        sync_data_bucket(resources.DataBucket)
    except BootstrapFailureException as ex:
        exit(254)
    return resources


if __name__ == "__main__":
    config = service_bootstrap()
    # Write config to current directory by default
    config.serialize(bootstrap_directory)
