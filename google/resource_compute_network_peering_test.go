package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccComputeNetworkPeering_basic(t *testing.T) {
	t.Parallel()

	primaryNetworkName := fmt.Sprintf("network-test-1-%d", randInt(t))
	peeringName := fmt.Sprintf("peering-test-1-%d", randInt(t))
	importId := fmt.Sprintf("%s/%s/%s", getTestProjectFromEnv(), primaryNetworkName, peeringName)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccComputeNetworkPeeringDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeNetworkPeering_basic(primaryNetworkName, peeringName, randString(t, 10)),
			},
			{
				ResourceName:      "google_compute_network_peering.foo",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     importId,
			},
		},
	})

}

func TestAccComputeNetworkPeering_subnetRoutes(t *testing.T) {
	t.Parallel()

	primaryNetworkName := fmt.Sprintf("network-test-1-%d", randInt(t))
	peeringName := fmt.Sprintf("peering-test-%d", randInt(t))
	importId := fmt.Sprintf("%s/%s/%s", getTestProjectFromEnv(), primaryNetworkName, peeringName)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccComputeNetworkPeeringDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeNetworkPeering_subnetRoutes(primaryNetworkName, peeringName, randString(t, 10)),
			},
			{
				ResourceName:      "google_compute_network_peering.bar",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     importId,
			},
		},
	})
}

func testAccComputeNetworkPeeringDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := googleProviderConfig(t)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_compute_network_peering" {
				continue
			}

			_, err := config.clientCompute.Networks.Get(
				config.Project, rs.Primary.ID).Do()
			if err == nil {
				return fmt.Errorf("Network peering still exists")
			}
		}

		return nil
	}
}

func testAccComputeNetworkPeering_basic(primaryNetworkName, peeringName, suffix string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "network1" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_network_peering" "foo" {
  name         = "%s"
  network      = google_compute_network.network1.self_link
  peer_network = google_compute_network.network2.self_link
}

resource "google_compute_network" "network2" {
  name                    = "network-test-2-%s"
  auto_create_subnetworks = false
}

resource "google_compute_network_peering" "bar" {
  network      = google_compute_network.network2.self_link
  peer_network = google_compute_network.network1.self_link
  name         = "peering-test-2-%s"
  import_custom_routes = true
  export_custom_routes = true		
}
`, primaryNetworkName, peeringName, suffix, suffix)
}

func testAccComputeNetworkPeering_subnetRoutes(primaryNetworkName, peeringName, suffix string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "network1" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_network" "network2" {
  name                    = "network-test-2-%s"
  auto_create_subnetworks = false
}

resource "google_compute_network_peering" "bar" {
  network      = google_compute_network.network1.self_link
  peer_network = google_compute_network.network2.self_link
  name         = "%s"
  import_subnet_routes_with_public_ip = true
  export_subnet_routes_with_public_ip = false
}
`, primaryNetworkName, suffix, peeringName)
}
