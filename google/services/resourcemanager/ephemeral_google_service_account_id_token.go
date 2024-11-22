// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager

import (
	"context"
	"fmt"

	"google.golang.org/api/idtoken"
	"google.golang.org/api/option"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-google/google/fwutils"
	"github.com/hashicorp/terraform-provider-google/google/fwvalidators"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/iamcredentials/v1"
)

var _ ephemeral.EphemeralResource = &googleEphemeralServiceAccountIdToken{}

func GoogleEphemeralServiceAccountIdToken() ephemeral.EphemeralResource {
	return &googleEphemeralServiceAccountIdToken{}
}

type googleEphemeralServiceAccountIdToken struct {
	providerConfig *transport_tpg.Config
}

func (p *googleEphemeralServiceAccountIdToken) Metadata(ctx context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_account_id_token"
}

type ephemeralServiceAccountIdTokenModel struct {
	TargetAudience       types.String `tfsdk:"target_audience"`
	TargetServiceAccount types.String `tfsdk:"target_service_account"`
	Delegates            types.Set    `tfsdk:"delegates"`
	IncludeEmail         types.Bool   `tfsdk:"include_email"`
	IdToken              types.String `tfsdk:"id_token"`
}

func (p *googleEphemeralServiceAccountIdToken) Schema(ctx context.Context, req ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema.Description = "This ephemeral resource provides a Google OpenID Connect (oidc) id_token."
	resp.Schema.MarkdownDescription = "This ephemeral resource provides a Google OpenID Connect (oidc) id_token."

	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"target_audience": schema.StringAttribute{
				Description: "The audience claim for the `id_token`.",
				Required:    true,
			},
			"target_service_account": schema.StringAttribute{
				Description: "The email of the service account being impersonated.  Used only when using impersonation mode.",
				Optional:    true,
				Validators: []validator.String{
					fwvalidators.ServiceAccountEmailValidator{},
				},
			},
			"delegates": schema.SetAttribute{
				Description: "Delegate chain of approvals needed to perform full impersonation. Specify the fully qualified service account name.  Used only when using impersonation mode.",
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(fwvalidators.ServiceAccountEmailValidator{}),
				},
			},
			"include_email": schema.BoolAttribute{
				Description: "Include the verified email in the claim. Used only when using impersonation mode.",
				Optional:    true, // Defaults to false when not set (Null / Unknown)
			},
			"id_token": schema.StringAttribute{
				Description: "The `id_token` representing the new generated identity.",
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *googleEphemeralServiceAccountIdToken) Configure(ctx context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	pd, ok := req.ProviderData.(*transport_tpg.Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *transport_tpg.Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	p.providerConfig = pd
}

func (p *googleEphemeralServiceAccountIdToken) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data ephemeralServiceAccountIdTokenModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	targetAudience := data.TargetAudience.ValueString()

	creds, err := p.providerConfig.GetCredentials([]string{userInfoScope}, false)
	if err != nil {
		resp.Diagnostics.AddError(
			"error calling GetCredentials()",
			err.Error(),
		)
		return
	}

	targetServiceAccount := data.TargetServiceAccount
	// If a target service account is provided, use the API to generate the idToken
	if !targetServiceAccount.IsNull() && !targetServiceAccount.IsUnknown() {
		service := p.providerConfig.NewIamCredentialsClient(p.providerConfig.UserAgent)
		name := fmt.Sprintf("projects/-/serviceAccounts/%s", targetServiceAccount.ValueString())

		tokenRequest := &iamcredentials.GenerateIdTokenRequest{
			Audience:     targetAudience,
			IncludeEmail: data.IncludeEmail.ValueBool(),
			Delegates:    fwutils.StringSet(data.Delegates),
		}
		at, err := service.Projects.ServiceAccounts.GenerateIdToken(name, tokenRequest).Do()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error calling iamcredentials.GenerateIdToken",
				err.Error(),
			)
			return
		}

		data.IdToken = types.StringValue(at.Token)
		resp.Diagnostics.Append(resp.Result.Set(ctx, data)...)
		return
	}

	// If no target service account, use the default credentials
	ctx = context.Background()
	co := []option.ClientOption{}
	if creds.JSON != nil {
		co = append(co, idtoken.WithCredentialsJSON(creds.JSON))
	}

	idTokenSource, err := idtoken.NewTokenSource(ctx, targetAudience, co...)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to retrieve TokenSource",
			err.Error(),
		)
		return
	}
	idToken, err := idTokenSource.Token()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to retrieve Token",
			err.Error(),
		)
		return
	}

	data.IdToken = types.StringValue(idToken.AccessToken)
	resp.Diagnostics.Append(resp.Result.Set(ctx, data)...)
}
