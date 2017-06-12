package shared

import (
	"google.golang.org/api/compute/v1"
)

type InstanceGroupManagersSetInstanceTemplateRequest struct {
	// InstanceTemplate: The URL of the instance template that is specified
	// for this managed instance group. The group uses this template to
	// create all new instances in the managed instance group.
	InstanceTemplate string `json:"instanceTemplate,omitempty"`

	// ForceSendFields is a list of field names (e.g. "InstanceTemplate") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "InstanceTemplate") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *InstanceGroupManagersSetInstanceTemplateRequest) ToProduction() *compute.InstanceGroupManagersSetInstanceTemplateRequest {
	if s == nil {
		return nil
	}

	n := compute.InstanceGroupManagersSetInstanceTemplateRequest(*s)
	return &n
}

func InstanceGroupManagersSetInstanceTemplateRequestFromProduction(s *compute.InstanceGroupManagersSetInstanceTemplateRequest) *InstanceGroupManagersSetInstanceTemplateRequest {
	if s == nil {
		return nil
	}

	n := InstanceGroupManagersSetInstanceTemplateRequest(*s)
	return &n
}
