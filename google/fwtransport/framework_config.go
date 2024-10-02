// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package fwtransport

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"
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
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"

	"github.com/hashicorp/terraform-provider-google/google/fwmodels"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/google/verify"

	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/sirupsen/logrus"
)

type FrameworkProviderConfig struct {
	// Temporary, as we'll replace use of FrameworkProviderConfig with transport_tpg.Config soon
	// transport_tpg.Config has a the fields below, hence these changes are needed
	Credentials                               types.String
	AccessToken                               types.String
	ImpersonateServiceAccount                 types.String
	ImpersonateServiceAccountDelegates        types.List
	RequestReason                             types.String
	AddTerraformAttributionLabel              types.Bool
	TerraformAttributionLabelAdditionStrategy types.String
	// End temporary

	BillingProject             types.String
	Client                     *http.Client
	Context                    context.Context
	gRPCLoggingOptions         []option.ClientOption
	PollInterval               time.Duration
	Project                    types.String
	Region                     types.String
	Zone                       types.String
	RequestBatcherIam          *transport_tpg.RequestBatcher
	RequestBatcherServiceUsage *transport_tpg.RequestBatcher
	Scopes                     types.List
	TokenSource                oauth2.TokenSource
	UniverseDomain             types.String
	UserAgent                  string
	UserProjectOverride        types.Bool
	DefaultLabels              types.Map

	// paths for client setup
	AccessApprovalBasePath           string
	AccessContextManagerBasePath     string
	ActiveDirectoryBasePath          string
	AlloydbBasePath                  string
	ApigeeBasePath                   string
	AppEngineBasePath                string
	ApphubBasePath                   string
	ArtifactRegistryBasePath         string
	BeyondcorpBasePath               string
	BiglakeBasePath                  string
	BigQueryBasePath                 string
	BigqueryAnalyticsHubBasePath     string
	BigqueryConnectionBasePath       string
	BigqueryDatapolicyBasePath       string
	BigqueryDataTransferBasePath     string
	BigqueryReservationBasePath      string
	BigtableBasePath                 string
	BillingBasePath                  string
	BinaryAuthorizationBasePath      string
	BlockchainNodeEngineBasePath     string
	CertificateManagerBasePath       string
	CloudAssetBasePath               string
	CloudBuildBasePath               string
	Cloudbuildv2BasePath             string
	ClouddeployBasePath              string
	ClouddomainsBasePath             string
	CloudFunctionsBasePath           string
	Cloudfunctions2BasePath          string
	CloudIdentityBasePath            string
	CloudIdsBasePath                 string
	CloudQuotasBasePath              string
	CloudRunBasePath                 string
	CloudRunV2BasePath               string
	CloudSchedulerBasePath           string
	CloudTasksBasePath               string
	ComposerBasePath                 string
	ComputeBasePath                  string
	ContainerAnalysisBasePath        string
	ContainerAttachedBasePath        string
	CoreBillingBasePath              string
	DatabaseMigrationServiceBasePath string
	DataCatalogBasePath              string
	DataFusionBasePath               string
	DataLossPreventionBasePath       string
	DataPipelineBasePath             string
	DataplexBasePath                 string
	DataprocBasePath                 string
	DataprocMetastoreBasePath        string
	DatastreamBasePath               string
	DeploymentManagerBasePath        string
	DialogflowBasePath               string
	DialogflowCXBasePath             string
	DiscoveryEngineBasePath          string
	DNSBasePath                      string
	DocumentAIBasePath               string
	DocumentAIWarehouseBasePath      string
	EdgecontainerBasePath            string
	EdgenetworkBasePath              string
	EssentialContactsBasePath        string
	FilestoreBasePath                string
	FirebaseAppCheckBasePath         string
	FirestoreBasePath                string
	GKEBackupBasePath                string
	GKEHubBasePath                   string
	GKEHub2BasePath                  string
	GkeonpremBasePath                string
	HealthcareBasePath               string
	IAM2BasePath                     string
	IAMBetaBasePath                  string
	IAMWorkforcePoolBasePath         string
	IapBasePath                      string
	IdentityPlatformBasePath         string
	IntegrationConnectorsBasePath    string
	IntegrationsBasePath             string
	KMSBasePath                      string
	LoggingBasePath                  string
	LookerBasePath                   string
	MemcacheBasePath                 string
	MigrationCenterBasePath          string
	MLEngineBasePath                 string
	MonitoringBasePath               string
	NetappBasePath                   string
	NetworkConnectivityBasePath      string
	NetworkManagementBasePath        string
	NetworkSecurityBasePath          string
	NetworkServicesBasePath          string
	NotebooksBasePath                string
	OrgPolicyBasePath                string
	OSConfigBasePath                 string
	OSLoginBasePath                  string
	PrivatecaBasePath                string
	PrivilegedAccessManagerBasePath  string
	PublicCABasePath                 string
	PubsubBasePath                   string
	PubsubLiteBasePath               string
	RedisBasePath                    string
	ResourceManagerBasePath          string
	SecretManagerBasePath            string
	SecretManagerRegionalBasePath    string
	SecureSourceManagerBasePath      string
	SecurityCenterBasePath           string
	SecurityCenterManagementBasePath string
	SecurityCenterV2BasePath         string
	SecuritypostureBasePath          string
	ServiceManagementBasePath        string
	ServiceNetworkingBasePath        string
	ServiceUsageBasePath             string
	SiteVerificationBasePath         string
	SourceRepoBasePath               string
	SpannerBasePath                  string
	SQLBasePath                      string
	StorageBasePath                  string
	StorageInsightsBasePath          string
	StorageTransferBasePath          string
	TagsBasePath                     string
	TPUBasePath                      string
	VertexAIBasePath                 string
	VmwareengineBasePath             string
	VPCAccessBasePath                string
	WorkbenchBasePath                string
	WorkflowsBasePath                string
}

// LoadAndValidateFramework handles the bulk of configuring the provider
// it is pulled out so that we can manually call this from our testing provider as well
func (p *FrameworkProviderConfig) LoadAndValidateFramework(ctx context.Context, data *fwmodels.ProviderModel, tfVersion string, diags *diag.Diagnostics, providerversion string) {

	// Set defaults if needed
	p.HandleDefaults(ctx, data, diags)
	if diags.HasError() {
		return
	}

	p.Context = ctx

	// Handle User Agent string
	p.UserAgent = CompileUserAgentString(ctx, "terraform-provider-google", tfVersion, providerversion)
	// opt in extension for adding to the User-Agent header
	if ext := os.Getenv("GOOGLE_TERRAFORM_USERAGENT_EXTENSION"); ext != "" {
		ua := p.UserAgent
		p.UserAgent = fmt.Sprintf("%s %s", ua, ext)
	}

	// Set up client configuration
	p.SetupClient(ctx, *data, diags)
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
	p.ApphubBasePath = data.ApphubCustomEndpoint.ValueString()
	p.ArtifactRegistryBasePath = data.ArtifactRegistryCustomEndpoint.ValueString()
	p.BeyondcorpBasePath = data.BeyondcorpCustomEndpoint.ValueString()
	p.BiglakeBasePath = data.BiglakeCustomEndpoint.ValueString()
	p.BigQueryBasePath = data.BigQueryCustomEndpoint.ValueString()
	p.BigqueryAnalyticsHubBasePath = data.BigqueryAnalyticsHubCustomEndpoint.ValueString()
	p.BigqueryConnectionBasePath = data.BigqueryConnectionCustomEndpoint.ValueString()
	p.BigqueryDatapolicyBasePath = data.BigqueryDatapolicyCustomEndpoint.ValueString()
	p.BigqueryDataTransferBasePath = data.BigqueryDataTransferCustomEndpoint.ValueString()
	p.BigqueryReservationBasePath = data.BigqueryReservationCustomEndpoint.ValueString()
	p.BigtableBasePath = data.BigtableCustomEndpoint.ValueString()
	p.BillingBasePath = data.BillingCustomEndpoint.ValueString()
	p.BinaryAuthorizationBasePath = data.BinaryAuthorizationCustomEndpoint.ValueString()
	p.BlockchainNodeEngineBasePath = data.BlockchainNodeEngineCustomEndpoint.ValueString()
	p.CertificateManagerBasePath = data.CertificateManagerCustomEndpoint.ValueString()
	p.CloudAssetBasePath = data.CloudAssetCustomEndpoint.ValueString()
	p.CloudBuildBasePath = data.CloudBuildCustomEndpoint.ValueString()
	p.Cloudbuildv2BasePath = data.Cloudbuildv2CustomEndpoint.ValueString()
	p.ClouddeployBasePath = data.ClouddeployCustomEndpoint.ValueString()
	p.ClouddomainsBasePath = data.ClouddomainsCustomEndpoint.ValueString()
	p.CloudFunctionsBasePath = data.CloudFunctionsCustomEndpoint.ValueString()
	p.Cloudfunctions2BasePath = data.Cloudfunctions2CustomEndpoint.ValueString()
	p.CloudIdentityBasePath = data.CloudIdentityCustomEndpoint.ValueString()
	p.CloudIdsBasePath = data.CloudIdsCustomEndpoint.ValueString()
	p.CloudQuotasBasePath = data.CloudQuotasCustomEndpoint.ValueString()
	p.CloudRunBasePath = data.CloudRunCustomEndpoint.ValueString()
	p.CloudRunV2BasePath = data.CloudRunV2CustomEndpoint.ValueString()
	p.CloudSchedulerBasePath = data.CloudSchedulerCustomEndpoint.ValueString()
	p.CloudTasksBasePath = data.CloudTasksCustomEndpoint.ValueString()
	p.ComposerBasePath = data.ComposerCustomEndpoint.ValueString()
	p.ComputeBasePath = data.ComputeCustomEndpoint.ValueString()
	p.ContainerAnalysisBasePath = data.ContainerAnalysisCustomEndpoint.ValueString()
	p.ContainerAttachedBasePath = data.ContainerAttachedCustomEndpoint.ValueString()
	p.CoreBillingBasePath = data.CoreBillingCustomEndpoint.ValueString()
	p.DatabaseMigrationServiceBasePath = data.DatabaseMigrationServiceCustomEndpoint.ValueString()
	p.DataCatalogBasePath = data.DataCatalogCustomEndpoint.ValueString()
	p.DataFusionBasePath = data.DataFusionCustomEndpoint.ValueString()
	p.DataLossPreventionBasePath = data.DataLossPreventionCustomEndpoint.ValueString()
	p.DataPipelineBasePath = data.DataPipelineCustomEndpoint.ValueString()
	p.DataplexBasePath = data.DataplexCustomEndpoint.ValueString()
	p.DataprocBasePath = data.DataprocCustomEndpoint.ValueString()
	p.DataprocMetastoreBasePath = data.DataprocMetastoreCustomEndpoint.ValueString()
	p.DatastreamBasePath = data.DatastreamCustomEndpoint.ValueString()
	p.DeploymentManagerBasePath = data.DeploymentManagerCustomEndpoint.ValueString()
	p.DialogflowBasePath = data.DialogflowCustomEndpoint.ValueString()
	p.DialogflowCXBasePath = data.DialogflowCXCustomEndpoint.ValueString()
	p.DiscoveryEngineBasePath = data.DiscoveryEngineCustomEndpoint.ValueString()
	p.DNSBasePath = data.DNSCustomEndpoint.ValueString()
	p.DocumentAIBasePath = data.DocumentAICustomEndpoint.ValueString()
	p.DocumentAIWarehouseBasePath = data.DocumentAIWarehouseCustomEndpoint.ValueString()
	p.EdgecontainerBasePath = data.EdgecontainerCustomEndpoint.ValueString()
	p.EdgenetworkBasePath = data.EdgenetworkCustomEndpoint.ValueString()
	p.EssentialContactsBasePath = data.EssentialContactsCustomEndpoint.ValueString()
	p.FilestoreBasePath = data.FilestoreCustomEndpoint.ValueString()
	p.FirebaseAppCheckBasePath = data.FirebaseAppCheckCustomEndpoint.ValueString()
	p.FirestoreBasePath = data.FirestoreCustomEndpoint.ValueString()
	p.GKEBackupBasePath = data.GKEBackupCustomEndpoint.ValueString()
	p.GKEHubBasePath = data.GKEHubCustomEndpoint.ValueString()
	p.GKEHub2BasePath = data.GKEHub2CustomEndpoint.ValueString()
	p.GkeonpremBasePath = data.GkeonpremCustomEndpoint.ValueString()
	p.HealthcareBasePath = data.HealthcareCustomEndpoint.ValueString()
	p.IAM2BasePath = data.IAM2CustomEndpoint.ValueString()
	p.IAMBetaBasePath = data.IAMBetaCustomEndpoint.ValueString()
	p.IAMWorkforcePoolBasePath = data.IAMWorkforcePoolCustomEndpoint.ValueString()
	p.IapBasePath = data.IapCustomEndpoint.ValueString()
	p.IdentityPlatformBasePath = data.IdentityPlatformCustomEndpoint.ValueString()
	p.IntegrationConnectorsBasePath = data.IntegrationConnectorsCustomEndpoint.ValueString()
	p.IntegrationsBasePath = data.IntegrationsCustomEndpoint.ValueString()
	p.KMSBasePath = data.KMSCustomEndpoint.ValueString()
	p.LoggingBasePath = data.LoggingCustomEndpoint.ValueString()
	p.LookerBasePath = data.LookerCustomEndpoint.ValueString()
	p.MemcacheBasePath = data.MemcacheCustomEndpoint.ValueString()
	p.MigrationCenterBasePath = data.MigrationCenterCustomEndpoint.ValueString()
	p.MLEngineBasePath = data.MLEngineCustomEndpoint.ValueString()
	p.MonitoringBasePath = data.MonitoringCustomEndpoint.ValueString()
	p.NetappBasePath = data.NetappCustomEndpoint.ValueString()
	p.NetworkConnectivityBasePath = data.NetworkConnectivityCustomEndpoint.ValueString()
	p.NetworkManagementBasePath = data.NetworkManagementCustomEndpoint.ValueString()
	p.NetworkSecurityBasePath = data.NetworkSecurityCustomEndpoint.ValueString()
	p.NetworkServicesBasePath = data.NetworkServicesCustomEndpoint.ValueString()
	p.NotebooksBasePath = data.NotebooksCustomEndpoint.ValueString()
	p.OrgPolicyBasePath = data.OrgPolicyCustomEndpoint.ValueString()
	p.OSConfigBasePath = data.OSConfigCustomEndpoint.ValueString()
	p.OSLoginBasePath = data.OSLoginCustomEndpoint.ValueString()
	p.PrivatecaBasePath = data.PrivatecaCustomEndpoint.ValueString()
	p.PrivilegedAccessManagerBasePath = data.PrivilegedAccessManagerCustomEndpoint.ValueString()
	p.PublicCABasePath = data.PublicCACustomEndpoint.ValueString()
	p.PubsubBasePath = data.PubsubCustomEndpoint.ValueString()
	p.PubsubLiteBasePath = data.PubsubLiteCustomEndpoint.ValueString()
	p.RedisBasePath = data.RedisCustomEndpoint.ValueString()
	p.ResourceManagerBasePath = data.ResourceManagerCustomEndpoint.ValueString()
	p.SecretManagerBasePath = data.SecretManagerCustomEndpoint.ValueString()
	p.SecretManagerRegionalBasePath = data.SecretManagerRegionalCustomEndpoint.ValueString()
	p.SecureSourceManagerBasePath = data.SecureSourceManagerCustomEndpoint.ValueString()
	p.SecurityCenterBasePath = data.SecurityCenterCustomEndpoint.ValueString()
	p.SecurityCenterManagementBasePath = data.SecurityCenterManagementCustomEndpoint.ValueString()
	p.SecurityCenterV2BasePath = data.SecurityCenterV2CustomEndpoint.ValueString()
	p.SecuritypostureBasePath = data.SecuritypostureCustomEndpoint.ValueString()
	p.ServiceManagementBasePath = data.ServiceManagementCustomEndpoint.ValueString()
	p.ServiceNetworkingBasePath = data.ServiceNetworkingCustomEndpoint.ValueString()
	p.ServiceUsageBasePath = data.ServiceUsageCustomEndpoint.ValueString()
	p.SiteVerificationBasePath = data.SiteVerificationCustomEndpoint.ValueString()
	p.SourceRepoBasePath = data.SourceRepoCustomEndpoint.ValueString()
	p.SpannerBasePath = data.SpannerCustomEndpoint.ValueString()
	p.SQLBasePath = data.SQLCustomEndpoint.ValueString()
	p.StorageBasePath = data.StorageCustomEndpoint.ValueString()
	p.StorageInsightsBasePath = data.StorageInsightsCustomEndpoint.ValueString()
	p.StorageTransferBasePath = data.StorageTransferCustomEndpoint.ValueString()
	p.TagsBasePath = data.TagsCustomEndpoint.ValueString()
	p.TPUBasePath = data.TPUCustomEndpoint.ValueString()
	p.VertexAIBasePath = data.VertexAICustomEndpoint.ValueString()
	p.VmwareengineBasePath = data.VmwareengineCustomEndpoint.ValueString()
	p.VPCAccessBasePath = data.VPCAccessCustomEndpoint.ValueString()
	p.WorkbenchBasePath = data.WorkbenchCustomEndpoint.ValueString()
	p.WorkflowsBasePath = data.WorkflowsCustomEndpoint.ValueString()

	// Temporary
	p.Credentials = data.Credentials
	p.AccessToken = data.AccessToken
	p.ImpersonateServiceAccount = data.ImpersonateServiceAccount
	p.ImpersonateServiceAccountDelegates = data.ImpersonateServiceAccountDelegates
	p.RequestReason = data.RequestReason
	p.AddTerraformAttributionLabel = data.AddTerraformAttributionLabel
	p.TerraformAttributionLabelAdditionStrategy = data.TerraformAttributionLabelAdditionStrategy
	// End temporary

	// Copy values from the ProviderModel struct containing data about the provider configuration (present only when responsing to ConfigureProvider rpc calls)
	// to the FrameworkProviderConfig struct that will be passed and available to all resources/data sources
	p.Context = ctx
	p.BillingProject = data.BillingProject
	p.DefaultLabels = data.DefaultLabels
	p.Project = data.Project
	p.Region = GetRegionFromRegionSelfLink(data.Region)
	p.Scopes = data.Scopes
	p.Zone = data.Zone
	p.UserProjectOverride = data.UserProjectOverride
	p.PollInterval = 10 * time.Second
	p.UniverseDomain = data.UniverseDomain
	p.RequestBatcherServiceUsage = transport_tpg.NewRequestBatcher("Service Usage", ctx, batchingConfig)
	p.RequestBatcherIam = transport_tpg.NewRequestBatcher("IAM", ctx, batchingConfig)
}

// HandleDefaults will handle all the defaults necessary in the provider
func (p *FrameworkProviderConfig) HandleDefaults(ctx context.Context, data *fwmodels.ProviderModel, diags *diag.Diagnostics) {
	if (data.AccessToken.IsNull() || data.AccessToken.IsUnknown()) && (data.Credentials.IsNull() || data.Credentials.IsUnknown()) {
		credentials := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_CREDENTIALS",
			"GOOGLE_CLOUD_KEYFILE_JSON",
			"GCLOUD_KEYFILE_JSON",
		}, nil)

		if credentials != nil {
			data.Credentials = types.StringValue(credentials.(string))
		}

		accessToken := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_OAUTH_ACCESS_TOKEN",
		}, nil)

		if accessToken != nil {
			data.AccessToken = types.StringValue(accessToken.(string))
		}
	}

	if (data.ImpersonateServiceAccount.IsNull() || data.ImpersonateServiceAccount.IsUnknown()) && os.Getenv("GOOGLE_IMPERSONATE_SERVICE_ACCOUNT") != "" {
		data.ImpersonateServiceAccount = types.StringValue(os.Getenv("GOOGLE_IMPERSONATE_SERVICE_ACCOUNT"))
	}

	if data.Project.IsNull() || data.Project.IsUnknown() {
		project := transport_tpg.MultiEnvDefault([]string{
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

	if data.Region.IsNull() || data.Region.IsUnknown() {
		region := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_REGION",
			"GCLOUD_REGION",
			"CLOUDSDK_COMPUTE_REGION",
		}, nil)

		if region != nil {
			data.Region = types.StringValue(region.(string))
		}
	}

	if data.Zone.IsNull() || data.Zone.IsUnknown() {
		zone := transport_tpg.MultiEnvDefault([]string{
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
		data.Scopes, d = types.ListValueFrom(ctx, types.StringType, transport_tpg.DefaultClientScopes)
		diags.Append(d...)
		if diags.HasError() {
			return
		}
	}

	if !data.Batching.IsNull() && !data.Batching.IsUnknown() {
		var pbConfigs []fwmodels.ProviderBatching
		d := data.Batching.ElementsAs(ctx, &pbConfigs, true)
		diags.Append(d...)
		if diags.HasError() {
			return
		}

		if pbConfigs[0].SendAfter.IsNull() || pbConfigs[0].SendAfter.IsUnknown() {
			pbConfigs[0].SendAfter = types.StringValue("10s")
		}

		if pbConfigs[0].EnableBatching.IsNull() || pbConfigs[0].EnableBatching.IsUnknown() {
			pbConfigs[0].EnableBatching = types.BoolValue(true)
		}

		data.Batching, d = types.ListValueFrom(ctx, types.ObjectType{}.WithAttributeTypes(fwmodels.ProviderBatchingAttributes), pbConfigs)
	}

	if (data.UserProjectOverride.IsNull() || data.UserProjectOverride.IsUnknown()) && os.Getenv("USER_PROJECT_OVERRIDE") != "" {
		override, err := strconv.ParseBool(os.Getenv("USER_PROJECT_OVERRIDE"))
		if err != nil {
			diags.AddError(
				"error parsing environment variable `USER_PROJECT_OVERRIDE` into bool", err.Error())
		}
		data.UserProjectOverride = types.BoolValue(override)
	}

	if (data.RequestReason.IsNull() || data.RequestReason.IsUnknown()) && os.Getenv("CLOUDSDK_CORE_REQUEST_REASON") != "" {
		data.RequestReason = types.StringValue(os.Getenv("CLOUDSDK_CORE_REQUEST_REASON"))
	}

	if data.RequestTimeout.IsNull() || data.RequestTimeout.IsUnknown() {
		data.RequestTimeout = types.StringValue("120s")
	}

	// Generated Products
	if data.AccessApprovalCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_ACCESS_APPROVAL_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.AccessApprovalBasePathKey])
		if customEndpoint != nil {
			data.AccessApprovalCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.AccessContextManagerCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_ACCESS_CONTEXT_MANAGER_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.AccessContextManagerBasePathKey])
		if customEndpoint != nil {
			data.AccessContextManagerCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.ActiveDirectoryCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_ACTIVE_DIRECTORY_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.ActiveDirectoryBasePathKey])
		if customEndpoint != nil {
			data.ActiveDirectoryCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.AlloydbCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_ALLOYDB_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.AlloydbBasePathKey])
		if customEndpoint != nil {
			data.AlloydbCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.ApigeeCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_APIGEE_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.ApigeeBasePathKey])
		if customEndpoint != nil {
			data.ApigeeCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.AppEngineCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_APP_ENGINE_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.AppEngineBasePathKey])
		if customEndpoint != nil {
			data.AppEngineCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.ApphubCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_APPHUB_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.ApphubBasePathKey])
		if customEndpoint != nil {
			data.ApphubCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.ArtifactRegistryCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_ARTIFACT_REGISTRY_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.ArtifactRegistryBasePathKey])
		if customEndpoint != nil {
			data.ArtifactRegistryCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.BeyondcorpCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_BEYONDCORP_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.BeyondcorpBasePathKey])
		if customEndpoint != nil {
			data.BeyondcorpCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.BiglakeCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_BIGLAKE_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.BiglakeBasePathKey])
		if customEndpoint != nil {
			data.BiglakeCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.BigQueryCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_BIG_QUERY_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.BigQueryBasePathKey])
		if customEndpoint != nil {
			data.BigQueryCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.BigqueryAnalyticsHubCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_BIGQUERY_ANALYTICS_HUB_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.BigqueryAnalyticsHubBasePathKey])
		if customEndpoint != nil {
			data.BigqueryAnalyticsHubCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.BigqueryConnectionCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_BIGQUERY_CONNECTION_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.BigqueryConnectionBasePathKey])
		if customEndpoint != nil {
			data.BigqueryConnectionCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.BigqueryDatapolicyCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_BIGQUERY_DATAPOLICY_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.BigqueryDatapolicyBasePathKey])
		if customEndpoint != nil {
			data.BigqueryDatapolicyCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.BigqueryDataTransferCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_BIGQUERY_DATA_TRANSFER_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.BigqueryDataTransferBasePathKey])
		if customEndpoint != nil {
			data.BigqueryDataTransferCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.BigqueryReservationCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_BIGQUERY_RESERVATION_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.BigqueryReservationBasePathKey])
		if customEndpoint != nil {
			data.BigqueryReservationCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.BigtableCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_BIGTABLE_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.BigtableBasePathKey])
		if customEndpoint != nil {
			data.BigtableCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.BillingCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_BILLING_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.BillingBasePathKey])
		if customEndpoint != nil {
			data.BillingCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.BinaryAuthorizationCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_BINARY_AUTHORIZATION_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.BinaryAuthorizationBasePathKey])
		if customEndpoint != nil {
			data.BinaryAuthorizationCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.BlockchainNodeEngineCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_BLOCKCHAIN_NODE_ENGINE_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.BlockchainNodeEngineBasePathKey])
		if customEndpoint != nil {
			data.BlockchainNodeEngineCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.CertificateManagerCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_CERTIFICATE_MANAGER_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.CertificateManagerBasePathKey])
		if customEndpoint != nil {
			data.CertificateManagerCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.CloudAssetCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_CLOUD_ASSET_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.CloudAssetBasePathKey])
		if customEndpoint != nil {
			data.CloudAssetCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.CloudBuildCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_CLOUD_BUILD_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.CloudBuildBasePathKey])
		if customEndpoint != nil {
			data.CloudBuildCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.Cloudbuildv2CustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_CLOUDBUILDV2_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.Cloudbuildv2BasePathKey])
		if customEndpoint != nil {
			data.Cloudbuildv2CustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.ClouddeployCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_CLOUDDEPLOY_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.ClouddeployBasePathKey])
		if customEndpoint != nil {
			data.ClouddeployCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.ClouddomainsCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_CLOUDDOMAINS_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.ClouddomainsBasePathKey])
		if customEndpoint != nil {
			data.ClouddomainsCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.CloudFunctionsCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_CLOUD_FUNCTIONS_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.CloudFunctionsBasePathKey])
		if customEndpoint != nil {
			data.CloudFunctionsCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.Cloudfunctions2CustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_CLOUDFUNCTIONS2_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.Cloudfunctions2BasePathKey])
		if customEndpoint != nil {
			data.Cloudfunctions2CustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.CloudIdentityCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_CLOUD_IDENTITY_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.CloudIdentityBasePathKey])
		if customEndpoint != nil {
			data.CloudIdentityCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.CloudIdsCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_CLOUD_IDS_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.CloudIdsBasePathKey])
		if customEndpoint != nil {
			data.CloudIdsCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.CloudQuotasCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_CLOUD_QUOTAS_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.CloudQuotasBasePathKey])
		if customEndpoint != nil {
			data.CloudQuotasCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.CloudRunCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_CLOUD_RUN_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.CloudRunBasePathKey])
		if customEndpoint != nil {
			data.CloudRunCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.CloudRunV2CustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_CLOUD_RUN_V2_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.CloudRunV2BasePathKey])
		if customEndpoint != nil {
			data.CloudRunV2CustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.CloudSchedulerCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_CLOUD_SCHEDULER_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.CloudSchedulerBasePathKey])
		if customEndpoint != nil {
			data.CloudSchedulerCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.CloudTasksCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_CLOUD_TASKS_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.CloudTasksBasePathKey])
		if customEndpoint != nil {
			data.CloudTasksCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.ComposerCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_COMPOSER_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.ComposerBasePathKey])
		if customEndpoint != nil {
			data.ComposerCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.ComputeCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_COMPUTE_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.ComputeBasePathKey])
		if customEndpoint != nil {
			data.ComputeCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.ContainerAnalysisCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_CONTAINER_ANALYSIS_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.ContainerAnalysisBasePathKey])
		if customEndpoint != nil {
			data.ContainerAnalysisCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.ContainerAttachedCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_CONTAINER_ATTACHED_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.ContainerAttachedBasePathKey])
		if customEndpoint != nil {
			data.ContainerAttachedCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.CoreBillingCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_CORE_BILLING_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.CoreBillingBasePathKey])
		if customEndpoint != nil {
			data.CoreBillingCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.DatabaseMigrationServiceCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_DATABASE_MIGRATION_SERVICE_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.DatabaseMigrationServiceBasePathKey])
		if customEndpoint != nil {
			data.DatabaseMigrationServiceCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.DataCatalogCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_DATA_CATALOG_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.DataCatalogBasePathKey])
		if customEndpoint != nil {
			data.DataCatalogCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.DataFusionCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_DATA_FUSION_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.DataFusionBasePathKey])
		if customEndpoint != nil {
			data.DataFusionCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.DataLossPreventionCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_DATA_LOSS_PREVENTION_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.DataLossPreventionBasePathKey])
		if customEndpoint != nil {
			data.DataLossPreventionCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.DataPipelineCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_DATA_PIPELINE_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.DataPipelineBasePathKey])
		if customEndpoint != nil {
			data.DataPipelineCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.DataplexCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_DATAPLEX_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.DataplexBasePathKey])
		if customEndpoint != nil {
			data.DataplexCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.DataprocCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_DATAPROC_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.DataprocBasePathKey])
		if customEndpoint != nil {
			data.DataprocCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.DataprocMetastoreCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_DATAPROC_METASTORE_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.DataprocMetastoreBasePathKey])
		if customEndpoint != nil {
			data.DataprocMetastoreCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.DatastreamCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_DATASTREAM_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.DatastreamBasePathKey])
		if customEndpoint != nil {
			data.DatastreamCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.DeploymentManagerCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_DEPLOYMENT_MANAGER_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.DeploymentManagerBasePathKey])
		if customEndpoint != nil {
			data.DeploymentManagerCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.DialogflowCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_DIALOGFLOW_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.DialogflowBasePathKey])
		if customEndpoint != nil {
			data.DialogflowCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.DialogflowCXCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_DIALOGFLOW_CX_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.DialogflowCXBasePathKey])
		if customEndpoint != nil {
			data.DialogflowCXCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.DiscoveryEngineCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_DISCOVERY_ENGINE_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.DiscoveryEngineBasePathKey])
		if customEndpoint != nil {
			data.DiscoveryEngineCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.DNSCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_DNS_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.DNSBasePathKey])
		if customEndpoint != nil {
			data.DNSCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.DocumentAICustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_DOCUMENT_AI_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.DocumentAIBasePathKey])
		if customEndpoint != nil {
			data.DocumentAICustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.DocumentAIWarehouseCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_DOCUMENT_AI_WAREHOUSE_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.DocumentAIWarehouseBasePathKey])
		if customEndpoint != nil {
			data.DocumentAIWarehouseCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.EdgecontainerCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_EDGECONTAINER_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.EdgecontainerBasePathKey])
		if customEndpoint != nil {
			data.EdgecontainerCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.EdgenetworkCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_EDGENETWORK_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.EdgenetworkBasePathKey])
		if customEndpoint != nil {
			data.EdgenetworkCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.EssentialContactsCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_ESSENTIAL_CONTACTS_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.EssentialContactsBasePathKey])
		if customEndpoint != nil {
			data.EssentialContactsCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.FilestoreCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_FILESTORE_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.FilestoreBasePathKey])
		if customEndpoint != nil {
			data.FilestoreCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.FirebaseAppCheckCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_FIREBASE_APP_CHECK_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.FirebaseAppCheckBasePathKey])
		if customEndpoint != nil {
			data.FirebaseAppCheckCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.FirestoreCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_FIRESTORE_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.FirestoreBasePathKey])
		if customEndpoint != nil {
			data.FirestoreCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.GKEBackupCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_GKE_BACKUP_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.GKEBackupBasePathKey])
		if customEndpoint != nil {
			data.GKEBackupCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.GKEHubCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_GKE_HUB_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.GKEHubBasePathKey])
		if customEndpoint != nil {
			data.GKEHubCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.GKEHub2CustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_GKE_HUB2_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.GKEHub2BasePathKey])
		if customEndpoint != nil {
			data.GKEHub2CustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.GkeonpremCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_GKEONPREM_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.GkeonpremBasePathKey])
		if customEndpoint != nil {
			data.GkeonpremCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.HealthcareCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_HEALTHCARE_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.HealthcareBasePathKey])
		if customEndpoint != nil {
			data.HealthcareCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.IAM2CustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_IAM2_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.IAM2BasePathKey])
		if customEndpoint != nil {
			data.IAM2CustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.IAMBetaCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_IAM_BETA_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.IAMBetaBasePathKey])
		if customEndpoint != nil {
			data.IAMBetaCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.IAMWorkforcePoolCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_IAM_WORKFORCE_POOL_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.IAMWorkforcePoolBasePathKey])
		if customEndpoint != nil {
			data.IAMWorkforcePoolCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.IapCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_IAP_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.IapBasePathKey])
		if customEndpoint != nil {
			data.IapCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.IdentityPlatformCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_IDENTITY_PLATFORM_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.IdentityPlatformBasePathKey])
		if customEndpoint != nil {
			data.IdentityPlatformCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.IntegrationConnectorsCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_INTEGRATION_CONNECTORS_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.IntegrationConnectorsBasePathKey])
		if customEndpoint != nil {
			data.IntegrationConnectorsCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.IntegrationsCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_INTEGRATIONS_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.IntegrationsBasePathKey])
		if customEndpoint != nil {
			data.IntegrationsCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.KMSCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_KMS_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.KMSBasePathKey])
		if customEndpoint != nil {
			data.KMSCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.LoggingCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_LOGGING_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.LoggingBasePathKey])
		if customEndpoint != nil {
			data.LoggingCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.LookerCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_LOOKER_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.LookerBasePathKey])
		if customEndpoint != nil {
			data.LookerCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.MemcacheCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_MEMCACHE_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.MemcacheBasePathKey])
		if customEndpoint != nil {
			data.MemcacheCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.MigrationCenterCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_MIGRATION_CENTER_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.MigrationCenterBasePathKey])
		if customEndpoint != nil {
			data.MigrationCenterCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.MLEngineCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_ML_ENGINE_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.MLEngineBasePathKey])
		if customEndpoint != nil {
			data.MLEngineCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.MonitoringCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_MONITORING_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.MonitoringBasePathKey])
		if customEndpoint != nil {
			data.MonitoringCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.NetappCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_NETAPP_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.NetappBasePathKey])
		if customEndpoint != nil {
			data.NetappCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.NetworkConnectivityCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_NETWORK_CONNECTIVITY_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.NetworkConnectivityBasePathKey])
		if customEndpoint != nil {
			data.NetworkConnectivityCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.NetworkManagementCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_NETWORK_MANAGEMENT_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.NetworkManagementBasePathKey])
		if customEndpoint != nil {
			data.NetworkManagementCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.NetworkSecurityCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_NETWORK_SECURITY_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.NetworkSecurityBasePathKey])
		if customEndpoint != nil {
			data.NetworkSecurityCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.NetworkServicesCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_NETWORK_SERVICES_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.NetworkServicesBasePathKey])
		if customEndpoint != nil {
			data.NetworkServicesCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.NotebooksCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_NOTEBOOKS_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.NotebooksBasePathKey])
		if customEndpoint != nil {
			data.NotebooksCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.OrgPolicyCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_ORG_POLICY_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.OrgPolicyBasePathKey])
		if customEndpoint != nil {
			data.OrgPolicyCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.OSConfigCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_OS_CONFIG_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.OSConfigBasePathKey])
		if customEndpoint != nil {
			data.OSConfigCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.OSLoginCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_OS_LOGIN_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.OSLoginBasePathKey])
		if customEndpoint != nil {
			data.OSLoginCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.PrivatecaCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_PRIVATECA_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.PrivatecaBasePathKey])
		if customEndpoint != nil {
			data.PrivatecaCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.PrivilegedAccessManagerCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_PRIVILEGED_ACCESS_MANAGER_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.PrivilegedAccessManagerBasePathKey])
		if customEndpoint != nil {
			data.PrivilegedAccessManagerCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.PublicCACustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_PUBLIC_CA_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.PublicCABasePathKey])
		if customEndpoint != nil {
			data.PublicCACustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.PubsubCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_PUBSUB_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.PubsubBasePathKey])
		if customEndpoint != nil {
			data.PubsubCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.PubsubLiteCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_PUBSUB_LITE_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.PubsubLiteBasePathKey])
		if customEndpoint != nil {
			data.PubsubLiteCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.RedisCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_REDIS_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.RedisBasePathKey])
		if customEndpoint != nil {
			data.RedisCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.ResourceManagerCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_RESOURCE_MANAGER_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.ResourceManagerBasePathKey])
		if customEndpoint != nil {
			data.ResourceManagerCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.SecretManagerCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_SECRET_MANAGER_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.SecretManagerBasePathKey])
		if customEndpoint != nil {
			data.SecretManagerCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.SecretManagerRegionalCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_SECRET_MANAGER_REGIONAL_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.SecretManagerRegionalBasePathKey])
		if customEndpoint != nil {
			data.SecretManagerRegionalCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.SecureSourceManagerCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_SECURE_SOURCE_MANAGER_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.SecureSourceManagerBasePathKey])
		if customEndpoint != nil {
			data.SecureSourceManagerCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.SecurityCenterCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_SECURITY_CENTER_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.SecurityCenterBasePathKey])
		if customEndpoint != nil {
			data.SecurityCenterCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.SecurityCenterManagementCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_SECURITY_CENTER_MANAGEMENT_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.SecurityCenterManagementBasePathKey])
		if customEndpoint != nil {
			data.SecurityCenterManagementCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.SecurityCenterV2CustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_SECURITY_CENTER_V2_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.SecurityCenterV2BasePathKey])
		if customEndpoint != nil {
			data.SecurityCenterV2CustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.SecuritypostureCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_SECURITYPOSTURE_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.SecuritypostureBasePathKey])
		if customEndpoint != nil {
			data.SecuritypostureCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.ServiceManagementCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_SERVICE_MANAGEMENT_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.ServiceManagementBasePathKey])
		if customEndpoint != nil {
			data.ServiceManagementCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.ServiceNetworkingCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_SERVICE_NETWORKING_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.ServiceNetworkingBasePathKey])
		if customEndpoint != nil {
			data.ServiceNetworkingCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.ServiceUsageCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_SERVICE_USAGE_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.ServiceUsageBasePathKey])
		if customEndpoint != nil {
			data.ServiceUsageCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.SiteVerificationCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_SITE_VERIFICATION_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.SiteVerificationBasePathKey])
		if customEndpoint != nil {
			data.SiteVerificationCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.SourceRepoCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_SOURCE_REPO_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.SourceRepoBasePathKey])
		if customEndpoint != nil {
			data.SourceRepoCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.SpannerCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_SPANNER_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.SpannerBasePathKey])
		if customEndpoint != nil {
			data.SpannerCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.SQLCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_SQL_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.SQLBasePathKey])
		if customEndpoint != nil {
			data.SQLCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.StorageCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_STORAGE_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.StorageBasePathKey])
		if customEndpoint != nil {
			data.StorageCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.StorageInsightsCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_STORAGE_INSIGHTS_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.StorageInsightsBasePathKey])
		if customEndpoint != nil {
			data.StorageInsightsCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.StorageTransferCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_STORAGE_TRANSFER_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.StorageTransferBasePathKey])
		if customEndpoint != nil {
			data.StorageTransferCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.TagsCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_TAGS_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.TagsBasePathKey])
		if customEndpoint != nil {
			data.TagsCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.TPUCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_TPU_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.TPUBasePathKey])
		if customEndpoint != nil {
			data.TPUCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.VertexAICustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_VERTEX_AI_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.VertexAIBasePathKey])
		if customEndpoint != nil {
			data.VertexAICustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.VmwareengineCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_VMWAREENGINE_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.VmwareengineBasePathKey])
		if customEndpoint != nil {
			data.VmwareengineCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.VPCAccessCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_VPC_ACCESS_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.VPCAccessBasePathKey])
		if customEndpoint != nil {
			data.VPCAccessCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.WorkbenchCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_WORKBENCH_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.WorkbenchBasePathKey])
		if customEndpoint != nil {
			data.WorkbenchCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
	if data.WorkflowsCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_WORKFLOWS_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.WorkflowsBasePathKey])
		if customEndpoint != nil {
			data.WorkflowsCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}

	// Handwritten Products / Versioned / Atypical Entries
	if data.CloudBillingCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_CLOUD_BILLING_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths["cloud_billing_custom_endpoint"])
		if customEndpoint != nil {
			data.CloudBillingCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}

	if data.ComposerCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_COMPOSER_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.ComposerBasePathKey])
		if customEndpoint != nil {
			data.ComposerCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}

	if data.ContainerCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_CONTAINER_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.ContainerBasePathKey])
		if customEndpoint != nil {
			data.ContainerCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}

	if data.DataflowCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_DATAFLOW_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.DataflowBasePathKey])
		if customEndpoint != nil {
			data.DataflowCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}

	if data.IamCredentialsCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_IAM_CREDENTIALS_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.IamCredentialsBasePathKey])
		if customEndpoint != nil {
			data.IamCredentialsCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}

	if data.ResourceManagerV3CustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_RESOURCE_MANAGER_V3_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.ResourceManagerV3BasePathKey])
		if customEndpoint != nil {
			data.ResourceManagerV3CustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}

	if data.IAMCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_IAM_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.IAMBasePathKey])
		if customEndpoint != nil {
			data.IAMCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}

	if data.ServiceNetworkingCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_SERVICE_NETWORKING_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.ServiceNetworkingBasePathKey])
		if customEndpoint != nil {
			data.ServiceNetworkingCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}

	if data.TagsLocationCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_TAGS_LOCATION_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.TagsLocationBasePathKey])
		if customEndpoint != nil {
			data.TagsLocationCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}

	// dcl
	if data.ContainerAwsCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_CONTAINERAWS_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.ContainerAwsBasePathKey])
		if customEndpoint != nil {
			data.ContainerAwsCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}

	if data.ContainerAzureCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_CONTAINERAZURE_CUSTOM_ENDPOINT",
		}, transport_tpg.DefaultBasePaths[transport_tpg.ContainerAzureBasePathKey])
		if customEndpoint != nil {
			data.ContainerAzureCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}

	// DCL generated defaults
	if data.ApikeysCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_APIKEYS_CUSTOM_ENDPOINT",
		}, "")
		if customEndpoint != nil {
			data.ApikeysCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}

	if data.AssuredWorkloadsCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_ASSURED_WORKLOADS_CUSTOM_ENDPOINT",
		}, "")
		if customEndpoint != nil {
			data.AssuredWorkloadsCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}

	if data.CloudBuildWorkerPoolCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_CLOUD_BUILD_WORKER_POOL_CUSTOM_ENDPOINT",
		}, "")
		if customEndpoint != nil {
			data.CloudBuildWorkerPoolCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}

	if data.CloudResourceManagerCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_CLOUD_RESOURCE_MANAGER_CUSTOM_ENDPOINT",
		}, "")
		if customEndpoint != nil {
			data.CloudResourceManagerCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}

	if data.DataplexCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_DATAPLEX_CUSTOM_ENDPOINT",
		}, "")
		if customEndpoint != nil {
			data.DataplexCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}

	if data.EventarcCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_EVENTARC_CUSTOM_ENDPOINT",
		}, "")
		if customEndpoint != nil {
			data.EventarcCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}

	if data.FirebaserulesCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_FIREBASERULES_CUSTOM_ENDPOINT",
		}, "")
		if customEndpoint != nil {
			data.FirebaserulesCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}

	if data.NetworkConnectivityCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_NETWORK_CONNECTIVITY_CUSTOM_ENDPOINT",
		}, "")
		if customEndpoint != nil {
			data.NetworkConnectivityCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}

	if data.RecaptchaEnterpriseCustomEndpoint.IsNull() {
		customEndpoint := transport_tpg.MultiEnvDefault([]string{
			"GOOGLE_RECAPTCHA_ENTERPRISE_CUSTOM_ENDPOINT",
		}, "")
		if customEndpoint != nil {
			data.RecaptchaEnterpriseCustomEndpoint = types.StringValue(customEndpoint.(string))
		}
	}
}

func (p *FrameworkProviderConfig) SetupClient(ctx context.Context, data fwmodels.ProviderModel, diags *diag.Diagnostics) {
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
	retryTransport := transport_tpg.NewTransportWithDefaultRetries(loggingTransport)

	// 4. Header Transport - outer wrapper to inject additional headers we want to apply
	// before making requests
	headerTransport := transport_tpg.NewTransportWithHeaders(retryTransport)
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

	p.TokenSource = tokenSource
	p.Client = client
}

func (p *FrameworkProviderConfig) SetupGrpcLogging() {
	logger := logrus.StandardLogger()

	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&transport_tpg.Formatter{
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

func (p *FrameworkProviderConfig) logGoogleIdentities(ctx context.Context, data fwmodels.ProviderModel, diags *diag.Diagnostics) {
	// GetCurrentUserEmailFramework doesn't pass an error back from logGoogleIdentities, so we want
	// a separate diagnostics here
	var d diag.Diagnostics

	if data.ImpersonateServiceAccount.IsNull() || data.ImpersonateServiceAccount.IsUnknown() {

		tokenSource := GetTokenSource(ctx, data, true, diags)
		if diags.HasError() {
			return
		}

		p.Client = oauth2.NewClient(ctx, tokenSource) // p.Client isn't initialised fully when this code is called.

		email := GetCurrentUserEmailFramework(p, p.UserAgent, &d)
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

	p.Client = oauth2.NewClient(ctx, tokenSource) // p.Client isn't initialised fully when this code is called.
	email := GetCurrentUserEmailFramework(p, p.UserAgent, &d)
	if d.HasError() {
		tflog.Info(ctx, "error retrieving userinfo for your provider credentials. have you enabled the 'https://www.googleapis.com/auth/userinfo.email' scope?")
	}

	tflog.Info(ctx, fmt.Sprintf("Terraform is configured with service account impersonation, original identity: %s, impersonated identity: %s", email, data.ImpersonateServiceAccount.ValueString()))

	// Add the Impersonated ClientOption back in to the OAuth2 TokenSource
	tokenSource = GetTokenSource(ctx, data, false, diags)
	if diags.HasError() {
		return
	}

	p.Client = oauth2.NewClient(ctx, tokenSource) // p.Client isn't initialised fully when this code is called.

	return
}

// Configuration helpers

// GetTokenSource gets token source based on the Google Credentials configured.
// If initialCredentialsOnly is true, don't follow the impersonation settings and return the initial set of creds.
func GetTokenSource(ctx context.Context, data fwmodels.ProviderModel, initialCredentialsOnly bool, diags *diag.Diagnostics) oauth2.TokenSource {
	creds := GetCredentials(ctx, data, initialCredentialsOnly, diags)

	return creds.TokenSource
}

// GetCredentials gets credentials with a given scope (clientScopes).
// If initialCredentialsOnly is true, don't follow the impersonation
// settings and return the initial set of creds instead.
func GetCredentials(ctx context.Context, data fwmodels.ProviderModel, initialCredentialsOnly bool, diags *diag.Diagnostics) googleoauth.Credentials {
	var clientScopes []string
	var delegates []string

	if !data.Scopes.IsNull() && !data.Scopes.IsUnknown() {
		d := data.Scopes.ElementsAs(ctx, &clientScopes, false)
		diags.Append(d...)
		if diags.HasError() {
			return googleoauth.Credentials{}
		}
	}

	if !data.ImpersonateServiceAccountDelegates.IsNull() && !data.ImpersonateServiceAccountDelegates.IsUnknown() {
		d := data.ImpersonateServiceAccountDelegates.ElementsAs(ctx, &delegates, false)
		diags.Append(d...)
		if diags.HasError() {
			return googleoauth.Credentials{}
		}
	}

	if !data.AccessToken.IsNull() && !data.AccessToken.IsUnknown() {
		contents, _, err := verify.PathOrContents(data.AccessToken.ValueString())
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
			TokenSource: transport_tpg.StaticTokenSource{oauth2.StaticTokenSource(token)},
		}
	}

	if !data.Credentials.IsNull() && !data.Credentials.IsUnknown() {
		contents, _, err := verify.PathOrContents(data.Credentials.ValueString())
		if err != nil {
			diags.AddError(fmt.Sprintf("error loading credentials: %s", err), err.Error())
			return googleoauth.Credentials{}
		}
		if len(contents) == 0 {
			diags.AddError("error loading credentials", "provided credentials are empty")
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

		creds, err := transport.Creds(ctx, option.WithCredentialsJSON([]byte(contents)), option.WithScopes(clientScopes...))
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
	creds, err := transport.Creds(context.Background(), option.WithScopes(clientScopes...))
	if err != nil {
		diags.AddError(fmt.Sprintf("Attempted to load application default credentials since neither `credentials` nor `access_token` was set in the provider block.  "+
			"No credentials loaded. To use your gcloud credentials, run 'gcloud auth application-default login'"), err.Error())
		return googleoauth.Credentials{}
	}

	return *creds
}

// GetBatchingConfig returns the batching config object given the
// provider configuration set for batching
func GetBatchingConfig(ctx context.Context, data types.List, diags *diag.Diagnostics) *transport_tpg.BatchingConfig {
	bc := &transport_tpg.BatchingConfig{
		SendAfter:      time.Second * transport_tpg.DefaultBatchSendIntervalSec,
		EnableBatching: true,
	}

	// Handle if entire batching block is null/unknown
	if data.IsNull() || data.IsUnknown() {
		return bc
	}

	var pbConfigs []fwmodels.ProviderBatching
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

func GetRegionFromRegionSelfLink(selfLink basetypes.StringValue) basetypes.StringValue {
	re := regexp.MustCompile("/compute/[a-zA-Z0-9]*/projects/[a-zA-Z0-9-]*/regions/([a-zA-Z0-9-]*)")
	value := selfLink.String()
	switch {
	case re.MatchString(value):
		if res := re.FindStringSubmatch(value); len(res) == 2 && res[1] != "" {
			region := res[1]
			return types.StringValue(region)
		}
	}
	return selfLink
}
