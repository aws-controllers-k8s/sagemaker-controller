# Notebook Instance Sample

## Prerequisites

This sample assumes that you have completed the [common prerequisites](/samples/README.md).

### Update the Notebook Instance Specification

Edit the roleARN value in my-notebook-instance.yaml to include the Sagemaker Execution permissions.

## Using the Notebook Instance operator

### Create a Notebook Instance

This command creates a Sagemaker notebook instance based on the specification provided in my-notebook-instance.yaml.
The Notebook Instance will start at the Pending state and will transition into InService once ready.

Note: The lifecycleConfigName is the name of the NotebookInstance Lifecycle configuration, it is created beforehand.

```
$ kubectl apply -f my-notebook-instance.yaml
```

### Lists Notebook Instances
This command lists all the notebook instances created using the ACK controller.
```
$ kubectl get NotebookInstance
```

### Describe a Notebook Instance
This command desribes a specific Notebook Instance, it is useful for checking items like the status, errors or parameters of the Notebook Instance.

Note: The status field NotebookInstanceURL returns a url in the form <name>.notebook.<region>.sagemaker.aws. To view the NotebookInstance in the browser, use https://<url>

```
$ kubectl describe NotebookInstance my-notebook
```

### Update a Notebook Instance
This commands updates the Notebook Instance with the updated spec provided in my-notebook-instance.yaml. The update command retains the state the controller was previously in. If the update command was called while the controller was in the InService state, it will end up in the InService state after the update. If the update command was called while the controller was in the Stopped/Stopping state, it will end up in the Stopped state after the update.
```
$ kubectl apply -f my-notebook-instance.yaml
```


### Delete a Notebook Instance
This command deletes the Notebook Instance.
```
$ kubectl delete NotebookInstance my-notebook
```

