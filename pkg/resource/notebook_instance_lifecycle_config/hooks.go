package notebook_instance_lifecycle_config

import (
	"errors"

	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
)

var (
	resourceName = GroupKind.Kind

	requeueWaitWhileUpdating = ackrequeue.NeededAfter(
		errors.New(resourceName+" is updating."),
		ackrequeue.DefaultRequeueAfterDuration,
	)
)
