package google

import (
	"fmt"
	"log"
	"strings"

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

		SchemaVersion: 1,
		MigrateState:  resourceContainerNodePoolMigrateState,

		Importer: &schema.ResourceImporter{
			State: resourceContainerNodePoolStateImporter,
		},

		Schema: map[string]*schema.Schema{
			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name_prefix"},
				ForceNew:      true,
			},

			"name_prefix": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"zone": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"cluster": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"initial_node_count": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"node_config": schemaNodeConfig,

			"autoscaling": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"min_node_count": &schema.Schema{
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntAtLeast(1),
						},

						"max_node_count": &schema.Schema{
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntAtLeast(1),
						},
					},
				},
			},
		},
	}
}

func resourceContainerNodePoolCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	zone := d.Get("zone").(string)
	cluster := d.Get("cluster").(string)
	nodeCount := d.Get("initial_node_count").(int)

	var name string
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else if v, ok := d.GetOk("name_prefix"); ok {
		name = resource.PrefixedUniqueId(v.(string))
	} else {
		name = resource.UniqueId()
	}

	nodePool := &container.NodePool{
		Name:             name,
		InitialNodeCount: int64(nodeCount),
	}

	if v, ok := d.GetOk("node_config"); ok {
		nodePool.Config = expandNodeConfig(v)
	}

	if v, ok := d.GetOk("autoscaling"); ok {
		autoscaling := v.([]interface{})[0].(map[string]interface{})
		nodePool.Autoscaling = &container.NodePoolAutoscaling{
			Enabled:      true,
			MinNodeCount: int64(autoscaling["min_node_count"].(int)),
			MaxNodeCount: int64(autoscaling["max_node_count"].(int)),
		}
	}

	req := &container.CreateNodePoolRequest{
		NodePool: nodePool,
	}

	op, err := config.clientContainer.Projects.Zones.Clusters.NodePools.Create(project, zone, cluster, req).Do()

	if err != nil {
		return fmt.Errorf("Error creating NodePool: %s", err)
	}

	waitErr := containerOperationWait(config, op, project, zone, "creating GKE NodePool", 10, 3)
	if waitErr != nil {
		// The resource didn't actually create
		d.SetId("")
		return waitErr
	}

	log.Printf("[INFO] GKE NodePool %s has been created", name)

	d.SetId(fmt.Sprintf("%s/%s/%s", zone, cluster, name))

	return resourceContainerNodePoolRead(d, meta)
}

func resourceContainerNodePoolRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	zone := d.Get("zone").(string)
	name := d.Get("name").(string)
	cluster := d.Get("cluster").(string)

	nodePool, err := config.clientContainer.Projects.Zones.Clusters.NodePools.Get(
		project, zone, cluster, name).Do()
	if err != nil {
		return fmt.Errorf("Error reading NodePool: %s", err)
	}

	d.Set("name", nodePool.Name)
	d.Set("initial_node_count", nodePool.InitialNodeCount)
	d.Set("node_config", flattenClusterNodeConfig(nodePool.Config))

	autoscaling := []map[string]interface{}{}
	if nodePool.Autoscaling != nil && nodePool.Autoscaling.Enabled {
		autoscaling = []map[string]interface{}{
			map[string]interface{}{
				"min_node_count": nodePool.Autoscaling.MinNodeCount,
				"max_node_count": nodePool.Autoscaling.MaxNodeCount,
			},
		}
	}
	d.Set("autoscaling", autoscaling)

	return nil
}

func resourceContainerNodePoolUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	zone := d.Get("zone").(string)
	name := d.Get("name").(string)
	cluster := d.Get("cluster").(string)

	if d.HasChange("autoscaling") {
		update := &container.ClusterUpdate{
			DesiredNodePoolId: name,
		}
		if v, ok := d.GetOk("autoscaling"); ok {
			autoscaling := v.([]interface{})[0].(map[string]interface{})
			update.DesiredNodePoolAutoscaling = &container.NodePoolAutoscaling{
				Enabled:      true,
				MinNodeCount: int64(autoscaling["min_node_count"].(int)),
				MaxNodeCount: int64(autoscaling["max_node_count"].(int)),
			}
		} else {
			update.DesiredNodePoolAutoscaling = &container.NodePoolAutoscaling{
				Enabled: false,
			}
		}

		req := &container.UpdateClusterRequest{
			Update: update,
		}
		op, err := config.clientContainer.Projects.Zones.Clusters.Update(
			project, zone, cluster, req).Do()
		if err != nil {
			return err
		}

		// Wait until it's updated
		waitErr := containerOperationWait(config, op, project, zone, "updating GKE node pool", 10, 2)
		if waitErr != nil {
			return waitErr
		}

		log.Printf("[INFO] Updated autoscaling in Node Pool %s", d.Id())
	}

	return resourceContainerNodePoolRead(d, meta)
}

func resourceContainerNodePoolDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	zone := d.Get("zone").(string)
	name := d.Get("name").(string)
	cluster := d.Get("cluster").(string)

	op, err := config.clientContainer.Projects.Zones.Clusters.NodePools.Delete(
		project, zone, cluster, name).Do()
	if err != nil {
		return fmt.Errorf("Error deleting NodePool: %s", err)
	}

	// Wait until it's deleted
	waitErr := containerOperationWait(config, op, project, zone, "deleting GKE NodePool", 10, 2)
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

	zone := d.Get("zone").(string)
	name := d.Get("name").(string)
	cluster := d.Get("cluster").(string)

	_, err = config.clientContainer.Projects.Zones.Clusters.NodePools.Get(
		project, zone, cluster, name).Do()
	if err != nil {
		if err = handleNotFoundError(err, d, fmt.Sprintf("Container NodePool %s", d.Get("name").(string))); err == nil {
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
