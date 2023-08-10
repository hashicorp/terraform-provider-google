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
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccComputePerInstanceConfig_statefulBasic(t *testing.T) {
	// Multiple fine-grained resources
	acctest.SkipIfVcr(t)
	t.Parallel()

	suffix := acctest.RandString(t, 10)
	igmName := fmt.Sprintf("tf-test-igm-%s", suffix)
	context := map[string]interface{}{
		"igm_name":      igmName,
		"random_suffix": suffix,
		"config_name":   fmt.Sprintf("instance-%s", acctest.RandString(t, 10)),
		"config_name2":  fmt.Sprintf("instance-%s", acctest.RandString(t, 10)),
		"config_name3":  fmt.Sprintf("instance-%s", acctest.RandString(t, 10)),
		"config_name4":  fmt.Sprintf("instance-%s", acctest.RandString(t, 10)),
	}
	igmId := fmt.Sprintf("projects/%s/zones/%s/instanceGroupManagers/%s",
		envvar.GetTestProjectFromEnv(), envvar.GetTestZoneFromEnv(), igmName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Create one endpoint
				Config: testAccComputePerInstanceConfig_statefulBasic(context),
			},
			{
				ResourceName:            "google_compute_per_instance_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_instance_state_on_destroy", "zone"},
			},
			{
				// Force-recreate old config
				Config: testAccComputePerInstanceConfig_statefulModified(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputePerInstanceConfigDestroyed(t, igmId, context["config_name"].(string)),
				),
			},
			{
				ResourceName:            "google_compute_per_instance_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_instance_state_on_destroy", "zone"},
			},
			{
				// Add two new endpoints
				Config: testAccComputePerInstanceConfig_statefulAdditional(context),
			},
			{
				ResourceName:            "google_compute_per_instance_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_instance_state_on_destroy", "zone"},
			},
			{
				ResourceName:            "google_compute_per_instance_config.with_disks",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"most_disruptive_allowed_action", "minimal_action", "remove_instance_state_on_destroy"},
			},
			{
				ResourceName:            "google_compute_per_instance_config.add2",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_instance_state_on_destroy", "zone"},
			},
			{
				// delete all configs
				Config: testAccComputePerInstanceConfig_igm(context),
				Check: resource.ComposeTestCheckFunc(
					// Config with remove_instance_state_on_destroy = false won't be destroyed (config4)
					testAccCheckComputePerInstanceConfigDestroyed(t, igmId, context["config_name2"].(string)),
					testAccCheckComputePerInstanceConfigDestroyed(t, igmId, context["config_name3"].(string)),
				),
			},
		},
	})
}

func TestAccComputePerInstanceConfig_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"igm_name":      fmt.Sprintf("tf-test-igm-%s", acctest.RandString(t, 10)),
		"config_name":   fmt.Sprintf("instance-%s", acctest.RandString(t, 10)),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Create one config
				Config: testAccComputePerInstanceConfig_statefulBasic(context),
			},
			{
				ResourceName:            "google_compute_per_instance_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_instance_state_on_destroy", "zone"},
			},
			{
				// Update an existing config
				Config: testAccComputePerInstanceConfig_update(context),
			},
			{
				ResourceName:            "google_compute_per_instance_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_instance_state_on_destroy", "zone"},
			},
		},
	})
}

func testAccComputePerInstanceConfig_statefulBasic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_per_instance_config" "default" {
	instance_group_manager = google_compute_instance_group_manager.igm.name
	name = "%{config_name}"
	remove_instance_state_on_destroy = true
	preserved_state {
		metadata = {
			asdf = "asdf"
		}
	}
}
`, context) + testAccComputePerInstanceConfig_igm(context)
}

func testAccComputePerInstanceConfig_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_per_instance_config" "default" {
	instance_group_manager = google_compute_instance_group_manager.igm.name
	name = "%{config_name}"
	remove_instance_state_on_destroy = true
	preserved_state {
		metadata = {
			asdf = "asdf"
			update = "12345"
		}
	}
}
`, context) + testAccComputePerInstanceConfig_igm(context)
}

func testAccComputePerInstanceConfig_statefulModified(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_per_instance_config" "default" {
	zone = google_compute_instance_group_manager.igm.zone
	instance_group_manager = google_compute_instance_group_manager.igm.name
	name = "%{config_name2}"
	remove_instance_state_on_destroy = true
	preserved_state {
		metadata = {
			asdf = "asdf"
		}
	}
}
`, context) + testAccComputePerInstanceConfig_igm(context)
}

func testAccComputePerInstanceConfig_statefulAdditional(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_per_instance_config" "default" {
	zone = google_compute_instance_group_manager.igm.zone
	instance_group_manager = google_compute_instance_group_manager.igm.name
	name = "%{config_name2}"
	remove_instance_state_on_destroy = true
	preserved_state {
		metadata = {
			asdf = "asdf"
		}
	}
}

resource "google_compute_per_instance_config" "with_disks" {
	zone = google_compute_instance_group_manager.igm.zone
	instance_group_manager = google_compute_instance_group_manager.igm.name
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

resource "google_compute_per_instance_config" "add2" {
	zone = google_compute_instance_group_manager.igm.zone
	instance_group_manager = google_compute_instance_group_manager.igm.name
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
  zone  = google_compute_instance_group_manager.igm.zone
  image = "debian-8-jessie-v20170523"
  physical_block_size_bytes = 4096
}

resource "google_compute_disk" "disk1" {
  name  = "test-disk2-%{random_suffix}"
  type  = "pd-ssd"
  zone  = google_compute_instance_group_manager.igm.zone
  image = "debian-cloud/debian-11"
  physical_block_size_bytes = 4096
}

resource "google_compute_disk" "disk2" {
  name  = "test-disk3-%{random_suffix}"
  type  = "pd-ssd"
  zone  = google_compute_instance_group_manager.igm.zone
  image = "https://www.googleapis.com/compute/v1/projects/centos-cloud/global/images/centos-7-v20210217"
  physical_block_size_bytes = 4096
}
`, context) + testAccComputePerInstanceConfig_igm(context)
}

func testAccComputePerInstanceConfig_igm(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance_template" "igm-basic" {
  name           = "tf-test-igm-%{random_suffix}"
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

resource "google_compute_instance_group_manager" "igm" {
  description = "Terraform test instance group manager"
  name        = "%{igm_name}"

  version {
    name              = "prod"
    instance_template = google_compute_instance_template.igm-basic.self_link
  }

  base_instance_name = "tf-test-igm-no-tp"
}
`, context)
}

// Checks that the per instance config with the given name was destroyed
func testAccCheckComputePerInstanceConfigDestroyed(t *testing.T, igmId, configName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		foundNames, err := testAccComputePerInstanceConfigListNames(t, igmId)
		if err != nil {
			return fmt.Errorf("unable to confirm config with name %s was destroyed: %v", configName, err)
		}
		if _, ok := foundNames[configName]; ok {
			return fmt.Errorf("config with name %s still exists", configName)
		}

		return nil
	}
}

func testAccComputePerInstanceConfigListNames(t *testing.T, igmId string) (map[string]struct{}, error) {
	config := acctest.GoogleProviderConfig(t)

	url := fmt.Sprintf("%s%s/listPerInstanceConfigs", config.ComputeBasePath, igmId)
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "POST",
		RawURL:    url,
		UserAgent: config.UserAgent,
	})
	if err != nil {
		return nil, err
	}

	v, ok := res["items"]
	if !ok || v == nil {
		return nil, nil
	}
	items := v.([]interface{})
	instanceConfigs := make(map[string]struct{})
	for _, item := range items {
		perInstanceConfig := item.(map[string]interface{})
		instanceConfigs[fmt.Sprintf("%v", perInstanceConfig["name"])] = struct{}{}
	}
	return instanceConfigs, nil
}
