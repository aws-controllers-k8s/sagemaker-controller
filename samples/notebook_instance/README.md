# Notebook Instance Sample

## Prerequisites

This sample assumes that you have completed the [common prerequisites](/samples/README.md).

### Update the Notebook Instance Specification

Edit the roleARN value in my-notebook-instance.yaml to include the Sagemaker Execution permissions.

## Using the Notebook Instance controller

### Create a Notebook Instance

This command creates a Sagemaker notebook instance based on the specification provided in my-notebook-instance.yaml.
The Notebook Instance will start at the Pending state and will transition into InService once ready.

```
$ kubectl apply -f my-notebook-instance.yaml
notebookinstance.sagemaker.services.k8s.aws/my-notebook-instance created
```
### Update a Notebook Instance
This commands updates the Notebook Instance with the updated spec provided in my-notebook-instance.yaml. The update command retains the state the Notebook Instance was previously in. If the update command was called while the Notebook Instance was in the InService state, it will end up in the InService state after the update. If the update command was called while the Notebook Instance was in the Stopped/Stopping state, it will end up in the Stopped state after the update.
Additionally this controller automatically takes care of setting the `Disassociate<field>` fields if the corresponding `<field>` field is removed from the spec and resource is updated. `<field>` in this case means the following fields: LifecycleConfigName, DefaultCodeRepository, AdditionalCodeRepositories, AcceleratorTypes.
```
$ kubectl apply -f my-notebook-instance.yaml
notebookinstance.sagemaker.services.k8s.aws/my-notebook-instance configured
```
### Describe a Notebook Instance
This command desribes a specific Notebook Instance, it is useful for checking items like the status, errors or parameters of the Notebook Instance.

Note: The status field Url returns a url in the form <name>.notebook.<region>.sagemaker.aws. To view the Jupyter Notebook in the browser, use https://<url>

```
$ kubectl describe NotebookInstance <YOUR NOTEBOOK INSTANCE NAME>
```
### List Notebook Instances
This command lists all the notebook instances created using the ACK controller.
```
$ kubectl get NotebookInstance
```
### Delete a Notebook Instance
This command deletes the Notebook Instance.
```
$ kubectl delete NotebookInstance <YOUR NOTEBOOK INSTANCE NAME>
```