# Model Quality Job Definition Sample

This sample demonstrates how to start model-quality job definitions using your own model-quality script, packaged in a SageMaker-compatible container, using the Amazon AWS Controllers for Kubernetes (ACK) service controller for Amazon SageMaker.                     

## Prerequisites

This sample assumes that you have completed the [common prerequisites](/samples/README.md).

You will also need an [Endpoint](/samples/endpoint/README.md) configured in SageMaker and you will need to run a baselining job to generate baseline constraints only.

### Get an Image

All SageMaker model-quality job definitions are run from within a container with all necessary dependencies and modules pre-installed and with the model-quality scripts referencing the acceptable input and output directories. Sample container images are [available](https://docs.aws.amazon.com/sagemaker/latest/dg/sagemaker-algo-docker-registry-paths.html).

A container image URL and tag looks has the following structure:
```
<account number>.dkr.ecr.<region>.amazonaws.com/<image name>:<tag>
```

### Updating the Model-Quality Specification

In the `my-model-quality-job-definition.yaml` file, modify the placeholder values with those associated with your account and model-quality job definition.

## Submitting your Model-Quality Job

### Create a Model-Quality Job Definition

To submit your prepared model-quality job definition specification, apply the specification to your Kubernetes cluster as such:
```
$ kubectl apply -f my-model-quality-job-definition.yaml
modelQualityjobdefinitions.sagemaker.services.k8s.aws.amazon.com/my-model-quality-job-definition created
```

### List Model-Quality Job Definitions

To list all Model-Quality Job Definition created using the ACK controller use the following command:
```
$ kubectl get modelQualityjobdefinitions
```

### Describe a Model-Quality Job Definition

To get more details about the Model-Quality Job Definition once it's submitted, like checking the status, errors or parameters of the Model-Quality Job Definition use the following command:
```
$ kubectl describe modelQualityjobdefinitions my-model-quality-job-definition
```
You can also check Status.ackResourceMetadata.Arn to verify the data quality job definition was created successfully.

### Delete a Model-Quality Job Definition

To delete the model-quality job definition, use the following command:
```
$ kubectl delete modelQualityjobdefinitions my-model-quality-job-definition
```