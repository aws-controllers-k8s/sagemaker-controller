package notebook_instance

import (
	ackcond "github.com/aws-controllers-k8s/runtime/pkg/condition"
	svcapitypes "github.com/aws-controllers-k8s/sagemaker-controller/apis/v1alpha1"
	svcsdk "github.com/aws/aws-sdk-go/service/sagemaker"
	corev1 "k8s.io/api/core/v1"
)

func (rm *resourceManager) customSetOutput(
	notebookInstanceStatus *string, ko *svcapitypes.NotebookInstance) {

	if notebookInstanceStatus == nil {
		return
	}
	pendingReason := "Notebook is currenty starting"
	if *notebookInstanceStatus == svcsdk.NotebookInstanceStatusDeleting || *notebookInstanceStatus == svcsdk.NotebookInstanceStatusFailed ||
		*notebookInstanceStatus == svcsdk.NotebookInstanceStatusInService || *notebookInstanceStatus == svcsdk.NotebookInstanceStatusStopped {
		ackcond.SetSynced(&resource{ko}, corev1.ConditionTrue, nil, nil)

	} else if *notebookInstanceStatus == svcsdk.NotebookInstanceStatusPending {
		ackcond.SetSynced(&resource{ko}, corev1.ConditionFalse, nil, &pendingReason)
	} else {
		ackcond.SetSynced(&resource{ko}, corev1.ConditionFalse, nil, nil)
	}
}

func (rm *resourceManager) customSetOutputDescribe(r *resource,
	ko *svcapitypes.NotebookInstance) {

	notebook_state := *ko.Status.NotebookInstanceStatus // Get the Notebook State
	if ko.Annotations != nil && ko.Annotations["done_updating"] == "true" {
		if notebook_state != svcsdk.NotebookInstanceStatusStopped {
			return //we want to keep requeing until update finishes
		}
		if ko.Status.StoppedByAck != nil && *ko.Status.StoppedByAck == "true" {
			nb_input := svcsdk.StartNotebookInstanceInput{}
			nb_input.NotebookInstanceName = &r.ko.Name
			rm.sdkapi.StartNotebookInstance(&nb_input)
			ko.Annotations["done_updating"] = "false"
			resetStoppedbyAck := "false"
			ko.Status.StoppedByAck = &resetStoppedbyAck
			return //dont want to set resource synced to true pre maturely.
		}
	}
	rm.customSetOutput(&notebook_state, ko) // We set the sync status here
}
