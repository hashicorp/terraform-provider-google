package google

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceStackdriverUptimeCheckIps_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccStackdriverUptimeCheckIpsConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.google_stackdriver_uptime_check_ips.some",
						"uptime_check_ips.#", regexp.MustCompile("^[1-9]+[0-9]*$")),
					resource.TestMatchResourceAttr("data.google_stackdriver_uptime_check_ips.some",
						"uptime_check_ips.0", regexp.MustCompile("^[0-9./]+$")),
				),
			},
		},
	})
}

const testAccStackdriverUptimeCheckIpsConfig = `
data "google_stackdriver_uptime_check_ips" "some" {}
`
