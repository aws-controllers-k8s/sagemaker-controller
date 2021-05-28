# Hyperparameter Tuning Job Sample

This sample demonstrates how to start hyperparameter jobs using your own hyperparameter script, packaged in a SageMaker-compatible container, using the Amazon AWS Controllers for Kubernetes (ACK) service controller for Amazon SageMaker.                     

## Prerequisites

This sample assumes that you have completed the [common prerequisites](/samples/README.md).

### Upload S3 Data

You will need training data uploaded to an S3 bucket. Make sure you have AWS credentials and and have the bucket in the same region where you plan to create SageMaker resources. Run the following python script to upload sample data to your S3 bucket.
```
python3 ../training/s3_sample_data.py $S3_BUCKET_NAME
```

### Get an Image

All SageMaker Hyperparameter jobs are run from within a container with all necessary dependencies and modules pre-installed and with the hyperparameter scripts referencing the acceptable input and output directories. Sample container images are [available](https://docs.aws.amazon.com/sagemaker/latest/dg/sagemaker-algo-docker-registry-paths.html).

A container image URL and tag looks has the following structure:
```
<account number>.dkr.ecr.<region>.amazonaws.com/<image name>:<tag>
```

### Updating the Hyperparameter Specification

In the `my-hyperparameter-job.yaml` file, modify the placeholder values with those associated with your account and hyperparameter job.

### Enabling Spot Training

In the `my-hyperparameter-job.yaml` file under `spec.trainingJobDefinition` add `enableManagedSpotTraining` and set the value to true. You will also need to specify a `spec.trainingJobDefinition.stoppingCondition.maxRuntimeInSeconds` and `spec.trainingJobDefinition.stoppingCondition.maxWaittimeInSeconds`

## Submitting your Hyperparameter Job

### Create a Hyperparameter Job

To submit your prepared hyperparameter job specification, apply the specification to your Kubernetes cluster as such:
```
$ kubectl apply -f my-hyperparameter-job.yaml
hyperparametertuningjob.sagemaker.services.k8s.aws.amazon.com/my-hyperparameter-job created
```

### List Hyperparameter Jobs

To list all Hyperparameter jobs created using the ACK controller use the following command:
```
$ kubectl get hyperparametertuningjob
```

### Describe a Hyperparameter Job

To get more details about the Hyperparameter job once it's submitted, like checking the status, errors or parameters of the Hyperparameter job use the following command:
```
$ kubectl describe hyperparametertuningjob my-hyperparameter-job
```

### Delete a Hyperparameter Job

To delete the hyperparameter job, use the following command:
```
$ kubectl delete hyperparametertuningjob my-hyperparameter-job
```