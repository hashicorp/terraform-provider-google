package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDnsRecordSet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDnsRecordSetRead,

		Schema: map[string]*schema.Schema{
			"managed_zone": {
				Type:     schema.TypeString,
				Required: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"rrdatas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"ttl": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"type": {
				Type:     schema.TypeString,
				Required: true,
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceDnsRecordSetRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	zone := d.Get("managed_zone").(string)
	name := d.Get("name").(string)
	dnsType := d.Get("type").(string)
	d.SetId(fmt.Sprintf("projects/%s/managedZones/%s/rrsets/%s/%s", project, zone, name, dnsType))

	resp, err := config.NewDnsClient(userAgent).ResourceRecordSets.List(project, zone).Name(name).Type(dnsType).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("dataSourceDnsRecordSet %q", name))
	}
	if len(resp.Rrsets) != 1 {
		return fmt.Errorf("Only expected 1 record set, got %d", len(resp.Rrsets))
	}

	if err := d.Set("rrdatas", resp.Rrsets[0].Rrdatas); err != nil {
		return fmt.Errorf("Error setting rrdatas: %s", err)
	}
	if err := d.Set("ttl", resp.Rrsets[0].Ttl); err != nil {
		return fmt.Errorf("Error setting ttl: %s", err)
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}

	return nil
}
