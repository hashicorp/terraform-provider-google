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
	"google.golang.org/api/cloudidentity/v1"
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
	"google.golang.org/api/transport"
)

type providerMeta struct {
	ModuleName string `cty:"module_name"`
}

// Config is the configuration structure used to instantiate the Google
// provider.
type Config struct {
	AccessToken                        string
	Credentials                        string
	ImpersonateServiceAccount          string
	ImpersonateServiceAccountDelegates []string
	Project                            string
	Region                             string
	BillingProject                     string
	Zone                               string
	Scopes                             []string
	BatchingConfig                     *batchingConfig
	UserProjectOverride                bool
	RequestTimeout                     time.Duration
	// PollInterval is passed to resource.StateChangeConf in common_operation.go
	// It controls the interval at which we poll for successful operations
	PollInterval time.Duration

	client    *http.Client
	context   context.Context
	userAgent string

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
	CloudIdentityBasePath        string
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

	CloudBillingBasePath           string
	ComposerBasePath               string
	ComputeBetaBasePath            string
	ContainerBasePath              string
	ContainerBetaBasePath          string
	DataprocBetaBasePath           string
	DataflowBasePath               string
	DnsBetaBasePath                string
	IamCredentialsBasePath         string
	ResourceManagerV2Beta1BasePath string
	IAMBasePath                    string
	CloudIoTBasePath               string
	ServiceNetworkingBasePath      string
	StorageTransferBasePath        string
	BigtableAdminBasePath          string

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
var CloudIdentityDefaultBasePath = "https://cloudidentity.googleapis.com/v1/"
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

// Methods to create new services from config
// Some base paths below need the version and possibly more of the path
// set on them. The client libraries are inconsistent about which values they need;
// while most only want the host URL, some older ones also want the version and some
// of those "projects" as well. You can find out if this is required by looking at
// the basePath value in the client library file.
func (c *Config) NewComputeClient(userAgent string) *compute.Service {
	computeClientBasePath := c.ComputeBasePath + "projects/"
	log.Printf("[INFO] Instantiating GCE client for path %s", computeClientBasePath)
	clientCompute, err := compute.NewService(c.context, option.WithHTTPClient(c.client))
	if err != nil {
		log.Printf("[WARN] Error creating client compute: %s", err)
		return nil
	}
	clientCompute.UserAgent = userAgent
	clientCompute.BasePath = computeClientBasePath

	return clientCompute
}

func (c *Config) NewComputeBetaClient(userAgent string) *computeBeta.Service {
	computeBetaClientBasePath := c.ComputeBetaBasePath + "projects/"
	log.Printf("[INFO] Instantiating GCE Beta client for path %s", computeBetaClientBasePath)
	clientComputeBeta, err := computeBeta.NewService(c.context, option.WithHTTPClient(c.client))
	if err != nil {
		log.Printf("[WARN] Error creating client compute beta: %s", err)
		return nil
	}
	clientComputeBeta.UserAgent = userAgent
	clientComputeBeta.BasePath = computeBetaClientBasePath

	return clientComputeBeta
}

func (c *Config) NewContainerClient(userAgent string) *container.Service {
	containerClientBasePath := removeBasePathVersion(c.ContainerBasePath)
	log.Printf("[INFO] Instantiating GKE client for path %s", containerClientBasePath)
	clientContainer, err := container.NewService(c.context, option.WithHTTPClient(c.client))
	if err != nil {
		log.Printf("[WARN] Error creating client container: %s", err)
		return nil
	}
	clientContainer.UserAgent = userAgent
	clientContainer.BasePath = containerClientBasePath

	return clientContainer
}

func (c *Config) NewContainerBetaClient(userAgent string) *containerBeta.Service {
	containerBetaClientBasePath := removeBasePathVersion(c.ContainerBetaBasePath)
	log.Printf("[INFO] Instantiating GKE Beta client for path %s", containerBetaClientBasePath)
	clientContainerBeta, err := containerBeta.NewService(c.context, option.WithHTTPClient(c.client))
	if err != nil {
		log.Printf("[WARN] Error creating client container beta: %s", err)
		return nil
	}
	clientContainerBeta.UserAgent = userAgent
	clientContainerBeta.BasePath = containerBetaClientBasePath

	return clientContainerBeta
}

func (c *Config) NewDnsClient(userAgent string) *dns.Service {
	dnsClientBasePath := removeBasePathVersion(c.DNSBasePath)
	dnsClientBasePath = strings.ReplaceAll(dnsClientBasePath, "/dns/", "")
	log.Printf("[INFO] Instantiating Google Cloud DNS client for path %s", dnsClientBasePath)
	clientDns, err := dns.NewService(c.context, option.WithHTTPClient(c.client))
	if err != nil {
		log.Printf("[WARN] Error creating client dns: %s", err)
		return nil
	}
	clientDns.UserAgent = userAgent
	clientDns.BasePath = dnsClientBasePath

	return clientDns
}

func (c *Config) NewDnsBetaClient(userAgent string) *dnsBeta.Service {
	dnsBetaClientBasePath := removeBasePathVersion(c.DnsBetaBasePath)
	dnsBetaClientBasePath = strings.ReplaceAll(dnsBetaClientBasePath, "/dns/", "")
	log.Printf("[INFO] Instantiating Google Cloud DNS Beta client for path %s", dnsBetaClientBasePath)
	clientDnsBeta, err := dnsBeta.NewService(c.context, option.WithHTTPClient(c.client))
	if err != nil {
		log.Printf("[WARN] Error creating client dns beta: %s", err)
		return nil
	}
	clientDnsBeta.UserAgent = userAgent
	clientDnsBeta.BasePath = dnsBetaClientBasePath

	return clientDnsBeta
}

func (c *Config) NewKmsClient(userAgent string) *cloudkms.Service {
	kmsClientBasePath := removeBasePathVersion(c.KMSBasePath)
	log.Printf("[INFO] Instantiating Google Cloud KMS client for path %s", kmsClientBasePath)
	clientKms, err := cloudkms.NewService(c.context, option.WithHTTPClient(c.client))
	if err != nil {
		log.Printf("[WARN] Error creating client kms: %s", err)
		return nil
	}
	clientKms.UserAgent = userAgent
	clientKms.BasePath = kmsClientBasePath

	return clientKms
}

func (c *Config) NewLoggingClient(userAgent string) *cloudlogging.Service {
	loggingClientBasePath := removeBasePathVersion(c.LoggingBasePath)
	log.Printf("[INFO] Instantiating Google Stackdriver Logging client for path %s", loggingClientBasePath)
	clientLogging, err := cloudlogging.NewService(c.context, option.WithHTTPClient(c.client))
	if err != nil {
		log.Printf("[WARN] Error creating client logging: %s", err)
		return nil
	}
	clientLogging.UserAgent = userAgent
	clientLogging.BasePath = loggingClientBasePath

	return clientLogging
}

func (c *Config) NewStorageClient(userAgent string) *storage.Service {
	storageClientBasePath := c.StorageBasePath
	log.Printf("[INFO] Instantiating Google Storage client for path %s", storageClientBasePath)
	clientStorage, err := storage.NewService(c.context, option.WithHTTPClient(c.client))
	if err != nil {
		log.Printf("[WARN] Error creating client storage: %s", err)
		return nil
	}
	clientStorage.UserAgent = userAgent
	clientStorage.BasePath = storageClientBasePath

	return clientStorage
}

func (c *Config) NewSqlAdminClient(userAgent string) *sqladmin.Service {
	sqlClientBasePath := removeBasePathVersion(removeBasePathVersion(c.SQLBasePath))
	log.Printf("[INFO] Instantiating Google SqlAdmin client for path %s", sqlClientBasePath)
	clientSqlAdmin, err := sqladmin.NewService(c.context, option.WithHTTPClient(c.client))
	if err != nil {
		log.Printf("[WARN] Error creating client storage: %s", err)
		return nil
	}
	clientSqlAdmin.UserAgent = userAgent
	clientSqlAdmin.BasePath = sqlClientBasePath

	return clientSqlAdmin
}

func (c *Config) NewPubsubClient(userAgent string) *pubsub.Service {
	pubsubClientBasePath := removeBasePathVersion(c.PubsubBasePath)
	log.Printf("[INFO] Instantiating Google Pubsub client for path %s", pubsubClientBasePath)
	wrappedPubsubClient := ClientWithAdditionalRetries(c.client, pubsubTopicProjectNotReady)
	clientPubsub, err := pubsub.NewService(c.context, option.WithHTTPClient(wrappedPubsubClient))
	if err != nil {
		log.Printf("[WARN] Error creating client pubsub: %s", err)
		return nil
	}
	clientPubsub.UserAgent = userAgent
	clientPubsub.BasePath = pubsubClientBasePath

	return clientPubsub
}

func (c *Config) NewDataflowClient(userAgent string) *dataflow.Service {
	dataflowClientBasePath := removeBasePathVersion(c.DataflowBasePath)
	log.Printf("[INFO] Instantiating Google Dataflow client for path %s", dataflowClientBasePath)
	clientDataflow, err := dataflow.NewService(c.context, option.WithHTTPClient(c.client))
	if err != nil {
		log.Printf("[WARN] Error creating client dataflow: %s", err)
		return nil
	}
	clientDataflow.UserAgent = userAgent
	clientDataflow.BasePath = dataflowClientBasePath

	return clientDataflow
}

func (c *Config) NewResourceManagerClient(userAgent string) *cloudresourcemanager.Service {
	resourceManagerBasePath := removeBasePathVersion(c.ResourceManagerBasePath)
	log.Printf("[INFO] Instantiating Google Cloud ResourceManager client for path %s", resourceManagerBasePath)
	clientResourceManager, err := cloudresourcemanager.NewService(c.context, option.WithHTTPClient(c.client))
	if err != nil {
		log.Printf("[WARN] Error creating client resource manager: %s", err)
		return nil
	}
	clientResourceManager.UserAgent = userAgent
	clientResourceManager.BasePath = resourceManagerBasePath

	return clientResourceManager
}

func (c *Config) NewResourceManagerV2Beta1Client(userAgent string) *resourceManagerV2Beta1.Service {
	resourceManagerV2Beta1BasePath := removeBasePathVersion(c.ResourceManagerV2Beta1BasePath)
	log.Printf("[INFO] Instantiating Google Cloud ResourceManager V client for path %s", resourceManagerV2Beta1BasePath)
	clientResourceManagerV2Beta1, err := resourceManagerV2Beta1.NewService(c.context, option.WithHTTPClient(c.client))
	if err != nil {
		log.Printf("[WARN] Error creating client resource manager v2beta1: %s", err)
		return nil
	}
	clientResourceManagerV2Beta1.UserAgent = userAgent
	clientResourceManagerV2Beta1.BasePath = resourceManagerV2Beta1BasePath

	return clientResourceManagerV2Beta1
}

func (c *Config) NewRuntimeconfigClient(userAgent string) *runtimeconfig.Service {
	runtimeConfigClientBasePath := removeBasePathVersion(c.RuntimeConfigBasePath)
	log.Printf("[INFO] Instantiating Google Cloud Runtimeconfig client for path %s", runtimeConfigClientBasePath)
	clientRuntimeconfig, err := runtimeconfig.NewService(c.context, option.WithHTTPClient(c.client))
	if err != nil {
		log.Printf("[WARN] Error creating client resource manager v2beta1: %s", err)
		return nil
	}
	clientRuntimeconfig.UserAgent = userAgent
	clientRuntimeconfig.BasePath = runtimeConfigClientBasePath

	return clientRuntimeconfig
}

func (c *Config) NewIamClient(userAgent string) *iam.Service {
	iamClientBasePath := removeBasePathVersion(c.IAMBasePath)
	log.Printf("[INFO] Instantiating Google Cloud IAM client for path %s", iamClientBasePath)
	clientIAM, err := iam.NewService(c.context, option.WithHTTPClient(c.client))
	if err != nil {
		log.Printf("[WARN] Error creating client iam: %s", err)
		return nil
	}
	clientIAM.UserAgent = userAgent
	clientIAM.BasePath = iamClientBasePath

	return clientIAM
}

func (c *Config) NewIamCredentialsClient(userAgent string) *iamcredentials.Service {
	iamCredentialsClientBasePath := removeBasePathVersion(c.IamCredentialsBasePath)
	log.Printf("[INFO] Instantiating Google Cloud IAMCredentials client for path %s", iamCredentialsClientBasePath)
	clientIamCredentials, err := iamcredentials.NewService(c.context, option.WithHTTPClient(c.client))
	if err != nil {
		log.Printf("[WARN] Error creating client iam credentials: %s", err)
		return nil
	}
	clientIamCredentials.UserAgent = userAgent
	clientIamCredentials.BasePath = iamCredentialsClientBasePath

	return clientIamCredentials
}

func (c *Config) NewServiceManClient(userAgent string) *servicemanagement.APIService {
	serviceManagementClientBasePath := removeBasePathVersion(c.ServiceManagementBasePath)
	log.Printf("[INFO] Instantiating Google Cloud Service Management client for path %s", serviceManagementClientBasePath)
	clientServiceMan, err := servicemanagement.NewService(c.context, option.WithHTTPClient(c.client))
	if err != nil {
		log.Printf("[WARN] Error creating client service management: %s", err)
		return nil
	}
	clientServiceMan.UserAgent = userAgent
	clientServiceMan.BasePath = serviceManagementClientBasePath

	return clientServiceMan
}

func (c *Config) NewServiceUsageClient(userAgent string) *serviceusage.Service {
	serviceUsageClientBasePath := removeBasePathVersion(c.ServiceUsageBasePath)
	log.Printf("[INFO] Instantiating Google Cloud Service Usage client for path %s", serviceUsageClientBasePath)
	clientServiceUsage, err := serviceusage.NewService(c.context, option.WithHTTPClient(c.client))
	if err != nil {
		log.Printf("[WARN] Error creating client service usage: %s", err)
		return nil
	}
	clientServiceUsage.UserAgent = userAgent
	clientServiceUsage.BasePath = serviceUsageClientBasePath

	return clientServiceUsage
}

func (c *Config) NewBillingClient(userAgent string) *cloudbilling.APIService {
	cloudBillingClientBasePath := removeBasePathVersion(c.CloudBillingBasePath)
	log.Printf("[INFO] Instantiating Google Cloud Billing client for path %s", cloudBillingClientBasePath)
	clientBilling, err := cloudbilling.NewService(c.context, option.WithHTTPClient(c.client))
	if err != nil {
		log.Printf("[WARN] Error creating client billing: %s", err)
		return nil
	}
	clientBilling.UserAgent = userAgent
	clientBilling.BasePath = cloudBillingClientBasePath

	return clientBilling
}

func (c *Config) NewBuildClient(userAgent string) *cloudbuild.Service {
	cloudBuildClientBasePath := removeBasePathVersion(c.CloudBuildBasePath)
	log.Printf("[INFO] Instantiating Google Cloud Build client for path %s", cloudBuildClientBasePath)
	clientBuild, err := cloudbuild.NewService(c.context, option.WithHTTPClient(c.client))
	if err != nil {
		log.Printf("[WARN] Error creating client build: %s", err)
		return nil
	}
	clientBuild.UserAgent = userAgent
	clientBuild.BasePath = cloudBuildClientBasePath

	return clientBuild
}

func (c *Config) NewCloudFunctionsClient(userAgent string) *cloudfunctions.Service {
	cloudFunctionsClientBasePath := removeBasePathVersion(c.CloudFunctionsBasePath)
	log.Printf("[INFO] Instantiating Google Cloud CloudFunctions Client for path %s", cloudFunctionsClientBasePath)
	clientCloudFunctions, err := cloudfunctions.NewService(c.context, option.WithHTTPClient(c.client))
	if err != nil {
		log.Printf("[WARN] Error creating client cloud functions: %s", err)
		return nil
	}
	clientCloudFunctions.UserAgent = userAgent
	clientCloudFunctions.BasePath = cloudFunctionsClientBasePath

	return clientCloudFunctions
}

func (c *Config) NewSourceRepoClient(userAgent string) *sourcerepo.Service {
	sourceRepoClientBasePath := removeBasePathVersion(c.SourceRepoBasePath)
	log.Printf("[INFO] Instantiating Google Cloud Source Repo client for path %s", sourceRepoClientBasePath)
	clientSourceRepo, err := sourcerepo.NewService(c.context, option.WithHTTPClient(c.client))
	if err != nil {
		log.Printf("[WARN] Error creating client source repo: %s", err)
		return nil
	}
	clientSourceRepo.UserAgent = userAgent
	clientSourceRepo.BasePath = sourceRepoClientBasePath

	return clientSourceRepo
}

func (c *Config) NewBigQueryClient(userAgent string) *bigquery.Service {
	bigQueryClientBasePath := c.BigQueryBasePath
	log.Printf("[INFO] Instantiating Google Cloud BigQuery client for path %s", bigQueryClientBasePath)
	wrappedBigQueryClient := ClientWithAdditionalRetries(c.client, iamMemberMissing)
	clientBigQuery, err := bigquery.NewService(c.context, option.WithHTTPClient(wrappedBigQueryClient))
	if err != nil {
		log.Printf("[WARN] Error creating client big query: %s", err)
		return nil
	}
	clientBigQuery.UserAgent = userAgent
	clientBigQuery.BasePath = bigQueryClientBasePath

	return clientBigQuery
}

func (c *Config) NewSpannerClient(userAgent string) *spanner.Service {
	spannerClientBasePath := removeBasePathVersion(c.SpannerBasePath)
	log.Printf("[INFO] Instantiating Google Cloud Spanner client for path %s", spannerClientBasePath)
	clientSpanner, err := spanner.NewService(c.context, option.WithHTTPClient(c.client))
	if err != nil {
		log.Printf("[WARN] Error creating client source repo: %s", err)
		return nil
	}
	clientSpanner.UserAgent = userAgent
	clientSpanner.BasePath = spannerClientBasePath

	return clientSpanner
}

func (c *Config) NewDataprocClient(userAgent string) *dataproc.Service {
	dataprocClientBasePath := removeBasePathVersion(c.DataprocBasePath)
	log.Printf("[INFO] Instantiating Google Cloud Dataproc client for path %s", dataprocClientBasePath)
	clientDataproc, err := dataproc.NewService(c.context, option.WithHTTPClient(c.client))
	if err != nil {
		log.Printf("[WARN] Error creating client dataproc: %s", err)
		return nil
	}
	clientDataproc.UserAgent = userAgent
	clientDataproc.BasePath = dataprocClientBasePath

	return clientDataproc
}

func (c *Config) NewDataprocBetaClient(userAgent string) *dataprocBeta.Service {
	dataprocBetaClientBasePath := removeBasePathVersion(c.DataprocBetaBasePath)
	log.Printf("[INFO] Instantiating Google Cloud Dataproc Beta client for path %s", dataprocBetaClientBasePath)
	clientDataprocBeta, err := dataprocBeta.NewService(c.context, option.WithHTTPClient(c.client))
	if err != nil {
		log.Printf("[WARN] Error creating client dataproc beta: %s", err)
		return nil
	}
	clientDataprocBeta.UserAgent = userAgent
	clientDataprocBeta.BasePath = dataprocBetaClientBasePath

	return clientDataprocBeta
}

func (c *Config) NewFilestoreClient(userAgent string) *file.Service {
	filestoreClientBasePath := removeBasePathVersion(c.FilestoreBasePath)
	log.Printf("[INFO] Instantiating Filestore client for path %s", filestoreClientBasePath)
	clientFilestore, err := file.NewService(c.context, option.WithHTTPClient(c.client))
	if err != nil {
		log.Printf("[WARN] Error creating client filestore: %s", err)
		return nil
	}
	clientFilestore.UserAgent = userAgent
	clientFilestore.BasePath = filestoreClientBasePath

	return clientFilestore
}

func (c *Config) NewCloudIoTClient(userAgent string) *cloudiot.Service {
	cloudIoTClientBasePath := removeBasePathVersion(c.CloudIoTBasePath)
	log.Printf("[INFO] Instantiating Google Cloud IoT Core client for path %s", cloudIoTClientBasePath)
	clientCloudIoT, err := cloudiot.NewService(c.context, option.WithHTTPClient(c.client))
	if err != nil {
		log.Printf("[WARN] Error creating client cloud iot: %s", err)
		return nil
	}
	clientCloudIoT.UserAgent = userAgent
	clientCloudIoT.BasePath = cloudIoTClientBasePath

	return clientCloudIoT
}

func (c *Config) NewAppEngineClient(userAgent string) *appengine.APIService {
	appEngineClientBasePath := removeBasePathVersion(c.AppEngineBasePath)
	log.Printf("[INFO] Instantiating App Engine client for path %s", appEngineClientBasePath)
	clientAppEngine, err := appengine.NewService(c.context, option.WithHTTPClient(c.client))
	if err != nil {
		log.Printf("[WARN] Error creating client appengine: %s", err)
		return nil
	}
	clientAppEngine.UserAgent = userAgent
	clientAppEngine.BasePath = appEngineClientBasePath

	return clientAppEngine
}

func (c *Config) NewComposerClient(userAgent string) *composer.Service {
	composerClientBasePath := removeBasePathVersion(c.ComposerBasePath)
	log.Printf("[INFO] Instantiating Cloud Composer client for path %s", composerClientBasePath)
	clientComposer, err := composer.NewService(c.context, option.WithHTTPClient(c.client))
	if err != nil {
		log.Printf("[WARN] Error creating client composer: %s", err)
		return nil
	}
	clientComposer.UserAgent = userAgent
	clientComposer.BasePath = composerClientBasePath

	return clientComposer
}

func (c *Config) NewServiceNetworkingClient(userAgent string) *servicenetworking.APIService {
	serviceNetworkingClientBasePath := removeBasePathVersion(c.ServiceNetworkingBasePath)
	log.Printf("[INFO] Instantiating Service Networking client for path %s", serviceNetworkingClientBasePath)
	clientServiceNetworking, err := servicenetworking.NewService(c.context, option.WithHTTPClient(c.client))
	if err != nil {
		log.Printf("[WARN] Error creating client service networking: %s", err)
		return nil
	}
	clientServiceNetworking.UserAgent = userAgent
	clientServiceNetworking.BasePath = serviceNetworkingClientBasePath

	return clientServiceNetworking
}

func (c *Config) NewStorageTransferClient(userAgent string) *storagetransfer.Service {
	storageTransferClientBasePath := removeBasePathVersion(c.StorageTransferBasePath)
	log.Printf("[INFO] Instantiating Google Cloud Storage Transfer client for path %s", storageTransferClientBasePath)
	clientStorageTransfer, err := storagetransfer.NewService(c.context, option.WithHTTPClient(c.client))
	if err != nil {
		log.Printf("[WARN] Error creating client storage transfer: %s", err)
		return nil
	}
	clientStorageTransfer.UserAgent = userAgent
	clientStorageTransfer.BasePath = storageTransferClientBasePath

	return clientStorageTransfer
}

func (c *Config) NewHealthcareClient(userAgent string) *healthcare.Service {
	healthcareClientBasePath := removeBasePathVersion(c.HealthcareBasePath)
	log.Printf("[INFO] Instantiating Google Cloud Healthcare client for path %s", healthcareClientBasePath)
	clientHealthcare, err := healthcare.NewService(c.context, option.WithHTTPClient(c.client))
	if err != nil {
		log.Printf("[WARN] Error creating client healthcare: %s", err)
		return nil
	}
	clientHealthcare.UserAgent = userAgent
	clientHealthcare.BasePath = healthcareClientBasePath

	return clientHealthcare
}

func (c *Config) NewCloudIdentityClient(userAgent string) *cloudidentity.Service {
	cloudidentityClientBasePath := removeBasePathVersion(c.CloudIdentityBasePath)
	log.Printf("[INFO] Instantiating Google Cloud CloudIdentity client for path %s", cloudidentityClientBasePath)
	clientCloudIdentity, err := cloudidentity.NewService(c.context, option.WithHTTPClient(c.client))
	if err != nil {
		log.Printf("[WARN] Error creating client cloud identity: %s", err)
		return nil
	}
	clientCloudIdentity.UserAgent = userAgent
	clientCloudIdentity.BasePath = cloudidentityClientBasePath

	return clientCloudIdentity
}

func (c *Config) BigTableClientFactory(userAgent string) *BigtableClientFactory {
	bigtableClientFactory := &BigtableClientFactory{
		UserAgent:   userAgent,
		TokenSource: c.tokenSource,
	}

	return bigtableClientFactory
}

// Unlike other clients, the Bigtable Admin client doesn't use a single
// service. Instead, there are several distinct services created off
// the base service object. To imitate most other handwritten clients,
// we expose those directly instead of providing the `Service` object
// as a factory.
func (c *Config) NewBigTableProjectsInstancesClient(userAgent string) *bigtableadmin.ProjectsInstancesService {
	bigtableAdminBasePath := removeBasePathVersion(c.BigtableAdminBasePath)
	log.Printf("[INFO] Instantiating Google Cloud BigtableAdmin for path %s", bigtableAdminBasePath)
	clientBigtable, err := bigtableadmin.NewService(c.context, option.WithHTTPClient(c.client))
	if err != nil {
		log.Printf("[WARN] Error creating client big table projects instances: %s", err)
		return nil
	}
	clientBigtable.UserAgent = userAgent
	clientBigtable.BasePath = bigtableAdminBasePath
	clientBigtableProjectsInstances := bigtableadmin.NewProjectsInstancesService(clientBigtable)

	return clientBigtableProjectsInstances
}

func (c *Config) NewBigTableProjectsInstancesTablesClient(userAgent string) *bigtableadmin.ProjectsInstancesTablesService {
	bigtableAdminBasePath := removeBasePathVersion(c.BigtableAdminBasePath)
	log.Printf("[INFO] Instantiating Google Cloud BigtableAdmin for path %s", bigtableAdminBasePath)
	clientBigtable, err := bigtableadmin.NewService(c.context, option.WithHTTPClient(c.client))
	if err != nil {
		log.Printf("[WARN] Error creating client projects instances tables: %s", err)
		return nil
	}
	clientBigtable.UserAgent = userAgent
	clientBigtable.BasePath = bigtableAdminBasePath
	clientBigtableProjectsInstancesTables := bigtableadmin.NewProjectsInstancesTablesService(clientBigtable)

	return clientBigtableProjectsInstancesTables
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
		token := &oauth2.Token{AccessToken: contents}

		if c.ImpersonateServiceAccount != "" {
			opts := []option.ClientOption{option.WithTokenSource(oauth2.StaticTokenSource(token)), option.ImpersonateCredentials(c.ImpersonateServiceAccount, c.ImpersonateServiceAccountDelegates...), option.WithScopes(clientScopes...)}
			creds, err := transport.Creds(context.TODO(), opts...)
			if err != nil {
				return googleoauth.Credentials{}, err
			}
			return *creds, nil
		}

		log.Printf("[INFO] Authenticating using configured Google JSON 'access_token'...")
		log.Printf("[INFO]   -- Scopes: %s", clientScopes)

		return googleoauth.Credentials{
			TokenSource: staticTokenSource{oauth2.StaticTokenSource(token)},
		}, nil
	}

	if c.Credentials != "" {
		contents, _, err := pathOrContents(c.Credentials)
		if err != nil {
			return googleoauth.Credentials{}, fmt.Errorf("error loading credentials: %s", err)
		}
		if c.ImpersonateServiceAccount != "" {
			opts := []option.ClientOption{option.WithCredentialsJSON([]byte(contents)), option.ImpersonateCredentials(c.ImpersonateServiceAccount, c.ImpersonateServiceAccountDelegates...), option.WithScopes(clientScopes...)}
			creds, err := transport.Creds(context.TODO(), opts...)
			if err != nil {
				return googleoauth.Credentials{}, err
			}
			return *creds, nil
		}
		creds, err := googleoauth.CredentialsFromJSON(c.context, []byte(contents), clientScopes...)
		if err != nil {
			return googleoauth.Credentials{}, fmt.Errorf("unable to parse credentials from '%s': %s", contents, err)
		}

		log.Printf("[INFO] Authenticating using configured Google JSON 'credentials'...")
		log.Printf("[INFO]   -- Scopes: %s", clientScopes)
		return *creds, nil
	}

	if c.ImpersonateServiceAccount != "" {
		opts := option.ImpersonateCredentials(c.ImpersonateServiceAccount, c.ImpersonateServiceAccountDelegates...)
		creds, err := transport.Creds(context.TODO(), opts, option.WithScopes(clientScopes...))
		if err != nil {
			return googleoauth.Credentials{}, err
		}
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
	c.CloudIdentityBasePath = CloudIdentityDefaultBasePath
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
