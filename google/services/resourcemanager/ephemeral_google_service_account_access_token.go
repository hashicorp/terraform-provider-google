// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-google/google/fwtransport"
	"github.com/hashicorp/terraform-provider-google/google/fwutils"
	"github.com/hashicorp/terraform-provider-google/google/fwvalidators"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	"google.golang.org/api/iamcredentials/v1"
)

var _ ephemeral.EphemeralResource = &googleEphemeralServiceAccountAccessToken{}

func GoogleEphemeralServiceAccountAccessToken() ephemeral.EphemeralResource {
	return &googleEphemeralServiceAccountAccessToken{}
}

type googleEphemeralServiceAccountAccessToken struct {
	providerConfig *fwtransport.FrameworkProviderConfig
}

func (p *googleEphemeralServiceAccountAccessToken) Metadata(ctx context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_account_access_token"
}

type ephemeralServiceAccountAccessTokenModel struct {
	TargetServiceAccount types.String `tfsdk:"target_service_account"`
	AccessToken          types.String `tfsdk:"access_token"`
	Scopes               types.Set    `tfsdk:"scopes"`
	Delegates            types.Set    `tfsdk:"delegates"`
	Lifetime             types.String `tfsdk:"lifetime"`
}

func (p *googleEphemeralServiceAccountAccessToken) Schema(ctx context.Context, req ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema.Description = "This ephemeral resource provides a google oauth2 access_token for a different service account than the one initially running the script."
	resp.Schema.MarkdownDescription = "This ephemeral resource provides a google oauth2 access_token for a different service account than the one initially running the script."

	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"target_service_account": schema.StringAttribute{
				Description: "The service account to impersonate (e.g. `service_B@your-project-id.iam.gserviceaccount.com`)",
				Required:    true,
				Validators: []validator.String{
					fwvalidators.ServiceAccountEmailValidator{},
				},
			},
			"access_token": schema.StringAttribute{
				Description: "The `access_token` representing the new generated identity.",
				Sensitive:   true,
				Computed:    true,
			},
			"lifetime": schema.StringAttribute{
				Description: "Lifetime of the impersonated token (defaults to its max: `3600s`)",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					fwvalidators.BoundedDuration{
						MinDuration: 0,
						MaxDuration: 3600 * time.Second,
					},
				},
			},
			"scopes": schema.SetAttribute{
				Description: "The scopes the new credential should have (e.g. `['cloud-platform']`)",
				Required:    true,
				ElementType: types.StringType,
			},
			"delegates": schema.SetAttribute{
				Description: "Delegate chain of approvals needed to perform full impersonation. Specify the fully qualified service account name.  (e.g. `['projects/-/serviceAccounts/delegate-svc-account@project-id.iam.gserviceaccount.com']`)",
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(fwvalidators.ServiceAccountEmailValidator{}),
				},
			},
		},
	}
}

func (p *googleEphemeralServiceAccountAccessToken) Configure(ctx context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	pd, ok := req.ProviderData.(*fwtransport.FrameworkProviderConfig)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *fwtransport.FrameworkProviderConfig, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	// Required for accessing userAgent and passing as an argument into a util function
	p.providerConfig = pd
}

func (p *googleEphemeralServiceAccountAccessToken) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data ephemeralServiceAccountAccessTokenModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// This is the default value for the lifetime of the access token
	// Both ephemeral resources and data sources do not allow you to set a value for this attribute in the schema
	if data.Lifetime.IsNull() {
		data.Lifetime = types.StringValue("3600s")
	}

	service := p.providerConfig.NewIamCredentialsClient(p.providerConfig.UserAgent)
	name := fmt.Sprintf("projects/-/serviceAccounts/%s", data.TargetServiceAccount.ValueString())

	ScopesSetValue, diags := data.Scopes.ToSetValue(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var delegates []string
	if !data.Delegates.IsNull() {
		delegates = fwutils.StringSet(data.Delegates)
	}

	tokenRequest := &iamcredentials.GenerateAccessTokenRequest{
		Lifetime:  data.Lifetime.ValueString(),
		Delegates: delegates,
		Scope:     tpgresource.CanonicalizeServiceScopes(fwutils.StringSet(ScopesSetValue)),
	}

	at, err := service.Projects.ServiceAccounts.GenerateAccessToken(name, tokenRequest).Do()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error generating access token",
			fmt.Sprintf("Error generating access token: %s", err),
		)
		return
	}

	data.AccessToken = types.StringValue(at.AccessToken)
	resp.Diagnostics.Append(resp.Result.Set(ctx, data)...)
}
