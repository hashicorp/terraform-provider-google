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
	composer "google.golang.org/api/composer/v1beta1"
)

const (
	composerEnvironmentEnvVariablesRegexp          = "[a-zA-Z_][a-zA-Z0-9_]*."
	composerEnvironmentReservedAirflowEnvVarRegexp = "AIRFLOW__[A-Z0-9_]+__[A-Z0-9_]+"
	composerEnvironmentVersionRegexp               = `composer-([0-9]+\.[0-9]+\.[0-9]+|latest)-airflow-([0-9]+\.[0-9]+(\.[0-9]+.*)?)`
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
	}

	composerConfigKeys = []string{
		"config.0.node_count",
		"config.0.node_config",
		"config.0.software_config",
		"config.0.private_environment_config",
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
				ValidateFunc: validateGCPName,
				Description:  `Name of the environment.`,
			},
			"region": {
				Type:        schema.TypeString,
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
							Description:  `The number of nodes in the Kubernetes Engine cluster that will be used to run this environment.`,
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
										Required:         true,
										ForceNew:         true,
										DiffSuppressFunc: compareSelfLinkOrResourceName,
										Description:      `The Compute Engine zone in which to deploy the VMs running the Apache Airflow software, specified as the zone name or relative resource name (e.g. "projects/{project}/zones/{zone}"). Must belong to the enclosing environment's project and region.`,
									},
									"machine_type": {
										Type:             schema.TypeString,
										Computed:         true,
										Optional:         true,
										ForceNew:         true,
										DiffSuppressFunc: compareSelfLinkOrResourceName,
										Description:      `The Compute Engine machine type used for cluster instances, specified as a name or relative resource name. For example: "projects/{project}/zones/{zone}/machineTypes/{machineType}". Must belong to the enclosing environment's project and region/zone.`,
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
										Description: `The disk size in GB used for node VMs. Minimum size is 20GB. If unspecified, defaults to 100GB. Cannot be updated.`,
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
										Description: `The set of Google API scopes to be made available on all node VMs. Cannot be updated. If empty, defaults to ["https://www.googleapis.com/auth/cloud-platform"].`,
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
									"tags": {
										Type:     schema.TypeSet,
										Optional: true,
										ForceNew: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Set:         schema.HashString,
										Description: `The list of instance tags applied to all node VMs. Tags are used to identify valid sources or targets for network firewalls. Each tag within the list must comply with RFC1035. Cannot be updated.`,
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
													Type:        schema.TypeBool,
													Required:    true,
													ForceNew:    true,
													Description: `Whether or not to enable Alias IPs in the GKE cluster. If true, a VPC-native cluster is created. Defaults to true if the ip_allocation_block is present in config.`,
												},
												"cluster_secondary_range_name": {
													Type:          schema.TypeString,
													Optional:      true,
													ForceNew:      true,
													Description:   `The name of the cluster's secondary range used to allocate IP addresses to pods. Specify either cluster_secondary_range_name or cluster_ipv4_cidr_block but not both. This field is applicable only when use_ip_aliases is true.`,
													ConflictsWith: []string{"config.0.node_config.0.ip_allocation_policy.0.cluster_ipv4_cidr_block"},
												},
												"services_secondary_range_name": {
													Type:          schema.TypeString,
													Optional:      true,
													ForceNew:      true,
													Description:   `The name of the services' secondary range used to allocate IP addresses to the cluster. Specify either services_secondary_range_name or services_ipv4_cidr_block but not both. This field is applicable only when use_ip_aliases is true.`,
													ConflictsWith: []string{"config.0.node_config.0.ip_allocation_policy.0.services_ipv4_cidr_block"},
												},
												"cluster_ipv4_cidr_block": {
													Type:             schema.TypeString,
													Optional:         true,
													ForceNew:         true,
													Description:      `The IP address range used to allocate IP addresses to pods in the cluster. Set to blank to have GKE choose a range with the default size. Set to /netmask (e.g. /14) to have GKE choose a range with a specific netmask. Set to a CIDR notation (e.g. 10.96.0.0/14) from the RFC-1918 private networks (e.g. 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16) to pick a specific range to use. Specify either cluster_secondary_range_name or cluster_ipv4_cidr_block but not both.`,
													DiffSuppressFunc: cidrOrSizeDiffSuppress,
													ConflictsWith:    []string{"config.0.node_config.0.ip_allocation_policy.0.cluster_secondary_range_name"},
												},
												"services_ipv4_cidr_block": {
													Type:             schema.TypeString,
													Optional:         true,
													ForceNew:         true,
													Description:      `The IP address range used to allocate IP addresses in this cluster. Set to blank to have GKE choose a range with the default size. Set to /netmask (e.g. /14) to have GKE choose a range with a specific netmask. Set to a CIDR notation (e.g. 10.96.0.0/14) from the RFC-1918 private networks (e.g. 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16) to pick a specific range to use. Specify either services_secondary_range_name or services_ipv4_cidr_block but not both.`,
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
										AtLeastOneOf:     composerSoftwareConfigKeys,
										ValidateFunc:     validateRegexp(composerEnvironmentVersionRegexp),
										DiffSuppressFunc: composerImageVersionDiffSuppress,
										Description:      `The version of the software running in the environment. This encapsulates both the version of Cloud Composer functionality and the version of Apache Airflow. It must match the regular expression composer-[0-9]+\.[0-9]+(\.[0-9]+)?-airflow-[0-9]+\.[0-9]+(\.[0-9]+.*)?. The Cloud Composer portion of the version is a semantic version. The portion of the image version following 'airflow-' is an official Apache Airflow repository release name. See documentation for allowed release names.`,
									},
									"python_version": {
										Type:         schema.TypeString,
										Optional:     true,
										AtLeastOneOf: composerSoftwareConfigKeys,
										Computed:     true,
										ForceNew:     true,
										Description:  `The major version of Python used to run the Apache Airflow scheduler, worker, and webserver processes. Can be set to '2' or '3'. If not specified, the default is '2'. Cannot be updated.`,
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
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
										AtLeastOneOf: []string{
											"config.0.private_environment_config.0.enable_private_endpoint",
											"config.0.private_environment_config.0.master_ipv4_cidr_block",
											"config.0.private_environment_config.0.cloud_sql_ipv4_cidr_block",
											"config.0.private_environment_config.0.web_server_ipv4_cidr_block",
										},
										ForceNew:    true,
										Description: `If true, access to the public endpoint of the GKE cluster is denied.`,
									},
									"master_ipv4_cidr_block": {
										Type:     schema.TypeString,
										Optional: true,
										AtLeastOneOf: []string{
											"config.0.private_environment_config.0.enable_private_endpoint",
											"config.0.private_environment_config.0.master_ipv4_cidr_block",
											"config.0.private_environment_config.0.cloud_sql_ipv4_cidr_block",
											"config.0.private_environment_config.0.web_server_ipv4_cidr_block",
										},
										ForceNew:    true,
										Default:     "172.16.0.0/28",
										Description: `The IP range in CIDR notation to use for the hosted master network. This range is used for assigning internal IP addresses to the cluster master or set of masters and to the internal load balancer virtual IP. This range must not overlap with any other ranges in use within the cluster's network. If left blank, the default value of '172.16.0.0/28' is used.`,
									},
									"web_server_ipv4_cidr_block": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										AtLeastOneOf: []string{
											"config.0.private_environment_config.0.enable_private_endpoint",
											"config.0.private_environment_config.0.master_ipv4_cidr_block",
											"config.0.private_environment_config.0.cloud_sql_ipv4_cidr_block",
											"config.0.private_environment_config.0.web_server_ipv4_cidr_block",
										},
										ForceNew:    true,
										Description: `The CIDR block from which IP range for web server will be reserved. Needs to be disjoint from master_ipv4_cidr_block and cloud_sql_ipv4_cidr_block.`,
									},
									"cloud_sql_ipv4_cidr_block": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										AtLeastOneOf: []string{
											"config.0.private_environment_config.0.enable_private_endpoint",
											"config.0.private_environment_config.0.master_ipv4_cidr_block",
											"config.0.private_environment_config.0.cloud_sql_ipv4_cidr_block",
											"config.0.private_environment_config.0.web_server_ipv4_cidr_block",
										},
										ForceNew:    true,
										Description: `The CIDR block from which IP range in tenant project will be reserved for Cloud SQL. Needs to be disjoint from web_server_ipv4_cidr_block.`,
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
	}
}

func resourceComposerEnvironmentCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

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
	op, err := config.clientComposer.Projects.Locations.Environments.Create(envName.parentName(), env).Do()
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
		config, op, envName.Project, "Creating Environment",
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

	if err := resourceComposerEnvironmentPostCreateUpdate(updateOnlyEnv, d, config); err != nil {
		return err
	}

	return resourceComposerEnvironmentRead(d, meta)
}

func resourceComposerEnvironmentRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	envName, err := resourceComposerEnvironmentName(d, config)
	if err != nil {
		return err
	}

	res, err := config.clientComposer.Projects.Locations.Environments.Get(envName.resourceName()).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("ComposerEnvironment %q", d.Id()))
	}

	// Set from getProject(d)
	if err := d.Set("project", envName.Project); err != nil {
		return fmt.Errorf("Error reading Environment: %s", err)
	}
	// Set from getRegion(d)
	if err := d.Set("region", envName.Region); err != nil {
		return fmt.Errorf("Error reading Environment: %s", err)
	}
	if err := d.Set("name", GetResourceNameFromSelfLink(res.Name)); err != nil {
		return fmt.Errorf("Error reading Environment: %s", err)
	}
	if err := d.Set("config", flattenComposerEnvironmentConfig(res.Config)); err != nil {
		return fmt.Errorf("Error reading Environment: %s", err)
	}
	if err := d.Set("labels", res.Labels); err != nil {
		return fmt.Errorf("Error reading Environment: %s", err)
	}
	return nil
}

func resourceComposerEnvironmentUpdate(d *schema.ResourceData, meta interface{}) error {
	tfConfig := meta.(*Config)

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

		if d.HasChange("config.0.software_config.0.image_version") {
			patchObj := &composer.Environment{
				Config: &composer.EnvironmentConfig{
					SoftwareConfig: &composer.SoftwareConfig{},
				},
			}
			if config != nil && config.SoftwareConfig != nil {
				patchObj.Config.SoftwareConfig.ImageVersion = config.SoftwareConfig.ImageVersion
			}
			err = resourceComposerEnvironmentPatchField("config.softwareConfig.imageVersion", patchObj, d, tfConfig)
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

			err = resourceComposerEnvironmentPatchField("config.softwareConfig.airflowConfigOverrides", patchObj, d, tfConfig)
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

			err = resourceComposerEnvironmentPatchField("config.softwareConfig.envVariables", patchObj, d, tfConfig)
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

			err = resourceComposerEnvironmentPatchField("config.softwareConfig.pypiPackages", patchObj, d, tfConfig)
			if err != nil {
				return err
			}
		}

		if d.HasChange("config.0.node_count") {
			patchObj := &composer.Environment{Config: &composer.EnvironmentConfig{}}
			if config != nil {
				patchObj.Config.NodeCount = config.NodeCount
			}
			err = resourceComposerEnvironmentPatchField("config.nodeCount", patchObj, d, tfConfig)
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
			err = resourceComposerEnvironmentPatchField("config.webServerNetworkAccessControl", patchObj, d, tfConfig)
			if err != nil {
				return err
			}
		}

	}

	if d.HasChange("labels") {
		patchEnv := &composer.Environment{Labels: expandLabels(d)}
		err := resourceComposerEnvironmentPatchField("labels", patchEnv, d, tfConfig)
		if err != nil {
			return err
		}
	}

	d.Partial(false)
	return resourceComposerEnvironmentRead(d, tfConfig)
}

func resourceComposerEnvironmentPostCreateUpdate(updateEnv *composer.Environment, d *schema.ResourceData, cfg *Config) error {
	if updateEnv == nil {
		return nil
	}

	d.Partial(true)

	if updateEnv.Config != nil && updateEnv.Config.SoftwareConfig != nil && len(updateEnv.Config.SoftwareConfig.PypiPackages) > 0 {
		log.Printf("[DEBUG] Running post-create update for Environment %q", d.Id())
		err := resourceComposerEnvironmentPatchField("config.softwareConfig.pypiPackages", updateEnv, d, cfg)
		if err != nil {
			return err
		}

		log.Printf("[DEBUG] Finish update to Environment %q post create for update only fields", d.Id())
	}
	d.Partial(false)
	return resourceComposerEnvironmentRead(d, cfg)
}

func resourceComposerEnvironmentPatchField(updateMask string, env *composer.Environment, d *schema.ResourceData, config *Config) error {
	envJson, _ := env.MarshalJSON()
	log.Printf("[DEBUG] Updating Environment %q (updateMask = %q): %s", d.Id(), updateMask, string(envJson))
	envName, err := resourceComposerEnvironmentName(d, config)
	if err != nil {
		return err
	}

	op, err := config.clientComposer.Projects.Locations.Environments.
		Patch(envName.resourceName(), env).
		UpdateMask(updateMask).Do()
	if err != nil {
		return err
	}

	waitErr := composerOperationWaitTime(
		config, op, envName.Project, "Updating newly created Environment",
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

	envName, err := resourceComposerEnvironmentName(d, config)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting Environment %q", d.Id())
	op, err := config.clientComposer.Projects.Locations.Environments.Delete(envName.resourceName()).Do()
	if err != nil {
		return err
	}

	err = composerOperationWaitTime(
		config, op, envName.Project, "Deleting Environment",
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
	return []interface{}{transformed}
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

	return transformed, nil
}

func expandComposerEnvironmentConfigNodeCount(v interface{}, d *schema.ResourceData, config *Config) (int64, error) {
	if v == nil {
		return 0, nil
	}
	return int64(v.(int)), nil
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
	if len(l) == 0 {
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
		if requiredZone == "" {
			return "", err
		}

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
	log.Printf("[WARNING] Creation operation for Composer Environment %q failed, check Environment isn't still running", id)
	// Try to get possible created but invalid environment.
	env, err := config.clientComposer.Projects.Locations.Environments.Get(envName.resourceName()).Do()
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
	op, err := config.clientComposer.Projects.Locations.Environments.Delete(envName.resourceName()).Do()
	if err != nil {
		return fmt.Errorf("Could not delete the invalid created environment with state %q: %s", env.State, err)
	}

	waitErr := composerOperationWaitTime(
		config, op, envName.Project,
		fmt.Sprintf("Deleting invalid created Environment with state %q", env.State),
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
	if oldVersions == nil || len(oldVersions) < 3 {
		// Somehow one of the versions didn't match the regexp or didn't
		// have values in the capturing groups. In that case, fall back to
		// an equality check.
		log.Printf("[WARN] Composer version didn't match regexp: %s", old)
		return old == new
	}
	if newVersions == nil || len(newVersions) < 3 {
		// Somehow one of the versions didn't match the regexp or didn't
		// have values in the capturing groups. In that case, fall back to
		// an equality check.
		log.Printf("[WARN] Composer version didn't match regexp: %s", new)
		return old == new
	}

	// Check airflow version using the version package to account for
	// diffs like 1.10 and 1.10.0
	eq, err := versionsEqual(oldVersions[2], newVersions[2])
	if err != nil {
		log.Printf("[WARN] Could not parse airflow version, %s", err)
	}
	if !eq {
		return false
	}

	// Check composer version. Assume that "latest" means we should
	// suppress the diff, because we don't have any other way of
	// knowing what the latest version actually is.
	if oldVersions[1] == "latest" || newVersions[1] == "latest" {
		return true
	}
	// If neither version is "latest", check them using the version
	// package like we did for airflow.
	eq, err = versionsEqual(oldVersions[1], newVersions[1])
	if err != nil {
		log.Printf("[WARN] Could not parse composer version, %s", err)
	}
	return eq
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
