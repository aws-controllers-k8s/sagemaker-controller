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
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
)

// Called in delta pre compare, sets spec default values to avoid reconciling errors.
func customSetDefaults(
	a *resource,
	b *resource,
) {
	// DisableGlueTableCreation is false by default, it cannot be null.
	if ackcompare.IsNotNil(a.ko.Spec.OfflineStoreConfig) && ackcompare.IsNotNil(b.ko.Spec.OfflineStoreConfig) {
		if ackcompare.IsNil(a.ko.Spec.OfflineStoreConfig.DisableGlueTableCreation) && ackcompare.IsNotNil(b.ko.Spec.OfflineStoreConfig.DisableGlueTableCreation) {
			a.ko.Spec.OfflineStoreConfig.DisableGlueTableCreation = b.ko.Spec.OfflineStoreConfig.DisableGlueTableCreation
		}
	}

	// DataCatalogConfig has a timestamped generated default value, it cannot be null.
	if ackcompare.IsNotNil(a.ko.Spec.OfflineStoreConfig) && ackcompare.IsNotNil(b.ko.Spec.OfflineStoreConfig) {
		if ackcompare.IsNil(a.ko.Spec.OfflineStoreConfig.DataCatalogConfig) && ackcompare.IsNotNil(b.ko.Spec.OfflineStoreConfig.DataCatalogConfig) {
			a.ko.Spec.OfflineStoreConfig.DataCatalogConfig = b.ko.Spec.OfflineStoreConfig.DataCatalogConfig
		}
	}

	// ResolvedOutputS3URI has a timestamped generated default value, it cannot be null.
	if ackcompare.IsNotNil(a.ko.Spec.OfflineStoreConfig) && ackcompare.IsNotNil(b.ko.Spec.OfflineStoreConfig) {
		if ackcompare.IsNotNil(a.ko.Spec.OfflineStoreConfig.S3StorageConfig) && ackcompare.IsNotNil(b.ko.Spec.OfflineStoreConfig.S3StorageConfig) {
			if ackcompare.IsNil(a.ko.Spec.OfflineStoreConfig.S3StorageConfig.ResolvedOutputS3URI) && ackcompare.IsNotNil(b.ko.Spec.OfflineStoreConfig.S3StorageConfig.ResolvedOutputS3URI) {
				a.ko.Spec.OfflineStoreConfig.S3StorageConfig.ResolvedOutputS3URI = b.ko.Spec.OfflineStoreConfig.S3StorageConfig.ResolvedOutputS3URI
			}
		}
	}
}
