package google

import (
	"encoding/json"
)

type ComputeApiVersion uint8

const (
	v1 ComputeApiVersion = iota
	v0beta
)

var ORDERED_COMPUTE_API_VERSIONS = []ComputeApiVersion{
	v0beta,
	v1,
}

// Convert between two types by converting to/from JSON. Intended to switch
// between multiple API versions, as they are strict supersets of one another.
func Convert(item, out interface{}) error {
	bytes, err := json.Marshal(item)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, out)
	if err != nil {
		return err
	}

	return nil
}

type TerraformResourceData interface {
	HasChange(string) bool
	GetOk(string) (interface{}, bool)
}

// Compare the keys set in schema against a list of features and their versions to determine
//what version of the API is required in order to manage the resource.
func getComputeApiVersion(d TerraformResourceData, baseVersion ComputeApiVersion, keys []Key) ComputeApiVersion {
	versions := map[ComputeApiVersion]struct{}{baseVersion: struct{}{}}
	for _, key := range keys {
		if key.InUseBy(d) {
			versions[key.Version] = struct{}{}
		}
	}

	return maxVersion(versions)
}

// Compare the keys set in schema against a list of features and their version, and a
// list of features that exist at the base version that can only be update at some other version,
// to determine what version of the API is required in order to update the resource.
func getComputeApiVersionUpdate(d TerraformResourceData, baseVersion ComputeApiVersion, keys, updateOnlyKeys []Key) ComputeApiVersion {
	versions := map[ComputeApiVersion]struct{}{baseVersion: struct{}{}}
	schemaVersion := getComputeApiVersion(d, baseVersion, keys)
	versions[schemaVersion] = struct{}{}

	for _, key := range updateOnlyKeys {
		if key.HasChangeBy(d) {
			versions[key.Version] = struct{}{}
		}
	}

	return maxVersion(versions)
}

type Key struct {
	Version ComputeApiVersion
	Item    string
}

func (s Key) HasChangeBy(d TerraformResourceData) bool {
	return d.HasChange(s.Item)
}

func (s Key) InUseBy(d TerraformResourceData) bool {
	_, ok := d.GetOk(s.Item)
	return ok && s.HasChangeBy(d)
}

func maxVersion(versionsInUse map[ComputeApiVersion]struct{}) ComputeApiVersion {
	for _, version := range ORDERED_COMPUTE_API_VERSIONS {
		if _, ok := versionsInUse[version]; ok {
			return version
		}
	}

	// Fallback to the final, most stable version
	return ORDERED_COMPUTE_API_VERSIONS[len(ORDERED_COMPUTE_API_VERSIONS)-1]
}
