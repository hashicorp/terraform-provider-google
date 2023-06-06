// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"errors"
	"sort"

	"google.golang.org/api/compute/v1"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
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
func MetadataUpdate(oldMDMap map[string]interface{}, newMDMap map[string]interface{}, serverMD *compute.Metadata) {
	curMDMap := make(map[string]string)
	// Load metadata on server into map
	for _, kv := range serverMD.Items {
		// If the server state has a key that we had in our old
		// state, but not in our new state, we should delete it
		_, okOld := oldMDMap[kv.Key]
		_, okNew := newMDMap[kv.Key]
		if okOld && !okNew {
			continue
		} else {
			curMDMap[kv.Key] = *kv.Value
		}
	}

	// Insert new metadata into existing metadata (overwriting when needed)
	for key, val := range newMDMap {
		curMDMap[key] = val.(string)
	}

	// Reformat old metadata into a list
	serverMD.Items = nil
	for key, val := range curMDMap {
		v := val
		serverMD.Items = append(serverMD.Items, &compute.MetadataItems{
			Key:   key,
			Value: &v,
		})
	}
}

// Update the beta metadata (serverMD) according to the provided diff (oldMDMap v
// newMDMap).
func BetaMetadataUpdate(oldMDMap map[string]interface{}, newMDMap map[string]interface{}, serverMD *compute.Metadata) {
	curMDMap := make(map[string]string)
	// Load metadata on server into map
	for _, kv := range serverMD.Items {
		// If the server state has a key that we had in our old
		// state, but not in our new state, we should delete it
		_, okOld := oldMDMap[kv.Key]
		_, okNew := newMDMap[kv.Key]
		if okOld && !okNew {
			continue
		} else {
			curMDMap[kv.Key] = *kv.Value
		}
	}

	// Insert new metadata into existing metadata (overwriting when needed)
	for key, val := range newMDMap {
		curMDMap[key] = val.(string)
	}

	// Reformat old metadata into a list
	serverMD.Items = nil
	for key, val := range curMDMap {
		v := val
		serverMD.Items = append(serverMD.Items, &compute.MetadataItems{
			Key:   key,
			Value: &v,
		})
	}
}

func expandComputeMetadata(m map[string]interface{}) []*compute.MetadataItems {
	metadata := make([]*compute.MetadataItems, len(m))
	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	// Append new metadata to existing metadata
	for _, key := range keys {
		v := m[key].(string)
		metadata = append(metadata, &compute.MetadataItems{
			Key:   key,
			Value: &v,
		})
	}

	return metadata
}

func flattenMetadataBeta(metadata *compute.Metadata) map[string]string {
	metadataMap := make(map[string]string)
	for _, item := range metadata.Items {
		metadataMap[item.Key] = *item.Value
	}
	return metadataMap
}

// This function differs from flattenMetadataBeta only in that it takes
// compute.metadata rather than compute.metadata as an argument. It should
// be removed in favour of flattenMetadataBeta if/when all resources using it get
// beta support.
func FlattenMetadata(metadata *compute.Metadata) map[string]interface{} {
	metadataMap := make(map[string]interface{})
	for _, item := range metadata.Items {
		metadataMap[item.Key] = *item.Value
	}
	return metadataMap
}

func resourceInstanceMetadata(d tpgresource.TerraformResourceData) (*compute.Metadata, error) {
	m := &compute.Metadata{}
	mdMap := d.Get("metadata").(map[string]interface{})
	if v, ok := d.GetOk("metadata_startup_script"); ok && v.(string) != "" {
		if w, ok := mdMap["startup-script"]; ok {
			// metadata.startup-script could be from metadata_startup_script in the first place
			if v != w {
				return nil, errors.New("Cannot provide both metadata_startup_script and metadata.startup-script.")
			}
		}
		mdMap["startup-script"] = v
	}
	if len(mdMap) > 0 {
		m.Items = make([]*compute.MetadataItems, 0, len(mdMap))
		var keys []string
		for k := range mdMap {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			v := mdMap[k].(string)
			m.Items = append(m.Items, &compute.MetadataItems{
				Key:   k,
				Value: &v,
			})
		}

		// Set the fingerprint. If the metadata has never been set before
		// then this will just be blank.
		m.Fingerprint = d.Get("metadata_fingerprint").(string)
	}

	return m, nil
}
