package google

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceGoogleMonitoringUptimeCheckIps_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleMonitoringUptimeCheckIps_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.google_monitoring_uptime_check_ips.foobar",
						"uptime_check_ips.0.ip_address", regexp.MustCompile("^\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}$")),
					resource.TestMatchResourceAttr("data.google_monitoring_uptime_check_ips.foobar",
						"uptime_check_ips.0.location", regexp.MustCompile("^[A-Z].+$")),
					resource.TestMatchResourceAttr("data.google_monitoring_uptime_check_ips.foobar",
						"uptime_check_ips.0.region", regexp.MustCompile("^[A-Z_]+$")),
				),
			},
		},
	})
}

const testAccDataSourceGoogleMonitoringUptimeCheckIps_basic = `
data "google_monitoring_uptime_check_ips" "foobar" {
}
`
