// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"google.golang.org/api/compute/v1"

	"github.com/hashicorp/terraform-provider-google/google/fwmodels"
	"github.com/hashicorp/terraform-provider-google/google/fwresource"
	"github.com/hashicorp/terraform-provider-google/google/fwtransport"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &ComputeNetworkFWDataSource{}
	_ datasource.DataSourceWithConfigure = &ComputeNetworkFWDataSource{}
)

// NewComputeNetworkFWDataSource is a helper function to simplify the provider implementation.
func NewComputeNetworkFWDataSource() datasource.DataSource {
	return &ComputeNetworkFWDataSource{}
}

// ComputeNetworkFWDataSource is the data source implementation.
type ComputeNetworkFWDataSource struct {
	client         *compute.Service
	providerConfig *transport_tpg.Config
}

type ComputeNetworkModel struct {
	Id                types.String `tfsdk:"id"`
	Project           types.String `tfsdk:"project"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	NetworkId         types.Int64  `tfsdk:"network_id"`
	NumericId         types.String `tfsdk:"numeric_id"`
	GatewayIpv4       types.String `tfsdk:"gateway_ipv4"`
	InternalIpv6Range types.String `tfsdk:"internal_ipv6_range"`
	SelfLink          types.String `tfsdk:"self_link"`
	// NetworkProfile  types.String `tfsdk:"network_profile"`
	// SubnetworksSelfLinks  types.List `tfsdk:"subnetworks_self_links"`
}

// Metadata returns the data source type name.
func (d *ComputeNetworkFWDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_fw_compute_network"
}

func (d *ComputeNetworkFWDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	p, ok := req.ProviderData.(*transport_tpg.Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *transport_tpg.Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = p.NewComputeClient(p.UserAgent)
	if resp.Diagnostics.HasError() {
		return
	}
	d.providerConfig = p
}

// Schema defines the schema for the data source.
func (d *ComputeNetworkFWDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A data source to get network details.",

		Attributes: map[string]schema.Attribute{
			"project": schema.StringAttribute{
				Description:         `The project name.`,
				MarkdownDescription: `The project name.`,
				Optional:            true,
			},
			"name": schema.StringAttribute{
				Description:         `The name of the Compute network.`,
				MarkdownDescription: `The name of the Compute network.`,
				Required:            true,
			},
			"description": schema.StringAttribute{
				Description:         `The description of the network.`,
				MarkdownDescription: `The description of the network.`,
				Computed:            true,
			},
			"network_id": schema.Int64Attribute{
				Description:         `The network ID.`,
				MarkdownDescription: `The network ID.`,
				Computed:            true,
			},
			"numeric_id": schema.StringAttribute{
				Description:         `The numeric ID of the network. Deprecated in favor of network_id.`,
				MarkdownDescription: `The numeric ID of the network. Deprecated in favor of network_id.`,
				Computed:            true,
				DeprecationMessage:  "`numeric_id` is deprecated and will be removed in a future major release. Use `network_id` instead.",
			},
			"gateway_ipv4": schema.StringAttribute{
				Description:         `The gateway address for default routing out of the network.`,
				MarkdownDescription: `The gateway address for default routing out of the network.`,
				Computed:            true,
			},
			"internal_ipv6_range": schema.StringAttribute{
				Description:         `The internal ipv6 address range of the network.`,
				MarkdownDescription: `The internal ipv6 address range of the network.`,
				Computed:            true,
			},
			"self_link": schema.StringAttribute{
				Description:         `The network self link.`,
				MarkdownDescription: `The network self link.`,
				Computed:            true,
			},
			// This is included for backwards compatibility with the original, SDK-implemented data source.
			"id": schema.StringAttribute{
				Description:         "Project identifier",
				MarkdownDescription: "Project identifier",
				Computed:            true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *ComputeNetworkFWDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ComputeNetworkModel
	var metaData *fwmodels.ProviderMetaModel

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

	// Use provider_meta to set User-Agent
	d.client.UserAgent = fwtransport.GenerateFrameworkUserAgentString(metaData, d.client.UserAgent)

	project := fwresource.GetProjectFramework(data.Project, types.StringValue(d.providerConfig.Project), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// GET Request
	clientResp, err := d.client.Networks.Get(project.ValueString(), data.Name.ValueString()).Do()
	if err != nil {
		fwtransport.HandleDatasourceNotFoundError(ctx, err, &resp.State, fmt.Sprintf("dataSourceComputeNetwork %q", data.Name.ValueString()), &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	tflog.Trace(ctx, "read compute network data source")

	// Put data in model
	id := fmt.Sprintf("projects/%s/global/networks/%s", project.ValueString(), clientResp.Name)
	data.Id = types.StringValue(id)
	data.Description = types.StringValue(clientResp.Description)
	data.NetworkId = types.Int64Value(int64(clientResp.Id))
	data.NumericId = types.StringValue(strconv.Itoa(int(clientResp.Id)))
	data.GatewayIpv4 = types.StringValue(clientResp.GatewayIPv4)
	data.InternalIpv6Range = types.StringValue(clientResp.InternalIpv6Range)
	data.SelfLink = types.StringValue(clientResp.SelfLink)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
