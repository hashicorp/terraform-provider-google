package google

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform/helper/mutexkv"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"

	googleoauth "golang.org/x/oauth2/google"
)

// Global MutexKV
var mutexKV = mutexkv.NewMutexKV()

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"credentials": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_CREDENTIALS",
					"GOOGLE_CLOUD_KEYFILE_JSON",
					"GCLOUD_KEYFILE_JSON",
				}, nil),
				ValidateFunc: validateCredentials,
			},

			"access_token": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_OAUTH_ACCESS_TOKEN",
				}, nil),
				ConflictsWith: []string{"credentials"},
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_PROJECT",
					"GOOGLE_CLOUD_PROJECT",
					"GCLOUD_PROJECT",
					"CLOUDSDK_CORE_PROJECT",
				}, nil),
			},

			"region": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_REGION",
					"GCLOUD_REGION",
					"CLOUDSDK_COMPUTE_REGION",
				}, nil),
			},

			"zone": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_ZONE",
					"GCLOUD_ZONE",
					"CLOUDSDK_COMPUTE_ZONE",
				}, nil),
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
							Default:      "10s",
							ValidateFunc: validateNonNegativeDuration(),
						},
						"enable_batching": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},

			"user_project_override": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			// Generated Products
			AccessContextManagerCustomEndpointEntryKey: AccessContextManagerCustomEndpointEntry,
			AppEngineCustomEndpointEntryKey:            AppEngineCustomEndpointEntry,
			BigqueryDataTransferCustomEndpointEntryKey: BigqueryDataTransferCustomEndpointEntry,
			BinaryAuthorizationCustomEndpointEntryKey:  BinaryAuthorizationCustomEndpointEntry,
			CloudBuildCustomEndpointEntryKey:           CloudBuildCustomEndpointEntry,
			CloudSchedulerCustomEndpointEntryKey:       CloudSchedulerCustomEndpointEntry,
			ComputeCustomEndpointEntryKey:              ComputeCustomEndpointEntry,
			DnsCustomEndpointEntryKey:                  DnsCustomEndpointEntry,
			FilestoreCustomEndpointEntryKey:            FilestoreCustomEndpointEntry,
			FirestoreCustomEndpointEntryKey:            FirestoreCustomEndpointEntry,
			KmsCustomEndpointEntryKey:                  KmsCustomEndpointEntry,
			LoggingCustomEndpointEntryKey:              LoggingCustomEndpointEntry,
			MonitoringCustomEndpointEntryKey:           MonitoringCustomEndpointEntry,
			PubsubCustomEndpointEntryKey:               PubsubCustomEndpointEntry,
			RedisCustomEndpointEntryKey:                RedisCustomEndpointEntry,
			ResourceManagerCustomEndpointEntryKey:      ResourceManagerCustomEndpointEntry,
			SecurityCenterCustomEndpointEntryKey:       SecurityCenterCustomEndpointEntry,
			SourceRepoCustomEndpointEntryKey:           SourceRepoCustomEndpointEntry,
			SpannerCustomEndpointEntryKey:              SpannerCustomEndpointEntry,
			SqlCustomEndpointEntryKey:                  SqlCustomEndpointEntry,
			StorageCustomEndpointEntryKey:              StorageCustomEndpointEntry,
			TpuCustomEndpointEntryKey:                  TpuCustomEndpointEntry,

			// Handwritten Products / Versioned / Atypical Entries
			CloudBillingCustomEndpointEntryKey:           CloudBillingCustomEndpointEntry,
			ComposerCustomEndpointEntryKey:               ComposerCustomEndpointEntry,
			ComputeBetaCustomEndpointEntryKey:            ComputeBetaCustomEndpointEntry,
			ContainerCustomEndpointEntryKey:              ContainerCustomEndpointEntry,
			ContainerBetaCustomEndpointEntryKey:          ContainerBetaCustomEndpointEntry,
			DataprocCustomEndpointEntryKey:               DataprocCustomEndpointEntry,
			DataprocBetaCustomEndpointEntryKey:           DataprocBetaCustomEndpointEntry,
			DataflowCustomEndpointEntryKey:               DataflowCustomEndpointEntry,
			DnsBetaCustomEndpointEntryKey:                DnsBetaCustomEndpointEntry,
			IamCredentialsCustomEndpointEntryKey:         IamCredentialsCustomEndpointEntry,
			ResourceManagerV2Beta1CustomEndpointEntryKey: ResourceManagerV2Beta1CustomEndpointEntry,
			RuntimeconfigCustomEndpointEntryKey:          RuntimeconfigCustomEndpointEntry,
			IAMCustomEndpointEntryKey:                    IAMCustomEndpointEntry,
			ServiceManagementCustomEndpointEntryKey:      ServiceManagementCustomEndpointEntry,
			ServiceNetworkingCustomEndpointEntryKey:      ServiceNetworkingCustomEndpointEntry,
			ServiceUsageCustomEndpointEntryKey:           ServiceUsageCustomEndpointEntry,
			BigQueryCustomEndpointEntryKey:               BigQueryCustomEndpointEntry,
			CloudFunctionsCustomEndpointEntryKey:         CloudFunctionsCustomEndpointEntry,
			CloudIoTCustomEndpointEntryKey:               CloudIoTCustomEndpointEntry,
			StorageTransferCustomEndpointEntryKey:        StorageTransferCustomEndpointEntry,
			BigtableAdminCustomEndpointEntryKey:          BigtableAdminCustomEndpointEntry,
		},

		DataSourcesMap: map[string]*schema.Resource{
			"google_active_folder":                            dataSourceGoogleActiveFolder(),
			"google_billing_account":                          dataSourceGoogleBillingAccount(),
			"google_dns_managed_zone":                         dataSourceDnsManagedZone(),
			"google_client_config":                            dataSourceGoogleClientConfig(),
			"google_client_openid_userinfo":                   dataSourceGoogleClientOpenIDUserinfo(),
			"google_cloudfunctions_function":                  dataSourceGoogleCloudFunctionsFunction(),
			"google_composer_image_versions":                  dataSourceGoogleComposerImageVersions(),
			"google_compute_address":                          dataSourceGoogleComputeAddress(),
			"google_compute_backend_service":                  dataSourceGoogleComputeBackendService(),
			"google_compute_default_service_account":          dataSourceGoogleComputeDefaultServiceAccount(),
			"google_compute_forwarding_rule":                  dataSourceGoogleComputeForwardingRule(),
			"google_compute_image":                            dataSourceGoogleComputeImage(),
			"google_compute_instance":                         dataSourceGoogleComputeInstance(),
			"google_compute_global_address":                   dataSourceGoogleComputeGlobalAddress(),
			"google_compute_instance_group":                   dataSourceGoogleComputeInstanceGroup(),
			"google_compute_lb_ip_ranges":                     dataSourceGoogleComputeLbIpRanges(),
			"google_compute_network":                          dataSourceGoogleComputeNetwork(),
			"google_compute_network_endpoint_group":           dataSourceGoogleComputeNetworkEndpointGroup(),
			"google_compute_node_types":                       dataSourceGoogleComputeNodeTypes(),
			"google_compute_regions":                          dataSourceGoogleComputeRegions(),
			"google_compute_region_instance_group":            dataSourceGoogleComputeRegionInstanceGroup(),
			"google_compute_subnetwork":                       dataSourceGoogleComputeSubnetwork(),
			"google_compute_zones":                            dataSourceGoogleComputeZones(),
			"google_compute_vpn_gateway":                      dataSourceGoogleComputeVpnGateway(),
			"google_compute_ssl_policy":                       dataSourceGoogleComputeSslPolicy(),
			"google_compute_ssl_certificate":                  dataSourceGoogleComputeSslCertificate(),
			"google_container_cluster":                        dataSourceGoogleContainerCluster(),
			"google_container_engine_versions":                dataSourceGoogleContainerEngineVersions(),
			"google_container_registry_repository":            dataSourceGoogleContainerRepo(),
			"google_container_registry_image":                 dataSourceGoogleContainerImage(),
			"google_iam_policy":                               dataSourceGoogleIamPolicy(),
			"google_iam_role":                                 dataSourceGoogleIamRole(),
			"google_kms_secret":                               dataSourceGoogleKmsSecret(),
			"google_kms_key_ring":                             dataSourceGoogleKmsKeyRing(),
			"google_kms_crypto_key":                           dataSourceGoogleKmsCryptoKey(),
			"google_kms_crypto_key_version":                   dataSourceGoogleKmsCryptoKeyVersion(),
			"google_folder":                                   dataSourceGoogleFolder(),
			"google_folder_organization_policy":               dataSourceGoogleFolderOrganizationPolicy(),
			"google_netblock_ip_ranges":                       dataSourceGoogleNetblockIpRanges(),
			"google_organization":                             dataSourceGoogleOrganization(),
			"google_project":                                  dataSourceGoogleProject(),
			"google_projects":                                 dataSourceGoogleProjects(),
			"google_project_organization_policy":              dataSourceGoogleProjectOrganizationPolicy(),
			"google_project_services":                         dataSourceGoogleProjectServices(),
			"google_service_account":                          dataSourceGoogleServiceAccount(),
			"google_service_account_access_token":             dataSourceGoogleServiceAccountAccessToken(),
			"google_service_account_key":                      dataSourceGoogleServiceAccountKey(),
			"google_storage_bucket_object":                    dataSourceGoogleStorageBucketObject(),
			"google_storage_object_signed_url":                dataSourceGoogleSignedUrl(),
			"google_storage_project_service_account":          dataSourceGoogleStorageProjectServiceAccount(),
			"google_storage_transfer_project_service_account": dataSourceGoogleStorageTransferProjectServiceAccount(),
			"google_tpu_tensorflow_versions":                  dataSourceTpuTensorflowVersions(),
		},

		ResourcesMap: ResourceMap(),

		ConfigureFunc: providerConfigure,
	}
}

func ResourceMap() map[string]*schema.Resource {
	resourceMap, _ := ResourceMapWithErrors()
	return resourceMap
}

func ResourceMapWithErrors() (map[string]*schema.Resource, error) {
	return mergeResourceMaps(
		GeneratedAccessContextManagerResourcesMap,
		GeneratedAppEngineResourcesMap,
		GeneratedBigqueryDataTransferResourcesMap,
		GeneratedBinaryAuthorizationResourcesMap,
		GeneratedCloudBuildResourcesMap,
		GeneratedCloudSchedulerResourcesMap,
		GeneratedComputeResourcesMap,
		GeneratedDnsResourcesMap,
		GeneratedFilestoreResourcesMap,
		GeneratedFirestoreResourcesMap,
		GeneratedKmsResourcesMap,
		GeneratedLoggingResourcesMap,
		GeneratedMonitoringResourcesMap,
		GeneratedPubsubResourcesMap,
		GeneratedRedisResourcesMap,
		GeneratedResourceManagerResourcesMap,
		GeneratedSecurityCenterResourcesMap,
		GeneratedSourceRepoResourcesMap,
		GeneratedSpannerResourcesMap,
		GeneratedSqlResourcesMap,
		GeneratedStorageResourcesMap,
		GeneratedTpuResourcesMap,
		map[string]*schema.Resource{
			"google_app_engine_application":                resourceAppEngineApplication(),
			"google_bigquery_dataset":                      resourceBigQueryDataset(),
			"google_bigquery_table":                        resourceBigQueryTable(),
			"google_bigtable_instance":                     resourceBigtableInstance(),
			"google_bigtable_instance_iam_binding":         ResourceIamBinding(IamBigtableInstanceSchema, NewBigtableInstanceUpdater, BigtableInstanceIdParseFunc),
			"google_bigtable_instance_iam_member":          ResourceIamMember(IamBigtableInstanceSchema, NewBigtableInstanceUpdater, BigtableInstanceIdParseFunc),
			"google_bigtable_instance_iam_policy":          ResourceIamPolicy(IamBigtableInstanceSchema, NewBigtableInstanceUpdater, BigtableInstanceIdParseFunc),
			"google_bigtable_table":                        resourceBigtableTable(),
			"google_billing_account_iam_binding":           ResourceIamBinding(IamBillingAccountSchema, NewBillingAccountIamUpdater, BillingAccountIdParseFunc),
			"google_billing_account_iam_member":            ResourceIamMember(IamBillingAccountSchema, NewBillingAccountIamUpdater, BillingAccountIdParseFunc),
			"google_billing_account_iam_policy":            ResourceIamPolicy(IamBillingAccountSchema, NewBillingAccountIamUpdater, BillingAccountIdParseFunc),
			"google_cloudfunctions_function":               resourceCloudFunctionsFunction(),
			"google_cloudiot_registry":                     resourceCloudIoTRegistry(),
			"google_composer_environment":                  resourceComposerEnvironment(),
			"google_compute_attached_disk":                 resourceComputeAttachedDisk(),
			"google_compute_instance":                      resourceComputeInstance(),
			"google_compute_instance_from_template":        resourceComputeInstanceFromTemplate(),
			"google_compute_instance_group":                resourceComputeInstanceGroup(),
			"google_compute_instance_group_manager":        resourceComputeInstanceGroupManager(),
			"google_compute_instance_iam_binding":          ResourceIamBinding(IamComputeInstanceSchema, NewComputeInstanceIamUpdater, ComputeInstanceIdParseFunc),
			"google_compute_instance_iam_member":           ResourceIamMember(IamComputeInstanceSchema, NewComputeInstanceIamUpdater, ComputeInstanceIdParseFunc),
			"google_compute_instance_iam_policy":           ResourceIamPolicy(IamComputeInstanceSchema, NewComputeInstanceIamUpdater, ComputeInstanceIdParseFunc),
			"google_compute_instance_template":             resourceComputeInstanceTemplate(),
			"google_compute_network_peering":               resourceComputeNetworkPeering(),
			"google_compute_project_default_network_tier":  resourceComputeProjectDefaultNetworkTier(),
			"google_compute_project_metadata":              resourceComputeProjectMetadata(),
			"google_compute_project_metadata_item":         resourceComputeProjectMetadataItem(),
			"google_compute_region_instance_group_manager": resourceComputeRegionInstanceGroupManager(),
			"google_compute_router_interface":              resourceComputeRouterInterface(),
			"google_compute_router_nat":                    resourceComputeRouterNat(),
			"google_compute_router_peer":                   resourceComputeRouterPeer(),
			"google_compute_security_policy":               resourceComputeSecurityPolicy(),
			"google_compute_shared_vpc_host_project":       resourceComputeSharedVpcHostProject(),
			"google_compute_shared_vpc_service_project":    resourceComputeSharedVpcServiceProject(),
			"google_compute_subnetwork_iam_binding":        ResourceIamBinding(IamComputeSubnetworkSchema, NewComputeSubnetworkIamUpdater, ComputeSubnetworkIdParseFunc),
			"google_compute_subnetwork_iam_member":         ResourceIamMember(IamComputeSubnetworkSchema, NewComputeSubnetworkIamUpdater, ComputeSubnetworkIdParseFunc),
			"google_compute_subnetwork_iam_policy":         ResourceIamPolicy(IamComputeSubnetworkSchema, NewComputeSubnetworkIamUpdater, ComputeSubnetworkIdParseFunc),
			"google_compute_target_pool":                   resourceComputeTargetPool(),
			"google_container_cluster":                     resourceContainerCluster(),
			"google_container_node_pool":                   resourceContainerNodePool(),
			"google_dataflow_job":                          resourceDataflowJob(),
			"google_dataproc_cluster":                      resourceDataprocCluster(),
			"google_dataproc_cluster_iam_binding":          ResourceIamBinding(IamDataprocClusterSchema, NewDataprocClusterUpdater, DataprocClusterIdParseFunc),
			"google_dataproc_cluster_iam_member":           ResourceIamMember(IamDataprocClusterSchema, NewDataprocClusterUpdater, DataprocClusterIdParseFunc),
			"google_dataproc_cluster_iam_policy":           ResourceIamPolicy(IamDataprocClusterSchema, NewDataprocClusterUpdater, DataprocClusterIdParseFunc),
			"google_dataproc_job":                          resourceDataprocJob(),
			"google_dataproc_job_iam_binding":              ResourceIamBinding(IamDataprocJobSchema, NewDataprocJobUpdater, DataprocJobIdParseFunc),
			"google_dataproc_job_iam_member":               ResourceIamMember(IamDataprocJobSchema, NewDataprocJobUpdater, DataprocJobIdParseFunc),
			"google_dataproc_job_iam_policy":               ResourceIamPolicy(IamDataprocJobSchema, NewDataprocJobUpdater, DataprocJobIdParseFunc),
			"google_dns_record_set":                        resourceDnsRecordSet(),
			"google_endpoints_service":                     resourceEndpointsService(),
			"google_folder":                                resourceGoogleFolder(),
			"google_folder_iam_binding":                    ResourceIamBinding(IamFolderSchema, NewFolderIamUpdater, FolderIdParseFunc),
			"google_folder_iam_member":                     ResourceIamMember(IamFolderSchema, NewFolderIamUpdater, FolderIdParseFunc),
			"google_folder_iam_policy":                     ResourceIamPolicy(IamFolderSchema, NewFolderIamUpdater, FolderIdParseFunc),
			"google_folder_organization_policy":            resourceGoogleFolderOrganizationPolicy(),
			"google_logging_billing_account_sink":          resourceLoggingBillingAccountSink(),
			"google_logging_billing_account_exclusion":     ResourceLoggingExclusion(BillingAccountLoggingExclusionSchema, NewBillingAccountLoggingExclusionUpdater, billingAccountLoggingExclusionIdParseFunc),
			"google_logging_organization_sink":             resourceLoggingOrganizationSink(),
			"google_logging_organization_exclusion":        ResourceLoggingExclusion(OrganizationLoggingExclusionSchema, NewOrganizationLoggingExclusionUpdater, organizationLoggingExclusionIdParseFunc),
			"google_logging_folder_sink":                   resourceLoggingFolderSink(),
			"google_logging_folder_exclusion":              ResourceLoggingExclusion(FolderLoggingExclusionSchema, NewFolderLoggingExclusionUpdater, folderLoggingExclusionIdParseFunc),
			"google_logging_project_sink":                  resourceLoggingProjectSink(),
			"google_logging_project_exclusion":             ResourceLoggingExclusion(ProjectLoggingExclusionSchema, NewProjectLoggingExclusionUpdater, projectLoggingExclusionIdParseFunc),
			"google_kms_key_ring_iam_binding":              ResourceIamBinding(IamKmsKeyRingSchema, NewKmsKeyRingIamUpdater, KeyRingIdParseFunc),
			"google_kms_key_ring_iam_member":               ResourceIamMember(IamKmsKeyRingSchema, NewKmsKeyRingIamUpdater, KeyRingIdParseFunc),
			"google_kms_key_ring_iam_policy":               ResourceIamPolicy(IamKmsKeyRingSchema, NewKmsKeyRingIamUpdater, KeyRingIdParseFunc),
			"google_kms_crypto_key_iam_binding":            ResourceIamBinding(IamKmsCryptoKeySchema, NewKmsCryptoKeyIamUpdater, CryptoIdParseFunc),
			"google_kms_crypto_key_iam_member":             ResourceIamMember(IamKmsCryptoKeySchema, NewKmsCryptoKeyIamUpdater, CryptoIdParseFunc),
			"google_service_networking_connection":         resourceServiceNetworkingConnection(),
			"google_spanner_instance_iam_binding":          ResourceIamBinding(IamSpannerInstanceSchema, NewSpannerInstanceIamUpdater, SpannerInstanceIdParseFunc),
			"google_spanner_instance_iam_member":           ResourceIamMember(IamSpannerInstanceSchema, NewSpannerInstanceIamUpdater, SpannerInstanceIdParseFunc),
			"google_spanner_instance_iam_policy":           ResourceIamPolicy(IamSpannerInstanceSchema, NewSpannerInstanceIamUpdater, SpannerInstanceIdParseFunc),
			"google_spanner_database_iam_binding":          ResourceIamBinding(IamSpannerDatabaseSchema, NewSpannerDatabaseIamUpdater, SpannerDatabaseIdParseFunc),
			"google_spanner_database_iam_member":           ResourceIamMember(IamSpannerDatabaseSchema, NewSpannerDatabaseIamUpdater, SpannerDatabaseIdParseFunc),
			"google_spanner_database_iam_policy":           ResourceIamPolicy(IamSpannerDatabaseSchema, NewSpannerDatabaseIamUpdater, SpannerDatabaseIdParseFunc),
			"google_sql_database_instance":                 resourceSqlDatabaseInstance(),
			"google_sql_ssl_cert":                          resourceSqlSslCert(),
			"google_sql_user":                              resourceSqlUser(),
			"google_organization_iam_binding":              ResourceIamBinding(IamOrganizationSchema, NewOrganizationIamUpdater, OrgIdParseFunc),
			"google_organization_iam_custom_role":          resourceGoogleOrganizationIamCustomRole(),
			"google_organization_iam_member":               ResourceIamMember(IamOrganizationSchema, NewOrganizationIamUpdater, OrgIdParseFunc),
			"google_organization_iam_policy":               ResourceIamPolicy(IamOrganizationSchema, NewOrganizationIamUpdater, OrgIdParseFunc),
			"google_organization_policy":                   resourceGoogleOrganizationPolicy(),
			"google_project":                               resourceGoogleProject(),
			"google_project_iam_policy":                    resourceGoogleProjectIamPolicy(),
			"google_project_iam_binding":                   ResourceIamBindingWithBatching(IamProjectSchema, NewProjectIamUpdater, ProjectIdParseFunc, IamBatchingEnabled),
			"google_project_iam_member":                    ResourceIamMemberWithBatching(IamProjectSchema, NewProjectIamUpdater, ProjectIdParseFunc, IamBatchingEnabled),
			"google_project_iam_audit_config":              ResourceIamAuditConfigWithBatching(IamProjectSchema, NewProjectIamUpdater, ProjectIdParseFunc, IamBatchingEnabled),
			"google_project_service":                       resourceGoogleProjectService(),
			"google_project_iam_custom_role":               resourceGoogleProjectIamCustomRole(),
			"google_project_organization_policy":           resourceGoogleProjectOrganizationPolicy(),
			"google_project_usage_export_bucket":           resourceProjectUsageBucket(),
			"google_project_services":                      resourceGoogleProjectServices(),
			"google_pubsub_subscription_iam_binding":       ResourceIamBinding(IamPubsubSubscriptionSchema, NewPubsubSubscriptionIamUpdater, PubsubSubscriptionIdParseFunc),
			"google_pubsub_subscription_iam_member":        ResourceIamMember(IamPubsubSubscriptionSchema, NewPubsubSubscriptionIamUpdater, PubsubSubscriptionIdParseFunc),
			"google_pubsub_subscription_iam_policy":        ResourceIamPolicy(IamPubsubSubscriptionSchema, NewPubsubSubscriptionIamUpdater, PubsubSubscriptionIdParseFunc),
			"google_runtimeconfig_config":                  resourceRuntimeconfigConfig(),
			"google_runtimeconfig_variable":                resourceRuntimeconfigVariable(),
			"google_service_account":                       resourceGoogleServiceAccount(),
			"google_service_account_iam_binding":           ResourceIamBinding(IamServiceAccountSchema, NewServiceAccountIamUpdater, ServiceAccountIdParseFunc),
			"google_service_account_iam_member":            ResourceIamMember(IamServiceAccountSchema, NewServiceAccountIamUpdater, ServiceAccountIdParseFunc),
			"google_service_account_iam_policy":            ResourceIamPolicy(IamServiceAccountSchema, NewServiceAccountIamUpdater, ServiceAccountIdParseFunc),
			"google_service_account_key":                   resourceGoogleServiceAccountKey(),
			"google_storage_bucket":                        resourceStorageBucket(),
			"google_storage_bucket_acl":                    resourceStorageBucketAcl(),
			// Legacy roles such as roles/storage.legacyBucketReader are automatically added
			// when creating a bucket. For this reason, it is better not to add the authoritative
			// google_storage_bucket_iam_policy resource.
			"google_storage_bucket_iam_binding": ResourceIamBinding(IamStorageBucketSchema, NewStorageBucketIamUpdater, StorageBucketIdParseFunc),
			"google_storage_bucket_iam_member":  ResourceIamMember(IamStorageBucketSchema, NewStorageBucketIamUpdater, StorageBucketIdParseFunc),
			"google_storage_bucket_iam_policy":  ResourceIamPolicy(IamStorageBucketSchema, NewStorageBucketIamUpdater, StorageBucketIdParseFunc),
			"google_storage_bucket_object":      resourceStorageBucketObject(),
			"google_storage_object_acl":         resourceStorageObjectAcl(),
			"google_storage_default_object_acl": resourceStorageDefaultObjectAcl(),
			"google_storage_notification":       resourceStorageNotification(),
			"google_storage_transfer_job":       resourceStorageTransferJob(),
		},
	)
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Project:             d.Get("project").(string),
		Region:              d.Get("region").(string),
		Zone:                d.Get("zone").(string),
		UserProjectOverride: d.Get("user_project_override").(bool),
	}

	// Add credential source
	if v, ok := d.GetOk("access_token"); ok {
		config.AccessToken = v.(string)
	} else if v, ok := d.GetOk("credentials"); ok {
		config.Credentials = v.(string)
	}

	scopes := d.Get("scopes").([]interface{})
	if len(scopes) > 0 {
		config.Scopes = make([]string, len(scopes), len(scopes))
	}
	for i, scope := range scopes {
		config.Scopes[i] = scope.(string)
	}

	batchCfg, err := expandProviderBatchingConfig(d.Get("batching"))
	if err != nil {
		return nil, err
	}
	config.BatchingConfig = batchCfg

	// Generated products
	config.AccessContextManagerBasePath = d.Get(AccessContextManagerCustomEndpointEntryKey).(string)
	config.AppEngineBasePath = d.Get(AppEngineCustomEndpointEntryKey).(string)
	config.BigqueryDataTransferBasePath = d.Get(BigqueryDataTransferCustomEndpointEntryKey).(string)
	config.BinaryAuthorizationBasePath = d.Get(BinaryAuthorizationCustomEndpointEntryKey).(string)
	config.CloudBuildBasePath = d.Get(CloudBuildCustomEndpointEntryKey).(string)
	config.CloudSchedulerBasePath = d.Get(CloudSchedulerCustomEndpointEntryKey).(string)
	config.ComputeBasePath = d.Get(ComputeCustomEndpointEntryKey).(string)
	config.DnsBasePath = d.Get(DnsCustomEndpointEntryKey).(string)
	config.FilestoreBasePath = d.Get(FilestoreCustomEndpointEntryKey).(string)
	config.FirestoreBasePath = d.Get(FirestoreCustomEndpointEntryKey).(string)
	config.KmsBasePath = d.Get(KmsCustomEndpointEntryKey).(string)
	config.LoggingBasePath = d.Get(LoggingCustomEndpointEntryKey).(string)
	config.MonitoringBasePath = d.Get(MonitoringCustomEndpointEntryKey).(string)
	config.PubsubBasePath = d.Get(PubsubCustomEndpointEntryKey).(string)
	config.RedisBasePath = d.Get(RedisCustomEndpointEntryKey).(string)
	config.ResourceManagerBasePath = d.Get(ResourceManagerCustomEndpointEntryKey).(string)
	config.SecurityCenterBasePath = d.Get(SecurityCenterCustomEndpointEntryKey).(string)
	config.SourceRepoBasePath = d.Get(SourceRepoCustomEndpointEntryKey).(string)
	config.SpannerBasePath = d.Get(SpannerCustomEndpointEntryKey).(string)
	config.SqlBasePath = d.Get(SqlCustomEndpointEntryKey).(string)
	config.StorageBasePath = d.Get(StorageCustomEndpointEntryKey).(string)
	config.TpuBasePath = d.Get(TpuCustomEndpointEntryKey).(string)

	// Handwritten Products / Versioned / Atypical Entries

	config.CloudBillingBasePath = d.Get(CloudBillingCustomEndpointEntryKey).(string)
	config.ComposerBasePath = d.Get(ComposerCustomEndpointEntryKey).(string)
	config.ComputeBetaBasePath = d.Get(ComputeBetaCustomEndpointEntryKey).(string)
	config.ContainerBasePath = d.Get(ContainerCustomEndpointEntryKey).(string)
	config.ContainerBetaBasePath = d.Get(ContainerBetaCustomEndpointEntryKey).(string)
	config.DataprocBasePath = d.Get(DataprocCustomEndpointEntryKey).(string)
	config.DataprocBetaBasePath = d.Get(DataprocBetaCustomEndpointEntryKey).(string)
	config.DataflowBasePath = d.Get(DataflowCustomEndpointEntryKey).(string)
	config.DnsBetaBasePath = d.Get(DnsBetaCustomEndpointEntryKey).(string)
	config.IamCredentialsBasePath = d.Get(IamCredentialsCustomEndpointEntryKey).(string)
	config.ResourceManagerV2Beta1BasePath = d.Get(ResourceManagerV2Beta1CustomEndpointEntryKey).(string)
	config.RuntimeconfigBasePath = d.Get(RuntimeconfigCustomEndpointEntryKey).(string)
	config.IAMBasePath = d.Get(IAMCustomEndpointEntryKey).(string)
	config.ServiceManagementBasePath = d.Get(ServiceManagementCustomEndpointEntryKey).(string)
	config.ServiceNetworkingBasePath = d.Get(ServiceNetworkingCustomEndpointEntryKey).(string)
	config.ServiceUsageBasePath = d.Get(ServiceUsageCustomEndpointEntryKey).(string)
	config.BigQueryBasePath = d.Get(BigQueryCustomEndpointEntryKey).(string)
	config.CloudFunctionsBasePath = d.Get(CloudFunctionsCustomEndpointEntryKey).(string)
	config.CloudIoTBasePath = d.Get(CloudIoTCustomEndpointEntryKey).(string)
	config.StorageTransferBasePath = d.Get(StorageTransferCustomEndpointEntryKey).(string)
	config.BigtableAdminBasePath = d.Get(BigtableAdminCustomEndpointEntryKey).(string)

	if err := config.LoadAndValidate(); err != nil {
		return nil, err
	}

	return &config, nil
}

// For a consumer of config.go that isn't a full fledged provider and doesn't
// have its own endpoint mechanism such as sweepers, init {{service}}BasePath
// values to a default. After using this, you should call config.LoadAndValidate.
func ConfigureBasePaths(c *Config) {
	// Generated Products
	c.AccessContextManagerBasePath = AccessContextManagerDefaultBasePath
	c.AppEngineBasePath = AppEngineDefaultBasePath
	c.BigqueryDataTransferBasePath = BigqueryDataTransferDefaultBasePath
	c.BinaryAuthorizationBasePath = BinaryAuthorizationDefaultBasePath
	c.CloudBuildBasePath = CloudBuildDefaultBasePath
	c.CloudSchedulerBasePath = CloudSchedulerDefaultBasePath
	c.ComputeBasePath = ComputeDefaultBasePath
	c.DnsBasePath = DnsDefaultBasePath
	c.FilestoreBasePath = FilestoreDefaultBasePath
	c.FirestoreBasePath = FirestoreDefaultBasePath
	c.KmsBasePath = KmsDefaultBasePath
	c.LoggingBasePath = LoggingDefaultBasePath
	c.MonitoringBasePath = MonitoringDefaultBasePath
	c.PubsubBasePath = PubsubDefaultBasePath
	c.RedisBasePath = RedisDefaultBasePath
	c.ResourceManagerBasePath = ResourceManagerDefaultBasePath
	c.SecurityCenterBasePath = SecurityCenterDefaultBasePath
	c.SourceRepoBasePath = SourceRepoDefaultBasePath
	c.SpannerBasePath = SpannerDefaultBasePath
	c.SqlBasePath = SqlDefaultBasePath
	c.StorageBasePath = StorageDefaultBasePath
	c.TpuBasePath = TpuDefaultBasePath

	// Handwritten Products / Versioned / Atypical Entries
	c.CloudBillingBasePath = CloudBillingDefaultBasePath
	c.ComposerBasePath = ComposerDefaultBasePath
	c.ComputeBetaBasePath = ComputeBetaDefaultBasePath
	c.ContainerBasePath = ContainerDefaultBasePath
	c.ContainerBetaBasePath = ContainerBetaDefaultBasePath
	c.DataprocBasePath = DataprocDefaultBasePath
	c.DataflowBasePath = DataflowDefaultBasePath
	c.DnsBetaBasePath = DnsBetaDefaultBasePath
	c.IamCredentialsBasePath = IamCredentialsDefaultBasePath
	c.ResourceManagerV2Beta1BasePath = ResourceManagerV2Beta1DefaultBasePath
	c.RuntimeconfigBasePath = RuntimeconfigDefaultBasePath
	c.IAMBasePath = IAMDefaultBasePath
	c.ServiceManagementBasePath = ServiceManagementDefaultBasePath
	c.ServiceNetworkingBasePath = ServiceNetworkingDefaultBasePath
	c.ServiceUsageBasePath = ServiceUsageDefaultBasePath
	c.BigQueryBasePath = BigQueryDefaultBasePath
	c.CloudFunctionsBasePath = CloudFunctionsDefaultBasePath
	c.CloudIoTBasePath = CloudIoTDefaultBasePath
	c.StorageTransferBasePath = StorageTransferDefaultBasePath
	c.BigtableAdminBasePath = BigtableAdminDefaultBasePath
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
			fmt.Errorf("JSON credentials in %q are not valid: %s", creds, err))
	}

	return
}
