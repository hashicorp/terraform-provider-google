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
		Update: resourceComputeNetworkUpdate,
		Delete: resourceComputeNetworkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"auto_create_subnetworks": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"routing_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"gateway_ipv4": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ipv4_range": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				// This needs to remain deprecated until the API is retired
				Deprecated: "Please use google_compute_subnetwork resources instead.",
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"self_link": {
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

	//
	// Possible modes:
	// - 1 Legacy mode - Create a network in the legacy mode. ipv4_range is set. auto_create_subnetworks must not be
	//     set (enforced by ConflictsWith schema attribute)
	// - 2 Distributed Mode - Create a new generation network that supports subnetworks:
	//   - 2.a - Auto subnet mode - auto_create_subnetworks = true, Google will generate 1 subnetwork per region
	//   - 2.b - Custom subnet mode - auto_create_subnetworks = false & ipv4_range not set,
	//
	autoCreateSubnetworks := d.Get("auto_create_subnetworks").(bool)
	if autoCreateSubnetworks && d.Get("ipv4_range").(string) != "" {
		return fmt.Errorf("ipv4_range can't be set if auto_create_subnetworks is true.")
	}

	// Build the network parameter
	network := &compute.Network{
		Name:                  d.Get("name").(string),
		AutoCreateSubnetworks: autoCreateSubnetworks,
		Description:           d.Get("description").(string),
	}

	if v, ok := d.GetOk("routing_mode"); ok {
		routingConfig := &compute.NetworkRoutingConfig{
			RoutingMode: v.(string),
		}
		network.RoutingConfig = routingConfig
	}

	if v, ok := d.GetOk("ipv4_range"); ok {
		log.Printf("[DEBUG] Setting IPv4Range (%#v) for legacy network mode", v.(string))
		network.IPv4Range = v.(string)
	} else {
		// custom subnet mode, so make sure AutoCreateSubnetworks field is included in request otherwise
		// google will create a network in legacy mode.
		network.ForceSendFields = []string{"AutoCreateSubnetworks"}
	}
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

	routingConfig := network.RoutingConfig

	d.Set("routing_mode", routingConfig.RoutingMode)
	d.Set("gateway_ipv4", network.GatewayIPv4)
	d.Set("ipv4_range", network.IPv4Range)
	d.Set("self_link", network.SelfLink)
	d.Set("name", network.Name)
	d.Set("description", network.Description)
	d.Set("auto_create_subnetworks", network.AutoCreateSubnetworks)
	d.Set("project", project)

	return nil
}

func resourceComputeNetworkUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	op, err := config.clientCompute.Networks.Patch(project, d.Id(), &compute.Network{
		RoutingConfig: &compute.NetworkRoutingConfig{
			RoutingMode: d.Get("routing_mode").(string),
		},
	}).Do()

	if err != nil {
		return fmt.Errorf("Error updating network: %s", err)
	}

	err = computeSharedOperationWait(config.clientCompute, op, project, "UpdateNetwork")
	if err != nil {
		return err
	}

	return resourceComputeNetworkRead(d, meta)
}

func resourceComputeNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	return deleteComputeNetwork(project, d.Id(), config)
}

func deleteComputeNetwork(project, network string, config *Config) error {
	op, err := config.clientCompute.Networks.Delete(
		project, network).Do()
	if err != nil {
		return fmt.Errorf("Error deleting network: %s", err)
	}

	err = computeOperationWaitTime(config.clientCompute, op, project, "Deleting Network", 10)
	if err != nil {
		return err
	}
	return nil
}
