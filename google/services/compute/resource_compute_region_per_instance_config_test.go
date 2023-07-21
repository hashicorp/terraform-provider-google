// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccComputeRegionPerInstanceConfig_statefulBasic(t *testing.T) {
	// Multiple fine-grained resources
	acctest.SkipIfVcr(t)
	t.Parallel()

	suffix := acctest.RandString(t, 10)
	rigmName := fmt.Sprintf("tf-test-rigm-%s", suffix)
	context := map[string]interface{}{
		"rigm_name":     rigmName,
		"random_suffix": suffix,
		"config_name":   fmt.Sprintf("instance-%s", acctest.RandString(t, 10)),
		"config_name2":  fmt.Sprintf("instance-%s", acctest.RandString(t, 10)),
		"config_name3":  fmt.Sprintf("instance-%s", acctest.RandString(t, 10)),
		"config_name4":  fmt.Sprintf("instance-%s", acctest.RandString(t, 10)),
	}
	rigmId := fmt.Sprintf("projects/%s/regions/%s/instanceGroupManagers/%s",
		envvar.GetTestProjectFromEnv(), envvar.GetTestRegionFromEnv(), rigmName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Create one endpoint
				Config: testAccComputeRegionPerInstanceConfig_statefulBasic(context),
			},
			{
				ResourceName:            "google_compute_region_per_instance_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_instance_state_on_destroy", "region"},
			},
			{
				// Force-recreate old config
				Config: testAccComputeRegionPerInstanceConfig_statefulModified(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeRegionPerInstanceConfigDestroyed(t, rigmId, context["config_name"].(string)),
				),
			},
			{
				ResourceName:            "google_compute_region_per_instance_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_instance_state_on_destroy", "region"},
			},
			{
				// Add two new endpoints
				Config: testAccComputeRegionPerInstanceConfig_statefulAdditional(context),
			},
			{
				ResourceName:            "google_compute_region_per_instance_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_instance_state_on_destroy", "region"},
			},
			{
				ResourceName:            "google_compute_region_per_instance_config.with_disks",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"most_disruptive_allowed_action", "minimal_action", "remove_instance_state_on_destroy"},
			},
			{
				ResourceName:            "google_compute_region_per_instance_config.add2",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_instance_state_on_destroy", "region"},
			},
			{
				// delete all configs
				Config: testAccComputeRegionPerInstanceConfig_rigm(context),
				Check: resource.ComposeTestCheckFunc(
					// Config with remove_instance_state_on_destroy = false won't be destroyed (config4)
					testAccCheckComputeRegionPerInstanceConfigDestroyed(t, rigmId, context["config_name2"].(string)),
					testAccCheckComputeRegionPerInstanceConfigDestroyed(t, rigmId, context["config_name3"].(string)),
				),
			},
		},
	})
}

func TestAccComputeRegionPerInstanceConfig_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"rigm_name":     fmt.Sprintf("tf-test-rigm-%s", acctest.RandString(t, 10)),
		"config_name":   fmt.Sprintf("instance-%s", acctest.RandString(t, 10)),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Create one config
				Config: testAccComputeRegionPerInstanceConfig_statefulBasic(context),
			},
			{
				ResourceName:            "google_compute_region_per_instance_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_instance_state_on_destroy", "region"},
			},
			{
				// Update an existing config
				Config: testAccComputeRegionPerInstanceConfig_update(context),
			},
			{
				ResourceName:            "google_compute_region_per_instance_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_instance_state_on_destroy", "region"},
			},
		},
	})
}

func testAccComputeRegionPerInstanceConfig_statefulBasic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_region_per_instance_config" "default" {
	region_instance_group_manager = google_compute_region_instance_group_manager.rigm.name
	name = "%{config_name}"
	remove_instance_state_on_destroy = true
	preserved_state {
		metadata = {
			asdf = "asdf"
		}
	}
}
`, context) + testAccComputeRegionPerInstanceConfig_rigm(context)
}

func testAccComputeRegionPerInstanceConfig_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_region_per_instance_config" "default" {
	region_instance_group_manager = google_compute_region_instance_group_manager.rigm.name
	name = "%{config_name}"
	remove_instance_state_on_destroy = true
	preserved_state {
		metadata = {
			asdf = "foo"
			updated = "12345"
		}
	}
}
`, context) + testAccComputeRegionPerInstanceConfig_rigm(context)
}

func testAccComputeRegionPerInstanceConfig_statefulModified(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_region_per_instance_config" "default" {
	region = google_compute_region_instance_group_manager.rigm.region
	region_instance_group_manager = google_compute_region_instance_group_manager.rigm.name
	name = "%{config_name2}"
	remove_instance_state_on_destroy = true
	preserved_state {
		metadata = {
			asdf = "asdf"
		}
	}
}
`, context) + testAccComputeRegionPerInstanceConfig_rigm(context)
}

func testAccComputeRegionPerInstanceConfig_statefulAdditional(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_region_per_instance_config" "default" {
	region = google_compute_region_instance_group_manager.rigm.region
	region_instance_group_manager = google_compute_region_instance_group_manager.rigm.name
	name = "%{config_name2}"
	remove_instance_state_on_destroy = true
	preserved_state {
		metadata = {
			asdf = "asdf"
		}
	}
}

resource "google_compute_region_per_instance_config" "with_disks" {
	region = google_compute_region_instance_group_manager.rigm.region
	region_instance_group_manager = google_compute_region_instance_group_manager.rigm.name
	name = "%{config_name3}"
	most_disruptive_allowed_action = "REFRESH"
	minimal_action = "REFRESH"
	remove_instance_state_on_destroy = true
	preserved_state {
		metadata = {
			meta = "123"
		}

		disk {
			device_name = "my-stateful-disk1"
			source      = google_compute_disk.disk.id
		}

		disk {
			device_name = "my-stateful-disk2"
			source      = google_compute_disk.disk1.id
		}

		disk {
			device_name = "my-stateful-disk3"
			source      = google_compute_disk.disk2.id
		}
	}
}

resource "google_compute_region_per_instance_config" "add2" {
	region = google_compute_region_instance_group_manager.rigm.region
	region_instance_group_manager = google_compute_region_instance_group_manager.rigm.name
	name = "%{config_name4}"
	preserved_state {
		metadata = {
			foo = "abc"
		}
	}
}

resource "google_compute_disk" "disk" {
  name  = "test-disk-%{random_suffix}"
  type  = "pd-ssd"
  zone  = "us-central1-c"
  image = "debian-8-jessie-v20170523"
  physical_block_size_bytes = 4096
}

resource "google_compute_disk" "disk1" {
  name  = "test-disk2-%{random_suffix}"
  type  = "pd-ssd"
  zone  = "us-central1-c"
  image = "debian-cloud/debian-11"
  physical_block_size_bytes = 4096
}

resource "google_compute_disk" "disk2" {
  name  = "test-disk3-%{random_suffix}"
  type  = "pd-ssd"
  zone  = "us-central1-c"
  image = "https://www.googleapis.com/compute/v1/projects/centos-cloud/global/images/centos-7-v20210217"
  physical_block_size_bytes = 4096
}
`, context) + testAccComputeRegionPerInstanceConfig_rigm(context)
}

func testAccComputeRegionPerInstanceConfig_rigm(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance_template" "rigm-basic" {
  name           = "tf-test-rigm-%{random_suffix}"
  machine_type   = "e2-medium"
  can_ip_forward = false
  tags           = ["foo", "bar"]

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    boot         = true
    device_name  = "my-stateful-disk"
  }

  network_interface {
    network = "default"
  }

  service_account {
    scopes = ["userinfo-email", "compute-ro", "storage-ro"]
  }
}

resource "google_compute_region_instance_group_manager" "rigm" {
  description = "Terraform test instance group manager"
  name        = "%{rigm_name}"

  version {
    name              = "prod"
    instance_template = google_compute_instance_template.rigm-basic.self_link
  }

  base_instance_name = "tf-test-rigm-no-tp"

  update_policy {
    instance_redistribution_type = "NONE"
    type                         = "OPPORTUNISTIC"
    minimal_action               = "REPLACE"
    max_surge_fixed              = 0
    max_unavailable_fixed        = 6
  }
}
`, context)
}

// Checks that the per instance config with the given name was destroyed
func testAccCheckComputeRegionPerInstanceConfigDestroyed(t *testing.T, rigmId, configName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		foundNames, err := testAccComputePerInstanceConfigListNames(t, rigmId)
		if err != nil {
			return fmt.Errorf("unable to confirm config with name %s was destroyed: %v", configName, err)
		}
		if _, ok := foundNames[configName]; ok {
			return fmt.Errorf("config with name %s still exists", configName)
		}

		return nil
	}
}
