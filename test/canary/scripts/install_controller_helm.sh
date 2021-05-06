# Deploy ACK Helm Charts 

# Inputs to this file as environment variables 
# SERVICE_REPO_PATH_DOCKER
# OIDC_ROLE_ARN
# SERVICE
# SERVICE_REGION

cd $SERVICE_REPO_PATH_DOCKER

yq w -i helm/values.yaml "serviceAccount.annotations" ""
yq w -i helm/values.yaml 'serviceAccount.annotations."eks.amazonaws.com/role-arn"' "$OIDC_ROLE_ARN"
yq w -i helm/values.yaml "aws.region" $SERVICE_REGION

export ACK_K8S_NAMESPACE=${NAMESPACE:-"ack-system"}
kubectl create namespace $ACK_K8S_NAMESPACE

helm delete -n $ACK_K8S_NAMESPACE ack-$SERVICE-controller 
helm install -n $ACK_K8S_NAMESPACE ack-$SERVICE-controller helm 

echo "Make sure helm charts are deployed properly"
kubectl -n $ACK_K8S_NAMESPACE get pods 
kubectl get crds