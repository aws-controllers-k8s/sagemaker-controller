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
  kubectl delete adoptedresources --all
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
  pytest_args=( -n 15 --dist loadfile --log-cli-level INFO --junitxml=report.xml )
  if [[ $SERVICE_REGION =~ ^(eu-north-1|eu-west-3)$  ]]; then
    # If select_regions_1 true we run the notebook_instance test
    pytest_args+=(-m "canary or select_regions_1")
  else
    pytest_args+=(-m "canary")
  pytest "${pytest_args[@]}"
  fi
popd
