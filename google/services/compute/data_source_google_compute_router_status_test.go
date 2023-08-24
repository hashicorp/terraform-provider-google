// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceComputeRouterStatus(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"suffix": acctest.RandString(t, 10),
		"region": "us-central1",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccCheckComputeRouterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceComputeRouterStatusConfig(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.google_compute_router_status.router1", "name", "google_compute_router.router1", "name"),
					resource.TestCheckResourceAttrPair("data.google_compute_router_status.router1", "region", "google_compute_router.router1", "region"),
					resource.TestCheckResourceAttrSet("data.google_compute_router_status.router1", "network"),
					resource.TestCheckResourceAttr("data.google_compute_router_status.router1", "best_routes.#", "2"),
					resource.TestCheckResourceAttr("data.google_compute_router_status.router1", "best_routes_for_router.#", "2"),
					resource.TestCheckResourceAttrPair("data.google_compute_router_status.router1", "best_routes.0.next_hop_ip", "google_compute_router_peer.router1_peer1", "peer_ip_address"),
					resource.TestCheckResourceAttrSet("data.google_compute_router_status.router1", "best_routes.0.next_hop_vpn_tunnel"),
				),
			},
		},
	})
}

func testAccDataSourceComputeRouterStatusConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "network1" {
  name                    = "network1-%{suffix}"
  routing_mode            = "GLOBAL"
  auto_create_subnetworks = false
}

resource "google_compute_network" "network2" {
  name                    = "network2-%{suffix}"
  routing_mode            = "GLOBAL"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "network2_subnet1" {
  name          = "ha-vpn-subnet-1-%{suffix}"
  ip_cidr_range = "192.168.1.0/24"
  region        = "%{region}"
  network       = google_compute_network.network2.id
}

resource "google_compute_subnetwork" "network2_subnet2" {
  name          = "ha-vpn-subnet-2-%{suffix}"
  ip_cidr_range = "192.168.2.0/24"
  region        = "us-east1"
  network       = google_compute_network.network2.id
}

resource "google_compute_router" "router1" {
  name    = "tf-test-ha-vpn-router1-%{suffix}"
  network = google_compute_network.network1.name
  region  = "%{region}"
  bgp {
    asn = 64514
  }
}

resource "google_compute_router" "router2" {
  name    = "tf-test-ha-vpn-router2-%{suffix}"
  network = google_compute_network.network2.name
  region  = "%{region}"
  bgp {
    asn = 64515
  }
}

resource "google_compute_ha_vpn_gateway" "ha_gateway1" {
  region  = "%{region}"
  name    = "tf-test-ha-vpn-1-%{suffix}"
  network = google_compute_network.network1.id
}

resource "google_compute_ha_vpn_gateway" "ha_gateway2" {
  region  = "%{region}"
  name    = "tf-test-ha-vpn-2-%{suffix}"
  network = google_compute_network.network2.id
}

resource "google_compute_vpn_tunnel" "tunnel1" {
  name                  = "ha-vpn-tunnel1-%{suffix}"
  region                = "%{region}"
  vpn_gateway           = google_compute_ha_vpn_gateway.ha_gateway1.id
  peer_gcp_gateway      = google_compute_ha_vpn_gateway.ha_gateway2.id
  shared_secret         = "a secret message"
  router                = google_compute_router.router1.id
  vpn_gateway_interface = 0
}

resource "google_compute_vpn_tunnel" "tunnel2" {
  name                  = "ha-vpn-tunnel2-%{suffix}"
  region                = "%{region}"
  vpn_gateway           = google_compute_ha_vpn_gateway.ha_gateway2.id
  peer_gcp_gateway      = google_compute_ha_vpn_gateway.ha_gateway1.id
  shared_secret         = "a secret message"
  router                = google_compute_router.router2.id
  vpn_gateway_interface = 0
}

resource "google_compute_router_interface" "router1_interface1" {
  name       = "router1-interface1-%{suffix}"
  router     = google_compute_router.router1.name
  region     = "%{region}"
  ip_range   = "169.254.0.1/30"
  vpn_tunnel = google_compute_vpn_tunnel.tunnel1.name
}

resource "google_compute_router_peer" "router1_peer1" {
  name                      = "router1-peer1-%{suffix}"
  router                    = google_compute_router.router1.name
  region                    = "%{region}"
  peer_ip_address           = "169.254.0.2"
  peer_asn                  = 64515
  advertised_route_priority = 100
  interface                 = google_compute_router_interface.router1_interface1.name
}

resource "google_compute_router_interface" "router2_interface1" {
  name       = "router2-interface1-%{suffix}"
  router     = google_compute_router.router2.name
  region     = "%{region}"
  ip_range   = "169.254.0.2/30"
  vpn_tunnel = google_compute_vpn_tunnel.tunnel2.name
}

resource "google_compute_router_peer" "router2_peer1" {
  name                      = "router2-peer1-%{suffix}"
  router                    = google_compute_router.router2.name
  region                    = "%{region}"
  peer_ip_address           = "169.254.0.1"
  peer_asn                  = 64514
  advertised_route_priority = 100
  interface                 = google_compute_router_interface.router2_interface1.name
}

resource "time_sleep" "wait_60_seconds" {
  create_duration = "60s"

  depends_on = [
    google_compute_router_peer.router1_peer1,
    google_compute_router_peer.router2_peer1,
  ]
}

data "google_compute_router_status" "router1" {
  name   = google_compute_router.router1.name
  region = google_compute_router.router1.region

  depends_on = [time_sleep.wait_60_seconds]
}
`, context)

}
