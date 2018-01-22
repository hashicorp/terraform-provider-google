package google

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"google.golang.org/api/container/v1"
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
					Required: true,
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
		Type:       schema.TypeInt,
		Optional:   true,
		ForceNew:   true,
		Computed:   true,
		Deprecated: "Use node_count instead",
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
		ForceNew: true,
	},

	"node_config": schemaNodeConfig,

	"node_count": {
		Type:         schema.TypeInt,
		Optional:     true,
		Computed:     true,
		ValidateFunc: validation.IntAtLeast(0),
	},
}

func resourceContainerNodePoolCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	nodePool, err := expandNodePool(d, "")
	if err != nil {
		return err
	}

	req := &container.CreateNodePoolRequest{
		NodePool: nodePool,
	}

	zone, err := getZone(d, config)
	if err != nil {
		return err
	}
	cluster := d.Get("cluster").(string)

	mutexKV.Lock(containerClusterMutexKey(project, zone, cluster))
	defer mutexKV.Unlock(containerClusterMutexKey(project, zone, cluster))
	op, err := config.clientContainer.Projects.Zones.Clusters.NodePools.Create(project, zone, cluster, req).Do()

	if err != nil {
		return fmt.Errorf("Error creating NodePool: %s", err)
	}

	timeoutInMinutes := int(d.Timeout(schema.TimeoutCreate).Minutes())
	waitErr := containerOperationWait(config, op, project, zone, "creating GKE NodePool", timeoutInMinutes, 3)
	if waitErr != nil {
		// The resource didn't actually create
		d.SetId("")
		return waitErr
	}

	log.Printf("[INFO] GKE NodePool %s has been created", nodePool.Name)

	d.SetId(fmt.Sprintf("%s/%s/%s", zone, cluster, nodePool.Name))

	return resourceContainerNodePoolRead(d, meta)
}

func resourceContainerNodePoolRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	zone, err := getZone(d, config)
	if err != nil {
		return err
	}
	cluster := d.Get("cluster").(string)
	name := getNodePoolName(d.Id())

	nodePool, err := config.clientContainer.Projects.Zones.Clusters.NodePools.Get(
		project, zone, cluster, name).Do()
	if err != nil {
		return fmt.Errorf("Error reading NodePool: %s", err)
	}

	npMap, err := flattenNodePool(d, config, nodePool, "")
	if err != nil {
		return err
	}

	for k, v := range npMap {
		d.Set(k, v)
	}

	d.Set("zone", zone)
	d.Set("project", project)

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

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	zone, err := getZone(d, config)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	cluster := d.Get("cluster").(string)
	timeoutInMinutes := int(d.Timeout(schema.TimeoutDelete).Minutes())

	mutexKV.Lock(containerClusterMutexKey(project, zone, cluster))
	defer mutexKV.Unlock(containerClusterMutexKey(project, zone, cluster))
	op, err := config.clientContainer.Projects.Zones.Clusters.NodePools.Delete(
		project, zone, cluster, name).Do()
	if err != nil {
		return fmt.Errorf("Error deleting NodePool: %s", err)
	}

	// Wait until it's deleted
	waitErr := containerOperationWait(config, op, project, zone, "deleting GKE NodePool", timeoutInMinutes, 2)
	if waitErr != nil {
		return waitErr
	}

	log.Printf("[INFO] GKE NodePool %s has been deleted", d.Id())

	d.SetId("")

	return nil
}

func resourceContainerNodePoolExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return false, err
	}

	zone, err := getZone(d, config)
	if err != nil {
		return false, err
	}
	cluster := d.Get("cluster").(string)
	name := getNodePoolName(d.Id())

	_, err = config.clientContainer.Projects.Zones.Clusters.NodePools.Get(
		project, zone, cluster, name).Do()
	if err != nil {
		if err = handleNotFoundError(err, d, fmt.Sprintf("Container NodePool %s", name)); err == nil {
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

	d.Set("zone", parts[0])
	d.Set("cluster", parts[1])
	d.Set("name", parts[2])

	return []*schema.ResourceData{d}, nil
}

func expandNodePool(d *schema.ResourceData, prefix string) (*container.NodePool, error) {
	var name string
	if v, ok := d.GetOk(prefix + "name"); ok {
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

	np := &container.NodePool{
		Name:             name,
		InitialNodeCount: int64(nodeCount),
		Config:           expandNodeConfig(d.Get(prefix + "node_config")),
	}

	if v, ok := d.GetOk(prefix + "autoscaling"); ok {
		autoscaling := v.([]interface{})[0].(map[string]interface{})
		np.Autoscaling = &container.NodePoolAutoscaling{
			Enabled:         true,
			MinNodeCount:    int64(autoscaling["min_node_count"].(int)),
			MaxNodeCount:    int64(autoscaling["max_node_count"].(int)),
			ForceSendFields: []string{"MinNodeCount"},
		}
	}

	if v, ok := d.GetOk(prefix + "management"); ok {
		managementConfig := v.([]interface{})[0].(map[string]interface{})
		np.Management = &container.NodeManagement{}

		if v, ok := managementConfig["auto_repair"]; ok {
			np.Management.AutoRepair = v.(bool)
		}

		if v, ok := managementConfig["auto_upgrade"]; ok {
			np.Management.AutoUpgrade = v.(bool)
		}
	}

	return np, nil
}

func flattenNodePool(d *schema.ResourceData, config *Config, np *container.NodePool, prefix string) (map[string]interface{}, error) {
	// Node pools don't expose the current node count in their API, so read the
	// instance groups instead. They should all have the same size, but in case a resize
	// failed or something else strange happened, we'll just use the average size.
	size := 0
	for _, url := range np.InstanceGroupUrls {
		// retrieve instance group manager (InstanceGroupUrls are actually URLs for InstanceGroupManagers)
		matches := instanceGroupManagerURL.FindStringSubmatch(url)
		igm, err := config.clientCompute.InstanceGroupManagers.Get(matches[1], matches[2], matches[3]).Do()
		if err != nil {
			return nil, fmt.Errorf("Error reading instance group manager returned as an instance group URL: %s", err)
		}
		size += int(igm.TargetSize)
	}
	nodePool := map[string]interface{}{
		"name":               np.Name,
		"name_prefix":        d.Get(prefix + "name_prefix"),
		"initial_node_count": np.InitialNodeCount,
		"node_count":         size / len(np.InstanceGroupUrls),
		"node_config":        flattenNodeConfig(np.Config),
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

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	zone, err := getZone(d, config)
	if err != nil {
		return err
	}

	npName := d.Get(prefix + "name").(string)

	if d.HasChange(prefix + "autoscaling") {
		update := &container.ClusterUpdate{
			DesiredNodePoolId: npName,
		}
		if v, ok := d.GetOk(prefix + "autoscaling"); ok {
			autoscaling := v.([]interface{})[0].(map[string]interface{})
			update.DesiredNodePoolAutoscaling = &container.NodePoolAutoscaling{
				Enabled:         true,
				MinNodeCount:    int64(autoscaling["min_node_count"].(int)),
				MaxNodeCount:    int64(autoscaling["max_node_count"].(int)),
				ForceSendFields: []string{"MinNodeCount"},
			}
		} else {
			update.DesiredNodePoolAutoscaling = &container.NodePoolAutoscaling{
				Enabled: false,
			}
		}

		req := &container.UpdateClusterRequest{
			Update: update,
		}
		mutexKV.Lock(containerClusterMutexKey(project, zone, clusterName))
		defer mutexKV.Unlock(containerClusterMutexKey(project, zone, clusterName))
		op, err := config.clientContainer.Projects.Zones.Clusters.Update(
			project, zone, clusterName, req).Do()
		if err != nil {
			return err
		}

		// Wait until it's updated
		waitErr := containerOperationWait(config, op, project, zone, "updating GKE node pool", timeoutInMinutes, 2)
		if waitErr != nil {
			return waitErr
		}

		log.Printf("[INFO] Updated autoscaling in Node Pool %s", d.Id())

		if prefix == "" {
			d.SetPartial("autoscaling")
		}
	}

	if d.HasChange(prefix + "node_count") {
		newSize := int64(d.Get(prefix + "node_count").(int))
		req := &container.SetNodePoolSizeRequest{
			NodeCount: newSize,
		}
		mutexKV.Lock(containerClusterMutexKey(project, zone, clusterName))
		defer mutexKV.Unlock(containerClusterMutexKey(project, zone, clusterName))
		op, err := config.clientContainer.Projects.Zones.Clusters.NodePools.SetSize(project, zone, clusterName, npName, req).Do()
		if err != nil {
			return err
		}

		// Wait until it's updated
		waitErr := containerOperationWait(config, op, project, zone, "updating GKE node pool size", timeoutInMinutes, 2)
		if waitErr != nil {
			return waitErr
		}

		log.Printf("[INFO] GKE node pool %s size has been updated to %d", npName, newSize)

		if prefix == "" {
			d.SetPartial("node_count")
		}
	}

	if d.HasChange(prefix + "management") {
		management := &container.NodeManagement{}
		if v, ok := d.GetOk(prefix + "management"); ok {
			managementConfig := v.([]interface{})[0].(map[string]interface{})
			management.AutoRepair = managementConfig["auto_repair"].(bool)
			management.AutoUpgrade = managementConfig["auto_upgrade"].(bool)
			management.ForceSendFields = []string{"AutoRepair", "AutoUpgrade"}
		}
		req := &container.SetNodePoolManagementRequest{
			Management: management,
		}
		mutexKV.Lock(containerClusterMutexKey(project, zone, clusterName))
		defer mutexKV.Unlock(containerClusterMutexKey(project, zone, clusterName))
		op, err := config.clientContainer.Projects.Zones.Clusters.NodePools.SetManagement(
			project, zone, clusterName, npName, req).Do()

		if err != nil {
			return err
		}

		// Wait until it's updated
		waitErr := containerOperationWait(config, op, project, zone, "updating GKE node pool management", timeoutInMinutes, 2)
		if waitErr != nil {
			return waitErr
		}

		log.Printf("[INFO] Updated management in Node Pool %s", npName)

		if prefix == "" {
			d.SetPartial("management")
		}
	}

	return nil
}

func getNodePoolName(id string) string {
	// name can be specified with name, name_prefix, or neither, so read it from the id.
	return strings.Split(id, "/")[2]
}
