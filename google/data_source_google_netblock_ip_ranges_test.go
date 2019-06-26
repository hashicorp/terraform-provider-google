package google

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceGoogleNetblockIpRanges_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetblockIpRangesConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.google_netblock_ip_ranges.some",
						"cidr_blocks.#", regexp.MustCompile(("^[1-9]+[0-9]*$"))),
					resource.TestMatchResourceAttr("data.google_netblock_ip_ranges.some",
						"cidr_blocks.0", regexp.MustCompile("^(?:[0-9a-fA-F./:]{1,4}){1,2}.*/[0-9]{1,3}$")),
					resource.TestMatchResourceAttr("data.google_netblock_ip_ranges.some",
						"cidr_blocks_ipv4.#", regexp.MustCompile(("^[1-9]+[0-9]*$"))),
					resource.TestMatchResourceAttr("data.google_netblock_ip_ranges.some",
						"cidr_blocks_ipv4.0", regexp.MustCompile("^(?:[0-9]{1,3}.){3}[0-9]{1,3}/[0-9]{1,2}$")),
					resource.TestMatchResourceAttr("data.google_netblock_ip_ranges.some",
						"cidr_blocks_ipv6.#", regexp.MustCompile(("^[1-9]+[0-9]*$"))),
					resource.TestMatchResourceAttr("data.google_netblock_ip_ranges.some",
						"cidr_blocks_ipv6.0", regexp.MustCompile("^(?:[0-9a-fA-F]{1,4}:){1,2}.*/[0-9]{1,3}$")),
				),
			},
		},
	})
}

const testAccNetblockIpRangesConfig = `
data "google_netblock_ip_ranges" "some" {}
`
