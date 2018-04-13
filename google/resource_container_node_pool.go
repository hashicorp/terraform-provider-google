package google

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	containerBeta "google.golang.org/api/container/v1beta1"
)

func resourceContainerNodePool() *schema.Resource {
	return &schema.Resource{
		Create: resourceContainerNodePoolCreate,
		Read:   resourceContainerNodePoolRead,
		Update: resourceContainerNodePoolUpdate,
		Delete: resourceContainerNodePoolDelete,
		Exists: resourceContainerNodePoolExists,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		SchemaVersion: 1,
		MigrateState:  resourceContainerNodePoolMigrateState,

		Importer: &schema.ResourceImporter{
			State: resourceContainerNodePoolStateImporter,
		},

		Schema: mergeSchemas(
			schemaNodePool,
			map[string]*schema.Schema{
				"project": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
				"zone": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
				"cluster": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
				},
				"region": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
				},
			}),
	}
}

var schemaNodePool = map[string]*schema.Schema{
	"autoscaling": &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"min_node_count": &schema.Schema{
					Type:         schema.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntAtLeast(0),
				},

				"max_node_count": &schema.Schema{
					Type:         schema.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntAtLeast(1),
				},
			},
		},
	},

	"initial_node_count": &schema.Schema{
		Type:     schema.TypeInt,
		Optional: true,
		ForceNew: true,
		Computed: true,
	},

	"instance_group_urls": {
		Type:     schema.TypeList,
		Computed: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	},

	"management": {
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"auto_repair": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},

				"auto_upgrade": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
			},
		},
	},

	"name": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
		ForceNew: true,
	},

	"name_prefix": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
		ForceNew: true,
		Deprecated: "Use the random provider instead. See migration instructions at " +
			"https://github.com/terraform-providers/terraform-provider-google/issues/1054#issuecomment-377390209",
	},

	"node_config": schemaNodeConfig,

	"node_count": {
		Type:         schema.TypeInt,
		Optional:     true,
		Computed:     true,
		ValidateFunc: validation.IntAtLeast(0),
	},

	"version": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	},
}

func getNodePoolRegion(d TerraformResourceData, config *Config) (string, error) {
	res, ok := d.GetOk("region")
	if !ok {
		if config.Zone != "" {
			return config.Zone, nil
		}
		return "", fmt.Errorf("need to set region")
	}
	return GetResourceNameFromSelfLink(res.(string)), nil
}

func generateLocation(d TerraformResourceData, config *Config) (string, error) {
	region, _ := getNodePoolRegion(d, config)
	zone, _ := getZone(d, config)

	if region == "" && zone == "" {
		return "", fmt.Errorf("need to set region or zone")
	}

	if region != "" && zone != "" {
		return "", fmt.Errorf("must only set region or zone")
	}

	if region != "" {
		return region, nil
	}
	return zone, nil
}

type NodePoolInformation struct {
	project  string
	nodePool *containerBeta.NodePool
	location string
	cluster  string
}

func (nodePoolInformation *NodePoolInformation) name() string {
	return nodePoolInformation.nodePool.Name
}

func (nodePoolInformation *NodePoolInformation) fullyQualifiedName() string {
	return fmt.Sprintf(
		"projects/%s/locations/%s/clusters/%s/nodePools/%s",
		nodePoolInformation.project,
		nodePoolInformation.location,
		nodePoolInformation.cluster,
		nodePoolInformation.name(),
	)
}

func (nodePoolInformation *NodePoolInformation) parent() string {
	return fmt.Sprintf(
		"projects/%s/locations/%s/clusters/%s",
		nodePoolInformation.project,
		nodePoolInformation.location,
		nodePoolInformation.cluster,
	)
}

func extractNodePoolInformation(d *schema.ResourceData, config *Config) (*NodePoolInformation, error) {
	project, err := getProject(d, config)
	if err != nil {
		return nil, err
	}

	nodePool, err := expandNodePool(d, "")
	if err != nil {
		return nil, err
	}

	location, err := generateLocation(d, config)
	if err != nil {
		return nil, err
	}

	return &NodePoolInformation{
		project:  project,
		nodePool: nodePool,
		location: location,
		cluster:  d.Get("cluster").(string),
	}, nil
}

func resourceContainerNodePoolCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	nodePoolInfo, err := extractNodePoolInformation(d, config)
	if err != nil {
		return err
	}

	mutexKV.Lock(containerClusterMutexKey(nodePoolInfo.project, nodePoolInfo.location, nodePoolInfo.cluster))
	defer mutexKV.Unlock(containerClusterMutexKey(nodePoolInfo.project, nodePoolInfo.location, nodePoolInfo.cluster))

	req := &containerBeta.CreateNodePoolRequest{
		NodePool: nodePoolInfo.nodePool,
	}

	operation, err := config.clientContainerBeta.
		Projects.Locations.Clusters.NodePools.Create(nodePoolInfo.parent(), req).Do()

	if err != nil {
		return fmt.Errorf("error creating NodePool: %s", err)
	}

	timeoutInMinutes := int(d.Timeout(schema.TimeoutCreate).Minutes())

	waitErr := containerBetaOperationWait(config,
		operation, nodePoolInfo.project,
		nodePoolInfo.location, "timeout creating GKE NodePool", timeoutInMinutes, 3)

	if waitErr != nil {
		// The resource didn't actually create
		d.SetId("")
		return waitErr
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", nodePoolInfo.location, nodePoolInfo.cluster, nodePoolInfo.name()))

	log.Printf("[INFO] GKE NodePool %s has been created", nodePoolInfo.name())

	return resourceContainerNodePoolRead(d, meta)
}

func resourceContainerNodePoolRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	nodePoolInfo, err := extractNodePoolInformation(d, config)
	if err != nil {
		return err
	}

	nodePool, err := config.clientContainerBeta.Projects.Locations.Clusters.NodePools.Get(nodePoolInfo.fullyQualifiedName()).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("NodePool %q from cluster %q", nodePoolInfo.name(), nodePoolInfo.cluster))
	}

	npMap, err := flattenNodePool(d, config, nodePool, "")
	if err != nil {
		return err
	}

	for k, v := range npMap {
		d.Set(k, v)
	}

	if isZone(nodePoolInfo.location) {
		d.Set("zone", nodePoolInfo.location)
	} else {
		d.Set("region", nodePoolInfo.location)
	}

	d.Set("project", nodePoolInfo.project)

	return nil
}

func resourceContainerNodePoolUpdate(d *schema.ResourceData, meta interface{}) error {
	cluster := d.Get("cluster").(string)
	timeoutInMinutes := int(d.Timeout(schema.TimeoutUpdate).Minutes())

	d.Partial(true)
	if err := nodePoolUpdate(d, meta, cluster, "", timeoutInMinutes); err != nil {
		return err
	}
	d.Partial(false)

	return resourceContainerNodePoolRead(d, meta)
}

func resourceContainerNodePoolDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	nodePoolInfo, err := extractNodePoolInformation(d, config)
	if err != nil {
		return err
	}

	timeoutInMinutes := int(d.Timeout(schema.TimeoutDelete).Minutes())

	mutexKV.Lock(containerClusterMutexKey(nodePoolInfo.project, nodePoolInfo.location, nodePoolInfo.cluster))
	defer mutexKV.Unlock(containerClusterMutexKey(nodePoolInfo.project, nodePoolInfo.location, nodePoolInfo.cluster))

	op, err := config.clientContainerBeta.Projects.Locations.Clusters.NodePools.Delete(nodePoolInfo.fullyQualifiedName()).Do()
	if err != nil {
		return fmt.Errorf("Error deleting NodePool: %s", err)
	}

	// Wait until it's deleted
	waitErr := containerBetaOperationWait(config, op, nodePoolInfo.project, nodePoolInfo.location, "deleting GKE NodePool", timeoutInMinutes, 2)
	if waitErr != nil {
		return waitErr
	}

	log.Printf("[INFO] GKE NodePool %s has been deleted", d.Id())

	d.SetId("")

	return nil
}

func resourceContainerNodePoolExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	config := meta.(*Config)

	nodePoolInfo, err := extractNodePoolInformation(d, config)
	if err != nil {
		return false, err
	}

	if err != nil {
		return false, err
	}

	_, err = config.clientContainerBeta.Projects.Locations.Clusters.NodePools.Get(nodePoolInfo.fullyQualifiedName()).Do()
	if err != nil {
		if err = handleNotFoundError(err, d, fmt.Sprintf("Container NodePool %s", nodePoolInfo.name())); err == nil {
			return false, nil
		}
		// There was some other error in reading the resource
		return true, err
	}
	return true, nil
}

func resourceContainerNodePoolStateImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("Invalid container cluster specifier. Expecting {zone}/{cluster}/{name}")
	}

	location := parts[0]
	if isZone(location) {
		d.Set("zone", location)
	} else {
		d.Set("region", location)
	}

	d.Set("cluster", parts[1])
	d.Set("name", parts[2])

	return []*schema.ResourceData{d}, nil
}

func expandNodePool(d *schema.ResourceData, prefix string) (*containerBeta.NodePool, error) {
	var name string
	if v, ok := d.GetOk(prefix + "name"); ok {
		if _, ok := d.GetOk(prefix + "name_prefix"); ok {
			return nil, fmt.Errorf("Cannot specify both name and name_prefix for a node_pool")
		}
		name = v.(string)
	} else if v, ok := d.GetOk(prefix + "name_prefix"); ok {
		name = resource.PrefixedUniqueId(v.(string))
	} else {
		name = resource.UniqueId()
	}

	nodeCount := 0
	if initialNodeCount, ok := d.GetOk(prefix + "initial_node_count"); ok {
		nodeCount = initialNodeCount.(int)
	}
	if nc, ok := d.GetOk(prefix + "node_count"); ok {
		if nodeCount != 0 {
			return nil, fmt.Errorf("Cannot set both initial_node_count and node_count on node pool %s", name)
		}
		nodeCount = nc.(int)
	}

	np := &containerBeta.NodePool{
		Name:             name,
		InitialNodeCount: int64(nodeCount),
		Config:           expandNodeConfig(d.Get(prefix + "node_config")),
		Version:          d.Get(prefix + "version").(string),
	}

	if v, ok := d.GetOk(prefix + "autoscaling"); ok {
		autoscaling := v.([]interface{})[0].(map[string]interface{})
		np.Autoscaling = &containerBeta.NodePoolAutoscaling{
			Enabled:         true,
			MinNodeCount:    int64(autoscaling["min_node_count"].(int)),
			MaxNodeCount:    int64(autoscaling["max_node_count"].(int)),
			ForceSendFields: []string{"MinNodeCount"},
		}
	}

	if v, ok := d.GetOk(prefix + "management"); ok {
		managementConfig := v.([]interface{})[0].(map[string]interface{})
		np.Management = &containerBeta.NodeManagement{}

		if v, ok := managementConfig["auto_repair"]; ok {
			np.Management.AutoRepair = v.(bool)
		}

		if v, ok := managementConfig["auto_upgrade"]; ok {
			np.Management.AutoUpgrade = v.(bool)
		}
	}

	return np, nil
}

func flattenNodePool(d *schema.ResourceData, config *Config, np *containerBeta.NodePool, prefix string) (map[string]interface{}, error) {
	// Node pools don't expose the current node count in their API, so read the
	// instance groups instead. They should all have the same size, but in case a resize
	// failed or something else strange happened, we'll just use the average size.
	size := 0
	for _, url := range np.InstanceGroupUrls {
		// retrieve instance group manager (InstanceGroupUrls are actually URLs for InstanceGroupManagers)
		matches := instanceGroupManagerURL.FindStringSubmatch(url)
		if len(matches) < 4 {
			return nil, fmt.Errorf("Error reading instance group manage URL '%q'", url)
		}
		igm, err := config.clientCompute.InstanceGroupManagers.Get(matches[1], matches[2], matches[3]).Do()
		if err != nil {
			return nil, fmt.Errorf("Error reading instance group manager returned as an instance group URL: %q", err)
		}
		size += int(igm.TargetSize)
	}
	nodePool := map[string]interface{}{
		"name":                np.Name,
		"name_prefix":         d.Get(prefix + "name_prefix"),
		"initial_node_count":  np.InitialNodeCount,
		"node_count":          size / len(np.InstanceGroupUrls),
		"node_config":         flattenNodeConfig(np.Config),
		"instance_group_urls": np.InstanceGroupUrls,
		"version":             np.Version,
	}

	if np.Autoscaling != nil && np.Autoscaling.Enabled {
		nodePool["autoscaling"] = []map[string]interface{}{
			map[string]interface{}{
				"min_node_count": np.Autoscaling.MinNodeCount,
				"max_node_count": np.Autoscaling.MaxNodeCount,
			},
		}
	}

	nodePool["management"] = []map[string]interface{}{
		{
			"auto_repair":  np.Management.AutoRepair,
			"auto_upgrade": np.Management.AutoUpgrade,
		},
	}

	return nodePool, nil
}

func nodePoolUpdate(d *schema.ResourceData, meta interface{}, clusterName, prefix string, timeoutInMinutes int) error {
	config := meta.(*Config)

	nodePoolInfo, err := extractNodePoolInformation(d, config)
	if err != nil {
		return err
	}

	npName := d.Get(prefix + "name").(string)
	lockKey := containerClusterMutexKey(nodePoolInfo.project, nodePoolInfo.location, clusterName)

	if d.HasChange(prefix + "autoscaling") {
		update := &containerBeta.ClusterUpdate{
			DesiredNodePoolId: npName,
		}
		if v, ok := d.GetOk(prefix + "autoscaling"); ok {
			autoscaling := v.([]interface{})[0].(map[string]interface{})
			update.DesiredNodePoolAutoscaling = &containerBeta.NodePoolAutoscaling{
				Enabled:         true,
				MinNodeCount:    int64(autoscaling["min_node_count"].(int)),
				MaxNodeCount:    int64(autoscaling["max_node_count"].(int)),
				ForceSendFields: []string{"MinNodeCount"},
			}
		} else {
			update.DesiredNodePoolAutoscaling = &containerBeta.NodePoolAutoscaling{
				Enabled: false,
			}
		}

		req := &containerBeta.UpdateClusterRequest{
			Update: update,
		}

		updateF := func() error {
			op, err := config.clientContainerBeta.Projects.Locations.Clusters.Update(nodePoolInfo.fullyQualifiedName(), req).Do()
			if err != nil {
				return err
			}

			// Wait until it's updated
			return containerBetaOperationWait(config, op,
				nodePoolInfo.project,
				nodePoolInfo.location, "updating GKE node pool",
				timeoutInMinutes, 2)
		}

		// Call update serially.
		if err := lockedCall(lockKey, updateF); err != nil {
			return err
		}

		log.Printf("[INFO] Updated autoscaling in Node Pool %s", d.Id())

		if prefix == "" {
			d.SetPartial("autoscaling")
		}
	}

	if d.HasChange(prefix + "node_count") {
		newSize := int64(d.Get(prefix + "node_count").(int))
		req := &containerBeta.SetNodePoolSizeRequest{
			NodeCount: newSize,
		}
		updateF := func() error {
			op, err := config.clientContainerBeta.Projects.Locations.Clusters.NodePools.SetSize(nodePoolInfo.fullyQualifiedName(), req).Do()

			if err != nil {
				return err
			}

			// Wait until it's updated
			return containerBetaOperationWait(config, op,
				nodePoolInfo.project,
				nodePoolInfo.location, "updating GKE node pool size",
				timeoutInMinutes, 2)
		}

		// Call update serially.
		if err := lockedCall(lockKey, updateF); err != nil {
			return err
		}

		log.Printf("[INFO] GKE node pool %s size has been updated to %d", npName, newSize)

		if prefix == "" {
			d.SetPartial("node_count")
		}
	}

	if d.HasChange(prefix + "management") {
		management := &containerBeta.NodeManagement{}
		if v, ok := d.GetOk(prefix + "management"); ok {
			managementConfig := v.([]interface{})[0].(map[string]interface{})
			management.AutoRepair = managementConfig["auto_repair"].(bool)
			management.AutoUpgrade = managementConfig["auto_upgrade"].(bool)
			management.ForceSendFields = []string{"AutoRepair", "AutoUpgrade"}
		}
		req := &containerBeta.SetNodePoolManagementRequest{
			Management: management,
		}

		updateF := func() error {
			op, err := config.clientContainerBeta.Projects.Locations.
				Clusters.NodePools.SetManagement(nodePoolInfo.fullyQualifiedName(), req).Do()

			if err != nil {
				return err
			}

			// Wait until it's updated
			return containerBetaOperationWait(config, op,
				nodePoolInfo.project,
				nodePoolInfo.location, "updating GKE node pool management", timeoutInMinutes, 2)
		}

		// Call update serially.
		if err := lockedCall(lockKey, updateF); err != nil {
			return err
		}

		log.Printf("[INFO] Updated management in Node Pool %s", npName)

		if prefix == "" {
			d.SetPartial("management")
		}
	}

	if d.HasChange(prefix + "version") {
		req := &containerBeta.UpdateNodePoolRequest{
			NodeVersion: d.Get("version").(string),
		}
		updateF := func() error {
			op, err := config.clientContainerBeta.Projects.
				Locations.Clusters.NodePools.Update(nodePoolInfo.fullyQualifiedName(), req).Do()

			if err != nil {
				return err
			}

			// Wait until it's updated
			return containerBetaOperationWait(config, op,
				nodePoolInfo.project,
				nodePoolInfo.location, "updating GKE node pool version", timeoutInMinutes, 2)
		}

		// Call update serially.
		if err := lockedCall(lockKey, updateF); err != nil {
			return err
		}

		log.Printf("[INFO] Updated version in Node Pool %s", npName)

		if prefix == "" {
			d.SetPartial("version")
		}
	}

	return nil
}
