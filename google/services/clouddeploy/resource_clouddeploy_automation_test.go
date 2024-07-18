// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package clouddeploy_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccClouddeployAutomation_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"service_account": envvar.GetTestServiceAccountFromEnv(t),
		"random_suffix":   acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckClouddeployAutomationDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccClouddeployAutomation_basic(context),
			},
			{
				ResourceName:            "google_clouddeploy_automation.automation",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "delivery_pipeline", "annotations", "labels", "terraform_labels"},
			},
			{
				Config: testAccClouddeployAutomation_update(context),
			},
			{
				ResourceName:            "google_clouddeploy_automation.automation",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "delivery_pipeline", "annotations", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccClouddeployAutomation_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_clouddeploy_automation" "automation" {
  name     = "tf-test-cd-automation%{random_suffix}"
  location = "us-central1"
  delivery_pipeline = google_clouddeploy_delivery_pipeline.pipeline.name
  service_account = "%{service_account}"
  selector {
    targets {
      id = "*"
      labels = {}
    }
  }
  rules {
    advance_rollout_rule {
      id                    = "advance-rollout"
      source_phases         = ["deploy"]
      wait                  = "200s"
    }
  }
}

resource "google_clouddeploy_delivery_pipeline" "pipeline" {
  name = "tf-test-cd-pipeline%{random_suffix}"
  location = "us-central1"
  serial_pipeline  {
    stages {
      target_id = "test"
      profiles = ["test-profile"]
    }
  }
 }
`, context)
}

func testAccClouddeployAutomation_update(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_clouddeploy_automation" "automation" {
  name     = "tf-test-cd-automation%{random_suffix}"
  location = "us-central1"
  delivery_pipeline = google_clouddeploy_delivery_pipeline.pipeline.name
  service_account = "%{service_account}"
  annotations = {
     first_annotation = "example-annotation-1"
     second_annotation = "example-annotation-2"
  }
  labels = {
    first_label = "example-label-1"
    second_label = "example-label-2"
  }
  description = "automation resource"
  selector {
    targets {
			id = "dev"
			labels = {
				foo = "bar2"
			}
    }
  }
	suspended = true
  rules {
    advance_rollout_rule {
      id                    = "advance-rollout"
      source_phases         = ["verify"]
      wait                  = "100s"
    }
  }
  rules {
    promote_release_rule{
      id = "promote-release"
      wait = "200s"
      destination_target_id = "@next"
      destination_phase = "stable"
    }
  }
}

resource "google_clouddeploy_delivery_pipeline" "pipeline" {
  name = "tf-test-cd-pipeline%{random_suffix}"
  location = "us-central1"
  serial_pipeline  {
    stages {
      target_id = "test"
      profiles = ["test-profile"]
    }
  }
 }
`, context)
}
