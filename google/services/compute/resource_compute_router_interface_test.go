// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccComputeRouterInterface_basic(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-router-%s", acctest.RandString(t, 10))
	context := map[string]interface{}{
		"name":   name,
		"region": "us-central1",
	}
	importIdFourPart := fmt.Sprintf("%s/%s/%s/%s", envvar.GetTestProjectFromEnv(), context["region"], context["name"], context["name"]) // name reused in config

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeRouterInterfaceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRouterInterfaceBasic(context),
				Check: testAccCheckComputeRouterInterfaceExists(
					t, "google_compute_router_interface.foobar"),
			},
			{
				ResourceName:      "google_compute_router_interface.foobar",
				ImportState:       true, // Will use the 3 part {{region}}/{{router}}/{{name}} import id by default as it's the id in state
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_compute_router_interface.foobar",
				ImportState:       true,
				ImportStateId:     importIdFourPart, // Make test step use 4 part {{project}}/{{region}}/{{router}}/{{name}} import id
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRouterInterfaceKeepRouter(name),
				Check: testAccCheckComputeRouterInterfaceDelete(
					t, "google_compute_router_interface.foobar"),
			},
		},
	})
}

func TestAccComputeRouterInterface_redundant(t *testing.T) {
	t.Parallel()

	routerName := fmt.Sprintf("tf-test-router-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeRouterInterfaceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRouterInterfaceRedundant(routerName),
				Check: testAccCheckComputeRouterInterfaceExists(
					t, "google_compute_router_interface.foobar_int2"),
			},
			{
				ResourceName:      "google_compute_router_interface.foobar_int2",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeRouterInterface_withTunnel(t *testing.T) {
	t.Parallel()

	routerName := fmt.Sprintf("tf-test-router-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeRouterInterfaceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRouterInterfaceWithTunnel(routerName),
				Check: testAccCheckComputeRouterInterfaceExists(
					t, "google_compute_router_interface.foobar"),
			},
			{
				ResourceName:      "google_compute_router_interface.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeRouterInterface_withPrivateIpAddress(t *testing.T) {
	t.Parallel()

	routerName := fmt.Sprintf("tf-test-router-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeRouterInterfaceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRouterInterfaceWithPrivateIpAddress(routerName),
				Check: testAccCheckComputeRouterInterfaceExists(
					t, "google_compute_router_interface.foobar"),
			},
			{
				ResourceName:      "google_compute_router_interface.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckComputeRouterInterfaceDestroyProducer(t *testing.T) func(s *terraform.State) error {
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

func testAccCheckComputeRouterInterfaceDelete(t *testing.T, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		routersService := config.NewComputeClient(config.UserAgent).Routers

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_compute_router_interface" {
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

			ifaces := router.Interfaces
			for _, iface := range ifaces {

				if iface.Name == name {
					return fmt.Errorf("Interface %s still exists on router %s/%s", name, region, router.Name)
				}
			}
		}

		return nil
	}
}

func testAccCheckComputeRouterInterfaceExists(t *testing.T, n string) resource.TestCheckFunc {
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

		for _, iface := range router.Interfaces {

			if iface.Name == name {
				return nil
			}
		}

		return fmt.Errorf("Interface %s not found for router %s", name, router.Name)
	}
}

func testAccComputeRouterInterfaceBasic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "foobar" {
  name = "%{name}-net"
}

resource "google_compute_subnetwork" "foobar" {
  name          = "%{name}-subnet"
  network       = google_compute_network.foobar.self_link
  ip_cidr_range = "10.0.0.0/16"
  region        = "%{region}"
}

resource "google_compute_address" "foobar" {
  name   = "%{name}-addr"
  region = google_compute_subnetwork.foobar.region
}

resource "google_compute_vpn_gateway" "foobar" {
  name    = "%{name}-gateway"
  network = google_compute_network.foobar.self_link
  region  = google_compute_subnetwork.foobar.region
}

resource "google_compute_forwarding_rule" "foobar_esp" {
  name        = "%{name}-fr1"
  region      = google_compute_vpn_gateway.foobar.region
  ip_protocol = "ESP"
  ip_address  = google_compute_address.foobar.address
  target      = google_compute_vpn_gateway.foobar.self_link
}

resource "google_compute_forwarding_rule" "foobar_udp500" {
  name        = "%{name}-fr2"
  region      = google_compute_forwarding_rule.foobar_esp.region
  ip_protocol = "UDP"
  port_range  = "500-500"
  ip_address  = google_compute_address.foobar.address
  target      = google_compute_vpn_gateway.foobar.self_link
}

resource "google_compute_forwarding_rule" "foobar_udp4500" {
  name        = "%{name}-fr3"
  region      = google_compute_forwarding_rule.foobar_udp500.region
  ip_protocol = "UDP"
  port_range  = "4500-4500"
  ip_address  = google_compute_address.foobar.address
  target      = google_compute_vpn_gateway.foobar.self_link
}

resource "google_compute_router" "foobar" {
  name    = "%{name}"
  region  = google_compute_forwarding_rule.foobar_udp500.region
  network = google_compute_network.foobar.self_link
  bgp {
    asn = 64514
  }
}

resource "google_compute_router_interface" "foobar" {
  name     = "%{name}"
  router   = google_compute_router.foobar.name
  region   = google_compute_router.foobar.region
  ip_range = "169.254.3.1/30"
}
`, context)
}

func testAccComputeRouterInterfaceRedundant(routerName string) string {
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

resource "google_compute_router" "foobar" {
  name    = "%s"
  region  = google_compute_subnetwork.foobar.region
  network = google_compute_network.foobar.self_link
  bgp {
    asn = 64514
  }
}

resource "google_compute_router_interface" "foobar_int1" {
  name     = "%s-int1"
  router   = google_compute_router.foobar.name
  region   = google_compute_router.foobar.region
  ip_range = "169.254.3.1/30"
}

resource "google_compute_router_interface" "foobar_int2" {
  name                = "%s-int2"
  router              = google_compute_router.foobar.name
  region              = google_compute_router.foobar.region
  ip_range            = "169.254.4.1/30"
  redundant_interface = google_compute_router_interface.foobar_int1.name
}
`, routerName, routerName, routerName, routerName, routerName)
}

func testAccComputeRouterInterfaceKeepRouter(routerName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name = "tf-test-%s"
}

resource "google_compute_subnetwork" "foobar" {
  name          = "tf-test-router-interface-subnetwork-%s"
  network       = google_compute_network.foobar.self_link
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
}

resource "google_compute_address" "foobar" {
  name   = "%s"
  region = google_compute_subnetwork.foobar.region
}

resource "google_compute_vpn_gateway" "foobar" {
  name    = "%s"
  network = google_compute_network.foobar.self_link
  region  = google_compute_subnetwork.foobar.region
}

resource "google_compute_forwarding_rule" "foobar_esp" {
  name        = "%s-1"
  region      = google_compute_vpn_gateway.foobar.region
  ip_protocol = "ESP"
  ip_address  = google_compute_address.foobar.address
  target      = google_compute_vpn_gateway.foobar.self_link
}

resource "google_compute_forwarding_rule" "foobar_udp500" {
  name        = "%s-2"
  region      = google_compute_forwarding_rule.foobar_esp.region
  ip_protocol = "UDP"
  port_range  = "500-500"
  ip_address  = google_compute_address.foobar.address
  target      = google_compute_vpn_gateway.foobar.self_link
}

resource "google_compute_forwarding_rule" "foobar_udp4500" {
  name        = "%s-3"
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
`, routerName, routerName, routerName, routerName, routerName, routerName, routerName, routerName)
}

func testAccComputeRouterInterfaceWithTunnel(routerName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name = "tf-test-%s"
}

resource "google_compute_subnetwork" "foobar" {
  name          = "tf-test-router-interface-subnetwork-%s"
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
`, routerName, routerName, routerName, routerName, routerName, routerName, routerName, routerName)
}

func testAccComputeRouterInterfaceWithPrivateIpAddress(routerName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name = "tf-test-%s"
}

resource "google_compute_subnetwork" "foobar" {
  name          = "tf-test-router-interface-subnetwork-%s"
  network       = google_compute_network.foobar.self_link
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
}

resource "google_compute_address" "foobar" {
  name  			 = "%s-addr"
  region 			 = google_compute_subnetwork.foobar.region
  subnetwork   = google_compute_subnetwork.foobar.id
  address_type = "INTERNAL"
}

resource "google_compute_router" "foobar" {
  name    = "%s"
  region  = google_compute_subnetwork.foobar.region
  network = google_compute_network.foobar.self_link
  bgp {
    asn = 64514
  }
}

resource "google_compute_router_interface" "foobar" {
  name               = "%s"
  router             = google_compute_router.foobar.name
  region             = google_compute_router.foobar.region
  subnetwork         = google_compute_subnetwork.foobar.self_link
  private_ip_address = google_compute_address.foobar.address
}
`, routerName, routerName, routerName, routerName, routerName)
}
