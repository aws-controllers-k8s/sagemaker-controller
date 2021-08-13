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

// Note: modifyingStatuses are all the statuses where the NotebookInstance requeues.
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
	notebookStatus := r.ko.Status.NotebookInstanceStatus
	return svccommon.RequeueIfModifying(notebookStatus, &resourceName, &modifyingStatuses)
}

// customSetOutput sets the ack syncedCondition depending on
// whether the latest status of the resource is one of the
// defined modifyingStatuses.
func (rm *resourceManager) customSetOutput(r *resource) {
	latestStatus := r.ko.Status.NotebookInstanceStatus
	svccommon.SetSyncedCondition(r, latestStatus, &resourceName, &modifyingStatuses)
}

//CustomSetOutputDescribe
func (rm *resourceManager) customSetOutputDescribe(r *resource,
	ko *svcapitypes.NotebookInstance) error {
	notebook_state := *ko.Status.NotebookInstanceStatus // Get the Notebook State
	if ko.Status.IsUpdating != nil && *ko.Status.IsUpdating == "true" {
		if notebook_state != svcsdk.NotebookInstanceStatusStopped {
			return nil //we want to keep requeing until update finishes
		}
		//TODO: Use annotations instead of status once the runtime supports updating metadata
		if ko.Status.StoppedByController != nil && *ko.Status.StoppedByController == "true" {
			err := rm.startNotebookInstance(r)
			if err != nil {
				return err
			}
			ko.Status.IsUpdating = aws.String("false")
			ko.Status.StoppedByController = aws.String("false")
			return nil //Indicating that StoppedByAck and IsUpdating got updated
		}
		ko.Status.IsUpdating = aws.String("false")
		return nil

	}
	rm.customSetOutput(&resource{ko}) // We set the sync status here
	return nil
}

// customSetOutputUpdate sets the isUpdating field and also sets the resource condition to false
// so that the controller requeues.
func (rm *resourceManager) customSetOutputUpdate(ko *svcapitypes.NotebookInstance) {
	//TODO: Replace the IsUpdating status with an annotation if the runtime can update annotations after a readOne call.
	ko.Status.IsUpdating = aws.String("true")
	ackcond.SetSynced(&resource{ko}, corev1.ConditionFalse, nil, aws.String("Notebook is currenty updating"))
}
