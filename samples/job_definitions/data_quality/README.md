# Data Quality Job Definition Sample

This sample demonstrates how to start data quality job definitions using your own data-quality-job-definitions script, packaged in a SageMaker-compatible container, using the Amazon AWS Controllers for Kubernetes (ACK) service controller for Amazon SageMaker.                     

## Prerequisites

This sample assumes that you have completed the [common prerequisites](/samples/README.md).

You will need an [Endpoint](/samples/endpoint/README.md) configured in SageMaker and you will need to run a baselining job to generate baseline statistics and constraints.

### Get an Image

All SageMaker data quality job definitions are run from within a container with all necessary dependencies and modules pre-installed and with the data-quality scripts referencing the acceptable input and output directories. Sample container images are [available](https://docs.aws.amazon.com/sagemaker/latest/dg/sagemaker-algo-docker-registry-paths.html).

A container image URL and tag looks has the following structure:
```
<account number>.dkr.ecr.<region>.amazonaws.com/<image name>:<tag>
```

### Updating the Data Quality Job Definition Specification

In the `my-data-quality-job-definition.yaml` file, modify the placeholder values with those associated with your account.

## Submitting your Data Quality Job Definition

### Create a Data Quality Job Definition 

To submit your prepared data quality job definition specification, apply the specification to your Kubernetes cluster as such:
```
$ kubectl apply -f my-data-quality-job-definition.yaml
dataqualityjobdefinitions.sagemaker.services.k8s.aws.amazon.com/my-data-quality-job-definition created
```

### List Data Quality Job Definitions

To monitor the data quality job definition status, you can use the following command:
```
$ kubectl get dataqualityjobdefinitions
```

### Describe a Data Quality Job Definition

To get more details about the Data Quality Job Definition once it's submitted, like checking the status, errors or parameters of the Data Quality Job Definition use the following command:
```
$ kubectl describe dataqualityjobdefinitions my-data-quality-job-definition
```
You can also check Status.ackResourceMetadata.Arn to verify the data quality job definition was created successfully.

### Delete a Data Quality Job Definition

To delete the data quality job definition, use the following command:
```
$ kubectl delete dataqualityjobdefinitions my-data-quality-job-definition
```