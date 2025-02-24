// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package memorystore_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccMemorystoreInstanceDatasourceConfig(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"network_name":  acctest.BootstrapSharedServiceNetworkingConnection(t, "memorystore-instance-ds"),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMemorystoreInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMemorystoreInstanceDatasourceConfig(context),
			},
		},
	})
}

func testAccMemorystoreInstanceDatasourceConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_memorystore_instance" "instance-basic" {
  instance_id                 = "tf-test-memorystore-instance%{random_suffix}"
  shard_count                 = 3
  desired_psc_auto_connections {
    network                   = google_compute_network.producer_net.id
    project_id                = data.google_project.project.project_id
  }
  location                    = "us-central1"
  deletion_protection_enabled = false
  depends_on                  = [google_network_connectivity_service_connection_policy.default]

}

resource "google_network_connectivity_service_connection_policy" "default" {
  name                        = "%{network_name}-policy"
  location                    = "us-central1"
  service_class               = "gcp-memorystore"
  description                 = "my basic service connection policy"
  network                     = google_compute_network.producer_net.id
  psc_config {
    subnetworks               = [google_compute_subnetwork.producer_subnet.id]
  }
}


resource "google_compute_subnetwork" "producer_subnet" {
	name                      = "%{network_name}-sn"
	ip_cidr_range             = "10.0.0.248/29"
	region                    = "us-central1"
	network                   = google_compute_network.producer_net.id
}

resource "google_compute_network" "producer_net" {
  name                        = "%{network_name}-vpc"
  auto_create_subnetworks     = false
}

 data "google_project" "project" {
 }

data "google_memorystore_instance" "default" {
  instance_id                 = google_memorystore_instance.instance-basic.instance_id
  location                    = "us-central1"

}
`, context)
}
