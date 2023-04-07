package google

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/version"

	googleoauth "golang.org/x/oauth2/google"
)

const TestEnvVar = "TF_ACC"

// Global MutexKV
var mutexKV = NewMutexKV()

// Provider returns a *schema.Provider.
func Provider() *schema.Provider {

	// The mtls service client gives the type of endpoint (mtls/regular)
	// at client creation. Since we use a shared client for requests we must
	// rewrite the endpoints to be mtls endpoints for the scenario where
	// mtls is enabled.
	if isMtls() {
		// if mtls is enabled switch all default endpoints to use the mtls endpoint
		for key, bp := range DefaultBasePaths {
			DefaultBasePaths[key] = getMtlsEndpoint(bp)
		}
	}

	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"credentials": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validateCredentials,
				ConflictsWith: []string{"access_token"},
			},

			"access_token": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"credentials"},
			},

			"impersonate_service_account": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"impersonate_service_account_delegates": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"billing_project": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"zone": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"scopes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
							ValidateFunc: validateNonNegativeDuration(),
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

			// Generated Products
			"access_approval_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"access_context_manager_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"active_directory_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"alloydb_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"apigee_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"app_engine_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"artifact_registry_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"beyondcorp_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"big_query_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"bigquery_analytics_hub_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"bigquery_connection_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"bigquery_datapolicy_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"bigquery_data_transfer_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"bigquery_reservation_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"bigtable_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"billing_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"binary_authorization_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"certificate_manager_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"cloud_asset_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"cloud_build_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"cloud_functions_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"cloudfunctions2_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"cloud_identity_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"cloud_ids_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"cloud_iot_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"cloud_run_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"cloud_run_v2_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"cloud_scheduler_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"cloud_tasks_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"compute_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"container_analysis_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"container_attached_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"data_catalog_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"data_fusion_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"data_loss_prevention_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"dataplex_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"dataproc_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"dataproc_metastore_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"datastore_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"datastream_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"deployment_manager_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"dialogflow_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"dialogflow_cx_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"dns_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"document_ai_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"essential_contacts_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"filestore_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"firestore_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"game_services_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"gke_backup_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"gke_hub_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"healthcare_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"iam2_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"iam_beta_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"iam_workforce_pool_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"iap_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"identity_platform_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"kms_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"logging_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"memcache_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"ml_engine_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"monitoring_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"network_management_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"network_services_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"notebooks_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"os_config_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"os_login_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"privateca_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"pubsub_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"pubsub_lite_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"redis_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"resource_manager_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"secret_manager_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"security_center_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"service_management_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"service_usage_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"source_repo_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"spanner_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"sql_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"storage_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"storage_transfer_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"tags_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"tpu_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"vertex_ai_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"vpc_access_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},
			"workflows_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
			},

			// Handwritten Products / Versioned / Atypical Entries
			CloudBillingCustomEndpointEntryKey:      CloudBillingCustomEndpointEntry,
			ComposerCustomEndpointEntryKey:          ComposerCustomEndpointEntry,
			ContainerCustomEndpointEntryKey:         ContainerCustomEndpointEntry,
			DataflowCustomEndpointEntryKey:          DataflowCustomEndpointEntry,
			IamCredentialsCustomEndpointEntryKey:    IamCredentialsCustomEndpointEntry,
			ResourceManagerV3CustomEndpointEntryKey: ResourceManagerV3CustomEndpointEntry,
			IAMCustomEndpointEntryKey:               IAMCustomEndpointEntry,
			ServiceNetworkingCustomEndpointEntryKey: ServiceNetworkingCustomEndpointEntry,
			TagsLocationCustomEndpointEntryKey:      TagsLocationCustomEndpointEntry,

			// dcl
			ContainerAwsCustomEndpointEntryKey:   ContainerAwsCustomEndpointEntry,
			ContainerAzureCustomEndpointEntryKey: ContainerAzureCustomEndpointEntry,
		},

		ProviderMetaSchema: map[string]*schema.Schema{
			"module_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			// ####### START datasources ###########
			"google_access_approval_folder_service_account":       DataSourceAccessApprovalFolderServiceAccount(),
			"google_access_approval_organization_service_account": DataSourceAccessApprovalOrganizationServiceAccount(),
			"google_access_approval_project_service_account":      DataSourceAccessApprovalProjectServiceAccount(),
			"google_active_folder":                                DataSourceGoogleActiveFolder(),
			"google_artifact_registry_repository":                 DataSourceArtifactRegistryRepository(),
			"google_app_engine_default_service_account":           DataSourceGoogleAppEngineDefaultServiceAccount(),
			"google_beyondcorp_app_connection":                    DataSourceGoogleBeyondcorpAppConnection(),
			"google_beyondcorp_app_connector":                     DataSourceGoogleBeyondcorpAppConnector(),
			"google_beyondcorp_app_gateway":                       DataSourceGoogleBeyondcorpAppGateway(),
			"google_billing_account":                              DataSourceGoogleBillingAccount(),
			"google_bigquery_default_service_account":             DataSourceGoogleBigqueryDefaultServiceAccount(),
			"google_cloudbuild_trigger":                           DataSourceGoogleCloudBuildTrigger(),
			"google_cloudfunctions_function":                      DataSourceGoogleCloudFunctionsFunction(),
			"google_cloudfunctions2_function":                     DataSourceGoogleCloudFunctions2Function(),
			"google_cloud_identity_groups":                        DataSourceGoogleCloudIdentityGroups(),
			"google_cloud_identity_group_memberships":             DataSourceGoogleCloudIdentityGroupMemberships(),
			"google_cloud_run_locations":                          DataSourceGoogleCloudRunLocations(),
			"google_cloud_run_service":                            DataSourceGoogleCloudRunService(),
			"google_composer_environment":                         DataSourceGoogleComposerEnvironment(),
			"google_composer_image_versions":                      DataSourceGoogleComposerImageVersions(),
			"google_compute_address":                              DataSourceGoogleComputeAddress(),
			"google_compute_addresses":                            DataSourceGoogleComputeAddresses(),
			"google_compute_backend_service":                      DataSourceGoogleComputeBackendService(),
			"google_compute_backend_bucket":                       DataSourceGoogleComputeBackendBucket(),
			"google_compute_default_service_account":              DataSourceGoogleComputeDefaultServiceAccount(),
			"google_compute_disk":                                 DataSourceGoogleComputeDisk(),
			"google_compute_forwarding_rule":                      DataSourceGoogleComputeForwardingRule(),
			"google_compute_global_address":                       DataSourceGoogleComputeGlobalAddress(),
			"google_compute_global_forwarding_rule":               DataSourceGoogleComputeGlobalForwardingRule(),
			"google_compute_ha_vpn_gateway":                       DataSourceGoogleComputeHaVpnGateway(),
			"google_compute_health_check":                         DataSourceGoogleComputeHealthCheck(),
			"google_compute_image":                                DataSourceGoogleComputeImage(),
			"google_compute_instance":                             DataSourceGoogleComputeInstance(),
			"google_compute_instance_group":                       DataSourceGoogleComputeInstanceGroup(),
			"google_compute_instance_group_manager":               DataSourceGoogleComputeInstanceGroupManager(),
			"google_compute_instance_serial_port":                 DataSourceGoogleComputeInstanceSerialPort(),
			"google_compute_instance_template":                    DataSourceGoogleComputeInstanceTemplate(),
			"google_compute_lb_ip_ranges":                         DataSourceGoogleComputeLbIpRanges(),
			"google_compute_network":                              DataSourceGoogleComputeNetwork(),
			"google_compute_network_endpoint_group":               DataSourceGoogleComputeNetworkEndpointGroup(),
			"google_compute_network_peering":                      DataSourceComputeNetworkPeering(),
			"google_compute_node_types":                           DataSourceGoogleComputeNodeTypes(),
			"google_compute_regions":                              DataSourceGoogleComputeRegions(),
			"google_compute_region_network_endpoint_group":        DataSourceGoogleComputeRegionNetworkEndpointGroup(),
			"google_compute_region_instance_group":                DataSourceGoogleComputeRegionInstanceGroup(),
			"google_compute_region_ssl_certificate":               DataSourceGoogleRegionComputeSslCertificate(),
			"google_compute_resource_policy":                      DataSourceGoogleComputeResourcePolicy(),
			"google_compute_router":                               DataSourceGoogleComputeRouter(),
			"google_compute_router_nat":                           DataSourceGoogleComputeRouterNat(),
			"google_compute_router_status":                        DataSourceGoogleComputeRouterStatus(),
			"google_compute_snapshot":                             DataSourceGoogleComputeSnapshot(),
			"google_compute_ssl_certificate":                      DataSourceGoogleComputeSslCertificate(),
			"google_compute_ssl_policy":                           DataSourceGoogleComputeSslPolicy(),
			"google_compute_subnetwork":                           DataSourceGoogleComputeSubnetwork(),
			"google_compute_vpn_gateway":                          DataSourceGoogleComputeVpnGateway(),
			"google_compute_zones":                                DataSourceGoogleComputeZones(),
			"google_container_azure_versions":                     DataSourceGoogleContainerAzureVersions(),
			"google_container_aws_versions":                       DataSourceGoogleContainerAwsVersions(),
			"google_container_attached_versions":                  DataSourceGoogleContainerAttachedVersions(),
			"google_container_attached_install_manifest":          DataSourceGoogleContainerAttachedInstallManifest(),
			"google_container_cluster":                            DataSourceGoogleContainerCluster(),
			"google_container_engine_versions":                    DataSourceGoogleContainerEngineVersions(),
			"google_container_registry_image":                     DataSourceGoogleContainerImage(),
			"google_container_registry_repository":                DataSourceGoogleContainerRepo(),
			"google_dataproc_metastore_service":                   DataSourceDataprocMetastoreService(),
			"google_game_services_game_server_deployment_rollout": DataSourceGameServicesGameServerDeploymentRollout(),
			"google_iam_policy":                                   DataSourceGoogleIamPolicy(),
			"google_iam_role":                                     DataSourceGoogleIamRole(),
			"google_iam_testable_permissions":                     DataSourceGoogleIamTestablePermissions(),
			"google_iap_client":                                   DataSourceGoogleIapClient(),
			"google_kms_crypto_key":                               DataSourceGoogleKmsCryptoKey(),
			"google_kms_crypto_key_version":                       DataSourceGoogleKmsCryptoKeyVersion(),
			"google_kms_key_ring":                                 DataSourceGoogleKmsKeyRing(),
			"google_kms_secret":                                   DataSourceGoogleKmsSecret(),
			"google_kms_secret_ciphertext":                        DataSourceGoogleKmsSecretCiphertext(),
			"google_folder":                                       DataSourceGoogleFolder(),
			"google_folders":                                      DataSourceGoogleFolders(),
			"google_folder_organization_policy":                   DataSourceGoogleFolderOrganizationPolicy(),
			"google_logging_project_cmek_settings":                DataSourceGoogleLoggingProjectCmekSettings(),
			"google_logging_sink":                                 DataSourceGoogleLoggingSink(),
			"google_monitoring_notification_channel":              DataSourceMonitoringNotificationChannel(),
			"google_monitoring_cluster_istio_service":             DataSourceMonitoringServiceClusterIstio(),
			"google_monitoring_istio_canonical_service":           DataSourceMonitoringIstioCanonicalService(),
			"google_monitoring_mesh_istio_service":                DataSourceMonitoringServiceMeshIstio(),
			"google_monitoring_app_engine_service":                DataSourceMonitoringServiceAppEngine(),
			"google_monitoring_uptime_check_ips":                  DataSourceGoogleMonitoringUptimeCheckIps(),
			"google_netblock_ip_ranges":                           DataSourceGoogleNetblockIpRanges(),
			"google_organization":                                 DataSourceGoogleOrganization(),
			"google_privateca_certificate_authority":              DataSourcePrivatecaCertificateAuthority(),
			"google_project":                                      DataSourceGoogleProject(),
			"google_projects":                                     DataSourceGoogleProjects(),
			"google_project_organization_policy":                  DataSourceGoogleProjectOrganizationPolicy(),
			"google_project_service":                              DataSourceGoogleProjectService(),
			"google_pubsub_subscription":                          DataSourceGooglePubsubSubscription(),
			"google_pubsub_topic":                                 DataSourceGooglePubsubTopic(),
			"google_secret_manager_secret":                        DataSourceSecretManagerSecret(),
			"google_secret_manager_secret_version":                DataSourceSecretManagerSecretVersion(),
			"google_secret_manager_secret_version_access":         DataSourceSecretManagerSecretVersionAccess(),
			"google_service_account":                              DataSourceGoogleServiceAccount(),
			"google_service_account_access_token":                 DataSourceGoogleServiceAccountAccessToken(),
			"google_service_account_id_token":                     DataSourceGoogleServiceAccountIdToken(),
			"google_service_account_jwt":                          DataSourceGoogleServiceAccountJwt(),
			"google_service_account_key":                          DataSourceGoogleServiceAccountKey(),
			"google_sourcerepo_repository":                        DataSourceGoogleSourceRepoRepository(),
			"google_spanner_instance":                             DataSourceSpannerInstance(),
			"google_sql_ca_certs":                                 DataSourceGoogleSQLCaCerts(),
			"google_sql_backup_run":                               DataSourceSqlBackupRun(),
			"google_sql_databases":                                DataSourceSqlDatabases(),
			"google_sql_database":                                 DataSourceSqlDatabase(),
			"google_sql_database_instance":                        DataSourceSqlDatabaseInstance(),
			"google_sql_database_instances":                       DataSourceSqlDatabaseInstances(),
			"google_service_networking_peered_dns_domain":         DataSourceGoogleServiceNetworkingPeeredDNSDomain(),
			"google_storage_bucket":                               DataSourceGoogleStorageBucket(),
			"google_storage_bucket_object":                        DataSourceGoogleStorageBucketObject(),
			"google_storage_bucket_object_content":                DataSourceGoogleStorageBucketObjectContent(),
			"google_storage_object_signed_url":                    DataSourceGoogleSignedUrl(),
			"google_storage_project_service_account":              DataSourceGoogleStorageProjectServiceAccount(),
			"google_storage_transfer_project_service_account":     DataSourceGoogleStorageTransferProjectServiceAccount(),
			"google_tags_tag_key":                                 DataSourceGoogleTagsTagKey(),
			"google_tags_tag_value":                               DataSourceGoogleTagsTagValue(),
			"google_tpu_tensorflow_versions":                      DataSourceTpuTensorflowVersions(),
			"google_vpc_access_connector":                         DataSourceVPCAccessConnector(),
			"google_redis_instance":                               DataSourceGoogleRedisInstance(),
			// ####### END datasources ###########
		},
		ResourcesMap: ResourceMap(),
	}

	provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return providerConfigure(ctx, d, provider)
	}

	ConfigureDCLProvider(provider)

	return provider
}

// Generated resources: 276
// Generated IAM resources: 186
// Total generated resources: 462
func ResourceMap() map[string]*schema.Resource {
	resourceMap, _ := ResourceMapWithErrors()
	return resourceMap
}

func ResourceMapWithErrors() (map[string]*schema.Resource, error) {
	return mergeResourceMaps(
		map[string]*schema.Resource{
			"google_folder_access_approval_settings":                       ResourceAccessApprovalFolderSettings(),
			"google_organization_access_approval_settings":                 ResourceAccessApprovalOrganizationSettings(),
			"google_project_access_approval_settings":                      ResourceAccessApprovalProjectSettings(),
			"google_access_context_manager_access_level":                   ResourceAccessContextManagerAccessLevel(),
			"google_access_context_manager_access_level_condition":         ResourceAccessContextManagerAccessLevelCondition(),
			"google_access_context_manager_access_levels":                  ResourceAccessContextManagerAccessLevels(),
			"google_access_context_manager_access_policy":                  ResourceAccessContextManagerAccessPolicy(),
			"google_access_context_manager_access_policy_iam_binding":      ResourceIamBinding(AccessContextManagerAccessPolicyIamSchema, AccessContextManagerAccessPolicyIamUpdaterProducer, AccessContextManagerAccessPolicyIdParseFunc),
			"google_access_context_manager_access_policy_iam_member":       ResourceIamMember(AccessContextManagerAccessPolicyIamSchema, AccessContextManagerAccessPolicyIamUpdaterProducer, AccessContextManagerAccessPolicyIdParseFunc),
			"google_access_context_manager_access_policy_iam_policy":       ResourceIamPolicy(AccessContextManagerAccessPolicyIamSchema, AccessContextManagerAccessPolicyIamUpdaterProducer, AccessContextManagerAccessPolicyIdParseFunc),
			"google_access_context_manager_authorized_orgs_desc":           ResourceAccessContextManagerAuthorizedOrgsDesc(),
			"google_access_context_manager_gcp_user_access_binding":        ResourceAccessContextManagerGcpUserAccessBinding(),
			"google_access_context_manager_service_perimeter":              ResourceAccessContextManagerServicePerimeter(),
			"google_access_context_manager_service_perimeter_resource":     ResourceAccessContextManagerServicePerimeterResource(),
			"google_access_context_manager_service_perimeters":             ResourceAccessContextManagerServicePerimeters(),
			"google_active_directory_domain":                               ResourceActiveDirectoryDomain(),
			"google_active_directory_domain_trust":                         ResourceActiveDirectoryDomainTrust(),
			"google_alloydb_backup":                                        ResourceAlloydbBackup(),
			"google_alloydb_cluster":                                       ResourceAlloydbCluster(),
			"google_alloydb_instance":                                      ResourceAlloydbInstance(),
			"google_apigee_addons_config":                                  ResourceApigeeAddonsConfig(),
			"google_apigee_endpoint_attachment":                            ResourceApigeeEndpointAttachment(),
			"google_apigee_env_keystore":                                   ResourceApigeeEnvKeystore(),
			"google_apigee_env_references":                                 ResourceApigeeEnvReferences(),
			"google_apigee_envgroup":                                       ResourceApigeeEnvgroup(),
			"google_apigee_envgroup_attachment":                            ResourceApigeeEnvgroupAttachment(),
			"google_apigee_environment":                                    ResourceApigeeEnvironment(),
			"google_apigee_environment_iam_binding":                        ResourceIamBinding(ApigeeEnvironmentIamSchema, ApigeeEnvironmentIamUpdaterProducer, ApigeeEnvironmentIdParseFunc),
			"google_apigee_environment_iam_member":                         ResourceIamMember(ApigeeEnvironmentIamSchema, ApigeeEnvironmentIamUpdaterProducer, ApigeeEnvironmentIdParseFunc),
			"google_apigee_environment_iam_policy":                         ResourceIamPolicy(ApigeeEnvironmentIamSchema, ApigeeEnvironmentIamUpdaterProducer, ApigeeEnvironmentIdParseFunc),
			"google_apigee_instance":                                       ResourceApigeeInstance(),
			"google_apigee_instance_attachment":                            ResourceApigeeInstanceAttachment(),
			"google_apigee_keystores_aliases_self_signed_cert":             ResourceApigeeKeystoresAliasesSelfSignedCert(),
			"google_apigee_nat_address":                                    ResourceApigeeNatAddress(),
			"google_apigee_organization":                                   ResourceApigeeOrganization(),
			"google_apigee_sync_authorization":                             ResourceApigeeSyncAuthorization(),
			"google_app_engine_application_url_dispatch_rules":             ResourceAppEngineApplicationUrlDispatchRules(),
			"google_app_engine_domain_mapping":                             ResourceAppEngineDomainMapping(),
			"google_app_engine_firewall_rule":                              ResourceAppEngineFirewallRule(),
			"google_app_engine_flexible_app_version":                       ResourceAppEngineFlexibleAppVersion(),
			"google_app_engine_service_network_settings":                   ResourceAppEngineServiceNetworkSettings(),
			"google_app_engine_service_split_traffic":                      ResourceAppEngineServiceSplitTraffic(),
			"google_app_engine_standard_app_version":                       ResourceAppEngineStandardAppVersion(),
			"google_artifact_registry_repository":                          ResourceArtifactRegistryRepository(),
			"google_artifact_registry_repository_iam_binding":              ResourceIamBinding(ArtifactRegistryRepositoryIamSchema, ArtifactRegistryRepositoryIamUpdaterProducer, ArtifactRegistryRepositoryIdParseFunc),
			"google_artifact_registry_repository_iam_member":               ResourceIamMember(ArtifactRegistryRepositoryIamSchema, ArtifactRegistryRepositoryIamUpdaterProducer, ArtifactRegistryRepositoryIdParseFunc),
			"google_artifact_registry_repository_iam_policy":               ResourceIamPolicy(ArtifactRegistryRepositoryIamSchema, ArtifactRegistryRepositoryIamUpdaterProducer, ArtifactRegistryRepositoryIdParseFunc),
			"google_beyondcorp_app_connection":                             ResourceBeyondcorpAppConnection(),
			"google_beyondcorp_app_connector":                              ResourceBeyondcorpAppConnector(),
			"google_beyondcorp_app_gateway":                                ResourceBeyondcorpAppGateway(),
			"google_bigquery_dataset":                                      ResourceBigQueryDataset(),
			"google_bigquery_dataset_access":                               ResourceBigQueryDatasetAccess(),
			"google_bigquery_job":                                          ResourceBigQueryJob(),
			"google_bigquery_routine":                                      ResourceBigQueryRoutine(),
			"google_bigquery_table_iam_binding":                            ResourceIamBinding(BigQueryTableIamSchema, BigQueryTableIamUpdaterProducer, BigQueryTableIdParseFunc),
			"google_bigquery_table_iam_member":                             ResourceIamMember(BigQueryTableIamSchema, BigQueryTableIamUpdaterProducer, BigQueryTableIdParseFunc),
			"google_bigquery_table_iam_policy":                             ResourceIamPolicy(BigQueryTableIamSchema, BigQueryTableIamUpdaterProducer, BigQueryTableIdParseFunc),
			"google_bigquery_analytics_hub_data_exchange":                  ResourceBigqueryAnalyticsHubDataExchange(),
			"google_bigquery_analytics_hub_data_exchange_iam_binding":      ResourceIamBinding(BigqueryAnalyticsHubDataExchangeIamSchema, BigqueryAnalyticsHubDataExchangeIamUpdaterProducer, BigqueryAnalyticsHubDataExchangeIdParseFunc),
			"google_bigquery_analytics_hub_data_exchange_iam_member":       ResourceIamMember(BigqueryAnalyticsHubDataExchangeIamSchema, BigqueryAnalyticsHubDataExchangeIamUpdaterProducer, BigqueryAnalyticsHubDataExchangeIdParseFunc),
			"google_bigquery_analytics_hub_data_exchange_iam_policy":       ResourceIamPolicy(BigqueryAnalyticsHubDataExchangeIamSchema, BigqueryAnalyticsHubDataExchangeIamUpdaterProducer, BigqueryAnalyticsHubDataExchangeIdParseFunc),
			"google_bigquery_analytics_hub_listing":                        ResourceBigqueryAnalyticsHubListing(),
			"google_bigquery_analytics_hub_listing_iam_binding":            ResourceIamBinding(BigqueryAnalyticsHubListingIamSchema, BigqueryAnalyticsHubListingIamUpdaterProducer, BigqueryAnalyticsHubListingIdParseFunc),
			"google_bigquery_analytics_hub_listing_iam_member":             ResourceIamMember(BigqueryAnalyticsHubListingIamSchema, BigqueryAnalyticsHubListingIamUpdaterProducer, BigqueryAnalyticsHubListingIdParseFunc),
			"google_bigquery_analytics_hub_listing_iam_policy":             ResourceIamPolicy(BigqueryAnalyticsHubListingIamSchema, BigqueryAnalyticsHubListingIamUpdaterProducer, BigqueryAnalyticsHubListingIdParseFunc),
			"google_bigquery_connection":                                   ResourceBigqueryConnectionConnection(),
			"google_bigquery_connection_iam_binding":                       ResourceIamBinding(BigqueryConnectionConnectionIamSchema, BigqueryConnectionConnectionIamUpdaterProducer, BigqueryConnectionConnectionIdParseFunc),
			"google_bigquery_connection_iam_member":                        ResourceIamMember(BigqueryConnectionConnectionIamSchema, BigqueryConnectionConnectionIamUpdaterProducer, BigqueryConnectionConnectionIdParseFunc),
			"google_bigquery_connection_iam_policy":                        ResourceIamPolicy(BigqueryConnectionConnectionIamSchema, BigqueryConnectionConnectionIamUpdaterProducer, BigqueryConnectionConnectionIdParseFunc),
			"google_bigquery_datapolicy_data_policy":                       ResourceBigqueryDatapolicyDataPolicy(),
			"google_bigquery_datapolicy_data_policy_iam_binding":           ResourceIamBinding(BigqueryDatapolicyDataPolicyIamSchema, BigqueryDatapolicyDataPolicyIamUpdaterProducer, BigqueryDatapolicyDataPolicyIdParseFunc),
			"google_bigquery_datapolicy_data_policy_iam_member":            ResourceIamMember(BigqueryDatapolicyDataPolicyIamSchema, BigqueryDatapolicyDataPolicyIamUpdaterProducer, BigqueryDatapolicyDataPolicyIdParseFunc),
			"google_bigquery_datapolicy_data_policy_iam_policy":            ResourceIamPolicy(BigqueryDatapolicyDataPolicyIamSchema, BigqueryDatapolicyDataPolicyIamUpdaterProducer, BigqueryDatapolicyDataPolicyIdParseFunc),
			"google_bigquery_data_transfer_config":                         ResourceBigqueryDataTransferConfig(),
			"google_bigquery_capacity_commitment":                          ResourceBigqueryReservationCapacityCommitment(),
			"google_bigquery_reservation":                                  ResourceBigqueryReservationReservation(),
			"google_bigtable_app_profile":                                  ResourceBigtableAppProfile(),
			"google_billing_budget":                                        ResourceBillingBudget(),
			"google_binary_authorization_attestor":                         ResourceBinaryAuthorizationAttestor(),
			"google_binary_authorization_attestor_iam_binding":             ResourceIamBinding(BinaryAuthorizationAttestorIamSchema, BinaryAuthorizationAttestorIamUpdaterProducer, BinaryAuthorizationAttestorIdParseFunc),
			"google_binary_authorization_attestor_iam_member":              ResourceIamMember(BinaryAuthorizationAttestorIamSchema, BinaryAuthorizationAttestorIamUpdaterProducer, BinaryAuthorizationAttestorIdParseFunc),
			"google_binary_authorization_attestor_iam_policy":              ResourceIamPolicy(BinaryAuthorizationAttestorIamSchema, BinaryAuthorizationAttestorIamUpdaterProducer, BinaryAuthorizationAttestorIdParseFunc),
			"google_binary_authorization_policy":                           ResourceBinaryAuthorizationPolicy(),
			"google_certificate_manager_certificate":                       ResourceCertificateManagerCertificate(),
			"google_certificate_manager_certificate_map":                   ResourceCertificateManagerCertificateMap(),
			"google_certificate_manager_certificate_map_entry":             ResourceCertificateManagerCertificateMapEntry(),
			"google_certificate_manager_dns_authorization":                 ResourceCertificateManagerDnsAuthorization(),
			"google_cloud_asset_folder_feed":                               ResourceCloudAssetFolderFeed(),
			"google_cloud_asset_organization_feed":                         ResourceCloudAssetOrganizationFeed(),
			"google_cloud_asset_project_feed":                              ResourceCloudAssetProjectFeed(),
			"google_cloudbuild_bitbucket_server_config":                    ResourceCloudBuildBitbucketServerConfig(),
			"google_cloudbuild_trigger":                                    ResourceCloudBuildTrigger(),
			"google_cloudfunctions_function_iam_binding":                   ResourceIamBinding(CloudFunctionsCloudFunctionIamSchema, CloudFunctionsCloudFunctionIamUpdaterProducer, CloudFunctionsCloudFunctionIdParseFunc),
			"google_cloudfunctions_function_iam_member":                    ResourceIamMember(CloudFunctionsCloudFunctionIamSchema, CloudFunctionsCloudFunctionIamUpdaterProducer, CloudFunctionsCloudFunctionIdParseFunc),
			"google_cloudfunctions_function_iam_policy":                    ResourceIamPolicy(CloudFunctionsCloudFunctionIamSchema, CloudFunctionsCloudFunctionIamUpdaterProducer, CloudFunctionsCloudFunctionIdParseFunc),
			"google_cloudfunctions2_function":                              ResourceCloudfunctions2function(),
			"google_cloudfunctions2_function_iam_binding":                  ResourceIamBinding(Cloudfunctions2functionIamSchema, Cloudfunctions2functionIamUpdaterProducer, Cloudfunctions2functionIdParseFunc),
			"google_cloudfunctions2_function_iam_member":                   ResourceIamMember(Cloudfunctions2functionIamSchema, Cloudfunctions2functionIamUpdaterProducer, Cloudfunctions2functionIdParseFunc),
			"google_cloudfunctions2_function_iam_policy":                   ResourceIamPolicy(Cloudfunctions2functionIamSchema, Cloudfunctions2functionIamUpdaterProducer, Cloudfunctions2functionIdParseFunc),
			"google_cloud_identity_group":                                  ResourceCloudIdentityGroup(),
			"google_cloud_identity_group_membership":                       ResourceCloudIdentityGroupMembership(),
			"google_cloud_ids_endpoint":                                    ResourceCloudIdsEndpoint(),
			"google_cloudiot_device":                                       ResourceCloudIotDevice(),
			"google_cloudiot_registry":                                     ResourceCloudIotDeviceRegistry(),
			"google_cloudiot_registry_iam_binding":                         ResourceIamBinding(CloudIotDeviceRegistryIamSchema, CloudIotDeviceRegistryIamUpdaterProducer, CloudIotDeviceRegistryIdParseFunc),
			"google_cloudiot_registry_iam_member":                          ResourceIamMember(CloudIotDeviceRegistryIamSchema, CloudIotDeviceRegistryIamUpdaterProducer, CloudIotDeviceRegistryIdParseFunc),
			"google_cloudiot_registry_iam_policy":                          ResourceIamPolicy(CloudIotDeviceRegistryIamSchema, CloudIotDeviceRegistryIamUpdaterProducer, CloudIotDeviceRegistryIdParseFunc),
			"google_cloud_run_domain_mapping":                              ResourceCloudRunDomainMapping(),
			"google_cloud_run_service":                                     ResourceCloudRunService(),
			"google_cloud_run_service_iam_binding":                         ResourceIamBinding(CloudRunServiceIamSchema, CloudRunServiceIamUpdaterProducer, CloudRunServiceIdParseFunc),
			"google_cloud_run_service_iam_member":                          ResourceIamMember(CloudRunServiceIamSchema, CloudRunServiceIamUpdaterProducer, CloudRunServiceIdParseFunc),
			"google_cloud_run_service_iam_policy":                          ResourceIamPolicy(CloudRunServiceIamSchema, CloudRunServiceIamUpdaterProducer, CloudRunServiceIdParseFunc),
			"google_cloud_run_v2_job":                                      ResourceCloudRunV2Job(),
			"google_cloud_run_v2_job_iam_binding":                          ResourceIamBinding(CloudRunV2JobIamSchema, CloudRunV2JobIamUpdaterProducer, CloudRunV2JobIdParseFunc),
			"google_cloud_run_v2_job_iam_member":                           ResourceIamMember(CloudRunV2JobIamSchema, CloudRunV2JobIamUpdaterProducer, CloudRunV2JobIdParseFunc),
			"google_cloud_run_v2_job_iam_policy":                           ResourceIamPolicy(CloudRunV2JobIamSchema, CloudRunV2JobIamUpdaterProducer, CloudRunV2JobIdParseFunc),
			"google_cloud_run_v2_service":                                  ResourceCloudRunV2Service(),
			"google_cloud_run_v2_service_iam_binding":                      ResourceIamBinding(CloudRunV2ServiceIamSchema, CloudRunV2ServiceIamUpdaterProducer, CloudRunV2ServiceIdParseFunc),
			"google_cloud_run_v2_service_iam_member":                       ResourceIamMember(CloudRunV2ServiceIamSchema, CloudRunV2ServiceIamUpdaterProducer, CloudRunV2ServiceIdParseFunc),
			"google_cloud_run_v2_service_iam_policy":                       ResourceIamPolicy(CloudRunV2ServiceIamSchema, CloudRunV2ServiceIamUpdaterProducer, CloudRunV2ServiceIdParseFunc),
			"google_cloud_scheduler_job":                                   ResourceCloudSchedulerJob(),
			"google_cloud_tasks_queue":                                     ResourceCloudTasksQueue(),
			"google_cloud_tasks_queue_iam_binding":                         ResourceIamBinding(CloudTasksQueueIamSchema, CloudTasksQueueIamUpdaterProducer, CloudTasksQueueIdParseFunc),
			"google_cloud_tasks_queue_iam_member":                          ResourceIamMember(CloudTasksQueueIamSchema, CloudTasksQueueIamUpdaterProducer, CloudTasksQueueIdParseFunc),
			"google_cloud_tasks_queue_iam_policy":                          ResourceIamPolicy(CloudTasksQueueIamSchema, CloudTasksQueueIamUpdaterProducer, CloudTasksQueueIdParseFunc),
			"google_compute_address":                                       ResourceComputeAddress(),
			"google_compute_autoscaler":                                    ResourceComputeAutoscaler(),
			"google_compute_backend_bucket":                                ResourceComputeBackendBucket(),
			"google_compute_backend_bucket_signed_url_key":                 ResourceComputeBackendBucketSignedUrlKey(),
			"google_compute_backend_service":                               ResourceComputeBackendService(),
			"google_compute_backend_service_signed_url_key":                ResourceComputeBackendServiceSignedUrlKey(),
			"google_compute_disk":                                          ResourceComputeDisk(),
			"google_compute_disk_iam_binding":                              ResourceIamBinding(ComputeDiskIamSchema, ComputeDiskIamUpdaterProducer, ComputeDiskIdParseFunc),
			"google_compute_disk_iam_member":                               ResourceIamMember(ComputeDiskIamSchema, ComputeDiskIamUpdaterProducer, ComputeDiskIdParseFunc),
			"google_compute_disk_iam_policy":                               ResourceIamPolicy(ComputeDiskIamSchema, ComputeDiskIamUpdaterProducer, ComputeDiskIdParseFunc),
			"google_compute_disk_resource_policy_attachment":               ResourceComputeDiskResourcePolicyAttachment(),
			"google_compute_external_vpn_gateway":                          ResourceComputeExternalVpnGateway(),
			"google_compute_firewall":                                      ResourceComputeFirewall(),
			"google_compute_forwarding_rule":                               ResourceComputeForwardingRule(),
			"google_compute_global_address":                                ResourceComputeGlobalAddress(),
			"google_compute_global_forwarding_rule":                        ResourceComputeGlobalForwardingRule(),
			"google_compute_global_network_endpoint":                       ResourceComputeGlobalNetworkEndpoint(),
			"google_compute_global_network_endpoint_group":                 ResourceComputeGlobalNetworkEndpointGroup(),
			"google_compute_ha_vpn_gateway":                                ResourceComputeHaVpnGateway(),
			"google_compute_health_check":                                  ResourceComputeHealthCheck(),
			"google_compute_http_health_check":                             ResourceComputeHttpHealthCheck(),
			"google_compute_https_health_check":                            ResourceComputeHttpsHealthCheck(),
			"google_compute_image":                                         ResourceComputeImage(),
			"google_compute_image_iam_binding":                             ResourceIamBinding(ComputeImageIamSchema, ComputeImageIamUpdaterProducer, ComputeImageIdParseFunc),
			"google_compute_image_iam_member":                              ResourceIamMember(ComputeImageIamSchema, ComputeImageIamUpdaterProducer, ComputeImageIdParseFunc),
			"google_compute_image_iam_policy":                              ResourceIamPolicy(ComputeImageIamSchema, ComputeImageIamUpdaterProducer, ComputeImageIdParseFunc),
			"google_compute_instance_iam_binding":                          ResourceIamBinding(ComputeInstanceIamSchema, ComputeInstanceIamUpdaterProducer, ComputeInstanceIdParseFunc),
			"google_compute_instance_iam_member":                           ResourceIamMember(ComputeInstanceIamSchema, ComputeInstanceIamUpdaterProducer, ComputeInstanceIdParseFunc),
			"google_compute_instance_iam_policy":                           ResourceIamPolicy(ComputeInstanceIamSchema, ComputeInstanceIamUpdaterProducer, ComputeInstanceIdParseFunc),
			"google_compute_instance_group_named_port":                     ResourceComputeInstanceGroupNamedPort(),
			"google_compute_interconnect_attachment":                       ResourceComputeInterconnectAttachment(),
			"google_compute_managed_ssl_certificate":                       ResourceComputeManagedSslCertificate(),
			"google_compute_network":                                       ResourceComputeNetwork(),
			"google_compute_network_endpoint":                              ResourceComputeNetworkEndpoint(),
			"google_compute_network_endpoint_group":                        ResourceComputeNetworkEndpointGroup(),
			"google_compute_network_peering_routes_config":                 ResourceComputeNetworkPeeringRoutesConfig(),
			"google_compute_node_group":                                    ResourceComputeNodeGroup(),
			"google_compute_node_template":                                 ResourceComputeNodeTemplate(),
			"google_compute_packet_mirroring":                              ResourceComputePacketMirroring(),
			"google_compute_per_instance_config":                           ResourceComputePerInstanceConfig(),
			"google_compute_region_autoscaler":                             ResourceComputeRegionAutoscaler(),
			"google_compute_region_backend_service":                        ResourceComputeRegionBackendService(),
			"google_compute_region_disk":                                   ResourceComputeRegionDisk(),
			"google_compute_region_disk_iam_binding":                       ResourceIamBinding(ComputeRegionDiskIamSchema, ComputeRegionDiskIamUpdaterProducer, ComputeRegionDiskIdParseFunc),
			"google_compute_region_disk_iam_member":                        ResourceIamMember(ComputeRegionDiskIamSchema, ComputeRegionDiskIamUpdaterProducer, ComputeRegionDiskIdParseFunc),
			"google_compute_region_disk_iam_policy":                        ResourceIamPolicy(ComputeRegionDiskIamSchema, ComputeRegionDiskIamUpdaterProducer, ComputeRegionDiskIdParseFunc),
			"google_compute_region_disk_resource_policy_attachment":        ResourceComputeRegionDiskResourcePolicyAttachment(),
			"google_compute_region_health_check":                           ResourceComputeRegionHealthCheck(),
			"google_compute_region_network_endpoint_group":                 ResourceComputeRegionNetworkEndpointGroup(),
			"google_compute_region_per_instance_config":                    ResourceComputeRegionPerInstanceConfig(),
			"google_compute_region_ssl_certificate":                        ResourceComputeRegionSslCertificate(),
			"google_compute_region_target_http_proxy":                      ResourceComputeRegionTargetHttpProxy(),
			"google_compute_region_target_https_proxy":                     ResourceComputeRegionTargetHttpsProxy(),
			"google_compute_region_target_tcp_proxy":                       ResourceComputeRegionTargetTcpProxy(),
			"google_compute_region_url_map":                                ResourceComputeRegionUrlMap(),
			"google_compute_reservation":                                   ResourceComputeReservation(),
			"google_compute_resource_policy":                               ResourceComputeResourcePolicy(),
			"google_compute_route":                                         ResourceComputeRoute(),
			"google_compute_router":                                        ResourceComputeRouter(),
			"google_compute_router_peer":                                   ResourceComputeRouterBgpPeer(),
			"google_compute_router_nat":                                    ResourceComputeRouterNat(),
			"google_compute_service_attachment":                            ResourceComputeServiceAttachment(),
			"google_compute_snapshot":                                      ResourceComputeSnapshot(),
			"google_compute_snapshot_iam_binding":                          ResourceIamBinding(ComputeSnapshotIamSchema, ComputeSnapshotIamUpdaterProducer, ComputeSnapshotIdParseFunc),
			"google_compute_snapshot_iam_member":                           ResourceIamMember(ComputeSnapshotIamSchema, ComputeSnapshotIamUpdaterProducer, ComputeSnapshotIdParseFunc),
			"google_compute_snapshot_iam_policy":                           ResourceIamPolicy(ComputeSnapshotIamSchema, ComputeSnapshotIamUpdaterProducer, ComputeSnapshotIdParseFunc),
			"google_compute_ssl_certificate":                               ResourceComputeSslCertificate(),
			"google_compute_ssl_policy":                                    ResourceComputeSslPolicy(),
			"google_compute_subnetwork":                                    ResourceComputeSubnetwork(),
			"google_compute_subnetwork_iam_binding":                        ResourceIamBinding(ComputeSubnetworkIamSchema, ComputeSubnetworkIamUpdaterProducer, ComputeSubnetworkIdParseFunc),
			"google_compute_subnetwork_iam_member":                         ResourceIamMember(ComputeSubnetworkIamSchema, ComputeSubnetworkIamUpdaterProducer, ComputeSubnetworkIdParseFunc),
			"google_compute_subnetwork_iam_policy":                         ResourceIamPolicy(ComputeSubnetworkIamSchema, ComputeSubnetworkIamUpdaterProducer, ComputeSubnetworkIdParseFunc),
			"google_compute_target_grpc_proxy":                             ResourceComputeTargetGrpcProxy(),
			"google_compute_target_http_proxy":                             ResourceComputeTargetHttpProxy(),
			"google_compute_target_https_proxy":                            ResourceComputeTargetHttpsProxy(),
			"google_compute_target_instance":                               ResourceComputeTargetInstance(),
			"google_compute_target_ssl_proxy":                              ResourceComputeTargetSslProxy(),
			"google_compute_target_tcp_proxy":                              ResourceComputeTargetTcpProxy(),
			"google_compute_url_map":                                       ResourceComputeUrlMap(),
			"google_compute_vpn_gateway":                                   ResourceComputeVpnGateway(),
			"google_compute_vpn_tunnel":                                    ResourceComputeVpnTunnel(),
			"google_container_analysis_note":                               ResourceContainerAnalysisNote(),
			"google_container_analysis_occurrence":                         ResourceContainerAnalysisOccurrence(),
			"google_container_attached_cluster":                            ResourceContainerAttachedCluster(),
			"google_data_catalog_entry":                                    ResourceDataCatalogEntry(),
			"google_data_catalog_entry_group":                              ResourceDataCatalogEntryGroup(),
			"google_data_catalog_entry_group_iam_binding":                  ResourceIamBinding(DataCatalogEntryGroupIamSchema, DataCatalogEntryGroupIamUpdaterProducer, DataCatalogEntryGroupIdParseFunc),
			"google_data_catalog_entry_group_iam_member":                   ResourceIamMember(DataCatalogEntryGroupIamSchema, DataCatalogEntryGroupIamUpdaterProducer, DataCatalogEntryGroupIdParseFunc),
			"google_data_catalog_entry_group_iam_policy":                   ResourceIamPolicy(DataCatalogEntryGroupIamSchema, DataCatalogEntryGroupIamUpdaterProducer, DataCatalogEntryGroupIdParseFunc),
			"google_data_catalog_policy_tag":                               ResourceDataCatalogPolicyTag(),
			"google_data_catalog_policy_tag_iam_binding":                   ResourceIamBinding(DataCatalogPolicyTagIamSchema, DataCatalogPolicyTagIamUpdaterProducer, DataCatalogPolicyTagIdParseFunc),
			"google_data_catalog_policy_tag_iam_member":                    ResourceIamMember(DataCatalogPolicyTagIamSchema, DataCatalogPolicyTagIamUpdaterProducer, DataCatalogPolicyTagIdParseFunc),
			"google_data_catalog_policy_tag_iam_policy":                    ResourceIamPolicy(DataCatalogPolicyTagIamSchema, DataCatalogPolicyTagIamUpdaterProducer, DataCatalogPolicyTagIdParseFunc),
			"google_data_catalog_tag":                                      ResourceDataCatalogTag(),
			"google_data_catalog_tag_template":                             ResourceDataCatalogTagTemplate(),
			"google_data_catalog_tag_template_iam_binding":                 ResourceIamBinding(DataCatalogTagTemplateIamSchema, DataCatalogTagTemplateIamUpdaterProducer, DataCatalogTagTemplateIdParseFunc),
			"google_data_catalog_tag_template_iam_member":                  ResourceIamMember(DataCatalogTagTemplateIamSchema, DataCatalogTagTemplateIamUpdaterProducer, DataCatalogTagTemplateIdParseFunc),
			"google_data_catalog_tag_template_iam_policy":                  ResourceIamPolicy(DataCatalogTagTemplateIamSchema, DataCatalogTagTemplateIamUpdaterProducer, DataCatalogTagTemplateIdParseFunc),
			"google_data_catalog_taxonomy":                                 ResourceDataCatalogTaxonomy(),
			"google_data_catalog_taxonomy_iam_binding":                     ResourceIamBinding(DataCatalogTaxonomyIamSchema, DataCatalogTaxonomyIamUpdaterProducer, DataCatalogTaxonomyIdParseFunc),
			"google_data_catalog_taxonomy_iam_member":                      ResourceIamMember(DataCatalogTaxonomyIamSchema, DataCatalogTaxonomyIamUpdaterProducer, DataCatalogTaxonomyIdParseFunc),
			"google_data_catalog_taxonomy_iam_policy":                      ResourceIamPolicy(DataCatalogTaxonomyIamSchema, DataCatalogTaxonomyIamUpdaterProducer, DataCatalogTaxonomyIdParseFunc),
			"google_data_fusion_instance":                                  ResourceDataFusionInstance(),
			"google_data_fusion_instance_iam_binding":                      ResourceIamBinding(DataFusionInstanceIamSchema, DataFusionInstanceIamUpdaterProducer, DataFusionInstanceIdParseFunc),
			"google_data_fusion_instance_iam_member":                       ResourceIamMember(DataFusionInstanceIamSchema, DataFusionInstanceIamUpdaterProducer, DataFusionInstanceIdParseFunc),
			"google_data_fusion_instance_iam_policy":                       ResourceIamPolicy(DataFusionInstanceIamSchema, DataFusionInstanceIamUpdaterProducer, DataFusionInstanceIdParseFunc),
			"google_data_loss_prevention_deidentify_template":              ResourceDataLossPreventionDeidentifyTemplate(),
			"google_data_loss_prevention_inspect_template":                 ResourceDataLossPreventionInspectTemplate(),
			"google_data_loss_prevention_job_trigger":                      ResourceDataLossPreventionJobTrigger(),
			"google_data_loss_prevention_stored_info_type":                 ResourceDataLossPreventionStoredInfoType(),
			"google_dataplex_asset_iam_binding":                            ResourceIamBinding(DataplexAssetIamSchema, DataplexAssetIamUpdaterProducer, DataplexAssetIdParseFunc),
			"google_dataplex_asset_iam_member":                             ResourceIamMember(DataplexAssetIamSchema, DataplexAssetIamUpdaterProducer, DataplexAssetIdParseFunc),
			"google_dataplex_asset_iam_policy":                             ResourceIamPolicy(DataplexAssetIamSchema, DataplexAssetIamUpdaterProducer, DataplexAssetIdParseFunc),
			"google_dataplex_lake_iam_binding":                             ResourceIamBinding(DataplexLakeIamSchema, DataplexLakeIamUpdaterProducer, DataplexLakeIdParseFunc),
			"google_dataplex_lake_iam_member":                              ResourceIamMember(DataplexLakeIamSchema, DataplexLakeIamUpdaterProducer, DataplexLakeIdParseFunc),
			"google_dataplex_lake_iam_policy":                              ResourceIamPolicy(DataplexLakeIamSchema, DataplexLakeIamUpdaterProducer, DataplexLakeIdParseFunc),
			"google_dataplex_zone_iam_binding":                             ResourceIamBinding(DataplexZoneIamSchema, DataplexZoneIamUpdaterProducer, DataplexZoneIdParseFunc),
			"google_dataplex_zone_iam_member":                              ResourceIamMember(DataplexZoneIamSchema, DataplexZoneIamUpdaterProducer, DataplexZoneIdParseFunc),
			"google_dataplex_zone_iam_policy":                              ResourceIamPolicy(DataplexZoneIamSchema, DataplexZoneIamUpdaterProducer, DataplexZoneIdParseFunc),
			"google_dataproc_autoscaling_policy":                           ResourceDataprocAutoscalingPolicy(),
			"google_dataproc_autoscaling_policy_iam_binding":               ResourceIamBinding(DataprocAutoscalingPolicyIamSchema, DataprocAutoscalingPolicyIamUpdaterProducer, DataprocAutoscalingPolicyIdParseFunc),
			"google_dataproc_autoscaling_policy_iam_member":                ResourceIamMember(DataprocAutoscalingPolicyIamSchema, DataprocAutoscalingPolicyIamUpdaterProducer, DataprocAutoscalingPolicyIdParseFunc),
			"google_dataproc_autoscaling_policy_iam_policy":                ResourceIamPolicy(DataprocAutoscalingPolicyIamSchema, DataprocAutoscalingPolicyIamUpdaterProducer, DataprocAutoscalingPolicyIdParseFunc),
			"google_dataproc_metastore_service":                            ResourceDataprocMetastoreService(),
			"google_dataproc_metastore_service_iam_binding":                ResourceIamBinding(DataprocMetastoreServiceIamSchema, DataprocMetastoreServiceIamUpdaterProducer, DataprocMetastoreServiceIdParseFunc),
			"google_dataproc_metastore_service_iam_member":                 ResourceIamMember(DataprocMetastoreServiceIamSchema, DataprocMetastoreServiceIamUpdaterProducer, DataprocMetastoreServiceIdParseFunc),
			"google_dataproc_metastore_service_iam_policy":                 ResourceIamPolicy(DataprocMetastoreServiceIamSchema, DataprocMetastoreServiceIamUpdaterProducer, DataprocMetastoreServiceIdParseFunc),
			"google_datastore_index":                                       ResourceDatastoreIndex(),
			"google_datastream_connection_profile":                         ResourceDatastreamConnectionProfile(),
			"google_datastream_private_connection":                         ResourceDatastreamPrivateConnection(),
			"google_datastream_stream":                                     ResourceDatastreamStream(),
			"google_deployment_manager_deployment":                         ResourceDeploymentManagerDeployment(),
			"google_dialogflow_agent":                                      ResourceDialogflowAgent(),
			"google_dialogflow_entity_type":                                ResourceDialogflowEntityType(),
			"google_dialogflow_fulfillment":                                ResourceDialogflowFulfillment(),
			"google_dialogflow_intent":                                     ResourceDialogflowIntent(),
			"google_dialogflow_cx_agent":                                   ResourceDialogflowCXAgent(),
			"google_dialogflow_cx_entity_type":                             ResourceDialogflowCXEntityType(),
			"google_dialogflow_cx_flow":                                    ResourceDialogflowCXFlow(),
			"google_dialogflow_cx_intent":                                  ResourceDialogflowCXIntent(),
			"google_dialogflow_cx_page":                                    ResourceDialogflowCXPage(),
			"google_dialogflow_cx_webhook":                                 ResourceDialogflowCXWebhook(),
			"google_dns_managed_zone":                                      ResourceDNSManagedZone(),
			"google_dns_managed_zone_iam_binding":                          ResourceIamBinding(DNSManagedZoneIamSchema, DNSManagedZoneIamUpdaterProducer, DNSManagedZoneIdParseFunc),
			"google_dns_managed_zone_iam_member":                           ResourceIamMember(DNSManagedZoneIamSchema, DNSManagedZoneIamUpdaterProducer, DNSManagedZoneIdParseFunc),
			"google_dns_managed_zone_iam_policy":                           ResourceIamPolicy(DNSManagedZoneIamSchema, DNSManagedZoneIamUpdaterProducer, DNSManagedZoneIdParseFunc),
			"google_dns_policy":                                            ResourceDNSPolicy(),
			"google_document_ai_processor":                                 ResourceDocumentAIProcessor(),
			"google_document_ai_processor_default_version":                 ResourceDocumentAIProcessorDefaultVersion(),
			"google_essential_contacts_contact":                            ResourceEssentialContactsContact(),
			"google_filestore_backup":                                      ResourceFilestoreBackup(),
			"google_filestore_instance":                                    ResourceFilestoreInstance(),
			"google_filestore_snapshot":                                    ResourceFilestoreSnapshot(),
			"google_firestore_database":                                    ResourceFirestoreDatabase(),
			"google_firestore_document":                                    ResourceFirestoreDocument(),
			"google_firestore_index":                                       ResourceFirestoreIndex(),
			"google_game_services_game_server_cluster":                     ResourceGameServicesGameServerCluster(),
			"google_game_services_game_server_config":                      ResourceGameServicesGameServerConfig(),
			"google_game_services_game_server_deployment":                  ResourceGameServicesGameServerDeployment(),
			"google_game_services_game_server_deployment_rollout":          ResourceGameServicesGameServerDeploymentRollout(),
			"google_game_services_realm":                                   ResourceGameServicesRealm(),
			"google_gke_backup_backup_plan":                                ResourceGKEBackupBackupPlan(),
			"google_gke_backup_backup_plan_iam_binding":                    ResourceIamBinding(GKEBackupBackupPlanIamSchema, GKEBackupBackupPlanIamUpdaterProducer, GKEBackupBackupPlanIdParseFunc),
			"google_gke_backup_backup_plan_iam_member":                     ResourceIamMember(GKEBackupBackupPlanIamSchema, GKEBackupBackupPlanIamUpdaterProducer, GKEBackupBackupPlanIdParseFunc),
			"google_gke_backup_backup_plan_iam_policy":                     ResourceIamPolicy(GKEBackupBackupPlanIamSchema, GKEBackupBackupPlanIamUpdaterProducer, GKEBackupBackupPlanIdParseFunc),
			"google_gke_hub_membership":                                    ResourceGKEHubMembership(),
			"google_gke_hub_membership_iam_binding":                        ResourceIamBinding(GKEHubMembershipIamSchema, GKEHubMembershipIamUpdaterProducer, GKEHubMembershipIdParseFunc),
			"google_gke_hub_membership_iam_member":                         ResourceIamMember(GKEHubMembershipIamSchema, GKEHubMembershipIamUpdaterProducer, GKEHubMembershipIdParseFunc),
			"google_gke_hub_membership_iam_policy":                         ResourceIamPolicy(GKEHubMembershipIamSchema, GKEHubMembershipIamUpdaterProducer, GKEHubMembershipIdParseFunc),
			"google_healthcare_consent_store":                              ResourceHealthcareConsentStore(),
			"google_healthcare_consent_store_iam_binding":                  ResourceIamBinding(HealthcareConsentStoreIamSchema, HealthcareConsentStoreIamUpdaterProducer, HealthcareConsentStoreIdParseFunc),
			"google_healthcare_consent_store_iam_member":                   ResourceIamMember(HealthcareConsentStoreIamSchema, HealthcareConsentStoreIamUpdaterProducer, HealthcareConsentStoreIdParseFunc),
			"google_healthcare_consent_store_iam_policy":                   ResourceIamPolicy(HealthcareConsentStoreIamSchema, HealthcareConsentStoreIamUpdaterProducer, HealthcareConsentStoreIdParseFunc),
			"google_healthcare_dataset":                                    ResourceHealthcareDataset(),
			"google_healthcare_dicom_store":                                ResourceHealthcareDicomStore(),
			"google_healthcare_fhir_store":                                 ResourceHealthcareFhirStore(),
			"google_healthcare_hl7_v2_store":                               ResourceHealthcareHl7V2Store(),
			"google_iam_access_boundary_policy":                            ResourceIAM2AccessBoundaryPolicy(),
			"google_iam_workload_identity_pool":                            ResourceIAMBetaWorkloadIdentityPool(),
			"google_iam_workload_identity_pool_provider":                   ResourceIAMBetaWorkloadIdentityPoolProvider(),
			"google_iam_workforce_pool":                                    ResourceIAMWorkforcePoolWorkforcePool(),
			"google_iam_workforce_pool_provider":                           ResourceIAMWorkforcePoolWorkforcePoolProvider(),
			"google_iap_app_engine_service_iam_binding":                    ResourceIamBinding(IapAppEngineServiceIamSchema, IapAppEngineServiceIamUpdaterProducer, IapAppEngineServiceIdParseFunc),
			"google_iap_app_engine_service_iam_member":                     ResourceIamMember(IapAppEngineServiceIamSchema, IapAppEngineServiceIamUpdaterProducer, IapAppEngineServiceIdParseFunc),
			"google_iap_app_engine_service_iam_policy":                     ResourceIamPolicy(IapAppEngineServiceIamSchema, IapAppEngineServiceIamUpdaterProducer, IapAppEngineServiceIdParseFunc),
			"google_iap_app_engine_version_iam_binding":                    ResourceIamBinding(IapAppEngineVersionIamSchema, IapAppEngineVersionIamUpdaterProducer, IapAppEngineVersionIdParseFunc),
			"google_iap_app_engine_version_iam_member":                     ResourceIamMember(IapAppEngineVersionIamSchema, IapAppEngineVersionIamUpdaterProducer, IapAppEngineVersionIdParseFunc),
			"google_iap_app_engine_version_iam_policy":                     ResourceIamPolicy(IapAppEngineVersionIamSchema, IapAppEngineVersionIamUpdaterProducer, IapAppEngineVersionIdParseFunc),
			"google_iap_brand":                                             ResourceIapBrand(),
			"google_iap_client":                                            ResourceIapClient(),
			"google_iap_tunnel_iam_binding":                                ResourceIamBinding(IapTunnelIamSchema, IapTunnelIamUpdaterProducer, IapTunnelIdParseFunc),
			"google_iap_tunnel_iam_member":                                 ResourceIamMember(IapTunnelIamSchema, IapTunnelIamUpdaterProducer, IapTunnelIdParseFunc),
			"google_iap_tunnel_iam_policy":                                 ResourceIamPolicy(IapTunnelIamSchema, IapTunnelIamUpdaterProducer, IapTunnelIdParseFunc),
			"google_iap_tunnel_instance_iam_binding":                       ResourceIamBinding(IapTunnelInstanceIamSchema, IapTunnelInstanceIamUpdaterProducer, IapTunnelInstanceIdParseFunc),
			"google_iap_tunnel_instance_iam_member":                        ResourceIamMember(IapTunnelInstanceIamSchema, IapTunnelInstanceIamUpdaterProducer, IapTunnelInstanceIdParseFunc),
			"google_iap_tunnel_instance_iam_policy":                        ResourceIamPolicy(IapTunnelInstanceIamSchema, IapTunnelInstanceIamUpdaterProducer, IapTunnelInstanceIdParseFunc),
			"google_iap_web_iam_binding":                                   ResourceIamBinding(IapWebIamSchema, IapWebIamUpdaterProducer, IapWebIdParseFunc),
			"google_iap_web_iam_member":                                    ResourceIamMember(IapWebIamSchema, IapWebIamUpdaterProducer, IapWebIdParseFunc),
			"google_iap_web_iam_policy":                                    ResourceIamPolicy(IapWebIamSchema, IapWebIamUpdaterProducer, IapWebIdParseFunc),
			"google_iap_web_backend_service_iam_binding":                   ResourceIamBinding(IapWebBackendServiceIamSchema, IapWebBackendServiceIamUpdaterProducer, IapWebBackendServiceIdParseFunc),
			"google_iap_web_backend_service_iam_member":                    ResourceIamMember(IapWebBackendServiceIamSchema, IapWebBackendServiceIamUpdaterProducer, IapWebBackendServiceIdParseFunc),
			"google_iap_web_backend_service_iam_policy":                    ResourceIamPolicy(IapWebBackendServiceIamSchema, IapWebBackendServiceIamUpdaterProducer, IapWebBackendServiceIdParseFunc),
			"google_iap_web_type_app_engine_iam_binding":                   ResourceIamBinding(IapWebTypeAppEngineIamSchema, IapWebTypeAppEngineIamUpdaterProducer, IapWebTypeAppEngineIdParseFunc),
			"google_iap_web_type_app_engine_iam_member":                    ResourceIamMember(IapWebTypeAppEngineIamSchema, IapWebTypeAppEngineIamUpdaterProducer, IapWebTypeAppEngineIdParseFunc),
			"google_iap_web_type_app_engine_iam_policy":                    ResourceIamPolicy(IapWebTypeAppEngineIamSchema, IapWebTypeAppEngineIamUpdaterProducer, IapWebTypeAppEngineIdParseFunc),
			"google_iap_web_type_compute_iam_binding":                      ResourceIamBinding(IapWebTypeComputeIamSchema, IapWebTypeComputeIamUpdaterProducer, IapWebTypeComputeIdParseFunc),
			"google_iap_web_type_compute_iam_member":                       ResourceIamMember(IapWebTypeComputeIamSchema, IapWebTypeComputeIamUpdaterProducer, IapWebTypeComputeIdParseFunc),
			"google_iap_web_type_compute_iam_policy":                       ResourceIamPolicy(IapWebTypeComputeIamSchema, IapWebTypeComputeIamUpdaterProducer, IapWebTypeComputeIdParseFunc),
			"google_identity_platform_config":                              ResourceIdentityPlatformConfig(),
			"google_identity_platform_default_supported_idp_config":        ResourceIdentityPlatformDefaultSupportedIdpConfig(),
			"google_identity_platform_inbound_saml_config":                 ResourceIdentityPlatformInboundSamlConfig(),
			"google_identity_platform_oauth_idp_config":                    ResourceIdentityPlatformOauthIdpConfig(),
			"google_identity_platform_project_default_config":              ResourceIdentityPlatformProjectDefaultConfig(),
			"google_identity_platform_tenant":                              ResourceIdentityPlatformTenant(),
			"google_identity_platform_tenant_default_supported_idp_config": ResourceIdentityPlatformTenantDefaultSupportedIdpConfig(),
			"google_identity_platform_tenant_inbound_saml_config":          ResourceIdentityPlatformTenantInboundSamlConfig(),
			"google_identity_platform_tenant_oauth_idp_config":             ResourceIdentityPlatformTenantOauthIdpConfig(),
			"google_kms_crypto_key":                                        ResourceKMSCryptoKey(),
			"google_kms_crypto_key_version":                                ResourceKMSCryptoKeyVersion(),
			"google_kms_key_ring":                                          ResourceKMSKeyRing(),
			"google_kms_key_ring_import_job":                               ResourceKMSKeyRingImportJob(),
			"google_kms_secret_ciphertext":                                 ResourceKMSSecretCiphertext(),
			"google_logging_log_view":                                      ResourceLoggingLogView(),
			"google_logging_metric":                                        ResourceLoggingMetric(),
			"google_memcache_instance":                                     ResourceMemcacheInstance(),
			"google_ml_engine_model":                                       ResourceMLEngineModel(),
			"google_monitoring_alert_policy":                               ResourceMonitoringAlertPolicy(),
			"google_monitoring_service":                                    ResourceMonitoringGenericService(),
			"google_monitoring_group":                                      ResourceMonitoringGroup(),
			"google_monitoring_metric_descriptor":                          ResourceMonitoringMetricDescriptor(),
			"google_monitoring_notification_channel":                       ResourceMonitoringNotificationChannel(),
			"google_monitoring_custom_service":                             ResourceMonitoringService(),
			"google_monitoring_slo":                                        ResourceMonitoringSlo(),
			"google_monitoring_uptime_check_config":                        ResourceMonitoringUptimeCheckConfig(),
			"google_network_management_connectivity_test":                  ResourceNetworkManagementConnectivityTest(),
			"google_network_services_edge_cache_keyset":                    ResourceNetworkServicesEdgeCacheKeyset(),
			"google_network_services_edge_cache_origin":                    ResourceNetworkServicesEdgeCacheOrigin(),
			"google_network_services_edge_cache_service":                   ResourceNetworkServicesEdgeCacheService(),
			"google_notebooks_environment":                                 ResourceNotebooksEnvironment(),
			"google_notebooks_instance":                                    ResourceNotebooksInstance(),
			"google_notebooks_instance_iam_binding":                        ResourceIamBinding(NotebooksInstanceIamSchema, NotebooksInstanceIamUpdaterProducer, NotebooksInstanceIdParseFunc),
			"google_notebooks_instance_iam_member":                         ResourceIamMember(NotebooksInstanceIamSchema, NotebooksInstanceIamUpdaterProducer, NotebooksInstanceIdParseFunc),
			"google_notebooks_instance_iam_policy":                         ResourceIamPolicy(NotebooksInstanceIamSchema, NotebooksInstanceIamUpdaterProducer, NotebooksInstanceIdParseFunc),
			"google_notebooks_location":                                    ResourceNotebooksLocation(),
			"google_notebooks_runtime":                                     ResourceNotebooksRuntime(),
			"google_notebooks_runtime_iam_binding":                         ResourceIamBinding(NotebooksRuntimeIamSchema, NotebooksRuntimeIamUpdaterProducer, NotebooksRuntimeIdParseFunc),
			"google_notebooks_runtime_iam_member":                          ResourceIamMember(NotebooksRuntimeIamSchema, NotebooksRuntimeIamUpdaterProducer, NotebooksRuntimeIdParseFunc),
			"google_notebooks_runtime_iam_policy":                          ResourceIamPolicy(NotebooksRuntimeIamSchema, NotebooksRuntimeIamUpdaterProducer, NotebooksRuntimeIdParseFunc),
			"google_os_config_patch_deployment":                            ResourceOSConfigPatchDeployment(),
			"google_os_login_ssh_public_key":                               ResourceOSLoginSSHPublicKey(),
			"google_privateca_ca_pool":                                     ResourcePrivatecaCaPool(),
			"google_privateca_ca_pool_iam_binding":                         ResourceIamBinding(PrivatecaCaPoolIamSchema, PrivatecaCaPoolIamUpdaterProducer, PrivatecaCaPoolIdParseFunc),
			"google_privateca_ca_pool_iam_member":                          ResourceIamMember(PrivatecaCaPoolIamSchema, PrivatecaCaPoolIamUpdaterProducer, PrivatecaCaPoolIdParseFunc),
			"google_privateca_ca_pool_iam_policy":                          ResourceIamPolicy(PrivatecaCaPoolIamSchema, PrivatecaCaPoolIamUpdaterProducer, PrivatecaCaPoolIdParseFunc),
			"google_privateca_certificate":                                 ResourcePrivatecaCertificate(),
			"google_privateca_certificate_authority":                       ResourcePrivatecaCertificateAuthority(),
			"google_privateca_certificate_template_iam_binding":            ResourceIamBinding(PrivatecaCertificateTemplateIamSchema, PrivatecaCertificateTemplateIamUpdaterProducer, PrivatecaCertificateTemplateIdParseFunc),
			"google_privateca_certificate_template_iam_member":             ResourceIamMember(PrivatecaCertificateTemplateIamSchema, PrivatecaCertificateTemplateIamUpdaterProducer, PrivatecaCertificateTemplateIdParseFunc),
			"google_privateca_certificate_template_iam_policy":             ResourceIamPolicy(PrivatecaCertificateTemplateIamSchema, PrivatecaCertificateTemplateIamUpdaterProducer, PrivatecaCertificateTemplateIdParseFunc),
			"google_pubsub_schema":                                         ResourcePubsubSchema(),
			"google_pubsub_subscription":                                   ResourcePubsubSubscription(),
			"google_pubsub_topic":                                          ResourcePubsubTopic(),
			"google_pubsub_topic_iam_binding":                              ResourceIamBinding(PubsubTopicIamSchema, PubsubTopicIamUpdaterProducer, PubsubTopicIdParseFunc),
			"google_pubsub_topic_iam_member":                               ResourceIamMember(PubsubTopicIamSchema, PubsubTopicIamUpdaterProducer, PubsubTopicIdParseFunc),
			"google_pubsub_topic_iam_policy":                               ResourceIamPolicy(PubsubTopicIamSchema, PubsubTopicIamUpdaterProducer, PubsubTopicIdParseFunc),
			"google_pubsub_lite_reservation":                               ResourcePubsubLiteReservation(),
			"google_pubsub_lite_subscription":                              ResourcePubsubLiteSubscription(),
			"google_pubsub_lite_topic":                                     ResourcePubsubLiteTopic(),
			"google_redis_instance":                                        ResourceRedisInstance(),
			"google_resource_manager_lien":                                 ResourceResourceManagerLien(),
			"google_secret_manager_secret":                                 ResourceSecretManagerSecret(),
			"google_secret_manager_secret_iam_binding":                     ResourceIamBinding(SecretManagerSecretIamSchema, SecretManagerSecretIamUpdaterProducer, SecretManagerSecretIdParseFunc),
			"google_secret_manager_secret_iam_member":                      ResourceIamMember(SecretManagerSecretIamSchema, SecretManagerSecretIamUpdaterProducer, SecretManagerSecretIdParseFunc),
			"google_secret_manager_secret_iam_policy":                      ResourceIamPolicy(SecretManagerSecretIamSchema, SecretManagerSecretIamUpdaterProducer, SecretManagerSecretIdParseFunc),
			"google_secret_manager_secret_version":                         ResourceSecretManagerSecretVersion(),
			"google_scc_mute_config":                                       ResourceSecurityCenterMuteConfig(),
			"google_scc_notification_config":                               ResourceSecurityCenterNotificationConfig(),
			"google_scc_source":                                            ResourceSecurityCenterSource(),
			"google_scc_source_iam_binding":                                ResourceIamBinding(SecurityCenterSourceIamSchema, SecurityCenterSourceIamUpdaterProducer, SecurityCenterSourceIdParseFunc),
			"google_scc_source_iam_member":                                 ResourceIamMember(SecurityCenterSourceIamSchema, SecurityCenterSourceIamUpdaterProducer, SecurityCenterSourceIdParseFunc),
			"google_scc_source_iam_policy":                                 ResourceIamPolicy(SecurityCenterSourceIamSchema, SecurityCenterSourceIamUpdaterProducer, SecurityCenterSourceIdParseFunc),
			"google_endpoints_service_iam_binding":                         ResourceIamBinding(ServiceManagementServiceIamSchema, ServiceManagementServiceIamUpdaterProducer, ServiceManagementServiceIdParseFunc),
			"google_endpoints_service_iam_member":                          ResourceIamMember(ServiceManagementServiceIamSchema, ServiceManagementServiceIamUpdaterProducer, ServiceManagementServiceIdParseFunc),
			"google_endpoints_service_iam_policy":                          ResourceIamPolicy(ServiceManagementServiceIamSchema, ServiceManagementServiceIamUpdaterProducer, ServiceManagementServiceIdParseFunc),
			"google_endpoints_service_consumers_iam_binding":               ResourceIamBinding(ServiceManagementServiceConsumersIamSchema, ServiceManagementServiceConsumersIamUpdaterProducer, ServiceManagementServiceConsumersIdParseFunc),
			"google_endpoints_service_consumers_iam_member":                ResourceIamMember(ServiceManagementServiceConsumersIamSchema, ServiceManagementServiceConsumersIamUpdaterProducer, ServiceManagementServiceConsumersIdParseFunc),
			"google_endpoints_service_consumers_iam_policy":                ResourceIamPolicy(ServiceManagementServiceConsumersIamSchema, ServiceManagementServiceConsumersIamUpdaterProducer, ServiceManagementServiceConsumersIdParseFunc),
			"google_sourcerepo_repository":                                 ResourceSourceRepoRepository(),
			"google_sourcerepo_repository_iam_binding":                     ResourceIamBinding(SourceRepoRepositoryIamSchema, SourceRepoRepositoryIamUpdaterProducer, SourceRepoRepositoryIdParseFunc),
			"google_sourcerepo_repository_iam_member":                      ResourceIamMember(SourceRepoRepositoryIamSchema, SourceRepoRepositoryIamUpdaterProducer, SourceRepoRepositoryIdParseFunc),
			"google_sourcerepo_repository_iam_policy":                      ResourceIamPolicy(SourceRepoRepositoryIamSchema, SourceRepoRepositoryIamUpdaterProducer, SourceRepoRepositoryIdParseFunc),
			"google_spanner_database":                                      ResourceSpannerDatabase(),
			"google_spanner_instance":                                      ResourceSpannerInstance(),
			"google_sql_database":                                          ResourceSQLDatabase(),
			"google_sql_source_representation_instance":                    ResourceSQLSourceRepresentationInstance(),
			"google_storage_bucket_iam_binding":                            ResourceIamBinding(StorageBucketIamSchema, StorageBucketIamUpdaterProducer, StorageBucketIdParseFunc),
			"google_storage_bucket_iam_member":                             ResourceIamMember(StorageBucketIamSchema, StorageBucketIamUpdaterProducer, StorageBucketIdParseFunc),
			"google_storage_bucket_iam_policy":                             ResourceIamPolicy(StorageBucketIamSchema, StorageBucketIamUpdaterProducer, StorageBucketIdParseFunc),
			"google_storage_bucket_access_control":                         ResourceStorageBucketAccessControl(),
			"google_storage_default_object_access_control":                 ResourceStorageDefaultObjectAccessControl(),
			"google_storage_hmac_key":                                      ResourceStorageHmacKey(),
			"google_storage_object_access_control":                         ResourceStorageObjectAccessControl(),
			"google_storage_transfer_agent_pool":                           ResourceStorageTransferAgentPool(),
			"google_tags_tag_binding":                                      ResourceTagsTagBinding(),
			"google_tags_tag_key":                                          ResourceTagsTagKey(),
			"google_tags_tag_key_iam_binding":                              ResourceIamBinding(TagsTagKeyIamSchema, TagsTagKeyIamUpdaterProducer, TagsTagKeyIdParseFunc),
			"google_tags_tag_key_iam_member":                               ResourceIamMember(TagsTagKeyIamSchema, TagsTagKeyIamUpdaterProducer, TagsTagKeyIdParseFunc),
			"google_tags_tag_key_iam_policy":                               ResourceIamPolicy(TagsTagKeyIamSchema, TagsTagKeyIamUpdaterProducer, TagsTagKeyIdParseFunc),
			"google_tags_tag_value":                                        ResourceTagsTagValue(),
			"google_tags_tag_value_iam_binding":                            ResourceIamBinding(TagsTagValueIamSchema, TagsTagValueIamUpdaterProducer, TagsTagValueIdParseFunc),
			"google_tags_tag_value_iam_member":                             ResourceIamMember(TagsTagValueIamSchema, TagsTagValueIamUpdaterProducer, TagsTagValueIdParseFunc),
			"google_tags_tag_value_iam_policy":                             ResourceIamPolicy(TagsTagValueIamSchema, TagsTagValueIamUpdaterProducer, TagsTagValueIdParseFunc),
			"google_tpu_node":                                              ResourceTPUNode(),
			"google_vertex_ai_dataset":                                     ResourceVertexAIDataset(),
			"google_vertex_ai_endpoint":                                    ResourceVertexAIEndpoint(),
			"google_vertex_ai_featurestore":                                ResourceVertexAIFeaturestore(),
			"google_vertex_ai_featurestore_entitytype":                     ResourceVertexAIFeaturestoreEntitytype(),
			"google_vertex_ai_featurestore_entitytype_feature":             ResourceVertexAIFeaturestoreEntitytypeFeature(),
			"google_vertex_ai_index":                                       ResourceVertexAIIndex(),
			"google_vertex_ai_tensorboard":                                 ResourceVertexAITensorboard(),
			"google_vpc_access_connector":                                  ResourceVPCAccessConnector(),
			"google_workflows_workflow":                                    ResourceWorkflowsWorkflow(),
		},
		map[string]*schema.Resource{
			// ####### START handwritten resources ###########
			"google_app_engine_application":                 ResourceAppEngineApplication(),
			"google_apigee_sharedflow":                      ResourceApigeeSharedFlow(),
			"google_apigee_sharedflow_deployment":           ResourceApigeeSharedFlowDeployment(),
			"google_apigee_flowhook":                        ResourceApigeeFlowhook(),
			"google_apigee_keystores_aliases_pkcs12":        ResourceApigeeKeystoresAliasesPkcs12(),
			"google_apigee_keystores_aliases_key_cert_file": ResourceApigeeKeystoresAliasesKeyCertFile(),
			"google_bigquery_table":                         ResourceBigQueryTable(),
			"google_bigtable_gc_policy":                     ResourceBigtableGCPolicy(),
			"google_bigtable_instance":                      ResourceBigtableInstance(),
			"google_bigtable_table":                         ResourceBigtableTable(),
			"google_billing_subaccount":                     ResourceBillingSubaccount(),
			"google_cloudfunctions_function":                ResourceCloudFunctionsFunction(),
			"google_composer_environment":                   ResourceComposerEnvironment(),
			"google_compute_attached_disk":                  ResourceComputeAttachedDisk(),
			"google_compute_instance":                       ResourceComputeInstance(),
			"google_compute_instance_from_template":         ResourceComputeInstanceFromTemplate(),
			"google_compute_instance_group":                 ResourceComputeInstanceGroup(),
			"google_compute_instance_group_manager":         ResourceComputeInstanceGroupManager(),
			"google_compute_instance_template":              ResourceComputeInstanceTemplate(),
			"google_compute_network_peering":                ResourceComputeNetworkPeering(),
			"google_compute_project_default_network_tier":   ResourceComputeProjectDefaultNetworkTier(),
			"google_compute_project_metadata":               ResourceComputeProjectMetadata(),
			"google_compute_project_metadata_item":          ResourceComputeProjectMetadataItem(),
			"google_compute_region_instance_group_manager":  ResourceComputeRegionInstanceGroupManager(),
			"google_compute_router_interface":               ResourceComputeRouterInterface(),
			"google_compute_security_policy":                ResourceComputeSecurityPolicy(),
			"google_compute_shared_vpc_host_project":        ResourceComputeSharedVpcHostProject(),
			"google_compute_shared_vpc_service_project":     ResourceComputeSharedVpcServiceProject(),
			"google_compute_target_pool":                    ResourceComputeTargetPool(),
			"google_container_cluster":                      ResourceContainerCluster(),
			"google_container_node_pool":                    ResourceContainerNodePool(),
			"google_container_registry":                     ResourceContainerRegistry(),
			"google_dataflow_job":                           ResourceDataflowJob(),
			"google_dataproc_cluster":                       ResourceDataprocCluster(),
			"google_dataproc_job":                           ResourceDataprocJob(),
			"google_dialogflow_cx_version":                  ResourceDialogflowCXVersion(),
			"google_dialogflow_cx_environment":              ResourceDialogflowCXEnvironment(),
			"google_dns_record_set":                         ResourceDnsRecordSet(),
			"google_endpoints_service":                      ResourceEndpointsService(),
			"google_folder":                                 ResourceGoogleFolder(),
			"google_folder_organization_policy":             ResourceGoogleFolderOrganizationPolicy(),
			"google_logging_billing_account_sink":           ResourceLoggingBillingAccountSink(),
			"google_logging_billing_account_exclusion":      ResourceLoggingExclusion(BillingAccountLoggingExclusionSchema, NewBillingAccountLoggingExclusionUpdater, BillingAccountLoggingExclusionIdParseFunc),
			"google_logging_billing_account_bucket_config":  ResourceLoggingBillingAccountBucketConfig(),
			"google_logging_organization_sink":              ResourceLoggingOrganizationSink(),
			"google_logging_organization_exclusion":         ResourceLoggingExclusion(OrganizationLoggingExclusionSchema, NewOrganizationLoggingExclusionUpdater, OrganizationLoggingExclusionIdParseFunc),
			"google_logging_organization_bucket_config":     ResourceLoggingOrganizationBucketConfig(),
			"google_logging_folder_sink":                    ResourceLoggingFolderSink(),
			"google_logging_folder_exclusion":               ResourceLoggingExclusion(FolderLoggingExclusionSchema, NewFolderLoggingExclusionUpdater, FolderLoggingExclusionIdParseFunc),
			"google_logging_folder_bucket_config":           ResourceLoggingFolderBucketConfig(),
			"google_logging_project_sink":                   ResourceLoggingProjectSink(),
			"google_logging_project_exclusion":              ResourceLoggingExclusion(ProjectLoggingExclusionSchema, NewProjectLoggingExclusionUpdater, ProjectLoggingExclusionIdParseFunc),
			"google_logging_project_bucket_config":          ResourceLoggingProjectBucketConfig(),
			"google_monitoring_dashboard":                   ResourceMonitoringDashboard(),
			"google_service_networking_connection":          ResourceServiceNetworkingConnection(),
			"google_sql_database_instance":                  ResourceSqlDatabaseInstance(),
			"google_sql_ssl_cert":                           ResourceSqlSslCert(),
			"google_sql_user":                               ResourceSqlUser(),
			"google_organization_iam_custom_role":           ResourceGoogleOrganizationIamCustomRole(),
			"google_organization_policy":                    ResourceGoogleOrganizationPolicy(),
			"google_project":                                ResourceGoogleProject(),
			"google_project_default_service_accounts":       ResourceGoogleProjectDefaultServiceAccounts(),
			"google_project_service":                        ResourceGoogleProjectService(),
			"google_project_iam_custom_role":                ResourceGoogleProjectIamCustomRole(),
			"google_project_organization_policy":            ResourceGoogleProjectOrganizationPolicy(),
			"google_project_usage_export_bucket":            ResourceProjectUsageBucket(),
			"google_service_account":                        ResourceGoogleServiceAccount(),
			"google_service_account_key":                    ResourceGoogleServiceAccountKey(),
			"google_service_networking_peered_dns_domain":   ResourceGoogleServiceNetworkingPeeredDNSDomain(),
			"google_storage_bucket":                         ResourceStorageBucket(),
			"google_storage_bucket_acl":                     ResourceStorageBucketAcl(),
			"google_storage_bucket_object":                  ResourceStorageBucketObject(),
			"google_storage_object_acl":                     ResourceStorageObjectAcl(),
			"google_storage_default_object_acl":             ResourceStorageDefaultObjectAcl(),
			"google_storage_notification":                   ResourceStorageNotification(),
			"google_storage_transfer_job":                   ResourceStorageTransferJob(),
			"google_tags_location_tag_binding":              ResourceTagsLocationTagBinding(),
			// ####### END handwritten resources ###########
		},
		map[string]*schema.Resource{
			// ####### START non-generated IAM resources ###########
			"google_bigtable_instance_iam_binding":       ResourceIamBinding(IamBigtableInstanceSchema, NewBigtableInstanceUpdater, BigtableInstanceIdParseFunc),
			"google_bigtable_instance_iam_member":        ResourceIamMember(IamBigtableInstanceSchema, NewBigtableInstanceUpdater, BigtableInstanceIdParseFunc),
			"google_bigtable_instance_iam_policy":        ResourceIamPolicy(IamBigtableInstanceSchema, NewBigtableInstanceUpdater, BigtableInstanceIdParseFunc),
			"google_bigtable_table_iam_binding":          ResourceIamBinding(IamBigtableTableSchema, NewBigtableTableUpdater, BigtableTableIdParseFunc),
			"google_bigtable_table_iam_member":           ResourceIamMember(IamBigtableTableSchema, NewBigtableTableUpdater, BigtableTableIdParseFunc),
			"google_bigtable_table_iam_policy":           ResourceIamPolicy(IamBigtableTableSchema, NewBigtableTableUpdater, BigtableTableIdParseFunc),
			"google_bigquery_dataset_iam_binding":        ResourceIamBinding(IamBigqueryDatasetSchema, NewBigqueryDatasetIamUpdater, BigqueryDatasetIdParseFunc),
			"google_bigquery_dataset_iam_member":         ResourceIamMember(IamBigqueryDatasetSchema, NewBigqueryDatasetIamUpdater, BigqueryDatasetIdParseFunc),
			"google_bigquery_dataset_iam_policy":         ResourceIamPolicy(IamBigqueryDatasetSchema, NewBigqueryDatasetIamUpdater, BigqueryDatasetIdParseFunc),
			"google_billing_account_iam_binding":         ResourceIamBinding(IamBillingAccountSchema, NewBillingAccountIamUpdater, BillingAccountIdParseFunc),
			"google_billing_account_iam_member":          ResourceIamMember(IamBillingAccountSchema, NewBillingAccountIamUpdater, BillingAccountIdParseFunc),
			"google_billing_account_iam_policy":          ResourceIamPolicy(IamBillingAccountSchema, NewBillingAccountIamUpdater, BillingAccountIdParseFunc),
			"google_dataproc_cluster_iam_binding":        ResourceIamBinding(IamDataprocClusterSchema, NewDataprocClusterUpdater, DataprocClusterIdParseFunc),
			"google_dataproc_cluster_iam_member":         ResourceIamMember(IamDataprocClusterSchema, NewDataprocClusterUpdater, DataprocClusterIdParseFunc),
			"google_dataproc_cluster_iam_policy":         ResourceIamPolicy(IamDataprocClusterSchema, NewDataprocClusterUpdater, DataprocClusterIdParseFunc),
			"google_dataproc_job_iam_binding":            ResourceIamBinding(IamDataprocJobSchema, NewDataprocJobUpdater, DataprocJobIdParseFunc),
			"google_dataproc_job_iam_member":             ResourceIamMember(IamDataprocJobSchema, NewDataprocJobUpdater, DataprocJobIdParseFunc),
			"google_dataproc_job_iam_policy":             ResourceIamPolicy(IamDataprocJobSchema, NewDataprocJobUpdater, DataprocJobIdParseFunc),
			"google_folder_iam_binding":                  ResourceIamBinding(IamFolderSchema, NewFolderIamUpdater, FolderIdParseFunc),
			"google_folder_iam_member":                   ResourceIamMember(IamFolderSchema, NewFolderIamUpdater, FolderIdParseFunc),
			"google_folder_iam_policy":                   ResourceIamPolicy(IamFolderSchema, NewFolderIamUpdater, FolderIdParseFunc),
			"google_folder_iam_audit_config":             ResourceIamAuditConfig(IamFolderSchema, NewFolderIamUpdater, FolderIdParseFunc),
			"google_healthcare_dataset_iam_binding":      ResourceIamBindingWithBatching(IamHealthcareDatasetSchema, NewHealthcareDatasetIamUpdater, DatasetIdParseFunc, IamBatchingEnabled),
			"google_healthcare_dataset_iam_member":       ResourceIamMemberWithBatching(IamHealthcareDatasetSchema, NewHealthcareDatasetIamUpdater, DatasetIdParseFunc, IamBatchingEnabled),
			"google_healthcare_dataset_iam_policy":       ResourceIamPolicy(IamHealthcareDatasetSchema, NewHealthcareDatasetIamUpdater, DatasetIdParseFunc),
			"google_healthcare_dicom_store_iam_binding":  ResourceIamBindingWithBatching(IamHealthcareDicomStoreSchema, NewHealthcareDicomStoreIamUpdater, DicomStoreIdParseFunc, IamBatchingEnabled),
			"google_healthcare_dicom_store_iam_member":   ResourceIamMemberWithBatching(IamHealthcareDicomStoreSchema, NewHealthcareDicomStoreIamUpdater, DicomStoreIdParseFunc, IamBatchingEnabled),
			"google_healthcare_dicom_store_iam_policy":   ResourceIamPolicy(IamHealthcareDicomStoreSchema, NewHealthcareDicomStoreIamUpdater, DicomStoreIdParseFunc),
			"google_healthcare_fhir_store_iam_binding":   ResourceIamBindingWithBatching(IamHealthcareFhirStoreSchema, NewHealthcareFhirStoreIamUpdater, FhirStoreIdParseFunc, IamBatchingEnabled),
			"google_healthcare_fhir_store_iam_member":    ResourceIamMemberWithBatching(IamHealthcareFhirStoreSchema, NewHealthcareFhirStoreIamUpdater, FhirStoreIdParseFunc, IamBatchingEnabled),
			"google_healthcare_fhir_store_iam_policy":    ResourceIamPolicy(IamHealthcareFhirStoreSchema, NewHealthcareFhirStoreIamUpdater, FhirStoreIdParseFunc),
			"google_healthcare_hl7_v2_store_iam_binding": ResourceIamBindingWithBatching(IamHealthcareHl7V2StoreSchema, NewHealthcareHl7V2StoreIamUpdater, Hl7V2StoreIdParseFunc, IamBatchingEnabled),
			"google_healthcare_hl7_v2_store_iam_member":  ResourceIamMemberWithBatching(IamHealthcareHl7V2StoreSchema, NewHealthcareHl7V2StoreIamUpdater, Hl7V2StoreIdParseFunc, IamBatchingEnabled),
			"google_healthcare_hl7_v2_store_iam_policy":  ResourceIamPolicy(IamHealthcareHl7V2StoreSchema, NewHealthcareHl7V2StoreIamUpdater, Hl7V2StoreIdParseFunc),
			"google_kms_key_ring_iam_binding":            ResourceIamBinding(IamKmsKeyRingSchema, NewKmsKeyRingIamUpdater, KeyRingIdParseFunc),
			"google_kms_key_ring_iam_member":             ResourceIamMember(IamKmsKeyRingSchema, NewKmsKeyRingIamUpdater, KeyRingIdParseFunc),
			"google_kms_key_ring_iam_policy":             ResourceIamPolicy(IamKmsKeyRingSchema, NewKmsKeyRingIamUpdater, KeyRingIdParseFunc),
			"google_kms_crypto_key_iam_binding":          ResourceIamBinding(IamKmsCryptoKeySchema, NewKmsCryptoKeyIamUpdater, CryptoIdParseFunc),
			"google_kms_crypto_key_iam_member":           ResourceIamMember(IamKmsCryptoKeySchema, NewKmsCryptoKeyIamUpdater, CryptoIdParseFunc),
			"google_kms_crypto_key_iam_policy":           ResourceIamPolicy(IamKmsCryptoKeySchema, NewKmsCryptoKeyIamUpdater, CryptoIdParseFunc),
			"google_spanner_instance_iam_binding":        ResourceIamBinding(IamSpannerInstanceSchema, NewSpannerInstanceIamUpdater, SpannerInstanceIdParseFunc),
			"google_spanner_instance_iam_member":         ResourceIamMember(IamSpannerInstanceSchema, NewSpannerInstanceIamUpdater, SpannerInstanceIdParseFunc),
			"google_spanner_instance_iam_policy":         ResourceIamPolicy(IamSpannerInstanceSchema, NewSpannerInstanceIamUpdater, SpannerInstanceIdParseFunc),
			"google_spanner_database_iam_binding":        ResourceIamBinding(IamSpannerDatabaseSchema, NewSpannerDatabaseIamUpdater, SpannerDatabaseIdParseFunc),
			"google_spanner_database_iam_member":         ResourceIamMember(IamSpannerDatabaseSchema, NewSpannerDatabaseIamUpdater, SpannerDatabaseIdParseFunc),
			"google_spanner_database_iam_policy":         ResourceIamPolicy(IamSpannerDatabaseSchema, NewSpannerDatabaseIamUpdater, SpannerDatabaseIdParseFunc),
			"google_organization_iam_binding":            ResourceIamBinding(IamOrganizationSchema, NewOrganizationIamUpdater, OrgIdParseFunc),
			"google_organization_iam_member":             ResourceIamMember(IamOrganizationSchema, NewOrganizationIamUpdater, OrgIdParseFunc),
			"google_organization_iam_policy":             ResourceIamPolicy(IamOrganizationSchema, NewOrganizationIamUpdater, OrgIdParseFunc),
			"google_organization_iam_audit_config":       ResourceIamAuditConfig(IamOrganizationSchema, NewOrganizationIamUpdater, OrgIdParseFunc),
			"google_project_iam_policy":                  ResourceIamPolicy(IamProjectSchema, NewProjectIamUpdater, ProjectIdParseFunc),
			"google_project_iam_binding":                 ResourceIamBindingWithBatching(IamProjectSchema, NewProjectIamUpdater, ProjectIdParseFunc, IamBatchingEnabled),
			"google_project_iam_member":                  ResourceIamMemberWithBatching(IamProjectSchema, NewProjectIamUpdater, ProjectIdParseFunc, IamBatchingEnabled),
			"google_project_iam_audit_config":            ResourceIamAuditConfigWithBatching(IamProjectSchema, NewProjectIamUpdater, ProjectIdParseFunc, IamBatchingEnabled),
			"google_pubsub_subscription_iam_binding":     ResourceIamBinding(IamPubsubSubscriptionSchema, NewPubsubSubscriptionIamUpdater, PubsubSubscriptionIdParseFunc),
			"google_pubsub_subscription_iam_member":      ResourceIamMember(IamPubsubSubscriptionSchema, NewPubsubSubscriptionIamUpdater, PubsubSubscriptionIdParseFunc),
			"google_pubsub_subscription_iam_policy":      ResourceIamPolicy(IamPubsubSubscriptionSchema, NewPubsubSubscriptionIamUpdater, PubsubSubscriptionIdParseFunc),
			"google_service_account_iam_binding":         ResourceIamBinding(IamServiceAccountSchema, NewServiceAccountIamUpdater, ServiceAccountIdParseFunc),
			"google_service_account_iam_member":          ResourceIamMember(IamServiceAccountSchema, NewServiceAccountIamUpdater, ServiceAccountIdParseFunc),
			"google_service_account_iam_policy":          ResourceIamPolicy(IamServiceAccountSchema, NewServiceAccountIamUpdater, ServiceAccountIdParseFunc),
			// ####### END non-generated IAM resources ###########
		},
		dclResources,
	)
}

func providerConfigure(ctx context.Context, d *schema.ResourceData, p *schema.Provider) (interface{}, diag.Diagnostics) {
	err := HandleSDKDefaults(d)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	HandleDCLCustomEndpointDefaults(d)

	config := Config{
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
		config.Credentials = MultiEnvSearch([]string{
			"GOOGLE_CREDENTIALS",
			"GOOGLE_CLOUD_KEYFILE_JSON",
			"GCLOUD_KEYFILE_JSON",
		})

		config.AccessToken = MultiEnvSearch([]string{
			"GOOGLE_OAUTH_ACCESS_TOKEN",
		})
	}

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

	batchCfg, err := ExpandProviderBatchingConfig(d.Get("batching"))
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
	config.ArtifactRegistryBasePath = d.Get("artifact_registry_custom_endpoint").(string)
	config.BeyondcorpBasePath = d.Get("beyondcorp_custom_endpoint").(string)
	config.BigQueryBasePath = d.Get("big_query_custom_endpoint").(string)
	config.BigqueryAnalyticsHubBasePath = d.Get("bigquery_analytics_hub_custom_endpoint").(string)
	config.BigqueryConnectionBasePath = d.Get("bigquery_connection_custom_endpoint").(string)
	config.BigqueryDatapolicyBasePath = d.Get("bigquery_datapolicy_custom_endpoint").(string)
	config.BigqueryDataTransferBasePath = d.Get("bigquery_data_transfer_custom_endpoint").(string)
	config.BigqueryReservationBasePath = d.Get("bigquery_reservation_custom_endpoint").(string)
	config.BigtableBasePath = d.Get("bigtable_custom_endpoint").(string)
	config.BillingBasePath = d.Get("billing_custom_endpoint").(string)
	config.BinaryAuthorizationBasePath = d.Get("binary_authorization_custom_endpoint").(string)
	config.CertificateManagerBasePath = d.Get("certificate_manager_custom_endpoint").(string)
	config.CloudAssetBasePath = d.Get("cloud_asset_custom_endpoint").(string)
	config.CloudBuildBasePath = d.Get("cloud_build_custom_endpoint").(string)
	config.CloudFunctionsBasePath = d.Get("cloud_functions_custom_endpoint").(string)
	config.Cloudfunctions2BasePath = d.Get("cloudfunctions2_custom_endpoint").(string)
	config.CloudIdentityBasePath = d.Get("cloud_identity_custom_endpoint").(string)
	config.CloudIdsBasePath = d.Get("cloud_ids_custom_endpoint").(string)
	config.CloudIotBasePath = d.Get("cloud_iot_custom_endpoint").(string)
	config.CloudRunBasePath = d.Get("cloud_run_custom_endpoint").(string)
	config.CloudRunV2BasePath = d.Get("cloud_run_v2_custom_endpoint").(string)
	config.CloudSchedulerBasePath = d.Get("cloud_scheduler_custom_endpoint").(string)
	config.CloudTasksBasePath = d.Get("cloud_tasks_custom_endpoint").(string)
	config.ComputeBasePath = d.Get("compute_custom_endpoint").(string)
	config.ContainerAnalysisBasePath = d.Get("container_analysis_custom_endpoint").(string)
	config.ContainerAttachedBasePath = d.Get("container_attached_custom_endpoint").(string)
	config.DataCatalogBasePath = d.Get("data_catalog_custom_endpoint").(string)
	config.DataFusionBasePath = d.Get("data_fusion_custom_endpoint").(string)
	config.DataLossPreventionBasePath = d.Get("data_loss_prevention_custom_endpoint").(string)
	config.DataplexBasePath = d.Get("dataplex_custom_endpoint").(string)
	config.DataprocBasePath = d.Get("dataproc_custom_endpoint").(string)
	config.DataprocMetastoreBasePath = d.Get("dataproc_metastore_custom_endpoint").(string)
	config.DatastoreBasePath = d.Get("datastore_custom_endpoint").(string)
	config.DatastreamBasePath = d.Get("datastream_custom_endpoint").(string)
	config.DeploymentManagerBasePath = d.Get("deployment_manager_custom_endpoint").(string)
	config.DialogflowBasePath = d.Get("dialogflow_custom_endpoint").(string)
	config.DialogflowCXBasePath = d.Get("dialogflow_cx_custom_endpoint").(string)
	config.DNSBasePath = d.Get("dns_custom_endpoint").(string)
	config.DocumentAIBasePath = d.Get("document_ai_custom_endpoint").(string)
	config.EssentialContactsBasePath = d.Get("essential_contacts_custom_endpoint").(string)
	config.FilestoreBasePath = d.Get("filestore_custom_endpoint").(string)
	config.FirestoreBasePath = d.Get("firestore_custom_endpoint").(string)
	config.GameServicesBasePath = d.Get("game_services_custom_endpoint").(string)
	config.GKEBackupBasePath = d.Get("gke_backup_custom_endpoint").(string)
	config.GKEHubBasePath = d.Get("gke_hub_custom_endpoint").(string)
	config.HealthcareBasePath = d.Get("healthcare_custom_endpoint").(string)
	config.IAM2BasePath = d.Get("iam2_custom_endpoint").(string)
	config.IAMBetaBasePath = d.Get("iam_beta_custom_endpoint").(string)
	config.IAMWorkforcePoolBasePath = d.Get("iam_workforce_pool_custom_endpoint").(string)
	config.IapBasePath = d.Get("iap_custom_endpoint").(string)
	config.IdentityPlatformBasePath = d.Get("identity_platform_custom_endpoint").(string)
	config.KMSBasePath = d.Get("kms_custom_endpoint").(string)
	config.LoggingBasePath = d.Get("logging_custom_endpoint").(string)
	config.MemcacheBasePath = d.Get("memcache_custom_endpoint").(string)
	config.MLEngineBasePath = d.Get("ml_engine_custom_endpoint").(string)
	config.MonitoringBasePath = d.Get("monitoring_custom_endpoint").(string)
	config.NetworkManagementBasePath = d.Get("network_management_custom_endpoint").(string)
	config.NetworkServicesBasePath = d.Get("network_services_custom_endpoint").(string)
	config.NotebooksBasePath = d.Get("notebooks_custom_endpoint").(string)
	config.OSConfigBasePath = d.Get("os_config_custom_endpoint").(string)
	config.OSLoginBasePath = d.Get("os_login_custom_endpoint").(string)
	config.PrivatecaBasePath = d.Get("privateca_custom_endpoint").(string)
	config.PubsubBasePath = d.Get("pubsub_custom_endpoint").(string)
	config.PubsubLiteBasePath = d.Get("pubsub_lite_custom_endpoint").(string)
	config.RedisBasePath = d.Get("redis_custom_endpoint").(string)
	config.ResourceManagerBasePath = d.Get("resource_manager_custom_endpoint").(string)
	config.SecretManagerBasePath = d.Get("secret_manager_custom_endpoint").(string)
	config.SecurityCenterBasePath = d.Get("security_center_custom_endpoint").(string)
	config.ServiceManagementBasePath = d.Get("service_management_custom_endpoint").(string)
	config.ServiceUsageBasePath = d.Get("service_usage_custom_endpoint").(string)
	config.SourceRepoBasePath = d.Get("source_repo_custom_endpoint").(string)
	config.SpannerBasePath = d.Get("spanner_custom_endpoint").(string)
	config.SQLBasePath = d.Get("sql_custom_endpoint").(string)
	config.StorageBasePath = d.Get("storage_custom_endpoint").(string)
	config.StorageTransferBasePath = d.Get("storage_transfer_custom_endpoint").(string)
	config.TagsBasePath = d.Get("tags_custom_endpoint").(string)
	config.TPUBasePath = d.Get("tpu_custom_endpoint").(string)
	config.VertexAIBasePath = d.Get("vertex_ai_custom_endpoint").(string)
	config.VPCAccessBasePath = d.Get("vpc_access_custom_endpoint").(string)
	config.WorkflowsBasePath = d.Get("workflows_custom_endpoint").(string)

	// Handwritten Products / Versioned / Atypical Entries
	config.CloudBillingBasePath = d.Get(CloudBillingCustomEndpointEntryKey).(string)
	config.ComposerBasePath = d.Get(ComposerCustomEndpointEntryKey).(string)
	config.ContainerBasePath = d.Get(ContainerCustomEndpointEntryKey).(string)
	config.DataflowBasePath = d.Get(DataflowCustomEndpointEntryKey).(string)
	config.IamCredentialsBasePath = d.Get(IamCredentialsCustomEndpointEntryKey).(string)
	config.ResourceManagerV3BasePath = d.Get(ResourceManagerV3CustomEndpointEntryKey).(string)
	config.IAMBasePath = d.Get(IAMCustomEndpointEntryKey).(string)
	config.ServiceNetworkingBasePath = d.Get(ServiceNetworkingCustomEndpointEntryKey).(string)
	config.ServiceUsageBasePath = d.Get(ServiceUsageCustomEndpointEntryKey).(string)
	config.BigtableAdminBasePath = d.Get(BigtableAdminCustomEndpointEntryKey).(string)
	config.TagsLocationBasePath = d.Get(TagsLocationCustomEndpointEntryKey).(string)

	// dcl
	config.ContainerAwsBasePath = d.Get(ContainerAwsCustomEndpointEntryKey).(string)
	config.ContainerAzureBasePath = d.Get(ContainerAzureCustomEndpointEntryKey).(string)

	stopCtx, ok := schema.StopContext(ctx)
	if !ok {
		stopCtx = ctx
	}
	if err := config.LoadAndValidate(stopCtx); err != nil {
		return nil, diag.FromErr(err)
	}

	return ProviderDCLConfigure(d, &config), nil
}

func validateCredentials(v interface{}, k string) (warnings []string, errors []error) {
	if v == nil || v.(string) == "" {
		return
	}
	creds := v.(string)
	// if this is a path and we can stat it, assume it's ok
	if _, err := os.Stat(creds); err == nil {
		return
	}
	if _, err := googleoauth.CredentialsFromJSON(context.Background(), []byte(creds)); err != nil {
		errors = append(errors,
			fmt.Errorf("JSON credentials are not valid: %s", err))
	}

	return
}
