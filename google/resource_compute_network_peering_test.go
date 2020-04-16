package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"google.golang.org/api/compute/v1"
)

func TestAccComputeNetworkPeering_basic(t *testing.T) {
	t.Parallel()
	var peering_beta compute.NetworkPeering

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
				Check: resource.ComposeTestCheckFunc(
					// network foo
					testAccCheckComputeNetworkPeeringExist(t, "google_compute_network_peering.foo", &peering_beta),
					testAccCheckComputeNetworkPeeringAutoCreateRoutes(true, &peering_beta),
					testAccCheckComputeNetworkPeeringImportCustomRoutes(false, &peering_beta),
					testAccCheckComputeNetworkPeeringExportCustomRoutes(false, &peering_beta),

					// network bar
					testAccCheckComputeNetworkPeeringExist(t, "google_compute_network_peering.bar", &peering_beta),
					testAccCheckComputeNetworkPeeringAutoCreateRoutes(true, &peering_beta),
					testAccCheckComputeNetworkPeeringImportCustomRoutes(true, &peering_beta),
					testAccCheckComputeNetworkPeeringExportCustomRoutes(true, &peering_beta),
				),
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

func testAccCheckComputeNetworkPeeringExist(t *testing.T, n string, peering *compute.NetworkPeering) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := googleProviderConfig(t)

		parts := strings.Split(rs.Primary.ID, "/")
		if len(parts) != 2 {
			return fmt.Errorf("Invalid network peering identifier: %s", rs.Primary.ID)
		}

		networkName, peeringName := parts[0], parts[1]

		network, err := config.clientCompute.Networks.Get(config.Project, networkName).Do()
		if err != nil {
			return err
		}

		found := findPeeringFromNetwork(network, peeringName)
		if found == nil {
			return fmt.Errorf("Network peering '%s' not found in network '%s'", peeringName, network.Name)
		}
		*peering = *found

		return nil
	}
}

func testAccCheckComputeNetworkPeeringAutoCreateRoutes(v bool, peering *compute.NetworkPeering) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if peering.ExchangeSubnetRoutes != v {
			return fmt.Errorf("should ExchangeSubnetRouts set to %t if AutoCreateRoutes is set to %t", v, v)
		}
		return nil
	}
}

func testAccCheckComputeNetworkPeeringImportCustomRoutes(v bool, peering *compute.NetworkPeering) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if peering.ImportCustomRoutes != v {
			return fmt.Errorf("should ImportCustomRoutes set to %t", v)
		}

		return nil
	}
}

func testAccCheckComputeNetworkPeeringExportCustomRoutes(v bool, peering *compute.NetworkPeering) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if peering.ExportCustomRoutes != v {
			return fmt.Errorf("should ExportCustomRoutes set to %t", v)
		}

		return nil
	}
}

func testAccComputeNetworkPeering_basic(primaryNetworkName, peeringName, suffix string) string {
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
  name                    = "network-test-2-%s"
  auto_create_subnetworks = false
}

resource "google_compute_network_peering" "bar" {
  network      = google_compute_network.network2.self_link
  peer_network = google_compute_network.network1.self_link
  name         = "peering-test-2-%s"
`

	s = s +
		`import_custom_routes = true
		export_custom_routes = true
		`
	s = s + `}`
	return fmt.Sprintf(s, primaryNetworkName, peeringName, suffix, suffix)
}
