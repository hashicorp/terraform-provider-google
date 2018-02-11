package google

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

func dataSourceGoogleComputeVpnGateway() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeVpnGatewayRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"self_link": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"network": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
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
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			// The resource doesn't exist anymore

			return fmt.Errorf("VPN Gateway Not Found : %s", d.Get("name"))
		}

		return fmt.Errorf("Error reading VPN Gateway: %s", err)
	}
	d.Set("network", gateway.Network)
	d.Set("region", gateway.Region)
	d.Set("self_link", gateway.SelfLink)
	d.Set("description", gateway.Description)
	d.SetId(gateway.Name)
	return nil
}
