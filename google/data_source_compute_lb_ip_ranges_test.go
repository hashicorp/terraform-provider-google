package google

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceComputeLbIpRanges_basic(t *testing.T) {
	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeLbIpRangesConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.google_compute_lb_ip_ranges.some",
						"network.#", regexp.MustCompile("^[1-9]+[0-9]*$")),
					resource.TestMatchResourceAttr("data.google_compute_lb_ip_ranges.some",
						"network.0", regexp.MustCompile("^[0-9./]+$")),
					resource.TestMatchResourceAttr("data.google_compute_lb_ip_ranges.some",
						"http_ssl_tcp_internal.#", regexp.MustCompile("^[1-9]+[0-9]*$")),
					resource.TestMatchResourceAttr("data.google_compute_lb_ip_ranges.some",
						"http_ssl_tcp_internal.0", regexp.MustCompile("^[0-9./]+$")),
				),
			},
		},
	})
}

const testAccComputeLbIpRangesConfig = `
data "google_compute_lb_ip_ranges" "some" {
}
`
