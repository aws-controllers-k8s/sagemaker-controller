#!/bin/bash

# cleanup on EXIT regardles of error

# Inputs to this file as environment variables
# SERVICE
# SERVICE_REGION
# CLUSTER_REGION
# CLUSTER_NAME
# SERVICE_REPO_PATH
# NAMESPACE

set -euo pipefail
export NAMESPACE=${NAMESPACE:-"ack-system"}
export AWS_DEFAULT_REGION=$SERVICE_REGION
export E2E_DIR=$SERVICE_REPO_PATH/test/e2e/
SCRIPTS_DIR=${SERVICE_REPO_PATH}/test/canary/scripts

source $SCRIPTS_DIR/setup_oidc.sh
source $SCRIPTS_DIR/install_controller_helm.sh

function print_controller_logs() {
  pod_id=$( kubectl get pods -n $NAMESPACE --field-selector="status.phase=Running" \
      --sort-by=.metadata.creationTimestamp \
      | grep ack-sagemaker-controller | awk '{print $1}' 2>/dev/null )

  kubectl -n $NAMESPACE logs "$pod_id"
}

function cleanup {
  #push to metrics to cloudwatch
  echo "Pushing Codebuild stats to Cloudwatch."
  cd $SCRIPTS_DIR
  python push_stats_to_cloudwatch.py

  echo "Cleaning up resources"
  set +e
  kubectl delete monitoringschedules --all
  kubectl delete endpoints.sagemaker --all
  kubectl delete endpointconfigs --all
  kubectl delete models --all
  kubectl delete trainingjobs --all
  kubectl delete processingjobs --all
  kubectl delete transformjobs --all
  kubectl delete hyperparametertuningjobs --all
  kubectl delete dataqualityjobdefinitions --all
  kubectl delete modelbiasjobdefinitions --all
  kubectl delete modelexplainabilityjobdefinitions --all
  kubectl delete modelqualityjobdefinitions --all
  # kubectl delete adoptedresources --all
  kubectl delete featuregroups --all
  kubectl delete modelpackages --all
  kubectl delete modelpackagegroups --all
  kubectl delete notebookinstances --all
  kubectl delete notebookinstancelifecycleconfig --all
  kubectl delete pipelineexecutions --all
  kubectl delete pipelines --all

  print_controller_logs

  helm delete -n $NAMESPACE ack-$SERVICE-controller
  pushd $SERVICE_REPO_PATH
    kubectl delete -f helm/crds
  popd
  kubectl delete namespace $NAMESPACE

  cd $E2E_DIR
  export PYTHONPATH=..
  python service_cleanup.py

}
trap cleanup EXIT


# Update kubeconfig
aws --region $CLUSTER_REGION eks update-kubeconfig --name $CLUSTER_NAME

# Setup OIDC
create_oidc_role "$CLUSTER_NAME" "$CLUSTER_REGION" "$NAMESPACE"

# Install service helm chart
install_helm_chart $SERVICE $OIDC_ROLE_ARN $SERVICE_REGION $NAMESPACE

echo "Log helm charts are deployed properly"
kubectl -n $NAMESPACE get pods
kubectl get crds

pushd $E2E_DIR
  export PYTHONPATH=..
  # create resources for test
  python service_bootstrap.py
  sleep 5m

  # run tests
  echo "Run Tests"
  pytest_args=( -n 15 --dist loadfile --log-cli-level INFO --junitxml ../canary/integration_tests.xml)
  declare pytest_marks
  if [[ $SERVICE_REGION =~ ^(eu-north-1|eu-west-3)$  ]]; then
    # If select_regions_1 true we run the notebook_instance test
    pytest_marks="canary or select_regions_1"
  elif [[ $SHALLOW_REGION = "shallow"  ]]; then
    pytest_marks="shallow_canary"
  else
    pytest_marks="canary"
  fi
  if [[ $SERVICE_REGION =~ ^(eu-north-1|ap-south-1|ap-southeast-3|us-east-2|me-central-1|eu-west-1|eu-central-1|sa-east-1|us-east-1|ap-northeast-2|eu-west-2|ap-northeast-1|us-west-2|ap-southeast-1|ap-southeast-2|ca-central-1)$  ]]; then
    # Above is the list of supported regions for Inference Component and if the current region is
    # included in this we will add the inference_component mark.
      pytest_marks+=" or inference_component"
  fi
  pytest_args+=(-m "$pytest_marks")
  pytest "${pytest_args[@]}"
popd
