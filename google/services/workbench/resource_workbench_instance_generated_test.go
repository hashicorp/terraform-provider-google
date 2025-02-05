// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package workbench_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccWorkbenchInstance_workbenchInstanceBasicExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckWorkbenchInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkbenchInstance_workbenchInstanceBasicExample(context),
			},
			{
				ResourceName:            "google_workbench_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_id", "instance_owners", "labels", "location", "name", "terraform_labels"},
			},
		},
	})
}

func testAccWorkbenchInstance_workbenchInstanceBasicExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_workbench_instance" "instance" {
  name = "tf-test-workbench-instance%{random_suffix}"
  location = "us-west1-a"
}
`, context)
}

func TestAccWorkbenchInstance_workbenchInstanceBasicContainerExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckWorkbenchInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkbenchInstance_workbenchInstanceBasicContainerExample(context),
			},
			{
				ResourceName:            "google_workbench_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_id", "instance_owners", "labels", "location", "name", "terraform_labels"},
			},
		},
	})
}

func testAccWorkbenchInstance_workbenchInstanceBasicContainerExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_workbench_instance" "instance" {
  name = "tf-test-workbench-instance%{random_suffix}"
  location = "us-west1-a"

  gce_setup {
    container_image {
      repository = "us-docker.pkg.dev/deeplearning-platform-release/gcr.io/base-cu113.py310"
      tag = "latest"
    }
  }
}
`, context)
}

func TestAccWorkbenchInstance_workbenchInstanceBasicGpuExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckWorkbenchInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkbenchInstance_workbenchInstanceBasicGpuExample(context),
			},
			{
				ResourceName:            "google_workbench_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"gce_setup.0.vm_image", "instance_id", "instance_owners", "labels", "location", "name", "terraform_labels"},
			},
		},
	})
}

func testAccWorkbenchInstance_workbenchInstanceBasicGpuExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_workbench_instance" "instance" {
  name = "tf-test-workbench-instance%{random_suffix}"
  location = "us-central1-a"
  gce_setup {
    machine_type = "n1-standard-1" // cant be e2 because of accelerator
    accelerator_configs {
      type         = "NVIDIA_TESLA_T4"
      core_count   = 1
    }
    vm_image {
      project      = "cloud-notebooks-managed"
      family       = "workbench-instances"
    }
  }
}
`, context)
}

func TestAccWorkbenchInstance_workbenchInstanceLabelsStoppedExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"service_account": envvar.GetTestServiceAccountFromEnv(t),
		"random_suffix":   acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckWorkbenchInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkbenchInstance_workbenchInstanceLabelsStoppedExample(context),
			},
			{
				ResourceName:            "google_workbench_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"desired_state", "instance_id", "instance_owners", "labels", "location", "name", "terraform_labels"},
			},
		},
	})
}

func testAccWorkbenchInstance_workbenchInstanceLabelsStoppedExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_workbench_instance" "instance" {
  name = "tf-test-workbench-instance%{random_suffix}"
  location = "us-central1-a"

  gce_setup {
    machine_type = "e2-standard-4"

    shielded_instance_config {
      enable_secure_boot = false
      enable_vtpm = false
      enable_integrity_monitoring = false
    }

    service_accounts {
      email = "%{service_account}"
    }

    metadata = {
      terraform = "true"
    }

  }

  labels = {
    k = "val"
  }

  desired_state = "STOPPED"

}
`, context)
}

func TestAccWorkbenchInstance_workbenchInstanceFullExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_id":      envvar.GetTestProjectFromEnv(),
		"service_account": envvar.GetTestServiceAccountFromEnv(t),
		"key_name":        acctest.BootstrapKMSKeyInLocation(t, "us-central1").CryptoKey.Name,
		"random_suffix":   acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckWorkbenchInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkbenchInstance_workbenchInstanceFullExample(context),
			},
			{
				ResourceName:            "google_workbench_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"gce_setup.0.boot_disk.0.disk_type", "gce_setup.0.data_disks.0.disk_type", "gce_setup.0.vm_image", "instance_id", "instance_owners", "labels", "location", "name", "terraform_labels"},
			},
		},
	})
}

func testAccWorkbenchInstance_workbenchInstanceFullExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "my_network" {
  name = "tf-test-wbi-test-default%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "my_subnetwork" {
  name   = "tf-test-wbi-test-default%{random_suffix}"
  network = google_compute_network.my_network.id
  region = "us-central1"
  ip_cidr_range = "10.0.1.0/24"
}

resource "google_compute_address" "static" {
  name = "tf-test-wbi-test-default%{random_suffix}"
}

resource "google_service_account_iam_binding" "act_as_permission" {
  service_account_id = "projects/%{project_id}/serviceAccounts/%{service_account}"
  role               = "roles/iam.serviceAccountUser"
  members = [
    "user:example@example.com",
  ]
}

resource "google_workbench_instance" "instance" {
  name = "tf-test-workbench-instance%{random_suffix}"
  location = "us-central1-a"

  gce_setup {
    machine_type = "n1-standard-4" // cant be e2 because of accelerator
    accelerator_configs {
      type         = "NVIDIA_TESLA_T4"
      core_count   = 1
    }

    shielded_instance_config {
      enable_secure_boot = true
      enable_vtpm = true
      enable_integrity_monitoring = true
    }

    disable_public_ip = false

    service_accounts {
      email = "%{service_account}"
    }

    boot_disk {
      disk_size_gb  = 310
      disk_type = "PD_SSD"
      disk_encryption = "CMEK"
      kms_key = "%{key_name}"
    }

    data_disks {
      disk_size_gb  = 330
      disk_type = "PD_SSD"
      disk_encryption = "CMEK"
      kms_key = "%{key_name}"
    }

    network_interfaces {
      network = google_compute_network.my_network.id
      subnet = google_compute_subnetwork.my_subnetwork.id
      nic_type = "GVNIC"
      access_configs {
        external_ip = google_compute_address.static.address
      }
    }

    metadata = {
      terraform = "true"
    }

    enable_ip_forwarding = true

    tags = ["abc", "def"]

  }

  disable_proxy_access = "true"

  instance_owners  = ["example@example.com"]

  labels = {
    k = "val"
  }

  desired_state = "ACTIVE"

  enable_third_party_identity = "true"

}
`, context)
}

func testAccCheckWorkbenchInstanceDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_workbench_instance" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{WorkbenchBasePath}}projects/{{project}}/locations/{{location}}/instances/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("WorkbenchInstance still exists at %s", url)
			}
		}

		return nil
	}
}
