#!/usr/bin/env bash
# OIDC Setup

# A function to get the OIDC_ID associated with an EKS cluster
function get_oidc_id() {
  local cluster_name="$1"
  local region="$2"
  eksctl utils associate-iam-oidc-provider --cluster $cluster_name --region $region --approve > /dev/null
  local oidc_url=$(aws eks describe-cluster --region $region --name $cluster_name  --query "cluster.identity.oidc.issuer" --output text | cut -c9-)
  echo "${oidc_url}"
}


function generate_trust_policy() {
  local oidc_url="$1"
  local namespace="$2"
  local account_id=$(aws sts get-caller-identity --output text --query "Account")

  cat <<EOF > trust.json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Federated": "arn:aws:iam::${account_id}:oidc-provider/${oidc_url}"
      },
      "Action": "sts:AssumeRoleWithWebIdentity",
      "Condition": {
        "StringEquals": {
          "${oidc_url}:aud": "sts.amazonaws.com",
          "${oidc_url}:sub": ["system:serviceaccount:${namespace}:ack-sagemaker-controller"]
        }
      }
    }
  ]
}
EOF
}

function create_oidc_role() {
  local cluster_name="$1"
  local region="$2"
  local namespace="$3"
  local oidc_role_name=ack-oidc-role-$cluster_name-$namespace
  
  # Create role only if it does not exist
  set +e
  aws iam get-role --role-name ${oidc_role_name}
  exit_code=$?
  set -euo pipefail

  if [[ $exit_code -eq 0 ]]; then
    echo "A role for this cluster and namespace already exists in this account, assuming sagemaker access and proceeding."
  else
    echo "Creating new IAM role: $oidc_role_name"
    local oidc_url=$(get_oidc_id "$cluster_name" "$region")
    local trustfile="trust.json"
    generate_trust_policy "$oidc_url" "$namespace"
    aws iam create-role --role-name "$oidc_role_name" --assume-role-policy-document file://${trustfile}
    aws iam attach-role-policy --role-name "$oidc_role_name" --policy-arn arn:aws:iam::aws:policy/AmazonSageMakerFullAccess
    aws iam attach-role-policy --role-name "$oidc_role_name" --policy-arn arn:aws:iam::aws:policy/AmazonS3FullAccess
    rm "${trustfile}" 
  fi
  local oidc_role_arn=$(aws iam get-role --role-name $oidc_role_name --output text --query 'Role.Arn')
  export OIDC_ROLE_ARN=$oidc_role_arn
}