# Endpoint Sample

This sample demonstrates how to create Endpoints using your own Endpoint_base/config script, packaged in a SageMaker-compatible container, using the Amazon AWS Controllers for Kubernetes (ACK) service controller for Amazon SageMaker.   

## Prerequisites

This sample assumes that you have completed the [common prerequisites](/samples/README.md).

You will also need a model in SageMaker for this sample. If you do not have one you must first create a [model](/samples/model/README.md)

In order to run [endpoint_base](/samples/endpoint/endpoint_base.yaml) you will need an endpoint_config which can be created by [endpoint_config](/samples/endpoint/endpoint_config.yaml)

### Updating the Endpoint Specification

In the `endpoint_config.yaml` file, modify the placeholder values with those associated with your account. The `spec.productionVariants.ModelName` should be the SageMaker model from the previous step.  

## Submitting your Endpoint Specification

### Create an Endpoint Config and Endpoint

To submit your prepared endpoint specification, apply the specification to your Kubernetes cluster as such:
```
$ kubectl apply -f my-endpoint.yaml
endpoints.sagemaker.services.k8s.aws.amazon.com/my-endpoint created
```
If it is a endpoint config:
```
$ kubectl apply -f my-endpoint-config.yaml
endpointsconfigs.sagemaker.services.k8s.aws /my-endpoint-config created
```

### List Endpoint Configs and Endpoints

To list all Endpoints created using the ACK controller use the following command:
```
$ kubectl get endpoints.sagemaker.services.k8s.aws
```
If it is a endpoint config it is endpointsconfigs.sagemaker.services.k8s.aws  
```
$ kubectl get endpointsconfigs.sagemaker.services.k8s.aws
```

### Describe an Endpoint Config and Endpoint

To get more details about the Endpoint once it's submitted, like checking the status, errors or parameters of the Endpoint use the following command:
```
$ kubectl describe endpoints.sagemaker.services.k8s.aws my-endpoint
```

If it is a endpoint config it is endpointsconfigs.sagemaker.services.k8s.aws  
```
$ kubectl describe endpointsconfigs.sagemaker.services.k8s.aws my-endpoint-config
```

### Delete an Endpoint Config and Endpoint

To delete the Endpoint, use the following command:
```
$ kubectl delete endpoints.sagemaker.services.k8s.aws my-endpoint
```

If it is a endpoint config it is endpointsconfigs.sagemaker.services.k8s.aws  
```
$ kubectl delete endpointsconfigs.sagemaker.services.k8s.aws  my-endpoint-config
```