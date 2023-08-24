// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccComputeNetworkEndpointGroup_networkEndpointGroup(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeNetworkEndpointGroupDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeNetworkEndpointGroup_networkEndpointGroup(context),
			},
			{
				ResourceName:            "google_compute_network_endpoint_group.neg",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"network", "subnetwork", "zone"},
			},
		},
	})
}

func TestAccComputeNetworkEndpointGroup_internalEndpoint(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeNetworkEndpointGroupDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeNetworkEndpointGroup_internalEndpoint(context),
			},
			{
				ResourceName:            "google_compute_network_endpoint_group.neg",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"network", "subnetwork", "zone"},
			},
		},
	})
}

func testAccComputeNetworkEndpointGroup_networkEndpointGroup(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network_endpoint_group" "neg" {
  name         = "tf-test-my-lb-neg%{random_suffix}"
  network      = google_compute_network.default.id
  default_port = "90"
  zone         = "us-central1-a"
}

resource "google_compute_network" "default" {
  name                    = "tf-test-neg-network%{random_suffix}"
  auto_create_subnetworks = true
}
`, context)
}

func testAccComputeNetworkEndpointGroup_internalEndpoint(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network_endpoint_group" "neg" {
  name                  = "tf-test-my-lb-neg%{random_suffix}"
  network               = google_compute_network.internal.id
  subnetwork            = google_compute_subnetwork.internal.id
  zone                  = "us-central1-a"
  network_endpoint_type = "GCE_VM_IP"
}

resource "google_compute_network_endpoint" "endpoint" {
  network_endpoint_group = google_compute_network_endpoint_group.neg.name
  #ip_address             = "127.0.0.1"
  instance   = google_compute_instance.default.name
  ip_address = google_compute_instance.default.network_interface[0].network_ip
}

resource "google_compute_network" "internal" {
  name                    = "tf-test-neg-network%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "internal"{
  name                    = "tf-test-my-subnetwork%{random_suffix}"
  network                 = google_compute_network.internal.id
  ip_cidr_range           = "10.128.0.0/20"
  region                  = "us-central1"
  private_ip_google_access= true
}

resource "google_compute_instance" "default" {
  name         = "tf-test-neg-%{random_suffix}"
  machine_type = "e2-medium"
  
  boot_disk {
    initialize_params {
      image = "debian-8-jessie-v20160803"
    }
   }
  
  network_interface {
    subnetwork = google_compute_subnetwork.internal.self_link
    access_config {
    }
  }
}

`, context)
}
