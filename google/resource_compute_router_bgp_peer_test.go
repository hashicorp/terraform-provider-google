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
			{
				Config: testAccComputeRouterPeerAdvertiseModeUpdate(routerName),
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

func TestAccComputeRouterPeer_enable(t *testing.T) {
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
				Config: testAccComputeRouterPeerEnable(routerName, false),
				Check: testAccCheckComputeRouterPeerExists(
					t, "google_compute_router_peer.foobar"),
			},
			{
				ResourceName:      "google_compute_router_peer.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRouterPeerEnable(routerName, true),
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

func TestAccComputeRouterPeer_bfd(t *testing.T) {
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
				Config: testAccComputeRouterPeerBfd(routerName, "DISABLED"),
				Check: testAccCheckComputeRouterPeerExists(
					t, "google_compute_router_peer.foobar"),
			},
			{
				ResourceName:      "google_compute_router_peer.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
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
		},
	})
}

func testAccCheckComputeRouterPeerDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := googleProviderConfig(t)

		routersService := config.NewComputeClient(config.userAgent).Routers

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

		routersService := config.NewComputeClient(config.userAgent).Routers

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

		routersService := config.NewComputeClient(config.userAgent).Routers
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
  auto_create_subnetworks = false
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

resource "google_compute_ha_vpn_gateway" "foobar" {
  name    = "%s-gateway"
  network = google_compute_network.foobar.self_link
  region  = google_compute_subnetwork.foobar.region
}

resource "google_compute_external_vpn_gateway" "external_gateway" {
  name            = "%s-external-gateway"
  redundancy_type = "SINGLE_IP_INTERNALLY_REDUNDANT"
  description     = "An externally managed VPN gateway"
  interface {
    id         = 0
    ip_address = "8.8.8.8"
  }
}

resource "google_compute_router" "foobar" {
  name    = "%s"
  region  = google_compute_subnetwork.foobar.region
  network = google_compute_network.foobar.self_link
  bgp {
    asn = 64514
  }
}

resource "google_compute_vpn_tunnel" "foobar" {
  name               = "%s"
  region             = google_compute_subnetwork.foobar.region
  vpn_gateway = google_compute_ha_vpn_gateway.foobar.id
  peer_external_gateway           = google_compute_external_vpn_gateway.external_gateway.id
  peer_external_gateway_interface = 0  
  shared_secret      = "unguessable"
  router             = google_compute_router.foobar.name
  vpn_gateway_interface           = 0
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
  ip_address                = "169.254.3.1"
  peer_ip_address           = "169.254.3.2"
  peer_asn                  = 65515
  advertised_route_priority = 100
  interface                 = google_compute_router_interface.foobar.name
}
`, routerName, routerName, routerName, routerName, routerName, routerName, routerName, routerName, routerName)
}

func testAccComputeRouterPeerKeepRouter(routerName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name = "%s-net"
  auto_create_subnetworks = false
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

resource "google_compute_ha_vpn_gateway" "foobar" {
  name    = "%s-gateway"
  network = google_compute_network.foobar.self_link
  region  = google_compute_subnetwork.foobar.region
}

resource "google_compute_external_vpn_gateway" "external_gateway" {
  name            = "%s-external-gateway"
  redundancy_type = "SINGLE_IP_INTERNALLY_REDUNDANT"
  description     = "An externally managed VPN gateway"
  interface {
    id         = 0
    ip_address = "8.8.8.8"
  }
}

resource "google_compute_router" "foobar" {
  name    = "%s"
  region  = google_compute_subnetwork.foobar.region
  network = google_compute_network.foobar.self_link
  bgp {
    asn = 64514
  }
}

resource "google_compute_vpn_tunnel" "foobar" {
  name               = "%s"
  region             = google_compute_subnetwork.foobar.region
  vpn_gateway = google_compute_ha_vpn_gateway.foobar.id
  peer_external_gateway           = google_compute_external_vpn_gateway.external_gateway.id
  peer_external_gateway_interface = 0  
  shared_secret      = "unguessable"
  router             = google_compute_router.foobar.name
  vpn_gateway_interface           = 0
}

resource "google_compute_router_interface" "foobar" {
  name       = "%s"
  router     = google_compute_router.foobar.name
  region     = google_compute_router.foobar.region
  ip_range   = "169.254.3.1/30"
  vpn_tunnel = google_compute_vpn_tunnel.foobar.name
}
`, routerName, routerName, routerName, routerName, routerName, routerName, routerName, routerName)
}

func testAccComputeRouterPeerAdvertiseMode(routerName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name = "%s-net"
  auto_create_subnetworks = false
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

resource "google_compute_ha_vpn_gateway" "foobar" {
  name    = "%s-gateway"
  network = google_compute_network.foobar.self_link
  region  = google_compute_subnetwork.foobar.region
}

resource "google_compute_external_vpn_gateway" "external_gateway" {
  name            = "%s-external-gateway"
  redundancy_type = "SINGLE_IP_INTERNALLY_REDUNDANT"
  description     = "An externally managed VPN gateway"
  interface {
    id         = 0
    ip_address = "8.8.8.8"
  }
}

resource "google_compute_router" "foobar" {
  name    = "%s"
  region  = google_compute_subnetwork.foobar.region
  network = google_compute_network.foobar.self_link
  bgp {
    asn = 64514
  }
}

resource "google_compute_vpn_tunnel" "foobar" {
  name               = "%s"
  region             = google_compute_subnetwork.foobar.region
  vpn_gateway = google_compute_ha_vpn_gateway.foobar.id
  peer_external_gateway           = google_compute_external_vpn_gateway.external_gateway.id
  peer_external_gateway_interface = 0
  shared_secret      = "unguessable"
  router             = google_compute_router.foobar.name
  vpn_gateway_interface           = 0  
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
  advertise_mode            = "DEFAULT"
  interface = google_compute_router_interface.foobar.name
}
`, routerName, routerName, routerName, routerName, routerName, routerName, routerName, routerName, routerName)
}

func testAccComputeRouterPeerAdvertiseModeUpdate(routerName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name = "%s-net"
  auto_create_subnetworks = false
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

resource "google_compute_ha_vpn_gateway" "foobar" {
  name    = "%s-gateway"
  network = google_compute_network.foobar.self_link
  region  = google_compute_subnetwork.foobar.region
}

resource "google_compute_external_vpn_gateway" "external_gateway" {
  name            = "%s-external-gateway"
  redundancy_type = "SINGLE_IP_INTERNALLY_REDUNDANT"
  description     = "An externally managed VPN gateway"
  interface {
    id         = 0
    ip_address = "8.8.8.8"
  }
}

resource "google_compute_router" "foobar" {
  name    = "%s"
  region  = google_compute_subnetwork.foobar.region
  network = google_compute_network.foobar.self_link
  bgp {
    asn = 64514
  }
}

resource "google_compute_vpn_tunnel" "foobar" {
  name               = "%s"
  region             = google_compute_subnetwork.foobar.region
  vpn_gateway = google_compute_ha_vpn_gateway.foobar.id
  peer_external_gateway           = google_compute_external_vpn_gateway.external_gateway.id
  peer_external_gateway_interface = 0
  shared_secret      = "unguessable"
  router             = google_compute_router.foobar.name
  vpn_gateway_interface           = 0  
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
  peer_ip_address           = "169.254.3.3"
  peer_asn                  = 65516
  advertised_route_priority = 0
  advertise_mode            = "CUSTOM"
  advertised_groups         = ["ALL_SUBNETS"]
  advertised_ip_ranges {
    range = "10.1.0.0/32"
  }
  interface = google_compute_router_interface.foobar.name
}
`, routerName, routerName, routerName, routerName, routerName, routerName, routerName, routerName, routerName)
}

func testAccComputeRouterPeerEnable(routerName string, enable bool) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name = "%s-net"
  auto_create_subnetworks = false
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

resource "google_compute_ha_vpn_gateway" "foobar" {
  name    = "%s-gateway"
  network = google_compute_network.foobar.self_link
  region  = google_compute_subnetwork.foobar.region
}

resource "google_compute_external_vpn_gateway" "external_gateway" {
  name            = "%s-external-gateway"
  redundancy_type = "SINGLE_IP_INTERNALLY_REDUNDANT"
  description     = "An externally managed VPN gateway"
  interface {
    id         = 0
    ip_address = "8.8.8.8"
  }
}

resource "google_compute_router" "foobar" {
  name    = "%s"
  region  = google_compute_subnetwork.foobar.region
  network = google_compute_network.foobar.self_link
  bgp {
    asn = 64514
  }
}

resource "google_compute_vpn_tunnel" "foobar" {
  name               = "%s"
  region             = google_compute_subnetwork.foobar.region
  vpn_gateway = google_compute_ha_vpn_gateway.foobar.id
  peer_external_gateway           = google_compute_external_vpn_gateway.external_gateway.id
  peer_external_gateway_interface = 0
  shared_secret      = "unguessable"
  router             = google_compute_router.foobar.name
  vpn_gateway_interface           = 0  
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
  interface = google_compute_router_interface.foobar.name
  enable                    = %v  
}
`, routerName, routerName, routerName, routerName, routerName, routerName, routerName, routerName, routerName, enable)
}

func testAccComputeRouterPeerBfd(routerName, bfdMode string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name = "%s-net"
  auto_create_subnetworks = false
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

resource "google_compute_ha_vpn_gateway" "foobar" {
  name    = "%s-gateway"
  network = google_compute_network.foobar.self_link
  region  = google_compute_subnetwork.foobar.region
}

resource "google_compute_external_vpn_gateway" "external_gateway" {
  name            = "%s-external-gateway"
  redundancy_type = "SINGLE_IP_INTERNALLY_REDUNDANT"
  description     = "An externally managed VPN gateway"
  interface {
    id         = 0
    ip_address = "8.8.8.8"
  }
}

resource "google_compute_router" "foobar" {
  name    = "%s"
  region  = google_compute_subnetwork.foobar.region
  network = google_compute_network.foobar.self_link
  bgp {
    asn = 64514
  }
}

resource "google_compute_vpn_tunnel" "foobar" {
  name               = "%s"
  region             = google_compute_subnetwork.foobar.region
  vpn_gateway = google_compute_ha_vpn_gateway.foobar.self_link
  peer_external_gateway           = google_compute_external_vpn_gateway.external_gateway.id
  peer_external_gateway_interface = 0  
  shared_secret      = "unguessable"
  router             = google_compute_router.foobar.name
  vpn_gateway_interface           = 0    
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
  ip_address                = "169.254.3.1"
  peer_ip_address           = "169.254.3.2"
  peer_asn                  = 65515
  advertised_route_priority = 100
  interface                 = google_compute_router_interface.foobar.name

  bfd {
    min_receive_interval        = 2000
    min_transmit_interval       = 2000
    multiplier                  = 6
    session_initialization_mode = "%s"
  }
}
`, routerName, routerName, routerName, routerName, routerName, routerName, routerName, routerName, routerName, bfdMode)
}
