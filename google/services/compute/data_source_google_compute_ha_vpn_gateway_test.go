// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceComputeHaVpnGateway(t *testing.T) {
	t.Parallel()

	gwName := fmt.Sprintf("tf-%s", acctest.RandString(t, 10))
	gatewayIpVersion := "IPV6"
	stackType := "IPV6_ONLY"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceComputeHaVpnGatewayConfig(gwName),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_compute_ha_vpn_gateway.ha_gateway", "google_compute_ha_vpn_gateway.ha_gateway"),
					resource.TestCheckResourceAttr("data.google_compute_ha_vpn_gateway.ha_gateway", "gateway_ip_version", "IPV4"),
					resource.TestCheckResourceAttr("data.google_compute_ha_vpn_gateway.ha_gateway", "stack_type", "IPV4_ONLY"),
				),
			}, {
				Config: testAccDataSourceComputeHaVpnGatewayFields(fmt.Sprintf("%s-2", gwName), gatewayIpVersion, stackType),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_compute_ha_vpn_gateway.ha_gateway", "google_compute_ha_vpn_gateway.ha_gateway"),
					resource.TestCheckResourceAttr("data.google_compute_ha_vpn_gateway.ha_gateway", "gateway_ip_version", gatewayIpVersion),
					resource.TestCheckResourceAttr("data.google_compute_ha_vpn_gateway.ha_gateway", "stack_type", stackType),
				),
			},
		},
	})
}

func testAccDataSourceComputeHaVpnGatewayConfig(gwName string) string {
	return fmt.Sprintf(`
resource "google_compute_ha_vpn_gateway" "ha_gateway" {
  name     = "%s"
  network  = google_compute_network.network1.id
}

resource "google_compute_network" "network1" {
  name                    = "%s"
  auto_create_subnetworks = false
}

data "google_compute_ha_vpn_gateway" "ha_gateway" {
  name = google_compute_ha_vpn_gateway.ha_gateway.name
}
`, gwName, gwName)
}

func testAccDataSourceComputeHaVpnGatewayFields(gwName, gatewayIpVersion, stackType string) string {
	return fmt.Sprintf(`
resource "google_compute_ha_vpn_gateway" "ha_gateway" {
  name     			 = "%s"
  network  			 = google_compute_network.network1.id
  gateway_ip_version = "%s"
  stack_type		 = "%s"
}

resource "google_compute_network" "network1" {
  name                    = "%s"
  auto_create_subnetworks = false
}

data "google_compute_ha_vpn_gateway" "ha_gateway" {
  name = google_compute_ha_vpn_gateway.ha_gateway.name
}
`, gwName, gatewayIpVersion, stackType, gwName)
}
