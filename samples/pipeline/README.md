# Pipeline Sample

This sample demonstrates how to submit a pipeline to Sagemaker for execution using your own JSON pipeline definition, using the AWS Controllers for Kubernetes (ACK) service controller for Amazon SageMaker.   

## Prerequisites

This sample assumes that you have completed the [common prerequisites](/samples/README.md).

### Updating the Pipeline Specification

In the `pipeline.yaml` file, modify the placeholder values with those associated with your account. 

## Submitting your Pipeline Specification

### Modify/Create a JSON pipeline definition

Create the JSON pipeline definition using the JSON schema documented at https://aws-sagemaker-mlops.github.io/sagemaker-model-building-pipeline-definition-JSON-schema/. You will need to convert your JSON pipeline definition into String format to pass to the *.spec.pipelineDefinition* key in the Kubernetes YAML spec. You may use online third-party tools to convert from JSON to String format. In this sample, you are provided a sample pipeline definition with one Training step.

### Submit pipeline to Sagemaker and start an execution

To submit your prepared pipeline specification, apply the specification to your Kubernetes cluster as such:
```
$ kubectl apply -f my-pipeline.yaml
pipeline.sagemaker.services.k8s.aws/my-kubernetes-pipeline created
```
To start an execution run of the pipeline:
```
$ kubectl apply -f pipeline-execution.yaml
pipelineexecution.sagemaker.services.k8s.aws/my-kubernetes-pipeline-execution created
```

### List pipelines and pipeline executions

To list all pipelines created using the ACK controller use the following command:
```
$ kubectl get pipeline.sagemaker.services.k8s.aws
```
If it is a pipeline executions it is endpointsconfigs.sagemaker.services.k8s.aws  
```
$ kubectl get pipelineexecution.sagemaker.services.k8s.aws
```

### Describe a pipeline and pipeline execution

To get more details about the pipeline once it's submitted, like checking the status, errors or parameters of the pipeline, use the following command:
```
$ kubectl describe pipeline.sagemaker.services.k8s.aws my-kubernetes-pipeline
```

If it is a endpoint config it is endpointsconfigs.sagemaker.services.k8s.aws  
```
$ kubectl describe pipelineexecution.sagemaker.services.k8s.aws my-kubernetes-pipeline-execution
```

### Delete a pipeline and a pipeline execution

To delete the pipeline, use the following command:
```
$ kubectl delete pipeline.sagemaker.services.k8s.aws my-kubernetes-pipeline
```

If it is a endpoint config it is endpointsconfigs.sagemaker.services.k8s.aws  
```
$ kubectl delete pipelineexecution.sagemaker.services.k8s.aws my-kubernetes-pipeline-execution
```