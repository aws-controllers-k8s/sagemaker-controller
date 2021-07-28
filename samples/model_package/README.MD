# Model Package Sample

This sample demonstrates how to create a Model Package using your own Model Package script, packaged in a SageMaker-compatible container, using the Amazon AWS Controllers for Kubernetes (ACK) service controller for Amazon SageMaker.   

## Prerequisites

This sample assumes that you have completed the [common prerequisites](/samples/README.md).

### Updating the Model Package Specification

In the `my-model-package-yaml` file, modify the placeholder values with those associated with your account, S3 path and Model Package  

## Submitting your Model Package 
### Create a Model Package To submit your prepared Model Package specification, apply the specification to your Kubernetes cluster as such:
```
$ kubectl apply -f my-model-package-yaml
modelpackages.sagemaker.services.k8s.aws.amazon.com/my-model-package created
```

### List Model Package 
To list all Model Package  created using the ACK controller use the following command:
```
$ kubectl get modelpackages
```

### Describe a Model Package To get more details about the Model Package once it's submitted, like checking the status, errors or parameters of the Model Package use the following command:
```
$ kubectl describe modelpackages my-model-package
```

### Delete a Model Package To delete the Model Package  use the following command:
```
$ kubectl delete modelpackages my-model-package```