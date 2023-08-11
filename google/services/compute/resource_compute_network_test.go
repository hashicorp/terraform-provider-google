// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"google.golang.org/api/compute/v1"
)

func TestAccComputeNetwork_explicitAutoSubnet(t *testing.T) {
	t.Parallel()

	var network compute.Network

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeNetworkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeNetwork_basic(acctest.RandString(t, 10)),
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

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeNetworkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeNetwork_custom_subnet(acctest.RandString(t, 10)),
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
	networkName := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeNetworkDestroyProducer(t),
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

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeNetworkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeNetwork_basic(acctest.RandString(t, 10)),
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

	var network compute.Network

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeNetworkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeNetwork_deleteDefaultRoute(acctest.RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeNetworkExists(
						t, "google_compute_network.bar", &network),
					testAccCheckComputeNetworkDefaultRoutesDeleted(
						t, "google_compute_network.bar", &network),
				),
			},
		},
	})
}

func TestAccComputeNetwork_networkFirewallPolicyEnforcementOrderAndUpdate(t *testing.T) {
	t.Parallel()

	var network compute.Network
	var updatedNetwork compute.Network
	networkName := acctest.RandString(t, 10)

	defaultNetworkFirewallPolicyEnforcementOrder := "AFTER_CLASSIC_FIREWALL"
	explicitNetworkFirewallPolicyEnforcementOrder := "BEFORE_CLASSIC_FIREWALL"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeNetworkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeNetwork_networkFirewallPolicyEnforcementOrderDefault(networkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeNetworkExists(
						t, "google_compute_network.acc_network_firewall_policy_enforcement_order", &network),
					testAccCheckComputeNetworkHasNetworkFirewallPolicyEnforcementOrder(
						t, "google_compute_network.acc_network_firewall_policy_enforcement_order", &network, defaultNetworkFirewallPolicyEnforcementOrder),
				),
			},
			{
				ResourceName:            "google_compute_network.acc_network_firewall_policy_enforcement_order",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			// Test updating the enforcement order works and updates in-place
			{
				Config: testAccComputeNetwork_networkFirewallPolicyEnforcementOrderUpdate(networkName, explicitNetworkFirewallPolicyEnforcementOrder),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeNetworkExists(
						t, "google_compute_network.acc_network_firewall_policy_enforcement_order", &updatedNetwork),
					testAccCheckComputeNetworkHasNetworkFirewallPolicyEnforcementOrder(
						t, "google_compute_network.acc_network_firewall_policy_enforcement_order", &updatedNetwork, explicitNetworkFirewallPolicyEnforcementOrder),
					testAccCheckComputeNetworkWasUpdated(&updatedNetwork, &network),
				),
			},
			{
				ResourceName:            "google_compute_network.acc_network_firewall_policy_enforcement_order",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
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

		config := acctest.GoogleProviderConfig(t)

		found, err := config.NewComputeClient(config.UserAgent).Networks.Get(
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

func testAccCheckComputeNetworkDefaultRoutesDeleted(t *testing.T, n string, network *compute.Network) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.Attributes["name"] == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acctest.GoogleProviderConfig(t)

		routes, err := config.NewComputeClient(config.UserAgent).Routes.List(config.Project).Filter(fmt.Sprintf("(network=\"%s\") AND (destRange=\"0.0.0.0/0\")", network.SelfLink)).Do()
		if err != nil {
			return err
		}

		if len(routes.Items) > 0 {
			return fmt.Errorf("Default routes were not deleted")
		}

		return nil
	}
}

func testAccCheckComputeNetworkIsAutoSubnet(t *testing.T, n string, network *compute.Network) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		found, err := config.NewComputeClient(config.UserAgent).Networks.Get(
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
		config := acctest.GoogleProviderConfig(t)

		found, err := config.NewComputeClient(config.UserAgent).Networks.Get(
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
		config := acctest.GoogleProviderConfig(t)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.Attributes["routing_mode"] == "" {
			return fmt.Errorf("Routing mode not found on resource")
		}

		found, err := config.NewComputeClient(config.UserAgent).Networks.Get(
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

func testAccCheckComputeNetworkHasNetworkFirewallPolicyEnforcementOrder(t *testing.T, n string, network *compute.Network, order string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.Attributes["network_firewall_policy_enforcement_order"] == "" {
			return fmt.Errorf("Network firewall policy enforcement order not found on resource")
		}

		found, err := config.NewComputeClient(config.UserAgent).Networks.Get(
			config.Project, network.Name).Do()
		if err != nil {
			return err
		}

		foundNetworkFirewallPolicyEnforcementOrder := found.NetworkFirewallPolicyEnforcementOrder

		if order != foundNetworkFirewallPolicyEnforcementOrder {
			return fmt.Errorf("Expected network firewall policy enforcement order %s to match %s", order, foundNetworkFirewallPolicyEnforcementOrder)
		}

		return nil
	}
}

func testAccCheckComputeNetworkWasUpdated(newNetwork *compute.Network, oldNetwork *compute.Network) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if oldNetwork.CreationTimestamp != newNetwork.CreationTimestamp {
			return fmt.Errorf("expected compute network to have been updated (had same creation time), instead was recreated - old creation time %s, new creation time %s", oldNetwork.CreationTimestamp, newNetwork.CreationTimestamp)
		}
		return nil
	}
}

func testAccComputeNetwork_basic(suffix string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "bar" {
  name                    = "tf-test-network-basic-%s"
  auto_create_subnetworks = true
}
`, suffix)
}

func testAccComputeNetwork_custom_subnet(suffix string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "baz" {
  name                    = "tf-test-network-custom-sn-%s"
  auto_create_subnetworks = false
}
`, suffix)
}

func testAccComputeNetwork_routing_mode(network, routingMode string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "acc_network_routing_mode" {
  name         = "tf-test-network-routing-mode-%s"
  routing_mode = "%s"
}
`, network, routingMode)
}

func testAccComputeNetwork_deleteDefaultRoute(suffix string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "bar" {
  name                            = "tf-test-network-delete-default-routes-%s"
  delete_default_routes_on_create = true
  auto_create_subnetworks         = false
}
`, suffix)
}

func testAccComputeNetwork_networkFirewallPolicyEnforcementOrderDefault(network string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "acc_network_firewall_policy_enforcement_order" {
  name = "tf-test-network-firewall-policy-enforcement-order-%s"
}
`, network)
}

func testAccComputeNetwork_networkFirewallPolicyEnforcementOrderUpdate(network, order string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "acc_network_firewall_policy_enforcement_order" {
  name                                      = "tf-test-network-firewall-policy-enforcement-order-%s"
  network_firewall_policy_enforcement_order = "%s"
}
`, network, order)
}
