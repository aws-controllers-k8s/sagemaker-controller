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

// CustomSetOutputDescribe has two functions
// 1. Setting the resource synced condition during sdk.Find
// 2. Starting the Notebook if the StoppedByController field is set to UpdateTriggerd
// 3. Reseting the StoppedByController field.
func (rm *resourceManager) customSetOutputDescribe(r *resource,
	ko *svcapitypes.NotebookInstance) error {
	notebookStatus := *ko.Status.NotebookInstanceStatus // Get the Notebook State
	if notebookStatus == svcsdk.NotebookInstanceStatusStopped {
		if ko.Status.StoppedByControllerMETA != nil && *ko.Status.StoppedByControllerMETA == "UpdateTriggered" {
			err := rm.startNotebookInstance(r)
			if err != nil {
				//We dont update StoppedByController here because we want it to try to start it again on the next reqeue.
				return err
			}
			ko.Status.StoppedByControllerMETA = nil
		}
	}
	rm.customSetOutput(&resource{ko}) // We set the sync status here
	return nil
}

// customSetOutputUpdate sets the StoppedByController field to UpdateTriggered if it was in UpdatePending before.
func (rm *resourceManager) customSetOutputUpdate(ko *svcapitypes.NotebookInstance, latest *resource) {
	//TODO: Replace the IsUpdating status with an annotation if the runtime can update annotations after a readOne call.
	if latest.ko.Status.StoppedByControllerMETA != nil && *latest.ko.Status.StoppedByControllerMETA == "UpdatePending" {
		ko.Status.StoppedByControllerMETA = aws.String("UpdateTriggered")
	}
	ackcond.SetSynced(&resource{ko}, corev1.ConditionFalse, nil, aws.String("Notebook is currenty updating"))
}
