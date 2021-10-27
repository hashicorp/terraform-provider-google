package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccComputeNetworkEndpointGroup_networkEndpointGroup(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeNetworkEndpointGroupDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeNetworkEndpointGroup_networkEndpointGroup(context),
			},
			{
				ResourceName:            "google_compute_network_endpoint_group.neg",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"network", "subnetwork", "zone"},
			},
		},
	})
}

func testAccComputeNetworkEndpointGroup_networkEndpointGroup(context map[string]interface{}) string {
	return Nprintf(`
resource "google_compute_network_endpoint_group" "neg" {
  name         = "tf-test-my-lb-neg%{random_suffix}"
  network      = google_compute_network.default.id
  default_port = "90"
  zone         = "us-central1-a"
}

resource "google_compute_network" "default" {
  name                    = "tf-test-neg-network%{random_suffix}"
  auto_create_subnetworks = true
}
`, context)
}
