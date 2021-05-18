## Prerequisites

This sample assumes that you have already configured an EKS cluster with the ACK operator. It also assumes that you have installed `kubectl` - you can find a link on our [installation page](TODO). It also assumes you have a job of one of these types already created `DataQuality|ModelQuality|ModelBias|ModelExplainability`.


### Updating the monitoring Specification

In the `my-monitoring-schedule.yaml` file, modify the placeholder values with those associated with your account and hyperparameter schedule. The `spec.monitoringScheduleConfig.monitoringType` can be any of these options `DataQuality|ModelQuality|ModelBias|ModelExplainability`.  

## Submitting your monitoring schedule

To submit your prepared monitoring schedule specification, apply the specification to your EKS cluster as such:
```
$ kubectl apply -f my-monitoring-schedule.yaml
monitorings.sagemaker.services.k8s.aws.amazon.com/my-monitoring-schedule created
```

To monitor the monitoring schedule status, you can use the following command:
```
$ kubectl get monitorings.sagemaker.services.k8s.aws  my-monitoring-schedule
```

To monitor the monitoring schedule in-depth once it has started, you can see the full status and any additional errors with the following command:
```
$ kubectl describe monitorings.sagemaker.services.k8s.aws   my-monitoring-schedule
```

To delete the monitoring schedule once it has started if errors occurred or for any reason with the following command:
```
$ kubectl delete monitorings.sagemaker.services.k8s.aws my-monitoring-schedule
```
