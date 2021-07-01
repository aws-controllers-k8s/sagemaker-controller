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

package feature_group

import (
	"errors"
	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
	svcsdk "github.com/aws/aws-sdk-go/service/sagemaker"
)

var (
	requeueWaitWhileCreating = ackrequeue.NeededAfter(
		errors.New("Still Creating, wait until we reach Created."),
		ackrequeue.DefaultRequeueAfterDuration,
	)
)

// isCreating returns true if supplied replication group resource state is 'creating'
func isCreating(r *resource) bool {
	if r == nil || r.ko.Status.FeatureGroupStatus == nil {
		return false
	}
	status := *r.ko.Status.FeatureGroupStatus
	return status == svcsdk.FeatureGroupStatusCreating
}
