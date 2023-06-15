// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"google.golang.org/api/compute/v1"

	tpgcompute "github.com/hashicorp/terraform-provider-google/google/services/compute"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// Since the google compute API uses optimistic locking, there is a chance
// we need to resubmit our updated metadata. To do this, you need to provide
// an update function that attempts to submit your metadata
func MetadataRetryWrapper(update func() error) error {
	return transport_tpg.MetadataRetryWrapper(update)
}

// Update the metadata (serverMD) according to the provided diff (oldMDMap v
// newMDMap).
//
// Deprecated: For backward compatibility MetadataUpdate is still working,
// but all new code should use PollCheckKnaMetadataUpdatetiveStatusFunc in the tpgcompute package instead.
func MetadataUpdate(oldMDMap map[string]interface{}, newMDMap map[string]interface{}, serverMD *compute.Metadata) {
	tpgcompute.MetadataUpdate(oldMDMap, newMDMap, serverMD)
}

// Update the beta metadata (serverMD) according to the provided diff (oldMDMap v
// newMDMap).
//
// Deprecated: For backward compatibility BetaMetadataUpdate is still working,
// but all new code should use BetaMetadataUpdate in the tpgcompute package instead.
func BetaMetadataUpdate(oldMDMap map[string]interface{}, newMDMap map[string]interface{}, serverMD *compute.Metadata) {
	tpgcompute.BetaMetadataUpdate(oldMDMap, newMDMap, serverMD)
}

// This function differs from flattenMetadataBeta only in that it takes
// compute.metadata rather than compute.metadata as an argument. It should
// be removed in favour of flattenMetadataBeta if/when all resources using it get
// beta support.
//
// Deprecated: For backward compatibility flattenMetadata is still working,
// but all new code should use FlattenMetadata in the tpgcompute package instead.
func flattenMetadata(metadata *compute.Metadata) map[string]interface{} {
	return tpgcompute.FlattenMetadata(metadata)
}
