// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package healthcare_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccHealthcareWorkspace_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckHealthcareWorkspaceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccHealthcareWorkspace_basic(context),
			},
			{
				ResourceName:            "google_healthcare_workspace.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dataset", "labels", "terraform_labels"},
			},
			{
				Config: testAccHealthcareWorkspace_update(context),
			},
			{
				ResourceName:            "google_healthcare_workspace.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dataset", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccHealthcareWorkspace_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_healthcare_workspace" "default" {
  name    = "tf-test-example-dm-workspace%{random_suffix}"
  dataset = google_healthcare_dataset.dataset.id

  settings {
    data_project_ids = ["tf-test-example-data-source-project-id%{random_suffix}"]
  }
  
  labels = {
    label1 = "labelvalue1"
  }
}


resource "google_healthcare_dataset" "dataset" {
  name     = "tf-test-example-dataset%{random_suffix}"
  location = "us-central1"
}
`, context)
}

func testAccHealthcareWorkspace_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_healthcare_workspace" "default" {
  name    = "tf-test-example-dm-workspace%{random_suffix}"
  dataset = google_healthcare_dataset.dataset.id

  settings {
    data_project_ids = ["tf-test-example-data-source-project-id%{random_suffix}"]
  }
}


resource "google_healthcare_dataset" "dataset" {
  name     = "tf-test-example-dataset%{random_suffix}"
  location = "us-central1"
}
`, context)
}
