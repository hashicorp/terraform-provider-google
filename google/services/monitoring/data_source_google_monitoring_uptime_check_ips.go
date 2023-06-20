// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package monitoring

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleMonitoringUptimeCheckIps() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleMonitoringUptimeCheckIpsRead,

		Schema: map[string]*schema.Schema{
			"uptime_check_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"location": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceGoogleMonitoringUptimeCheckIpsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url := "https://monitoring.googleapis.com/v3/uptimeCheckIps"

	uptimeCheckIps, err := tpgresource.PaginatedListRequest("", url, userAgent, config, flattenUptimeCheckIpsList)
	if err != nil {
		return fmt.Errorf("Error retrieving monitoring uptime check ips: %s", err)
	}

	if err := d.Set("uptime_check_ips", uptimeCheckIps); err != nil {
		return fmt.Errorf("Error retrieving monitoring uptime check ips: %s", err)
	}
	d.SetId("uptime_check_ips_id")
	return nil
}

func flattenUptimeCheckIpsList(resp map[string]interface{}) []interface{} {
	ipObjList := resp["uptimeCheckIps"].([]interface{})
	uptimeCheckIps := make([]interface{}, len(ipObjList))
	for i, u := range ipObjList {
		ipObj := u.(map[string]interface{})
		uptimeCheckIps[i] = map[string]interface{}{
			"region":     ipObj["region"],
			"location":   ipObj["location"],
			"ip_address": ipObj["ipAddress"],
		}
	}
	return uptimeCheckIps
}
