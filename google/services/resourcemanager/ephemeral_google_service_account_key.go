// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-google/google/fwtransport"
	"github.com/hashicorp/terraform-provider-google/google/verify"
)

var _ ephemeral.EphemeralResource = &googleEphemeralServiceAccountKey{}

func GoogleEphemeralServiceAccountKey() ephemeral.EphemeralResource {
	return &googleEphemeralServiceAccountKey{}
}

type googleEphemeralServiceAccountKey struct {
	providerConfig *fwtransport.FrameworkProviderConfig
}

func (p *googleEphemeralServiceAccountKey) Metadata(ctx context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_account_key"
}

type ephemeralServiceAccountKeyModel struct {
	Name          types.String `tfsdk:"name"`
	PublicKeyType types.String `tfsdk:"public_key_type"`
	Project       types.String `tfsdk:"project"`
	KeyAlgorithm  types.String `tfsdk:"key_algorithm"`
	PublicKey     types.String `tfsdk:"public_key"`
}

func (p *googleEphemeralServiceAccountKey) Schema(ctx context.Context, req ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(verify.ServiceAccountKeyNameRegex),
						"must match regex: "+verify.ServiceAccountKeyNameRegex,
					),
				}},
			"public_key_type": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"TYPE_X509_PEM_FILE",
						"TYPE_RAW",
					),
				},
			},
			"project": schema.StringAttribute{
				Optional: true,
			},
			"key_algorithm": schema.StringAttribute{
				Computed: true,
			},
			"public_key": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (p *googleEphemeralServiceAccountKey) Configure(ctx context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
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

	p.providerConfig = pd
}

func (p *googleEphemeralServiceAccountKey) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data ephemeralServiceAccountKeyModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	keyName := data.Name.ValueString()

	// Validate name
	r := regexp.MustCompile(verify.ServiceAccountKeyNameRegex)
	if !r.MatchString(keyName) {
		resp.Diagnostics.AddError(
			"Invalid key name",
			fmt.Sprintf("Invalid key name %q does not match regexp %q", keyName, verify.ServiceAccountKeyNameRegex),
		)
		return
	}

	publicKeyType := data.PublicKeyType.ValueString()
	if publicKeyType == "" {
		publicKeyType = "TYPE_X509_PEM_FILE"
	}

	sak, err := p.providerConfig.NewIamClient(p.providerConfig.UserAgent).Projects.ServiceAccounts.Keys.Get(keyName).PublicKeyType(publicKeyType).Do()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving Service Account Key",
			fmt.Sprintf("Error retrieving Service Account Key %q: %s", keyName, err),
		)
		return
	}

	data.Name = types.StringValue(sak.Name)
	data.KeyAlgorithm = types.StringValue(sak.KeyAlgorithm)
	data.PublicKey = types.StringValue(sak.PublicKeyData)

	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
