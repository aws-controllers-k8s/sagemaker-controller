# Model Explainability Job Definition Sample

This sample demonstrates how to start model-explainability job definitions using your own model-explainability script, packaged in a SageMaker-compatible container, using the Amazon AWS Controllers for Kubernetes (ACK) service controller for Amazon SageMaker.                     

## Prerequisites

This sample assumes that you have completed the [common prerequisites](/samples/README.md).

You will also need an [Endpoint](/samples/endpoint/README.md) configured in SageMaker and you will need to run a baselining job to generate baseline constraints only.

### Get an Image

All SageMaker model-explainability job definitions are run from within a container with all necessary dependencies and modules pre-installed and with the model-explainability scripts referencing the acceptable input and output directories. Sample container images are [available](https://docs.aws.amazon.com/sagemaker/latest/dg/sagemaker-algo-docker-registry-paths.html).

A container image URL and tag looks has the following structure:
```
<account number>.dkr.ecr.<region>.amazonaws.com/<image name>:<tag>
```

### Updating the Model-Explainability Job Definition Specification

In the `my-model-explainability-job-definition.yaml` file, modify the placeholder values with those associated with your account and model-explainability job definition.

## Submitting your Model-Explainability Job Definition

### Create a Model-Explainability Job Definiton

To submit your prepared model-explainability job definition specification, apply the specification to your Kubernetes cluster as such:
```
$ kubectl apply -f my-model-explainability-job-definition.yaml
modelexplainabilityjobdefinitions.sagemaker.services.k8s.aws.amazon.com/my-model-explainability-job-definition created
```

### List Model-Explainability Job Definitons

To list all Model-Explainablity Job Definition created using the ACK controller use the following command:
```
$ kubectl get modelexplainabilityjobdefinitions
```

### Describe a Model-Explainability Job Definiton

To get more details about the Model-Explainability Job Definition once it's submitted, like checking the status, errors or parameters of the Model-Explainability Job Definition use the following command:
```
$ kubectl describe modelexplainabilityjobdefinitions my-model-explainability-job-definition
```
You can also check Status.ackResourceMetadata.Arn to verify the data quality job definition was created successfully.

### Delete a Model-Explainability Job Definiton

To delete the model-explainability job definition, use the following command:
```
$ kubectl delete modelexplainabilityjobdefinitions my-model-explainability-job-definition
```