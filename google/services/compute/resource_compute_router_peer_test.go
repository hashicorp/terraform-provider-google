// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccComputeRouterBgpPeer_routerPeerRouterAppliance(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeRouterBgpPeerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRouterBgpPeer_routerPeerRouterAppliance(context),
			},
			{
				ResourceName:            "google_compute_router_peer.peer",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"router_appliance_instance", "router", "region"},
			},
		},
	})
}

func testAccComputeRouterBgpPeer_routerPeerRouterAppliance(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "network" {
  name                    = "tf-test-my-router%{random_suffix}-net"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "subnetwork" {
  name          = "tf-test-my-router%{random_suffix}-sub"
  network       = google_compute_network.network.self_link
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
}

resource "google_compute_address" "addr_intf" {
  name         = "tf-test-my-router%{random_suffix}-addr-intf"
  region       = google_compute_subnetwork.subnetwork.region
  subnetwork   = google_compute_subnetwork.subnetwork.id
  address_type = "INTERNAL"
}

resource "google_compute_address" "addr_intf_redundant" {
  name         = "tf-test-my-router%{random_suffix}-addr-intf-red"
  region       = google_compute_subnetwork.subnetwork.region
  subnetwork   = google_compute_subnetwork.subnetwork.id
  address_type = "INTERNAL"
}

resource "google_compute_address" "addr_peer" {
  name         = "tf-test-my-router%{random_suffix}-addr-peer"
  region       = google_compute_subnetwork.subnetwork.region
  subnetwork   = google_compute_subnetwork.subnetwork.id
  address_type = "INTERNAL"
}

resource "google_compute_instance" "instance" {
  name           = "router-appliance"
  zone           = "us-central1-a"
  machine_type   = "e2-medium"
  can_ip_forward = true

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
    }
  }

  network_interface {
    network_ip = google_compute_address.addr_peer.address
    subnetwork = google_compute_subnetwork.subnetwork.self_link
  }
}

resource "google_network_connectivity_hub" "hub" {
  name = "tf-test-my-router%{random_suffix}-hub"
}

resource "google_network_connectivity_spoke" "spoke" {
  name     = "tf-test-my-router%{random_suffix}-spoke"
  location = google_compute_subnetwork.subnetwork.region
  hub      = google_network_connectivity_hub.hub.id

  linked_router_appliance_instances {
    instances {
      virtual_machine = google_compute_instance.instance.self_link
      ip_address      = google_compute_address.addr_peer.address
    }
    site_to_site_data_transfer = false
  }
}

resource "google_compute_router" "router" {
  name    = "tf-test-my-router%{random_suffix}-router"
  region  = google_compute_subnetwork.subnetwork.region
  network = google_compute_network.network.self_link
  bgp {
    asn = 64514
  }
}

resource "google_compute_router_interface" "interface_redundant" {
  name               = "tf-test-my-router%{random_suffix}-intf-red"
  region             = google_compute_router.router.region
  router             = google_compute_router.router.name
  subnetwork         = google_compute_subnetwork.subnetwork.self_link
  private_ip_address = google_compute_address.addr_intf_redundant.address
}

resource "google_compute_router_interface" "interface" {
  name                = "tf-test-my-router%{random_suffix}-intf"
  region              = google_compute_router.router.region
  router              = google_compute_router.router.name
  subnetwork          = google_compute_subnetwork.subnetwork.self_link
  private_ip_address  = google_compute_address.addr_intf.address
  redundant_interface = google_compute_router_interface.interface_redundant.name
}

resource "google_compute_router_peer" "peer" {
  name                      = "tf-test-my-router-peer%{random_suffix}"
  router                    = google_compute_router.router.name
  region                    = google_compute_router.router.region
  interface                 = google_compute_router_interface.interface.name
  router_appliance_instance = google_compute_instance.instance.self_link
  peer_asn                  = 65513
  peer_ip_address           = google_compute_address.addr_peer.address
}
`, context)
}

func testAccCheckComputeRouterBgpPeerDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_compute_router_peer" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{ComputeBasePath}}projects/{{project}}/regions/{{region}}/routers/{{router}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("ComputeRouterBgpPeer still exists at %s", url)
			}
		}

		return nil
	}
}
