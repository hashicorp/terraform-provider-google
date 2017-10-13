package google

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/compute/v1"
)

func resourceComputeRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeRouteCreate,
		Read:   resourceComputeRouteRead,
		Delete: resourceComputeRouteDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"dest_range": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"network": &schema.Schema{
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
			},

			"priority": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"next_hop_gateway": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
			},

			"next_hop_instance": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"next_hop_instance_zone": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"next_hop_ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"next_hop_vpn_tunnel": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"tags": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"self_link": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"next_hop_network": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceComputeRouteCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	network, err := ParseNetworkFieldValue(d.Get("network").(string), d, config)
	if err != nil {
		return err
	}

	// Next hop data
	var nextHopInstance, nextHopIp, nextHopGateway,
		nextHopVpnTunnel string
	if v, ok := d.GetOk("next_hop_ip"); ok {
		nextHopIp = v.(string)
	}
	if v, ok := d.GetOk("next_hop_gateway"); ok {
		if v == "default-internet-gateway" {
			nextHopGateway = fmt.Sprintf("projects/%s/global/gateways/default-internet-gateway", project)
		} else {
			nextHopGateway = v.(string)
		}
	}
	if v, ok := d.GetOk("next_hop_vpn_tunnel"); ok {
		nextHopVpnTunnel = v.(string)
	}
	if v, ok := d.GetOk("next_hop_instance"); ok {
		nextHopInstanceFieldValue, err := parseComputeRouteNextHopInstanceFieldValue(v.(string), d, config)
		if err != nil {
			return fmt.Errorf("Invalid next_hop_instance: %s", err)
		}

		nextInstance, err := config.clientCompute.Instances.Get(
			project,
			nextHopInstanceFieldValue.Zone,
			nextHopInstanceFieldValue.Name).Do()
		if err != nil {
			return fmt.Errorf("Error reading instance: %s", err)
		}

		nextHopInstance = nextInstance.SelfLink
	}

	// Tags
	var tags []string
	if v := d.Get("tags").(*schema.Set); v.Len() > 0 {
		tags = make([]string, v.Len())
		for i, v := range v.List() {
			tags[i] = v.(string)
		}
	}

	// Build the route parameter
	route := &compute.Route{
		Name:             d.Get("name").(string),
		DestRange:        d.Get("dest_range").(string),
		Network:          network.RelativeLink(),
		NextHopInstance:  nextHopInstance,
		NextHopVpnTunnel: nextHopVpnTunnel,
		NextHopIp:        nextHopIp,
		NextHopGateway:   nextHopGateway,
		Priority:         int64(d.Get("priority").(int)),
		Tags:             tags,
	}
	log.Printf("[DEBUG] Route insert request: %#v", route)
	op, err := config.clientCompute.Routes.Insert(
		project, route).Do()
	if err != nil {
		return fmt.Errorf("Error creating route: %s", err)
	}

	// It probably maybe worked, so store the ID now
	d.SetId(route.Name)

	err = computeOperationWait(config.clientCompute, op, project, "Creating Route")
	if err != nil {
		return err
	}

	return resourceComputeRouteRead(d, meta)
}

func resourceComputeRouteRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	route, err := config.clientCompute.Routes.Get(
		project, d.Id()).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Route %q", d.Get("name").(string)))
	}

	nextHopInstanceFieldValue, err := parseComputeRouteNextHopInstanceFieldValue(route.NextHopInstance, d, config)
	if err != nil {
		return fmt.Errorf("Invalid next_hop_instance: %s", err)
	}

	d.Set("name", route.Name)
	d.Set("dest_range", route.DestRange)
	d.Set("network", route.Network)
	d.Set("priority", route.Priority)
	d.Set("next_hop_gateway", route.NextHopGateway)
	d.Set("next_hop_instance", nextHopInstanceFieldValue.Name)
	d.Set("next_hop_instance_zone", nextHopInstanceFieldValue.Zone)
	d.Set("next_hop_ip", route.NextHopIp)
	d.Set("next_hop_vpn_tunnel", route.NextHopVpnTunnel)
	d.Set("tags", route.Tags)
	d.Set("next_hop_network", route.NextHopNetwork)
	d.Set("self_link", route.SelfLink)

	return nil
}

func resourceComputeRouteDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	// Delete the route
	op, err := config.clientCompute.Routes.Delete(
		project, d.Id()).Do()
	if err != nil {
		return fmt.Errorf("Error deleting route: %s", err)
	}

	err = computeOperationWait(config.clientCompute, op, project, "Deleting Route")
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func parseComputeRouteNextHopInstanceFieldValue(nextHopInstance string, d TerraformResourceData, config *Config) (*ZonalFieldValue, error) {
	return parseZonalFieldValue("instances", nextHopInstance, "project", "next_hop_instance_zone", d, config, true)
}
