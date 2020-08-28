package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGoogleComputeNetwork() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeNetworkRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"gateway_ipv4": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"subnetworks_self_links": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceGoogleComputeNetworkRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	network, err := config.clientCompute.Networks.Get(project, name).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Network Not Found : %s", name))
	}
	if err := d.Set("gateway_ipv4", network.GatewayIPv4); err != nil {
		return fmt.Errorf("Error reading gateway_ipv4: %s", err)
	}
	if err := d.Set("self_link", network.SelfLink); err != nil {
		return fmt.Errorf("Error reading self_link: %s", err)
	}
	if err := d.Set("description", network.Description); err != nil {
		return fmt.Errorf("Error reading description: %s", err)
	}
	if err := d.Set("subnetworks_self_links", network.Subnetworks); err != nil {
		return fmt.Errorf("Error reading subnetworks_self_links: %s", err)
	}
	d.SetId(fmt.Sprintf("projects/%s/global/networks/%s", project, network.Name))
	return nil
}
