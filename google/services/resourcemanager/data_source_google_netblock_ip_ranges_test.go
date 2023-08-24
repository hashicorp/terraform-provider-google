// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleNetblockIpRanges_basic(t *testing.T) {
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetblockIpRangesConfig,
				Check: resource.ComposeTestCheckFunc(
					// Cloud netblocks
					resource.TestMatchResourceAttr("data.google_netblock_ip_ranges.cloud",
						"cidr_blocks.#", regexp.MustCompile(("^[1-9]+[0-9]*$"))),
					resource.TestMatchResourceAttr("data.google_netblock_ip_ranges.cloud",
						"cidr_blocks.0", regexp.MustCompile("^(?:[0-9a-fA-F./:]{1,4}){1,2}.*/[0-9]{1,3}$")),
					resource.TestMatchResourceAttr("data.google_netblock_ip_ranges.cloud",
						"cidr_blocks_ipv4.#", regexp.MustCompile(("^[1-9]+[0-9]*$"))),
					resource.TestMatchResourceAttr("data.google_netblock_ip_ranges.cloud",
						"cidr_blocks_ipv4.0", regexp.MustCompile("^(?:[0-9]{1,3}.){3}[0-9]{1,3}/[0-9]{1,2}$")),
					resource.TestMatchResourceAttr("data.google_netblock_ip_ranges.cloud",
						"cidr_blocks_ipv6.#", regexp.MustCompile(("^[1-9]+[0-9]*$"))),
					resource.TestMatchResourceAttr("data.google_netblock_ip_ranges.cloud",
						"cidr_blocks_ipv6.0", regexp.MustCompile("^(?:[0-9a-fA-F]{1,4}:){1,2}.*/[0-9]{1,3}$")),
				),
			},
			{
				Config: testAccNetblockIpRangesConfig_google,
				Check: resource.ComposeTestCheckFunc(
					// Google netblocks
					resource.TestMatchResourceAttr("data.google_netblock_ip_ranges.google",
						"cidr_blocks.#", regexp.MustCompile(("^[1-9]+[0-9]*$"))),
					resource.TestMatchResourceAttr("data.google_netblock_ip_ranges.google",
						"cidr_blocks.0", regexp.MustCompile("^(?:[0-9a-fA-F./:]{1,4}){1,2}.*/[0-9]{1,3}$")),
					resource.TestMatchResourceAttr("data.google_netblock_ip_ranges.google",
						"cidr_blocks_ipv4.#", regexp.MustCompile(("^[1-9]+[0-9]*$"))),
					resource.TestMatchResourceAttr("data.google_netblock_ip_ranges.google",
						"cidr_blocks_ipv4.0", regexp.MustCompile("^(?:[0-9]{1,3}.){3}[0-9]{1,3}/[0-9]{1,2}$")),
					resource.TestMatchResourceAttr("data.google_netblock_ip_ranges.google",
						"cidr_blocks_ipv6.#", regexp.MustCompile(("^[1-9]+[0-9]*$"))),
					resource.TestMatchResourceAttr("data.google_netblock_ip_ranges.google",
						"cidr_blocks_ipv6.0", regexp.MustCompile("^(?:[0-9a-fA-F]{1,4}:){1,2}.*/[0-9]{1,3}$")),
				),
			},
			{
				Config: testAccNetblockIpRangesConfig_restricted,
				Check: resource.ComposeTestCheckFunc(
					// Private Google Access Restricted VIP
					resource.TestCheckResourceAttr("data.google_netblock_ip_ranges.restricted", "cidr_blocks.#", "1"),
					resource.TestMatchResourceAttr("data.google_netblock_ip_ranges.restricted",
						"cidr_blocks.0", regexp.MustCompile("^(?:[0-9a-fA-F./:]{1,4}){1,2}.*/[0-9]{1,3}$")),
					resource.TestCheckResourceAttr("data.google_netblock_ip_ranges.restricted", "cidr_blocks_ipv4.#", "1"),
					resource.TestMatchResourceAttr("data.google_netblock_ip_ranges.restricted",
						"cidr_blocks_ipv4.0", regexp.MustCompile("^(?:[0-9]{1,3}.){3}[0-9]{1,3}/[0-9]{1,2}$")),
					resource.TestCheckResourceAttr("data.google_netblock_ip_ranges.restricted", "cidr_blocks_ipv6.#", "0"),
				),
			},
			{
				Config: testAccNetblockIpRangesConfig_private,
				Check: resource.ComposeTestCheckFunc(
					// Private Google Access Unrestricted VIP
					resource.TestCheckResourceAttr("data.google_netblock_ip_ranges.private", "cidr_blocks.#", "1"),
					resource.TestMatchResourceAttr("data.google_netblock_ip_ranges.private",
						"cidr_blocks.0", regexp.MustCompile("^(?:[0-9a-fA-F./:]{1,4}){1,2}.*/[0-9]{1,3}$")),
					resource.TestCheckResourceAttr("data.google_netblock_ip_ranges.private", "cidr_blocks_ipv4.#", "1"),
					resource.TestMatchResourceAttr("data.google_netblock_ip_ranges.private",
						"cidr_blocks_ipv4.0", regexp.MustCompile("^(?:[0-9]{1,3}.){3}[0-9]{1,3}/[0-9]{1,2}$")),
					resource.TestCheckResourceAttr("data.google_netblock_ip_ranges.private", "cidr_blocks_ipv6.#", "0"),
				),
			},
			{
				Config: testAccNetblockIpRangesConfig_dns,
				Check: resource.ComposeTestCheckFunc(
					// DNS outbound forwarding
					resource.TestCheckResourceAttr("data.google_netblock_ip_ranges.dns", "cidr_blocks.#", "1"),
					resource.TestMatchResourceAttr("data.google_netblock_ip_ranges.dns",
						"cidr_blocks.0", regexp.MustCompile("^(?:[0-9a-fA-F./:]{1,4}){1,2}.*/[0-9]{1,3}$")),
					resource.TestCheckResourceAttr("data.google_netblock_ip_ranges.dns", "cidr_blocks_ipv4.#", "1"),
					resource.TestMatchResourceAttr("data.google_netblock_ip_ranges.dns",
						"cidr_blocks_ipv4.0", regexp.MustCompile("^(?:[0-9]{1,3}.){3}[0-9]{1,3}/[0-9]{1,2}$")),
					resource.TestCheckResourceAttr("data.google_netblock_ip_ranges.dns", "cidr_blocks_ipv6.#", "0"),
				),
			},
			{
				Config: testAccNetblockIpRangesConfig_iap,
				Check: resource.ComposeTestCheckFunc(
					// IAP sources
					resource.TestCheckResourceAttr("data.google_netblock_ip_ranges.iap", "cidr_blocks.#", "1"),
					resource.TestMatchResourceAttr("data.google_netblock_ip_ranges.iap",
						"cidr_blocks.0", regexp.MustCompile("^(?:[0-9a-fA-F./:]{1,4}){1,2}.*/[0-9]{1,3}$")),
					resource.TestCheckResourceAttr("data.google_netblock_ip_ranges.iap", "cidr_blocks_ipv4.#", "1"),
					resource.TestMatchResourceAttr("data.google_netblock_ip_ranges.iap",
						"cidr_blocks_ipv4.0", regexp.MustCompile("^(?:[0-9]{1,3}.){3}[0-9]{1,3}/[0-9]{1,2}$")),
					resource.TestCheckResourceAttr("data.google_netblock_ip_ranges.iap", "cidr_blocks_ipv6.#", "0"),
				),
			},
			{
				Config: testAccNetblockIpRangesConfig_hc,
				Check: resource.ComposeTestCheckFunc(
					// Modern health checkers
					resource.TestCheckResourceAttr("data.google_netblock_ip_ranges.hc", "cidr_blocks.#", "2"),
					resource.TestMatchResourceAttr("data.google_netblock_ip_ranges.hc",
						"cidr_blocks.0", regexp.MustCompile("^(?:[0-9a-fA-F./:]{1,4}){1,2}.*/[0-9]{1,3}$")),
					resource.TestCheckResourceAttr("data.google_netblock_ip_ranges.hc", "cidr_blocks_ipv4.#", "2"),
					resource.TestMatchResourceAttr("data.google_netblock_ip_ranges.hc",
						"cidr_blocks_ipv4.0", regexp.MustCompile("^(?:[0-9]{1,3}.){3}[0-9]{1,3}/[0-9]{1,2}$")),
					resource.TestCheckResourceAttr("data.google_netblock_ip_ranges.hc", "cidr_blocks_ipv6.#", "0"),
				),
			},
			{
				Config: testAccNetblockIpRangesConfig_lhc,
				Check: resource.ComposeTestCheckFunc(
					// Legacy health checkers
					resource.TestCheckResourceAttr("data.google_netblock_ip_ranges.lhc", "cidr_blocks.#", "3"),
					resource.TestMatchResourceAttr("data.google_netblock_ip_ranges.lhc",
						"cidr_blocks.0", regexp.MustCompile("^(?:[0-9a-fA-F./:]{1,4}){1,2}.*/[0-9]{1,3}$")),
					resource.TestCheckResourceAttr("data.google_netblock_ip_ranges.lhc", "cidr_blocks_ipv4.#", "3"),
					resource.TestMatchResourceAttr("data.google_netblock_ip_ranges.lhc",
						"cidr_blocks_ipv4.0", regexp.MustCompile("^(?:[0-9]{1,3}.){3}[0-9]{1,3}/[0-9]{1,2}$")),
					resource.TestCheckResourceAttr("data.google_netblock_ip_ranges.lhc", "cidr_blocks_ipv6.#", "0"),
				),
			},
		},
	})
}

const testAccNetblockIpRangesConfig = `
data "google_netblock_ip_ranges" "cloud" {}
`

const testAccNetblockIpRangesConfig_google = `
data "google_netblock_ip_ranges" "google" {
  range_type = "google-netblocks"
}
`

const testAccNetblockIpRangesConfig_restricted = `
data "google_netblock_ip_ranges" "restricted" {
  range_type = "restricted-googleapis"
}
`

const testAccNetblockIpRangesConfig_private = `
data "google_netblock_ip_ranges" "private" {
  range_type = "private-googleapis"
}
`

const testAccNetblockIpRangesConfig_dns = `
data "google_netblock_ip_ranges" "dns" {
  range_type = "dns-forwarders"
}
`

const testAccNetblockIpRangesConfig_iap = `
data "google_netblock_ip_ranges" "iap" {
  range_type = "iap-forwarders"
}
`

const testAccNetblockIpRangesConfig_hc = `
data "google_netblock_ip_ranges" "hc" {
  range_type = "health-checkers"
}
`

const testAccNetblockIpRangesConfig_lhc = `
data "google_netblock_ip_ranges" "lhc" {
  range_type = "legacy-health-checkers"
}
`
