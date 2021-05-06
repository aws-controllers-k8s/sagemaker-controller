#!/bin/bash

# cleanup on EXIT regardles of error 

# Inputs to this file as environment variables
# SERVICE
# SERVICE_REGION

function cleanup {
  echo "Cleaning up resources"
  cd $SERVICE_REPO_PATH_DOCKER/test/e2e/
  python ./cleanup.py $SERVICE
}
trap cleanup EXIT

# Setup OIDC
. ./test/e2e/canary/scripts/setup_oidc.sh

# Install service helm chart
. ./test/e2e/canary/scripts/install_controller_helm.sh

# create resources for test
cd $SERVICE_REPO_PATH_DOCKER/test/e2e/

export AWS_ROLE_ARN=$(aws sts get-caller-identity --query "Arn")
export AWS_DEFAULT_REGION=$SERVICE_REGION 

python ./bootstrap.py $SERVICE
sleep 10m

# TOOODOOOOO: RUN ALL TESTS run tests
echo "Run Tests"
PYTHONPATH=. pytest -n 10 --dist loadfile --log-cli-level INFO $SERVICE -m canary tests/test_model.py