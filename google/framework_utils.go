// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-provider-google/google/fwmodels"
	"github.com/hashicorp/terraform-provider-google/google/fwresource"
	"github.com/hashicorp/terraform-provider-google/google/fwtransport"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

func CompileUserAgentString(ctx context.Context, name, tfVersion, provVersion string) string {
	return fwtransport.CompileUserAgentString(ctx, name, tfVersion, provVersion)
}

func GetCurrentUserEmailFramework(p *fwtransport.FrameworkProviderConfig, userAgent string, diags *diag.Diagnostics) string {
	return fwtransport.GetCurrentUserEmailFramework(p, userAgent, diags)
}

func generateFrameworkUserAgentString(metaData *fwmodels.ProviderMetaModel, currUserAgent string) string {
	return fwtransport.GenerateFrameworkUserAgentString(metaData, currUserAgent)
}

// GetProject reads the "project" field from the given resource and falls
// back to the provider's value if not given. If the provider's value is not
// given, an error is returned.
func getProjectFramework(rVal, pVal types.String, diags *diag.Diagnostics) types.String {
	return fwresource.GetProjectFramework(rVal, pVal, diags)
}

func handleDatasourceNotFoundError(ctx context.Context, err error, state *tfsdk.State, resource string, diags *diag.Diagnostics) {
	fwtransport.HandleDatasourceNotFoundError(ctx, err, state, resource, diags)
}

// field helpers

// Parses a project field with the following formats:
// - projects/{my_projects}/{resource_type}/{resource_name}
func parseProjectFieldValueFramework(resourceType, fieldValue, projectSchemaField string, rVal, pVal types.String, isEmptyValid bool, diags *diag.Diagnostics) *tpgresource.ProjectFieldValue {
	return fwresource.ParseProjectFieldValueFramework(resourceType, fieldValue, projectSchemaField, rVal, pVal, isEmptyValid, diags)
}
