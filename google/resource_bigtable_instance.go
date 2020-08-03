package google

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

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
			resourceBigtableInstanceClusterReorderTypeList,
		),

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceBigtableInstanceResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceBigtableInstanceUpgradeV0,
				Version: 0,
			},
		},

		// ----------------------------------------------------------------------
		// IMPORTANT: Do not add any additional ForceNew fields to this resource.
		// Destroying/recreating instances can lead to data loss for users.
		// ----------------------------------------------------------------------
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name (also called Instance Id in the Cloud Console) of the Cloud Bigtable instance.`,
			},

			"cluster": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: `A block of cluster configuration options. This can be specified at least once, and up to 4 times.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The ID of the Cloud Bigtable cluster.`,
						},
						"zone": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The zone to create the Cloud Bigtable cluster in. Each cluster must have a different zone in the same region. Zones that support Bigtable instances are noted on the Cloud Bigtable locations page.`,
						},
						"num_nodes": {
							Type:     schema.TypeInt,
							Optional: true,
							// DEVELOPMENT instances could get returned with either zero or one node,
							// so mark as computed.
							Computed:     true,
							ValidateFunc: validation.IntAtLeast(1),
							Description:  `The number of nodes in your Cloud Bigtable cluster. Required, with a minimum of 1 for a PRODUCTION instance. Must be left unset for a DEVELOPMENT instance.`,
						},
						"storage_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "SSD",
							ValidateFunc: validation.StringInSlice([]string{"SSD", "HDD"}, false),
							Description:  `The storage type to use. One of "SSD" or "HDD". Defaults to "SSD".`,
						},
					},
				},
			},
			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The human-readable display name of the Bigtable instance. Defaults to the instance name.`,
			},

			"instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "PRODUCTION",
				ValidateFunc: validation.StringInSlice([]string{"DEVELOPMENT", "PRODUCTION"}, false),
				Description:  `The instance type to create. One of "DEVELOPMENT" or "PRODUCTION". Defaults to "PRODUCTION".`,
			},

			"deletion_protection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: `Whether or not to allow Terraform to destroy the instance. Unless this field is set to false in Terraform state, a terraform destroy or terraform apply that would delete the instance will fail.`,
			},

			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `A mapping of labels to assign to the resource.`,
			},

			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The ID of the project in which the resource belongs. If it is not provided, the provider project is used.`,
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

	if _, ok := d.GetOk("labels"); ok {
		conf.Labels = expandLabels(d)
	}

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

	id, err := replaceVars(d, config, "projects/{{project}}/instances/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

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

	instanceName := d.Get("name").(string)

	instance, err := c.InstanceInfo(ctx, instanceName)
	if err != nil {
		log.Printf("[WARN] Removing %s because it's gone", instanceName)
		d.SetId("")
		return nil
	}

	d.Set("project", project)

	clusters, err := c.Clusters(ctx, instance.Name)
	if err != nil {
		return fmt.Errorf("Error retrieving instance clusters. %s", err)
	}

	clustersNewState := []map[string]interface{}{}
	for _, cluster := range clusters {
		clustersNewState = append(clustersNewState, flattenBigtableCluster(cluster))
	}

	log.Printf("[DEBUG] Setting clusters in state: %#v", clustersNewState)
	err = d.Set("cluster", clustersNewState)
	if err != nil {
		return fmt.Errorf("Error setting clusters in state: %s", err.Error())
	}

	d.Set("name", instance.Name)
	d.Set("display_name", instance.DisplayName)
	d.Set("labels", instance.Labels)
	// Don't set instance_type: we don't want to detect drift on it because it can
	// change under-the-hood.

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

	conf := &bigtable.InstanceWithClustersConfig{
		InstanceID: d.Get("name").(string),
	}

	displayName, ok := d.GetOk("display_name")
	if !ok {
		displayName = conf.InstanceID
	}
	conf.DisplayName = displayName.(string)

	if d.HasChange("labels") {
		conf.Labels = expandLabels(d)
	}

	switch d.Get("instance_type").(string) {
	case "DEVELOPMENT":
		conf.InstanceType = bigtable.DEVELOPMENT
	case "PRODUCTION":
		conf.InstanceType = bigtable.PRODUCTION
	}

	conf.Clusters = expandBigtableClusters(d.Get("cluster").([]interface{}), conf.InstanceID)

	_, err = bigtable.UpdateInstanceAndSyncClusters(ctx, c, conf)
	if err != nil {
		return fmt.Errorf("Error updating instance. %s", err)
	}

	return resourceBigtableInstanceRead(d, meta)
}

func resourceBigtableInstanceDestroy(d *schema.ResourceData, meta interface{}) error {
	if d.Get("deletion_protection").(bool) {
		return fmt.Errorf("cannot destroy instance without setting deletion_protection=false and running `terraform apply`")
	}
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

	name := d.Get("name").(string)
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

// resourceBigtableInstanceClusterReorderTypeList causes the cluster block to
// act like a TypeSet while it's a TypeList underneath. It preserves state
// ordering on updates, and causes the resource to get recreated if it would
// attempt to perform an impossible change.
// This doesn't use the standard unordered list utility (https://github.com/GoogleCloudPlatform/magic-modules/blob/master/templates/terraform/unordered_list_customize_diff.erb)
// because some fields can't be modified using the API and we recreate the instance
// when they're changed.
func resourceBigtableInstanceClusterReorderTypeList(diff *schema.ResourceDiff, meta interface{}) error {
	oldCount, newCount := diff.GetChange("cluster.#")

	// simulate Required:true, MinItems:1, MaxItems:4 for "cluster"
	if newCount.(int) < 1 {
		return fmt.Errorf("config is invalid: Too few cluster blocks: Should have at least 1 \"cluster\" block")
	}
	if newCount.(int) > 4 {
		return fmt.Errorf("config is invalid: Too many cluster blocks: No more than 4 \"cluster\" blocks are allowed")
	}

	// exit early if we're in create (name's old value is nil)
	n, _ := diff.GetChange("name")
	if n == nil || n == "" {
		return nil
	}

	oldIds := []string{}
	clusters := make(map[string]interface{}, newCount.(int))

	for i := 0; i < oldCount.(int); i++ {
		oldId, _ := diff.GetChange(fmt.Sprintf("cluster.%d.cluster_id", i))
		if oldId != nil && oldId != "" {
			oldIds = append(oldIds, oldId.(string))
		}
	}
	log.Printf("[DEBUG] Saw old ids: %#v", oldIds)

	for i := 0; i < newCount.(int); i++ {
		_, newId := diff.GetChange(fmt.Sprintf("cluster.%d.cluster_id", i))
		_, c := diff.GetChange(fmt.Sprintf("cluster.%d", i))
		clusters[newId.(string)] = c
	}

	// create a list of clusters using the old order when possible to minimise
	// diffs
	// initially, add matching clusters to their index by id (nil otherwise)
	// then, fill in nils with new clusters.
	// [a, b, c, e] -> [c, a, d] becomes [a, nil, c] followed by [a, d, c]
	var orderedClusters []interface{}
	for i := 0; i < newCount.(int); i++ {
		// when i is out of range of old, all values are nil
		if i >= len(oldIds) {
			orderedClusters = append(orderedClusters, nil)
			continue
		}

		oldId := oldIds[i]
		if c, ok := clusters[oldId]; ok {
			log.Printf("[DEBUG] Matched: %#v", oldId)
			orderedClusters = append(orderedClusters, c)
			delete(clusters, oldId)
		} else {
			orderedClusters = append(orderedClusters, nil)
		}
	}

	log.Printf("[DEBUG] Remaining clusters: %#v", clusters)
	for _, elem := range clusters {
		for i, e := range orderedClusters {
			if e == nil {
				orderedClusters[i] = elem
			}
		}
	}

	err := diff.SetNew("cluster", orderedClusters)
	if err != nil {
		return fmt.Errorf("Error setting cluster diff: %s", err)
	}

	// Clusters can't have their zone / storage_type updated, ForceNew if it's
	// changed. This will show a diff with the old state on the left side and
	// the unmodified new state on the right and the ForceNew attributed to the
	// _old state index_ even if the diff appears to have moved.
	// This depends on the clusters having been reordered already by the prior
	// SetNew call.
	// We've implemented it here because it doesn't return an error in the
	// client and silently fails.
	for i := 0; i < newCount.(int); i++ {
		oldId, newId := diff.GetChange(fmt.Sprintf("cluster.%d.cluster_id", i))
		if oldId != newId {
			continue
		}

		oZone, nZone := diff.GetChange(fmt.Sprintf("cluster.%d.zone", i))
		if oZone != nZone {
			err := diff.ForceNew(fmt.Sprintf("cluster.%d.zone", i))
			if err != nil {
				return fmt.Errorf("Error setting cluster diff: %s", err)
			}
		}

		oST, nST := diff.GetChange(fmt.Sprintf("cluster.%d.storage_type", i))
		if oST != nST {
			err := diff.ForceNew(fmt.Sprintf("cluster.%d.storage_type", i))
			if err != nil {
				return fmt.Errorf("Error setting cluster diff: %s", err)
			}
		}
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
	id, err := replaceVars(d, config, "projects/{{project}}/instances/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
