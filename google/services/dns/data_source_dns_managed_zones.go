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
	_ datasource.DataSource              = &GoogleDnsManagedZonesDataSource{}
	_ datasource.DataSourceWithConfigure = &GoogleDnsManagedZonesDataSource{}
)

func NewGoogleDnsManagedZonesDataSource() datasource.DataSource {
	return &GoogleDnsManagedZonesDataSource{}
}

// GoogleDnsManagedZonesDataSource defines the data source implementation
type GoogleDnsManagedZonesDataSource struct {
	client  *dns.Service
	project types.String
}

type GoogleDnsManagedZonesModel struct {
	Id           types.String `tfsdk:"id"`
	Project      types.String `tfsdk:"project"`
	ManagedZones types.List   `tfsdk:"managed_zones"`
}

func (d *GoogleDnsManagedZonesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_managed_zones"
}

func (d *GoogleDnsManagedZonesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Provides access to all zones for a given project within Google Cloud DNS",

		Attributes: map[string]schema.Attribute{

			"project": schema.StringAttribute{
				Description:         "The ID of the project for the Google Cloud.",
				MarkdownDescription: "The ID of the project for the Google Cloud.",
				Optional:            true,
			},

			// Id field is added to match plugin-framework migrated google_dns_managed_zone data source
			// Whilst ID fields are required in the SDK, they're not needed in the plugin-framework.
			"id": schema.StringAttribute{
				Description:         "foobar",
				MarkdownDescription: "foobar",
				Computed:            true,
			},
		},

		Blocks: map[string]schema.Block{
			"managed_zones": schema.ListNestedBlock{
				Description:         "The list of managed zones in the given project.",
				MarkdownDescription: "The list of managed zones in the given project.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description:         "A unique name for the resource.",
							MarkdownDescription: "A unique name for the resource.",
							Computed:            true,
						},

						// Google Cloud DNS ManagedZone resources do not have a SelfLink attribute.
						"project": schema.StringAttribute{
							Description:         "The ID of the project for the Google Cloud.",
							MarkdownDescription: "The ID of the project for the Google Cloud.",
							Computed:            true,
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
				},
			},
		},
	}
}

func (d *GoogleDnsManagedZonesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *GoogleDnsManagedZonesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data GoogleDnsManagedZonesModel
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

	data.Id = types.StringValue(fmt.Sprintf("projects/%s/managedZones", data.Project.ValueString()))

	tflog.Debug(ctx, fmt.Sprintf("fetching managed zones from project %s", data.Project.ValueString()))

	clientResp, err := d.client.ManagedZones.List(data.Project.ValueString()).Do()
	if err != nil {
		fwtransport.HandleDatasourceNotFoundError(ctx, err, &resp.State, fmt.Sprintf("dataSourceDnsManagedZones %q", data.Project.ValueString()), &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	tflog.Trace(ctx, "read dns managed zones data source")

	zones, di := flattenManagedZones(ctx, clientResp.ManagedZones, data.Project.ValueString())
	diags.Append(di...)

	if len(zones) > 0 {
		mzObjType := types.ObjectType{}.WithAttributeTypes(getDnsManagedZoneAttrs())
		data.ManagedZones, di = types.ListValueFrom(ctx, mzObjType, zones)
		diags.Append(di...)
	}

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenManagedZones(ctx context.Context, managedZones []*dns.ManagedZone, project string) ([]types.Object, diag.Diagnostics) {
	var zones []types.Object
	var diags diag.Diagnostics

	for _, zone := range managedZones {

		data := GoogleDnsManagedZoneModel{
			// Id is not an API value but we assemble it here to match the google_dns_managed_zone data source
			// and fulfil the GoogleDnsManagedZoneModel's fields.
			// IDs are not required in the plugin-framework (vs the SDK)
			Id:      types.StringValue(fmt.Sprintf("projects/%s/managedZones/%s", project, zone.Name)),
			Project: types.StringValue(project),

			DnsName:       types.StringValue(zone.DnsName),
			Name:          types.StringValue(zone.Name),
			Description:   types.StringValue(zone.Description),
			ManagedZoneId: types.Int64Value(int64(zone.Id)),
			Visibility:    types.StringValue(zone.Visibility),
		}

		data.NameServers, diags = types.ListValueFrom(ctx, types.StringType, zone.NameServers)
		diags.Append(diags...)

		obj, d := types.ObjectValueFrom(ctx, getDnsManagedZoneAttrs(), data)
		diags.Append(d...)

		zones = append(zones, obj)
	}

	return zones, diags
}
