// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccComputeRouterPeer_basic(t *testing.T) {
	t.Parallel()

	routerName := fmt.Sprintf("tf-test-router-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeRouterPeerDestroyProducer(t),
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

	routerName := fmt.Sprintf("tf-test-router-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeRouterPeerDestroyProducer(t),
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

	routerName := fmt.Sprintf("tf-test-router-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeRouterPeerDestroyProducer(t),
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

	routerName := fmt.Sprintf("tf-test-router-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeRouterPeerDestroyProducer(t),
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

func TestAccComputeRouterPeer_routerApplianceInstance(t *testing.T) {
	t.Parallel()

	routerName := fmt.Sprintf("tf-test-router-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeRouterPeerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRouterPeerRouterApplianceInstance(routerName),
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

func TestAccComputeRouterPeer_Ipv6Basic(t *testing.T) {
	t.Parallel()

	routerName := fmt.Sprintf("tf-test-router-%s", acctest.RandString(t, 10))
	resourceName := "google_compute_router_peer.foobar"
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeRouterPeerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRouterPeerIpv6(routerName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeRouterPeerExists(
						t, resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_ipv6", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeRouterPeer_UpdateIpv6Address(t *testing.T) {
	t.Parallel()

	routerName := fmt.Sprintf("tf-test-router-%s", acctest.RandString(t, 10))
	resourceName := "google_compute_router_peer.foobar"
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeRouterPeerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRouterPeerIpv6(routerName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeRouterPeerExists(
						t, resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_ipv6", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRouterPeerUpdateIpv6Address(routerName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeRouterPeerExists(
						t, resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_ipv6", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeRouterPeer_EnableDisableIpv6(t *testing.T) {
	t.Parallel()

	routerName := fmt.Sprintf("tf-test-router-%s", acctest.RandString(t, 10))
	resourceName := "google_compute_router_peer.foobar"
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeRouterPeerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRouterPeerNoIpv6(routerName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeRouterPeerExists(
						t, resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_ipv6", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRouterPeerIpv6(routerName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeRouterPeerExists(
						t, resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_ipv6", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRouterPeerIpv6(routerName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeRouterPeerExists(
						t, resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_ipv6", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckComputeRouterPeerDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		routersService := config.NewComputeClient(config.UserAgent).Routers

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_compute_router" {
				continue
			}

			project, err := acctest.GetTestProject(rs.Primary, config)
			if err != nil {
				return err
			}

			region, err := acctest.GetTestRegion(rs.Primary, config)
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
		config := acctest.GoogleProviderConfig(t)

		routersService := config.NewComputeClient(config.UserAgent).Routers

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_compute_router_peer" {
				continue
			}

			project, err := acctest.GetTestProject(rs.Primary, config)
			if err != nil {
				return err
			}

			region, err := acctest.GetTestRegion(rs.Primary, config)
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

		config := acctest.GoogleProviderConfig(t)

		project, err := acctest.GetTestProject(rs.Primary, config)
		if err != nil {
			return err
		}

		region, err := acctest.GetTestRegion(rs.Primary, config)
		if err != nil {
			return err
		}

		name := rs.Primary.Attributes["name"]
		routerName := rs.Primary.Attributes["router"]

		routersService := config.NewComputeClient(config.UserAgent).Routers
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
  vpn_tunnel = google_compute_vpn_tunnel.foobar.name
}

resource "google_compute_router_peer" "foobar" {
  name                      = "%s"
  router                    = google_compute_router.foobar.name
  region                    = google_compute_router.foobar.region
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

func testAccComputeRouterPeerRouterApplianceInstance(routerName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name                    = "%s-net"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "foobar" {
  name          = "%s-sub"
  network       = google_compute_network.foobar.self_link
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
}

resource "google_compute_address" "addr_intf" {
  name         = "%s-addr-intf"
  region       = google_compute_subnetwork.foobar.region
  subnetwork   = google_compute_subnetwork.foobar.id
  address_type = "INTERNAL"
}

resource "google_compute_address" "addr_intf_red" {
  name         = "%s-addr-intf-red"
  region       = google_compute_subnetwork.foobar.region
  subnetwork   = google_compute_subnetwork.foobar.id
  address_type = "INTERNAL"
}

resource "google_compute_address" "addr_peer" {
  name         = "%s-addr-peer"
  region       = google_compute_subnetwork.foobar.region
  subnetwork   = google_compute_subnetwork.foobar.id
  address_type = "INTERNAL"
}

resource "google_compute_instance" "foobar" {
  name           = "%s-vm"
  machine_type   = "e2-medium"
  zone           = "us-central1-a"
  can_ip_forward = true

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
    }
  }

  network_interface {
    network_ip = google_compute_address.addr_peer.address
    subnetwork = google_compute_subnetwork.foobar.self_link
  }
}

resource "google_network_connectivity_hub" "foobar" {
  name = "%s-hub"
}

resource "google_network_connectivity_spoke" "foobar" {
  name     = "%s-spoke"
  location = google_compute_subnetwork.foobar.region
  hub      = google_network_connectivity_hub.foobar.id

  linked_router_appliance_instances {
    instances {
      virtual_machine = google_compute_instance.foobar.self_link
      ip_address      = google_compute_address.addr_peer.address
    }
    site_to_site_data_transfer = false
  }
}

resource "google_compute_router" "foobar" {
  name    = "%s-ra"
  region  = google_compute_subnetwork.foobar.region
  network = google_compute_network.foobar.self_link
  bgp {
    asn = 64514
  }
}

resource "google_compute_router_interface" "foobar_redundant" {
  name                = "%s-intf-red"
  region              = google_compute_router.foobar.region
  router              = google_compute_router.foobar.name
  subnetwork          = google_compute_subnetwork.foobar.self_link
  private_ip_address  = google_compute_address.addr_intf_red.address
}

resource "google_compute_router_interface" "foobar" {
  name                = "%s-intf"
  region              = google_compute_router.foobar.region
  router              = google_compute_router.foobar.name
  subnetwork          = google_compute_subnetwork.foobar.self_link
  private_ip_address  = google_compute_address.addr_intf.address
  redundant_interface = google_compute_router_interface.foobar_redundant.name
}

resource "google_compute_router_peer" "foobar" {
  name                      = "%s-peer"
  router                    = google_compute_router.foobar.name
  region                    = google_compute_router.foobar.region
  peer_ip_address           = google_compute_address.addr_peer.address
  peer_asn                  = 65515
  interface                 = google_compute_router_interface.foobar.name
  router_appliance_instance = google_compute_instance.foobar.self_link
}
`, routerName, routerName, routerName, routerName, routerName, routerName, routerName, routerName, routerName, routerName, routerName, routerName)
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
  vpn_tunnel = google_compute_vpn_tunnel.foobar.name
}

resource "google_compute_router_peer" "foobar" {
  name                      = "%s"
  router                    = google_compute_router.foobar.name
  region                    = google_compute_router.foobar.region
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
  vpn_tunnel = google_compute_vpn_tunnel.foobar.name
}

resource "google_compute_router_peer" "foobar" {
  name                      = "%s"
  router                    = google_compute_router.foobar.name
  region                    = google_compute_router.foobar.region
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

func testAccComputeRouterPeerUpdateIpv6Address(routerName string, enableIpv6 bool) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name = "%s-net"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "foobar" {
  name          = "%s-subnet"
  network       = google_compute_network.foobar.id
  ip_cidr_range = "10.2.0.0/16"
  region        = "us-central1"
  stack_type = "IPV4_IPV6"
  ipv6_access_type = "EXTERNAL"
}

resource "google_compute_ha_vpn_gateway" "foobar" {
  name    = "%s-gateway"
  network = google_compute_network.foobar.id
  region  = google_compute_subnetwork.foobar.region
  stack_type = "IPV4_IPV6"
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
  network = google_compute_network.foobar.id
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

  enable_ipv6               = %v
  ipv6_nexthop_address      = "2600:2d00:0000:0002:0000:0000:0000:0002"
  peer_ipv6_nexthop_address = "2600:2d00:0:2::1"
}
`, routerName, routerName, routerName, routerName, routerName, routerName, routerName, routerName, enableIpv6)
}

func testAccComputeRouterPeerNoIpv6(routerName string, enableIpv6 bool) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name = "%s-net"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "foobar" {
  name          = "%s-subnet"
  network       = google_compute_network.foobar.id
  ip_cidr_range = "10.2.0.0/16"
  region        = "us-central1"
  stack_type = "IPV4_IPV6"
  ipv6_access_type = "EXTERNAL"
}

resource "google_compute_ha_vpn_gateway" "foobar" {
  name    = "%s-gateway"
  network = google_compute_network.foobar.id
  region  = google_compute_subnetwork.foobar.region
  stack_type = "IPV4_IPV6"
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
  network = google_compute_network.foobar.id
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

  enable_ipv6               = %v
}
`, routerName, routerName, routerName, routerName, routerName, routerName, routerName, routerName, enableIpv6)
}

func testAccComputeRouterPeerIpv6(routerName string, enableIpv6 bool) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name = "%s-net"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "foobar" {
  name          = "%s-subnet"
  network       = google_compute_network.foobar.id
  ip_cidr_range = "10.2.0.0/16"
  region        = "us-central1"
  stack_type = "IPV4_IPV6"
  ipv6_access_type = "EXTERNAL"
}

resource "google_compute_ha_vpn_gateway" "foobar" {
  name    = "%s-gateway"
  network = google_compute_network.foobar.id
  region  = google_compute_subnetwork.foobar.region
  stack_type = "IPV4_IPV6"
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
  network = google_compute_network.foobar.id
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

  enable_ipv6               = %v
  ipv6_nexthop_address      = "2600:2d00:0000:0002:0000:0000:0000:0001"
  peer_ipv6_nexthop_address = "2600:2d00:0:2::2"
}
`, routerName, routerName, routerName, routerName, routerName, routerName, routerName, routerName, enableIpv6)
}
