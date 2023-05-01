package google

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleMonitoringUptimeCheckIps_basic(t *testing.T) {
	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
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
