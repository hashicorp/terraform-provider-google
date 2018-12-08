package openstack

import (
	"fmt"
	"log"
	"time"

	"github.com/gophercloud/gophercloud/openstack/containerinfra/v1/clusters"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceContainerInfraCluster() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceContainerInfraClusterRead,
		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"project_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"api_address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"coe_version": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_template_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"container_version": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_timeout": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"discovery_url": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"docker_volume_size": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"flavor": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"master_flavor": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"keypair": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"labels": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
			"master_count": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"node_count": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"master_addresses": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_addresses": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"stack_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceContainerInfraClusterRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	containerInfraClient, err := config.containerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack container infra client: %s", err)
	}

	name := d.Get("name").(string)
	c, err := clusters.Get(containerInfraClient, name).Extract()
	if err != nil {
		return fmt.Errorf("Error getting OpenStack container infra cluster: %s", err)
	}

	d.SetId(c.UUID)

	d.Set("project_id", c.ProjectID)
	d.Set("user_id", c.UserID)
	d.Set("api_address", c.APIAddress)
	d.Set("coe_version", c.COEVersion)
	d.Set("cluster_template_id", c.ClusterTemplateID)
	d.Set("container_version", c.ContainerVersion)
	d.Set("create_timeout", c.CreateTimeout)
	d.Set("discovery_url", c.DiscoveryURL)
	d.Set("docker_volume_size", c.DockerVolumeSize)
	d.Set("flavor", c.FlavorID)
	d.Set("master_flavor", c.MasterFlavorID)
	d.Set("keypair", c.KeyPair)
	d.Set("master_count", c.MasterCount)
	d.Set("node_count", c.NodeCount)
	d.Set("master_addresses", c.MasterAddresses)
	d.Set("node_addresses", c.NodeAddresses)
	d.Set("stack_id", c.StackID)

	if err := d.Set("labels", c.Labels); err != nil {
		log.Printf("[DEBUG] Unable to set labels for cluster %s: %s", c.UUID, err)
	}
	if err := d.Set("created_at", c.CreatedAt.Format(time.RFC3339)); err != nil {
		log.Printf("[DEBUG] Unable to set created_at for cluster %s: %s", c.UUID, err)
	}
	if err := d.Set("updated_at", c.UpdatedAt.Format(time.RFC3339)); err != nil {
		log.Printf("[DEBUG] Unable to set updated_at for cluster %s: %s", c.UUID, err)
	}

	d.Set("region", GetRegion(d, config))

	return nil
}
