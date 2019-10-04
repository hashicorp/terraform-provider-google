package google

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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

		CustomizeDiff: customdiff.All(
			resourceNodeConfigEmptyGuestAccelerator,
		),

		Schema: mergeSchemas(
			schemaNodePool,
			map[string]*schema.Schema{
				"project": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
				"cluster": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
				},
				"zone": {
					Type:       schema.TypeString,
					Optional:   true,
					Computed:   true,
					Deprecated: "use location instead",
					ForceNew:   true,
				},
				"region": {
					Type:       schema.TypeString,
					Optional:   true,
					Computed:   true,
					Deprecated: "use location instead",
					ForceNew:   true,
				},
				"location": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
			}),
	}
}

var schemaNodePool = map[string]*schema.Schema{
	"autoscaling": {
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"min_node_count": {
					Type:         schema.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntAtLeast(0),
				},

				"max_node_count": {
					Type:         schema.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntAtLeast(1),
				},
			},
		},
	},

	"max_pods_per_node": {
		Type:     schema.TypeInt,
		Optional: true,
		ForceNew: true,
		Computed: true,
	},

	"initial_node_count": {
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

	"name": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
		ForceNew: true,
	},

	"name_prefix": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
		ForceNew: true,
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

type NodePoolInformation struct {
	project  string
	location string
	cluster  string
}

func (nodePoolInformation *NodePoolInformation) fullyQualifiedName(nodeName string) string {
	return fmt.Sprintf(
		"projects/%s/locations/%s/clusters/%s/nodePools/%s",
		nodePoolInformation.project,
		nodePoolInformation.location,
		nodePoolInformation.cluster,
		nodeName,
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

func (nodePoolInformation *NodePoolInformation) lockKey() string {
	return containerClusterMutexKey(nodePoolInformation.project,
		nodePoolInformation.location, nodePoolInformation.cluster)
}

func extractNodePoolInformation(d *schema.ResourceData, config *Config) (*NodePoolInformation, error) {
	project, err := getProject(d, config)
	if err != nil {
		return nil, err
	}

	location, err := getLocation(d, config)
	if err != nil {
		return nil, err
	}

	return &NodePoolInformation{
		project:  project,
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

	nodePool, err := expandNodePool(d, "")
	if err != nil {
		return err
	}

	mutexKV.Lock(nodePoolInfo.lockKey())
	defer mutexKV.Unlock(nodePoolInfo.lockKey())

	req := &containerBeta.CreateNodePoolRequest{
		NodePool: nodePool,
	}

	timeout := d.Timeout(schema.TimeoutCreate)
	startTime := time.Now()

	var operation *containerBeta.Operation
	err = resource.Retry(timeout, func() *resource.RetryError {
		operation, err = config.clientContainerBeta.
			Projects.Locations.Clusters.NodePools.Create(nodePoolInfo.parent(), req).Do()

		if err != nil {
			if isFailedPreconditionError(err) {
				// We get failed precondition errors if the cluster is updating
				// while we try to add the node pool.
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error creating NodePool: %s", err)
	}
	timeout -= time.Since(startTime)

	d.SetId(fmt.Sprintf("%s/%s/%s", nodePoolInfo.location, nodePoolInfo.cluster, nodePool.Name))

	waitErr := containerOperationWait(config,
		operation, nodePoolInfo.project,
		nodePoolInfo.location, "creating GKE NodePool", int(timeout.Minutes()))

	if waitErr != nil {
		// The resource didn't actually create
		d.SetId("")
		return waitErr
	}

	log.Printf("[INFO] GKE NodePool %s has been created", nodePool.Name)

	return resourceContainerNodePoolRead(d, meta)
}

func resourceContainerNodePoolRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	nodePoolInfo, err := extractNodePoolInformation(d, config)

	name := getNodePoolName(d.Id())

	if err != nil {
		return err
	}

	var nodePool = &containerBeta.NodePool{}
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		nodePool, err = config.clientContainerBeta.
			Projects.Locations.Clusters.NodePools.Get(nodePoolInfo.fullyQualifiedName(name)).Do()

		if err != nil {
			return resource.NonRetryableError(err)
		}
		if nodePool.Status != "RUNNING" {
			return resource.RetryableError(fmt.Errorf("Nodepool %q has status %q with message %q", d.Get("name"), nodePool.Status, nodePool.StatusMessage))
		}
		return nil
	})

	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("NodePool %q from cluster %q", name, nodePoolInfo.cluster))
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

	d.Set("location", nodePoolInfo.location)
	d.Set("project", nodePoolInfo.project)

	return nil
}

func resourceContainerNodePoolUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	timeoutInMinutes := int(d.Timeout(schema.TimeoutUpdate).Minutes())

	nodePoolInfo, err := extractNodePoolInformation(d, config)
	if err != nil {
		return err
	}

	d.Partial(true)
	if err := nodePoolUpdate(d, meta, nodePoolInfo, "", timeoutInMinutes); err != nil {
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

	name := getNodePoolName(d.Id())

	timeoutInMinutes := int(d.Timeout(schema.TimeoutDelete).Minutes())

	mutexKV.Lock(nodePoolInfo.lockKey())
	defer mutexKV.Unlock(nodePoolInfo.lockKey())

	var op = &containerBeta.Operation{}
	var count = 0
	err = resource.Retry(30*time.Second, func() *resource.RetryError {
		count++
		op, err = config.clientContainerBeta.Projects.Locations.
			Clusters.NodePools.Delete(nodePoolInfo.fullyQualifiedName(name)).Do()

		if err != nil {
			return resource.RetryableError(err)
		}

		if count == 15 {
			return resource.NonRetryableError(fmt.Errorf("Error retrying to delete node pool %s", name))
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error deleting NodePool: %s", err)
	}

	// Wait until it's deleted
	waitErr := containerOperationWait(config, op, nodePoolInfo.project, nodePoolInfo.location, "deleting GKE NodePool", timeoutInMinutes)
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

	name := getNodePoolName(d.Id())

	_, err = config.clientContainerBeta.Projects.Locations.Clusters.NodePools.Get(nodePoolInfo.fullyQualifiedName(name)).Do()
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

	switch len(parts) {
	case 3:
		location := parts[0]
		if isZone(location) {
			d.Set("zone", location)
		} else {
			d.Set("region", location)
		}

		d.Set("location", location)
		d.Set("cluster", parts[1])
		d.Set("name", parts[2])
	case 4:
		d.Set("project", parts[0])

		location := parts[1]
		if isZone(location) {
			d.Set("zone", location)
		} else {
			d.Set("region", location)
		}

		d.Set("location", location)
		d.Set("cluster", parts[2])
		d.Set("name", parts[3])

		// override the inputted ID with the <location>/<cluster>/<name> format
		d.SetId(strings.Join(parts[1:], "/"))
	default:
		return nil, fmt.Errorf("Invalid container cluster specifier. Expecting {location}/{cluster}/{name} or {project}/{location}/{cluster}/{name}")
	}

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

	if v, ok := d.GetOk(prefix + "max_pods_per_node"); ok {
		np.MaxPodsConstraint = &containerBeta.MaxPodsConstraint{
			MaxPodsPerNode: int64(v.(int)),
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
		igm, err := config.clientComputeBeta.InstanceGroupManagers.Get(matches[1], matches[2], matches[3]).Do()
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

	if np.Autoscaling != nil {
		if np.Autoscaling.Enabled {
			nodePool["autoscaling"] = []map[string]interface{}{
				{
					"min_node_count": np.Autoscaling.MinNodeCount,
					"max_node_count": np.Autoscaling.MaxNodeCount,
				},
			}
		} else {
			nodePool["autoscaling"] = []map[string]interface{}{}
		}
	}

	if np.MaxPodsConstraint != nil {
		nodePool["max_pods_per_node"] = np.MaxPodsConstraint.MaxPodsPerNode
	}

	nodePool["management"] = []map[string]interface{}{
		{
			"auto_repair":  np.Management.AutoRepair,
			"auto_upgrade": np.Management.AutoUpgrade,
		},
	}

	return nodePool, nil
}

func nodePoolUpdate(d *schema.ResourceData, meta interface{}, nodePoolInfo *NodePoolInformation, prefix string, timeoutInMinutes int) error {
	config := meta.(*Config)

	name := d.Get(prefix + "name").(string)

	lockKey := nodePoolInfo.lockKey()

	if d.HasChange(prefix + "autoscaling") {
		update := &containerBeta.ClusterUpdate{
			DesiredNodePoolId: name,
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
			op, err := config.clientContainerBeta.Projects.Locations.Clusters.Update(nodePoolInfo.parent(), req).Do()
			if err != nil {
				return err
			}

			// Wait until it's updated
			return containerOperationWait(config, op,
				nodePoolInfo.project,
				nodePoolInfo.location, "updating GKE node pool",
				timeoutInMinutes)
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

	if d.HasChange(prefix + "node_config") {
		if d.HasChange(prefix + "node_config.0.image_type") {
			req := &containerBeta.UpdateClusterRequest{
				Update: &containerBeta.ClusterUpdate{
					DesiredNodePoolId: name,
					DesiredImageType:  d.Get(prefix + "node_config.0.image_type").(string),
				},
			}

			updateF := func() error {
				op, err := config.clientContainerBeta.Projects.Locations.Clusters.Update(nodePoolInfo.parent(), req).Do()
				if err != nil {
					return err
				}

				// Wait until it's updated
				return containerOperationWait(config, op,
					nodePoolInfo.project,
					nodePoolInfo.location, "updating GKE node pool",
					timeoutInMinutes)
			}

			// Call update serially.
			if err := lockedCall(lockKey, updateF); err != nil {
				return err
			}

			log.Printf("[INFO] Updated image type in Node Pool %s", d.Id())
		}

		if prefix == "" {
			d.SetPartial("node_config")
		}
	}

	if d.HasChange(prefix + "node_count") {
		newSize := int64(d.Get(prefix + "node_count").(int))
		req := &containerBeta.SetNodePoolSizeRequest{
			NodeCount: newSize,
		}
		updateF := func() error {
			op, err := config.clientContainerBeta.Projects.Locations.Clusters.NodePools.SetSize(nodePoolInfo.fullyQualifiedName(name), req).Do()

			if err != nil {
				return err
			}

			// Wait until it's updated
			return containerOperationWait(config, op,
				nodePoolInfo.project,
				nodePoolInfo.location, "updating GKE node pool size",
				timeoutInMinutes)
		}

		// Call update serially.
		if err := lockedCall(lockKey, updateF); err != nil {
			return err
		}

		log.Printf("[INFO] GKE node pool %s size has been updated to %d", name, newSize)

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
				Clusters.NodePools.SetManagement(nodePoolInfo.fullyQualifiedName(name), req).Do()

			if err != nil {
				return err
			}

			// Wait until it's updated
			return containerOperationWait(config, op,
				nodePoolInfo.project,
				nodePoolInfo.location, "updating GKE node pool management", timeoutInMinutes)
		}

		// Call update serially.
		if err := lockedCall(lockKey, updateF); err != nil {
			return err
		}

		log.Printf("[INFO] Updated management in Node Pool %s", name)

		if prefix == "" {
			d.SetPartial("management")
		}
	}

	if d.HasChange(prefix + "version") {
		req := &containerBeta.UpdateNodePoolRequest{
			NodePoolId:  name,
			NodeVersion: d.Get(prefix + "version").(string),
		}
		updateF := func() error {
			op, err := config.clientContainerBeta.Projects.
				Locations.Clusters.NodePools.Update(nodePoolInfo.fullyQualifiedName(name), req).Do()

			if err != nil {
				return err
			}

			// Wait until it's updated
			return containerOperationWait(config, op,
				nodePoolInfo.project,
				nodePoolInfo.location, "updating GKE node pool version", timeoutInMinutes)
		}

		// Call update serially.
		if err := lockedCall(lockKey, updateF); err != nil {
			return err
		}

		log.Printf("[INFO] Updated version in Node Pool %s", name)

		if prefix == "" {
			d.SetPartial("version")
		}
	}

	return nil
}

func getNodePoolName(id string) string {
	// name can be specified with name, name_prefix, or neither, so read it from the id.
	return strings.Split(id, "/")[2]
}
