package notebook_instance

import (
	ackcond "github.com/aws-controllers-k8s/runtime/pkg/condition"
	svcapitypes "github.com/aws-controllers-k8s/sagemaker-controller/apis/v1alpha1"
	svccommon "github.com/aws-controllers-k8s/sagemaker-controller/pkg/common"
	corev1 "k8s.io/api/core/v1"
)

// customSetOutput sets the ack syncedCondition depending on
// whether the latest status of the resource is one of the
// defined modifyingStatuses.
func (rm *resourceManager) customSetOutput(r *resource) {
	latestStatus := r.ko.Status.NotebookInstanceStatus
	svccommon.SetSyncedCondition(r, latestStatus, &resourceName, &modifyingStatuses)
}

//Create does
func (rm *resourceManager) customSetOutputCreate(
	notebookInstanceStatus *string, ko *svcapitypes.NotebookInstance) {
	if notebookInstanceStatus == nil {
		return
	}
	pendingReason := "Notebook is currenty starting"
	ackcond.SetSynced(&resource{ko}, corev1.ConditionFalse, nil, &pendingReason)

}

//Update uses this function
func (rm *resourceManager) customSetOutputDescribe(r *resource, ko *svcapitypes.NotebookInstance) {
	rm.customSetOutput(&resource{ko}) // We set the sync status here

}
