package shared

import (
	"google.golang.org/api/compute/v1"
)

// NamedPort: The named port. For example: .
type NamedPort struct {
	// Name: The name for this named port. The name must be 1-63 characters
	// long, and comply with RFC1035.
	Name string `json:"name,omitempty"`

	// Port: The port number, which can be a value between 1 and 65535.
	Port int64 `json:"port,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Name") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Name") to include in API
	// requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *NamedPort) ToProduction() *compute.NamedPort {
	if s == nil {
		return nil
	}

	n := compute.NamedPort(*s)
	return &n
}

func NamedPortFromProduction(s *compute.NamedPort) *NamedPort {
	if s == nil {
		return nil
	}

	n := NamedPort(*s)
	return &n
}

func NamedPortsToProduction(namedPorts []*NamedPort) []*compute.NamedPort {
	arr := make([]*compute.NamedPort, 0, len(namedPorts))
	for _, v := range namedPorts {
		arr = append(arr, v.ToProduction())
	}
	return arr
}

func NamedPortsFromProduction(namedPorts []*compute.NamedPort) []*NamedPort {
	arr := make([]*NamedPort, 0, len(namedPorts))
	for _, v := range namedPorts {
		arr = append(arr, NamedPortFromProduction(v))
	}
	return arr
}
