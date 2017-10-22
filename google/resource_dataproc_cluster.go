package google

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	"google.golang.org/api/dataproc/v1"
	"google.golang.org/api/googleapi"
)

func resourceDataprocCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceDataprocClusterCreate,
		Read:   resourceDataprocClusterRead,
		Update: resourceDataprocClusterUpdate,
		Delete: resourceDataprocClusterDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
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
				Elem:     schema.TypeString,
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
										Type:          schema.TypeString,
										Optional:      true,
										Computed:      true,
										ForceNew:      true,
										ConflictsWith: []string{"cluster_config.0.gce_cluster_config.0.subnetwork"},
										StateFunc: func(s interface{}) string {
											return extractLastResourceFromUri(s.(string))
										},
									},

									"subnetwork": {
										Type:          schema.TypeString,
										Optional:      true,
										ForceNew:      true,
										ConflictsWith: []string{"cluster_config.0.gce_cluster_config.0.network"},
										StateFunc: func(s interface{}) string {
											return extractLastResourceFromUri(s.(string))
										},
									},

									"tags": {
										Type:     schema.TypeList,
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

												// API does not honour this if set ...
												// It simply ignores it completely
												// "num_local_ssds": { ... }

												"boot_disk_size_gb": {
													Type:         schema.TypeInt,
													Optional:     true,
													Computed:     true,
													ForceNew:     true,
													ValidateFunc: validation.IntAtLeast(10),
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
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},

									"override_properties": {
										Type:     schema.TypeMap,
										Optional: true,
										ForceNew: true,
										Elem:     schema.TypeString,
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
									// is overriden, this will be empty.
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

	cluster.Config = expandClusterConfig(d)
	if _, ok := d.GetOk("labels"); ok {
		cluster.Labels = expandLabels(d)
	}

	// Checking here caters for the case where the user does not specify cluster_config
	// at all, as well where it is simply missing from the gce_cluster_config
	if region == "global" && cluster.Config.GceClusterConfig.ZoneUri == "" {
		return errors.New("zone is mandatory when region is set to 'global'")
	}

	// Create the cluster
	op, err := config.clientDataproc.Projects.Regions.Clusters.Create(
		project, region, cluster).Do()
	if err != nil {
		return err
	}

	d.SetId(cluster.ClusterName)

	// Wait until it's created
	timeoutInMinutes := int(d.Timeout(schema.TimeoutCreate).Minutes())
	waitErr := dataprocClusterOperationWait(config, op, "creating Dataproc cluster", timeoutInMinutes, 3)
	if waitErr != nil {
		// The resource didn't actually create
		d.SetId("")
		return waitErr
	}

	log.Printf("[INFO] Dataproc cluster %s has been created", cluster.ClusterName)
	return resourceDataprocClusterRead(d, meta)

}

func expandClusterConfig(d *schema.ResourceData) *dataproc.ClusterConfig {
	conf := &dataproc.ClusterConfig{
		// SDK requires GceClusterConfig to be specified,
		// even if no explicit values specified
		GceClusterConfig: &dataproc.GceClusterConfig{},
	}

	if v, ok := d.GetOk("cluster_config"); ok {
		confs := v.([]interface{})
		if (len(confs)) == 0 {
			return conf
		}
	}

	if v, ok := d.GetOk("cluster_config.0.staging_bucket"); ok {
		conf.ConfigBucket = v.(string)
	}

	if cfg, ok := configOptions(d, "cluster_config.0.gce_cluster_config"); ok {
		conf.GceClusterConfig = expandGceClusterConfig(cfg)
	}

	if cfg, ok := configOptions(d, "cluster_config.0.software_config"); ok {
		conf.SoftwareConfig = expandSoftwareConfig(cfg)
	}

	if v, ok := d.GetOk("cluster_config.0.initialization_action"); ok {
		conf.InitializationActions = expandInitializationActions(v)
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
		log.Println("[INFO] got preemtible worker config")
		conf.SecondaryWorkerConfig = expandPreemptibleInstanceGroupConfig(cfg)
		if conf.SecondaryWorkerConfig.NumInstances > 0 {
			conf.SecondaryWorkerConfig.IsPreemptible = true
		}
	}
	return conf
}

func expandGceClusterConfig(cfg map[string]interface{}) *dataproc.GceClusterConfig {
	conf := &dataproc.GceClusterConfig{}

	if v, ok := cfg["zone"]; ok {
		conf.ZoneUri = v.(string)
	}
	if v, ok := cfg["network"]; ok {
		conf.NetworkUri = extractLastResourceFromUri(v.(string))
	}
	if v, ok := cfg["subnetwork"]; ok {
		conf.SubnetworkUri = extractLastResourceFromUri(v.(string))
	}
	if v, ok := cfg["tags"]; ok {
		conf.Tags = convertStringArr(v.([]interface{}))
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
	return conf
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
		icg.MachineTypeUri = extractLastResourceFromUri(v.(string))
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
		}
	}
	return icg
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
		patch := config.clientDataproc.Projects.Regions.Clusters.Patch(
			project, region, clusterName, cluster)
		op, err := patch.UpdateMask(strings.Join(updMask, ",")).Do()
		if err != nil {
			return err
		}

		// Wait until it's updated
		waitErr := dataprocClusterOperationWait(config, op, "updating Dataproc cluster ", timeoutInMinutes, 2)
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

	cluster, err := config.clientDataproc.Projects.Regions.Clusters.Get(
		project, region, clusterName).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Dataproc Cluster %q", clusterName))
	}

	d.Set("name", cluster.ClusterName)
	d.Set("region", region)
	d.Set("labels", cluster.Labels)

	cfg, err := flattenClusterConfig(d, cluster.Config)
	if err != nil {
		return err
	}

	d.Set("cluster_config", cfg)
	return nil
}

func flattenClusterConfig(d *schema.ResourceData, cfg *dataproc.ClusterConfig) ([]map[string]interface{}, error) {

	data := map[string]interface{}{
		"delete_autogen_bucket": d.Get("cluster_config.0.delete_autogen_bucket").(bool),
		"staging_bucket":        d.Get("cluster_config.0.staging_bucket").(string),

		"bucket":                    cfg.ConfigBucket,
		"gce_cluster_config":        flattenGceClusterConfig(d, cfg.GceClusterConfig),
		"software_config":           flattenSoftwareConfig(d, cfg.SoftwareConfig),
		"master_config":             flattenInstanceGroupConfig(d, cfg.MasterConfig),
		"worker_config":             flattenInstanceGroupConfig(d, cfg.WorkerConfig),
		"preemptible_worker_config": flattenPreemptibleInstanceGroupConfig(d, cfg.SecondaryWorkerConfig),
	}

	if len(cfg.InitializationActions) > 0 {
		val, err := flattenInitializationActions(cfg.InitializationActions)
		if err != nil {
			return nil, err
		}
		data["intialization_action"] = val
	}
	return []map[string]interface{}{data}, nil
}

func flattenSoftwareConfig(d *schema.ResourceData, sc *dataproc.SoftwareConfig) []map[string]interface{} {
	data := map[string]interface{}{
		"image_version":       sc.ImageVersion,
		"properties":          sc.Properties,
		"override_properties": d.Get("cluster_config.0.software_config.0.override_properties").(map[string]interface{}),
	}

	return []map[string]interface{}{data}
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
		"tags":            gcc.Tags,
		"service_account": gcc.ServiceAccount,
		"zone":            extractLastResourceFromUri(gcc.ZoneUri),
	}

	if gcc.NetworkUri != "" {
		gceConfig["network"] = extractLastResourceFromUri(gcc.NetworkUri)
	}
	if gcc.SubnetworkUri != "" {
		gceConfig["subnetwork"] = extractLastResourceFromUri(gcc.SubnetworkUri)
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
		}
	}

	data["disk_config"] = []map[string]interface{}{disk}
	return []map[string]interface{}{data}
}

func flattenInstanceGroupConfig(d *schema.ResourceData, icg *dataproc.InstanceGroupConfig) []map[string]interface{} {
	disk := map[string]interface{}{}
	data := map[string]interface{}{
	//"instance_names": []string{},
	}

	if icg != nil {
		data["num_instances"] = icg.NumInstances
		data["machine_type"] = extractLastResourceFromUri(icg.MachineTypeUri)
		data["instance_names"] = icg.InstanceNames
		if icg.DiskConfig != nil {
			disk["boot_disk_size_gb"] = icg.DiskConfig.BootDiskSizeGb
			disk["num_local_ssds"] = icg.DiskConfig.NumLocalSsds
		}
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
	deleteAutoGenBucket := d.Get("cluster_config.0.delete_autogen_bucket").(bool)
	timeoutInMinutes := int(d.Timeout(schema.TimeoutDelete).Minutes())

	if deleteAutoGenBucket {
		if err := deleteAutogenBucketIfExists(d, meta); err != nil {
			return err
		}
	}

	log.Printf("[DEBUG] Deleting Dataproc cluster %s", clusterName)
	op, err := config.clientDataproc.Projects.Regions.Clusters.Delete(
		project, region, clusterName).Do()
	if err != nil {
		return err
	}

	// Wait until it's deleted
	waitErr := dataprocClusterOperationWait(config, op, "deleting Dataproc cluster", timeoutInMinutes, 3)
	if waitErr != nil {
		return waitErr
	}

	log.Printf("[INFO] Dataproc cluster %s has been deleted", d.Id())
	d.SetId("")
	return nil
}

func deleteAutogenBucketIfExists(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	// If the user did not specify a specific override staging bucket, then GCP
	// creates one automatically. Clean this up to avoid dangling resources.
	if v, ok := d.GetOk("cluster_config.0.staging_bucket"); ok {
		log.Printf("[DEBUG] staging bucket %s (for dataproc cluster) has explicitly been set, leaving it...", v)
		return nil
	}
	bucket := d.Get("cluster_config.0.bucket").(string)

	log.Printf("[DEBUG] Attempting to delete autogenerated bucket %s (for dataproc cluster)", bucket)
	return emptyAndDeleteStorageBucket(config, bucket)
}

func emptyAndDeleteStorageBucket(config *Config, bucket string) error {
	err := deleteStorageBucketContents(config, bucket)
	if err != nil {
		return err
	}

	err = deleteEmptyBucket(config, bucket)
	if err != nil {
		return err
	}
	return nil
}

func deleteEmptyBucket(config *Config, bucket string) error {
	// remove empty bucket
	err := resource.Retry(1*time.Minute, func() *resource.RetryError {
		err := config.clientStorage.Buckets.Delete(bucket).Do()
		if err == nil {
			return nil
		}
		gerr, ok := err.(*googleapi.Error)
		if gerr.Code == http.StatusNotFound {
			// Bucket may be gone already ignore
			return nil
		}
		if ok && gerr.Code == http.StatusTooManyRequests {
			return resource.RetryableError(gerr)
		}
		return resource.NonRetryableError(err)
	})
	if err != nil {
		fmt.Printf("[ERROR] Attempting to delete autogenerated bucket (for dataproc cluster): Error deleting bucket %s: %v\n\n", bucket, err)
		return err
	}
	log.Printf("[DEBUG] Attempting to delete autogenerated bucket (for dataproc cluster): Deleted bucket %v\n\n", bucket)

	return nil

}

func deleteStorageBucketContents(config *Config, bucket string) error {

	res, err := config.clientStorage.Objects.List(bucket).Do()
	if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == http.StatusNotFound {
		// Bucket is already gone ...
		return nil
	}
	if err != nil {
		log.Fatalf("[DEBUG] Attempting to delete autogenerated bucket %s (for dataproc cluster). Error Objects.List failed: %v", bucket, err)
		return err
	}

	if len(res.Items) > 0 {
		// purge the bucket...
		log.Printf("[DEBUG] Attempting to delete autogenerated bucket (for dataproc cluster). \n\n")

		for _, object := range res.Items {
			log.Printf("[DEBUG] Attempting to delete autogenerated bucket (for dataproc cluster). Found %s", object.Name)

			err := config.clientStorage.Objects.Delete(bucket, object.Name).Do()
			if err != nil {
				if gerr, ok := err.(*googleapi.Error); ok && gerr.Code != http.StatusNotFound {
					log.Printf("[DEBUG] Attempting to delete autogenerated bucket (for dataproc cluster): Error trying to delete object: %s %s\n\n", object.Name, err)
					return err
				}
			}
			log.Printf("[DEBUG] Attempting to delete autogenerated bucket (for dataproc cluster): Object deleted: %s \n\n", object.Name)
		}
	}

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
