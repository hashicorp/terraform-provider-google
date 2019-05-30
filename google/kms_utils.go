package google

import (
	"fmt"
	"regexp"
	"strings"
)

type kmsKeyRingId struct {
	Project  string
	Location string
	Name     string
}

func (s *kmsKeyRingId) keyRingId() string {
	return fmt.Sprintf("projects/%s/locations/%s/keyRings/%s", s.Project, s.Location, s.Name)
}

func (s *kmsKeyRingId) terraformId() string {
	return fmt.Sprintf("%s/%s/%s", s.Project, s.Location, s.Name)
}

func parseKmsKeyRingId(id string, config *Config) (*kmsKeyRingId, error) {
	parts := strings.Split(id, "/")

	keyRingIdRegex := regexp.MustCompile("^(" + ProjectRegex + ")/([a-z0-9-])+/([a-zA-Z0-9_-]{1,63})$")
	keyRingIdWithoutProjectRegex := regexp.MustCompile("^([a-z0-9-])+/([a-zA-Z0-9_-]{1,63})$")
	keyRingRelativeLinkRegex := regexp.MustCompile("^projects/(" + ProjectRegex + ")/locations/([a-z0-9-]+)/keyRings/([a-zA-Z0-9_-]{1,63})$")

	if keyRingIdRegex.MatchString(id) {
		return &kmsKeyRingId{
			Project:  parts[0],
			Location: parts[1],
			Name:     parts[2],
		}, nil
	}

	if keyRingIdWithoutProjectRegex.MatchString(id) {
		if config.Project == "" {
			return nil, fmt.Errorf("The default project for the provider must be set when using the `{location}/{keyRingName}` id format.")
		}

		return &kmsKeyRingId{
			Project:  config.Project,
			Location: parts[0],
			Name:     parts[1],
		}, nil
	}

	if parts := keyRingRelativeLinkRegex.FindStringSubmatch(id); parts != nil {
		return &kmsKeyRingId{
			Project:  parts[1],
			Location: parts[2],
			Name:     parts[3],
		}, nil
	}
	return nil, fmt.Errorf("Invalid KeyRing id format, expecting `{projectId}/{locationId}/{keyRingName}` or `{locationId}/{keyRingName}.`")
}
