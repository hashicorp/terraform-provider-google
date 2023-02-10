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
			"google_access_approval_folder_service_account":       dataSourceAccessApprovalFolderServiceAccount(),
			"google_access_approval_organization_service_account": dataSourceAccessApprovalOrganizationServiceAccount(),
			"google_access_approval_project_service_account":      dataSourceAccessApprovalProjectServiceAccount(),
			"google_active_folder":                                dataSourceGoogleActiveFolder(),
			"google_artifact_registry_repository":                 dataSourceArtifactRegistryRepository(),
			"google_app_engine_default_service_account":           dataSourceGoogleAppEngineDefaultServiceAccount(),
			"google_beyondcorp_app_connection":                    dataSourceGoogleBeyondcorpAppConnection(),
			"google_beyondcorp_app_connector":                     dataSourceGoogleBeyondcorpAppConnector(),
			"google_beyondcorp_app_gateway":                       dataSourceGoogleBeyondcorpAppGateway(),
			"google_billing_account":                              dataSourceGoogleBillingAccount(),
			"google_bigquery_default_service_account":             dataSourceGoogleBigqueryDefaultServiceAccount(),
			"google_client_config":                                dataSourceGoogleClientConfig(),
			"google_client_openid_userinfo":                       dataSourceGoogleClientOpenIDUserinfo(),
			"google_cloudbuild_trigger":                           dataSourceGoogleCloudBuildTrigger(),
			"google_cloudfunctions_function":                      dataSourceGoogleCloudFunctionsFunction(),
			"google_cloudfunctions2_function":                     dataSourceGoogleCloudFunctions2Function(),
			"google_cloud_identity_groups":                        dataSourceGoogleCloudIdentityGroups(),
			"google_cloud_identity_group_memberships":             dataSourceGoogleCloudIdentityGroupMemberships(),
			"google_cloud_run_locations":                          dataSourceGoogleCloudRunLocations(),
			"google_cloud_run_service":                            dataSourceGoogleCloudRunService(),
			"google_composer_environment":                         dataSourceGoogleComposerEnvironment(),
			"google_composer_image_versions":                      dataSourceGoogleComposerImageVersions(),
			"google_compute_address":                              dataSourceGoogleComputeAddress(),
			"google_compute_addresses":                            dataSourceGoogleComputeAddresses(),
			"google_compute_backend_service":                      dataSourceGoogleComputeBackendService(),
			"google_compute_backend_bucket":                       dataSourceGoogleComputeBackendBucket(),
			"google_compute_default_service_account":              dataSourceGoogleComputeDefaultServiceAccount(),
			"google_compute_disk":                                 dataSourceGoogleComputeDisk(),
			"google_compute_forwarding_rule":                      dataSourceGoogleComputeForwardingRule(),
			"google_compute_global_address":                       dataSourceGoogleComputeGlobalAddress(),
			"google_compute_global_forwarding_rule":               dataSourceGoogleComputeGlobalForwardingRule(),
			"google_compute_ha_vpn_gateway":                       dataSourceGoogleComputeHaVpnGateway(),
			"google_compute_health_check":                         dataSourceGoogleComputeHealthCheck(),
			"google_compute_image":                                dataSourceGoogleComputeImage(),
			"google_compute_instance":                             dataSourceGoogleComputeInstance(),
			"google_compute_instance_group":                       dataSourceGoogleComputeInstanceGroup(),
			"google_compute_instance_group_manager":               dataSourceGoogleComputeInstanceGroupManager(),
			"google_compute_instance_serial_port":                 dataSourceGoogleComputeInstanceSerialPort(),
			"google_compute_instance_template":                    dataSourceGoogleComputeInstanceTemplate(),
			"google_compute_lb_ip_ranges":                         dataSourceGoogleComputeLbIpRanges(),
			"google_compute_network":                              dataSourceGoogleComputeNetwork(),
			"google_compute_network_endpoint_group":               dataSourceGoogleComputeNetworkEndpointGroup(),
			"google_compute_network_peering":                      dataSourceComputeNetworkPeering(),
			"google_compute_node_types":                           dataSourceGoogleComputeNodeTypes(),
			"google_compute_regions":                              dataSourceGoogleComputeRegions(),
			"google_compute_region_network_endpoint_group":        dataSourceGoogleComputeRegionNetworkEndpointGroup(),
			"google_compute_region_instance_group":                dataSourceGoogleComputeRegionInstanceGroup(),
			"google_compute_region_ssl_certificate":               dataSourceGoogleRegionComputeSslCertificate(),
			"google_compute_resource_policy":                      dataSourceGoogleComputeResourcePolicy(),
			"google_compute_router":                               dataSourceGoogleComputeRouter(),
			"google_compute_router_nat":                           dataSourceGoogleComputeRouterNat(),
			"google_compute_router_status":                        dataSourceGoogleComputeRouterStatus(),
			"google_compute_snapshot":                             dataSourceGoogleComputeSnapshot(),
			"google_compute_ssl_certificate":                      dataSourceGoogleComputeSslCertificate(),
			"google_compute_ssl_policy":                           dataSourceGoogleComputeSslPolicy(),
			"google_compute_subnetwork":                           dataSourceGoogleComputeSubnetwork(),
			"google_compute_vpn_gateway":                          dataSourceGoogleComputeVpnGateway(),
			"google_compute_zones":                                dataSourceGoogleComputeZones(),
			"google_container_azure_versions":                     dataSourceGoogleContainerAzureVersions(),
			"google_container_aws_versions":                       dataSourceGoogleContainerAwsVersions(),
			"google_container_attached_versions":                  dataSourceGoogleContainerAttachedVersions(),
			"google_container_attached_install_manifest":          dataSourceGoogleContainerAttachedInstallManifest(),
			"google_container_cluster":                            dataSourceGoogleContainerCluster(),
			"google_container_engine_versions":                    dataSourceGoogleContainerEngineVersions(),
			"google_container_registry_image":                     dataSourceGoogleContainerImage(),
			"google_container_registry_repository":                dataSourceGoogleContainerRepo(),
			"google_dataproc_metastore_service":                   dataSourceDataprocMetastoreService(),
			"google_game_services_game_server_deployment_rollout": dataSourceGameServicesGameServerDeploymentRollout(),
			"google_iam_policy":                                   dataSourceGoogleIamPolicy(),
			"google_iam_role":                                     dataSourceGoogleIamRole(),
			"google_iam_testable_permissions":                     dataSourceGoogleIamTestablePermissions(),
			"google_iap_client":                                   dataSourceGoogleIapClient(),
			"google_kms_crypto_key":                               dataSourceGoogleKmsCryptoKey(),
			"google_kms_crypto_key_version":                       dataSourceGoogleKmsCryptoKeyVersion(),
			"google_kms_key_ring":                                 dataSourceGoogleKmsKeyRing(),
			"google_kms_secret":                                   dataSourceGoogleKmsSecret(),
			"google_kms_secret_ciphertext":                        dataSourceGoogleKmsSecretCiphertext(),
			"google_folder":                                       dataSourceGoogleFolder(),
			"google_folders":                                      dataSourceGoogleFolders(),
			"google_folder_organization_policy":                   dataSourceGoogleFolderOrganizationPolicy(),
			"google_logging_project_cmek_settings":                dataSourceGoogleLoggingProjectCmekSettings(),
			"google_monitoring_notification_channel":              dataSourceMonitoringNotificationChannel(),
			"google_monitoring_cluster_istio_service":             dataSourceMonitoringServiceClusterIstio(),
			"google_monitoring_istio_canonical_service":           dataSourceMonitoringIstioCanonicalService(),
			"google_monitoring_mesh_istio_service":                dataSourceMonitoringServiceMeshIstio(),
			"google_monitoring_app_engine_service":                dataSourceMonitoringServiceAppEngine(),
			"google_monitoring_uptime_check_ips":                  dataSourceGoogleMonitoringUptimeCheckIps(),
			"google_netblock_ip_ranges":                           dataSourceGoogleNetblockIpRanges(),
			"google_organization":                                 dataSourceGoogleOrganization(),
			"google_privateca_certificate_authority":              dataSourcePrivatecaCertificateAuthority(),
			"google_project":                                      dataSourceGoogleProject(),
			"google_projects":                                     dataSourceGoogleProjects(),
			"google_project_organization_policy":                  dataSourceGoogleProjectOrganizationPolicy(),
			"google_project_service":                              dataSourceGoogleProjectService(),
			"google_pubsub_subscription":                          dataSourceGooglePubsubSubscription(),
			"google_pubsub_topic":                                 dataSourceGooglePubsubTopic(),
			"google_secret_manager_secret":                        dataSourceSecretManagerSecret(),
			"google_secret_manager_secret_version":                dataSourceSecretManagerSecretVersion(),
			"google_secret_manager_secret_version_access":         dataSourceSecretManagerSecretVersionAccess(),
			"google_service_account":                              dataSourceGoogleServiceAccount(),
			"google_service_account_access_token":                 dataSourceGoogleServiceAccountAccessToken(),
			"google_service_account_id_token":                     dataSourceGoogleServiceAccountIdToken(),
			"google_service_account_jwt":                          dataSourceGoogleServiceAccountJwt(),
			"google_service_account_key":                          dataSourceGoogleServiceAccountKey(),
			"google_sourcerepo_repository":                        dataSourceGoogleSourceRepoRepository(),
			"google_spanner_instance":                             dataSourceSpannerInstance(),
			"google_sql_ca_certs":                                 dataSourceGoogleSQLCaCerts(),
			"google_sql_backup_run":                               dataSourceSqlBackupRun(),
			"google_sql_database":                                 dataSourceSqlDatabase(),
			"google_sql_database_instance":                        dataSourceSqlDatabaseInstance(),
			"google_sql_database_instances":                       dataSourceSqlDatabaseInstances(),
			"google_service_networking_peered_dns_domain":         dataSourceGoogleServiceNetworkingPeeredDNSDomain(),
			"google_storage_bucket":                               dataSourceGoogleStorageBucket(),
			"google_storage_bucket_object":                        dataSourceGoogleStorageBucketObject(),
			"google_storage_bucket_object_content":                dataSourceGoogleStorageBucketObjectContent(),
			"google_storage_object_signed_url":                    dataSourceGoogleSignedUrl(),
			"google_storage_project_service_account":              dataSourceGoogleStorageProjectServiceAccount(),
			"google_storage_transfer_project_service_account":     dataSourceGoogleStorageTransferProjectServiceAccount(),
			"google_tags_tag_key":                                 dataSourceGoogleTagsTagKey(),
			"google_tags_tag_value":                               dataSourceGoogleTagsTagValue(),
			"google_tpu_tensorflow_versions":                      dataSourceTpuTensorflowVersions(),
			"google_vpc_access_connector":                         dataSourceVPCAccessConnector(),
			"google_redis_instance":                               dataSourceGoogleRedisInstance(),
			// ####### END datasources ###########
		},
		ResourcesMap: ResourceMap(),
	}

	provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return providerConfigure(ctx, d, provider)
	}

	configureDCLProvider(provider)

	return provider
}

// Generated resources: 264
// Generated IAM resources: 168
// Total generated resources: 432
func ResourceMap() map[string]*schema.Resource {
	resourceMap, _ := ResourceMapWithErrors()
	return resourceMap
}

func ResourceMapWithErrors() (map[string]*schema.Resource, error) {
	return mergeResourceMaps(
		map[string]*schema.Resource{
			"google_folder_access_approval_settings":                       resourceAccessApprovalFolderSettings(),
			"google_project_access_approval_settings":                      resourceAccessApprovalProjectSettings(),
			"google_organization_access_approval_settings":                 resourceAccessApprovalOrganizationSettings(),
			"google_access_context_manager_access_level":                   resourceAccessContextManagerAccessLevel(),
			"google_access_context_manager_access_level_condition":         resourceAccessContextManagerAccessLevelCondition(),
			"google_access_context_manager_access_levels":                  resourceAccessContextManagerAccessLevels(),
			"google_access_context_manager_access_policy":                  resourceAccessContextManagerAccessPolicy(),
			"google_access_context_manager_access_policy_iam_binding":      ResourceIamBinding(AccessContextManagerAccessPolicyIamSchema, AccessContextManagerAccessPolicyIamUpdaterProducer, AccessContextManagerAccessPolicyIdParseFunc),
			"google_access_context_manager_access_policy_iam_member":       ResourceIamMember(AccessContextManagerAccessPolicyIamSchema, AccessContextManagerAccessPolicyIamUpdaterProducer, AccessContextManagerAccessPolicyIdParseFunc),
			"google_access_context_manager_access_policy_iam_policy":       ResourceIamPolicy(AccessContextManagerAccessPolicyIamSchema, AccessContextManagerAccessPolicyIamUpdaterProducer, AccessContextManagerAccessPolicyIdParseFunc),
			"google_access_context_manager_gcp_user_access_binding":        resourceAccessContextManagerGcpUserAccessBinding(),
			"google_access_context_manager_service_perimeter":              resourceAccessContextManagerServicePerimeter(),
			"google_access_context_manager_service_perimeter_resource":     resourceAccessContextManagerServicePerimeterResource(),
			"google_access_context_manager_service_perimeters":             resourceAccessContextManagerServicePerimeters(),
			"google_active_directory_domain":                               resourceActiveDirectoryDomain(),
			"google_active_directory_domain_trust":                         resourceActiveDirectoryDomainTrust(),
			"google_alloydb_backup":                                        resourceAlloydbBackup(),
			"google_alloydb_cluster":                                       resourceAlloydbCluster(),
			"google_alloydb_instance":                                      resourceAlloydbInstance(),
			"google_apigee_addons_config":                                  resourceApigeeAddonsConfig(),
			"google_apigee_endpoint_attachment":                            resourceApigeeEndpointAttachment(),
			"google_apigee_envgroup":                                       resourceApigeeEnvgroup(),
			"google_apigee_envgroup_attachment":                            resourceApigeeEnvgroupAttachment(),
			"google_apigee_environment":                                    resourceApigeeEnvironment(),
			"google_apigee_environment_iam_binding":                        ResourceIamBinding(ApigeeEnvironmentIamSchema, ApigeeEnvironmentIamUpdaterProducer, ApigeeEnvironmentIdParseFunc),
			"google_apigee_environment_iam_member":                         ResourceIamMember(ApigeeEnvironmentIamSchema, ApigeeEnvironmentIamUpdaterProducer, ApigeeEnvironmentIdParseFunc),
			"google_apigee_environment_iam_policy":                         ResourceIamPolicy(ApigeeEnvironmentIamSchema, ApigeeEnvironmentIamUpdaterProducer, ApigeeEnvironmentIdParseFunc),
			"google_apigee_instance":                                       resourceApigeeInstance(),
			"google_apigee_instance_attachment":                            resourceApigeeInstanceAttachment(),
			"google_apigee_nat_address":                                    resourceApigeeNatAddress(),
			"google_apigee_organization":                                   resourceApigeeOrganization(),
			"google_apigee_sync_authorization":                             resourceApigeeSyncAuthorization(),
			"google_app_engine_application_url_dispatch_rules":             resourceAppEngineApplicationUrlDispatchRules(),
			"google_app_engine_domain_mapping":                             resourceAppEngineDomainMapping(),
			"google_app_engine_firewall_rule":                              resourceAppEngineFirewallRule(),
			"google_app_engine_flexible_app_version":                       resourceAppEngineFlexibleAppVersion(),
			"google_app_engine_service_network_settings":                   resourceAppEngineServiceNetworkSettings(),
			"google_app_engine_service_split_traffic":                      resourceAppEngineServiceSplitTraffic(),
			"google_app_engine_standard_app_version":                       resourceAppEngineStandardAppVersion(),
			"google_artifact_registry_repository":                          resourceArtifactRegistryRepository(),
			"google_artifact_registry_repository_iam_binding":              ResourceIamBinding(ArtifactRegistryRepositoryIamSchema, ArtifactRegistryRepositoryIamUpdaterProducer, ArtifactRegistryRepositoryIdParseFunc),
			"google_artifact_registry_repository_iam_member":               ResourceIamMember(ArtifactRegistryRepositoryIamSchema, ArtifactRegistryRepositoryIamUpdaterProducer, ArtifactRegistryRepositoryIdParseFunc),
			"google_artifact_registry_repository_iam_policy":               ResourceIamPolicy(ArtifactRegistryRepositoryIamSchema, ArtifactRegistryRepositoryIamUpdaterProducer, ArtifactRegistryRepositoryIdParseFunc),
			"google_beyondcorp_app_connection":                             resourceBeyondcorpAppConnection(),
			"google_beyondcorp_app_connector":                              resourceBeyondcorpAppConnector(),
			"google_beyondcorp_app_gateway":                                resourceBeyondcorpAppGateway(),
			"google_bigquery_dataset":                                      resourceBigQueryDataset(),
			"google_bigquery_dataset_access":                               resourceBigQueryDatasetAccess(),
			"google_bigquery_job":                                          resourceBigQueryJob(),
			"google_bigquery_table_iam_binding":                            ResourceIamBinding(BigQueryTableIamSchema, BigQueryTableIamUpdaterProducer, BigQueryTableIdParseFunc),
			"google_bigquery_table_iam_member":                             ResourceIamMember(BigQueryTableIamSchema, BigQueryTableIamUpdaterProducer, BigQueryTableIdParseFunc),
			"google_bigquery_table_iam_policy":                             ResourceIamPolicy(BigQueryTableIamSchema, BigQueryTableIamUpdaterProducer, BigQueryTableIdParseFunc),
			"google_bigquery_routine":                                      resourceBigQueryRoutine(),
			"google_bigquery_analytics_hub_data_exchange":                  resourceBigqueryAnalyticsHubDataExchange(),
			"google_bigquery_analytics_hub_data_exchange_iam_binding":      ResourceIamBinding(BigqueryAnalyticsHubDataExchangeIamSchema, BigqueryAnalyticsHubDataExchangeIamUpdaterProducer, BigqueryAnalyticsHubDataExchangeIdParseFunc),
			"google_bigquery_analytics_hub_data_exchange_iam_member":       ResourceIamMember(BigqueryAnalyticsHubDataExchangeIamSchema, BigqueryAnalyticsHubDataExchangeIamUpdaterProducer, BigqueryAnalyticsHubDataExchangeIdParseFunc),
			"google_bigquery_analytics_hub_data_exchange_iam_policy":       ResourceIamPolicy(BigqueryAnalyticsHubDataExchangeIamSchema, BigqueryAnalyticsHubDataExchangeIamUpdaterProducer, BigqueryAnalyticsHubDataExchangeIdParseFunc),
			"google_bigquery_analytics_hub_listing":                        resourceBigqueryAnalyticsHubListing(),
			"google_bigquery_analytics_hub_listing_iam_binding":            ResourceIamBinding(BigqueryAnalyticsHubListingIamSchema, BigqueryAnalyticsHubListingIamUpdaterProducer, BigqueryAnalyticsHubListingIdParseFunc),
			"google_bigquery_analytics_hub_listing_iam_member":             ResourceIamMember(BigqueryAnalyticsHubListingIamSchema, BigqueryAnalyticsHubListingIamUpdaterProducer, BigqueryAnalyticsHubListingIdParseFunc),
			"google_bigquery_analytics_hub_listing_iam_policy":             ResourceIamPolicy(BigqueryAnalyticsHubListingIamSchema, BigqueryAnalyticsHubListingIamUpdaterProducer, BigqueryAnalyticsHubListingIdParseFunc),
			"google_bigquery_connection":                                   resourceBigqueryConnectionConnection(),
			"google_bigquery_connection_iam_binding":                       ResourceIamBinding(BigqueryConnectionConnectionIamSchema, BigqueryConnectionConnectionIamUpdaterProducer, BigqueryConnectionConnectionIdParseFunc),
			"google_bigquery_connection_iam_member":                        ResourceIamMember(BigqueryConnectionConnectionIamSchema, BigqueryConnectionConnectionIamUpdaterProducer, BigqueryConnectionConnectionIdParseFunc),
			"google_bigquery_connection_iam_policy":                        ResourceIamPolicy(BigqueryConnectionConnectionIamSchema, BigqueryConnectionConnectionIamUpdaterProducer, BigqueryConnectionConnectionIdParseFunc),
			"google_bigquery_data_transfer_config":                         resourceBigqueryDataTransferConfig(),
			"google_bigquery_reservation":                                  resourceBigqueryReservationReservation(),
			"google_bigtable_app_profile":                                  resourceBigtableAppProfile(),
			"google_billing_budget":                                        resourceBillingBudget(),
			"google_binary_authorization_attestor":                         resourceBinaryAuthorizationAttestor(),
			"google_binary_authorization_attestor_iam_binding":             ResourceIamBinding(BinaryAuthorizationAttestorIamSchema, BinaryAuthorizationAttestorIamUpdaterProducer, BinaryAuthorizationAttestorIdParseFunc),
			"google_binary_authorization_attestor_iam_member":              ResourceIamMember(BinaryAuthorizationAttestorIamSchema, BinaryAuthorizationAttestorIamUpdaterProducer, BinaryAuthorizationAttestorIdParseFunc),
			"google_binary_authorization_attestor_iam_policy":              ResourceIamPolicy(BinaryAuthorizationAttestorIamSchema, BinaryAuthorizationAttestorIamUpdaterProducer, BinaryAuthorizationAttestorIdParseFunc),
			"google_binary_authorization_policy":                           resourceBinaryAuthorizationPolicy(),
			"google_certificate_manager_dns_authorization":                 resourceCertificateManagerDnsAuthorization(),
			"google_certificate_manager_certificate":                       resourceCertificateManagerCertificate(),
			"google_certificate_manager_certificate_map":                   resourceCertificateManagerCertificateMap(),
			"google_certificate_manager_certificate_map_entry":             resourceCertificateManagerCertificateMapEntry(),
			"google_cloud_asset_project_feed":                              resourceCloudAssetProjectFeed(),
			"google_cloud_asset_folder_feed":                               resourceCloudAssetFolderFeed(),
			"google_cloud_asset_organization_feed":                         resourceCloudAssetOrganizationFeed(),
			"google_cloudbuild_trigger":                                    resourceCloudBuildTrigger(),
			"google_cloudfunctions_function_iam_binding":                   ResourceIamBinding(CloudFunctionsCloudFunctionIamSchema, CloudFunctionsCloudFunctionIamUpdaterProducer, CloudFunctionsCloudFunctionIdParseFunc),
			"google_cloudfunctions_function_iam_member":                    ResourceIamMember(CloudFunctionsCloudFunctionIamSchema, CloudFunctionsCloudFunctionIamUpdaterProducer, CloudFunctionsCloudFunctionIdParseFunc),
			"google_cloudfunctions_function_iam_policy":                    ResourceIamPolicy(CloudFunctionsCloudFunctionIamSchema, CloudFunctionsCloudFunctionIamUpdaterProducer, CloudFunctionsCloudFunctionIdParseFunc),
			"google_cloudfunctions2_function":                              resourceCloudfunctions2function(),
			"google_cloudfunctions2_function_iam_binding":                  ResourceIamBinding(Cloudfunctions2functionIamSchema, Cloudfunctions2functionIamUpdaterProducer, Cloudfunctions2functionIdParseFunc),
			"google_cloudfunctions2_function_iam_member":                   ResourceIamMember(Cloudfunctions2functionIamSchema, Cloudfunctions2functionIamUpdaterProducer, Cloudfunctions2functionIdParseFunc),
			"google_cloudfunctions2_function_iam_policy":                   ResourceIamPolicy(Cloudfunctions2functionIamSchema, Cloudfunctions2functionIamUpdaterProducer, Cloudfunctions2functionIdParseFunc),
			"google_cloud_identity_group":                                  resourceCloudIdentityGroup(),
			"google_cloud_identity_group_membership":                       resourceCloudIdentityGroupMembership(),
			"google_cloud_ids_endpoint":                                    resourceCloudIdsEndpoint(),
			"google_cloudiot_registry":                                     resourceCloudIotDeviceRegistry(),
			"google_cloudiot_registry_iam_binding":                         ResourceIamBinding(CloudIotDeviceRegistryIamSchema, CloudIotDeviceRegistryIamUpdaterProducer, CloudIotDeviceRegistryIdParseFunc),
			"google_cloudiot_registry_iam_member":                          ResourceIamMember(CloudIotDeviceRegistryIamSchema, CloudIotDeviceRegistryIamUpdaterProducer, CloudIotDeviceRegistryIdParseFunc),
			"google_cloudiot_registry_iam_policy":                          ResourceIamPolicy(CloudIotDeviceRegistryIamSchema, CloudIotDeviceRegistryIamUpdaterProducer, CloudIotDeviceRegistryIdParseFunc),
			"google_cloudiot_device":                                       resourceCloudIotDevice(),
			"google_cloud_run_domain_mapping":                              resourceCloudRunDomainMapping(),
			"google_cloud_run_service":                                     resourceCloudRunService(),
			"google_cloud_run_service_iam_binding":                         ResourceIamBinding(CloudRunServiceIamSchema, CloudRunServiceIamUpdaterProducer, CloudRunServiceIdParseFunc),
			"google_cloud_run_service_iam_member":                          ResourceIamMember(CloudRunServiceIamSchema, CloudRunServiceIamUpdaterProducer, CloudRunServiceIdParseFunc),
			"google_cloud_run_service_iam_policy":                          ResourceIamPolicy(CloudRunServiceIamSchema, CloudRunServiceIamUpdaterProducer, CloudRunServiceIdParseFunc),
			"google_cloud_run_v2_job":                                      resourceCloudRunV2Job(),
			"google_cloud_run_v2_job_iam_binding":                          ResourceIamBinding(CloudRunV2JobIamSchema, CloudRunV2JobIamUpdaterProducer, CloudRunV2JobIdParseFunc),
			"google_cloud_run_v2_job_iam_member":                           ResourceIamMember(CloudRunV2JobIamSchema, CloudRunV2JobIamUpdaterProducer, CloudRunV2JobIdParseFunc),
			"google_cloud_run_v2_job_iam_policy":                           ResourceIamPolicy(CloudRunV2JobIamSchema, CloudRunV2JobIamUpdaterProducer, CloudRunV2JobIdParseFunc),
			"google_cloud_run_v2_service":                                  resourceCloudRunV2Service(),
			"google_cloud_run_v2_service_iam_binding":                      ResourceIamBinding(CloudRunV2ServiceIamSchema, CloudRunV2ServiceIamUpdaterProducer, CloudRunV2ServiceIdParseFunc),
			"google_cloud_run_v2_service_iam_member":                       ResourceIamMember(CloudRunV2ServiceIamSchema, CloudRunV2ServiceIamUpdaterProducer, CloudRunV2ServiceIdParseFunc),
			"google_cloud_run_v2_service_iam_policy":                       ResourceIamPolicy(CloudRunV2ServiceIamSchema, CloudRunV2ServiceIamUpdaterProducer, CloudRunV2ServiceIdParseFunc),
			"google_cloud_scheduler_job":                                   resourceCloudSchedulerJob(),
			"google_cloud_tasks_queue":                                     resourceCloudTasksQueue(),
			"google_cloud_tasks_queue_iam_binding":                         ResourceIamBinding(CloudTasksQueueIamSchema, CloudTasksQueueIamUpdaterProducer, CloudTasksQueueIdParseFunc),
			"google_cloud_tasks_queue_iam_member":                          ResourceIamMember(CloudTasksQueueIamSchema, CloudTasksQueueIamUpdaterProducer, CloudTasksQueueIdParseFunc),
			"google_cloud_tasks_queue_iam_policy":                          ResourceIamPolicy(CloudTasksQueueIamSchema, CloudTasksQueueIamUpdaterProducer, CloudTasksQueueIdParseFunc),
			"google_compute_address":                                       resourceComputeAddress(),
			"google_compute_autoscaler":                                    resourceComputeAutoscaler(),
			"google_compute_backend_bucket":                                resourceComputeBackendBucket(),
			"google_compute_backend_bucket_signed_url_key":                 resourceComputeBackendBucketSignedUrlKey(),
			"google_compute_backend_service":                               resourceComputeBackendService(),
			"google_compute_region_backend_service":                        resourceComputeRegionBackendService(),
			"google_compute_backend_service_signed_url_key":                resourceComputeBackendServiceSignedUrlKey(),
			"google_compute_region_disk_resource_policy_attachment":        resourceComputeRegionDiskResourcePolicyAttachment(),
			"google_compute_disk_resource_policy_attachment":               resourceComputeDiskResourcePolicyAttachment(),
			"google_compute_disk":                                          resourceComputeDisk(),
			"google_compute_disk_iam_binding":                              ResourceIamBinding(ComputeDiskIamSchema, ComputeDiskIamUpdaterProducer, ComputeDiskIdParseFunc),
			"google_compute_disk_iam_member":                               ResourceIamMember(ComputeDiskIamSchema, ComputeDiskIamUpdaterProducer, ComputeDiskIdParseFunc),
			"google_compute_disk_iam_policy":                               ResourceIamPolicy(ComputeDiskIamSchema, ComputeDiskIamUpdaterProducer, ComputeDiskIdParseFunc),
			"google_compute_firewall":                                      resourceComputeFirewall(),
			"google_compute_forwarding_rule":                               resourceComputeForwardingRule(),
			"google_compute_global_address":                                resourceComputeGlobalAddress(),
			"google_compute_global_forwarding_rule":                        resourceComputeGlobalForwardingRule(),
			"google_compute_http_health_check":                             resourceComputeHttpHealthCheck(),
			"google_compute_https_health_check":                            resourceComputeHttpsHealthCheck(),
			"google_compute_health_check":                                  resourceComputeHealthCheck(),
			"google_compute_image":                                         resourceComputeImage(),
			"google_compute_image_iam_binding":                             ResourceIamBinding(ComputeImageIamSchema, ComputeImageIamUpdaterProducer, ComputeImageIdParseFunc),
			"google_compute_image_iam_member":                              ResourceIamMember(ComputeImageIamSchema, ComputeImageIamUpdaterProducer, ComputeImageIdParseFunc),
			"google_compute_image_iam_policy":                              ResourceIamPolicy(ComputeImageIamSchema, ComputeImageIamUpdaterProducer, ComputeImageIdParseFunc),
			"google_compute_instance_iam_binding":                          ResourceIamBinding(ComputeInstanceIamSchema, ComputeInstanceIamUpdaterProducer, ComputeInstanceIdParseFunc),
			"google_compute_instance_iam_member":                           ResourceIamMember(ComputeInstanceIamSchema, ComputeInstanceIamUpdaterProducer, ComputeInstanceIdParseFunc),
			"google_compute_instance_iam_policy":                           ResourceIamPolicy(ComputeInstanceIamSchema, ComputeInstanceIamUpdaterProducer, ComputeInstanceIdParseFunc),
			"google_compute_instance_group_named_port":                     resourceComputeInstanceGroupNamedPort(),
			"google_compute_interconnect_attachment":                       resourceComputeInterconnectAttachment(),
			"google_compute_network":                                       resourceComputeNetwork(),
			"google_compute_network_endpoint":                              resourceComputeNetworkEndpoint(),
			"google_compute_network_endpoint_group":                        resourceComputeNetworkEndpointGroup(),
			"google_compute_global_network_endpoint":                       resourceComputeGlobalNetworkEndpoint(),
			"google_compute_global_network_endpoint_group":                 resourceComputeGlobalNetworkEndpointGroup(),
			"google_compute_region_network_endpoint_group":                 resourceComputeRegionNetworkEndpointGroup(),
			"google_compute_node_group":                                    resourceComputeNodeGroup(),
			"google_compute_network_peering_routes_config":                 resourceComputeNetworkPeeringRoutesConfig(),
			"google_compute_node_template":                                 resourceComputeNodeTemplate(),
			"google_compute_packet_mirroring":                              resourceComputePacketMirroring(),
			"google_compute_per_instance_config":                           resourceComputePerInstanceConfig(),
			"google_compute_region_per_instance_config":                    resourceComputeRegionPerInstanceConfig(),
			"google_compute_region_autoscaler":                             resourceComputeRegionAutoscaler(),
			"google_compute_region_disk":                                   resourceComputeRegionDisk(),
			"google_compute_region_disk_iam_binding":                       ResourceIamBinding(ComputeRegionDiskIamSchema, ComputeRegionDiskIamUpdaterProducer, ComputeRegionDiskIdParseFunc),
			"google_compute_region_disk_iam_member":                        ResourceIamMember(ComputeRegionDiskIamSchema, ComputeRegionDiskIamUpdaterProducer, ComputeRegionDiskIdParseFunc),
			"google_compute_region_disk_iam_policy":                        ResourceIamPolicy(ComputeRegionDiskIamSchema, ComputeRegionDiskIamUpdaterProducer, ComputeRegionDiskIdParseFunc),
			"google_compute_region_url_map":                                resourceComputeRegionUrlMap(),
			"google_compute_region_health_check":                           resourceComputeRegionHealthCheck(),
			"google_compute_resource_policy":                               resourceComputeResourcePolicy(),
			"google_compute_route":                                         resourceComputeRoute(),
			"google_compute_router":                                        resourceComputeRouter(),
			"google_compute_router_nat":                                    resourceComputeRouterNat(),
			"google_compute_router_peer":                                   resourceComputeRouterBgpPeer(),
			"google_compute_snapshot":                                      resourceComputeSnapshot(),
			"google_compute_snapshot_iam_binding":                          ResourceIamBinding(ComputeSnapshotIamSchema, ComputeSnapshotIamUpdaterProducer, ComputeSnapshotIdParseFunc),
			"google_compute_snapshot_iam_member":                           ResourceIamMember(ComputeSnapshotIamSchema, ComputeSnapshotIamUpdaterProducer, ComputeSnapshotIdParseFunc),
			"google_compute_snapshot_iam_policy":                           ResourceIamPolicy(ComputeSnapshotIamSchema, ComputeSnapshotIamUpdaterProducer, ComputeSnapshotIdParseFunc),
			"google_compute_ssl_certificate":                               resourceComputeSslCertificate(),
			"google_compute_managed_ssl_certificate":                       resourceComputeManagedSslCertificate(),
			"google_compute_region_ssl_certificate":                        resourceComputeRegionSslCertificate(),
			"google_compute_reservation":                                   resourceComputeReservation(),
			"google_compute_service_attachment":                            resourceComputeServiceAttachment(),
			"google_compute_ssl_policy":                                    resourceComputeSslPolicy(),
			"google_compute_subnetwork":                                    resourceComputeSubnetwork(),
			"google_compute_subnetwork_iam_binding":                        ResourceIamBinding(ComputeSubnetworkIamSchema, ComputeSubnetworkIamUpdaterProducer, ComputeSubnetworkIdParseFunc),
			"google_compute_subnetwork_iam_member":                         ResourceIamMember(ComputeSubnetworkIamSchema, ComputeSubnetworkIamUpdaterProducer, ComputeSubnetworkIdParseFunc),
			"google_compute_subnetwork_iam_policy":                         ResourceIamPolicy(ComputeSubnetworkIamSchema, ComputeSubnetworkIamUpdaterProducer, ComputeSubnetworkIdParseFunc),
			"google_compute_target_http_proxy":                             resourceComputeTargetHttpProxy(),
			"google_compute_target_https_proxy":                            resourceComputeTargetHttpsProxy(),
			"google_compute_region_target_http_proxy":                      resourceComputeRegionTargetHttpProxy(),
			"google_compute_region_target_https_proxy":                     resourceComputeRegionTargetHttpsProxy(),
			"google_compute_region_target_tcp_proxy":                       resourceComputeRegionTargetTcpProxy(),
			"google_compute_target_instance":                               resourceComputeTargetInstance(),
			"google_compute_target_ssl_proxy":                              resourceComputeTargetSslProxy(),
			"google_compute_target_tcp_proxy":                              resourceComputeTargetTcpProxy(),
			"google_compute_vpn_gateway":                                   resourceComputeVpnGateway(),
			"google_compute_ha_vpn_gateway":                                resourceComputeHaVpnGateway(),
			"google_compute_external_vpn_gateway":                          resourceComputeExternalVpnGateway(),
			"google_compute_url_map":                                       resourceComputeUrlMap(),
			"google_compute_vpn_tunnel":                                    resourceComputeVpnTunnel(),
			"google_compute_target_grpc_proxy":                             resourceComputeTargetGrpcProxy(),
			"google_container_analysis_note":                               resourceContainerAnalysisNote(),
			"google_container_analysis_occurrence":                         resourceContainerAnalysisOccurrence(),
			"google_container_attached_cluster":                            resourceContainerAttachedCluster(),
			"google_data_catalog_entry_group":                              resourceDataCatalogEntryGroup(),
			"google_data_catalog_entry_group_iam_binding":                  ResourceIamBinding(DataCatalogEntryGroupIamSchema, DataCatalogEntryGroupIamUpdaterProducer, DataCatalogEntryGroupIdParseFunc),
			"google_data_catalog_entry_group_iam_member":                   ResourceIamMember(DataCatalogEntryGroupIamSchema, DataCatalogEntryGroupIamUpdaterProducer, DataCatalogEntryGroupIdParseFunc),
			"google_data_catalog_entry_group_iam_policy":                   ResourceIamPolicy(DataCatalogEntryGroupIamSchema, DataCatalogEntryGroupIamUpdaterProducer, DataCatalogEntryGroupIdParseFunc),
			"google_data_catalog_entry":                                    resourceDataCatalogEntry(),
			"google_data_catalog_tag_template":                             resourceDataCatalogTagTemplate(),
			"google_data_catalog_tag_template_iam_binding":                 ResourceIamBinding(DataCatalogTagTemplateIamSchema, DataCatalogTagTemplateIamUpdaterProducer, DataCatalogTagTemplateIdParseFunc),
			"google_data_catalog_tag_template_iam_member":                  ResourceIamMember(DataCatalogTagTemplateIamSchema, DataCatalogTagTemplateIamUpdaterProducer, DataCatalogTagTemplateIdParseFunc),
			"google_data_catalog_tag_template_iam_policy":                  ResourceIamPolicy(DataCatalogTagTemplateIamSchema, DataCatalogTagTemplateIamUpdaterProducer, DataCatalogTagTemplateIdParseFunc),
			"google_data_catalog_tag":                                      resourceDataCatalogTag(),
			"google_data_fusion_instance":                                  resourceDataFusionInstance(),
			"google_data_fusion_instance_iam_binding":                      ResourceIamBinding(DataFusionInstanceIamSchema, DataFusionInstanceIamUpdaterProducer, DataFusionInstanceIdParseFunc),
			"google_data_fusion_instance_iam_member":                       ResourceIamMember(DataFusionInstanceIamSchema, DataFusionInstanceIamUpdaterProducer, DataFusionInstanceIdParseFunc),
			"google_data_fusion_instance_iam_policy":                       ResourceIamPolicy(DataFusionInstanceIamSchema, DataFusionInstanceIamUpdaterProducer, DataFusionInstanceIdParseFunc),
			"google_data_loss_prevention_job_trigger":                      resourceDataLossPreventionJobTrigger(),
			"google_data_loss_prevention_inspect_template":                 resourceDataLossPreventionInspectTemplate(),
			"google_data_loss_prevention_stored_info_type":                 resourceDataLossPreventionStoredInfoType(),
			"google_data_loss_prevention_deidentify_template":              resourceDataLossPreventionDeidentifyTemplate(),
			"google_dataproc_autoscaling_policy":                           resourceDataprocAutoscalingPolicy(),
			"google_dataproc_autoscaling_policy_iam_binding":               ResourceIamBinding(DataprocAutoscalingPolicyIamSchema, DataprocAutoscalingPolicyIamUpdaterProducer, DataprocAutoscalingPolicyIdParseFunc),
			"google_dataproc_autoscaling_policy_iam_member":                ResourceIamMember(DataprocAutoscalingPolicyIamSchema, DataprocAutoscalingPolicyIamUpdaterProducer, DataprocAutoscalingPolicyIdParseFunc),
			"google_dataproc_autoscaling_policy_iam_policy":                ResourceIamPolicy(DataprocAutoscalingPolicyIamSchema, DataprocAutoscalingPolicyIamUpdaterProducer, DataprocAutoscalingPolicyIdParseFunc),
			"google_dataproc_metastore_service":                            resourceDataprocMetastoreService(),
			"google_dataproc_metastore_service_iam_binding":                ResourceIamBinding(DataprocMetastoreServiceIamSchema, DataprocMetastoreServiceIamUpdaterProducer, DataprocMetastoreServiceIdParseFunc),
			"google_dataproc_metastore_service_iam_member":                 ResourceIamMember(DataprocMetastoreServiceIamSchema, DataprocMetastoreServiceIamUpdaterProducer, DataprocMetastoreServiceIdParseFunc),
			"google_dataproc_metastore_service_iam_policy":                 ResourceIamPolicy(DataprocMetastoreServiceIamSchema, DataprocMetastoreServiceIamUpdaterProducer, DataprocMetastoreServiceIdParseFunc),
			"google_datastore_index":                                       resourceDatastoreIndex(),
			"google_datastream_connection_profile":                         resourceDatastreamConnectionProfile(),
			"google_datastream_private_connection":                         resourceDatastreamPrivateConnection(),
			"google_datastream_stream":                                     resourceDatastreamStream(),
			"google_deployment_manager_deployment":                         resourceDeploymentManagerDeployment(),
			"google_dialogflow_agent":                                      resourceDialogflowAgent(),
			"google_dialogflow_intent":                                     resourceDialogflowIntent(),
			"google_dialogflow_entity_type":                                resourceDialogflowEntityType(),
			"google_dialogflow_fulfillment":                                resourceDialogflowFulfillment(),
			"google_dialogflow_cx_agent":                                   resourceDialogflowCXAgent(),
			"google_dialogflow_cx_intent":                                  resourceDialogflowCXIntent(),
			"google_dialogflow_cx_flow":                                    resourceDialogflowCXFlow(),
			"google_dialogflow_cx_page":                                    resourceDialogflowCXPage(),
			"google_dialogflow_cx_entity_type":                             resourceDialogflowCXEntityType(),
			"google_dialogflow_cx_webhook":                                 resourceDialogflowCXWebhook(),
			"google_dns_managed_zone":                                      resourceDNSManagedZone(),
			"google_dns_managed_zone_iam_binding":                          ResourceIamBinding(DNSManagedZoneIamSchema, DNSManagedZoneIamUpdaterProducer, DNSManagedZoneIdParseFunc),
			"google_dns_managed_zone_iam_member":                           ResourceIamMember(DNSManagedZoneIamSchema, DNSManagedZoneIamUpdaterProducer, DNSManagedZoneIdParseFunc),
			"google_dns_managed_zone_iam_policy":                           ResourceIamPolicy(DNSManagedZoneIamSchema, DNSManagedZoneIamUpdaterProducer, DNSManagedZoneIdParseFunc),
			"google_dns_policy":                                            resourceDNSPolicy(),
			"google_document_ai_processor":                                 resourceDocumentAIProcessor(),
			"google_document_ai_processor_default_version":                 resourceDocumentAIProcessorDefaultVersion(),
			"google_essential_contacts_contact":                            resourceEssentialContactsContact(),
			"google_filestore_instance":                                    resourceFilestoreInstance(),
			"google_filestore_snapshot":                                    resourceFilestoreSnapshot(),
			"google_filestore_backup":                                      resourceFilestoreBackup(),
			"google_firestore_index":                                       resourceFirestoreIndex(),
			"google_firestore_document":                                    resourceFirestoreDocument(),
			"google_game_services_realm":                                   resourceGameServicesRealm(),
			"google_game_services_game_server_cluster":                     resourceGameServicesGameServerCluster(),
			"google_game_services_game_server_deployment":                  resourceGameServicesGameServerDeployment(),
			"google_game_services_game_server_config":                      resourceGameServicesGameServerConfig(),
			"google_game_services_game_server_deployment_rollout":          resourceGameServicesGameServerDeploymentRollout(),
			"google_gke_backup_backup_plan":                                resourceGKEBackupBackupPlan(),
			"google_gke_backup_backup_plan_iam_binding":                    ResourceIamBinding(GKEBackupBackupPlanIamSchema, GKEBackupBackupPlanIamUpdaterProducer, GKEBackupBackupPlanIdParseFunc),
			"google_gke_backup_backup_plan_iam_member":                     ResourceIamMember(GKEBackupBackupPlanIamSchema, GKEBackupBackupPlanIamUpdaterProducer, GKEBackupBackupPlanIdParseFunc),
			"google_gke_backup_backup_plan_iam_policy":                     ResourceIamPolicy(GKEBackupBackupPlanIamSchema, GKEBackupBackupPlanIamUpdaterProducer, GKEBackupBackupPlanIdParseFunc),
			"google_gke_hub_membership":                                    resourceGKEHubMembership(),
			"google_gke_hub_membership_iam_binding":                        ResourceIamBinding(GKEHubMembershipIamSchema, GKEHubMembershipIamUpdaterProducer, GKEHubMembershipIdParseFunc),
			"google_gke_hub_membership_iam_member":                         ResourceIamMember(GKEHubMembershipIamSchema, GKEHubMembershipIamUpdaterProducer, GKEHubMembershipIdParseFunc),
			"google_gke_hub_membership_iam_policy":                         ResourceIamPolicy(GKEHubMembershipIamSchema, GKEHubMembershipIamUpdaterProducer, GKEHubMembershipIdParseFunc),
			"google_healthcare_dataset":                                    resourceHealthcareDataset(),
			"google_healthcare_dicom_store":                                resourceHealthcareDicomStore(),
			"google_healthcare_fhir_store":                                 resourceHealthcareFhirStore(),
			"google_healthcare_hl7_v2_store":                               resourceHealthcareHl7V2Store(),
			"google_healthcare_consent_store":                              resourceHealthcareConsentStore(),
			"google_healthcare_consent_store_iam_binding":                  ResourceIamBinding(HealthcareConsentStoreIamSchema, HealthcareConsentStoreIamUpdaterProducer, HealthcareConsentStoreIdParseFunc),
			"google_healthcare_consent_store_iam_member":                   ResourceIamMember(HealthcareConsentStoreIamSchema, HealthcareConsentStoreIamUpdaterProducer, HealthcareConsentStoreIdParseFunc),
			"google_healthcare_consent_store_iam_policy":                   ResourceIamPolicy(HealthcareConsentStoreIamSchema, HealthcareConsentStoreIamUpdaterProducer, HealthcareConsentStoreIdParseFunc),
			"google_iam_access_boundary_policy":                            resourceIAM2AccessBoundaryPolicy(),
			"google_iam_workload_identity_pool":                            resourceIAMBetaWorkloadIdentityPool(),
			"google_iam_workload_identity_pool_provider":                   resourceIAMBetaWorkloadIdentityPoolProvider(),
			"google_iam_workforce_pool":                                    resourceIAMWorkforcePoolWorkforcePool(),
			"google_iam_workforce_pool_provider":                           resourceIAMWorkforcePoolWorkforcePoolProvider(),
			"google_iap_web_iam_binding":                                   ResourceIamBinding(IapWebIamSchema, IapWebIamUpdaterProducer, IapWebIdParseFunc),
			"google_iap_web_iam_member":                                    ResourceIamMember(IapWebIamSchema, IapWebIamUpdaterProducer, IapWebIdParseFunc),
			"google_iap_web_iam_policy":                                    ResourceIamPolicy(IapWebIamSchema, IapWebIamUpdaterProducer, IapWebIdParseFunc),
			"google_iap_web_type_compute_iam_binding":                      ResourceIamBinding(IapWebTypeComputeIamSchema, IapWebTypeComputeIamUpdaterProducer, IapWebTypeComputeIdParseFunc),
			"google_iap_web_type_compute_iam_member":                       ResourceIamMember(IapWebTypeComputeIamSchema, IapWebTypeComputeIamUpdaterProducer, IapWebTypeComputeIdParseFunc),
			"google_iap_web_type_compute_iam_policy":                       ResourceIamPolicy(IapWebTypeComputeIamSchema, IapWebTypeComputeIamUpdaterProducer, IapWebTypeComputeIdParseFunc),
			"google_iap_web_type_app_engine_iam_binding":                   ResourceIamBinding(IapWebTypeAppEngineIamSchema, IapWebTypeAppEngineIamUpdaterProducer, IapWebTypeAppEngineIdParseFunc),
			"google_iap_web_type_app_engine_iam_member":                    ResourceIamMember(IapWebTypeAppEngineIamSchema, IapWebTypeAppEngineIamUpdaterProducer, IapWebTypeAppEngineIdParseFunc),
			"google_iap_web_type_app_engine_iam_policy":                    ResourceIamPolicy(IapWebTypeAppEngineIamSchema, IapWebTypeAppEngineIamUpdaterProducer, IapWebTypeAppEngineIdParseFunc),
			"google_iap_app_engine_version_iam_binding":                    ResourceIamBinding(IapAppEngineVersionIamSchema, IapAppEngineVersionIamUpdaterProducer, IapAppEngineVersionIdParseFunc),
			"google_iap_app_engine_version_iam_member":                     ResourceIamMember(IapAppEngineVersionIamSchema, IapAppEngineVersionIamUpdaterProducer, IapAppEngineVersionIdParseFunc),
			"google_iap_app_engine_version_iam_policy":                     ResourceIamPolicy(IapAppEngineVersionIamSchema, IapAppEngineVersionIamUpdaterProducer, IapAppEngineVersionIdParseFunc),
			"google_iap_app_engine_service_iam_binding":                    ResourceIamBinding(IapAppEngineServiceIamSchema, IapAppEngineServiceIamUpdaterProducer, IapAppEngineServiceIdParseFunc),
			"google_iap_app_engine_service_iam_member":                     ResourceIamMember(IapAppEngineServiceIamSchema, IapAppEngineServiceIamUpdaterProducer, IapAppEngineServiceIdParseFunc),
			"google_iap_app_engine_service_iam_policy":                     ResourceIamPolicy(IapAppEngineServiceIamSchema, IapAppEngineServiceIamUpdaterProducer, IapAppEngineServiceIdParseFunc),
			"google_iap_web_backend_service_iam_binding":                   ResourceIamBinding(IapWebBackendServiceIamSchema, IapWebBackendServiceIamUpdaterProducer, IapWebBackendServiceIdParseFunc),
			"google_iap_web_backend_service_iam_member":                    ResourceIamMember(IapWebBackendServiceIamSchema, IapWebBackendServiceIamUpdaterProducer, IapWebBackendServiceIdParseFunc),
			"google_iap_web_backend_service_iam_policy":                    ResourceIamPolicy(IapWebBackendServiceIamSchema, IapWebBackendServiceIamUpdaterProducer, IapWebBackendServiceIdParseFunc),
			"google_iap_tunnel_instance_iam_binding":                       ResourceIamBinding(IapTunnelInstanceIamSchema, IapTunnelInstanceIamUpdaterProducer, IapTunnelInstanceIdParseFunc),
			"google_iap_tunnel_instance_iam_member":                        ResourceIamMember(IapTunnelInstanceIamSchema, IapTunnelInstanceIamUpdaterProducer, IapTunnelInstanceIdParseFunc),
			"google_iap_tunnel_instance_iam_policy":                        ResourceIamPolicy(IapTunnelInstanceIamSchema, IapTunnelInstanceIamUpdaterProducer, IapTunnelInstanceIdParseFunc),
			"google_iap_tunnel_iam_binding":                                ResourceIamBinding(IapTunnelIamSchema, IapTunnelIamUpdaterProducer, IapTunnelIdParseFunc),
			"google_iap_tunnel_iam_member":                                 ResourceIamMember(IapTunnelIamSchema, IapTunnelIamUpdaterProducer, IapTunnelIdParseFunc),
			"google_iap_tunnel_iam_policy":                                 ResourceIamPolicy(IapTunnelIamSchema, IapTunnelIamUpdaterProducer, IapTunnelIdParseFunc),
			"google_iap_brand":                                             resourceIapBrand(),
			"google_iap_client":                                            resourceIapClient(),
			"google_identity_platform_config":                              resourceIdentityPlatformConfig(),
			"google_identity_platform_default_supported_idp_config":        resourceIdentityPlatformDefaultSupportedIdpConfig(),
			"google_identity_platform_tenant_default_supported_idp_config": resourceIdentityPlatformTenantDefaultSupportedIdpConfig(),
			"google_identity_platform_inbound_saml_config":                 resourceIdentityPlatformInboundSamlConfig(),
			"google_identity_platform_tenant_inbound_saml_config":          resourceIdentityPlatformTenantInboundSamlConfig(),
			"google_identity_platform_oauth_idp_config":                    resourceIdentityPlatformOauthIdpConfig(),
			"google_identity_platform_tenant_oauth_idp_config":             resourceIdentityPlatformTenantOauthIdpConfig(),
			"google_identity_platform_tenant":                              resourceIdentityPlatformTenant(),
			"google_identity_platform_project_default_config":              resourceIdentityPlatformProjectDefaultConfig(),
			"google_kms_key_ring":                                          resourceKMSKeyRing(),
			"google_kms_crypto_key":                                        resourceKMSCryptoKey(),
			"google_kms_crypto_key_version":                                resourceKMSCryptoKeyVersion(),
			"google_kms_key_ring_import_job":                               resourceKMSKeyRingImportJob(),
			"google_kms_secret_ciphertext":                                 resourceKMSSecretCiphertext(),
			"google_logging_metric":                                        resourceLoggingMetric(),
			"google_memcache_instance":                                     resourceMemcacheInstance(),
			"google_ml_engine_model":                                       resourceMLEngineModel(),
			"google_monitoring_alert_policy":                               resourceMonitoringAlertPolicy(),
			"google_monitoring_group":                                      resourceMonitoringGroup(),
			"google_monitoring_notification_channel":                       resourceMonitoringNotificationChannel(),
			"google_monitoring_custom_service":                             resourceMonitoringService(),
			"google_monitoring_service":                                    resourceMonitoringGenericService(),
			"google_monitoring_slo":                                        resourceMonitoringSlo(),
			"google_monitoring_uptime_check_config":                        resourceMonitoringUptimeCheckConfig(),
			"google_monitoring_metric_descriptor":                          resourceMonitoringMetricDescriptor(),
			"google_network_management_connectivity_test":                  resourceNetworkManagementConnectivityTest(),
			"google_network_services_edge_cache_keyset":                    resourceNetworkServicesEdgeCacheKeyset(),
			"google_network_services_edge_cache_origin":                    resourceNetworkServicesEdgeCacheOrigin(),
			"google_network_services_edge_cache_service":                   resourceNetworkServicesEdgeCacheService(),
			"google_notebooks_environment":                                 resourceNotebooksEnvironment(),
			"google_notebooks_instance":                                    resourceNotebooksInstance(),
			"google_notebooks_instance_iam_binding":                        ResourceIamBinding(NotebooksInstanceIamSchema, NotebooksInstanceIamUpdaterProducer, NotebooksInstanceIdParseFunc),
			"google_notebooks_instance_iam_member":                         ResourceIamMember(NotebooksInstanceIamSchema, NotebooksInstanceIamUpdaterProducer, NotebooksInstanceIdParseFunc),
			"google_notebooks_instance_iam_policy":                         ResourceIamPolicy(NotebooksInstanceIamSchema, NotebooksInstanceIamUpdaterProducer, NotebooksInstanceIdParseFunc),
			"google_notebooks_runtime":                                     resourceNotebooksRuntime(),
			"google_notebooks_runtime_iam_binding":                         ResourceIamBinding(NotebooksRuntimeIamSchema, NotebooksRuntimeIamUpdaterProducer, NotebooksRuntimeIdParseFunc),
			"google_notebooks_runtime_iam_member":                          ResourceIamMember(NotebooksRuntimeIamSchema, NotebooksRuntimeIamUpdaterProducer, NotebooksRuntimeIdParseFunc),
			"google_notebooks_runtime_iam_policy":                          ResourceIamPolicy(NotebooksRuntimeIamSchema, NotebooksRuntimeIamUpdaterProducer, NotebooksRuntimeIdParseFunc),
			"google_notebooks_location":                                    resourceNotebooksLocation(),
			"google_os_config_patch_deployment":                            resourceOSConfigPatchDeployment(),
			"google_os_login_ssh_public_key":                               resourceOSLoginSSHPublicKey(),
			"google_privateca_certificate_authority":                       resourcePrivatecaCertificateAuthority(),
			"google_privateca_certificate":                                 resourcePrivatecaCertificate(),
			"google_privateca_ca_pool":                                     resourcePrivatecaCaPool(),
			"google_privateca_ca_pool_iam_binding":                         ResourceIamBinding(PrivatecaCaPoolIamSchema, PrivatecaCaPoolIamUpdaterProducer, PrivatecaCaPoolIdParseFunc),
			"google_privateca_ca_pool_iam_member":                          ResourceIamMember(PrivatecaCaPoolIamSchema, PrivatecaCaPoolIamUpdaterProducer, PrivatecaCaPoolIdParseFunc),
			"google_privateca_ca_pool_iam_policy":                          ResourceIamPolicy(PrivatecaCaPoolIamSchema, PrivatecaCaPoolIamUpdaterProducer, PrivatecaCaPoolIdParseFunc),
			"google_privateca_certificate_template_iam_binding":            ResourceIamBinding(PrivatecaCertificateTemplateIamSchema, PrivatecaCertificateTemplateIamUpdaterProducer, PrivatecaCertificateTemplateIdParseFunc),
			"google_privateca_certificate_template_iam_member":             ResourceIamMember(PrivatecaCertificateTemplateIamSchema, PrivatecaCertificateTemplateIamUpdaterProducer, PrivatecaCertificateTemplateIdParseFunc),
			"google_privateca_certificate_template_iam_policy":             ResourceIamPolicy(PrivatecaCertificateTemplateIamSchema, PrivatecaCertificateTemplateIamUpdaterProducer, PrivatecaCertificateTemplateIdParseFunc),
			"google_pubsub_topic":                                          resourcePubsubTopic(),
			"google_pubsub_topic_iam_binding":                              ResourceIamBinding(PubsubTopicIamSchema, PubsubTopicIamUpdaterProducer, PubsubTopicIdParseFunc),
			"google_pubsub_topic_iam_member":                               ResourceIamMember(PubsubTopicIamSchema, PubsubTopicIamUpdaterProducer, PubsubTopicIdParseFunc),
			"google_pubsub_topic_iam_policy":                               ResourceIamPolicy(PubsubTopicIamSchema, PubsubTopicIamUpdaterProducer, PubsubTopicIdParseFunc),
			"google_pubsub_subscription":                                   resourcePubsubSubscription(),
			"google_pubsub_schema":                                         resourcePubsubSchema(),
			"google_pubsub_lite_reservation":                               resourcePubsubLiteReservation(),
			"google_pubsub_lite_topic":                                     resourcePubsubLiteTopic(),
			"google_pubsub_lite_subscription":                              resourcePubsubLiteSubscription(),
			"google_redis_instance":                                        resourceRedisInstance(),
			"google_resource_manager_lien":                                 resourceResourceManagerLien(),
			"google_secret_manager_secret":                                 resourceSecretManagerSecret(),
			"google_secret_manager_secret_iam_binding":                     ResourceIamBinding(SecretManagerSecretIamSchema, SecretManagerSecretIamUpdaterProducer, SecretManagerSecretIdParseFunc),
			"google_secret_manager_secret_iam_member":                      ResourceIamMember(SecretManagerSecretIamSchema, SecretManagerSecretIamUpdaterProducer, SecretManagerSecretIdParseFunc),
			"google_secret_manager_secret_iam_policy":                      ResourceIamPolicy(SecretManagerSecretIamSchema, SecretManagerSecretIamUpdaterProducer, SecretManagerSecretIdParseFunc),
			"google_secret_manager_secret_version":                         resourceSecretManagerSecretVersion(),
			"google_scc_source":                                            resourceSecurityCenterSource(),
			"google_scc_source_iam_binding":                                ResourceIamBinding(SecurityCenterSourceIamSchema, SecurityCenterSourceIamUpdaterProducer, SecurityCenterSourceIdParseFunc),
			"google_scc_source_iam_member":                                 ResourceIamMember(SecurityCenterSourceIamSchema, SecurityCenterSourceIamUpdaterProducer, SecurityCenterSourceIdParseFunc),
			"google_scc_source_iam_policy":                                 ResourceIamPolicy(SecurityCenterSourceIamSchema, SecurityCenterSourceIamUpdaterProducer, SecurityCenterSourceIdParseFunc),
			"google_scc_notification_config":                               resourceSecurityCenterNotificationConfig(),
			"google_endpoints_service_iam_binding":                         ResourceIamBinding(ServiceManagementServiceIamSchema, ServiceManagementServiceIamUpdaterProducer, ServiceManagementServiceIdParseFunc),
			"google_endpoints_service_iam_member":                          ResourceIamMember(ServiceManagementServiceIamSchema, ServiceManagementServiceIamUpdaterProducer, ServiceManagementServiceIdParseFunc),
			"google_endpoints_service_iam_policy":                          ResourceIamPolicy(ServiceManagementServiceIamSchema, ServiceManagementServiceIamUpdaterProducer, ServiceManagementServiceIdParseFunc),
			"google_endpoints_service_consumers_iam_binding":               ResourceIamBinding(ServiceManagementServiceConsumersIamSchema, ServiceManagementServiceConsumersIamUpdaterProducer, ServiceManagementServiceConsumersIdParseFunc),
			"google_endpoints_service_consumers_iam_member":                ResourceIamMember(ServiceManagementServiceConsumersIamSchema, ServiceManagementServiceConsumersIamUpdaterProducer, ServiceManagementServiceConsumersIdParseFunc),
			"google_endpoints_service_consumers_iam_policy":                ResourceIamPolicy(ServiceManagementServiceConsumersIamSchema, ServiceManagementServiceConsumersIamUpdaterProducer, ServiceManagementServiceConsumersIdParseFunc),
			"google_sourcerepo_repository":                                 resourceSourceRepoRepository(),
			"google_sourcerepo_repository_iam_binding":                     ResourceIamBinding(SourceRepoRepositoryIamSchema, SourceRepoRepositoryIamUpdaterProducer, SourceRepoRepositoryIdParseFunc),
			"google_sourcerepo_repository_iam_member":                      ResourceIamMember(SourceRepoRepositoryIamSchema, SourceRepoRepositoryIamUpdaterProducer, SourceRepoRepositoryIdParseFunc),
			"google_sourcerepo_repository_iam_policy":                      ResourceIamPolicy(SourceRepoRepositoryIamSchema, SourceRepoRepositoryIamUpdaterProducer, SourceRepoRepositoryIdParseFunc),
			"google_spanner_instance":                                      resourceSpannerInstance(),
			"google_spanner_database":                                      resourceSpannerDatabase(),
			"google_sql_database":                                          resourceSQLDatabase(),
			"google_sql_source_representation_instance":                    resourceSQLSourceRepresentationInstance(),
			"google_storage_bucket_iam_binding":                            ResourceIamBinding(StorageBucketIamSchema, StorageBucketIamUpdaterProducer, StorageBucketIdParseFunc),
			"google_storage_bucket_iam_member":                             ResourceIamMember(StorageBucketIamSchema, StorageBucketIamUpdaterProducer, StorageBucketIdParseFunc),
			"google_storage_bucket_iam_policy":                             ResourceIamPolicy(StorageBucketIamSchema, StorageBucketIamUpdaterProducer, StorageBucketIdParseFunc),
			"google_storage_bucket_access_control":                         resourceStorageBucketAccessControl(),
			"google_storage_object_access_control":                         resourceStorageObjectAccessControl(),
			"google_storage_default_object_access_control":                 resourceStorageDefaultObjectAccessControl(),
			"google_storage_hmac_key":                                      resourceStorageHmacKey(),
			"google_storage_transfer_agent_pool":                           resourceStorageTransferAgentPool(),
			"google_tags_tag_key":                                          resourceTagsTagKey(),
			"google_tags_tag_key_iam_binding":                              ResourceIamBinding(TagsTagKeyIamSchema, TagsTagKeyIamUpdaterProducer, TagsTagKeyIdParseFunc),
			"google_tags_tag_key_iam_member":                               ResourceIamMember(TagsTagKeyIamSchema, TagsTagKeyIamUpdaterProducer, TagsTagKeyIdParseFunc),
			"google_tags_tag_key_iam_policy":                               ResourceIamPolicy(TagsTagKeyIamSchema, TagsTagKeyIamUpdaterProducer, TagsTagKeyIdParseFunc),
			"google_tags_tag_value":                                        resourceTagsTagValue(),
			"google_tags_tag_value_iam_binding":                            ResourceIamBinding(TagsTagValueIamSchema, TagsTagValueIamUpdaterProducer, TagsTagValueIdParseFunc),
			"google_tags_tag_value_iam_member":                             ResourceIamMember(TagsTagValueIamSchema, TagsTagValueIamUpdaterProducer, TagsTagValueIdParseFunc),
			"google_tags_tag_value_iam_policy":                             ResourceIamPolicy(TagsTagValueIamSchema, TagsTagValueIamUpdaterProducer, TagsTagValueIdParseFunc),
			"google_tags_tag_binding":                                      resourceTagsTagBinding(),
			"google_tpu_node":                                              resourceTPUNode(),
			"google_vertex_ai_tensorboard":                                 resourceVertexAITensorboard(),
			"google_vertex_ai_dataset":                                     resourceVertexAIDataset(),
			"google_vertex_ai_endpoint":                                    resourceVertexAIEndpoint(),
			"google_vertex_ai_featurestore":                                resourceVertexAIFeaturestore(),
			"google_vertex_ai_featurestore_entitytype":                     resourceVertexAIFeaturestoreEntitytype(),
			"google_vertex_ai_featurestore_entitytype_feature":             resourceVertexAIFeaturestoreEntitytypeFeature(),
			"google_vertex_ai_index":                                       resourceVertexAIIndex(),
			"google_vpc_access_connector":                                  resourceVPCAccessConnector(),
			"google_workflows_workflow":                                    resourceWorkflowsWorkflow(),
		},
		map[string]*schema.Resource{
			// ####### START handwritten resources ###########
			"google_app_engine_application":                resourceAppEngineApplication(),
			"google_bigquery_table":                        resourceBigQueryTable(),
			"google_bigtable_gc_policy":                    resourceBigtableGCPolicy(),
			"google_bigtable_instance":                     resourceBigtableInstance(),
			"google_bigtable_table":                        resourceBigtableTable(),
			"google_billing_subaccount":                    resourceBillingSubaccount(),
			"google_cloudfunctions_function":               resourceCloudFunctionsFunction(),
			"google_composer_environment":                  resourceComposerEnvironment(),
			"google_compute_attached_disk":                 resourceComputeAttachedDisk(),
			"google_compute_instance":                      resourceComputeInstance(),
			"google_compute_instance_from_template":        resourceComputeInstanceFromTemplate(),
			"google_compute_instance_group":                resourceComputeInstanceGroup(),
			"google_compute_instance_group_manager":        resourceComputeInstanceGroupManager(),
			"google_compute_instance_template":             resourceComputeInstanceTemplate(),
			"google_compute_network_peering":               resourceComputeNetworkPeering(),
			"google_compute_project_default_network_tier":  resourceComputeProjectDefaultNetworkTier(),
			"google_compute_project_metadata":              resourceComputeProjectMetadata(),
			"google_compute_project_metadata_item":         resourceComputeProjectMetadataItem(),
			"google_compute_region_instance_group_manager": resourceComputeRegionInstanceGroupManager(),
			"google_compute_router_interface":              resourceComputeRouterInterface(),
			"google_compute_security_policy":               resourceComputeSecurityPolicy(),
			"google_compute_shared_vpc_host_project":       resourceComputeSharedVpcHostProject(),
			"google_compute_shared_vpc_service_project":    resourceComputeSharedVpcServiceProject(),
			"google_compute_target_pool":                   resourceComputeTargetPool(),
			"google_container_cluster":                     resourceContainerCluster(),
			"google_container_node_pool":                   resourceContainerNodePool(),
			"google_container_registry":                    resourceContainerRegistry(),
			"google_dataflow_job":                          resourceDataflowJob(),
			"google_dataproc_cluster":                      resourceDataprocCluster(),
			"google_dataproc_job":                          resourceDataprocJob(),
			"google_dialogflow_cx_version":                 resourceDialogflowCXVersion(),
			"google_dialogflow_cx_environment":             resourceDialogflowCXEnvironment(),
			"google_dns_record_set":                        resourceDnsRecordSet(),
			"google_endpoints_service":                     resourceEndpointsService(),
			"google_folder":                                resourceGoogleFolder(),
			"google_folder_organization_policy":            resourceGoogleFolderOrganizationPolicy(),
			"google_logging_billing_account_sink":          resourceLoggingBillingAccountSink(),
			"google_logging_billing_account_exclusion":     ResourceLoggingExclusion(BillingAccountLoggingExclusionSchema, NewBillingAccountLoggingExclusionUpdater, billingAccountLoggingExclusionIdParseFunc),
			"google_logging_billing_account_bucket_config": ResourceLoggingBillingAccountBucketConfig(),
			"google_logging_organization_sink":             resourceLoggingOrganizationSink(),
			"google_logging_organization_exclusion":        ResourceLoggingExclusion(OrganizationLoggingExclusionSchema, NewOrganizationLoggingExclusionUpdater, organizationLoggingExclusionIdParseFunc),
			"google_logging_organization_bucket_config":    ResourceLoggingOrganizationBucketConfig(),
			"google_logging_folder_sink":                   resourceLoggingFolderSink(),
			"google_logging_folder_exclusion":              ResourceLoggingExclusion(FolderLoggingExclusionSchema, NewFolderLoggingExclusionUpdater, folderLoggingExclusionIdParseFunc),
			"google_logging_folder_bucket_config":          ResourceLoggingFolderBucketConfig(),
			"google_logging_project_sink":                  resourceLoggingProjectSink(),
			"google_logging_project_exclusion":             ResourceLoggingExclusion(ProjectLoggingExclusionSchema, NewProjectLoggingExclusionUpdater, projectLoggingExclusionIdParseFunc),
			"google_logging_project_bucket_config":         ResourceLoggingProjectBucketConfig(),
			"google_monitoring_dashboard":                  resourceMonitoringDashboard(),
			"google_service_networking_connection":         resourceServiceNetworkingConnection(),
			"google_sql_database_instance":                 resourceSqlDatabaseInstance(),
			"google_sql_ssl_cert":                          resourceSqlSslCert(),
			"google_sql_user":                              resourceSqlUser(),
			"google_organization_iam_custom_role":          resourceGoogleOrganizationIamCustomRole(),
			"google_organization_policy":                   resourceGoogleOrganizationPolicy(),
			"google_project":                               resourceGoogleProject(),
			"google_project_default_service_accounts":      resourceGoogleProjectDefaultServiceAccounts(),
			"google_project_service":                       resourceGoogleProjectService(),
			"google_project_iam_custom_role":               resourceGoogleProjectIamCustomRole(),
			"google_project_organization_policy":           resourceGoogleProjectOrganizationPolicy(),
			"google_project_usage_export_bucket":           resourceProjectUsageBucket(),
			"google_service_account":                       resourceGoogleServiceAccount(),
			"google_service_account_key":                   resourceGoogleServiceAccountKey(),
			"google_service_networking_peered_dns_domain":  resourceGoogleServiceNetworkingPeeredDNSDomain(),
			"google_storage_bucket":                        resourceStorageBucket(),
			"google_storage_bucket_acl":                    resourceStorageBucketAcl(),
			"google_storage_bucket_object":                 resourceStorageBucketObject(),
			"google_storage_object_acl":                    resourceStorageObjectAcl(),
			"google_storage_default_object_acl":            resourceStorageDefaultObjectAcl(),
			"google_storage_notification":                  resourceStorageNotification(),
			"google_storage_transfer_job":                  resourceStorageTransferJob(),
			"google_tags_location_tag_binding":             resourceTagsLocationTagBinding(),
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
	HandleSDKDefaults(d)
	HandleDCLProviderDefaults(d)

	config := Config{
		Project:             d.Get("project").(string),
		Region:              d.Get("region").(string),
		Zone:                d.Get("zone").(string),
		UserProjectOverride: d.Get("user_project_override").(bool),
		BillingProject:      d.Get("billing_project").(string),
		userAgent:           p.UserAgent("terraform-provider-google", version.ProviderVersion),
	}

	// opt in extension for adding to the User-Agent header
	if ext := os.Getenv("GOOGLE_TERRAFORM_USERAGENT_EXTENSION"); ext != "" {
		ua := config.userAgent
		config.userAgent = fmt.Sprintf("%s %s", ua, ext)
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
		config.Credentials = multiEnvSearch([]string{
			"GOOGLE_CREDENTIALS",
			"GOOGLE_CLOUD_KEYFILE_JSON",
			"GCLOUD_KEYFILE_JSON",
		})

		config.AccessToken = multiEnvSearch([]string{
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

	batchCfg, err := expandProviderBatchingConfig(d.Get("batching"))
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

	return providerDCLConfigure(d, &config), nil
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
