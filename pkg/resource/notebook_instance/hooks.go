package notebook_instance

import (
	"context"
	"errors"
	"time"

	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
	svccommon "github.com/aws-controllers-k8s/sagemaker-controller/pkg/common"
	svcsdk "github.com/aws/aws-sdk-go/service/sagemaker"
)

var (
	requeueWaitWhileStopping = ackrequeue.NeededAfter(
		errors.New("NotebookInstance in 'Stopping' state, cannot be modified or deleted"),
		20*time.Second,
	)
	modifyingStatuses = []string{svcsdk.NotebookInstanceStatusPending, svcsdk.NotebookInstanceStatusUpdating,
		svcsdk.NotebookInstanceStatusDeleting, svcsdk.NotebookInstanceStatusStopping}
	resourceName             = resourceGK.Kind
	requeueWaitWhileDeleting = ackrequeue.NeededAfter(
		errors.New(resourceName+" is deleting."),
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
