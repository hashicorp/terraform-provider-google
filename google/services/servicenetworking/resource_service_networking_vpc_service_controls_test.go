// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package servicenetworking_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccServiceNetworkingVPCServiceControls_update(t *testing.T) {
	t.Parallel()
	suffix := acctest.RandString(t, 10)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceNetworkingVPCServiceControls_full(suffix, "true"),
			},
			{
				ResourceName:            "google_service_networking_vpc_service_controls.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"network", "project", "service"},
			},
			{
				Config: testAccServiceNetworkingVPCServiceControls_full(suffix, "false"),
			},
			{
				ResourceName:            "google_service_networking_vpc_service_controls.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"network", "project", "service"},
			},
		},
	})
}

func testAccServiceNetworkingVPCServiceControls_full(suffix, enabled string) string {
	return fmt.Sprintf(`
# Create a VPC
resource "google_compute_network" "default" {
  name = "tf-test-example-network%s"
}

# Create an IP address
resource "google_compute_global_address" "default" {
  name          = "tf-test-psa-range%s"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.default.id
}

# Create a private connection
resource "google_service_networking_connection" "default" {
  network                 = google_compute_network.default.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.default.name]
}

# Enable VPC-SC on the producer network
resource "google_service_networking_vpc_service_controls" "default" {
  network    = google_compute_network.default.name
  service    = "servicenetworking.googleapis.com"
  enabled    = %s
  depends_on = [google_service_networking_connection.default]
}
`, suffix, suffix, enabled)
}
