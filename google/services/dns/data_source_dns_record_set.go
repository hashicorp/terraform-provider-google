// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package dns

import (
	"context"
	"fmt"

	"google.golang.org/api/dns/v1"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-provider-google/google/fwmodels"
	"github.com/hashicorp/terraform-provider-google/google/fwresource"
	"github.com/hashicorp/terraform-provider-google/google/fwtransport"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ datasource.DataSource              = &GoogleDnsRecordSetDataSource{}
	_ datasource.DataSourceWithConfigure = &GoogleDnsRecordSetDataSource{}
)

func NewGoogleDnsRecordSetDataSource() datasource.DataSource {
	return &GoogleDnsRecordSetDataSource{}
}

// GoogleDnsRecordSetDataSource defines the data source implementation
type GoogleDnsRecordSetDataSource struct {
	client  *dns.Service
	project types.String
}

type GoogleDnsRecordSetModel struct {
	Id          types.String `tfsdk:"id"`
	ManagedZone types.String `tfsdk:"managed_zone"`
	Name        types.String `tfsdk:"name"`
	Rrdatas     types.List   `tfsdk:"rrdatas"`
	Ttl         types.Int64  `tfsdk:"ttl"`
	Type        types.String `tfsdk:"type"`
	Project     types.String `tfsdk:"project"`
}

func (d *GoogleDnsRecordSetDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_record_set"
}

func (d *GoogleDnsRecordSetDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "A DNS record set within Google Cloud DNS",

		Attributes: map[string]schema.Attribute{
			"managed_zone": schema.StringAttribute{
				MarkdownDescription: "The Name of the zone.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The DNS name for the resource.",
				Required:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The identifier of a supported record type. See the list of Supported DNS record types.",
				Required:            true,
			},
			"project": schema.StringAttribute{
				MarkdownDescription: "The ID of the project for the Google Cloud.",
				Optional:            true,
			},
			"rrdatas": schema.ListAttribute{
				MarkdownDescription: "The string data for the records in this record set.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"ttl": schema.Int64Attribute{
				MarkdownDescription: "The time-to-live of this record set (seconds).",
				Computed:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "DNS record set identifier",
				Computed:            true,
			},
		},
	}
}

func (d *GoogleDnsRecordSetDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	p, ok := req.ProviderData.(*fwtransport.FrameworkProviderConfig)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *fwtransport.FrameworkProviderConfig, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = p.NewDnsClient(p.UserAgent, &resp.Diagnostics)
	d.project = p.Project
}

func (d *GoogleDnsRecordSetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data GoogleDnsRecordSetModel
	var metaData *fwmodels.ProviderMetaModel
	var diags diag.Diagnostics

	// Read Provider meta into the meta model
	resp.Diagnostics.Append(req.ProviderMeta.Get(ctx, &metaData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	d.client.UserAgent = fwtransport.GenerateFrameworkUserAgentString(metaData, d.client.UserAgent)

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Project = fwresource.GetProjectFramework(data.Project, d.project, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Id = types.StringValue(fmt.Sprintf("projects/%s/managedZones/%s/rrsets/%s/%s", data.Project.ValueString(), data.ManagedZone.ValueString(), data.Name.ValueString(), data.Type.ValueString()))
	clientResp, err := d.client.ResourceRecordSets.List(data.Project.ValueString(), data.ManagedZone.ValueString()).Name(data.Name.ValueString()).Type(data.Type.ValueString()).Do()
	if err != nil {
		fwtransport.HandleDatasourceNotFoundError(ctx, err, &resp.State, fmt.Sprintf("dataSourceDnsRecordSet %q", data.Name.ValueString()), &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	if len(clientResp.Rrsets) != 1 {
		resp.Diagnostics.AddError("only expected 1 record set", fmt.Sprintf("%d record sets were returned", len(clientResp.Rrsets)))
	}

	tflog.Trace(ctx, "read dns record set data source")

	data.Type = types.StringValue(clientResp.Rrsets[0].Type)
	data.Ttl = types.Int64Value(clientResp.Rrsets[0].Ttl)
	data.Rrdatas, diags = types.ListValueFrom(ctx, types.StringType, clientResp.Rrsets[0].Rrdatas)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
