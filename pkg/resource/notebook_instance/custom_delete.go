package notebook_instance

import (
	svcsdk "github.com/aws/aws-sdk-go/service/sagemaker"
)

/*
This code stops the NotebookInstance right before its about to be deleted.
*/
func (rm *resourceManager) customPreDelete(
	r *resource) {

	latestStatus := *r.ko.Status.NotebookInstanceStatus

	if &latestStatus == nil {
		return
	}

	//We only want to stop the Notebook if its not already stopped/stopping or not in a failed state.
	if latestStatus != svcsdk.NotebookInstanceStatusStopped && latestStatus != svcsdk.NotebookInstanceStatusFailed &&
		latestStatus != svcsdk.NotebookInstanceStatusStopping {
		nb_input := svcsdk.StopNotebookInstanceInput{}
		nb_input.NotebookInstanceName = &r.ko.Name
		rm.sdkapi.StopNotebookInstance(&nb_input)
	}
}
