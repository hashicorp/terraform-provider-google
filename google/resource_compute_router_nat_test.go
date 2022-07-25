package google

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccComputeRouterNat_basic(t *testing.T) {
	t.Parallel()

	project := getTestProjectFromEnv()
	region := getTestRegionFromEnv()

	testId := randString(t, 10)
	routerName := fmt.Sprintf("tf-test-router-nat-%s", testId)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeRouterNatDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRouterNatBasic(routerName),
			},
			{
				// implicitly full ImportStateId
				ResourceName:      "google_compute_router_nat.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_compute_router_nat.foobar",
				ImportStateId:     fmt.Sprintf("%s/%s/%s/%s", project, region, routerName, routerName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_compute_router_nat.foobar",
				ImportStateId:     fmt.Sprintf("%s/%s/%s", region, routerName, routerName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_compute_router_nat.foobar",
				ImportStateId:     fmt.Sprintf("%s/%s", routerName, routerName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRouterNatKeepRouter(routerName),
				Check: testAccCheckComputeRouterNatDelete(
					t, "google_compute_router_nat.foobar"),
			},
		},
	})
}

func TestAccComputeRouterNat_tcpTimeWaitTimeoutSec(t *testing.T) {
	t.Parallel()

	testId := randString(t, 10)
	routerName := fmt.Sprintf("tf-test-router-nat-%s", testId)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeRouterNatDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRouterNatTcpTimeWaitTimeoutSec(routerName, 180),
			},
			{
				// implicitly full ImportStateId
				ResourceName:      "google_compute_router_nat.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRouterNatTcpTimeWaitTimeoutSec(routerName, 150),
			},
			{
				// implicitly full ImportStateId
				ResourceName:      "google_compute_router_nat.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeRouterNat_rule(t *testing.T) {
	t.Parallel()

	testId := randString(t, 10)
	routerName := fmt.Sprintf("tf-test-router-nat-%s", testId)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeRouterNatDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRouterNatRule(routerName),
			},
			{
				// implicitly full ImportStateId
				ResourceName:      "google_compute_router_nat.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRouterNatRuleUpdate(routerName),
			},
			{
				// implicitly full ImportStateId
				ResourceName:      "google_compute_router_nat.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeRouterNat_update(t *testing.T) {
	t.Parallel()

	testId := randString(t, 10)
	routerName := fmt.Sprintf("tf-test-router-nat-%s", testId)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeRouterNatDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRouterNatBasicBeforeUpdate(routerName),
			},
			{
				ResourceName:      "google_compute_router_nat.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRouterNatUpdated(routerName),
			},
			{
				ResourceName:      "google_compute_router_nat.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRouterNatUpdateToNatIPsId(routerName),
			},
			{
				ResourceName:      "google_compute_router_nat.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRouterNatUpdateToNatIPsName(routerName),
			},
			{
				ResourceName:      "google_compute_router_nat.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRouterNatBasicBeforeUpdate(routerName),
			},
			{
				ResourceName:      "google_compute_router_nat.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeRouterNat_removeLogConfig(t *testing.T) {
	t.Parallel()

	testId := randString(t, 10)
	routerName := fmt.Sprintf("tf-test-router-nat-%s", testId)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeRouterNatDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRouterNatLogConfig(routerName),
			},
			{
				ResourceName:      "google_compute_router_nat.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRouterNatLogConfigRemoved(routerName),
			},
			{
				ResourceName:      "google_compute_router_nat.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeRouterNat_withManualIpAndSubnetConfiguration(t *testing.T) {
	t.Parallel()

	testId := randString(t, 10)
	routerName := fmt.Sprintf("tf-test-router-nat-%s", testId)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeRouterNatDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRouterNatWithManualIpAndSubnetConfiguration(routerName),
			},
			{
				ResourceName:      "google_compute_router_nat.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeRouterNat_withPortAllocationMethods(t *testing.T) {
	t.Parallel()

	testId := randString(t, 10)
	routerName := fmt.Sprintf("tf-test-router-nat-%s", testId)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeRouterNatDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRouterNatWithAllocationMethod(routerName, true, false),
			},
			{
				ResourceName:      "google_compute_router_nat.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRouterNatWithAllocationMethod(routerName, false, false),
			},
			{
				ResourceName:      "google_compute_router_nat.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRouterNatWithAllocationMethod(routerName, true, false),
			},
			{
				ResourceName:      "google_compute_router_nat.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRouterNatWithAllocationMethod(routerName, false, true),
			},
			{
				ResourceName:      "google_compute_router_nat.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRouterNatWithAllocationMethodWithParameters(routerName, false, true, 256, 8192),
			},
			{
				ResourceName:      "google_compute_router_nat.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeRouterNat_withNatIpsAndDrainNatIps(t *testing.T) {
	t.Parallel()

	testId := randString(t, 10)
	routerName := fmt.Sprintf("tf-test-router-nat-%s", testId)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeRouterNatDestroyProducer(t),
		Steps: []resource.TestStep{
			// (ERROR): Creation with drain nat IPs should fail
			{
				Config:      testAccComputeRouterNatWithOneDrainOneRemovedNatIps(routerName),
				ExpectError: regexp.MustCompile("New RouterNat cannot have drain_nat_ips"),
			},
			// Create NAT with three nat IPs
			{
				Config: testAccComputeRouterNatWithNatIps(routerName),
			},
			{
				ResourceName:      "google_compute_router_nat.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// (ERROR) - Should not allow draining IPs still in natIps
			{
				Config:      testAccComputeRouterNatWithInvalidDrainNatIpsStillInNatIps(routerName),
				ExpectError: regexp.MustCompile("cannot be drained if still set in nat_ips"),
			},
			// natIps #1, #2, #3--> natIp #2, drainNatIp #3
			{
				Config: testAccComputeRouterNatWithOneDrainOneRemovedNatIps(routerName),
			},
			{
				ResourceName:      "google_compute_router_nat.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// (ERROR): Should not be able to drain previously removed natIps (#1)
			{
				Config:      testAccComputeRouterNatWithInvalidDrainMissingNatIp(routerName),
				ExpectError: regexp.MustCompile("was not previously set in nat_ips"),
			},
		},
	})
}

func testAccCheckComputeRouterNatDestroyProducer(t *testing.T) func(s *terraform.State) error {
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
				return fmt.Errorf("Error, Router %s in region %s still exists", routerName, region)
			}
		}

		return nil
	}
}

func testAccCheckComputeRouterNatDelete(t *testing.T, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := googleProviderConfig(t)

		routersService := config.NewComputeClient(config.userAgent).Routers

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_compute_router_nat" {
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

			nats := router.Nats
			for _, nat := range nats {
				if nat.Name == name {
					return fmt.Errorf("Nat %s still exists on router %s/%s", name, region, router.Name)
				}
			}
		}

		return nil
	}
}

func testAccComputeRouterNatBasic(routerName string) string {
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
}

resource "google_compute_router_nat" "foobar" {
  name                               = "%s"
  router                             = google_compute_router.foobar.name
  region                             = google_compute_router.foobar.region
  nat_ip_allocate_option             = "AUTO_ONLY"
  source_subnetwork_ip_ranges_to_nat = "ALL_SUBNETWORKS_ALL_IP_RANGES"
  log_config {
    enable = true
    filter = "ERRORS_ONLY"
  }
}
`, routerName, routerName, routerName, routerName)
}

// Like basic but with extra resources
func testAccComputeRouterNatBasicBeforeUpdate(routerName string) string {
	return fmt.Sprintf(`
resource "google_compute_router" "foobar" {
  name    = "%s"
  region  = google_compute_subnetwork.foobar.region
  network = google_compute_network.foobar.self_link
}

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

resource "google_compute_router_nat" "foobar" {
  name                               = "%s"
  router                             = google_compute_router.foobar.name
  region                             = google_compute_router.foobar.region
  nat_ip_allocate_option             = "AUTO_ONLY"
  source_subnetwork_ip_ranges_to_nat = "ALL_SUBNETWORKS_ALL_IP_RANGES"

  log_config {
    enable = true
    filter = "ERRORS_ONLY"
  }
}
`, routerName, routerName, routerName, routerName, routerName)
}

func testAccComputeRouterNatUpdated(routerName string) string {
	return fmt.Sprintf(`
resource "google_compute_router" "foobar" {
  name    = "%s"
  region  = google_compute_subnetwork.foobar.region
  network = google_compute_network.foobar.self_link
}

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

resource "google_compute_router_nat" "foobar" {
  name   = "%s"
  router = google_compute_router.foobar.name
  region = google_compute_router.foobar.region

  nat_ip_allocate_option = "MANUAL_ONLY"
  nat_ips                = [google_compute_address.foobar.self_link]

  source_subnetwork_ip_ranges_to_nat = "LIST_OF_SUBNETWORKS"

  subnetwork {
    name                    = google_compute_subnetwork.foobar.self_link
    source_ip_ranges_to_nat = ["ALL_IP_RANGES"]
  }

  udp_idle_timeout_sec             = 60
  icmp_idle_timeout_sec            = 60
  tcp_established_idle_timeout_sec = 1600
  tcp_transitory_idle_timeout_sec  = 60

  log_config {
    enable = true
    filter = "TRANSLATIONS_ONLY"
  }
}
`, routerName, routerName, routerName, routerName, routerName)
}

func testAccComputeRouterNatUpdateToNatIPsId(routerName string) string {
	return fmt.Sprintf(`
resource "google_compute_router" "foobar" {
name    = "%s"
region  = google_compute_subnetwork.foobar.region
network = google_compute_network.foobar.self_link
}

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

resource "google_compute_router_nat" "foobar" {
  name   = "%s"
  router = google_compute_router.foobar.name
  region = google_compute_router.foobar.region

  nat_ip_allocate_option = "MANUAL_ONLY"
  nat_ips                = [google_compute_address.foobar.id]

  source_subnetwork_ip_ranges_to_nat = "LIST_OF_SUBNETWORKS"

  subnetwork {
    name                    = google_compute_subnetwork.foobar.self_link
    source_ip_ranges_to_nat = ["ALL_IP_RANGES"]
  }

  udp_idle_timeout_sec             = 60
  icmp_idle_timeout_sec            = 60
  tcp_established_idle_timeout_sec = 1600
  tcp_transitory_idle_timeout_sec  = 60

  log_config {
    enable = true
    filter = "TRANSLATIONS_ONLY"
  }
}
`, routerName, routerName, routerName, routerName, routerName)
}

func testAccComputeRouterNatUpdateToNatIPsName(routerName string) string {
	return fmt.Sprintf(`
resource "google_compute_router" "foobar" {
name    = "%s"
region  = google_compute_subnetwork.foobar.region
network = google_compute_network.foobar.self_link
}

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

resource "google_compute_router_nat" "foobar" {
  name   = "%s"
  router = google_compute_router.foobar.name
  region = google_compute_router.foobar.region

  nat_ip_allocate_option = "MANUAL_ONLY"
  nat_ips                = [google_compute_address.foobar.name]

  source_subnetwork_ip_ranges_to_nat = "LIST_OF_SUBNETWORKS"

  subnetwork {
    name                    = google_compute_subnetwork.foobar.self_link
    source_ip_ranges_to_nat = ["ALL_IP_RANGES"]
  }

  udp_idle_timeout_sec             = 60
  icmp_idle_timeout_sec            = 60
  tcp_established_idle_timeout_sec = 1600
  tcp_transitory_idle_timeout_sec  = 60

  log_config {
    enable = true
    filter = "TRANSLATIONS_ONLY"
  }
}
`, routerName, routerName, routerName, routerName, routerName)
}

func testAccComputeRouterNatWithManualIpAndSubnetConfiguration(routerName string) string {
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

resource "google_compute_address" "foobar" {
  name   = "router-nat-%s-addr"
  region = google_compute_subnetwork.foobar.region
}

resource "google_compute_router" "foobar" {
  name    = "%s"
  region  = google_compute_subnetwork.foobar.region
  network = google_compute_network.foobar.self_link
  bgp {
    asn = 64514
  }
}

resource "google_compute_router_nat" "foobar" {
  name                               = "%s"
  router                             = google_compute_router.foobar.name
  region                             = google_compute_router.foobar.region
  nat_ip_allocate_option             = "MANUAL_ONLY"
  nat_ips                            = [google_compute_address.foobar.self_link]
  source_subnetwork_ip_ranges_to_nat = "LIST_OF_SUBNETWORKS"
  subnetwork {
    name                    = google_compute_subnetwork.foobar.name
    source_ip_ranges_to_nat = ["ALL_IP_RANGES"]
  }
}
`, routerName, routerName, routerName, routerName, routerName)
}

func testAccComputeRouterNatWithAllocationMethod(routerName string, enableEndpointIndependentMapping, enableDynamicPortAllocation bool) string {
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

resource "google_compute_address" "foobar" {
  name   = "router-nat-%s-addr"
  region = google_compute_subnetwork.foobar.region
}

resource "google_compute_router" "foobar" {
  name    = "%s"
  region  = google_compute_subnetwork.foobar.region
  network = google_compute_network.foobar.self_link
  bgp {
    asn = 64514
  }
}

resource "google_compute_router_nat" "foobar" {
  name                               = "%s"
  router                             = google_compute_router.foobar.name
  region                             = google_compute_router.foobar.region
  nat_ip_allocate_option             = "MANUAL_ONLY"
  nat_ips                            = [google_compute_address.foobar.self_link]
  source_subnetwork_ip_ranges_to_nat = "LIST_OF_SUBNETWORKS"
  subnetwork {
    name                    = google_compute_subnetwork.foobar.name
    source_ip_ranges_to_nat = ["ALL_IP_RANGES"]
  }
  enable_endpoint_independent_mapping = %t
  enable_dynamic_port_allocation = %t
}
`, routerName, routerName, routerName, routerName, routerName, enableEndpointIndependentMapping, enableDynamicPortAllocation)
}

func testAccComputeRouterNatWithAllocationMethodWithParameters(routerName string, enableEndpointIndependentMapping, enableDynamicPortAllocation bool, minPortsPerVm, maxPortsPerVm uint32) string {
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

resource "google_compute_address" "foobar" {
  name   = "router-nat-%s-addr"
  region = google_compute_subnetwork.foobar.region
}

resource "google_compute_router" "foobar" {
  name    = "%s"
  region  = google_compute_subnetwork.foobar.region
  network = google_compute_network.foobar.self_link
  bgp {
    asn = 64514
  }
}

resource "google_compute_router_nat" "foobar" {
  name                               = "%s"
  router                             = google_compute_router.foobar.name
  region                             = google_compute_router.foobar.region
  nat_ip_allocate_option             = "MANUAL_ONLY"
  nat_ips                            = [google_compute_address.foobar.self_link]
  source_subnetwork_ip_ranges_to_nat = "LIST_OF_SUBNETWORKS"
  subnetwork {
    name                    = google_compute_subnetwork.foobar.name
    source_ip_ranges_to_nat = ["ALL_IP_RANGES"]
  }
  enable_endpoint_independent_mapping = %t
  enable_dynamic_port_allocation = %t
  min_ports_per_vm = %d
  max_ports_per_vm = %d
}
`, routerName, routerName, routerName, routerName, routerName, enableEndpointIndependentMapping, enableDynamicPortAllocation, minPortsPerVm, maxPortsPerVm)
}

func testAccComputeRouterNatBaseResourcesWithNatIps(routerName string) string {
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

resource "google_compute_router" "foobar" {
  name     = "%s"
  region   = google_compute_subnetwork.foobar.region
  network  = google_compute_network.foobar.self_link
}
`, routerName, routerName, routerName, routerName, routerName, routerName)
}

func testAccComputeRouterNatWithNatIps(routerName string) string {
	return fmt.Sprintf(`
%s

resource "google_compute_router_nat" "foobar" {
  name     = "%s"
  router   = google_compute_router.foobar.name
  region   = google_compute_router.foobar.region

  nat_ip_allocate_option = "MANUAL_ONLY"
  nat_ips = [
    google_compute_address.addr1.self_link,
    google_compute_address.addr2.self_link,
    google_compute_address.addr3.self_link,
  ]

  source_subnetwork_ip_ranges_to_nat = "LIST_OF_SUBNETWORKS"
  subnetwork {
    name                    = google_compute_subnetwork.foobar.self_link
    source_ip_ranges_to_nat = ["ALL_IP_RANGES"]
  }
}
`, testAccComputeRouterNatBaseResourcesWithNatIps(routerName), routerName)
}

func testAccComputeRouterNatWithOneDrainOneRemovedNatIps(routerName string) string {
	return fmt.Sprintf(`
%s

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
  nat_ips = [
    google_compute_address.addr2.self_link,
  ]

  drain_nat_ips = [
    google_compute_address.addr3.self_link,
  ]
}
`, testAccComputeRouterNatBaseResourcesWithNatIps(routerName), routerName)
}

func testAccComputeRouterNatWithInvalidDrainMissingNatIp(routerName string) string {
	return fmt.Sprintf(`
%s

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
  nat_ips = [
    google_compute_address.addr2.self_link,
  ]

  drain_nat_ips = [
    google_compute_address.addr1.self_link,
    google_compute_address.addr3.self_link,
  ]
}
`, testAccComputeRouterNatBaseResourcesWithNatIps(routerName), routerName)
}

func testAccComputeRouterNatWithInvalidDrainNatIpsStillInNatIps(routerName string) string {
	return fmt.Sprintf(`
%s

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
  nat_ips = [
    google_compute_address.addr1.self_link,
    google_compute_address.addr2.self_link,
    google_compute_address.addr3.self_link,
  ]

  drain_nat_ips = [
    google_compute_address.addr3.self_link,
  ]
}
`, testAccComputeRouterNatBaseResourcesWithNatIps(routerName), routerName)
}

func testAccComputeRouterNatKeepRouter(routerName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name                    = "%s"
  auto_create_subnetworks = "false"
}

resource "google_compute_subnetwork" "foobar" {
  name          = "%s"
  network       = google_compute_network.foobar.self_link
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
}

resource "google_compute_router" "foobar" {
  name    = "%s"
  region  = google_compute_subnetwork.foobar.region
  network = google_compute_network.foobar.self_link
}
`, routerName, routerName, routerName)
}

func testAccComputeRouterNatLogConfig(routerName string) string {
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
}

resource "google_compute_router_nat" "foobar" {
  name                               = "%s"
  router                             = google_compute_router.foobar.name
  region                             = google_compute_router.foobar.region
  nat_ip_allocate_option             = "AUTO_ONLY"
  source_subnetwork_ip_ranges_to_nat = "ALL_SUBNETWORKS_ALL_IP_RANGES"
  log_config {
    enable = false
    filter = "ALL"
  }
}
`, routerName, routerName, routerName, routerName)
}

func testAccComputeRouterNatLogConfigRemoved(routerName string) string {
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
}

resource "google_compute_router_nat" "foobar" {
  name                               = "%s"
  router                             = google_compute_router.foobar.name
  region                             = google_compute_router.foobar.region
  nat_ip_allocate_option             = "AUTO_ONLY"
  source_subnetwork_ip_ranges_to_nat = "ALL_SUBNETWORKS_ALL_IP_RANGES"
}
`, routerName, routerName, routerName, routerName)
}

func testAccComputeRouterNatRule(routerName string) string {
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
}

resource "google_compute_address" "foobar" {
  name   = "%s-addr"
  region = google_compute_subnetwork.foobar.region
}

resource "google_compute_address" "foobar2" {
  name   = "%s-addr-2"
  region = google_compute_subnetwork.foobar.region
}

resource "google_compute_router_nat" "foobar" {
  name                                = "%s"
  router                              = google_compute_router.foobar.name
  region                              = google_compute_router.foobar.region
  nat_ip_allocate_option              = "MANUAL_ONLY"
  nat_ips                             = [google_compute_address.foobar.id]
  source_subnetwork_ip_ranges_to_nat  = "ALL_SUBNETWORKS_ALL_IP_RANGES"
  enable_endpoint_independent_mapping = false

  rules {
	rule_number = 1
	match       = "inIpRange(destination.ip, '1.1.0.0/16') || inIpRange(destination.ip, '2.2.0.0/16')"
	action {
	  source_nat_active_ips = [google_compute_address.foobar2.id]
	}
  }
}
`, routerName, routerName, routerName, routerName, routerName, routerName)
}

func testAccComputeRouterNatRuleUpdate(routerName string) string {
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
}

resource "google_compute_address" "foobar" {
  name   = "%s-addr"
  region = google_compute_subnetwork.foobar.region
}

resource "google_compute_address" "foobar2" {
  name   = "%s-addr-2"
  region = google_compute_subnetwork.foobar.region
}

resource "google_compute_address" "foobar3" {
  name   = "%s-addr-3"
  region = google_compute_subnetwork.foobar.region
}

resource "google_compute_router_nat" "foobar" {
  name                                = "%s"
  router                              = google_compute_router.foobar.name
  region                              = google_compute_router.foobar.region
  nat_ip_allocate_option              = "MANUAL_ONLY"
  nat_ips                             = [google_compute_address.foobar.id]
  source_subnetwork_ip_ranges_to_nat  = "ALL_SUBNETWORKS_ALL_IP_RANGES"
  enable_endpoint_independent_mapping = false

  rules {
	rule_number = 1
	match       = "inIpRange(destination.ip, '1.1.0.0/16') || inIpRange(destination.ip, '5.5.0.0/16')"
	action {
	  source_nat_active_ips = [google_compute_address.foobar2.id]
	}
  }

  rules {
	rule_number = 2
	match       = "inIpRange(destination.ip, '3.3.0.0/16') || inIpRange(destination.ip, '4.4.0.0/16')"
	action {
	  source_nat_active_ips = [google_compute_address.foobar3.id]
	}
  }
}
`, routerName, routerName, routerName, routerName, routerName, routerName, routerName)
}

func testAccComputeRouterNatTcpTimeWaitTimeoutSec(routerName string, timeout int) string {
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
}

resource "google_compute_router_nat" "foobar" {
  name                               = "%s"
  router                             = google_compute_router.foobar.name
  region                             = google_compute_router.foobar.region
  nat_ip_allocate_option             = "AUTO_ONLY"
  source_subnetwork_ip_ranges_to_nat = "ALL_SUBNETWORKS_ALL_IP_RANGES"
  tcp_time_wait_timeout_sec          = "%d"
}
`, routerName, routerName, routerName, routerName, timeout)
}
