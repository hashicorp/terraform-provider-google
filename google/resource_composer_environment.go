package google

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"google.golang.org/api/composer/v1"
)

const (
	composerEnvironmentEnvVariablesRegexp          = "[a-zA-Z_][a-zA-Z0-9_]*."
	composerEnvironmentReservedAirflowEnvVarRegexp = "AIRFLOW__[A-Z0-9_]+__[A-Z0-9_]+"
	composerEnvironmentVersionRegexp               = `composer-(([0-9]+)(\.[0-9]+\.[0-9]+(-preview\.[0-9]+)?)?|latest)-airflow-(([0-9]+)((\.[0-9]+)(\.[0-9]+)?)?)`
)

var composerEnvironmentReservedEnvVar = map[string]struct{}{
	"AIRFLOW_HOME":     {},
	"C_FORCE_ROOT":     {},
	"CONTAINER_NAME":   {},
	"DAGS_FOLDER":      {},
	"GCP_PROJECT":      {},
	"GCS_BUCKET":       {},
	"GKE_CLUSTER_NAME": {},
	"SQL_DATABASE":     {},
	"SQL_INSTANCE":     {},
	"SQL_PASSWORD":     {},
	"SQL_PROJECT":      {},
	"SQL_REGION":       {},
	"SQL_USER":         {},
}

var (
	composerSoftwareConfigKeys = []string{
		"config.0.software_config.0.airflow_config_overrides",
		"config.0.software_config.0.pypi_packages",
		"config.0.software_config.0.env_variables",
		"config.0.software_config.0.image_version",
		"config.0.software_config.0.python_version",
		"config.0.software_config.0.scheduler_count",
	}

	composerConfigKeys = []string{
		"config.0.node_count",
		"config.0.node_config",
		"config.0.software_config",
		"config.0.private_environment_config",
		"config.0.web_server_network_access_control",
		"config.0.database_config",
		"config.0.web_server_config",
		"config.0.encryption_config",
		"config.0.maintenance_window",
		"config.0.workloads_config",
		"config.0.environment_size",
		"config.0.master_authorized_networks_config",
	}

	composerPrivateEnvironmentConfig = []string{
		"config.0.private_environment_config.0.enable_private_endpoint",
		"config.0.private_environment_config.0.master_ipv4_cidr_block",
		"config.0.private_environment_config.0.cloud_sql_ipv4_cidr_block",
		"config.0.private_environment_config.0.web_server_ipv4_cidr_block",
		"config.0.private_environment_config.0.cloud_composer_network_ipv4_cidr_block",
		"config.0.private_environment_config.0.enable_privately_used_public_ips",
		"config.0.private_environment_config.0.cloud_composer_connection_subnetwork",
	}

	composerIpAllocationPolicyKeys = []string{
		"config.0.node_config.0.ip_allocation_policy.0.use_ip_aliases",
		"config.0.node_config.0.ip_allocation_policy.0.cluster_secondary_range_name",
		"config.0.node_config.0.ip_allocation_policy.0.services_secondary_range_name",
		"config.0.node_config.0.ip_allocation_policy.0.cluster_ipv4_cidr_block",
		"config.0.node_config.0.ip_allocation_policy.0.services_ipv4_cidr_block",
	}

	allowedIpRangesConfig = &schema.Resource{
		Schema: map[string]*schema.Schema{
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `IP address or range, defined using CIDR notation, of requests that this rule applies to. Examples: 192.168.1.1 or 192.168.0.0/16 or 2001:db8::/32 or 2001:0db8:0000:0042:0000:8a2e:0370:7334. IP range prefixes should be properly truncated. For example, 1.2.3.4/24 should be truncated to 1.2.3.0/24. Similarly, for IPv6, 2001:db8::1/32 should be truncated to 2001:db8::/32.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `A description of this ip range.`,
			},
		},
	}

	cidrBlocks = &schema.Resource{
		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `display_name is a field for users to identify CIDR blocks.`,
			},
			"cidr_block": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `cidr_block must be specified in CIDR notation.`,
			},
		},
	}
)

func resourceComposerEnvironment() *schema.Resource {
	return &schema.Resource{
		Create: resourceComposerEnvironmentCreate,
		Read:   resourceComposerEnvironmentRead,
		Update: resourceComposerEnvironmentUpdate,
		Delete: resourceComposerEnvironmentDelete,

		Importer: &schema.ResourceImporter{
			State: resourceComposerEnvironmentImport,
		},

		Timeouts: &schema.ResourceTimeout{
			// Composer takes <= 1 hr for create/update.
			Create: schema.DefaultTimeout(120 * time.Minute),
			Update: schema.DefaultTimeout(120 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateGCEName,
				Description:  `Name of the environment.`,
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: `The location or Compute Engine region for the environment.`,
			},
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The ID of the project in which the resource belongs. If it is not provided, the provider project is used.`,
			},
			"config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: `Configuration parameters for this environment.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_count": {
							Type:         schema.TypeInt,
							Computed:     true,
							Optional:     true,
							AtLeastOneOf: composerConfigKeys,
							ValidateFunc: validation.IntAtLeast(3),
							Description:  `The number of nodes in the Kubernetes Engine cluster that will be used to run this environment. This field is supported for Cloud Composer environments in versions composer-1.*.*-airflow-*.*.*.`,
						},
						"node_config": {
							Type:         schema.TypeList,
							Computed:     true,
							Optional:     true,
							AtLeastOneOf: composerConfigKeys,
							MaxItems:     1,
							Description:  `The configuration used for the Kubernetes Engine cluster.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"zone": {
										Type:             schema.TypeString,
										Optional:         true,
										Computed:         true,
										ForceNew:         true,
										DiffSuppressFunc: compareSelfLinkOrResourceName,
										Description:      `The Compute Engine zone in which to deploy the VMs running the Apache Airflow software, specified as the zone name or relative resource name (e.g. "projects/{project}/zones/{zone}"). Must belong to the enclosing environment's project and region. This field is supported for Cloud Composer environments in versions composer-1.*.*-airflow-*.*.*.`,
									},
									"machine_type": {
										Type:             schema.TypeString,
										Computed:         true,
										Optional:         true,
										ForceNew:         true,
										DiffSuppressFunc: compareSelfLinkOrResourceName,
										Description:      `The Compute Engine machine type used for cluster instances, specified as a name or relative resource name. For example: "projects/{project}/zones/{zone}/machineTypes/{machineType}". Must belong to the enclosing environment's project and region/zone. This field is supported for Cloud Composer environments in versions composer-1.*.*-airflow-*.*.*.`,
									},
									"network": {
										Type:             schema.TypeString,
										Computed:         true,
										Optional:         true,
										ForceNew:         true,
										DiffSuppressFunc: compareSelfLinkOrResourceName,
										Description:      `The Compute Engine machine type used for cluster instances, specified as a name or relative resource name. For example: "projects/{project}/zones/{zone}/machineTypes/{machineType}". Must belong to the enclosing environment's project and region/zone. The network must belong to the environment's project. If unspecified, the "default" network ID in the environment's project is used. If a Custom Subnet Network is provided, subnetwork must also be provided.`,
									},
									"subnetwork": {
										Type:             schema.TypeString,
										Optional:         true,
										ForceNew:         true,
										DiffSuppressFunc: compareSelfLinkOrResourceName,
										Description:      `The Compute Engine subnetwork to be used for machine communications, , specified as a self-link, relative resource name (e.g. "projects/{project}/regions/{region}/subnetworks/{subnetwork}"), or by name. If subnetwork is provided, network must also be provided and the subnetwork must belong to the enclosing environment's project and region.`,
									},
									"disk_size_gb": {
										Type:        schema.TypeInt,
										Computed:    true,
										Optional:    true,
										ForceNew:    true,
										Description: `The disk size in GB used for node VMs. Minimum size is 20GB. If unspecified, defaults to 100GB. Cannot be updated. This field is supported for Cloud Composer environments in versions composer-1.*.*-airflow-*.*.*.`,
									},
									"oauth_scopes": {
										Type:     schema.TypeSet,
										Computed: true,
										Optional: true,
										ForceNew: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Set:         schema.HashString,
										Description: `The set of Google API scopes to be made available on all node VMs. Cannot be updated. If empty, defaults to ["https://www.googleapis.com/auth/cloud-platform"]. This field is supported for Cloud Composer environments in versions composer-1.*.*-airflow-*.*.*.`,
									},
									"service_account": {
										Type:             schema.TypeString,
										Computed:         true,
										Optional:         true,
										ForceNew:         true,
										ValidateFunc:     validateServiceAccountRelativeNameOrEmail,
										DiffSuppressFunc: compareServiceAccountEmailToLink,
										Description:      `The Google Cloud Platform Service Account to be used by the node VMs. If a service account is not specified, the "default" Compute Engine service account is used. Cannot be updated. If given, note that the service account must have roles/composer.worker for any GCP resources created under the Cloud Composer Environment.`,
									},
									"enable_ip_masq_agent": {
										Type:        schema.TypeBool,
										Computed:    true,
										Optional:    true,
										ForceNew:    true,
										Description: `Deploys 'ip-masq-agent' daemon set in the GKE cluster and defines nonMasqueradeCIDRs equals to pod IP range so IP masquerading is used for all destination addresses, except between pods traffic. See: https://cloud.google.com/kubernetes-engine/docs/how-to/ip-masquerade-agent`,
									},

									"tags": {
										Type:     schema.TypeSet,
										Optional: true,
										ForceNew: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Set:         schema.HashString,
										Description: `The list of instance tags applied to all node VMs. Tags are used to identify valid sources or targets for network firewalls. Each tag within the list must comply with RFC1035. Cannot be updated. This field is supported for Cloud Composer environments in versions composer-1.*.*-airflow-*.*.*.`,
									},
									"ip_allocation_policy": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										ForceNew:    true,
										ConfigMode:  schema.SchemaConfigModeAttr,
										MaxItems:    1,
										Description: `Configuration for controlling how IPs are allocated in the GKE cluster. Cannot be updated.`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"use_ip_aliases": {
													Type:         schema.TypeBool,
													Optional:     true,
													ForceNew:     true,
													AtLeastOneOf: composerIpAllocationPolicyKeys,
													Description:  `Whether or not to enable Alias IPs in the GKE cluster. If true, a VPC-native cluster is created. Defaults to true if the ip_allocation_policy block is present in config. This field is only supported for Cloud Composer environments in versions composer-1.*.*-airflow-*.*.*. Environments in newer versions always use VPC-native GKE clusters.`,
												},
												"cluster_secondary_range_name": {
													Type:          schema.TypeString,
													Optional:      true,
													ForceNew:      true,
													AtLeastOneOf:  composerIpAllocationPolicyKeys,
													Description:   `The name of the cluster's secondary range used to allocate IP addresses to pods. Specify either cluster_secondary_range_name or cluster_ipv4_cidr_block but not both. For Cloud Composer environments in versions composer-1.*.*-airflow-*.*.*, this field is applicable only when use_ip_aliases is true.`,
													ConflictsWith: []string{"config.0.node_config.0.ip_allocation_policy.0.cluster_ipv4_cidr_block"},
												},
												"services_secondary_range_name": {
													Type:          schema.TypeString,
													Optional:      true,
													ForceNew:      true,
													AtLeastOneOf:  composerIpAllocationPolicyKeys,
													Description:   `The name of the services' secondary range used to allocate IP addresses to the cluster. Specify either services_secondary_range_name or services_ipv4_cidr_block but not both. For Cloud Composer environments in versions composer-1.*.*-airflow-*.*.*, this field is applicable only when use_ip_aliases is true.`,
													ConflictsWith: []string{"config.0.node_config.0.ip_allocation_policy.0.services_ipv4_cidr_block"},
												},
												"cluster_ipv4_cidr_block": {
													Type:             schema.TypeString,
													Optional:         true,
													ForceNew:         true,
													AtLeastOneOf:     composerIpAllocationPolicyKeys,
													Description:      `The IP address range used to allocate IP addresses to pods in the cluster. For Cloud Composer environments in versions composer-1.*.*-airflow-*.*.*, this field is applicable only when use_ip_aliases is true. Set to blank to have GKE choose a range with the default size. Set to /netmask (e.g. /14) to have GKE choose a range with a specific netmask. Set to a CIDR notation (e.g. 10.96.0.0/14) from the RFC-1918 private networks (e.g. 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16) to pick a specific range to use. Specify either cluster_secondary_range_name or cluster_ipv4_cidr_block but not both.`,
													DiffSuppressFunc: cidrOrSizeDiffSuppress,
													ConflictsWith:    []string{"config.0.node_config.0.ip_allocation_policy.0.cluster_secondary_range_name"},
												},
												"services_ipv4_cidr_block": {
													Type:             schema.TypeString,
													Optional:         true,
													ForceNew:         true,
													AtLeastOneOf:     composerIpAllocationPolicyKeys,
													Description:      `The IP address range used to allocate IP addresses in this cluster. For Cloud Composer environments in versions composer-1.*.*-airflow-*.*.*, this field is applicable only when use_ip_aliases is true. Set to blank to have GKE choose a range with the default size. Set to /netmask (e.g. /14) to have GKE choose a range with a specific netmask. Set to a CIDR notation (e.g. 10.96.0.0/14) from the RFC-1918 private networks (e.g. 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16) to pick a specific range to use. Specify either services_secondary_range_name or services_ipv4_cidr_block but not both.`,
													DiffSuppressFunc: cidrOrSizeDiffSuppress,
													ConflictsWith:    []string{"config.0.node_config.0.ip_allocation_policy.0.services_secondary_range_name"},
												},
											},
										},
									},
								},
							},
						},
						"software_config": {
							Type:         schema.TypeList,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: composerConfigKeys,
							MaxItems:     1,
							Description:  `The configuration settings for software inside the environment.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"airflow_config_overrides": {
										Type:         schema.TypeMap,
										Optional:     true,
										AtLeastOneOf: composerSoftwareConfigKeys,
										Elem:         &schema.Schema{Type: schema.TypeString},
										Description:  `Apache Airflow configuration properties to override. Property keys contain the section and property names, separated by a hyphen, for example "core-dags_are_paused_at_creation". Section names must not contain hyphens ("-"), opening square brackets ("["), or closing square brackets ("]"). The property name must not be empty and cannot contain "=" or ";". Section and property names cannot contain characters: "." Apache Airflow configuration property names must be written in snake_case. Property values can contain any character, and can be written in any lower/upper case format. Certain Apache Airflow configuration property values are blacklisted, and cannot be overridden.`,
									},
									"pypi_packages": {
										Type:         schema.TypeMap,
										Optional:     true,
										AtLeastOneOf: composerSoftwareConfigKeys,
										Elem:         &schema.Schema{Type: schema.TypeString},
										ValidateFunc: validateComposerEnvironmentPypiPackages,
										Description:  `Custom Python Package Index (PyPI) packages to be installed in the environment. Keys refer to the lowercase package name (e.g. "numpy"). Values are the lowercase extras and version specifier (e.g. "==1.12.0", "[devel,gcp_api]", "[devel]>=1.8.2, <1.9.2"). To specify a package without pinning it to a version specifier, use the empty string as the value.`,
									},
									"env_variables": {
										Type:         schema.TypeMap,
										Optional:     true,
										AtLeastOneOf: composerSoftwareConfigKeys,
										Elem:         &schema.Schema{Type: schema.TypeString},
										ValidateFunc: validateComposerEnvironmentEnvVariables,
										Description:  `Additional environment variables to provide to the Apache Airflow scheduler, worker, and webserver processes. Environment variable names must match the regular expression [a-zA-Z_][a-zA-Z0-9_]*. They cannot specify Apache Airflow software configuration overrides (they cannot match the regular expression AIRFLOW__[A-Z0-9_]+__[A-Z0-9_]+), and they cannot match any of the following reserved names: AIRFLOW_HOME C_FORCE_ROOT CONTAINER_NAME DAGS_FOLDER GCP_PROJECT GCS_BUCKET GKE_CLUSTER_NAME SQL_DATABASE SQL_INSTANCE SQL_PASSWORD SQL_PROJECT SQL_REGION SQL_USER.`,
									},
									"image_version": {
										Type:             schema.TypeString,
										Computed:         true,
										Optional:         true,
										ForceNew:         true,
										AtLeastOneOf:     composerSoftwareConfigKeys,
										ValidateFunc:     validateRegexp(composerEnvironmentVersionRegexp),
										DiffSuppressFunc: composerImageVersionDiffSuppress,
										Description:      `The version of the software running in the environment. This encapsulates both the version of Cloud Composer functionality and the version of Apache Airflow. It must match the regular expression composer-([0-9]+(\.[0-9]+\.[0-9]+(-preview\.[0-9]+)?)?|latest)-airflow-([0-9]+(\.[0-9]+(\.[0-9]+)?)?). The Cloud Composer portion of the image version is a full semantic version, or an alias in the form of major version number or 'latest'. The Apache Airflow portion of the image version is a full semantic version that points to one of the supported Apache Airflow versions, or an alias in the form of only major or major.minor versions specified. See documentation for more details and version list.`,
									},
									"python_version": {
										Type:         schema.TypeString,
										Optional:     true,
										AtLeastOneOf: composerSoftwareConfigKeys,
										Computed:     true,
										ForceNew:     true,
										Description:  `The major version of Python used to run the Apache Airflow scheduler, worker, and webserver processes. Can be set to '2' or '3'. If not specified, the default is '2'. Cannot be updated. This field is supported for Cloud Composer environments in versions composer-1.*.*-airflow-*.*.*. Environments in newer versions always use Python major version 3.`,
									},
									"scheduler_count": {
										Type:         schema.TypeInt,
										Optional:     true,
										AtLeastOneOf: composerSoftwareConfigKeys,
										Computed:     true,
										Description:  `The number of schedulers for Airflow. This field is supported for Cloud Composer environments in versions composer-1.*.*-airflow-2.*.*.`,
									},
								},
							},
						},
						"private_environment_config": {
							Type:         schema.TypeList,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: composerConfigKeys,
							MaxItems:     1,
							ForceNew:     true,
							Description:  `The configuration used for the Private IP Cloud Composer environment.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enable_private_endpoint": {
										Type:         schema.TypeBool,
										Optional:     true,
										Default:      true,
										AtLeastOneOf: composerPrivateEnvironmentConfig,
										ForceNew:     true,
										Description:  `If true, access to the public endpoint of the GKE cluster is denied. If this field is set to true, ip_allocation_policy.use_ip_aliases must be set to true for Cloud Composer environments in versions composer-1.*.*-airflow-*.*.*.`,
									},
									"master_ipv4_cidr_block": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										AtLeastOneOf: composerPrivateEnvironmentConfig,
										ForceNew:     true,
										Description:  `The IP range in CIDR notation to use for the hosted master network. This range is used for assigning internal IP addresses to the cluster master or set of masters and to the internal load balancer virtual IP. This range must not overlap with any other ranges in use within the cluster's network. If left blank, the default value of '172.16.0.0/28' is used.`,
									},
									"web_server_ipv4_cidr_block": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										AtLeastOneOf: composerPrivateEnvironmentConfig,
										ForceNew:     true,
										Description:  `The CIDR block from which IP range for web server will be reserved. Needs to be disjoint from master_ipv4_cidr_block and cloud_sql_ipv4_cidr_block. This field is supported for Cloud Composer environments in versions composer-1.*.*-airflow-*.*.*.`,
									},
									"cloud_sql_ipv4_cidr_block": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										AtLeastOneOf: composerPrivateEnvironmentConfig,
										ForceNew:     true,
										Description:  `The CIDR block from which IP range in tenant project will be reserved for Cloud SQL. Needs to be disjoint from web_server_ipv4_cidr_block.`,
									},
									"cloud_composer_network_ipv4_cidr_block": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										AtLeastOneOf: composerPrivateEnvironmentConfig,
										ForceNew:     true,
										Description:  `The CIDR block from which IP range for Cloud Composer Network in tenant project will be reserved. Needs to be disjoint from private_cluster_config.master_ipv4_cidr_block and cloud_sql_ipv4_cidr_block. This field is supported for Cloud Composer environments in versions composer-2.*.*-airflow-*.*.* and newer.`,
									},
									"enable_privately_used_public_ips": {
										Type:         schema.TypeBool,
										Optional:     true,
										Computed:     true,
										AtLeastOneOf: composerPrivateEnvironmentConfig,
										ForceNew:     true,
										Description:  `When enabled, IPs from public (non-RFC1918) ranges can be used for ip_allocation_policy.cluster_ipv4_cidr_block and ip_allocation_policy.service_ipv4_cidr_block.`,
									},
									"cloud_composer_connection_subnetwork": {
										Type:             schema.TypeString,
										Optional:         true,
										Computed:         true,
										AtLeastOneOf:     composerPrivateEnvironmentConfig,
										ForceNew:         true,
										DiffSuppressFunc: compareSelfLinkRelativePaths,
										Description:      `When specified, the environment will use Private Service Connect instead of VPC peerings to connect to Cloud SQL in the Tenant Project, and the PSC endpoint in the Customer Project will use an IP address from this subnetwork. This field is supported for Cloud Composer environments in versions composer-2.*.*-airflow-*.*.* and newer.`,
									},
								},
							},
						},
						"web_server_network_access_control": {
							Type:         schema.TypeList,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: composerConfigKeys,
							MaxItems:     1,
							Description:  `The network-level access control policy for the Airflow web server. If unspecified, no network-level access restrictions will be applied. This field is supported for Cloud Composer environments in versions composer-1.*.*-airflow-*.*.*.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allowed_ip_range": {
										Type:        schema.TypeSet,
										Computed:    true,
										Optional:    true,
										Elem:        allowedIpRangesConfig,
										Description: `A collection of allowed IP ranges with descriptions.`,
									},
								},
							},
						},
						"database_config": {
							Type:         schema.TypeList,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: composerConfigKeys,
							MaxItems:     1,
							Description:  `The configuration of Cloud SQL instance that is used by the Apache Airflow software. This field is supported for Cloud Composer environments in versions composer-1.*.*-airflow-*.*.*.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"machine_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `Optional. Cloud SQL machine type used by Airflow database. It has to be one of: db-n1-standard-2, db-n1-standard-4, db-n1-standard-8 or db-n1-standard-16. If not specified, db-n1-standard-2 will be used.`,
									},
								},
							},
						},
						"web_server_config": {
							Type:         schema.TypeList,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: composerConfigKeys,
							MaxItems:     1,
							Description:  `The configuration settings for the Airflow web server App Engine instance. This field is supported for Cloud Composer environments in versions composer-1.*.*-airflow-*.*.*.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"machine_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `Optional. Machine type on which Airflow web server is running. It has to be one of: composer-n1-webserver-2, composer-n1-webserver-4 or composer-n1-webserver-8. If not specified, composer-n1-webserver-2 will be used. Value custom is returned only in response, if Airflow web server parameters were manually changed to a non-standard values.`,
									},
								},
							},
						},
						"encryption_config": {
							Type:         schema.TypeList,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: composerConfigKeys,
							MaxItems:     1,
							Description:  `The encryption options for the Composer environment and its dependencies.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kms_key_name": {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: `Optional. Customer-managed Encryption Key available through Google's Key Management Service. Cannot be updated.`,
									},
								},
							},
						},
						"maintenance_window": {
							Type:         schema.TypeList,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: composerConfigKeys,
							MaxItems:     1,
							Description:  `The configuration for Cloud Composer maintenance window.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_time": {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    false,
										Description: `Start time of the first recurrence of the maintenance window.`,
									},
									"end_time": {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    false,
										Description: `Maintenance window end time. It is used only to calculate the duration of the maintenance window. The value for end-time must be in the future, relative to 'start_time'.`,
									},
									"recurrence": {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    false,
										Description: `Maintenance window recurrence. Format is a subset of RFC-5545 (https://tools.ietf.org/html/rfc5545) 'RRULE'. The only allowed values for 'FREQ' field are 'FREQ=DAILY' and 'FREQ=WEEKLY;BYDAY=...'. Example values: 'FREQ=WEEKLY;BYDAY=TU,WE', 'FREQ=DAILY'.`,
									},
								},
							},
						},
						"workloads_config": {
							Type:         schema.TypeList,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: composerConfigKeys,
							MaxItems:     1,
							Description:  `The workloads configuration settings for the GKE cluster associated with the Cloud Composer environment. Supported for Cloud Composer environments in versions composer-2.*.*-airflow-*.*.* and newer.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"scheduler": {
										Type:        schema.TypeList,
										Optional:    true,
										ForceNew:    false,
										Description: `Configuration for resources used by Airflow schedulers.`,
										MaxItems:    1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"cpu": {
													Type:         schema.TypeFloat,
													Optional:     true,
													ForceNew:     false,
													ValidateFunc: validation.FloatAtLeast(0),
													Description:  `CPU request and limit for a single Airflow scheduler replica`,
												},
												"memory_gb": {
													Type:         schema.TypeFloat,
													Optional:     true,
													ForceNew:     false,
													ValidateFunc: validation.FloatAtLeast(0),
													Description:  `Memory (GB) request and limit for a single Airflow scheduler replica.`,
												},
												"storage_gb": {
													Type:         schema.TypeFloat,
													Optional:     true,
													ForceNew:     false,
													ValidateFunc: validation.FloatAtLeast(0),
													Description:  `Storage (GB) request and limit for a single Airflow scheduler replica.`,
												},
												"count": {
													Type:         schema.TypeInt,
													Optional:     true,
													ForceNew:     false,
													ValidateFunc: validation.IntAtLeast(0),
													Description:  `The number of schedulers.`,
												},
											},
										},
									},
									"web_server": {
										Type:        schema.TypeList,
										Optional:    true,
										ForceNew:    false,
										Description: `Configuration for resources used by Airflow web server.`,
										MaxItems:    1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"cpu": {
													Type:         schema.TypeFloat,
													Optional:     true,
													ForceNew:     false,
													ValidateFunc: validation.FloatAtLeast(0),
													Description:  `CPU request and limit for Airflow web server.`,
												},
												"memory_gb": {
													Type:         schema.TypeFloat,
													Optional:     true,
													ForceNew:     false,
													ValidateFunc: validation.FloatAtLeast(0),
													Description:  `Memory (GB) request and limit for Airflow web server.`,
												},
												"storage_gb": {
													Type:         schema.TypeFloat,
													Optional:     true,
													ForceNew:     false,
													ValidateFunc: validation.FloatAtLeast(0),
													Description:  `Storage (GB) request and limit for Airflow web server.`,
												},
											},
										},
									},
									"worker": {
										Type:        schema.TypeList,
										Optional:    true,
										ForceNew:    false,
										Description: `Configuration for resources used by Airflow workers.`,
										MaxItems:    1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"cpu": {
													Type:         schema.TypeFloat,
													Optional:     true,
													ForceNew:     false,
													ValidateFunc: validation.FloatAtLeast(0),
													Description:  `CPU request and limit for a single Airflow worker replica.`,
												},
												"memory_gb": {
													Type:         schema.TypeFloat,
													Optional:     true,
													ForceNew:     false,
													ValidateFunc: validation.FloatAtLeast(0),
													Description:  `Memory (GB) request and limit for a single Airflow worker replica.`,
												},
												"storage_gb": {
													Type:         schema.TypeFloat,
													Optional:     true,
													ForceNew:     false,
													ValidateFunc: validation.FloatAtLeast(0),
													Description:  `Storage (GB) request and limit for a single Airflow worker replica.`,
												},
												"min_count": {
													Type:         schema.TypeInt,
													Optional:     true,
													ForceNew:     false,
													ValidateFunc: validation.IntAtLeast(0),
													Description:  `Minimum number of workers for autoscaling.`,
												},
												"max_count": {
													Type:         schema.TypeInt,
													Optional:     true,
													ForceNew:     false,
													ValidateFunc: validation.IntAtLeast(0),
													Description:  `Maximum number of workers for autoscaling.`,
												},
											},
										},
									},
								},
							},
						},
						"environment_size": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     false,
							AtLeastOneOf: composerConfigKeys,
							ValidateFunc: validation.StringInSlice([]string{"ENVIRONMENT_SIZE_SMALL", "ENVIRONMENT_SIZE_MEDIUM", "ENVIRONMENT_SIZE_LARGE"}, false),
							Description:  `The size of the Cloud Composer environment. This field is supported for Cloud Composer environments in versions composer-2.*.*-airflow-*.*.* and newer.`,
						},
						"master_authorized_networks_config": {
							Type:         schema.TypeList,
							Optional:     true,
							AtLeastOneOf: composerConfigKeys,
							MaxItems:     1,
							Description:  `Configuration options for the master authorized networks feature. Enabled master authorized networks will disallow all external traffic to access Kubernetes master through HTTPS except traffic from the given CIDR blocks, Google Compute Engine Public IPs and Google Prod IPs.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: `Whether or not master authorized networks is enabled.`,
									},
									"cidr_blocks": {
										Type:        schema.TypeSet,
										Optional:    true,
										Elem:        cidrBlocks,
										Description: `cidr_blocks define up to 50 external networks that could access Kubernetes master through HTTPS.`,
									},
								},
							},
						},
						"airflow_uri": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The URI of the Apache Airflow Web UI hosted within this environment.`,
						},
						"dag_gcs_prefix": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The Cloud Storage prefix of the DAGs for this environment. Although Cloud Storage objects reside in a flat namespace, a hierarchical file tree can be simulated using '/'-delimited object name prefixes. DAG objects for this environment reside in a simulated directory with this prefix.`,
						},
						"gke_cluster": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The Kubernetes Engine cluster used to run this environment.`,
						},
					},
				},
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `User-defined labels for this environment. The labels map can contain no more than 64 entries. Entries of the labels map are UTF8 strings that comply with the following restrictions: Label keys must be between 1 and 63 characters long and must conform to the following regular expression: [a-z]([-a-z0-9]*[a-z0-9])?. Label values must be between 0 and 63 characters long and must conform to the regular expression ([a-z]([-a-z0-9]*[a-z0-9])?)?. No more than 64 labels can be associated with a given environment. Both keys and values must be <= 128 bytes in size.`,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceComposerEnvironmentCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	envName, err := resourceComposerEnvironmentName(d, config)
	if err != nil {
		return err
	}

	transformedConfig, err := expandComposerEnvironmentConfig(d.Get("config"), d, config)
	if err != nil {
		return err
	}

	env := &composer.Environment{
		Name:   envName.resourceName(),
		Labels: expandLabels(d),
		Config: transformedConfig,
	}

	// Some fields cannot be specified during create and must be updated post-creation.
	updateOnlyEnv := getComposerEnvironmentPostCreateUpdateObj(env)

	log.Printf("[DEBUG] Creating new Environment %q", envName.parentName())
	op, err := config.NewComposerClient(userAgent).Projects.Locations.Environments.Create(envName.parentName(), env).Do()
	if err != nil {
		return err
	}

	// Store the ID now
	id, err := replaceVars(d, config, "projects/{{project}}/locations/{{region}}/environments/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	waitErr := composerOperationWaitTime(
		config, op, envName.Project, "Creating Environment", userAgent,
		d.Timeout(schema.TimeoutCreate))

	if waitErr != nil {
		// The resource didn't actually get created, remove from state.
		d.SetId("")

		errMsg := fmt.Sprintf("Error waiting to create Environment: %s", waitErr)
		if err := handleComposerEnvironmentCreationOpFailure(id, envName, d, config); err != nil {
			return fmt.Errorf("Error waiting to create Environment: %s. An initial "+
				"environment was or is still being created, and clean up failed with "+
				"error: %s.", errMsg, err)
		}

		return fmt.Errorf("Error waiting to create Environment: %s", waitErr)
	}

	log.Printf("[DEBUG] Finished creating Environment %q: %#v", d.Id(), op)

	if err := resourceComposerEnvironmentPostCreateUpdate(updateOnlyEnv, d, config, userAgent); err != nil {
		return err
	}

	return resourceComposerEnvironmentRead(d, meta)
}

func resourceComposerEnvironmentRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	envName, err := resourceComposerEnvironmentName(d, config)
	if err != nil {
		return err
	}

	res, err := config.NewComposerClient(userAgent).Projects.Locations.Environments.Get(envName.resourceName()).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("ComposerEnvironment %q", d.Id()))
	}

	// Set from getProject(d)
	if err := d.Set("project", envName.Project); err != nil {
		return fmt.Errorf("Error setting Environment: %s", err)
	}
	// Set from getRegion(d)
	if err := d.Set("region", envName.Region); err != nil {
		return fmt.Errorf("Error setting Environment: %s", err)
	}
	if err := d.Set("name", GetResourceNameFromSelfLink(res.Name)); err != nil {
		return fmt.Errorf("Error setting Environment: %s", err)
	}
	if err := d.Set("config", flattenComposerEnvironmentConfig(res.Config)); err != nil {
		return fmt.Errorf("Error setting Environment: %s", err)
	}
	if err := d.Set("labels", res.Labels); err != nil {
		return fmt.Errorf("Error setting Environment: %s", err)
	}
	return nil
}

func resourceComposerEnvironmentUpdate(d *schema.ResourceData, meta interface{}) error {
	tfConfig := meta.(*Config)
	userAgent, err := generateUserAgentString(d, tfConfig.userAgent)
	if err != nil {
		return err
	}

	d.Partial(true)

	// Composer only allows PATCHing one field at a time, so for each updatable field, we
	// 1. determine if it needs to be updated
	// 2. construct a PATCH object with only that field populated
	// 3. call resourceComposerEnvironmentPatchField(...)to update that single field.
	if d.HasChange("config") {
		config, err := expandComposerEnvironmentConfig(d.Get("config"), d, tfConfig)
		if err != nil {
			return err
		}

		if d.HasChange("config.0.software_config.0.scheduler_count") {
			patchObj := &composer.Environment{
				Config: &composer.EnvironmentConfig{
					SoftwareConfig: &composer.SoftwareConfig{},
				},
			}
			if config != nil && config.SoftwareConfig != nil {
				patchObj.Config.SoftwareConfig.SchedulerCount = config.SoftwareConfig.SchedulerCount
			}
			err = resourceComposerEnvironmentPatchField("config.softwareConfig.schedulerCount", userAgent, patchObj, d, tfConfig)
			if err != nil {
				return err
			}
		}

		if d.HasChange("config.0.software_config.0.airflow_config_overrides") {
			patchObj := &composer.Environment{
				Config: &composer.EnvironmentConfig{
					SoftwareConfig: &composer.SoftwareConfig{
						AirflowConfigOverrides: make(map[string]string),
					},
				},
			}

			if config != nil && config.SoftwareConfig != nil && len(config.SoftwareConfig.AirflowConfigOverrides) > 0 {
				patchObj.Config.SoftwareConfig.AirflowConfigOverrides = config.SoftwareConfig.AirflowConfigOverrides
			}

			err = resourceComposerEnvironmentPatchField("config.softwareConfig.airflowConfigOverrides", userAgent, patchObj, d, tfConfig)
			if err != nil {
				return err
			}
		}

		if d.HasChange("config.0.software_config.0.env_variables") {
			patchObj := &composer.Environment{
				Config: &composer.EnvironmentConfig{
					SoftwareConfig: &composer.SoftwareConfig{
						EnvVariables: make(map[string]string),
					},
				},
			}
			if config != nil && config.SoftwareConfig != nil && len(config.SoftwareConfig.EnvVariables) > 0 {
				patchObj.Config.SoftwareConfig.EnvVariables = config.SoftwareConfig.EnvVariables
			}

			err = resourceComposerEnvironmentPatchField("config.softwareConfig.envVariables", userAgent, patchObj, d, tfConfig)
			if err != nil {
				return err
			}
		}

		if d.HasChange("config.0.software_config.0.pypi_packages") {
			patchObj := &composer.Environment{
				Config: &composer.EnvironmentConfig{
					SoftwareConfig: &composer.SoftwareConfig{
						PypiPackages: make(map[string]string),
					},
				},
			}
			if config != nil && config.SoftwareConfig != nil && config.SoftwareConfig.PypiPackages != nil {
				patchObj.Config.SoftwareConfig.PypiPackages = config.SoftwareConfig.PypiPackages
			}

			err = resourceComposerEnvironmentPatchField("config.softwareConfig.pypiPackages", userAgent, patchObj, d, tfConfig)
			if err != nil {
				return err
			}
		}

		if d.HasChange("config.0.node_count") {
			patchObj := &composer.Environment{Config: &composer.EnvironmentConfig{}}
			if config != nil {
				patchObj.Config.NodeCount = config.NodeCount
			}
			err = resourceComposerEnvironmentPatchField("config.nodeCount", userAgent, patchObj, d, tfConfig)
			if err != nil {
				return err
			}
		}

		// If web_server_network_access_control has more fields added it may require changes here.
		// This is scoped specifically to allowed_ip_range due to https://github.com/hashicorp/terraform-plugin-sdk/issues/98
		if d.HasChange("config.0.web_server_network_access_control.0.allowed_ip_range") {
			patchObj := &composer.Environment{Config: &composer.EnvironmentConfig{}}
			if config != nil {
				patchObj.Config.WebServerNetworkAccessControl = config.WebServerNetworkAccessControl
			}
			err = resourceComposerEnvironmentPatchField("config.webServerNetworkAccessControl", userAgent, patchObj, d, tfConfig)
			if err != nil {
				return err
			}
		}

		if d.HasChange("config.0.database_config.0.machine_type") {
			patchObj := &composer.Environment{Config: &composer.EnvironmentConfig{}}
			if config != nil {
				patchObj.Config.DatabaseConfig = config.DatabaseConfig
			}
			err = resourceComposerEnvironmentPatchField("config.databaseConfig.machineType", userAgent, patchObj, d, tfConfig)
			if err != nil {
				return err
			}
		}

		if d.HasChange("config.0.web_server_config.0.machine_type") {
			patchObj := &composer.Environment{Config: &composer.EnvironmentConfig{}}
			if config != nil {
				patchObj.Config.WebServerConfig = config.WebServerConfig
			}
			err = resourceComposerEnvironmentPatchField("config.webServerConfig.machineType", userAgent, patchObj, d, tfConfig)
			if err != nil {
				return err
			}
		}

		if d.HasChange("config.0.maintenance_window") {
			patchObj := &composer.Environment{Config: &composer.EnvironmentConfig{}}
			if config != nil {
				patchObj.Config.MaintenanceWindow = config.MaintenanceWindow
			}
			err = resourceComposerEnvironmentPatchField("config.maintenanceWindow", userAgent, patchObj, d, tfConfig)
			if err != nil {
				return err
			}
		}

		if d.HasChange("config.0.workloads_config") {
			patchObj := &composer.Environment{Config: &composer.EnvironmentConfig{}}
			if config != nil {
				patchObj.Config.WorkloadsConfig = config.WorkloadsConfig
			}
			err = resourceComposerEnvironmentPatchField("config.WorkloadsConfig", userAgent, patchObj, d, tfConfig)
			if err != nil {
				return err
			}
		}
		if d.HasChange("config.0.environment_size") {
			patchObj := &composer.Environment{Config: &composer.EnvironmentConfig{}}
			if config != nil {
				patchObj.Config.EnvironmentSize = config.EnvironmentSize
			}
			err = resourceComposerEnvironmentPatchField("config.EnvironmentSize", userAgent, patchObj, d, tfConfig)
			if err != nil {
				return err
			}
		}
		if d.HasChange("config.0.master_authorized_networks_config") {
			patchObj := &composer.Environment{Config: &composer.EnvironmentConfig{}}
			if config != nil {
				patchObj.Config.MasterAuthorizedNetworksConfig = config.MasterAuthorizedNetworksConfig
			}
			err = resourceComposerEnvironmentPatchField("config.MasterAuthorizedNetworksConfig", userAgent, patchObj, d, tfConfig)
			if err != nil {
				return err
			}
		}
	}

	if d.HasChange("labels") {
		patchEnv := &composer.Environment{Labels: expandLabels(d)}
		err := resourceComposerEnvironmentPatchField("labels", userAgent, patchEnv, d, tfConfig)
		if err != nil {
			return err
		}
	}

	d.Partial(false)
	return resourceComposerEnvironmentRead(d, tfConfig)
}

func resourceComposerEnvironmentPostCreateUpdate(updateEnv *composer.Environment, d *schema.ResourceData, cfg *Config, userAgent string) error {
	if updateEnv == nil {
		return nil
	}

	d.Partial(true)

	if updateEnv.Config != nil && updateEnv.Config.SoftwareConfig != nil && len(updateEnv.Config.SoftwareConfig.PypiPackages) > 0 {
		log.Printf("[DEBUG] Running post-create update for Environment %q", d.Id())
		err := resourceComposerEnvironmentPatchField("config.softwareConfig.pypiPackages", userAgent, updateEnv, d, cfg)
		if err != nil {
			return err
		}

		log.Printf("[DEBUG] Finish update to Environment %q post create for update only fields", d.Id())
	}
	d.Partial(false)
	return resourceComposerEnvironmentRead(d, cfg)
}

func resourceComposerEnvironmentPatchField(updateMask, userAgent string, env *composer.Environment, d *schema.ResourceData, config *Config) error {
	envJson, _ := env.MarshalJSON()
	log.Printf("[DEBUG] Updating Environment %q (updateMask = %q): %s", d.Id(), updateMask, string(envJson))
	envName, err := resourceComposerEnvironmentName(d, config)
	if err != nil {
		return err
	}

	op, err := config.NewComposerClient(userAgent).Projects.Locations.Environments.
		Patch(envName.resourceName(), env).
		UpdateMask(updateMask).Do()
	if err != nil {
		return err
	}

	waitErr := composerOperationWaitTime(
		config, op, envName.Project, "Updating newly created Environment", userAgent,
		d.Timeout(schema.TimeoutCreate))
	if waitErr != nil {
		// The resource didn't actually update.
		return fmt.Errorf("Error waiting to update Environment: %s", waitErr)
	}

	log.Printf("[DEBUG] Finished updating Environment %q (updateMask = %q)", d.Id(), updateMask)
	return nil
}

func resourceComposerEnvironmentDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	envName, err := resourceComposerEnvironmentName(d, config)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting Environment %q", d.Id())
	op, err := config.NewComposerClient(userAgent).Projects.Locations.Environments.Delete(envName.resourceName()).Do()
	if err != nil {
		return err
	}

	err = composerOperationWaitTime(
		config, op, envName.Project, "Deleting Environment", userAgent,
		d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Finished deleting Environment %q: %#v", d.Id(), op)
	return nil
}

func resourceComposerEnvironmentImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{"projects/(?P<project>[^/]+)/locations/(?P<region>[^/]+)/environments/(?P<name>[^/]+)", "(?P<project>[^/]+)/(?P<region>[^/]+)/(?P<name>[^/]+)", "(?P<name>[^/]+)"}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "projects/{{project}}/locations/{{region}}/environments/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenComposerEnvironmentConfig(envCfg *composer.EnvironmentConfig) interface{} {
	if envCfg == nil {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["gke_cluster"] = envCfg.GkeCluster
	transformed["dag_gcs_prefix"] = envCfg.DagGcsPrefix
	transformed["node_count"] = envCfg.NodeCount
	transformed["airflow_uri"] = envCfg.AirflowUri
	transformed["node_config"] = flattenComposerEnvironmentConfigNodeConfig(envCfg.NodeConfig)
	transformed["software_config"] = flattenComposerEnvironmentConfigSoftwareConfig(envCfg.SoftwareConfig)
	transformed["private_environment_config"] = flattenComposerEnvironmentConfigPrivateEnvironmentConfig(envCfg.PrivateEnvironmentConfig)
	transformed["web_server_network_access_control"] = flattenComposerEnvironmentConfigWebServerNetworkAccessControl(envCfg.WebServerNetworkAccessControl)
	transformed["database_config"] = flattenComposerEnvironmentConfigDatabaseConfig(envCfg.DatabaseConfig)
	transformed["web_server_config"] = flattenComposerEnvironmentConfigWebServerConfig(envCfg.WebServerConfig)
	transformed["encryption_config"] = flattenComposerEnvironmentConfigEncryptionConfig(envCfg.EncryptionConfig)
	transformed["maintenance_window"] = flattenComposerEnvironmentConfigMaintenanceWindow(envCfg.MaintenanceWindow)
	transformed["workloads_config"] = flattenComposerEnvironmentConfigWorkloadsConfig(envCfg.WorkloadsConfig)
	transformed["environment_size"] = envCfg.EnvironmentSize
	transformed["master_authorized_networks_config"] = flattenComposerEnvironmentConfigMasterAuthorizedNetworksConfig(envCfg.MasterAuthorizedNetworksConfig)
	return []interface{}{transformed}
}

func flattenComposerEnvironmentConfigWebServerNetworkAccessControl(accessControl *composer.WebServerNetworkAccessControl) interface{} {
	if accessControl == nil || accessControl.AllowedIpRanges == nil {
		return nil
	}

	transformed := make([]interface{}, 0, len(accessControl.AllowedIpRanges))
	for _, ipRange := range accessControl.AllowedIpRanges {
		data := map[string]interface{}{
			"value":       ipRange.Value,
			"description": ipRange.Description,
		}
		transformed = append(transformed, data)
	}

	webServerNetworkAccessControl := make(map[string]interface{})

	webServerNetworkAccessControl["allowed_ip_range"] = schema.NewSet(schema.HashResource(allowedIpRangesConfig), transformed)

	return []interface{}{webServerNetworkAccessControl}
}

func flattenComposerEnvironmentConfigDatabaseConfig(databaseCfg *composer.DatabaseConfig) interface{} {
	if databaseCfg == nil {
		return nil
	}

	transformed := make(map[string]interface{})
	transformed["machine_type"] = databaseCfg.MachineType

	return []interface{}{transformed}
}

func flattenComposerEnvironmentConfigWebServerConfig(webServerCfg *composer.WebServerConfig) interface{} {
	if webServerCfg == nil {
		return nil
	}

	transformed := make(map[string]interface{})
	transformed["machine_type"] = webServerCfg.MachineType

	return []interface{}{transformed}
}

func flattenComposerEnvironmentConfigEncryptionConfig(encryptionCfg *composer.EncryptionConfig) interface{} {
	if encryptionCfg == nil {
		return nil
	}

	transformed := make(map[string]interface{})
	transformed["kms_key_name"] = encryptionCfg.KmsKeyName

	return []interface{}{transformed}
}

func flattenComposerEnvironmentConfigMaintenanceWindow(maintenanceWindow *composer.MaintenanceWindow) interface{} {
	if maintenanceWindow == nil {
		return nil
	}

	transformed := make(map[string]interface{})
	transformed["start_time"] = maintenanceWindow.StartTime
	transformed["end_time"] = maintenanceWindow.EndTime
	transformed["recurrence"] = maintenanceWindow.Recurrence

	return []interface{}{transformed}
}

func flattenComposerEnvironmentConfigWorkloadsConfig(workloadsConfig *composer.WorkloadsConfig) interface{} {
	if workloadsConfig == nil {
		return nil
	}

	transformed := make(map[string]interface{})
	transformedScheduler := make(map[string]interface{})
	transformedWebServer := make(map[string]interface{})
	transformedWorker := make(map[string]interface{})

	wlCfgScheduler := workloadsConfig.Scheduler
	wlCfgWebServer := workloadsConfig.WebServer
	wlCfgWorker := workloadsConfig.Worker

	if wlCfgScheduler == nil {
		transformedScheduler = nil
	} else {
		transformedScheduler["cpu"] = wlCfgScheduler.Cpu
		transformedScheduler["memory_gb"] = wlCfgScheduler.MemoryGb
		transformedScheduler["storage_gb"] = wlCfgScheduler.StorageGb
		transformedScheduler["count"] = wlCfgScheduler.Count
	}

	if wlCfgWebServer == nil {
		transformedWebServer = nil
	} else {
		transformedWebServer["cpu"] = wlCfgWebServer.Cpu
		transformedWebServer["memory_gb"] = wlCfgWebServer.MemoryGb
		transformedWebServer["storage_gb"] = wlCfgWebServer.StorageGb
	}

	if wlCfgWorker == nil {
		transformedWorker = nil
	} else {
		transformedWorker["cpu"] = wlCfgWorker.Cpu
		transformedWorker["memory_gb"] = wlCfgWorker.MemoryGb
		transformedWorker["storage_gb"] = wlCfgWorker.StorageGb
		transformedWorker["min_count"] = wlCfgWorker.MinCount
		transformedWorker["max_count"] = wlCfgWorker.MaxCount
	}

	transformed["scheduler"] = []interface{}{transformedScheduler}
	transformed["web_server"] = []interface{}{transformedWebServer}
	transformed["worker"] = []interface{}{transformedWorker}

	return []interface{}{transformed}
}

func flattenComposerEnvironmentConfigPrivateEnvironmentConfig(envCfg *composer.PrivateEnvironmentConfig) interface{} {
	if envCfg == nil {
		return nil
	}

	transformed := make(map[string]interface{})
	transformed["enable_private_endpoint"] = envCfg.PrivateClusterConfig.EnablePrivateEndpoint
	transformed["master_ipv4_cidr_block"] = envCfg.PrivateClusterConfig.MasterIpv4CidrBlock
	transformed["cloud_sql_ipv4_cidr_block"] = envCfg.CloudSqlIpv4CidrBlock
	transformed["web_server_ipv4_cidr_block"] = envCfg.WebServerIpv4CidrBlock
	transformed["cloud_composer_network_ipv4_cidr_block"] = envCfg.CloudComposerNetworkIpv4CidrBlock
	transformed["enable_privately_used_public_ips"] = envCfg.EnablePrivatelyUsedPublicIps
	transformed["cloud_composer_connection_subnetwork"] = envCfg.CloudComposerConnectionSubnetwork

	return []interface{}{transformed}
}

func flattenComposerEnvironmentConfigNodeConfig(nodeCfg *composer.NodeConfig) interface{} {
	if nodeCfg == nil {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["zone"] = nodeCfg.Location
	transformed["machine_type"] = nodeCfg.MachineType
	transformed["network"] = nodeCfg.Network
	transformed["subnetwork"] = nodeCfg.Subnetwork
	transformed["disk_size_gb"] = nodeCfg.DiskSizeGb
	transformed["service_account"] = nodeCfg.ServiceAccount
	transformed["oauth_scopes"] = flattenComposerEnvironmentConfigNodeConfigOauthScopes(nodeCfg.OauthScopes)
	transformed["enable_ip_masq_agent"] = nodeCfg.EnableIpMasqAgent
	transformed["tags"] = flattenComposerEnvironmentConfigNodeConfigTags(nodeCfg.Tags)
	transformed["ip_allocation_policy"] = flattenComposerEnvironmentConfigNodeConfigIPAllocationPolicy(nodeCfg.IpAllocationPolicy)
	return []interface{}{transformed}
}

func flattenComposerEnvironmentConfigNodeConfigIPAllocationPolicy(ipPolicy *composer.IPAllocationPolicy) interface{} {
	if ipPolicy == nil {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["use_ip_aliases"] = ipPolicy.UseIpAliases
	transformed["cluster_ipv4_cidr_block"] = ipPolicy.ClusterIpv4CidrBlock
	transformed["cluster_secondary_range_name"] = ipPolicy.ClusterSecondaryRangeName
	transformed["services_ipv4_cidr_block"] = ipPolicy.ServicesIpv4CidrBlock
	transformed["services_secondary_range_name"] = ipPolicy.ServicesSecondaryRangeName

	return []interface{}{transformed}
}

func flattenComposerEnvironmentConfigNodeConfigOauthScopes(v interface{}) interface{} {
	if v == nil {
		return v
	}
	return schema.NewSet(schema.HashString, convertStringArrToInterface(v.([]string)))
}

func flattenComposerEnvironmentConfigNodeConfigTags(v interface{}) interface{} {
	if v == nil {
		return v
	}
	return schema.NewSet(schema.HashString, convertStringArrToInterface(v.([]string)))
}

func flattenComposerEnvironmentConfigSoftwareConfig(softwareCfg *composer.SoftwareConfig) interface{} {
	if softwareCfg == nil {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["image_version"] = softwareCfg.ImageVersion
	transformed["python_version"] = softwareCfg.PythonVersion
	transformed["airflow_config_overrides"] = softwareCfg.AirflowConfigOverrides
	transformed["pypi_packages"] = softwareCfg.PypiPackages
	transformed["env_variables"] = softwareCfg.EnvVariables
	transformed["scheduler_count"] = softwareCfg.SchedulerCount
	return []interface{}{transformed}
}

func flattenComposerEnvironmentConfigMasterAuthorizedNetworksConfig(masterAuthNetsCfg *composer.MasterAuthorizedNetworksConfig) interface{} {
	if masterAuthNetsCfg == nil {
		return nil
	}

	transformed := make([]interface{}, 0, len(masterAuthNetsCfg.CidrBlocks))
	for _, cidrBlock := range masterAuthNetsCfg.CidrBlocks {
		data := map[string]interface{}{
			"display_name": cidrBlock.DisplayName,
			"cidr_block":   cidrBlock.CidrBlock,
		}
		transformed = append(transformed, data)
	}

	masterAuthorizedNetworksConfig := make(map[string]interface{})
	masterAuthorizedNetworksConfig["enabled"] = masterAuthNetsCfg.Enabled
	masterAuthorizedNetworksConfig["cidr_blocks"] = schema.NewSet(schema.HashResource(cidrBlocks), transformed)

	return []interface{}{masterAuthorizedNetworksConfig}
}

func expandComposerEnvironmentConfig(v interface{}, d *schema.ResourceData, config *Config) (*composer.EnvironmentConfig, error) {
	l := v.([]interface{})
	if len(l) == 0 {
		return nil, nil
	}
	original := l[0].(map[string]interface{})
	transformed := &composer.EnvironmentConfig{}

	if nodeCountRaw, ok := original["node_count"]; ok {
		transformedNodeCount, err := expandComposerEnvironmentConfigNodeCount(nodeCountRaw, d, config)
		if err != nil {
			return nil, err
		}
		transformed.NodeCount = transformedNodeCount
	}

	transformedNodeConfig, err := expandComposerEnvironmentConfigNodeConfig(original["node_config"], d, config)
	if err != nil {
		return nil, err
	}
	transformed.NodeConfig = transformedNodeConfig

	transformedSoftwareConfig, err := expandComposerEnvironmentConfigSoftwareConfig(original["software_config"], d, config)
	if err != nil {
		return nil, err
	}
	transformed.SoftwareConfig = transformedSoftwareConfig

	transformedPrivateEnvironmentConfig, err := expandComposerEnvironmentConfigPrivateEnvironmentConfig(original["private_environment_config"], d, config)
	if err != nil {
		return nil, err
	}
	transformed.PrivateEnvironmentConfig = transformedPrivateEnvironmentConfig

	transformedWebServerNetworkAccessControl, err := expandComposerEnvironmentConfigWebServerNetworkAccessControl(original["web_server_network_access_control"], d, config)
	if err != nil {
		return nil, err
	}
	transformed.WebServerNetworkAccessControl = transformedWebServerNetworkAccessControl

	transformedDatabaseConfig, err := expandComposerEnvironmentConfigDatabaseConfig(original["database_config"], d, config)
	if err != nil {
		return nil, err
	}
	transformed.DatabaseConfig = transformedDatabaseConfig

	transformedWebServerConfig, err := expandComposerEnvironmentConfigWebServerConfig(original["web_server_config"], d, config)
	if err != nil {
		return nil, err
	}
	transformed.WebServerConfig = transformedWebServerConfig

	transformedEncryptionConfig, err := expandComposerEnvironmentConfigEncryptionConfig(original["encryption_config"], d, config)
	if err != nil {
		return nil, err
	}
	transformed.EncryptionConfig = transformedEncryptionConfig

	transformedMaintenanceWindow, err := expandComposerEnvironmentConfigMaintenanceWindow(original["maintenance_window"], d, config)
	if err != nil {
		return nil, err
	}
	transformed.MaintenanceWindow = transformedMaintenanceWindow
	transformedWorkloadsConfig, err := expandComposerEnvironmentConfigWorkloadsConfig(original["workloads_config"], d, config)
	if err != nil {
		return nil, err
	}
	transformed.WorkloadsConfig = transformedWorkloadsConfig

	transformedEnvironmentSize, err := expandComposerEnvironmentConfigEnvironmentSize(original["environment_size"], d, config)
	if err != nil {
		return nil, err
	}
	transformed.EnvironmentSize = transformedEnvironmentSize
	transformedMasterAuthorizedNetworksConfig, err := expandComposerEnvironmentConfigMasterAuthorizedNetworksConfig(original["master_authorized_networks_config"], d, config)
	if err != nil {
		return nil, err
	}
	transformed.MasterAuthorizedNetworksConfig = transformedMasterAuthorizedNetworksConfig
	return transformed, nil
}

func expandComposerEnvironmentConfigNodeCount(v interface{}, d *schema.ResourceData, config *Config) (int64, error) {
	if v == nil {
		return 0, nil
	}
	return int64(v.(int)), nil
}

func expandComposerEnvironmentConfigWebServerNetworkAccessControl(v interface{}, d *schema.ResourceData, config *Config) (*composer.WebServerNetworkAccessControl, error) {
	l := v.([]interface{})
	if len(l) == 0 {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})

	allowedIpRangesRaw := original["allowed_ip_range"].(*schema.Set).List()
	if len(allowedIpRangesRaw) == 0 {
		return nil, nil
	}

	transformed := &composer.WebServerNetworkAccessControl{}
	allowedIpRanges := make([]*composer.AllowedIpRange, 0, len(original))

	for _, originalIpRange := range allowedIpRangesRaw {
		originalRangeRaw := originalIpRange.(map[string]interface{})
		transformedRange := &composer.AllowedIpRange{Value: originalRangeRaw["value"].(string)}
		if v, ok := originalRangeRaw["description"]; ok {
			transformedRange.Description = v.(string)
		}
		allowedIpRanges = append(allowedIpRanges, transformedRange)
	}

	transformed.AllowedIpRanges = allowedIpRanges
	return transformed, nil
}

func expandComposerEnvironmentConfigMasterAuthorizedNetworksConfig(v interface{}, d *schema.ResourceData, config *Config) (*composer.MasterAuthorizedNetworksConfig, error) {
	l := v.([]interface{})
	if len(l) == 0 {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})

	cidrBlocksRaw := original["cidr_blocks"].(*schema.Set).List()

	transformed := &composer.MasterAuthorizedNetworksConfig{}
	cidrBlocks := make([]*composer.CidrBlock, 0, len(original))

	for _, originalCidrBlock := range cidrBlocksRaw {
		originalCidrBlockRaw := originalCidrBlock.(map[string]interface{})
		transformedCidrBlock := &composer.CidrBlock{}
		if v, ok := originalCidrBlockRaw["display_name"]; ok {
			transformedCidrBlock.DisplayName = v.(string)
		}
		transformedCidrBlock.CidrBlock = originalCidrBlockRaw["cidr_block"].(string)
		cidrBlocks = append(cidrBlocks, transformedCidrBlock)
	}
	transformed.Enabled = original["enabled"].(bool)
	transformed.CidrBlocks = cidrBlocks
	return transformed, nil
}

func expandComposerEnvironmentConfigDatabaseConfig(v interface{}, d *schema.ResourceData, config *Config) (*composer.DatabaseConfig, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})

	transformed := &composer.DatabaseConfig{}
	transformed.MachineType = original["machine_type"].(string)

	return transformed, nil
}

func expandComposerEnvironmentConfigWebServerConfig(v interface{}, d *schema.ResourceData, config *Config) (*composer.WebServerConfig, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})

	transformed := &composer.WebServerConfig{}
	transformed.MachineType = original["machine_type"].(string)

	return transformed, nil
}

func expandComposerEnvironmentConfigEncryptionConfig(v interface{}, d *schema.ResourceData, config *Config) (*composer.EncryptionConfig, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})

	transformed := &composer.EncryptionConfig{}
	transformed.KmsKeyName = original["kms_key_name"].(string)

	return transformed, nil
}

func expandComposerEnvironmentConfigMaintenanceWindow(v interface{}, d *schema.ResourceData, config *Config) (*composer.MaintenanceWindow, error) {
	l := v.([]interface{})
	if len(l) == 0 {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := &composer.MaintenanceWindow{}

	if v, ok := original["start_time"]; ok {
		transformed.StartTime = v.(string)
	}

	if v, ok := original["end_time"]; ok {
		transformed.EndTime = v.(string)
	}

	if v, ok := original["recurrence"]; ok {
		transformed.Recurrence = v.(string)
	}

	return transformed, nil
}

func expandComposerEnvironmentConfigWorkloadsConfig(v interface{}, d *schema.ResourceData, config *Config) (*composer.WorkloadsConfig, error) {
	l := v.([]interface{})
	if len(l) == 0 {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := &composer.WorkloadsConfig{}

	if v, ok := original["scheduler"]; ok {
		if len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
			transformedScheduler := &composer.SchedulerResource{}
			originalSchedulerRaw := v.([]interface{})[0].(map[string]interface{})
			transformedScheduler.Count = int64(originalSchedulerRaw["count"].(int))
			transformedScheduler.Cpu = originalSchedulerRaw["cpu"].(float64)
			transformedScheduler.MemoryGb = originalSchedulerRaw["memory_gb"].(float64)
			transformedScheduler.StorageGb = originalSchedulerRaw["storage_gb"].(float64)
			transformed.Scheduler = transformedScheduler
		}
	}

	if v, ok := original["web_server"]; ok {
		if len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
			transformedWebServer := &composer.WebServerResource{}
			originalWebServerRaw := v.([]interface{})[0].(map[string]interface{})
			transformedWebServer.Cpu = originalWebServerRaw["cpu"].(float64)
			transformedWebServer.MemoryGb = originalWebServerRaw["memory_gb"].(float64)
			transformedWebServer.StorageGb = originalWebServerRaw["storage_gb"].(float64)
			transformed.WebServer = transformedWebServer
		}
	}

	if v, ok := original["worker"]; ok {
		if len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
			transformedWorker := &composer.WorkerResource{}
			originalWorkerRaw := v.([]interface{})[0].(map[string]interface{})
			transformedWorker.Cpu = originalWorkerRaw["cpu"].(float64)
			transformedWorker.MemoryGb = originalWorkerRaw["memory_gb"].(float64)
			transformedWorker.StorageGb = originalWorkerRaw["storage_gb"].(float64)
			transformedWorker.MinCount = int64(originalWorkerRaw["min_count"].(int))
			transformedWorker.MaxCount = int64(originalWorkerRaw["max_count"].(int))
			transformed.Worker = transformedWorker
		}
	}

	return transformed, nil
}

func expandComposerEnvironmentConfigEnvironmentSize(v interface{}, d *schema.ResourceData, config *Config) (string, error) {
	if v == nil {
		return "", nil
	}
	return v.(string), nil
}

func expandComposerEnvironmentConfigPrivateEnvironmentConfig(v interface{}, d *schema.ResourceData, config *Config) (*composer.PrivateEnvironmentConfig, error) {
	l := v.([]interface{})
	if len(l) == 0 {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := &composer.PrivateEnvironmentConfig{
		EnablePrivateEnvironment: true,
	}

	subBlock := &composer.PrivateClusterConfig{}

	if v, ok := original["enable_private_endpoint"]; ok {
		subBlock.EnablePrivateEndpoint = v.(bool)
	}

	if v, ok := original["master_ipv4_cidr_block"]; ok {
		subBlock.MasterIpv4CidrBlock = v.(string)
	}

	if v, ok := original["cloud_sql_ipv4_cidr_block"]; ok {
		transformed.CloudSqlIpv4CidrBlock = v.(string)
	}

	if v, ok := original["web_server_ipv4_cidr_block"]; ok {
		transformed.WebServerIpv4CidrBlock = v.(string)
	}

	if v, ok := original["cloud_composer_network_ipv4_cidr_block"]; ok {
		transformed.CloudComposerNetworkIpv4CidrBlock = v.(string)
	}
	if v, ok := original["enable_privately_used_public_ips"]; ok {
		transformed.EnablePrivatelyUsedPublicIps = v.(bool)
	}
	if v, ok := original["cloud_composer_connection_subnetwork"]; ok {
		transformed.CloudComposerConnectionSubnetwork = v.(string)
	}

	transformed.PrivateClusterConfig = subBlock

	return transformed, nil
}

func expandComposerEnvironmentConfigNodeConfig(v interface{}, d *schema.ResourceData, config *Config) (*composer.NodeConfig, error) {
	l := v.([]interface{})
	if len(l) == 0 {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := &composer.NodeConfig{}

	if transformedDiskSizeGb, ok := original["disk_size_gb"]; ok {
		transformed.DiskSizeGb = int64(transformedDiskSizeGb.(int))
	}

	if v, ok := original["service_account"]; ok {
		transformedServiceAccount, err := expandComposerEnvironmentServiceAccount(v, d, config)
		if err != nil {
			return nil, err
		}
		transformed.ServiceAccount = transformedServiceAccount
	}

	if transformedEnableIpMasqAgent, ok := original["enable_ip_masq_agent"]; ok {
		transformed.EnableIpMasqAgent = transformedEnableIpMasqAgent.(bool)
	}

	var nodeConfigZone string
	if v, ok := original["zone"]; ok {
		transformedZone, err := expandComposerEnvironmentZone(v, d, config)
		if err != nil {
			return nil, err
		}
		transformed.Location = transformedZone
		nodeConfigZone = transformedZone
	}

	if v, ok := original["machine_type"]; ok {
		transformedMachineType, err := expandComposerEnvironmentMachineType(v, d, config, nodeConfigZone)
		if err != nil {
			return nil, err
		}
		transformed.MachineType = transformedMachineType
	}

	if v, ok := original["network"]; ok {
		transformedNetwork, err := expandComposerEnvironmentNetwork(v, d, config)
		if err != nil {
			return nil, err
		}
		transformed.Network = transformedNetwork
	}

	if v, ok := original["subnetwork"]; ok {
		transformedSubnetwork, err := expandComposerEnvironmentSubnetwork(v, d, config)
		if err != nil {
			return nil, err
		}
		transformed.Subnetwork = transformedSubnetwork
	}
	transformedIPAllocationPolicy, err := expandComposerEnvironmentIPAllocationPolicy(original["ip_allocation_policy"], d, config)
	if err != nil {
		return nil, err
	}
	transformed.IpAllocationPolicy = transformedIPAllocationPolicy

	transformedOauthScopes, err := expandComposerEnvironmentSetList(original["oauth_scopes"], d, config)
	if err != nil {
		return nil, err
	}
	transformed.OauthScopes = transformedOauthScopes

	transformedTags, err := expandComposerEnvironmentSetList(original["tags"], d, config)
	if err != nil {
		return nil, err
	}
	transformed.Tags = transformedTags

	return transformed, nil
}

func expandComposerEnvironmentIPAllocationPolicy(v interface{}, d *schema.ResourceData, config *Config) (*composer.IPAllocationPolicy, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := &composer.IPAllocationPolicy{}

	if v, ok := original["use_ip_aliases"]; ok {
		transformed.UseIpAliases = v.(bool)
	}

	if v, ok := original["cluster_ipv4_cidr_block"]; ok {
		transformed.ClusterIpv4CidrBlock = v.(string)
	}

	if v, ok := original["cluster_secondary_range_name"]; ok {
		transformed.ClusterSecondaryRangeName = v.(string)
	}

	if v, ok := original["services_ipv4_cidr_block"]; ok {
		transformed.ServicesIpv4CidrBlock = v.(string)
	}

	if v, ok := original["services_secondary_range_name"]; ok {
		transformed.ServicesSecondaryRangeName = v.(string)
	}
	return transformed, nil

}

func expandComposerEnvironmentServiceAccount(v interface{}, d *schema.ResourceData, config *Config) (string, error) {
	serviceAccount := v.(string)
	if len(serviceAccount) == 0 {
		return "", nil
	}

	return GetResourceNameFromSelfLink(serviceAccount), nil
}

func expandComposerEnvironmentZone(v interface{}, d *schema.ResourceData, config *Config) (string, error) {
	zone := v.(string)
	if len(zone) == 0 {
		return zone, nil
	}
	if !strings.Contains(zone, "/") {
		project, err := getProject(d, config)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("projects/%s/zones/%s", project, zone), nil
	}

	return getRelativePath(zone)
}

func expandComposerEnvironmentMachineType(v interface{}, d *schema.ResourceData, config *Config, nodeCfgZone string) (string, error) {
	machineType := v.(string)
	requiredZone := GetResourceNameFromSelfLink(nodeCfgZone)

	fv, err := ParseMachineTypesFieldValue(v.(string), d, config)
	if err != nil {

		// Try to construct machine type with zone/project given in config.
		project, err := getProject(d, config)
		if err != nil {
			return "", err
		}

		fv = &ZonalFieldValue{
			Project:      project,
			Zone:         requiredZone,
			Name:         GetResourceNameFromSelfLink(machineType),
			resourceType: "machineTypes",
		}
	}

	// Make sure zone in node_config.machineType matches node_config.zone if
	// given.
	if requiredZone != "" && fv.Zone != requiredZone {
		return "", fmt.Errorf("node_config machine_type %q must be in node_config zone %q", machineType, requiredZone)
	}
	return fv.RelativeLink(), nil
}

func expandComposerEnvironmentNetwork(v interface{}, d *schema.ResourceData, config *Config) (string, error) {
	fv, err := ParseNetworkFieldValue(v.(string), d, config)
	if err != nil {
		return "", err
	}
	return fv.RelativeLink(), nil
}

func expandComposerEnvironmentSubnetwork(v interface{}, d *schema.ResourceData, config *Config) (string, error) {
	fv, err := ParseSubnetworkFieldValue(v.(string), d, config)
	if err != nil {
		return "", err
	}
	return fv.RelativeLink(), nil
}

func expandComposerEnvironmentSetList(v interface{}, d *schema.ResourceData, config *Config) ([]string, error) {
	if v == nil {
		return nil, nil
	}
	return convertStringArr(v.(*schema.Set).List()), nil
}

func expandComposerEnvironmentConfigSoftwareConfig(v interface{}, d *schema.ResourceData, config *Config) (*composer.SoftwareConfig, error) {
	l := v.([]interface{})
	if len(l) == 0 {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := &composer.SoftwareConfig{}

	transformed.ImageVersion = original["image_version"].(string)
	transformed.PythonVersion = original["python_version"].(string)
	transformed.AirflowConfigOverrides = expandComposerEnvironmentConfigSoftwareConfigStringMap(original, "airflow_config_overrides")
	transformed.PypiPackages = expandComposerEnvironmentConfigSoftwareConfigStringMap(original, "pypi_packages")
	transformed.EnvVariables = expandComposerEnvironmentConfigSoftwareConfigStringMap(original, "env_variables")
	transformed.SchedulerCount = int64(original["scheduler_count"].(int))
	return transformed, nil
}

func expandComposerEnvironmentConfigSoftwareConfigStringMap(softwareConfig map[string]interface{}, k string) map[string]string {
	v, ok := softwareConfig[k]
	if ok && v != nil {
		return convertStringMap(v.(map[string]interface{}))
	}
	return map[string]string{}
}

func validateComposerEnvironmentPypiPackages(v interface{}, k string) (ws []string, errors []error) {
	if v == nil {
		return ws, errors
	}
	for pkgName := range v.(map[string]interface{}) {
		if pkgName != strings.ToLower(pkgName) {
			errors = append(errors,
				fmt.Errorf("PYPI package %q can only contain lowercase characters", pkgName))
		}
	}

	return ws, errors
}

func validateComposerEnvironmentEnvVariables(v interface{}, k string) (ws []string, errors []error) {
	if v == nil {
		return ws, errors
	}

	reEnvVarName := regexp.MustCompile(composerEnvironmentEnvVariablesRegexp)
	reAirflowReserved := regexp.MustCompile(composerEnvironmentReservedAirflowEnvVarRegexp)

	for envVarName := range v.(map[string]interface{}) {
		if !reEnvVarName.MatchString(envVarName) {
			errors = append(errors,
				fmt.Errorf("env_variable %q must match regexp %q", envVarName, composerEnvironmentEnvVariablesRegexp))
		} else if _, ok := composerEnvironmentReservedEnvVar[envVarName]; ok {
			errors = append(errors,
				fmt.Errorf("env_variable %q is a reserved name and cannot be used", envVarName))
		} else if reAirflowReserved.MatchString(envVarName) {
			errors = append(errors,
				fmt.Errorf("env_variable %q cannot match reserved Airflow variable names with regexp %q",
					envVarName, composerEnvironmentReservedAirflowEnvVarRegexp))
		}
	}

	return ws, errors
}

func handleComposerEnvironmentCreationOpFailure(id string, envName *composerEnvironmentName, d *schema.ResourceData, config *Config) error {
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	log.Printf("[WARNING] Creation operation for Composer Environment %q failed, check Environment isn't still running", id)
	// Try to get possible created but invalid environment.
	env, err := config.NewComposerClient(userAgent).Projects.Locations.Environments.Get(envName.resourceName()).Do()
	if err != nil {
		// If error is 401, we don't have to clean up environment, return nil.
		// Otherwise, we encountered another error.
		return handleNotFoundError(err, d, fmt.Sprintf("Composer Environment %q", envName.resourceName()))
	}

	if env.State == "CREATING" {
		return fmt.Errorf(
			"Getting creation operation state failed while waiting for environment to finish creating, "+
				"but environment seems to still be in 'CREATING' state. Wait for operation to finish and either "+
				"manually delete environment or import %q into your state", id)
	}

	log.Printf("[WARNING] Environment %q from failed creation operation was created, deleting.", id)
	op, err := config.NewComposerClient(userAgent).Projects.Locations.Environments.Delete(envName.resourceName()).Do()
	if err != nil {
		return fmt.Errorf("Could not delete the invalid created environment with state %q: %s", env.State, err)
	}

	waitErr := composerOperationWaitTime(
		config, op, envName.Project,
		fmt.Sprintf("Deleting invalid created Environment with state %q", env.State), userAgent,
		d.Timeout(schema.TimeoutCreate))
	if waitErr != nil {
		return fmt.Errorf("Error waiting to delete invalid Environment with state %q: %s", env.State, waitErr)
	}

	return nil
}

func getComposerEnvironmentPostCreateUpdateObj(env *composer.Environment) (updateEnv *composer.Environment) {
	// pypiPackages can only be added via update
	if env != nil && env.Config != nil && env.Config.SoftwareConfig != nil {
		if len(env.Config.SoftwareConfig.PypiPackages) > 0 {
			updateEnv = &composer.Environment{
				Config: &composer.EnvironmentConfig{
					SoftwareConfig: &composer.SoftwareConfig{
						PypiPackages: env.Config.SoftwareConfig.PypiPackages,
					},
				},
			}
			// Clear PYPI packages - otherwise, API will return error
			// that the create request is invalid.
			env.Config.SoftwareConfig.PypiPackages = make(map[string]string)
		}
	}

	return updateEnv
}

func resourceComposerEnvironmentName(d *schema.ResourceData, config *Config) (*composerEnvironmentName, error) {
	project, err := getProject(d, config)
	if err != nil {
		return nil, err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return nil, err
	}

	return &composerEnvironmentName{
		Project:     project,
		Region:      region,
		Environment: d.Get("name").(string),
	}, nil
}

type composerEnvironmentName struct {
	Project     string
	Region      string
	Environment string
}

func (n *composerEnvironmentName) resourceName() string {
	return fmt.Sprintf("projects/%s/locations/%s/environments/%s", n.Project, n.Region, n.Environment)
}

func (n *composerEnvironmentName) parentName() string {
	return fmt.Sprintf("projects/%s/locations/%s", n.Project, n.Region)
}

// The value we store (i.e. `old` in this method), might be only the service account email,
// but we expect either the email or the name (projects/.../serviceAccounts/...)
func compareServiceAccountEmailToLink(_, old, new string, _ *schema.ResourceData) bool {
	// old is the service account email returned from the server.
	if !strings.HasPrefix("projects/", old) {
		return old == GetResourceNameFromSelfLink(new)
	}
	return compareSelfLinkRelativePaths("", old, new, nil)
}

func validateServiceAccountRelativeNameOrEmail(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	serviceAccountRe := "(" + strings.Join(PossibleServiceAccountNames, "|") + ")"
	if strings.HasPrefix(value, "projects/") {
		serviceAccountRe = fmt.Sprintf("projects/(.+)/serviceAccounts/%s", serviceAccountRe)
	}
	r := regexp.MustCompile(serviceAccountRe)
	if !r.MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q (%q) doesn't match regexp %q", k, value, serviceAccountRe))
	}

	return
}

func composerImageVersionDiffSuppress(_, old, new string, _ *schema.ResourceData) bool {
	versionRe := regexp.MustCompile(composerEnvironmentVersionRegexp)
	oldVersions := versionRe.FindStringSubmatch(old)
	newVersions := versionRe.FindStringSubmatch(new)
	if oldVersions == nil || len(oldVersions) < 10 {
		// Somehow one of the versions didn't match the regexp or didn't
		// have values in the capturing groups. In that case, fall back to
		// an equality check.
		if old != "" {
			log.Printf("[WARN] Image version didn't match regexp: %s", old)
		}
		return old == new
	}
	if newVersions == nil || len(newVersions) < 10 {
		// Somehow one of the versions didn't match the regexp or didn't
		// have values in the capturing groups. In that case, fall back to
		// an equality check.
		if new != "" {
			log.Printf("[WARN] Image version didn't match regexp: %s", new)
		}
		return old == new
	}

	oldAirflow := oldVersions[5]
	oldAirflowMajor := oldVersions[6]
	oldAirflowMajorMinor := oldVersions[6] + oldVersions[8]
	newAirflow := newVersions[5]
	newAirflowMajor := newVersions[6]
	newAirflowMajorMinor := newVersions[6] + newVersions[8]
	// Check Airflow versions.
	if oldAirflow == oldAirflowMajor || newAirflow == newAirflowMajor {
		// If one of the Airflow versions specifies only major version
		// (like 1), we can only compare major versions.
		eq, err := versionsEqual(oldAirflowMajor, newAirflowMajor)
		if err != nil {
			log.Printf("[WARN] Could not parse airflow version, %s", err)
		}
		if !eq {
			return false
		}
	} else if oldAirflow == oldAirflowMajorMinor || newAirflow == newAirflowMajorMinor {
		// If one of the Airflow versions specifies only major and minor version
		// (like 1.10), we can only compare major and minor versions.
		eq, err := versionsEqual(oldAirflowMajorMinor, newAirflowMajorMinor)
		if err != nil {
			log.Printf("[WARN] Could not parse airflow version, %s", err)
		}
		if !eq {
			return false
		}
	} else {
		// Otherwise, we compare the full Airflow versions (like 1.10.15).
		eq, err := versionsEqual(oldAirflow, newAirflow)
		if err != nil {
			log.Printf("[WARN] Could not parse airflow version, %s", err)
		}
		if !eq {
			return false
		}
	}

	oldComposer := oldVersions[1]
	oldComposerMajor := oldVersions[2]
	newComposer := newVersions[1]
	newComposerMajor := newVersions[2]
	// Check Composer versions.
	if oldComposer == "latest" || newComposer == "latest" {
		// We don't know what the latest version is so we suppress the diff.
		return true
	} else if oldComposer == oldComposerMajor || newComposer == newComposerMajor {
		// If one of the Composer versions specifies only major version
		// (like 1), we can only compare major versions.
		eq, err := versionsEqual(oldComposerMajor, newComposerMajor)
		if err != nil {
			log.Printf("[WARN] Could not parse composer version, %s", err)
		}
		return eq
	} else {
		// Otherwise, we compare the full Composer versions (like 1.18.1).
		eq, err := versionsEqual(oldComposer, newComposer)
		if err != nil {
			log.Printf("[WARN] Could not parse composer version, %s", err)
		}
		return eq
	}
}

func versionsEqual(old, new string) (bool, error) {
	o, err := version.NewVersion(old)
	if err != nil {
		return false, err
	}
	n, err := version.NewVersion(new)
	if err != nil {
		return false, err
	}
	return o.Equal(n), nil
}
