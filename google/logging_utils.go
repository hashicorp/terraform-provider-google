package google

import (
	"fmt"
	"regexp"
)

// loggingSinkTypes contains all the possible Stackdriver Logging resource types. Used to parse ids safely.
var loggingSinkTypes = []string{
	"billingAccount",
	"folders",
	"organizations",
	"projects",
}

// LoggingSinkId represents the parts that make up the canonical id used within terraform for a logging resource.
type LoggingSinkId struct {
	typ     string
	typName string
	name    string
}

// loggingSinkIdRegex matches valid logging sink canonical ids
var loggingSinkIdRegex = regexp.MustCompile("(.+)/(.+)/sinks/(.+)")

// canonicalId returns the LoggingSinkId as the canonical id used within terraform.
func (l LoggingSinkId) canonicalId() string {
	return fmt.Sprintf("%s/%s/sinks/%s", l.typ, l.typName, l.name)
}

// parent returns the "parent-level" resource that the sink is in (e.g. `folders/foo` for id `folders/foo/sinks/bar`)
func (l LoggingSinkId) parent() string {
	return fmt.Sprintf("%s/%s", l.typ, l.typName)
}

// parseLoggingSinkId parses a canonical id into a LoggingSinkId, or returns an error on failure.
func parseLoggingSinkId(id string) (*LoggingSinkId, error) {
	parts := loggingSinkIdRegex.FindStringSubmatch(id)
	if parts == nil {
		return nil, fmt.Errorf("unable to parse logging sink id %#v", id)
	}
	// If our type is not a valid logging sink type, complain loudly
	validLoggingSinkType := false
	for _, v := range loggingSinkTypes {
		if v == parts[1] {
			validLoggingSinkType = true
			break
		}
	}

	if !validLoggingSinkType {
		return nil, fmt.Errorf("Logging type %s is not valid. Valid types: %#v", parts[1], loggingSinkTypes)
	}
	return &LoggingSinkId{
		typ:     parts[1],
		typName: parts[2],
		name:    parts[3],
	}, nil
}
