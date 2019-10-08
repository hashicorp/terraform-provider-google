package google

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	dataproc "google.golang.org/api/dataproc/v1beta2"
)

var resolveDataprocImageVersion = regexp.MustCompile(`(?P<Major>[^\s.-]+)\.(?P<Minor>[^\s.-]+)(?:\.(?P<Subminor>[^\s.-]+))?(?:\-(?P<Distr>[^\s.-]+))?`)

func resourceDataprocCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceDataprocClusterCreate,
		Read:   resourceDataprocClusterRead,
		Update: resourceDataprocClusterUpdate,
		Delete: resourceDataprocClusterDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)

					if len(value) > 55 {
						errors = append(errors, fmt.Errorf(
							"%q cannot be longer than 55 characters", k))
					}
					if !regexp.MustCompile("^[a-z0-9-]+$").MatchString(value) {
						errors = append(errors, fmt.Errorf(
							"%q can only contain lowercase letters, numbers and hyphens", k))
					}
					if !regexp.MustCompile("^[a-z]").MatchString(value) {
						errors = append(errors, fmt.Errorf(
							"%q must start with a letter", k))
					}
					if !regexp.MustCompile("[a-z0-9]$").MatchString(value) {
						errors = append(errors, fmt.Errorf(
							"%q must end with a number or a letter", k))
					}
					return
				},
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "global",
				ForceNew: true,
			},

			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				// GCP automatically adds two labels
				//    'goog-dataproc-cluster-uuid'
				//    'goog-dataproc-cluster-name'
				Computed: true,
			},

			"cluster_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"delete_autogen_bucket": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
							Removed: "If you need a bucket that can be deleted, please create" +
								"a new one and set the `staging_bucket` field",
						},

						"staging_bucket": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						// If the user does not specify a staging bucket, GCP will allocate one automatically.
						// The staging_bucket field provides a way for the user to supply their own
						// staging bucket. The bucket field is purely a computed field which details
						// the definitive bucket allocated and in use (either the user supplied one via
						// staging_bucket, or the GCP generated one)
						"bucket": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"gce_cluster_config": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{

									"zone": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},

									"network": {
										Type:             schema.TypeString,
										Optional:         true,
										Computed:         true,
										ForceNew:         true,
										ConflictsWith:    []string{"cluster_config.0.gce_cluster_config.0.subnetwork"},
										DiffSuppressFunc: compareSelfLinkOrResourceName,
									},

									"subnetwork": {
										Type:             schema.TypeString,
										Optional:         true,
										ForceNew:         true,
										ConflictsWith:    []string{"cluster_config.0.gce_cluster_config.0.network"},
										DiffSuppressFunc: compareSelfLinkOrResourceName,
									},

									"tags": {
										Type:     schema.TypeSet,
										Optional: true,
										ForceNew: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},

									"service_account": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},

									"service_account_scopes": {
										Type:     schema.TypeSet,
										Optional: true,
										Computed: true,
										ForceNew: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
											StateFunc: func(v interface{}) string {
												return canonicalizeServiceScope(v.(string))
											},
										},
										Set: stringScopeHashcode,
									},

									"internal_ip_only": {
										Type:     schema.TypeBool,
										Optional: true,
										ForceNew: true,
										Default:  false,
									},

									"metadata": {
										Type:     schema.TypeMap,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
										ForceNew: true,
									},
								},
							},
						},

						"master_config": instanceConfigSchema(),
						"worker_config": instanceConfigSchema(),
						// preemptible_worker_config has a slightly different config
						"preemptible_worker_config": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"num_instances": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},

									// API does not honour this if set ...
									// It always uses whatever is specified for the worker_config
									// "machine_type": { ... }
									"disk_config": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										MaxItems: 1,

										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"num_local_ssds": {
													Type:     schema.TypeInt,
													Optional: true,
													Computed: true,
													ForceNew: true,
												},

												"boot_disk_size_gb": {
													Type:         schema.TypeInt,
													Optional:     true,
													Computed:     true,
													ForceNew:     true,
													ValidateFunc: validation.IntAtLeast(10),
												},

												"boot_disk_type": {
													Type:         schema.TypeString,
													Optional:     true,
													ForceNew:     true,
													ValidateFunc: validation.StringInSlice([]string{"pd-standard", "pd-ssd", ""}, false),
													Default:      "pd-standard",
												},
											},
										},
									},

									"instance_names": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
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
									"image_version": {
										Type:             schema.TypeString,
										Optional:         true,
										Computed:         true,
										ForceNew:         true,
										DiffSuppressFunc: dataprocImageVersionDiffSuppress,
									},

									"override_properties": {
										Type:     schema.TypeMap,
										Optional: true,
										ForceNew: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},

									"properties": {
										Type:     schema.TypeMap,
										Computed: true,
									},

									// We have two versions of the properties field here because by default
									// dataproc will set a number of default properties for you out of the
									// box. If you want to override one or more, if we only had one field,
									// you would need to add in all these values as well otherwise you would
									// get a diff. To make this easier, 'properties' simply contains the computed
									// values (including overrides) for all properties, whilst override_properties
									// is only for properties the user specifically wants to override. If nothing
									// is overridden, this will be empty.

									"optional_components": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
											ValidateFunc: validation.StringInSlice([]string{"COMPONENT_UNSPECIFIED", "ANACONDA", "DRUID", "HIVE_WEBHCAT",
												"JUPYTER", "KERBEROS", "PRESTO", "ZEPPELIN", "ZOOKEEPER"}, false),
										},
									},
								},
							},
						},

						"initialization_action": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"script": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},

									"timeout_sec": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  300,
										ForceNew: true,
									},
								},
							},
						},
						"encryption_config": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kms_key_name": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func instanceConfigSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"num_instances": {
					Type:     schema.TypeInt,
					Optional: true,
					Computed: true,
				},

				"image_uri": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},

				"machine_type": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},

				"disk_config": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					MaxItems: 1,

					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"num_local_ssds": {
								Type:     schema.TypeInt,
								Optional: true,
								Computed: true,
								ForceNew: true,
							},

							"boot_disk_size_gb": {
								Type:         schema.TypeInt,
								Optional:     true,
								Computed:     true,
								ForceNew:     true,
								ValidateFunc: validation.IntAtLeast(10),
							},

							"boot_disk_type": {
								Type:         schema.TypeString,
								Optional:     true,
								ForceNew:     true,
								ValidateFunc: validation.StringInSlice([]string{"pd-standard", "pd-ssd", ""}, false),
								Default:      "pd-standard",
							},
						},
					},
				},

				// Note: preemptible workers don't support accelerators
				"accelerators": {
					Type:     schema.TypeSet,
					Optional: true,
					ForceNew: true,
					Elem:     acceleratorsSchema(),
				},

				"instance_names": {
					Type:     schema.TypeList,
					Computed: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
			},
		},
	}
}

// We need to pull accelerators' schema out so we can use it to make a set hash func
func acceleratorsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"accelerator_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"accelerator_count": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceDataprocClusterCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region := d.Get("region").(string)
	cluster := &dataproc.Cluster{
		ClusterName: d.Get("name").(string),
		ProjectId:   project,
	}

	cluster.Config, err = expandClusterConfig(d, config)
	if err != nil {
		return err
	}

	if _, ok := d.GetOk("labels"); ok {
		cluster.Labels = expandLabels(d)
	}

	// Checking here caters for the case where the user does not specify cluster_config
	// at all, as well where it is simply missing from the gce_cluster_config
	if region == "global" && cluster.Config.GceClusterConfig.ZoneUri == "" {
		return errors.New("zone is mandatory when region is set to 'global'")
	}

	// Create the cluster
	op, err := config.clientDataprocBeta.Projects.Regions.Clusters.Create(
		project, region, cluster).Do()
	if err != nil {
		return fmt.Errorf("Error creating Dataproc cluster: %s", err)
	}

	d.SetId(cluster.ClusterName)

	// Wait until it's created
	timeoutInMinutes := int(d.Timeout(schema.TimeoutCreate).Minutes())
	waitErr := dataprocClusterOperationWait(config, op, "creating Dataproc cluster", timeoutInMinutes)
	if waitErr != nil {
		// The resource didn't actually create
		// Note that we do not remove the ID here - this resource tends to leave
		// partially created clusters behind, so we'll let the next Read remove
		// it.
		return waitErr
	}

	log.Printf("[INFO] Dataproc cluster %s has been created", cluster.ClusterName)
	return resourceDataprocClusterRead(d, meta)

}

func expandClusterConfig(d *schema.ResourceData, config *Config) (*dataproc.ClusterConfig, error) {
	conf := &dataproc.ClusterConfig{
		// SDK requires GceClusterConfig to be specified,
		// even if no explicit values specified
		GceClusterConfig: &dataproc.GceClusterConfig{},
	}

	if v, ok := d.GetOk("cluster_config"); ok {
		confs := v.([]interface{})
		if (len(confs)) == 0 {
			return conf, nil
		}
	}

	if v, ok := d.GetOk("cluster_config.0.staging_bucket"); ok {
		conf.ConfigBucket = v.(string)
	}

	c, err := expandGceClusterConfig(d, config)
	if err != nil {
		return nil, err
	}
	conf.GceClusterConfig = c

	if cfg, ok := configOptions(d, "cluster_config.0.software_config"); ok {
		conf.SoftwareConfig = expandSoftwareConfig(cfg)
	}

	if v, ok := d.GetOk("cluster_config.0.initialization_action"); ok {
		conf.InitializationActions = expandInitializationActions(v)
	}

	if cfg, ok := configOptions(d, "cluster_config.0.encryption_config"); ok {
		conf.EncryptionConfig = expandEncryptionConfig(cfg)
	}

	if cfg, ok := configOptions(d, "cluster_config.0.master_config"); ok {
		log.Println("[INFO] got master_config")
		conf.MasterConfig = expandInstanceGroupConfig(cfg)
	}

	if cfg, ok := configOptions(d, "cluster_config.0.worker_config"); ok {
		log.Println("[INFO] got worker config")
		conf.WorkerConfig = expandInstanceGroupConfig(cfg)
	}

	if cfg, ok := configOptions(d, "cluster_config.0.preemptible_worker_config"); ok {
		log.Println("[INFO] got preemptible worker config")
		conf.SecondaryWorkerConfig = expandPreemptibleInstanceGroupConfig(cfg)
		if conf.SecondaryWorkerConfig.NumInstances > 0 {
			conf.SecondaryWorkerConfig.IsPreemptible = true
		}
	}
	return conf, nil
}

func expandGceClusterConfig(d *schema.ResourceData, config *Config) (*dataproc.GceClusterConfig, error) {
	conf := &dataproc.GceClusterConfig{}

	v, ok := d.GetOk("cluster_config.0.gce_cluster_config")
	if !ok {
		return conf, nil
	}
	cfg := v.([]interface{})[0].(map[string]interface{})

	if v, ok := cfg["zone"]; ok {
		conf.ZoneUri = v.(string)
	}
	if v, ok := cfg["network"]; ok {
		nf, err := ParseNetworkFieldValue(v.(string), d, config)
		if err != nil {
			return nil, fmt.Errorf("cannot determine self_link for network %q: %s", v, err)
		}

		conf.NetworkUri = nf.RelativeLink()
	}
	if v, ok := cfg["subnetwork"]; ok {
		snf, err := ParseSubnetworkFieldValue(v.(string), d, config)
		if err != nil {
			return nil, fmt.Errorf("cannot determine self_link for subnetwork %q: %s", v, err)
		}

		conf.SubnetworkUri = snf.RelativeLink()
	}
	if v, ok := cfg["tags"]; ok {
		conf.Tags = convertStringSet(v.(*schema.Set))
	}
	if v, ok := cfg["service_account"]; ok {
		conf.ServiceAccount = v.(string)
	}
	if scopes, ok := cfg["service_account_scopes"]; ok {
		scopesSet := scopes.(*schema.Set)
		scopes := make([]string, scopesSet.Len())
		for i, scope := range scopesSet.List() {
			scopes[i] = canonicalizeServiceScope(scope.(string))
		}
		conf.ServiceAccountScopes = scopes
	}
	if v, ok := cfg["internal_ip_only"]; ok {
		conf.InternalIpOnly = v.(bool)
	}
	if v, ok := cfg["metadata"]; ok {
		conf.Metadata = convertStringMap(v.(map[string]interface{}))
	}
	return conf, nil
}

func expandSoftwareConfig(cfg map[string]interface{}) *dataproc.SoftwareConfig {
	conf := &dataproc.SoftwareConfig{}
	if v, ok := cfg["override_properties"]; ok {
		m := make(map[string]string)
		for k, val := range v.(map[string]interface{}) {
			m[k] = val.(string)
		}
		conf.Properties = m
	}
	if v, ok := cfg["image_version"]; ok {
		conf.ImageVersion = v.(string)
	}
	if components, ok := cfg["optional_components"]; ok {
		compSet := components.(*schema.Set)
		components := make([]string, compSet.Len())
		for i, component := range compSet.List() {
			components[i] = component.(string)
		}
		conf.OptionalComponents = components
	}
	return conf
}

func expandEncryptionConfig(cfg map[string]interface{}) *dataproc.EncryptionConfig {
	conf := &dataproc.EncryptionConfig{}
	if v, ok := cfg["kms_key_name"]; ok {
		conf.GcePdKmsKeyName = v.(string)
	}
	return conf
}

func expandInitializationActions(v interface{}) []*dataproc.NodeInitializationAction {
	actionList := v.([]interface{})

	actions := []*dataproc.NodeInitializationAction{}
	for _, v1 := range actionList {
		actionItem := v1.(map[string]interface{})
		action := &dataproc.NodeInitializationAction{
			ExecutableFile: actionItem["script"].(string),
		}
		if x, ok := actionItem["timeout_sec"]; ok {
			action.ExecutionTimeout = strconv.Itoa(x.(int)) + "s"
		}
		actions = append(actions, action)
	}

	return actions
}

func expandPreemptibleInstanceGroupConfig(cfg map[string]interface{}) *dataproc.InstanceGroupConfig {
	icg := &dataproc.InstanceGroupConfig{}

	if v, ok := cfg["num_instances"]; ok {
		icg.NumInstances = int64(v.(int))
	}
	if dc, ok := cfg["disk_config"]; ok {
		d := dc.([]interface{})
		if len(d) > 0 {
			dcfg := d[0].(map[string]interface{})
			icg.DiskConfig = &dataproc.DiskConfig{}

			if v, ok := dcfg["boot_disk_size_gb"]; ok {
				icg.DiskConfig.BootDiskSizeGb = int64(v.(int))
			}
			if v, ok := dcfg["num_local_ssds"]; ok {
				icg.DiskConfig.NumLocalSsds = int64(v.(int))
			}
			if v, ok := dcfg["boot_disk_type"]; ok {
				icg.DiskConfig.BootDiskType = v.(string)
			}
		}
	}
	return icg
}

func expandInstanceGroupConfig(cfg map[string]interface{}) *dataproc.InstanceGroupConfig {
	icg := &dataproc.InstanceGroupConfig{}

	if v, ok := cfg["num_instances"]; ok {
		icg.NumInstances = int64(v.(int))
	}
	if v, ok := cfg["machine_type"]; ok {
		icg.MachineTypeUri = GetResourceNameFromSelfLink(v.(string))
	}
	if v, ok := cfg["image_uri"]; ok {
		icg.ImageUri = v.(string)
	}

	if dc, ok := cfg["disk_config"]; ok {
		d := dc.([]interface{})
		if len(d) > 0 {
			dcfg := d[0].(map[string]interface{})
			icg.DiskConfig = &dataproc.DiskConfig{}

			if v, ok := dcfg["boot_disk_size_gb"]; ok {
				icg.DiskConfig.BootDiskSizeGb = int64(v.(int))
			}
			if v, ok := dcfg["num_local_ssds"]; ok {
				icg.DiskConfig.NumLocalSsds = int64(v.(int))
			}
			if v, ok := dcfg["boot_disk_type"]; ok {
				icg.DiskConfig.BootDiskType = v.(string)
			}
		}
	}

	icg.Accelerators = expandAccelerators(cfg["accelerators"].(*schema.Set).List())
	return icg
}

func expandAccelerators(configured []interface{}) []*dataproc.AcceleratorConfig {
	accelerators := make([]*dataproc.AcceleratorConfig, 0, len(configured))
	for _, raw := range configured {
		data := raw.(map[string]interface{})
		accelerator := dataproc.AcceleratorConfig{
			AcceleratorTypeUri: data["accelerator_type"].(string),
			AcceleratorCount:   int64(data["accelerator_count"].(int)),
		}

		accelerators = append(accelerators, &accelerator)
	}

	return accelerators
}

func resourceDataprocClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region := d.Get("region").(string)
	clusterName := d.Get("name").(string)
	timeoutInMinutes := int(d.Timeout(schema.TimeoutUpdate).Minutes())

	cluster := &dataproc.Cluster{
		ClusterName: clusterName,
		ProjectId:   project,
		Config:      &dataproc.ClusterConfig{},
	}

	updMask := []string{}

	if d.HasChange("labels") {
		v := d.Get("labels")
		m := make(map[string]string)
		for k, val := range v.(map[string]interface{}) {
			m[k] = val.(string)
		}
		cluster.Labels = m

		updMask = append(updMask, "labels")
	}

	if d.HasChange("cluster_config.0.worker_config.0.num_instances") {
		desiredNumWorks := d.Get("cluster_config.0.worker_config.0.num_instances").(int)
		cluster.Config.WorkerConfig = &dataproc.InstanceGroupConfig{
			NumInstances: int64(desiredNumWorks),
		}

		updMask = append(updMask, "config.worker_config.num_instances")
	}

	if d.HasChange("cluster_config.0.preemptible_worker_config.0.num_instances") {
		desiredNumWorks := d.Get("cluster_config.0.preemptible_worker_config.0.num_instances").(int)
		cluster.Config.SecondaryWorkerConfig = &dataproc.InstanceGroupConfig{
			NumInstances: int64(desiredNumWorks),
		}

		updMask = append(updMask, "config.secondary_worker_config.num_instances")
	}

	if len(updMask) > 0 {
		patch := config.clientDataprocBeta.Projects.Regions.Clusters.Patch(
			project, region, clusterName, cluster)
		op, err := patch.UpdateMask(strings.Join(updMask, ",")).Do()
		if err != nil {
			return err
		}

		// Wait until it's updated
		waitErr := dataprocClusterOperationWait(config, op, "updating Dataproc cluster ", timeoutInMinutes)
		if waitErr != nil {
			return waitErr
		}

		log.Printf("[INFO] Dataproc cluster %s has been updated ", d.Id())
	}

	return resourceDataprocClusterRead(d, meta)
}

func resourceDataprocClusterRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region := d.Get("region").(string)
	clusterName := d.Get("name").(string)

	cluster, err := config.clientDataprocBeta.Projects.Regions.Clusters.Get(
		project, region, clusterName).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Dataproc Cluster %q", clusterName))
	}

	d.Set("name", cluster.ClusterName)
	d.Set("project", project)
	d.Set("region", region)
	d.Set("labels", cluster.Labels)

	cfg, err := flattenClusterConfig(d, cluster.Config)
	if err != nil {
		return err
	}

	err = d.Set("cluster_config", cfg)
	if err != nil {
		return err
	}
	return nil
}

func flattenClusterConfig(d *schema.ResourceData, cfg *dataproc.ClusterConfig) ([]map[string]interface{}, error) {

	data := map[string]interface{}{
		"staging_bucket": d.Get("cluster_config.0.staging_bucket").(string),

		"bucket":                    cfg.ConfigBucket,
		"gce_cluster_config":        flattenGceClusterConfig(d, cfg.GceClusterConfig),
		"software_config":           flattenSoftwareConfig(d, cfg.SoftwareConfig),
		"master_config":             flattenInstanceGroupConfig(d, cfg.MasterConfig),
		"worker_config":             flattenInstanceGroupConfig(d, cfg.WorkerConfig),
		"preemptible_worker_config": flattenPreemptibleInstanceGroupConfig(d, cfg.SecondaryWorkerConfig),
		"encryption_config":         flattenEncryptionConfig(d, cfg.EncryptionConfig),
	}

	if len(cfg.InitializationActions) > 0 {
		val, err := flattenInitializationActions(cfg.InitializationActions)
		if err != nil {
			return nil, err
		}
		data["initialization_action"] = val
	}
	return []map[string]interface{}{data}, nil
}

func flattenSoftwareConfig(d *schema.ResourceData, sc *dataproc.SoftwareConfig) []map[string]interface{} {
	data := map[string]interface{}{
		"image_version":       sc.ImageVersion,
		"optional_components": sc.OptionalComponents,
		"properties":          sc.Properties,
		"override_properties": d.Get("cluster_config.0.software_config.0.override_properties").(map[string]interface{}),
	}

	return []map[string]interface{}{data}
}

func flattenEncryptionConfig(d *schema.ResourceData, ec *dataproc.EncryptionConfig) []map[string]interface{} {
	if ec == nil {
		return nil
	}

	data := map[string]interface{}{
		"kms_key_name": ec.GcePdKmsKeyName,
	}

	return []map[string]interface{}{data}
}

func flattenAccelerators(accelerators []*dataproc.AcceleratorConfig) interface{} {
	acceleratorsTypeSet := schema.NewSet(schema.HashResource(acceleratorsSchema()), []interface{}{})
	for _, accelerator := range accelerators {
		data := map[string]interface{}{
			"accelerator_type":  GetResourceNameFromSelfLink(accelerator.AcceleratorTypeUri),
			"accelerator_count": int(accelerator.AcceleratorCount),
		}

		acceleratorsTypeSet.Add(data)
	}

	return acceleratorsTypeSet
}

func flattenInitializationActions(nia []*dataproc.NodeInitializationAction) ([]map[string]interface{}, error) {

	actions := []map[string]interface{}{}
	for _, v := range nia {
		action := map[string]interface{}{
			"script": v.ExecutableFile,
		}
		if len(v.ExecutionTimeout) > 0 {
			tsec, err := extractInitTimeout(v.ExecutionTimeout)
			if err != nil {
				return nil, err
			}
			action["timeout_sec"] = tsec
		}

		actions = append(actions, action)
	}
	return actions, nil

}

func flattenGceClusterConfig(d *schema.ResourceData, gcc *dataproc.GceClusterConfig) []map[string]interface{} {

	gceConfig := map[string]interface{}{
		"tags":             schema.NewSet(schema.HashString, convertStringArrToInterface(gcc.Tags)),
		"service_account":  gcc.ServiceAccount,
		"zone":             GetResourceNameFromSelfLink(gcc.ZoneUri),
		"internal_ip_only": gcc.InternalIpOnly,
		"metadata":         gcc.Metadata,
	}

	if gcc.NetworkUri != "" {
		gceConfig["network"] = gcc.NetworkUri
	}
	if gcc.SubnetworkUri != "" {
		gceConfig["subnetwork"] = gcc.SubnetworkUri
	}
	if len(gcc.ServiceAccountScopes) > 0 {
		gceConfig["service_account_scopes"] = schema.NewSet(stringScopeHashcode, convertStringArrToInterface(gcc.ServiceAccountScopes))
	}

	return []map[string]interface{}{gceConfig}
}

func flattenPreemptibleInstanceGroupConfig(d *schema.ResourceData, icg *dataproc.InstanceGroupConfig) []map[string]interface{} {
	disk := map[string]interface{}{}
	data := map[string]interface{}{}

	if icg != nil {
		data["num_instances"] = icg.NumInstances
		data["instance_names"] = icg.InstanceNames
		if icg.DiskConfig != nil {
			disk["boot_disk_size_gb"] = icg.DiskConfig.BootDiskSizeGb
			disk["num_local_ssds"] = icg.DiskConfig.NumLocalSsds
			disk["boot_disk_type"] = icg.DiskConfig.BootDiskType
		}
	}

	data["disk_config"] = []map[string]interface{}{disk}
	return []map[string]interface{}{data}
}

func flattenInstanceGroupConfig(d *schema.ResourceData, icg *dataproc.InstanceGroupConfig) []map[string]interface{} {
	disk := map[string]interface{}{}
	data := map[string]interface{}{}

	if icg != nil {
		data["num_instances"] = icg.NumInstances
		data["machine_type"] = GetResourceNameFromSelfLink(icg.MachineTypeUri)
		data["image_uri"] = icg.ImageUri
		data["instance_names"] = icg.InstanceNames
		if icg.DiskConfig != nil {
			disk["boot_disk_size_gb"] = icg.DiskConfig.BootDiskSizeGb
			disk["num_local_ssds"] = icg.DiskConfig.NumLocalSsds
			disk["boot_disk_type"] = icg.DiskConfig.BootDiskType
		}

		data["accelerators"] = flattenAccelerators(icg.Accelerators)
	}

	data["disk_config"] = []map[string]interface{}{disk}
	return []map[string]interface{}{data}
}

func extractInitTimeout(t string) (int, error) {
	d, err := time.ParseDuration(t)
	if err != nil {
		return 0, err
	}
	return int(d.Seconds()), nil
}

func resourceDataprocClusterDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region := d.Get("region").(string)
	clusterName := d.Get("name").(string)
	timeoutInMinutes := int(d.Timeout(schema.TimeoutDelete).Minutes())

	log.Printf("[DEBUG] Deleting Dataproc cluster %s", clusterName)
	op, err := config.clientDataprocBeta.Projects.Regions.Clusters.Delete(
		project, region, clusterName).Do()
	if err != nil {
		return err
	}

	// Wait until it's deleted
	waitErr := dataprocClusterOperationWait(config, op, "deleting Dataproc cluster", timeoutInMinutes)
	if waitErr != nil {
		return waitErr
	}

	log.Printf("[INFO] Dataproc cluster %s has been deleted", d.Id())
	d.SetId("")

	return nil
}

func configOptions(d *schema.ResourceData, option string) (map[string]interface{}, bool) {
	if v, ok := d.GetOk(option); ok {
		clist := v.([]interface{})
		if len(clist) == 0 {
			return nil, false
		}

		if clist[0] != nil {
			return clist[0].(map[string]interface{}), true
		}
	}
	return nil, false
}

func dataprocImageVersionDiffSuppress(_, old, new string, _ *schema.ResourceData) bool {
	oldV, err := parseDataprocImageVersion(old)
	if err != nil {
		return false
	}
	newV, err := parseDataprocImageVersion(new)
	if err != nil {
		return false
	}

	if newV.major != oldV.major {
		return false
	}
	if newV.minor != oldV.minor {
		return false
	}
	// Only compare subminor version if set in config version.
	if newV.subminor != "" && newV.subminor != oldV.subminor {
		return false
	}
	// Only compare os if it is set in config version.
	if newV.osName != "" && newV.osName != oldV.osName {
		return false
	}
	return true
}

type dataprocImageVersion struct {
	major    string
	minor    string
	subminor string
	osName   string
}

func parseDataprocImageVersion(version string) (*dataprocImageVersion, error) {
	matches := resolveDataprocImageVersion.FindStringSubmatch(version)
	if len(matches) != 5 {
		return nil, fmt.Errorf("invalid image version %q", version)
	}

	return &dataprocImageVersion{
		major:    matches[1],
		minor:    matches[2],
		subminor: matches[3],
		osName:   matches[4],
	}, nil
}
