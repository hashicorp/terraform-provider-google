package google

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"cloud.google.com/go/bigtable"
)

func resourceBigtableInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigtableInstanceCreate,
		Read:   resourceBigtableInstanceRead,
		Update: resourceBigtableInstanceUpdate,
		Delete: resourceBigtableInstanceDestroy,

		Importer: &schema.ResourceImporter{
			State: resourceBigtableInstanceImport,
		},

		CustomizeDiff: customdiff.All(
			resourceBigtableInstanceValidateDevelopment,
			resourceBigtableInstanceClusterReorderTypeList,
		),

		MigrateState: resourceBigtableInstanceMigrateState,

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"cluster": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"zone": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"num_nodes": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(3),
						},
						"storage_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "SSD",
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"SSD", "HDD"}, false),
						},
					},
				},
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "PRODUCTION",
				ValidateFunc: validation.StringInSlice([]string{"DEVELOPMENT", "PRODUCTION"}, false),
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Removed:  "Use cluster instead.",
			},

			"zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Removed:  "Use cluster instead.",
			},

			"num_nodes": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Removed:  "Use cluster instead.",
			},

			"storage_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Removed:  "Use cluster instead.",
			},
		},
	}
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

	conf.Clusters = expandBigtableClusters(d.Get("cluster").([]interface{}), conf.InstanceID)

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

	clusters, err := c.Clusters(ctx, instance.Name)
	if err != nil {
		return fmt.Errorf("Error retrieving instance clusters. %s", err)
	}

	clustersNewState := []map[string]interface{}{}
	for i, cluster := range clusters {
		// DEVELOPMENT clusters have num_nodes = 0 on their first (and only)
		// cluster while PRODUCTION clusters will have at least 3.
		if i == 0 {
			var instanceType string
			if cluster.ServeNodes == 0 {
				instanceType = "DEVELOPMENT"
			} else {
				instanceType = "PRODUCTION"
			}

			d.Set("instance_type", instanceType)
		}

		clustersNewState = append(clustersNewState, flattenBigtableCluster(cluster))
	}

	log.Printf("[DEBUG] Setting clusters in state: %#v", clustersNewState)
	err = d.Set("cluster", clustersNewState)
	if err != nil {
		return fmt.Errorf("Error setting clusters in state: %s", err.Error())
	}

	d.Set("name", instance.Name)
	d.Set("display_name", instance.DisplayName)

	return nil
}

func resourceBigtableInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
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

	clusters, err := c.Clusters(ctx, d.Get("name").(string))
	if err != nil {
		return fmt.Errorf("Error retrieving clusters for instance %s", err.Error())
	}

	clusterMap := make(map[string]*bigtable.ClusterInfo, len(clusters))
	for _, cluster := range clusters {
		clusterMap[cluster.Name] = cluster
	}

	for _, cluster := range d.Get("cluster").([]interface{}) {
		config := cluster.(map[string]interface{})
		cluster_id := config["cluster_id"].(string)
		if cluster, ok := clusterMap[cluster_id]; ok {
			if cluster.ServeNodes != config["num_nodes"].(int) {
				err = c.UpdateCluster(ctx, d.Get("name").(string), cluster.Name, int32(config["num_nodes"].(int)))
				if err != nil {
					return fmt.Errorf("Error updating cluster %s for instance %s", cluster.Name, d.Get("name").(string))
				}
			}
		}
	}

	return resourceBigtableInstanceRead(d, meta)
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

func flattenBigtableCluster(c *bigtable.ClusterInfo) map[string]interface{} {
	var storageType string
	switch c.StorageType {
	case bigtable.SSD:
		storageType = "SSD"
	case bigtable.HDD:
		storageType = "HDD"
	}

	return map[string]interface{}{
		"zone":         c.Zone,
		"num_nodes":    c.ServeNodes,
		"cluster_id":   c.Name,
		"storage_type": storageType,
	}
}

func expandBigtableClusters(clusters []interface{}, instanceID string) []bigtable.ClusterConfig {
	results := make([]bigtable.ClusterConfig, 0, len(clusters))
	for _, c := range clusters {
		cluster := c.(map[string]interface{})
		zone := cluster["zone"].(string)
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

func resourceBigtableInstanceValidateDevelopment(diff *schema.ResourceDiff, meta interface{}) error {
	if diff.Get("instance_type").(string) != "DEVELOPMENT" {
		return nil
	}
	if diff.Get("cluster.#").(int) != 1 {
		return fmt.Errorf("config is invalid: instance with instance_type=\"DEVELOPMENT\" should have exactly one \"cluster\" block")
	}
	if diff.Get("cluster.0.num_nodes").(int) != 0 {
		return fmt.Errorf("config is invalid: num_nodes cannot be set for instance_type=\"DEVELOPMENT\"")
	}
	return nil
}

func resourceBigtableInstanceClusterReorderTypeList(diff *schema.ResourceDiff, meta interface{}) error {
	old_count, new_count := diff.GetChange("cluster.#")

	// simulate Required:true, MinItems:1, MaxItems:4 for "cluster"
	if new_count.(int) < 1 {
		return fmt.Errorf("config is invalid: Too few cluster blocks: Should have at least 1 \"cluster\" block")
	}
	if new_count.(int) > 4 {
		return fmt.Errorf("config is invalid: Too many cluster blocks: No more than 4 \"cluster\" blocks are allowed")
	}

	if old_count.(int) != new_count.(int) {
		return nil
	}

	var old_ids []string
	clusters := make(map[string]interface{}, new_count.(int))

	for i := 0; i < new_count.(int); i++ {
		old_id, new_id := diff.GetChange(fmt.Sprintf("cluster.%d.cluster_id", i))
		if old_id != nil && old_id != "" {
			old_ids = append(old_ids, old_id.(string))
		}
		_, c := diff.GetChange(fmt.Sprintf("cluster.%d", i))
		clusters[new_id.(string)] = c
	}

	// reorder clusters according to the old cluster order
	var old_cluster_order []interface{}
	for _, id := range old_ids {
		if c, ok := clusters[id]; ok {
			old_cluster_order = append(old_cluster_order, c)
		}
	}

	err := diff.SetNew("cluster", old_cluster_order)
	if err != nil {
		return fmt.Errorf("Error setting cluster diff: %s", err)
	}

	return nil
}

func resourceBigtableInstanceImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{
		"projects/(?P<project>[^/]+)/instances/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<name>[^/]+)",
		"(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func resourceBigtableInstanceMigrateState(
	v int,
	is *terraform.InstanceState,
	meta interface{},
) (*terraform.InstanceState, error) {
	switch v {
	case 0:
		// Pre-2.14.0 version; assume "cluster" is a set
		// Pull out clusters and reindex to store as a list
		// TODO: Anything we can do about ordering?
		//log.Printf("hashes (len: %d): %v", len(hashes), hashes)
		//log.Printf("hashes[0]: %s; hashes[1]: %s", hashes[0], hashes[1])
		hashes := make(map[string]bool)
		for k := range is.Attributes {
			if strings.HasPrefix(k, "cluster.") && !strings.Contains(k, "#") {
				//log.Printf("key: %s", k)
				parts := strings.Split(k, ".")
				hash := parts[1]
				hashes[hash] = true
				//log.Printf("hashes: %v", hashes)
			}
		}
		newAttributes := make(map[string]string)
		idx := 0
		fields := []string{"cluster_id", "num_nodes", "storage_type", "zone"}
		for hash, _ := range hashes {
			//log.Printf("hash: %s", hash)
			for _, field := range fields {
				//log.Printf("field: %s", field)
				oldAttrKey := fmt.Sprintf("cluster.%s.%s", hash, field)
				newAttrKey := fmt.Sprintf("cluster.%d.%s", idx, field)
				if _, exists := is.Attributes[oldAttrKey]; exists {
					newAttributes[newAttrKey] = is.Attributes[oldAttrKey]
					delete(is.Attributes, oldAttrKey)
				}
			}
			idx++
		}
		for k, v := range newAttributes {
			is.Attributes[k] = v
		}
		// Also remove legacy cluster_id, zone, num_nodes, and storage_type attributes
		// in favor of nested cluster object
		for _, field := range fields {
			if _, exists := is.Attributes[field]; exists {
				delete(is.Attributes, field)
			}
		}
	default:
		return nil, fmt.Errorf("invalid schema version %d", v)
	}

	return is, nil
}
