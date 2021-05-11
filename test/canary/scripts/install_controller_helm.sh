#!/usr/bin/env bash

# Deploy ACK Helm Charts

function install_helm_chart() {
    local service="$1"
    local oidc_role_arn="$2"
    local region="$3"
    local namespace="$4"

    yq w -i helm/values.yaml "serviceAccount.annotations" ""
    yq w -i helm/values.yaml 'serviceAccount.annotations."eks.amazonaws.com/role-arn"' "$oidc_role_arn"
    yq w -i helm/values.yaml "aws.region" $region

    kubectl create namespace $namespace
    kubectl apply -f helm/crds
    helm install -n $namespace ack-$service-controller --skip-crds helm
}