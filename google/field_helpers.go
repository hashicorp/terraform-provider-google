// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

const (
	// Deprecated: For backward compatibility globalLinkTemplate is still working,
	// but all new code should use GlobalLinkTemplate in the tpgresource package instead.
	globalLinkTemplate = tpgresource.GlobalLinkTemplate
	// Deprecated: For backward compatibility globalLinkBasePattern is still working,
	// but all new code should use GlobalLinkBasePattern in the tpgresource package instead.
	globalLinkBasePattern = tpgresource.GlobalLinkBasePattern
	// Deprecated: For backward compatibility zonalLinkTemplate is still working,
	// but all new code should use ZonalLinkTemplate in the tpgresource package instead.
	zonalLinkTemplate = tpgresource.ZonalLinkTemplate
	// Deprecated: For backward compatibility zonalLinkBasePattern is still working,
	// but all new code should use ZonalLinkBasePattern in the tpgresource package instead.
	zonalLinkBasePattern = tpgresource.ZonalLinkBasePattern
	// Deprecated: For backward compatibility zonalPartialLinkBasePattern is still working,
	// but all new code should use ZonalPartialLinkBasePattern in the tpgresource package instead.
	zonalPartialLinkBasePattern = tpgresource.ZonalPartialLinkBasePattern
	// Deprecated: For backward compatibility regionalLinkTemplate is still working,
	// but all new code should use RegionalLinkTemplate in the tpgresource package instead.
	regionalLinkTemplate = tpgresource.RegionalLinkTemplate
	// Deprecated: For backward compatibility regionalLinkBasePattern is still working,
	// but all new code should use RegionalLinkBasePattern in the tpgresource package instead.
	regionalLinkBasePattern = tpgresource.RegionalLinkBasePattern
	// Deprecated: For backward compatibility regionalPartialLinkBasePattern is still working,
	// but all new code should use RegionalPartialLinkBasePattern in the tpgresource package instead.
	regionalPartialLinkBasePattern = tpgresource.RegionalPartialLinkBasePattern
	// Deprecated: For backward compatibility projectLinkTemplate is still working,
	// but all new code should use ProjectLinkTemplate in the tpgresource package instead.
	projectLinkTemplate = tpgresource.ProjectLinkTemplate
	// Deprecated: For backward compatibility projectBasePattern is still working,
	// but all new code should use ProjectBasePattern in the tpgresource package instead.
	projectBasePattern = tpgresource.ProjectBasePattern
	// Deprecated: For backward compatibility organizationLinkTemplate is still working,
	// but all new code should use OrganizationLinkTemplate in the tpgresource package instead.
	organizationLinkTemplate = tpgresource.OrganizationLinkTemplate
	// Deprecated: For backward compatibility organizationBasePattern is still working,
	// but all new code should use OrganizationBasePattern in the tpgresource package instead.
	organizationBasePattern = tpgresource.OrganizationBasePattern
)

// ------------------------------------------------------------
// Field helpers
// ------------------------------------------------------------

// Deprecated: For backward compatibility ParseNetworkFieldValue is still working,
// but all new code should use ParseNetworkFieldValue in the tpgresource package instead.
func ParseNetworkFieldValue(network string, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (*tpgresource.GlobalFieldValue, error) {
	return tpgresource.ParseNetworkFieldValue(network, d, config)
}

// Deprecated: For backward compatibility ParseSubnetworkFieldValue is still working,
// but all new code should use ParseSubnetworkFieldValue in the tpgresource package instead.
func ParseSubnetworkFieldValue(subnetwork string, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (*tpgresource.RegionalFieldValue, error) {
	return tpgresource.ParseSubnetworkFieldValue(subnetwork, d, config)
}

// Deprecated: For backward compatibility ParseSubnetworkFieldValueWithProjectField is still working,
// but all new code should use ParseSubnetworkFieldValueWithProjectField in the tpgresource package instead.
func ParseSubnetworkFieldValueWithProjectField(subnetwork, projectField string, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (*tpgresource.RegionalFieldValue, error) {
	return tpgresource.ParseSubnetworkFieldValueWithProjectField(subnetwork, projectField, d, config)
}

// Deprecated: For backward compatibility ParseSslCertificateFieldValue is still working,
// but all new code should use ParseSslCertificateFieldValue in the tpgresource package instead.
func ParseSslCertificateFieldValue(sslCertificate string, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (*tpgresource.GlobalFieldValue, error) {
	return tpgresource.ParseSslCertificateFieldValue(sslCertificate, d, config)
}

// Deprecated: For backward compatibility ParseHttpHealthCheckFieldValue is still working,
// but all new code should use ParseHttpHealthCheckFieldValue in the tpgresource package instead.
func ParseHttpHealthCheckFieldValue(healthCheck string, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (*tpgresource.GlobalFieldValue, error) {
	return tpgresource.ParseHttpHealthCheckFieldValue(healthCheck, d, config)
}

// Deprecated: For backward compatibility ParseDiskFieldValue is still working,
// but all new code should use ParseDiskFieldValue in the tpgresource package instead.
func ParseDiskFieldValue(disk string, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (*tpgresource.ZonalFieldValue, error) {
	return tpgresource.ParseDiskFieldValue(disk, d, config)
}

// Deprecated: For backward compatibility ParseRegionDiskFieldValue is still working,
// but all new code should use ParseRegionDiskFieldValue in the tpgresource package instead.
func ParseRegionDiskFieldValue(disk string, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (*tpgresource.RegionalFieldValue, error) {
	return tpgresource.ParseRegionDiskFieldValue(disk, d, config)
}

// Deprecated: For backward compatibility ParseOrganizationCustomRoleName is still working,
// but all new code should use ParseOrganizationCustomRoleName in the tpgresource package instead.
func ParseOrganizationCustomRoleName(role string) (*tpgresource.OrganizationFieldValue, error) {
	return tpgresource.ParseOrganizationCustomRoleName(role)
}

// Deprecated: For backward compatibility ParseAcceleratorFieldValue is still working,
// but all new code should use ParseAcceleratorFieldValue in the tpgresource package instead.
func ParseAcceleratorFieldValue(accelerator string, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (*tpgresource.ZonalFieldValue, error) {
	return tpgresource.ParseAcceleratorFieldValue(accelerator, d, config)
}

// Deprecated: For backward compatibility ParseMachineTypesFieldValue is still working,
// but all new code should use ParseMachineTypesFieldValue in the tpgresource package instead.
func ParseMachineTypesFieldValue(machineType string, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (*tpgresource.ZonalFieldValue, error) {
	return tpgresource.ParseMachineTypesFieldValue(machineType, d, config)
}

// Deprecated: For backward compatibility ParseInstanceFieldValue is still working,
// but all new code should use ParseInstanceFieldValue in the tpgresource package instead.
func ParseInstanceFieldValue(instance string, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (*tpgresource.ZonalFieldValue, error) {
	return tpgresource.ParseInstanceFieldValue(instance, d, config)
}

// Deprecated: For backward compatibility ParseInstanceGroupFieldValue is still working,
// but all new code should use ParseInstanceGroupFieldValue in the tpgresource package instead.
func ParseInstanceGroupFieldValue(instanceGroup string, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (*tpgresource.ZonalFieldValue, error) {
	return tpgresource.ParseInstanceGroupFieldValue(instanceGroup, d, config)
}

// Deprecated: For backward compatibility ParseInstanceTemplateFieldValue is still working,
// but all new code should use ParseInstanceTemplateFieldValue in the tpgresource package instead.
func ParseInstanceTemplateFieldValue(instanceTemplate string, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (*tpgresource.GlobalFieldValue, error) {
	return tpgresource.ParseInstanceTemplateFieldValue(instanceTemplate, d, config)
}

// Deprecated: For backward compatibility ParseMachineImageFieldValue is still working,
// but all new code should use ParseMachineImageFieldValue in the tpgresource package instead.
func ParseMachineImageFieldValue(machineImage string, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (*tpgresource.GlobalFieldValue, error) {
	return tpgresource.ParseMachineImageFieldValue(machineImage, d, config)
}

// Deprecated: For backward compatibility ParseSecurityPolicyFieldValue is still working,
// but all new code should use ParseSecurityPolicyFieldValue in the tpgresource package instead.
func ParseSecurityPolicyFieldValue(securityPolicy string, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (*tpgresource.GlobalFieldValue, error) {
	return tpgresource.ParseSecurityPolicyFieldValue(securityPolicy, d, config)
}

// Deprecated: For backward compatibility ParseNetworkEndpointGroupFieldValue is still working,
// but all new code should use ParseNetworkEndpointGroupFieldValue in the tpgresource package instead.
func ParseNetworkEndpointGroupFieldValue(networkEndpointGroup string, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (*tpgresource.ZonalFieldValue, error) {
	return tpgresource.ParseNetworkEndpointGroupFieldValue(networkEndpointGroup, d, config)
}

// Deprecated: For backward compatibility ParseNetworkEndpointGroupRegionalFieldValue is still working,
// but all new code should use ParseNetworkEndpointGroupRegionalFieldValue in the tpgresource package instead.
func ParseNetworkEndpointGroupRegionalFieldValue(networkEndpointGroup string, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (*tpgresource.RegionalFieldValue, error) {
	return tpgresource.ParseNetworkEndpointGroupRegionalFieldValue(networkEndpointGroup, d, config)
}

// ------------------------------------------------------------
// Base helpers used to create helpers for specific fields.
// ------------------------------------------------------------

// Parses a global field supporting 5 different formats:
// - https://www.googleapis.com/compute/ANY_VERSION/projects/{my_project}/global/{resource_type}/{resource_name}
// - projects/{my_project}/global/{resource_type}/{resource_name}
// - global/{resource_type}/{resource_name}
// - resource_name
// - "" (empty string). RelativeLink() returns empty if isEmptyValid is true.
//
// If the project is not specified, it first tries to get the project from the `projectSchemaField` and then fallback on the default project.
//
// Deprecated: For backward compatibility parseGlobalFieldValue is still working,
// but all new code should use ParseGlobalFieldValue in the tpgresource package instead.
func parseGlobalFieldValue(resourceType, fieldValue, projectSchemaField string, d tpgresource.TerraformResourceData, config *transport_tpg.Config, isEmptyValid bool) (*tpgresource.GlobalFieldValue, error) {
	return tpgresource.ParseGlobalFieldValue(resourceType, fieldValue, projectSchemaField, d, config, isEmptyValid)
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
//
// Deprecated: For backward compatibility parseZonalFieldValue is still working,
// but all new code should use ParseZonalFieldValue in the tpgresource package instead.
func parseZonalFieldValue(resourceType, fieldValue, projectSchemaField, zoneSchemaField string, d tpgresource.TerraformResourceData, config *transport_tpg.Config, isEmptyValid bool) (*tpgresource.ZonalFieldValue, error) {
	return tpgresource.ParseZonalFieldValue(resourceType, fieldValue, projectSchemaField, zoneSchemaField, d, config, isEmptyValid)
}

// Deprecated: For backward compatibility getProjectFromSchema is still working,
// but all new code should use GetProjectFromSchema in the tpgresource package instead.
func getProjectFromSchema(projectSchemaField string, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (string, error) {
	return tpgresource.GetProjectFromSchema(projectSchemaField, d, config)
}

// Deprecated: For backward compatibility getBillingProjectFromSchema is still working,
// but all new code should use GetBillingProjectFromSchema in the tpgresource package instead.
func getBillingProjectFromSchema(billingProjectSchemaField string, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (string, error) {
	return tpgresource.GetBillingProjectFromSchema(billingProjectSchemaField, d, config)
}

// Parses an organization field with the following formats:
// - organizations/{my_organizations}/{resource_type}/{resource_name}
//
// Deprecated: For backward compatibility parseOrganizationFieldValue is still working,
// but all new code should use ParseOrganizationFieldValue in the tpgresource package instead.
func parseOrganizationFieldValue(resourceType, fieldValue string, isEmptyValid bool) (*tpgresource.OrganizationFieldValue, error) {
	return tpgresource.ParseOrganizationFieldValue(resourceType, fieldValue, isEmptyValid)
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
//
// Deprecated: For backward compatibility parseRegionalFieldValue is still working,
// but all new code should use ParseRegionalFieldValue in the tpgresource package instead.
func parseRegionalFieldValue(resourceType, fieldValue, projectSchemaField, regionSchemaField, zoneSchemaField string, d tpgresource.TerraformResourceData, config *transport_tpg.Config, isEmptyValid bool) (*tpgresource.RegionalFieldValue, error) {
	return tpgresource.ParseRegionalFieldValue(resourceType, fieldValue, projectSchemaField, regionSchemaField, zoneSchemaField, d, config, isEmptyValid)
}

// Infers the region based on the following (in order of priority):
// - `regionSchemaField` in resource schema
// - region extracted from the `zoneSchemaField` in resource schema
// - provider-level region
// - region extracted from the provider-level zone
//
// Deprecated: For backward compatibility getRegionFromSchema is still working,
// but all new code should use GetRegionFromSchema in the tpgresource package instead.
func getRegionFromSchema(regionSchemaField, zoneSchemaField string, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (string, error) {
	return tpgresource.GetRegionFromSchema(regionSchemaField, zoneSchemaField, d, config)
}

// Parses a project field with the following formats:
// - projects/{my_projects}/{resource_type}/{resource_name}
//
// Deprecated: For backward compatibility parseProjectFieldValue is still working,
// but all new code should use ParseProjectFieldValue in the tpgresource package instead.
func parseProjectFieldValue(resourceType, fieldValue, projectSchemaField string, d tpgresource.TerraformResourceData, config *transport_tpg.Config, isEmptyValid bool) (*tpgresource.ProjectFieldValue, error) {
	return tpgresource.ParseProjectFieldValue(resourceType, fieldValue, projectSchemaField, d, config, isEmptyValid)
}
