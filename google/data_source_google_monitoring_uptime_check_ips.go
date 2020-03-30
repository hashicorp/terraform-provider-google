package google

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceGoogleMonitoringUptimeCheckIps() *schema.Resource {
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
	config := meta.(*Config)

	url := "https://monitoring.googleapis.com/v3/uptimeCheckIps"

	uptimeCheckIps, err := paginatedListRequest("", url, config, flattenUptimeCheckIpsList)
	if err != nil {
		return fmt.Errorf("Error retrieving monitoring uptime check ips: %s", err)
	}

	if err := d.Set("uptime_check_ips", uptimeCheckIps); err != nil {
		return fmt.Errorf("Error retrieving monitoring uptime check ips: %s", err)
	}
	d.SetId(time.Now().UTC().String())
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
