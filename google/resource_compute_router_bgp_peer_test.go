package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccComputeRouterPeer_basic(t *testing.T) {
	t.Parallel()

	routerName := fmt.Sprintf("tf-test-router-%s", randString(t, 10))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeRouterPeerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRouterPeerBasic(routerName),
				Check: testAccCheckComputeRouterPeerExists(
					t, "google_compute_router_peer.foobar"),
			},
			{
				ResourceName:      "google_compute_router_peer.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRouterPeerKeepRouter(routerName),
				Check: testAccCheckComputeRouterPeerDelete(
					t, "google_compute_router_peer.foobar"),
			},
		},
	})
}

func TestAccComputeRouterPeer_advertiseMode(t *testing.T) {
	t.Parallel()

	routerName := fmt.Sprintf("tf-test-router-%s", randString(t, 10))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeRouterPeerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRouterPeerAdvertiseMode(routerName),
				Check: testAccCheckComputeRouterPeerExists(
					t, "google_compute_router_peer.foobar"),
			},
			{
				ResourceName:      "google_compute_router_peer.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckComputeRouterPeerDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := googleProviderConfig(t)

		routersService := config.clientCompute.Routers

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_compute_router" {
				continue
			}

			project, err := getTestProject(rs.Primary, config)
			if err != nil {
				return err
			}

			region, err := getTestRegion(rs.Primary, config)
			if err != nil {
				return err
			}

			routerName := rs.Primary.Attributes["router"]

			_, err = routersService.Get(project, region, routerName).Do()

			if err == nil {
				return fmt.Errorf("Error, Router %s in region %s still exists",
					routerName, region)
			}
		}

		return nil
	}
}

func testAccCheckComputeRouterPeerDelete(t *testing.T, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := googleProviderConfig(t)

		routersService := config.clientCompute.Routers

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_compute_router_peer" {
				continue
			}

			project, err := getTestProject(rs.Primary, config)
			if err != nil {
				return err
			}

			region, err := getTestRegion(rs.Primary, config)
			if err != nil {
				return err
			}

			name := rs.Primary.Attributes["name"]
			routerName := rs.Primary.Attributes["router"]

			router, err := routersService.Get(project, region, routerName).Do()

			if err != nil {
				return fmt.Errorf("Error Reading Router %s: %s", routerName, err)
			}

			peers := router.BgpPeers
			for _, peer := range peers {

				if peer.Name == name {
					return fmt.Errorf("Peer %s still exists on router %s/%s", name, region, router.Name)
				}
			}
		}

		return nil
	}
}

func testAccCheckComputeRouterPeerExists(t *testing.T, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := googleProviderConfig(t)

		project, err := getTestProject(rs.Primary, config)
		if err != nil {
			return err
		}

		region, err := getTestRegion(rs.Primary, config)
		if err != nil {
			return err
		}

		name := rs.Primary.Attributes["name"]
		routerName := rs.Primary.Attributes["router"]

		routersService := config.clientCompute.Routers
		router, err := routersService.Get(project, region, routerName).Do()

		if err != nil {
			return fmt.Errorf("Error Reading Router %s: %s", routerName, err)
		}

		for _, peer := range router.BgpPeers {

			if peer.Name == name {
				return nil
			}
		}

		return fmt.Errorf("Peer %s not found for router %s", name, router.Name)
	}
}

func testAccComputeRouterPeerBasic(routerName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name = "%s-net"
}

resource "google_compute_subnetwork" "foobar" {
  name          = "%s-subnet"
  network       = google_compute_network.foobar.self_link
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
}

resource "google_compute_address" "foobar" {
  name   = "%s"
  region = google_compute_subnetwork.foobar.region
}

resource "google_compute_vpn_gateway" "foobar" {
  name    = "%s-gateway"
  network = google_compute_network.foobar.self_link
  region  = google_compute_subnetwork.foobar.region
}

resource "google_compute_forwarding_rule" "foobar_esp" {
  name        = "%s-frfr1"
  region      = google_compute_vpn_gateway.foobar.region
  ip_protocol = "ESP"
  ip_address  = google_compute_address.foobar.address
  target      = google_compute_vpn_gateway.foobar.self_link
}

resource "google_compute_forwarding_rule" "foobar_udp500" {
  name        = "%s-fr2"
  region      = google_compute_forwarding_rule.foobar_esp.region
  ip_protocol = "UDP"
  port_range  = "500-500"
  ip_address  = google_compute_address.foobar.address
  target      = google_compute_vpn_gateway.foobar.self_link
}

resource "google_compute_forwarding_rule" "foobar_udp4500" {
  name        = "%s-fr3"
  region      = google_compute_forwarding_rule.foobar_udp500.region
  ip_protocol = "UDP"
  port_range  = "4500-4500"
  ip_address  = google_compute_address.foobar.address
  target      = google_compute_vpn_gateway.foobar.self_link
}

resource "google_compute_router" "foobar" {
  name    = "%s"
  region  = google_compute_forwarding_rule.foobar_udp500.region
  network = google_compute_network.foobar.self_link
  bgp {
    asn = 64514
  }
}

resource "google_compute_vpn_tunnel" "foobar" {
  name               = "%s"
  region             = google_compute_forwarding_rule.foobar_udp4500.region
  target_vpn_gateway = google_compute_vpn_gateway.foobar.self_link
  shared_secret      = "unguessable"
  peer_ip            = "8.8.8.8"
  router             = google_compute_router.foobar.name
}

resource "google_compute_router_interface" "foobar" {
  name       = "%s"
  router     = google_compute_router.foobar.name
  region     = google_compute_router.foobar.region
  ip_range   = "169.254.3.1/30"
  vpn_tunnel = google_compute_vpn_tunnel.foobar.name
}

resource "google_compute_router_peer" "foobar" {
  name                      = "%s"
  router                    = google_compute_router.foobar.name
  region                    = google_compute_router.foobar.region
  peer_ip_address           = "169.254.3.2"
  peer_asn                  = 65515
  advertised_route_priority = 100
  interface                 = google_compute_router_interface.foobar.name
}
`, routerName, routerName, routerName, routerName, routerName, routerName, routerName, routerName, routerName, routerName, routerName)
}

func testAccComputeRouterPeerKeepRouter(routerName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name = "%s-net"
}

resource "google_compute_subnetwork" "foobar" {
  name          = "%s-subnet"
  network       = google_compute_network.foobar.self_link
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
}

resource "google_compute_address" "foobar" {
  name   = "%s"
  region = google_compute_subnetwork.foobar.region
}

resource "google_compute_vpn_gateway" "foobar" {
  name    = "%s-gateway"
  network = google_compute_network.foobar.self_link
  region  = google_compute_subnetwork.foobar.region
}

resource "google_compute_forwarding_rule" "foobar_esp" {
  name        = "%s-fr1"
  region      = google_compute_vpn_gateway.foobar.region
  ip_protocol = "ESP"
  ip_address  = google_compute_address.foobar.address
  target      = google_compute_vpn_gateway.foobar.self_link
}

resource "google_compute_forwarding_rule" "foobar_udp500" {
  name        = "%s-fr2"
  region      = google_compute_forwarding_rule.foobar_esp.region
  ip_protocol = "UDP"
  port_range  = "500-500"
  ip_address  = google_compute_address.foobar.address
  target      = google_compute_vpn_gateway.foobar.self_link
}

resource "google_compute_forwarding_rule" "foobar_udp4500" {
  name        = "%s-fr3"
  region      = google_compute_forwarding_rule.foobar_udp500.region
  ip_protocol = "UDP"
  port_range  = "4500-4500"
  ip_address  = google_compute_address.foobar.address
  target      = google_compute_vpn_gateway.foobar.self_link
}

resource "google_compute_router" "foobar" {
  name    = "%s"
  region  = google_compute_forwarding_rule.foobar_udp500.region
  network = google_compute_network.foobar.self_link
  bgp {
    asn = 64514
  }
}

resource "google_compute_vpn_tunnel" "foobar" {
  name               = "%s"
  region             = google_compute_forwarding_rule.foobar_udp4500.region
  target_vpn_gateway = google_compute_vpn_gateway.foobar.self_link
  shared_secret      = "unguessable"
  peer_ip            = "8.8.8.8"
  router             = google_compute_router.foobar.name
}

resource "google_compute_router_interface" "foobar" {
  name       = "%s"
  router     = google_compute_router.foobar.name
  region     = google_compute_router.foobar.region
  ip_range   = "169.254.3.1/30"
  vpn_tunnel = google_compute_vpn_tunnel.foobar.name
}
`, routerName, routerName, routerName, routerName, routerName, routerName, routerName, routerName, routerName, routerName)
}

func testAccComputeRouterPeerAdvertiseMode(routerName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name = "%s-net"
}

resource "google_compute_subnetwork" "foobar" {
  name          = "%s-subnet"
  network       = google_compute_network.foobar.self_link
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
}

resource "google_compute_address" "foobar" {
  name   = "%s-addr"
  region = google_compute_subnetwork.foobar.region
}

resource "google_compute_vpn_gateway" "foobar" {
  name    = "%s-gateway"
  network = google_compute_network.foobar.self_link
  region  = google_compute_subnetwork.foobar.region
}

resource "google_compute_forwarding_rule" "foobar_esp" {
  name        = "%s-fr1"
  region      = google_compute_vpn_gateway.foobar.region
  ip_protocol = "ESP"
  ip_address  = google_compute_address.foobar.address
  target      = google_compute_vpn_gateway.foobar.self_link
}

resource "google_compute_forwarding_rule" "foobar_udp500" {
  name        = "%s-fr2"
  region      = google_compute_forwarding_rule.foobar_esp.region
  ip_protocol = "UDP"
  port_range  = "500-500"
  ip_address  = google_compute_address.foobar.address
  target      = google_compute_vpn_gateway.foobar.self_link
}

resource "google_compute_forwarding_rule" "foobar_udp4500" {
  name        = "%s-fr3"
  region      = google_compute_forwarding_rule.foobar_udp500.region
  ip_protocol = "UDP"
  port_range  = "4500-4500"
  ip_address  = google_compute_address.foobar.address
  target      = google_compute_vpn_gateway.foobar.self_link
}

resource "google_compute_router" "foobar" {
  name    = "%s"
  region  = google_compute_forwarding_rule.foobar_udp500.region
  network = google_compute_network.foobar.self_link
  bgp {
    asn = 64514
  }
}

resource "google_compute_vpn_tunnel" "foobar" {
  name               = "%s"
  region             = google_compute_forwarding_rule.foobar_udp4500.region
  target_vpn_gateway = google_compute_vpn_gateway.foobar.self_link
  shared_secret      = "unguessable"
  peer_ip            = "8.8.8.8"
  router             = google_compute_router.foobar.name
}

resource "google_compute_router_interface" "foobar" {
  name       = "%s"
  router     = google_compute_router.foobar.name
  region     = google_compute_router.foobar.region
  ip_range   = "169.254.3.1/30"
  vpn_tunnel = google_compute_vpn_tunnel.foobar.name
}

resource "google_compute_router_peer" "foobar" {
  name                      = "%s"
  router                    = google_compute_router.foobar.name
  region                    = google_compute_router.foobar.region
  peer_ip_address           = "169.254.3.2"
  peer_asn                  = 65515
  advertised_route_priority = 100
  advertise_mode            = "CUSTOM"
  advertised_groups         = ["ALL_SUBNETS"]
  advertised_ip_ranges {
    range = "10.1.0.0/32"
  }
  interface = google_compute_router_interface.foobar.name
}
`, routerName, routerName, routerName, routerName, routerName, routerName, routerName, routerName, routerName, routerName, routerName)
}
