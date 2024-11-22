// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
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

var _ ephemeral.EphemeralResource = &googleEphemeralServiceAccountJwt{}

func GoogleEphemeralServiceAccountJwt() ephemeral.EphemeralResource {
	return &googleEphemeralServiceAccountJwt{}
}

type googleEphemeralServiceAccountJwt struct {
	providerConfig *transport_tpg.Config
}

func (p *googleEphemeralServiceAccountJwt) Metadata(ctx context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_account_jwt"
}

type ephemeralServiceAccountJwtModel struct {
	Payload              types.String `tfsdk:"payload"`
	ExpiresIn            types.Int64  `tfsdk:"expires_in"`
	TargetServiceAccount types.String `tfsdk:"target_service_account"`
	Delegates            types.Set    `tfsdk:"delegates"`
	Jwt                  types.String `tfsdk:"jwt"`
}

func (p *googleEphemeralServiceAccountJwt) Schema(ctx context.Context, req ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Produces an arbitrary self-signed JWT for service accounts.",
		Attributes: map[string]schema.Attribute{
			"payload": schema.StringAttribute{
				Required:    true,
				Description: `A JSON-encoded JWT claims set that will be included in the signed JWT.`,
			},
			"expires_in": schema.Int64Attribute{
				Optional:    true,
				Description: "Number of seconds until the JWT expires. If set and non-zero an `exp` claim will be added to the payload derived from the current timestamp plus expires_in seconds.",
				Validators: []validator.Int64{
					int64validator.AtLeast(1), // Must be greater than 0
				},
			},
			"target_service_account": schema.StringAttribute{
				Description: "The email of the service account that will sign the JWT.",
				Required:    true,
				Validators: []validator.String{
					fwvalidators.ServiceAccountEmailValidator{},
				},
			},
			"delegates": schema.SetAttribute{
				Description: "Delegate chain of approvals needed to perform full impersonation. Specify the fully qualified service account name.",
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(fwvalidators.ServiceAccountEmailValidator{}),
				},
			},
			"jwt": schema.StringAttribute{
				Description: "The signed JWT containing the JWT Claims Set from the `payload`.",
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *googleEphemeralServiceAccountJwt) Configure(ctx context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
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

func (p *googleEphemeralServiceAccountJwt) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data ephemeralServiceAccountJwtModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	payload := data.Payload.ValueString()

	if !data.ExpiresIn.IsNull() {
		expiresIn := data.ExpiresIn.ValueInt64()
		var decoded map[string]interface{}
		if err := json.Unmarshal([]byte(payload), &decoded); err != nil {
			resp.Diagnostics.AddError("Error decoding payload", err.Error())
			return
		}

		decoded["exp"] = time.Now().Add(time.Duration(expiresIn) * time.Second).Unix()

		payloadBytesWithExp, err := json.Marshal(decoded)
		if err != nil {
			resp.Diagnostics.AddError("Error re-encoding payload", err.Error())
			return
		}

		payload = string(payloadBytesWithExp)

	}

	name := fmt.Sprintf("projects/-/serviceAccounts/%s", data.TargetServiceAccount.ValueString())

	service := p.providerConfig.NewIamCredentialsClient(p.providerConfig.UserAgent)
	jwtRequest := &iamcredentials.SignJwtRequest{
		Payload:   payload,
		Delegates: fwutils.StringSet(data.Delegates),
	}

	jwtResponse, err := service.Projects.ServiceAccounts.SignJwt(name, jwtRequest).Do()
	if err != nil {
		resp.Diagnostics.AddError("Error calling iamcredentials.SignJwt", err.Error())
		return
	}

	data.Jwt = types.StringValue(jwtResponse.SignedJwt)

	resp.Diagnostics.Append(resp.Result.Set(ctx, data)...)
}
