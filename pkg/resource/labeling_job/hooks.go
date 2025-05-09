// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package labeling_job

import (
	"errors"

	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
	svccommon "github.com/aws-controllers-k8s/sagemaker-controller/pkg/common"
	svcsdktypes "github.com/aws/aws-sdk-go-v2/service/sagemaker/types"
)

var (
	modifyingStatuses = []string{
		string(svcsdktypes.LabelingJobStatusInProgress),
		string(svcsdktypes.LabelingJobStatusStopping),
	}
	resourceName = GroupKind.Kind

	requeueWaitWhileDeleting = ackrequeue.NeededAfter(
		errors.New(resourceName+" is Stopping."),
		ackrequeue.DefaultRequeueAfterDuration,
	)
)

// customSetOutput sets ConditionTypeResourceSynced condition to True or False
// based on the labelingJobStatus on AWS so the reconciler can determine if a
// requeue is needed
func (rm *resourceManager) customSetOutput(r *resource) {
	jobStatus := r.ko.Status.LabelingJobStatus
	svccommon.SetSyncedCondition(r, jobStatus, &resourceName, &modifyingStatuses)
}
