# Training Job Sample

This sample demonstrates how to start training jobs using your own training script, packaged in a SageMaker-compatible container, using the Amazon AWS Controllers for Kubernetes (ACK) service controller for Amazon SageMaker.                     

## Prerequisites

This sample assumes that you have completed the [common prerequisites](/samples/README.md).

###  Upload S3 Data

You will need training data uploaded to an S3 bucket. Make sure you have AWS credentials and and have the bucket in the same region where you plan to create SageMaker resources. Run the following python script to upload sample data to your S3 bucket.
```
python3 s3_sample_data.py $S3_BUCKET_NAME
```

### Get an Image

All SageMaker training jobs are run from within a container with all necessary dependencies and modules pre-installed and with the training scripts referencing the acceptable input and output directories. Sample container images are [available](https://docs.aws.amazon.com/sagemaker/latest/dg/sagemaker-algo-docker-registry-paths.html).

A container image URL and tag looks has the following structure:
```
<account number>.dkr.ecr.<region>.amazonaws.com/<image name>:<tag>
```

### Updating the Training Specification

In the `my-training-job.yaml` file, modify the placeholder values with those associated with your account and training job.

## Submitting your Training Job

### Create a Training Job
To submit your prepared training job specification, apply the specification to your Kubernetes cluster as such:
```
$ kubectl apply -f my-training-job.yaml
trainingjob.sagemaker.services.k8s.aws.amazon.com/my-training-job created
```

### List Training Jobs
To list all training jobs created using the ACK controller use the following command:
```
$ kubectl get trainingjob
```

### Describe a Training Job
To get more details about the training job once it's submitted, like checking the status, errors or parameters of the training job use the following command:
```
$ kubectl describe trainingjob my-training-job
```

### Delete a Training Job
To delete the training job, use the following command:
```
$ kubectl delete trainingjob my-training-job
```