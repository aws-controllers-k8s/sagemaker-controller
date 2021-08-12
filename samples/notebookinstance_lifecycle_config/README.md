# Notebook Instance Lifecycle Configuration sample

## Prerequisites

This sample assumes that you have completed the [common prerequisites](/samples/README.md).

### Create a Notebook Instance Lifecycle Configuration
This command creates a Notebook Instance Lifecycle configuration.
Note: The onCreate and onStart fields must be in base64.

In the example crd the value of the onCreate field decodes to: `echo "Entering onCreate"` and the value of the onStart field decodes to: `echo "Entering onStart"`
```
$ kubectl apply -f my-notebookinstance_lifecycle_config.yaml
```
### Update a Notebook Instance Lifecycle Configuration
This command updates a Notebook Instance Lifecycle configuration.
Note: The onCreate and onStart fields must be in base64.
```
$ kubectl apply -f my-notebookinstance_lifecycle_config.yaml
```
### Describe a Notebook Instance Lifecycle Configuration
This command describes a Notebook Instance Lifecycle Configuration.
```
$ kubectl describe NotebookInstanceLifecycleConfig <YOUR LIFECYCLE CONFIG NAME>
```
### List Notebook Instance Lifecycle Configurations
This command lists Notebook Instance Lifecycle Configurations and their last modified times.
```
$ kubectl get NotebookInstanceLifecycleConfig 
```
### Delete a Notebook Instance Lifecycle Configuration
This command deletes a Notebook Instance Lifecycle Configuration.
```
$ kubectl delete NotebookInstanceLifecycleConfig <YOUR LIFECYCLE CONFIG NAME>
```