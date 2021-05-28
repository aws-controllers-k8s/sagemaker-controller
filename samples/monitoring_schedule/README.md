# Monitoring Schedule Sample

This sample demonstrates how to start Monitoring scheduling using your own Monitor scheduling script, packaged in a SageMaker-compatible container, using the Amazon AWS Controllers for Kubernetes (ACK) service controller for Amazon SageMaker.   

## Prerequisites

This sample assumes that you have completed the [common prerequisites](/samples/README.md). It also assumes you have a job definition of one of these [types](/samples/job_definitions).

### Updating the Monitoring Specification

In the `my-monitoring-schedule.yaml` file, modify the placeholder values with those associated with your account and Monitoring schedule. The `spec.monitoringScheduleConfig.monitoringType` can be any of these [types](/samples/job_definitions).

## Submitting your Monitoring Schedule

### Create a Monitoring Schedule
To submit your prepared monitoring schedule specification, apply the specification to your Kubernetes cluster as such:
```
$ kubectl apply -f my-monitoring-schedule.yaml
monitorings.sagemaker.services.k8s.aws.amazon.com/my-monitoring-schedule created
```

### List Monitoring Schedules
To list all Monitoring Schedules created using the ACK controller use the following command:
```
$ kubectl get monitorings.sagemaker.services.k8s.aws
```

### Describe a Monitoring Schedule
To get more details about the Monitoring Schedule once it's submitted, like checking the status, errors or parameters of the Monitoring Schedule use the following command:
```
$ kubectl describe monitorings.sagemaker.services.k8s.aws my-monitoring-schedule
```

### Delete a Monitoring Schedule
To delete the Monitoring Schedule, use the following command:
```
$ kubectl delete monitorings.sagemaker.services.k8s.aws my-monitoring-schedule
```