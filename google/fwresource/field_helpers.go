// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package fwresource

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// GetProject reads the "project" field from the given resource and falls
// back to the provider's value if not given. If the provider's value is not
// given, an error is returned.
func GetProjectFramework(rVal, pVal types.String, diags *diag.Diagnostics) types.String {
	return getProviderDefaultFromFrameworkSchema("project", rVal, pVal, diags)
}

func GetRegionFramework(rVal, pVal types.String, diags *diag.Diagnostics) types.String {
	return getProviderDefaultFromFrameworkSchema("region", rVal, pVal, diags)
}

func GetZoneFramework(rVal, pVal types.String, diags *diag.Diagnostics) types.String {
	return getProviderDefaultFromFrameworkSchema("zone", rVal, pVal, diags)
}

func getProviderDefaultFromFrameworkSchema(schemaField string, rVal, pVal types.String, diags *diag.Diagnostics) types.String {
	if !rVal.IsNull() && rVal.ValueString() != "" {
		return rVal
	}

	if !pVal.IsNull() && pVal.ValueString() != "" {
		return pVal
	}

	diags.AddError("required field is not set", fmt.Sprintf("%s is not set", schemaField))
	return types.String{}
}

// Parses a project field with the following formats:
// - projects/{my_projects}/{resource_type}/{resource_name}
func ParseProjectFieldValueFramework(resourceType, fieldValue, projectSchemaField string, rVal, pVal types.String, isEmptyValid bool, diags *diag.Diagnostics) *tpgresource.ProjectFieldValue {
	if len(fieldValue) == 0 {
		if isEmptyValid {
			return &tpgresource.ProjectFieldValue{ResourceType: resourceType}
		}
		diags.AddError("field can not be empty", fmt.Sprintf("The project field for resource %s cannot be empty", resourceType))
		return nil
	}

	r := regexp.MustCompile(fmt.Sprintf(tpgresource.ProjectBasePattern, resourceType))
	if parts := r.FindStringSubmatch(fieldValue); parts != nil {
		return &tpgresource.ProjectFieldValue{
			Project: parts[1],
			Name:    parts[2],

			ResourceType: resourceType,
		}
	}

	project := getProviderDefaultFromFrameworkSchema(projectSchemaField, rVal, pVal, diags)
	if diags.HasError() {
		return nil
	}

	return &tpgresource.ProjectFieldValue{
		Project: project.ValueString(),
		Name:    tpgresource.GetResourceNameFromSelfLink(fieldValue),

		ResourceType: resourceType,
	}
}

// This function isn't a test of transport.go; instead, it is used as an alternative
// to ReplaceVars inside tests.
func ReplaceVarsForFrameworkTest(prov *transport_tpg.Config, rs *terraform.ResourceState, linkTmpl string) (string, error) {
	re := regexp.MustCompile("{{([[:word:]]+)}}")
	var project, region, zone string

	if strings.Contains(linkTmpl, "{{project}}") {
		project = rs.Primary.Attributes["project"]
	}

	if strings.Contains(linkTmpl, "{{region}}") {
		region = tpgresource.GetResourceNameFromSelfLink(rs.Primary.Attributes["region"])
	}

	if strings.Contains(linkTmpl, "{{zone}}") {
		zone = tpgresource.GetResourceNameFromSelfLink(rs.Primary.Attributes["zone"])
	}

	replaceFunc := func(s string) string {
		m := re.FindStringSubmatch(s)[1]
		if m == "project" {
			return project
		}
		if m == "region" {
			return region
		}
		if m == "zone" {
			return zone
		}

		if v, ok := rs.Primary.Attributes[m]; ok {
			return v
		}

		// Attempt to draw values from the provider
		if f := reflect.Indirect(reflect.ValueOf(prov)).FieldByName(m); f.IsValid() {
			return f.String()
		}

		return ""
	}

	return re.ReplaceAllStringFunc(linkTmpl, replaceFunc), nil
}
