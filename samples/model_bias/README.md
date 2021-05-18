# Bring Your Own Container Sample

This sample demonstrates how to start model-bias jobs using your own model-bias script, packaged in a SageMaker-compatible container, using the Amazon AWS Controllers for Kubernetes (ACK) service controller for Amazon SageMaker.                     

## Prerequisites

This sample assumes that you have already configured an EKS cluster with the ACK operator. It also assumes that you have installed `kubectl` - you can find a link on our [installation page](TODO).

In order to follow this script, you must first create a model-bias script packaged in a Dockerfile that is [compatible with Amazon SageMaker](https://docs.aws.amazon.com/sagemaker/latest/dg/amazon-sagemaker-containers.html).

You will also need an Endpoint configured in SageMaker you can create them using the `endpoint_config.yaml` and `endpoint_base.yaml` or through the SageMaker console.

## Preparing the model-bias Script

### Uploading your Script

All SageMaker model-bias jobs are run from within a container with all necessary dependencies and modules pre-installed and with the model-bias scripts referencing the acceptable input and output directories. This container should be uploaded to an [ECR repository](https://aws.amazon.com/ecr/) accessible from within your AWS account. When uploaded correctly, you should have a repository URL and tag associated with the container image - this will be needed for the next step. Sample container images are [available](https://docs.aws.amazon.com/sagemaker/latest/dg/ecr-us-west-2.html).


A container image URL and tag looks has the following structure:
```
<account number>.dkr.ecr.<region>.amazonaws.com/<image name>:<tag>
```

### Updating the model-bias Specification

In the `my-model-bias-job.yaml` file, modify the placeholder values with those associated with your account and model-bias job. The `modelBiasAppSpecification.imageURI` should be the container image from the previous step. The `spec.roleARN` field should be the ARN of an IAM role which has permissions to access your S3 resources. The `modelBiasAppSpecification.modelBiasJobInput.endpointInput.endpointName` should be the name of your Endpoint in SageMaker. If you have not yet created a role with these permissions, you can find an example policy at [Amazon SageMaker Roles](https://docs.aws.amazon.com/sagemaker/latest/dg/sagemaker-roles.html#sagemaker-roles-createmodel-biasjob-perms).

## Submitting your model-bias Job

To submit your prepared model-bias job specification, apply the specification to your EKS cluster as such:
```
$ kubectl apply -f my-model-bias-job.yaml
modelbiasjobdefinitions.sagemaker.services.k8s.aws.amazon.com/my-model-bias-job created
```

To monitor the model-bias job status, you can use the following command:
```
$ kubectl get modelbiasjobdefinitions my-model-bias-job
```

To monitor the model-bias job in-depth once it has started, you can see the full status and any additional errors with the following command:
```
$ kubectl describe modelbiasjobdefinitions my-model-bias-job
```

To delete the model-bias job once it has started if errors occurred or for any reason with the following command:
```
$ kubectl delete modelbiasjobdefinitions my-model-bias-job
```
