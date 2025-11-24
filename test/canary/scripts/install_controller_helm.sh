#!/usr/bin/env bash

# Deploy ACK Helm Charts

function install_helm_chart() {
    local service="$1"
    local oidc_role_arn="$2"
    local region="$3"
    local namespace="$4"

    yq -i e '.serviceAccount.annotations = {}' helm/values.yaml
    yq -i e '.serviceAccount.annotations."eks.amazonaws.com" = "'"$oidc_role_arn"'"' helm/values.yaml
    yq -i e '.aws.region = "'"$region"'"' helm/values.yaml
    yq -i e '.log.level = "debug"' helm/values.yaml
    yq -i e '.log.enable_development_logging = true' helm/values.yaml

    kubectl apply -f helm/crds
    helm install -n $namespace --create-namespace ack-$service-controller --skip-crds helm
}
