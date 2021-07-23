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

package common

import (
	"errors"
	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
)

// RequeueIfModifying creates and returns an
// ackrequeue if a resource's latest status matches
// one of the provided modifying statuses.
func RequeueIfModifying(
	latestStatus *string,
	resourceName *string,
	modifyingStatuses *[]string,
) error {
	if latestStatus == nil || !IsModifyingStatus(latestStatus, modifyingStatuses) {
		return nil
	}

	errMsg := *resourceName + " in " + *latestStatus + " state cannot be modified or deleted."
	requeueWaitWhileModifying := ackrequeue.NeededAfter(
		errors.New(errMsg),
		ackrequeue.DefaultRequeueAfterDuration,
	)
	return requeueWaitWhileModifying
}
