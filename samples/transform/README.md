# Bring Your Own Container Sample

This sample demonstrates how to start transform jobs using your own transform script, packaged in a SageMaker-compatible container, using the Amazon AWS Controllers for Kubernetes (ACK) service controller for Amazon SageMaker.                     

## Prerequisites

This sample assumes that you have already configured an EKS cluster with the ACK operator. It also assumes that you have installed `kubectl` - you can find a link on our [installation page](TODO).

In order to follow this script, you must first create a transform script packaged in a Dockerfile that is [compatible with Amazon SageMaker](https://docs.aws.amazon.com/sagemaker/latest/dg/amazon-sagemaker-containers.html). 

You will also need a model in SageMaker for this sample. If you do not have one you must first create a [model](https://docs.aws.amazon.com/sagemaker/latest/dg/sagemaker-mkt-model-pkg-model.html)
## Preparing the transform Script


### Updating the transform Specification

In the `my-transform-job.yaml` file, modify the placeholder values with those associated with your account and transform job. The `spec.algorithmSpecification.modelName` should be the model from the previous step. The `spec.roleARN` field should be the ARN of an IAM role which has permissions to access your S3 resources. If you have not yet created a role with these permissions, you can find an example policy at [Amazon SageMaker Roles](https://docs.aws.amazon.com/sagemaker/latest/dg/sagemaker-roles.html#sagemaker-roles-createtransformjob-perms).

## Submitting your transform Job

To submit your prepared transform job specification, apply the specification to your EKS cluster as such:
```
$ kubectl apply -f my-transform-job.yaml
transformjob.sagemaker.services.k8s.aws.amazon.com/my-transform-job created
```

To monitor the transform job status, you can use the following command:
```
$ kubectl get transformjob my-transform-job
```

To monitor the transform job in-depth once it has started, you can see the full status and any additional errors with the following command:
```
$ kubectl describe transformjob my-transform-job
```

To delete the transform job once it has started if errors occurred or for any reason with the following command:
```
$ kubectl delete transformjob my-transform-job
```
