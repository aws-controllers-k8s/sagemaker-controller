#!/usr/bin/env bash

# Deploy ACK Helm Charts

function install_helm_chart() {
    local service="$1"
    local oidc_role_arn="$2"
    local region="$3"
    local namespace="$4"

    yq eval ".serviceAccount.annotations = {}" -i helm/values.yaml
    yq eval ".serviceAccount.annotations.\"eks.amazonaws.com/role-arn\" = \"$oidc_role_arn\"" -i helm/values.yaml
    yq eval ".aws.region = \"$region\"" -i helm/values.yaml
    yq eval '.log.level = "debug"' -i helm/values.yaml
    yq eval '.log.enable_development_logging = true' -i helm/values.yaml

    kubectl apply -f helm/crds
    helm install -n $namespace --create-namespace ack-$service-controller --skip-crds helm
}
