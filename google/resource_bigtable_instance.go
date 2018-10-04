package google

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/customdiff"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	"cloud.google.com/go/bigtable"
)

func resourceBigtableInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigtableInstanceCreate,
		Read:   resourceBigtableInstanceRead,
		// TODO: Update is only needed because we're doing forcenew in customizediff
		// when we're done with the deprecation, we can drop customizediff and make cluster forcenew
		Update: schema.Noop,
		Delete: resourceBigtableInstanceDestroy,
		CustomizeDiff: customdiff.All(
			resourceBigTableInstanceClusterCustomizeDiff,
		),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"cluster_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Deprecated:    "Use cluster instead.",
				ConflictsWith: []string{"cluster"},
			},

			"cluster": {
				Type:          schema.TypeSet,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"cluster_id", "zone", "num_nodes", "storage_type"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"zone": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"num_nodes": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"storage_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "SSD",
							ValidateFunc: validation.StringInSlice([]string{"SSD", "HDD"}, false),
						},
					},
				},
			},

			"zone": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				Deprecated:    "Use cluster instead.",
				ConflictsWith: []string{"cluster"},
			},

			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"num_nodes": {
				Type:          schema.TypeInt,
				Optional:      true,
				Deprecated:    "Use cluster instead.",
				ConflictsWith: []string{"cluster"},
			},

			"instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "PRODUCTION",
				ValidateFunc: validation.StringInSlice([]string{"DEVELOPMENT", "PRODUCTION"}, false),
			},

			"storage_type": {
				Type:          schema.TypeString,
				Optional:      true,
				Default:       "SSD",
				ValidateFunc:  validation.StringInSlice([]string{"SSD", "HDD"}, false),
				Deprecated:    "Use cluster instead.",
				ConflictsWith: []string{"cluster"},
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceBigTableInstanceClusterCustomizeDiff(d *schema.ResourceDiff, meta interface{}) error {
	if d.Get("cluster_id").(string) == "" && d.Get("cluster.#").(int) == 0 {
		return fmt.Errorf("At least one cluster must be set.")
	}
	if !d.HasChange("cluster_id") && !d.HasChange("zone") && !d.HasChange("num_nodes") &&
		!d.HasChange("storage_type") && !d.HasChange("cluster") {
		return nil
	}
	if d.Get("cluster.#").(int) == 1 {
		// if we have exactly one cluster, and it has the same values as the old top-level
		// values, we can assume the user is trying to go from the deprecated values to the
		// new values, and we shouldn't ForceNew. We know that the top-level values aren't
		// set, because they ConflictWith cluster.
		oldID, _ := d.GetChange("cluster_id")
		oldNodes, _ := d.GetChange("num_nodes")
		oldZone, _ := d.GetChange("zone")
		oldStorageType, _ := d.GetChange("storage_type")
		new := d.Get("cluster").(*schema.Set).List()[0].(map[string]interface{})

		if oldID.(string) == new["cluster_id"].(string) &&
			oldNodes.(int) == new["num_nodes"].(int) &&
			oldZone.(string) == new["zone"].(string) &&
			oldStorageType.(string) == new["storage_type"].(string) {
			return nil
		}
	}
	if d.HasChange("cluster_id") {
		d.ForceNew("cluster_id")
	}
	if d.HasChange("cluster") {
		d.ForceNew("cluster")
	}
	if d.HasChange("zone") {
		d.ForceNew("zone")
	}
	if d.HasChange("num_nodes") {
		d.ForceNew("num_nodes")
	}
	if d.HasChange("storage_type") {
		d.ForceNew("storage_type")
	}
	return nil
}

func resourceBigtableInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ctx := context.Background()

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	conf := &bigtable.InstanceWithClustersConfig{
		InstanceID: d.Get("name").(string),
	}

	displayName, ok := d.GetOk("display_name")
	if !ok {
		displayName = conf.InstanceID
	}
	conf.DisplayName = displayName.(string)

	switch d.Get("instance_type").(string) {
	case "DEVELOPMENT":
		conf.InstanceType = bigtable.DEVELOPMENT
	case "PRODUCTION":
		conf.InstanceType = bigtable.PRODUCTION
	}

	if d.Get("cluster.#").(int) > 0 {
		// expand cluster
		conf.Clusters = expandBigtableClusters(d.Get("cluster").(*schema.Set).List(), conf.InstanceID, config.Zone)
		if err != nil {
			return fmt.Errorf("error expanding clusters: %s", err.Error())
		}
	} else {
		// TODO: remove this when we're done with the deprecation period
		zone, err := getZone(d, config)
		if err != nil {
			return err
		}
		cluster := bigtable.ClusterConfig{
			InstanceID: conf.InstanceID,
			NumNodes:   int32(d.Get("num_nodes").(int)),
			Zone:       zone,
			ClusterID:  d.Get("cluster_id").(string),
		}
		switch d.Get("storage_type").(string) {
		case "HDD":
			cluster.StorageType = bigtable.HDD
		case "SSD":
			cluster.StorageType = bigtable.SSD
		}
		conf.Clusters = append(conf.Clusters, cluster)
	}

	c, err := config.bigtableClientFactory.NewInstanceAdminClient(project)
	if err != nil {
		return fmt.Errorf("Error starting instance admin client. %s", err)
	}

	defer c.Close()

	err = c.CreateInstanceWithClusters(ctx, conf)
	if err != nil {
		return fmt.Errorf("Error creating instance. %s", err)
	}

	d.SetId(conf.InstanceID)

	return resourceBigtableInstanceRead(d, meta)
}

func resourceBigtableInstanceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ctx := context.Background()

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	c, err := config.bigtableClientFactory.NewInstanceAdminClient(project)
	if err != nil {
		return fmt.Errorf("Error starting instance admin client. %s", err)
	}

	defer c.Close()

	instance, err := c.InstanceInfo(ctx, d.Id())
	if err != nil {
		log.Printf("[WARN] Removing %s because it's gone", d.Id())
		d.SetId("")
		return fmt.Errorf("Error retrieving instance. Could not find %s. %s", d.Id(), err)
	}

	d.Set("project", project)
	if d.Get("cluster.#").(int) > 0 {
		clusters := d.Get("cluster").(*schema.Set).List()
		clusterState := []map[string]interface{}{}
		for _, cl := range clusters {
			cluster := cl.(map[string]interface{})
			clus, err := c.GetCluster(ctx, instance.Name, cluster["cluster_id"].(string))
			if err != nil {
				if isGoogleApiErrorWithCode(err, 404) {
					log.Printf("[WARN] Cluster %q not found, not setting it in state", cluster["cluster_id"].(string))
					continue
				}
				return fmt.Errorf("Error retrieving cluster %q: %s", cluster["cluster_id"].(string), err.Error())
			}
			clusterState = append(clusterState, flattenBigtableCluster(clus, cluster["storage_type"].(string)))
		}
		err = d.Set("cluster", clusterState)
		if err != nil {
			return fmt.Errorf("Error setting clusters in state: %s", err.Error())
		}
		d.Set("cluster_id", "")
		d.Set("zone", "")
		d.Set("num_nodes", 0)
		d.Set("storage_type", "SSD")
	} else {
		// TODO remove this when we're done with our deprecation period
		zone, err := getZone(d, config)
		if err != nil {
			return err
		}
		d.Set("zone", zone)
	}
	d.Set("name", instance.Name)
	d.Set("display_name", instance.DisplayName)

	return nil
}

func resourceBigtableInstanceDestroy(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ctx := context.Background()

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	c, err := config.bigtableClientFactory.NewInstanceAdminClient(project)
	if err != nil {
		return fmt.Errorf("Error starting instance admin client. %s", err)
	}

	defer c.Close()

	name := d.Id()
	err = c.DeleteInstance(ctx, name)
	if err != nil {
		return fmt.Errorf("Error deleting instance. %s", err)
	}

	d.SetId("")

	return nil
}

func flattenBigtableCluster(c *bigtable.ClusterInfo, storageType string) map[string]interface{} {
	return map[string]interface{}{
		"zone":         c.Zone,
		"num_nodes":    c.ServeNodes,
		"cluster_id":   c.Name,
		"storage_type": storageType,
	}
}

func expandBigtableClusters(clusters []interface{}, instanceID string, defaultZone string) []bigtable.ClusterConfig {
	results := make([]bigtable.ClusterConfig, 0, len(clusters))
	for _, c := range clusters {
		cluster := c.(map[string]interface{})
		zone := defaultZone
		if confZone, ok := cluster["zone"]; ok {
			zone = confZone.(string)
		}
		var storageType bigtable.StorageType
		switch cluster["storage_type"].(string) {
		case "SSD":
			storageType = bigtable.SSD
		case "HDD":
			storageType = bigtable.HDD
		}
		results = append(results, bigtable.ClusterConfig{
			InstanceID:  instanceID,
			Zone:        zone,
			ClusterID:   cluster["cluster_id"].(string),
			NumNodes:    int32(cluster["num_nodes"].(int)),
			StorageType: storageType,
		})
	}
	return results
}
