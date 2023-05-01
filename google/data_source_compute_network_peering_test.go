package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceComputeNetworkPeering_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
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
	return Nprintf(`
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
