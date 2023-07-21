// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccComputeNetworkPeering_basic(t *testing.T) {
	t.Parallel()

	primaryNetworkName := fmt.Sprintf("tf-test-network-peering-1-%d", acctest.RandInt(t))
	peeringName := fmt.Sprintf("peering-test-1-%d", acctest.RandInt(t))
	importId := fmt.Sprintf("%s/%s/%s", envvar.GetTestProjectFromEnv(), primaryNetworkName, peeringName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComputeNetworkPeeringDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeNetworkPeering_basic(primaryNetworkName, peeringName, acctest.RandString(t, 10)),
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

	primaryNetworkName := fmt.Sprintf("tf-test-network-peering-1-%d", acctest.RandInt(t))
	peeringName := fmt.Sprintf("peering-test-%d", acctest.RandInt(t))
	importId := fmt.Sprintf("%s/%s/%s", envvar.GetTestProjectFromEnv(), primaryNetworkName, peeringName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComputeNetworkPeeringDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeNetworkPeering_subnetRoutes(primaryNetworkName, peeringName, acctest.RandString(t, 10)),
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

	primaryNetworkName := fmt.Sprintf("tf-test-network-peering-1-%d", acctest.RandInt(t))
	peeringName := fmt.Sprintf("peering-test-%d", acctest.RandInt(t))
	importId := fmt.Sprintf("%s/%s/%s", envvar.GetTestProjectFromEnv(), primaryNetworkName, peeringName)
	suffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
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

func TestAccComputeNetworkPeering_stackType(t *testing.T) {
	t.Parallel()

	primaryNetworkName := fmt.Sprintf("tf-test-network-1-%d", acctest.RandInt(t))
	peeringNetworkName := fmt.Sprintf("tf-test-network-2-%d", acctest.RandInt(t))
	peeringName := fmt.Sprintf("tf-test-peering-%d", acctest.RandInt(t))
	importId := fmt.Sprintf("%s/%s/%s", envvar.GetTestProjectFromEnv(), primaryNetworkName, peeringName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccComputeNetworkPeeringDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeNetworkPeering_stackTypeDefault(primaryNetworkName, peeringNetworkName, peeringName),
			},
			{
				ResourceName:      "google_compute_network_peering.foo",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     importId,
			},
			{
				Config: testAccComputeNetworkPeering_stackTypeUpdate(primaryNetworkName, peeringNetworkName, peeringName),
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
		config := acctest.GoogleProviderConfig(t)

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

func testAccComputeNetworkPeering_stackTypeDefault(primaryNetworkName, peeringNetworkName, peeringName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "network1" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_network" "network2" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_network_peering" "foo" {
  name         = "%s"
  network      = google_compute_network.network1.self_link
  peer_network = google_compute_network.network2.self_link
}
`, primaryNetworkName, peeringNetworkName, peeringName)
}

func testAccComputeNetworkPeering_stackTypeUpdate(primaryNetworkName, peeringNetworkName, peeringName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "network1" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_network" "network2" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_network_peering" "foo" {
  name         = "%s"
  network      = google_compute_network.network1.self_link
  peer_network = google_compute_network.network2.self_link
  stack_type   = "IPV4_IPV6"
}
`, primaryNetworkName, peeringNetworkName, peeringName)
}
