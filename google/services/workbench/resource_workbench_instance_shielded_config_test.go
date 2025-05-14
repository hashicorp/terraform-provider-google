// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package workbench_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccWorkbenchInstance_shielded_config_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkbenchInstance_shielded_config_false(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_workbench_instance.instance", "state", "ACTIVE"),
				),
			},
			{
				ResourceName:            "google_workbench_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "instance_owners", "location", "instance_id", "request_id", "labels", "terraform_labels", "desired_state", "update_time", "health_info", "health_state"},
			},
			{
				Config: testAccWorkbenchInstance_shielded_config_true(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_workbench_instance.instance", "state", "ACTIVE"),
				),
			},
			{
				ResourceName:            "google_workbench_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "instance_owners", "location", "instance_id", "request_id", "labels", "terraform_labels", "desired_state", "update_time", "health_info", "health_state"},
			},
		},
	})
}

func TestAccWorkbenchInstance_shielded_config_remove(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkbenchInstance_shielded_config_true(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_workbench_instance.instance", "state", "ACTIVE"),
				),
			},
			{
				ResourceName:            "google_workbench_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "instance_owners", "location", "instance_id", "request_id", "labels", "terraform_labels", "desired_state", "update_time", "health_info", "health_state"},
			},
			{
				Config: testAccWorkbenchInstance_shielded_config_none(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_workbench_instance.instance", "state", "ACTIVE"),
				),
			},
			{
				ResourceName:            "google_workbench_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "instance_owners", "location", "instance_id", "request_id", "labels", "terraform_labels", "desired_state", "update_time", "health_info", "health_state"},
			},
		},
	})
}

func TestAccWorkbenchInstance_shielded_config_double_apply(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkbenchInstance_shielded_config_none(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_workbench_instance.instance", "state", "ACTIVE"),
				),
			},
			{
				ResourceName:            "google_workbench_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "instance_owners", "location", "instance_id", "request_id", "labels", "terraform_labels", "desired_state", "update_time", "health_info", "health_state"},
			},
			{
				Config: testAccWorkbenchInstance_shielded_config_none(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_workbench_instance.instance", "state", "ACTIVE"),
				),
			},
			{
				ResourceName:            "google_workbench_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "instance_owners", "location", "instance_id", "request_id", "labels", "terraform_labels", "desired_state", "update_time", "health_info", "health_state"},
			},
			{
				Config: testAccWorkbenchInstance_shielded_config_false(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_workbench_instance.instance", "state", "ACTIVE"),
				),
			},
			{
				ResourceName:            "google_workbench_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "instance_owners", "location", "instance_id", "request_id", "labels", "terraform_labels", "desired_state", "update_time", "health_info", "health_state"},
			},
			{
				Config: testAccWorkbenchInstance_shielded_config_false(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_workbench_instance.instance", "state", "ACTIVE"),
				),
			},
			{
				ResourceName:            "google_workbench_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "instance_owners", "location", "instance_id", "request_id", "labels", "terraform_labels", "desired_state", "update_time", "health_info", "health_state"},
			},
			{
				Config: testAccWorkbenchInstance_shielded_config_true(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_workbench_instance.instance", "state", "ACTIVE"),
				),
			},
			{
				ResourceName:            "google_workbench_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "instance_owners", "location", "instance_id", "request_id", "labels", "terraform_labels", "desired_state", "update_time", "health_info", "health_state"},
			},
			{
				Config: testAccWorkbenchInstance_shielded_config_true(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_workbench_instance.instance", "state", "ACTIVE"),
				),
			},
			{
				ResourceName:            "google_workbench_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "instance_owners", "location", "instance_id", "request_id", "labels", "terraform_labels", "desired_state", "update_time", "health_info", "health_state"},
			},
		},
	})
}

func testAccWorkbenchInstance_shielded_config_true(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_workbench_instance" "instance" {
  name = "tf-test-workbench-instance%{random_suffix}"
  location = "us-central1-a"

  gce_setup {
    shielded_instance_config {
      enable_secure_boot = true
      enable_vtpm = true
      enable_integrity_monitoring = true
    }
  }
}
`, context)
}

func testAccWorkbenchInstance_shielded_config_false(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_workbench_instance" "instance" {
  name = "tf-test-workbench-instance%{random_suffix}"
  location = "us-central1-a"

  gce_setup {
    shielded_instance_config {
      enable_secure_boot = false
      enable_vtpm = false
      enable_integrity_monitoring = false
    }
  }

}
`, context)
}

func testAccWorkbenchInstance_shielded_config_none(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_workbench_instance" "instance" {
  name = "tf-test-workbench-instance%{random_suffix}"
  location = "us-central1-a"
}
`, context)
}
