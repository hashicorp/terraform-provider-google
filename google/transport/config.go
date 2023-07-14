// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"

	"github.com/hashicorp/terraform-provider-google/google/verify"

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
	runadminv2 "google.golang.org/api/run/v2"
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

type ProviderMeta struct {
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
	BatchingConfig                     *BatchingConfig
	UserProjectOverride                bool
	RequestReason                      string
	RequestTimeout                     time.Duration
	// PollInterval is passed to resource.StateChangeConf in common_operation.go
	// It controls the interval at which we poll for successful operations
	PollInterval time.Duration

	Client             *http.Client
	Context            context.Context
	UserAgent          string
	gRPCLoggingOptions []option.ClientOption

	tokenSource oauth2.TokenSource

	AccessApprovalBasePath           string
	AccessContextManagerBasePath     string
	ActiveDirectoryBasePath          string
	AlloydbBasePath                  string
	ApigeeBasePath                   string
	AppEngineBasePath                string
	ArtifactRegistryBasePath         string
	BeyondcorpBasePath               string
	BigQueryBasePath                 string
	BigqueryAnalyticsHubBasePath     string
	BigqueryConnectionBasePath       string
	BigqueryDatapolicyBasePath       string
	BigqueryDataTransferBasePath     string
	BigqueryReservationBasePath      string
	BigtableBasePath                 string
	BillingBasePath                  string
	BinaryAuthorizationBasePath      string
	CertificateManagerBasePath       string
	CloudAssetBasePath               string
	CloudBuildBasePath               string
	Cloudbuildv2BasePath             string
	CloudFunctionsBasePath           string
	Cloudfunctions2BasePath          string
	CloudIdentityBasePath            string
	CloudIdsBasePath                 string
	CloudIotBasePath                 string
	CloudRunBasePath                 string
	CloudRunV2BasePath               string
	CloudSchedulerBasePath           string
	CloudTasksBasePath               string
	ComputeBasePath                  string
	ContainerAnalysisBasePath        string
	ContainerAttachedBasePath        string
	DatabaseMigrationServiceBasePath string
	DataCatalogBasePath              string
	DataFusionBasePath               string
	DataLossPreventionBasePath       string
	DataplexBasePath                 string
	DataprocBasePath                 string
	DataprocMetastoreBasePath        string
	DatastoreBasePath                string
	DatastreamBasePath               string
	DeploymentManagerBasePath        string
	DialogflowBasePath               string
	DialogflowCXBasePath             string
	DNSBasePath                      string
	DocumentAIBasePath               string
	EssentialContactsBasePath        string
	FilestoreBasePath                string
	FirestoreBasePath                string
	GameServicesBasePath             string
	GKEBackupBasePath                string
	GKEHubBasePath                   string
	GKEHub2BasePath                  string
	HealthcareBasePath               string
	IAM2BasePath                     string
	IAMBetaBasePath                  string
	IAMWorkforcePoolBasePath         string
	IapBasePath                      string
	IdentityPlatformBasePath         string
	KMSBasePath                      string
	LoggingBasePath                  string
	LookerBasePath                   string
	MemcacheBasePath                 string
	MLEngineBasePath                 string
	MonitoringBasePath               string
	NetworkManagementBasePath        string
	NetworkSecurityBasePath          string
	NetworkServicesBasePath          string
	NotebooksBasePath                string
	OSConfigBasePath                 string
	OSLoginBasePath                  string
	PrivatecaBasePath                string
	PublicCABasePath                 string
	PubsubBasePath                   string
	PubsubLiteBasePath               string
	RedisBasePath                    string
	ResourceManagerBasePath          string
	SecretManagerBasePath            string
	SecurityCenterBasePath           string
	ServiceManagementBasePath        string
	ServiceUsageBasePath             string
	SourceRepoBasePath               string
	SpannerBasePath                  string
	SQLBasePath                      string
	StorageBasePath                  string
	StorageTransferBasePath          string
	TagsBasePath                     string
	TPUBasePath                      string
	VertexAIBasePath                 string
	VPCAccessBasePath                string
	WorkflowsBasePath                string

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
	TagsLocationBasePath      string

	// dcl
	ContainerAwsBasePath   string
	ContainerAzureBasePath string

	RequestBatcherServiceUsage *RequestBatcher
	RequestBatcherIam          *RequestBatcher
}

const AccessApprovalBasePathKey = "AccessApproval"
const AccessContextManagerBasePathKey = "AccessContextManager"
const ActiveDirectoryBasePathKey = "ActiveDirectory"
const AlloydbBasePathKey = "Alloydb"
const ApigeeBasePathKey = "Apigee"
const AppEngineBasePathKey = "AppEngine"
const ArtifactRegistryBasePathKey = "ArtifactRegistry"
const BeyondcorpBasePathKey = "Beyondcorp"
const BigQueryBasePathKey = "BigQuery"
const BigqueryAnalyticsHubBasePathKey = "BigqueryAnalyticsHub"
const BigqueryConnectionBasePathKey = "BigqueryConnection"
const BigqueryDatapolicyBasePathKey = "BigqueryDatapolicy"
const BigqueryDataTransferBasePathKey = "BigqueryDataTransfer"
const BigqueryReservationBasePathKey = "BigqueryReservation"
const BigtableBasePathKey = "Bigtable"
const BillingBasePathKey = "Billing"
const BinaryAuthorizationBasePathKey = "BinaryAuthorization"
const CertificateManagerBasePathKey = "CertificateManager"
const CloudAssetBasePathKey = "CloudAsset"
const CloudBuildBasePathKey = "CloudBuild"
const Cloudbuildv2BasePathKey = "Cloudbuildv2"
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
const ContainerAttachedBasePathKey = "ContainerAttached"
const DatabaseMigrationServiceBasePathKey = "DatabaseMigrationService"
const DataCatalogBasePathKey = "DataCatalog"
const DataFusionBasePathKey = "DataFusion"
const DataLossPreventionBasePathKey = "DataLossPrevention"
const DataplexBasePathKey = "Dataplex"
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
const GKEBackupBasePathKey = "GKEBackup"
const GKEHubBasePathKey = "GKEHub"
const GKEHub2BasePathKey = "GKEHub2"
const HealthcareBasePathKey = "Healthcare"
const IAM2BasePathKey = "IAM2"
const IAMBetaBasePathKey = "IAMBeta"
const IAMWorkforcePoolBasePathKey = "IAMWorkforcePool"
const IapBasePathKey = "Iap"
const IdentityPlatformBasePathKey = "IdentityPlatform"
const KMSBasePathKey = "KMS"
const LoggingBasePathKey = "Logging"
const LookerBasePathKey = "Looker"
const MemcacheBasePathKey = "Memcache"
const MLEngineBasePathKey = "MLEngine"
const MonitoringBasePathKey = "Monitoring"
const NetworkManagementBasePathKey = "NetworkManagement"
const NetworkSecurityBasePathKey = "NetworkSecurity"
const NetworkServicesBasePathKey = "NetworkServices"
const NotebooksBasePathKey = "Notebooks"
const OSConfigBasePathKey = "OSConfig"
const OSLoginBasePathKey = "OSLogin"
const PrivatecaBasePathKey = "Privateca"
const PublicCABasePathKey = "PublicCA"
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
const TagsLocationBasePathKey = "TagsLocation"

// Generated product base paths
var DefaultBasePaths = map[string]string{
	AccessApprovalBasePathKey:           "https://accessapproval.googleapis.com/v1/",
	AccessContextManagerBasePathKey:     "https://accesscontextmanager.googleapis.com/v1/",
	ActiveDirectoryBasePathKey:          "https://managedidentities.googleapis.com/v1/",
	AlloydbBasePathKey:                  "https://alloydb.googleapis.com/v1/",
	ApigeeBasePathKey:                   "https://apigee.googleapis.com/v1/",
	AppEngineBasePathKey:                "https://appengine.googleapis.com/v1/",
	ArtifactRegistryBasePathKey:         "https://artifactregistry.googleapis.com/v1/",
	BeyondcorpBasePathKey:               "https://beyondcorp.googleapis.com/v1/",
	BigQueryBasePathKey:                 "https://bigquery.googleapis.com/bigquery/v2/",
	BigqueryAnalyticsHubBasePathKey:     "https://analyticshub.googleapis.com/v1/",
	BigqueryConnectionBasePathKey:       "https://bigqueryconnection.googleapis.com/v1/",
	BigqueryDatapolicyBasePathKey:       "https://bigquerydatapolicy.googleapis.com/v1/",
	BigqueryDataTransferBasePathKey:     "https://bigquerydatatransfer.googleapis.com/v1/",
	BigqueryReservationBasePathKey:      "https://bigqueryreservation.googleapis.com/v1/",
	BigtableBasePathKey:                 "https://bigtableadmin.googleapis.com/v2/",
	BillingBasePathKey:                  "https://billingbudgets.googleapis.com/v1/",
	BinaryAuthorizationBasePathKey:      "https://binaryauthorization.googleapis.com/v1/",
	CertificateManagerBasePathKey:       "https://certificatemanager.googleapis.com/v1/",
	CloudAssetBasePathKey:               "https://cloudasset.googleapis.com/v1/",
	CloudBuildBasePathKey:               "https://cloudbuild.googleapis.com/v1/",
	Cloudbuildv2BasePathKey:             "https://cloudbuild.googleapis.com/v2/",
	CloudFunctionsBasePathKey:           "https://cloudfunctions.googleapis.com/v1/",
	Cloudfunctions2BasePathKey:          "https://cloudfunctions.googleapis.com/v2/",
	CloudIdentityBasePathKey:            "https://cloudidentity.googleapis.com/v1/",
	CloudIdsBasePathKey:                 "https://ids.googleapis.com/v1/",
	CloudIotBasePathKey:                 "https://cloudiot.googleapis.com/v1/",
	CloudRunBasePathKey:                 "https://{{location}}-run.googleapis.com/",
	CloudRunV2BasePathKey:               "https://run.googleapis.com/v2/",
	CloudSchedulerBasePathKey:           "https://cloudscheduler.googleapis.com/v1/",
	CloudTasksBasePathKey:               "https://cloudtasks.googleapis.com/v2/",
	ComputeBasePathKey:                  "https://compute.googleapis.com/compute/v1/",
	ContainerAnalysisBasePathKey:        "https://containeranalysis.googleapis.com/v1/",
	ContainerAttachedBasePathKey:        "https://{{location}}-gkemulticloud.googleapis.com/v1/",
	DatabaseMigrationServiceBasePathKey: "https://datamigration.googleapis.com/v1/",
	DataCatalogBasePathKey:              "https://datacatalog.googleapis.com/v1/",
	DataFusionBasePathKey:               "https://datafusion.googleapis.com/v1/",
	DataLossPreventionBasePathKey:       "https://dlp.googleapis.com/v2/",
	DataplexBasePathKey:                 "https://dataplex.googleapis.com/v1/",
	DataprocBasePathKey:                 "https://dataproc.googleapis.com/v1/",
	DataprocMetastoreBasePathKey:        "https://metastore.googleapis.com/v1/",
	DatastoreBasePathKey:                "https://datastore.googleapis.com/v1/",
	DatastreamBasePathKey:               "https://datastream.googleapis.com/v1/",
	DeploymentManagerBasePathKey:        "https://www.googleapis.com/deploymentmanager/v2/",
	DialogflowBasePathKey:               "https://dialogflow.googleapis.com/v2/",
	DialogflowCXBasePathKey:             "https://{{location}}-dialogflow.googleapis.com/v3/",
	DNSBasePathKey:                      "https://dns.googleapis.com/dns/v1/",
	DocumentAIBasePathKey:               "https://{{location}}-documentai.googleapis.com/v1/",
	EssentialContactsBasePathKey:        "https://essentialcontacts.googleapis.com/v1/",
	FilestoreBasePathKey:                "https://file.googleapis.com/v1/",
	FirestoreBasePathKey:                "https://firestore.googleapis.com/v1/",
	GameServicesBasePathKey:             "https://gameservices.googleapis.com/v1/",
	GKEBackupBasePathKey:                "https://gkebackup.googleapis.com/v1/",
	GKEHubBasePathKey:                   "https://gkehub.googleapis.com/v1/",
	GKEHub2BasePathKey:                  "https://gkehub.googleapis.com/v1/",
	HealthcareBasePathKey:               "https://healthcare.googleapis.com/v1/",
	IAM2BasePathKey:                     "https://iam.googleapis.com/v2/",
	IAMBetaBasePathKey:                  "https://iam.googleapis.com/v1/",
	IAMWorkforcePoolBasePathKey:         "https://iam.googleapis.com/v1/",
	IapBasePathKey:                      "https://iap.googleapis.com/v1/",
	IdentityPlatformBasePathKey:         "https://identitytoolkit.googleapis.com/v2/",
	KMSBasePathKey:                      "https://cloudkms.googleapis.com/v1/",
	LoggingBasePathKey:                  "https://logging.googleapis.com/v2/",
	LookerBasePathKey:                   "https://looker.googleapis.com/v1/",
	MemcacheBasePathKey:                 "https://memcache.googleapis.com/v1/",
	MLEngineBasePathKey:                 "https://ml.googleapis.com/v1/",
	MonitoringBasePathKey:               "https://monitoring.googleapis.com/",
	NetworkManagementBasePathKey:        "https://networkmanagement.googleapis.com/v1/",
	NetworkSecurityBasePathKey:          "https://networksecurity.googleapis.com/v1/",
	NetworkServicesBasePathKey:          "https://networkservices.googleapis.com/v1/",
	NotebooksBasePathKey:                "https://notebooks.googleapis.com/v1/",
	OSConfigBasePathKey:                 "https://osconfig.googleapis.com/v1/",
	OSLoginBasePathKey:                  "https://oslogin.googleapis.com/v1/",
	PrivatecaBasePathKey:                "https://privateca.googleapis.com/v1/",
	PublicCABasePathKey:                 "https://publicca.googleapis.com/v1/",
	PubsubBasePathKey:                   "https://pubsub.googleapis.com/v1/",
	PubsubLiteBasePathKey:               "https://{{region}}-pubsublite.googleapis.com/v1/admin/",
	RedisBasePathKey:                    "https://redis.googleapis.com/v1/",
	ResourceManagerBasePathKey:          "https://cloudresourcemanager.googleapis.com/v1/",
	SecretManagerBasePathKey:            "https://secretmanager.googleapis.com/v1/",
	SecurityCenterBasePathKey:           "https://securitycenter.googleapis.com/v1/",
	ServiceManagementBasePathKey:        "https://servicemanagement.googleapis.com/v1/",
	ServiceUsageBasePathKey:             "https://serviceusage.googleapis.com/v1/",
	SourceRepoBasePathKey:               "https://sourcerepo.googleapis.com/v1/",
	SpannerBasePathKey:                  "https://spanner.googleapis.com/v1/",
	SQLBasePathKey:                      "https://sqladmin.googleapis.com/sql/v1beta4/",
	StorageBasePathKey:                  "https://storage.googleapis.com/storage/v1/",
	StorageTransferBasePathKey:          "https://storagetransfer.googleapis.com/v1/",
	TagsBasePathKey:                     "https://cloudresourcemanager.googleapis.com/v3/",
	TPUBasePathKey:                      "https://tpu.googleapis.com/v1/",
	VertexAIBasePathKey:                 "https://{{region}}-aiplatform.googleapis.com/v1/",
	VPCAccessBasePathKey:                "https://vpcaccess.googleapis.com/v1/",
	WorkflowsBasePathKey:                "https://workflows.googleapis.com/v1/",
	CloudBillingBasePathKey:             "https://cloudbilling.googleapis.com/v1/",
	ComposerBasePathKey:                 "https://composer.googleapis.com/v1/",
	ContainerBasePathKey:                "https://container.googleapis.com/v1/",
	DataflowBasePathKey:                 "https://dataflow.googleapis.com/v1b3/",
	IAMBasePathKey:                      "https://iam.googleapis.com/v1/",
	IamCredentialsBasePathKey:           "https://iamcredentials.googleapis.com/v1/",
	ResourceManagerV3BasePathKey:        "https://cloudresourcemanager.googleapis.com/v3/",
	ServiceNetworkingBasePathKey:        "https://servicenetworking.googleapis.com/v1/",
	BigtableAdminBasePathKey:            "https://bigtableadmin.googleapis.com/v2/",
	ContainerAwsBasePathKey:             "https://{{location}}-gkemulticloud.googleapis.com/v1/",
	ContainerAzureBasePathKey:           "https://{{location}}-gkemulticloud.googleapis.com/v1/",
	TagsLocationBasePathKey:             "https://{{location}}-cloudresourcemanager.googleapis.com/v3/",
}

var DefaultClientScopes = []string{
	"https://www.googleapis.com/auth/cloud-platform",
	"https://www.googleapis.com/auth/userinfo.email",
}

func HandleSDKDefaults(d *schema.ResourceData) error {
	if d.Get("impersonate_service_account") == "" {
		d.Set("impersonate_service_account", MultiEnvDefault([]string{
			"GOOGLE_IMPERSONATE_SERVICE_ACCOUNT",
		}, nil))
	}

	if d.Get("project") == "" {
		d.Set("project", MultiEnvDefault([]string{
			"GOOGLE_PROJECT",
			"GOOGLE_CLOUD_PROJECT",
			"GCLOUD_PROJECT",
			"CLOUDSDK_CORE_PROJECT",
		}, nil))
	}

	if d.Get("billing_project") == "" {
		d.Set("billing_project", MultiEnvDefault([]string{
			"GOOGLE_BILLING_PROJECT",
		}, nil))
	}

	if d.Get("region") == "" {
		d.Set("region", MultiEnvDefault([]string{
			"GOOGLE_REGION",
			"GCLOUD_REGION",
			"CLOUDSDK_COMPUTE_REGION",
		}, nil))
	}

	if d.Get("zone") == "" {
		d.Set("zone", MultiEnvDefault([]string{
			"GOOGLE_ZONE",
			"GCLOUD_ZONE",
			"CLOUDSDK_COMPUTE_ZONE",
		}, nil))
	}

	if _, ok := d.GetOkExists("user_project_override"); !ok {
		override := MultiEnvDefault([]string{
			"USER_PROJECT_OVERRIDE",
		}, nil)

		if override != nil {
			b, err := strconv.ParseBool(override.(string))
			if err != nil {
				return err
			}
			d.Set("user_project_override", b)
		}
	}

	if d.Get("request_reason") == "" {
		d.Set("request_reason", MultiEnvDefault([]string{
			"CLOUDSDK_CORE_REQUEST_REASON",
		}, nil))
	}

	// Generated Products
	if d.Get("access_approval_custom_endpoint") == "" {
		d.Set("access_approval_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_ACCESS_APPROVAL_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[AccessApprovalBasePathKey]))
	}
	if d.Get("access_context_manager_custom_endpoint") == "" {
		d.Set("access_context_manager_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_ACCESS_CONTEXT_MANAGER_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[AccessContextManagerBasePathKey]))
	}
	if d.Get("active_directory_custom_endpoint") == "" {
		d.Set("active_directory_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_ACTIVE_DIRECTORY_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[ActiveDirectoryBasePathKey]))
	}
	if d.Get("alloydb_custom_endpoint") == "" {
		d.Set("alloydb_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_ALLOYDB_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[AlloydbBasePathKey]))
	}
	if d.Get("apigee_custom_endpoint") == "" {
		d.Set("apigee_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_APIGEE_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[ApigeeBasePathKey]))
	}
	if d.Get("app_engine_custom_endpoint") == "" {
		d.Set("app_engine_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_APP_ENGINE_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[AppEngineBasePathKey]))
	}
	if d.Get("artifact_registry_custom_endpoint") == "" {
		d.Set("artifact_registry_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_ARTIFACT_REGISTRY_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[ArtifactRegistryBasePathKey]))
	}
	if d.Get("beyondcorp_custom_endpoint") == "" {
		d.Set("beyondcorp_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_BEYONDCORP_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[BeyondcorpBasePathKey]))
	}
	if d.Get("big_query_custom_endpoint") == "" {
		d.Set("big_query_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_BIG_QUERY_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[BigQueryBasePathKey]))
	}
	if d.Get("bigquery_analytics_hub_custom_endpoint") == "" {
		d.Set("bigquery_analytics_hub_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_BIGQUERY_ANALYTICS_HUB_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[BigqueryAnalyticsHubBasePathKey]))
	}
	if d.Get("bigquery_connection_custom_endpoint") == "" {
		d.Set("bigquery_connection_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_BIGQUERY_CONNECTION_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[BigqueryConnectionBasePathKey]))
	}
	if d.Get("bigquery_datapolicy_custom_endpoint") == "" {
		d.Set("bigquery_datapolicy_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_BIGQUERY_DATAPOLICY_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[BigqueryDatapolicyBasePathKey]))
	}
	if d.Get("bigquery_data_transfer_custom_endpoint") == "" {
		d.Set("bigquery_data_transfer_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_BIGQUERY_DATA_TRANSFER_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[BigqueryDataTransferBasePathKey]))
	}
	if d.Get("bigquery_reservation_custom_endpoint") == "" {
		d.Set("bigquery_reservation_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_BIGQUERY_RESERVATION_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[BigqueryReservationBasePathKey]))
	}
	if d.Get("bigtable_custom_endpoint") == "" {
		d.Set("bigtable_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_BIGTABLE_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[BigtableBasePathKey]))
	}
	if d.Get("billing_custom_endpoint") == "" {
		d.Set("billing_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_BILLING_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[BillingBasePathKey]))
	}
	if d.Get("binary_authorization_custom_endpoint") == "" {
		d.Set("binary_authorization_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_BINARY_AUTHORIZATION_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[BinaryAuthorizationBasePathKey]))
	}
	if d.Get("certificate_manager_custom_endpoint") == "" {
		d.Set("certificate_manager_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_CERTIFICATE_MANAGER_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[CertificateManagerBasePathKey]))
	}
	if d.Get("cloud_asset_custom_endpoint") == "" {
		d.Set("cloud_asset_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_CLOUD_ASSET_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[CloudAssetBasePathKey]))
	}
	if d.Get("cloud_build_custom_endpoint") == "" {
		d.Set("cloud_build_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_CLOUD_BUILD_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[CloudBuildBasePathKey]))
	}
	if d.Get("cloudbuildv2_custom_endpoint") == "" {
		d.Set("cloudbuildv2_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_CLOUDBUILDV2_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[Cloudbuildv2BasePathKey]))
	}
	if d.Get("cloud_functions_custom_endpoint") == "" {
		d.Set("cloud_functions_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_CLOUD_FUNCTIONS_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[CloudFunctionsBasePathKey]))
	}
	if d.Get("cloudfunctions2_custom_endpoint") == "" {
		d.Set("cloudfunctions2_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_CLOUDFUNCTIONS2_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[Cloudfunctions2BasePathKey]))
	}
	if d.Get("cloud_identity_custom_endpoint") == "" {
		d.Set("cloud_identity_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_CLOUD_IDENTITY_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[CloudIdentityBasePathKey]))
	}
	if d.Get("cloud_ids_custom_endpoint") == "" {
		d.Set("cloud_ids_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_CLOUD_IDS_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[CloudIdsBasePathKey]))
	}
	if d.Get("cloud_iot_custom_endpoint") == "" {
		d.Set("cloud_iot_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_CLOUD_IOT_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[CloudIotBasePathKey]))
	}
	if d.Get("cloud_run_custom_endpoint") == "" {
		d.Set("cloud_run_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_CLOUD_RUN_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[CloudRunBasePathKey]))
	}
	if d.Get("cloud_run_v2_custom_endpoint") == "" {
		d.Set("cloud_run_v2_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_CLOUD_RUN_V2_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[CloudRunV2BasePathKey]))
	}
	if d.Get("cloud_scheduler_custom_endpoint") == "" {
		d.Set("cloud_scheduler_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_CLOUD_SCHEDULER_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[CloudSchedulerBasePathKey]))
	}
	if d.Get("cloud_tasks_custom_endpoint") == "" {
		d.Set("cloud_tasks_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_CLOUD_TASKS_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[CloudTasksBasePathKey]))
	}
	if d.Get("compute_custom_endpoint") == "" {
		d.Set("compute_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_COMPUTE_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[ComputeBasePathKey]))
	}
	if d.Get("container_analysis_custom_endpoint") == "" {
		d.Set("container_analysis_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_CONTAINER_ANALYSIS_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[ContainerAnalysisBasePathKey]))
	}
	if d.Get("container_attached_custom_endpoint") == "" {
		d.Set("container_attached_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_CONTAINER_ATTACHED_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[ContainerAttachedBasePathKey]))
	}
	if d.Get("database_migration_service_custom_endpoint") == "" {
		d.Set("database_migration_service_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_DATABASE_MIGRATION_SERVICE_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[DatabaseMigrationServiceBasePathKey]))
	}
	if d.Get("data_catalog_custom_endpoint") == "" {
		d.Set("data_catalog_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_DATA_CATALOG_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[DataCatalogBasePathKey]))
	}
	if d.Get("data_fusion_custom_endpoint") == "" {
		d.Set("data_fusion_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_DATA_FUSION_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[DataFusionBasePathKey]))
	}
	if d.Get("data_loss_prevention_custom_endpoint") == "" {
		d.Set("data_loss_prevention_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_DATA_LOSS_PREVENTION_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[DataLossPreventionBasePathKey]))
	}
	if d.Get("dataplex_custom_endpoint") == "" {
		d.Set("dataplex_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_DATAPLEX_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[DataplexBasePathKey]))
	}
	if d.Get("dataproc_custom_endpoint") == "" {
		d.Set("dataproc_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_DATAPROC_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[DataprocBasePathKey]))
	}
	if d.Get("dataproc_metastore_custom_endpoint") == "" {
		d.Set("dataproc_metastore_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_DATAPROC_METASTORE_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[DataprocMetastoreBasePathKey]))
	}
	if d.Get("datastore_custom_endpoint") == "" {
		d.Set("datastore_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_DATASTORE_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[DatastoreBasePathKey]))
	}
	if d.Get("datastream_custom_endpoint") == "" {
		d.Set("datastream_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_DATASTREAM_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[DatastreamBasePathKey]))
	}
	if d.Get("deployment_manager_custom_endpoint") == "" {
		d.Set("deployment_manager_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_DEPLOYMENT_MANAGER_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[DeploymentManagerBasePathKey]))
	}
	if d.Get("dialogflow_custom_endpoint") == "" {
		d.Set("dialogflow_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_DIALOGFLOW_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[DialogflowBasePathKey]))
	}
	if d.Get("dialogflow_cx_custom_endpoint") == "" {
		d.Set("dialogflow_cx_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_DIALOGFLOW_CX_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[DialogflowCXBasePathKey]))
	}
	if d.Get("dns_custom_endpoint") == "" {
		d.Set("dns_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_DNS_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[DNSBasePathKey]))
	}
	if d.Get("document_ai_custom_endpoint") == "" {
		d.Set("document_ai_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_DOCUMENT_AI_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[DocumentAIBasePathKey]))
	}
	if d.Get("essential_contacts_custom_endpoint") == "" {
		d.Set("essential_contacts_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_ESSENTIAL_CONTACTS_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[EssentialContactsBasePathKey]))
	}
	if d.Get("filestore_custom_endpoint") == "" {
		d.Set("filestore_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_FILESTORE_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[FilestoreBasePathKey]))
	}
	if d.Get("firestore_custom_endpoint") == "" {
		d.Set("firestore_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_FIRESTORE_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[FirestoreBasePathKey]))
	}
	if d.Get("game_services_custom_endpoint") == "" {
		d.Set("game_services_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_GAME_SERVICES_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[GameServicesBasePathKey]))
	}
	if d.Get("gke_backup_custom_endpoint") == "" {
		d.Set("gke_backup_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_GKE_BACKUP_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[GKEBackupBasePathKey]))
	}
	if d.Get("gke_hub_custom_endpoint") == "" {
		d.Set("gke_hub_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_GKE_HUB_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[GKEHubBasePathKey]))
	}
	if d.Get("gke_hub2_custom_endpoint") == "" {
		d.Set("gke_hub2_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_GKE_HUB2_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[GKEHub2BasePathKey]))
	}
	if d.Get("healthcare_custom_endpoint") == "" {
		d.Set("healthcare_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_HEALTHCARE_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[HealthcareBasePathKey]))
	}
	if d.Get("iam2_custom_endpoint") == "" {
		d.Set("iam2_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_IAM2_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[IAM2BasePathKey]))
	}
	if d.Get("iam_beta_custom_endpoint") == "" {
		d.Set("iam_beta_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_IAM_BETA_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[IAMBetaBasePathKey]))
	}
	if d.Get("iam_workforce_pool_custom_endpoint") == "" {
		d.Set("iam_workforce_pool_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_IAM_WORKFORCE_POOL_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[IAMWorkforcePoolBasePathKey]))
	}
	if d.Get("iap_custom_endpoint") == "" {
		d.Set("iap_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_IAP_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[IapBasePathKey]))
	}
	if d.Get("identity_platform_custom_endpoint") == "" {
		d.Set("identity_platform_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_IDENTITY_PLATFORM_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[IdentityPlatformBasePathKey]))
	}
	if d.Get("kms_custom_endpoint") == "" {
		d.Set("kms_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_KMS_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[KMSBasePathKey]))
	}
	if d.Get("logging_custom_endpoint") == "" {
		d.Set("logging_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_LOGGING_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[LoggingBasePathKey]))
	}
	if d.Get("looker_custom_endpoint") == "" {
		d.Set("looker_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_LOOKER_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[LookerBasePathKey]))
	}
	if d.Get("memcache_custom_endpoint") == "" {
		d.Set("memcache_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_MEMCACHE_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[MemcacheBasePathKey]))
	}
	if d.Get("ml_engine_custom_endpoint") == "" {
		d.Set("ml_engine_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_ML_ENGINE_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[MLEngineBasePathKey]))
	}
	if d.Get("monitoring_custom_endpoint") == "" {
		d.Set("monitoring_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_MONITORING_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[MonitoringBasePathKey]))
	}
	if d.Get("network_management_custom_endpoint") == "" {
		d.Set("network_management_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_NETWORK_MANAGEMENT_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[NetworkManagementBasePathKey]))
	}
	if d.Get("network_security_custom_endpoint") == "" {
		d.Set("network_security_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_NETWORK_SECURITY_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[NetworkSecurityBasePathKey]))
	}
	if d.Get("network_services_custom_endpoint") == "" {
		d.Set("network_services_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_NETWORK_SERVICES_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[NetworkServicesBasePathKey]))
	}
	if d.Get("notebooks_custom_endpoint") == "" {
		d.Set("notebooks_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_NOTEBOOKS_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[NotebooksBasePathKey]))
	}
	if d.Get("os_config_custom_endpoint") == "" {
		d.Set("os_config_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_OS_CONFIG_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[OSConfigBasePathKey]))
	}
	if d.Get("os_login_custom_endpoint") == "" {
		d.Set("os_login_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_OS_LOGIN_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[OSLoginBasePathKey]))
	}
	if d.Get("privateca_custom_endpoint") == "" {
		d.Set("privateca_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_PRIVATECA_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[PrivatecaBasePathKey]))
	}
	if d.Get("public_ca_custom_endpoint") == "" {
		d.Set("public_ca_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_PUBLIC_CA_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[PublicCABasePathKey]))
	}
	if d.Get("pubsub_custom_endpoint") == "" {
		d.Set("pubsub_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_PUBSUB_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[PubsubBasePathKey]))
	}
	if d.Get("pubsub_lite_custom_endpoint") == "" {
		d.Set("pubsub_lite_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_PUBSUB_LITE_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[PubsubLiteBasePathKey]))
	}
	if d.Get("redis_custom_endpoint") == "" {
		d.Set("redis_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_REDIS_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[RedisBasePathKey]))
	}
	if d.Get("resource_manager_custom_endpoint") == "" {
		d.Set("resource_manager_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_RESOURCE_MANAGER_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[ResourceManagerBasePathKey]))
	}
	if d.Get("secret_manager_custom_endpoint") == "" {
		d.Set("secret_manager_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_SECRET_MANAGER_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[SecretManagerBasePathKey]))
	}
	if d.Get("security_center_custom_endpoint") == "" {
		d.Set("security_center_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_SECURITY_CENTER_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[SecurityCenterBasePathKey]))
	}
	if d.Get("service_management_custom_endpoint") == "" {
		d.Set("service_management_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_SERVICE_MANAGEMENT_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[ServiceManagementBasePathKey]))
	}
	if d.Get("service_usage_custom_endpoint") == "" {
		d.Set("service_usage_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_SERVICE_USAGE_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[ServiceUsageBasePathKey]))
	}
	if d.Get("source_repo_custom_endpoint") == "" {
		d.Set("source_repo_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_SOURCE_REPO_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[SourceRepoBasePathKey]))
	}
	if d.Get("spanner_custom_endpoint") == "" {
		d.Set("spanner_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_SPANNER_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[SpannerBasePathKey]))
	}
	if d.Get("sql_custom_endpoint") == "" {
		d.Set("sql_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_SQL_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[SQLBasePathKey]))
	}
	if d.Get("storage_custom_endpoint") == "" {
		d.Set("storage_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_STORAGE_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[StorageBasePathKey]))
	}
	if d.Get("storage_transfer_custom_endpoint") == "" {
		d.Set("storage_transfer_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_STORAGE_TRANSFER_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[StorageTransferBasePathKey]))
	}
	if d.Get("tags_custom_endpoint") == "" {
		d.Set("tags_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_TAGS_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[TagsBasePathKey]))
	}
	if d.Get("tpu_custom_endpoint") == "" {
		d.Set("tpu_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_TPU_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[TPUBasePathKey]))
	}
	if d.Get("vertex_ai_custom_endpoint") == "" {
		d.Set("vertex_ai_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_VERTEX_AI_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[VertexAIBasePathKey]))
	}
	if d.Get("vpc_access_custom_endpoint") == "" {
		d.Set("vpc_access_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_VPC_ACCESS_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[VPCAccessBasePathKey]))
	}
	if d.Get("workflows_custom_endpoint") == "" {
		d.Set("workflows_custom_endpoint", MultiEnvDefault([]string{
			"GOOGLE_WORKFLOWS_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[WorkflowsBasePathKey]))
	}

	if d.Get(CloudBillingCustomEndpointEntryKey) == "" {
		d.Set(CloudBillingCustomEndpointEntryKey, MultiEnvDefault([]string{
			"GOOGLE_CLOUD_BILLING_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[CloudBillingBasePathKey]))
	}

	if d.Get(ComposerCustomEndpointEntryKey) == "" {
		d.Set(ComposerCustomEndpointEntryKey, MultiEnvDefault([]string{
			"GOOGLE_COMPOSER_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[ComposerBasePathKey]))
	}

	if d.Get(ContainerCustomEndpointEntryKey) == "" {
		d.Set(ContainerCustomEndpointEntryKey, MultiEnvDefault([]string{
			"GOOGLE_CONTAINER_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[ContainerBasePathKey]))
	}

	if d.Get(DataflowCustomEndpointEntryKey) == "" {
		d.Set(DataflowCustomEndpointEntryKey, MultiEnvDefault([]string{
			"GOOGLE_DATAFLOW_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[DataflowBasePathKey]))
	}

	if d.Get(IamCredentialsCustomEndpointEntryKey) == "" {
		d.Set(IamCredentialsCustomEndpointEntryKey, MultiEnvDefault([]string{
			"GOOGLE_IAM_CREDENTIALS_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[IamCredentialsBasePathKey]))
	}

	if d.Get(ResourceManagerV3CustomEndpointEntryKey) == "" {
		d.Set(ResourceManagerV3CustomEndpointEntryKey, MultiEnvDefault([]string{
			"GOOGLE_RESOURCE_MANAGER_V3_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[ResourceManagerV3BasePathKey]))
	}

	if d.Get(IAMCustomEndpointEntryKey) == "" {
		d.Set(IAMCustomEndpointEntryKey, MultiEnvDefault([]string{
			"GOOGLE_IAM_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[IAMBasePathKey]))
	}

	if d.Get(ServiceNetworkingCustomEndpointEntryKey) == "" {
		d.Set(ServiceNetworkingCustomEndpointEntryKey, MultiEnvDefault([]string{
			"GOOGLE_SERVICE_NETWORKING_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[ServiceNetworkingBasePathKey]))
	}

	if d.Get(TagsLocationCustomEndpointEntryKey) == "" {
		d.Set(TagsLocationCustomEndpointEntryKey, MultiEnvDefault([]string{
			"GOOGLE_TAGS_LOCATION_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[TagsLocationBasePathKey]))
	}

	if d.Get(ContainerAwsCustomEndpointEntryKey) == "" {
		d.Set(ContainerAwsCustomEndpointEntryKey, MultiEnvDefault([]string{
			"GOOGLE_CONTAINERAWS_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[ContainerAwsBasePathKey]))
	}

	if d.Get(ContainerAzureCustomEndpointEntryKey) == "" {
		d.Set(ContainerAzureCustomEndpointEntryKey, MultiEnvDefault([]string{
			"GOOGLE_CONTAINERAZURE_CUSTOM_ENDPOINT",
		}, DefaultBasePaths[ContainerAzureBasePathKey]))
	}

	return nil
}

func (c *Config) LoadAndValidate(ctx context.Context) error {
	if len(c.Scopes) == 0 {
		c.Scopes = DefaultClientScopes
	}

	c.Context = ctx

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
	headerTransport := NewTransportWithHeaders(retryTransport)
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

	c.Client = client
	c.Context = ctx
	c.Region = GetRegionFromRegionSelfLink(c.Region)
	c.RequestBatcherServiceUsage = NewRequestBatcher("Service Usage", ctx, c.BatchingConfig)
	c.RequestBatcherIam = NewRequestBatcher("IAM", ctx, c.BatchingConfig)
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

func ExpandProviderBatchingConfig(v interface{}) (*BatchingConfig, error) {
	config := &BatchingConfig{
		SendAfter:      time.Second * DefaultBatchSendIntervalSec,
		EnableBatching: true,
	}

	if v == nil {
		return config, nil
	}
	ls := v.([]interface{})
	if len(ls) == 0 || ls[0] == nil {
		return config, nil
	}

	cfgV := ls[0].(map[string]interface{})
	if sendAfterV, ok := cfgV["send_after"]; ok && sendAfterV != "" {
		SendAfter, err := time.ParseDuration(sendAfterV.(string))
		if err != nil {
			return nil, fmt.Errorf("unable to parse duration from 'send_after' value %q", sendAfterV)
		}
		config.SendAfter = SendAfter
	}

	if enable, ok := cfgV["enable_batching"]; ok {
		config.EnableBatching = enable.(bool)
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
		c.Client = oauth2.NewClient(c.Context, tokenSource) // c.Client isn't initialised fully when this code is called.

		email, err := GetCurrentUserEmail(c, c.UserAgent)
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
	c.Client = oauth2.NewClient(c.Context, tokenSource) // c.Client isn't initialised fully when this code is called.

	email, err := GetCurrentUserEmail(c, c.UserAgent)
	if err != nil {
		log.Printf("[INFO] error retrieving userinfo for your provider credentials. have you enabled the 'https://www.googleapis.com/auth/userinfo.email' scope? error: %s", err)
	}

	log.Printf("[INFO] Terraform is configured with service account impersonation, original identity: %s, impersonated identity: %s", email, c.ImpersonateServiceAccount)

	// Add the Impersonated ClientOption back in to the OAuth2 TokenSource

	tokenSource, err = c.getTokenSource(c.Scopes, false)
	if err != nil {
		return err
	}
	c.Client = oauth2.NewClient(c.Context, tokenSource) // c.Client isn't initialised fully when this code is called.

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
	clientCompute, err := compute.NewService(c.Context, option.WithHTTPClient(c.Client))
	if err != nil {
		log.Printf("[WARN] Error creating client compute: %s", err)
		return nil
	}
	clientCompute.UserAgent = userAgent
	clientCompute.BasePath = c.ComputeBasePath

	return clientCompute
}

func (c *Config) NewContainerClient(userAgent string) *container.Service {
	containerClientBasePath := RemoveBasePathVersion(c.ContainerBasePath)
	log.Printf("[INFO] Instantiating GKE client for path %s", containerClientBasePath)
	clientContainer, err := container.NewService(c.Context, option.WithHTTPClient(c.Client))
	if err != nil {
		log.Printf("[WARN] Error creating client container: %s", err)
		return nil
	}
	clientContainer.UserAgent = userAgent
	clientContainer.BasePath = containerClientBasePath

	return clientContainer
}

func (c *Config) NewDnsClient(userAgent string) *dns.Service {
	dnsClientBasePath := RemoveBasePathVersion(c.DNSBasePath)
	dnsClientBasePath = strings.ReplaceAll(dnsClientBasePath, "/dns/", "")
	log.Printf("[INFO] Instantiating Google Cloud DNS client for path %s", dnsClientBasePath)
	clientDns, err := dns.NewService(c.Context, option.WithHTTPClient(c.Client))
	if err != nil {
		log.Printf("[WARN] Error creating client dns: %s", err)
		return nil
	}
	clientDns.UserAgent = userAgent
	clientDns.BasePath = dnsClientBasePath

	return clientDns
}

func (c *Config) NewKmsClientWithCtx(ctx context.Context, userAgent string) *cloudkms.Service {
	kmsClientBasePath := RemoveBasePathVersion(c.KMSBasePath)
	log.Printf("[INFO] Instantiating Google Cloud KMS client for path %s", kmsClientBasePath)
	clientKms, err := cloudkms.NewService(ctx, option.WithHTTPClient(c.Client))
	if err != nil {
		log.Printf("[WARN] Error creating client kms: %s", err)
		return nil
	}
	clientKms.UserAgent = userAgent
	clientKms.BasePath = kmsClientBasePath

	return clientKms
}

func (c *Config) NewKmsClient(userAgent string) *cloudkms.Service {
	return c.NewKmsClientWithCtx(c.Context, userAgent)
}

func (c *Config) NewLoggingClient(userAgent string) *cloudlogging.Service {
	loggingClientBasePath := RemoveBasePathVersion(c.LoggingBasePath)
	log.Printf("[INFO] Instantiating Google Stackdriver Logging client for path %s", loggingClientBasePath)
	clientLogging, err := cloudlogging.NewService(c.Context, option.WithHTTPClient(c.Client))
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
	clientStorage, err := storage.NewService(c.Context, option.WithHTTPClient(c.Client))
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
		Transport:     c.Client.Transport,
		CheckRedirect: c.Client.CheckRedirect,
		Jar:           c.Client.Jar,
		Timeout:       timeout,
	}
	clientStorage, err := storage.NewService(c.Context, option.WithHTTPClient(httpClient))
	if err != nil {
		log.Printf("[WARN] Error creating client storage: %s", err)
		return nil
	}
	clientStorage.UserAgent = userAgent
	clientStorage.BasePath = storageClientBasePath

	return clientStorage
}

func (c *Config) NewSqlAdminClient(userAgent string) *sqladmin.Service {
	sqlClientBasePath := RemoveBasePathVersion(RemoveBasePathVersion(c.SQLBasePath))
	log.Printf("[INFO] Instantiating Google SqlAdmin client for path %s", sqlClientBasePath)
	clientSqlAdmin, err := sqladmin.NewService(c.Context, option.WithHTTPClient(c.Client))
	if err != nil {
		log.Printf("[WARN] Error creating client storage: %s", err)
		return nil
	}
	clientSqlAdmin.UserAgent = userAgent
	clientSqlAdmin.BasePath = sqlClientBasePath

	return clientSqlAdmin
}

func (c *Config) NewPubsubClient(userAgent string) *pubsub.Service {
	pubsubClientBasePath := RemoveBasePathVersion(c.PubsubBasePath)
	log.Printf("[INFO] Instantiating Google Pubsub client for path %s", pubsubClientBasePath)
	wrappedPubsubClient := ClientWithAdditionalRetries(c.Client, PubsubTopicProjectNotReady)
	clientPubsub, err := pubsub.NewService(c.Context, option.WithHTTPClient(wrappedPubsubClient))
	if err != nil {
		log.Printf("[WARN] Error creating client pubsub: %s", err)
		return nil
	}
	clientPubsub.UserAgent = userAgent
	clientPubsub.BasePath = pubsubClientBasePath

	return clientPubsub
}

func (c *Config) NewDataflowClient(userAgent string) *dataflow.Service {
	dataflowClientBasePath := RemoveBasePathVersion(c.DataflowBasePath)
	log.Printf("[INFO] Instantiating Google Dataflow client for path %s", dataflowClientBasePath)
	clientDataflow, err := dataflow.NewService(c.Context, option.WithHTTPClient(c.Client))
	if err != nil {
		log.Printf("[WARN] Error creating client dataflow: %s", err)
		return nil
	}
	clientDataflow.UserAgent = userAgent
	clientDataflow.BasePath = dataflowClientBasePath

	return clientDataflow
}

func (c *Config) NewResourceManagerClient(userAgent string) *cloudresourcemanager.Service {
	resourceManagerBasePath := RemoveBasePathVersion(c.ResourceManagerBasePath)
	log.Printf("[INFO] Instantiating Google Cloud ResourceManager client for path %s", resourceManagerBasePath)
	clientResourceManager, err := cloudresourcemanager.NewService(c.Context, option.WithHTTPClient(c.Client))
	if err != nil {
		log.Printf("[WARN] Error creating client resource manager: %s", err)
		return nil
	}
	clientResourceManager.UserAgent = userAgent
	clientResourceManager.BasePath = resourceManagerBasePath

	return clientResourceManager
}

func (c *Config) NewResourceManagerV3Client(userAgent string) *resourceManagerV3.Service {
	resourceManagerV3BasePath := RemoveBasePathVersion(c.ResourceManagerV3BasePath)
	log.Printf("[INFO] Instantiating Google Cloud ResourceManager V3 client for path %s", resourceManagerV3BasePath)
	clientResourceManagerV3, err := resourceManagerV3.NewService(c.Context, option.WithHTTPClient(c.Client))
	if err != nil {
		log.Printf("[WARN] Error creating client resource manager v3: %s", err)
		return nil
	}
	clientResourceManagerV3.UserAgent = userAgent
	clientResourceManagerV3.BasePath = resourceManagerV3BasePath

	return clientResourceManagerV3
}

func (c *Config) NewIamClient(userAgent string) *iam.Service {
	iamClientBasePath := RemoveBasePathVersion(c.IAMBasePath)
	log.Printf("[INFO] Instantiating Google Cloud IAM client for path %s", iamClientBasePath)
	clientIAM, err := iam.NewService(c.Context, option.WithHTTPClient(c.Client))
	if err != nil {
		log.Printf("[WARN] Error creating client iam: %s", err)
		return nil
	}
	clientIAM.UserAgent = userAgent
	clientIAM.BasePath = iamClientBasePath

	return clientIAM
}

func (c *Config) NewIamCredentialsClient(userAgent string) *iamcredentials.Service {
	iamCredentialsClientBasePath := RemoveBasePathVersion(c.IamCredentialsBasePath)
	log.Printf("[INFO] Instantiating Google Cloud IAMCredentials client for path %s", iamCredentialsClientBasePath)
	clientIamCredentials, err := iamcredentials.NewService(c.Context, option.WithHTTPClient(c.Client))
	if err != nil {
		log.Printf("[WARN] Error creating client iam credentials: %s", err)
		return nil
	}
	clientIamCredentials.UserAgent = userAgent
	clientIamCredentials.BasePath = iamCredentialsClientBasePath

	return clientIamCredentials
}

func (c *Config) NewServiceManClient(userAgent string) *servicemanagement.APIService {
	serviceManagementClientBasePath := RemoveBasePathVersion(c.ServiceManagementBasePath)
	log.Printf("[INFO] Instantiating Google Cloud Service Management client for path %s", serviceManagementClientBasePath)
	clientServiceMan, err := servicemanagement.NewService(c.Context, option.WithHTTPClient(c.Client))
	if err != nil {
		log.Printf("[WARN] Error creating client service management: %s", err)
		return nil
	}
	clientServiceMan.UserAgent = userAgent
	clientServiceMan.BasePath = serviceManagementClientBasePath

	return clientServiceMan
}

func (c *Config) NewServiceUsageClient(userAgent string) *serviceusage.Service {
	serviceUsageClientBasePath := RemoveBasePathVersion(c.ServiceUsageBasePath)
	log.Printf("[INFO] Instantiating Google Cloud Service Usage client for path %s", serviceUsageClientBasePath)
	clientServiceUsage, err := serviceusage.NewService(c.Context, option.WithHTTPClient(c.Client))
	if err != nil {
		log.Printf("[WARN] Error creating client service usage: %s", err)
		return nil
	}
	clientServiceUsage.UserAgent = userAgent
	clientServiceUsage.BasePath = serviceUsageClientBasePath

	return clientServiceUsage
}

func (c *Config) NewBillingClient(userAgent string) *cloudbilling.APIService {
	cloudBillingClientBasePath := RemoveBasePathVersion(c.CloudBillingBasePath)
	log.Printf("[INFO] Instantiating Google Cloud Billing client for path %s", cloudBillingClientBasePath)
	clientBilling, err := cloudbilling.NewService(c.Context, option.WithHTTPClient(c.Client))
	if err != nil {
		log.Printf("[WARN] Error creating client billing: %s", err)
		return nil
	}
	clientBilling.UserAgent = userAgent
	clientBilling.BasePath = cloudBillingClientBasePath

	return clientBilling
}

func (c *Config) NewBuildClient(userAgent string) *cloudbuild.Service {
	cloudBuildClientBasePath := RemoveBasePathVersion(c.CloudBuildBasePath)
	log.Printf("[INFO] Instantiating Google Cloud Build client for path %s", cloudBuildClientBasePath)
	clientBuild, err := cloudbuild.NewService(c.Context, option.WithHTTPClient(c.Client))
	if err != nil {
		log.Printf("[WARN] Error creating client build: %s", err)
		return nil
	}
	clientBuild.UserAgent = userAgent
	clientBuild.BasePath = cloudBuildClientBasePath

	return clientBuild
}

func (c *Config) NewCloudFunctionsClient(userAgent string) *cloudfunctions.Service {
	cloudFunctionsClientBasePath := RemoveBasePathVersion(c.CloudFunctionsBasePath)
	log.Printf("[INFO] Instantiating Google Cloud CloudFunctions Client for path %s", cloudFunctionsClientBasePath)
	clientCloudFunctions, err := cloudfunctions.NewService(c.Context, option.WithHTTPClient(c.Client))
	if err != nil {
		log.Printf("[WARN] Error creating client cloud functions: %s", err)
		return nil
	}
	clientCloudFunctions.UserAgent = userAgent
	clientCloudFunctions.BasePath = cloudFunctionsClientBasePath

	return clientCloudFunctions
}

func (c *Config) NewSourceRepoClient(userAgent string) *sourcerepo.Service {
	sourceRepoClientBasePath := RemoveBasePathVersion(c.SourceRepoBasePath)
	log.Printf("[INFO] Instantiating Google Cloud Source Repo client for path %s", sourceRepoClientBasePath)
	clientSourceRepo, err := sourcerepo.NewService(c.Context, option.WithHTTPClient(c.Client))
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
	wrappedBigQueryClient := ClientWithAdditionalRetries(c.Client, IamMemberMissing)
	clientBigQuery, err := bigquery.NewService(c.Context, option.WithHTTPClient(wrappedBigQueryClient))
	if err != nil {
		log.Printf("[WARN] Error creating client big query: %s", err)
		return nil
	}
	clientBigQuery.UserAgent = userAgent
	clientBigQuery.BasePath = bigQueryClientBasePath

	return clientBigQuery
}

func (c *Config) NewSpannerClient(userAgent string) *spanner.Service {
	spannerClientBasePath := RemoveBasePathVersion(c.SpannerBasePath)
	log.Printf("[INFO] Instantiating Google Cloud Spanner client for path %s", spannerClientBasePath)
	clientSpanner, err := spanner.NewService(c.Context, option.WithHTTPClient(c.Client))
	if err != nil {
		log.Printf("[WARN] Error creating client source repo: %s", err)
		return nil
	}
	clientSpanner.UserAgent = userAgent
	clientSpanner.BasePath = spannerClientBasePath

	return clientSpanner
}

func (c *Config) NewDataprocClient(userAgent string) *dataproc.Service {
	dataprocClientBasePath := RemoveBasePathVersion(c.DataprocBasePath)
	log.Printf("[INFO] Instantiating Google Cloud Dataproc client for path %s", dataprocClientBasePath)
	clientDataproc, err := dataproc.NewService(c.Context, option.WithHTTPClient(c.Client))
	if err != nil {
		log.Printf("[WARN] Error creating client dataproc: %s", err)
		return nil
	}
	clientDataproc.UserAgent = userAgent
	clientDataproc.BasePath = dataprocClientBasePath

	return clientDataproc
}

func (c *Config) NewCloudIoTClient(userAgent string) *cloudiot.Service {
	cloudIoTClientBasePath := RemoveBasePathVersion(c.CloudIoTBasePath)
	log.Printf("[INFO] Instantiating Google Cloud IoT Core client for path %s", cloudIoTClientBasePath)
	clientCloudIoT, err := cloudiot.NewService(c.Context, option.WithHTTPClient(c.Client))
	if err != nil {
		log.Printf("[WARN] Error creating client cloud iot: %s", err)
		return nil
	}
	clientCloudIoT.UserAgent = userAgent
	clientCloudIoT.BasePath = cloudIoTClientBasePath

	return clientCloudIoT
}

func (c *Config) NewAppEngineClient(userAgent string) *appengine.APIService {
	appEngineClientBasePath := RemoveBasePathVersion(c.AppEngineBasePath)
	log.Printf("[INFO] Instantiating App Engine client for path %s", appEngineClientBasePath)
	clientAppEngine, err := appengine.NewService(c.Context, option.WithHTTPClient(c.Client))
	if err != nil {
		log.Printf("[WARN] Error creating client appengine: %s", err)
		return nil
	}
	clientAppEngine.UserAgent = userAgent
	clientAppEngine.BasePath = appEngineClientBasePath

	return clientAppEngine
}

func (c *Config) NewComposerClient(userAgent string) *composer.Service {
	composerClientBasePath := RemoveBasePathVersion(c.ComposerBasePath)
	log.Printf("[INFO] Instantiating Cloud Composer client for path %s", composerClientBasePath)
	clientComposer, err := composer.NewService(c.Context, option.WithHTTPClient(c.Client))
	if err != nil {
		log.Printf("[WARN] Error creating client composer: %s", err)
		return nil
	}
	clientComposer.UserAgent = userAgent
	clientComposer.BasePath = composerClientBasePath

	return clientComposer
}

func (c *Config) NewServiceNetworkingClient(userAgent string) *servicenetworking.APIService {
	serviceNetworkingClientBasePath := RemoveBasePathVersion(c.ServiceNetworkingBasePath)
	log.Printf("[INFO] Instantiating Service Networking client for path %s", serviceNetworkingClientBasePath)
	clientServiceNetworking, err := servicenetworking.NewService(c.Context, option.WithHTTPClient(c.Client))
	if err != nil {
		log.Printf("[WARN] Error creating client service networking: %s", err)
		return nil
	}
	clientServiceNetworking.UserAgent = userAgent
	clientServiceNetworking.BasePath = serviceNetworkingClientBasePath

	return clientServiceNetworking
}

func (c *Config) NewStorageTransferClient(userAgent string) *storagetransfer.Service {
	storageTransferClientBasePath := RemoveBasePathVersion(c.StorageTransferBasePath)
	log.Printf("[INFO] Instantiating Google Cloud Storage Transfer client for path %s", storageTransferClientBasePath)
	clientStorageTransfer, err := storagetransfer.NewService(c.Context, option.WithHTTPClient(c.Client))
	if err != nil {
		log.Printf("[WARN] Error creating client storage transfer: %s", err)
		return nil
	}
	clientStorageTransfer.UserAgent = userAgent
	clientStorageTransfer.BasePath = storageTransferClientBasePath

	return clientStorageTransfer
}

func (c *Config) NewHealthcareClient(userAgent string) *healthcare.Service {
	healthcareClientBasePath := RemoveBasePathVersion(c.HealthcareBasePath)
	log.Printf("[INFO] Instantiating Google Cloud Healthcare client for path %s", healthcareClientBasePath)
	clientHealthcare, err := healthcare.NewService(c.Context, option.WithHTTPClient(c.Client))
	if err != nil {
		log.Printf("[WARN] Error creating client healthcare: %s", err)
		return nil
	}
	clientHealthcare.UserAgent = userAgent
	clientHealthcare.BasePath = healthcareClientBasePath

	return clientHealthcare
}

func (c *Config) NewCloudIdentityClient(userAgent string) *cloudidentity.Service {
	cloudidentityClientBasePath := RemoveBasePathVersion(c.CloudIdentityBasePath)
	log.Printf("[INFO] Instantiating Google Cloud CloudIdentity client for path %s", cloudidentityClientBasePath)
	clientCloudIdentity, err := cloudidentity.NewService(c.Context, option.WithHTTPClient(c.Client))
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
	bigtableAdminBasePath := RemoveBasePathVersion(c.BigtableAdminBasePath)
	log.Printf("[INFO] Instantiating Google Cloud BigtableAdmin for path %s", bigtableAdminBasePath)
	clientBigtable, err := bigtableadmin.NewService(c.Context, option.WithHTTPClient(c.Client))
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
	bigtableAdminBasePath := RemoveBasePathVersion(c.BigtableAdminBasePath)
	log.Printf("[INFO] Instantiating Google Cloud BigtableAdmin for path %s", bigtableAdminBasePath)
	clientBigtable, err := bigtableadmin.NewService(c.Context, option.WithHTTPClient(c.Client))
	if err != nil {
		log.Printf("[WARN] Error creating client projects instances tables: %s", err)
		return nil
	}
	clientBigtable.UserAgent = userAgent
	clientBigtable.BasePath = bigtableAdminBasePath
	clientBigtableProjectsInstancesTables := bigtableadmin.NewProjectsInstancesTablesService(clientBigtable)

	return clientBigtableProjectsInstancesTables
}

func (c *Config) NewCloudRunV2Client(userAgent string) *runadminv2.Service {
	runAdminV2ClientBasePath := RemoveBasePathVersion(RemoveBasePathVersion(c.CloudRunV2BasePath))
	log.Printf("[INFO] Instantiating Google Cloud Run Admin v2 client for path %s", runAdminV2ClientBasePath)
	clientRunAdminV2, err := runadminv2.NewService(c.Context, option.WithHTTPClient(c.Client))
	if err != nil {
		log.Printf("[WARN] Error creating client run admin: %s", err)
		return nil
	}
	clientRunAdminV2.UserAgent = userAgent
	clientRunAdminV2.BasePath = runAdminV2ClientBasePath

	return clientRunAdminV2
}

// StaticTokenSource is used to be able to identify static token sources without reflection.
type StaticTokenSource struct {
	oauth2.TokenSource
}

// Get a set of credentials with a given scope (clientScopes) based on the Config object.
// If initialCredentialsOnly is true, don't follow the impersonation settings and return the initial set of creds
// instead.
func (c *Config) GetCredentials(clientScopes []string, initialCredentialsOnly bool) (googleoauth.Credentials, error) {
	if c.AccessToken != "" {
		contents, _, err := verify.PathOrContents(c.AccessToken)
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
			TokenSource: StaticTokenSource{oauth2.StaticTokenSource(token)},
		}, nil
	}

	if c.Credentials != "" {
		contents, _, err := verify.PathOrContents(c.Credentials)
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

		creds, err := transport.Creds(c.Context, option.WithCredentialsJSON([]byte(contents)), option.WithScopes(clientScopes...))
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
	creds, err := transport.Creds(context.Background(), option.WithScopes(clientScopes...))
	if err != nil {
		return googleoauth.Credentials{}, fmt.Errorf("Attempted to load application default credentials since neither `credentials` nor `access_token` was set in the provider block.  No credentials loaded. To use your gcloud credentials, run 'gcloud auth application-default login'.  Original error: %w", err)
	}

	return *creds, nil
}

// Remove the `/{{version}}/` from a base path if present.
func RemoveBasePathVersion(url string) string {
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
	c.AlloydbBasePath = DefaultBasePaths[AlloydbBasePathKey]
	c.ApigeeBasePath = DefaultBasePaths[ApigeeBasePathKey]
	c.AppEngineBasePath = DefaultBasePaths[AppEngineBasePathKey]
	c.ArtifactRegistryBasePath = DefaultBasePaths[ArtifactRegistryBasePathKey]
	c.BeyondcorpBasePath = DefaultBasePaths[BeyondcorpBasePathKey]
	c.BigQueryBasePath = DefaultBasePaths[BigQueryBasePathKey]
	c.BigqueryAnalyticsHubBasePath = DefaultBasePaths[BigqueryAnalyticsHubBasePathKey]
	c.BigqueryConnectionBasePath = DefaultBasePaths[BigqueryConnectionBasePathKey]
	c.BigqueryDatapolicyBasePath = DefaultBasePaths[BigqueryDatapolicyBasePathKey]
	c.BigqueryDataTransferBasePath = DefaultBasePaths[BigqueryDataTransferBasePathKey]
	c.BigqueryReservationBasePath = DefaultBasePaths[BigqueryReservationBasePathKey]
	c.BigtableBasePath = DefaultBasePaths[BigtableBasePathKey]
	c.BillingBasePath = DefaultBasePaths[BillingBasePathKey]
	c.BinaryAuthorizationBasePath = DefaultBasePaths[BinaryAuthorizationBasePathKey]
	c.CertificateManagerBasePath = DefaultBasePaths[CertificateManagerBasePathKey]
	c.CloudAssetBasePath = DefaultBasePaths[CloudAssetBasePathKey]
	c.CloudBuildBasePath = DefaultBasePaths[CloudBuildBasePathKey]
	c.Cloudbuildv2BasePath = DefaultBasePaths[Cloudbuildv2BasePathKey]
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
	c.ContainerAttachedBasePath = DefaultBasePaths[ContainerAttachedBasePathKey]
	c.DatabaseMigrationServiceBasePath = DefaultBasePaths[DatabaseMigrationServiceBasePathKey]
	c.DataCatalogBasePath = DefaultBasePaths[DataCatalogBasePathKey]
	c.DataFusionBasePath = DefaultBasePaths[DataFusionBasePathKey]
	c.DataLossPreventionBasePath = DefaultBasePaths[DataLossPreventionBasePathKey]
	c.DataplexBasePath = DefaultBasePaths[DataplexBasePathKey]
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
	c.GKEBackupBasePath = DefaultBasePaths[GKEBackupBasePathKey]
	c.GKEHubBasePath = DefaultBasePaths[GKEHubBasePathKey]
	c.GKEHub2BasePath = DefaultBasePaths[GKEHub2BasePathKey]
	c.HealthcareBasePath = DefaultBasePaths[HealthcareBasePathKey]
	c.IAM2BasePath = DefaultBasePaths[IAM2BasePathKey]
	c.IAMBetaBasePath = DefaultBasePaths[IAMBetaBasePathKey]
	c.IAMWorkforcePoolBasePath = DefaultBasePaths[IAMWorkforcePoolBasePathKey]
	c.IapBasePath = DefaultBasePaths[IapBasePathKey]
	c.IdentityPlatformBasePath = DefaultBasePaths[IdentityPlatformBasePathKey]
	c.KMSBasePath = DefaultBasePaths[KMSBasePathKey]
	c.LoggingBasePath = DefaultBasePaths[LoggingBasePathKey]
	c.LookerBasePath = DefaultBasePaths[LookerBasePathKey]
	c.MemcacheBasePath = DefaultBasePaths[MemcacheBasePathKey]
	c.MLEngineBasePath = DefaultBasePaths[MLEngineBasePathKey]
	c.MonitoringBasePath = DefaultBasePaths[MonitoringBasePathKey]
	c.NetworkManagementBasePath = DefaultBasePaths[NetworkManagementBasePathKey]
	c.NetworkSecurityBasePath = DefaultBasePaths[NetworkSecurityBasePathKey]
	c.NetworkServicesBasePath = DefaultBasePaths[NetworkServicesBasePathKey]
	c.NotebooksBasePath = DefaultBasePaths[NotebooksBasePathKey]
	c.OSConfigBasePath = DefaultBasePaths[OSConfigBasePathKey]
	c.OSLoginBasePath = DefaultBasePaths[OSLoginBasePathKey]
	c.PrivatecaBasePath = DefaultBasePaths[PrivatecaBasePathKey]
	c.PublicCABasePath = DefaultBasePaths[PublicCABasePathKey]
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
	c.TagsLocationBasePath = DefaultBasePaths[TagsLocationBasePathKey]
}

func GetCurrentUserEmail(config *Config, userAgent string) (string, error) {
	// When environment variables UserProjectOverride and BillingProject are set for the provider,
	// the header X-Goog-User-Project is set for the API requests.
	// But it causes an error when calling GetCurrentUserEmail. Set the project to be "NO_BILLING_PROJECT_OVERRIDE".
	// And then it triggers the header X-Goog-User-Project to be set to empty string.

	// See https://github.com/golang/oauth2/issues/306 for a recommendation to do this from a Go maintainer
	// URL retrieved from https://accounts.google.com/.well-known/openid-configuration
	res, err := SendRequest(SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   "NO_BILLING_PROJECT_OVERRIDE",
		RawURL:    "https://openidconnect.googleapis.com/v1/userinfo",
		UserAgent: userAgent,
	})

	if err != nil {
		return "", fmt.Errorf("error retrieving userinfo for your provider credentials. have you enabled the 'https://www.googleapis.com/auth/userinfo.email' scope? error: %s", err)
	}
	if res["email"] == nil {
		return "", fmt.Errorf("error retrieving email from userinfo. email was nil in the response.")
	}
	return res["email"].(string), nil
}

func MultiEnvSearch(ks []string) string {
	for _, k := range ks {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return ""
}

// MultiEnvDefault is a helper function that returns the value of the first
// environment variable in the given list that returns a non-empty value. If
// none of the environment variables return a value, the default value is
// returned.
func MultiEnvDefault(ks []string, dv interface{}) interface{} {
	for _, k := range ks {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return dv
}

func CustomEndpointValidator() validator.String {
	return stringvalidator.RegexMatches(regexp.MustCompile(`.*/[^/]+/$`), "")
}

// return the region a selfLink is referring to
func GetRegionFromRegionSelfLink(selfLink string) string {
	re := regexp.MustCompile("/compute/[a-zA-Z0-9]*/projects/[a-zA-Z0-9-]*/regions/([a-zA-Z0-9-]*)")
	switch {
	case re.MatchString(selfLink):
		if res := re.FindStringSubmatch(selfLink); len(res) == 2 && res[1] != "" {
			return res[1]
		}
	}
	return selfLink
}
