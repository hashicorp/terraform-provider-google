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
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/google/verify"
)

var _ ephemeral.EphemeralResource = &googleEphemeralServiceAccountKey{}

func GoogleEphemeralServiceAccountKey() ephemeral.EphemeralResource {
	return &googleEphemeralServiceAccountKey{}
}

type googleEphemeralServiceAccountKey struct {
	providerConfig *transport_tpg.Config
}

func (p *googleEphemeralServiceAccountKey) Metadata(ctx context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_account_key"
}

type ephemeralServiceAccountKeyModel struct {
	Name          types.String `tfsdk:"name"`
	PublicKeyType types.String `tfsdk:"public_key_type"`
	KeyAlgorithm  types.String `tfsdk:"key_algorithm"`
	PublicKey     types.String `tfsdk:"public_key"`
}

func (p *googleEphemeralServiceAccountKey) Schema(ctx context.Context, req ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get an ephemeral service account public key.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the service account key. This must have format `projects/{PROJECT_ID}/serviceAccounts/{ACCOUNT}/keys/{KEYID}`, where `{ACCOUNT}` is the email address or unique id of the service account.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(verify.ServiceAccountKeyNameRegex),
						"must match regex: "+verify.ServiceAccountKeyNameRegex,
					),
				}},
			"public_key_type": schema.StringAttribute{
				Description: "The output format of the public key requested. TYPE_X509_PEM_FILE is the default output format.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"TYPE_X509_PEM_FILE",
						"TYPE_RAW_PUBLIC_KEY",
					),
				},
			},
			"key_algorithm": schema.StringAttribute{
				Description: "The algorithm used to generate the key.",
				Computed:    true,
			},
			"public_key": schema.StringAttribute{
				Description: "The public key, base64 encoded.",
				Computed:    true,
			},
		},
	}
}

func (p *googleEphemeralServiceAccountKey) Configure(ctx context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
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
