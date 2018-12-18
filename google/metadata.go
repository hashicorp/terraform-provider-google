package google

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	computeBeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/compute/v1"
)

const FINGERPRINT_RETRIES = 10

var FINGERPRINT_FAIL_ERRORS = []string{"Invalid fingerprint.", "Supplied fingerprint does not match current metadata fingerprint."}

// Since the google compute API uses optimistic locking, there is a chance
// we need to resubmit our updated metadata. To do this, you need to provide
// an update function that attempts to submit your metadata
func MetadataRetryWrapper(update func() error) error {
	attempt := 0
	for attempt < FINGERPRINT_RETRIES {
		err := update()
		if err == nil {
			return nil
		}

		// Check to see if the error matches any of our fingerprint-related failure messages
		var fingerprintError bool
		for _, msg := range FINGERPRINT_FAIL_ERRORS {
			if strings.Contains(err.Error(), msg) {
				fingerprintError = true
				break
			}
		}

		if !fingerprintError {
			// Something else went wrong, don't retry
			return err
		}

		attempt++
	}

	return fmt.Errorf("Failed to update metadata after %d retries", attempt)
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
func BetaMetadataUpdate(oldMDMap map[string]interface{}, newMDMap map[string]interface{}, serverMD *computeBeta.Metadata) {
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
		serverMD.Items = append(serverMD.Items, &computeBeta.MetadataItems{
			Key:   key,
			Value: &v,
		})
	}
}

func expandComputeMetadata(m map[string]interface{}) []*compute.MetadataItems {
	metadata := make([]*compute.MetadataItems, len(m))
	// Append new metadata to existing metadata
	for key, val := range m {
		v := val.(string)
		metadata = append(metadata, &compute.MetadataItems{
			Key:   key,
			Value: &v,
		})
	}

	return metadata
}

func flattenMetadataBeta(metadata *computeBeta.Metadata) map[string]string {
	metadataMap := make(map[string]string)
	for _, item := range metadata.Items {
		metadataMap[item.Key] = *item.Value
	}
	return metadataMap
}

// This function differs from flattenMetadataBeta only in that it takes
// compute.metadata rather than computeBeta.metadata as an argument. It should
// be removed in favour of flattenMetadataBeta if/when all resources using it get
// beta support.
func flattenMetadata(metadata *compute.Metadata) map[string]interface{} {
	metadataMap := make(map[string]interface{})
	for _, item := range metadata.Items {
		metadataMap[item.Key] = *item.Value
	}
	return metadataMap
}

func resourceInstanceMetadata(d *schema.ResourceData) (*computeBeta.Metadata, error) {
	m := &computeBeta.Metadata{}
	mdMap := d.Get("metadata").(map[string]interface{})
	if v, ok := d.GetOk("metadata_startup_script"); ok && v.(string) != "" {
		if ss, ok := mdMap["startup-script"]; ok && ss != "" {
			return nil, errors.New("Cannot provide both metadata_startup_script and metadata.startup-script.")
		}
		mdMap["startup-script"] = v
	}
	if len(mdMap) > 0 {
		m.Items = make([]*computeBeta.MetadataItems, 0, len(mdMap))
		for key, val := range mdMap {
			v := val.(string)
			m.Items = append(m.Items, &computeBeta.MetadataItems{
				Key:   key,
				Value: &v,
			})
		}

		// Set the fingerprint. If the metadata has never been set before
		// then this will just be blank.
		m.Fingerprint = d.Get("metadata_fingerprint").(string)
	}

	return m, nil
}
