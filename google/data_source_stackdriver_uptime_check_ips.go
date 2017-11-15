package google

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/monitoring/v3"
)

func dataSourceStackdriverUptimeCheckIps() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceStackdriverUptimeCheckIpsRead,

		Schema: map[string]*schema.Schema{
			"uptime_check_ips": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		},
	}
}

func dataSourceStackdriverUptimeCheckIpsRead(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)

	// https://cloud.google.com/monitoring/api/ref_v3/rest/v3/uptimeCheckIps/list

	resp, err := config.clientMonitoring.UptimeCheckIps.List().Do()
	if err != nil {
		return fmt.Errorf("Error retrieving Stackdriver Uptime Check IPs: %s", err.Error())
	}

	uptimeCheckIps := flattenUptimeCheckIpAddresses(resp.UptimeCheckIps)
	d.Set("uptime_check_ips", uptimeCheckIps)

	log.Printf("[DEBUG] Received Stackdriver Uptime Check IPs: %q", uptimeCheckIps)

	d.SetId(time.Now().UTC().String())

	return nil
}

func flattenUptimeCheckIpAddresses(uptimeCheckIps []*monitoring.UptimeCheckIp) []string {
	result := make([]string, len(uptimeCheckIps), len(uptimeCheckIps))
	for i, uptimeCheckIp := range uptimeCheckIps {
		result[i] = uptimeCheckIp.IpAddress
	}
	sort.Strings(result)
	return result
}
