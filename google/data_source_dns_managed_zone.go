// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package google

import (
	"context"
	"fmt"

	"google.golang.org/api/dns/v1"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ datasource.DataSource              = &GoogleDnsManagedZoneDataSource{}
	_ datasource.DataSourceWithConfigure = &GoogleDnsManagedZoneDataSource{}
)

func NewGoogleDnsManagedZoneDataSource() datasource.DataSource {
	return &GoogleDnsManagedZoneDataSource{}
}

// GoogleDnsManagedZoneDataSource defines the data source implementation
type GoogleDnsManagedZoneDataSource struct {
	client  *dns.Service
	project types.String
}

type GoogleDnsManagedZoneModel struct {
	Id            types.String `tfsdk:"id"`
	DnsName       types.String `tfsdk:"dns_name"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	ManagedZoneId types.Int64  `tfsdk:"managed_zone_id"`
	NameServers   types.List   `tfsdk:"name_servers"`
	Visibility    types.String `tfsdk:"visibility"`
	Project       types.String `tfsdk:"project"`
}

func (d *GoogleDnsManagedZoneDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_managed_zone"
}

func (d *GoogleDnsManagedZoneDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Provides access to a zone's attributes within Google Cloud DNS",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description:         "A unique name for the resource.",
				MarkdownDescription: "A unique name for the resource.",
				Required:            true,
			},

			// Google Cloud DNS ManagedZone resources do not have a SelfLink attribute.
			"project": schema.StringAttribute{
				Description:         "The ID of the project for the Google Cloud.",
				MarkdownDescription: "The ID of the project for the Google Cloud.",
				Optional:            true,
			},

			"dns_name": schema.StringAttribute{
				Description:         "The fully qualified DNS name of this zone.",
				MarkdownDescription: "The fully qualified DNS name of this zone.",
				Computed:            true,
			},

			"description": schema.StringAttribute{
				Description:         "A textual description field.",
				MarkdownDescription: "A textual description field.",
				Computed:            true,
			},

			"managed_zone_id": schema.Int64Attribute{
				Description:         "Unique identifier for the resource; defined by the server.",
				MarkdownDescription: "Unique identifier for the resource; defined by the server.",
				Computed:            true,
			},

			"name_servers": schema.ListAttribute{
				Description: "The list of nameservers that will be authoritative for this " +
					"domain. Use NS records to redirect from your DNS provider to these names, " +
					"thus making Google Cloud DNS authoritative for this zone.",
				MarkdownDescription: "The list of nameservers that will be authoritative for this " +
					"domain. Use NS records to redirect from your DNS provider to these names, " +
					"thus making Google Cloud DNS authoritative for this zone.",
				Computed:    true,
				ElementType: types.StringType,
			},

			"visibility": schema.StringAttribute{
				Description: "The zone's visibility: public zones are exposed to the Internet, " +
					"while private zones are visible only to Virtual Private Cloud resources.",
				MarkdownDescription: "The zone's visibility: public zones are exposed to the Internet, " +
					"while private zones are visible only to Virtual Private Cloud resources.",
				Computed: true,
			},

			"id": schema.StringAttribute{
				Description:         "DNS managed zone identifier",
				MarkdownDescription: "DNS managed zone identifier",
				Computed:            true,
			},
		},
	}
}

func (d *GoogleDnsManagedZoneDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = p.NewDnsClient(p.userAgent, &resp.Diagnostics)
	d.project = p.project
}

func (d *GoogleDnsManagedZoneDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data GoogleDnsManagedZoneModel
	var metaData *ProviderMetaModel
	var diags diag.Diagnostics

	// Read Provider meta into the meta model
	resp.Diagnostics.Append(req.ProviderMeta.Get(ctx, &metaData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	d.client.UserAgent = generateFrameworkUserAgentString(metaData, d.client.UserAgent)

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Project = getProjectFramework(data.Project, d.project, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Id = types.StringValue(fmt.Sprintf("projects/%s/managedZones/%s", data.Project.ValueString(), data.Name.ValueString()))
	clientResp, err := d.client.ManagedZones.Get(data.Project.ValueString(), data.Name.ValueString()).Do()
	if err != nil {
		handleDatasourceNotFoundError(ctx, err, &resp.State, fmt.Sprintf("dataSourceDnsManagedZone %q", data.Name.ValueString()), &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	tflog.Trace(ctx, "read dns record set data source")

	data.DnsName = types.StringValue(clientResp.DnsName)
	data.Description = types.StringValue(clientResp.Description)
	data.ManagedZoneId = types.Int64Value(int64(clientResp.Id))
	data.Visibility = types.StringValue(clientResp.Visibility)
	data.NameServers, diags = types.ListValueFrom(ctx, types.StringType, clientResp.NameServers)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
