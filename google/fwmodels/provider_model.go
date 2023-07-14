// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package fwmodels

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ProviderModel describes the provider config data model.
type ProviderModel struct {
	Credentials                        types.String `tfsdk:"credentials"`
	AccessToken                        types.String `tfsdk:"access_token"`
	ImpersonateServiceAccount          types.String `tfsdk:"impersonate_service_account"`
	ImpersonateServiceAccountDelegates types.List   `tfsdk:"impersonate_service_account_delegates"`
	Project                            types.String `tfsdk:"project"`
	BillingProject                     types.String `tfsdk:"billing_project"`
	Region                             types.String `tfsdk:"region"`
	Zone                               types.String `tfsdk:"zone"`
	Scopes                             types.List   `tfsdk:"scopes"`
	Batching                           types.List   `tfsdk:"batching"`
	UserProjectOverride                types.Bool   `tfsdk:"user_project_override"`
	RequestTimeout                     types.String `tfsdk:"request_timeout"`
	RequestReason                      types.String `tfsdk:"request_reason"`

	// Generated Products
	AccessApprovalCustomEndpoint           types.String `tfsdk:"access_approval_custom_endpoint"`
	AccessContextManagerCustomEndpoint     types.String `tfsdk:"access_context_manager_custom_endpoint"`
	ActiveDirectoryCustomEndpoint          types.String `tfsdk:"active_directory_custom_endpoint"`
	AlloydbCustomEndpoint                  types.String `tfsdk:"alloydb_custom_endpoint"`
	ApigeeCustomEndpoint                   types.String `tfsdk:"apigee_custom_endpoint"`
	AppEngineCustomEndpoint                types.String `tfsdk:"app_engine_custom_endpoint"`
	ArtifactRegistryCustomEndpoint         types.String `tfsdk:"artifact_registry_custom_endpoint"`
	BeyondcorpCustomEndpoint               types.String `tfsdk:"beyondcorp_custom_endpoint"`
	BigQueryCustomEndpoint                 types.String `tfsdk:"big_query_custom_endpoint"`
	BigqueryAnalyticsHubCustomEndpoint     types.String `tfsdk:"bigquery_analytics_hub_custom_endpoint"`
	BigqueryConnectionCustomEndpoint       types.String `tfsdk:"bigquery_connection_custom_endpoint"`
	BigqueryDatapolicyCustomEndpoint       types.String `tfsdk:"bigquery_datapolicy_custom_endpoint"`
	BigqueryDataTransferCustomEndpoint     types.String `tfsdk:"bigquery_data_transfer_custom_endpoint"`
	BigqueryReservationCustomEndpoint      types.String `tfsdk:"bigquery_reservation_custom_endpoint"`
	BigtableCustomEndpoint                 types.String `tfsdk:"bigtable_custom_endpoint"`
	BillingCustomEndpoint                  types.String `tfsdk:"billing_custom_endpoint"`
	BinaryAuthorizationCustomEndpoint      types.String `tfsdk:"binary_authorization_custom_endpoint"`
	CertificateManagerCustomEndpoint       types.String `tfsdk:"certificate_manager_custom_endpoint"`
	CloudAssetCustomEndpoint               types.String `tfsdk:"cloud_asset_custom_endpoint"`
	CloudBuildCustomEndpoint               types.String `tfsdk:"cloud_build_custom_endpoint"`
	Cloudbuildv2CustomEndpoint             types.String `tfsdk:"cloudbuildv2_custom_endpoint"`
	CloudFunctionsCustomEndpoint           types.String `tfsdk:"cloud_functions_custom_endpoint"`
	Cloudfunctions2CustomEndpoint          types.String `tfsdk:"cloudfunctions2_custom_endpoint"`
	CloudIdentityCustomEndpoint            types.String `tfsdk:"cloud_identity_custom_endpoint"`
	CloudIdsCustomEndpoint                 types.String `tfsdk:"cloud_ids_custom_endpoint"`
	CloudIotCustomEndpoint                 types.String `tfsdk:"cloud_iot_custom_endpoint"`
	CloudRunCustomEndpoint                 types.String `tfsdk:"cloud_run_custom_endpoint"`
	CloudRunV2CustomEndpoint               types.String `tfsdk:"cloud_run_v2_custom_endpoint"`
	CloudSchedulerCustomEndpoint           types.String `tfsdk:"cloud_scheduler_custom_endpoint"`
	CloudTasksCustomEndpoint               types.String `tfsdk:"cloud_tasks_custom_endpoint"`
	ComputeCustomEndpoint                  types.String `tfsdk:"compute_custom_endpoint"`
	ContainerAnalysisCustomEndpoint        types.String `tfsdk:"container_analysis_custom_endpoint"`
	ContainerAttachedCustomEndpoint        types.String `tfsdk:"container_attached_custom_endpoint"`
	DatabaseMigrationServiceCustomEndpoint types.String `tfsdk:"database_migration_service_custom_endpoint"`
	DataCatalogCustomEndpoint              types.String `tfsdk:"data_catalog_custom_endpoint"`
	DataFusionCustomEndpoint               types.String `tfsdk:"data_fusion_custom_endpoint"`
	DataLossPreventionCustomEndpoint       types.String `tfsdk:"data_loss_prevention_custom_endpoint"`
	DataplexCustomEndpoint                 types.String `tfsdk:"dataplex_custom_endpoint"`
	DataprocCustomEndpoint                 types.String `tfsdk:"dataproc_custom_endpoint"`
	DataprocMetastoreCustomEndpoint        types.String `tfsdk:"dataproc_metastore_custom_endpoint"`
	DatastoreCustomEndpoint                types.String `tfsdk:"datastore_custom_endpoint"`
	DatastreamCustomEndpoint               types.String `tfsdk:"datastream_custom_endpoint"`
	DeploymentManagerCustomEndpoint        types.String `tfsdk:"deployment_manager_custom_endpoint"`
	DialogflowCustomEndpoint               types.String `tfsdk:"dialogflow_custom_endpoint"`
	DialogflowCXCustomEndpoint             types.String `tfsdk:"dialogflow_cx_custom_endpoint"`
	DNSCustomEndpoint                      types.String `tfsdk:"dns_custom_endpoint"`
	DocumentAICustomEndpoint               types.String `tfsdk:"document_ai_custom_endpoint"`
	EssentialContactsCustomEndpoint        types.String `tfsdk:"essential_contacts_custom_endpoint"`
	FilestoreCustomEndpoint                types.String `tfsdk:"filestore_custom_endpoint"`
	FirestoreCustomEndpoint                types.String `tfsdk:"firestore_custom_endpoint"`
	GameServicesCustomEndpoint             types.String `tfsdk:"game_services_custom_endpoint"`
	GKEBackupCustomEndpoint                types.String `tfsdk:"gke_backup_custom_endpoint"`
	GKEHubCustomEndpoint                   types.String `tfsdk:"gke_hub_custom_endpoint"`
	GKEHub2CustomEndpoint                  types.String `tfsdk:"gke_hub2_custom_endpoint"`
	HealthcareCustomEndpoint               types.String `tfsdk:"healthcare_custom_endpoint"`
	IAM2CustomEndpoint                     types.String `tfsdk:"iam2_custom_endpoint"`
	IAMBetaCustomEndpoint                  types.String `tfsdk:"iam_beta_custom_endpoint"`
	IAMWorkforcePoolCustomEndpoint         types.String `tfsdk:"iam_workforce_pool_custom_endpoint"`
	IapCustomEndpoint                      types.String `tfsdk:"iap_custom_endpoint"`
	IdentityPlatformCustomEndpoint         types.String `tfsdk:"identity_platform_custom_endpoint"`
	KMSCustomEndpoint                      types.String `tfsdk:"kms_custom_endpoint"`
	LoggingCustomEndpoint                  types.String `tfsdk:"logging_custom_endpoint"`
	LookerCustomEndpoint                   types.String `tfsdk:"looker_custom_endpoint"`
	MemcacheCustomEndpoint                 types.String `tfsdk:"memcache_custom_endpoint"`
	MLEngineCustomEndpoint                 types.String `tfsdk:"ml_engine_custom_endpoint"`
	MonitoringCustomEndpoint               types.String `tfsdk:"monitoring_custom_endpoint"`
	NetworkManagementCustomEndpoint        types.String `tfsdk:"network_management_custom_endpoint"`
	NetworkSecurityCustomEndpoint          types.String `tfsdk:"network_security_custom_endpoint"`
	NetworkServicesCustomEndpoint          types.String `tfsdk:"network_services_custom_endpoint"`
	NotebooksCustomEndpoint                types.String `tfsdk:"notebooks_custom_endpoint"`
	OSConfigCustomEndpoint                 types.String `tfsdk:"os_config_custom_endpoint"`
	OSLoginCustomEndpoint                  types.String `tfsdk:"os_login_custom_endpoint"`
	PrivatecaCustomEndpoint                types.String `tfsdk:"privateca_custom_endpoint"`
	PublicCACustomEndpoint                 types.String `tfsdk:"public_ca_custom_endpoint"`
	PubsubCustomEndpoint                   types.String `tfsdk:"pubsub_custom_endpoint"`
	PubsubLiteCustomEndpoint               types.String `tfsdk:"pubsub_lite_custom_endpoint"`
	RedisCustomEndpoint                    types.String `tfsdk:"redis_custom_endpoint"`
	ResourceManagerCustomEndpoint          types.String `tfsdk:"resource_manager_custom_endpoint"`
	SecretManagerCustomEndpoint            types.String `tfsdk:"secret_manager_custom_endpoint"`
	SecurityCenterCustomEndpoint           types.String `tfsdk:"security_center_custom_endpoint"`
	ServiceManagementCustomEndpoint        types.String `tfsdk:"service_management_custom_endpoint"`
	ServiceUsageCustomEndpoint             types.String `tfsdk:"service_usage_custom_endpoint"`
	SourceRepoCustomEndpoint               types.String `tfsdk:"source_repo_custom_endpoint"`
	SpannerCustomEndpoint                  types.String `tfsdk:"spanner_custom_endpoint"`
	SQLCustomEndpoint                      types.String `tfsdk:"sql_custom_endpoint"`
	StorageCustomEndpoint                  types.String `tfsdk:"storage_custom_endpoint"`
	StorageTransferCustomEndpoint          types.String `tfsdk:"storage_transfer_custom_endpoint"`
	TagsCustomEndpoint                     types.String `tfsdk:"tags_custom_endpoint"`
	TPUCustomEndpoint                      types.String `tfsdk:"tpu_custom_endpoint"`
	VertexAICustomEndpoint                 types.String `tfsdk:"vertex_ai_custom_endpoint"`
	VPCAccessCustomEndpoint                types.String `tfsdk:"vpc_access_custom_endpoint"`
	WorkflowsCustomEndpoint                types.String `tfsdk:"workflows_custom_endpoint"`

	// Handwritten Products / Versioned / Atypical Entries
	CloudBillingCustomEndpoint      types.String `tfsdk:"cloud_billing_custom_endpoint"`
	ComposerCustomEndpoint          types.String `tfsdk:"composer_custom_endpoint"`
	ContainerCustomEndpoint         types.String `tfsdk:"container_custom_endpoint"`
	DataflowCustomEndpoint          types.String `tfsdk:"dataflow_custom_endpoint"`
	IamCredentialsCustomEndpoint    types.String `tfsdk:"iam_credentials_custom_endpoint"`
	ResourceManagerV3CustomEndpoint types.String `tfsdk:"resource_manager_v3_custom_endpoint"`
	IAMCustomEndpoint               types.String `tfsdk:"iam_custom_endpoint"`
	ServiceNetworkingCustomEndpoint types.String `tfsdk:"service_networking_custom_endpoint"`
	TagsLocationCustomEndpoint      types.String `tfsdk:"tags_location_custom_endpoint"`

	// dcl
	ContainerAwsCustomEndpoint   types.String `tfsdk:"container_aws_custom_endpoint"`
	ContainerAzureCustomEndpoint types.String `tfsdk:"container_azure_custom_endpoint"`

	// dcl generated
	ApikeysCustomEndpoint              types.String `tfsdk:"apikeys_custom_endpoint"`
	AssuredWorkloadsCustomEndpoint     types.String `tfsdk:"assured_workloads_custom_endpoint"`
	CloudBuildWorkerPoolCustomEndpoint types.String `tfsdk:"cloud_build_worker_pool_custom_endpoint"`
	CloudDeployCustomEndpoint          types.String `tfsdk:"clouddeploy_custom_endpoint"`
	CloudResourceManagerCustomEndpoint types.String `tfsdk:"cloud_resource_manager_custom_endpoint"`
	EventarcCustomEndpoint             types.String `tfsdk:"eventarc_custom_endpoint"`
	FirebaserulesCustomEndpoint        types.String `tfsdk:"firebaserules_custom_endpoint"`
	NetworkConnectivityCustomEndpoint  types.String `tfsdk:"network_connectivity_custom_endpoint"`
	OrgPolicyCustomEndpoint            types.String `tfsdk:"org_policy_custom_endpoint"`
	RecaptchaEnterpriseCustomEndpoint  types.String `tfsdk:"recaptcha_enterprise_custom_endpoint"`
}

type ProviderBatching struct {
	SendAfter      types.String `tfsdk:"send_after"`
	EnableBatching types.Bool   `tfsdk:"enable_batching"`
}

var ProviderBatchingAttributes = map[string]attr.Type{
	"send_after":      types.StringType,
	"enable_batching": types.BoolType,
}

// ProviderMetaModel describes the provider meta model
type ProviderMetaModel struct {
	ModuleName types.String `tfsdk:"module_name"`
}
