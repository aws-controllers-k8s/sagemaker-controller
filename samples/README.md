# Job Sample Overview

This sample demonstrates how to start jobs using your own script, packaged in a SageMaker-compatible container, using the Amazon AWS Controllers for Kubernetes (ACK) service controller for Amazon SageMaker.                     

## Prerequisites    

Follow the instructions on our [installation page](/README.md#getting-started) to create a Kubernetes cluster and install sagemaker controller.

Run the following commands to create a SageMaker execution IAM role which is used by SageMaker service to access AWS resources. 

```
export SAGEMAKER_EXECUTION_ROLE_NAME=ack-sagemaker-execution-role-$CLUSTER_NAME

TRUST="{ \"Version\": \"2012-10-17\", \"Statement\": [ { \"Effect\": \"Allow\", \"Principal\": { \"Service\": \"sagemaker.amazonaws.com\" }, \"Action\": \"sts:AssumeRole\" } ] }"
aws iam create-role --role-name ${SAGEMAKER_EXECUTION_ROLE_NAME} --assume-role-policy-document "$TRUST"
aws iam attach-role-policy --role-name ${SAGEMAKER_EXECUTION_ROLE_NAME} --policy-arn arn:aws:iam::aws:policy/AmazonSageMakerFullAccess
aws iam attach-role-policy --role-name ${SAGEMAKER_EXECUTION_ROLE_NAME} --policy-arn arn:aws:iam::aws:policy/AmazonS3FullAccess

SAGEMAKER_EXECUTION_ROLE_ARN=$(aws iam get-role --role-name ${SAGEMAKER_EXECUTION_ROLE_NAME} --output text --query 'Role.Arn')

echo $SAGEMAKER_EXECUTION_ROLE_ARN
```
Note down the execution role ARN to use in samples
### Creating your first Job

The easiest way to start is taking a look at the sample training jobs and its corresponding [README](/samples/training/README.md)