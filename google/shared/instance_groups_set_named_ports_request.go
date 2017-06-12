package shared

import (
	"google.golang.org/api/compute/v1"
)

type InstanceGroupsSetNamedPortsRequest struct {
	// Fingerprint: The fingerprint of the named ports information for this
	// instance group. Use this optional property to prevent conflicts when
	// multiple users change the named ports settings concurrently. Obtain
	// the fingerprint with the instanceGroups.get method. Then, include the
	// fingerprint in your request to ensure that you do not overwrite
	// changes that were applied from another concurrent request.
	Fingerprint string `json:"fingerprint,omitempty"`

	// NamedPorts: The list of named ports to set for this instance group.
	NamedPorts []*NamedPort `json:"namedPorts,omitempty"`

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

func (s *InstanceGroupsSetNamedPortsRequest) ToProduction() *compute.InstanceGroupsSetNamedPortsRequest {
	if s == nil {
		return nil
	}

	return &compute.InstanceGroupsSetNamedPortsRequest{
		Fingerprint:     s.Fingerprint,
		NamedPorts:      NamedPortsToProduction(s.NamedPorts),
		ForceSendFields: s.ForceSendFields,
		NullFields:      s.NullFields,
	}
}

func InstanceGroupsSetNamedPortsRequestFromProduction(s *compute.InstanceGroupsSetNamedPortsRequest) *InstanceGroupsSetNamedPortsRequest {
	if s == nil {
		return nil
	}

	return &InstanceGroupsSetNamedPortsRequest{
		Fingerprint:     s.Fingerprint,
		NamedPorts:      NamedPortsFromProduction(s.NamedPorts),
		ForceSendFields: s.ForceSendFields,
		NullFields:      s.NullFields,
	}
}
