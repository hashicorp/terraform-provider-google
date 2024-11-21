// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package fwprovider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/metaschema"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-provider-google/google/functions"
	"github.com/hashicorp/terraform-provider-google/google/fwmodels"
	"github.com/hashicorp/terraform-provider-google/google/fwtransport"
	"github.com/hashicorp/terraform-provider-google/google/fwvalidators"
	"github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"
	"github.com/hashicorp/terraform-provider-google/version"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider                       = &FrameworkProvider{}
	_ provider.ProviderWithMetaSchema         = &FrameworkProvider{}
	_ provider.ProviderWithFunctions          = &FrameworkProvider{}
	_ provider.ProviderWithEphemeralResources = &FrameworkProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New() provider.ProviderWithMetaSchema {
	return &FrameworkProvider{
		Version: version.ProviderVersion,
	}
}

// FrameworkProvider is the provider implementation.
type FrameworkProvider struct {
	fwtransport.FrameworkProviderConfig
	Version string
}

// Metadata returns
// - the provider type name : this controls how "google" is present at the start of all resource type names
// - the provider version   : this is currently unused by Terraform core
func (p *FrameworkProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "google"
	resp.Version = p.Version
}

// MetaSchema returns the provider meta schema.
func (p *FrameworkProvider) MetaSchema(_ context.Context, _ provider.MetaSchemaRequest, resp *provider.MetaSchemaResponse) {
	resp.Schema = metaschema.Schema{
		Attributes: map[string]metaschema.Attribute{
			"module_name": metaschema.StringAttribute{
				Optional: true,
			},
		},
	}
}

// Schema defines the provider-level schema for configuration data.
func (p *FrameworkProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"credentials": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("access_token"),
					}...),
					fwvalidators.CredentialsValidator(),
					fwvalidators.NonEmptyStringValidator(),
				},
			},
			"access_token": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("credentials"),
					}...),
					fwvalidators.NonEmptyStringValidator(),
				},
			},
			"impersonate_service_account": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					fwvalidators.NonEmptyStringValidator(),
				},
			},
			"impersonate_service_account_delegates": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"project": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					fwvalidators.NonEmptyStringValidator(),
				},
			},
			"billing_project": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					fwvalidators.NonEmptyStringValidator(),
				},
			},
			"region": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					fwvalidators.NonEmptyStringValidator(),
				},
			},
			"zone": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					fwvalidators.NonEmptyStringValidator(),
				},
			},
			"scopes": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"user_project_override": schema.BoolAttribute{
				Optional: true,
			},
			"request_timeout": schema.StringAttribute{
				Optional: true,
			},
			"request_reason": schema.StringAttribute{
				Optional: true,
			},
			"universe_domain": schema.StringAttribute{
				Optional: true,
			},
			"default_labels": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"add_terraform_attribution_label": schema.BoolAttribute{
				Optional: true,
			},
			"terraform_attribution_label_addition_strategy": schema.StringAttribute{
				Optional: true,
			},
			// Generated Products
			"access_approval_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"access_context_manager_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"active_directory_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"alloydb_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"apigee_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"app_engine_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"apphub_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"artifact_registry_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"beyondcorp_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"biglake_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"big_query_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"bigquery_analytics_hub_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"bigquery_connection_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"bigquery_datapolicy_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"bigquery_data_transfer_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"bigquery_reservation_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"bigtable_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"billing_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"binary_authorization_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"blockchain_node_engine_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"certificate_manager_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"cloud_asset_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"cloud_build_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"cloudbuildv2_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"clouddeploy_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"clouddomains_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"cloud_functions_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"cloudfunctions2_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"cloud_identity_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"cloud_ids_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"cloud_quotas_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"cloud_run_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"cloud_run_v2_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"cloud_scheduler_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"cloud_tasks_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"composer_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"compute_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"container_analysis_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"container_attached_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"core_billing_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"database_migration_service_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"data_catalog_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"data_fusion_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"data_loss_prevention_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"data_pipeline_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"dataplex_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"dataproc_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"dataproc_gdc_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"dataproc_metastore_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"datastream_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"deployment_manager_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"dialogflow_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"dialogflow_cx_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"discovery_engine_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"dns_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"document_ai_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"document_ai_warehouse_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"edgecontainer_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"edgenetwork_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"essential_contacts_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"filestore_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"firebase_app_check_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"firestore_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"gke_backup_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"gke_hub_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"gke_hub2_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"gkeonprem_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"healthcare_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"iam2_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"iam_beta_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"iam_workforce_pool_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"iap_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"identity_platform_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"integration_connectors_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"integrations_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"kms_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"logging_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"looker_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"memcache_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"memorystore_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"migration_center_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"ml_engine_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"monitoring_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"netapp_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"network_connectivity_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"network_management_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"network_security_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"network_services_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"notebooks_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"oracle_database_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"org_policy_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"os_config_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"os_login_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"privateca_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"privileged_access_manager_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"public_ca_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"pubsub_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"pubsub_lite_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"redis_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"resource_manager_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"secret_manager_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"secret_manager_regional_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"secure_source_manager_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"security_center_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"security_center_management_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"security_center_v2_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"securityposture_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"service_management_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"service_networking_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"service_usage_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"site_verification_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"source_repo_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"spanner_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"sql_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"storage_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"storage_insights_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"storage_transfer_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"tags_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"tpu_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"transcoder_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"vertex_ai_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"vmwareengine_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"vpc_access_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"workbench_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"workflows_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},

			// Handwritten Products / Versioned / Atypical Entries
			"cloud_billing_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"container_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"dataflow_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"iam_credentials_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"resource_manager_v3_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"iam_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"tags_location_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},

			// dcl
			"container_aws_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
			"container_azure_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					transport_tpg.CustomEndpointValidator(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"batching": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"send_after": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								fwvalidators.NonNegativeDurationValidator(),
							},
						},
						"enable_batching": schema.BoolAttribute{
							Optional: true,
						},
					},
				},
			},
		},
	}

	transport_tpg.ConfigureDCLCustomEndpointAttributesFramework(&resp.Schema)
}

// Configure prepares the metadata/'meta' required for data sources and resources to function.
// Configuration logic implemented here should take user inputs and use them to populate a struct
// with that necessary metadata, e.g. default project value, configured client, etc.
// That prepared 'meta' struct is then returned in the response, and that value will later be supplied to all resources/data sources when they need to configure themselves.
func (p *FrameworkProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data fwmodels.ProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	p.LoadAndValidateFramework(ctx, &data, req.TerraformVersion, &resp.Diagnostics, p.Version)
	if resp.Diagnostics.HasError() {
		return
	}

	// This is how we make provider configuration info (configured clients, default project, etc) available to resources and data sources
	// implemented using the plugin-framework. The resources' Configure functions receive this data in the ConfigureRequest argument.
	resp.DataSourceData = &p.FrameworkProviderConfig
	resp.ResourceData = &p.FrameworkProviderConfig
	resp.EphemeralResourceData = &p.FrameworkProviderConfig
}

// DataSources defines the data sources implemented in the provider.
func (p *FrameworkProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		resourcemanager.NewGoogleClientConfigDataSource,
		resourcemanager.NewGoogleClientOpenIDUserinfoDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *FrameworkProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}

// Functions defines the provider functions implemented in the provider.
func (p *FrameworkProvider) Functions(_ context.Context) []func() function.Function {
	return []func() function.Function{
		functions.NewLocationFromIdFunction,
		functions.NewNameFromIdFunction,
		functions.NewProjectFromIdFunction,
		functions.NewRegionFromIdFunction,
		functions.NewRegionFromZoneFunction,
		functions.NewZoneFromIdFunction,
	}
}

// EphemeralResources defines the resources that are of ephemeral type implemented in the provider.
func (p *FrameworkProvider) EphemeralResources(_ context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{
		resourcemanager.GoogleEphemeralServiceAccountAccessToken,
		resourcemanager.GoogleEphemeralServiceAccountIdToken,
		resourcemanager.GoogleEphemeralServiceAccountJwt,
	}
}
