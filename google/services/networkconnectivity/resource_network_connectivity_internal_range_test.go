// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package networkconnectivity_test

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"testing"
)

func TestAccNetworkConnectivityInternalRange_networkConnectivityInternalRangesBasicExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkConnectivityInternalRangeDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkConnectivityInternalRange_networkConnectivityInternalRangesBasicExample_full(context),
			},
			{
				ResourceName:            "google_network_connectivity_internal_range.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "network", "labels", "terraform_labels"},
			},
			{
				Config: testAccNetworkConnectivityInternalRange_networkConnectivityInternalRangesBasicExample_update(context),
			},
			{
				ResourceName:            "google_network_connectivity_internal_range.default",
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
  network = google_compute_network.default.self_link
  usage   = "FOR_VPC"
  peering = "FOR_SELF"
  target_cidr_range = ["10.0.0.0/8"]
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
  name    = "updated-internal-range%{random_suffix}"
  description = "Update internal range"
  network = google_compute_network.default.self_link
  usage   = "FOR_VPC"
  peering = "NOT_SHARED"
  target_cidr_range = ["192.168.0.0/16"]
  prefix_length = 22
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

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkConnectivityInternalRangeDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkConnectivityInternalRange_networkConnectivityInternalRangesExternalRangesExample_full(context),
			},
			{
				ResourceName:            "google_network_connectivity_internal_range.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "network", "labels", "terraform_labels"},
			},
			{
				Config: testAccNetworkConnectivityInternalRange_networkConnectivityInternalRangesExternalRangesExample_update(context),
			},
			{
				ResourceName:            "google_network_connectivity_internal_range.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "network", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccNetworkConnectivityInternalRange_networkConnectivityInternalRangesExternalRangesExample_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_network_connectivity_internal_range" "default" {
  name    = "basic%{random_suffix}"
  description = "Test internal range for resources outside the VPC"
  network = google_compute_network.default.self_link
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

func testAccNetworkConnectivityInternalRange_networkConnectivityInternalRangesExternalRangesExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_network_connectivity_internal_range" "default" {
  name    = "updated-internal-range%{random_suffix}"
  description = "Update internal range"
  network = google_compute_network.default.self_link
  usage   = "FOR_VPC"
  peering = "FOR_SELF"
  ip_cidr_range = "10.0.0.0/24"
}

resource "google_compute_network" "default" {
  name                    = "tf-test-internal-ranges%{random_suffix}"
  auto_create_subnetworks = false
}
`, context)
}
