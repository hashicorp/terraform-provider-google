package google

import (
	"encoding/json"
	"fmt"
	"strings"
)

type ComputeApiVersion uint8

const (
	v1 ComputeApiVersion = iota
	v0beta
)

var OrderedComputeApiVersions = []ComputeApiVersion{
	v0beta,
	v1,
}

// Convert between two types by converting to/from JSON. Intended to switch
// between multiple API versions, as they are strict supersets of one another.
// Convert loses information about ForceSendFields and NullFields.
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

// Compare the fields set in schema against a list of features and their versions to determine
// what version of the API is required in order to manage the resource.
func getComputeApiVersion(d TerraformResourceData, resourceVersion ComputeApiVersion, features []Feature) ComputeApiVersion {
	versions := map[ComputeApiVersion]struct{}{resourceVersion: struct{}{}}
	for _, feature := range features {
		if feature.InUseByDefault(d) {
			versions[feature.Version] = struct{}{}
		}
	}

	return maxVersion(versions)
}

// Compare the fields set in schema against a list of features and their version, and a
// list of features that exist at the base resource version that can only be update at some other
// version, to determine what version of the API is required in order to update the resource.
func getComputeApiVersionUpdate(d TerraformResourceData, resourceVersion ComputeApiVersion, features, updateOnlyFields []Feature) ComputeApiVersion {
	versions := map[ComputeApiVersion]struct{}{resourceVersion: struct{}{}}

	for _, feature := range features {
		if feature.InUseByUpdate(d) {
			versions[feature.Version] = struct{}{}
		}
	}

	for _, feature := range updateOnlyFields {
		if feature.HasChangeBy(d) {
			versions[feature.Version] = struct{}{}
		}
	}

	return maxVersion(versions)
}

// A field of a resource and the version of the Compute API required to use it.
type Feature struct {
	Version ComputeApiVersion
	// Path to the beta field.
	//
	// The feature is considered to be in-use if the field referenced by "Item" is set in the state.
	// The path can reference:
	// - a beta field at the top-level (e.g. "min_cpu_platform").
	// - a beta field nested inside a list (e.g. "network_interface.*.alias_ip_range" is considered to be
	// 		in-use if the "alias_ip_range" field is set in the state for any of the network interfaces).
	//
	// Note: beta field nested inside a SET are NOT supported at the moment.
	Item string

	// Optional, only set if your field has a default value.
	// If the value for the field is equal to the DefaultValue, we assume the beta feature is not activated.
	DefaultValue interface{}
}

// Returns true when a feature has been modified.
// This is most important when updating a resource to remove versioned feature usage; if the
// resource is reverting to its base version, it needs to perform a final update at the higher
// version in order to remove high version features.
func (s Feature) HasChangeBy(d TerraformResourceData) bool {
	return d.HasChange(s.Item)
}

type InUseFunc func(d TerraformResourceData, path string, defaultValue interface{}) bool

func defaultInUseFunc(d TerraformResourceData, path string, defaultValue interface{}) bool {
	// At read and delete time, there is no change.
	// At create time, all fields are marked has changed. We should only consider the feature active if the field has
	// a value set and that this value is not the default value.
	value, ok := d.GetOk(path)
	return ok && value != defaultValue
}

func updateInUseFunc(d TerraformResourceData, path string, defaultValue interface{}) bool {
	// During a resource update, if the beta field has changes, the feature is considered active even if the new value
	// is the default value. This is because the beta API must be called to change the value of the field back to the
	// default value.
	value, ok := d.GetOk(path)
	return (ok && value != defaultValue) || d.HasChange(path)
}

// Return true when a feature appears in schema and doesn't hold the default value.
func (s Feature) InUseByDefault(d TerraformResourceData) bool {
	return inUseBy(d, s.Item, s.DefaultValue, defaultInUseFunc)
}

func (s Feature) InUseByUpdate(d TerraformResourceData) bool {
	return inUseBy(d, s.Item, s.DefaultValue, updateInUseFunc)
}

func inUseBy(d TerraformResourceData, path string, defaultValue interface{}, inUseFunc InUseFunc) bool {
	pos := strings.Index(path, "*")
	if pos == -1 {
		return inUseFunc(d, path, defaultValue)
	}

	prefix := path[0:pos]
	suffix := path[pos+1:]

	v, ok := d.GetOk(prefix + "#")

	if !ok {
		return false
	}

	count := v.(int)
	for i := 0; i < count; i++ {
		nestedPath := fmt.Sprintf("%s%d%s", prefix, i, suffix)
		if inUseBy(d, nestedPath, defaultValue, inUseFunc) {
			return true
		}
	}

	return false
}

func maxVersion(versionsInUse map[ComputeApiVersion]struct{}) ComputeApiVersion {
	for _, version := range OrderedComputeApiVersions {
		if _, ok := versionsInUse[version]; ok {
			return version
		}
	}

	// Fallback to the final, most stable version
	return OrderedComputeApiVersions[len(OrderedComputeApiVersions)-1]
}
