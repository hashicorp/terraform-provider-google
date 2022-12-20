package google

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/sirupsen/logrus"
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
	resourceManagerV3 "google.golang.org/api/cloudresourcemanager/v3"
	"google.golang.org/api/composer/v1"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/container/v1"
	dataflow "google.golang.org/api/dataflow/v1b3"
	"google.golang.org/api/dataproc/v1"
	"google.golang.org/api/dns/v1"
	healthcare "google.golang.org/api/healthcare/v1"
	"google.golang.org/api/iam/v1"
	iamcredentials "google.golang.org/api/iamcredentials/v1"
	cloudlogging "google.golang.org/api/logging/v2"
	"google.golang.org/api/pubsub/v1"
	"google.golang.org/api/servicemanagement/v1"
	"google.golang.org/api/servicenetworking/v1"
	"google.golang.org/api/serviceusage/v1"
	"google.golang.org/api/sourcerepo/v1"
	"google.golang.org/api/spanner/v1"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
	"google.golang.org/api/storage/v1"
	"google.golang.org/api/storagetransfer/v1"
	"google.golang.org/api/transport"
	"google.golang.org/grpc"
)

type providerMeta struct {
	ModuleName string `cty:"module_name"`
}

type Formatter struct {
	TimestampFormat string
	LogFormat       string
}

// Borrowed logic from https://github.com/sirupsen/logrus/blob/master/json_formatter.go and https://github.com/t-tomalak/logrus-easy-formatter/blob/master/formatter.go
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Suppress logs if TF_LOG is not DEBUG or TRACE
	if !logging.IsDebugOrHigher() {
		return nil, nil
	}

	// Also suppress based on log content
	// - frequent transport spam
	// - ListenSocket logs from gRPC
	isTransportSpam := strings.Contains(entry.Message, "transport is closing")
	listenSocketRegex := regexp.MustCompile(`\[Server #\d+( ListenSocket #\d+)*\]`) // Match patterns like `[Server #00]` or `[Server #00 ListenSocket #00]`
	isListenSocketLog := listenSocketRegex.MatchString(entry.Message)
	if isTransportSpam || isListenSocketLog {
		return nil, nil
	}

	output := f.LogFormat
	entry.Level = logrus.DebugLevel // Force Entries to be Debug

	timestampFormat := f.TimestampFormat

	output = strings.Replace(output, "%time%", entry.Time.Format(timestampFormat), 1)

	output = strings.Replace(output, "%msg%", entry.Message, 1)

	level := strings.ToUpper(entry.Level.String())
	output = strings.Replace(output, "%lvl%", level, 1)

	var gRPCMessageFlag bool
	for k, val := range entry.Data {
		switch v := val.(type) {
		case string:
			output = strings.Replace(output, "%"+k+"%", v, 1)
		case int:
			s := strconv.Itoa(v)
			output = strings.Replace(output, "%"+k+"%", s, 1)
		case bool:
			s := strconv.FormatBool(v)
			output = strings.Replace(output, "%"+k+"%", s, 1)
		}

		if k != "system" {
			gRPCMessageFlag = true
		}
	}

	if gRPCMessageFlag {
		data := make(logrus.Fields, len(entry.Data)+4)
		for k, v := range entry.Data {
			switch v := v.(type) {
			case error:
				// Otherwise errors are ignored by `encoding/json`
				// https://github.com/sirupsen/logrus/issues/137
				data[k] = v.Error()
			default:
				data[k] = v
			}
		}

		var b *bytes.Buffer
		if entry.Buffer != nil {
			b = entry.Buffer
		} else {
			b = &bytes.Buffer{}
		}

		encoder := json.NewEncoder(b)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(data); err != nil {
			return nil, fmt.Errorf("failed to marshal fields to JSON, %w", err)
		}

		finalOutput := append([]byte(output), b.Bytes()...)
		return finalOutput, nil
	}

	return []byte(output), nil
}

// Config is the configuration structure used to instantiate the Google
// provider.
type Config struct {
	DCLConfig
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
	RequestReason                      string
	RequestTimeout                     time.Duration
	// PollInterval is passed to resource.StateChangeConf in common_operation.go
	// It controls the interval at which we poll for successful operations
	PollInterval time.Duration

	client             *http.Client
	context            context.Context
	userAgent          string
	gRPCLoggingOptions []option.ClientOption

	tokenSource oauth2.TokenSource

	AccessApprovalBasePath       string
	AccessContextManagerBasePath string
	ActiveDirectoryBasePath      string
	ApigeeBasePath               string
	AppEngineBasePath            string
	ArtifactRegistryBasePath     string
	BeyondcorpBasePath           string
	BigQueryBasePath             string
	BigqueryAnalyticsHubBasePath string
	BigqueryConnectionBasePath   string
	BigqueryDataTransferBasePath string
	BigqueryReservationBasePath  string
	BigtableBasePath             string
	BillingBasePath              string
	BinaryAuthorizationBasePath  string
	CertificateManagerBasePath   string
	CloudAssetBasePath           string
	CloudBuildBasePath           string
	CloudFunctionsBasePath       string
	Cloudfunctions2BasePath      string
	CloudIdentityBasePath        string
	CloudIdsBasePath             string
	CloudIotBasePath             string
	CloudRunBasePath             string
	CloudRunV2BasePath           string
	CloudSchedulerBasePath       string
	CloudTasksBasePath           string
	ComputeBasePath              string
	ContainerAnalysisBasePath    string
	DataCatalogBasePath          string
	DataFusionBasePath           string
	DataLossPreventionBasePath   string
	DataprocBasePath             string
	DataprocMetastoreBasePath    string
	DatastoreBasePath            string
	DatastreamBasePath           string
	DeploymentManagerBasePath    string
	DialogflowBasePath           string
	DialogflowCXBasePath         string
	DNSBasePath                  string
	DocumentAIBasePath           string
	EssentialContactsBasePath    string
	FilestoreBasePath            string
	FirestoreBasePath            string
	GameServicesBasePath         string
	GKEHubBasePath               string
	HealthcareBasePath           string
	IAMBetaBasePath              string
	IAMWorkforcePoolBasePath     string
	IapBasePath                  string
	IdentityPlatformBasePath     string
	KMSBasePath                  string
	LoggingBasePath              string
	MemcacheBasePath             string
	MLEngineBasePath             string
	MonitoringBasePath           string
	NetworkManagementBasePath    string
	NetworkServicesBasePath      string
	NotebooksBasePath            string
	OSConfigBasePath             string
	OSLoginBasePath              string
	PrivatecaBasePath            string
	PubsubBasePath               string
	PubsubLiteBasePath           string
	RedisBasePath                string
	ResourceManagerBasePath      string
	SecretManagerBasePath        string
	SecurityCenterBasePath       string
	ServiceManagementBasePath    string
	ServiceUsageBasePath         string
	SourceRepoBasePath           string
	SpannerBasePath              string
	SQLBasePath                  string
	StorageBasePath              string
	StorageTransferBasePath      string
	TagsBasePath                 string
	TPUBasePath                  string
	VertexAIBasePath             string
	VPCAccessBasePath            string
	WorkflowsBasePath            string

	CloudBillingBasePath      string
	ComposerBasePath          string
	ContainerBasePath         string
	DataflowBasePath          string
	IamCredentialsBasePath    string
	ResourceManagerV3BasePath string
	IAMBasePath               string
	CloudIoTBasePath          string
	ServiceNetworkingBasePath string
	BigtableAdminBasePath     string

	// dcl
	ContainerAwsBasePath   string
	ContainerAzureBasePath string

	requestBatcherServiceUsage *RequestBatcher
	requestBatcherIam          *RequestBatcher
}

const AccessApprovalBasePathKey = "AccessApproval"
const AccessContextManagerBasePathKey = "AccessContextManager"
const ActiveDirectoryBasePathKey = "ActiveDirectory"
const ApigeeBasePathKey = "Apigee"
const AppEngineBasePathKey = "AppEngine"
const ArtifactRegistryBasePathKey = "ArtifactRegistry"
const BeyondcorpBasePathKey = "Beyondcorp"
const BigQueryBasePathKey = "BigQuery"
const BigqueryAnalyticsHubBasePathKey = "BigqueryAnalyticsHub"
const BigqueryConnectionBasePathKey = "BigqueryConnection"
const BigqueryDataTransferBasePathKey = "BigqueryDataTransfer"
const BigqueryReservationBasePathKey = "BigqueryReservation"
const BigtableBasePathKey = "Bigtable"
const BillingBasePathKey = "Billing"
const BinaryAuthorizationBasePathKey = "BinaryAuthorization"
const CertificateManagerBasePathKey = "CertificateManager"
const CloudAssetBasePathKey = "CloudAsset"
const CloudBuildBasePathKey = "CloudBuild"
const CloudFunctionsBasePathKey = "CloudFunctions"
const Cloudfunctions2BasePathKey = "Cloudfunctions2"
const CloudIdentityBasePathKey = "CloudIdentity"
const CloudIdsBasePathKey = "CloudIds"
const CloudIotBasePathKey = "CloudIot"
const CloudRunBasePathKey = "CloudRun"
const CloudRunV2BasePathKey = "CloudRunV2"
const CloudSchedulerBasePathKey = "CloudScheduler"
const CloudTasksBasePathKey = "CloudTasks"
const ComputeBasePathKey = "Compute"
const ContainerAnalysisBasePathKey = "ContainerAnalysis"
const DataCatalogBasePathKey = "DataCatalog"
const DataFusionBasePathKey = "DataFusion"
const DataLossPreventionBasePathKey = "DataLossPrevention"
const DataprocBasePathKey = "Dataproc"
const DataprocMetastoreBasePathKey = "DataprocMetastore"
const DatastoreBasePathKey = "Datastore"
const DatastreamBasePathKey = "Datastream"
const DeploymentManagerBasePathKey = "DeploymentManager"
const DialogflowBasePathKey = "Dialogflow"
const DialogflowCXBasePathKey = "DialogflowCX"
const DNSBasePathKey = "DNS"
const DocumentAIBasePathKey = "DocumentAI"
const EssentialContactsBasePathKey = "EssentialContacts"
const FilestoreBasePathKey = "Filestore"
const FirestoreBasePathKey = "Firestore"
const GameServicesBasePathKey = "GameServices"
const GKEHubBasePathKey = "GKEHub"
const HealthcareBasePathKey = "Healthcare"
const IAMBetaBasePathKey = "IAMBeta"
const IAMWorkforcePoolBasePathKey = "IAMWorkforcePool"
const IapBasePathKey = "Iap"
const IdentityPlatformBasePathKey = "IdentityPlatform"
const KMSBasePathKey = "KMS"
const LoggingBasePathKey = "Logging"
const MemcacheBasePathKey = "Memcache"
const MLEngineBasePathKey = "MLEngine"
const MonitoringBasePathKey = "Monitoring"
const NetworkManagementBasePathKey = "NetworkManagement"
const NetworkServicesBasePathKey = "NetworkServices"
const NotebooksBasePathKey = "Notebooks"
const OSConfigBasePathKey = "OSConfig"
const OSLoginBasePathKey = "OSLogin"
const PrivatecaBasePathKey = "Privateca"
const PubsubBasePathKey = "Pubsub"
const PubsubLiteBasePathKey = "PubsubLite"
const RedisBasePathKey = "Redis"
const ResourceManagerBasePathKey = "ResourceManager"
const SecretManagerBasePathKey = "SecretManager"
const SecurityCenterBasePathKey = "SecurityCenter"
const ServiceManagementBasePathKey = "ServiceManagement"
const ServiceUsageBasePathKey = "ServiceUsage"
const SourceRepoBasePathKey = "SourceRepo"
const SpannerBasePathKey = "Spanner"
const SQLBasePathKey = "SQL"
const StorageBasePathKey = "Storage"
const StorageTransferBasePathKey = "StorageTransfer"
const TagsBasePathKey = "Tags"
const TPUBasePathKey = "TPU"
const VertexAIBasePathKey = "VertexAI"
const VPCAccessBasePathKey = "VPCAccess"
const WorkflowsBasePathKey = "Workflows"
const CloudBillingBasePathKey = "CloudBilling"
const ComposerBasePathKey = "Composer"
const ContainerBasePathKey = "Container"
const DataflowBasePathKey = "Dataflow"
const IAMBasePathKey = "IAM"
const IamCredentialsBasePathKey = "IamCredentials"
const ResourceManagerV3BasePathKey = "ResourceManagerV3"
const ServiceNetworkingBasePathKey = "ServiceNetworking"
const BigtableAdminBasePathKey = "BigtableAdmin"
const ContainerAwsBasePathKey = "ContainerAws"
const ContainerAzureBasePathKey = "ContainerAzure"

// Generated product base paths
var DefaultBasePaths = map[string]string{
	AccessApprovalBasePathKey:       "https://accessapproval.googleapis.com/v1/",
	AccessContextManagerBasePathKey: "https://accesscontextmanager.googleapis.com/v1/",
	ActiveDirectoryBasePathKey:      "https://managedidentities.googleapis.com/v1/",
	ApigeeBasePathKey:               "https://apigee.googleapis.com/v1/",
	AppEngineBasePathKey:            "https://appengine.googleapis.com/v1/",
	ArtifactRegistryBasePathKey:     "https://artifactregistry.googleapis.com/v1/",
	BeyondcorpBasePathKey:           "https://beyondcorp.googleapis.com/v1/",
	BigQueryBasePathKey:             "https://bigquery.googleapis.com/bigquery/v2/",
	BigqueryAnalyticsHubBasePathKey: "https://analyticshub.googleapis.com/v1/",
	BigqueryConnectionBasePathKey:   "https://bigqueryconnection.googleapis.com/v1/",
	BigqueryDataTransferBasePathKey: "https://bigquerydatatransfer.googleapis.com/v1/",
	BigqueryReservationBasePathKey:  "https://bigqueryreservation.googleapis.com/v1/",
	BigtableBasePathKey:             "https://bigtableadmin.googleapis.com/v2/",
	BillingBasePathKey:              "https://billingbudgets.googleapis.com/v1/",
	BinaryAuthorizationBasePathKey:  "https://binaryauthorization.googleapis.com/v1/",
	CertificateManagerBasePathKey:   "https://certificatemanager.googleapis.com/v1/",
	CloudAssetBasePathKey:           "https://cloudasset.googleapis.com/v1/",
	CloudBuildBasePathKey:           "https://cloudbuild.googleapis.com/v1/",
	CloudFunctionsBasePathKey:       "https://cloudfunctions.googleapis.com/v1/",
	Cloudfunctions2BasePathKey:      "https://cloudfunctions.googleapis.com/v2/",
	CloudIdentityBasePathKey:        "https://cloudidentity.googleapis.com/v1/",
	CloudIdsBasePathKey:             "https://ids.googleapis.com/v1/",
	CloudIotBasePathKey:             "https://cloudiot.googleapis.com/v1/",
	CloudRunBasePathKey:             "https://{{location}}-run.googleapis.com/",
	CloudRunV2BasePathKey:           "https://run.googleapis.com/v2/",
	CloudSchedulerBasePathKey:       "https://cloudscheduler.googleapis.com/v1/",
	CloudTasksBasePathKey:           "https://cloudtasks.googleapis.com/v2/",
	ComputeBasePathKey:              "https://compute.googleapis.com/compute/v1/",
	ContainerAnalysisBasePathKey:    "https://containeranalysis.googleapis.com/v1/",
	DataCatalogBasePathKey:          "https://datacatalog.googleapis.com/v1/",
	DataFusionBasePathKey:           "https://datafusion.googleapis.com/v1/",
	DataLossPreventionBasePathKey:   "https://dlp.googleapis.com/v2/",
	DataprocBasePathKey:             "https://dataproc.googleapis.com/v1/",
	DataprocMetastoreBasePathKey:    "https://metastore.googleapis.com/v1/",
	DatastoreBasePathKey:            "https://datastore.googleapis.com/v1/",
	DatastreamBasePathKey:           "https://datastream.googleapis.com/v1/",
	DeploymentManagerBasePathKey:    "https://www.googleapis.com/deploymentmanager/v2/",
	DialogflowBasePathKey:           "https://dialogflow.googleapis.com/v2/",
	DialogflowCXBasePathKey:         "https://{{location}}-dialogflow.googleapis.com/v3/",
	DNSBasePathKey:                  "https://dns.googleapis.com/dns/v1/",
	DocumentAIBasePathKey:           "https://{{location}}-documentai.googleapis.com/v1/",
	EssentialContactsBasePathKey:    "https://essentialcontacts.googleapis.com/v1/",
	FilestoreBasePathKey:            "https://file.googleapis.com/v1/",
	FirestoreBasePathKey:            "https://firestore.googleapis.com/v1/",
	GameServicesBasePathKey:         "https://gameservices.googleapis.com/v1/",
	GKEHubBasePathKey:               "https://gkehub.googleapis.com/v1/",
	HealthcareBasePathKey:           "https://healthcare.googleapis.com/v1/",
	IAMBetaBasePathKey:              "https://iam.googleapis.com/v1/",
	IAMWorkforcePoolBasePathKey:     "https://iam.googleapis.com/v1/",
	IapBasePathKey:                  "https://iap.googleapis.com/v1/",
	IdentityPlatformBasePathKey:     "https://identitytoolkit.googleapis.com/v2/",
	KMSBasePathKey:                  "https://cloudkms.googleapis.com/v1/",
	LoggingBasePathKey:              "https://logging.googleapis.com/v2/",
	MemcacheBasePathKey:             "https://memcache.googleapis.com/v1/",
	MLEngineBasePathKey:             "https://ml.googleapis.com/v1/",
	MonitoringBasePathKey:           "https://monitoring.googleapis.com/",
	NetworkManagementBasePathKey:    "https://networkmanagement.googleapis.com/v1/",
	NetworkServicesBasePathKey:      "https://networkservices.googleapis.com/v1/",
	NotebooksBasePathKey:            "https://notebooks.googleapis.com/v1/",
	OSConfigBasePathKey:             "https://osconfig.googleapis.com/v1/",
	OSLoginBasePathKey:              "https://oslogin.googleapis.com/v1/",
	PrivatecaBasePathKey:            "https://privateca.googleapis.com/v1/",
	PubsubBasePathKey:               "https://pubsub.googleapis.com/v1/",
	PubsubLiteBasePathKey:           "https://{{region}}-pubsublite.googleapis.com/v1/admin/",
	RedisBasePathKey:                "https://redis.googleapis.com/v1/",
	ResourceManagerBasePathKey:      "https://cloudresourcemanager.googleapis.com/v1/",
	SecretManagerBasePathKey:        "https://secretmanager.googleapis.com/v1/",
	SecurityCenterBasePathKey:       "https://securitycenter.googleapis.com/v1/",
	ServiceManagementBasePathKey:    "https://servicemanagement.googleapis.com/v1/",
	ServiceUsageBasePathKey:         "https://serviceusage.googleapis.com/v1/",
	SourceRepoBasePathKey:           "https://sourcerepo.googleapis.com/v1/",
	SpannerBasePathKey:              "https://spanner.googleapis.com/v1/",
	SQLBasePathKey:                  "https://sqladmin.googleapis.com/sql/v1beta4/",
	StorageBasePathKey:              "https://storage.googleapis.com/storage/v1/",
	StorageTransferBasePathKey:      "https://storagetransfer.googleapis.com/v1/",
	TagsBasePathKey:                 "https://cloudresourcemanager.googleapis.com/v3/",
	TPUBasePathKey:                  "https://tpu.googleapis.com/v1/",
	VertexAIBasePathKey:             "https://{{region}}-aiplatform.googleapis.com/v1/",
	VPCAccessBasePathKey:            "https://vpcaccess.googleapis.com/v1/",
	WorkflowsBasePathKey:            "https://workflows.googleapis.com/v1/",
	CloudBillingBasePathKey:         "https://cloudbilling.googleapis.com/v1/",
	ComposerBasePathKey:             "https://composer.googleapis.com/v1/",
	ContainerBasePathKey:            "https://container.googleapis.com/v1/",
	DataflowBasePathKey:             "https://dataflow.googleapis.com/v1b3/",
	IAMBasePathKey:                  "https://iam.googleapis.com/v1/",
	IamCredentialsBasePathKey:       "https://iamcredentials.googleapis.com/v1/",
	ResourceManagerV3BasePathKey:    "https://cloudresourcemanager.googleapis.com/v3/",
	ServiceNetworkingBasePathKey:    "https://servicenetworking.googleapis.com/v1/",
	BigtableAdminBasePathKey:        "https://bigtableadmin.googleapis.com/v2/",
	ContainerAwsBasePathKey:         "https://{{location}}-gkemulticloud.googleapis.com/v1/",
	ContainerAzureBasePathKey:       "https://{{location}}-gkemulticloud.googleapis.com/v1/",
}

var DefaultClientScopes = []string{
	"https://www.googleapis.com/auth/cloud-platform",
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
	if c.RequestReason != "" {
		headerTransport.Set("X-Goog-Request-Reason", c.RequestReason)
	}

	// Ensure $userProject is set for all HTTP requests using the client if specified by the provider config
	// See https://cloud.google.com/apis/docs/system-parameters
	if c.UserProjectOverride && c.BillingProject != "" {
		headerTransport.Set("X-Goog-User-Project", c.BillingProject)
	}

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

	// gRPC Logging setup
	logger := logrus.StandardLogger()

	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&Formatter{
		TimestampFormat: "2006/01/02 15:04:05",
		LogFormat:       "%time% [%lvl%] %msg% \n",
	})

	alwaysLoggingDeciderClient := func(ctx context.Context, fullMethodName string) bool { return true }
	grpc_logrus.ReplaceGrpcLogger(logrus.NewEntry(logger))

	c.gRPCLoggingOptions = append(
		c.gRPCLoggingOptions, option.WithGRPCDialOption(grpc.WithUnaryInterceptor(
			grpc_logrus.PayloadUnaryClientInterceptor(logrus.NewEntry(logger), alwaysLoggingDeciderClient))),
		option.WithGRPCDialOption(grpc.WithStreamInterceptor(
			grpc_logrus.PayloadStreamClientInterceptor(logrus.NewEntry(logger), alwaysLoggingDeciderClient))),
	)

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
	log.Printf("[INFO] Instantiating GCE client for path %s", c.ComputeBasePath)
	clientCompute, err := compute.NewService(c.context, option.WithHTTPClient(c.client))
	if err != nil {
		log.Printf("[WARN] Error creating client compute: %s", err)
		return nil
	}
	clientCompute.UserAgent = userAgent
	clientCompute.BasePath = c.ComputeBasePath

	return clientCompute
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

// For object uploads, we need to override the specific timeout because they are long, synchronous operations.
func (c *Config) NewStorageClientWithTimeoutOverride(userAgent string, timeout time.Duration) *storage.Service {
	storageClientBasePath := c.StorageBasePath
	log.Printf("[INFO] Instantiating Google Storage client for path %s", storageClientBasePath)
	// Copy the existing HTTP client (which has no unexported fields [as of Oct 2021 at least], so this is safe).
	// We have to do this because otherwise we will accidentally change the timeout for all other
	// synchronous operations, which would not be desirable.
	httpClient := &http.Client{
		Transport:     c.client.Transport,
		CheckRedirect: c.client.CheckRedirect,
		Jar:           c.client.Jar,
		Timeout:       timeout,
	}
	clientStorage, err := storage.NewService(c.context, option.WithHTTPClient(httpClient))
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

func (c *Config) NewResourceManagerV3Client(userAgent string) *resourceManagerV3.Service {
	resourceManagerV3BasePath := removeBasePathVersion(c.ResourceManagerV3BasePath)
	log.Printf("[INFO] Instantiating Google Cloud ResourceManager V3 client for path %s", resourceManagerV3BasePath)
	clientResourceManagerV3, err := resourceManagerV3.NewService(c.context, option.WithHTTPClient(c.client))
	if err != nil {
		log.Printf("[WARN] Error creating client resource manager v3: %s", err)
		return nil
	}
	clientResourceManagerV3.UserAgent = userAgent
	clientResourceManagerV3.BasePath = resourceManagerV3BasePath

	return clientResourceManagerV3
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
		UserAgent:           userAgent,
		TokenSource:         c.tokenSource,
		gRPCLoggingOptions:  c.gRPCLoggingOptions,
		BillingProject:      c.BillingProject,
		UserProjectOverride: c.UserProjectOverride,
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
	c.ArtifactRegistryBasePath = DefaultBasePaths[ArtifactRegistryBasePathKey]
	c.BeyondcorpBasePath = DefaultBasePaths[BeyondcorpBasePathKey]
	c.BigQueryBasePath = DefaultBasePaths[BigQueryBasePathKey]
	c.BigqueryAnalyticsHubBasePath = DefaultBasePaths[BigqueryAnalyticsHubBasePathKey]
	c.BigqueryConnectionBasePath = DefaultBasePaths[BigqueryConnectionBasePathKey]
	c.BigqueryDataTransferBasePath = DefaultBasePaths[BigqueryDataTransferBasePathKey]
	c.BigqueryReservationBasePath = DefaultBasePaths[BigqueryReservationBasePathKey]
	c.BigtableBasePath = DefaultBasePaths[BigtableBasePathKey]
	c.BillingBasePath = DefaultBasePaths[BillingBasePathKey]
	c.BinaryAuthorizationBasePath = DefaultBasePaths[BinaryAuthorizationBasePathKey]
	c.CertificateManagerBasePath = DefaultBasePaths[CertificateManagerBasePathKey]
	c.CloudAssetBasePath = DefaultBasePaths[CloudAssetBasePathKey]
	c.CloudBuildBasePath = DefaultBasePaths[CloudBuildBasePathKey]
	c.CloudFunctionsBasePath = DefaultBasePaths[CloudFunctionsBasePathKey]
	c.Cloudfunctions2BasePath = DefaultBasePaths[Cloudfunctions2BasePathKey]
	c.CloudIdentityBasePath = DefaultBasePaths[CloudIdentityBasePathKey]
	c.CloudIdsBasePath = DefaultBasePaths[CloudIdsBasePathKey]
	c.CloudIotBasePath = DefaultBasePaths[CloudIotBasePathKey]
	c.CloudRunBasePath = DefaultBasePaths[CloudRunBasePathKey]
	c.CloudRunV2BasePath = DefaultBasePaths[CloudRunV2BasePathKey]
	c.CloudSchedulerBasePath = DefaultBasePaths[CloudSchedulerBasePathKey]
	c.CloudTasksBasePath = DefaultBasePaths[CloudTasksBasePathKey]
	c.ComputeBasePath = DefaultBasePaths[ComputeBasePathKey]
	c.ContainerAnalysisBasePath = DefaultBasePaths[ContainerAnalysisBasePathKey]
	c.DataCatalogBasePath = DefaultBasePaths[DataCatalogBasePathKey]
	c.DataFusionBasePath = DefaultBasePaths[DataFusionBasePathKey]
	c.DataLossPreventionBasePath = DefaultBasePaths[DataLossPreventionBasePathKey]
	c.DataprocBasePath = DefaultBasePaths[DataprocBasePathKey]
	c.DataprocMetastoreBasePath = DefaultBasePaths[DataprocMetastoreBasePathKey]
	c.DatastoreBasePath = DefaultBasePaths[DatastoreBasePathKey]
	c.DatastreamBasePath = DefaultBasePaths[DatastreamBasePathKey]
	c.DeploymentManagerBasePath = DefaultBasePaths[DeploymentManagerBasePathKey]
	c.DialogflowBasePath = DefaultBasePaths[DialogflowBasePathKey]
	c.DialogflowCXBasePath = DefaultBasePaths[DialogflowCXBasePathKey]
	c.DNSBasePath = DefaultBasePaths[DNSBasePathKey]
	c.DocumentAIBasePath = DefaultBasePaths[DocumentAIBasePathKey]
	c.EssentialContactsBasePath = DefaultBasePaths[EssentialContactsBasePathKey]
	c.FilestoreBasePath = DefaultBasePaths[FilestoreBasePathKey]
	c.FirestoreBasePath = DefaultBasePaths[FirestoreBasePathKey]
	c.GameServicesBasePath = DefaultBasePaths[GameServicesBasePathKey]
	c.GKEHubBasePath = DefaultBasePaths[GKEHubBasePathKey]
	c.HealthcareBasePath = DefaultBasePaths[HealthcareBasePathKey]
	c.IAMBetaBasePath = DefaultBasePaths[IAMBetaBasePathKey]
	c.IAMWorkforcePoolBasePath = DefaultBasePaths[IAMWorkforcePoolBasePathKey]
	c.IapBasePath = DefaultBasePaths[IapBasePathKey]
	c.IdentityPlatformBasePath = DefaultBasePaths[IdentityPlatformBasePathKey]
	c.KMSBasePath = DefaultBasePaths[KMSBasePathKey]
	c.LoggingBasePath = DefaultBasePaths[LoggingBasePathKey]
	c.MemcacheBasePath = DefaultBasePaths[MemcacheBasePathKey]
	c.MLEngineBasePath = DefaultBasePaths[MLEngineBasePathKey]
	c.MonitoringBasePath = DefaultBasePaths[MonitoringBasePathKey]
	c.NetworkManagementBasePath = DefaultBasePaths[NetworkManagementBasePathKey]
	c.NetworkServicesBasePath = DefaultBasePaths[NetworkServicesBasePathKey]
	c.NotebooksBasePath = DefaultBasePaths[NotebooksBasePathKey]
	c.OSConfigBasePath = DefaultBasePaths[OSConfigBasePathKey]
	c.OSLoginBasePath = DefaultBasePaths[OSLoginBasePathKey]
	c.PrivatecaBasePath = DefaultBasePaths[PrivatecaBasePathKey]
	c.PubsubBasePath = DefaultBasePaths[PubsubBasePathKey]
	c.PubsubLiteBasePath = DefaultBasePaths[PubsubLiteBasePathKey]
	c.RedisBasePath = DefaultBasePaths[RedisBasePathKey]
	c.ResourceManagerBasePath = DefaultBasePaths[ResourceManagerBasePathKey]
	c.SecretManagerBasePath = DefaultBasePaths[SecretManagerBasePathKey]
	c.SecurityCenterBasePath = DefaultBasePaths[SecurityCenterBasePathKey]
	c.ServiceManagementBasePath = DefaultBasePaths[ServiceManagementBasePathKey]
	c.ServiceUsageBasePath = DefaultBasePaths[ServiceUsageBasePathKey]
	c.SourceRepoBasePath = DefaultBasePaths[SourceRepoBasePathKey]
	c.SpannerBasePath = DefaultBasePaths[SpannerBasePathKey]
	c.SQLBasePath = DefaultBasePaths[SQLBasePathKey]
	c.StorageBasePath = DefaultBasePaths[StorageBasePathKey]
	c.StorageTransferBasePath = DefaultBasePaths[StorageTransferBasePathKey]
	c.TagsBasePath = DefaultBasePaths[TagsBasePathKey]
	c.TPUBasePath = DefaultBasePaths[TPUBasePathKey]
	c.VertexAIBasePath = DefaultBasePaths[VertexAIBasePathKey]
	c.VPCAccessBasePath = DefaultBasePaths[VPCAccessBasePathKey]
	c.WorkflowsBasePath = DefaultBasePaths[WorkflowsBasePathKey]

	// Handwritten Products / Versioned / Atypical Entries
	c.CloudBillingBasePath = DefaultBasePaths[CloudBillingBasePathKey]
	c.ComposerBasePath = DefaultBasePaths[ComposerBasePathKey]
	c.ContainerBasePath = DefaultBasePaths[ContainerBasePathKey]
	c.DataprocBasePath = DefaultBasePaths[DataprocBasePathKey]
	c.DataflowBasePath = DefaultBasePaths[DataflowBasePathKey]
	c.IamCredentialsBasePath = DefaultBasePaths[IamCredentialsBasePathKey]
	c.ResourceManagerV3BasePath = DefaultBasePaths[ResourceManagerV3BasePathKey]
	c.IAMBasePath = DefaultBasePaths[IAMBasePathKey]
	c.ServiceNetworkingBasePath = DefaultBasePaths[ServiceNetworkingBasePathKey]
	c.BigQueryBasePath = DefaultBasePaths[BigQueryBasePathKey]
	c.BigtableAdminBasePath = DefaultBasePaths[BigtableAdminBasePathKey]
}
