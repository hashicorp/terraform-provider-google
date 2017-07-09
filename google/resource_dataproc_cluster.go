package google

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"

	"google.golang.org/api/dataproc/v1"

	"errors"
	"log"
	"regexp"
	"strings"
	"time"
)

func resourceDataprocCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceDataprocClusterCreate,
		Read:   resourceDataprocClusterRead,
		Update: resourceDataprocClusterUpdate,
		Delete: resourceDataprocClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of this cluster",
				Required:    true,
				ForceNew:    true,
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
			"bucket": {
				Type:        schema.TypeString,
				Description: "The Google Cloud Storage bucket to use with the Google Cloud Storage connector",
				Required:    false,
				Computed:    true,
				ForceNew:    true,
			},

			"image_version": {
				Type:        schema.TypeString,
				Description: "The image version to use for the cluster",
				Required:    false,
				Computed:    true,
				ForceNew:    true,
			},

			"no_address": {
				Type:        schema.TypeBool,
				Description: "If set to true, the instances in the cluster will not be assigned external IP addresses",
				Required:    false,
				Computed:    true,
				ForceNew:    true,
			},

			// At most one of [ network | subnetwork ] may be specified:
			"network": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"subnetwork": {
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

			"zone": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				ForceNew: true,
			},

			"labels": {
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
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"num_masters": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ForceNew: false,
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
					},
				},
			},

			"worker_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: false,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"num_workers": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ForceNew: false,
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

						"preemptible_num_workers": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},

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
					},
				},
			},

			"service_account": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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

	timeoutInMinutes := int(d.Timeout(schema.TimeoutCreate).Minutes())

	// Mandatory
	clusterName := d.Get("name").(string)
	region := d.Get("region").(string)
	zone := d.Get("zone").(string)

	gceConfig := &dataproc.GceClusterConfig{}

	if _, ok := d.GetOk("network"); ok {
		network, err := getNetworkName(d, "network")
		if err != nil {
			return err
		}
		gceConfig.NetworkUri = network
	}

	if v, ok := d.GetOk("subnetwork"); ok {
		gceConfig.SubnetworkUri = v.(string)
	}

	if v, ok := d.GetOk("tags"); ok {
		tagsList := v.([]interface{})
		tags := []string{}
		for _, v := range tagsList {
			tags = append(tags, v.(string))
		}
		gceConfig.Tags = tags
	}

	if v, ok := d.GetOk("service_account"); ok {
		gceConfig.ServiceAccount = v.(string)
	}

	if v, ok := d.GetOk("service_account_scopes"); ok {
		scopesList := v.([]interface{})
		scopes := []string{}
		for _, v := range scopesList {
			scopes = append(scopes, canonicalizeServiceScope(v.(string)))
		}

		gceConfig.ServiceAccountScopes = scopes
	}

	gceConfig.ZoneUri = zone

	clusterConfig := &dataproc.ClusterConfig{
		GceClusterConfig: gceConfig,
	}

	if v, ok := d.GetOk("master_config"); ok {
		masterConfigs := v.([]interface{})
		if len(masterConfigs) > 1 {
			return fmt.Errorf("Cannot specify more than one master_config.")
		}
		masterConfig := masterConfigs[0].(map[string]interface{})

		clusterConfig.MasterConfig = &dataproc.InstanceGroupConfig{
			DiskConfig: &dataproc.DiskConfig{},
		}

		if v, ok = masterConfig["num_masters"]; ok {
			clusterConfig.MasterConfig.NumInstances = int64(v.(int))
		}
		if v, ok = masterConfig["machine_type"]; ok {
			clusterConfig.MasterConfig.MachineTypeUri = v.(string)
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
		if len(configs) > 1 {
			return fmt.Errorf("Cannot specify more than one worker_config.")
		}
		config := configs[0].(map[string]interface{})

		clusterConfig.WorkerConfig = &dataproc.InstanceGroupConfig{
			DiskConfig: &dataproc.DiskConfig{},
		}

		if v, ok = config["num_workers"]; ok {
			clusterConfig.WorkerConfig.NumInstances = int64(v.(int))
		}
		if v, ok = config["machine_type"]; ok {
			clusterConfig.WorkerConfig.MachineTypeUri = v.(string)
		}
		if v, ok = config["boot_disk_size_gb"]; ok {
			clusterConfig.WorkerConfig.DiskConfig.BootDiskSizeGb = int64(v.(int))
		}
		if v, ok = config["num_local_ssds"]; ok {
			clusterConfig.WorkerConfig.DiskConfig.NumLocalSsds = int64(v.(int))
		}
	}

	cluster := &dataproc.Cluster{
		ClusterName: clusterName,
		ProjectId:   project,
		Config:      clusterConfig,
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

	log.Println(op)
	// Wait until it's created
	waitErr := dataprocClusterOperationWait(config, op, "creating Dataproc cluster", timeoutInMinutes, 3)
	if waitErr != nil {
		// The resource didn't actually create
		d.SetId("")
		return waitErr
	}

	log.Printf("[INFO] Dataproc cluster %s has been created", clusterName)

	d.SetId(clusterName)

	e := resourceDataprocClusterRead(d, meta)
	if e != nil {
		log.Printf("Got an error reading back after creating, \n", e)
	}
	return e
}

func resourceDataprocClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	// Only caters for certain changes to exit if not applicable
	if !(d.HasChange("labels") ||
		d.HasChange("worker_config.0.num_workers") ||
		d.HasChange("worker_config.0.preemptible_num_workers")) {
		return errors.New("update resource called, but nothing is allowed to be changed - programmer issue")
	}

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
		Config: &dataproc.ClusterConfig{
			WorkerConfig: &dataproc.InstanceGroupConfig{},
		},
	}

	d.Partial(true)

	if d.HasChange("worker_config.0.num_workers") {

		wconfigs := d.Get("worker_config").([]interface{})
		conf := wconfigs[0].(map[string]interface{})

		desiredNumWorks := conf["num_workers"].(int)
		cluster.Config.WorkerConfig.NumInstances = int64(desiredNumWorks)

		patch := config.clientDataproc.Projects.Regions.Clusters.Patch(
			project, region, clusterName, cluster)
		patch.UpdateMask("config.worker_config.num_instances")

		op, err := patch.Do()
		if err != nil {
			return err
		}

		// Wait until it's updated
		waitErr := dataprocClusterOperationWait(config, op, "updating Dataproc cluster ", timeoutInMinutes, 2)
		if waitErr != nil {
			return waitErr
		}

		log.Printf("[INFO] Dataproc cluster %s has num_workers updated to %d", d.Id(),
			desiredNumWorks)

		d.SetPartial("worker_config")
	}

	d.Partial(false)

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

	d.Set("labels", cluster.Labels)
	d.Set("name", cluster.ClusterName)
	d.Set("bucket", cluster.Config.ConfigBucket)

	if cluster.Config.GceClusterConfig != nil {
		zoneUrl := strings.Split(cluster.Config.GceClusterConfig.ZoneUri, "/")
		d.Set("zone", zoneUrl[len(zoneUrl)-1])

		d.Set("network", cluster.Config.GceClusterConfig.NetworkUri)
		d.Set("subnet", cluster.Config.GceClusterConfig.SubnetworkUri)
		d.Set("tags", cluster.Config.GceClusterConfig.Tags)
		d.Set("metadata", cluster.Config.GceClusterConfig.Metadata)
		d.Set("service_account", cluster.Config.GceClusterConfig.ServiceAccount)
		d.Set("service_account_scopes", cluster.Config.GceClusterConfig.ServiceAccountScopes)
		if len(cluster.Config.GceClusterConfig.ServiceAccountScopes) > 0 {
			d.Set("service_account_scopes", cluster.Config.GceClusterConfig.ServiceAccountScopes)
		}

	}

	if cluster.Config.SoftwareConfig != nil {
		d.Set("image_version", cluster.Config.SoftwareConfig.ImageVersion)
	}

	if cluster.Config.MasterConfig != nil {
		mType := strings.Split(cluster.Config.MasterConfig.MachineTypeUri, "/")
		masterConfig := []map[string]interface{}{
			{
				"num_masters":       cluster.Config.MasterConfig.NumInstances,
				"boot_disk_size_gb": cluster.Config.MasterConfig.DiskConfig.BootDiskSizeGb,
				"machine_type":      mType[len(mType)-1],
				"num_local_ssds":    cluster.Config.MasterConfig.DiskConfig.NumLocalSsds,
			},
		}
		d.Set("master_config", masterConfig)
	}

	if cluster.Config.WorkerConfig != nil {
		mType := strings.Split(cluster.Config.WorkerConfig.MachineTypeUri, "/")
		workerConfig := []map[string]interface{}{
			{
				"num_workers":       cluster.Config.WorkerConfig.NumInstances,
				"boot_disk_size_gb": cluster.Config.WorkerConfig.DiskConfig.BootDiskSizeGb,
				"machine_type":      mType[len(mType)-1],
				"num_local_ssds":    cluster.Config.WorkerConfig.DiskConfig.NumLocalSsds,
			},
		}
		if cluster.Config.SecondaryWorkerConfig != nil {
			workerConfig[0]["preemptible_num_workers"] = cluster.Config.SecondaryWorkerConfig.NumInstances
			workerConfig[0]["preemptible_boot_disk_size_gb"] = cluster.Config.SecondaryWorkerConfig.DiskConfig.BootDiskSizeGb
		}

		d.Set("worker_config", workerConfig)
	}

	return nil
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

	log.Printf("[DEBUG] Deleting Dataproc cluster %s", d.Get("name").(string))
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
