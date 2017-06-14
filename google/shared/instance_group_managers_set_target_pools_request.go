package shared

import (
	"google.golang.org/api/compute/v1"
)

type InstanceGroupManagersSetTargetPoolsRequest struct {
	// Fingerprint: The fingerprint of the target pools information. Use
	// this optional property to prevent conflicts when multiple users
	// change the target pools settings concurrently. Obtain the fingerprint
	// with the instanceGroupManagers.get method. Then, include the
	// fingerprint in your request to ensure that you do not overwrite
	// changes that were applied from another concurrent request.
	Fingerprint string `json:"fingerprint,omitempty"`

	// TargetPools: The list of target pool URLs that instances in this
	// managed instance group belong to. The managed instance group applies
	// these target pools to all of the instances in the group. Existing
	// instances and new instances in the group all receive these target
	// pool settings.
	TargetPools []string `json:"targetPools,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Fingerprint") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Fingerprint") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *InstanceGroupManagersSetTargetPoolsRequest) ToProduction() *compute.InstanceGroupManagersSetTargetPoolsRequest {
	if s == nil {
		return nil
	}

	n := compute.InstanceGroupManagersSetTargetPoolsRequest(*s)
	return &n
}

func InstanceGroupManagersSetTargetPoolsRequestFromProduction(s *compute.InstanceGroupManagersSetTargetPoolsRequest) *InstanceGroupManagersSetTargetPoolsRequest {
	if s == nil {
		return nil
	}

	n := InstanceGroupManagersSetTargetPoolsRequest(*s)
	return &n
}
