// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package networkconnectivity_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccNetworkConnectivityInternalRange_networkConnectivityInternalRangesBasicExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	resourceName := "google_network_connectivity_internal_range.default"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkConnectivityInternalRangeDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkConnectivityInternalRange_networkConnectivityInternalRangesBasicExample_full(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "description", "Test internal range"),
					resource.TestCheckResourceAttr(
						resourceName, "target_cidr_range.0", "192.168.0.0/24"),
					resource.TestCheckResourceAttr(
						resourceName, "prefix_length", "24"),
					resource.TestCheckResourceAttr(
						resourceName, "overlaps.0", "OVERLAP_ROUTE_RANGE"),
					resource.TestCheckResourceAttr(
						resourceName, "labels.label-a", "b"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "network", "labels", "terraform_labels"},
			},
			{
				Config: testAccNetworkConnectivityInternalRange_networkConnectivityInternalRangesBasicExample_update(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "description", "Updated description"),
					resource.TestCheckResourceAttr(
						resourceName, "target_cidr_range.0", "192.168.0.0/24"),
					resource.TestCheckResourceAttr(
						resourceName, "prefix_length", "24"),
					resource.TestCheckResourceAttr(
						resourceName, "overlaps.0", "OVERLAP_ROUTE_RANGE"),
					resource.TestCheckResourceAttr(
						resourceName, "overlaps.1", "OVERLAP_EXISTING_SUBNET_RANGE"),
					resource.TestCheckResourceAttr(
						resourceName, "labels.label-b", "c"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "network", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccNetworkConnectivityInternalRange_networkConnectivityInternalRangesBasicExample_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_network_connectivity_internal_range" "default" {
  name    = "basic%{random_suffix}"
  description = "Test internal range"
  network = google_compute_network.default.name
  usage   = "FOR_VPC"
  peering = "FOR_SELF"
  target_cidr_range = ["192.168.0.0/24"]
  prefix_length = 24
  overlaps = ["OVERLAP_ROUTE_RANGE"]
  
  labels  = {
    label-a: "b"
  }
}

resource "google_compute_network" "default" {
  name                    = "tf-test-internal-ranges%{random_suffix}"
  auto_create_subnetworks = false
}
`, context)
}

func testAccNetworkConnectivityInternalRange_networkConnectivityInternalRangesBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_network_connectivity_internal_range" "default" {
  name    = "basic%{random_suffix}"
  description = "Updated description"
  network = google_compute_network.default.name
  usage   = "FOR_VPC"
  peering = "FOR_SELF"
  target_cidr_range = ["192.168.0.0/24"]
  prefix_length = 24
  overlaps = ["OVERLAP_ROUTE_RANGE", "OVERLAP_EXISTING_SUBNET_RANGE"]
  
  labels  = {
    label-b: "c"
  }
}

resource "google_compute_network" "default" {
  name                    = "tf-test-internal-ranges%{random_suffix}"
  auto_create_subnetworks = false
}
`, context)
}

func TestAccNetworkConnectivityInternalRange_networkConnectivityInternalRangesExternalRangesExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	resourceName := "google_network_connectivity_internal_range.default"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkConnectivityInternalRangeDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkConnectivityInternalRange_networkConnectivityInternalRangesExternalRangesExample_full(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "description", "Test internal range for resources outside the VPC"),
					resource.TestCheckResourceAttr(
						resourceName, "ip_cidr_range", "192.16.0.0/24"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "network", "labels", "terraform_labels", "usage"},
			},
			{
				Config: testAccNetworkConnectivityInternalRange_networkConnectivityInternalRangesExternalRangesExample_update(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "description", "Updated description"),
					resource.TestCheckResourceAttr(
						resourceName, "ip_cidr_range", "192.16.0.0/16"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "network", "labels", "terraform_labels", "usage"},
			},
		},
	})
}

func testAccNetworkConnectivityInternalRange_networkConnectivityInternalRangesExternalRangesExample_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_network_connectivity_internal_range" "default" {
  name    = "basic%{random_suffix}"
  description = "Test internal range for resources outside the VPC"
  network = google_compute_network.default.name
  usage   = "EXTERNAL_TO_VPC"
  peering = "FOR_SELF"
  ip_cidr_range = "192.16.0.0/24"
}

resource "google_compute_network" "default" {
  name                    = "tf-test-internal-ranges%{random_suffix}"
  auto_create_subnetworks = false
}
`, context)
}

func testAccNetworkConnectivityInternalRange_networkConnectivityInternalRangesExternalRangesExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_network_connectivity_internal_range" "default" {
  name    = "basic%{random_suffix}"
  description = "Updated description"
  network = google_compute_network.default.name
  usage   = "EXTERNAL_TO_VPC"
  peering = "FOR_SELF"
  ip_cidr_range = "192.16.0.0/16"
}

resource "google_compute_network" "default" {
  name                    = "tf-test-internal-ranges%{random_suffix}"
  auto_create_subnetworks = false
}
`, context)
}

func TestAccNetworkConnectivityInternalRange_networkConnectivityInternalRangesExcludeCIDRExample_full(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	resourceName := "google_network_connectivity_internal_range.default"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkConnectivityInternalRangeDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkConnectivityInternalRange_networkConnectivityInternalRangesExcludeCIDRExample_full(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "description", "Test internal range exclude CIDR"),
					resource.TestCheckResourceAttr(
						resourceName, "target_cidr_range.0", "10.4.0.0/16"),
					resource.TestCheckResourceAttr(
						resourceName, "target_cidr_range.1", "10.5.0.0/16"),
					resource.TestCheckResourceAttr(
						resourceName, "prefix_length", "24"),
					resource.TestCheckResourceAttr(
						resourceName, "exclude_cidr_ranges.#", "6"),
					resource.TestCheckResourceAttr(
						resourceName, "exclude_cidr_ranges.0", "10.5.0.0/24"),
					resource.TestCheckResourceAttr(
						resourceName, "exclude_cidr_ranges.1", "10.4.1.0/24"),
					resource.TestCheckResourceAttr(
						resourceName, "exclude_cidr_ranges.2", "10.4.0.0/24"),
					resource.TestCheckResourceAttr(
						resourceName, "exclude_cidr_ranges.3", "10.4.12.0/24"),
					resource.TestCheckResourceAttr(
						resourceName, "exclude_cidr_ranges.4", "10.4.32.0/24"),
					resource.TestCheckResourceAttr(
						resourceName, "exclude_cidr_ranges.5", "10.6.0.0/24"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "network", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccNetworkConnectivityInternalRange_networkConnectivityInternalRangesExcludeCIDRExample_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_network_connectivity_internal_range" "default" {
  name    = "basic%{random_suffix}"
  description = "Test internal range exclude CIDR"
  network = google_compute_network.default.name

  prefix_length = 24
  target_cidr_range = [
    "10.4.0.0/16",
	"10.5.0.0/16",
  ]
  exclude_cidr_ranges = [
    "10.5.0.0/24",
    "10.4.1.0/24",
    "10.4.0.0/24",
    "10.4.12.0/24",
	"10.4.32.0/24",
	"10.6.0.0/24",
  ]
  usage   = "FOR_VPC"
  peering = "FOR_SELF"
}

resource "google_compute_network" "default" {
  name                    = "tf-test-internal-ranges%{random_suffix}"
  auto_create_subnetworks = false
}
`, context)
}
