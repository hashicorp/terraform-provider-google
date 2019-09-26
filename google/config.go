package google

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/logging"
	"github.com/hashicorp/terraform-plugin-sdk/helper/pathorcontents"
	"github.com/hashicorp/terraform-plugin-sdk/httpclient"
	"github.com/terraform-providers/terraform-provider-google/version"

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
	"google.golang.org/api/iam/v1"
	iamcredentials "google.golang.org/api/iamcredentials/v1"
	cloudlogging "google.golang.org/api/logging/v2"
	"google.golang.org/api/option"
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
	Region              string
	Zone                string
	Scopes              []string
	BatchingConfig      *batchingConfig
	UserProjectOverride bool

	client           *http.Client
	terraformVersion string
	userAgent        string

	tokenSource oauth2.TokenSource

	AccessContextManagerBasePath string
	AppEngineBasePath            string
	BigQueryBasePath             string
	BigqueryDataTransferBasePath string
	BigtableBasePath             string
	BinaryAuthorizationBasePath  string
	CloudBuildBasePath           string
	CloudFunctionsBasePath       string
	CloudSchedulerBasePath       string
	ComputeBasePath              string
	ContainerAnalysisBasePath    string
	DataprocBasePath             string
	DNSBasePath                  string
	FilestoreBasePath            string
	FirestoreBasePath            string
	IapBasePath                  string
	KMSBasePath                  string
	LoggingBasePath              string
	MLEngineBasePath             string
	MonitoringBasePath           string
	PubsubBasePath               string
	RedisBasePath                string
	ResourceManagerBasePath      string
	RuntimeConfigBasePath        string
	SecurityCenterBasePath       string
	SourceRepoBasePath           string
	SpannerBasePath              string
	SQLBasePath                  string
	StorageBasePath              string
	TPUBasePath                  string

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

	ServiceManagementBasePath string
	clientServiceMan          *servicemanagement.APIService

	ServiceUsageBasePath string
	clientServiceUsage   *serviceusage.Service

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
var AccessContextManagerDefaultBasePath = "https://accesscontextmanager.googleapis.com/v1/"
var AppEngineDefaultBasePath = "https://appengine.googleapis.com/v1/"
var BigQueryDefaultBasePath = "https://www.googleapis.com/bigquery/v2/"
var BigqueryDataTransferDefaultBasePath = "https://bigquerydatatransfer.googleapis.com/v1/"
var BigtableDefaultBasePath = "https://bigtableadmin.googleapis.com/v2/"
var BinaryAuthorizationDefaultBasePath = "https://binaryauthorization.googleapis.com/v1/"
var CloudBuildDefaultBasePath = "https://cloudbuild.googleapis.com/v1/"
var CloudFunctionsDefaultBasePath = "https://cloudfunctions.googleapis.com/v1/"
var CloudSchedulerDefaultBasePath = "https://cloudscheduler.googleapis.com/v1/"
var ComputeDefaultBasePath = "https://www.googleapis.com/compute/v1/"
var ContainerAnalysisDefaultBasePath = "https://containeranalysis.googleapis.com/v1/"
var DataprocDefaultBasePath = "https://dataproc.googleapis.com/v1/"
var DNSDefaultBasePath = "https://www.googleapis.com/dns/v1/"
var FilestoreDefaultBasePath = "https://file.googleapis.com/v1/"
var FirestoreDefaultBasePath = "https://firestore.googleapis.com/v1/"
var IapDefaultBasePath = "https://iap.googleapis.com/v1/"
var KMSDefaultBasePath = "https://cloudkms.googleapis.com/v1/"
var LoggingDefaultBasePath = "https://logging.googleapis.com/v2/"
var MLEngineDefaultBasePath = "https://ml.googleapis.com/v1/"
var MonitoringDefaultBasePath = "https://monitoring.googleapis.com/v3/"
var PubsubDefaultBasePath = "https://pubsub.googleapis.com/v1/"
var RedisDefaultBasePath = "https://redis.googleapis.com/v1/"
var ResourceManagerDefaultBasePath = "https://cloudresourcemanager.googleapis.com/v1/"
var RuntimeConfigDefaultBasePath = "https://runtimeconfig.googleapis.com/v1beta1/"
var SecurityCenterDefaultBasePath = "https://securitycenter.googleapis.com/v1/"
var SourceRepoDefaultBasePath = "https://sourcerepo.googleapis.com/v1/"
var SpannerDefaultBasePath = "https://spanner.googleapis.com/v1/"
var SQLDefaultBasePath = "https://www.googleapis.com/sql/v1beta4/"
var StorageDefaultBasePath = "https://www.googleapis.com/storage/v1/"
var TPUDefaultBasePath = "https://tpu.googleapis.com/v1/"

var defaultClientScopes = []string{
	"https://www.googleapis.com/auth/compute",
	"https://www.googleapis.com/auth/cloud-platform",
	"https://www.googleapis.com/auth/ndev.clouddns.readwrite",
	"https://www.googleapis.com/auth/devstorage.full_control",
}

func (c *Config) LoadAndValidate() error {
	if len(c.Scopes) == 0 {
		c.Scopes = defaultClientScopes
	}

	tokenSource, err := c.getTokenSource(c.Scopes)
	if err != nil {
		return err
	}
	c.tokenSource = tokenSource

	client := oauth2.NewClient(context.Background(), tokenSource)
	client.Transport = logging.NewTransport("Google", client.Transport)
	// Each individual request should return within 30s - timeouts will be retried.
	// This is a timeout for, e.g. a single GET request of an operation - not a
	// timeout for the maximum amount of time a logical request can take.
	client.Timeout, _ = time.ParseDuration("30s")

	tfUserAgent := httpclient.TerraformUserAgent(c.terraformVersion)
	providerVersion := fmt.Sprintf("terraform-provider-google/%s", version.ProviderVersion)
	userAgent := fmt.Sprintf("%s %s", tfUserAgent, providerVersion)

	c.client = client
	c.userAgent = userAgent

	context := context.Background()

	// This base path and some others below need the version and possibly more of the path
	// set on them. The client libraries are inconsistent about which values they need;
	// while most only want the host URL, some older ones also want the version and some
	// of those "projects" as well. You can find out if this is required by looking at
	// the basePath value in the client library file.
	computeClientBasePath := removeBasePathVersion(c.ComputeBasePath) + "v1/projects/"
	log.Printf("[INFO] Instantiating GCE client for path %s", computeClientBasePath)
	c.clientCompute, err = compute.NewService(context, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientCompute.UserAgent = userAgent
	c.clientCompute.BasePath = computeClientBasePath

	computeBetaClientBasePath := removeBasePathVersion(c.ComputeBetaBasePath) + "beta/projects/"
	log.Printf("[INFO] Instantiating GCE Beta client for path %s", computeBetaClientBasePath)
	c.clientComputeBeta, err = computeBeta.NewService(context, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientComputeBeta.UserAgent = userAgent
	c.clientComputeBeta.BasePath = computeBetaClientBasePath

	containerClientBasePath := removeBasePathVersion(c.ContainerBasePath)
	log.Printf("[INFO] Instantiating GKE client for path %s", containerClientBasePath)
	c.clientContainer, err = container.NewService(context, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientContainer.UserAgent = userAgent
	c.clientContainer.BasePath = containerClientBasePath

	containerBetaClientBasePath := removeBasePathVersion(c.ContainerBetaBasePath)
	log.Printf("[INFO] Instantiating GKE Beta client for path %s", containerBetaClientBasePath)
	c.clientContainerBeta, err = containerBeta.NewService(context, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientContainerBeta.UserAgent = userAgent
	c.clientContainerBeta.BasePath = containerBetaClientBasePath

	dnsClientBasePath := removeBasePathVersion(c.DNSBasePath) + "v1/projects/"
	log.Printf("[INFO] Instantiating Google Cloud DNS client for path %s", dnsClientBasePath)
	c.clientDns, err = dns.NewService(context, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientDns.UserAgent = userAgent
	c.clientDns.BasePath = dnsClientBasePath

	dnsBetaClientBasePath := removeBasePathVersion(c.DnsBetaBasePath) + "v1beta2/projects/"
	log.Printf("[INFO] Instantiating Google Cloud DNS Beta client for path %s", dnsBetaClientBasePath)
	c.clientDnsBeta, err = dnsBeta.NewService(context, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientDnsBeta.UserAgent = userAgent
	c.clientDnsBeta.BasePath = dnsBetaClientBasePath

	kmsClientBasePath := removeBasePathVersion(c.KMSBasePath)
	log.Printf("[INFO] Instantiating Google Cloud KMS client for path %s", kmsClientBasePath)
	c.clientKms, err = cloudkms.NewService(context, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientKms.UserAgent = userAgent
	c.clientKms.BasePath = kmsClientBasePath

	loggingClientBasePath := removeBasePathVersion(c.LoggingBasePath)
	log.Printf("[INFO] Instantiating Google Stackdriver Logging client for path %s", loggingClientBasePath)
	c.clientLogging, err = cloudlogging.NewService(context, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientLogging.UserAgent = userAgent
	c.clientLogging.BasePath = loggingClientBasePath

	storageClientBasePath := removeBasePathVersion(c.StorageBasePath) + "v1/"
	log.Printf("[INFO] Instantiating Google Storage client for path %s", storageClientBasePath)
	c.clientStorage, err = storage.NewService(context, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientStorage.UserAgent = userAgent
	c.clientStorage.BasePath = storageClientBasePath

	sqlClientBasePath := removeBasePathVersion(c.SQLBasePath) + "v1beta4/"
	log.Printf("[INFO] Instantiating Google SqlAdmin client for path %s", sqlClientBasePath)
	c.clientSqlAdmin, err = sqladmin.NewService(context, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientSqlAdmin.UserAgent = userAgent
	c.clientSqlAdmin.BasePath = sqlClientBasePath

	pubsubClientBasePath := removeBasePathVersion(c.PubsubBasePath)
	log.Printf("[INFO] Instantiating Google Pubsub client for path %s", pubsubClientBasePath)
	c.clientPubsub, err = pubsub.NewService(context, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientPubsub.UserAgent = userAgent
	c.clientPubsub.BasePath = pubsubClientBasePath

	dataflowClientBasePath := removeBasePathVersion(c.DataflowBasePath)
	log.Printf("[INFO] Instantiating Google Dataflow client for path %s", dataflowClientBasePath)
	c.clientDataflow, err = dataflow.NewService(context, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientDataflow.UserAgent = userAgent
	c.clientDataflow.BasePath = dataflowClientBasePath

	resourceManagerBasePath := removeBasePathVersion(c.ResourceManagerBasePath)
	log.Printf("[INFO] Instantiating Google Cloud ResourceManager client for path %s", resourceManagerBasePath)
	c.clientResourceManager, err = cloudresourcemanager.NewService(context, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientResourceManager.UserAgent = userAgent
	c.clientResourceManager.BasePath = resourceManagerBasePath

	resourceManagerV2Beta1BasePath := removeBasePathVersion(c.ResourceManagerV2Beta1BasePath)
	log.Printf("[INFO] Instantiating Google Cloud ResourceManager V client for path %s", resourceManagerV2Beta1BasePath)
	c.clientResourceManagerV2Beta1, err = resourceManagerV2Beta1.NewService(context, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientResourceManagerV2Beta1.UserAgent = userAgent
	c.clientResourceManagerV2Beta1.BasePath = resourceManagerV2Beta1BasePath

	runtimeConfigClientBasePath := removeBasePathVersion(c.RuntimeConfigBasePath)
	log.Printf("[INFO] Instantiating Google Cloud Runtimeconfig client for path %s", runtimeConfigClientBasePath)
	c.clientRuntimeconfig, err = runtimeconfig.NewService(context, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientRuntimeconfig.UserAgent = userAgent
	c.clientRuntimeconfig.BasePath = runtimeConfigClientBasePath

	iamClientBasePath := removeBasePathVersion(c.IAMBasePath)
	log.Printf("[INFO] Instantiating Google Cloud IAM client for path %s", iamClientBasePath)
	c.clientIAM, err = iam.NewService(context, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientIAM.UserAgent = userAgent
	c.clientIAM.BasePath = iamClientBasePath

	iamCredentialsClientBasePath := removeBasePathVersion(c.IamCredentialsBasePath)
	log.Printf("[INFO] Instantiating Google Cloud IAMCredentials client for path %s", iamCredentialsClientBasePath)
	c.clientIamCredentials, err = iamcredentials.NewService(context, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientIamCredentials.UserAgent = userAgent
	c.clientIamCredentials.BasePath = iamCredentialsClientBasePath

	serviceManagementClientBasePath := removeBasePathVersion(c.ServiceManagementBasePath)
	log.Printf("[INFO] Instantiating Google Cloud Service Management client for path %s", serviceManagementClientBasePath)
	c.clientServiceMan, err = servicemanagement.NewService(context, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientServiceMan.UserAgent = userAgent
	c.clientServiceMan.BasePath = serviceManagementClientBasePath

	serviceUsageClientBasePath := removeBasePathVersion(c.ServiceUsageBasePath)
	log.Printf("[INFO] Instantiating Google Cloud Service Usage client for path %s", serviceUsageClientBasePath)
	c.clientServiceUsage, err = serviceusage.NewService(context, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientServiceUsage.UserAgent = userAgent
	c.clientServiceUsage.BasePath = serviceUsageClientBasePath

	cloudBillingClientBasePath := removeBasePathVersion(c.CloudBillingBasePath)
	log.Printf("[INFO] Instantiating Google Cloud Billing client for path %s", cloudBillingClientBasePath)
	c.clientBilling, err = cloudbilling.NewService(context, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientBilling.UserAgent = userAgent
	c.clientBilling.BasePath = cloudBillingClientBasePath

	cloudBuildClientBasePath := removeBasePathVersion(c.CloudBuildBasePath)
	log.Printf("[INFO] Instantiating Google Cloud Build client for path %s", cloudBuildClientBasePath)
	c.clientBuild, err = cloudbuild.NewService(context, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientBuild.UserAgent = userAgent
	c.clientBuild.BasePath = cloudBuildClientBasePath

	bigQueryClientBasePath := removeBasePathVersion(c.BigQueryBasePath) + "v2/"
	log.Printf("[INFO] Instantiating Google Cloud BigQuery client for path %s", bigQueryClientBasePath)
	c.clientBigQuery, err = bigquery.NewService(context, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientBigQuery.UserAgent = userAgent
	c.clientBigQuery.BasePath = bigQueryClientBasePath

	cloudFunctionsClientBasePath := removeBasePathVersion(c.CloudFunctionsBasePath)
	log.Printf("[INFO] Instantiating Google Cloud CloudFunctions Client for path %s", cloudFunctionsClientBasePath)
	c.clientCloudFunctions, err = cloudfunctions.NewService(context, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientCloudFunctions.UserAgent = userAgent
	c.clientCloudFunctions.BasePath = cloudFunctionsClientBasePath

	c.bigtableClientFactory = &BigtableClientFactory{
		UserAgent:   userAgent,
		TokenSource: tokenSource,
	}

	bigtableAdminBasePath := removeBasePathVersion(c.BigtableAdminBasePath)
	log.Printf("[INFO] Instantiating Google Cloud BigtableAdmin for path %s", bigtableAdminBasePath)

	clientBigtable, err := bigtableadmin.NewService(context, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	clientBigtable.UserAgent = userAgent
	clientBigtable.BasePath = bigtableAdminBasePath
	c.clientBigtableProjectsInstances = bigtableadmin.NewProjectsInstancesService(clientBigtable)

	sourceRepoClientBasePath := removeBasePathVersion(c.SourceRepoBasePath)
	log.Printf("[INFO] Instantiating Google Cloud Source Repo client for path %s", sourceRepoClientBasePath)
	c.clientSourceRepo, err = sourcerepo.NewService(context, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientSourceRepo.UserAgent = userAgent
	c.clientSourceRepo.BasePath = sourceRepoClientBasePath

	spannerClientBasePath := removeBasePathVersion(c.SpannerBasePath)
	log.Printf("[INFO] Instantiating Google Cloud Spanner client for path %s", spannerClientBasePath)
	c.clientSpanner, err = spanner.NewService(context, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientSpanner.UserAgent = userAgent
	c.clientSpanner.BasePath = spannerClientBasePath

	dataprocClientBasePath := removeBasePathVersion(c.DataprocBasePath)
	log.Printf("[INFO] Instantiating Google Cloud Dataproc client for path %s", dataprocClientBasePath)
	c.clientDataproc, err = dataproc.NewService(context, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientDataproc.UserAgent = userAgent
	c.clientDataproc.BasePath = dataprocClientBasePath

	dataprocBetaClientBasePath := removeBasePathVersion(c.DataprocBetaBasePath)
	log.Printf("[INFO] Instantiating Google Cloud Dataproc Beta client for path %s", dataprocBetaClientBasePath)
	c.clientDataprocBeta, err = dataprocBeta.NewService(context, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientDataprocBeta.UserAgent = userAgent
	c.clientDataprocBeta.BasePath = dataprocClientBasePath

	filestoreClientBasePath := removeBasePathVersion(c.FilestoreBasePath)
	log.Printf("[INFO] Instantiating Filestore client for path %s", filestoreClientBasePath)
	c.clientFilestore, err = file.NewService(context, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientFilestore.UserAgent = userAgent
	c.clientFilestore.BasePath = filestoreClientBasePath

	cloudIoTClientBasePath := removeBasePathVersion(c.CloudIoTBasePath)
	log.Printf("[INFO] Instantiating Google Cloud IoT Core client for path %s", cloudIoTClientBasePath)
	c.clientCloudIoT, err = cloudiot.NewService(context, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientCloudIoT.UserAgent = userAgent
	c.clientCloudIoT.BasePath = cloudIoTClientBasePath

	appEngineClientBasePath := removeBasePathVersion(c.AppEngineBasePath)
	log.Printf("[INFO] Instantiating App Engine client for path %s", appEngineClientBasePath)
	c.clientAppEngine, err = appengine.NewService(context, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientAppEngine.UserAgent = userAgent
	c.clientAppEngine.BasePath = appEngineClientBasePath

	composerClientBasePath := removeBasePathVersion(c.ComposerBasePath)
	log.Printf("[INFO] Instantiating Cloud Composer client for path %s", composerClientBasePath)
	c.clientComposer, err = composer.NewService(context, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientComposer.UserAgent = userAgent
	c.clientComposer.BasePath = composerClientBasePath

	serviceNetworkingClientBasePath := removeBasePathVersion(c.ServiceNetworkingBasePath)
	log.Printf("[INFO] Instantiating Service Networking client for path %s", serviceNetworkingClientBasePath)
	c.clientServiceNetworking, err = servicenetworking.NewService(context, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientServiceNetworking.UserAgent = userAgent
	c.clientServiceNetworking.BasePath = serviceNetworkingClientBasePath

	storageTransferClientBasePath := removeBasePathVersion(c.StorageTransferBasePath)
	log.Printf("[INFO] Instantiating Google Cloud Storage Transfer client for path %s", storageTransferClientBasePath)
	c.clientStorageTransfer, err = storagetransfer.NewService(context, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	c.clientStorageTransfer.UserAgent = userAgent
	c.clientStorageTransfer.BasePath = storageTransferClientBasePath

	c.Region = GetRegionFromRegionSelfLink(c.Region)

	c.requestBatcherServiceUsage = NewRequestBatcher("Service Usage", context, c.BatchingConfig)
	c.requestBatcherIam = NewRequestBatcher("IAM", context, c.BatchingConfig)

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

func (c *Config) getTokenSource(clientScopes []string) (oauth2.TokenSource, error) {
	if c.AccessToken != "" {
		contents, _, err := pathorcontents.Read(c.AccessToken)
		if err != nil {
			return nil, fmt.Errorf("Error loading access token: %s", err)
		}

		log.Printf("[INFO] Authenticating using configured Google JSON 'access_token'...")
		log.Printf("[INFO]   -- Scopes: %s", clientScopes)
		token := &oauth2.Token{AccessToken: contents}
		return oauth2.StaticTokenSource(token), nil
	}

	if c.Credentials != "" {
		contents, _, err := pathorcontents.Read(c.Credentials)
		if err != nil {
			return nil, fmt.Errorf("Error loading credentials: %s", err)
		}

		creds, err := googleoauth.CredentialsFromJSON(context.Background(), []byte(contents), clientScopes...)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse credentials from '%s': %s", contents, err)
		}

		log.Printf("[INFO] Authenticating using configured Google JSON 'credentials'...")
		log.Printf("[INFO]   -- Scopes: %s", clientScopes)
		return creds.TokenSource, nil
	}

	log.Printf("[INFO] Authenticating using DefaultClient...")
	log.Printf("[INFO]   -- Scopes: %s", clientScopes)
	return googleoauth.DefaultTokenSource(context.Background(), clientScopes...)
}

// Remove the `/{{version}}/` from a base path, replacing it with `/`
func removeBasePathVersion(url string) string {
	return regexp.MustCompile(`/[^/]+/$`).ReplaceAllString(url, "/")
}

// For a consumer of config.go that isn't a full fledged provider and doesn't
// have its own endpoint mechanism such as sweepers, init {{service}}BasePath
// values to a default. After using this, you should call config.LoadAndValidate.
func ConfigureBasePaths(c *Config) {
	// Generated Products
	c.AccessContextManagerBasePath = AccessContextManagerDefaultBasePath
	c.AppEngineBasePath = AppEngineDefaultBasePath
	c.BigQueryBasePath = BigQueryDefaultBasePath
	c.BigqueryDataTransferBasePath = BigqueryDataTransferDefaultBasePath
	c.BigtableBasePath = BigtableDefaultBasePath
	c.BinaryAuthorizationBasePath = BinaryAuthorizationDefaultBasePath
	c.CloudBuildBasePath = CloudBuildDefaultBasePath
	c.CloudFunctionsBasePath = CloudFunctionsDefaultBasePath
	c.CloudSchedulerBasePath = CloudSchedulerDefaultBasePath
	c.ComputeBasePath = ComputeDefaultBasePath
	c.ContainerAnalysisBasePath = ContainerAnalysisDefaultBasePath
	c.DataprocBasePath = DataprocDefaultBasePath
	c.DNSBasePath = DNSDefaultBasePath
	c.FilestoreBasePath = FilestoreDefaultBasePath
	c.FirestoreBasePath = FirestoreDefaultBasePath
	c.IapBasePath = IapDefaultBasePath
	c.KMSBasePath = KMSDefaultBasePath
	c.LoggingBasePath = LoggingDefaultBasePath
	c.MLEngineBasePath = MLEngineDefaultBasePath
	c.MonitoringBasePath = MonitoringDefaultBasePath
	c.PubsubBasePath = PubsubDefaultBasePath
	c.RedisBasePath = RedisDefaultBasePath
	c.ResourceManagerBasePath = ResourceManagerDefaultBasePath
	c.RuntimeConfigBasePath = RuntimeConfigDefaultBasePath
	c.SecurityCenterBasePath = SecurityCenterDefaultBasePath
	c.SourceRepoBasePath = SourceRepoDefaultBasePath
	c.SpannerBasePath = SpannerDefaultBasePath
	c.SQLBasePath = SQLDefaultBasePath
	c.StorageBasePath = StorageDefaultBasePath
	c.TPUBasePath = TPUDefaultBasePath

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
	c.ServiceManagementBasePath = ServiceManagementDefaultBasePath
	c.ServiceNetworkingBasePath = ServiceNetworkingDefaultBasePath
	c.ServiceUsageBasePath = ServiceUsageDefaultBasePath
	c.BigQueryBasePath = BigQueryDefaultBasePath
	c.CloudIoTBasePath = CloudIoTDefaultBasePath
	c.StorageTransferBasePath = StorageTransferDefaultBasePath
	c.BigtableAdminBasePath = BigtableAdminDefaultBasePath
}
