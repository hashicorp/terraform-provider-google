// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package dns

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/dns/v1"
)

func DataSourceDnsManagedZones() *schema.Resource {

	mzSchema := DataSourceDnsManagedZone().Schema
	tpgresource.AddOptionalFieldsToSchema(mzSchema, "name")

	return &schema.Resource{
		Read: dataSourceDnsManagedZonesRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"managed_zones": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: mzSchema,
				},
			},

			// Google Cloud DNS ManagedZone resources do not have a SelfLink attribute.
			"project": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceDnsManagedZonesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("projects/%s/managedZones", project))

	zones, err := config.NewDnsClient(userAgent).ManagedZones.List(project).Do()
	if err != nil {
		return err
	}

	if err := d.Set("managed_zones", flattenZones(zones.ManagedZones, project)); err != nil {
		return fmt.Errorf("error setting managed_zones: %s", err)
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("error setting project: %s", err)
	}

	return nil
}

// flattenZones flattens the list of managed zones into a format that can be assigned to the managed_zones field
// on the plural datasource. This includes setting the project value for each item, as this isn't returned by the API.
func flattenZones(items []*dns.ManagedZone, project string) []map[string]interface{} {
	var zones []map[string]interface{}

	for _, item := range items {
		if item != nil {
			data := map[string]interface{}{
				"id":              fmt.Sprintf("projects/%s/managedZones/%s", project, item.Name), // Matches construction in singlur data source
				"dns_name":        item.DnsName,
				"name":            item.Name,
				"managed_zone_id": item.Id,
				"description":     item.Description,
				"visibility":      item.Visibility,
				"name_servers":    item.NameServers,
				"project":         project,
			}

			zones = append(zones, data)
		}
	}

	return zones
}
