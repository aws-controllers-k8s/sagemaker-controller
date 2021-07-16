# Feature Group Sample

This sample demonstrates how to create a feature group using the Amazon AWS Controllers for Kubernetes (ACK) service controller for Amazon SageMaker.

Inspiration for this sample was taken from the notebook on [Fraud Detection with Amazon SageMaker FeatureStore](https://sagemaker-examples.readthedocs.io/en/latest/sagemaker-featurestore/sagemaker_featurestore_fraud_detection_python_sdk.html).

## Prerequisites

This sample assumes that you have completed the [common prerequisites](https://github.com/aws-controllers-k8s/sagemaker-controller/blob/main/samples/README.md).

### Create an S3 bucket:

Since we are using the offline store in this example, you need to set up an s3 bucket. [Here are directions](https://docs.aws.amazon.com/AmazonS3/latest/userguide/create-bucket-overview.html) to set up your s3 bucket through the S3 Console, AWS SDK, or AWS CLI.

### Updating the Feature Group Specification:

In the `my-feature-group.yaml` file, modify the placeholder values with those associated with your account and feature group.

## Creating your Feature Group

### Create a Feature Group:

To submit your prepared feature group specification, apply the specification to your Kubernetes cluster as such:

```
$ kubectl apply -f my-feature-group.yaml
featuregroup.sagemaker.services.k8s.aws/my-feature-group created
```

### List Feature Groups:

To list all feature groups created using the ACK controller use the following command:

```
$ kubectl get featuregroup
```

### Describe a Feature Group:

To get more details about the feature group once it's submitted, like checking the status, errors or parameters of the feature group use the following command:

```
$ kubectl describe featuregroup my-feature-group
```

## Ingesting Data into your Feature Group

Note that ingestion is **not** supported in the controller.
To ingest data from the my-sample-data.csv file into your feature group, use the following command:

```
$ python3 data_ingestion.py -i my-sample-data.csv -fg my-feature-group
```

## Deleting your Feature Group

To delete the feature group, use the following command:

```
$ kubectl delete featuregroup my-feature-group
featuregroup.sagemaker.services.k8s.aws "my-feature-group" deleted
```
