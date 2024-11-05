// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package dataprocgdc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataprocGdcApplicationEnvironment_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocGdcApplicationEnvironment_full(context),
			},
			{
				ResourceName:            "google_dataproc_gdc_application_environment.application-environment",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "application_environment_id", "labels", "location", "serviceinstance", "terraform_labels"},
			},
			{
				Config: testAccDataprocGdcApplicationEnvironment_update(context),
			},
			{
				ResourceName:            "google_dataproc_gdc_application_environment.application-environment",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "application_environment_id", "labels", "location", "serviceinstance", "terraform_labels"},
			},
		},
	})
}

func testAccDataprocGdcApplicationEnvironment_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataproc_gdc_application_environment" "application-environment" {
  application_environment_id = "tf-test-dp-tf-e2e-application-environment%{random_suffix}"
  serviceinstance = "do-not-delete-dataproc-gdc-instance"
  project         = "gdce-cluster-monitoring"
  location        = "us-west2"
  namespace = "default"
  display_name = "An application environment"
  labels = {
    "test-label": "label-value"
  }
  annotations = {
    "an_annotation": "annotation_value"
  }
  spark_application_environment_config {
    default_properties = {
      "spark.executor.memory": "4g"
    }
    default_version = "1.2"
  }
}
`, context)
}

func testAccDataprocGdcApplicationEnvironment_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataproc_gdc_application_environment" "application-environment" {
  application_environment_id = "tf-test-dp-tf-e2e-application-environment%{random_suffix}"
  serviceinstance = "do-not-delete-dataproc-gdc-instance"
  project         = "gdce-cluster-monitoring"
  location        = "us-west2"
  namespace = "default"
  display_name = "An application environment"
  labels = {
    "test-label": "new-label-value"
  }
  annotations = {
    "an_annotation": "new_annotation_value"
    "another_annotation": "ok"
  }
  spark_application_environment_config {
    default_properties = {
      "spark.executor.memory": "2g"
    }
    default_version = "1.1"
  }
}
`, context)
}
