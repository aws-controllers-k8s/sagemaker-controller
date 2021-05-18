## Prerequisites

This sample assumes that you have already configured an EKS cluster with the ACK operator. It also assumes that you have installed `kubectl` - you can find a link on our [installation page](TODO).

You will also need a model in SageMaker for this sample. If you do not have one you must first create a [model](https://docs.aws.amazon.com/sagemaker/latest/dg/sagemaker-mkt-model-pkg-model.html)

In order to run `endpoint_base.yaml` you will need an endpoint_config which can be created by the `endpoint_config.yaml`.



### Updating the Endpoint Specification

In the `endpoint_config.yaml` file, modify the placeholder values with those associated with your account and hyperparameter job. The `spec.algorithmSpecification.modelName` should be the SageMaker model from the previous step.  

## Submitting your endpoint Job


To submit your prepared endpoint job specification, apply the specification to your EKS cluster as such:
```
$ kubectl apply -f my-endpoint-job.yaml
endpoints.sagemaker.services.k8s.aws.amazon.com/my-endpoint-job created
```



To monitor the endpoint job status, you can use the following command:
```
$ kubectl get endpoints.sagemaker.services.k8s.aws   my-endpoint-job
```
If it is a endpoint config job it is endpointsconfigs.sagemaker.services.k8s.aws  
```
$ kubectl get endpointsconfigs.sagemaker.services.k8s.aws  my-endpoint-config-job
```

To monitor the endpoint job in-depth once it has started, you can see the full status and any additional errors with the following command:
```
$ kubectl describe endpoints.sagemaker.services.k8s.aws   my-endpoint-job
```

If it is a endpoint config job it is endpointsconfigs.sagemaker.services.k8s.aws  
```
$ kubectl describe endpointsconfigs.sagemaker.services.k8s.aws  my-endpoint-config-job
```

To delete the endpoint job once it has started if errors occurred or for any reason with the following command:
```
$ kubectl delete endpoints.sagemaker.services.k8s.aws my-endpoint-job
```

If it is a endpoint config job it is endpointsconfigs.sagemaker.services.k8s.aws  
```
$ kubectl delete endpointsconfigs.sagemaker.services.k8s.aws  my-endpoint-config-job
```