package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"google.golang.org/api/compute/v1"
)

func TestAccComputeNetwork_explicitAutoSubnet(t *testing.T) {
	t.Parallel()

	var network compute.Network

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeNetwork_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeNetworkExists(
						"google_compute_network.bar", &network),
					testAccCheckComputeNetworkIsAutoSubnet(
						"google_compute_network.bar", &network),
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

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeNetwork_custom_subnet(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeNetworkExists(
						"google_compute_network.baz", &network),
					testAccCheckComputeNetworkIsCustomSubnet(
						"google_compute_network.baz", &network),
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
	networkName := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeNetwork_routing_mode(networkName, "GLOBAL"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeNetworkExists(
						"google_compute_network.acc_network_routing_mode", &network),
					testAccCheckComputeNetworkHasRoutingMode(
						"google_compute_network.acc_network_routing_mode", &network, "GLOBAL"),
				),
			},
			// Test updating the routing field (only updatable field).
			{
				Config: testAccComputeNetwork_routing_mode(networkName, "REGIONAL"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeNetworkExists(
						"google_compute_network.acc_network_routing_mode", &network),
					testAccCheckComputeNetworkHasRoutingMode(
						"google_compute_network.acc_network_routing_mode", &network, "REGIONAL"),
				),
			},
		},
	})
}

func TestAccComputeNetwork_default_routing_mode(t *testing.T) {
	t.Parallel()

	var network compute.Network

	expectedRoutingMode := "REGIONAL"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeNetwork_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeNetworkExists(
						"google_compute_network.bar", &network),
					testAccCheckComputeNetworkHasRoutingMode(
						"google_compute_network.bar", &network, expectedRoutingMode),
				),
			},
		},
	})
}

func TestAccComputeNetwork_networkDeleteDefaultRoute(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeNetwork_deleteDefaultRoute(),
			},
		},
	})
}

func testAccCheckComputeNetworkExists(n string, network *compute.Network) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.Attributes["name"] == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

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

func testAccCheckComputeNetworkIsAutoSubnet(n string, network *compute.Network) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

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

func testAccCheckComputeNetworkIsCustomSubnet(n string, network *compute.Network) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

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

func testAccCheckComputeNetworkHasRoutingMode(n string, network *compute.Network, routingMode string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

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

func testAccComputeNetwork_basic() string {
	return fmt.Sprintf(`
resource "google_compute_network" "bar" {
  name                    = "network-test-%s"
  auto_create_subnetworks = true
}
`, acctest.RandString(10))
}

func testAccComputeNetwork_custom_subnet() string {
	return fmt.Sprintf(`
resource "google_compute_network" "baz" {
  name                    = "network-test-%s"
  auto_create_subnetworks = false
}
`, acctest.RandString(10))
}

func testAccComputeNetwork_routing_mode(network, routingMode string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "acc_network_routing_mode" {
  name         = "network-test-%s"
  routing_mode = "%s"
}
`, network, routingMode)
}

func testAccComputeNetwork_deleteDefaultRoute() string {
	return fmt.Sprintf(`
resource "google_compute_network" "bar" {
  name                            = "network-test-%s"
  delete_default_routes_on_create = true
  auto_create_subnetworks         = false
}
`, acctest.RandString(10))
}
