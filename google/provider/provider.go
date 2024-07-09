// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package provider

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/google/verify"
	"github.com/hashicorp/terraform-provider-google/version"
)

// Provider returns a *schema.Provider.
func Provider() *schema.Provider {

	// The mtls service client gives the type of endpoint (mtls/regular)
	// at client creation. Since we use a shared client for requests we must
	// rewrite the endpoints to be mtls endpoints for the scenario where
	// mtls is enabled.
	if isMtls() {
		// if mtls is enabled switch all default endpoints to use the mtls endpoint
		for key, bp := range transport_tpg.DefaultBasePaths {
			transport_tpg.DefaultBasePaths[key] = getMtlsEndpoint(bp)
		}
	}

	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"credentials": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  ValidateCredentials,
				ConflictsWith: []string{"access_token"},
			},

			"access_token": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  ValidateEmptyStrings,
				ConflictsWith: []string{"credentials"},
			},

			"impersonate_service_account": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: ValidateEmptyStrings,
			},

			"impersonate_service_account_delegates": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"project": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: ValidateEmptyStrings,
			},

			"billing_project": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: ValidateEmptyStrings,
			},

			"region": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: ValidateEmptyStrings,
			},

			"zone": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: ValidateEmptyStrings,
			},

			"scopes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"universe_domain": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"batching": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"send_after": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: verify.ValidateNonNegativeDuration(),
						},
						"enable_batching": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},

			"user_project_override": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"request_timeout": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"request_reason": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"default_labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"add_terraform_attribution_label": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"terraform_attribution_label_addition_strategy": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Generated Products
			"access_approval_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"access_context_manager_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"active_directory_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"alloydb_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"apigee_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"app_engine_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"apphub_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"artifact_registry_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"beyondcorp_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"biglake_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"big_query_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"bigquery_analytics_hub_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"bigquery_connection_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"bigquery_datapolicy_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"bigquery_data_transfer_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"bigquery_reservation_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"bigtable_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"billing_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"binary_authorization_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"blockchain_node_engine_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"certificate_manager_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"cloud_asset_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"cloud_build_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"cloudbuildv2_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"clouddeploy_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"clouddomains_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"cloud_functions_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"cloudfunctions2_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"cloud_identity_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"cloud_ids_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"cloud_quotas_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"cloud_run_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"cloud_run_v2_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"cloud_scheduler_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"cloud_tasks_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"composer_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"compute_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"container_analysis_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"container_attached_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"core_billing_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"database_migration_service_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"data_catalog_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"data_fusion_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"data_loss_prevention_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"data_pipeline_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"dataplex_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"dataproc_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"dataproc_metastore_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"datastore_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"datastream_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"deployment_manager_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"dialogflow_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"dialogflow_cx_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"discovery_engine_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"dns_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"document_ai_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"document_ai_warehouse_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"edgecontainer_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"edgenetwork_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"essential_contacts_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"filestore_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"firebase_app_check_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"firestore_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"gke_backup_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"gke_hub_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"gke_hub2_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"gkeonprem_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"healthcare_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"iam2_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"iam_beta_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"iam_workforce_pool_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"iap_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"identity_platform_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"integration_connectors_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"integrations_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"kms_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"logging_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"looker_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"memcache_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"migration_center_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"ml_engine_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"monitoring_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"netapp_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"network_connectivity_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"network_management_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"network_security_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"network_services_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"notebooks_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"org_policy_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"os_config_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"os_login_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"privateca_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"privileged_access_manager_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"public_ca_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"pubsub_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"pubsub_lite_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"redis_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"resource_manager_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"secret_manager_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"secure_source_manager_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"security_center_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"security_center_management_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"security_center_v2_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"securityposture_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"service_management_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"service_networking_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"service_usage_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"site_verification_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"source_repo_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"spanner_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"sql_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"storage_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"storage_insights_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"storage_transfer_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"tags_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"tpu_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"vertex_ai_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"vmwareengine_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"vpc_access_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"workbench_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},
			"workflows_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: transport_tpg.ValidateCustomEndpoint,
			},

			// Handwritten Products / Versioned / Atypical Entries
			transport_tpg.CloudBillingCustomEndpointEntryKey:      transport_tpg.CloudBillingCustomEndpointEntry,
			transport_tpg.ComposerCustomEndpointEntryKey:          transport_tpg.ComposerCustomEndpointEntry,
			transport_tpg.ContainerCustomEndpointEntryKey:         transport_tpg.ContainerCustomEndpointEntry,
			transport_tpg.DataflowCustomEndpointEntryKey:          transport_tpg.DataflowCustomEndpointEntry,
			transport_tpg.IamCredentialsCustomEndpointEntryKey:    transport_tpg.IamCredentialsCustomEndpointEntry,
			transport_tpg.ResourceManagerV3CustomEndpointEntryKey: transport_tpg.ResourceManagerV3CustomEndpointEntry,
			transport_tpg.IAMCustomEndpointEntryKey:               transport_tpg.IAMCustomEndpointEntry,
			transport_tpg.ServiceNetworkingCustomEndpointEntryKey: transport_tpg.ServiceNetworkingCustomEndpointEntry,
			transport_tpg.TagsLocationCustomEndpointEntryKey:      transport_tpg.TagsLocationCustomEndpointEntry,

			// dcl
			transport_tpg.ContainerAwsCustomEndpointEntryKey:   transport_tpg.ContainerAwsCustomEndpointEntry,
			transport_tpg.ContainerAzureCustomEndpointEntryKey: transport_tpg.ContainerAzureCustomEndpointEntry,
		},

		ProviderMetaSchema: map[string]*schema.Schema{
			"module_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},

		DataSourcesMap: DatasourceMap(),
		ResourcesMap:   ResourceMap(),
	}

	provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return ProviderConfigure(ctx, d, provider)
	}

	transport_tpg.ConfigureDCLProvider(provider)

	return provider
}

func DatasourceMap() map[string]*schema.Resource {
	datasourceMap, _ := DatasourceMapWithErrors()
	return datasourceMap
}

func DatasourceMapWithErrors() (map[string]*schema.Resource, error) {
	return mergeResourceMaps(
		handwrittenDatasources,
		generatedIAMDatasources,
		handwrittenIAMDatasources,
	)
}

func ResourceMap() map[string]*schema.Resource {
	resourceMap, _ := ResourceMapWithErrors()
	return resourceMap
}

func ResourceMapWithErrors() (map[string]*schema.Resource, error) {
	return mergeResourceMaps(
		generatedResources,
		handwrittenResources,
		handwrittenIAMResources,
		dclResources,
	)
}

func ProviderConfigure(ctx context.Context, d *schema.ResourceData, p *schema.Provider) (interface{}, diag.Diagnostics) {
	err := transport_tpg.HandleSDKDefaults(d)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	config := transport_tpg.Config{
		Project:             d.Get("project").(string),
		Region:              d.Get("region").(string),
		Zone:                d.Get("zone").(string),
		UserProjectOverride: d.Get("user_project_override").(bool),
		BillingProject:      d.Get("billing_project").(string),
		UserAgent:           p.UserAgent("terraform-provider-google", version.ProviderVersion),
	}

	// opt in extension for adding to the User-Agent header
	if ext := os.Getenv("GOOGLE_TERRAFORM_USERAGENT_EXTENSION"); ext != "" {
		ua := config.UserAgent
		config.UserAgent = fmt.Sprintf("%s %s", ua, ext)
	}

	if v, ok := d.GetOk("request_timeout"); ok {
		var err error
		config.RequestTimeout, err = time.ParseDuration(v.(string))
		if err != nil {
			return nil, diag.FromErr(err)
		}
	}

	if v, ok := d.GetOk("request_reason"); ok {
		config.RequestReason = v.(string)
	}

	// Check for primary credentials in config. Note that if neither is set, ADCs
	// will be used if available.
	if v, ok := d.GetOk("access_token"); ok {
		config.AccessToken = v.(string)
	}

	if v, ok := d.GetOk("credentials"); ok {
		config.Credentials = v.(string)
	}

	// only check environment variables if neither value was set in config- this
	// means config beats env var in all cases.
	if config.AccessToken == "" && config.Credentials == "" {
		config.Credentials = transport_tpg.MultiEnvSearch([]string{
			"GOOGLE_CREDENTIALS",
			"GOOGLE_CLOUD_KEYFILE_JSON",
			"GCLOUD_KEYFILE_JSON",
		})

		config.AccessToken = transport_tpg.MultiEnvSearch([]string{
			"GOOGLE_OAUTH_ACCESS_TOKEN",
		})
	}

	// Set the universe domain to the configured value, if any
	if v, ok := d.GetOk("universe_domain"); ok {
		config.UniverseDomain = v.(string)
	}

	// Configure DCL basePath
	transport_tpg.ProviderDCLConfigure(d, &config)

	// Replace hostname by the universe_domain field.
	if config.UniverseDomain != "" && config.UniverseDomain != "googleapis.com" {
		for key, basePath := range transport_tpg.DefaultBasePaths {
			transport_tpg.DefaultBasePaths[key] = strings.ReplaceAll(basePath, "googleapis.com", config.UniverseDomain)
		}
	}

	err = transport_tpg.SetEndpointDefaults(d)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	transport_tpg.HandleDCLCustomEndpointDefaults(d)

	// Given that impersonate_service_account is a secondary auth method, it has
	// no conflicts to worry about. We pull the env var in a DefaultFunc.
	if v, ok := d.GetOk("impersonate_service_account"); ok {
		config.ImpersonateServiceAccount = v.(string)
	}

	delegates := d.Get("impersonate_service_account_delegates").([]interface{})
	if len(delegates) > 0 {
		config.ImpersonateServiceAccountDelegates = make([]string, len(delegates))
	}
	for i, delegate := range delegates {
		config.ImpersonateServiceAccountDelegates[i] = delegate.(string)
	}

	scopes := d.Get("scopes").([]interface{})
	if len(scopes) > 0 {
		config.Scopes = make([]string, len(scopes))
	}
	for i, scope := range scopes {
		config.Scopes[i] = scope.(string)
	}

	config.DefaultLabels = make(map[string]string)
	defaultLabels := d.Get("default_labels").(map[string]interface{})

	for k, v := range defaultLabels {
		config.DefaultLabels[k] = v.(string)
	}

	// Attribution label is opt-in; if unset, the default for AddTerraformAttributionLabel is false.
	config.AddTerraformAttributionLabel = d.Get("add_terraform_attribution_label").(bool)
	if config.AddTerraformAttributionLabel {
		config.TerraformAttributionLabelAdditionStrategy = transport_tpg.CreateOnlyAttributionStrategy
		if v, ok := d.GetOk("terraform_attribution_label_addition_strategy"); ok {
			config.TerraformAttributionLabelAdditionStrategy = v.(string)
		}
		switch config.TerraformAttributionLabelAdditionStrategy {
		case transport_tpg.CreateOnlyAttributionStrategy, transport_tpg.ProactiveAttributionStrategy:
		default:
			return nil, diag.FromErr(fmt.Errorf("unrecognized terraform_attribution_label_addition_strategy %q", config.TerraformAttributionLabelAdditionStrategy))
		}
	}

	batchCfg, err := transport_tpg.ExpandProviderBatchingConfig(d.Get("batching"))
	if err != nil {
		return nil, diag.FromErr(err)
	}
	config.BatchingConfig = batchCfg

	// Generated products
	config.AccessApprovalBasePath = d.Get("access_approval_custom_endpoint").(string)
	config.AccessContextManagerBasePath = d.Get("access_context_manager_custom_endpoint").(string)
	config.ActiveDirectoryBasePath = d.Get("active_directory_custom_endpoint").(string)
	config.AlloydbBasePath = d.Get("alloydb_custom_endpoint").(string)
	config.ApigeeBasePath = d.Get("apigee_custom_endpoint").(string)
	config.AppEngineBasePath = d.Get("app_engine_custom_endpoint").(string)
	config.ApphubBasePath = d.Get("apphub_custom_endpoint").(string)
	config.ArtifactRegistryBasePath = d.Get("artifact_registry_custom_endpoint").(string)
	config.BeyondcorpBasePath = d.Get("beyondcorp_custom_endpoint").(string)
	config.BiglakeBasePath = d.Get("biglake_custom_endpoint").(string)
	config.BigQueryBasePath = d.Get("big_query_custom_endpoint").(string)
	config.BigqueryAnalyticsHubBasePath = d.Get("bigquery_analytics_hub_custom_endpoint").(string)
	config.BigqueryConnectionBasePath = d.Get("bigquery_connection_custom_endpoint").(string)
	config.BigqueryDatapolicyBasePath = d.Get("bigquery_datapolicy_custom_endpoint").(string)
	config.BigqueryDataTransferBasePath = d.Get("bigquery_data_transfer_custom_endpoint").(string)
	config.BigqueryReservationBasePath = d.Get("bigquery_reservation_custom_endpoint").(string)
	config.BigtableBasePath = d.Get("bigtable_custom_endpoint").(string)
	config.BillingBasePath = d.Get("billing_custom_endpoint").(string)
	config.BinaryAuthorizationBasePath = d.Get("binary_authorization_custom_endpoint").(string)
	config.BlockchainNodeEngineBasePath = d.Get("blockchain_node_engine_custom_endpoint").(string)
	config.CertificateManagerBasePath = d.Get("certificate_manager_custom_endpoint").(string)
	config.CloudAssetBasePath = d.Get("cloud_asset_custom_endpoint").(string)
	config.CloudBuildBasePath = d.Get("cloud_build_custom_endpoint").(string)
	config.Cloudbuildv2BasePath = d.Get("cloudbuildv2_custom_endpoint").(string)
	config.ClouddeployBasePath = d.Get("clouddeploy_custom_endpoint").(string)
	config.ClouddomainsBasePath = d.Get("clouddomains_custom_endpoint").(string)
	config.CloudFunctionsBasePath = d.Get("cloud_functions_custom_endpoint").(string)
	config.Cloudfunctions2BasePath = d.Get("cloudfunctions2_custom_endpoint").(string)
	config.CloudIdentityBasePath = d.Get("cloud_identity_custom_endpoint").(string)
	config.CloudIdsBasePath = d.Get("cloud_ids_custom_endpoint").(string)
	config.CloudQuotasBasePath = d.Get("cloud_quotas_custom_endpoint").(string)
	config.CloudRunBasePath = d.Get("cloud_run_custom_endpoint").(string)
	config.CloudRunV2BasePath = d.Get("cloud_run_v2_custom_endpoint").(string)
	config.CloudSchedulerBasePath = d.Get("cloud_scheduler_custom_endpoint").(string)
	config.CloudTasksBasePath = d.Get("cloud_tasks_custom_endpoint").(string)
	config.ComposerBasePath = d.Get("composer_custom_endpoint").(string)
	config.ComputeBasePath = d.Get("compute_custom_endpoint").(string)
	config.ContainerAnalysisBasePath = d.Get("container_analysis_custom_endpoint").(string)
	config.ContainerAttachedBasePath = d.Get("container_attached_custom_endpoint").(string)
	config.CoreBillingBasePath = d.Get("core_billing_custom_endpoint").(string)
	config.DatabaseMigrationServiceBasePath = d.Get("database_migration_service_custom_endpoint").(string)
	config.DataCatalogBasePath = d.Get("data_catalog_custom_endpoint").(string)
	config.DataFusionBasePath = d.Get("data_fusion_custom_endpoint").(string)
	config.DataLossPreventionBasePath = d.Get("data_loss_prevention_custom_endpoint").(string)
	config.DataPipelineBasePath = d.Get("data_pipeline_custom_endpoint").(string)
	config.DataplexBasePath = d.Get("dataplex_custom_endpoint").(string)
	config.DataprocBasePath = d.Get("dataproc_custom_endpoint").(string)
	config.DataprocMetastoreBasePath = d.Get("dataproc_metastore_custom_endpoint").(string)
	config.DatastoreBasePath = d.Get("datastore_custom_endpoint").(string)
	config.DatastreamBasePath = d.Get("datastream_custom_endpoint").(string)
	config.DeploymentManagerBasePath = d.Get("deployment_manager_custom_endpoint").(string)
	config.DialogflowBasePath = d.Get("dialogflow_custom_endpoint").(string)
	config.DialogflowCXBasePath = d.Get("dialogflow_cx_custom_endpoint").(string)
	config.DiscoveryEngineBasePath = d.Get("discovery_engine_custom_endpoint").(string)
	config.DNSBasePath = d.Get("dns_custom_endpoint").(string)
	config.DocumentAIBasePath = d.Get("document_ai_custom_endpoint").(string)
	config.DocumentAIWarehouseBasePath = d.Get("document_ai_warehouse_custom_endpoint").(string)
	config.EdgecontainerBasePath = d.Get("edgecontainer_custom_endpoint").(string)
	config.EdgenetworkBasePath = d.Get("edgenetwork_custom_endpoint").(string)
	config.EssentialContactsBasePath = d.Get("essential_contacts_custom_endpoint").(string)
	config.FilestoreBasePath = d.Get("filestore_custom_endpoint").(string)
	config.FirebaseAppCheckBasePath = d.Get("firebase_app_check_custom_endpoint").(string)
	config.FirestoreBasePath = d.Get("firestore_custom_endpoint").(string)
	config.GKEBackupBasePath = d.Get("gke_backup_custom_endpoint").(string)
	config.GKEHubBasePath = d.Get("gke_hub_custom_endpoint").(string)
	config.GKEHub2BasePath = d.Get("gke_hub2_custom_endpoint").(string)
	config.GkeonpremBasePath = d.Get("gkeonprem_custom_endpoint").(string)
	config.HealthcareBasePath = d.Get("healthcare_custom_endpoint").(string)
	config.IAM2BasePath = d.Get("iam2_custom_endpoint").(string)
	config.IAMBetaBasePath = d.Get("iam_beta_custom_endpoint").(string)
	config.IAMWorkforcePoolBasePath = d.Get("iam_workforce_pool_custom_endpoint").(string)
	config.IapBasePath = d.Get("iap_custom_endpoint").(string)
	config.IdentityPlatformBasePath = d.Get("identity_platform_custom_endpoint").(string)
	config.IntegrationConnectorsBasePath = d.Get("integration_connectors_custom_endpoint").(string)
	config.IntegrationsBasePath = d.Get("integrations_custom_endpoint").(string)
	config.KMSBasePath = d.Get("kms_custom_endpoint").(string)
	config.LoggingBasePath = d.Get("logging_custom_endpoint").(string)
	config.LookerBasePath = d.Get("looker_custom_endpoint").(string)
	config.MemcacheBasePath = d.Get("memcache_custom_endpoint").(string)
	config.MigrationCenterBasePath = d.Get("migration_center_custom_endpoint").(string)
	config.MLEngineBasePath = d.Get("ml_engine_custom_endpoint").(string)
	config.MonitoringBasePath = d.Get("monitoring_custom_endpoint").(string)
	config.NetappBasePath = d.Get("netapp_custom_endpoint").(string)
	config.NetworkConnectivityBasePath = d.Get("network_connectivity_custom_endpoint").(string)
	config.NetworkManagementBasePath = d.Get("network_management_custom_endpoint").(string)
	config.NetworkSecurityBasePath = d.Get("network_security_custom_endpoint").(string)
	config.NetworkServicesBasePath = d.Get("network_services_custom_endpoint").(string)
	config.NotebooksBasePath = d.Get("notebooks_custom_endpoint").(string)
	config.OrgPolicyBasePath = d.Get("org_policy_custom_endpoint").(string)
	config.OSConfigBasePath = d.Get("os_config_custom_endpoint").(string)
	config.OSLoginBasePath = d.Get("os_login_custom_endpoint").(string)
	config.PrivatecaBasePath = d.Get("privateca_custom_endpoint").(string)
	config.PrivilegedAccessManagerBasePath = d.Get("privileged_access_manager_custom_endpoint").(string)
	config.PublicCABasePath = d.Get("public_ca_custom_endpoint").(string)
	config.PubsubBasePath = d.Get("pubsub_custom_endpoint").(string)
	config.PubsubLiteBasePath = d.Get("pubsub_lite_custom_endpoint").(string)
	config.RedisBasePath = d.Get("redis_custom_endpoint").(string)
	config.ResourceManagerBasePath = d.Get("resource_manager_custom_endpoint").(string)
	config.SecretManagerBasePath = d.Get("secret_manager_custom_endpoint").(string)
	config.SecureSourceManagerBasePath = d.Get("secure_source_manager_custom_endpoint").(string)
	config.SecurityCenterBasePath = d.Get("security_center_custom_endpoint").(string)
	config.SecurityCenterManagementBasePath = d.Get("security_center_management_custom_endpoint").(string)
	config.SecurityCenterV2BasePath = d.Get("security_center_v2_custom_endpoint").(string)
	config.SecuritypostureBasePath = d.Get("securityposture_custom_endpoint").(string)
	config.ServiceManagementBasePath = d.Get("service_management_custom_endpoint").(string)
	config.ServiceNetworkingBasePath = d.Get("service_networking_custom_endpoint").(string)
	config.ServiceUsageBasePath = d.Get("service_usage_custom_endpoint").(string)
	config.SiteVerificationBasePath = d.Get("site_verification_custom_endpoint").(string)
	config.SourceRepoBasePath = d.Get("source_repo_custom_endpoint").(string)
	config.SpannerBasePath = d.Get("spanner_custom_endpoint").(string)
	config.SQLBasePath = d.Get("sql_custom_endpoint").(string)
	config.StorageBasePath = d.Get("storage_custom_endpoint").(string)
	config.StorageInsightsBasePath = d.Get("storage_insights_custom_endpoint").(string)
	config.StorageTransferBasePath = d.Get("storage_transfer_custom_endpoint").(string)
	config.TagsBasePath = d.Get("tags_custom_endpoint").(string)
	config.TPUBasePath = d.Get("tpu_custom_endpoint").(string)
	config.VertexAIBasePath = d.Get("vertex_ai_custom_endpoint").(string)
	config.VmwareengineBasePath = d.Get("vmwareengine_custom_endpoint").(string)
	config.VPCAccessBasePath = d.Get("vpc_access_custom_endpoint").(string)
	config.WorkbenchBasePath = d.Get("workbench_custom_endpoint").(string)
	config.WorkflowsBasePath = d.Get("workflows_custom_endpoint").(string)

	// Handwritten Products / Versioned / Atypical Entries
	config.CloudBillingBasePath = d.Get(transport_tpg.CloudBillingCustomEndpointEntryKey).(string)
	config.ComposerBasePath = d.Get(transport_tpg.ComposerCustomEndpointEntryKey).(string)
	config.ContainerBasePath = d.Get(transport_tpg.ContainerCustomEndpointEntryKey).(string)
	config.DataflowBasePath = d.Get(transport_tpg.DataflowCustomEndpointEntryKey).(string)
	config.IamCredentialsBasePath = d.Get(transport_tpg.IamCredentialsCustomEndpointEntryKey).(string)
	config.ResourceManagerV3BasePath = d.Get(transport_tpg.ResourceManagerV3CustomEndpointEntryKey).(string)
	config.IAMBasePath = d.Get(transport_tpg.IAMCustomEndpointEntryKey).(string)
	config.ServiceUsageBasePath = d.Get(transport_tpg.ServiceUsageCustomEndpointEntryKey).(string)
	config.BigtableAdminBasePath = d.Get(transport_tpg.BigtableAdminCustomEndpointEntryKey).(string)
	config.TagsLocationBasePath = d.Get(transport_tpg.TagsLocationCustomEndpointEntryKey).(string)

	// dcl
	config.ContainerAwsBasePath = d.Get(transport_tpg.ContainerAwsCustomEndpointEntryKey).(string)
	config.ContainerAzureBasePath = d.Get(transport_tpg.ContainerAzureCustomEndpointEntryKey).(string)

	stopCtx, ok := schema.StopContext(ctx)
	if !ok {
		stopCtx = ctx
	}
	if err := config.LoadAndValidate(stopCtx); err != nil {
		return nil, diag.FromErr(err)
	}

	// Verify that universe domains match between credentials and configuration
	if v, ok := d.GetOk("universe_domain"); ok {
		if config.UniverseDomain == "" && v.(string) != "googleapis.com" { // v can't be "", as it wouldn't pass `ok` above
			return nil, diag.FromErr(fmt.Errorf("Universe domain mismatch: '%s' supplied directly to Terraform with no matching universe domain in credentials. Credentials with no 'universe_domain' set are assumed to be in the default universe.", v))
		} else if v.(string) != config.UniverseDomain && !(config.UniverseDomain == "" && v.(string) == "googleapis.com") {
			return nil, diag.FromErr(fmt.Errorf("Universe domain mismatch: '%s' does not match the universe domain '%s' supplied directly to Terraform. The 'universe_domain' provider configuration must match the universe domain supplied by credentials.", config.UniverseDomain, v))
		}
	} else if config.UniverseDomain != "" && config.UniverseDomain != "googleapis.com" {
		return nil, diag.FromErr(fmt.Errorf("Universe domain mismatch: Universe domain '%s' was found in credentials without a corresponding 'universe_domain' provider configuration set. Please set 'universe_domain' to '%s' or use different credentials.", config.UniverseDomain, config.UniverseDomain))
	}

	return &config, nil
}

func mergeResourceMaps(ms ...map[string]*schema.Resource) (map[string]*schema.Resource, error) {
	merged := make(map[string]*schema.Resource)
	duplicates := []string{}

	for _, m := range ms {
		for k, v := range m {
			if _, ok := merged[k]; ok {
				duplicates = append(duplicates, k)
			}

			merged[k] = v
		}
	}

	var err error
	if len(duplicates) > 0 {
		err = fmt.Errorf("saw duplicates in mergeResourceMaps: %v", duplicates)
	}

	return merged, err
}
