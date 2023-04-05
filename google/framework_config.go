package google

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"golang.org/x/oauth2"
	googleoauth "golang.org/x/oauth2/google"

	"google.golang.org/api/option"
	"google.golang.org/api/transport"
	"google.golang.org/grpc"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"

	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/sirupsen/logrus"
)

// Provider methods

// LoadAndValidateFramework handles the bulk of configuring the provider
// it is pulled out so that we can manually call this from our testing provider as well
func (p *frameworkProvider) LoadAndValidateFramework(ctx context.Context, data ProviderModel, tfVersion string, diags *diag.Diagnostics) {
	// Set defaults if needed
	p.HandleDefaults(ctx, &data, diags)
	if diags.HasError() {
		return
	}

	p.context = ctx

	// Handle User Agent string
	p.userAgent = CompileUserAgentString(ctx, "terraform-provider-google", tfVersion, p.version)
	// opt in extension for adding to the User-Agent header
	if ext := os.Getenv("GOOGLE_TERRAFORM_USERAGENT_EXTENSION"); ext != "" {
		ua := p.userAgent
		p.userAgent = fmt.Sprintf("%s %s", ua, ext)
	}

	// Set up client configuration
	p.SetupClient(ctx, data, diags)
	if diags.HasError() {
		return
	}

	// gRPC Logging setup
	p.SetupGrpcLogging()

	// Handle Batching Config
	batchingConfig := GetBatchingConfig(ctx, data.Batching, diags)
	if diags.HasError() {
		return
	}

	// Setup Base Paths for clients
	// Generated products
	p.AccessApprovalBasePath = data.AccessApprovalCustomEndpoint.ValueString()
	p.AccessContextManagerBasePath = data.AccessContextManagerCustomEndpoint.ValueString()
	p.ActiveDirectoryBasePath = data.ActiveDirectoryCustomEndpoint.ValueString()
	p.AlloydbBasePath = data.AlloydbCustomEndpoint.ValueString()
	p.ApigeeBasePath = data.ApigeeCustomEndpoint.ValueString()
	p.AppEngineBasePath = data.AppEngineCustomEndpoint.ValueString()
	p.ArtifactRegistryBasePath = data.ArtifactRegistryCustomEndpoint.ValueString()
	p.BeyondcorpBasePath = data.BeyondcorpCustomEndpoint.ValueString()
	p.BigQueryBasePath = data.BigQueryCustomEndpoint.ValueString()
	p.BigqueryAnalyticsHubBasePath = data.BigqueryAnalyticsHubCustomEndpoint.ValueString()
	p.BigqueryConnectionBasePath = data.BigqueryConnectionCustomEndpoint.ValueString()
	p.BigqueryDatapolicyBasePath = data.BigqueryDatapolicyCustomEndpoint.ValueString()
	p.BigqueryDataTransferBasePath = data.BigqueryDataTransferCustomEndpoint.ValueString()
	p.BigqueryReservationBasePath = data.BigqueryReservationCustomEndpoint.ValueString()
	p.BigtableBasePath = data.BigtableCustomEndpoint.ValueString()
	p.BillingBasePath = data.BillingCustomEndpoint.ValueString()
	p.BinaryAuthorizationBasePath = data.BinaryAuthorizationCustomEndpoint.ValueString()
	p.CertificateManagerBasePath = data.CertificateManagerCustomEndpoint.ValueString()
	p.CloudAssetBasePath = data.CloudAssetCustomEndpoint.ValueString()
	p.CloudBuildBasePath = data.CloudBuildCustomEndpoint.ValueString()
	p.CloudFunctionsBasePath = data.CloudFunctionsCustomEndpoint.ValueString()
	p.Cloudfunctions2BasePath = data.Cloudfunctions2CustomEndpoint.ValueString()
	p.CloudIdentityBasePath = data.CloudIdentityCustomEndpoint.ValueString()
	p.CloudIdsBasePath = data.CloudIdsCustomEndpoint.ValueString()
	p.CloudIotBasePath = data.CloudIotCustomEndpoint.ValueString()
	p.CloudRunBasePath = data.CloudRunCustomEndpoint.ValueString()
	p.CloudRunV2BasePath = data.CloudRunV2CustomEndpoint.ValueString()
	p.CloudSchedulerBasePath = data.CloudSchedulerCustomEndpoint.ValueString()
	p.CloudTasksBasePath = data.CloudTasksCustomEndpoint.ValueString()
	p.ComputeBasePath = data.ComputeCustomEndpoint.ValueString()
	p.ContainerAnalysisBasePath = data.ContainerAnalysisCustomEndpoint.ValueString()
	p.ContainerAttachedBasePath = data.ContainerAttachedCustomEndpoint.ValueString()
	p.DataCatalogBasePath = data.DataCatalogCustomEndpoint.ValueString()
	p.DataFusionBasePath = data.DataFusionCustomEndpoint.ValueString()
	p.DataLossPreventionBasePath = data.DataLossPreventionCustomEndpoint.ValueString()
	p.DataplexBasePath = data.DataplexCustomEndpoint.ValueString()
	p.DataprocBasePath = data.DataprocCustomEndpoint.ValueString()
	p.DataprocMetastoreBasePath = data.DataprocMetastoreCustomEndpoint.ValueString()
	p.DatastoreBasePath = data.DatastoreCustomEndpoint.ValueString()
	p.DatastreamBasePath = data.DatastreamCustomEndpoint.ValueString()
	p.DeploymentManagerBasePath = data.DeploymentManagerCustomEndpoint.ValueString()
	p.DialogflowBasePath = data.DialogflowCustomEndpoint.ValueString()
	p.DialogflowCXBasePath = data.DialogflowCXCustomEndpoint.ValueString()
	p.DNSBasePath = data.DNSCustomEndpoint.ValueString()
	p.DocumentAIBasePath = data.DocumentAICustomEndpoint.ValueString()
	p.EssentialContactsBasePath = data.EssentialContactsCustomEndpoint.ValueString()
	p.FilestoreBasePath = data.FilestoreCustomEndpoint.ValueString()
	p.FirestoreBasePath = data.FirestoreCustomEndpoint.ValueString()
	p.GameServicesBasePath = data.GameServicesCustomEndpoint.ValueString()
	p.GKEBackupBasePath = data.GKEBackupCustomEndpoint.ValueString()
	p.GKEHubBasePath = data.GKEHubCustomEndpoint.ValueString()
	p.HealthcareBasePath = data.HealthcareCustomEndpoint.ValueString()
	p.IAM2BasePath = data.IAM2CustomEndpoint.ValueString()
	p.IAMBetaBasePath = data.IAMBetaCustomEndpoint.ValueString()
	p.IAMWorkforcePoolBasePath = data.IAMWorkforcePoolCustomEndpoint.ValueString()
	p.IapBasePath = data.IapCustomEndpoint.ValueString()
	p.IdentityPlatformBasePath = data.IdentityPlatformCustomEndpoint.ValueString()
	p.KMSBasePath = data.KMSCustomEndpoint.ValueString()
	p.LoggingBasePath = data.LoggingCustomEndpoint.ValueString()
	p.MemcacheBasePath = data.MemcacheCustomEndpoint.ValueString()
	p.MLEngineBasePath = data.MLEngineCustomEndpoint.ValueString()
	p.MonitoringBasePath = data.MonitoringCustomEndpoint.ValueString()
	p.NetworkManagementBasePath = data.NetworkManagementCustomEndpoint.ValueString()
	p.NetworkServicesBasePath = data.NetworkServicesCustomEndpoint.ValueString()
	p.NotebooksBasePath = data.NotebooksCustomEndpoint.ValueString()
	p.OSConfigBasePath = data.OSConfigCustomEndpoint.ValueString()
	p.OSLoginBasePath = data.OSLoginCustomEndpoint.ValueString()
	p.PrivatecaBasePath = data.PrivatecaCustomEndpoint.ValueString()
	p.PubsubBasePath = data.PubsubCustomEndpoint.ValueString()
	p.PubsubLiteBasePath = data.PubsubLiteCustomEndpoint.ValueString()
	p.RedisBasePath = data.RedisCustomEndpoint.ValueString()
	p.ResourceManagerBasePath = data.ResourceManagerCustomEndpoint.ValueString()
	p.SecretManagerBasePath = data.SecretManagerCustomEndpoint.ValueString()
	p.SecurityCenterBasePath = data.SecurityCenterCustomEndpoint.ValueString()
	p.ServiceManagementBasePath = data.ServiceManagementCustomEndpoint.ValueString()
	p.ServiceUsageBasePath = data.ServiceUsageCustomEndpoint.ValueString()
	p.SourceRepoBasePath = data.SourceRepoCustomEndpoint.ValueString()
	p.SpannerBasePath = data.SpannerCustomEndpoint.ValueString()
	p.SQLBasePath = data.SQLCustomEndpoint.ValueString()
	p.StorageBasePath = data.StorageCustomEndpoint.ValueString()
	p.StorageTransferBasePath = data.StorageTransferCustomEndpoint.ValueString()
	p.TagsBasePath = data.TagsCustomEndpoint.ValueString()
	p.TPUBasePath = data.TPUCustomEndpoint.ValueString()
	p.VertexAIBasePath = data.VertexAICustomEndpoint.ValueString()
	p.VPCAccessBasePath = data.VPCAccessCustomEndpoint.ValueString()
	p.WorkflowsBasePath = data.WorkflowsCustomEndpoint.ValueString()

	p.context = ctx
	p.region = data.Region
	p.zone = data.Zone
	p.pollInterval = 10 * time.Second
	p.project = data.Project
	p.requestBatcherServiceUsage = NewRequestBatcher("Service Usage", ctx, batchingConfig)
	p.requestBatcherIam = NewRequestBatcher("IAM", ctx, batchingConfig)
}

// HandleDefaults will handle all the defaults necessary in the provider
func (p *frameworkProvider) HandleDefaults(ctx context.Context, data *ProviderModel, diags *diag.Diagnostics) {
	if data.AccessToken.IsNull() && data.Credentials.IsNull() {
		credentials := MultiEnvDefault([]string{
			"GOOGLE_CREDENTIALS",
			"GOOGLE_CLOUD_KEYFILE_JSON",
			"GCLOUD_KEYFILE_JSON",
		}, nil)

		if credentials != nil {
			data.Credentials = types.StringValue(credentials.(string))
		}

		accessToken := MultiEnvDefault([]string{
			"GOOGLE_OAUTH_ACCESS_TOKEN",
		}, nil)

		if accessToken != nil {
			data.AccessToken = types.StringValue(accessToken.(string))
		}
	}

	if data.ImpersonateServiceAccount.IsNull() && os.Getenv("GOOGLE_IMPERSONATE_SERVICE_ACCOUNT") != "" {
		data.ImpersonateServiceAccount = types.StringValue(os.Getenv("GOOGLE_IMPERSONATE_SERVICE_ACCOUNT"))
	}

	if data.Project.IsNull() {
		project := MultiEnvDefault([]string{
			"GOOGLE_PROJECT",
			"GOOGLE_CLOUD_PROJECT",
			"GCLOUD_PROJECT",
			"CLOUDSDK_CORE_PROJECT",
		}, nil)
		if project != nil {
			data.Project = types.StringValue(project.(string))
		}
	}

	if data.BillingProject.IsNull() && os.Getenv("GOOGLE_BILLING_PROJECT") != "" {
		data.BillingProject = types.StringValue(os.Getenv("GOOGLE_BILLING_PROJECT"))
	}

	if data.Region.IsNull() {
		region := MultiEnvDefault([]string{
			"GOOGLE_REGION",
			"GCLOUD_REGION",
			"CLOUDSDK_COMPUTE_REGION",
		}, nil)

		if region != nil {
			data.Region = types.StringValue(region.(string))
		}
	}

	if data.Zone.IsNull() {
		zone := MultiEnvDefault([]string{
			"GOOGLE_ZONE",
			"GCLOUD_ZONE",
			"CLOUDSDK_COMPUTE_ZONE",
		}, nil)

		if zone != nil {
			data.Zone = types.StringValue(zone.(string))
		}
	}

	if len(data.Scopes.Elements()) == 0 {
		var d diag.Diagnostics
		data.Scopes, d = types.ListValueFrom(ctx, types.StringType, defaultClientScopes)
		diags.Append(d...)
		if diags.HasError() {
			return
		}
	}

	if !data.Batching.IsNull() {
		var pbConfigs []ProviderBatching
		d := data.Batching.ElementsAs(ctx, &pbConfigs, true)
		diags.Append(d...)
		if diags.HasError() {
			return
		}

		if pbConfigs[0].SendAfter.IsNull() {
			pbConfigs[0].SendAfter = types.StringValue("10s")
		}

		if pbConfigs[0].EnableBatching.IsNull() {
			pbConfigs[0].EnableBatching = types.BoolValue(true)
		}

		data.Batching, d = types.ListValueFrom(ctx, types.ObjectType{}.WithAttributeTypes(ProviderBatchingAttributes), pbConfigs)
	}

	if data.UserProjectOverride.IsNull() && os.Getenv("USER_PROJECT_OVERRIDE") != "" {
		override, err := strconv.ParseBool(os.Getenv("USER_PROJECT_OVERRIDE"))
		if err != nil {
			diags.AddError(
				"error parsing environment variable `USER_PROJECT_OVERRIDE` into bool", err.Error())
		}
		data.UserProjectOverride = types.BoolValue(override)
	}

	if data.RequestReason.IsNull() && os.Getenv("CLOUDSDK_CORE_REQUEST_REASON") != "" {
		data.RequestReason = types.StringValue(os.Getenv("CLOUDSDK_CORE_REQUEST_REASON"))
	}

	if data.RequestTimeout.IsNull() {
		data.RequestTimeout = types.StringValue("120s")
	}

	// Generated Products
	if data.AccessApprovalCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_ACCESS_APPROVAL_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[AccessApprovalBasePathKey])
		if customEndpoint != nil {
			data.AccessApprovalCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.AccessContextManagerCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_ACCESS_CONTEXT_MANAGER_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[AccessContextManagerBasePathKey])
		if customEndpoint != nil {
			data.AccessContextManagerCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.ActiveDirectoryCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_ACTIVE_DIRECTORY_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[ActiveDirectoryBasePathKey])
		if customEndpoint != nil {
			data.ActiveDirectoryCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.AlloydbCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_ALLOYDB_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[AlloydbBasePathKey])
		if customEndpoint != nil {
			data.AlloydbCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.ApigeeCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_APIGEE_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[ApigeeBasePathKey])
		if customEndpoint != nil {
			data.ApigeeCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.AppEngineCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_APP_ENGINE_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[AppEngineBasePathKey])
		if customEndpoint != nil {
			data.AppEngineCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.ArtifactRegistryCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_ARTIFACT_REGISTRY_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[ArtifactRegistryBasePathKey])
		if customEndpoint != nil {
			data.ArtifactRegistryCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.BeyondcorpCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_BEYONDCORP_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[BeyondcorpBasePathKey])
		if customEndpoint != nil {
			data.BeyondcorpCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.BigQueryCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_BIG_QUERY_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[BigQueryBasePathKey])
		if customEndpoint != nil {
			data.BigQueryCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.BigqueryAnalyticsHubCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_BIGQUERY_ANALYTICS_HUB_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[BigqueryAnalyticsHubBasePathKey])
		if customEndpoint != nil {
			data.BigqueryAnalyticsHubCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.BigqueryConnectionCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_BIGQUERY_CONNECTION_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[BigqueryConnectionBasePathKey])
		if customEndpoint != nil {
			data.BigqueryConnectionCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.BigqueryDatapolicyCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_BIGQUERY_DATAPOLICY_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[BigqueryDatapolicyBasePathKey])
		if customEndpoint != nil {
			data.BigqueryDatapolicyCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.BigqueryDataTransferCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_BIGQUERY_DATA_TRANSFER_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[BigqueryDataTransferBasePathKey])
		if customEndpoint != nil {
			data.BigqueryDataTransferCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.BigqueryReservationCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_BIGQUERY_RESERVATION_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[BigqueryReservationBasePathKey])
		if customEndpoint != nil {
			data.BigqueryReservationCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.BigtableCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_BIGTABLE_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[BigtableBasePathKey])
		if customEndpoint != nil {
			data.BigtableCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.BillingCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_BILLING_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[BillingBasePathKey])
		if customEndpoint != nil {
			data.BillingCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.BinaryAuthorizationCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_BINARY_AUTHORIZATION_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[BinaryAuthorizationBasePathKey])
		if customEndpoint != nil {
			data.BinaryAuthorizationCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.CertificateManagerCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_CERTIFICATE_MANAGER_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[CertificateManagerBasePathKey])
		if customEndpoint != nil {
			data.CertificateManagerCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.CloudAssetCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_CLOUD_ASSET_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[CloudAssetBasePathKey])
		if customEndpoint != nil {
			data.CloudAssetCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.CloudBuildCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_CLOUD_BUILD_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[CloudBuildBasePathKey])
		if customEndpoint != nil {
			data.CloudBuildCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.CloudFunctionsCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_CLOUD_FUNCTIONS_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[CloudFunctionsBasePathKey])
		if customEndpoint != nil {
			data.CloudFunctionsCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.Cloudfunctions2CustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_CLOUDFUNCTIONS2_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[Cloudfunctions2BasePathKey])
		if customEndpoint != nil {
			data.Cloudfunctions2CustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.CloudIdentityCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_CLOUD_IDENTITY_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[CloudIdentityBasePathKey])
		if customEndpoint != nil {
			data.CloudIdentityCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.CloudIdsCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_CLOUD_IDS_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[CloudIdsBasePathKey])
		if customEndpoint != nil {
			data.CloudIdsCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.CloudIotCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_CLOUD_IOT_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[CloudIotBasePathKey])
		if customEndpoint != nil {
			data.CloudIotCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.CloudRunCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_CLOUD_RUN_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[CloudRunBasePathKey])
		if customEndpoint != nil {
			data.CloudRunCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.CloudRunV2CustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_CLOUD_RUN_V2_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[CloudRunV2BasePathKey])
		if customEndpoint != nil {
			data.CloudRunV2CustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.CloudSchedulerCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_CLOUD_SCHEDULER_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[CloudSchedulerBasePathKey])
		if customEndpoint != nil {
			data.CloudSchedulerCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.CloudTasksCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_CLOUD_TASKS_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[CloudTasksBasePathKey])
		if customEndpoint != nil {
			data.CloudTasksCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.ComputeCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_COMPUTE_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[ComputeBasePathKey])
		if customEndpoint != nil {
			data.ComputeCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.ContainerAnalysisCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_CONTAINER_ANALYSIS_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[ContainerAnalysisBasePathKey])
		if customEndpoint != nil {
			data.ContainerAnalysisCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.ContainerAttachedCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_CONTAINER_ATTACHED_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[ContainerAttachedBasePathKey])
		if customEndpoint != nil {
			data.ContainerAttachedCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.DataCatalogCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_DATA_CATALOG_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[DataCatalogBasePathKey])
		if customEndpoint != nil {
			data.DataCatalogCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.DataFusionCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_DATA_FUSION_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[DataFusionBasePathKey])
		if customEndpoint != nil {
			data.DataFusionCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.DataLossPreventionCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_DATA_LOSS_PREVENTION_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[DataLossPreventionBasePathKey])
		if customEndpoint != nil {
			data.DataLossPreventionCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.DataplexCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_DATAPLEX_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[DataplexBasePathKey])
		if customEndpoint != nil {
			data.DataplexCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.DataprocCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_DATAPROC_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[DataprocBasePathKey])
		if customEndpoint != nil {
			data.DataprocCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.DataprocMetastoreCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_DATAPROC_METASTORE_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[DataprocMetastoreBasePathKey])
		if customEndpoint != nil {
			data.DataprocMetastoreCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.DatastoreCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_DATASTORE_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[DatastoreBasePathKey])
		if customEndpoint != nil {
			data.DatastoreCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.DatastreamCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_DATASTREAM_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[DatastreamBasePathKey])
		if customEndpoint != nil {
			data.DatastreamCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.DeploymentManagerCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_DEPLOYMENT_MANAGER_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[DeploymentManagerBasePathKey])
		if customEndpoint != nil {
			data.DeploymentManagerCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.DialogflowCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_DIALOGFLOW_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[DialogflowBasePathKey])
		if customEndpoint != nil {
			data.DialogflowCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.DialogflowCXCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_DIALOGFLOW_CX_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[DialogflowCXBasePathKey])
		if customEndpoint != nil {
			data.DialogflowCXCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.DNSCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_DNS_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[DNSBasePathKey])
		if customEndpoint != nil {
			data.DNSCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.DocumentAICustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_DOCUMENT_AI_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[DocumentAIBasePathKey])
		if customEndpoint != nil {
			data.DocumentAICustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.EssentialContactsCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_ESSENTIAL_CONTACTS_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[EssentialContactsBasePathKey])
		if customEndpoint != nil {
			data.EssentialContactsCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.FilestoreCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_FILESTORE_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[FilestoreBasePathKey])
		if customEndpoint != nil {
			data.FilestoreCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.FirestoreCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_FIRESTORE_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[FirestoreBasePathKey])
		if customEndpoint != nil {
			data.FirestoreCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.GameServicesCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_GAME_SERVICES_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[GameServicesBasePathKey])
		if customEndpoint != nil {
			data.GameServicesCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.GKEBackupCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_GKE_BACKUP_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[GKEBackupBasePathKey])
		if customEndpoint != nil {
			data.GKEBackupCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.GKEHubCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_GKE_HUB_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[GKEHubBasePathKey])
		if customEndpoint != nil {
			data.GKEHubCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.HealthcareCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_HEALTHCARE_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[HealthcareBasePathKey])
		if customEndpoint != nil {
			data.HealthcareCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.IAM2CustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_IAM2_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[IAM2BasePathKey])
		if customEndpoint != nil {
			data.IAM2CustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.IAMBetaCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_IAM_BETA_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[IAMBetaBasePathKey])
		if customEndpoint != nil {
			data.IAMBetaCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.IAMWorkforcePoolCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_IAM_WORKFORCE_POOL_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[IAMWorkforcePoolBasePathKey])
		if customEndpoint != nil {
			data.IAMWorkforcePoolCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.IapCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_IAP_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[IapBasePathKey])
		if customEndpoint != nil {
			data.IapCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.IdentityPlatformCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_IDENTITY_PLATFORM_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[IdentityPlatformBasePathKey])
		if customEndpoint != nil {
			data.IdentityPlatformCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.KMSCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_KMS_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[KMSBasePathKey])
		if customEndpoint != nil {
			data.KMSCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.LoggingCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_LOGGING_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[LoggingBasePathKey])
		if customEndpoint != nil {
			data.LoggingCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.MemcacheCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_MEMCACHE_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[MemcacheBasePathKey])
		if customEndpoint != nil {
			data.MemcacheCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.MLEngineCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_ML_ENGINE_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[MLEngineBasePathKey])
		if customEndpoint != nil {
			data.MLEngineCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.MonitoringCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_MONITORING_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[MonitoringBasePathKey])
		if customEndpoint != nil {
			data.MonitoringCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.NetworkManagementCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_NETWORK_MANAGEMENT_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[NetworkManagementBasePathKey])
		if customEndpoint != nil {
			data.NetworkManagementCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.NetworkServicesCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_NETWORK_SERVICES_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[NetworkServicesBasePathKey])
		if customEndpoint != nil {
			data.NetworkServicesCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.NotebooksCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_NOTEBOOKS_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[NotebooksBasePathKey])
		if customEndpoint != nil {
			data.NotebooksCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.OSConfigCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_OS_CONFIG_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[OSConfigBasePathKey])
		if customEndpoint != nil {
			data.OSConfigCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.OSLoginCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_OS_LOGIN_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[OSLoginBasePathKey])
		if customEndpoint != nil {
			data.OSLoginCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.PrivatecaCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_PRIVATECA_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[PrivatecaBasePathKey])
		if customEndpoint != nil {
			data.PrivatecaCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.PubsubCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_PUBSUB_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[PubsubBasePathKey])
		if customEndpoint != nil {
			data.PubsubCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.PubsubLiteCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_PUBSUB_LITE_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[PubsubLiteBasePathKey])
		if customEndpoint != nil {
			data.PubsubLiteCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.RedisCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_REDIS_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[RedisBasePathKey])
		if customEndpoint != nil {
			data.RedisCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.ResourceManagerCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_RESOURCE_MANAGER_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[ResourceManagerBasePathKey])
		if customEndpoint != nil {
			data.ResourceManagerCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.SecretManagerCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_SECRET_MANAGER_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[SecretManagerBasePathKey])
		if customEndpoint != nil {
			data.SecretManagerCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.SecurityCenterCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_SECURITY_CENTER_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[SecurityCenterBasePathKey])
		if customEndpoint != nil {
			data.SecurityCenterCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.ServiceManagementCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_SERVICE_MANAGEMENT_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[ServiceManagementBasePathKey])
		if customEndpoint != nil {
			data.ServiceManagementCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.ServiceUsageCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_SERVICE_USAGE_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[ServiceUsageBasePathKey])
		if customEndpoint != nil {
			data.ServiceUsageCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.SourceRepoCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_SOURCE_REPO_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[SourceRepoBasePathKey])
		if customEndpoint != nil {
			data.SourceRepoCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.SpannerCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_SPANNER_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[SpannerBasePathKey])
		if customEndpoint != nil {
			data.SpannerCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.SQLCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_SQL_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[SQLBasePathKey])
		if customEndpoint != nil {
			data.SQLCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.StorageCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_STORAGE_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[StorageBasePathKey])
		if customEndpoint != nil {
			data.StorageCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.StorageTransferCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_STORAGE_TRANSFER_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[StorageTransferBasePathKey])
		if customEndpoint != nil {
			data.StorageTransferCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.TagsCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_TAGS_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[TagsBasePathKey])
		if customEndpoint != nil {
			data.TagsCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.TPUCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_TPU_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[TPUBasePathKey])
		if customEndpoint != nil {
			data.TPUCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.VertexAICustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_VERTEX_AI_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[VertexAIBasePathKey])
		if customEndpoint != nil {
			data.VertexAICustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.VPCAccessCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_VPC_ACCESS_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[VPCAccessBasePathKey])
		if customEndpoint != nil {
			data.VPCAccessCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.WorkflowsCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_WORKFLOWS_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[WorkflowsBasePathKey])
		if customEndpoint != nil {
			data.WorkflowsCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}

	// Handwritten Products / Versioned / Atypical Entries
	if data.CloudBillingCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_CLOUD_BILLING_CUSTOM_ENDPOINT",
		}, DefaultBasePaths["cloud_billing_custom_endpoint"])
		if customEndpoint != nil {
			data.CloudBillingCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}

	if data.ComposerCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_COMPOSER_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[ComposerBasePathKey])
		if customEndpoint != nil {
			data.ComposerCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}

	if data.ContainerCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_CONTAINER_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[ContainerBasePathKey])
		if customEndpoint != nil {
			data.ContainerCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}

	if data.DataflowCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_DATAFLOW_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[DataflowBasePathKey])
		if customEndpoint != nil {
			data.DataflowCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}

	if data.IamCredentialsCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_IAM_CREDENTIALS_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[IamCredentialsBasePathKey])
		if customEndpoint != nil {
			data.IamCredentialsCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}

	if data.ResourceManagerV3CustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_RESOURCE_MANAGER_V3_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[ResourceManagerV3BasePathKey])
		if customEndpoint != nil {
			data.ResourceManagerV3CustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}

	if data.IAMCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_IAM_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[IAMBasePathKey])
		if customEndpoint != nil {
			data.IAMCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}

	if data.ServiceNetworkingCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_SERVICE_NETWORKING_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[ServiceNetworkingBasePathKey])
		if customEndpoint != nil {
			data.ServiceNetworkingCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}

	if data.TagsLocationCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_TAGS_LOCATION_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[TagsLocationBasePathKey])
		if customEndpoint != nil {
			data.TagsLocationCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}

	// dcl
	if data.ContainerAwsCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_CONTAINERAWS_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[ContainerAwsBasePathKey])
		if customEndpoint != nil {
			data.ContainerAwsCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}

	if data.ContainerAzureCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_CONTAINERAZURE_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[ContainerAzureBasePathKey])
		if customEndpoint != nil {
			data.ContainerAzureCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}

	// DCL generated defaults
	if data.ApikeysCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_APIKEYS_CUSTOM_ENDPOINT",
		}, "")
		if customEndpoint != nil {
			data.ApikeysCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}

	if data.AssuredWorkloadsCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_ASSURED_WORKLOADS_CUSTOM_ENDPOINT",
		}, "")
		if customEndpoint != nil {
			data.AssuredWorkloadsCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}

	if data.CloudBuildWorkerPoolCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_CLOUD_BUILD_WORKER_POOL_CUSTOM_ENDPOINT",
		}, "")
		if customEndpoint != nil {
			data.CloudBuildWorkerPoolCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}

	if data.CloudDeployCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_CLOUDDEPLOY_CUSTOM_ENDPOINT",
		}, "")
		if customEndpoint != nil {
			data.CloudDeployCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}

	if data.CloudResourceManagerCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_CLOUD_RESOURCE_MANAGER_CUSTOM_ENDPOINT",
		}, "")
		if customEndpoint != nil {
			data.CloudResourceManagerCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}

	if data.DataplexCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_DATAPLEX_CUSTOM_ENDPOINT",
		}, "")
		if customEndpoint != nil {
			data.DataplexCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}

	if data.EventarcCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_EVENTARC_CUSTOM_ENDPOINT",
		}, "")
		if customEndpoint != nil {
			data.EventarcCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}

	if data.FirebaserulesCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_FIREBASERULES_CUSTOM_ENDPOINT",
		}, "")
		if customEndpoint != nil {
			data.FirebaserulesCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}

	if data.NetworkConnectivityCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_NETWORK_CONNECTIVITY_CUSTOM_ENDPOINT",
		}, "")
		if customEndpoint != nil {
			data.NetworkConnectivityCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}

	if data.RecaptchaEnterpriseCustomEndpoint.IsNull() {
		customEndpoint := MultiEnvDefault([]string{
			"GOOGLE_RECAPTCHA_ENTERPRISE_CUSTOM_ENDPOINT",
		}, "")
		if customEndpoint != nil {
			data.RecaptchaEnterpriseCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
}

func (p *frameworkProvider) SetupClient(ctx context.Context, data ProviderModel, diags *diag.Diagnostics) {
	tokenSource := GetTokenSource(ctx, data, false, diags)
	if diags.HasError() {
		return
	}

	cleanCtx := context.WithValue(ctx, oauth2.HTTPClient, cleanhttp.DefaultClient())

	// 1. MTLS TRANSPORT/CLIENT - sets up proper auth headers
	client, _, err := transport.NewHTTPClient(cleanCtx, option.WithTokenSource(tokenSource))
	if err != nil {
		diags.AddError("error creating new http client", err.Error())
		return
	}

	// Userinfo is fetched before request logging is enabled to reduce additional noise.
	p.logGoogleIdentities(ctx, data, diags)
	if diags.HasError() {
		return
	}

	// 2. Logging Transport - ensure we log HTTP requests to GCP APIs.
	loggingTransport := logging.NewTransport("Google", client.Transport)

	// 3. Retry Transport - retries common temporary errors
	// Keep order for wrapping logging so we log each retried request as well.
	// This value should be used if needed to create shallow copies with additional retry predicates.
	// See ClientWithAdditionalRetries
	retryTransport := NewTransportWithDefaultRetries(loggingTransport)

	// 4. Header Transport - outer wrapper to inject additional headers we want to apply
	// before making requests
	headerTransport := newTransportWithHeaders(retryTransport)
	if !data.RequestReason.IsNull() {
		headerTransport.Set("X-Goog-Request-Reason", data.RequestReason.ValueString())
	}

	// Ensure $userProject is set for all HTTP requests using the client if specified by the provider config
	// See https://cloud.google.com/apis/docs/system-parameters
	if data.UserProjectOverride.ValueBool() && !data.BillingProject.IsNull() {
		headerTransport.Set("X-Goog-User-Project", data.BillingProject.ValueString())
	}

	// Set final transport value.
	client.Transport = headerTransport

	// This timeout is a timeout per HTTP request, not per logical operation.
	timeout, err := time.ParseDuration(data.RequestTimeout.ValueString())
	if err != nil {
		diags.AddError("error parsing request timeout", err.Error())
	}
	client.Timeout = timeout

	p.tokenSource = tokenSource
	p.client = client
}

func (p *frameworkProvider) SetupGrpcLogging() {
	logger := logrus.StandardLogger()

	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&Formatter{
		TimestampFormat: "2006/01/02 15:04:05",
		LogFormat:       "%time% [%lvl%] %msg% \n",
	})

	alwaysLoggingDeciderClient := func(ctx context.Context, fullMethodName string) bool { return true }
	grpc_logrus.ReplaceGrpcLogger(logrus.NewEntry(logger))

	p.gRPCLoggingOptions = append(
		p.gRPCLoggingOptions, option.WithGRPCDialOption(grpc.WithUnaryInterceptor(
			grpc_logrus.PayloadUnaryClientInterceptor(logrus.NewEntry(logger), alwaysLoggingDeciderClient))),
		option.WithGRPCDialOption(grpc.WithStreamInterceptor(
			grpc_logrus.PayloadStreamClientInterceptor(logrus.NewEntry(logger), alwaysLoggingDeciderClient))),
	)
}

func (p *frameworkProvider) logGoogleIdentities(ctx context.Context, data ProviderModel, diags *diag.Diagnostics) {
	// GetCurrentUserEmailFramework doesn't pass an error back from logGoogleIdentities, so we want
	// a separate diagnostics here
	var d diag.Diagnostics

	if data.ImpersonateServiceAccount.IsNull() {

		tokenSource := GetTokenSource(ctx, data, true, diags)
		if diags.HasError() {
			return
		}

		p.client = oauth2.NewClient(ctx, tokenSource) // p.client isn't initialised fully when this code is called.

		email := GetCurrentUserEmailFramework(p, p.userAgent, &d)
		if d.HasError() {
			tflog.Info(ctx, "error retrieving userinfo for your provider credentials. have you enabled the 'https://www.googleapis.com/auth/userinfo.email' scope?")
		}

		tflog.Info(ctx, fmt.Sprintf("Terraform is using this identity: %s", email))
		return
	}

	// Drop Impersonated ClientOption from OAuth2 TokenSource to infer original identity
	tokenSource := GetTokenSource(ctx, data, true, diags)
	if diags.HasError() {
		return
	}

	p.client = oauth2.NewClient(ctx, tokenSource) // p.client isn't initialised fully when this code is called.
	email := GetCurrentUserEmailFramework(p, p.userAgent, &d)
	if d.HasError() {
		tflog.Info(ctx, "error retrieving userinfo for your provider credentials. have you enabled the 'https://www.googleapis.com/auth/userinfo.email' scope?")
	}

	tflog.Info(ctx, fmt.Sprintf("Terraform is configured with service account impersonation, original identity: %s, impersonated identity: %s", email, data.ImpersonateServiceAccount.ValueString()))

	// Add the Impersonated ClientOption back in to the OAuth2 TokenSource
	tokenSource = GetTokenSource(ctx, data, false, diags)
	if diags.HasError() {
		return
	}

	p.client = oauth2.NewClient(ctx, tokenSource) // p.client isn't initialised fully when this code is called.

	return
}

// Configuration helpers

// GetTokenSource gets token source based on the Google Credentials configured.
// If initialCredentialsOnly is true, don't follow the impersonation settings and return the initial set of creds.
func GetTokenSource(ctx context.Context, data ProviderModel, initialCredentialsOnly bool, diags *diag.Diagnostics) oauth2.TokenSource {
	creds := GetCredentials(ctx, data, initialCredentialsOnly, diags)

	return creds.TokenSource
}

// GetCredentials gets credentials with a given scope (clientScopes).
// If initialCredentialsOnly is true, don't follow the impersonation
// settings and return the initial set of creds instead.
func GetCredentials(ctx context.Context, data ProviderModel, initialCredentialsOnly bool, diags *diag.Diagnostics) googleoauth.Credentials {
	var clientScopes []string
	var delegates []string

	d := data.Scopes.ElementsAs(ctx, &clientScopes, false)
	diags.Append(d...)
	if diags.HasError() {
		return googleoauth.Credentials{}
	}

	d = data.ImpersonateServiceAccountDelegates.ElementsAs(ctx, &delegates, false)
	diags.Append(d...)
	if diags.HasError() {
		return googleoauth.Credentials{}
	}

	if !data.AccessToken.IsNull() {
		contents, _, err := pathOrContents(data.AccessToken.ValueString())
		if err != nil {
			diags.AddError("error loading access token", err.Error())
			return googleoauth.Credentials{}
		}

		token := &oauth2.Token{AccessToken: contents}
		if !data.ImpersonateServiceAccount.IsNull() && !initialCredentialsOnly {
			opts := []option.ClientOption{option.WithTokenSource(oauth2.StaticTokenSource(token)), option.ImpersonateCredentials(data.ImpersonateServiceAccount.ValueString(), delegates...), option.WithScopes(clientScopes...)}
			creds, err := transport.Creds(context.TODO(), opts...)
			if err != nil {
				diags.AddError("error impersonating credentials", err.Error())
				return googleoauth.Credentials{}
			}
			return *creds
		}

		tflog.Info(ctx, "Authenticating using configured Google JSON 'access_token'...")
		tflog.Info(ctx, fmt.Sprintf("  -- Scopes: %s", clientScopes))
		return googleoauth.Credentials{
			TokenSource: staticTokenSource{oauth2.StaticTokenSource(token)},
		}
	}

	if !data.Credentials.IsNull() {
		contents, _, err := pathOrContents(data.Credentials.ValueString())
		if err != nil {
			diags.AddError(fmt.Sprintf("error loading credentials: %s", err), err.Error())
			return googleoauth.Credentials{}
		}

		if !data.ImpersonateServiceAccount.IsNull() && !initialCredentialsOnly {
			opts := []option.ClientOption{option.WithCredentialsJSON([]byte(contents)), option.ImpersonateCredentials(data.ImpersonateServiceAccount.ValueString(), delegates...), option.WithScopes(clientScopes...)}
			creds, err := transport.Creds(context.TODO(), opts...)
			if err != nil {
				diags.AddError("error impersonating credentials", err.Error())
				return googleoauth.Credentials{}
			}
			return *creds
		}

		creds, err := googleoauth.CredentialsFromJSON(ctx, []byte(contents), clientScopes...)
		if err != nil {
			diags.AddError("unable to parse credentials", err.Error())
			return googleoauth.Credentials{}
		}

		tflog.Info(ctx, "Authenticating using configured Google JSON 'credentials'...")
		tflog.Info(ctx, fmt.Sprintf("  -- Scopes: %s", clientScopes))
		return *creds
	}

	if !data.ImpersonateServiceAccount.IsNull() && !initialCredentialsOnly {
		opts := option.ImpersonateCredentials(data.ImpersonateServiceAccount.ValueString(), delegates...)
		creds, err := transport.Creds(context.TODO(), opts, option.WithScopes(clientScopes...))
		if err != nil {
			diags.AddError("error impersonating credentials", err.Error())
			return googleoauth.Credentials{}
		}

		return *creds
	}

	tflog.Info(ctx, "Authenticating using DefaultClient...")
	tflog.Info(ctx, fmt.Sprintf("  -- Scopes: %s", clientScopes))
	defaultTS, err := googleoauth.DefaultTokenSource(context.Background(), clientScopes...)
	if err != nil {
		diags.AddError(fmt.Sprintf("Attempted to load application default credentials since neither `credentials` nor `access_token` was set in the provider block.  "+
			"No credentials loaded. To use your gcloud credentials, run 'gcloud auth application-default login'"), err.Error())
		return googleoauth.Credentials{}
	}

	return googleoauth.Credentials{
		TokenSource: defaultTS,
	}
}

// GetBatchingConfig returns the batching config object given the
// provider configuration set for batching
func GetBatchingConfig(ctx context.Context, data types.List, diags *diag.Diagnostics) *batchingConfig {
	bc := &batchingConfig{
		SendAfter:      time.Second * DefaultBatchSendIntervalSec,
		EnableBatching: true,
	}

	if data.IsNull() {
		return bc
	}

	var pbConfigs []ProviderBatching
	d := data.ElementsAs(ctx, &pbConfigs, true)
	diags.Append(d...)
	if diags.HasError() {
		return bc
	}

	sendAfter, err := time.ParseDuration(pbConfigs[0].SendAfter.ValueString())
	if err != nil {
		diags.AddError("error parsing send after time duration", err.Error())
		return bc
	}

	bc.SendAfter = sendAfter

	if !pbConfigs[0].EnableBatching.IsNull() {
		bc.EnableBatching = pbConfigs[0].EnableBatching.ValueBool()
	}

	return bc
}
