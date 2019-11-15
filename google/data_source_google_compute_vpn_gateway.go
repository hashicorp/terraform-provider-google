package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"google.golang.org/api/compute/v1"
)

func dataSourceGoogleComputeVpnGateway() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeVpnGatewayRead,

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

			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"network": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGoogleComputeVpnGatewayRead(d *schema.ResourceData, meta interface{}) error {
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

	vpnGatewaysService := compute.NewTargetVpnGatewaysService(config.clientCompute)

	gateway, err := vpnGatewaysService.Get(project, region, name).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("VPN Gateway Not Found : %s", name))
	}
	d.Set("network", gateway.Network)
	d.Set("region", gateway.Region)
	d.Set("self_link", gateway.SelfLink)
	d.Set("description", gateway.Description)
	d.Set("project", project)
	d.SetId(fmt.Sprintf("projects/%s/regions/%s/targetVpnGateways/%s", project, region, name))
	return nil
}
