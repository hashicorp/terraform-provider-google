package google

import (
	"fmt"
	"regexp"
)

const (
	globalLinkTemplate          = "projects/%s/global/%s/%s"
	globalLinkBasePattern       = "projects/(.+)/global/%s/(.+)"
	zonalLinkTemplate           = "projects/%s/zones/%s/%s/%s"
	zonalLinkBasePattern        = "projects/(.+)/zones/(.+)/%s/(.+)"
	zonalPartialLinkBasePattern = "zones/(.+)/%s/(.+)"
)

// ------------------------------------------------------------
// Field helpers
// ------------------------------------------------------------

func ParseNetworkFieldValue(network string, d TerraformResourceData, config *Config) (*GlobalFieldValue, error) {
	return parseGlobalFieldValue("networks", network, "project", d, config, true)
}

func ParseSslCertificateFieldValue(sslCertificate string, d TerraformResourceData, config *Config) (*GlobalFieldValue, error) {
	return parseGlobalFieldValue("sslCertificates", sslCertificate, "project", d, config, false)
}

func ParseHttpHealthCheckFieldValue(healthCheck string, d TerraformResourceData, config *Config) (*GlobalFieldValue, error) {
	return parseGlobalFieldValue("httpHealthChecks", healthCheck, "project", d, config, false)
}

func ParseDiskFieldValue(disk string, d TerraformResourceData, config *Config) (*ZonalFieldValue, error) {
	return parseZonalFieldValue("disks", disk, "project", "zone", d, config, false)
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

// Parses a global field supporting 5 different formats:
// - https://www.googleapis.com/compute/ANY_VERSION/projects/{my_project}/global/{resource_type}/{resource_name}
// - projects/{my_project}/global/{resource_type}/{resource_name}
// - global/{resource_type}/{resource_name}
// - resource_name
// - "" (empty string). RelativeLink() returns empty if isEmptyValid is true.
//
// If the project is not specified, it first tries to get the project from the `projectSchemaField` and then fallback on the default project.
func parseGlobalFieldValue(resourceType, fieldValue, projectSchemaField string, d TerraformResourceData, config *Config, isEmptyValid bool) (*GlobalFieldValue, error) {
	if len(fieldValue) == 0 {
		if isEmptyValid {
			return &GlobalFieldValue{resourceType: resourceType}, nil
		}
		return nil, fmt.Errorf("The global field for resource %s cannot be empty", resourceType)
	}

	r := regexp.MustCompile(fmt.Sprintf(globalLinkBasePattern, resourceType))
	if parts := r.FindStringSubmatch(fieldValue); parts != nil {
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

type ZonalFieldValue struct {
	Project string
	Zone    string
	Name    string

	resourceType string
}

func (f ZonalFieldValue) RelativeLink() string {
	if len(f.Name) == 0 {
		return ""
	}

	return fmt.Sprintf(zonalLinkTemplate, f.Project, f.Zone, f.resourceType, f.Name)
}

// Parses a zonal field supporting 5 different formats:
// - https://www.googleapis.com/compute/ANY_VERSION/projects/{my_project}/zones/{zone}/{resource_type}/{resource_name}
// - projects/{my_project}/zones/{zone}/{resource_type}/{resource_name}
// - zones/{zone}/{resource_type}/{resource_name}
// - resource_name
// - "" (empty string). RelativeLink() returns empty if isEmptyValid is true.
//
// If the project is not specified, it first tries to get the project from the `projectSchemaField` and then fallback on the default project.
// If the zone is not specified, it takes the value of `zoneSchemaField`.
func parseZonalFieldValue(resourceType, fieldValue, projectSchemaField, zoneSchemaField string, d TerraformResourceData, config *Config, isEmptyValid bool) (*ZonalFieldValue, error) {
	if len(fieldValue) == 0 {
		if isEmptyValid {
			return &ZonalFieldValue{resourceType: resourceType}, nil
		}
		return nil, fmt.Errorf("The zonal field for resource %s cannot be empty.", resourceType)
	}

	r := regexp.MustCompile(fmt.Sprintf(zonalLinkBasePattern, resourceType))
	if parts := r.FindStringSubmatch(fieldValue); parts != nil {
		return &ZonalFieldValue{
			Project:      parts[1],
			Zone:         parts[2],
			Name:         parts[3],
			resourceType: resourceType,
		}, nil
	}

	project, err := getProjectFromSchema(projectSchemaField, d, config)
	if err != nil {
		return nil, err
	}

	r = regexp.MustCompile(fmt.Sprintf(zonalPartialLinkBasePattern, resourceType))
	if parts := r.FindStringSubmatch(fieldValue); parts != nil {
		return &ZonalFieldValue{
			Project:      project,
			Zone:         parts[1],
			Name:         parts[2],
			resourceType: resourceType,
		}, nil
	}

	if len(zoneSchemaField) == 0 {
		return nil, fmt.Errorf("Invalid field format. Got '%s', expected format '%s'", fieldValue, fmt.Sprintf(globalLinkTemplate, "{project}", resourceType, "{name}"))
	}

	zone, ok := d.GetOk(zoneSchemaField)
	if !ok {
		return nil, fmt.Errorf("A zone must be specified")
	}

	return &ZonalFieldValue{
		Project:      project,
		Zone:         zone.(string),
		Name:         GetResourceNameFromSelfLink(fieldValue),
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
