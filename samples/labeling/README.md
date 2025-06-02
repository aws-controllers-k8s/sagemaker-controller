# Labeling Job Sample

This sample demonstrates how to start processing jobs using your own labeling artifacts, packaged in a SageMaker-compatible container, using the Amazon AWS Controllers for Kubernetes (ACK) service controller for Amazon SageMaker.

## Prerequisites

This sample assumes that you have completed the [common prerequisites](/samples/README.md).

You will need to upload [images/*.jpg](/samples/labeling/images/), [instructions.template](/samples/labeling/instrucitons.template), and [class_labels.json](/samples/labeling/class_labels.json) to an S3 bucket and update the s3Input s3URI path.

You will also need to replace `<YOUR S3 PATH>` within [input.manifest](/samples/labeling/input.manifest) and upload the file to the same S3 bucket.

Make sure that your S3 bucket contains necessary CORS permission.
```
# Define the configuration rules
cors_configuration = {
    'CORSRules': [{
        'AllowedHeaders': [],
        'AllowedMethods': ['GET'],
        'AllowedOrigins': ['*'],
        'ExposeHeaders': []
    }]
}


# Set the CORS configuration
s3 = boto3.client("s3")
s3.put_bucket_cors(Bucket=<BUCKET_NAME>,
                   CORSConfiguration=cors_configuration)
```

### Updating the Processing Job Specification

In the `my-labeling-job.yaml` file, modify the placeholder values with those associated with your account and labeling job.

## Submitting your Processing Job

### Create a Processing Job

To submit your prepared processing job specification, apply the specification to your Kubernetes cluster as such:
```
$ kubectl apply -f my-labeling-job.yaml
processingjob.sagemaker.services.k8s.aws.amazon.com/my-labeling-job created
```

### List Processing Jobs
To list all processing jobs created using the ACK controller use the following command:
```
$ kubectl get labelingjobs
```

### Describe a Processing Job
To get more details about the processing job once it's submitted, like checking the status, errors or parameters of the processing job use the following command:
```
$ kubectl describe labelingjobs my-labeling-job
```

### Delete a Processing Job
To delete the processing job, use the following command:
```
$ kubectl delete labelingjobs my-labeling-job
```
