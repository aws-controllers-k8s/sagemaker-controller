# Bring Your Own Container Sample

This sample demonstrates how to start data-quality jobs using your own data-quality script, packaged in a SageMaker-compatible container, using the Amazon AWS Controllers for Kubernetes (ACK) service controller for Amazon SageMaker.                     

## Prerequisites

This sample assumes that you have already configured an EKS cluster with the ACK operator. It also assumes that you have installed `kubectl` - you can find a link on our [installation page](TODO).

In order to follow this script, you must first create a data-quality script packaged in a Dockerfile that is [compatible with Amazon SageMaker](https://docs.aws.amazon.com/sagemaker/latest/dg/amazon-sagemaker-containers.html).

You will also need an Endpoint configured in SageMaker you can create them using the `endpoint_config.yaml` and `endpoint_base.yaml` or through the SageMaker console.

## Preparing the data-quality Script

### Uploading your Script

All SageMaker data-quality jobs are run from within a container with all necessary dependencies and modules pre-installed and with the data-quality scripts referencing the acceptable input and output directories. This container should be uploaded to an [ECR repository](https://aws.amazon.com/ecr/) accessible from within your AWS account. When uploaded correctly, you should have a repository URL and tag associated with the container image - this will be needed for the next step. Sample container images are [available](https://docs.aws.amazon.com/sagemaker/latest/dg/ecr-us-west-2.html).


A container image URL and tag looks has the following structure:
```
<account number>.dkr.ecr.<region>.amazonaws.com/<image name>:<tag>
```

### Updating the data-quality Specification

In the `my-data-quality-job.yaml` file, modify the placeholder values with those associated with your account and data-quality job. The `dataqualityAppSpecification.imageURI` should be the container image from the previous step. The `spec.roleARN` field should be the ARN of an IAM role which has permissions to access your S3 resources. The `dataqualityAppSpecification.dataQualityJobInput.endpointInput.endpointName` should be the name of your Endpoint in SageMaker. If you have not yet created a role with these permissions, you can find an example policy at [Amazon SageMaker Roles](https://docs.aws.amazon.com/sagemaker/latest/dg/sagemaker-roles.html#sagemaker-roles-createdata-qualityjob-perms).

## Submitting your data-quality Job

To submit your prepared data-quality job specification, apply the specification to your EKS cluster as such:
```
$ kubectl apply -f my-data-quality-job.yaml
dataqualityjobdefinitions.sagemaker.services.k8s.aws.amazon.com/my-data-quality-job created
```

To monitor the data-quality job status, you can use the following command:
```
$ kubectl get dataqualityjobdefinitions my-data-quality-job
```

To monitor the data-quality job in-depth once it has started, you can see the full status and any additional errors with the following command:
```
$ kubectl describe dataqualityjobdefinitions my-data-quality-job
```

To delete the data-quality job once it has started if errors occurred or for any reason with the following command:
```
$ kubectl delete dataqualityjobdefinitions my-data-quality-job
```
