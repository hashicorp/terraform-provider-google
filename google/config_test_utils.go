package google

import (
	"net/http/httptest"
	"strings"
)

// NewTestConfig create a config using the http test server.
func NewTestConfig(server *httptest.Server) *Config {
	cfg := &Config{}
	cfg.client = server.Client()
	configureTestBasePaths(cfg, server.URL)
	return cfg
}

func configureTestBasePaths(c *Config, url string) {
	if !strings.HasSuffix(url, "/") {
		url = url + "/"
	}
	// Generated Products
	c.AccessApprovalBasePath = url
	c.AccessContextManagerBasePath = url
	c.ActiveDirectoryBasePath = url
	c.ApigeeBasePath = url
	c.AppEngineBasePath = url
	c.ArtifactRegistryBasePath = url
	c.BeyondcorpBasePath = url
	c.BigQueryBasePath = url
	c.BigqueryAnalyticsHubBasePath = url
	c.BigqueryConnectionBasePath = url
	c.BigqueryDataTransferBasePath = url
	c.BigqueryReservationBasePath = url
	c.BigtableBasePath = url
	c.BillingBasePath = url
	c.BinaryAuthorizationBasePath = url
	c.CertificateManagerBasePath = url
	c.CloudAssetBasePath = url
	c.CloudBuildBasePath = url
	c.CloudFunctionsBasePath = url
	c.Cloudfunctions2BasePath = url
	c.CloudIdentityBasePath = url
	c.CloudIdsBasePath = url
	c.CloudIotBasePath = url
	c.CloudRunBasePath = url
	c.CloudRunV2BasePath = url
	c.CloudSchedulerBasePath = url
	c.CloudTasksBasePath = url
	c.ComputeBasePath = url
	c.ContainerAnalysisBasePath = url
	c.DataCatalogBasePath = url
	c.DataFusionBasePath = url
	c.DataLossPreventionBasePath = url
	c.DataprocBasePath = url
	c.DataprocMetastoreBasePath = url
	c.DatastoreBasePath = url
	c.DatastreamBasePath = url
	c.DeploymentManagerBasePath = url
	c.DialogflowBasePath = url
	c.DialogflowCXBasePath = url
	c.DNSBasePath = url
	c.DocumentAIBasePath = url
	c.EssentialContactsBasePath = url
	c.FilestoreBasePath = url
	c.FirestoreBasePath = url
	c.GameServicesBasePath = url
	c.GKEHubBasePath = url
	c.HealthcareBasePath = url
	c.IAMBetaBasePath = url
	c.IAMWorkforcePoolBasePath = url
	c.IapBasePath = url
	c.IdentityPlatformBasePath = url
	c.KMSBasePath = url
	c.LoggingBasePath = url
	c.MemcacheBasePath = url
	c.MLEngineBasePath = url
	c.MonitoringBasePath = url
	c.NetworkManagementBasePath = url
	c.NetworkServicesBasePath = url
	c.NotebooksBasePath = url
	c.OSConfigBasePath = url
	c.OSLoginBasePath = url
	c.PrivatecaBasePath = url
	c.PubsubBasePath = url
	c.PubsubLiteBasePath = url
	c.RedisBasePath = url
	c.ResourceManagerBasePath = url
	c.SecretManagerBasePath = url
	c.SecurityCenterBasePath = url
	c.ServiceManagementBasePath = url
	c.ServiceUsageBasePath = url
	c.SourceRepoBasePath = url
	c.SpannerBasePath = url
	c.SQLBasePath = url
	c.StorageBasePath = url
	c.StorageTransferBasePath = url
	c.TagsBasePath = url
	c.TPUBasePath = url
	c.VertexAIBasePath = url
	c.VPCAccessBasePath = url
	c.WorkflowsBasePath = url

	// Handwritten Products / Versioned / Atypical Entries
	c.CloudBillingBasePath = url
	c.ComposerBasePath = url
	c.ContainerBasePath = url
	c.DataprocBasePath = url
	c.DataflowBasePath = url
	c.IamCredentialsBasePath = url
	c.ResourceManagerV3BasePath = url
	c.IAMBasePath = url
	c.ServiceNetworkingBasePath = url
	c.BigQueryBasePath = url
	c.StorageTransferBasePath = url
	c.BigtableAdminBasePath = url
}
