#!/usr/bin/env bash

# Deploy ACK Helm Charts

function install_helm_chart() {
    local service="$1"
    local oidc_role_arn="$2"
    local region="$3"
    local namespace="$4"
    local account_id=$(aws sts get-caller-identity --output text --query "Account")

    yq w -i helm/values.yaml "serviceAccount.annotations" ""
    yq w -i helm/values.yaml 'serviceAccount.annotations."eks.amazonaws.com/role-arn"' "$oidc_role_arn"
    yq w -i helm/values.yaml "aws.region" $region
    yq w -i helm/values.yaml "aws.account_id" $account_id

    kubectl create namespace $namespace
    kubectl apply -f helm/crds
    helm install -n $namespace ack-$service-controller --skip-crds helm
}