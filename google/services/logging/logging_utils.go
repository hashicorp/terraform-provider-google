// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package logging

import (
	"fmt"
	"regexp"
)

// LoggingSinkResourceTypes contains all the possible Stackdriver Logging resource types. Used to parse ids safely.
var LoggingSinkResourceTypes = []string{
	"billingAccounts",
	"folders",
	"organizations",
	"projects",
}

// LoggingSinkId represents the parts that make up the canonical id used within terraform for a logging resource.
type LoggingSinkId struct {
	resourceType string
	resourceId   string
	name         string
}

// loggingSinkIdRegex matches valid logging sink canonical ids
var loggingSinkIdRegex = regexp.MustCompile("(.+)/(.+)/sinks/(.+)")

// canonicalId returns the LoggingSinkId as the canonical id used within terraform.
func (l LoggingSinkId) canonicalId() string {
	return fmt.Sprintf("%s/%s/sinks/%s", l.resourceType, l.resourceId, l.name)
}

// parent returns the "parent-level" resource that the sink is in (e.g. `folders/foo` for id `folders/foo/sinks/bar`)
func (l LoggingSinkId) parent() string {
	return fmt.Sprintf("%s/%s", l.resourceType, l.resourceId)
}

// ParseLoggingSinkId parses a canonical id into a LoggingSinkId, or returns an error on failure.
func ParseLoggingSinkId(id string) (*LoggingSinkId, error) {
	parts := loggingSinkIdRegex.FindStringSubmatch(id)
	if parts == nil {
		return nil, fmt.Errorf("unable to parse logging sink id %#v", id)
	}
	// If our resourceType is not a valid logging sink resource type, complain loudly
	validLoggingSinkResourceType := false
	for _, v := range LoggingSinkResourceTypes {
		if v == parts[1] {
			validLoggingSinkResourceType = true
			break
		}
	}

	if !validLoggingSinkResourceType {
		return nil, fmt.Errorf("Logging resource type %s is not valid. Valid resource types: %#v", parts[1],
			LoggingSinkResourceTypes)
	}
	return &LoggingSinkId{
		resourceType: parts[1],
		resourceId:   parts[2],
		name:         parts[3],
	}, nil
}
