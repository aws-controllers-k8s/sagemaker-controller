# Batch Transform Job Sample

This sample demonstrates how to start batch transform jobs using your own batch-transform script, packaged in a SageMaker-compatible container, using the Amazon AWS Controllers for Kubernetes (ACK) service controller for Amazon SageMaker.                     

## Prerequisites

This sample assumes that you have completed the [common prerequisites](/samples/README.md).

You will also need a model in SageMaker for this sample. If you do not have one you must first create a [model](/samples/model/README.md).

### Updating the Batch Transform Job Specification

In the `my-batch-transform-job.yaml` file, modify the placeholder values with those associated with your account and batchtransform job. 

## Submitting your Batch Transform Job

### Create a Batch Transform Job

To submit your prepared batch transform job specification, apply the specification to your Kubernetes cluster as such:
```
$ kubectl apply -f my-batch-transform-job.yaml
batch-transformjob.sagemaker.services.k8s.aws.amazon.com/my-batch-transform-job created
```

### List Batch Transform Jobs

To list all Batch Transform Jobs created using the ACK controller use the following command:
```
$ kubectl get batch-transformjob
```

### Describe a Batch Transform Job

To get more details about the Batch Transform Job once it's submitted, like checking the status, errors or parameters of the Batch Transform Job use the following command:
```
$ kubectl describe batch-transformjob my-batch-transform-job
```

### Delete a Batch Transform Job
To delete the Batch Transform Job, use the following command:
```
$ kubectl delete batch-transformjob my-batch-transform-job
```