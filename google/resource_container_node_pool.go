package google

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"google.golang.org/api/container/v1"
)

var clusterIdRegex = regexp.MustCompile("projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/clusters/(?P<name>[^/]+)")

func resourceContainerNodePool() *schema.Resource {
	return &schema.Resource{
		Create: resourceContainerNodePoolCreate,
		Read:   resourceContainerNodePoolRead,
		Update: resourceContainerNodePoolUpdate,
		Delete: resourceContainerNodePoolDelete,
		Exists: resourceContainerNodePoolExists,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		SchemaVersion: 1,
		MigrateState:  resourceContainerNodePoolMigrateState,

		Importer: &schema.ResourceImporter{
			State: resourceContainerNodePoolStateImporter,
		},

		CustomizeDiff: customdiff.All(
			resourceNodeConfigEmptyGuestAccelerator,
		),

		UseJSONNumber: true,

		Schema: mergeSchemas(
			schemaNodePool,
			map[string]*schema.Schema{
				"project": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
					ForceNew:    true,
					Description: `The ID of the project in which to create the node pool. If blank, the provider-configured project will be used.`,
				},
				"cluster": {
					Type:        schema.TypeString,
					Required:    true,
					ForceNew:    true,
					Description: `The cluster to create the node pool for. Cluster must be present in location provided for zonal clusters.`,
				},
				"location": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
					ForceNew:    true,
					Description: `The location (region or zone) of the cluster.`,
				},
				"operation": {
					Type:     schema.TypeString,
					Computed: true,
				},
			}),
	}
}

var schemaNodePool = map[string]*schema.Schema{
	"autoscaling": {
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: `Configuration required by cluster autoscaler to adjust the size of the node pool to the current cluster usage.`,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"min_node_count": {
					Type:         schema.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntAtLeast(0),
					Description:  `Minimum number of nodes in the NodePool. Must be >=0 and <= max_node_count.`,
				},

				"max_node_count": {
					Type:         schema.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntAtLeast(1),
					Description:  `Maximum number of nodes in the NodePool. Must be >= min_node_count.`,
				},
			},
		},
	},

	"max_pods_per_node": {
		Type:        schema.TypeInt,
		Optional:    true,
		ForceNew:    true,
		Computed:    true,
		Description: `The maximum number of pods per node in this node pool. Note that this does not work on node pools which are "route-based" - that is, node pools belonging to clusters that do not have IP Aliasing enabled.`,
	},

	"node_locations": {
		Type:        schema.TypeSet,
		Optional:    true,
		Computed:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: `The list of zones in which the node pool's nodes should be located. Nodes must be in the region of their regional cluster or in the same region as their cluster's zone for zonal clusters. If unspecified, the cluster-level node_locations will be used.`,
	},

	"upgrade_settings": {
		Type:        schema.TypeList,
		Optional:    true,
		Computed:    true,
		MaxItems:    1,
		Description: `Specify node upgrade settings to change how many nodes GKE attempts to upgrade at once. The number of nodes upgraded simultaneously is the sum of max_surge and max_unavailable. The maximum number of nodes upgraded simultaneously is limited to 20.`,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"max_surge": {
					Type:         schema.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntAtLeast(0),
					Description:  `The number of additional nodes that can be added to the node pool during an upgrade. Increasing max_surge raises the number of nodes that can be upgraded simultaneously. Can be set to 0 or greater.`,
				},

				"max_unavailable": {
					Type:         schema.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntAtLeast(0),
					Description:  `The number of nodes that can be simultaneously unavailable during an upgrade. Increasing max_unavailable raises the number of nodes that can be upgraded in parallel. Can be set to 0 or greater.`,
				},
			},
		},
	},

	"initial_node_count": {
		Type:        schema.TypeInt,
		Optional:    true,
		ForceNew:    true,
		Computed:    true,
		Description: `The initial number of nodes for the pool. In regional or multi-zonal clusters, this is the number of nodes per zone. Changing this will force recreation of the resource.`,
	},

	"instance_group_urls": {
		Type:        schema.TypeList,
		Computed:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: `The resource URLs of the managed instance groups associated with this node pool.`,
	},

	"managed_instance_group_urls": {
		Type:        schema.TypeList,
		Computed:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: `List of instance group URLs which have been assigned to this node pool.`,
	},

	"management": {
		Type:        schema.TypeList,
		Optional:    true,
		Computed:    true,
		MaxItems:    1,
		Description: `Node management configuration, wherein auto-repair and auto-upgrade is configured.`,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"auto_repair": {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     false,
					Description: `Whether the nodes will be automatically repaired.`,
				},

				"auto_upgrade": {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     false,
					Description: `Whether the nodes will be automatically upgraded.`,
				},
			},
		},
	},

	"name": {
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		ForceNew:    true,
		Description: `The name of the node pool. If left blank, Terraform will auto-generate a unique name.`,
	},

	"name_prefix": {
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		ForceNew:    true,
		Description: `Creates a unique name for the node pool beginning with the specified prefix. Conflicts with name.`,
	},

	"node_config": schemaNodeConfig(),

	"node_count": {
		Type:         schema.TypeInt,
		Optional:     true,
		Computed:     true,
		ValidateFunc: validation.IntAtLeast(0),
		Description:  `The number of nodes per instance group. This field can be used to update the number of nodes per instance group but should not be used alongside autoscaling.`,
	},

	"version": {
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: `The Kubernetes version for the nodes in this pool. Note that if this field and auto_upgrade are both specified, they will fight each other for what the node version should be, so setting both is highly discouraged. While a fuzzy version can be specified, it's recommended that you specify explicit versions as Terraform will see spurious diffs when fuzzy versions are used. See the google_container_engine_versions data source's version_prefix field to approximate fuzzy versions in a Terraform-compatible way.`,
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
	cluster := d.Get("cluster").(string)

	if fieldValues := clusterIdRegex.FindStringSubmatch(cluster); fieldValues != nil {
		log.Printf("[DEBUG] matching parent cluster %s to regex %s", cluster, clusterIdRegex.String())
		return &NodePoolInformation{
			project:  fieldValues[1],
			location: fieldValues[2],
			cluster:  fieldValues[3],
		}, nil
	}
	log.Printf("[DEBUG] parent cluster %s does not match regex %s", cluster, clusterIdRegex.String())

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
		cluster:  cluster,
	}, nil
}

func resourceContainerNodePoolCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

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

	req := &container.CreateNodePoolRequest{
		NodePool: nodePool,
	}

	timeout := d.Timeout(schema.TimeoutCreate)
	startTime := time.Now()

	// we attempt to prefetch the node pool to make sure it doesn't exist before creation
	var id = fmt.Sprintf("projects/%s/locations/%s/clusters/%s/nodePools/%s", nodePoolInfo.project, nodePoolInfo.location, nodePoolInfo.cluster, nodePool.Name)
	name := getNodePoolName(id)
	clusterNodePoolsGetCall := config.NewContainerClient(userAgent).Projects.Locations.Clusters.NodePools.Get(nodePoolInfo.fullyQualifiedName(name))
	if config.UserProjectOverride {
		clusterNodePoolsGetCall.Header().Add("X-Goog-User-Project", nodePoolInfo.project)
	}
	_, err = clusterNodePoolsGetCall.Do()
	if err != nil && isGoogleApiErrorWithCode(err, 404) {
		// Set the ID before we attempt to create if the resource doesn't exist. That
		// way, if we receive an error but the resource is created anyway, it will be
		// refreshed on the next call to apply.
		d.SetId(id)
	} else if err == nil {
		return fmt.Errorf("resource - %s - already exists", id)
	}

	var operation *container.Operation
	err = resource.Retry(timeout, func() *resource.RetryError {
		clusterNodePoolsCreateCall := config.NewContainerClient(userAgent).Projects.Locations.Clusters.NodePools.Create(nodePoolInfo.parent(), req)
		if config.UserProjectOverride {
			clusterNodePoolsCreateCall.Header().Add("X-Goog-User-Project", nodePoolInfo.project)
		}
		operation, err = clusterNodePoolsCreateCall.Do()

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

	waitErr := containerOperationWait(config,
		operation, nodePoolInfo.project,
		nodePoolInfo.location, "creating GKE NodePool", userAgent, timeout)

	if waitErr != nil {
		// Check if the create operation failed because Terraform was prematurely terminated. If it was we can persist the
		// operation id to state so that a subsequent refresh of this resource will wait until the operation has terminated
		// before attempting to Read the state of the cluster. This allows a graceful resumption of a Create that was killed
		// by the upstream Terraform process exiting early such as a sigterm.
		select {
		case <-config.context.Done():
			log.Printf("[DEBUG] Persisting %s so this operation can be resumed \n", operation.Name)
			if err := d.Set("operation", operation.Name); err != nil {
				return fmt.Errorf("Error setting operation: %s", err)
			}
			return nil
		default:
			// leaving default case to ensure this is non blocking
		}
		// Check if resource was created but apply timed out.
		// Common cause for that is GCE_STOCKOUT which will wait for resources and return error after timeout,
		// but in fact nodepool will be created so we have to capture that in state.
		_, err = clusterNodePoolsGetCall.Do()
		if err != nil {
			d.SetId("")
			return waitErr
		}
	}

	log.Printf("[INFO] GKE NodePool %s has been created", nodePool.Name)

	if err = resourceContainerNodePoolRead(d, meta); err != nil {
		return err
	}

	//Check cluster is in running state
	_, err = containerClusterAwaitRestingState(config, nodePoolInfo.project, nodePoolInfo.location, nodePoolInfo.cluster, userAgent, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}

	state, err := containerNodePoolAwaitRestingState(config, d.Id(), nodePoolInfo.project, userAgent, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}

	if containerNodePoolRestingStates[state] == ErrorState {
		return fmt.Errorf("NodePool %s was created in the error state %q", nodePool.Name, state)
	}

	return nil
}

func resourceContainerNodePoolRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	nodePoolInfo, err := extractNodePoolInformation(d, config)
	if err != nil {
		return err
	}

	operation := d.Get("operation").(string)
	if operation != "" {
		log.Printf("[DEBUG] in progress operation detected at %v, attempting to resume", operation)
		op := &container.Operation{
			Name: operation,
		}
		if err := d.Set("operation", ""); err != nil {
			return fmt.Errorf("Error setting operation: %s", err)
		}
		waitErr := containerOperationWait(config, op, nodePoolInfo.project, nodePoolInfo.location, "resuming GKE node pool", userAgent, d.Timeout(schema.TimeoutRead))
		if waitErr != nil {
			return waitErr
		}
	}

	name := getNodePoolName(d.Id())

	clusterNodePoolsGetCall := config.NewContainerClient(userAgent).Projects.Locations.Clusters.NodePools.Get(nodePoolInfo.fullyQualifiedName(name))
	if config.UserProjectOverride {
		clusterNodePoolsGetCall.Header().Add("X-Goog-User-Project", nodePoolInfo.project)
	}
	nodePool, err := clusterNodePoolsGetCall.Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("NodePool %q from cluster %q", name, nodePoolInfo.cluster))
	}

	npMap, err := flattenNodePool(d, config, nodePool, "")
	if err != nil {
		return err
	}

	for k, v := range npMap {
		if err := d.Set(k, v); err != nil {
			return fmt.Errorf("Error setting %s: %s", k, err)
		}
	}

	if err := d.Set("location", nodePoolInfo.location); err != nil {
		return fmt.Errorf("Error setting location: %s", err)
	}
	if err := d.Set("project", nodePoolInfo.project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}

	return nil
}

func resourceContainerNodePoolUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	nodePoolInfo, err := extractNodePoolInformation(d, config)
	if err != nil {
		return err
	}
	name := getNodePoolName(d.Id())

	//Check cluster is in running state
	_, err = containerClusterAwaitRestingState(config, nodePoolInfo.project, nodePoolInfo.location, nodePoolInfo.cluster, userAgent, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}

	_, err = containerNodePoolAwaitRestingState(config, nodePoolInfo.fullyQualifiedName(name), nodePoolInfo.project, userAgent, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}

	d.Partial(true)
	if err := nodePoolUpdate(d, meta, nodePoolInfo, "", d.Timeout(schema.TimeoutUpdate)); err != nil {
		return err
	}
	d.Partial(false)

	//Check cluster is in running state
	_, err = containerClusterAwaitRestingState(config, nodePoolInfo.project, nodePoolInfo.location, nodePoolInfo.cluster, userAgent, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}
	_, err = containerNodePoolAwaitRestingState(config, nodePoolInfo.fullyQualifiedName(name), nodePoolInfo.project, userAgent, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}

	return resourceContainerNodePoolRead(d, meta)
}

func resourceContainerNodePoolDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	nodePoolInfo, err := extractNodePoolInformation(d, config)
	if err != nil {
		return err
	}

	name := getNodePoolName(d.Id())

	//Check cluster is in running state
	_, err = containerClusterAwaitRestingState(config, nodePoolInfo.project, nodePoolInfo.location, nodePoolInfo.cluster, userAgent, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		if isGoogleApiErrorWithCode(err, 404) {
			log.Printf("[INFO] GKE cluster %s doesn't exist, skipping node pool %s deletion", nodePoolInfo.cluster, d.Id())
			return nil
		}
		return err
	}

	_, err = containerNodePoolAwaitRestingState(config, nodePoolInfo.fullyQualifiedName(name), nodePoolInfo.project, userAgent, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		// If the node pool doesn't get created and then we try to delete it, we get an error,
		// but I don't think we need an error during delete if it doesn't exist
		if isGoogleApiErrorWithCode(err, 404) {
			log.Printf("node pool %q not found, doesn't need to be cleaned up", name)
			return nil
		} else {
			return err
		}
	}

	mutexKV.Lock(nodePoolInfo.lockKey())
	defer mutexKV.Unlock(nodePoolInfo.lockKey())

	timeout := d.Timeout(schema.TimeoutDelete)
	startTime := time.Now()

	var operation *container.Operation
	err = resource.Retry(timeout, func() *resource.RetryError {
		clusterNodePoolsDeleteCall := config.NewContainerClient(userAgent).Projects.Locations.Clusters.NodePools.Delete(nodePoolInfo.fullyQualifiedName(name))
		if config.UserProjectOverride {
			clusterNodePoolsDeleteCall.Header().Add("X-Goog-User-Project", nodePoolInfo.project)
		}
		operation, err = clusterNodePoolsDeleteCall.Do()

		if err != nil {
			if isFailedPreconditionError(err) {
				// We get failed precondition errors if the cluster is updating
				// while we try to delete the node pool.
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("Error deleting NodePool: %s", err)
	}

	timeout -= time.Since(startTime)

	// Wait until it's deleted
	waitErr := containerOperationWait(config, operation, nodePoolInfo.project, nodePoolInfo.location, "deleting GKE NodePool", userAgent, timeout)
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

	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return false, err
	}

	name := getNodePoolName(d.Id())
	clusterNodePoolsGetCall := config.NewContainerClient(userAgent).Projects.Locations.Clusters.NodePools.Get(nodePoolInfo.fullyQualifiedName(name))
	if config.UserProjectOverride {
		clusterNodePoolsGetCall.Header().Add("X-Goog-User-Project", nodePoolInfo.project)
	}
	_, err = clusterNodePoolsGetCall.Do()

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
	config := meta.(*Config)

	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return nil, err
	}

	if err := parseImportId([]string{"projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/clusters/(?P<cluster>[^/]+)/nodePools/(?P<name>[^/]+)", "(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<cluster>[^/]+)/(?P<name>[^/]+)", "(?P<location>[^/]+)/(?P<cluster>[^/]+)/(?P<name>[^/]+)"}, d, config); err != nil {
		return nil, err
	}

	id, err := replaceVars(d, config, "projects/{{project}}/locations/{{location}}/clusters/{{cluster}}/nodePools/{{name}}")
	if err != nil {
		return nil, err
	}

	d.SetId(id)

	project, err := getProject(d, config)
	if err != nil {
		return nil, err
	}

	nodePoolInfo, err := extractNodePoolInformation(d, config)
	if err != nil {
		return nil, err
	}

	//Check cluster is in running state
	_, err = containerClusterAwaitRestingState(config, nodePoolInfo.project, nodePoolInfo.location, nodePoolInfo.cluster, userAgent, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return nil, err
	}

	if _, err := containerNodePoolAwaitRestingState(config, d.Id(), project, userAgent, d.Timeout(schema.TimeoutCreate)); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func expandNodePool(d *schema.ResourceData, prefix string) (*container.NodePool, error) {
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

	var locations []string
	if v, ok := d.GetOk("node_locations"); ok && v.(*schema.Set).Len() > 0 {
		locations = convertStringSet(v.(*schema.Set))
	}

	np := &container.NodePool{
		Name:             name,
		InitialNodeCount: int64(nodeCount),
		Config:           expandNodeConfig(d.Get(prefix + "node_config")),
		Locations:        locations,
		Version:          d.Get(prefix + "version").(string),
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

	if v, ok := d.GetOk(prefix + "max_pods_per_node"); ok {
		np.MaxPodsConstraint = &container.MaxPodsConstraint{
			MaxPodsPerNode: int64(v.(int)),
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

	if v, ok := d.GetOk(prefix + "upgrade_settings"); ok {
		upgradeSettingsConfig := v.([]interface{})[0].(map[string]interface{})
		np.UpgradeSettings = &container.UpgradeSettings{}

		if v, ok := upgradeSettingsConfig["max_surge"]; ok {
			np.UpgradeSettings.MaxSurge = int64(v.(int))
		}

		if v, ok := upgradeSettingsConfig["max_unavailable"]; ok {
			np.UpgradeSettings.MaxUnavailable = int64(v.(int))
		}
	}

	return np, nil
}

func flattenNodePool(d *schema.ResourceData, config *Config, np *container.NodePool, prefix string) (map[string]interface{}, error) {
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return nil, err
	}

	// Node pools don't expose the current node count in their API, so read the
	// instance groups instead. They should all have the same size, but in case a resize
	// failed or something else strange happened, we'll just use the average size.
	size := 0
	igmUrls := []string{}
	managedIgmUrls := []string{}
	for _, url := range np.InstanceGroupUrls {
		// retrieve instance group manager (InstanceGroupUrls are actually URLs for InstanceGroupManagers)
		matches := instanceGroupManagerURL.FindStringSubmatch(url)
		if len(matches) < 4 {
			return nil, fmt.Errorf("Error reading instance group manage URL '%q'", url)
		}
		igm, err := config.NewComputeClient(userAgent).InstanceGroupManagers.Get(matches[1], matches[2], matches[3]).Do()
		if isGoogleApiErrorWithCode(err, 404) {
			// The IGM URL in is stale; don't include it
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("Error reading instance group manager returned as an instance group URL: %q", err)
		}
		size += int(igm.TargetSize)
		igmUrls = append(igmUrls, url)
		managedIgmUrls = append(managedIgmUrls, igm.InstanceGroup)
	}
	nodeCount := 0
	if len(igmUrls) > 0 {
		nodeCount = size / len(igmUrls)
	}
	nodePool := map[string]interface{}{
		"name":                        np.Name,
		"name_prefix":                 d.Get(prefix + "name_prefix"),
		"initial_node_count":          np.InitialNodeCount,
		"node_locations":              schema.NewSet(schema.HashString, convertStringArrToInterface(np.Locations)),
		"node_count":                  nodeCount,
		"node_config":                 flattenNodeConfig(np.Config),
		"instance_group_urls":         igmUrls,
		"managed_instance_group_urls": managedIgmUrls,
		"version":                     np.Version,
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

	if np.UpgradeSettings != nil {
		nodePool["upgrade_settings"] = []map[string]interface{}{
			{
				"max_surge":       np.UpgradeSettings.MaxSurge,
				"max_unavailable": np.UpgradeSettings.MaxUnavailable,
			},
		}
	} else {
		delete(nodePool, "upgrade_settings")
	}

	return nodePool, nil
}

func nodePoolUpdate(d *schema.ResourceData, meta interface{}, nodePoolInfo *NodePoolInformation, prefix string, timeout time.Duration) error {
	config := meta.(*Config)
	name := d.Get(prefix + "name").(string)

	lockKey := nodePoolInfo.lockKey()

	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	if d.HasChange(prefix + "autoscaling") {
		update := &container.ClusterUpdate{
			DesiredNodePoolId: name,
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

		updateF := func() error {
			clusterUpdateCall := config.NewContainerClient(userAgent).Projects.Locations.Clusters.Update(nodePoolInfo.parent(), req)
			if config.UserProjectOverride {
				clusterUpdateCall.Header().Add("X-Goog-User-Project", nodePoolInfo.project)
			}
			op, err := clusterUpdateCall.Do()
			if err != nil {
				return err
			}

			// Wait until it's updated
			return containerOperationWait(config, op,
				nodePoolInfo.project,
				nodePoolInfo.location, "updating GKE node pool", userAgent,
				timeout)
		}

		// Call update serially.
		if err := lockedCall(lockKey, updateF); err != nil {
			return err
		}

		log.Printf("[INFO] Updated autoscaling in Node Pool %s", d.Id())
	}

	if d.HasChange(prefix + "node_config") {
		if d.HasChange(prefix + "node_config.0.image_type") {
			req := &container.UpdateClusterRequest{
				Update: &container.ClusterUpdate{
					DesiredNodePoolId: name,
					DesiredImageType:  d.Get(prefix + "node_config.0.image_type").(string),
				},
			}

			updateF := func() error {
				clusterUpdateCall := config.NewContainerClient(userAgent).Projects.Locations.Clusters.Update(nodePoolInfo.parent(), req)
				if config.UserProjectOverride {
					clusterUpdateCall.Header().Add("X-Goog-User-Project", nodePoolInfo.project)
				}
				op, err := clusterUpdateCall.Do()
				if err != nil {
					return err
				}

				// Wait until it's updated
				return containerOperationWait(config, op,
					nodePoolInfo.project,
					nodePoolInfo.location, "updating GKE node pool", userAgent,
					timeout)
			}

			// Call update serially.
			if err := lockedCall(lockKey, updateF); err != nil {
				return err
			}

			log.Printf("[INFO] Updated image type in Node Pool %s", d.Id())
		}

		if d.HasChange(prefix + "node_config.0.workload_metadata_config") {
			req := &container.UpdateNodePoolRequest{
				NodePoolId: name,
				WorkloadMetadataConfig: expandWorkloadMetadataConfig(
					d.Get(prefix + "node_config.0.workload_metadata_config")),
			}
			if req.WorkloadMetadataConfig == nil {
				req.ForceSendFields = []string{"WorkloadMetadataConfig"}
			}
			updateF := func() error {
				clusterNodePoolsUpdateCall := config.NewContainerClient(userAgent).Projects.Locations.Clusters.NodePools.Update(nodePoolInfo.fullyQualifiedName(name), req)
				if config.UserProjectOverride {
					clusterNodePoolsUpdateCall.Header().Add("X-Goog-User-Project", nodePoolInfo.project)
				}
				op, err := clusterNodePoolsUpdateCall.Do()

				if err != nil {
					return err
				}

				// Wait until it's updated
				return containerOperationWait(config, op,
					nodePoolInfo.project,
					nodePoolInfo.location,
					"updating GKE node pool workload_metadata_config", userAgent,
					timeout)
			}

			// Call update serially.
			if err := lockedCall(lockKey, updateF); err != nil {
				return err
			}

			log.Printf("[INFO] Updated workload_metadata_config for node pool %s", name)
		}

	}

	if d.HasChange(prefix + "node_count") {
		newSize := int64(d.Get(prefix + "node_count").(int))
		req := &container.SetNodePoolSizeRequest{
			NodeCount: newSize,
		}
		updateF := func() error {
			clusterNodePoolsSetSizeCall := config.NewContainerClient(userAgent).Projects.Locations.Clusters.NodePools.SetSize(nodePoolInfo.fullyQualifiedName(name), req)
			if config.UserProjectOverride {
				clusterNodePoolsSetSizeCall.Header().Add("X-Goog-User-Project", nodePoolInfo.project)
			}
			op, err := clusterNodePoolsSetSizeCall.Do()

			if err != nil {
				return err
			}

			// Wait until it's updated
			return containerOperationWait(config, op,
				nodePoolInfo.project,
				nodePoolInfo.location, "updating GKE node pool size", userAgent,
				timeout)
		}

		// Call update serially.
		if err := lockedCall(lockKey, updateF); err != nil {
			return err
		}

		log.Printf("[INFO] GKE node pool %s size has been updated to %d", name, newSize)
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

		updateF := func() error {
			clusterNodePoolsSetManagementCall := config.NewContainerClient(userAgent).Projects.Locations.Clusters.NodePools.SetManagement(nodePoolInfo.fullyQualifiedName(name), req)
			if config.UserProjectOverride {
				clusterNodePoolsSetManagementCall.Header().Add("X-Goog-User-Project", nodePoolInfo.project)
			}
			op, err := clusterNodePoolsSetManagementCall.Do()

			if err != nil {
				return err
			}

			// Wait until it's updated
			return containerOperationWait(config, op,
				nodePoolInfo.project,
				nodePoolInfo.location, "updating GKE node pool management", userAgent, timeout)
		}

		// Call update serially.
		if err := lockedCall(lockKey, updateF); err != nil {
			return err
		}

		log.Printf("[INFO] Updated management in Node Pool %s", name)
	}

	if d.HasChange(prefix + "version") {
		req := &container.UpdateNodePoolRequest{
			NodePoolId:  name,
			NodeVersion: d.Get(prefix + "version").(string),
		}
		updateF := func() error {
			clusterNodePoolsUpdateCall := config.NewContainerClient(userAgent).Projects.Locations.Clusters.NodePools.Update(nodePoolInfo.fullyQualifiedName(name), req)
			if config.UserProjectOverride {
				clusterNodePoolsUpdateCall.Header().Add("X-Goog-User-Project", nodePoolInfo.project)
			}
			op, err := clusterNodePoolsUpdateCall.Do()

			if err != nil {
				return err
			}

			// Wait until it's updated
			return containerOperationWait(config, op,
				nodePoolInfo.project,
				nodePoolInfo.location, "updating GKE node pool version", userAgent, timeout)
		}

		// Call update serially.
		if err := lockedCall(lockKey, updateF); err != nil {
			return err
		}

		log.Printf("[INFO] Updated version in Node Pool %s", name)
	}

	if d.HasChange(prefix + "node_locations") {
		req := &container.UpdateNodePoolRequest{
			Locations: convertStringSet(d.Get(prefix + "node_locations").(*schema.Set)),
		}
		updateF := func() error {
			clusterNodePoolsUpdateCall := config.NewContainerClient(userAgent).Projects.Locations.Clusters.NodePools.Update(nodePoolInfo.fullyQualifiedName(name), req)
			if config.UserProjectOverride {
				clusterNodePoolsUpdateCall.Header().Add("X-Goog-User-Project", nodePoolInfo.project)
			}
			op, err := clusterNodePoolsUpdateCall.Do()

			if err != nil {
				return err
			}

			// Wait until it's updated
			return containerOperationWait(config, op, nodePoolInfo.project, nodePoolInfo.location, "updating GKE node pool node locations", userAgent, timeout)
		}

		// Call update serially.
		if err := lockedCall(lockKey, updateF); err != nil {
			return err
		}

		log.Printf("[INFO] Updated node locations in Node Pool %s", name)
	}

	if d.HasChange(prefix + "upgrade_settings") {
		upgradeSettings := &container.UpgradeSettings{}
		if v, ok := d.GetOk(prefix + "upgrade_settings"); ok {
			upgradeSettingsConfig := v.([]interface{})[0].(map[string]interface{})
			upgradeSettings.MaxSurge = int64(upgradeSettingsConfig["max_surge"].(int))
			upgradeSettings.MaxUnavailable = int64(upgradeSettingsConfig["max_unavailable"].(int))
		}
		req := &container.UpdateNodePoolRequest{
			UpgradeSettings: upgradeSettings,
		}
		updateF := func() error {
			clusterNodePoolsUpdateCall := config.NewContainerClient(userAgent).Projects.Locations.Clusters.NodePools.Update(nodePoolInfo.fullyQualifiedName(name), req)
			if config.UserProjectOverride {
				clusterNodePoolsUpdateCall.Header().Add("X-Goog-User-Project", nodePoolInfo.project)
			}
			op, err := clusterNodePoolsUpdateCall.Do()

			if err != nil {
				return err
			}

			// Wait until it's updated
			return containerOperationWait(config, op, nodePoolInfo.project, nodePoolInfo.location, "updating GKE node pool upgrade settings", userAgent, timeout)
		}

		// Call update serially.
		if err := lockedCall(lockKey, updateF); err != nil {
			return err
		}

		log.Printf("[INFO] Updated upgrade settings in Node Pool %s", name)
	}

	return nil
}

func getNodePoolName(id string) string {
	// name can be specified with name, name_prefix, or neither, so read it from the id.
	splits := strings.Split(id, "/")
	return splits[len(splits)-1]
}

var containerNodePoolRestingStates = RestingStates{
	"RUNNING":            ReadyState,
	"RUNNING_WITH_ERROR": ErrorState,
	"ERROR":              ErrorState,
}

// takes in a config object, full node pool name, project name and the current CRUD action timeout
// returns a state with no error if the state is a resting state, and the last state with an error otherwise
func containerNodePoolAwaitRestingState(config *Config, name, project, userAgent string, timeout time.Duration) (state string, err error) {
	err = resource.Retry(timeout, func() *resource.RetryError {
		clusterNodePoolsGetCall := config.NewContainerClient(userAgent).Projects.Locations.Clusters.NodePools.Get(name)
		if config.UserProjectOverride {
			clusterNodePoolsGetCall.Header().Add("X-Goog-User-Project", project)
		}
		nodePool, gErr := clusterNodePoolsGetCall.Do()
		if gErr != nil {
			return resource.NonRetryableError(gErr)
		}

		state = nodePool.Status
		switch stateType := containerNodePoolRestingStates[state]; stateType {
		case ReadyState:
			log.Printf("[DEBUG] NodePool %q has status %q with message %q.", name, state, nodePool.StatusMessage)
			return nil
		case ErrorState:
			log.Printf("[DEBUG] NodePool %q has error state %q with message %q.", name, state, nodePool.StatusMessage)
			return nil
		default:
			return resource.RetryableError(fmt.Errorf("NodePool %q has state %q with message %q", name, state, nodePool.StatusMessage))
		}
	})

	return state, err
}
