package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDnsManagedZone() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDnsManagedZoneRead,

		Schema: map[string]*schema.Schema{
			"dns_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name_servers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"visibility": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// Google Cloud DNS ManagedZone resources do not have a SelfLink attribute.
			"project": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceDnsManagedZoneRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	d.SetId(fmt.Sprintf("projects/%s/managedZones/%s", project, name))

	zone, err := config.clientDns.ManagedZones.Get(
		project, name).Do()
	if err != nil {
		return err
	}

	if err := d.Set("name_servers", zone.NameServers); err != nil {
		return fmt.Errorf("Error reading name_servers: %s", err)
	}
	if err := d.Set("name", zone.Name); err != nil {
		return fmt.Errorf("Error reading name: %s", err)
	}
	if err := d.Set("dns_name", zone.DnsName); err != nil {
		return fmt.Errorf("Error reading dns_name: %s", err)
	}
	if err := d.Set("description", zone.Description); err != nil {
		return fmt.Errorf("Error reading description: %s", err)
	}
	if err := d.Set("visibility", zone.Visibility); err != nil {
		return fmt.Errorf("Error reading visibility: %s", err)
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading project: %s", err)
	}

	return nil
}
