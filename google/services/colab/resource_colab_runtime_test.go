// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package colab_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccColabRuntime_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"key_name":      acctest.BootstrapKMSKeyInLocation(t, "us-central1").CryptoKey.Name,
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckColabRuntimeDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccColabRuntime_full(context),
			},
			{
				ResourceName:            "google_colab_runtime.runtime",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"desired_state", "location", "name", "auto_upgrade"},
			},
			{
				Config: testAccColabRuntime_no_state(context),
			},
			{
				ResourceName:            "google_colab_runtime.runtime",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"desired_state", "location", "name", "auto_upgrade"},
			},
			{
				Config: testAccColabRuntime_stopped(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_colab_runtime.runtime", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_colab_runtime.runtime",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"desired_state", "location", "name", "auto_upgrade"},
			},
			{
				Config: testAccColabRuntime_full(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_colab_runtime.runtime", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_colab_runtime.runtime",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"desired_state", "location", "name", "auto_upgrade"},
			},
		},
	})
}

func testAccColabRuntime_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_colab_runtime_template" "my_template" {
  name        = "tf-test-colab-runtime%{random_suffix}"
  display_name = "Runtime template full"
  location    = "us-central1"
  description = "Full runtime template"
  machine_spec {
    machine_type     = "n1-standard-2"
    accelerator_type = "NVIDIA_TESLA_T4"
    accelerator_count = "1"
  }

  data_persistent_disk_spec {
    disk_type    = "pd-standard"
    disk_size_gb = 200
  }

  network_spec {
    enable_internet_access = true
  }

  labels = {
    k = "val"
  }

  idle_shutdown_config {
    idle_timeout = "3600s"
  }

  euc_config {
    euc_disabled = true
  }

  shielded_vm_config {
    enable_secure_boot = true
  }

  network_tags = ["abc", "def"]

  encryption_spec {
    kms_key_name = "%{key_name}"
  }
}

resource "google_colab_runtime" "runtime" {
  name = "tf-test-colab-runtime%{random_suffix}"
  location = "us-central1" 
  
  notebook_runtime_template_ref {
    notebook_runtime_template = google_colab_runtime_template.my_template.id
  }
  
  display_name = "Runtime full"
  runtime_user = "gterraformtestuser@gmail.com"
  description = "Full runtime"

  desired_state = "RUNNING"

  auto_upgrade = true

  depends_on = [
    google_colab_runtime_template.my_template
  ]
}
`, context)
}

func testAccColabRuntime_stopped(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_colab_runtime_template" "my_template" {
  name        = "tf-test-colab-runtime%{random_suffix}"
  display_name = "Runtime template full"
  location    = "us-central1"
  description = "Full runtime template"
  machine_spec {
    machine_type     = "n1-standard-2"
    accelerator_type = "NVIDIA_TESLA_T4"
    accelerator_count = "1"
  }

  data_persistent_disk_spec {
    disk_type    = "pd-standard"
    disk_size_gb = 200
  }

  network_spec {
    enable_internet_access = true
  }

  labels = {
    k = "val"
  }

  idle_shutdown_config {
    idle_timeout = "3600s"
  }

  euc_config {
    euc_disabled = true
  }

  shielded_vm_config {
    enable_secure_boot = true
  }

  network_tags = ["abc", "def"]

  encryption_spec {
    kms_key_name = "%{key_name}"
  }
}

resource "google_colab_runtime" "runtime" {
  name = "tf-test-colab-runtime%{random_suffix}"
  location = "us-central1" 
  
  notebook_runtime_template_ref {
    notebook_runtime_template = google_colab_runtime_template.my_template.id
  }
  
  display_name = "Runtime full"
  runtime_user = "gterraformtestuser@gmail.com"
  description = "Full runtime"

  desired_state = "STOPPED"

  depends_on = [
    google_colab_runtime_template.my_template
  ]
}
`, context)
}

func testAccColabRuntime_no_state(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_colab_runtime_template" "my_template" {
  name        = "tf-test-colab-runtime%{random_suffix}"
  display_name = "Runtime template full"
  location    = "us-central1"
  description = "Full runtime template"
  machine_spec {
    machine_type     = "n1-standard-2"
    accelerator_type = "NVIDIA_TESLA_T4"
    accelerator_count = "1"
  }

  data_persistent_disk_spec {
    disk_type    = "pd-standard"
    disk_size_gb = 200
  }

  network_spec {
    enable_internet_access = true
  }

  labels = {
    k = "val"
  }

  idle_shutdown_config {
    idle_timeout = "3600s"
  }

  euc_config {
    euc_disabled = true
  }

  shielded_vm_config {
    enable_secure_boot = true
  }

  network_tags = ["abc", "def"]

  encryption_spec {
    kms_key_name = "%{key_name}"
  }
}

resource "google_colab_runtime" "runtime" {
  name = "tf-test-colab-runtime%{random_suffix}"
  location = "us-central1" 
  
  notebook_runtime_template_ref {
    notebook_runtime_template = google_colab_runtime_template.my_template.id
  }
  
  display_name = "Runtime full"
  runtime_user = "gterraformtestuser@gmail.com"
  description = "Full runtime"

  depends_on = [
    google_colab_runtime_template.my_template
  ]
}
`, context)
}
