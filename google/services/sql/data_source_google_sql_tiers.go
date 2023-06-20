// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package sql

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

func DataSourceGoogleSQLTiers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleSQLTiersRead,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: `Project ID of the project for which to list tiers.`,
			},
			"tiers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tier": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `An identifier for the machine type, for example, db-custom-1-3840.`,
						},
						"ram": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The maximum ram usage of this tier in bytes.`,
						},
						"disk_quota": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The maximum disk size of this tier in bytes.`,
						},
						"region": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: `The applicable regions for this tier.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceGoogleSQLTiersRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Fetching tiers for project %s", project)

	response, err := config.NewSqlAdminClient(userAgent).Tiers.List(project).Do()
	if err != nil {
		return fmt.Errorf("error retrieving tiers: %s", err)
	}

	log.Printf("[DEBUG] Fetched available tiers for project %s", project)

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("error setting project: %s", err)
	}
	if err := d.Set("tiers", flattenTiers(response.Items)); err != nil {
		return fmt.Errorf("error setting tiers: %s", err)
	}

	d.SetId(fmt.Sprintf("projects/%s", project))

	return nil
}

func flattenTiers(items []*sqladmin.Tier) []map[string]interface{} {
	var tiers []map[string]interface{}

	for _, item := range items {
		if item != nil {
			data := map[string]interface{}{
				"tier":       item.Tier,
				"ram":        item.RAM,
				"disk_quota": item.DiskQuota,
				"region":     item.Region,
			}

			tiers = append(tiers, data)
		}
	}

	return tiers
}
