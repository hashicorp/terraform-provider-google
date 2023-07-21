// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"testing"
)

func TestAccDataSourceRegionInstanceGroup(t *testing.T) {
	// Randomness in instance template
	acctest.SkipIfVcr(t)
	t.Parallel()
	name := "tf-test-" + acctest.RandString(t, 6)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRegionInstanceGroup_basic(fmt.Sprintf("tf-test-rigm--%d", acctest.RandInt(t)), name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_compute_region_instance_group.data_source", "name", name),
					resource.TestCheckResourceAttr("data.google_compute_region_instance_group.data_source", "project", envvar.GetTestProjectFromEnv()),
					resource.TestCheckResourceAttr("data.google_compute_region_instance_group.data_source", "instances.#", "1")),
			},
		},
	})
}

func testAccDataSourceRegionInstanceGroup_basic(rigmName, instanceManagerName string) string {
	return fmt.Sprintf(`
resource "google_compute_target_pool" "foo" {
  name = "%s"
}

data "google_compute_image" "debian" {
  project = "debian-cloud"
  name    = "debian-11-bullseye-v20220719"
}

resource "google_compute_instance_template" "foo" {
  machine_type = "e2-medium"
  disk {
    source_image = data.google_compute_image.debian.self_link
  }
  network_interface {
    access_config {
    }
    network = "default"
  }
}

resource "google_compute_region_instance_group_manager" "foo" {
  name               = "%s"
  base_instance_name = "foo"
  version {
    instance_template = google_compute_instance_template.foo.self_link
    name              = "primary"
  }
  region       = "us-central1"
  target_pools = [google_compute_target_pool.foo.self_link]
  target_size  = 1

  named_port {
    name = "web"
    port = 80
  }
  wait_for_instances = true
}

data "google_compute_region_instance_group" "data_source" {
  self_link = google_compute_region_instance_group_manager.foo.instance_group
}
`, rigmName, instanceManagerName)
}
