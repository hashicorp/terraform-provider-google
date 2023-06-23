// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package fwresource

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

// GetProject reads the "project" field from the given resource and falls
// back to the provider's value if not given. If the provider's value is not
// given, an error is returned.
func GetProjectFramework(rVal, pVal types.String, diags *diag.Diagnostics) types.String {
	return getProjectFromFrameworkSchema("project", rVal, pVal, diags)
}

func getProjectFromFrameworkSchema(projectSchemaField string, rVal, pVal types.String, diags *diag.Diagnostics) types.String {
	if !rVal.IsNull() && rVal.ValueString() != "" {
		return rVal
	}

	if !pVal.IsNull() && pVal.ValueString() != "" {
		return pVal
	}

	diags.AddError("required field is not set", fmt.Sprintf("%s is not set", projectSchemaField))
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

	project := getProjectFromFrameworkSchema(projectSchemaField, rVal, pVal, diags)
	if diags.HasError() {
		return nil
	}

	return &tpgresource.ProjectFieldValue{
		Project: project.ValueString(),
		Name:    tpgresource.GetResourceNameFromSelfLink(fieldValue),

		ResourceType: resourceType,
	}
}
