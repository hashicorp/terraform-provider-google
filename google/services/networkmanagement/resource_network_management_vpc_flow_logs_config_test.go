// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package networkmanagement_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccNetworkManagementVpcFlowLogsConfig_updateInterconnect(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkManagementVpcFlowLogsConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkManagementVpcFlowLogsConfig_fullInterconnect(context),
			},
			{
				ResourceName:            "google_network_management_vpc_flow_logs_config.interconnect-test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "terraform_labels", "vpc_flow_logs_config_id"},
			},
			{
				Config: testAccNetworkManagementVpcFlowLogsConfig_updateInterconnect(context),
			},
			{
				ResourceName:            "google_network_management_vpc_flow_logs_config.interconnect-test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "terraform_labels", "vpc_flow_logs_config_id"},
			},
		},
	})
}

func testAccNetworkManagementVpcFlowLogsConfig_fullInterconnect(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}

resource "google_network_management_vpc_flow_logs_config" "interconnect-test" {
  vpc_flow_logs_config_id = "tf-test-full-interconnect-test-id%{random_suffix}"
  location                = "global"
  interconnect_attachment = "projects/${data.google_project.project.number}/regions/us-east4/interconnectAttachments/${google_compute_interconnect_attachment.attachment.name}"
}

resource "google_compute_network" "network" {
  name     = "tf-test-full-interconnect-test-network%{random_suffix}"
}

resource "google_compute_router" "router" {
  name    = "tf-test-full-interconnect-test-router%{random_suffix}"
  network = google_compute_network.network.name
  bgp {
    asn = 16550
  }
}

resource "google_compute_interconnect_attachment" "attachment" {
  name                     = "tf-test-full-interconnect-test-id%{random_suffix}"
  edge_availability_domain = "AVAILABILITY_DOMAIN_1"
  type                     = "PARTNER"
  router                   = google_compute_router.router.id
  mtu                      = 1500
}

`, context)
}

func testAccNetworkManagementVpcFlowLogsConfig_updateInterconnect(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}

resource "google_network_management_vpc_flow_logs_config" "interconnect-test" {
  vpc_flow_logs_config_id = "tf-test-full-interconnect-test-id%{random_suffix}"
  location                = "global"
  interconnect_attachment = "projects/${data.google_project.project.number}/regions/us-east4/interconnectAttachments/${google_compute_interconnect_attachment.attachment.name}"
  state                   = "DISABLED"
  aggregation_interval    = "INTERVAL_30_SEC"
  description             = "This is an updated description"
  flow_sampling           = 0.5
  metadata                = "EXCLUDE_ALL_METADATA"
}

resource "google_compute_network" "network" {
  name     = "tf-test-full-interconnect-test-network%{random_suffix}"
}

resource "google_compute_router" "router" {
  name    = "tf-test-full-interconnect-test-router%{random_suffix}"
  network = google_compute_network.network.name
  bgp {
    asn = 16550
  }
}

resource "google_compute_interconnect_attachment" "attachment" {
  name                     = "tf-test-full-interconnect-test-id%{random_suffix}"
  edge_availability_domain = "AVAILABILITY_DOMAIN_1"
  type                     = "PARTNER"
  router                   = google_compute_router.router.id
  mtu                      = 1500
}

`, context)
}

func TestAccNetworkManagementVpcFlowLogsConfig_updateVpn(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkManagementVpcFlowLogsConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkManagementVpcFlowLogsConfig_fullVpn(context),
			},
			{
				ResourceName:            "google_network_management_vpc_flow_logs_config.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "terraform_labels", "vpc_flow_logs_config_id"},
			},
			{
				Config: testAccNetworkManagementVpcFlowLogsConfig_updateVpn(context),
			},
			{
				ResourceName:            "google_network_management_vpc_flow_logs_config.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "terraform_labels", "vpc_flow_logs_config_id"},
			},
		},
	})
}

func testAccNetworkManagementVpcFlowLogsConfig_fullVpn(context map[string]interface{}) string {
	vpcFlowLogsCfg := acctest.Nprintf(`
data "google_project" "project" {
}

resource "google_network_management_vpc_flow_logs_config" "example" {
  vpc_flow_logs_config_id = "id-example-%{random_suffix}"
  location                = "global"
  vpn_tunnel              = "projects/${data.google_project.project.number}/regions/us-central1/vpnTunnels/${google_compute_vpn_tunnel.tunnel.name}"
}
`, context)
	return fmt.Sprintf("%s\n\n%s\n\n", vpcFlowLogsCfg, testAccNetworkManagementVpcFlowLogsConfig_baseResources(context))
}

func testAccNetworkManagementVpcFlowLogsConfig_updateVpn(context map[string]interface{}) string {
	vpcFlowLogsCfg := acctest.Nprintf(`
data "google_project" "project" {
}

resource "google_network_management_vpc_flow_logs_config" "example" {
  vpc_flow_logs_config_id = "id-example-%{random_suffix}"
  location                = "global"
  vpn_tunnel              = "projects/${data.google_project.project.number}/regions/us-central1/vpnTunnels/${google_compute_vpn_tunnel.tunnel.name}"
  state                   = "DISABLED"
  aggregation_interval    = "INTERVAL_30_SEC"
  description             = "This is an updated description"
  flow_sampling           = 0.5
  metadata                = "EXCLUDE_ALL_METADATA"
}
`, context)
	return fmt.Sprintf("%s\n\n%s\n\n", vpcFlowLogsCfg, testAccNetworkManagementVpcFlowLogsConfig_baseResources(context))
}

func testAccNetworkManagementVpcFlowLogsConfig_baseResources(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_vpn_tunnel" "tunnel" {
  name               = "tf-test-example-tunnel%{random_suffix}"
  peer_ip            = "15.0.0.120"
  shared_secret      = "a secret message"
  target_vpn_gateway = google_compute_vpn_gateway.target_gateway.id

  depends_on = [
    google_compute_forwarding_rule.fr_esp,
    google_compute_forwarding_rule.fr_udp500,
    google_compute_forwarding_rule.fr_udp4500,
  ]
}

resource "google_compute_vpn_gateway" "target_gateway" {
  name     = "tf-test-example-gateway%{random_suffix}"
  network  = google_compute_network.network.id
}

resource "google_compute_network" "network" {
  name     = "tf-test-example-network%{random_suffix}"
}

resource "google_compute_address" "vpn_static_ip" {
  name     = "tf-test-example-address%{random_suffix}"
}

resource "google_compute_forwarding_rule" "fr_esp" {
  name        = "tf-test-example-fresp%{random_suffix}"
  ip_protocol = "ESP"
  ip_address  = google_compute_address.vpn_static_ip.address
  target      = google_compute_vpn_gateway.target_gateway.id
}

resource "google_compute_forwarding_rule" "fr_udp500" {
  name        = "tf-test-example-fr500%{random_suffix}"
  ip_protocol = "UDP"
  port_range  = "500"
  ip_address  = google_compute_address.vpn_static_ip.address
  target      = google_compute_vpn_gateway.target_gateway.id
}

resource "google_compute_forwarding_rule" "fr_udp4500" {
  name        = "tf-test-example-fr4500%{random_suffix}"
  ip_protocol = "UDP"
  port_range  = "4500"
  ip_address  = google_compute_address.vpn_static_ip.address
  target      = google_compute_vpn_gateway.target_gateway.id
}

resource "google_compute_route" "route" {
  name                = "tf-test-example-route%{random_suffix}"
  network             = google_compute_network.network.name
  dest_range          = "15.0.0.0/24"
  priority            = 1000
  next_hop_vpn_tunnel = google_compute_vpn_tunnel.tunnel.id
}
`, context)
}
