package google

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceComputeLbIpRanges_basic(t *testing.T) {
	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
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
