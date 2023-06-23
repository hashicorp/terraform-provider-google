// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package fwtransport

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-provider-google/google/fwmodels"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

const uaEnvVar = "TF_APPEND_USER_AGENT"

func CompileUserAgentString(ctx context.Context, name, tfVersion, provVersion string) string {
	ua := fmt.Sprintf("Terraform/%s (+https://www.terraform.io) Terraform-Plugin-SDK/%s %s/%s", tfVersion, "terraform-plugin-framework", name, provVersion)

	if add := os.Getenv(uaEnvVar); add != "" {
		add = strings.TrimSpace(add)
		if len(add) > 0 {
			ua += " " + add
			tflog.Debug(ctx, fmt.Sprintf("Using modified User-Agent: %s", ua))
		}
	}

	return ua
}

func GetCurrentUserEmailFramework(p *FrameworkProviderConfig, userAgent string, diags *diag.Diagnostics) string {
	// When environment variables UserProjectOverride and BillingProject are set for the provider,
	// the header X-Goog-User-Project is set for the API requests.
	// But it causes an error when calling GetCurrUserEmail. Set the project to be "NO_BILLING_PROJECT_OVERRIDE".
	// And then it triggers the header X-Goog-User-Project to be set to empty string.

	// See https://github.com/golang/oauth2/issues/306 for a recommendation to do this from a Go maintainer
	// URL retrieved from https://accounts.google.com/.well-known/openid-configuration
	res, d := SendFrameworkRequest(p, "GET", "NO_BILLING_PROJECT_OVERRIDE", "https://openidconnect.googleapis.com/v1/userinfo", userAgent, nil)
	diags.Append(d...)

	if diags.HasError() {
		tflog.Info(p.Context, "error retrieving userinfo for your provider credentials. have you enabled the 'https://www.googleapis.com/auth/userinfo.email' scope?")
		return ""
	}
	if res["email"] == nil {
		diags.AddError("error retrieving email from userinfo.", "email was nil in the response.")
		return ""
	}
	return res["email"].(string)
}

func GenerateFrameworkUserAgentString(metaData *fwmodels.ProviderMetaModel, currUserAgent string) string {
	if metaData != nil && !metaData.ModuleName.IsNull() && metaData.ModuleName.ValueString() != "" {
		return strings.Join([]string{currUserAgent, metaData.ModuleName.ValueString()}, " ")
	}

	return currUserAgent
}

func HandleDatasourceNotFoundError(ctx context.Context, err error, state *tfsdk.State, resource string, diags *diag.Diagnostics) {
	if transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
		tflog.Warn(ctx, fmt.Sprintf("Removing %s because it's gone", resource))
		// The resource doesn't exist anymore
		state.RemoveResource(ctx)
	}

	diags.AddError(fmt.Sprintf("Error when reading or editing %s", resource), err.Error())
}
