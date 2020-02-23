package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

	d.Set("name_servers", zone.NameServers)
	d.Set("name", zone.Name)
	d.Set("dns_name", zone.DnsName)
	d.Set("description", zone.Description)
	d.Set("visibility", zone.Visibility)
	d.Set("project", project)

	return nil
}
