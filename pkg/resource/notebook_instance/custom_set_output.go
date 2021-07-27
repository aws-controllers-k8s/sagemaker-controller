package notebook_instance

import (
	ackcond "github.com/aws-controllers-k8s/runtime/pkg/condition"
	svcapitypes "github.com/aws-controllers-k8s/sagemaker-controller/apis/v1alpha1"
	svccommon "github.com/aws-controllers-k8s/sagemaker-controller/pkg/common"
	"github.com/aws/aws-sdk-go/aws"
	svcsdk "github.com/aws/aws-sdk-go/service/sagemaker"
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
func (rm *resourceManager) customSetOutputCreateUpdate(ko *svcapitypes.NotebookInstance) {
	if ko == nil {
		return
	}
	ackcond.SetSynced(&resource{ko}, corev1.ConditionFalse, nil, aws.String("Notebook is currenty starting or updating"))
}

func (rm *resourceManager) customSetOutputDescribe(r *resource,
	ko *svcapitypes.NotebookInstance) bool {

	notebook_state := *ko.Status.NotebookInstanceStatus // Get the Notebook State
	if ko.Annotations != nil && ko.Annotations["done_updating"] == "true" {
		if notebook_state != svcsdk.NotebookInstanceStatusStopped {
			return false //we want to keep requeing until update finishes
		}
		if ko.Status.StoppedByAck != nil && *ko.Status.StoppedByAck == "true" {
			nb_input := svcsdk.StartNotebookInstanceInput{}
			nb_input.NotebookInstanceName = &r.ko.Name
			rm.sdkapi.StartNotebookInstance(&nb_input)
			ko.Annotations["done_updating"] = "false"
			resetStoppedbyAck := "false"
			ko.Status.StoppedByAck = &resetStoppedbyAck
			return true //dont want to set resource synced to true pre maturely.
		}
	}
	rm.customSetOutput(&resource{ko}) // We set the sync status here
	return false
}
