// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package google

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the data source satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &GoogleClientOpenIDUserinfoDataSource{}
	_ datasource.DataSourceWithConfigure = &GoogleClientOpenIDUserinfoDataSource{}
)

func NewGoogleClientOpenIDUserinfoDataSource() datasource.DataSource {
	return &GoogleClientOpenIDUserinfoDataSource{}
}

type GoogleClientOpenIDUserinfoDataSource struct {
	providerConfig *frameworkProvider
}

type GoogleClientOpenIDUserinfoModel struct {
	// Id could/should be removed in future as it's not necessary in the plugin framework
	// https://github.com/hashicorp/terraform-plugin-testing/issues/84
	Id    types.String `tfsdk:"id"`
	Email types.String `tfsdk:"email"`
}

func (d *GoogleClientOpenIDUserinfoDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_client_openid_userinfo"
}

func (d *GoogleClientOpenIDUserinfoDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `Get OpenID userinfo about the credentials used with the Google provider, specifically the email.
This datasource enables you to export the email of the account you've authenticated the provider with; this can be used alongside data.google_client_config's access_token to perform OpenID Connect authentication with GKE and configure an RBAC role for the email used.

Note: This resource will only work as expected if the provider is configured to use the https://www.googleapis.com/auth/userinfo.email scope! You will receive an error otherwise. The provider uses this scope by default.`,
		MarkdownDescription: `Get OpenID userinfo about the credentials used with the Google provider, specifically the email.
This datasource enables you to export the email of the account you've authenticated the provider with; this can be used alongside data.google_client_config's access_token to perform OpenID Connect authentication with GKE and configure an RBAC role for the email used.

~> This resource will only work as expected if the provider is configured to use the https://www.googleapis.com/auth/userinfo.email scope! You will receive an error otherwise. The provider uses this scope by default.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "The ID of this data source in Terraform state. Its value is the same as the email attribute. Do not use this field, use the email attribute instead.",
				MarkdownDescription: "The ID of this data source in Terraform state. Its value is the same as the `email` attribute. Do not use this field, use the `email` attribute instead.",
				Computed:            true,
			},
			"email": schema.StringAttribute{
				Description:         "The email of the account used by the provider to authenticate with GCP.",
				MarkdownDescription: "The email of the account used by the provider to authenticate with GCP.",
				Computed:            true,
			},
		},
	}
}

func (d *GoogleClientOpenIDUserinfoDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	p, ok := req.ProviderData.(*frameworkProvider)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *frameworkProvider, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	// Required for accessing userAgent and passing as an argument into a util function
	d.providerConfig = p
}

func (d *GoogleClientOpenIDUserinfoDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data GoogleClientOpenIDUserinfoModel
	var metaData *ProviderMetaModel
	var diags diag.Diagnostics

	// Read Provider meta into the meta model
	resp.Diagnostics.Append(req.ProviderMeta.Get(ctx, &metaData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	userAgent := generateFrameworkUserAgentString(metaData, d.providerConfig.userAgent)
	email := GetCurrentUserEmailFramework(d.providerConfig, userAgent, &diags)

	data.Email = types.StringValue(email)
	data.Id = types.StringValue(email)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
