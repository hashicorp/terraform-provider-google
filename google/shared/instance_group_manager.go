package shared

import (
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

// InstanceGroupManager: An Instance Group Manager resource.
type InstanceGroupManager struct {
	// BaseInstanceName: The base instance name to use for instances in this
	// group. The value must be 1-58 characters long. Instances are named by
	// appending a hyphen and a random four-character string to the base
	// instance name. The base instance name must comply with RFC1035.
	BaseInstanceName string `json:"baseInstanceName,omitempty"`

	// CreationTimestamp: [Output Only] The creation timestamp for this
	// managed instance group in RFC3339 text format.
	CreationTimestamp string `json:"creationTimestamp,omitempty"`

	// CurrentActions: [Output Only] The list of instance actions and the
	// number of instances in this managed instance group that are scheduled
	// for each of those actions.
	CurrentActions *InstanceGroupManagerActionsSummary `json:"currentActions,omitempty"`

	// Description: An optional description of this resource. Provide this
	// property when you create the resource.
	Description string `json:"description,omitempty"`

	// Fingerprint: [Output Only] The fingerprint of the resource data. You
	// can use this optional field for optimistic locking when you update
	// the resource.
	Fingerprint string `json:"fingerprint,omitempty"`

	// Id: [Output Only] A unique identifier for this resource type. The
	// server generates this identifier.
	Id uint64 `json:"id,omitempty,string"`

	// InstanceGroup: [Output Only] The URL of the Instance Group resource.
	InstanceGroup string `json:"instanceGroup,omitempty"`

	// InstanceTemplate: The URL of the instance template that is specified
	// for this managed instance group. The group uses this template to
	// create all new instances in the managed instance group.
	InstanceTemplate string `json:"instanceTemplate,omitempty"`

	// Kind: [Output Only] The resource type, which is always
	// compute#instanceGroupManager for managed instance groups.
	Kind string `json:"kind,omitempty"`

	// Name: The name of the managed instance group. The name must be 1-63
	// characters long, and comply with RFC1035.
	Name string `json:"name,omitempty"`

	// NamedPorts: Named ports configured for the Instance Groups
	// complementary to this Instance Group Manager.
	NamedPorts []*NamedPort `json:"namedPorts,omitempty"`

	// Region: [Output Only] The URL of the region where the managed
	// instance group resides (for regional resources).
	Region string `json:"region,omitempty"`

	// SelfLink: [Output Only] The URL for this managed instance group. The
	// server defines this URL.
	SelfLink string `json:"selfLink,omitempty"`

	// TargetPools: The URLs for all TargetPool resources to which instances
	// in the instanceGroup field are added. The target pools automatically
	// apply to all of the instances in the managed instance group.
	TargetPools []string `json:"targetPools,omitempty"`

	// TargetSize: The target number of running instances for this managed
	// instance group. Deleting or abandoning instances reduces this number.
	// Resizing the group changes this number.
	TargetSize int64 `json:"targetSize,omitempty"`

	// Zone: [Output Only] The URL of the zone where the managed instance
	// group is located (for zonal resources).
	Zone string `json:"zone,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "BaseInstanceName") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "BaseInstanceName") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *InstanceGroupManager) ToProduction() *compute.InstanceGroupManager {
	if s == nil {
		return nil
	}

	return &compute.InstanceGroupManager{
		BaseInstanceName:  s.BaseInstanceName,
		CreationTimestamp: s.CreationTimestamp,
		CurrentActions:    s.CurrentActions.ToProduction(),
		Description:       s.Description,
		Fingerprint:       s.Fingerprint,
		Id:                s.Id,
		InstanceGroup:     s.InstanceGroup,
		InstanceTemplate:  s.InstanceTemplate,
		Kind:              s.Kind,
		Name:              s.Name,
		NamedPorts:        NamedPortsToProduction(s.NamedPorts),
		Region:            s.Region,
		SelfLink:          s.SelfLink,
		TargetPools:       s.TargetPools,
		TargetSize:        s.TargetSize,
		Zone:              s.Zone,
		ServerResponse:    s.ServerResponse,
		ForceSendFields:   s.ForceSendFields,
		NullFields:        s.NullFields,
	}
}

func InstanceGroupManagerFromProduction(s *compute.InstanceGroupManager) *InstanceGroupManager {
	if s == nil {
		return nil
	}

	return &InstanceGroupManager{
		BaseInstanceName:  s.BaseInstanceName,
		CreationTimestamp: s.CreationTimestamp,
		CurrentActions:    InstanceGroupManagerActionsSummaryFromProduction(s.CurrentActions),
		Description:       s.Description,
		Fingerprint:       s.Fingerprint,
		Id:                s.Id,
		InstanceGroup:     s.InstanceGroup,
		InstanceTemplate:  s.InstanceTemplate,
		Kind:              s.Kind,
		Name:              s.Name,
		NamedPorts:        NamedPortsFromProduction(s.NamedPorts),
		Region:            s.Region,
		SelfLink:          s.SelfLink,
		TargetPools:       s.TargetPools,
		TargetSize:        s.TargetSize,
		Zone:              s.Zone,
		ServerResponse:    s.ServerResponse,
		ForceSendFields:   s.ForceSendFields,
		NullFields:        s.NullFields,
	}
}
