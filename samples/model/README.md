# Model Sample

This sample demonstrates how to start models/create a model using your own model script, packaged in a SageMaker-compatible container, using the Amazon AWS Controllers for Kubernetes (ACK) service controller for Amazon SageMaker.                     

## Prerequisites

This sample assumes that you have completed the [common prerequisites](/samples/README.md).

### Upload S3 Data

You will need a model uploaded to an S3 bucket. Make sure you have AWS credentials and and have the bucket in the same region where you plan to create SageMaker resources. Run the following python script to upload sample data to your S3 bucket.
```
python3 s3_sample_model.py $S3_BUCKET_NAME
```

### Get an Image

All SageMaker models are run from within a container with all necessary dependencies and modules pre-installed and with the model scripts referencing the acceptable input and output directories. Sample container images are [available](https://docs.aws.amazon.com/sagemaker/latest/dg/sagemaker-algo-docker-registry-paths.html).

A container image URL and tag looks has the following structure:
```
<account number>.dkr.ecr.<region>.amazonaws.com/<image name>:<tag>
```

### Updating the Model Specification

In the `my-model.yaml` file, modify the placeholder values with those associated with your account and model.

## Submitting your Model

### Create your Model
To submit your prepared model specification, apply the specification to your Kubernetes cluster as such:
```
$ kubectl apply -f my-model.yaml
models.sagemaker.services.k8s.aws.amazon.com/my-model created
```

### List Models
To list all Models created using the ACK controller use the following command:
```
$ kubectl get models.sagemaker.services.k8s.aws
```

### Describe your Model
To get more details about the Model once it's submitted, like checking the status, errors or parameters of the Model use the following command:
```
$ kubectl describe models.sagemaker.services.k8s.aws my-model
```

### Delete your Model
To delete the model, use the following command:
```
$ kubectl delete models.sagemaker.services.k8s.aws my-model
```