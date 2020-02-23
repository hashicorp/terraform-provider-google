package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccComputeRouterNat_basic(t *testing.T) {
	t.Parallel()

	project := getTestProjectFromEnv()
	region := getTestRegionFromEnv()

	testId := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeRouterNatDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRouterNatBasic(testId),
			},
			{
				ResourceName: "google_compute_router_nat.foobar",
				// implicitly: ImportStateId:     fmt.Sprintf("%s/%s/router-nat-test-%s/router-nat-test-%s", project, region, testId, testId),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_compute_router_nat.foobar",
				ImportStateId:     fmt.Sprintf("%s/%s/router-nat-test-%s/router-nat-test-%s", project, region, testId, testId),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_compute_router_nat.foobar",
				ImportStateId:     fmt.Sprintf("%s/router-nat-test-%s/router-nat-test-%s", region, testId, testId),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_compute_router_nat.foobar",
				ImportStateId:     fmt.Sprintf("router-nat-test-%s/router-nat-test-%s", testId, testId),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRouterNatKeepRouter(testId),
				Check: testAccCheckComputeRouterNatDelete(
					"google_compute_router_nat.foobar"),
			},
		},
	})
}

func TestAccComputeRouterNat_update(t *testing.T) {
	t.Parallel()

	testId := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeRouterNatDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRouterNatBasicBeforeUpdate(testId),
			},
			{
				ResourceName:      "google_compute_router_nat.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRouterNatUpdated(testId),
			},
			{
				ResourceName:      "google_compute_router_nat.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRouterNatBasicBeforeUpdate(testId),
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

	testId := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeRouterNatDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRouterNatWithManualIpAndSubnetConfiguration(testId),
			},
			{
				ResourceName:      "google_compute_router_nat.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckComputeRouterNatDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

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
			return fmt.Errorf("Error, Router %s in region %s still exists", routerName, region)
		}
	}

	return nil
}

func testAccCheckComputeRouterNatDelete(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		routersService := config.clientComputeBeta.Routers

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

func testAccComputeRouterNatBasic(testId string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name = "router-nat-test-%s"
}

resource "google_compute_subnetwork" "foobar" {
  name          = "router-nat-test-subnetwork-%s"
  network       = google_compute_network.foobar.self_link
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
}

resource "google_compute_router" "foobar" {
  name    = "router-nat-test-%s"
  region  = google_compute_subnetwork.foobar.region
  network = google_compute_network.foobar.self_link
}

resource "google_compute_router_nat" "foobar" {
  name                               = "router-nat-test-%s"
  router                             = google_compute_router.foobar.name
  region                             = google_compute_router.foobar.region
  nat_ip_allocate_option             = "AUTO_ONLY"
  source_subnetwork_ip_ranges_to_nat = "ALL_SUBNETWORKS_ALL_IP_RANGES"
  log_config {
    enable = true
    filter = "ERRORS_ONLY"
  }
}
`, testId, testId, testId, testId)
}

// Like basic but with extra resources
func testAccComputeRouterNatBasicBeforeUpdate(randPrefix string) string {
	return fmt.Sprintf(`
resource "google_compute_router" "foobar" {
  name    = "router-nat-test-%s"
  region  = google_compute_subnetwork.foobar.region
  network = google_compute_network.foobar.self_link
}

resource "google_compute_network" "foobar" {
  name = "router-nat-test-%s"
}

resource "google_compute_subnetwork" "foobar" {
  name          = "router-nat-test-subnetwork-%s"
  network       = google_compute_network.foobar.self_link
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
}

resource "google_compute_address" "foobar" {
  name   = "router-nat-test-%s"
  region = google_compute_subnetwork.foobar.region
}

resource "google_compute_router_nat" "foobar" {
  name                               = "router-nat-test-%s"
  router                             = google_compute_router.foobar.name
  region                             = google_compute_router.foobar.region
  nat_ip_allocate_option             = "AUTO_ONLY"
  source_subnetwork_ip_ranges_to_nat = "ALL_SUBNETWORKS_ALL_IP_RANGES"

  log_config {
    enable = true
    filter = "ERRORS_ONLY"
  }
}
`, randPrefix, randPrefix, randPrefix, randPrefix, randPrefix)
}

func testAccComputeRouterNatUpdated(randPrefix string) string {
	return fmt.Sprintf(`
resource "google_compute_router" "foobar" {
  name    = "router-nat-test-%s"
  region  = google_compute_subnetwork.foobar.region
  network = google_compute_network.foobar.self_link
}

resource "google_compute_network" "foobar" {
  name = "router-nat-test-%s"
}

resource "google_compute_subnetwork" "foobar" {
  name          = "router-nat-test-subnetwork-%s"
  network       = google_compute_network.foobar.self_link
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
}

resource "google_compute_address" "foobar" {
  name   = "router-nat-test-%s"
  region = google_compute_subnetwork.foobar.region
}

resource "google_compute_router_nat" "foobar" {
  name   = "router-nat-test-%s"
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
`, randPrefix, randPrefix, randPrefix, randPrefix, randPrefix)
}

func testAccComputeRouterNatWithManualIpAndSubnetConfiguration(testId string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name                    = "router-nat-test-%s"
  auto_create_subnetworks = "false"
}

resource "google_compute_subnetwork" "foobar" {
  name          = "router-nat-test-subnetwork-%s"
  network       = google_compute_network.foobar.self_link
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
}

resource "google_compute_address" "foobar" {
  name   = "router-nat-test-%s"
  region = google_compute_subnetwork.foobar.region
}

resource "google_compute_router" "foobar" {
  name    = "router-nat-test-%s"
  region  = google_compute_subnetwork.foobar.region
  network = google_compute_network.foobar.self_link
  bgp {
    asn = 64514
  }
}

resource "google_compute_router_nat" "foobar" {
  name                               = "router-nat-test-%s"
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
`, testId, testId, testId, testId, testId)
}

func testAccComputeRouterNatKeepRouter(testId string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name                    = "router-nat-test-%s"
  auto_create_subnetworks = "false"
}

resource "google_compute_subnetwork" "foobar" {
  name          = "router-nat-test-subnetwork-%s"
  network       = google_compute_network.foobar.self_link
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
}

resource "google_compute_router" "foobar" {
  name    = "router-nat-test-%s"
  region  = google_compute_subnetwork.foobar.region
  network = google_compute_network.foobar.self_link
}
`, testId, testId, testId)
}
