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

// Code generated by ack-generate. DO NOT EDIT.

package user_profile

import (
	"bytes"
	"reflect"

	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	acktags "github.com/aws-controllers-k8s/runtime/pkg/tags"
)

// Hack to avoid import errors during build...
var (
	_ = &bytes.Buffer{}
	_ = &reflect.Method{}
	_ = &acktags.Tags{}
)

// newResourceDelta returns a new `ackcompare.Delta` used to compare two
// resources
func newResourceDelta(
	a *resource,
	b *resource,
) *ackcompare.Delta {
	delta := ackcompare.NewDelta()
	if (a == nil && b != nil) ||
		(a != nil && b == nil) {
		delta.Add("", a, b)
		return delta
	}

	if ackcompare.HasNilDifference(a.ko.Spec.DomainID, b.ko.Spec.DomainID) {
		delta.Add("Spec.DomainID", a.ko.Spec.DomainID, b.ko.Spec.DomainID)
	} else if a.ko.Spec.DomainID != nil && b.ko.Spec.DomainID != nil {
		if *a.ko.Spec.DomainID != *b.ko.Spec.DomainID {
			delta.Add("Spec.DomainID", a.ko.Spec.DomainID, b.ko.Spec.DomainID)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.SingleSignOnUserIdentifier, b.ko.Spec.SingleSignOnUserIdentifier) {
		delta.Add("Spec.SingleSignOnUserIdentifier", a.ko.Spec.SingleSignOnUserIdentifier, b.ko.Spec.SingleSignOnUserIdentifier)
	} else if a.ko.Spec.SingleSignOnUserIdentifier != nil && b.ko.Spec.SingleSignOnUserIdentifier != nil {
		if *a.ko.Spec.SingleSignOnUserIdentifier != *b.ko.Spec.SingleSignOnUserIdentifier {
			delta.Add("Spec.SingleSignOnUserIdentifier", a.ko.Spec.SingleSignOnUserIdentifier, b.ko.Spec.SingleSignOnUserIdentifier)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.SingleSignOnUserValue, b.ko.Spec.SingleSignOnUserValue) {
		delta.Add("Spec.SingleSignOnUserValue", a.ko.Spec.SingleSignOnUserValue, b.ko.Spec.SingleSignOnUserValue)
	} else if a.ko.Spec.SingleSignOnUserValue != nil && b.ko.Spec.SingleSignOnUserValue != nil {
		if *a.ko.Spec.SingleSignOnUserValue != *b.ko.Spec.SingleSignOnUserValue {
			delta.Add("Spec.SingleSignOnUserValue", a.ko.Spec.SingleSignOnUserValue, b.ko.Spec.SingleSignOnUserValue)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.UserProfileName, b.ko.Spec.UserProfileName) {
		delta.Add("Spec.UserProfileName", a.ko.Spec.UserProfileName, b.ko.Spec.UserProfileName)
	} else if a.ko.Spec.UserProfileName != nil && b.ko.Spec.UserProfileName != nil {
		if *a.ko.Spec.UserProfileName != *b.ko.Spec.UserProfileName {
			delta.Add("Spec.UserProfileName", a.ko.Spec.UserProfileName, b.ko.Spec.UserProfileName)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.UserSettings, b.ko.Spec.UserSettings) {
		delta.Add("Spec.UserSettings", a.ko.Spec.UserSettings, b.ko.Spec.UserSettings)
	} else if a.ko.Spec.UserSettings != nil && b.ko.Spec.UserSettings != nil {
		if ackcompare.HasNilDifference(a.ko.Spec.UserSettings.ExecutionRole, b.ko.Spec.UserSettings.ExecutionRole) {
			delta.Add("Spec.UserSettings.ExecutionRole", a.ko.Spec.UserSettings.ExecutionRole, b.ko.Spec.UserSettings.ExecutionRole)
		} else if a.ko.Spec.UserSettings.ExecutionRole != nil && b.ko.Spec.UserSettings.ExecutionRole != nil {
			if *a.ko.Spec.UserSettings.ExecutionRole != *b.ko.Spec.UserSettings.ExecutionRole {
				delta.Add("Spec.UserSettings.ExecutionRole", a.ko.Spec.UserSettings.ExecutionRole, b.ko.Spec.UserSettings.ExecutionRole)
			}
		}
		if ackcompare.HasNilDifference(a.ko.Spec.UserSettings.JupyterServerAppSettings, b.ko.Spec.UserSettings.JupyterServerAppSettings) {
			delta.Add("Spec.UserSettings.JupyterServerAppSettings", a.ko.Spec.UserSettings.JupyterServerAppSettings, b.ko.Spec.UserSettings.JupyterServerAppSettings)
		} else if a.ko.Spec.UserSettings.JupyterServerAppSettings != nil && b.ko.Spec.UserSettings.JupyterServerAppSettings != nil {
			if ackcompare.HasNilDifference(a.ko.Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec, b.ko.Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec) {
				delta.Add("Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec", a.ko.Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec, b.ko.Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec)
			} else if a.ko.Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec != nil && b.ko.Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec != nil {
				if ackcompare.HasNilDifference(a.ko.Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.InstanceType, b.ko.Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.InstanceType) {
					delta.Add("Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.InstanceType", a.ko.Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.InstanceType, b.ko.Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.InstanceType)
				} else if a.ko.Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.InstanceType != nil && b.ko.Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.InstanceType != nil {
					if *a.ko.Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.InstanceType != *b.ko.Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.InstanceType {
						delta.Add("Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.InstanceType", a.ko.Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.InstanceType, b.ko.Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.InstanceType)
					}
				}
				if ackcompare.HasNilDifference(a.ko.Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.LifecycleConfigARN, b.ko.Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.LifecycleConfigARN) {
					delta.Add("Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.LifecycleConfigARN", a.ko.Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.LifecycleConfigARN, b.ko.Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.LifecycleConfigARN)
				} else if a.ko.Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.LifecycleConfigARN != nil && b.ko.Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.LifecycleConfigARN != nil {
					if *a.ko.Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.LifecycleConfigARN != *b.ko.Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.LifecycleConfigARN {
						delta.Add("Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.LifecycleConfigARN", a.ko.Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.LifecycleConfigARN, b.ko.Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.LifecycleConfigARN)
					}
				}
				if ackcompare.HasNilDifference(a.ko.Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.SageMakerImageARN, b.ko.Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.SageMakerImageARN) {
					delta.Add("Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.SageMakerImageARN", a.ko.Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.SageMakerImageARN, b.ko.Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.SageMakerImageARN)
				} else if a.ko.Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.SageMakerImageARN != nil && b.ko.Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.SageMakerImageARN != nil {
					if *a.ko.Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.SageMakerImageARN != *b.ko.Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.SageMakerImageARN {
						delta.Add("Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.SageMakerImageARN", a.ko.Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.SageMakerImageARN, b.ko.Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.SageMakerImageARN)
					}
				}
				if ackcompare.HasNilDifference(a.ko.Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.SageMakerImageVersionARN, b.ko.Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.SageMakerImageVersionARN) {
					delta.Add("Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.SageMakerImageVersionARN", a.ko.Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.SageMakerImageVersionARN, b.ko.Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.SageMakerImageVersionARN)
				} else if a.ko.Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.SageMakerImageVersionARN != nil && b.ko.Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.SageMakerImageVersionARN != nil {
					if *a.ko.Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.SageMakerImageVersionARN != *b.ko.Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.SageMakerImageVersionARN {
						delta.Add("Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.SageMakerImageVersionARN", a.ko.Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.SageMakerImageVersionARN, b.ko.Spec.UserSettings.JupyterServerAppSettings.DefaultResourceSpec.SageMakerImageVersionARN)
					}
				}
			}
			if !ackcompare.SliceStringPEqual(a.ko.Spec.UserSettings.JupyterServerAppSettings.LifecycleConfigARNs, b.ko.Spec.UserSettings.JupyterServerAppSettings.LifecycleConfigARNs) {
				delta.Add("Spec.UserSettings.JupyterServerAppSettings.LifecycleConfigARNs", a.ko.Spec.UserSettings.JupyterServerAppSettings.LifecycleConfigARNs, b.ko.Spec.UserSettings.JupyterServerAppSettings.LifecycleConfigARNs)
			}
		}
		if ackcompare.HasNilDifference(a.ko.Spec.UserSettings.KernelGatewayAppSettings, b.ko.Spec.UserSettings.KernelGatewayAppSettings) {
			delta.Add("Spec.UserSettings.KernelGatewayAppSettings", a.ko.Spec.UserSettings.KernelGatewayAppSettings, b.ko.Spec.UserSettings.KernelGatewayAppSettings)
		} else if a.ko.Spec.UserSettings.KernelGatewayAppSettings != nil && b.ko.Spec.UserSettings.KernelGatewayAppSettings != nil {
			if !reflect.DeepEqual(a.ko.Spec.UserSettings.KernelGatewayAppSettings.CustomImages, b.ko.Spec.UserSettings.KernelGatewayAppSettings.CustomImages) {
				delta.Add("Spec.UserSettings.KernelGatewayAppSettings.CustomImages", a.ko.Spec.UserSettings.KernelGatewayAppSettings.CustomImages, b.ko.Spec.UserSettings.KernelGatewayAppSettings.CustomImages)
			}
			if ackcompare.HasNilDifference(a.ko.Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec, b.ko.Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec) {
				delta.Add("Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec", a.ko.Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec, b.ko.Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec)
			} else if a.ko.Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec != nil && b.ko.Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec != nil {
				if ackcompare.HasNilDifference(a.ko.Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.InstanceType, b.ko.Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.InstanceType) {
					delta.Add("Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.InstanceType", a.ko.Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.InstanceType, b.ko.Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.InstanceType)
				} else if a.ko.Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.InstanceType != nil && b.ko.Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.InstanceType != nil {
					if *a.ko.Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.InstanceType != *b.ko.Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.InstanceType {
						delta.Add("Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.InstanceType", a.ko.Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.InstanceType, b.ko.Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.InstanceType)
					}
				}
				if ackcompare.HasNilDifference(a.ko.Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.LifecycleConfigARN, b.ko.Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.LifecycleConfigARN) {
					delta.Add("Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.LifecycleConfigARN", a.ko.Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.LifecycleConfigARN, b.ko.Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.LifecycleConfigARN)
				} else if a.ko.Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.LifecycleConfigARN != nil && b.ko.Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.LifecycleConfigARN != nil {
					if *a.ko.Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.LifecycleConfigARN != *b.ko.Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.LifecycleConfigARN {
						delta.Add("Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.LifecycleConfigARN", a.ko.Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.LifecycleConfigARN, b.ko.Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.LifecycleConfigARN)
					}
				}
				if ackcompare.HasNilDifference(a.ko.Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.SageMakerImageARN, b.ko.Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.SageMakerImageARN) {
					delta.Add("Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.SageMakerImageARN", a.ko.Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.SageMakerImageARN, b.ko.Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.SageMakerImageARN)
				} else if a.ko.Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.SageMakerImageARN != nil && b.ko.Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.SageMakerImageARN != nil {
					if *a.ko.Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.SageMakerImageARN != *b.ko.Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.SageMakerImageARN {
						delta.Add("Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.SageMakerImageARN", a.ko.Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.SageMakerImageARN, b.ko.Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.SageMakerImageARN)
					}
				}
				if ackcompare.HasNilDifference(a.ko.Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.SageMakerImageVersionARN, b.ko.Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.SageMakerImageVersionARN) {
					delta.Add("Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.SageMakerImageVersionARN", a.ko.Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.SageMakerImageVersionARN, b.ko.Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.SageMakerImageVersionARN)
				} else if a.ko.Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.SageMakerImageVersionARN != nil && b.ko.Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.SageMakerImageVersionARN != nil {
					if *a.ko.Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.SageMakerImageVersionARN != *b.ko.Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.SageMakerImageVersionARN {
						delta.Add("Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.SageMakerImageVersionARN", a.ko.Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.SageMakerImageVersionARN, b.ko.Spec.UserSettings.KernelGatewayAppSettings.DefaultResourceSpec.SageMakerImageVersionARN)
					}
				}
			}
			if !ackcompare.SliceStringPEqual(a.ko.Spec.UserSettings.KernelGatewayAppSettings.LifecycleConfigARNs, b.ko.Spec.UserSettings.KernelGatewayAppSettings.LifecycleConfigARNs) {
				delta.Add("Spec.UserSettings.KernelGatewayAppSettings.LifecycleConfigARNs", a.ko.Spec.UserSettings.KernelGatewayAppSettings.LifecycleConfigARNs, b.ko.Spec.UserSettings.KernelGatewayAppSettings.LifecycleConfigARNs)
			}
		}
		if ackcompare.HasNilDifference(a.ko.Spec.UserSettings.RStudioServerProAppSettings, b.ko.Spec.UserSettings.RStudioServerProAppSettings) {
			delta.Add("Spec.UserSettings.RStudioServerProAppSettings", a.ko.Spec.UserSettings.RStudioServerProAppSettings, b.ko.Spec.UserSettings.RStudioServerProAppSettings)
		} else if a.ko.Spec.UserSettings.RStudioServerProAppSettings != nil && b.ko.Spec.UserSettings.RStudioServerProAppSettings != nil {
			if ackcompare.HasNilDifference(a.ko.Spec.UserSettings.RStudioServerProAppSettings.AccessStatus, b.ko.Spec.UserSettings.RStudioServerProAppSettings.AccessStatus) {
				delta.Add("Spec.UserSettings.RStudioServerProAppSettings.AccessStatus", a.ko.Spec.UserSettings.RStudioServerProAppSettings.AccessStatus, b.ko.Spec.UserSettings.RStudioServerProAppSettings.AccessStatus)
			} else if a.ko.Spec.UserSettings.RStudioServerProAppSettings.AccessStatus != nil && b.ko.Spec.UserSettings.RStudioServerProAppSettings.AccessStatus != nil {
				if *a.ko.Spec.UserSettings.RStudioServerProAppSettings.AccessStatus != *b.ko.Spec.UserSettings.RStudioServerProAppSettings.AccessStatus {
					delta.Add("Spec.UserSettings.RStudioServerProAppSettings.AccessStatus", a.ko.Spec.UserSettings.RStudioServerProAppSettings.AccessStatus, b.ko.Spec.UserSettings.RStudioServerProAppSettings.AccessStatus)
				}
			}
			if ackcompare.HasNilDifference(a.ko.Spec.UserSettings.RStudioServerProAppSettings.UserGroup, b.ko.Spec.UserSettings.RStudioServerProAppSettings.UserGroup) {
				delta.Add("Spec.UserSettings.RStudioServerProAppSettings.UserGroup", a.ko.Spec.UserSettings.RStudioServerProAppSettings.UserGroup, b.ko.Spec.UserSettings.RStudioServerProAppSettings.UserGroup)
			} else if a.ko.Spec.UserSettings.RStudioServerProAppSettings.UserGroup != nil && b.ko.Spec.UserSettings.RStudioServerProAppSettings.UserGroup != nil {
				if *a.ko.Spec.UserSettings.RStudioServerProAppSettings.UserGroup != *b.ko.Spec.UserSettings.RStudioServerProAppSettings.UserGroup {
					delta.Add("Spec.UserSettings.RStudioServerProAppSettings.UserGroup", a.ko.Spec.UserSettings.RStudioServerProAppSettings.UserGroup, b.ko.Spec.UserSettings.RStudioServerProAppSettings.UserGroup)
				}
			}
		}
		if !ackcompare.SliceStringPEqual(a.ko.Spec.UserSettings.SecurityGroups, b.ko.Spec.UserSettings.SecurityGroups) {
			delta.Add("Spec.UserSettings.SecurityGroups", a.ko.Spec.UserSettings.SecurityGroups, b.ko.Spec.UserSettings.SecurityGroups)
		}
		if ackcompare.HasNilDifference(a.ko.Spec.UserSettings.SharingSettings, b.ko.Spec.UserSettings.SharingSettings) {
			delta.Add("Spec.UserSettings.SharingSettings", a.ko.Spec.UserSettings.SharingSettings, b.ko.Spec.UserSettings.SharingSettings)
		} else if a.ko.Spec.UserSettings.SharingSettings != nil && b.ko.Spec.UserSettings.SharingSettings != nil {
			if ackcompare.HasNilDifference(a.ko.Spec.UserSettings.SharingSettings.NotebookOutputOption, b.ko.Spec.UserSettings.SharingSettings.NotebookOutputOption) {
				delta.Add("Spec.UserSettings.SharingSettings.NotebookOutputOption", a.ko.Spec.UserSettings.SharingSettings.NotebookOutputOption, b.ko.Spec.UserSettings.SharingSettings.NotebookOutputOption)
			} else if a.ko.Spec.UserSettings.SharingSettings.NotebookOutputOption != nil && b.ko.Spec.UserSettings.SharingSettings.NotebookOutputOption != nil {
				if *a.ko.Spec.UserSettings.SharingSettings.NotebookOutputOption != *b.ko.Spec.UserSettings.SharingSettings.NotebookOutputOption {
					delta.Add("Spec.UserSettings.SharingSettings.NotebookOutputOption", a.ko.Spec.UserSettings.SharingSettings.NotebookOutputOption, b.ko.Spec.UserSettings.SharingSettings.NotebookOutputOption)
				}
			}
			if ackcompare.HasNilDifference(a.ko.Spec.UserSettings.SharingSettings.S3KMSKeyID, b.ko.Spec.UserSettings.SharingSettings.S3KMSKeyID) {
				delta.Add("Spec.UserSettings.SharingSettings.S3KMSKeyID", a.ko.Spec.UserSettings.SharingSettings.S3KMSKeyID, b.ko.Spec.UserSettings.SharingSettings.S3KMSKeyID)
			} else if a.ko.Spec.UserSettings.SharingSettings.S3KMSKeyID != nil && b.ko.Spec.UserSettings.SharingSettings.S3KMSKeyID != nil {
				if *a.ko.Spec.UserSettings.SharingSettings.S3KMSKeyID != *b.ko.Spec.UserSettings.SharingSettings.S3KMSKeyID {
					delta.Add("Spec.UserSettings.SharingSettings.S3KMSKeyID", a.ko.Spec.UserSettings.SharingSettings.S3KMSKeyID, b.ko.Spec.UserSettings.SharingSettings.S3KMSKeyID)
				}
			}
			if ackcompare.HasNilDifference(a.ko.Spec.UserSettings.SharingSettings.S3OutputPath, b.ko.Spec.UserSettings.SharingSettings.S3OutputPath) {
				delta.Add("Spec.UserSettings.SharingSettings.S3OutputPath", a.ko.Spec.UserSettings.SharingSettings.S3OutputPath, b.ko.Spec.UserSettings.SharingSettings.S3OutputPath)
			} else if a.ko.Spec.UserSettings.SharingSettings.S3OutputPath != nil && b.ko.Spec.UserSettings.SharingSettings.S3OutputPath != nil {
				if *a.ko.Spec.UserSettings.SharingSettings.S3OutputPath != *b.ko.Spec.UserSettings.SharingSettings.S3OutputPath {
					delta.Add("Spec.UserSettings.SharingSettings.S3OutputPath", a.ko.Spec.UserSettings.SharingSettings.S3OutputPath, b.ko.Spec.UserSettings.SharingSettings.S3OutputPath)
				}
			}
		}
		if ackcompare.HasNilDifference(a.ko.Spec.UserSettings.TensorBoardAppSettings, b.ko.Spec.UserSettings.TensorBoardAppSettings) {
			delta.Add("Spec.UserSettings.TensorBoardAppSettings", a.ko.Spec.UserSettings.TensorBoardAppSettings, b.ko.Spec.UserSettings.TensorBoardAppSettings)
		} else if a.ko.Spec.UserSettings.TensorBoardAppSettings != nil && b.ko.Spec.UserSettings.TensorBoardAppSettings != nil {
			if ackcompare.HasNilDifference(a.ko.Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec, b.ko.Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec) {
				delta.Add("Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec", a.ko.Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec, b.ko.Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec)
			} else if a.ko.Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec != nil && b.ko.Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec != nil {
				if ackcompare.HasNilDifference(a.ko.Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.InstanceType, b.ko.Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.InstanceType) {
					delta.Add("Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.InstanceType", a.ko.Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.InstanceType, b.ko.Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.InstanceType)
				} else if a.ko.Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.InstanceType != nil && b.ko.Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.InstanceType != nil {
					if *a.ko.Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.InstanceType != *b.ko.Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.InstanceType {
						delta.Add("Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.InstanceType", a.ko.Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.InstanceType, b.ko.Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.InstanceType)
					}
				}
				if ackcompare.HasNilDifference(a.ko.Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.LifecycleConfigARN, b.ko.Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.LifecycleConfigARN) {
					delta.Add("Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.LifecycleConfigARN", a.ko.Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.LifecycleConfigARN, b.ko.Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.LifecycleConfigARN)
				} else if a.ko.Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.LifecycleConfigARN != nil && b.ko.Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.LifecycleConfigARN != nil {
					if *a.ko.Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.LifecycleConfigARN != *b.ko.Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.LifecycleConfigARN {
						delta.Add("Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.LifecycleConfigARN", a.ko.Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.LifecycleConfigARN, b.ko.Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.LifecycleConfigARN)
					}
				}
				if ackcompare.HasNilDifference(a.ko.Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.SageMakerImageARN, b.ko.Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.SageMakerImageARN) {
					delta.Add("Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.SageMakerImageARN", a.ko.Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.SageMakerImageARN, b.ko.Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.SageMakerImageARN)
				} else if a.ko.Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.SageMakerImageARN != nil && b.ko.Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.SageMakerImageARN != nil {
					if *a.ko.Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.SageMakerImageARN != *b.ko.Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.SageMakerImageARN {
						delta.Add("Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.SageMakerImageARN", a.ko.Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.SageMakerImageARN, b.ko.Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.SageMakerImageARN)
					}
				}
				if ackcompare.HasNilDifference(a.ko.Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.SageMakerImageVersionARN, b.ko.Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.SageMakerImageVersionARN) {
					delta.Add("Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.SageMakerImageVersionARN", a.ko.Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.SageMakerImageVersionARN, b.ko.Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.SageMakerImageVersionARN)
				} else if a.ko.Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.SageMakerImageVersionARN != nil && b.ko.Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.SageMakerImageVersionARN != nil {
					if *a.ko.Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.SageMakerImageVersionARN != *b.ko.Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.SageMakerImageVersionARN {
						delta.Add("Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.SageMakerImageVersionARN", a.ko.Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.SageMakerImageVersionARN, b.ko.Spec.UserSettings.TensorBoardAppSettings.DefaultResourceSpec.SageMakerImageVersionARN)
					}
				}
			}
		}
	}

	return delta
}
