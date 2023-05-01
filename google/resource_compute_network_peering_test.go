package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccComputeNetworkPeering_basic(t *testing.T) {
	t.Parallel()

	primaryNetworkName := fmt.Sprintf("tf-test-network-peering-1-%d", RandInt(t))
	peeringName := fmt.Sprintf("peering-test-1-%d", RandInt(t))
	importId := fmt.Sprintf("%s/%s/%s", acctest.GetTestProjectFromEnv(), primaryNetworkName, peeringName)

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComputeNetworkPeeringDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeNetworkPeering_basic(primaryNetworkName, peeringName, RandString(t, 10)),
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

	primaryNetworkName := fmt.Sprintf("tf-test-network-peering-1-%d", RandInt(t))
	peeringName := fmt.Sprintf("peering-test-%d", RandInt(t))
	importId := fmt.Sprintf("%s/%s/%s", acctest.GetTestProjectFromEnv(), primaryNetworkName, peeringName)

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComputeNetworkPeeringDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeNetworkPeering_subnetRoutes(primaryNetworkName, peeringName, RandString(t, 10)),
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

func TestAccComputeNetworkPeering_customRoutesUpdate(t *testing.T) {
	t.Parallel()

	primaryNetworkName := fmt.Sprintf("tf-test-network-peering-1-%d", RandInt(t))
	peeringName := fmt.Sprintf("peering-test-%d", RandInt(t))
	importId := fmt.Sprintf("%s/%s/%s", acctest.GetTestProjectFromEnv(), primaryNetworkName, peeringName)
	suffix := RandString(t, 10)

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComputeNetworkPeeringDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeNetworkPeeringDefaultCustomRoutes(primaryNetworkName, peeringName, suffix),
			},
			{
				ResourceName:      "google_compute_network_peering.bar",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     importId,
			},
			{
				Config: testAccComputeNetworkPeering_basic(primaryNetworkName, peeringName, suffix),
			},
			{
				ResourceName:      "google_compute_network_peering.bar",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     importId,
			},
			{
				Config: testAccComputeNetworkPeeringDefaultCustomRoutes(primaryNetworkName, peeringName, suffix),
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
		config := GoogleProviderConfig(t)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_compute_network_peering" {
				continue
			}

			_, err := config.NewComputeClient(config.UserAgent).Networks.Get(
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
  name                    = "tf-test-network-peering-2-%s"
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
  name                    = "tf-test-network-peering-2-%s"
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

func testAccComputeNetworkPeeringDefaultCustomRoutes(primaryNetworkName, peeringName, suffix string) string {
	s := `
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
  name                    = "tf-test-network-peering-2-%s"
  auto_create_subnetworks = false
}

resource "google_compute_network_peering" "bar" {
  network      = google_compute_network.network2.self_link
  peer_network = google_compute_network.network1.self_link
  name         = "peering-test-2-%s"
}`
	return fmt.Sprintf(s, primaryNetworkName, peeringName, suffix, suffix)
}
