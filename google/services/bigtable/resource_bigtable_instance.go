// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package bigtable

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"cloud.google.com/go/bigtable"
)

func ResourceBigtableInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigtableInstanceCreate,
		Read:   resourceBigtableInstanceRead,
		Update: resourceBigtableInstanceUpdate,
		Delete: resourceBigtableInstanceDestroy,

		Importer: &schema.ResourceImporter{
			State: resourceBigtableInstanceImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			resourceBigtableInstanceClusterReorderTypeList,
			resourceBigtableInstanceUniqueClusterID,
		),

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceBigtableInstanceResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: ResourceBigtableInstanceUpgradeV0,
				Version: 0,
			},
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name (also called Instance Id in the Cloud Console) of the Cloud Bigtable instance. Must be 6-33 characters and must only contain hyphens, lowercase letters and numbers.`,
			},

			"cluster": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: `A block of cluster configuration options. This can be specified at least once.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The ID of the Cloud Bigtable cluster. Must be 6-30 characters and must only contain hyphens, lowercase letters and numbers.`,
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: `The zone to create the Cloud Bigtable cluster in. Each cluster must have a different zone in the same region. Zones that support Bigtable instances are noted on the Cloud Bigtable locations page.`,
						},
						"num_nodes": {
							Type:     schema.TypeInt,
							Optional: true,
							// DEVELOPMENT instances could get returned with either zero or one node,
							// so mark as computed.
							Computed:     true,
							ValidateFunc: validation.IntAtLeast(1),
							Description:  `The number of nodes in the cluster. If no value is set, Cloud Bigtable automatically allocates nodes based on your data footprint and optimized for 50% storage utilization.`,
						},
						"storage_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "SSD",
							ValidateFunc: validation.StringInSlice([]string{"SSD", "HDD"}, false),
							Description:  `The storage type to use. One of "SSD" or "HDD". Defaults to "SSD".`,
						},
						"kms_key_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `Describes the Cloud KMS encryption key that will be used to protect the destination Bigtable cluster. The requirements for this key are: 1) The Cloud Bigtable service account associated with the project that contains this cluster must be granted the cloudkms.cryptoKeyEncrypterDecrypter role on the CMEK key. 2) Only regional keys can be used and the region of the CMEK key must match the region of the cluster. 3) All clusters within an instance must use the same CMEK key. Values are of the form projects/{project}/locations/{location}/keyRings/{keyring}/cryptoKeys/{key}`,
						},
						"autoscaling_config": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "A list of Autoscaling configurations. Only one element is used and allowed.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"min_nodes": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: `The minimum number of nodes for autoscaling.`,
									},
									"max_nodes": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: `The maximum number of nodes for autoscaling.`,
									},
									"cpu_target": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: `The target CPU utilization for autoscaling. Value must be between 10 and 80.`,
									},
									"storage_target": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: `The target storage utilization for autoscaling, in GB, for each node in a cluster. This number is limited between 2560 (2.5TiB) and 5120 (5TiB) for a SSD cluster and between 8192 (8TiB) and 16384 (16 TiB) for an HDD cluster. If not set, whatever is already set for the cluster will not change, or if the cluster is just being created, it will use the default value of 2560 for SSD clusters and 8192 for HDD clusters.`,
									},
								},
							},
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
				Deprecated:   `It is recommended to leave this field unspecified since the distinction between "DEVELOPMENT" and "PRODUCTION" instances is going away, and all instances will become "PRODUCTION" instances. This means that new and existing "DEVELOPMENT" instances will be converted to "PRODUCTION" instances. It is recommended for users to use "PRODUCTION" instances in any case, since a 1-node "PRODUCTION" instance is functionally identical to a "DEVELOPMENT" instance, but without the accompanying restrictions.`,
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
		UseJSONNumber: true,
	}
}

func resourceBigtableInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	ctx := context.Background()

	project, err := tpgresource.GetProject(d, config)
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
		conf.Labels = tpgresource.ExpandLabels(d)
	}

	switch d.Get("instance_type").(string) {
	case "DEVELOPMENT":
		conf.InstanceType = bigtable.DEVELOPMENT
	case "PRODUCTION":
		conf.InstanceType = bigtable.PRODUCTION
	}

	conf.Clusters, err = expandBigtableClusters(d.Get("cluster").([]interface{}), conf.InstanceID, config)
	if err != nil {
		return err
	}

	c, err := config.BigTableClientFactory(userAgent).NewInstanceAdminClient(project)
	if err != nil {
		return fmt.Errorf("Error starting instance admin client. %s", err)
	}

	defer c.Close()

	ctxWithTimeout, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutCreate))
	defer cancel()
	if err := c.CreateInstanceWithClusters(ctxWithTimeout, conf); err != nil {
		return fmt.Errorf("Error creating instance. %s", err)
	}

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/instances/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return resourceBigtableInstanceRead(d, meta)
}

func resourceBigtableInstanceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	ctx := context.Background()

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	c, err := config.BigTableClientFactory(userAgent).NewInstanceAdminClient(project)
	if err != nil {
		return fmt.Errorf("Error starting instance admin client. %s", err)
	}

	defer c.Close()

	instanceName := d.Get("name").(string)

	instance, err := c.InstanceInfo(ctx, instanceName)
	if err != nil {
		if tpgresource.IsNotFoundGrpcError(err) {
			log.Printf("[WARN] Removing %s because it's gone", instanceName)
			d.SetId("")
			return nil
		}
		return err
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}

	clusters, err := c.Clusters(ctx, instance.Name)
	if err != nil {
		partiallyUnavailableErr, ok := err.(bigtable.ErrPartiallyUnavailable)

		if !ok {
			return fmt.Errorf("Error retrieving instance clusters. %s", err)
		}

		unavailableClusterZones := getUnavailableClusterZones(d.Get("cluster").([]interface{}), partiallyUnavailableErr.Locations)

		if len(unavailableClusterZones) > 0 {
			return fmt.Errorf("Error retrieving instance clusters. The following zones are unavailable: %s", strings.Join(unavailableClusterZones, ", "))
		}
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

	if err := d.Set("name", instance.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err := d.Set("display_name", instance.DisplayName); err != nil {
		return fmt.Errorf("Error setting display_name: %s", err)
	}
	if err := d.Set("labels", instance.Labels); err != nil {
		return fmt.Errorf("Error setting labels: %s", err)
	}
	// Don't set instance_type: we don't want to detect drift on it because it can
	// change under-the-hood.

	return nil
}

func resourceBigtableInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	ctx := context.Background()

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	c, err := config.BigTableClientFactory(userAgent).NewInstanceAdminClient(project)
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
		conf.Labels = tpgresource.ExpandLabels(d)
	}

	switch d.Get("instance_type").(string) {
	case "DEVELOPMENT":
		conf.InstanceType = bigtable.DEVELOPMENT
	case "PRODUCTION":
		conf.InstanceType = bigtable.PRODUCTION
	}

	conf.Clusters, err = expandBigtableClusters(d.Get("cluster").([]interface{}), conf.InstanceID, config)
	if err != nil {
		return err
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutUpdate))
	defer cancel()
	if _, err := bigtable.UpdateInstanceAndSyncClusters(ctxWithTimeout, c, conf); err != nil {
		return fmt.Errorf("Error updating instance. %s", err)
	}

	return resourceBigtableInstanceRead(d, meta)
}

func resourceBigtableInstanceDestroy(d *schema.ResourceData, meta interface{}) error {
	if d.Get("deletion_protection").(bool) {
		return fmt.Errorf("cannot destroy instance without setting deletion_protection=false and running `terraform apply`")
	}
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	ctx := context.Background()

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	c, err := config.BigTableClientFactory(userAgent).NewInstanceAdminClient(project)
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

	cluster := map[string]interface{}{
		"zone":         c.Zone,
		"num_nodes":    c.ServeNodes,
		"cluster_id":   c.Name,
		"storage_type": storageType,
		"kms_key_name": c.KMSKeyName,
	}
	if c.AutoscalingConfig != nil {
		cluster["autoscaling_config"] = make([]map[string]interface{}, 1)
		autoscaling_config := cluster["autoscaling_config"].([]map[string]interface{})
		autoscaling_config[0] = make(map[string]interface{})
		autoscaling_config[0]["min_nodes"] = c.AutoscalingConfig.MinNodes
		autoscaling_config[0]["max_nodes"] = c.AutoscalingConfig.MaxNodes
		autoscaling_config[0]["cpu_target"] = c.AutoscalingConfig.CPUTargetPercent
		autoscaling_config[0]["storage_target"] = c.AutoscalingConfig.StorageUtilizationPerNode
	}
	return cluster
}

func getUnavailableClusterZones(clusters []interface{}, unavailableZones []string) []string {
	var zones []string

	for _, c := range clusters {
		cluster := c.(map[string]interface{})
		zone := cluster["zone"].(string)

		for _, unavailableZone := range unavailableZones {
			if zone == unavailableZone {
				zones = append(zones, zone)
				break
			}
		}
	}
	return zones
}

func expandBigtableClusters(clusters []interface{}, instanceID string, config *transport_tpg.Config) ([]bigtable.ClusterConfig, error) {
	results := make([]bigtable.ClusterConfig, 0, len(clusters))
	for _, c := range clusters {
		cluster := c.(map[string]interface{})
		zone, err := getBigtableZone(cluster["zone"].(string), config)
		if err != nil {
			return nil, err
		}
		var storageType bigtable.StorageType
		switch cluster["storage_type"].(string) {
		case "SSD":
			storageType = bigtable.SSD
		case "HDD":
			storageType = bigtable.HDD
		}

		cluster_config := bigtable.ClusterConfig{
			InstanceID:  instanceID,
			Zone:        zone,
			ClusterID:   cluster["cluster_id"].(string),
			NumNodes:    int32(cluster["num_nodes"].(int)),
			StorageType: storageType,
			KMSKeyName:  cluster["kms_key_name"].(string),
		}
		autoscaling_configs := cluster["autoscaling_config"].([]interface{})
		if len(autoscaling_configs) > 0 {
			autoscaling_config := autoscaling_configs[0].(map[string]interface{})
			cluster_config.AutoscalingConfig = &bigtable.AutoscalingConfig{
				MinNodes:                  autoscaling_config["min_nodes"].(int),
				MaxNodes:                  autoscaling_config["max_nodes"].(int),
				CPUTargetPercent:          autoscaling_config["cpu_target"].(int),
				StorageUtilizationPerNode: autoscaling_config["storage_target"].(int),
			}
		}
		results = append(results, cluster_config)
	}
	return results, nil
}

// getBigtableZone reads the "zone" value from the given resource data and falls back
// to provider's value if not given.  If neither is provided, returns an error.
func getBigtableZone(z string, config *transport_tpg.Config) (string, error) {
	if z == "" {
		if config.Zone != "" {
			return config.Zone, nil
		}
		return "", fmt.Errorf("cannot determine zone: set in cluster.0.zone, or set provider-level zone")
	}
	return tpgresource.GetResourceNameFromSelfLink(z), nil
}

// resourceBigtableInstanceUniqueClusterID asserts cluster ID uniqueness.
func resourceBigtableInstanceUniqueClusterID(_ context.Context, diff *schema.ResourceDiff, meta interface{}) error {
	_, newCount := diff.GetChange("cluster.#")
	clusters := map[string]bool{}

	for i := 0; i < newCount.(int); i++ {
		_, newId := diff.GetChange(fmt.Sprintf("cluster.%d.cluster_id", i))
		clusterID := newId.(string)
		if clusters[clusterID] {
			return fmt.Errorf("duplicated cluster_id: %q", clusterID)
		}
		clusters[clusterID] = true
	}

	return nil
}

// resourceBigtableInstanceClusterReorderTypeList causes the cluster block to
// act like a TypeSet while it's a TypeList underneath. It preserves state
// ordering on updates, and causes the resource to get recreated if it would
// attempt to perform an impossible change.
// This doesn't use the standard unordered list utility (https://github.com/GoogleCloudPlatform/magic-modules/blob/main/templates/terraform/unordered_list_customize_diff.erb)
// because some fields can't be modified using the API and we recreate the instance
// when they're changed.
func resourceBigtableInstanceClusterReorderTypeList(_ context.Context, diff *schema.ResourceDiff, meta interface{}) error {
	oldCount, newCount := diff.GetChange("cluster.#")

	// Simulate Required:true, MinItems:1 for "cluster". This doesn't work
	// when the whole `cluster` field is removed on update.
	if newCount.(int) < 1 {
		return fmt.Errorf("config is invalid: Too few cluster blocks: Should have at least 1 \"cluster\" block")
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
				break
			}
		}
	}

	err := diff.SetNew("cluster", orderedClusters)
	if err != nil {
		return fmt.Errorf("Error setting cluster diff: %s", err)
	}

	// Clusters can't have their zone, storage_type or kms_key_name updated,
	// ForceNew if it's changed. This will show a diff with the old state on
	// the left side and the unmodified new state on the right and the ForceNew
	// attributed to the _old state index_ even if the diff appears to have moved.
	// This depends on the clusters having been reordered already by the prior
	// SetNew call.
	// We've implemented it here because it doesn't return an error in the
	// client and silently fails.
	for i := 0; i < newCount.(int); i++ {
		oldId, newId := diff.GetChange(fmt.Sprintf("cluster.%d.cluster_id", i))
		if oldId != newId {
			continue
		}

		// ForceNew only if the old and the new clusters have the matching cluster ID.
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

		oKey, nKey := diff.GetChange(fmt.Sprintf("cluster.%d.kms_key_name", i))
		if oKey != nKey {
			err := diff.ForceNew(fmt.Sprintf("cluster.%d.kms_key_name", i))
			if err != nil {
				return fmt.Errorf("Error setting cluster diff: %s", err)
			}
		}
	}

	return nil
}

func resourceBigtableInstanceImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"projects/(?P<project>[^/]+)/instances/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<name>[^/]+)",
		"(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/instances/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
