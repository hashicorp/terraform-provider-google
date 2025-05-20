// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package workbench_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccWorkbenchInstance_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkbenchInstance_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_workbench_instance.instance", "state", "ACTIVE"),
				),
			},
			{
				ResourceName:            "google_workbench_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "instance_owners", "location", "instance_id", "request_id", "labels", "terraform_labels"},
			},
			{
				Config: testAccWorkbenchInstance_update(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_workbench_instance.instance", "state", "ACTIVE"),
				),
			},
			{
				ResourceName:            "google_workbench_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "instance_owners", "location", "instance_id", "request_id", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccWorkbenchInstance_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_workbench_instance" "instance" {
  name = "tf-test-workbench-instance%{random_suffix}"
  location = "us-central1-a"
}
`, context)
}

func testAccWorkbenchInstance_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_workbench_instance" "instance" {
  name = "tf-test-workbench-instance%{random_suffix}"
  location = "us-central1-a"

  gce_setup {
    machine_type = "n1-standard-16"

    accelerator_configs{
      type         = "NVIDIA_TESLA_T4"
      core_count   = 1
    }

    shielded_instance_config {
      enable_secure_boot = false
      enable_vtpm = true
      enable_integrity_monitoring = false
    }

	boot_disk {
		disk_size_gb  = 310
	  }
  
	  data_disks {
		disk_size_gb  = 330
	  }

    metadata = {
      terraform = "true"
    }

  }

  labels = {
    k = "val"
  }

}
`, context)
}

func TestAccWorkbenchInstance_updateGpu(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkbenchInstance_basicGpu(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_workbench_instance.instance", "state", "ACTIVE"),
				),
			},
			{
				ResourceName:            "google_workbench_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "instance_owners", "location", "instance_id", "request_id", "labels", "terraform_labels"},
			},
			{
				Config: testAccWorkbenchInstance_updateGpu(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_workbench_instance.instance", "state", "ACTIVE"),
				),
			},
			{
				ResourceName:            "google_workbench_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "instance_owners", "location", "instance_id", "request_id", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccWorkbenchInstance_basicGpu(context map[string]interface{}) string {
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

  }
}
`, context)
}

func testAccWorkbenchInstance_updateGpu(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_workbench_instance" "instance" {
  name = "tf-test-workbench-instance%{random_suffix}"
  location = "us-central1-a"

  gce_setup {
    machine_type = "n1-standard-16"

    accelerator_configs{
      type         = "NVIDIA_TESLA_P4"
      core_count   = 1
    }

    shielded_instance_config {
      enable_secure_boot = false
      enable_vtpm = true
      enable_integrity_monitoring = false
    }

  }

  labels = {
    k = "val"
  }

}
`, context)
}

func TestAccWorkbenchInstance_removeGpu(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkbenchInstance_Gpu(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_workbench_instance.instance", "state", "ACTIVE"),
				),
			},
			{
				ResourceName:            "google_workbench_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "instance_owners", "location", "instance_id", "request_id", "labels", "terraform_labels"},
			},
			{
				Config: testAccWorkbenchInstance_removeGpu(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_workbench_instance.instance", "state", "ACTIVE"),
				),
			},
			{
				ResourceName:            "google_workbench_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "instance_owners", "location", "instance_id", "request_id", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccWorkbenchInstance_Gpu(context map[string]interface{}) string {
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

  }
}
`, context)
}

func testAccWorkbenchInstance_removeGpu(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_workbench_instance" "instance" {
  name = "tf-test-workbench-instance%{random_suffix}"
  location = "us-central1-a"

  gce_setup {
    machine_type = "n1-standard-16"

  }

}
`, context)
}

func TestAccWorkbenchInstance_updateMetadata(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkbenchInstance_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_workbench_instance.instance", "state", "ACTIVE"),
				),
			},
			{
				ResourceName:            "google_workbench_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "instance_owners", "location", "instance_id", "request_id", "labels", "terraform_labels", "update_time"},
			},
			{
				Config: testAccWorkbenchInstance_updateMetadata(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_workbench_instance.instance", "state", "ACTIVE"),
				),
			},
			{
				ResourceName:            "google_workbench_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "instance_owners", "location", "instance_id", "request_id", "labels", "terraform_labels", "update_time"},
			},
			{
				Config: testAccWorkbenchInstance_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_workbench_instance.instance", "state", "ACTIVE"),
				),
			},
			{
				ResourceName:            "google_workbench_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "instance_owners", "location", "instance_id", "request_id", "labels", "terraform_labels", "update_time"},
			},
		},
	})
}

func TestAccWorkbenchInstance_updateMetadataKey(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkbenchInstance_updateMetadata(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_workbench_instance.instance", "state", "ACTIVE"),
				),
			},
			{
				ResourceName:            "google_workbench_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "instance_owners", "location", "instance_id", "request_id", "labels", "terraform_labels", "update_time", "health_info", "health_state"},
			},
			{
				Config: testAccWorkbenchInstance_updateMetadataKey(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_workbench_instance.instance", "state", "ACTIVE"),
				),
			},
			{
				ResourceName:            "google_workbench_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "instance_owners", "location", "instance_id", "request_id", "labels", "terraform_labels", "update_time", "health_info", "health_state"},
			},
			{
				Config: testAccWorkbenchInstance_updateMetadata(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_workbench_instance.instance", "state", "ACTIVE"),
				),
			},
			{
				ResourceName:            "google_workbench_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "instance_owners", "location", "instance_id", "request_id", "labels", "terraform_labels", "update_time", "health_info", "health_state"},
			},
		},
	})
}

func testAccWorkbenchInstance_updateMetadata(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_workbench_instance" "instance" {
  name = "tf-test-workbench-instance%{random_suffix}"
  location = "us-central1-a"

  gce_setup {
    metadata = {
      terraform = "true"
      "resource-url" = "new-fake-value",
    }
  }

  labels = {
    k = "val"
  }

}
`, context)
}

func testAccWorkbenchInstance_updateMetadataKey(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_workbench_instance" "instance" {
  name = "tf-test-workbench-instance%{random_suffix}"
  location = "us-central1-a"

  gce_setup {
    metadata = {
      terraform = "true",
      "idle-timeout-seconds" = "10800",
      "image-url" = "fake-value",
    }
  }

  labels = {
    k = "val"
  }

}
`, context)
}

func TestAccWorkbenchInstance_updateState(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkbenchInstance_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_workbench_instance.instance", "state", "ACTIVE"),
				),
			},
			{
				ResourceName:            "google_workbench_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "instance_owners", "location", "instance_id", "request_id", "labels", "terraform_labels", "desired_state", "update_time"},
			},
			{
				Config: testAccWorkbenchInstance_updateState(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_workbench_instance.instance", "state", "STOPPED"),
				),
			},
			{
				ResourceName:            "google_workbench_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "instance_owners", "location", "instance_id", "request_id", "labels", "terraform_labels", "desired_state", "update_time"},
			},
			{
				Config: testAccWorkbenchInstance_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_workbench_instance.instance", "state", "ACTIVE"),
				),
			},
			{
				ResourceName:            "google_workbench_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "instance_owners", "location", "instance_id", "request_id", "labels", "terraform_labels", "desired_state", "update_time"},
			},
		},
	})
}

func testAccWorkbenchInstance_updateState(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_workbench_instance" "instance" {
  name = "tf-test-workbench-instance%{random_suffix}"
  location = "us-central1-a"

  desired_state = "STOPPED"

}
`, context)
}

func TestAccWorkbenchInstance_empty_accelerator(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkbenchInstance_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_workbench_instance.instance", "state", "ACTIVE"),
				),
			},
			{
				ResourceName:            "google_workbench_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "instance_owners", "location", "instance_id", "request_id", "labels", "terraform_labels"},
			},
			{
				Config: testAccWorkbenchInstance_empty_accelerator(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_workbench_instance.instance", "state", "ACTIVE"),
				),
			},
			{
				ResourceName:            "google_workbench_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "instance_owners", "location", "instance_id", "request_id", "labels", "terraform_labels"},
			},
			{
				Config: testAccWorkbenchInstance_empty_accelerator(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_workbench_instance.instance", "state", "ACTIVE"),
				),
			},
			{
				ResourceName:            "google_workbench_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "instance_owners", "location", "instance_id", "request_id", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccWorkbenchInstance_empty_accelerator(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_workbench_instance" "instance" {
  name = "tf-test-workbench-instance%{random_suffix}"
  location = "us-central1-a"

  gce_setup {
    accelerator_configs{
    }
  }
}
`, context)
}

func TestAccWorkbenchInstance_updateBootDisk(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkbenchInstance_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_workbench_instance.instance", "state", "ACTIVE"),
				),
			},
			{
				ResourceName:            "google_workbench_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "instance_owners", "location", "instance_id", "request_id", "labels", "terraform_labels", "desired_state"},
			},
			{
				Config: testAccWorkbenchInstance_updateBootDisk(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_workbench_instance.instance", "state", "ACTIVE"),
				),
			},
			{
				ResourceName:            "google_workbench_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "instance_owners", "location", "instance_id", "request_id", "labels", "terraform_labels", "desired_state"},
			},
		},
	})
}

func TestAccWorkbenchInstance_updateDataDisk(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkbenchInstance_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_workbench_instance.instance", "state", "ACTIVE"),
				),
			},
			{
				ResourceName:            "google_workbench_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "instance_owners", "location", "instance_id", "request_id", "labels", "terraform_labels", "desired_state"},
			},
			{
				Config: testAccWorkbenchInstance_updateDataDisk(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_workbench_instance.instance", "state", "ACTIVE"),
				),
			},
			{
				ResourceName:            "google_workbench_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "instance_owners", "location", "instance_id", "request_id", "labels", "terraform_labels", "desired_state"},
			},
		},
	})
}

func TestAccWorkbenchInstance_updateBothDisks(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkbenchInstance_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_workbench_instance.instance", "state", "ACTIVE"),
				),
			},
			{
				ResourceName:            "google_workbench_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "instance_owners", "location", "instance_id", "request_id", "labels", "terraform_labels", "desired_state"},
			},
			{
				Config: testAccWorkbenchInstance_updateBothDisks(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_workbench_instance.instance", "state", "ACTIVE"),
				),
			},
			{
				ResourceName:            "google_workbench_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "instance_owners", "location", "instance_id", "request_id", "labels", "terraform_labels", "desired_state"},
			},
		},
	})
}

func testAccWorkbenchInstance_updateBootDisk(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_workbench_instance" "instance" {
  name = "tf-test-workbench-instance%{random_suffix}"
  location = "us-central1-a"
  gce_setup {
	boot_disk {
		disk_size_gb  = 310
	  }
	}
}
`, context)
}

func testAccWorkbenchInstance_updateDataDisk(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_workbench_instance" "instance" {
  name = "tf-test-workbench-instance%{random_suffix}"
  location = "us-central1-a"
  gce_setup {  
	  data_disks {
		disk_size_gb  = 330
	  }
	}
}
`, context)
}

func testAccWorkbenchInstance_updateBothDisks(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_workbench_instance" "instance" {
  name = "tf-test-workbench-instance%{random_suffix}"
  location = "us-central1-a"
  gce_setup {
	boot_disk {
		disk_size_gb  = 310
	  }

	  data_disks {
		disk_size_gb  = 330
	  }
	}
}
`, context)
}

func TestAccWorkbenchInstance_updatelabels(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkbenchInstance_label(context),
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
				Config: testAccWorkbenchInstance_basic(context),
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
				Config: testAccWorkbenchInstance_label(context),
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

func testAccWorkbenchInstance_label(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_workbench_instance" "instance" {
  name = "tf-test-workbench-instance%{random_suffix}"
  location = "us-central1-a"
  labels = {
    k = "val"
  }
}
`, context)
}

func TestAccWorkbenchInstance_updateCustomContainers(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkbenchInstance_customcontainer(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_workbench_instance.instance", "state", "ACTIVE"),
				),
			},
			{
				ResourceName:            "google_workbench_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "instance_owners", "location", "instance_id", "request_id", "labels", "terraform_labels", "desired_state"},
			},
			{
				Config: testAccWorkbenchInstance_updatedcustomcontainer(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_workbench_instance.instance", "state", "ACTIVE"),
				),
			},
			{
				ResourceName:            "google_workbench_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "instance_owners", "location", "instance_id", "request_id", "labels", "terraform_labels", "desired_state"},
			},
		},
	})
}

func testAccWorkbenchInstance_customcontainer(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_workbench_instance" "instance" {
  name = "tf-test-workbench-instance%{random_suffix}"
  location = "us-central1-a"
  gce_setup {
    container_image {
      repository = "us-docker.pkg.dev/deeplearning-platform-release/gcr.io/base-cu113.py310"
      tag = "latest"
    }
  }
}
`, context)
}

func testAccWorkbenchInstance_updatedcustomcontainer(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_workbench_instance" "instance" {
  name = "tf-test-workbench-instance%{random_suffix}"
  location = "us-central1-a"
  gce_setup {
    container_image {
      repository = "gcr.io/deeplearning-platform-release/workbench-container"
      tag = "20241117-2200-rc0"
    }
  }
}
`, context)
}
