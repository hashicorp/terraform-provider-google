// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package memcache_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccMemcacheInstance_update(t *testing.T) {
	t.Parallel()

	prefix := fmt.Sprintf("%d", acctest.RandInt(t))
	name := fmt.Sprintf("tf-test-%s", prefix)
	network := acctest.BootstrapSharedServiceNetworkingConnection(t, "memcache-instance-update-1")

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMemcacheInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMemcacheInstance_update(prefix, name, network),
			},
			{
				ResourceName:      "google_memcache_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMemcacheInstance_update2(prefix, name, network),
			},
			{
				ResourceName:      "google_memcache_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMemcacheInstance_update(prefix, name, network string) string {
	return fmt.Sprintf(`
resource "google_memcache_instance" "test" {
  name = "%s"
  region = "us-central1"
  authorized_network = data.google_compute_network.memcache_network.id

  node_config {
    cpu_count      = 1
    memory_size_mb = 1024
  }
  node_count = 1

  memcache_parameters {
    params = {
      "listen-backlog" = "2048"
      "max-item-size" = "8388608"
    }
  }
}

data "google_compute_network" "memcache_network" {
  name = "%s"
}
`, name, network)
}

func testAccMemcacheInstance_update2(prefix, name, network string) string {
	return fmt.Sprintf(`
resource "google_memcache_instance" "test" {
  name = "%s"
  region = "us-central1"
  authorized_network = data.google_compute_network.memcache_network.id

  node_config {
    cpu_count      = 1
    memory_size_mb = 1024
  }
  node_count = 2

  memcache_parameters {
    params = {
      "listen-backlog" = "2048"
      "max-item-size" = "8388608"
    }
  }
}

data "google_compute_network" "memcache_network" {
  name = "%s"
}
`, name, network)
}
