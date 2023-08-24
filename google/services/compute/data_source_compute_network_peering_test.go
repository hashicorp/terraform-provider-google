// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceComputeNetworkPeering_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComputeNetworkPeeringDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceComputeNetworkPeering_basic(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_compute_network_peering.peering1_ds", "google_compute_network_peering.peering1"),
				),
			},
		},
	})
}

func testAccDataSourceComputeNetworkPeering_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network_peering" "peering1" {
  name         = "peering1-%{random_suffix}"
  network      = google_compute_network.default.self_link
  peer_network = google_compute_network.other.self_link
}

resource "google_compute_network_peering" "peering2" {
  name         = "peering2-%{random_suffix}"
  network      = google_compute_network.other.self_link
  peer_network = google_compute_network.default.self_link
}

resource "google_compute_network" "default" {
  name                    = "foobar-%{random_suffix}"
  auto_create_subnetworks = "false"
}

resource "google_compute_network" "other" {
  name                    = "other-%{random_suffix}"
  auto_create_subnetworks = "false"
}

data "google_compute_network_peering" "peering1_ds" {
  name       = google_compute_network_peering.peering1.name
  network    = google_compute_network_peering.peering1.network
}
`, context)
}
