package shared

import (
	"google.golang.org/api/compute/v1"
)

type InstanceGroupManagerActionsSummary struct {
	// Abandoning: [Output Only] The total number of instances in the
	// managed instance group that are scheduled to be abandoned. Abandoning
	// an instance removes it from the managed instance group without
	// deleting it.
	Abandoning int64 `json:"abandoning,omitempty"`

	// Creating: [Output Only] The number of instances in the managed
	// instance group that are scheduled to be created or are currently
	// being created. If the group fails to create any of these instances,
	// it tries again until it creates the instance successfully.
	//
	// If you have disabled creation retries, this field will not be
	// populated; instead, the creatingWithoutRetries field will be
	// populated.
	Creating int64 `json:"creating,omitempty"`

	// CreatingWithoutRetries: [Output Only] The number of instances that
	// the managed instance group will attempt to create. The group attempts
	// to create each instance only once. If the group fails to create any
	// of these instances, it decreases the group's targetSize value
	// accordingly.
	CreatingWithoutRetries int64 `json:"creatingWithoutRetries,omitempty"`

	// Deleting: [Output Only] The number of instances in the managed
	// instance group that are scheduled to be deleted or are currently
	// being deleted.
	Deleting int64 `json:"deleting,omitempty"`

	// None: [Output Only] The number of instances in the managed instance
	// group that are running and have no scheduled actions.
	None int64 `json:"none,omitempty"`

	// Recreating: [Output Only] The number of instances in the managed
	// instance group that are scheduled to be recreated or are currently
	// being being recreated. Recreating an instance deletes the existing
	// root persistent disk and creates a new disk from the image that is
	// defined in the instance template.
	Recreating int64 `json:"recreating,omitempty"`

	// Refreshing: [Output Only] The number of instances in the managed
	// instance group that are being reconfigured with properties that do
	// not require a restart or a recreate action. For example, setting or
	// removing target pools for the instance.
	Refreshing int64 `json:"refreshing,omitempty"`

	// Restarting: [Output Only] The number of instances in the managed
	// instance group that are scheduled to be restarted or are currently
	// being restarted.
	Restarting int64 `json:"restarting,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Abandoning") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Abandoning") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *InstanceGroupManagerActionsSummary) ToProduction() *compute.InstanceGroupManagerActionsSummary {
	if s == nil {
		return nil
	}

	n := compute.InstanceGroupManagerActionsSummary(*s)
	return &n
}

func InstanceGroupManagerActionsSummaryFromProduction(s *compute.InstanceGroupManagerActionsSummary) *InstanceGroupManagerActionsSummary {
	if s == nil {
		return nil
	}

	n := InstanceGroupManagerActionsSummary(*s)
	return &n
}
