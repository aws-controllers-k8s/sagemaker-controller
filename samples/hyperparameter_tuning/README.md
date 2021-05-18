# Bring Your Own Container Sample

This sample demonstrates how to start hyperparameter jobs using your own hyperparameter script, packaged in a SageMaker-compatible container, using the Amazon AWS Controllers for Kubernetes (ACK) service controller for Amazon SageMaker.                     

## Prerequisites

This sample assumes that you have already configured an EKS cluster with the ACK operator. It also assumes that you have installed `kubectl` - you can find a link on our [installation page](TODO).

In order to follow this script, you must first create a hyperparameter script packaged in a Dockerfile that is [compatible with Amazon SageMaker](https://docs.aws.amazon.com/sagemaker/latest/dg/amazon-sagemaker-containers.html). 

## Preparing the HyperparameterTuning Script

### Uploading your Script

All SageMaker Hyperparameter jobs are run from within a container with all necessary dependencies and modules pre-installed and with the hyperparameter scripts referencing the acceptable input and output directories. This container should be uploaded to an [ECR repository](https://aws.amazon.com/ecr/) accessible from within your AWS account. When uploaded correctly, you should have a repository URL and tag associated with the container image - this will be needed for the next step. Sample container images are [available](https://docs.aws.amazon.com/sagemaker/latest/dg/ecr-us-west-2.html).

A container image URL and tag looks has the following structure:
```
<account number>.dkr.ecr.<region>.amazonaws.com/<image name>:<tag>
```

### Updating the Hyperparameter Specification

In the `my-hyperparameter-job.yaml` file, modify the placeholder values with those associated with your account and hyperparameter job. The `spec.algorithmSpecification.hyperparameterImage` should be the container image from the previous step. The `spec.roleARN` field should be the ARN of an IAM role which has permissions to access your S3 resources. If you have not yet created a role with these permissions, you can find an example policy at [Amazon SageMaker Roles](https://docs.aws.amazon.com/sagemaker/latest/dg/sagemaker-roles.html#sagemaker-roles-createhyperparametertuningjob-perms). 


### Enabling Spot Training
In the `my-hyperparameter-job.yaml` file under `spec.trainingJobDefinition` add `enableManagedSpotTraining` and set the value to true. You will also need to specify a `spec.trainingJobDefinition.stoppingCondition.maxRuntimeInSeconds` and `spec.trainingJobDefinition.stoppingCondition.maxWaittimeInSeconds`

## Submitting your Hyperparameter Job

To submit your prepared hyperparameter job specification, apply the specification to your EKS cluster as such:
```
$ kubectl apply -f my-hyperparameter-job.yaml
hyperparametertuningjob.sagemaker.services.k8s.aws.amazon.com/my-hyperparameter-job created
```

To monitor the hyperparameter job status, you can use the following command:
```
$ kubectl get hyperparametertuningjob my-hyperparameter-job
```

To monitor the hyperparameter job in-depth once it has started, you can see the full status and any additional errors with the following command:
```
$ kubectl describe hyperparametertuningjob my-hyperparameter-job
```

To delete the hyperparameter job once it has started if errors occurred or for any reason with the following command:
```
$ kubectl delete hyperparametertuningjob my-hyperparameter-job
```
