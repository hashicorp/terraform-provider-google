package google

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceGoogleComputeForwardingRule() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeForwardingRuleRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"target": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"backend_service": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"ip_protocol": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"load_balancing_scheme": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"network": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"port_range": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"ports": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},

			"self_link": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"subnetwork": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGoogleComputeForwardingRuleRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	name := d.Get("name").(string)

	frule, err := config.clientCompute.ForwardingRules.Get(
		project, region, name).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Forwarding Rule Not Found : %s", name))
	}
	d.SetId(frule.Name)

	d.Set("self_link", frule.SelfLink)
	d.Set("description", frule.Description)
	d.Set("backend_service", frule.BackendService)
	d.Set("ip_address", frule.IPAddress)
	d.Set("ip_protocol", frule.IPProtocol)
	d.Set("load_balancing_scheme", frule.LoadBalancingScheme)
	d.Set("name", frule.Name)
	d.Set("port_range", frule.PortRange)
	d.Set("ports", frule.Ports)
	d.Set("subnetwork", frule.Subnetwork)
	d.Set("network", frule.Network)
	d.Set("target", frule.Target)
	d.Set("project", project)
	d.Set("region", region)

	return nil
}
