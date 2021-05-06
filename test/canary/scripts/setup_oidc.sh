# OIDC Setup 

# Inputs to this file as environment variables 
# CLUSTER_REGION
# CLUSTER_NAME

NAMESPACE=${NAMESPACE:-"ack-system"}

AWS_ACC_NUM=$(aws sts get-caller-identity --output text --query "Account")
aws --region $CLUSTER_REGION eks update-kubeconfig --name $CLUSTER_NAME
eksctl utils associate-iam-oidc-provider --cluster $CLUSTER_NAME --region $CLUSTER_REGION --approve

OIDC_URL=$(aws eks describe-cluster --region $CLUSTER_REGION --name $CLUSTER_NAME  --query "cluster.identity.oidc.issuer" --output text | cut -c9-)

cat <<EOF > trust.json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Federated": "arn:aws:iam::$AWS_ACC_NUM:oidc-provider/$OIDC_URL"
      },
      "Action": "sts:AssumeRoleWithWebIdentity",
      "Condition": {
        "StringEquals": {
          "$OIDC_URL:aud": "sts.amazonaws.com",
          "$OIDC_URL:sub": ["system:serviceaccount:${NAMESPACE}:ack-sagemaker-controller"]
        }
      }
    }
  ]
}
EOF


# TODO : check if iam role exists 
aws iam create-role --role-name ack-oidc-role-$CLUSTER_NAME --assume-role-policy-document file://trust.json
aws iam attach-role-policy --role-name ack-oidc-role-$CLUSTER_NAME --policy-arn arn:aws:iam::aws:policy/AmazonSageMakerFullAccess
aws iam attach-role-policy --role-name ack-oidc-role-$CLUSTER_NAME --policy-arn arn:aws:iam::aws:policy/AmazonS3FullAccess

export OIDC_ROLE_ARN=$(aws iam get-role --role-name ack-oidc-role-$CLUSTER_NAME --output text --query 'Role.Arn')