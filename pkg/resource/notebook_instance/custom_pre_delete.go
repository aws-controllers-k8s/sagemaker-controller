package notebook_instance

import (
	svcsdk "github.com/aws/aws-sdk-go/service/sagemaker"
)

/*
This code stops the NotebookInstance right before its about to be deleted/updated.
Returns True if the Notebook was stopped
*/
func (rm *resourceManager) customStopNotebook(
	r *resource) bool {
	latestStatus := *r.ko.Status.NotebookInstanceStatus
	if &latestStatus == nil {
		return false
	}

	//We only want to stop the Notebook if its not already stopped/stopping or not in a failed state.
	if rm.isInServiceStatus(latestStatus) {
		err := rm.stopNotebookInstance(r)
		if err == nil {
			return true
		}
	}
	return false
}

func (rm *resourceManager) isInServiceStatus(latestStatus string) bool {
	return latestStatus == svcsdk.NotebookInstanceStatusInService
}
