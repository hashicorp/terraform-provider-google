package google

import (
	"context"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/metaschema"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/oauth2"

	"google.golang.org/api/option"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.ProviderWithMetaSchema = &frameworkProvider{}

	defaultClientScopes = []string{
		"https://www.googleapis.com/auth/cloud-platform",
		"https://www.googleapis.com/auth/userinfo.email",
	}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) provider.ProviderWithMetaSchema {
	return &frameworkProvider{
		version: version,
	}
}

// frameworkProvider is the provider implementation.
type frameworkProvider struct {
	billingProject             types.String
	client                     *http.Client
	context                    context.Context
	gRPCLoggingOptions         []option.ClientOption
	pollInterval               time.Duration
	project                    types.String
	region                     types.String
	zone                       types.String
	requestBatcherIam          *RequestBatcher
	requestBatcherServiceUsage *RequestBatcher
	scopes                     []string
	tokenSource                oauth2.TokenSource
	userAgent                  string
	userProjectOverride        bool
	version                    string

	// paths for client setup
	AccessApprovalBasePath       string
	AccessContextManagerBasePath string
	ActiveDirectoryBasePath      string
	AlloydbBasePath              string
	ApigeeBasePath               string
	AppEngineBasePath            string
	ArtifactRegistryBasePath     string
	BeyondcorpBasePath           string
	BigQueryBasePath             string
	BigqueryAnalyticsHubBasePath string
	BigqueryConnectionBasePath   string
	BigqueryDatapolicyBasePath   string
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
	ContainerAttachedBasePath    string
	DataCatalogBasePath          string
	DataFusionBasePath           string
	DataLossPreventionBasePath   string
	DataplexBasePath             string
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
	GKEBackupBasePath            string
	GKEHubBasePath               string
	HealthcareBasePath           string
	IAM2BasePath                 string
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
}

// Metadata returns the provider type name.
func (p *frameworkProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "google"
	resp.Version = p.version
}

// MetaSchema returns the provider meta schema.
func (p *frameworkProvider) MetaSchema(_ context.Context, _ provider.MetaSchemaRequest, resp *provider.MetaSchemaResponse) {
	resp.Schema = metaschema.Schema{
		Attributes: map[string]metaschema.Attribute{
			"module_name": metaschema.StringAttribute{
				Optional: true,
			},
		},
	}
}

// Schema defines the provider-level schema for configuration data.
func (p *frameworkProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"credentials": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("access_token"),
					}...),
					CredentialsValidator(),
				},
			},
			"access_token": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("credentials"),
					}...),
				},
			},
			"impersonate_service_account": schema.StringAttribute{
				Optional: true,
			},
			"impersonate_service_account_delegates": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"project": schema.StringAttribute{
				Optional: true,
			},
			"billing_project": schema.StringAttribute{
				Optional: true,
			},
			"region": schema.StringAttribute{
				Optional: true,
			},
			"zone": schema.StringAttribute{
				Optional: true,
			},
			"scopes": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"user_project_override": schema.BoolAttribute{
				Optional: true,
			},
			"request_timeout": schema.StringAttribute{
				Optional: true,
			},
			"request_reason": schema.StringAttribute{
				Optional: true,
			},

			// Generated Products
			"access_approval_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"access_context_manager_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"active_directory_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"alloydb_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"apigee_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"app_engine_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"artifact_registry_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"beyondcorp_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"big_query_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"bigquery_analytics_hub_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"bigquery_connection_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"bigquery_datapolicy_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"bigquery_data_transfer_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"bigquery_reservation_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"bigtable_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"billing_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"binary_authorization_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"certificate_manager_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"cloud_asset_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"cloud_build_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"cloud_functions_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"cloudfunctions2_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"cloud_identity_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"cloud_ids_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"cloud_iot_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"cloud_run_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"cloud_run_v2_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"cloud_scheduler_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"cloud_tasks_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"compute_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"container_analysis_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"container_attached_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"data_catalog_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"data_fusion_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"data_loss_prevention_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"dataplex_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"dataproc_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"dataproc_metastore_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"datastore_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"datastream_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"deployment_manager_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"dialogflow_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"dialogflow_cx_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"dns_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"document_ai_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"essential_contacts_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"filestore_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"firestore_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"game_services_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"gke_backup_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"gke_hub_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"healthcare_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"iam2_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"iam_beta_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"iam_workforce_pool_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"iap_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"identity_platform_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"kms_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"logging_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"memcache_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"ml_engine_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"monitoring_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"network_management_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"network_services_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"notebooks_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"os_config_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"os_login_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"privateca_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"pubsub_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"pubsub_lite_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"redis_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"resource_manager_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"secret_manager_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"security_center_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"service_management_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"service_usage_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"source_repo_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"spanner_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"sql_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"storage_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"storage_transfer_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"tags_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"tpu_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"vertex_ai_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"vpc_access_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"workflows_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},

			// Handwritten Products / Versioned / Atypical Entries
			"cloud_billing_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"composer_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"container_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"dataflow_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"iam_credentials_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"resource_manager_v3_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"iam_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"service_networking_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"tags_location_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},

			// dcl
			"container_aws_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
			"container_azure_custom_endpoint": &schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					CustomEndpointValidator(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"batching": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"send_after": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								NonNegativeDurationValidator(),
							},
						},
						"enable_batching": schema.BoolAttribute{
							Optional: true,
						},
					},
				},
			},
		},
	}

	configureDCLCustomEndpointAttributesFramework(&resp.Schema)
}

// Configure prepares an API client for data sources and resources.
func (p *frameworkProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data ProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	p.LoadAndValidateFramework(ctx, data, req.TerraformVersion, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Example client configuration for data sources and resources
	resp.DataSourceData = p
	resp.ResourceData = p
}

// DataSources defines the data sources implemented in the provider.
func (p *frameworkProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewGoogleClientConfigDataSource,
		NewGoogleClientOpenIDUserinfoDataSource,
		NewGoogleDnsManagedZoneDataSource,
		NewGoogleDnsRecordSetDataSource,
		NewGoogleDnsKeysDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *frameworkProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}
