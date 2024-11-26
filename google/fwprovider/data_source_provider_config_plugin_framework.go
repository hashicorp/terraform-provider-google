// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package fwprovider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-google/google/fwmodels"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// Ensure the data source satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &GoogleProviderConfigPluginFrameworkDataSource{}
	_ datasource.DataSourceWithConfigure = &GoogleProviderConfigPluginFrameworkDataSource{}
)

func NewGoogleProviderConfigPluginFrameworkDataSource() datasource.DataSource {
	return &GoogleProviderConfigPluginFrameworkDataSource{}
}

type GoogleProviderConfigPluginFrameworkDataSource struct {
	providerConfig *transport_tpg.Config
}

// GoogleProviderConfigPluginFrameworkModel describes the data source and matches the schema. Its fields match those in a Config struct (google/transport/config.go) but uses a different type system.
//   - In the original Config struct old SDK/Go primitives types are used, e.g. `string`
//   - Here in the GoogleProviderConfigPluginFrameworkModel struct we need to use  the terraform-plugin-framework/types type system, e.g. `types.String`
//   - This is needed because the PF type system is 'baked into' how we define schemas. The schema will expect a nullable type.
//   - See terraform-plugin-framework/datasource/schema#StringAttribute's CustomType: https://pkg.go.dev/github.com/hashicorp/terraform-plugin-framework@v1.7.0/datasource/schema#StringAttribute
//   - Due to the different type systems of Config versus GoogleProviderConfigPluginFrameworkModel struct, we need to convert from primitive types to terraform-plugin-framework/types when we populate
//     GoogleProviderConfigPluginFrameworkModel structs with data in this data source's Read method.
type GoogleProviderConfigPluginFrameworkModel struct {
	Credentials                        types.String `tfsdk:"credentials"`
	AccessToken                        types.String `tfsdk:"access_token"`
	ImpersonateServiceAccount          types.String `tfsdk:"impersonate_service_account"`
	ImpersonateServiceAccountDelegates types.List   `tfsdk:"impersonate_service_account_delegates"`
	Project                            types.String `tfsdk:"project"`
	BillingProject                     types.String `tfsdk:"billing_project"`
	Region                             types.String `tfsdk:"region"`
	Zone                               types.String `tfsdk:"zone"`
	Scopes                             types.List   `tfsdk:"scopes"`
	//	omit Batching
	UserProjectOverride                       types.Bool   `tfsdk:"user_project_override"`
	RequestTimeout                            types.String `tfsdk:"request_timeout"`
	RequestReason                             types.String `tfsdk:"request_reason"`
	UniverseDomain                            types.String `tfsdk:"universe_domain"`
	DefaultLabels                             types.Map    `tfsdk:"default_labels"`
	AddTerraformAttributionLabel              types.Bool   `tfsdk:"add_terraform_attribution_label"`
	TerraformAttributionLabelAdditionStrategy types.String `tfsdk:"terraform_attribution_label_addition_strategy"`
}

func (d *GoogleProviderConfigPluginFrameworkDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_provider_config_plugin_framework"
}

func (d *GoogleProviderConfigPluginFrameworkDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	resp.Schema = schema.Schema{

		Description:         "Use this data source to access the configuration of the Google Cloud provider. This data source is implemented with the SDK.",
		MarkdownDescription: "Use this data source to access the configuration of the Google Cloud provider. This data source is implemented with the SDK.",
		Attributes: map[string]schema.Attribute{
			// Start of user inputs
			"access_token": schema.StringAttribute{
				Description:         "The access_token argument used to configure the provider",
				MarkdownDescription: "The access_token argument used to configure the provider",
				Computed:            true,
				Sensitive:           true,
			},
			"credentials": schema.StringAttribute{
				Description:         "The credentials argument used to configure the provider",
				MarkdownDescription: "The credentials argument used to configure the provider",
				Computed:            true,
				Sensitive:           true,
			},
			"impersonate_service_account": schema.StringAttribute{
				Description:         "The impersonate_service_account argument used to configure the provider",
				MarkdownDescription: "The impersonate_service_account argument used to configure the provider.",
				Computed:            true,
			},
			"impersonate_service_account_delegates": schema.ListAttribute{
				ElementType:         types.StringType,
				Description:         "The impersonate_service_account_delegates argument used to configure the provider",
				MarkdownDescription: "The impersonate_service_account_delegates argument used to configure the provider.",
				Computed:            true,
			},
			"project": schema.StringAttribute{
				Description:         "The project argument used to configure the provider",
				MarkdownDescription: "The project argument used to configure the provider.",
				Computed:            true,
			},
			"region": schema.StringAttribute{
				Description:         "The region argument used to configure the provider.",
				MarkdownDescription: "The region argument used to configure the provider.",
				Computed:            true,
			},
			"billing_project": schema.StringAttribute{
				Description:         "The billing_project argument used to configure the provider.",
				MarkdownDescription: "The billing_project argument used to configure the provider.",
				Computed:            true,
			},
			"zone": schema.StringAttribute{
				Description:         "The zone argument used to configure the provider.",
				MarkdownDescription: "The zone argument used to configure the provider.",
				Computed:            true,
			},
			"universe_domain": schema.StringAttribute{
				Description:         "The universe_domain argument used to configure the provider.",
				MarkdownDescription: "The universe_domain argument used to configure the provider.",
				Computed:            true,
			},
			"scopes": schema.ListAttribute{
				ElementType:         types.StringType,
				Description:         "The scopes argument used to configure the provider.",
				MarkdownDescription: "The scopes argument used to configure the provider.",
				Computed:            true,
			},
			"user_project_override": schema.BoolAttribute{
				Description:         "The user_project_override argument used to configure the provider.",
				MarkdownDescription: "The user_project_override argument used to configure the provider.",
				Computed:            true,
			},
			"request_reason": schema.StringAttribute{
				Description:         "The request_reason argument used to configure the provider.",
				MarkdownDescription: "The request_reason argument used to configure the provider.",
				Computed:            true,
			},
			"request_timeout": schema.StringAttribute{
				Description:         "The request_timeout argument used to configure the provider.",
				MarkdownDescription: "The request_timeout argument used to configure the provider.",
				Computed:            true,
			},
			"default_labels": schema.MapAttribute{
				ElementType:         types.StringType,
				Description:         "The default_labels argument used to configure the provider.",
				MarkdownDescription: "The default_labels argument used to configure the provider.",
				Computed:            true,
			},
			"add_terraform_attribution_label": schema.BoolAttribute{
				Description:         "The add_terraform_attribution_label argument used to configure the provider.",
				MarkdownDescription: "The add_terraform_attribution_label argument used to configure the provider.",
				Computed:            true,
			},
			"terraform_attribution_label_addition_strategy": schema.StringAttribute{
				Description:         "The terraform_attribution_label_addition_strategy argument used to configure the provider.",
				MarkdownDescription: "The terraform_attribution_label_addition_strategy argument used to configure the provider.",
				Computed:            true,
			},
			// End of user inputs

			// Note - this data source excludes the default and custom endpoints for individual services
		},
	}
}

func (d *GoogleProviderConfigPluginFrameworkDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	// Required for accessing project, region, zone and tokenSource
	d.providerConfig = p
}

func (d *GoogleProviderConfigPluginFrameworkDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data GoogleProviderConfigPluginFrameworkModel
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

	// Copy all values from the provider config into this data source
	//    - The 'meta' from the provider configuration process uses Go primitive types (e.g. `string`) but this data source needs to use the plugin-framework type system due to being PF-implemented
	//    - As a result we have to make conversions between type systems in the value assignments below
	data.Credentials = types.StringValue(d.providerConfig.Credentials)
	data.AccessToken = types.StringValue(d.providerConfig.AccessToken)
	data.ImpersonateServiceAccount = types.StringValue(d.providerConfig.ImpersonateServiceAccount)

	delegateAttrs := make([]attr.Value, len(d.providerConfig.ImpersonateServiceAccountDelegates))
	for i, delegate := range d.providerConfig.ImpersonateServiceAccountDelegates {
		delegateAttrs[i] = types.StringValue(delegate)
	}
	delegates, di := types.ListValue(types.StringType, delegateAttrs)
	if di.HasError() {
		resp.Diagnostics.Append(di...)
	}
	data.ImpersonateServiceAccountDelegates = delegates

	data.Project = types.StringValue(d.providerConfig.Project)
	data.Region = types.StringValue(d.providerConfig.Region)
	data.BillingProject = types.StringValue(d.providerConfig.BillingProject)
	data.Zone = types.StringValue(d.providerConfig.Zone)
	data.UniverseDomain = types.StringValue(d.providerConfig.UniverseDomain)

	scopeAttrs := make([]attr.Value, len(d.providerConfig.Scopes))
	for i, scope := range d.providerConfig.Scopes {
		scopeAttrs[i] = types.StringValue(scope)
	}
	scopes, di := types.ListValue(types.StringType, scopeAttrs)
	if di.HasError() {
		resp.Diagnostics.Append(di...)
	}
	data.Scopes = scopes

	data.UserProjectOverride = types.BoolValue(d.providerConfig.UserProjectOverride)
	data.RequestReason = types.StringValue(d.providerConfig.RequestReason)
	data.RequestTimeout = types.StringValue(d.providerConfig.RequestTimeout.String())

	lbs := make(map[string]attr.Value, len(d.providerConfig.DefaultLabels))
	for k, v := range d.providerConfig.DefaultLabels {
		lbs[k] = types.StringValue(v)
	}
	labels, di := types.MapValueFrom(ctx, types.StringType, lbs)
	if di.HasError() {
		resp.Diagnostics.Append(di...)
	}
	data.DefaultLabels = labels

	data.AddTerraformAttributionLabel = types.BoolValue(d.providerConfig.AddTerraformAttributionLabel)
	data.TerraformAttributionLabelAdditionStrategy = types.StringValue(d.providerConfig.TerraformAttributionLabelAdditionStrategy)

	// Warn users against using this data source
	resp.Diagnostics.Append(diag.NewWarningDiagnostic(
		"Data source google_provider_config_plugin_framework should not be used",
		"Data source google_provider_config_plugin_framework is intended to be used only in acceptance tests for the provider. Instead, please use the google_client_config data source to access provider configuration details, or open a GitHub issue requesting new features in that datasource. Please go to: https://github.com/hashicorp/terraform-provider-google/issues/new/choose",
	))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
