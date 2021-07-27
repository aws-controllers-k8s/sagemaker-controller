package notebook_instance

import (
	ackcond "github.com/aws-controllers-k8s/runtime/pkg/condition"
	svcapitypes "github.com/aws-controllers-k8s/sagemaker-controller/apis/v1alpha1"
	svccommon "github.com/aws-controllers-k8s/sagemaker-controller/pkg/common"
	"github.com/aws/aws-sdk-go/aws"
	corev1 "k8s.io/api/core/v1"
)

// customSetOutput sets the ack syncedCondition depending on
// whether the latest status of the resource is one of the
// defined modifyingStatuses.
func (rm *resourceManager) customSetOutput(r *resource) {
	latestStatus := r.ko.Status.NotebookInstanceStatus
	svccommon.SetSyncedCondition(r, latestStatus, &resourceName, &modifyingStatuses)
}

//The resource from create does not have a state in the status field
func (rm *resourceManager) customSetOutputCreate(ko *svcapitypes.NotebookInstance) {
	if ko == nil {
		return
	}
	ackcond.SetSynced(&resource{ko}, corev1.ConditionFalse, nil, aws.String("Notebook is currenty starting"))
}

//Update uses this function
func (rm *resourceManager) customSetOutputDescribe(r *resource, ko *svcapitypes.NotebookInstance) {
	rm.customSetOutput(&resource{ko}) // We set the sync status here
}
