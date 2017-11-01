package google

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/compute/v1"
)

func resourceComputeNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeNetworkCreate,
		Read:   resourceComputeNetworkRead,
		Delete: resourceComputeNetworkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"auto_create_subnetworks": &schema.Schema{
				Type:          schema.TypeBool,
				Optional:      true,
				ForceNew:      true,
				Default:       true,
				ConflictsWith: []string{"ipv4_range"},
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"gateway_ipv4": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"ipv4_range": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Removed:  "Please use google_compute_subnetwork resources instead.",
			},

			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"self_link": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceComputeNetworkCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	// Build the network parameter
	network := &compute.Network{
		Name: d.Get("name").(string),
		AutoCreateSubnetworks: d.Get("auto_create_subnetworks").(bool),
		Description:           d.Get("description").(string),
	}

	// make sure AutoCreateSubnetworks field is included in request
	network.ForceSendFields = []string{"AutoCreateSubnetworks"}
	log.Printf("[DEBUG] Network insert request: %#v", network)
	op, err := config.clientCompute.Networks.Insert(
		project, network).Do()
	if err != nil {
		return fmt.Errorf("Error creating network: %s", err)
	}

	// It probably maybe worked, so store the ID now
	d.SetId(network.Name)

	err = computeOperationWait(config.clientCompute, op, project, "Creating Network")
	if err != nil {
		return err
	}

	return resourceComputeNetworkRead(d, meta)
}

func resourceComputeNetworkRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	network, err := config.clientCompute.Networks.Get(
		project, d.Id()).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Network %q", d.Get("name").(string)))
	}

	d.Set("gateway_ipv4", network.GatewayIPv4)
	d.Set("self_link", network.SelfLink)
	d.Set("name", network.Name)
	d.Set("auto_create_subnetworks", network.AutoCreateSubnetworks)

	return nil
}

func resourceComputeNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	// Delete the network
	op, err := config.clientCompute.Networks.Delete(
		project, d.Id()).Do()
	if err != nil {
		return fmt.Errorf("Error deleting network: %s", err)
	}

	err = computeOperationWaitTime(config.clientCompute, op, project, "Deleting Network", 10)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
