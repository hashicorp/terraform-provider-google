// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleComputeRouterNat_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeRouterNatDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleComputeRouterNat_basic(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_compute_router_nat.foo", "google_compute_router_nat.nat"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleComputeRouterNat_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "net" {
  name = "my-network%{random_suffix}"
}
	  
resource "google_compute_subnetwork" "subnet" {
  name          = "my-subnetwork%{random_suffix}"
  network       = google_compute_network.net.id
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
}
	  
resource "google_compute_router" "router" {
  name    = "my-router%{random_suffix}"
  region  = google_compute_subnetwork.subnet.region
  network = google_compute_network.net.id

  bgp {
   asn = 64514
  }
}
	  
resource "google_compute_router_nat" "nat" {
  name                               = "my-router-nat%{random_suffix}"
  router                             = google_compute_router.router.name
  region                             = google_compute_router.router.region
  nat_ip_allocate_option             = "AUTO_ONLY"
  source_subnetwork_ip_ranges_to_nat = "ALL_SUBNETWORKS_ALL_IP_RANGES"

  log_config {
    enable = true
    filter = "ERRORS_ONLY"
  }
}
	  
data "google_compute_router_nat" "foo" {
  name = google_compute_router_nat.nat.name
  router = google_compute_router_nat.nat.router
  region = google_compute_router.router.region
}`, context)

}
