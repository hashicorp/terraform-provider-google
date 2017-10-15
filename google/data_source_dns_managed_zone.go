package google

import "github.com/hashicorp/terraform/helper/schema"

func dataSourceDnsManagedZone() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDnsManagedZoneRead,

		Schema: map[string]*schema.Schema{
			"dns_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"name_servers": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			// Google Cloud DNS ManagedZone resources do not have a SelfLink attribute.

			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceDnsManagedZoneRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	d.SetId(d.Get("name").(string))

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	zone, err := config.clientDns.ManagedZones.Get(
		project, d.Id()).Do()
	if err != nil {
		return err
	}

	d.Set("name_servers", zone.NameServers)
	d.Set("name", zone.Name)
	d.Set("dns_name", zone.DnsName)
	d.Set("description", zone.Description)

	return nil
}
