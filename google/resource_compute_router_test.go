package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccComputeRouter_basic(t *testing.T) {
	t.Parallel()

	testId := randString(t, 10)
	routerName := fmt.Sprintf("tf-test-router-%s", testId)
	resourceRegion := "europe-west1"
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeRouterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRouterBasic(routerName, resourceRegion),
			},
			{
				ResourceName:      "google_compute_router.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeRouter_noRegion(t *testing.T) {
	t.Parallel()

	testId := randString(t, 10)
	routerName := fmt.Sprintf("tf-test-router-%s", testId)
	providerRegion := "us-central1"
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeRouterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRouterNoRegion(routerName, providerRegion),
			},
			{
				ResourceName:      "google_compute_router.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeRouter_full(t *testing.T) {
	t.Parallel()

	testId := randString(t, 10)
	routerName := fmt.Sprintf("tf-test-router-%s", testId)
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeRouterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRouterFull(routerName),
			},
			{
				ResourceName:      "google_compute_router.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeRouter_update(t *testing.T) {
	t.Parallel()

	testId := randString(t, 10)
	routerName := fmt.Sprintf("tf-test-router-%s", testId)
	region := getTestRegionFromEnv()
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeRouterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRouterBasic(routerName, region),
			},
			{
				ResourceName:      "google_compute_router.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRouterFull(routerName),
			},
			{
				ResourceName:      "google_compute_router.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRouterBasic(routerName, region),
			},
			{
				ResourceName:      "google_compute_router.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeRouter_updateAddRemoveBGP(t *testing.T) {
	t.Parallel()

	testId := randString(t, 10)
	routerName := fmt.Sprintf("tf-test-router-%s", testId)
	region := getTestRegionFromEnv()
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeRouterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRouterBasic(routerName, region),
			},
			{
				ResourceName:      "google_compute_router.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRouter_noBGP(routerName, region),
			},
			{
				ResourceName:      "google_compute_router.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRouterBasic(routerName, region),
			},
			{
				ResourceName:      "google_compute_router.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeRouterBasic(routerName, resourceRegion string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name                    = "%s-net"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "foobar" {
  name          = "%s-subnet"
  network       = google_compute_network.foobar.self_link
  ip_cidr_range = "10.0.0.0/16"
  region        = "%s"
}

resource "google_compute_router" "foobar" {
  name    = "%s"
  region  = google_compute_subnetwork.foobar.region
  network = google_compute_network.foobar.name
  bgp {
    asn = 4294967294
  }
}
`, routerName, routerName, resourceRegion, routerName)
}

func testAccComputeRouterNoRegion(routerName, providerRegion string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name                    = "%s-net"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "foobar" {
  name          = "%s-subnet"
  network       = google_compute_network.foobar.self_link
  ip_cidr_range = "10.0.0.0/16"
  region        = "%s"
}

resource "google_compute_router" "foobar" {
  name    = "%s"
  network = google_compute_network.foobar.name
  bgp {
    asn = 64514
  }
}
`, routerName, routerName, providerRegion, routerName)
}

func testAccComputeRouterFull(routerName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name                    = "%s-net"
  auto_create_subnetworks = false
}

resource "google_compute_router" "foobar" {
  name    = "%s"
  network = google_compute_network.foobar.name
  bgp {
    asn               = 64514
    advertise_mode    = "CUSTOM"
    advertised_groups = ["ALL_SUBNETS"]
    advertised_ip_ranges {
      range = "1.2.3.4"
    }
    advertised_ip_ranges {
      range = "6.7.0.0/16"
    }
  }
}
`, routerName, routerName)
}

func testAccComputeRouter_noBGP(routerName, resourceRegion string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name                    = "%s-net"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "foobar" {
  name          = "%s-subnet"
  network       = google_compute_network.foobar.self_link
  ip_cidr_range = "10.0.0.0/16"
  region        = "%s"
}

resource "google_compute_router" "foobar" {
  name    = "%s"
  region  = google_compute_subnetwork.foobar.region
  network = google_compute_network.foobar.name
}
`, routerName, routerName, resourceRegion, routerName)
}
