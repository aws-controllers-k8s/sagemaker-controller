# Processing Job Sample

This sample demonstrates how to start processing jobs using your own processing script, packaged in a SageMaker-compatible container, using the Amazon AWS Controllers for Kubernetes (ACK) service controller for Amazon SageMaker.                     

## Prerequisites

This sample assumes that you have completed the [common prerequisites](/samples/README.md).

You will need to upload [kmeans_preprocessing.py](/samples/processing/kmeans_preprocessing.py) to an S3 bucket and update the s3Input s3URI path.

### Get an Image

All SageMaker processing jobs are run from within a container with all necessary dependencies and modules pre-installed and with the processing job scripts referencing the acceptable input and output directories. Sample container images are [available](https://github.com/aws/deep-learning-containers/blob/master/available_images.md).

A container image URL and tag looks has the following structure:
```
<account number>.dkr.ecr.<region>.amazonaws.com/<image name>:<tag>
```

### Updating the Processing Job Specification

In the `my-processing-job.yaml` file, modify the placeholder values with those associated with your account and processing job. 

## Submitting your Processing Job

### Create a Processing Job

To submit your prepared processing job specification, apply the specification to your Kubernetes cluster as such:
```
$ kubectl apply -f my-processing-job.yaml
processingjob.sagemaker.services.k8s.aws.amazon.com/my-processing-job created
```

### List Processing Jobs
To list all processing jobs created using the ACK controller use the following command:
```
$ kubectl get processingjob
```

### Describe a Processing Job
To get more details about the processing job once it's submitted, like checking the status, errors or parameters of the processing job use the following command:
```
$ kubectl describe processingjob my-processing-job
```

### Delete a Processing Job
To delete the processing job, use the following command:
```
$ kubectl delete processingjob my-processing-job
```