# Job Sample Overview

This sample demonstrates how to start jobs using your own script, packaged in a SageMaker-compatible container, using the Amazon AWS Controllers for Kubernetes (ACK) service controller for Amazon SageMaker.                     

## Prerequisites    

Follow the instructions on our [installation page](/README.md#getting-started) to create a Kubernetes cluster and install sagemaker controller.

### SageMaker execution IAM role
Run the following commands to create a SageMaker execution IAM role which is used by SageMaker service to access AWS resources. 

```
export SAGEMAKER_EXECUTION_ROLE_NAME=ack-sagemaker-execution-role

TRUST="{ \"Version\": \"2012-10-17\", \"Statement\": [ { \"Effect\": \"Allow\", \"Principal\": { \"Service\": \"sagemaker.amazonaws.com\" }, \"Action\": \"sts:AssumeRole\" } ] }"
aws iam create-role --role-name ${SAGEMAKER_EXECUTION_ROLE_NAME} --assume-role-policy-document "$TRUST"
aws iam attach-role-policy --role-name ${SAGEMAKER_EXECUTION_ROLE_NAME} --policy-arn arn:aws:iam::aws:policy/AmazonSageMakerFullAccess
aws iam attach-role-policy --role-name ${SAGEMAKER_EXECUTION_ROLE_NAME} --policy-arn arn:aws:iam::aws:policy/AmazonS3FullAccess

SAGEMAKER_EXECUTION_ROLE_ARN=$(aws iam get-role --role-name ${SAGEMAKER_EXECUTION_ROLE_NAME} --output text --query 'Role.Arn')

echo $SAGEMAKER_EXECUTION_ROLE_ARN
```
Note down the execution role ARN to use in samples.

### S3 Bucket
Run the following commands to create a S3 bucket which is used by SageMaker service to access and upload data. Use the region you used during installation unless you are trying cross region. SageMaker resources will be created in this region.
```
export AWS_DEFAULT_REGION=<REGION>
export S3_BUCKET_NAME="ack-data-bucket-$AWS_DEFAULT_REGION"

# [Option 1] if your region is us-east-1
aws s3api create-bucket --bucket $S3_BUCKET_NAME --region $AWS_DEFAULT_REGION

# [Option 2] if your region is NOT us-east-1
aws s3api create-bucket --bucket $S3_BUCKET_NAME --region $AWS_DEFAULT_REGION \
--create-bucket-configuration LocationConstraint=$AWS_DEFAULT_REGION

echo $S3_BUCKET_NAME
```
Note down your S3 bucket name which will be used in the samples.

### Creating your first Job

The easiest way to start is taking a look at the sample training jobs and its corresponding [README](/samples/training/README.md)