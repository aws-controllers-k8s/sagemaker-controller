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

package pipeline_execution

import (
	"errors"
	"strings"

	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
	svcapitypes "github.com/aws-controllers-k8s/sagemaker-controller/apis/v1alpha1"
	svccommon "github.com/aws-controllers-k8s/sagemaker-controller/pkg/common"
	svcsdk "github.com/aws/aws-sdk-go-v2/service/sagemaker"
	svcsdktypes "github.com/aws/aws-sdk-go-v2/service/sagemaker/types"
)

var (
	modifyingStatuses = []string{string(svcsdktypes.PipelineExecutionStatusExecuting),
		string(svcsdktypes.PipelineExecutionStatusStopping)}

	resourceName = GroupKind.Kind

	requeueWaitWhileDeleting = ackrequeue.NeededAfter(
		errors.New(resourceName+" is Stopping."),
		ackrequeue.DefaultRequeueAfterDuration,
	)
)

// customSetOutput sets ConditionTypeResourceSynced condition to True or False
// based on the pipelineExecutionStatus on AWS so the reconciler can determine if a
// requeue is needed
func (rm *resourceManager) customSetOutput(r *resource) {
	latestStatus := r.ko.Status.PipelineExecutionStatus
	svccommon.SetSyncedCondition(r, latestStatus, &resourceName, &modifyingStatuses)
}

func (rm *resourceManager) customSetSpec(ko *svcapitypes.PipelineExecution, resp *svcsdk.DescribePipelineExecutionOutput) {
	if ko.Spec.PipelineName == nil {
		var name = strings.SplitAfter(*resp.PipelineArn, "/")
		ko.Spec.PipelineName = &name[0]
	}
}
