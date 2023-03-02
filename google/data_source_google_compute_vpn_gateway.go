package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/compute/v1"
)

func DataSourceGoogleComputeVpnGateway() *schema.Resource {
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
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	name := d.Get("name").(string)

	vpnGatewaysService := compute.NewTargetVpnGatewaysService(config.NewComputeClient(userAgent))

	gateway, err := vpnGatewaysService.Get(project, region, name).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("VPN Gateway Not Found : %s", name))
	}
	if err := d.Set("network", ConvertSelfLinkToV1(gateway.Network)); err != nil {
		return fmt.Errorf("Error setting network: %s", err)
	}
	if err := d.Set("region", gateway.Region); err != nil {
		return fmt.Errorf("Error setting region: %s", err)
	}
	if err := d.Set("self_link", gateway.SelfLink); err != nil {
		return fmt.Errorf("Error setting self_link: %s", err)
	}
	if err := d.Set("description", gateway.Description); err != nil {
		return fmt.Errorf("Error setting description: %s", err)
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	d.SetId(fmt.Sprintf("projects/%s/regions/%s/targetVpnGateways/%s", project, region, name))
	return nil
}
