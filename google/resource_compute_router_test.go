package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccComputeRouter_basic(t *testing.T) {
	t.Parallel()

	testId := randString(t, 10)
	resourceRegion := "europe-west1"
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeRouterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRouterBasic(testId, resourceRegion),
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
	providerRegion := "us-central1"
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeRouterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRouterNoRegion(testId, providerRegion),
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
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeRouterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRouterFull(testId),
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
	region := getTestRegionFromEnv()
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeRouterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRouterBasic(testId, region),
			},
			{
				ResourceName:      "google_compute_router.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRouterFull(testId),
			},
			{
				ResourceName:      "google_compute_router.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRouterBasic(testId, region),
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
	region := getTestRegionFromEnv()
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeRouterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRouterBasic(testId, region),
			},
			{
				ResourceName:      "google_compute_router.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRouter_noBGP(testId, region),
			},
			{
				ResourceName:      "google_compute_router.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRouterBasic(testId, region),
			},
			{
				ResourceName:      "google_compute_router.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeRouterBasic(testId, resourceRegion string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name                    = "router-test-%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "foobar" {
  name          = "router-test-subnetwork-%s"
  network       = google_compute_network.foobar.self_link
  ip_cidr_range = "10.0.0.0/16"
  region        = "%s"
}

resource "google_compute_router" "foobar" {
  name    = "router-test-%s"
  region  = google_compute_subnetwork.foobar.region
  network = google_compute_network.foobar.name
  bgp {
    asn = 4294967294
  }
}
`, testId, testId, resourceRegion, testId)
}

func testAccComputeRouterNoRegion(testId, providerRegion string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name                    = "router-test-%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "foobar" {
  name          = "router-test-subnetwork-%s"
  network       = google_compute_network.foobar.self_link
  ip_cidr_range = "10.0.0.0/16"
  region        = "%s"
}

resource "google_compute_router" "foobar" {
  name    = "router-test-%s"
  network = google_compute_network.foobar.name
  bgp {
    asn = 64514
  }
}
`, testId, testId, providerRegion, testId)
}

func testAccComputeRouterFull(testId string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name                    = "router-test-%s"
  auto_create_subnetworks = false
}

resource "google_compute_router" "foobar" {
  name    = "router-test-%s"
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
`, testId, testId)
}

func testAccComputeRouter_noBGP(testId, resourceRegion string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name                    = "router-test-%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "foobar" {
  name          = "router-test-subnetwork-%s"
  network       = google_compute_network.foobar.self_link
  ip_cidr_range = "10.0.0.0/16"
  region        = "%s"
}

resource "google_compute_router" "foobar" {
  name    = "router-test-%s"
  region  = google_compute_subnetwork.foobar.region
  network = google_compute_network.foobar.name
}
`, testId, testId, resourceRegion, testId)
}
