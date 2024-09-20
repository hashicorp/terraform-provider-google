// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func testAccCheckComputeRouterNatAddressDestroyProducer(t *testing.T) func(s *terraform.State) error {
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
				return fmt.Errorf("Error, Router %s in region %s still exists", routerName, region)
			}
		}

		return nil
	}
}

func TestAccComputeRouterNatAddress_withAddressCountDecrease(t *testing.T) {
	t.Parallel()

	testId := acctest.RandString(t, 10)
	routerName := fmt.Sprintf("tf-test-router-nat-%s", testId)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeRouterNatAddressDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRouterNatAddress_withAddressCount(routerName, "2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_compute_router_nat.foobar", "initial_nat_ips.#", "1"),
					resource.TestCheckResourceAttr("google_compute_router_nat.foobar", "nat_ips.#", "1"),
					resource.TestCheckResourceAttr("google_compute_router_nat_address.foobar", "nat_ips.#", "2"),
					resource.TestCheckResourceAttr("data.google_compute_router_nat.foo", "nat_ips.#", "2"),
				),
			},
			{
				ResourceName:      "google_compute_router_nat_address.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRouterNatAddress_withAddressCount(routerName, "3"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_compute_router_nat.foobar", "initial_nat_ips.#", "1"),
					resource.TestCheckResourceAttr("google_compute_router_nat.foobar", "nat_ips.#", "2"),
					resource.TestCheckResourceAttr("google_compute_router_nat_address.foobar", "nat_ips.#", "3"),
					resource.TestCheckResourceAttr("data.google_compute_router_nat.foo", "nat_ips.#", "3"),
				),
			},
			{
				ResourceName:      "google_compute_router_nat_address.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRouterNatAddress_withAddressCount(routerName, "2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_compute_router_nat.foobar", "initial_nat_ips.#", "1"),
					resource.TestCheckResourceAttr("google_compute_router_nat.foobar", "nat_ips.#", "3"),
					resource.TestCheckResourceAttr("google_compute_router_nat_address.foobar", "nat_ips.#", "2"),
					resource.TestCheckResourceAttr("data.google_compute_router_nat.foo", "nat_ips.#", "2"),
				),
			},
			{
				ResourceName:      "google_compute_router_nat_address.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeRouterNatAddress_withAddressRemoved(t *testing.T) {
	t.Parallel()

	testId := acctest.RandString(t, 10)
	routerName := fmt.Sprintf("tf-test-router-nat-%s", testId)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
		},
		CheckDestroy: testAccCheckComputeRouterNatAddressDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRouterNatAddressWithNatIps(routerName),
			},
			{
				ResourceName:      "google_compute_router_nat_address.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRouterNatAddressWithAddressRemoved(routerName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_compute_router_nat.foo", "nat_ips.#", "1"),
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_compute_router_nat.foo",
						"google_compute_router_nat.foobar",
						map[string]struct{}{
							"initial_nat_ips":   {},
							"initial_nat_ips.#": {},
							"initial_nat_ips.0": {},
						},
					),
				),
			},
			{
				ResourceName:      "google_compute_router_nat.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeRouterNatAddress_withAutoAllocateAndAddressRemoved(t *testing.T) {
	t.Parallel()

	testId := acctest.RandString(t, 10)
	routerName := fmt.Sprintf("tf-test-router-nat-%s", testId)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
		},
		CheckDestroy: testAccCheckComputeRouterNatAddressDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRouterNatAddressWithNatIps(routerName),
			},
			{
				ResourceName:      "google_compute_router_nat_address.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRouterNatAddressWithAutoAllocateAndAddressRemoved(routerName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_compute_router_nat.foo", "nat_ips.#", "0"),
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_compute_router_nat.foo",
						"google_compute_router_nat.foobar",
						map[string]struct{}{
							"initial_nat_ips":   {},
							"initial_nat_ips.#": {},
							"initial_nat_ips.0": {},
						},
					),
				),
			},
			{
				ResourceName:      "google_compute_router_nat.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeRouterNatAddress_withNatIpsAndDrainNatIps(t *testing.T) {
	t.Parallel()

	testId := acctest.RandString(t, 10)
	routerName := fmt.Sprintf("tf-test-router-nat-%s", testId)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeRouterNatDestroyProducer(t),
		Steps: []resource.TestStep{
			// (ERROR): Creation with drain nat IPs should fail
			{
				Config:      testAccComputeRouterNatAddressWithOneDrainOneRemovedNatIps(routerName),
				ExpectError: regexp.MustCompile("New RouterNat cannot have drain_nat_ips"),
			},
			// Create NAT with three nat IPs
			{
				Config: testAccComputeRouterNatAddressWithNatIps(routerName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_compute_router_nat.foobar", "initial_nat_ips.#", "1"),
					resource.TestCheckResourceAttr("google_compute_router_nat.foobar", "nat_ips.#", "1"),
					resource.TestCheckResourceAttr("google_compute_router_nat_address.foobar", "nat_ips.#", "3"),
					resource.TestCheckResourceAttr("data.google_compute_router_nat.foo", "nat_ips.#", "3"),
				),
			},
			{
				ResourceName:      "google_compute_router_nat_address.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// (ERROR) - Should not allow draining IPs still in natIps
			{
				Config:      testAccComputeRouterNatAddressWithInvalidDrainNatIpsStillInNatIps(routerName),
				ExpectError: regexp.MustCompile("cannot be drained if still set in nat_ips"),
			},
			// natIps #1, #2, #3--> natIp #2, drainNatIp #3
			{
				Config: testAccComputeRouterNatAddressWithOneDrainOneRemovedNatIps(routerName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_compute_router_nat.foobar", "initial_nat_ips.#", "1"),
					resource.TestCheckResourceAttr("google_compute_router_nat.foobar", "nat_ips.#", "3"),
					resource.TestCheckResourceAttr("google_compute_router_nat_address.foobar", "nat_ips.#", "1"),
					resource.TestCheckResourceAttr("data.google_compute_router_nat.foo", "nat_ips.#", "1"),
				),
			},
			{
				ResourceName:      "google_compute_router_nat_address.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// (ERROR): Should not be able to drain previously removed natIps (#1)
			{
				Config:      testAccComputeRouterNatAddressWithInvalidDrainMissingNatIp(routerName),
				ExpectError: regexp.MustCompile("was not previously set in nat_ips"),
			},
			{
				Config: testAccComputeRouterNatAddressWithAddressRemoved(routerName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("google_compute_router_nat.foobar", "initial_nat_ips"),
					resource.TestCheckResourceAttr("google_compute_router_nat.foobar", "nat_ips.#", "1"),
					resource.TestCheckResourceAttr("data.google_compute_router_nat.foo", "nat_ips.#", "1"),
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_compute_router_nat.foo",
						"google_compute_router_nat.foobar",
						map[string]struct{}{
							"initial_nat_ips":   {},
							"initial_nat_ips.#": {},
							"initial_nat_ips.0": {},
						},
					),
				),
			},
			{
				ResourceName:      "google_compute_router_nat.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeRouterNatAddress_withAddressCount(routerName, routerCount string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name = "%s-net"
}

resource "google_compute_subnetwork" "foobar" {
  name          = "%s-subnet"
  network       = google_compute_network.foobar.self_link
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-east1"
}

resource "google_compute_router" "foobar" {
  name    = "%s"
  region  = google_compute_subnetwork.foobar.region
  network = google_compute_network.foobar.self_link
}

resource "google_compute_address" "foobar" {
  count = %s
  name = "%s-address-${count.index}"
  region = google_compute_subnetwork.foobar.region

  lifecycle {
    create_before_destroy = true
  }
}

resource "google_compute_router_nat_address" "foobar" {
  nat_ips = google_compute_address.foobar.*.self_link
  router = google_compute_router.foobar.name
  router_nat = google_compute_router_nat.foobar.name
  region = google_compute_router_nat.foobar.region
}

resource "google_compute_router_nat" "foobar" {
  name                               = "%s-nat"
  router                             = google_compute_router.foobar.name
  region                             = google_compute_router.foobar.region

  nat_ip_allocate_option             = "MANUAL_ONLY"
  source_subnetwork_ip_ranges_to_nat = "LIST_OF_SUBNETWORKS"

  initial_nat_ips = [google_compute_address.foobar[0].self_link]
  
  subnetwork {
    name  = google_compute_subnetwork.foobar.self_link
    source_ip_ranges_to_nat = ["ALL_IP_RANGES"]
  }

  min_ports_per_vm = 1024

  log_config {
    enable = true
    filter = "ERRORS_ONLY"
  }
}

data "google_compute_router_nat" "foo" {
  name = google_compute_router_nat.foobar.name
  router = google_compute_router_nat.foobar.router
  region = google_compute_router.foobar.region

  depends_on = [google_compute_router_nat_address.foobar]
}
`, routerName, routerName, routerName, routerCount, routerName, routerName)
}

func testAccComputeRouterNatAddressBaseResourcesWithNatIps(routerName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name                    = "%s-net"
  auto_create_subnetworks = "false"
}

resource "google_compute_subnetwork" "foobar" {
  name          = "%s-subnet"
  network       = google_compute_network.foobar.self_link
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
}

resource "google_compute_address" "addr1" {
  name   = "%s-addr1"
  region = google_compute_subnetwork.foobar.region
}

resource "google_compute_address" "addr2" {
  name   = "%s-addr2"
  region = google_compute_subnetwork.foobar.region
}

resource "google_compute_address" "addr3" {
  name   = "%s-addr3"
  region = google_compute_subnetwork.foobar.region
}

resource "google_compute_address" "addr4" {
  name   = "%s-addr4"
  region = google_compute_subnetwork.foobar.region
}

resource "google_compute_router" "foobar" {
  name     = "%s"
  region   = google_compute_subnetwork.foobar.region
  network  = google_compute_network.foobar.self_link
}
`, routerName, routerName, routerName, routerName, routerName, routerName, routerName)
}

func testAccComputeRouterNatAddressWithNatIps(routerName string) string {
	return fmt.Sprintf(`
%s

resource "google_compute_router_nat_address" "foobar" {
  nat_ips = [
    google_compute_address.addr1.self_link,
    google_compute_address.addr2.self_link,
    google_compute_address.addr3.self_link,
  ]
  router = google_compute_router.foobar.name
  router_nat = google_compute_router_nat.foobar.name
  region = google_compute_router_nat.foobar.region
}

resource "google_compute_router_nat" "foobar" {
  name     = "%s"
  router   = google_compute_router.foobar.name
  region   = google_compute_router.foobar.region

  nat_ip_allocate_option = "MANUAL_ONLY"
  initial_nat_ips = [google_compute_address.addr4.self_link]

  source_subnetwork_ip_ranges_to_nat = "LIST_OF_SUBNETWORKS"
  subnetwork {
    name                    = google_compute_subnetwork.foobar.self_link
    source_ip_ranges_to_nat = ["ALL_IP_RANGES"]
  }
}

data "google_compute_router_nat" "foo" {
  name = google_compute_router_nat.foobar.name
  router = google_compute_router_nat.foobar.router
  region = google_compute_router.foobar.region

  depends_on = [google_compute_router_nat_address.foobar]
}
`, testAccComputeRouterNatAddressBaseResourcesWithNatIps(routerName), routerName)
}

func testAccComputeRouterNatAddressWithAddressRemoved(routerName string) string {
	return fmt.Sprintf(`
%s

resource "google_compute_router_nat" "foobar" {
  name     = "%s"
  router   = google_compute_router.foobar.name
  region   = google_compute_router.foobar.region

  nat_ip_allocate_option = "MANUAL_ONLY"
  nat_ips = [google_compute_address.addr4.self_link]

  source_subnetwork_ip_ranges_to_nat = "LIST_OF_SUBNETWORKS"
  subnetwork {
    name                    = google_compute_subnetwork.foobar.self_link
    source_ip_ranges_to_nat = ["ALL_IP_RANGES"]
  }
}

data "google_compute_router_nat" "foo" {
  name = google_compute_router_nat.foobar.name
  router = google_compute_router_nat.foobar.router
  region = google_compute_router.foobar.region
}
`, testAccComputeRouterNatAddressBaseResourcesWithNatIps(routerName), routerName)
}

func testAccComputeRouterNatAddressWithAutoAllocateAndAddressRemoved(routerName string) string {
	return fmt.Sprintf(`
%s

resource "google_compute_router_nat" "foobar" {
  name     = "%s"
  router   = google_compute_router.foobar.name
  region   = google_compute_router.foobar.region

  nat_ip_allocate_option = "AUTO_ONLY"
  nat_ips = []

  source_subnetwork_ip_ranges_to_nat = "ALL_SUBNETWORKS_ALL_IP_RANGES"
}

data "google_compute_router_nat" "foo" {
  name = google_compute_router_nat.foobar.name
  router = google_compute_router_nat.foobar.router
  region = google_compute_router.foobar.region
}
`, testAccComputeRouterNatAddressBaseResourcesWithNatIps(routerName), routerName)
}

func testAccComputeRouterNatAddressWithOneDrainOneRemovedNatIps(routerName string) string {
	return fmt.Sprintf(`
%s

resource "google_compute_router_nat_address" "foobar" {
  nat_ips = [
    google_compute_address.addr2.self_link,
  ]

  drain_nat_ips = [
    google_compute_address.addr3.self_link,
  ]
  router = google_compute_router.foobar.name
  router_nat = google_compute_router_nat.foobar.name
  region = google_compute_router_nat.foobar.region
}

resource "google_compute_router_nat" "foobar" {
  name     = "%s"
  router   = google_compute_router.foobar.name
  region   = google_compute_router.foobar.region

  source_subnetwork_ip_ranges_to_nat = "LIST_OF_SUBNETWORKS"
  subnetwork {
    name                    = google_compute_subnetwork.foobar.self_link
    source_ip_ranges_to_nat = ["ALL_IP_RANGES"]
  }

  nat_ip_allocate_option = "MANUAL_ONLY"
  initial_nat_ips = [google_compute_address.addr4.self_link]
}

data "google_compute_router_nat" "foo" {
  name = google_compute_router_nat.foobar.name
  router = google_compute_router_nat.foobar.router
  region = google_compute_router.foobar.region

  depends_on = [google_compute_router_nat_address.foobar]
}
`, testAccComputeRouterNatAddressBaseResourcesWithNatIps(routerName), routerName)
}

func testAccComputeRouterNatAddressWithInvalidDrainNatIpsStillInNatIps(routerName string) string {
	return fmt.Sprintf(`
%s

resource "google_compute_router_nat_address" "foobar" {
  nat_ips = [
    google_compute_address.addr1.self_link,
    google_compute_address.addr2.self_link,
    google_compute_address.addr3.self_link,
  ]

  drain_nat_ips = [
    google_compute_address.addr3.self_link,
  ]
  router = google_compute_router.foobar.name
  router_nat = google_compute_router_nat.foobar.name
  region = google_compute_router_nat.foobar.region
}


resource "google_compute_router_nat" "foobar" {
  name     = "%s"
  router   = google_compute_router.foobar.name
  region   = google_compute_router.foobar.region

  source_subnetwork_ip_ranges_to_nat = "LIST_OF_SUBNETWORKS"
  subnetwork {
    name                    = google_compute_subnetwork.foobar.self_link
    source_ip_ranges_to_nat = ["ALL_IP_RANGES"]
  }

  nat_ip_allocate_option = "MANUAL_ONLY"
  initial_nat_ips = [google_compute_address.addr4.self_link]
}

data "google_compute_router_nat" "foo" {
  name = google_compute_router_nat.foobar.name
  router = google_compute_router_nat.foobar.router
  region = google_compute_router.foobar.region

  depends_on = [google_compute_router_nat_address.foobar]
}
`, testAccComputeRouterNatAddressBaseResourcesWithNatIps(routerName), routerName)
}

func testAccComputeRouterNatAddressWithInvalidDrainMissingNatIp(routerName string) string {
	return fmt.Sprintf(`
%s

resource "google_compute_router_nat_address" "foobar" {
  nat_ips = [
    google_compute_address.addr2.self_link,
  ]

  drain_nat_ips = [
    google_compute_address.addr1.self_link,
    google_compute_address.addr3.self_link,
  ]
  router = google_compute_router.foobar.name
  router_nat = google_compute_router_nat.foobar.name
  region = google_compute_router_nat.foobar.region
}

resource "google_compute_router_nat" "foobar" {
  name     = "%s"
  router   = google_compute_router.foobar.name
  region   = google_compute_router.foobar.region

  source_subnetwork_ip_ranges_to_nat = "LIST_OF_SUBNETWORKS"
  subnetwork {
    name                    = google_compute_subnetwork.foobar.self_link
    source_ip_ranges_to_nat = ["ALL_IP_RANGES"]
  }

  nat_ip_allocate_option = "MANUAL_ONLY"
  initial_nat_ips = [google_compute_address.addr4.self_link]
}

data "google_compute_router_nat" "foo" {
  name = google_compute_router_nat.foobar.name
  router = google_compute_router_nat.foobar.router
  region = google_compute_router.foobar.region

  depends_on = [google_compute_router_nat_address.foobar]
}
`, testAccComputeRouterNatAddressBaseResourcesWithNatIps(routerName), routerName)
}
