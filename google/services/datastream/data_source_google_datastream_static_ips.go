// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package datastream

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceGoogleDatastreamStaticIps() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleDatastreamStaticIpsRead,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"location": {
				Type:     schema.TypeString,
				Required: true,
			},
			"static_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceGoogleDatastreamStaticIpsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{DatastreamBasePath}}projects/{{project}}/locations/{{location}}:fetchStaticIps")
	if err != nil {
		return err
	}

	staticIps, err := tpgresource.PaginatedListRequest(project, url, userAgent, config, flattenStaticIpsList)
	if err != nil {
		return fmt.Errorf("Error retrieving monitoring uptime check ips: %s", err)
	}

	if err := d.Set("static_ips", staticIps); err != nil {
		return fmt.Errorf("Error retrieving monitoring uptime check ips: %s", err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}:fetchStaticIps")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)
	return nil
}

func flattenStaticIpsList(resp map[string]interface{}) []interface{} {
	ipList := resp["staticIps"].([]interface{})
	staticIps := make([]interface{}, len(ipList))
	for i, u := range ipList {
		staticIps[i] = u
	}
	return staticIps
}
