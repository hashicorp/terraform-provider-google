package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccComputeVpnTunnel_regionFromGateway(t *testing.T) {
	t.Parallel()
	region := "us-central1"
	if getTestRegionFromEnv() == region {
		// Make sure we choose a region that isn't the provider default
		// in order to test getting the region from the gateway and not the
		// provider.
		region = "us-west1"
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeVpnTunnelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeVpnTunnel_regionFromGateway(region),
			},
			{
				ResourceName:            "google_compute_vpn_tunnel.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdPrefix:     fmt.Sprintf("%s/%s/", getTestProjectFromEnv(), region),
				ImportStateVerifyIgnore: []string{"shared_secret", "detailed_status"},
			},
		},
	})
}

func TestAccComputeVpnTunnel_router(t *testing.T) {
	t.Parallel()

	router := fmt.Sprintf("tf-test-tunnel-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeVpnTunnelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeVpnTunnelRouter(router),
			},
			{
				ResourceName:            "google_compute_vpn_tunnel.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"shared_secret", "detailed_status"},
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
			{
				Config: testAccComputeVpnTunnelDefaultTrafficSelectors(),
			},
			{
				ResourceName:            "google_compute_vpn_tunnel.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"shared_secret", "detailed_status"},
			},
		},
	})
}

func testAccComputeVpnTunnel_regionFromGateway(region string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name                    = "tf-test-%[1]s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "foobar" {
  name          = "tf-test-%[1]s"
  network       = google_compute_network.foobar.self_link
  ip_cidr_range = "10.0.0.0/16"
  region        = "%[2]s"
}

resource "google_compute_address" "foobar" {
  name   = "tf-test-%[1]s"
  region = google_compute_subnetwork.foobar.region
}

resource "google_compute_vpn_gateway" "foobar" {
  name    = "tf-test-%[1]s"
  network = google_compute_network.foobar.self_link
  region  = google_compute_subnetwork.foobar.region
}

resource "google_compute_forwarding_rule" "foobar_esp" {
  name        = "tf-test-%[1]s-esp"
  region      = google_compute_vpn_gateway.foobar.region
  ip_protocol = "ESP"
  ip_address  = google_compute_address.foobar.address
  target      = google_compute_vpn_gateway.foobar.self_link
}

resource "google_compute_forwarding_rule" "foobar_udp500" {
  name        = "tf-test-%[1]s-udp500"
  region      = google_compute_forwarding_rule.foobar_esp.region
  ip_protocol = "UDP"
  port_range  = "500-500"
  ip_address  = google_compute_address.foobar.address
  target      = google_compute_vpn_gateway.foobar.self_link
}

resource "google_compute_forwarding_rule" "foobar_udp4500" {
  name        = "tf-test-%[1]s-udp4500"
  region      = google_compute_forwarding_rule.foobar_udp500.region
  ip_protocol = "UDP"
  port_range  = "4500-4500"
  ip_address  = google_compute_address.foobar.address
  target      = google_compute_vpn_gateway.foobar.self_link
}

resource "google_compute_vpn_tunnel" "foobar" {
  name                    = "tf-test-%[1]s"
  target_vpn_gateway      = google_compute_vpn_gateway.foobar.self_link
  shared_secret           = "unguessable"
  peer_ip                 = "8.8.8.8"
  local_traffic_selector  = [google_compute_subnetwork.foobar.ip_cidr_range]
  remote_traffic_selector = ["192.168.0.0/24", "192.168.1.0/24"]

  depends_on = [google_compute_forwarding_rule.foobar_udp4500]
}
`, acctest.RandString(10), region)
}

func testAccComputeVpnTunnelRouter(router string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name                    = "tf-test-%[1]s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "foobar" {
  name          = "tf-test-subnetwork-%[1]s"
  network       = google_compute_network.foobar.self_link
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
}

resource "google_compute_address" "foobar" {
  name   = "tf-test-%[1]s"
  region = google_compute_subnetwork.foobar.region
}

resource "google_compute_vpn_gateway" "foobar" {
  name    = "tf-test-%[1]s"
  network = google_compute_network.foobar.self_link
  region  = google_compute_subnetwork.foobar.region
}

resource "google_compute_forwarding_rule" "foobar_esp" {
  name        = "tf-test-%[1]s-1"
  region      = google_compute_vpn_gateway.foobar.region
  ip_protocol = "ESP"
  ip_address  = google_compute_address.foobar.address
  target      = google_compute_vpn_gateway.foobar.self_link
}

resource "google_compute_forwarding_rule" "foobar_udp500" {
  name        = "tf-test-%[1]s-2"
  region      = google_compute_forwarding_rule.foobar_esp.region
  ip_protocol = "UDP"
  port_range  = "500-500"
  ip_address  = google_compute_address.foobar.address
  target      = google_compute_vpn_gateway.foobar.self_link
}

resource "google_compute_forwarding_rule" "foobar_udp4500" {
  name        = "tf-test-%[1]s-3"
  region      = google_compute_forwarding_rule.foobar_udp500.region
  ip_protocol = "UDP"
  port_range  = "4500-4500"
  ip_address  = google_compute_address.foobar.address
  target      = google_compute_vpn_gateway.foobar.self_link
}

resource "google_compute_router" "foobar" {
  name    = "%[2]s"
  region  = google_compute_forwarding_rule.foobar_udp500.region
  network = google_compute_network.foobar.self_link
  bgp {
    asn = 64514
  }
}

resource "google_compute_vpn_tunnel" "foobar" {
  name               = "tf-test-%[1]s"
  region             = google_compute_forwarding_rule.foobar_udp4500.region
  target_vpn_gateway = google_compute_vpn_gateway.foobar.self_link
  shared_secret      = "unguessable"
  peer_ip            = "8.8.8.8"
  router             = google_compute_router.foobar.self_link
}
`, acctest.RandString(10), router)
}

func testAccComputeVpnTunnelDefaultTrafficSelectors() string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name                    = "tf-test-%[1]s"
  auto_create_subnetworks = "true"
}

resource "google_compute_address" "foobar" {
  name   = "tf-test-%[1]s"
  region = "us-central1"
}

resource "google_compute_vpn_gateway" "foobar" {
  name    = "tf-test-%[1]s"
  network = google_compute_network.foobar.self_link
  region  = google_compute_address.foobar.region
}

resource "google_compute_forwarding_rule" "foobar_esp" {
  name        = "tf-test-%[1]s-esp"
  region      = google_compute_vpn_gateway.foobar.region
  ip_protocol = "ESP"
  ip_address  = google_compute_address.foobar.address
  target      = google_compute_vpn_gateway.foobar.self_link
}

resource "google_compute_forwarding_rule" "foobar_udp500" {
  name        = "tf-test-%[1]s-udp500"
  region      = google_compute_forwarding_rule.foobar_esp.region
  ip_protocol = "UDP"
  port_range  = "500-500"
  ip_address  = google_compute_address.foobar.address
  target      = google_compute_vpn_gateway.foobar.self_link
}

resource "google_compute_forwarding_rule" "foobar_udp4500" {
  name        = "tf-test-%[1]s-udp4500"
  region      = google_compute_forwarding_rule.foobar_udp500.region
  ip_protocol = "UDP"
  port_range  = "4500-4500"
  ip_address  = google_compute_address.foobar.address
  target      = google_compute_vpn_gateway.foobar.self_link
}

resource "google_compute_vpn_tunnel" "foobar" {
  name               = "tf-test-%[1]s"
  region             = google_compute_forwarding_rule.foobar_udp4500.region
  target_vpn_gateway = google_compute_vpn_gateway.foobar.self_link
  shared_secret      = "unguessable"
  peer_ip            = "8.8.8.8"
}
`, acctest.RandString(10))
}
