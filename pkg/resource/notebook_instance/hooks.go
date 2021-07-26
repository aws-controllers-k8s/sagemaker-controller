package notebook_instance

import (
	"context"

	svccommon "github.com/aws-controllers-k8s/sagemaker-controller/pkg/common"

	svcsdk "github.com/aws/aws-sdk-go/service/sagemaker"
)

var (
	modifyingStatuses = []string{svcsdk.NotebookInstanceStatusPending, svcsdk.NotebookInstanceStatusUpdating,
		svcsdk.NotebookInstanceStatusDeleting, svcsdk.NotebookInstanceStatusStopping}
	resourceName = resourceGK.Kind
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
