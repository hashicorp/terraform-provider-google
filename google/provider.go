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
		},

		DataSourcesMap: map[string]*schema.Resource{
			"google_active_folder":                            dataSourceGoogleActiveFolder(),
			"google_billing_account":                          dataSourceGoogleBillingAccount(),
			"google_dns_managed_zone":                         dataSourceDnsManagedZone(),
			"google_client_config":                            dataSourceGoogleClientConfig(),
			"google_cloudfunctions_function":                  dataSourceGoogleCloudFunctionsFunction(),
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
			"google_compute_regions":                          dataSourceGoogleComputeRegions(),
			"google_compute_region_instance_group":            dataSourceGoogleComputeRegionInstanceGroup(),
			"google_compute_subnetwork":                       dataSourceGoogleComputeSubnetwork(),
			"google_compute_zones":                            dataSourceGoogleComputeZones(),
			"google_compute_vpn_gateway":                      dataSourceGoogleComputeVpnGateway(),
			"google_compute_ssl_policy":                       dataSourceGoogleComputeSslPolicy(),
			"google_container_cluster":                        dataSourceGoogleContainerCluster(),
			"google_container_engine_versions":                dataSourceGoogleContainerEngineVersions(),
			"google_container_registry_repository":            dataSourceGoogleContainerRepo(),
			"google_container_registry_image":                 dataSourceGoogleContainerImage(),
			"google_iam_policy":                               dataSourceGoogleIamPolicy(),
			"google_iam_role":                                 dataSourceGoogleIamRole(),
			"google_kms_secret":                               dataSourceGoogleKmsSecret(),
			"google_kms_key_ring":                             dataSourceGoogleKmsKeyRing(),
			"google_kms_crypto_key":                           dataSourceGoogleKmsCryptoKey(),
			"google_folder":                                   dataSourceGoogleFolder(),
			"google_netblock_ip_ranges":                       dataSourceGoogleNetblockIpRanges(),
			"google_organization":                             dataSourceGoogleOrganization(),
			"google_project":                                  dataSourceGoogleProject(),
			"google_project_services":                         dataSourceGoogleProjectServices(),
			"google_service_account":                          dataSourceGoogleServiceAccount(),
			"google_service_account_key":                      dataSourceGoogleServiceAccountKey(),
			"google_storage_bucket_object":                    dataSourceGoogleStorageBucketObject(),
			"google_storage_object_signed_url":                dataSourceGoogleSignedUrl(),
			"google_storage_project_service_account":          dataSourceGoogleStorageProjectServiceAccount(),
			"google_storage_transfer_project_service_account": dataSourceGoogleStorageTransferProjectServiceAccount(),
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
		GeneratedAppEngineResourcesMap,
		GeneratedComputeResourcesMap,
		GeneratedCloudBuildResourcesMap,
		GeneratedDnsResourcesMap,
		GeneratedRedisResourcesMap,
		GeneratedResourceManagerResourcesMap,
		GeneratedSourceRepoResourcesMap,
		GeneratedSpannerResourcesMap,
		GeneratedStorageResourcesMap,
		GeneratedMonitoringResourcesMap,
		map[string]*schema.Resource{
			"google_app_engine_application":                resourceAppEngineApplication(),
			"google_bigquery_dataset":                      resourceBigQueryDataset(),
			"google_bigquery_table":                        resourceBigQueryTable(),
			"google_bigtable_instance":                     resourceBigtableInstance(),
			"google_bigtable_table":                        resourceBigtableTable(),
			"google_billing_account_iam_binding":           ResourceIamBindingWithImport(IamBillingAccountSchema, NewBillingAccountIamUpdater, BillingAccountIdParseFunc),
			"google_billing_account_iam_member":            ResourceIamMemberWithImport(IamBillingAccountSchema, NewBillingAccountIamUpdater, BillingAccountIdParseFunc),
			"google_billing_account_iam_policy":            ResourceIamPolicyWithImport(IamBillingAccountSchema, NewBillingAccountIamUpdater, BillingAccountIdParseFunc),
			"google_cloudfunctions_function":               resourceCloudFunctionsFunction(),
			"google_cloudiot_registry":                     resourceCloudIoTRegistry(),
			"google_composer_environment":                  resourceComposerEnvironment(),
			"google_compute_attached_disk":                 resourceComputeAttachedDisk(),
			"google_compute_backend_service":               resourceComputeBackendService(),
			"google_compute_global_forwarding_rule":        resourceComputeGlobalForwardingRule(),
			"google_compute_instance":                      resourceComputeInstance(),
			"google_compute_instance_from_template":        resourceComputeInstanceFromTemplate(),
			"google_compute_instance_group":                resourceComputeInstanceGroup(),
			"google_compute_instance_group_manager":        resourceComputeInstanceGroupManager(),
			"google_compute_instance_template":             resourceComputeInstanceTemplate(),
			"google_compute_network":                       resourceComputeNetwork(),
			"google_compute_network_peering":               resourceComputeNetworkPeering(),
			"google_compute_project_metadata":              resourceComputeProjectMetadata(),
			"google_compute_project_metadata_item":         resourceComputeProjectMetadataItem(),
			"google_compute_region_backend_service":        resourceComputeRegionBackendService(),
			"google_compute_region_instance_group_manager": resourceComputeRegionInstanceGroupManager(),
			"google_compute_router_interface":              resourceComputeRouterInterface(),
			"google_compute_router_nat":                    resourceComputeRouterNat(),
			"google_compute_router_peer":                   resourceComputeRouterPeer(),
			"google_compute_security_policy":               resourceComputeSecurityPolicy(),
			"google_compute_shared_vpc_host_project":       resourceComputeSharedVpcHostProject(),
			"google_compute_shared_vpc_service_project":    resourceComputeSharedVpcServiceProject(),
			"google_compute_target_pool":                   resourceComputeTargetPool(),
			"google_container_cluster":                     resourceContainerCluster(),
			"google_container_node_pool":                   resourceContainerNodePool(),
			"google_dataflow_job":                          resourceDataflowJob(),
			"google_dataproc_cluster":                      resourceDataprocCluster(),
			"google_dataproc_job":                          resourceDataprocJob(),
			"google_dns_record_set":                        resourceDnsRecordSet(),
			"google_endpoints_service":                     resourceEndpointsService(),
			"google_folder":                                resourceGoogleFolder(),
			"google_folder_iam_binding":                    ResourceIamBindingWithImport(IamFolderSchema, NewFolderIamUpdater, FolderIdParseFunc),
			"google_folder_iam_member":                     ResourceIamMemberWithImport(IamFolderSchema, NewFolderIamUpdater, FolderIdParseFunc),
			"google_folder_iam_policy":                     ResourceIamPolicyWithImport(IamFolderSchema, NewFolderIamUpdater, FolderIdParseFunc),
			"google_folder_organization_policy":            resourceGoogleFolderOrganizationPolicy(),
			"google_logging_billing_account_sink":          resourceLoggingBillingAccountSink(),
			"google_logging_billing_account_exclusion":     ResourceLoggingExclusion(BillingAccountLoggingExclusionSchema, NewBillingAccountLoggingExclusionUpdater, billingAccountLoggingExclusionIdParseFunc),
			"google_logging_organization_sink":             resourceLoggingOrganizationSink(),
			"google_logging_organization_exclusion":        ResourceLoggingExclusion(OrganizationLoggingExclusionSchema, NewOrganizationLoggingExclusionUpdater, organizationLoggingExclusionIdParseFunc),
			"google_logging_folder_sink":                   resourceLoggingFolderSink(),
			"google_logging_folder_exclusion":              ResourceLoggingExclusion(FolderLoggingExclusionSchema, NewFolderLoggingExclusionUpdater, folderLoggingExclusionIdParseFunc),
			"google_logging_project_sink":                  resourceLoggingProjectSink(),
			"google_logging_project_exclusion":             ResourceLoggingExclusion(ProjectLoggingExclusionSchema, NewProjectLoggingExclusionUpdater, projectLoggingExclusionIdParseFunc),
			"google_kms_key_ring":                          resourceKmsKeyRing(),
			"google_kms_key_ring_iam_binding":              ResourceIamBindingWithImport(IamKmsKeyRingSchema, NewKmsKeyRingIamUpdater, KeyRingIdParseFunc),
			"google_kms_key_ring_iam_member":               ResourceIamMemberWithImport(IamKmsKeyRingSchema, NewKmsKeyRingIamUpdater, KeyRingIdParseFunc),
			"google_kms_key_ring_iam_policy":               ResourceIamPolicyWithImport(IamKmsKeyRingSchema, NewKmsKeyRingIamUpdater, KeyRingIdParseFunc),
			"google_kms_crypto_key":                        resourceKmsCryptoKey(),
			"google_kms_crypto_key_iam_binding":            ResourceIamBindingWithImport(IamKmsCryptoKeySchema, NewKmsCryptoKeyIamUpdater, CryptoIdParseFunc),
			"google_kms_crypto_key_iam_member":             ResourceIamMemberWithImport(IamKmsCryptoKeySchema, NewKmsCryptoKeyIamUpdater, CryptoIdParseFunc),
			"google_spanner_instance_iam_binding":          ResourceIamBindingWithImport(IamSpannerInstanceSchema, NewSpannerInstanceIamUpdater, SpannerInstanceIdParseFunc),
			"google_spanner_instance_iam_member":           ResourceIamMemberWithImport(IamSpannerInstanceSchema, NewSpannerInstanceIamUpdater, SpannerInstanceIdParseFunc),
			"google_spanner_instance_iam_policy":           ResourceIamPolicyWithImport(IamSpannerInstanceSchema, NewSpannerInstanceIamUpdater, SpannerInstanceIdParseFunc),
			"google_spanner_database_iam_binding":          ResourceIamBindingWithImport(IamSpannerDatabaseSchema, NewSpannerDatabaseIamUpdater, SpannerDatabaseIdParseFunc),
			"google_spanner_database_iam_member":           ResourceIamMemberWithImport(IamSpannerDatabaseSchema, NewSpannerDatabaseIamUpdater, SpannerDatabaseIdParseFunc),
			"google_spanner_database_iam_policy":           ResourceIamPolicyWithImport(IamSpannerDatabaseSchema, NewSpannerDatabaseIamUpdater, SpannerDatabaseIdParseFunc),
			"google_sql_database":                          resourceSqlDatabase(),
			"google_sql_database_instance":                 resourceSqlDatabaseInstance(),
			"google_sql_ssl_cert":                          resourceSqlSslCert(),
			"google_sql_user":                              resourceSqlUser(),
			"google_organization_iam_binding":              ResourceIamBindingWithImport(IamOrganizationSchema, NewOrganizationIamUpdater, OrgIdParseFunc),
			"google_organization_iam_custom_role":          resourceGoogleOrganizationIamCustomRole(),
			"google_organization_iam_member":               ResourceIamMemberWithImport(IamOrganizationSchema, NewOrganizationIamUpdater, OrgIdParseFunc),
			"google_organization_iam_policy":               ResourceIamPolicyWithImport(IamOrganizationSchema, NewOrganizationIamUpdater, OrgIdParseFunc),
			"google_organization_policy":                   resourceGoogleOrganizationPolicy(),
			"google_project":                               resourceGoogleProject(),
			"google_project_iam_policy":                    resourceGoogleProjectIamPolicy(),
			"google_project_iam_binding":                   ResourceIamBindingWithImport(IamProjectSchema, NewProjectIamUpdater, ProjectIdParseFunc),
			"google_project_iam_member":                    ResourceIamMemberWithImport(IamProjectSchema, NewProjectIamUpdater, ProjectIdParseFunc),
			"google_project_iam_audit_config":              ResourceIamAuditConfigWithImport(IamProjectSchema, NewProjectIamUpdater, ProjectIdParseFunc),
			"google_project_service":                       resourceGoogleProjectService(),
			"google_project_iam_custom_role":               resourceGoogleProjectIamCustomRole(),
			"google_project_organization_policy":           resourceGoogleProjectOrganizationPolicy(),
			"google_project_usage_export_bucket":           resourceProjectUsageBucket(),
			"google_project_services":                      resourceGoogleProjectServices(),
			"google_pubsub_topic":                          resourcePubsubTopic(),
			"google_pubsub_topic_iam_binding":              ResourceIamBindingWithImport(IamPubsubTopicSchema, NewPubsubTopicIamUpdater, PubsubTopicIdParseFunc),
			"google_pubsub_topic_iam_member":               ResourceIamMemberWithImport(IamPubsubTopicSchema, NewPubsubTopicIamUpdater, PubsubTopicIdParseFunc),
			"google_pubsub_topic_iam_policy":               ResourceIamPolicyWithImport(IamPubsubTopicSchema, NewPubsubTopicIamUpdater, PubsubTopicIdParseFunc),
			"google_pubsub_subscription":                   resourcePubsubSubscription(),
			"google_pubsub_subscription_iam_binding":       ResourceIamBindingWithImport(IamPubsubSubscriptionSchema, NewPubsubSubscriptionIamUpdater, PubsubSubscriptionIdParseFunc),
			"google_pubsub_subscription_iam_member":        ResourceIamMemberWithImport(IamPubsubSubscriptionSchema, NewPubsubSubscriptionIamUpdater, PubsubSubscriptionIdParseFunc),
			"google_pubsub_subscription_iam_policy":        ResourceIamPolicyWithImport(IamPubsubSubscriptionSchema, NewPubsubSubscriptionIamUpdater, PubsubSubscriptionIdParseFunc),
			"google_runtimeconfig_config":                  resourceRuntimeconfigConfig(),
			"google_runtimeconfig_variable":                resourceRuntimeconfigVariable(),
			"google_service_account":                       resourceGoogleServiceAccount(),
			"google_service_account_iam_binding":           ResourceIamBindingWithImport(IamServiceAccountSchema, NewServiceAccountIamUpdater, ServiceAccountIdParseFunc),
			"google_service_account_iam_member":            ResourceIamMemberWithImport(IamServiceAccountSchema, NewServiceAccountIamUpdater, ServiceAccountIdParseFunc),
			"google_service_account_iam_policy":            ResourceIamPolicyWithImport(IamServiceAccountSchema, NewServiceAccountIamUpdater, ServiceAccountIdParseFunc),
			"google_service_account_key":                   resourceGoogleServiceAccountKey(),
			"google_storage_bucket":                        resourceStorageBucket(),
			"google_storage_bucket_acl":                    resourceStorageBucketAcl(),
			// Legacy roles such as roles/storage.legacyBucketReader are automatically added
			// when creating a bucket. For this reason, it is better not to add the authoritative
			// google_storage_bucket_iam_policy resource.
			"google_storage_bucket_iam_binding": ResourceIamBinding(IamStorageBucketSchema, NewStorageBucketIamUpdater),
			"google_storage_bucket_iam_member":  ResourceIamMember(IamStorageBucketSchema, NewStorageBucketIamUpdater),
			"google_storage_bucket_iam_policy":  ResourceIamPolicy(IamStorageBucketSchema, NewStorageBucketIamUpdater),
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
		Project: d.Get("project").(string),
		Region:  d.Get("region").(string),
		Zone:    d.Get("zone").(string),
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

	if err := config.loadAndValidate(); err != nil {
		return nil, err
	}

	return &config, nil
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
