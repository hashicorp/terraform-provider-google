package google

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	dns "google.golang.org/api/dns/v1"
)

func dataSourceDnsManagedZone() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDnsManagedZoneRead,

		Schema: map[string]*schema.Schema{
			"dns_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
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

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	dnsName := d.Get("dns_name").(string)
	if name == "" && dnsName == "" {
		return errors.New("Either name or dns_name must be provided.")
	}

	var zone *dns.ManagedZone

	if name != "" {
		zone, err = config.clientDns.ManagedZones.Get(project, name).Do()
		if err != nil {
			return err
		}
	} else {
		zones, err := config.clientDns.ManagedZones.List(project).DnsName(dnsName).Do()
		if err != nil {
			return err
		}

		if len(zones.ManagedZones) == 0 {
			return fmt.Errorf("No DNS Managed Zones found for %s DNS Name", dnsName)
		}

		zone = zones.ManagedZones[0]
	}

	d.SetId(zone.Name)

	d.Set("name_servers", zone.NameServers)
	d.Set("name", zone.Name)
	d.Set("dns_name", zone.DnsName)
	d.Set("description", zone.Description)

	return nil
}
