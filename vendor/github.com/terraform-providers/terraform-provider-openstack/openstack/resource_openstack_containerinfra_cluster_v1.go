package openstack

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/containerinfra/v1/clusters"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceContainerInfraClusterV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceContainerInfraClusterV1Create,
		Read:   resourceContainerInfraClusterV1Read,
		Update: resourceContainerInfraClusterV1Update,
		Delete: resourceContainerInfraClusterV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"project_id": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
			"user_id": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
			"created_at": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: false,
				Computed: true,
			},
			"updated_at": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: false,
				Computed: true,
			},
			"api_address": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: false,
				Computed: true,
			},
			"coe_version": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: false,
				Computed: true,
			},
			"cluster_template_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_MAGNUM_CLUSTER_TEMPLATE", nil),
			},
			"container_version": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: false,
				Computed: true,
			},
			"create_timeout": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"discovery_url": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"docker_volume_size": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"flavor": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_MAGNUM_FLAVOR", nil),
			},
			"master_flavor": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_MAGNUM_MASTER_FLAVOR", nil),
			},
			"keypair": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"labels": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"master_count": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"node_count": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},
			"master_addresses": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: false,
				Computed: true,
			},
			"node_addresses": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: false,
				Computed: true,
			},
			"stack_id": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: false,
				Computed: true,
			},
		},
	}
}

func resourceContainerInfraClusterV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	containerInfraClient, err := config.containerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack container infra client: %s", err)
	}

	// Get and check labels map.
	labels, err := containerInfraLabelsMapV1(d)
	if err != nil {
		return err
	}

	createOpts := clusters.CreateOpts{
		ClusterTemplateID: d.Get("cluster_template_id").(string),
		DiscoveryURL:      d.Get("discovery_url").(string),
		FlavorID:          d.Get("flavor").(string),
		Keypair:           d.Get("keypair").(string),
		Labels:            labels,
		MasterFlavorID:    d.Get("master_flavor").(string),
		Name:              d.Get("name").(string),
	}

	// Set int parameters that will be passed by reference.
	createTimeout := d.Get("create_timeout").(int)
	if createTimeout > 0 {
		createOpts.CreateTimeout = &createTimeout
	}
	dockerVolumeSize := d.Get("docker_volume_size").(int)
	if dockerVolumeSize > 0 {
		createOpts.DockerVolumeSize = &dockerVolumeSize
	}
	masterCount := d.Get("master_count").(int)
	if masterCount > 0 {
		createOpts.MasterCount = &masterCount
	}
	nodeCount := d.Get("node_count").(int)
	if nodeCount > 0 {
		createOpts.NodeCount = &nodeCount
	}

	s, err := clusters.Create(containerInfraClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating OpenStack container infra Cluster: %s", err)
	}

	// Store the Cluster ID.
	d.SetId(s)

	log.Printf("[DEBUG] Waiting for Cluster (%s) to become ready", s)
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"CREATE_IN_PROGRESS"},
		Target:       []string{"CREATE_COMPLETE"},
		Refresh:      ContainerInfraClusterV1StateRefreshFunc(containerInfraClient, s),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        1 * time.Minute,
		PollInterval: 20 * time.Second,
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for container infra Cluster (%s) to become ready: %s",
			s, err)
	}

	d.SetId(s)

	log.Printf("[DEBUG] Created Cluster %s", s)
	return resourceContainerInfraClusterV1Read(d, meta)
}

func resourceContainerInfraClusterV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	containerInfraClient, err := config.containerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack container infra client: %s", err)
	}

	s, err := clusters.Get(containerInfraClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "cluster")
	}

	log.Printf("[DEBUG] Retrieved Cluster %s: %#v", d.Id(), s)

	if err := d.Set("labels", s.Labels); err != nil {
		return fmt.Errorf("Unable to set labels: %s", err)
	}

	d.Set("name", s.Name)
	d.Set("api_address", s.APIAddress)
	d.Set("coe_version", s.COEVersion)
	d.Set("cluster_template_id", s.ClusterTemplateID)
	d.Set("container_version", s.ContainerVersion)
	d.Set("create_timeout", s.CreateTimeout)
	d.Set("discovery_url", s.DiscoveryURL)
	d.Set("docker_volume_size", s.DockerVolumeSize)
	d.Set("flavor", s.FlavorID)
	d.Set("master_flavor", s.MasterFlavorID)
	d.Set("keypair", s.KeyPair)
	d.Set("master_count", s.MasterCount)
	d.Set("node_count", s.NodeCount)
	d.Set("master_addresses", s.MasterAddresses)
	d.Set("node_addresses", s.NodeAddresses)
	d.Set("stack_id", s.StackID)

	if err := d.Set("created_at", s.CreatedAt.Format(time.RFC3339)); err != nil {
		log.Printf("[DEBUG] created_at: %s", err)
	}
	if err := d.Set("updated_at", s.UpdatedAt.Format(time.RFC3339)); err != nil {
		log.Printf("[DEBUG] updated_at: %s", err)
	}

	return nil
}

func resourceContainerInfraClusterV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	containerInfraClient, err := config.containerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack container infra client: %s", err)
	}

	updateOpts := []clusters.UpdateOptsBuilder{}

	if d.HasChange("node_count") {
		v := d.Get("node_count").(int)
		nodeCount := strconv.Itoa(v)
		updateOpts = append(updateOpts, clusters.UpdateOpts{
			Op:    clusters.ReplaceOp,
			Path:  strings.Join([]string{"/", "node_count"}, ""),
			Value: nodeCount,
		})
	}

	log.Printf("[DEBUG] Updating Cluster %s with options: %+v", d.Id(), updateOpts)

	_, err = clusters.Update(containerInfraClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating OpenStack container infra Cluster: %s", err)
	}

	log.Printf("[DEBUG] Waiting for Cluster (%s) to become updated", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"UPDATE_IN_PROGRESS"},
		Target:       []string{"UPDATE_COMPLETE"},
		Refresh:      ContainerInfraClusterV1StateRefreshFunc(containerInfraClient, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        1 * time.Minute,
		PollInterval: 20 * time.Second,
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for container infra Cluster (%s) to become updated: %s",
			d.Id(), err)
	}

	return resourceContainerInfraClusterV1Read(d, meta)
}

func resourceContainerInfraClusterV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	containerInfraClient, err := config.containerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack container infra client: %s", err)
	}

	if err := clusters.Delete(containerInfraClient, d.Id()).ExtractErr(); err != nil {
		return fmt.Errorf("Error deleting Cluster: %v", err)
	}

	log.Printf("[DEBUG] Waiting for Cluster (%s) to become deleted", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"DELETE_IN_PROGRESS"},
		Target:       []string{"DELETE_COMPLETE"},
		Refresh:      ContainerInfraClusterV1StateRefreshFunc(containerInfraClient, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        30 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for container infra Cluster (%s) to become deleted: %s",
			d.Id(), err)
	}

	return nil
}

// ContainerInfraClusterV1StateRefreshFunc returns a resource.StateRefreshFunc
// that is used to watch a container infra Cluster.
func ContainerInfraClusterV1StateRefreshFunc(client *gophercloud.ServiceClient, clusterID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		c, err := clusters.Get(client, clusterID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return c, "DELETE_COMPLETE", nil
			}
			return nil, "", err
		}

		errorStatuses := []string{
			"CREATE_FAILED",
			"UPDATE_FAILED",
			"DELETE_FAILED",
			"RESUME_FAILED",
			"ROLLBACK_FAILED",
		}
		for _, errorStatus := range errorStatuses {
			if c.Status == errorStatus {
				err = fmt.Errorf("There was an error creating the container infra cluster: %s", c.StatusReason)
				return c, c.Status, err
			}
		}

		return c, c.Status, nil
	}
}
