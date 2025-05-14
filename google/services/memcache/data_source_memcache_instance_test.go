// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package memcache_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccMemcacheInstanceDatasourceConfig(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMemcacheInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMemcacheInstanceDatasourceConfig(context),
				Check: acctest.CheckDataSourceStateMatchesResourceState(
					"data.google_memcache_instance.default",
					"google_memcache_instance.instance",
				),
			},
		},
	})
}

func testAccMemcacheInstanceDatasourceConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "memcache_network" {
  name                     = "test-network"
}

resource "google_compute_global_address" "service_range" {
  name                     = "address"
  purpose                  = "VPC_PEERING"
  address_type             = "INTERNAL"
  prefix_length            = 16
  network                  = google_compute_network.memcache_network.id
}

resource "google_service_networking_connection" "private_service_connection" {
  network                  = google_compute_network.memcache_network.id
  service                  = "servicenetworking.googleapis.com"
  reserved_peering_ranges  = [google_compute_global_address.service_range.name]
}

resource "google_memcache_instance" "instance" {
  name                     = "test-instance"
  authorized_network       = google_service_networking_connection.private_service_connection.network
  region                   = "us-central1"
  node_config {
    cpu_count              = 1
    memory_size_mb         = 1024
  }
  node_count               = 1
}
data "google_memcache_instance" "default" {
name                       = google_memcache_instance.instance.name
region                     = "us-central1"
}
`, context)
}
