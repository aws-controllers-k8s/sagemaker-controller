package notebook_instance

import (
	"context"
	"errors"

	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
	svcapitypes "github.com/aws-controllers-k8s/sagemaker-controller/apis/v1alpha1"
	svccommon "github.com/aws-controllers-k8s/sagemaker-controller/pkg/common"
	"github.com/aws/aws-sdk-go/aws"
	svcsdk "github.com/aws/aws-sdk-go/service/sagemaker"
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

// CustomSetOutputDescribe has three functions
// 1. Setting the resource synced condition during sdk.Find
// 2. Starting the Notebook if the StoppedByController field is set to UpdateTriggerd
// 3. Reseting the StoppedByControllerMetadata field if the Notebook starts.
func (rm *resourceManager) customSetOutputDescribe(r *resource) error {
	// Get the Notebook State
	notebookStatus := *r.ko.Status.NotebookInstanceStatus
	if notebookStatus == svcsdk.NotebookInstanceStatusStopped {
		if r.ko.Status.StoppedByControllerMetadata != nil && *r.ko.Status.StoppedByControllerMetadata == "UpdateTriggered" {
			err := rm.startNotebookInstance(r)
			if err != nil {
				// It's best to not update StoppedByControllerMetadata here because the controller can always try again.
				return err
			}
			r.ko.Status.StoppedByControllerMetadata = nil
			// Notebook after rm.startNotebookInstance(r) is in pending state
			// Manually change status since otherwise synced condition will be set to True as
			// of a result Stopped status.
			r.ko.Status.NotebookInstanceStatus = aws.String(svcsdk.NotebookInstanceStatusPending)
		}
	}
	// The resource synced status is set here.
	rm.customSetOutput(r)
	return nil
}

// customSetOutputUpdate sets the StoppedByControllerMeta field to UpdateTriggered if it was in UpdatePending before.
func (rm *resourceManager) customSetOutputUpdate(ko *svcapitypes.NotebookInstance, latest *resource) {
	//TODO: Replace the StoppedByControllerMeta status with an annotation if the runtime can update annotations after a readOne call.
	if latest.ko.Status.StoppedByControllerMetadata != nil && *latest.ko.Status.StoppedByControllerMetadata == "UpdatePending" {
		ko.Status.StoppedByControllerMetadata = aws.String("UpdateTriggered")
	}
	svccommon.SetSyncedCondition(&resource{ko}, aws.String(svcsdk.NotebookInstanceStatusUpdating), &resourceName, &modifyingStatuses)
}
