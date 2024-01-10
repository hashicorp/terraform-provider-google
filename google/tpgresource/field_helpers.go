// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package tpgresource

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

const (
	GlobalLinkTemplate             = "projects/%s/global/%s/%s"
	GlobalLinkBasePattern          = "projects/(.+)/global/%s/(.+)"
	ZonalLinkTemplate              = "projects/%s/zones/%s/%s/%s"
	ZonalLinkBasePattern           = "projects/(.+)/zones/(.+)/%s/(.+)"
	ZonalPartialLinkBasePattern    = "zones/(.+)/%s/(.+)"
	RegionalLinkTemplate           = "projects/%s/regions/%s/%s/%s"
	RegionalLinkBasePattern        = "projects/(.+)/regions/(.+)/%s/(.+)"
	RegionalPartialLinkBasePattern = "regions/(.+)/%s/(.+)"
	ProjectLinkTemplate            = "projects/%s/%s/%s"
	ProjectBasePattern             = "projects/(.+)/%s/(.+)"
	OrganizationLinkTemplate       = "organizations/%s/%s/%s"
	OrganizationBasePattern        = "organizations/(.+)/%s/(.+)"
)

// ------------------------------------------------------------
// Field helpers
// ------------------------------------------------------------

func ParseNetworkFieldValue(network string, d TerraformResourceData, config *transport_tpg.Config) (*GlobalFieldValue, error) {
	return ParseGlobalFieldValue("networks", network, "project", d, config, true)
}

func ParseSubnetworkFieldValue(subnetwork string, d TerraformResourceData, config *transport_tpg.Config) (*RegionalFieldValue, error) {
	return ParseRegionalFieldValue("subnetworks", subnetwork, "project", "region", "zone", d, config, true)
}

func ParseSubnetworkFieldValueWithProjectField(subnetwork, projectField string, d TerraformResourceData, config *transport_tpg.Config) (*RegionalFieldValue, error) {
	return ParseRegionalFieldValue("subnetworks", subnetwork, projectField, "region", "zone", d, config, true)
}

func ParseSslCertificateFieldValue(sslCertificate string, d TerraformResourceData, config *transport_tpg.Config) (*GlobalFieldValue, error) {
	return ParseGlobalFieldValue("sslCertificates", sslCertificate, "project", d, config, false)
}

func ParseHttpHealthCheckFieldValue(healthCheck string, d TerraformResourceData, config *transport_tpg.Config) (*GlobalFieldValue, error) {
	return ParseGlobalFieldValue("httpHealthChecks", healthCheck, "project", d, config, false)
}

func ParseDiskFieldValue(disk string, d TerraformResourceData, config *transport_tpg.Config) (*ZonalFieldValue, error) {
	return ParseZonalFieldValue("disks", disk, "project", "zone", d, config, false)
}

func ParseRegionDiskFieldValue(disk string, d TerraformResourceData, config *transport_tpg.Config) (*RegionalFieldValue, error) {
	return ParseRegionalFieldValue("disks", disk, "project", "region", "zone", d, config, false)
}

func ParseOrganizationCustomRoleName(role string) (*OrganizationFieldValue, error) {
	return ParseOrganizationFieldValue("roles", role, false)
}

func ParseAcceleratorFieldValue(accelerator string, d TerraformResourceData, config *transport_tpg.Config) (*ZonalFieldValue, error) {
	return ParseZonalFieldValue("acceleratorTypes", accelerator, "project", "zone", d, config, false)
}

func ParseMachineTypesFieldValue(machineType string, d TerraformResourceData, config *transport_tpg.Config) (*ZonalFieldValue, error) {
	return ParseZonalFieldValue("machineTypes", machineType, "project", "zone", d, config, false)
}

func ParseInstanceFieldValue(instance string, d TerraformResourceData, config *transport_tpg.Config) (*ZonalFieldValue, error) {
	return ParseZonalFieldValue("instances", instance, "project", "zone", d, config, false)
}

func ParseInstanceGroupFieldValue(instanceGroup string, d TerraformResourceData, config *transport_tpg.Config) (*ZonalFieldValue, error) {
	return ParseZonalFieldValue("instanceGroups", instanceGroup, "project", "zone", d, config, false)
}

func ParseInstanceTemplateFieldValue(instanceTemplate string, d TerraformResourceData, config *transport_tpg.Config) (*GlobalFieldValue, error) {
	return ParseGlobalFieldValue("instanceTemplates", instanceTemplate, "project", d, config, false)
}

func ParseMachineImageFieldValue(machineImage string, d TerraformResourceData, config *transport_tpg.Config) (*GlobalFieldValue, error) {
	return ParseGlobalFieldValue("machineImages", machineImage, "project", d, config, false)
}

func ParseSecurityPolicyFieldValue(securityPolicy string, d TerraformResourceData, config *transport_tpg.Config) (*GlobalFieldValue, error) {
	return ParseGlobalFieldValue("securityPolicies", securityPolicy, "project", d, config, true)
}

func ParseSecurityPolicyRegionalFieldValue(securityPolicy string, d TerraformResourceData, config *transport_tpg.Config) (*RegionalFieldValue, error) {
	return ParseRegionalFieldValue("securityPolicies", securityPolicy, "project", "region", "zone", d, config, true)
}

func ParseNetworkEndpointGroupFieldValue(networkEndpointGroup string, d TerraformResourceData, config *transport_tpg.Config) (*ZonalFieldValue, error) {
	return ParseZonalFieldValue("networkEndpointGroups", networkEndpointGroup, "project", "zone", d, config, false)
}

func ParseNetworkEndpointGroupRegionalFieldValue(networkEndpointGroup string, d TerraformResourceData, config *transport_tpg.Config) (*RegionalFieldValue, error) {
	return ParseRegionalFieldValue("networkEndpointGroups", networkEndpointGroup, "project", "region", "zone", d, config, false)
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

	return fmt.Sprintf(GlobalLinkTemplate, f.Project, f.resourceType, f.Name)
}

// Parses a global field supporting 5 different formats:
// - https://www.googleapis.com/compute/ANY_VERSION/projects/{my_project}/global/{resource_type}/{resource_name}
// - projects/{my_project}/global/{resource_type}/{resource_name}
// - global/{resource_type}/{resource_name}
// - resource_name
// - "" (empty string). RelativeLink() returns empty if isEmptyValid is true.
//
// If the project is not specified, it first tries to get the project from the `projectSchemaField` and then fallback on the default project.
func ParseGlobalFieldValue(resourceType, fieldValue, projectSchemaField string, d TerraformResourceData, config *transport_tpg.Config, isEmptyValid bool) (*GlobalFieldValue, error) {
	if len(fieldValue) == 0 {
		if isEmptyValid {
			return &GlobalFieldValue{resourceType: resourceType}, nil
		}
		return nil, fmt.Errorf("The global field for resource %s cannot be empty", resourceType)
	}

	r := regexp.MustCompile(fmt.Sprintf(GlobalLinkBasePattern, resourceType))
	if parts := r.FindStringSubmatch(fieldValue); parts != nil {
		return &GlobalFieldValue{
			Project: parts[1],
			Name:    parts[2],

			resourceType: resourceType,
		}, nil
	}

	project, err := GetProjectFromSchema(projectSchemaField, d, config)
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

	ResourceType string
}

func (f ZonalFieldValue) RelativeLink() string {
	if len(f.Name) == 0 {
		return ""
	}

	return fmt.Sprintf(ZonalLinkTemplate, f.Project, f.Zone, f.ResourceType, f.Name)
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
func ParseZonalFieldValue(resourceType, fieldValue, projectSchemaField, zoneSchemaField string, d TerraformResourceData, config *transport_tpg.Config, isEmptyValid bool) (*ZonalFieldValue, error) {
	if len(fieldValue) == 0 {
		if isEmptyValid {
			return &ZonalFieldValue{ResourceType: resourceType}, nil
		}
		return nil, fmt.Errorf("The zonal field for resource %s cannot be empty.", resourceType)
	}

	r := regexp.MustCompile(fmt.Sprintf(ZonalLinkBasePattern, resourceType))
	if parts := r.FindStringSubmatch(fieldValue); parts != nil {
		return &ZonalFieldValue{
			Project:      parts[1],
			Zone:         parts[2],
			Name:         parts[3],
			ResourceType: resourceType,
		}, nil
	}

	project, err := GetProjectFromSchema(projectSchemaField, d, config)
	if err != nil {
		return nil, err
	}

	r = regexp.MustCompile(fmt.Sprintf(ZonalPartialLinkBasePattern, resourceType))
	if parts := r.FindStringSubmatch(fieldValue); parts != nil {
		return &ZonalFieldValue{
			Project:      project,
			Zone:         parts[1],
			Name:         parts[2],
			ResourceType: resourceType,
		}, nil
	}

	if len(zoneSchemaField) == 0 {
		return nil, fmt.Errorf("Invalid field format. Got '%s', expected format '%s'", fieldValue, fmt.Sprintf(GlobalLinkTemplate, "{project}", resourceType, "{name}"))
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
		ResourceType: resourceType,
	}, nil
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
func ParseZonalFieldValueDiff(resourceType, fieldValue, projectSchemaField, zoneSchemaField string, d *schema.ResourceDiff, config *transport_tpg.Config, isEmptyValid bool) (*ZonalFieldValue, error) {
	r := regexp.MustCompile(fmt.Sprintf(ZonalLinkBasePattern, resourceType))
	if parts := r.FindStringSubmatch(fieldValue); parts != nil {
		return &ZonalFieldValue{
			Project:      parts[1],
			Zone:         parts[2],
			Name:         parts[3],
			ResourceType: resourceType,
		}, nil
	}

	project, err := GetProjectFromDiff(d, config)
	if err != nil {
		return nil, err
	}

	r = regexp.MustCompile(fmt.Sprintf(ZonalPartialLinkBasePattern, resourceType))
	if parts := r.FindStringSubmatch(fieldValue); parts != nil {
		return &ZonalFieldValue{
			Project:      project,
			Zone:         parts[1],
			Name:         parts[2],
			ResourceType: resourceType,
		}, nil
	}

	if len(zoneSchemaField) == 0 {
		return nil, fmt.Errorf("Invalid field format. Got '%s', expected format '%s'", fieldValue, fmt.Sprintf(GlobalLinkTemplate, "{project}", resourceType, "{name}"))
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
		ResourceType: resourceType,
	}, nil
}

func GetProjectFromSchema(projectSchemaField string, d TerraformResourceData, config *transport_tpg.Config) (string, error) {
	res, ok := d.GetOk(projectSchemaField)
	if ok && projectSchemaField != "" {
		return res.(string), nil
	}
	if config.Project != "" {
		return config.Project, nil
	}
	return "", fmt.Errorf("%s: required field is not set", projectSchemaField)
}

func GetUniverseDomainFromSchema(universeSchemaField string, d TerraformResourceData, config *transport_tpg.Config) (string, error) {
	res, ok := d.GetOk(universeSchemaField)
	if ok && universeSchemaField != "" {
		return res.(string), nil
	}
	if config.UniverseDomain != "" {
		return config.UniverseDomain, nil
	}
	if config.UniverseDomain == "" {
		return "googleapis.com", nil
	}
	return "", fmt.Errorf("%s: Error getting the provider field ", universeSchemaField)
}

func GetBillingProjectFromSchema(billingProjectSchemaField string, d TerraformResourceData, config *transport_tpg.Config) (string, error) {
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

	return fmt.Sprintf(OrganizationLinkTemplate, f.OrgId, f.resourceType, f.Name)
}

// Parses an organization field with the following formats:
// - organizations/{my_organizations}/{resource_type}/{resource_name}
func ParseOrganizationFieldValue(resourceType, fieldValue string, isEmptyValid bool) (*OrganizationFieldValue, error) {
	if len(fieldValue) == 0 {
		if isEmptyValid {
			return &OrganizationFieldValue{resourceType: resourceType}, nil
		}
		return nil, fmt.Errorf("The organization field for resource %s cannot be empty", resourceType)
	}

	r := regexp.MustCompile(fmt.Sprintf(OrganizationBasePattern, resourceType))
	if parts := r.FindStringSubmatch(fieldValue); parts != nil {
		return &OrganizationFieldValue{
			OrgId: parts[1],
			Name:  parts[2],

			resourceType: resourceType,
		}, nil
	}

	return nil, fmt.Errorf("Invalid field format. Got '%s', expected format '%s'", fieldValue, fmt.Sprintf(OrganizationLinkTemplate, "{org_id}", resourceType, "{name}"))
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

	return fmt.Sprintf(RegionalLinkTemplate, f.Project, f.Region, f.resourceType, f.Name)
}

// Parses a regional field supporting 5 different formats:
// - https://www.googleapis.com/compute/ANY_VERSION/projects/{my_project}/regions/{region}/{resource_type}/{resource_name}
// - projects/{my_project}/regions/{region}/{resource_type}/{resource_name}
// - regions/{region}/{resource_type}/{resource_name}
// - resource_name
// - "" (empty string). RelativeLink() returns empty if isEmptyValid is true.
//
// If the project is not specified, it first tries to get the project from the `projectSchemaField` and then fallback on the default project.
// If the region is not specified, see function documentation for `GetRegionFromSchema`.
func ParseRegionalFieldValue(resourceType, fieldValue, projectSchemaField, regionSchemaField, zoneSchemaField string, d TerraformResourceData, config *transport_tpg.Config, isEmptyValid bool) (*RegionalFieldValue, error) {
	if len(fieldValue) == 0 {
		if isEmptyValid {
			return &RegionalFieldValue{resourceType: resourceType}, nil
		}
		return nil, fmt.Errorf("The regional field for resource %s cannot be empty.", resourceType)
	}

	r := regexp.MustCompile(fmt.Sprintf(RegionalLinkBasePattern, resourceType))
	if parts := r.FindStringSubmatch(fieldValue); parts != nil {
		return &RegionalFieldValue{
			Project:      parts[1],
			Region:       parts[2],
			Name:         parts[3],
			resourceType: resourceType,
		}, nil
	}

	project, err := GetProjectFromSchema(projectSchemaField, d, config)
	if err != nil {
		return nil, err
	}

	r = regexp.MustCompile(fmt.Sprintf(RegionalPartialLinkBasePattern, resourceType))
	if parts := r.FindStringSubmatch(fieldValue); parts != nil {
		return &RegionalFieldValue{
			Project:      project,
			Region:       parts[1],
			Name:         parts[2],
			resourceType: resourceType,
		}, nil
	}

	region, err := GetRegionFromSchema(regionSchemaField, zoneSchemaField, d, config)
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
func GetRegionFromSchema(regionSchemaField, zoneSchemaField string, d TerraformResourceData, config *transport_tpg.Config) (string, error) {
	// if identical such as GKE location, check if it's a zone first and find
	// the region if so. Otherwise, return as it's a region.
	if regionSchemaField == zoneSchemaField {
		if v, ok := d.GetOk(regionSchemaField); ok {
			if IsZone(v.(string)) {
				return GetRegionFromZone(v.(string)), nil
			}

			return v.(string), nil
		}
	}

	if v, ok := d.GetOk(regionSchemaField); ok && regionSchemaField != "" {
		return GetResourceNameFromSelfLink(v.(string)), nil
	}
	if v, ok := d.GetOk(zoneSchemaField); ok && zoneSchemaField != "" {
		zone := GetResourceNameFromSelfLink(v.(string))
		return GetRegionFromZone(zone), nil
	}
	if config.Region != "" {
		return config.Region, nil
	}
	if config.Zone != "" {
		return GetRegionFromZone(config.Zone), nil
	}

	return "", fmt.Errorf("Cannot determine region: set in this resource, or set provider-level 'region' or 'zone'.")
}

type ProjectFieldValue struct {
	Project string
	Name    string

	ResourceType string
}

func (f ProjectFieldValue) RelativeLink() string {
	if len(f.Name) == 0 {
		return ""
	}

	return fmt.Sprintf(ProjectLinkTemplate, f.Project, f.ResourceType, f.Name)
}

// Parses a project field with the following formats:
// - projects/{my_projects}/{resource_type}/{resource_name}
func ParseProjectFieldValue(resourceType, fieldValue, projectSchemaField string, d TerraformResourceData, config *transport_tpg.Config, isEmptyValid bool) (*ProjectFieldValue, error) {
	if len(fieldValue) == 0 {
		if isEmptyValid {
			return &ProjectFieldValue{ResourceType: resourceType}, nil
		}
		return nil, fmt.Errorf("The project field for resource %s cannot be empty", resourceType)
	}

	r := regexp.MustCompile(fmt.Sprintf(ProjectBasePattern, resourceType))
	if parts := r.FindStringSubmatch(fieldValue); parts != nil {
		return &ProjectFieldValue{
			Project: parts[1],
			Name:    parts[2],

			ResourceType: resourceType,
		}, nil
	}

	project, err := GetProjectFromSchema(projectSchemaField, d, config)
	if err != nil {
		return nil, err
	}

	return &ProjectFieldValue{
		Project: project,
		Name:    GetResourceNameFromSelfLink(fieldValue),

		ResourceType: resourceType,
	}, nil
}

// ExtractFieldByPattern returns the value of a field extracted from a parent field according to the given regular expression pattern.
// An error is returned if the field already has a value different than the value extracted.
func ExtractFieldByPattern(fieldName, fieldValue, parentFieldValue, pattern string) (string, error) {
	var extractedValue string
	// Fetch value from container if the container exists.
	if parentFieldValue != "" {
		r := regexp.MustCompile(pattern)
		m := r.FindStringSubmatch(parentFieldValue)
		if m != nil && len(m) >= 2 {
			extractedValue = m[1]
		} else if fieldValue == "" {
			// The pattern didn't match and the value doesn't exist.
			return "", fmt.Errorf("parent of %q has no matching values from pattern %q in value %q", fieldName, pattern, parentFieldValue)
		}
	}

	// If both values exist and are different, error
	if fieldValue != "" && extractedValue != "" && fieldValue != extractedValue {
		return "", fmt.Errorf("%q has conflicting values of %q (from parent) and %q (from self)", fieldName, extractedValue, fieldValue)
	}

	// If value does not exist, use the value in container.
	if fieldValue == "" {
		return extractedValue, nil
	}

	return fieldValue, nil
}
