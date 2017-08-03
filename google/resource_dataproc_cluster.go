package google

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

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

			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "global",
				ForceNew: true,
			},

			"zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

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

			"image_version": {
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
				ConflictsWith: []string{"subnetwork"},
				StateFunc: func(s interface{}) string {
					return extractLastResourceFromUri(s.(string))
				},
			},

			"subnetwork": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"network"},
				StateFunc: func(s interface{}) string {
					return extractLastResourceFromUri(s.(string))
				},
			},

			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     schema.TypeString,
			},

			"properties": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     schema.TypeString,
			},

			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"master_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"num_masters": {
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

						"num_local_ssds": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},

						"boot_disk_size_gb": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ForceNew: true,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := v.(int)

								if value < 10 {
									errors = append(errors, fmt.Errorf(
										"%q cannot be less than 10", k))
								}
								return
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

			"worker_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"num_workers": {
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

						"num_local_ssds": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},

						"boot_disk_size_gb": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ForceNew: true,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := v.(int)

								if value < 10 {
									errors = append(errors, fmt.Errorf(
										"%q cannot be less than 10", k))
								}
								return
							},
						},

						"preemptible_num_workers": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},

						// "preemptible_machine_type" cannot be specified directly, it takes its
						// value from the standard worker "machine_type" field

						"preemptible_boot_disk_size_gb": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ForceNew: true,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := v.(int)

								if value < 10 {
									errors = append(errors, fmt.Errorf(
										"%q cannot be less than 10", k))
								}
								return
							},
						},

						"instance_names": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},

						"preemptible_instance_names": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},

			"service_account": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"service_account_scopes": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					StateFunc: func(v interface{}) string {
						return canonicalizeServiceScope(v.(string))
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
						},
					},
				},
			},

			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     schema.TypeString,
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

	// Mandatory
	clusterName := d.Get("name").(string)
	region := d.Get("region").(string)
	zone, zok := d.GetOk("zone")

	if region == "global" && !zok {
		return errors.New("zone is mandatory when region is set to 'global'")
	}

	gceConfig := &dataproc.GceClusterConfig{
		ZoneUri: zone.(string),
	}

	if v, ok := d.GetOk("network"); ok {
		gceConfig.NetworkUri = extractLastResourceFromUri(v.(string))
	}

	if v, ok := d.GetOk("subnetwork"); ok {
		gceConfig.SubnetworkUri = extractLastResourceFromUri(v.(string))
	}

	if v, ok := d.GetOk("tags"); ok {
		gceConfig.Tags = convertStringArr(v.([]interface{}))
	}

	if v, ok := d.GetOk("service_account"); ok {
		gceConfig.ServiceAccount = v.(string)
	}

	if v, ok := d.GetOk("service_account_scopes"); ok {
		gceConfig.ServiceAccountScopes = convertAndMapStringArr(v.([]interface{}), canonicalizeServiceScope)
		sort.Strings(gceConfig.ServiceAccountScopes)
	}

	clusterConfig := &dataproc.ClusterConfig{
		GceClusterConfig: gceConfig,
		SoftwareConfig:   &dataproc.SoftwareConfig{},
	}

	if v, ok := d.GetOk("initialization_action"); ok {
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
		clusterConfig.InitializationActions = actions
	}

	if v, ok := d.GetOk("staging_bucket"); ok {
		clusterConfig.ConfigBucket = v.(string)
	}

	if v, ok := d.GetOk("master_config"); ok {
		masterConfigs := v.([]interface{})
		masterConfig := masterConfigs[0].(map[string]interface{})

		clusterConfig.MasterConfig = &dataproc.InstanceGroupConfig{
			DiskConfig: &dataproc.DiskConfig{},
		}

		if v, ok = masterConfig["num_masters"]; ok {
			clusterConfig.MasterConfig.NumInstances = int64(v.(int))
		}
		if v, ok = masterConfig["machine_type"]; ok {
			clusterConfig.MasterConfig.MachineTypeUri = extractLastResourceFromUri(v.(string))
		}
		if v, ok = masterConfig["boot_disk_size_gb"]; ok {
			clusterConfig.MasterConfig.DiskConfig.BootDiskSizeGb = int64(v.(int))
		}
		if v, ok = masterConfig["num_local_ssds"]; ok {
			clusterConfig.MasterConfig.DiskConfig.NumLocalSsds = int64(v.(int))
		}
	}

	if v, ok := d.GetOk("worker_config"); ok {
		configs := v.([]interface{})
		config := configs[0].(map[string]interface{})

		clusterConfig.WorkerConfig = &dataproc.InstanceGroupConfig{
			DiskConfig: &dataproc.DiskConfig{},
		}

		if v, ok = config["num_workers"]; ok {
			clusterConfig.WorkerConfig.NumInstances = int64(v.(int))
		}
		if v, ok = config["machine_type"]; ok {
			clusterConfig.WorkerConfig.MachineTypeUri = extractLastResourceFromUri(v.(string))
		}
		if v, ok = config["boot_disk_size_gb"]; ok {
			clusterConfig.WorkerConfig.DiskConfig.BootDiskSizeGb = int64(v.(int))
		}
		if v, ok = config["num_local_ssds"]; ok {
			clusterConfig.WorkerConfig.DiskConfig.NumLocalSsds = int64(v.(int))
		}

		clusterConfig.SecondaryWorkerConfig = &dataproc.InstanceGroupConfig{
			DiskConfig: &dataproc.DiskConfig{},
		}

		if v, ok = config["preemptible_num_workers"]; ok {
			clusterConfig.SecondaryWorkerConfig.NumInstances = int64(v.(int))
			if clusterConfig.SecondaryWorkerConfig.NumInstances > 0 {
				clusterConfig.SecondaryWorkerConfig.IsPreemptible = true
			}
		}
		if v, ok = config["preemptible_boot_disk_size_gb"]; ok {
			clusterConfig.SecondaryWorkerConfig.DiskConfig.BootDiskSizeGb = int64(v.(int))
		}
	}

	cluster := &dataproc.Cluster{
		ClusterName: clusterName,
		ProjectId:   project,
		Config:      clusterConfig,
	}

	if v, ok := d.GetOk("properties"); ok {
		m := make(map[string]string)
		for k, val := range v.(map[string]interface{}) {
			m[k] = val.(string)
		}
		cluster.Config.SoftwareConfig.Properties = m
	}

	if v, ok := d.GetOk("image_version"); ok {
		cluster.Config.SoftwareConfig.ImageVersion = v.(string)
	}

	if v, ok := d.GetOk("labels"); ok {
		m := make(map[string]string)
		for k, val := range v.(map[string]interface{}) {
			m[k] = val.(string)
		}
		cluster.Labels = m
	}
	if v, ok := d.GetOk("metadata"); ok {
		m := make(map[string]string)
		for k, val := range v.(map[string]interface{}) {
			m[k] = val.(string)
		}
		cluster.Config.GceClusterConfig.Metadata = m
	}

	// Create the cluster
	op, err := config.clientDataproc.Projects.Regions.Clusters.Create(
		project, region, cluster).Do()
	if err != nil {
		return err
	}

	d.SetId(clusterName)

	// Wait until it's created
	timeoutInMinutes := int(d.Timeout(schema.TimeoutCreate).Minutes())
	waitErr := dataprocClusterOperationWait(config, op, "creating Dataproc cluster", timeoutInMinutes, 3)
	if waitErr != nil {
		// The resource didn't actually create
		d.SetId("")
		return waitErr
	}

	log.Printf("[INFO] Dataproc cluster %s has been created", clusterName)
	return resourceDataprocClusterRead(d, meta)

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

	if d.HasChange("worker_config.0.num_workers") {
		wconfigs := d.Get("worker_config").([]interface{})
		conf := wconfigs[0].(map[string]interface{})

		desiredNumWorks := conf["num_workers"].(int)
		cluster.Config.WorkerConfig = &dataproc.InstanceGroupConfig{
			NumInstances: int64(desiredNumWorks),
		}

		updMask = append(updMask, "config.worker_config.num_instances")
	}

	if d.HasChange("worker_config.0.preemptible_num_workers") {
		wconfigs := d.Get("worker_config").([]interface{})
		conf := wconfigs[0].(map[string]interface{})

		desiredNumWorks := conf["preemptible_num_workers"].(int)
		cluster.Config.SecondaryWorkerConfig = &dataproc.InstanceGroupConfig{
			NumInstances: int64(desiredNumWorks),
		}

		updMask = append(updMask, "config.secondary_worker_config.num_instances")
	}

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
	return resourceDataprocClusterRead(d, meta)
}

func extractLastResourceFromUri(rUri string) string {
	rUris := strings.Split(rUri, "/")
	return rUris[len(rUris)-1]
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

	d.Set("labels", cluster.Labels)
	d.Set("name", cluster.ClusterName)
	d.Set("bucket", cluster.Config.ConfigBucket)

	if len(cluster.Config.InitializationActions) > 0 {
		actions := []map[string]interface{}{}
		for _, v := range cluster.Config.InitializationActions {

			action := []map[string]interface{}{
				{
					"script": v.ExecutableFile,
				},
			}
			if len(v.ExecutionTimeout) > 0 {
				tsec, err := extractInitTimeout(v.ExecutionTimeout)
				if err != nil {
					return err
				}
				action[0]["timeout_sec"] = tsec
			}

			actions = append(actions, action[0])
		}
		d.Set("initialization_action", actions)
	}

	if cluster.Config.GceClusterConfig != nil {
		d.Set("zone", extractLastResourceFromUri(cluster.Config.GceClusterConfig.ZoneUri))
		d.Set("network", extractLastResourceFromUri(cluster.Config.GceClusterConfig.NetworkUri))
		d.Set("subnet", extractLastResourceFromUri(cluster.Config.GceClusterConfig.SubnetworkUri))
		d.Set("tags", cluster.Config.GceClusterConfig.Tags)
		d.Set("metadata", cluster.Config.GceClusterConfig.Metadata)
		d.Set("service_account", cluster.Config.GceClusterConfig.ServiceAccount)
		if len(cluster.Config.GceClusterConfig.ServiceAccountScopes) > 0 {
			sort.Strings(cluster.Config.GceClusterConfig.ServiceAccountScopes)
			d.Set("service_account_scopes", cluster.Config.GceClusterConfig.ServiceAccountScopes)
		}

	}

	if cluster.Config.SoftwareConfig != nil {
		d.Set("image_version", cluster.Config.SoftwareConfig.ImageVersion)
		//We only want our overriden values here for now
		//d.Set("properties", cluster.Config.SoftwareConfig.Properties)
	}

	if cluster.Config.MasterConfig != nil {
		masterConfig := []map[string]interface{}{
			{
				"num_masters":       cluster.Config.MasterConfig.NumInstances,
				"boot_disk_size_gb": cluster.Config.MasterConfig.DiskConfig.BootDiskSizeGb,
				"machine_type":      extractLastResourceFromUri(cluster.Config.MasterConfig.MachineTypeUri),
				"num_local_ssds":    cluster.Config.MasterConfig.DiskConfig.NumLocalSsds,
				"instance_names":    cluster.Config.MasterConfig.InstanceNames,
			},
		}
		d.Set("master_config", masterConfig)
	}

	if cluster.Config.WorkerConfig != nil {
		workerConfig := []map[string]interface{}{
			{
				"num_workers":       cluster.Config.WorkerConfig.NumInstances,
				"boot_disk_size_gb": cluster.Config.WorkerConfig.DiskConfig.BootDiskSizeGb,
				"machine_type":      extractLastResourceFromUri(cluster.Config.WorkerConfig.MachineTypeUri),
				"num_local_ssds":    cluster.Config.WorkerConfig.DiskConfig.NumLocalSsds,
				"instance_names":    cluster.Config.WorkerConfig.InstanceNames,
			},
		}

		if cluster.Config.SecondaryWorkerConfig != nil {
			workerConfig[0]["preemptible_num_workers"] = cluster.Config.SecondaryWorkerConfig.NumInstances
			workerConfig[0]["preemptible_boot_disk_size_gb"] = cluster.Config.SecondaryWorkerConfig.DiskConfig.BootDiskSizeGb
			workerConfig[0]["preemptible_instance_names"] = cluster.Config.SecondaryWorkerConfig.InstanceNames
		}

		d.Set("worker_config", workerConfig)
	}

	return nil
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
	deleteAutoGenBucket := d.Get("delete_autogen_bucket").(bool)
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
	if v, ok := d.GetOk("staging_bucket"); ok {
		log.Printf("[DEBUG] staging bucket %s (for dataproc cluster) has explicitly been set, leaving it...", v)
		return nil
	}
	bucket := d.Get("bucket").(string)

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
