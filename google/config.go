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
	resourceManagerV2 "google.golang.org/api/cloudresourcemanager/v2"
	composer "google.golang.org/api/composer/v1beta1"
	computeBeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/container/v1"
	containerBeta "google.golang.org/api/container/v1beta1"
	dataflow "google.golang.org/api/dataflow/v1b3"
	"google.golang.org/api/dataproc/v1"
	"google.golang.org/api/dns/v1"
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
	ApigeeBasePath               string
	AppEngineBasePath            string
	BigQueryBasePath             string
	BigqueryDataTransferBasePath string
	BigqueryReservationBasePath  string
	BigtableBasePath             string
	BillingBasePath              string
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
	DialogflowCXBasePath         string
	DNSBasePath                  string
	FilestoreBasePath            string
	FirestoreBasePath            string
	GameServicesBasePath         string
	HealthcareBasePath           string
	IapBasePath                  string
	IdentityPlatformBasePath     string
	KMSBasePath                  string
	LoggingBasePath              string
	MemcacheBasePath             string
	MLEngineBasePath             string
	MonitoringBasePath           string
	NetworkManagementBasePath    string
	NotebooksBasePath            string
	OSConfigBasePath             string
	OSLoginBasePath              string
	PubsubBasePath               string
	PubsubLiteBasePath           string
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
	TagsBasePath                 string
	TPUBasePath                  string
	VertexAIBasePath             string
	VPCAccessBasePath            string
	WorkflowsBasePath            string

	CloudBillingBasePath      string
	ComposerBasePath          string
	ComputeBetaBasePath       string
	ContainerBasePath         string
	ContainerBetaBasePath     string
	DataprocBetaBasePath      string
	DataflowBasePath          string
	IamCredentialsBasePath    string
	ResourceManagerV2BasePath string
	IAMBasePath               string
	CloudIoTBasePath          string
	ServiceNetworkingBasePath string
	StorageTransferBasePath   string
	BigtableAdminBasePath     string

	requestBatcherServiceUsage *RequestBatcher
	requestBatcherIam          *RequestBatcher

	// start DCLBasePaths
	// dataprocBasePath is implemented in mm
	EventarcBasePath string
	GkeHubBasePath   string
}

const AccessApprovalBasePathKey = "AccessApproval"
const AccessContextManagerBasePathKey = "AccessContextManager"
const ActiveDirectoryBasePathKey = "ActiveDirectory"
const ApigeeBasePathKey = "Apigee"
const AppEngineBasePathKey = "AppEngine"
const BigQueryBasePathKey = "BigQuery"
const BigqueryDataTransferBasePathKey = "BigqueryDataTransfer"
const BigqueryReservationBasePathKey = "BigqueryReservation"
const BigtableBasePathKey = "Bigtable"
const BillingBasePathKey = "Billing"
const BinaryAuthorizationBasePathKey = "BinaryAuthorization"
const CloudAssetBasePathKey = "CloudAsset"
const CloudBuildBasePathKey = "CloudBuild"
const CloudFunctionsBasePathKey = "CloudFunctions"
const CloudIdentityBasePathKey = "CloudIdentity"
const CloudIotBasePathKey = "CloudIot"
const CloudRunBasePathKey = "CloudRun"
const CloudSchedulerBasePathKey = "CloudScheduler"
const CloudTasksBasePathKey = "CloudTasks"
const ComputeBasePathKey = "Compute"
const ContainerAnalysisBasePathKey = "ContainerAnalysis"
const DataCatalogBasePathKey = "DataCatalog"
const DataLossPreventionBasePathKey = "DataLossPrevention"
const DataprocBasePathKey = "Dataproc"
const DatastoreBasePathKey = "Datastore"
const DeploymentManagerBasePathKey = "DeploymentManager"
const DialogflowBasePathKey = "Dialogflow"
const DialogflowCXBasePathKey = "DialogflowCX"
const DNSBasePathKey = "DNS"
const FilestoreBasePathKey = "Filestore"
const FirestoreBasePathKey = "Firestore"
const GameServicesBasePathKey = "GameServices"
const HealthcareBasePathKey = "Healthcare"
const IapBasePathKey = "Iap"
const IdentityPlatformBasePathKey = "IdentityPlatform"
const KMSBasePathKey = "KMS"
const LoggingBasePathKey = "Logging"
const MemcacheBasePathKey = "Memcache"
const MLEngineBasePathKey = "MLEngine"
const MonitoringBasePathKey = "Monitoring"
const NetworkManagementBasePathKey = "NetworkManagement"
const NotebooksBasePathKey = "Notebooks"
const OSConfigBasePathKey = "OSConfig"
const OSLoginBasePathKey = "OSLogin"
const PubsubBasePathKey = "Pubsub"
const PubsubLiteBasePathKey = "PubsubLite"
const RedisBasePathKey = "Redis"
const ResourceManagerBasePathKey = "ResourceManager"
const RuntimeConfigBasePathKey = "RuntimeConfig"
const SecretManagerBasePathKey = "SecretManager"
const SecurityCenterBasePathKey = "SecurityCenter"
const ServiceManagementBasePathKey = "ServiceManagement"
const ServiceUsageBasePathKey = "ServiceUsage"
const SourceRepoBasePathKey = "SourceRepo"
const SpannerBasePathKey = "Spanner"
const SQLBasePathKey = "SQL"
const StorageBasePathKey = "Storage"
const TagsBasePathKey = "Tags"
const TPUBasePathKey = "TPU"
const VertexAIBasePathKey = "VertexAI"
const VPCAccessBasePathKey = "VPCAccess"
const WorkflowsBasePathKey = "Workflows"
const CloudBillingBasePathKey = "CloudBilling"
const ComposerBasePathKey = "Composer"
const ComputeBetaBasePathKey = "ComputeBeta"
const ContainerBasePathKey = "Container"
const DataprocBetaBasePathKey = "DataprocBeta"
const ContainerBetaBasePathKey = "ContainerBeta"
const DataflowBasePathKey = "Dataflow"
const IAMBasePathKey = "IAM"
const IamCredentialsBasePathKey = "IamCredentials"
const ResourceManagerV2BasePathKey = "ResourceManagerV2"
const ServiceNetworkingBasePathKey = "ServiceNetworking"
const StorageTransferBasePathKey = "StorageTransfer"
const BigtableAdminBasePathKey = "BigtableAdmin"
const GkeHubFeatureBasePathKey = "GkeHubFeatureBasePathKey"

// Generated product base paths
var DefaultBasePaths = map[string]string{
	AccessApprovalBasePathKey:       "https://accessapproval.googleapis.com/v1/",
	AccessContextManagerBasePathKey: "https://accesscontextmanager.googleapis.com/v1/",
	ActiveDirectoryBasePathKey:      "https://managedidentities.googleapis.com/v1/",
	ApigeeBasePathKey:               "https://apigee.googleapis.com/v1/",
	AppEngineBasePathKey:            "https://appengine.googleapis.com/v1/",
	BigQueryBasePathKey:             "https://bigquery.googleapis.com/bigquery/v2/",
	BigqueryDataTransferBasePathKey: "https://bigquerydatatransfer.googleapis.com/v1/",
	BigqueryReservationBasePathKey:  "https://bigqueryreservation.googleapis.com/v1/",
	BigtableBasePathKey:             "https://bigtableadmin.googleapis.com/v2/",
	BillingBasePathKey:              "https://billingbudgets.googleapis.com/v1/",
	BinaryAuthorizationBasePathKey:  "https://binaryauthorization.googleapis.com/v1/",
	CloudAssetBasePathKey:           "https://cloudasset.googleapis.com/v1/",
	CloudBuildBasePathKey:           "https://cloudbuild.googleapis.com/v1/",
	CloudFunctionsBasePathKey:       "https://cloudfunctions.googleapis.com/v1/",
	CloudIdentityBasePathKey:        "https://cloudidentity.googleapis.com/v1/",
	CloudIotBasePathKey:             "https://cloudiot.googleapis.com/v1/",
	CloudRunBasePathKey:             "https://{{location}}-run.googleapis.com/",
	CloudSchedulerBasePathKey:       "https://cloudscheduler.googleapis.com/v1/",
	CloudTasksBasePathKey:           "https://cloudtasks.googleapis.com/v2/",
	ComputeBasePathKey:              "https://compute.googleapis.com/compute/v1/",
	ContainerAnalysisBasePathKey:    "https://containeranalysis.googleapis.com/v1/",
	DataCatalogBasePathKey:          "https://datacatalog.googleapis.com/v1/",
	DataLossPreventionBasePathKey:   "https://dlp.googleapis.com/v2/",
	DataprocBasePathKey:             "https://dataproc.googleapis.com/v1/",
	DatastoreBasePathKey:            "https://datastore.googleapis.com/v1/",
	DeploymentManagerBasePathKey:    "https://www.googleapis.com/deploymentmanager/v2/",
	DialogflowBasePathKey:           "https://dialogflow.googleapis.com/v2/",
	DialogflowCXBasePathKey:         "https://{{location}}-dialogflow.googleapis.com/v3/",
	DNSBasePathKey:                  "https://dns.googleapis.com/dns/v1/",
	FilestoreBasePathKey:            "https://file.googleapis.com/v1/",
	FirestoreBasePathKey:            "https://firestore.googleapis.com/v1/",
	GameServicesBasePathKey:         "https://gameservices.googleapis.com/v1/",
	HealthcareBasePathKey:           "https://healthcare.googleapis.com/v1/",
	IapBasePathKey:                  "https://iap.googleapis.com/v1/",
	IdentityPlatformBasePathKey:     "https://identitytoolkit.googleapis.com/v2/",
	KMSBasePathKey:                  "https://cloudkms.googleapis.com/v1/",
	LoggingBasePathKey:              "https://logging.googleapis.com/v2/",
	MemcacheBasePathKey:             "https://memcache.googleapis.com/v1/",
	MLEngineBasePathKey:             "https://ml.googleapis.com/v1/",
	MonitoringBasePathKey:           "https://monitoring.googleapis.com/",
	NetworkManagementBasePathKey:    "https://networkmanagement.googleapis.com/v1/",
	NotebooksBasePathKey:            "https://notebooks.googleapis.com/v1/",
	OSConfigBasePathKey:             "https://osconfig.googleapis.com/v1/",
	OSLoginBasePathKey:              "https://oslogin.googleapis.com/v1/",
	PubsubBasePathKey:               "https://pubsub.googleapis.com/v1/",
	PubsubLiteBasePathKey:           "https://{{region}}-pubsublite.googleapis.com/v1/admin/",
	RedisBasePathKey:                "https://redis.googleapis.com/v1/",
	ResourceManagerBasePathKey:      "https://cloudresourcemanager.googleapis.com/v1/",
	RuntimeConfigBasePathKey:        "https://runtimeconfig.googleapis.com/v1beta1/",
	SecretManagerBasePathKey:        "https://secretmanager.googleapis.com/v1/",
	SecurityCenterBasePathKey:       "https://securitycenter.googleapis.com/v1/",
	ServiceManagementBasePathKey:    "https://servicemanagement.googleapis.com/v1/",
	ServiceUsageBasePathKey:         "https://serviceusage.googleapis.com/v1/",
	SourceRepoBasePathKey:           "https://sourcerepo.googleapis.com/v1/",
	SpannerBasePathKey:              "https://spanner.googleapis.com/v1/",
	SQLBasePathKey:                  "https://sqladmin.googleapis.com/sql/v1beta4/",
	StorageBasePathKey:              "https://storage.googleapis.com/storage/v1/",
	TagsBasePathKey:                 "https://cloudresourcemanager.googleapis.com/v3/",
	TPUBasePathKey:                  "https://tpu.googleapis.com/v1/",
	VertexAIBasePathKey:             "https://{{region}}-aiplatform.googleapis.com/v1/",
	VPCAccessBasePathKey:            "https://vpcaccess.googleapis.com/v1/",
	WorkflowsBasePathKey:            "https://workflows.googleapis.com/v1/",
	CloudBillingBasePathKey:         "https://cloudbilling.googleapis.com/v1/",
	ComposerBasePathKey:             "https://composer.googleapis.com/v1/",
	ComputeBetaBasePathKey:          "https://www.googleapis.com/compute/beta/",
	ContainerBasePathKey:            "https://container.googleapis.com/v1/",
	ContainerBetaBasePathKey:        "https://container.googleapis.com/v1beta1/",
	DataprocBetaBasePathKey:         "https://dataproc.googleapis.com/v1beta2/",
	DataflowBasePathKey:             "https://dataflow.googleapis.com/v1b3/",
	IAMBasePathKey:                  "https://iam.googleapis.com/v1/",
	IamCredentialsBasePathKey:       "https://iamcredentials.googleapis.com/v1/",
	ResourceManagerV2BasePathKey:    "https://cloudresourcemanager.googleapis.com/v2/",
	ServiceNetworkingBasePathKey:    "https://servicenetworking.googleapis.com/v1/",
	StorageTransferBasePathKey:      "https://storagetransfer.googleapis.com/v1/",
	BigtableAdminBasePathKey:        "https://bigtableadmin.googleapis.com/v2/",
	GkeHubFeatureBasePathKey:        "https://gkehub.googleapis.com/v1beta/",
}

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

	tokenSource, err := c.getTokenSource(c.Scopes, false)
	if err != nil {
		return err
	}

	c.tokenSource = tokenSource

	cleanCtx := context.WithValue(ctx, oauth2.HTTPClient, cleanhttp.DefaultClient())

	// 1. MTLS TRANSPORT/CLIENT - sets up proper auth headers
	client, _, err := transport.NewHTTPClient(cleanCtx, option.WithTokenSource(tokenSource))
	if err != nil {
		return err
	}
	// Userinfo is fetched before request logging is enabled to reduce additional noise.
	err = c.logGoogleIdentities()
	if err != nil {
		return err
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

	// Set final transport value.
	client.Transport = headerTransport

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
		return 120 * time.Second
	}
	return c.RequestTimeout
}

// Print Identities executing terraform API Calls.
func (c *Config) logGoogleIdentities() error {
	if c.ImpersonateServiceAccount == "" {

		tokenSource, err := c.getTokenSource(c.Scopes, true)
		if err != nil {
			return err
		}
		c.client = oauth2.NewClient(c.context, tokenSource) // c.client isn't initialised fully when this code is called.

		email, err := GetCurrentUserEmail(c, c.userAgent)
		if err != nil {
			log.Printf("[INFO] error retrieving userinfo for your provider credentials. have you enabled the 'https://www.googleapis.com/auth/userinfo.email' scope? error: %s", err)
		}

		log.Printf("[INFO] Terraform is using this identity: %s", email)

		return nil

	}

	// Drop Impersonated ClientOption from OAuth2 TokenSource to infer original identity

	tokenSource, err := c.getTokenSource(c.Scopes, true)
	if err != nil {
		return err
	}
	c.client = oauth2.NewClient(c.context, tokenSource) // c.client isn't initialised fully when this code is called.

	email, err := GetCurrentUserEmail(c, c.userAgent)
	if err != nil {
		log.Printf("[INFO] error retrieving userinfo for your provider credentials. have you enabled the 'https://www.googleapis.com/auth/userinfo.email' scope? error: %s", err)
	}

	log.Printf("[INFO] Terraform is configured with service account impersonation, original identity: %s, impersonated identity: %s", email, c.ImpersonateServiceAccount)

	// Add the Impersonated ClientOption back in to the OAuth2 TokenSource

	tokenSource, err = c.getTokenSource(c.Scopes, false)
	if err != nil {
		return err
	}
	c.client = oauth2.NewClient(c.context, tokenSource) // c.client isn't initialised fully when this code is called.

	return nil
}

// Get a TokenSource based on the Google Credentials configured.
// If initialCredentialsOnly is true, don't follow the impersonation settings and return the initial set of creds.
func (c *Config) getTokenSource(clientScopes []string, initialCredentialsOnly bool) (oauth2.TokenSource, error) {
	creds, err := c.GetCredentials(clientScopes, initialCredentialsOnly)
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

func (c *Config) NewKmsClientWithCtx(ctx context.Context, userAgent string) *cloudkms.Service {
	kmsClientBasePath := removeBasePathVersion(c.KMSBasePath)
	log.Printf("[INFO] Instantiating Google Cloud KMS client for path %s", kmsClientBasePath)
	clientKms, err := cloudkms.NewService(ctx, option.WithHTTPClient(c.client))
	if err != nil {
		log.Printf("[WARN] Error creating client kms: %s", err)
		return nil
	}
	clientKms.UserAgent = userAgent
	clientKms.BasePath = kmsClientBasePath

	return clientKms
}

func (c *Config) NewKmsClient(userAgent string) *cloudkms.Service {
	return c.NewKmsClientWithCtx(c.context, userAgent)
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

func (c *Config) NewResourceManagerV2Client(userAgent string) *resourceManagerV2.Service {
	resourceManagerV2BasePath := removeBasePathVersion(c.ResourceManagerV2BasePath)
	log.Printf("[INFO] Instantiating Google Cloud ResourceManager V client for path %s", resourceManagerV2BasePath)
	clientResourceManagerV2, err := resourceManagerV2.NewService(c.context, option.WithHTTPClient(c.client))
	if err != nil {
		log.Printf("[WARN] Error creating client resource manager v2: %s", err)
		return nil
	}
	clientResourceManagerV2.UserAgent = userAgent
	clientResourceManagerV2.BasePath = resourceManagerV2BasePath

	return clientResourceManagerV2
}

func (c *Config) NewRuntimeconfigClient(userAgent string) *runtimeconfig.Service {
	runtimeConfigClientBasePath := removeBasePathVersion(c.RuntimeConfigBasePath)
	log.Printf("[INFO] Instantiating Google Cloud Runtimeconfig client for path %s", runtimeConfigClientBasePath)
	clientRuntimeconfig, err := runtimeconfig.NewService(c.context, option.WithHTTPClient(c.client))
	if err != nil {
		log.Printf("[WARN] Error creating client runtime config: %s", err)
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

// Get a set of credentials with a given scope (clientScopes) based on the Config object.
// If initialCredentialsOnly is true, don't follow the impersonation settings and return the initial set of creds
// instead.
func (c *Config) GetCredentials(clientScopes []string, initialCredentialsOnly bool) (googleoauth.Credentials, error) {

	if c.AccessToken != "" {
		contents, _, err := pathOrContents(c.AccessToken)
		if err != nil {
			return googleoauth.Credentials{}, fmt.Errorf("Error loading access token: %s", err)
		}
		token := &oauth2.Token{AccessToken: contents}

		if c.ImpersonateServiceAccount != "" && !initialCredentialsOnly {
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
		if c.ImpersonateServiceAccount != "" && !initialCredentialsOnly {
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

	if c.ImpersonateServiceAccount != "" && !initialCredentialsOnly {
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
	c.AccessApprovalBasePath = DefaultBasePaths[AccessApprovalBasePathKey]
	c.AccessContextManagerBasePath = DefaultBasePaths[AccessContextManagerBasePathKey]
	c.ActiveDirectoryBasePath = DefaultBasePaths[ActiveDirectoryBasePathKey]
	c.ApigeeBasePath = DefaultBasePaths[ApigeeBasePathKey]
	c.AppEngineBasePath = DefaultBasePaths[AppEngineBasePathKey]
	c.BigQueryBasePath = DefaultBasePaths[BigQueryBasePathKey]
	c.BigqueryDataTransferBasePath = DefaultBasePaths[BigqueryDataTransferBasePathKey]
	c.BigqueryReservationBasePath = DefaultBasePaths[BigqueryReservationBasePathKey]
	c.BigtableBasePath = DefaultBasePaths[BigtableBasePathKey]
	c.BillingBasePath = DefaultBasePaths[BillingBasePathKey]
	c.BinaryAuthorizationBasePath = DefaultBasePaths[BinaryAuthorizationBasePathKey]
	c.CloudAssetBasePath = DefaultBasePaths[CloudAssetBasePathKey]
	c.CloudBuildBasePath = DefaultBasePaths[CloudBuildBasePathKey]
	c.CloudFunctionsBasePath = DefaultBasePaths[CloudFunctionsBasePathKey]
	c.CloudIdentityBasePath = DefaultBasePaths[CloudIdentityBasePathKey]
	c.CloudIotBasePath = DefaultBasePaths[CloudIotBasePathKey]
	c.CloudRunBasePath = DefaultBasePaths[CloudRunBasePathKey]
	c.CloudSchedulerBasePath = DefaultBasePaths[CloudSchedulerBasePathKey]
	c.CloudTasksBasePath = DefaultBasePaths[CloudTasksBasePathKey]
	c.ComputeBasePath = DefaultBasePaths[ComputeBasePathKey]
	c.ContainerAnalysisBasePath = DefaultBasePaths[ContainerAnalysisBasePathKey]
	c.DataCatalogBasePath = DefaultBasePaths[DataCatalogBasePathKey]
	c.DataLossPreventionBasePath = DefaultBasePaths[DataLossPreventionBasePathKey]
	c.DataprocBasePath = DefaultBasePaths[DataprocBasePathKey]
	c.DatastoreBasePath = DefaultBasePaths[DatastoreBasePathKey]
	c.DeploymentManagerBasePath = DefaultBasePaths[DeploymentManagerBasePathKey]
	c.DialogflowBasePath = DefaultBasePaths[DialogflowBasePathKey]
	c.DialogflowCXBasePath = DefaultBasePaths[DialogflowCXBasePathKey]
	c.DNSBasePath = DefaultBasePaths[DNSBasePathKey]
	c.FilestoreBasePath = DefaultBasePaths[FilestoreBasePathKey]
	c.FirestoreBasePath = DefaultBasePaths[FirestoreBasePathKey]
	c.GameServicesBasePath = DefaultBasePaths[GameServicesBasePathKey]
	c.HealthcareBasePath = DefaultBasePaths[HealthcareBasePathKey]
	c.IapBasePath = DefaultBasePaths[IapBasePathKey]
	c.IdentityPlatformBasePath = DefaultBasePaths[IdentityPlatformBasePathKey]
	c.KMSBasePath = DefaultBasePaths[KMSBasePathKey]
	c.LoggingBasePath = DefaultBasePaths[LoggingBasePathKey]
	c.MemcacheBasePath = DefaultBasePaths[MemcacheBasePathKey]
	c.MLEngineBasePath = DefaultBasePaths[MLEngineBasePathKey]
	c.MonitoringBasePath = DefaultBasePaths[MonitoringBasePathKey]
	c.NetworkManagementBasePath = DefaultBasePaths[NetworkManagementBasePathKey]
	c.NotebooksBasePath = DefaultBasePaths[NotebooksBasePathKey]
	c.OSConfigBasePath = DefaultBasePaths[OSConfigBasePathKey]
	c.OSLoginBasePath = DefaultBasePaths[OSLoginBasePathKey]
	c.PubsubBasePath = DefaultBasePaths[PubsubBasePathKey]
	c.PubsubLiteBasePath = DefaultBasePaths[PubsubLiteBasePathKey]
	c.RedisBasePath = DefaultBasePaths[RedisBasePathKey]
	c.ResourceManagerBasePath = DefaultBasePaths[ResourceManagerBasePathKey]
	c.RuntimeConfigBasePath = DefaultBasePaths[RuntimeConfigBasePathKey]
	c.SecretManagerBasePath = DefaultBasePaths[SecretManagerBasePathKey]
	c.SecurityCenterBasePath = DefaultBasePaths[SecurityCenterBasePathKey]
	c.ServiceManagementBasePath = DefaultBasePaths[ServiceManagementBasePathKey]
	c.ServiceUsageBasePath = DefaultBasePaths[ServiceUsageBasePathKey]
	c.SourceRepoBasePath = DefaultBasePaths[SourceRepoBasePathKey]
	c.SpannerBasePath = DefaultBasePaths[SpannerBasePathKey]
	c.SQLBasePath = DefaultBasePaths[SQLBasePathKey]
	c.StorageBasePath = DefaultBasePaths[StorageBasePathKey]
	c.TagsBasePath = DefaultBasePaths[TagsBasePathKey]
	c.TPUBasePath = DefaultBasePaths[TPUBasePathKey]
	c.VertexAIBasePath = DefaultBasePaths[VertexAIBasePathKey]
	c.VPCAccessBasePath = DefaultBasePaths[VPCAccessBasePathKey]
	c.WorkflowsBasePath = DefaultBasePaths[WorkflowsBasePathKey]

	// Handwritten Products / Versioned / Atypical Entries
	c.CloudBillingBasePath = DefaultBasePaths[CloudBillingBasePathKey]
	c.ComposerBasePath = DefaultBasePaths[ComposerBasePathKey]
	c.ComputeBetaBasePath = DefaultBasePaths[ComputeBetaBasePathKey]
	c.ContainerBasePath = DefaultBasePaths[ContainerBasePathKey]
	c.ContainerBetaBasePath = DefaultBasePaths[ContainerBetaBasePathKey]
	c.DataprocBasePath = DefaultBasePaths[DataprocBasePathKey]
	c.DataflowBasePath = DefaultBasePaths[DataflowBasePathKey]
	c.IamCredentialsBasePath = DefaultBasePaths[IamCredentialsBasePathKey]
	c.ResourceManagerV2BasePath = DefaultBasePaths[ResourceManagerV2BasePathKey]
	c.IAMBasePath = DefaultBasePaths[IAMBasePathKey]
	c.ServiceNetworkingBasePath = DefaultBasePaths[ServiceNetworkingBasePathKey]
	c.BigQueryBasePath = DefaultBasePaths[BigQueryBasePathKey]
	c.StorageTransferBasePath = DefaultBasePaths[StorageTransferBasePathKey]
	c.BigtableAdminBasePath = DefaultBasePaths[BigtableAdminBasePathKey]
}
