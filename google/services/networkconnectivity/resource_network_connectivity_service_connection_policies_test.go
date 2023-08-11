// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package networkconnectivity_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccNetworkConnectivityServiceConnectionPolicy_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"networkProducerName":         fmt.Sprintf("tf-test-network-%s", acctest.RandString(t, 10)),
		"subnetworkProducerName1":     fmt.Sprintf("tf-test-subnet-producer-%s", acctest.RandString(t, 10)),
		"subnetworkProducerName2":     fmt.Sprintf("tf-test-subnet-producer-%s", acctest.RandString(t, 10)),
		"serviceConnectionPolicyName": fmt.Sprintf("tf-test-service-connection-policy-%s", acctest.RandString(t, 10)),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkConnectivityServiceConnectionPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkConnectivityServiceConnectionPolicy_basic(context),
			},
			{
				ResourceName:      "google_network_connectivity_service_connection_policy.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNetworkConnectivityServiceConnectionPolicy_update(context),
			},
			{
				ResourceName:      "google_network_connectivity_service_connection_policy.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNetworkConnectivityServiceConnectionPolicy_basic(context),
			},
			{
				ResourceName:      "google_network_connectivity_service_connection_policy.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNetworkConnectivityServiceConnectionPolicy_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
  resource "google_compute_network" "producer_net" {
    name                    = "%{networkProducerName}"
    auto_create_subnetworks = false
  }
  
  resource "google_compute_subnetwork" "producer_subnet" {
    name          = "%{subnetworkProducerName1}"
    ip_cidr_range = "10.0.0.0/16"
    region        = "us-central1"
    network       = google_compute_network.producer_net.id
  }
  
  resource "google_network_connectivity_service_connection_policy" "default" {
    name = "%{serviceConnectionPolicyName}"
    location = "us-central1"
    service_class = "gcp-memorystore-redis"
    network = google_compute_network.producer_net.id
    psc_config {
      subnetworks = [google_compute_subnetwork.producer_subnet.id]
      limit = 2
    }
  }
`, context)
}

func testAccNetworkConnectivityServiceConnectionPolicy_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "producer_net" {
  name                    = "%{networkProducerName}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "producer_subnet1" {
  name          = "%{subnetworkProducerName2}"
  ip_cidr_range = "10.1.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.producer_net.id
}

resource "google_network_connectivity_service_connection_policy" "default" {
  name = "%{serviceConnectionPolicyName}"
  location = "us-central1"
  service_class = "gcp-memorystore-redis"
  network = google_compute_network.producer_net.id
  psc_config {
    subnetworks = [google_compute_subnetwork.producer_subnet1.id]
    limit = 4
  }
  labels      = {
    foo = "bar"
  }
}
`, context)
}
