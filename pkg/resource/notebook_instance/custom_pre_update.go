package notebook_instance

import (
	"context"

	svcsdk "github.com/aws/aws-sdk-go/service/sagemaker"
)

/*
This function stops the notebook instance(if its running) before the update build request.
TODO: Set an annotation here or in sdk.go once runtime can update annotations while returning an error.
*/
func (rm *resourceManager) customPreUpdate(
	ctx context.Context,
	desired *resource,
	latest *resource,
) bool {

	latestStatus := *latest.ko.Status.NotebookInstanceStatus
	if &latestStatus == nil {
		return false
	}

	if latestStatus == svcsdk.NotebookInstanceStatusInService {
		nb_input := svcsdk.StopNotebookInstanceInput{}
		nb_input.NotebookInstanceName = &desired.ko.Name
		rm.sdkapi.StopNotebookInstance(&nb_input)
		return true
	}
	return false
}
