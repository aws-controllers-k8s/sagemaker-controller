# ACK service controller for Amazon SageMaker

This repository contains source code for the AWS Controllers for Kubernetes
(ACK) service controller for Amazon SageMaker.

Please [log issues][ack-issues] and feedback on the main AWS Controllers for
Kubernetes Github project.

[ack-issues]: https://github.com/aws/aws-controllers-k8s/issues

## Contributing

We welcome community contributions and pull requests.

See our [contribution guide](/CONTRIBUTING.md) for more information on how to
report issues, set up a development environment, and submit code.

We adhere to the [Amazon Open Source Code of Conduct][coc].

You can also learn more about our [Governance](/GOVERNANCE.md) structure.

[coc]: https://aws.github.io/code-of-conduct

## License

This project is [licensed](/LICENSE) under the Apache-2.0 License.

## Supported SageMaker Resources
For a list of supported resources, refer to the [SageMaker API Reference](https://aws-controllers-k8s.github.io/community/reference/).

Find the helm charts and controller images on Amazon ECR Public Gallery.
- Helm Chart: https://gallery.ecr.aws/aws-controllers-k8s/sagemaker-chart

- Controller Image: https://gallery.ecr.aws/aws-controllers-k8s/sagemaker-controller

## Getting Started

### 1.0 ACK SageMaker Controller

For a step-by-step tutorial head over to [Machine Learning with the ACK SageMaker Controller](https://aws-controllers-k8s.github.io/community/docs/tutorials/sagemaker-example/)

### 2.0 ACK Application Auto Scaling Controller

For a step-by-step tutorial on how to use Application Auto Scaling Controller with SageMaker head over to [Scale SageMaker Workloads with Application Auto Scaling](https://aws-controllers-k8s.github.io/community/docs/tutorials/autoscaling-example/)

### 3.0 Samples 

#### 3.1 SageMaker samples

Head over to the [samples directory](/samples) and follow the README to create resources. 

#### 3.2 Application-autoscaling samples

Head over to the [samples directory in application-autoscaling controller repository](https://github.com/aws-controllers-k8s/applicationautoscaling-controller/tree/main/samples/hosting-autoscaling-on-sagemaker) and follow the README to create resources. 

### 4.0 Cross Region Resource Management

Head over to the [Manage Resources In Multiple Regions](https://aws-controllers-k8s.github.io/community/docs/user-docs/multi-region-resource-management/)

### 5.0 Cross Account Resource Management

Head over to the [Manage Resources In Multiple AWS Accounts](https://aws-controllers-k8s.github.io/community/docs/user-docs/cross-account-resource-management/)

### 6.0 Adopt Resources

ACK controller provides to provide the ability to “adopt” resources that were not originally created by ACK service controller. Given the user configures the controller with permissions which has access to existing resource, the controller will be able to determine the current specification and status of the AWS resource and reconcile said resource as if the ACK controller had originally created it.

Sample:
```yaml
apiVersion: services.k8s.aws/v1alpha1
kind: AdoptedResource
metadata:
  name: adopt-endpoint-sample
spec:  
  aws:
    # resource to adopt, not created by ACK
    nameOrID: xgboost-endpoint
  kubernetes:
    group: sagemaker.services.k8s.aws
    kind: Endpoint
    metadata:
      # target K8s CR name
      name: xgboost-endpoint
```
Save the above to a file name adopt-endpoint-sample.yaml.

Submit the CR
```sh
kubectl apply -f adopt-endpoint-sample.yaml
```

Check for `ACK.Adopted` condition to be true under `status.conditions`

```sh
kubectl describe adoptedresource adopt-endpoint-sample
```

Output should look similar to this:
```yaml
---
kind: AdoptedResource
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: '{"apiVersion":"services.k8s.aws/v1alpha1","kind":"AdoptedResource","metadata":{"annotations":{},"name":"xgboost-endpoint","namespace":"default"},"spec":{"aws":{"nameOrID":"xgboost-endpoint"},"kubernetes":{"group":"sagemaker.services.k8s.aws","kind":"Endpoint","metadata":{"name":"xgboost-endpoint"}}}}'
  creationTimestamp: '2021-04-27T02:49:14Z'
  finalizers:
  - finalizers.services.k8s.aws/AdoptedResource
  generation: 1
  name: adopt-endpoint-sample
  namespace: default
  resourceVersion: '12669876'
  selfLink: "/apis/services.k8s.aws/v1alpha1/namespaces/default/adoptedresources/adopt-endpoint-sample"
  uid: 35f8fa92-29dd-4040-9d0d-0b07bbd7ca0b
spec:
  aws:
    nameOrID: xgboost-endpoint
  kubernetes:
    group: sagemaker.services.k8s.aws
    kind: Endpoint
    metadata:
      name: xgboost-endpoint
status:
  conditions:
  - status: 'True'
    type: ACK.Adopted
```

Check resource exists in cluster
```sh
kubectl describe endpoints.sagemaker xgboost-endpoint
```


