package notebook_instance

import (
	svcsdk "github.com/aws/aws-sdk-go/service/sagemaker"
)

// Note: This file is used for functions calling Sagemaker APIs that are non-CRUD.

// This function stops the Notebook Instance.
func (rm *resourceManager) stopNotebookInstance(r *resource) error {
	input := &svcsdk.StopNotebookInstanceInput{}
	input.SetNotebookInstanceName(*r.ko.Spec.NotebookInstanceName)

	// Stop Notebook Instance does not return a response
	_, err := rm.sdkapi.StopNotebookInstance(input)
	rm.metrics.RecordAPICall("STOP", "StopNotebookInstance", err)
	return err
}
