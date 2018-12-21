package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccComputeVpnTunnel_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeVpnTunnelDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeVpnTunnel_basic(),
			},
			resource.TestStep{
				ResourceName:            "google_compute_vpn_tunnel.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"shared_secret"},
			},
		},
	})
}

func TestAccComputeVpnTunnel_router(t *testing.T) {
	t.Parallel()

	router := fmt.Sprintf("tunnel-test-router-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeVpnTunnelDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeVpnTunnelRouter(router),
			},
			resource.TestStep{
				ResourceName:            "google_compute_vpn_tunnel.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"shared_secret"},
			},
		},
	})
}

func TestAccComputeVpnTunnel_defaultTrafficSelectors(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeVpnTunnelDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeVpnTunnelDefaultTrafficSelectors(),
			},
			resource.TestStep{
				ResourceName:            "google_compute_vpn_tunnel.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"shared_secret"},
			},
		},
	})
}

func testAccComputeVpnTunnel_basic() string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
	name = "tunnel-test-%s"
}
resource "google_compute_subnetwork" "foobar" {
	name = "tunnel-test-subnetwork-%s"
	network = "${google_compute_network.foobar.self_link}"
	ip_cidr_range = "10.0.0.0/16"
	region = "us-central1"
}
resource "google_compute_address" "foobar" {
	name = "tunnel-test-%s"
	region = "${google_compute_subnetwork.foobar.region}"
}
resource "google_compute_vpn_gateway" "foobar" {
	name = "tunnel-test-%s"
	network = "${google_compute_network.foobar.self_link}"
	region = "${google_compute_subnetwork.foobar.region}"
}
resource "google_compute_forwarding_rule" "foobar_esp" {
	name = "tunnel-test-%s"
	region = "${google_compute_vpn_gateway.foobar.region}"
	ip_protocol = "ESP"
	ip_address = "${google_compute_address.foobar.address}"
	target = "${google_compute_vpn_gateway.foobar.self_link}"
}
resource "google_compute_forwarding_rule" "foobar_udp500" {
	name = "tunnel-test-%s"
	region = "${google_compute_forwarding_rule.foobar_esp.region}"
	ip_protocol = "UDP"
	port_range = "500-500"
	ip_address = "${google_compute_address.foobar.address}"
	target = "${google_compute_vpn_gateway.foobar.self_link}"
}
resource "google_compute_forwarding_rule" "foobar_udp4500" {
	name = "tunnel-test-%s"
	region = "${google_compute_forwarding_rule.foobar_udp500.region}"
	ip_protocol = "UDP"
	port_range = "4500-4500"
	ip_address = "${google_compute_address.foobar.address}"
	target = "${google_compute_vpn_gateway.foobar.self_link}"
}
resource "google_compute_vpn_tunnel" "foobar" {
	name = "tunnel-test-%s"
	region = "${google_compute_forwarding_rule.foobar_udp4500.region}"
	target_vpn_gateway = "${google_compute_vpn_gateway.foobar.self_link}"
	shared_secret = "unguessable"
	peer_ip = "8.8.8.8"
	local_traffic_selector = ["${google_compute_subnetwork.foobar.ip_cidr_range}"]
	remote_traffic_selector = ["192.168.0.0/24", "192.168.1.0/24"]
}`, acctest.RandString(10), acctest.RandString(10), acctest.RandString(10),
		acctest.RandString(10), acctest.RandString(10), acctest.RandString(10),
		acctest.RandString(10), acctest.RandString(10))
}

func testAccComputeVpnTunnelRouter(router string) string {
	testId := acctest.RandString(10)
	return fmt.Sprintf(`
		resource "google_compute_network" "foobar" {
			name = "tunnel-test-%s"
		}
		resource "google_compute_subnetwork" "foobar" {
			name = "tunnel-test-subnetwork-%s"
			network = "${google_compute_network.foobar.self_link}"
			ip_cidr_range = "10.0.0.0/16"
			region = "us-central1"
		}
		resource "google_compute_address" "foobar" {
			name = "tunnel-test-%s"
			region = "${google_compute_subnetwork.foobar.region}"
		}
		resource "google_compute_vpn_gateway" "foobar" {
			name = "tunnel-test-%s"
			network = "${google_compute_network.foobar.self_link}"
			region = "${google_compute_subnetwork.foobar.region}"
		}
		resource "google_compute_forwarding_rule" "foobar_esp" {
			name = "tunnel-test-%s-1"
			region = "${google_compute_vpn_gateway.foobar.region}"
			ip_protocol = "ESP"
			ip_address = "${google_compute_address.foobar.address}"
			target = "${google_compute_vpn_gateway.foobar.self_link}"
		}
		resource "google_compute_forwarding_rule" "foobar_udp500" {
			name = "tunnel-test-%s-2"
			region = "${google_compute_forwarding_rule.foobar_esp.region}"
			ip_protocol = "UDP"
			port_range = "500-500"
			ip_address = "${google_compute_address.foobar.address}"
			target = "${google_compute_vpn_gateway.foobar.self_link}"
		}
		resource "google_compute_forwarding_rule" "foobar_udp4500" {
			name = "tunnel-test-%s-3"
			region = "${google_compute_forwarding_rule.foobar_udp500.region}"
			ip_protocol = "UDP"
			port_range = "4500-4500"
			ip_address = "${google_compute_address.foobar.address}"
			target = "${google_compute_vpn_gateway.foobar.self_link}"
		}
		resource "google_compute_router" "foobar"{
			name = "%s"
			region = "${google_compute_forwarding_rule.foobar_udp500.region}"
			network = "${google_compute_network.foobar.self_link}"
			bgp {
				asn = 64514
			}
		}
		resource "google_compute_vpn_tunnel" "foobar" {
			name = "tunnel-test-%s"
			region = "${google_compute_forwarding_rule.foobar_udp4500.region}"
			target_vpn_gateway = "${google_compute_vpn_gateway.foobar.self_link}"
			shared_secret = "unguessable"
			peer_ip = "8.8.8.8"
			router = "${google_compute_router.foobar.self_link}"
		}
	`, testId, testId, testId, testId, testId, testId, testId, router, testId)
}

func testAccComputeVpnTunnelDefaultTrafficSelectors() string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
	name = "tunnel-test-%s"
	auto_create_subnetworks = "true"
}
resource "google_compute_address" "foobar" {
	name = "tunnel-test-%s"
	region = "us-central1"
}
resource "google_compute_vpn_gateway" "foobar" {
	name = "tunnel-test-%s"
	network = "${google_compute_network.foobar.self_link}"
	region = "${google_compute_address.foobar.region}"
}
resource "google_compute_forwarding_rule" "foobar_esp" {
	name = "tunnel-test-%s"
	region = "${google_compute_vpn_gateway.foobar.region}"
	ip_protocol = "ESP"
	ip_address = "${google_compute_address.foobar.address}"
	target = "${google_compute_vpn_gateway.foobar.self_link}"
}
resource "google_compute_forwarding_rule" "foobar_udp500" {
	name = "tunnel-test-%s"
	region = "${google_compute_forwarding_rule.foobar_esp.region}"
	ip_protocol = "UDP"
	port_range = "500-500"
	ip_address = "${google_compute_address.foobar.address}"
	target = "${google_compute_vpn_gateway.foobar.self_link}"
}
resource "google_compute_forwarding_rule" "foobar_udp4500" {
	name = "tunnel-test-%s"
	region = "${google_compute_forwarding_rule.foobar_udp500.region}"
	ip_protocol = "UDP"
	port_range = "4500-4500"
	ip_address = "${google_compute_address.foobar.address}"
	target = "${google_compute_vpn_gateway.foobar.self_link}"
}
resource "google_compute_vpn_tunnel" "foobar" {
	name = "tunnel-test-%s"
	region = "${google_compute_forwarding_rule.foobar_udp4500.region}"
	target_vpn_gateway = "${google_compute_vpn_gateway.foobar.self_link}"
	shared_secret = "unguessable"
	peer_ip = "8.8.8.8"
}`, acctest.RandString(10), acctest.RandString(10), acctest.RandString(10),
		acctest.RandString(10), acctest.RandString(10), acctest.RandString(10),
		acctest.RandString(10))
}
