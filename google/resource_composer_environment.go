package google

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"google.golang.org/api/composer/v1"
)

const (
	composerEnvironmentEnvVariablesRegexp          = "[a-zA-Z_][a-zA-Z0-9_]*."
	composerEnvironmentReservedAirflowEnvVarRegexp = "AIRFLOW__[A-Z0-9_]+__[A-Z0-9_]+"
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
			Create: schema.DefaultTimeout(3600 * time.Second),
			Update: schema.DefaultTimeout(3600 * time.Second),
			Delete: schema.DefaultTimeout(360 * time.Second),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
							ValidateFunc: validation.IntAtLeast(3),
						},
						"node_config": {
							Type:     schema.TypeList,
							Computed: true,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"zone": {
										Type:             schema.TypeString,
										Computed:         true,
										Optional:         true,
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
								},
							},
						},
						"software_config": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"airflow_config_overrides": {
										Type:     schema.TypeMap,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"pypi_packages": {
										Type:         schema.TypeMap,
										Optional:     true,
										Elem:         &schema.Schema{Type: schema.TypeString},
										ValidateFunc: validateComposerEnvironmentPypiPackages,
									},
									"env_variables": {
										Type:         schema.TypeMap,
										Optional:     true,
										Elem:         &schema.Schema{Type: schema.TypeString},
										ValidateFunc: validateComposerEnvironmentEnvVariables,
									},
									"image_version": {
										Type:     schema.TypeString,
										Computed: true,
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
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
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
	id, err := replaceVars(d, config, "{{project}}/{{region}}/{{name}}")
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

	return resourceComposerEnvironmentRead(d, 272)
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

	// Set from d.GetProject
	if err := d.Set("project", envName.Project); err != nil {
		return fmt.Errorf("Error reading Environment: %s", err)
	}
	// Set from d.GetRegion
	if err := d.Set("region", envName.Region); err != nil {
		return fmt.Errorf("Error reading Environment: %s", err)
	}
	if err := d.Set("name", GetResourceNameFromSelfLink(res.Name)); err != nil {
		return fmt.Errorf("Error reading Environment: %s", err)
	}
	if err := d.Set("config", flattenComposerEnvironmentConfig(res.Config)); err != nil {
		return fmt.Errorf("Error reading Environment: %s", err)
	}
	if err := d.Set("uuid", res.Uuid); err != nil {
		return fmt.Errorf("Error reading Environment: %s", err)
	}
	if err := d.Set("state", res.State); err != nil {
		return fmt.Errorf("Error reading Environment: %s", err)
	}
	if err := d.Set("create_time", res.CreateTime); err != nil {
		return fmt.Errorf("Error reading Environment: %s", err)
	}
	if err := d.Set("update_time", res.UpdateTime); err != nil {
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
	if d.HasChange("config") {
		config, err := expandComposerEnvironmentConfig(d.Get("config"), d, tfConfig)
		if err != nil {
			return err
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
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error waiting to update Environment: %s", waitErr)
	}

	log.Printf("[DEBUG] Finished updating Environment %q (updateMask = %q)", d.Id(), updateMask)

	return resourceComposerEnvironmentRead(d, config)
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
	parseImportId([]string{"projects/(?P<project>[^/]+)/locations/(?P<region>[^/]+)/environments/(?P<name>[^/]+)", "(?P<project>[^/]+)/(?P<region>[^/]+)/(?P<name>[^/]+)", "(?P<name>[^/]+)"}, d, config)

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "{{project}}/{{region}}/{{name}}")
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
	return transformed, nil
}

func expandComposerEnvironmentConfigNodeCount(v interface{}, d *schema.ResourceData, config *Config) (int64, error) {
	if v == nil {
		return 0, nil
	}
	return int64(v.(int)), nil
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

func expandComposerEnvironmentServiceAccount(v interface{}, d *schema.ResourceData, config *Config) (string, error) {
	serviceAccount := v.(string)
	if len(serviceAccount) == 0 {
		return "", nil
	}

	serviceAccountEmail := GetResourceNameFromSelfLink(serviceAccount)
	r := regexp.MustCompile("(" + strings.Join(PossibleServiceAccountNames, "|") + ")")
	if !r.MatchString(serviceAccountEmail) {
		return "", fmt.Errorf("service account %q must match one of: %+v", serviceAccountEmail, PossibleServiceAccountNames)
	}

	return serviceAccountEmail, nil
}

func expandComposerEnvironmentZone(v interface{}, d *schema.ResourceData, config *Config) (string, error) {
	zone := v.(string)
	if len(zone) == 0 {
		return "", nil
	}

	r := regexp.MustCompile("projects/(.+)/zones/(.+)")
	if r.MatchString(zone) {
		return getRelativePath(zone)
	}

	project, err := getProject(d, config)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("projects/%s/zones/%s", project, zone), nil
}

func expandComposerEnvironmentMachineType(v interface{}, d *schema.ResourceData, config *Config, nodeCfgZone interface{}) (string, error) {
	fullMachineType := v.(string)
	if len(fullMachineType) == 0 {
		return "", nil
	}

	var machineType, zone, project string

	tkns := strings.Split(fullMachineType, "/")
	if len(tkns) == 6 && tkns[0] == "projects" {
		// projects/{project}/zones/{zone}/machineTypes/{machineTypeId}
		project = tkns[1]

		// Continue parsing rest of machine_type.
		tkns = tkns[2:]
	}
	if len(tkns) == 4 && tkns[0] == "zones" && tkns[2] == "machineTypes" {
		// zones/{zone}/machineTypes/{machineTypeId}
		zone = tkns[1]
	}
	if len(tkns) > 0 {
		machineType = tkns[len(tkns)-1]
	}

	if len(project) == 0 {
		// Try to infer project from config.
		var err error
		project, err = getProject(d, config)
		if err != nil {
			return "", err
		}
	}

	// Check if zone was given in node_config.
	if nodeCfgZone != nil {
		z := GetResourceNameFromSelfLink(nodeCfgZone.(string))
		if len(z) > 0 {
			if len(zone) > 0 && zone != z {
				// Make sure machine type and zone locations do not conflict.
				return "", fmt.Errorf(`invalid machine_type: "%s" for specified zone "%s" in node_config`, machineType, zone)
			}
			zone = z
		}
	}
	if len(zone) == 0 {
		return "", fmt.Errorf(`zone must be specified for machine_type "%s" as separate zone field or in machine_type`, machineType)
	}
	if len(machineType) == 0 {
		return "", fmt.Errorf(`invalid machine_type "%s"`, fullMachineType)
	}

	return fmt.Sprintf("projects/%s/zones/%s/machineTypes/%s", project, zone, machineType), nil
}

func expandComposerEnvironmentNetwork(v interface{}, d *schema.ResourceData, config *Config) (string, error) {
	network := v.(string)
	if len(network) == 0 {
		return "", nil
	}

	r := regexp.MustCompile("projects/(.+)/networks/(.+)")
	if r.MatchString(network) {
		return getRelativePath(network)
	}

	project, err := getProject(d, config)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("projects/%s/global/networks/%s", project, network), nil
}

func expandComposerEnvironmentSubnetwork(v interface{}, d *schema.ResourceData, config *Config) (string, error) {
	subnet := v.(string)
	if len(subnet) == 0 {
		return "", nil
	}

	r := regexp.MustCompile("projects/(.+)/regions/(.+)/subnetworks/(.+)")
	if r.MatchString(subnet) {
		return getRelativePath(subnet)
	}

	project, err := getProject(d, config)
	if err != nil {
		return "", err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("projects/%s/regions/%s/subnetworks/%s", project, region, subnet), nil
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
	transformed.AirflowConfigOverrides = expandComposerEnvironmentStringMap(original["airflow_config_overrides"], d, config)
	transformed.PypiPackages = expandComposerEnvironmentStringMap(original["pypi_packages"], d, config)
	transformed.EnvVariables = expandComposerEnvironmentStringMap(original["env_variables"], d, config)

	return transformed, nil
}

func expandComposerEnvironmentStringMap(v interface{}, d *schema.ResourceData, config *Config) map[string]string {
	if v == nil {
		return map[string]string{}
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m
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

	envName := d.Get("name").(string)
	if len(envName) == 0 {
		return nil, fmt.Errorf("cannot have empty environment name")
	}
	return &composerEnvironmentName{
		Project:     project,
		Region:      region,
		Environment: envName,
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
