package google

import (
	"fmt"
	"log"
	"strings"

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

// flattenComputeMetadata transforms a list of MetadataItems (as returned via the GCP client) into a simple map from key
// to value.
func flattenComputeMetadata(metadata []*compute.MetadataItems) map[string]string {
	m := map[string]string{}

	for _, item := range metadata {
		// check for duplicates
		if item.Value == nil {
			continue
		}
		if val, ok := m[item.Key]; ok {
			// warn loudly!
			log.Printf("[WARN] Key '%s' already has value '%s' when flattening - ignoring incoming value '%s'",
				item.Key,
				val,
				*item.Value)
		}
		m[item.Key] = *item.Value
	}

	return m
}

// expandComputeMetadata transforms a map representing computing metadata into a list of compute.MetadataItems suitable
// for the GCP client.
func expandComputeMetadata(m map[string]string) []*compute.MetadataItems {
	metadata := make([]*compute.MetadataItems, len(m))

	idx := 0
	for key, value := range m {
		// Make a copy of value as we need a ptr type; if we directly use 'value' then all items will reference the same
		// memory address
		vtmp := value
		metadata[idx] = &compute.MetadataItems{Key: key, Value: &vtmp}
		idx++
	}

	return metadata
}
