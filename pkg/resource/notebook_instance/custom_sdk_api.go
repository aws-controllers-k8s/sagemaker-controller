package notebook_instance

import (
	svcsdk "github.com/aws/aws-sdk-go/service/sagemaker"
)

// Note: This file is used for functions calling Sagemaker APIs that are non-CRUD.

// stopNotebookInstance stops the Notebook Instance and returns
// an error if there was an error with stopping the Notebook Instance.
func (rm *resourceManager) stopNotebookInstance(r *resource) error {
	input := &svcsdk.StopNotebookInstanceInput{}
	input.SetNotebookInstanceName(*r.ko.Spec.NotebookInstanceName)

	// Stop Notebook Instance does not return a response
	_, err := rm.sdkapi.StopNotebookInstance(input)
	rm.metrics.RecordAPICall("STOP", "StopNotebookInstance", err)
	return err
}

// startNotebookInstance starts the Notebook Instance and returns an
// error if there was an error with starting the Notebook Instance.
func (rm *resourceManager) startNotebookInstance(r *resource) error {
	input := &svcsdk.StartNotebookInstanceInput{}
	input.NotebookInstanceName = r.ko.Spec.NotebookInstanceName

	//Start Notebook Instance does not return a response
	_, err := rm.sdkapi.StartNotebookInstance(input)
	rm.metrics.RecordAPICall("START", "StartNotebookInstance", err)
	return err
}
