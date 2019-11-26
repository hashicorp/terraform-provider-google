package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	computeBeta "google.golang.org/api/compute/v0.beta"
)

func TestAccComputeNetworkPeering_basic(t *testing.T) {
	t.Parallel()
	var peering_beta computeBeta.NetworkPeering

	primaryNetworkName := acctest.RandomWithPrefix("network-test-1")
	peeringName := acctest.RandomWithPrefix("peering-test-1")
	importId := fmt.Sprintf("%s/%s/%s", getTestProjectFromEnv(), primaryNetworkName, peeringName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccComputeNetworkPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeNetworkPeering_basic(primaryNetworkName, peeringName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeNetworkPeeringExist("google_compute_network_peering.foo", &peering_beta),
					testAccCheckComputeNetworkPeeringAutoCreateRoutes(true, &peering_beta),
					testAccCheckComputeNetworkPeeringExist("google_compute_network_peering.bar", &peering_beta),
					testAccCheckComputeNetworkPeeringAutoCreateRoutes(true, &peering_beta),
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

func testAccComputeNetworkPeeringDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_compute_network_peering" {
			continue
		}

		_, err := config.clientComputeBeta.Networks.Get(
			config.Project, rs.Primary.ID).Do()
		if err == nil {
			return fmt.Errorf("Network peering still exists")
		}
	}

	return nil
}

func testAccCheckComputeNetworkPeeringExist(n string, peering *computeBeta.NetworkPeering) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		parts := strings.Split(rs.Primary.ID, "/")
		if len(parts) != 2 {
			return fmt.Errorf("Invalid network peering identifier: %s", rs.Primary.ID)
		}

		networkName, peeringName := parts[0], parts[1]

		network, err := config.clientComputeBeta.Networks.Get(config.Project, networkName).Do()
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

func testAccCheckComputeNetworkPeeringAutoCreateRoutes(v bool, peering *computeBeta.NetworkPeering) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if peering.ExchangeSubnetRoutes != v {
			return fmt.Errorf("should ExchangeSubnetRouts set to %t if AutoCreateRoutes is set to %t", v, v)
		}
		return nil
	}
}

func testAccComputeNetworkPeering_basic(primaryNetworkName, peeringName string) string {
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

	s = s + `}`
	return fmt.Sprintf(s, primaryNetworkName, peeringName, acctest.RandString(10), acctest.RandString(10))
}
