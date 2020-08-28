package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGoogleComputeForwardingRule() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeForwardingRuleRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"target": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"backend_service": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ip_protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"load_balancing_scheme": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"network": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"port_range": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ports": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},

			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"subnetwork": {
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
	d.SetId(fmt.Sprintf("projects/%s/regions/%s/forwardingRules/%s", project, region, name))

	if err := d.Set("self_link", frule.SelfLink); err != nil {
		return fmt.Errorf("Error reading self_link: %s", err)
	}
	if err := d.Set("description", frule.Description); err != nil {
		return fmt.Errorf("Error reading description: %s", err)
	}
	if err := d.Set("backend_service", frule.BackendService); err != nil {
		return fmt.Errorf("Error reading backend_service: %s", err)
	}
	if err := d.Set("ip_address", frule.IPAddress); err != nil {
		return fmt.Errorf("Error reading ip_address: %s", err)
	}
	if err := d.Set("ip_protocol", frule.IPProtocol); err != nil {
		return fmt.Errorf("Error reading ip_protocol: %s", err)
	}
	if err := d.Set("load_balancing_scheme", frule.LoadBalancingScheme); err != nil {
		return fmt.Errorf("Error reading load_balancing_scheme: %s", err)
	}
	if err := d.Set("name", frule.Name); err != nil {
		return fmt.Errorf("Error reading name: %s", err)
	}
	if err := d.Set("port_range", frule.PortRange); err != nil {
		return fmt.Errorf("Error reading port_range: %s", err)
	}
	if err := d.Set("ports", frule.Ports); err != nil {
		return fmt.Errorf("Error reading ports: %s", err)
	}
	if err := d.Set("subnetwork", frule.Subnetwork); err != nil {
		return fmt.Errorf("Error reading subnetwork: %s", err)
	}
	if err := d.Set("network", frule.Network); err != nil {
		return fmt.Errorf("Error reading network: %s", err)
	}
	if err := d.Set("target", frule.Target); err != nil {
		return fmt.Errorf("Error reading target: %s", err)
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading project: %s", err)
	}
	if err := d.Set("region", region); err != nil {
		return fmt.Errorf("Error reading region: %s", err)
	}

	return nil
}
