package google

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"google.golang.org/api/compute/v1"
	"strings"
	"testing"
)

func TestAccComputeNetworkPeering_basic(t *testing.T) {
	var peering compute.NetworkPeering

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccComputeNetworkPeeringDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeNetworkPeering_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeNetworkPeeringExist("google_compute_network_peering.foo", &peering),
					testAccCheckComputeNetworkPeeringAutoCreateRoutes(true, &peering),
					testAccCheckComputeNetworkPeeringExist("google_compute_network_peering.bar", &peering),
					testAccCheckComputeNetworkPeeringAutoCreateRoutes(true, &peering),
				),
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

		_, err := config.clientCompute.Networks.Get(
			config.Project, rs.Primary.ID).Do()
		if err == nil {
			return fmt.Errorf("Network peering still exists")
		}
	}

	return nil
}

func testAccCheckComputeNetworkPeeringExist(n string, peering *compute.NetworkPeering) resource.TestCheckFunc {
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
		if peering.AutoCreateRoutes != v {
			return fmt.Errorf("should AutoCreateRoutes set to %t", v)
		}

		return nil
	}
}

var testAccComputeNetworkPeering_basic = fmt.Sprintf(`
resource "google_compute_network" "network1" {
	name = "network-test-1-%s"
	auto_create_subnetworks = false
}

resource "google_compute_network" "network2" {
	name = "network-test-2-%s"
	auto_create_subnetworks = false
}

resource "google_compute_network_peering" "foo" {
	name = "peering-test-1-%s"
	network = "${google_compute_network.network1.self_link}"
	peer_network = "${google_compute_network.network2.self_link}"
}

resource "google_compute_network_peering" "bar" {
	name = "peering-test-2-%s"
	auto_create_routes = true
	network = "${google_compute_network.network2.self_link}"
	peer_network = "${google_compute_network.network1.self_link}"
}
`, acctest.RandString(10), acctest.RandString(10), acctest.RandString(10), acctest.RandString(10))
