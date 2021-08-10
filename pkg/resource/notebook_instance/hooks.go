package notebook_instance

import (
	"context"
	"errors"

	ackcond "github.com/aws-controllers-k8s/runtime/pkg/condition"
	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
	svcapitypes "github.com/aws-controllers-k8s/sagemaker-controller/apis/v1alpha1"
	svccommon "github.com/aws-controllers-k8s/sagemaker-controller/pkg/common"
	"github.com/aws/aws-sdk-go/aws"
	svcsdk "github.com/aws/aws-sdk-go/service/sagemaker"
	corev1 "k8s.io/api/core/v1"
)

var (
	modifyingStatuses = []string{
		svcsdk.NotebookInstanceStatusPending,
		svcsdk.NotebookInstanceStatusUpdating,
		svcsdk.NotebookInstanceStatusDeleting,
		svcsdk.NotebookInstanceStatusStopping,
	}

	resourceName = resourceGK.Kind

	requeueWaitWhileDeleting = ackrequeue.NeededAfter(
		errors.New(resourceName+" is deleting."),
		ackrequeue.DefaultRequeueAfterDuration,
	)
	requeueWaitWhileStopping = ackrequeue.NeededAfter(
		errors.New(resourceName+" is stopping."),
		ackrequeue.DefaultRequeueAfterDuration,
	)
)

// requeueUntilCanModify creates and returns an
// ackrequeue error if a resource's latest status matches
// any of the defined modifying statuses below.
func (rm *resourceManager) requeueUntilCanModify(
	ctx context.Context,
	r *resource,
) error {
	latestStatus := r.ko.Status.NotebookInstanceStatus
	return svccommon.RequeueIfModifying(latestStatus, &resourceName, &modifyingStatuses)
}

func isNotebookUpdating(r *resource) bool {
	if r.ko.Status.NotebookInstanceStatus == nil {
		return false
	}
	notebookInstanceStatus := r.ko.Status.NotebookInstanceStatus
	return *notebookInstanceStatus == svcsdk.NotebookInstanceStatusUpdating
}

// customSetOutput sets the ack syncedCondition depending on
// whether the latest status of the resource is one of the
// defined modifyingStatuses.
func (rm *resourceManager) customSetOutput(r *resource) {
	latestStatus := r.ko.Status.NotebookInstanceStatus
	svccommon.SetSyncedCondition(r, latestStatus, &resourceName, &modifyingStatuses)
}

func (rm *resourceManager) customSetOutputDescribe(r *resource,
	ko *svcapitypes.NotebookInstance) bool {
	notebook_state := *ko.Status.NotebookInstanceStatus // Get the Notebook State
	if ko.Annotations != nil && ko.Annotations["done_updating"] == "true" {
		if notebook_state != svcsdk.NotebookInstanceStatusStopped {
			return false //we want to keep requeing until update finishes
		}
		if ko.Status.StoppedByAck != nil && *ko.Status.StoppedByAck == "true" {
			rm.startNotebookInstance(r)
			ko.Annotations["done_updating"] = "false"
			resetStoppedbyAck := "false"
			ko.Status.StoppedByAck = &resetStoppedbyAck
			return true //dont want to set resource synced to true pre maturely.
		}
	}
	rm.customSetOutput(&resource{ko}) // We set the sync status here
	return false
}

//The resource from create does not have a state in the status field
func (rm *resourceManager) customSetOutputUpdate(ko *svcapitypes.NotebookInstance) {
	if ko == nil {
		return
	}
	ackcond.SetSynced(&resource{ko}, corev1.ConditionFalse, nil, aws.String("Notebook is currenty updating"))
}
