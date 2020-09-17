package google

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"google.golang.org/api/option"

	"golang.org/x/oauth2"
	googleoauth "golang.org/x/oauth2/google"
	appengine "google.golang.org/api/appengine/v1"
	"google.golang.org/api/bigquery/v2"
	"google.golang.org/api/bigtableadmin/v2"
	"google.golang.org/api/cloudbilling/v1"
	"google.golang.org/api/cloudbuild/v1"
	"google.golang.org/api/cloudfunctions/v1"
	"google.golang.org/api/cloudiot/v1"
	"google.golang.org/api/cloudkms/v1"
	"google.golang.org/api/cloudresourcemanager/v1"
	resourceManagerV2Beta1 "google.golang.org/api/cloudresourcemanager/v2beta1"
	composer "google.golang.org/api/composer/v1beta1"
	computeBeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/container/v1"
	containerBeta "google.golang.org/api/container/v1beta1"
	dataflow "google.golang.org/api/dataflow/v1b3"
	"google.golang.org/api/dataproc/v1"
	dataprocBeta "google.golang.org/api/dataproc/v1beta2"
	"google.golang.org/api/dns/v1"
	dnsBeta "google.golang.org/api/dns/v1beta2"
	file "google.golang.org/api/file/v1beta1"
	healthcare "google.golang.org/api/healthcare/v1"
	"google.golang.org/api/iam/v1"
	iamcredentials "google.golang.org/api/iamcredentials/v1"
	cloudlogging "google.golang.org/api/logging/v2"
	"google.golang.org/api/pubsub/v1"
	runtimeconfig "google.golang.org/api/runtimeconfig/v1beta1"
	"google.golang.org/api/servicemanagement/v1"
	"google.golang.org/api/servicenetworking/v1"
	"google.golang.org/api/serviceusage/v1"
	"google.golang.org/api/sourcerepo/v1"
	"google.golang.org/api/spanner/v1"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
	"google.golang.org/api/storage/v1"
	"google.golang.org/api/storagetransfer/v1"
)

// Config is the configuration structure used to instantiate the Google
// provider.
type Config struct {
	Credentials         string
	AccessToken         string
	Project             string
	BillingProject      string
	Region              string
	Zone                string
	Scopes              []string
	BatchingConfig      *batchingConfig
	UserProjectOverride bool
	RequestTimeout      time.Duration
	// PollInterval is passed to resource.StateChangeConf in common_operation.go
	// It controls the interval at which we poll for successful operations
	PollInterval time.Duration

	client                *http.Client
	wrappedBigQueryClient *http.Client
	wrappedPubsubClient   *http.Client
	context               context.Context
	userAgent             string

	tokenSource oauth2.TokenSource

	AccessApprovalBasePath       string
	AccessContextManagerBasePath string
	ActiveDirectoryBasePath      string
	AppEngineBasePath            string
	BigQueryBasePath             string
	BigqueryDataTransferBasePath string
	BigtableBasePath             string
	BinaryAuthorizationBasePath  string
	CloudAssetBasePath           string
	CloudBuildBasePath           string
	CloudFunctionsBasePath       string
	CloudIotBasePath             string
	CloudRunBasePath             string
	CloudSchedulerBasePath       string
	CloudTasksBasePath           string
	ComputeBasePath              string
	ContainerAnalysisBasePath    string
	DataCatalogBasePath          string
	DataLossPreventionBasePath   string
	DataprocBasePath             string
	DatastoreBasePath            string
	DeploymentManagerBasePath    string
	DialogflowBasePath           string
	DNSBasePath                  string
	FilestoreBasePath            string
	FirestoreBasePath            string
	GameServicesBasePath         string
	HealthcareBasePath           string
	IapBasePath                  string
	IdentityPlatformBasePath     string
	KMSBasePath                  string
	LoggingBasePath              string
	MLEngineBasePath             string
	MonitoringBasePath           string
	NetworkManagementBasePath    string
	OSConfigBasePath             string
	OSLoginBasePath              string
	PubsubBasePath               string
	RedisBasePath                string
	ResourceManagerBasePath      string
	RuntimeConfigBasePath        string
	SecretManagerBasePath        string
	SecurityCenterBasePath       string
	ServiceManagementBasePath    string
	ServiceUsageBasePath         string
	SourceRepoBasePath           string
	SpannerBasePath              string
	SQLBasePath                  string
	StorageBasePath              string
	TPUBasePath                  string
	VPCAccessBasePath            string

	CloudBillingBasePath string
	clientBilling        *cloudbilling.APIService

	clientBuild *cloudbuild.Service

	ComposerBasePath string
	clientComposer   *composer.Service

	clientCompute *compute.Service

	ComputeBetaBasePath string
	clientComputeBeta   *computeBeta.Service

	ContainerBasePath string
	clientContainer   *container.Service

	ContainerBetaBasePath string
	clientContainerBeta   *containerBeta.Service

	clientDataproc *dataproc.Service

	DataprocBetaBasePath string
	clientDataprocBeta   *dataprocBeta.Service

	DataflowBasePath string
	clientDataflow   *dataflow.Service

	clientDns *dns.Service

	DnsBetaBasePath string
	clientDnsBeta   *dnsBeta.Service

	clientFilestore *file.Service

	IamCredentialsBasePath string
	clientIamCredentials   *iamcredentials.Service

	clientKms *cloudkms.Service

	clientLogging *cloudlogging.Service

	clientPubsub *pubsub.Service

	clientResourceManager *cloudresourcemanager.Service

	ResourceManagerV2Beta1BasePath string
	clientResourceManagerV2Beta1   *resourceManagerV2Beta1.Service

	clientRuntimeconfig *runtimeconfig.Service

	clientSpanner *spanner.Service

	clientSourceRepo *sourcerepo.Service

	clientStorage *storage.Service

	clientSqlAdmin *sqladmin.Service

	IAMBasePath string
	clientIAM   *iam.Service

	clientHealthcare *healthcare.Service

	clientServiceMan *servicemanagement.APIService

	clientServiceUsage *serviceusage.Service

	clientBigQuery *bigquery.Service

	clientCloudFunctions *cloudfunctions.Service

	CloudIoTBasePath string
	clientCloudIoT   *cloudiot.Service

	clientAppEngine *appengine.APIService

	ServiceNetworkingBasePath string
	clientServiceNetworking   *servicenetworking.APIService

	StorageTransferBasePath string
	clientStorageTransfer   *storagetransfer.Service

	bigtableClientFactory *BigtableClientFactory
	BigtableAdminBasePath string
	// Unlike other clients, the Bigtable Admin client doesn't use a single
	// service. Instead, there are several distinct services created off
	// the base service object. To imitate most other handwritten clients,
	// we expose those directly instead of providing the `Service` object
	// as a factory.
	clientBigtableProjectsInstances *bigtableadmin.ProjectsInstancesService

	requestBatcherServiceUsage *RequestBatcher
	requestBatcherIam          *RequestBatcher
}

// Generated product base paths
var AccessApprovalDefaultBasePath = "https://accessapproval.googleapis.com/v1/"
var AccessContextManagerDefaultBasePath = "https://accesscontextmanager.googleapis.com/v1/"
var ActiveDirectoryDefaultBasePath = "https://managedidentities.googleapis.com/v1/"
var AppEngineDefaultBasePath = "https://appengine.googleapis.com/v1/"
var BigQueryDefaultBasePath = "https://bigquery.googleapis.com/bigquery/v2/"
var BigqueryDataTransferDefaultBasePath = "https://bigquerydatatransfer.googleapis.com/v1/"
var BigtableDefaultBasePath = "https://bigtableadmin.googleapis.com/v2/"
var BinaryAuthorizationDefaultBasePath = "https://binaryauthorization.googleapis.com/v1/"
var CloudAssetDefaultBasePath = "https://cloudasset.googleapis.com/v1/"
var CloudBuildDefaultBasePath = "https://cloudbuild.googleapis.com/v1/"
var CloudFunctionsDefaultBasePath = "https://cloudfunctions.googleapis.com/v1/"
var CloudIotDefaultBasePath = "https://cloudiot.googleapis.com/v1/"
var CloudRunDefaultBasePath = "https://{{location}}-run.googleapis.com/"
var CloudSchedulerDefaultBasePath = "https://cloudscheduler.googleapis.com/v1/"
var CloudTasksDefaultBasePath = "https://cloudtasks.googleapis.com/v2/"
var ComputeDefaultBasePath = "https://compute.googleapis.com/compute/v1/"
var ContainerAnalysisDefaultBasePath = "https://containeranalysis.googleapis.com/v1/"
var DataCatalogDefaultBasePath = "https://datacatalog.googleapis.com/v1/"
var DataLossPreventionDefaultBasePath = "https://dlp.googleapis.com/v2/"
var DataprocDefaultBasePath = "https://dataproc.googleapis.com/v1/"
var DatastoreDefaultBasePath = "https://datastore.googleapis.com/v1/"
var DeploymentManagerDefaultBasePath = "https://www.googleapis.com/deploymentmanager/v2/"
var DialogflowDefaultBasePath = "https://dialogflow.googleapis.com/v2/"
var DNSDefaultBasePath = "https://dns.googleapis.com/dns/v1/"
var FilestoreDefaultBasePath = "https://file.googleapis.com/v1/"
var FirestoreDefaultBasePath = "https://firestore.googleapis.com/v1/"
var GameServicesDefaultBasePath = "https://gameservices.googleapis.com/v1/"
var HealthcareDefaultBasePath = "https://healthcare.googleapis.com/v1/"
var IapDefaultBasePath = "https://iap.googleapis.com/v1/"
var IdentityPlatformDefaultBasePath = "https://identitytoolkit.googleapis.com/v2/"
var KMSDefaultBasePath = "https://cloudkms.googleapis.com/v1/"
var LoggingDefaultBasePath = "https://logging.googleapis.com/v2/"
var MLEngineDefaultBasePath = "https://ml.googleapis.com/v1/"
var MonitoringDefaultBasePath = "https://monitoring.googleapis.com/"
var NetworkManagementDefaultBasePath = "https://networkmanagement.googleapis.com/v1/"
var OSConfigDefaultBasePath = "https://osconfig.googleapis.com/v1/"
var OSLoginDefaultBasePath = "https://oslogin.googleapis.com/v1/"
var PubsubDefaultBasePath = "https://pubsub.googleapis.com/v1/"
var RedisDefaultBasePath = "https://redis.googleapis.com/v1/"
var ResourceManagerDefaultBasePath = "https://cloudresourcemanager.googleapis.com/v1/"
var RuntimeConfigDefaultBasePath = "https://runtimeconfig.googleapis.com/v1beta1/"
var SecretManagerDefaultBasePath = "https://secretmanager.googleapis.com/v1/"
var SecurityCenterDefaultBasePath = "https://securitycenter.googleapis.com/v1/"
var ServiceManagementDefaultBasePath = "https://servicemanagement.googleapis.com/v1/"
var ServiceUsageDefaultBasePath = "https://serviceusage.googleapis.com/v1/"
var SourceRepoDefaultBasePath = "https://sourcerepo.googleapis.com/v1/"
var SpannerDefaultBasePath = "https://spanner.googleapis.com/v1/"
var SQLDefaultBasePath = "https://sqladmin.googleapis.com/sql/v1beta4/"
var StorageDefaultBasePath = "https://storage.googleapis.com/storage/v1/"
var TPUDefaultBasePath = "https://tpu.googleapis.com/v1/"
var VPCAccessDefaultBasePath = "https://vpcaccess.googleapis.com/v1/"

var DefaultClientScopes = []string{
	"https://www.googleapis.com/auth/compute",
	"https://www.googleapis.com/auth/cloud-platform",
	"https://www.googleapis.com/auth/cloud-identity",
	"https://www.googleapis.com/auth/ndev.clouddns.readwrite",
	"https://www.googleapis.com/auth/devstorage.full_control",
	"https://www.googleapis.com/auth/userinfo.email",
}

func (c *Config) LoadAndValidate(ctx context.Context) error {
	if len(c.Scopes) == 0 {
		c.Scopes = DefaultClientScopes
	}

	c.context = ctx

	tokenSource, err := c.getTokenSource(c.Scopes)
	if err != nil {
		return err
	}
	c.tokenSource = tokenSource

	cleanCtx := context.WithValue(ctx, oauth2.HTTPClient, cleanhttp.DefaultClient())

	// 1. OAUTH2 TRANSPORT/CLIENT - sets up proper auth headers
	client := oauth2.NewClient(cleanCtx, tokenSource)

	// 2. Logging Transport - ensure we log HTTP requests to GCP APIs.
	loggingTransport := logging.NewTransport("Google", client.Transport)

	// 3. Retry Transport - retries common temporary errors
	// Keep order for wrapping logging so we log each retried request as well.
	// This value should be used if needed to create shallow copies with additional retry predicates.
	// See ClientWithAdditionalRetries
	retryTransport := NewTransportWithDefaultRetries(loggingTransport)

	// Set final transport value.
	client.Transport = retryTransport

	// This timeout is a timeout per HTTP request, not per logical operation.
	client.Timeout = c.synchronousTimeout()

	c.client = client
	c.context = ctx

	// This base path and some others below need the version and possibly more of the path
	// set on them. The client libraries are inconsistent about which values they need;
	// while most only want the host URL, some older ones also want the version and some
	// of those "projects" as well. You can find out if this is required by looking at
	// the basePath value in the client library file.
	computeClientBasePath := c.ComputeBasePath + "projects/"
	log.Printf("[INFO] Instantiating GCE client for path %s", computeClientBasePath)
	c.clientCompute, err = compute.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientCompute.UserAgent = c.userAgent
	c.clientCompute.BasePath = computeClientBasePath

	computeBetaClientBasePath := c.ComputeBetaBasePath + "projects/"
	log.Printf("[INFO] Instantiating GCE Beta client for path %s", computeBetaClientBasePath)
	c.clientComputeBeta, err = computeBeta.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientComputeBeta.UserAgent = c.userAgent
	c.clientComputeBeta.BasePath = computeBetaClientBasePath

	containerClientBasePath := removeBasePathVersion(c.ContainerBasePath)
	log.Printf("[INFO] Instantiating GKE client for path %s", containerClientBasePath)
	c.clientContainer, err = container.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientContainer.UserAgent = c.userAgent
	c.clientContainer.BasePath = containerClientBasePath

	containerBetaClientBasePath := removeBasePathVersion(c.ContainerBetaBasePath)
	log.Printf("[INFO] Instantiating GKE Beta client for path %s", containerBetaClientBasePath)
	c.clientContainerBeta, err = containerBeta.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientContainerBeta.UserAgent = c.userAgent
	c.clientContainerBeta.BasePath = containerBetaClientBasePath

	dnsClientBasePath := removeBasePathVersion(c.DNSBasePath)
	dnsClientBasePath = strings.ReplaceAll(dnsClientBasePath, "/dns/", "")
	log.Printf("[INFO] Instantiating Google Cloud DNS client for path %s", dnsClientBasePath)
	c.clientDns, err = dns.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientDns.UserAgent = c.userAgent
	c.clientDns.BasePath = dnsClientBasePath

	dnsBetaClientBasePath := removeBasePathVersion(c.DnsBetaBasePath)
	dnsBetaClientBasePath = strings.ReplaceAll(dnsBetaClientBasePath, "/dns/", "")
	log.Printf("[INFO] Instantiating Google Cloud DNS Beta client for path %s", dnsBetaClientBasePath)
	c.clientDnsBeta, err = dnsBeta.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientDnsBeta.UserAgent = c.userAgent
	c.clientDnsBeta.BasePath = dnsBetaClientBasePath

	kmsClientBasePath := removeBasePathVersion(c.KMSBasePath)
	log.Printf("[INFO] Instantiating Google Cloud KMS client for path %s", kmsClientBasePath)
	c.clientKms, err = cloudkms.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientKms.UserAgent = c.userAgent
	c.clientKms.BasePath = kmsClientBasePath

	loggingClientBasePath := removeBasePathVersion(c.LoggingBasePath)
	log.Printf("[INFO] Instantiating Google Stackdriver Logging client for path %s", loggingClientBasePath)
	c.clientLogging, err = cloudlogging.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientLogging.UserAgent = c.userAgent
	c.clientLogging.BasePath = loggingClientBasePath

	storageClientBasePath := c.StorageBasePath
	log.Printf("[INFO] Instantiating Google Storage client for path %s", storageClientBasePath)
	c.clientStorage, err = storage.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientStorage.UserAgent = c.userAgent
	c.clientStorage.BasePath = storageClientBasePath

	sqlClientBasePath := removeBasePathVersion(removeBasePathVersion(c.SQLBasePath))
	log.Printf("[INFO] Instantiating Google SqlAdmin client for path %s", sqlClientBasePath)
	c.clientSqlAdmin, err = sqladmin.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientSqlAdmin.UserAgent = c.userAgent
	c.clientSqlAdmin.BasePath = sqlClientBasePath

	pubsubClientBasePath := removeBasePathVersion(c.PubsubBasePath)
	log.Printf("[INFO] Instantiating Google Pubsub client for path %s", pubsubClientBasePath)
	wrappedPubsubClient := ClientWithAdditionalRetries(client, retryTransport, pubsubTopicProjectNotReady)
	c.wrappedPubsubClient = wrappedPubsubClient
	c.clientPubsub, err = pubsub.NewService(ctx, option.WithHTTPClient(wrappedPubsubClient))
	if err != nil {
		return err
	}
	c.clientPubsub.UserAgent = c.userAgent
	c.clientPubsub.BasePath = pubsubClientBasePath

	dataflowClientBasePath := removeBasePathVersion(c.DataflowBasePath)
	log.Printf("[INFO] Instantiating Google Dataflow client for path %s", dataflowClientBasePath)
	c.clientDataflow, err = dataflow.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientDataflow.UserAgent = c.userAgent
	c.clientDataflow.BasePath = dataflowClientBasePath

	resourceManagerBasePath := removeBasePathVersion(c.ResourceManagerBasePath)
	log.Printf("[INFO] Instantiating Google Cloud ResourceManager client for path %s", resourceManagerBasePath)
	c.clientResourceManager, err = cloudresourcemanager.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientResourceManager.UserAgent = c.userAgent
	c.clientResourceManager.BasePath = resourceManagerBasePath

	resourceManagerV2Beta1BasePath := removeBasePathVersion(c.ResourceManagerV2Beta1BasePath)
	log.Printf("[INFO] Instantiating Google Cloud ResourceManager V client for path %s", resourceManagerV2Beta1BasePath)
	c.clientResourceManagerV2Beta1, err = resourceManagerV2Beta1.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientResourceManagerV2Beta1.UserAgent = c.userAgent
	c.clientResourceManagerV2Beta1.BasePath = resourceManagerV2Beta1BasePath

	runtimeConfigClientBasePath := removeBasePathVersion(c.RuntimeConfigBasePath)
	log.Printf("[INFO] Instantiating Google Cloud Runtimeconfig client for path %s", runtimeConfigClientBasePath)
	c.clientRuntimeconfig, err = runtimeconfig.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientRuntimeconfig.UserAgent = c.userAgent
	c.clientRuntimeconfig.BasePath = runtimeConfigClientBasePath

	iamClientBasePath := removeBasePathVersion(c.IAMBasePath)
	log.Printf("[INFO] Instantiating Google Cloud IAM client for path %s", iamClientBasePath)
	c.clientIAM, err = iam.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientIAM.UserAgent = c.userAgent
	c.clientIAM.BasePath = iamClientBasePath

	iamCredentialsClientBasePath := removeBasePathVersion(c.IamCredentialsBasePath)
	log.Printf("[INFO] Instantiating Google Cloud IAMCredentials client for path %s", iamCredentialsClientBasePath)
	c.clientIamCredentials, err = iamcredentials.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientIamCredentials.UserAgent = c.userAgent
	c.clientIamCredentials.BasePath = iamCredentialsClientBasePath

	serviceManagementClientBasePath := removeBasePathVersion(c.ServiceManagementBasePath)
	log.Printf("[INFO] Instantiating Google Cloud Service Management client for path %s", serviceManagementClientBasePath)
	c.clientServiceMan, err = servicemanagement.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientServiceMan.UserAgent = c.userAgent
	c.clientServiceMan.BasePath = serviceManagementClientBasePath

	serviceUsageClientBasePath := removeBasePathVersion(c.ServiceUsageBasePath)
	log.Printf("[INFO] Instantiating Google Cloud Service Usage client for path %s", serviceUsageClientBasePath)
	c.clientServiceUsage, err = serviceusage.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientServiceUsage.UserAgent = c.userAgent
	c.clientServiceUsage.BasePath = serviceUsageClientBasePath

	cloudBillingClientBasePath := removeBasePathVersion(c.CloudBillingBasePath)
	log.Printf("[INFO] Instantiating Google Cloud Billing client for path %s", cloudBillingClientBasePath)
	c.clientBilling, err = cloudbilling.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientBilling.UserAgent = c.userAgent
	c.clientBilling.BasePath = cloudBillingClientBasePath

	cloudBuildClientBasePath := removeBasePathVersion(c.CloudBuildBasePath)
	log.Printf("[INFO] Instantiating Google Cloud Build client for path %s", cloudBuildClientBasePath)
	c.clientBuild, err = cloudbuild.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientBuild.UserAgent = c.userAgent
	c.clientBuild.BasePath = cloudBuildClientBasePath

	bigQueryClientBasePath := c.BigQueryBasePath
	log.Printf("[INFO] Instantiating Google Cloud BigQuery client for path %s", bigQueryClientBasePath)
	wrappedBigQueryClient := ClientWithAdditionalRetries(client, retryTransport, iamMemberMissing)
	c.wrappedBigQueryClient = wrappedBigQueryClient
	c.clientBigQuery, err = bigquery.NewService(ctx, option.WithHTTPClient(wrappedBigQueryClient))
	if err != nil {
		return err
	}
	c.clientBigQuery.UserAgent = c.userAgent
	c.clientBigQuery.BasePath = bigQueryClientBasePath

	cloudFunctionsClientBasePath := removeBasePathVersion(c.CloudFunctionsBasePath)
	log.Printf("[INFO] Instantiating Google Cloud CloudFunctions Client for path %s", cloudFunctionsClientBasePath)
	c.clientCloudFunctions, err = cloudfunctions.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientCloudFunctions.UserAgent = c.userAgent
	c.clientCloudFunctions.BasePath = cloudFunctionsClientBasePath

	c.bigtableClientFactory = &BigtableClientFactory{
		UserAgent:   c.userAgent,
		TokenSource: tokenSource,
	}

	bigtableAdminBasePath := removeBasePathVersion(c.BigtableAdminBasePath)
	log.Printf("[INFO] Instantiating Google Cloud BigtableAdmin for path %s", bigtableAdminBasePath)

	clientBigtable, err := bigtableadmin.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	clientBigtable.UserAgent = c.userAgent
	clientBigtable.BasePath = bigtableAdminBasePath
	c.clientBigtableProjectsInstances = bigtableadmin.NewProjectsInstancesService(clientBigtable)

	sourceRepoClientBasePath := removeBasePathVersion(c.SourceRepoBasePath)
	log.Printf("[INFO] Instantiating Google Cloud Source Repo client for path %s", sourceRepoClientBasePath)
	c.clientSourceRepo, err = sourcerepo.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientSourceRepo.UserAgent = c.userAgent
	c.clientSourceRepo.BasePath = sourceRepoClientBasePath

	spannerClientBasePath := removeBasePathVersion(c.SpannerBasePath)
	log.Printf("[INFO] Instantiating Google Cloud Spanner client for path %s", spannerClientBasePath)
	c.clientSpanner, err = spanner.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientSpanner.UserAgent = c.userAgent
	c.clientSpanner.BasePath = spannerClientBasePath

	dataprocClientBasePath := removeBasePathVersion(c.DataprocBasePath)
	log.Printf("[INFO] Instantiating Google Cloud Dataproc client for path %s", dataprocClientBasePath)
	c.clientDataproc, err = dataproc.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientDataproc.UserAgent = c.userAgent
	c.clientDataproc.BasePath = dataprocClientBasePath

	dataprocBetaClientBasePath := removeBasePathVersion(c.DataprocBetaBasePath)
	log.Printf("[INFO] Instantiating Google Cloud Dataproc Beta client for path %s", dataprocBetaClientBasePath)
	c.clientDataprocBeta, err = dataprocBeta.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientDataprocBeta.UserAgent = c.userAgent
	c.clientDataprocBeta.BasePath = dataprocClientBasePath

	filestoreClientBasePath := removeBasePathVersion(c.FilestoreBasePath)
	log.Printf("[INFO] Instantiating Filestore client for path %s", filestoreClientBasePath)
	c.clientFilestore, err = file.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientFilestore.UserAgent = c.userAgent
	c.clientFilestore.BasePath = filestoreClientBasePath

	cloudIoTClientBasePath := removeBasePathVersion(c.CloudIoTBasePath)
	log.Printf("[INFO] Instantiating Google Cloud IoT Core client for path %s", cloudIoTClientBasePath)
	c.clientCloudIoT, err = cloudiot.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientCloudIoT.UserAgent = c.userAgent
	c.clientCloudIoT.BasePath = cloudIoTClientBasePath

	appEngineClientBasePath := removeBasePathVersion(c.AppEngineBasePath)
	log.Printf("[INFO] Instantiating App Engine client for path %s", appEngineClientBasePath)
	c.clientAppEngine, err = appengine.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientAppEngine.UserAgent = c.userAgent
	c.clientAppEngine.BasePath = appEngineClientBasePath

	composerClientBasePath := removeBasePathVersion(c.ComposerBasePath)
	log.Printf("[INFO] Instantiating Cloud Composer client for path %s", composerClientBasePath)
	c.clientComposer, err = composer.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientComposer.UserAgent = c.userAgent
	c.clientComposer.BasePath = composerClientBasePath

	serviceNetworkingClientBasePath := removeBasePathVersion(c.ServiceNetworkingBasePath)
	log.Printf("[INFO] Instantiating Service Networking client for path %s", serviceNetworkingClientBasePath)
	c.clientServiceNetworking, err = servicenetworking.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientServiceNetworking.UserAgent = c.userAgent
	c.clientServiceNetworking.BasePath = serviceNetworkingClientBasePath

	storageTransferClientBasePath := removeBasePathVersion(c.StorageTransferBasePath)
	log.Printf("[INFO] Instantiating Google Cloud Storage Transfer client for path %s", storageTransferClientBasePath)
	c.clientStorageTransfer, err = storagetransfer.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientStorageTransfer.UserAgent = c.userAgent
	c.clientStorageTransfer.BasePath = storageTransferClientBasePath

	healthcareClientBasePath := removeBasePathVersion(c.HealthcareBasePath)
	log.Printf("[INFO] Instantiating Google Cloud Healthcare client for path %s", healthcareClientBasePath)

	c.clientHealthcare, err = healthcare.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientHealthcare.UserAgent = c.userAgent
	c.clientHealthcare.BasePath = healthcareClientBasePath

	c.Region = GetRegionFromRegionSelfLink(c.Region)

	c.requestBatcherServiceUsage = NewRequestBatcher("Service Usage", ctx, c.BatchingConfig)
	c.requestBatcherIam = NewRequestBatcher("IAM", ctx, c.BatchingConfig)

	c.PollInterval = 10 * time.Second

	return nil
}

func expandProviderBatchingConfig(v interface{}) (*batchingConfig, error) {
	config := &batchingConfig{
		sendAfter:      time.Second * defaultBatchSendIntervalSec,
		enableBatching: true,
	}

	if v == nil {
		return config, nil
	}
	ls := v.([]interface{})
	if len(ls) == 0 || ls[0] == nil {
		return config, nil
	}

	cfgV := ls[0].(map[string]interface{})
	if sendAfterV, ok := cfgV["send_after"]; ok {
		sendAfter, err := time.ParseDuration(sendAfterV.(string))
		if err != nil {
			return nil, fmt.Errorf("unable to parse duration from 'send_after' value %q", sendAfterV)
		}
		config.sendAfter = sendAfter
	}

	if enable, ok := cfgV["enable_batching"]; ok {
		config.enableBatching = enable.(bool)
	}

	return config, nil
}

func (c *Config) synchronousTimeout() time.Duration {
	if c.RequestTimeout == 0 {
		return 30 * time.Second
	}
	return c.RequestTimeout
}

func (c *Config) getTokenSource(clientScopes []string) (oauth2.TokenSource, error) {
	creds, err := c.GetCredentials(clientScopes)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	return creds.TokenSource, nil
}

// staticTokenSource is used to be able to identify static token sources without reflection.
type staticTokenSource struct {
	oauth2.TokenSource
}

func (c *Config) GetCredentials(clientScopes []string) (googleoauth.Credentials, error) {
	if c.AccessToken != "" {
		contents, _, err := pathOrContents(c.AccessToken)
		if err != nil {
			return googleoauth.Credentials{}, fmt.Errorf("Error loading access token: %s", err)
		}

		log.Printf("[INFO] Authenticating using configured Google JSON 'access_token'...")
		log.Printf("[INFO]   -- Scopes: %s", clientScopes)
		token := &oauth2.Token{AccessToken: contents}

		return googleoauth.Credentials{
			TokenSource: staticTokenSource{oauth2.StaticTokenSource(token)},
		}, nil
	}

	if c.Credentials != "" {
		contents, _, err := pathOrContents(c.Credentials)
		if err != nil {
			return googleoauth.Credentials{}, fmt.Errorf("error loading credentials: %s", err)
		}

		creds, err := googleoauth.CredentialsFromJSON(c.context, []byte(contents), clientScopes...)
		if err != nil {
			return googleoauth.Credentials{}, fmt.Errorf("unable to parse credentials from '%s': %s", contents, err)
		}

		log.Printf("[INFO] Authenticating using configured Google JSON 'credentials'...")
		log.Printf("[INFO]   -- Scopes: %s", clientScopes)
		return *creds, nil
	}

	log.Printf("[INFO] Authenticating using DefaultClient...")
	log.Printf("[INFO]   -- Scopes: %s", clientScopes)

	defaultTS, err := googleoauth.DefaultTokenSource(context.Background(), clientScopes...)
	if err != nil {
		return googleoauth.Credentials{}, fmt.Errorf("Attempted to load application default credentials since neither `credentials` nor `access_token` was set in the provider block.  No credentials loaded. To use your gcloud credentials, run 'gcloud auth application-default login'.  Original error: %w", err)
	}
	return googleoauth.Credentials{
		TokenSource: defaultTS,
	}, err
}

// Remove the `/{{version}}/` from a base path if present.
func removeBasePathVersion(url string) string {
	re := regexp.MustCompile(`(?P<base>http[s]://.*)(?P<version>/[^/]+?/$)`)
	return re.ReplaceAllString(url, "$1/")
}

// For a consumer of config.go that isn't a full fledged provider and doesn't
// have its own endpoint mechanism such as sweepers, init {{service}}BasePath
// values to a default. After using this, you should call config.LoadAndValidate.
func ConfigureBasePaths(c *Config) {
	// Generated Products
	c.AccessApprovalBasePath = AccessApprovalDefaultBasePath
	c.AccessContextManagerBasePath = AccessContextManagerDefaultBasePath
	c.ActiveDirectoryBasePath = ActiveDirectoryDefaultBasePath
	c.AppEngineBasePath = AppEngineDefaultBasePath
	c.BigQueryBasePath = BigQueryDefaultBasePath
	c.BigqueryDataTransferBasePath = BigqueryDataTransferDefaultBasePath
	c.BigtableBasePath = BigtableDefaultBasePath
	c.BinaryAuthorizationBasePath = BinaryAuthorizationDefaultBasePath
	c.CloudAssetBasePath = CloudAssetDefaultBasePath
	c.CloudBuildBasePath = CloudBuildDefaultBasePath
	c.CloudFunctionsBasePath = CloudFunctionsDefaultBasePath
	c.CloudIotBasePath = CloudIotDefaultBasePath
	c.CloudRunBasePath = CloudRunDefaultBasePath
	c.CloudSchedulerBasePath = CloudSchedulerDefaultBasePath
	c.CloudTasksBasePath = CloudTasksDefaultBasePath
	c.ComputeBasePath = ComputeDefaultBasePath
	c.ContainerAnalysisBasePath = ContainerAnalysisDefaultBasePath
	c.DataCatalogBasePath = DataCatalogDefaultBasePath
	c.DataLossPreventionBasePath = DataLossPreventionDefaultBasePath
	c.DataprocBasePath = DataprocDefaultBasePath
	c.DatastoreBasePath = DatastoreDefaultBasePath
	c.DeploymentManagerBasePath = DeploymentManagerDefaultBasePath
	c.DialogflowBasePath = DialogflowDefaultBasePath
	c.DNSBasePath = DNSDefaultBasePath
	c.FilestoreBasePath = FilestoreDefaultBasePath
	c.FirestoreBasePath = FirestoreDefaultBasePath
	c.GameServicesBasePath = GameServicesDefaultBasePath
	c.HealthcareBasePath = HealthcareDefaultBasePath
	c.IapBasePath = IapDefaultBasePath
	c.IdentityPlatformBasePath = IdentityPlatformDefaultBasePath
	c.KMSBasePath = KMSDefaultBasePath
	c.LoggingBasePath = LoggingDefaultBasePath
	c.MLEngineBasePath = MLEngineDefaultBasePath
	c.MonitoringBasePath = MonitoringDefaultBasePath
	c.NetworkManagementBasePath = NetworkManagementDefaultBasePath
	c.OSConfigBasePath = OSConfigDefaultBasePath
	c.OSLoginBasePath = OSLoginDefaultBasePath
	c.PubsubBasePath = PubsubDefaultBasePath
	c.RedisBasePath = RedisDefaultBasePath
	c.ResourceManagerBasePath = ResourceManagerDefaultBasePath
	c.RuntimeConfigBasePath = RuntimeConfigDefaultBasePath
	c.SecretManagerBasePath = SecretManagerDefaultBasePath
	c.SecurityCenterBasePath = SecurityCenterDefaultBasePath
	c.ServiceManagementBasePath = ServiceManagementDefaultBasePath
	c.ServiceUsageBasePath = ServiceUsageDefaultBasePath
	c.SourceRepoBasePath = SourceRepoDefaultBasePath
	c.SpannerBasePath = SpannerDefaultBasePath
	c.SQLBasePath = SQLDefaultBasePath
	c.StorageBasePath = StorageDefaultBasePath
	c.TPUBasePath = TPUDefaultBasePath
	c.VPCAccessBasePath = VPCAccessDefaultBasePath

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
	c.IAMBasePath = IAMDefaultBasePath
	c.ServiceNetworkingBasePath = ServiceNetworkingDefaultBasePath
	c.BigQueryBasePath = BigQueryDefaultBasePath
	c.StorageTransferBasePath = StorageTransferDefaultBasePath
	c.BigtableAdminBasePath = BigtableAdminDefaultBasePath
}
