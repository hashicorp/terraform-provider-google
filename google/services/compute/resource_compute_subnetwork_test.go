// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	tpgcompute "github.com/hashicorp/terraform-provider-google/google/services/compute"

	"google.golang.org/api/compute/v1"
)

// Unit tests

func TestIsShrinkageIpCidr(t *testing.T) {
	cases := map[string]struct {
		Old, New  string
		Shrinkage bool
	}{
		"Expansion same network ip": {
			Old:       "10.0.0.0/24",
			New:       "10.0.0.0/16",
			Shrinkage: false,
		},
		"Expansion different network ip": {
			Old:       "10.0.1.0/24",
			New:       "10.0.0.0/16",
			Shrinkage: false,
		},
		"Shrinkage same network ip": {
			Old:       "10.0.0.0/16",
			New:       "10.0.0.0/24",
			Shrinkage: true,
		},
		"Shrinkage different network ip": {
			Old:       "10.0.0.0/16",
			New:       "10.1.0.0/16",
			Shrinkage: true,
		},
	}

	for tn, tc := range cases {
		if tpgcompute.IsShrinkageIpCidr(context.Background(), tc.Old, tc.New, nil) != tc.Shrinkage {
			t.Errorf("%s failed: Shrinkage should be %t", tn, tc.Shrinkage)
		}
	}
}

// Acceptance tests

func TestAccComputeSubnetwork_basic(t *testing.T) {
	t.Parallel()

	var subnetwork1 compute.Subnetwork
	var subnetwork2 compute.Subnetwork

	cnName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	subnetwork1Name := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	subnetwork2Name := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	subnetwork3Name := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeSubnetworkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSubnetwork_basic(cnName, subnetwork1Name, subnetwork2Name, subnetwork3Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSubnetworkExists(
						t, "google_compute_subnetwork.network-ref-by-url", &subnetwork1),
					testAccCheckComputeSubnetworkExists(
						t, "google_compute_subnetwork.network-ref-by-name", &subnetwork2),
				),
			},
			{
				ResourceName:      "google_compute_subnetwork.network-ref-by-url",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_compute_subnetwork.network-with-private-google-access",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeSubnetwork_update(t *testing.T) {
	t.Parallel()

	var subnetwork compute.Subnetwork

	cnName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	subnetworkName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeSubnetworkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSubnetwork_update1(cnName, "10.2.0.0/24", subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSubnetworkExists(
						t, "google_compute_subnetwork.network-with-private-google-access", &subnetwork),
				),
			},
			{
				// Expand IP CIDR range and update private_ip_google_access
				Config: testAccComputeSubnetwork_update2(cnName, "10.2.0.0/16", subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSubnetworkExists(
						t, "google_compute_subnetwork.network-with-private-google-access", &subnetwork),
				),
			},
			{
				// Shrink IP CIDR range and update private_ip_google_access
				Config: testAccComputeSubnetwork_update2(cnName, "10.2.0.0/24", subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSubnetworkExists(
						t, "google_compute_subnetwork.network-with-private-google-access", &subnetwork),
				),
			},
			{
				// Add a secondary range and enable flow logs at once
				Config: testAccComputeSubnetwork_update3(cnName, "10.2.0.0/24", subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSubnetworkExists(
						t, "google_compute_subnetwork.network-with-private-google-access", &subnetwork),
				),
			},
			{
				ResourceName:      "google_compute_subnetwork.network-with-private-google-access",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

	if subnetwork.PrivateIpGoogleAccess {
		t.Errorf("Expected PrivateIpGoogleAccess to be false, got %v", subnetwork.PrivateIpGoogleAccess)
	}
}

func TestAccComputeSubnetwork_secondaryIpRanges(t *testing.T) {
	t.Parallel()

	var subnetwork compute.Subnetwork

	cnName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	subnetworkName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeSubnetworkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSubnetwork_secondaryIpRanges_update1(cnName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSubnetworkExists(t, "google_compute_subnetwork.network-with-private-secondary-ip-ranges", &subnetwork),
					testAccCheckComputeSubnetworkHasSecondaryIpRange(&subnetwork, "tf-test-secondary-range-update1", "192.168.10.0/24"),
				),
			},
			{
				Config: testAccComputeSubnetwork_secondaryIpRanges_update2(cnName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSubnetworkExists(t, "google_compute_subnetwork.network-with-private-secondary-ip-ranges", &subnetwork),
					testAccCheckComputeSubnetworkHasSecondaryIpRange(&subnetwork, "tf-test-secondary-range-update1", "192.168.10.0/24"),
					testAccCheckComputeSubnetworkHasSecondaryIpRange(&subnetwork, "tf-test-secondary-range-update2", "192.168.11.0/24"),
				),
			},
			{
				Config: testAccComputeSubnetwork_secondaryIpRanges_update3(cnName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSubnetworkExists(t, "google_compute_subnetwork.network-with-private-secondary-ip-ranges", &subnetwork),
					testAccCheckComputeSubnetworkHasSecondaryIpRange(&subnetwork, "tf-test-secondary-range-update1", "192.168.10.0/24"),
					testAccCheckComputeSubnetworkHasSecondaryIpRange(&subnetwork, "tf-test-secondary-range-update2", "192.168.11.0/24"),
				),
			},
			{
				Config: testAccComputeSubnetwork_secondaryIpRanges_update4(cnName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSubnetworkExists(t, "google_compute_subnetwork.network-with-private-secondary-ip-ranges", &subnetwork),
					testAccCheckComputeSubnetworkHasNotSecondaryIpRange(&subnetwork, "tf-test-secondary-range-update1", "192.168.10.0/24"),
					testAccCheckComputeSubnetworkHasNotSecondaryIpRange(&subnetwork, "tf-test-secondary-range-update2", "192.168.11.0/24"),
				),
			},
			{
				Config: testAccComputeSubnetwork_secondaryIpRanges_update1(cnName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSubnetworkExists(t, "google_compute_subnetwork.network-with-private-secondary-ip-ranges", &subnetwork),
					testAccCheckComputeSubnetworkHasSecondaryIpRange(&subnetwork, "tf-test-secondary-range-update1", "192.168.10.0/24"),
					testAccCheckComputeSubnetworkHasNotSecondaryIpRange(&subnetwork, "tf-test-secondary-range-update2", "192.168.11.0/24"),
				),
			},
		},
	})
}

func TestAccComputeSubnetwork_flowLogs(t *testing.T) {
	t.Parallel()

	var subnetwork compute.Subnetwork

	cnName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	subnetworkName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeSubnetworkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSubnetwork_flowLogs(cnName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSubnetworkExists(
						t, "google_compute_subnetwork.network-with-flow-logs", &subnetwork),
				),
			},
			{
				ResourceName:      "google_compute_subnetwork.network-with-flow-logs",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeSubnetwork_flowLogsUpdate1(cnName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSubnetworkExists(
						t, "google_compute_subnetwork.network-with-flow-logs", &subnetwork),
				),
			},
			{
				ResourceName:      "google_compute_subnetwork.network-with-flow-logs",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeSubnetwork_flowLogsUpdate2(cnName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSubnetworkExists(
						t, "google_compute_subnetwork.network-with-flow-logs", &subnetwork),
				),
			},
			{
				ResourceName:      "google_compute_subnetwork.network-with-flow-logs",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeSubnetwork_flowLogsUpdate3(cnName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSubnetworkExists(
						t, "google_compute_subnetwork.network-with-flow-logs", &subnetwork),
				),
			},
			{
				ResourceName:      "google_compute_subnetwork.network-with-flow-logs",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeSubnetwork_flowLogsDelete(cnName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSubnetworkExists(
						t, "google_compute_subnetwork.network-with-flow-logs", &subnetwork),
				),
			},
			{
				ResourceName:      "google_compute_subnetwork.network-with-flow-logs",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeSubnetwork_flowLogsMigrate(t *testing.T) {
	t.Parallel()

	var subnetwork compute.Subnetwork

	cnName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	subnetworkName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeSubnetworkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSubnetwork_flowLogsMigrate(cnName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSubnetworkExists(
						t, "google_compute_subnetwork.network-with-flow-logs", &subnetwork),
				),
			},
			{
				ResourceName:      "google_compute_subnetwork.network-with-flow-logs",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeSubnetwork_flowLogsMigrate2(cnName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSubnetworkExists(
						t, "google_compute_subnetwork.network-with-flow-logs", &subnetwork),
				),
			},
			{
				ResourceName:      "google_compute_subnetwork.network-with-flow-logs",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeSubnetwork_flowLogsMigrate3(cnName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSubnetworkExists(
						t, "google_compute_subnetwork.network-with-flow-logs", &subnetwork),
				),
			},
			{
				ResourceName:      "google_compute_subnetwork.network-with-flow-logs",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeSubnetwork_ipv6(t *testing.T) {
	t.Parallel()

	cnName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	subnetworkName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeSubnetworkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSubnetwork_ipv4(cnName, subnetworkName),
			},
			{
				ResourceName:      "google_compute_subnetwork.subnetwork",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeSubnetwork_ipv6(cnName, subnetworkName),
			},
			{
				ResourceName:      "google_compute_subnetwork.subnetwork",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckComputeSubnetworkExists(t *testing.T, n string, subnetwork *compute.Subnetwork) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acctest.GoogleProviderConfig(t)
		region := rs.Primary.Attributes["region"]
		subnet_name := rs.Primary.Attributes["name"]

		found, err := config.NewComputeClient(config.UserAgent).Subnetworks.Get(
			config.Project, region, subnet_name).Do()
		if err != nil {
			return err
		}

		if found.Name != subnet_name {
			return fmt.Errorf("Subnetwork not found")
		}

		*subnetwork = *found

		return nil
	}
}

func testAccCheckComputeSubnetworkHasSecondaryIpRange(subnetwork *compute.Subnetwork, rangeName, ipCidrRange string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, secondaryRange := range subnetwork.SecondaryIpRanges {
			if secondaryRange.RangeName == rangeName {
				if secondaryRange.IpCidrRange == ipCidrRange {
					return nil
				}
				return fmt.Errorf("Secondary range %s has the wrong ip_cidr_range. Expected %s, got %s", rangeName, ipCidrRange, secondaryRange.IpCidrRange)
			}
		}

		return fmt.Errorf("Secondary range %s not found", rangeName)
	}
}

func testAccCheckComputeSubnetworkHasNotSecondaryIpRange(subnetwork *compute.Subnetwork, rangeName, ipCidrRange string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, secondaryRange := range subnetwork.SecondaryIpRanges {
			if secondaryRange.RangeName == rangeName {
				if secondaryRange.IpCidrRange == ipCidrRange {
					return fmt.Errorf("Secondary range %s has the wrong ip_cidr_range. Expected %s, got %s", rangeName, ipCidrRange, secondaryRange.IpCidrRange)
				}
			}
		}

		return nil
	}
}

func testAccComputeSubnetwork_basic(cnName, subnetwork1Name, subnetwork2Name, subnetwork3Name string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "custom-test" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "network-ref-by-url" {
  name          = "%s"
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.custom-test.self_link
}

resource "google_compute_subnetwork" "network-ref-by-name" {
  name          = "%s"
  ip_cidr_range = "10.1.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.custom-test.name
}

resource "google_compute_subnetwork" "network-with-private-google-access" {
  name                     = "%s"
  ip_cidr_range            = "10.2.0.0/16"
  region                   = "us-central1"
  network                  = google_compute_network.custom-test.self_link
  private_ip_google_access = true
}
`, cnName, subnetwork1Name, subnetwork2Name, subnetwork3Name)
}

func testAccComputeSubnetwork_update1(cnName, cidrRange, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "custom-test" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "network-with-private-google-access" {
  name                     = "%s"
  ip_cidr_range            = "%s"
  region                   = "us-central1"
  network                  = google_compute_network.custom-test.self_link
  private_ip_google_access = true
}
`, cnName, subnetworkName, cidrRange)
}

func testAccComputeSubnetwork_update2(cnName, cidrRange, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "custom-test" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "network-with-private-google-access" {
  name          = "%s"
  ip_cidr_range = "%s"
  region        = "us-central1"
  network       = google_compute_network.custom-test.self_link
}
`, cnName, subnetworkName, cidrRange)
}

func testAccComputeSubnetwork_update3(cnName, cidrRange, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "custom-test" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "network-with-private-google-access" {
  name          = "%s"
  ip_cidr_range = "%s"
  region        = "us-central1"
  network       = google_compute_network.custom-test.self_link

  secondary_ip_range {
    range_name    = "tf-test-secondary-range-update"
    ip_cidr_range = "192.168.10.0/24"
  }
}
`, cnName, subnetworkName, cidrRange)
}

func testAccComputeSubnetwork_secondaryIpRanges_update1(cnName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "custom-test" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "network-with-private-secondary-ip-ranges" {
  name          = "%s"
  ip_cidr_range = "10.2.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.custom-test.self_link
  secondary_ip_range {
    range_name    = "tf-test-secondary-range-update1"
    ip_cidr_range = "192.168.10.0/24"
  }
}
`, cnName, subnetworkName)
}

func testAccComputeSubnetwork_secondaryIpRanges_update2(cnName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "custom-test" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "network-with-private-secondary-ip-ranges" {
  name          = "%s"
  ip_cidr_range = "10.2.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.custom-test.self_link
  secondary_ip_range {
    range_name    = "tf-test-secondary-range-update1"
    ip_cidr_range = "192.168.10.0/24"
  }
  secondary_ip_range {
    range_name    = "tf-test-secondary-range-update2"
    ip_cidr_range = "192.168.11.0/24"
  }
}
`, cnName, subnetworkName)
}

func testAccComputeSubnetwork_secondaryIpRanges_update3(cnName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "custom-test" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "network-with-private-secondary-ip-ranges" {
  name          = "%s"
  ip_cidr_range = "10.2.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.custom-test.self_link
  secondary_ip_range {
    range_name    = "tf-test-secondary-range-update2"
    ip_cidr_range = "192.168.11.0/24"
  }
  secondary_ip_range {
    range_name    = "tf-test-secondary-range-update1"
    ip_cidr_range = "192.168.10.0/24"
  }
}
`, cnName, subnetworkName)
}

func testAccComputeSubnetwork_secondaryIpRanges_update4(cnName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "custom-test" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "network-with-private-secondary-ip-ranges" {
  name               = "%s"
  ip_cidr_range      = "10.2.0.0/16"
  region             = "us-central1"
  network            = google_compute_network.custom-test.self_link
  secondary_ip_range = []
}
`, cnName, subnetworkName)
}

func testAccComputeSubnetwork_flowLogs(cnName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "custom-test" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "network-with-flow-logs" {
  name          = "%s"
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.custom-test.self_link
  log_config {
    aggregation_interval = "INTERVAL_5_SEC"
    flow_sampling        = 0.5
    metadata             = "INCLUDE_ALL_METADATA"
  }
}
`, cnName, subnetworkName)
}

func testAccComputeSubnetwork_flowLogsUpdate1(cnName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "custom-test" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "network-with-flow-logs" {
  name          = "%s"
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.custom-test.self_link
  log_config {
    aggregation_interval = "INTERVAL_30_SEC"
    flow_sampling        = 0.8
    metadata             = "EXCLUDE_ALL_METADATA"
  }
}
`, cnName, subnetworkName)
}

func testAccComputeSubnetwork_flowLogsUpdate2(cnName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "custom-test" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "network-with-flow-logs" {
  name          = "%s"
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.custom-test.self_link
  log_config {
    aggregation_interval = "INTERVAL_30_SEC"
    flow_sampling        = 0.8
    metadata             = "CUSTOM_METADATA"
    metadata_fields      = [
        "src_gke_details",
        "dest_gke_details",
    ]
    filter_expr          = "inIpRange(connection.src_ip, '10.0.0.0/8')"
  }
}
`, cnName, subnetworkName)
}

func testAccComputeSubnetwork_flowLogsUpdate3(cnName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "custom-test" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "network-with-flow-logs" {
  name          = "%s"
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.custom-test.self_link
  log_config {
    aggregation_interval = "INTERVAL_30_SEC"
    flow_sampling        = 0.8
    metadata             = "INCLUDE_ALL_METADATA"
  }
}
`, cnName, subnetworkName)
}

func testAccComputeSubnetwork_flowLogsDelete(cnName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "custom-test" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "network-with-flow-logs" {
  name          = "%s"
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.custom-test.self_link
}
`, cnName, subnetworkName)
}

func testAccComputeSubnetwork_flowLogsMigrate(cnName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "custom-test" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "network-with-flow-logs" {
  name          = "%s"
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.custom-test.self_link
  log_config {
    aggregation_interval = "INTERVAL_30_SEC"
    flow_sampling        = 0.6
    metadata             = "INCLUDE_ALL_METADATA"
  }
}
`, cnName, subnetworkName)
}

func testAccComputeSubnetwork_flowLogsMigrate2(cnName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "custom-test" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "network-with-flow-logs" {
  name          = "%s"
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.custom-test.self_link
  log_config {
    aggregation_interval = "INTERVAL_30_SEC"
    flow_sampling        = 0.7
    metadata             = "INCLUDE_ALL_METADATA"
  }
}
`, cnName, subnetworkName)
}

func testAccComputeSubnetwork_flowLogsMigrate3(cnName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "custom-test" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "network-with-flow-logs" {
  name          = "%s"
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.custom-test.self_link
  log_config {
    aggregation_interval = "INTERVAL_30_SEC"
    flow_sampling        = 0.8
    metadata             = "INCLUDE_ALL_METADATA"
  }
}
`, cnName, subnetworkName)
}

func testAccComputeSubnetwork_ipv4(cnName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "custom-test" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "subnetwork" {
  name          = "%s"
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.custom-test.self_link
}
`, cnName, subnetworkName)
}

func testAccComputeSubnetwork_ipv6(cnName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "custom-test" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "subnetwork" {
  name             = "%s"
  ip_cidr_range    = "10.0.0.0/16"
  region           = "us-central1"
  network          = google_compute_network.custom-test.self_link
  stack_type       = "IPV4_IPV6"
  ipv6_access_type = "EXTERNAL"
}
`, cnName, subnetworkName)
}
