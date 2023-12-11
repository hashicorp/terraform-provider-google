// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package workbench_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

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
			},
			{
				ResourceName:            "google_workbench_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "instance_owners", "location", "instance_id", "request_id", "labels", "terraform_labels"},
			},
			{
				Config: testAccWorkbenchInstance_update(context),
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
			},
			{
				ResourceName:            "google_workbench_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "instance_owners", "location", "instance_id", "request_id", "labels", "terraform_labels"},
			},
			{
				Config: testAccWorkbenchInstance_updateGpu(context),
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
			},
			{
				ResourceName:            "google_workbench_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "instance_owners", "location", "instance_id", "request_id", "labels", "terraform_labels"},
			},
			{
				Config: testAccWorkbenchInstance_updateGpu(context),
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
			},
			{
				ResourceName:            "google_workbench_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "instance_owners", "location", "instance_id", "request_id", "labels", "terraform_labels"},
			},
			{
				Config: testAccWorkbenchInstance_updateMetadata(context),
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

func testAccWorkbenchInstance_updateMetadata(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_workbench_instance" "instance" {
  name = "tf-test-workbench-instance%{random_suffix}"
  location = "us-central1-a"

  gce_setup {
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
