package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"google.golang.org/api/compute/v1"
)

func TestAccComputeNetwork_explicitAutoSubnet(t *testing.T) {
	t.Parallel()

	var network compute.Network

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeNetworkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeNetwork_basic(randString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeNetworkExists(
						t, "google_compute_network.bar", &network),
					testAccCheckComputeNetworkIsAutoSubnet(
						t, "google_compute_network.bar", &network),
				),
			},
			{
				ResourceName:      "google_compute_network.bar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeNetwork_customSubnet(t *testing.T) {
	t.Parallel()

	var network compute.Network

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeNetworkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeNetwork_custom_subnet(randString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeNetworkExists(
						t, "google_compute_network.baz", &network),
					testAccCheckComputeNetworkIsCustomSubnet(
						t, "google_compute_network.baz", &network),
				),
			},
			{
				ResourceName:      "google_compute_network.baz",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeNetwork_routingModeAndUpdate(t *testing.T) {
	t.Parallel()

	var network compute.Network
	networkName := randString(t, 10)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeNetworkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeNetwork_routing_mode(networkName, "GLOBAL"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeNetworkExists(
						t, "google_compute_network.acc_network_routing_mode", &network),
					testAccCheckComputeNetworkHasRoutingMode(
						t, "google_compute_network.acc_network_routing_mode", &network, "GLOBAL"),
				),
			},
			// Test updating the routing field (only updatable field).
			{
				Config: testAccComputeNetwork_routing_mode(networkName, "REGIONAL"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeNetworkExists(
						t, "google_compute_network.acc_network_routing_mode", &network),
					testAccCheckComputeNetworkHasRoutingMode(
						t, "google_compute_network.acc_network_routing_mode", &network, "REGIONAL"),
				),
			},
		},
	})
}

func TestAccComputeNetwork_default_routing_mode(t *testing.T) {
	t.Parallel()

	var network compute.Network

	expectedRoutingMode := "REGIONAL"

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeNetworkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeNetwork_basic(randString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeNetworkExists(
						t, "google_compute_network.bar", &network),
					testAccCheckComputeNetworkHasRoutingMode(
						t, "google_compute_network.bar", &network, expectedRoutingMode),
				),
			},
		},
	})
}

func TestAccComputeNetwork_networkDeleteDefaultRoute(t *testing.T) {
	t.Parallel()

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeNetworkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeNetwork_deleteDefaultRoute(randString(t, 10)),
			},
		},
	})
}

func testAccCheckComputeNetworkExists(t *testing.T, n string, network *compute.Network) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.Attributes["name"] == "" {
			return fmt.Errorf("No ID is set")
		}

		config := googleProviderConfig(t)

		found, err := config.clientCompute.Networks.Get(
			config.Project, rs.Primary.Attributes["name"]).Do()
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.Attributes["name"] {
			return fmt.Errorf("Network not found")
		}

		*network = *found

		return nil
	}
}

func testAccCheckComputeNetworkIsAutoSubnet(t *testing.T, n string, network *compute.Network) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := googleProviderConfig(t)

		found, err := config.clientCompute.Networks.Get(
			config.Project, network.Name).Do()
		if err != nil {
			return err
		}

		if !found.AutoCreateSubnetworks {
			return fmt.Errorf("should have AutoCreateSubnetworks = true")
		}

		if found.IPv4Range != "" {
			return fmt.Errorf("should not have IPv4Range")
		}

		return nil
	}
}

func testAccCheckComputeNetworkIsCustomSubnet(t *testing.T, n string, network *compute.Network) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := googleProviderConfig(t)

		found, err := config.clientCompute.Networks.Get(
			config.Project, network.Name).Do()
		if err != nil {
			return err
		}

		if found.AutoCreateSubnetworks {
			return fmt.Errorf("should have AutoCreateSubnetworks = false")
		}

		if found.IPv4Range != "" {
			return fmt.Errorf("should not have IPv4Range")
		}

		return nil
	}
}

func testAccCheckComputeNetworkHasRoutingMode(t *testing.T, n string, network *compute.Network, routingMode string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := googleProviderConfig(t)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.Attributes["routing_mode"] == "" {
			return fmt.Errorf("Routing mode not found on resource")
		}

		found, err := config.clientCompute.Networks.Get(
			config.Project, network.Name).Do()
		if err != nil {
			return err
		}

		foundRoutingMode := found.RoutingConfig.RoutingMode

		if routingMode != foundRoutingMode {
			return fmt.Errorf("Expected routing mode %s to match actual routing mode %s", routingMode, foundRoutingMode)
		}

		return nil
	}
}

func testAccComputeNetwork_basic(suffix string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "bar" {
  name                    = "network-test-%s"
  auto_create_subnetworks = true
}
`, suffix)
}

func testAccComputeNetwork_custom_subnet(suffix string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "baz" {
  name                    = "network-test-%s"
  auto_create_subnetworks = false
}
`, suffix)
}

func testAccComputeNetwork_routing_mode(network, routingMode string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "acc_network_routing_mode" {
  name         = "network-test-%s"
  routing_mode = "%s"
}
`, network, routingMode)
}

func testAccComputeNetwork_deleteDefaultRoute(suffix string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "bar" {
  name                            = "network-test-%s"
  delete_default_routes_on_create = true
  auto_create_subnetworks         = false
}
`, suffix)
}
