package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGoogleComputeGlobalAddress() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeGlobalAddressRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"project": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func dataSourceGoogleComputeGlobalAddressRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	address, err := config.clientCompute.GlobalAddresses.Get(project, name).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Global Address Not Found : %s", name))
	}

	if err := d.Set("address", address.Address); err != nil {
		return fmt.Errorf("Error reading address: %s", err)
	}
	if err := d.Set("status", address.Status); err != nil {
		return fmt.Errorf("Error reading status: %s", err)
	}
	if err := d.Set("self_link", address.SelfLink); err != nil {
		return fmt.Errorf("Error reading self_link: %s", err)
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading project: %s", err)
	}
	d.SetId(fmt.Sprintf("projects/%s/global/addresses/%s", project, name))
	return nil
}
