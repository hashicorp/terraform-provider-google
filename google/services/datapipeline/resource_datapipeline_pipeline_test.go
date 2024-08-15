// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package datapipeline_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func setTestCheckDataPipelinePipelineId(res string, pipelineId *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		updateId, err := getTestResourceDataPipelinePipelineId(res, s)
		if err != nil {
			return err
		}
		*pipelineId = updateId
		return nil
	}
}

func testCheckDataPipelinePipelineIdAfterUpdate(res string, pipelineId *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		updateId, err := getTestResourceDataPipelinePipelineId(res, s)
		if err != nil {
			return err
		}

		if pipelineId == nil {
			return fmt.Errorf("unexpected error, pipeline ID was not set")
		}

		if *pipelineId != updateId {
			return fmt.Errorf("unexpected mismatch in pipeline ID after update, resource was recreated. Initial %q, Updated %q",
				*pipelineId, updateId)
		}
		return nil
	}
}

func getTestResourceDataPipelinePipelineId(res string, s *terraform.State) (string, error) {
	rs, ok := s.RootModule().Resources[res]
	if !ok {
		return "", fmt.Errorf("not found: %s", res)
	}

	if rs.Primary.ID == "" {
		return "", fmt.Errorf("no ID is set for %s", res)
	}

	if v, ok := rs.Primary.Attributes["id"]; ok {
		return v, nil
	}

	return "", fmt.Errorf("id not set on resource %s", res)
}

func TestAccDataPipelinePipeline_basicLaunchTemplate(t *testing.T) {
	t.Parallel()

	var generatedId string
	suffix := acctest.RandString(t, 10)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataPipelinePipelineDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataPipelinePipeline_basicLaunchTemplate(suffix),
				Check:  setTestCheckDataPipelinePipelineId("google_data_pipeline_pipeline.primary", &generatedId),
			},
			{
				ResourceName:      "google_data_pipeline_pipeline.primary",
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore input-only field for import
				ImportStateVerifyIgnore: []string{"region"},
			},
			{
				Config: testAccDataPipelinePipeline_basicLaunchTemplateUpdate(suffix),
				Check:  testCheckDataPipelinePipelineIdAfterUpdate("google_data_pipeline_pipeline.primary", &generatedId),
			},
			{
				ResourceName:      "google_data_pipeline_pipeline.primary",
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore input-only field for import
				ImportStateVerifyIgnore: []string{"region"},
			},
		},
	})
}

func TestAccDataPipelinePipeline_basicFlexTemplate(t *testing.T) {
	t.Parallel()

	var generatedId string
	suffix := acctest.RandString(t, 10)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataPipelinePipelineDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataPipelinePipeline_basicFlexTemplate(suffix),
				Check:  setTestCheckDataPipelinePipelineId("google_data_pipeline_pipeline.primary", &generatedId),
			},
			{
				ResourceName:      "google_data_pipeline_pipeline.primary",
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore input-only field for import
				ImportStateVerifyIgnore: []string{"region"},
			},
			{
				Config: testAccDataPipelinePipeline_basicFlexTemplateUpdate(suffix),
				Check:  testCheckDataPipelinePipelineIdAfterUpdate("google_data_pipeline_pipeline.primary", &generatedId),
			},
			{
				ResourceName:      "google_data_pipeline_pipeline.primary",
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore input-only field for import
				ImportStateVerifyIgnore: []string{"region"},
			},
		},
	})
}

func testAccDataPipelinePipeline_basicFlexTemplateUpdate(suffix string) string {
	return fmt.Sprintf(`
resource "google_service_account" "service_account" {
  account_id   = "tf-test-service-%s"
  display_name = "Service Account"
}

resource "google_data_pipeline_pipeline" "primary" {
  name         = "tf-test-pipeline-%s"
  display_name = "update-pipeline"
  type         = "PIPELINE_TYPE_BATCH"
  state        = "STATE_ACTIVE"

  workload {
    dataflow_flex_template_request {
      project_id = "my-project"
      launch_parameter {
        job_name = "my-job"
        parameters = {
          "name": "wrench"
        }
        environment {
          num_workers                     = 5
          max_workers                     = 5
          zone                            = "us-centra1-a"
          service_account_email           = google_service_account.service_account.email
          network                         = "default"
          temp_location                   = "gs://my-bucket/tmp_dir"
          machine_type                    = "E2"
          additional_experiments          = ["test"]
          additional_user_labels = {
            "context" : "test"
          }
          worker_region    = "us-central1"
          worker_zone      = "us-central1-a"

          enable_streaming_engine = "false"
        }
        container_spec_gcs_path = "gs://my-bucket/template"
        update                  = false
        transform_name_mappings = {"name": "wrench"}
      }
      location = "us-central1"
    }
  }
  schedule_info {
    schedule  = "0 * * * *"
    time_zone = "UTC"
  }
  pipeline_sources = {
    "name": "wrench"
  }
}
`, suffix, suffix)
}
func testAccDataPipelinePipeline_basicFlexTemplate(suffix string) string {
	return fmt.Sprintf(`
resource "google_service_account" "service_account" {
  account_id   = "tf-test-service-%s"
  display_name = "Service Account"
}

resource "google_data_pipeline_pipeline" "primary" {
  name         = "tf-test-pipeline-%s"
  type         = "PIPELINE_TYPE_BATCH"
  state        = "STATE_ACTIVE"

  workload {
    dataflow_flex_template_request {
      project_id = "my-project"
      launch_parameter {
        job_name = "my-job"
        parameters = {
          "name": "wrench"
        }
        environment {
          num_workers                     = 5
          max_workers                     = 5
          zone                            = "us-centra1-a"
          service_account_email           = google_service_account.service_account.email
          network                         = "default"
          temp_location                   = "gs://my-bucket/tmp_dir"
          machine_type                    = "E2"
          additional_experiments          = []
          additional_user_labels = {
            "context" : "test"
          }
          worker_region    = "us-central1"
          worker_zone      = "us-central1-a"

          enable_streaming_engine = "false"
        }
        container_spec_gcs_path = "gs://my-bucket/template"
        update                  = false
        transform_name_mappings = {"name": "wrench"}
      }
      location = "us-central1"
    }
  }
  schedule_info {
    schedule  = "0 * * * *"
    time_zone = "UTC"
  }
  pipeline_sources = {
    "name": "wrench"
  }
}
`, suffix, suffix)
}

func testAccDataPipelinePipeline_basicLaunchTemplateUpdate(suffix string) string {
	return fmt.Sprintf(`
resource "google_service_account" "service_account" {
  account_id   = "tf-test-service-%s"
  display_name = "Service Account"
}

resource "google_data_pipeline_pipeline" "primary" {
  name         = "tf-test-pipeline-%s"
  display_name = "update-pipeline"
  type         = "PIPELINE_TYPE_BATCH"
  state        = "STATE_ACTIVE"

  workload {
    dataflow_launch_template_request {
      project_id = "my-project"
      gcs_path = "gs://my-bucket/path"
      launch_parameters {
        job_name = "my-job"
        parameters = {
          "name": "wrench"
        }
        environment {
          num_workers                     = 5
          max_workers                     = 5
          zone                            = "us-centra1-a"
          service_account_email           = google_service_account.service_account.email
          network                         = "default"
          temp_location                   = "gs://my-bucket/tmp_dir"
          bypass_temp_dir_validation      = false
          machine_type                    = "E2"
          additional_experiments          = ["test"]
          additional_user_labels = {
            "context" : "test"
          }
          worker_region    = "us-central1"
          worker_zone      = "us-central1-a"

          enable_streaming_engine = "false"
        }
        update                  = false
        transform_name_mapping = {"name": "wrench"}
      }
      location = "us-central1"
    }
  }
  schedule_info {
    schedule  = "0 * * * *"
    time_zone = "UTC"
  }
  pipeline_sources = {
    "name": "wrench"
  }
}
`, suffix, suffix)
}
func testAccDataPipelinePipeline_basicLaunchTemplate(suffix string) string {
	return fmt.Sprintf(`
resource "google_service_account" "service_account" {
  account_id   = "tf-test-service-%s"
  display_name = "Service Account"
}

resource "google_data_pipeline_pipeline" "primary" {
  name         = "tf-test-pipeline-%s"
  type         = "PIPELINE_TYPE_BATCH"
  state        = "STATE_ACTIVE"

  workload {
    dataflow_launch_template_request {
      project_id = "my-project"
      gcs_path = "gs://my-bucket/path"
      launch_parameters {
        job_name = "my-job"
        parameters = {
          "name": "wrench"
        }
        environment {
          num_workers                     = 5
          max_workers                     = 5
          zone                            = "us-centra1-a"
          service_account_email           = google_service_account.service_account.email
          network                         = "default"
          temp_location                   = "gs://my-bucket/tmp_dir"
          bypass_temp_dir_validation      = false
          machine_type                    = "E2"
          additional_experiments          = []
          additional_user_labels = {
            "context" : "test"
          }
          worker_region    = "us-central1"
          worker_zone      = "us-central1-a"

          enable_streaming_engine = "false"
        }
        update                  = false
        transform_name_mapping = {"name": "wrench"}
      }
      location = "us-central1"
    }
  }
  schedule_info {
    schedule  = "0 * * * *"
    time_zone = "UTC"
  }
  pipeline_sources = {
    "name": "wrench"
  }
}
`, suffix, suffix)
}
