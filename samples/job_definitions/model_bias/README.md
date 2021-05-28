# Model Bias Job Definition Sample

This sample demonstrates how to start model-bias job definitions using your own model-bias job definition script, packaged in a SageMaker-compatible container, using the Amazon AWS Controllers for Kubernetes (ACK) service controller for Amazon SageMaker.                     

## Prerequisites

This sample assumes that you have completed the [common prerequisites](/samples/README.md).

You will also need an [Endpoint](/samples/endpoint/README.md) configured in SageMaker and you will need to run a baselining job to generate baseline constraints only.

### Get an Image

All SageMaker model-bias job definitions are run from within a container with all necessary dependencies and modules pre-installed and with the model-bias scripts referencing the acceptable input and output directories. Sample container images are [available](https://docs.aws.amazon.com/sagemaker/latest/dg/sagemaker-algo-docker-registry-paths.html).

A container image URL and tag looks has the following structure:
```
<account number>.dkr.ecr.<region>.amazonaws.com/<image name>:<tag>
```

### Updating the Model Bias Job Definition Specification

In the `my-model-bias-job-definition.yaml` file, modify the placeholder values with those associated with your account and model-bias job definition.

## Submitting your Model Bias Job Definition

### Create a Model Bias Job Definition

To submit your prepared model-bias job definition specification, apply the specification to your Kubernetes cluster as such:
```
$ kubectl apply -f my-model-bias-job-definition.yaml
modelbiasjobdefinitions.sagemaker.services.k8s.aws.amazon.com/my-model-bias-job-definition created
```

### List Model Bias Job Definitions

To list all Model Bias Job Definition created using the ACK controller use the following command:
```
$ kubectl get modelbiasjobdefinitions
```

### Describe a Model Bias Job Definition

To get more details about the Model Bias Job Definition once it's submitted, like checking the status, errors or parameters of the Model Bias Job Definition use the following command:
```
$ kubectl describe modelbiasjobdefinitions my-model-bias-job-definition
```
You can also check Status.ackResourceMetadata.Arn to verify the data quality job definition was created successfully.

### Delete a Model Bias Job Definition

To delete the model-bias job definition, use the following command:
```
$ kubectl delete modelbiasjobdefinitions my-model-bias-job-definition
```