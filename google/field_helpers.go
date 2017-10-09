package google

import (
	"fmt"
	"regexp"
)

const (
	globalLinkTemplate    = "projects/%s/global/%s/%s"
	globalLinkBasePattern = "projects/(.+)/global/%s/(.+)"
)

// ------------------------------------------------------------
// Field helpers
// ------------------------------------------------------------

func ParseNetworkFieldValue(network string, d TerraformResourceData, config *Config) (*GlobalFieldValue, error) {
	return parseGlobalFieldValue("networks", network, "project", d, config, true)
}

// ------------------------------------------------------------
// Base helpers used to create helpers for specific fields.
// ------------------------------------------------------------

type GlobalFieldValue struct {
	Project string
	Name    string

	resourceType string
}

func (f GlobalFieldValue) RelativeLink() string {
	if len(f.Name) == 0 {
		return ""
	}

	return fmt.Sprintf(globalLinkTemplate, f.Project, f.resourceType, f.Name)
}

// Parses a global field supporting 4 different formats:
// - https://www.googleapis.com/compute/ANY_VERSION/projects/{my-project}/global/{resource_type}/{resource_name}
// - projects/{my-project}/global/{resource_type}/{resource_name}
// - global/{resource_type}/{resource_name} (default project is used)
// - resource_name (default project is used)
// - "" (empty string). RelativeLink() returns empty if isEmptyValid is true.
func parseGlobalFieldValue(resourceType, fieldValue, projectSchemaField string, d TerraformResourceData, config *Config, isEmptyValid bool) (*GlobalFieldValue, error) {
	if len(fieldValue) == 0 {
		if isEmptyValid {
			return &GlobalFieldValue{resourceType: resourceType}, nil
		}
		return nil, fmt.Errorf("The global field for resource %s cannot be empty", resourceType)
	}

	r := regexp.MustCompile(fmt.Sprintf(globalLinkBasePattern, resourceType))

	if r.MatchString(fieldValue) {
		parts := r.FindStringSubmatch(fieldValue)

		return &GlobalFieldValue{
			Project: parts[1],
			Name:    parts[2],

			resourceType: resourceType,
		}, nil
	}

	project, err := getProjectFromSchema(projectSchemaField, d, config)
	if err != nil {
		return nil, err
	}

	return &GlobalFieldValue{
		Project: project,
		Name:    GetResourceNameFromSelfLink(fieldValue),

		resourceType: resourceType,
	}, nil
}

func getProjectFromSchema(projectSchemaField string, d TerraformResourceData, config *Config) (string, error) {
	res, ok := d.GetOk(projectSchemaField)
	if !ok || len(projectSchemaField) == 0 {
		if config.Project != "" {
			return config.Project, nil
		}
		return "", fmt.Errorf("project: required field is not set")
	}
	return res.(string), nil
}
