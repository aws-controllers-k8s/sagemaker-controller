# Model Package Group Sample

This sample demonstrates how to create a Model Package Group using your own Model Package Group script, packaged in a SageMaker-compatible container, using the Amazon AWS Controllers for Kubernetes (ACK) service controller for Amazon SageMaker.   

## Prerequisites

This sample assumes that you have completed the [common prerequisites](/samples/README.md).

### Updating the Model Package Group Specification

In the `my-model-package-group.yaml` file, modify the placeholder values with those associated with your account and Model Package Group. 

## Submitting your Model Package Group

### Create a Model Package Group
To submit your prepared Model Package Group specification, apply the specification to your Kubernetes cluster as such:
```
$ kubectl apply -f my-model-package-group.yaml
modelpackagegroups.sagemaker.services.k8s.aws.amazon.com/my-model-package-group created
```

### List Model Package Groups
To list all Model Package Groups created using the ACK controller use the following command:
```
$ kubectl get modelpackagegroups
```

### Describe a Model Package Group
To get more details about the Model Package Group once it's submitted, like checking the status, errors or parameters of the Model Package Group use the following command:
```
$ kubectl describe modelpackagegroups my-model-package-group
```

### Delete a Model Package Group
To delete the Model Package Group, use the following command:
```
$ kubectl delete modelpackagegroups my-model-package-group
```