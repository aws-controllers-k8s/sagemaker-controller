package notebook_instance

import (
	"errors"
	"time"

	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
	svcsdk "github.com/aws/aws-sdk-go/service/sagemaker"
)

var (
	requeueWaitWhileStopping = ackrequeue.NeededAfter(
		errors.New("NotebookInstance in 'Stopping' state, currently transitioning to 'Stopped' so Notebook can be modified or deleted"),
		10*time.Second,
	)
	requeueWaitWhilePending = ackrequeue.NeededAfter(
		errors.New("NotebookInstance in 'Pending' state, cannot be modified or deleted"),
		10*time.Second,
	)
	requeueWaitWhileDeleting = ackrequeue.NeededAfter(
		errors.New("NotebookInstance in 'Deleting' state, cannot be modified or deleted"),
		10*time.Second,
	)
	requeueWaitWhileUpdating = ackrequeue.NeededAfter(
		errors.New("NotebookInstance in 'Updating' state, cannot be modified or deleted"),
		20*time.Second,
	)
)

func isNotebookStopping(r *resource) bool {
	if r.ko.Status.NotebookInstanceStatus == nil {
		return false
	}
	notebookInstanceStatus := r.ko.Status.NotebookInstanceStatus

	return *notebookInstanceStatus == svcsdk.NotebookInstanceStatusStopping
}

func isNotebookPending(r *resource) bool {
	if r.ko.Status.NotebookInstanceStatus == nil {
		return false
	}
	notebookInstanceStatus := r.ko.Status.NotebookInstanceStatus

	return *notebookInstanceStatus == svcsdk.NotebookInstanceStatusPending
}

func isNotebookDeleting(r *resource) bool {
	if r.ko.Status.NotebookInstanceStatus == nil {
		return false
	}
	notebookInstanceStatus := r.ko.Status.NotebookInstanceStatus

	return *notebookInstanceStatus == svcsdk.NotebookInstanceStatusDeleting
}
func isNotebookUpdating(r *resource) bool {
	if r.ko.Status.NotebookInstanceStatus == nil {
		return false
	}
	notebookInstanceStatus := r.ko.Status.NotebookInstanceStatus

	return *notebookInstanceStatus == svcsdk.NotebookInstanceStatusUpdating
}
