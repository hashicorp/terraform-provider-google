package google

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateGCPName,
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_count": {
							Type:         schema.TypeInt,
							Computed:     true,
							Optional:     true,
							AtLeastOneOf: composerConfigKeys,
							ValidateFunc: validation.IntAtLeast(3),
						},
						"node_config": {
							Type:         schema.TypeList,
							Computed:     true,
							Optional:     true,
							AtLeastOneOf: composerConfigKeys,
							MaxItems:     1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"zone": {
										Type:             schema.TypeString,
										Required:         true,
										ForceNew:         true,
										DiffSuppressFunc: compareSelfLinkOrResourceName,
									},
									"machine_type": {
										Type:             schema.TypeString,
										Computed:         true,
										Optional:         true,
										ForceNew:         true,
										DiffSuppressFunc: compareSelfLinkOrResourceName,
									},
									"network": {
										Type:             schema.TypeString,
										Computed:         true,
										Optional:         true,
										ForceNew:         true,
										DiffSuppressFunc: compareSelfLinkOrResourceName,
									},
									"subnetwork": {
										Type:             schema.TypeString,
										Optional:         true,
										ForceNew:         true,
										DiffSuppressFunc: compareSelfLinkOrResourceName,
									},
									"disk_size_gb": {
										Type:     schema.TypeInt,
										Computed: true,
										Optional: true,
										ForceNew: true,
									},
									"oauth_scopes": {
										Type:     schema.TypeSet,
										Computed: true,
										Optional: true,
										ForceNew: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Set: schema.HashString,
									},
									"service_account": {
										Type:             schema.TypeString,
										Computed:         true,
										Optional:         true,
										ForceNew:         true,
										ValidateFunc:     validateServiceAccountRelativeNameOrEmail,
										DiffSuppressFunc: compareServiceAccountEmailToLink,
									},
									"tags": {
										Type:     schema.TypeSet,
										Optional: true,
										ForceNew: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Set: schema.HashString,
									},
									"ip_allocation_policy": {
										Type:       schema.TypeList,
										Optional:   true,
										Computed:   true,
										ForceNew:   true,
										ConfigMode: schema.SchemaConfigModeAttr,
										MaxItems:   1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"use_ip_aliases": {
													Type:     schema.TypeBool,
													Required: true,
													ForceNew: true,
												},
												"cluster_secondary_range_name": {
													Type:          schema.TypeString,
													Optional:      true,
													ForceNew:      true,
													ConflictsWith: []string{"config.0.node_config.0.ip_allocation_policy.0.cluster_ipv4_cidr_block"},
												},
												"services_secondary_range_name": {
													Type:          schema.TypeString,
													Optional:      true,
													ForceNew:      true,
													ConflictsWith: []string{"config.0.node_config.0.ip_allocation_policy.0.services_ipv4_cidr_block"},
												},
												"cluster_ipv4_cidr_block": {
													Type:             schema.TypeString,
													Optional:         true,
													ForceNew:         true,
													DiffSuppressFunc: cidrOrSizeDiffSuppress,
													ConflictsWith:    []string{"config.0.node_config.0.ip_allocation_policy.0.cluster_secondary_range_name"},
												},
												"services_ipv4_cidr_block": {
													Type:             schema.TypeString,
													Optional:         true,
													ForceNew:         true,
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
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"airflow_config_overrides": {
										Type:         schema.TypeMap,
										Optional:     true,
										AtLeastOneOf: composerSoftwareConfigKeys,
										Elem:         &schema.Schema{Type: schema.TypeString},
									},
									"pypi_packages": {
										Type:         schema.TypeMap,
										Optional:     true,
										AtLeastOneOf: composerSoftwareConfigKeys,
										Elem:         &schema.Schema{Type: schema.TypeString},
										ValidateFunc: validateComposerEnvironmentPypiPackages,
									},
									"env_variables": {
										Type:         schema.TypeMap,
										Optional:     true,
										AtLeastOneOf: composerSoftwareConfigKeys,
										Elem:         &schema.Schema{Type: schema.TypeString},
										ValidateFunc: validateComposerEnvironmentEnvVariables,
									},
									"image_version": {
										Type:             schema.TypeString,
										Computed:         true,
										Optional:         true,
										AtLeastOneOf:     composerSoftwareConfigKeys,
										ValidateFunc:     validateRegexp(composerEnvironmentVersionRegexp),
										DiffSuppressFunc: composerImageVersionDiffSuppress,
									},
									"python_version": {
										Type:         schema.TypeString,
										Optional:     true,
										AtLeastOneOf: composerSoftwareConfigKeys,
										Computed:     true,
										ForceNew:     true,
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
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enable_private_endpoint": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
										AtLeastOneOf: []string{
											"config.0.private_environment_config.0.enable_private_endpoint",
											"config.0.private_environment_config.0.master_ipv4_cidr_block",
										},
										ForceNew: true,
									},
									"master_ipv4_cidr_block": {
										Type:     schema.TypeString,
										Optional: true,
										AtLeastOneOf: []string{
											"config.0.private_environment_config.0.enable_private_endpoint",
											"config.0.private_environment_config.0.master_ipv4_cidr_block",
										},
										ForceNew: true,
										Default:  "172.16.0.0/28",
									},
								},
							},
						},
						"airflow_uri": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dag_gcs_prefix": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gke_cluster": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
		config.clientComposer, op, envName.Project, "Creating Environment",
		int(d.Timeout(schema.TimeoutCreate).Minutes()))

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
			d.SetPartial("config")
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
			d.SetPartial("config")
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
			d.SetPartial("config")
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
			d.SetPartial("config")
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
			d.SetPartial("config")
		}
	}

	if d.HasChange("labels") {
		patchEnv := &composer.Environment{Labels: expandLabels(d)}
		err := resourceComposerEnvironmentPatchField("labels", patchEnv, d, tfConfig)
		if err != nil {
			return err
		}
		d.SetPartial("labels")
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
		d.SetPartial("config")
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
		config.clientComposer, op, envName.Project, "Updating newly created Environment",
		int(d.Timeout(schema.TimeoutCreate).Minutes()))
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
		config.clientComposer, op, envName.Project, "Deleting Environment",
		int(d.Timeout(schema.TimeoutDelete).Minutes()))
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
		config.clientComposer, op, envName.Project,
		fmt.Sprintf("Deleting invalid created Environment with state %q", env.State),
		int(d.Timeout(schema.TimeoutCreate).Minutes()))
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
