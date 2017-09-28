package google

import (
	"fmt"
	"regexp"
)

const networkLinkTemplate = "projects/%s/global/networks/%s"

var networkLinkRegex = regexp.MustCompile("projects/(.+)/global/networks/(.+)")

type NetworkFieldValue struct {
	Project string
	Name    string
}

// Parses a `network` supporting 4 different formats:
// - https://www.googleapis.com/compute/{version}/projects/myproject/global/networks/my-network
// - projects/myproject/global/networks/my-network
// - global/networks/my-network (default project is used)
// - my-network (default project is used)
func ParseNetworkFieldValue(network string, config *Config) *NetworkFieldValue {
	if networkLinkRegex.MatchString(network) {
		parts := networkLinkRegex.FindStringSubmatch(network)

		return &NetworkFieldValue{
			Project: parts[1],
			Name:    parts[2],
		}
	}

	return &NetworkFieldValue{
		Project: config.Project,
		Name:    GetResourceNameFromSelfLink(network),
	}
}

func (f NetworkFieldValue) RelativeLink() string {
	return fmt.Sprintf(networkLinkTemplate, f.Project, f.Name)
}
