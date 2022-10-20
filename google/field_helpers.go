package google

import (
	"fmt"
	"regexp"
)

const (
	globalLinkTemplate             = "projects/%s/global/%s/%s"
	globalLinkBasePattern          = "projects/(.+)/global/%s/(.+)"
	zonalLinkTemplate              = "projects/%s/zones/%s/%s/%s"
	zonalLinkBasePattern           = "projects/(.+)/zones/(.+)/%s/(.+)"
	zonalPartialLinkBasePattern    = "zones/(.+)/%s/(.+)"
	regionalLinkTemplate           = "projects/%s/regions/%s/%s/%s"
	regionalLinkBasePattern        = "projects/(.+)/regions/(.+)/%s/(.+)"
	regionalPartialLinkBasePattern = "regions/(.+)/%s/(.+)"
	projectLinkTemplate            = "projects/%s/%s/%s"
	projectBasePattern             = "projects/(.+)/%s/(.+)"
	organizationLinkTemplate       = "organizations/%s/%s/%s"
	organizationBasePattern        = "organizations/(.+)/%s/(.+)"
)

// ------------------------------------------------------------
// Field helpers
// ------------------------------------------------------------

func ParseNetworkFieldValue(network string, d TerraformResourceData, config *Config) (*GlobalFieldValue, error) {
	return parseGlobalFieldValue("networks", network, "project", d, config, true)
}

func ParseSubnetworkFieldValue(subnetwork string, d TerraformResourceData, config *Config) (*RegionalFieldValue, error) {
	return parseRegionalFieldValue("subnetworks", subnetwork, "project", "region", "zone", d, config, true)
}

func ParseSubnetworkFieldValueWithProjectField(subnetwork, projectField string, d TerraformResourceData, config *Config) (*RegionalFieldValue, error) {
	return parseRegionalFieldValue("subnetworks", subnetwork, projectField, "region", "zone", d, config, true)
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

func ParseRegionDiskFieldValue(disk string, d TerraformResourceData, config *Config) (*RegionalFieldValue, error) {
	return parseRegionalFieldValue("disks", disk, "project", "region", "zone", d, config, false)
}

func ParseOrganizationCustomRoleName(role string) (*OrganizationFieldValue, error) {
	return parseOrganizationFieldValue("roles", role, false)
}

func ParseAcceleratorFieldValue(accelerator string, d TerraformResourceData, config *Config) (*ZonalFieldValue, error) {
	return parseZonalFieldValue("acceleratorTypes", accelerator, "project", "zone", d, config, false)
}

func ParseMachineTypesFieldValue(machineType string, d TerraformResourceData, config *Config) (*ZonalFieldValue, error) {
	return parseZonalFieldValue("machineTypes", machineType, "project", "zone", d, config, false)
}

func ParseInstanceFieldValue(instance string, d TerraformResourceData, config *Config) (*ZonalFieldValue, error) {
	return parseZonalFieldValue("instances", instance, "project", "zone", d, config, false)
}

func ParseInstanceGroupFieldValue(instanceGroup string, d TerraformResourceData, config *Config) (*ZonalFieldValue, error) {
	return parseZonalFieldValue("instanceGroups", instanceGroup, "project", "zone", d, config, false)
}

func ParseInstanceTemplateFieldValue(instanceTemplate string, d TerraformResourceData, config *Config) (*GlobalFieldValue, error) {
	return parseGlobalFieldValue("instanceTemplates", instanceTemplate, "project", d, config, false)
}

func ParseMachineImageFieldValue(machineImage string, d TerraformResourceData, config *Config) (*GlobalFieldValue, error) {
	return parseGlobalFieldValue("machineImages", machineImage, "project", d, config, false)
}

func ParseSecurityPolicyFieldValue(securityPolicy string, d TerraformResourceData, config *Config) (*GlobalFieldValue, error) {
	return parseGlobalFieldValue("securityPolicies", securityPolicy, "project", d, config, true)
}

func ParseNetworkEndpointGroupFieldValue(networkEndpointGroup string, d TerraformResourceData, config *Config) (*ZonalFieldValue, error) {
	return parseZonalFieldValue("networkEndpointGroups", networkEndpointGroup, "project", "zone", d, config, false)
}

func ParseNetworkEndpointGroupRegionalFieldValue(networkEndpointGroup string, d TerraformResourceData, config *Config) (*RegionalFieldValue, error) {
	return parseRegionalFieldValue("networkEndpointGroups", networkEndpointGroup, "project", "region", "zone", d, config, false)
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
		zone = config.Zone
		if zone == "" {
			return nil, fmt.Errorf("A zone must be specified")
		}
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
	if ok && projectSchemaField != "" {
		return res.(string), nil
	}
	if config.Project != "" {
		return config.Project, nil
	}
	return "", fmt.Errorf("%s: required field is not set", projectSchemaField)
}

func getBillingProjectFromSchema(billingProjectSchemaField string, d TerraformResourceData, config *Config) (string, error) {
	res, ok := d.GetOk(billingProjectSchemaField)
	if ok && billingProjectSchemaField != "" {
		return res.(string), nil
	}
	if config.BillingProject != "" {
		return config.BillingProject, nil
	}
	return "", fmt.Errorf("%s: required field is not set", billingProjectSchemaField)
}

type OrganizationFieldValue struct {
	OrgId string
	Name  string

	resourceType string
}

func (f OrganizationFieldValue) RelativeLink() string {
	if len(f.Name) == 0 {
		return ""
	}

	return fmt.Sprintf(organizationLinkTemplate, f.OrgId, f.resourceType, f.Name)
}

// Parses an organization field with the following formats:
// - organizations/{my_organizations}/{resource_type}/{resource_name}
func parseOrganizationFieldValue(resourceType, fieldValue string, isEmptyValid bool) (*OrganizationFieldValue, error) {
	if len(fieldValue) == 0 {
		if isEmptyValid {
			return &OrganizationFieldValue{resourceType: resourceType}, nil
		}
		return nil, fmt.Errorf("The organization field for resource %s cannot be empty", resourceType)
	}

	r := regexp.MustCompile(fmt.Sprintf(organizationBasePattern, resourceType))
	if parts := r.FindStringSubmatch(fieldValue); parts != nil {
		return &OrganizationFieldValue{
			OrgId: parts[1],
			Name:  parts[2],

			resourceType: resourceType,
		}, nil
	}

	return nil, fmt.Errorf("Invalid field format. Got '%s', expected format '%s'", fieldValue, fmt.Sprintf(organizationLinkTemplate, "{org_id}", resourceType, "{name}"))
}

type RegionalFieldValue struct {
	Project string
	Region  string
	Name    string

	resourceType string
}

func (f RegionalFieldValue) RelativeLink() string {
	if len(f.Name) == 0 {
		return ""
	}

	return fmt.Sprintf(regionalLinkTemplate, f.Project, f.Region, f.resourceType, f.Name)
}

// Parses a regional field supporting 5 different formats:
// - https://www.googleapis.com/compute/ANY_VERSION/projects/{my_project}/regions/{region}/{resource_type}/{resource_name}
// - projects/{my_project}/regions/{region}/{resource_type}/{resource_name}
// - regions/{region}/{resource_type}/{resource_name}
// - resource_name
// - "" (empty string). RelativeLink() returns empty if isEmptyValid is true.
//
// If the project is not specified, it first tries to get the project from the `projectSchemaField` and then fallback on the default project.
// If the region is not specified, see function documentation for `getRegionFromSchema`.
func parseRegionalFieldValue(resourceType, fieldValue, projectSchemaField, regionSchemaField, zoneSchemaField string, d TerraformResourceData, config *Config, isEmptyValid bool) (*RegionalFieldValue, error) {
	if len(fieldValue) == 0 {
		if isEmptyValid {
			return &RegionalFieldValue{resourceType: resourceType}, nil
		}
		return nil, fmt.Errorf("The regional field for resource %s cannot be empty.", resourceType)
	}

	r := regexp.MustCompile(fmt.Sprintf(regionalLinkBasePattern, resourceType))
	if parts := r.FindStringSubmatch(fieldValue); parts != nil {
		return &RegionalFieldValue{
			Project:      parts[1],
			Region:       parts[2],
			Name:         parts[3],
			resourceType: resourceType,
		}, nil
	}

	project, err := getProjectFromSchema(projectSchemaField, d, config)
	if err != nil {
		return nil, err
	}

	r = regexp.MustCompile(fmt.Sprintf(regionalPartialLinkBasePattern, resourceType))
	if parts := r.FindStringSubmatch(fieldValue); parts != nil {
		return &RegionalFieldValue{
			Project:      project,
			Region:       parts[1],
			Name:         parts[2],
			resourceType: resourceType,
		}, nil
	}

	region, err := getRegionFromSchema(regionSchemaField, zoneSchemaField, d, config)
	if err != nil {
		return nil, err
	}

	return &RegionalFieldValue{
		Project:      project,
		Region:       region,
		Name:         GetResourceNameFromSelfLink(fieldValue),
		resourceType: resourceType,
	}, nil
}

// Infers the region based on the following (in order of priority):
// - `regionSchemaField` in resource schema
// - region extracted from the `zoneSchemaField` in resource schema
// - provider-level region
// - region extracted from the provider-level zone
func getRegionFromSchema(regionSchemaField, zoneSchemaField string, d TerraformResourceData, config *Config) (string, error) {
	// if identical such as GKE location, check if it's a zone first and find
	// the region if so. Otherwise, return as it's a region.
	if regionSchemaField == zoneSchemaField {
		if v, ok := d.GetOk(regionSchemaField); ok {
			if isZone(v.(string)) {
				return getRegionFromZone(v.(string)), nil
			}

			return v.(string), nil
		}
	}

	if v, ok := d.GetOk(regionSchemaField); ok && regionSchemaField != "" {
		return GetResourceNameFromSelfLink(v.(string)), nil
	}
	if v, ok := d.GetOk(zoneSchemaField); ok && zoneSchemaField != "" {
		return getRegionFromZone(v.(string)), nil
	}
	if config.Region != "" {
		return config.Region, nil
	}
	if config.Zone != "" {
		return getRegionFromZone(config.Zone), nil
	}

	return "", fmt.Errorf("Cannot determine region: set in this resource, or set provider-level 'region' or 'zone'.")
}

type ProjectFieldValue struct {
	Project string
	Name    string

	resourceType string
}

func (f ProjectFieldValue) RelativeLink() string {
	if len(f.Name) == 0 {
		return ""
	}

	return fmt.Sprintf(projectLinkTemplate, f.Project, f.resourceType, f.Name)
}

// Parses a project field with the following formats:
// - projects/{my_projects}/{resource_type}/{resource_name}
func parseProjectFieldValue(resourceType, fieldValue, projectSchemaField string, d TerraformResourceData, config *Config, isEmptyValid bool) (*ProjectFieldValue, error) {
	if len(fieldValue) == 0 {
		if isEmptyValid {
			return &ProjectFieldValue{resourceType: resourceType}, nil
		}
		return nil, fmt.Errorf("The project field for resource %s cannot be empty", resourceType)
	}

	r := regexp.MustCompile(fmt.Sprintf(projectBasePattern, resourceType))
	if parts := r.FindStringSubmatch(fieldValue); parts != nil {
		return &ProjectFieldValue{
			Project: parts[1],
			Name:    parts[2],

			resourceType: resourceType,
		}, nil
	}

	project, err := getProjectFromSchema(projectSchemaField, d, config)
	if err != nil {
		return nil, err
	}

	return &ProjectFieldValue{
		Project: project,
		Name:    GetResourceNameFromSelfLink(fieldValue),

		resourceType: resourceType,
	}, nil
}
