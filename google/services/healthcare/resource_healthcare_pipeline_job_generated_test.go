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

package healthcare_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccHealthcarePipelineJob_healthcarePipelineJobReconciliationExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckHealthcarePipelineJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccHealthcarePipelineJob_healthcarePipelineJobReconciliationExample(context),
			},
			{
				ResourceName:            "google_healthcare_pipeline_job.example-pipeline",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dataset", "labels", "location", "self_link", "terraform_labels"},
			},
		},
	})
}

func testAccHealthcarePipelineJob_healthcarePipelineJobReconciliationExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}

resource "google_healthcare_pipeline_job" "example-pipeline" {
  name  = "tf_test_example_pipeline_job%{random_suffix}"
  location = "us-central1"
  dataset = google_healthcare_dataset.dataset.id
  disable_lineage = true
  reconciliation_pipeline_job {
    merge_config {
      description = "sample description for reconciliation rules"
      whistle_config_source {
        uri = "gs://${google_storage_bucket.bucket.name}/${google_storage_bucket_object.merge_file.name}"
        import_uri_prefix = "gs://${google_storage_bucket.bucket.name}"
      }
    }
    matching_uri_prefix = "gs://${google_storage_bucket.bucket.name}"
    fhir_store_destination = "${google_healthcare_dataset.dataset.id}/fhirStores/${google_healthcare_fhir_store.fhirstore.name}"
  }
}

resource "google_healthcare_dataset" "dataset" {
  name     = "tf_test_example_dataset%{random_suffix}"
  location = "us-central1"
}

resource "google_healthcare_fhir_store" "fhirstore" {
  name    = "tf_test_fhir_store%{random_suffix}"
  dataset = google_healthcare_dataset.dataset.id
  version = "R4"
  enable_update_create          = true
  disable_referential_integrity = true
}

resource "google_storage_bucket" "bucket" {
    name          = "tf_test_example_bucket_name%{random_suffix}"
    location      = "us-central1"
    uniform_bucket_level_access = true
}

resource "google_storage_bucket_object" "merge_file" {
  name    = "merge.wstl"
  content = " "
  bucket  = google_storage_bucket.bucket.name
}

resource "google_storage_bucket_iam_member" "hsa" {
    bucket = google_storage_bucket.bucket.name
    role   = "roles/storage.objectUser"
    member = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-healthcare.iam.gserviceaccount.com"
}
`, context)
}

func TestAccHealthcarePipelineJob_healthcarePipelineJobBackfillExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckHealthcarePipelineJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccHealthcarePipelineJob_healthcarePipelineJobBackfillExample(context),
			},
			{
				ResourceName:            "google_healthcare_pipeline_job.example-pipeline",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dataset", "labels", "location", "self_link", "terraform_labels"},
			},
		},
	})
}

func testAccHealthcarePipelineJob_healthcarePipelineJobBackfillExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_healthcare_pipeline_job" "example-pipeline" {
  name  = "tf_test_example_backfill_pipeline%{random_suffix}"
  location = "us-central1"
  dataset = google_healthcare_dataset.dataset.id
  backfill_pipeline_job {
    mapping_pipeline_job = "${google_healthcare_dataset.dataset.id}/pipelineJobs/tf_test_example_mapping_pipeline_job%{random_suffix}"
  }      
}

resource "google_healthcare_dataset" "dataset" {
  name     = "tf_test_example_dataset%{random_suffix}"
  location = "us-central1"
}
`, context)
}

func TestAccHealthcarePipelineJob_healthcarePipelineJobWhistleMappingExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckHealthcarePipelineJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccHealthcarePipelineJob_healthcarePipelineJobWhistleMappingExample(context),
			},
			{
				ResourceName:            "google_healthcare_pipeline_job.example-mapping-pipeline",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dataset", "labels", "location", "self_link", "terraform_labels"},
			},
		},
	})
}

func testAccHealthcarePipelineJob_healthcarePipelineJobWhistleMappingExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}

resource "google_healthcare_pipeline_job" "example-mapping-pipeline" {
  name  = "tf_test_example_mapping_pipeline_job%{random_suffix}"
  location = "us-central1"
  dataset = google_healthcare_dataset.dataset.id
  disable_lineage = true
  labels = {
    example_label_key = "example_label_value"
  }  
  mapping_pipeline_job {
    mapping_config {
      whistle_config_source {
        uri = "gs://${google_storage_bucket.bucket.name}/${google_storage_bucket_object.mapping_file.name}"
        import_uri_prefix = "gs://${google_storage_bucket.bucket.name}"
      }
      description = "example description for mapping configuration"
    }
    fhir_streaming_source {
      fhir_store = "${google_healthcare_dataset.dataset.id}/fhirStores/${google_healthcare_fhir_store.source_fhirstore.name}"
      description = "example description for streaming fhirstore"
    }
    fhir_store_destination = "${google_healthcare_dataset.dataset.id}/fhirStores/${google_healthcare_fhir_store.dest_fhirstore.name}"
  }
}

resource "google_healthcare_dataset" "dataset" {
  name     = "tf_test_example_dataset%{random_suffix}"
  location = "us-central1"
}

resource "google_healthcare_fhir_store" "source_fhirstore" {
  name    = "tf_test_source_fhir_store%{random_suffix}"
  dataset = google_healthcare_dataset.dataset.id
  version = "R4"
  enable_update_create          = true
  disable_referential_integrity = true
}

resource "google_healthcare_fhir_store" "dest_fhirstore" {
  name    = "tf_test_dest_fhir_store%{random_suffix}"
  dataset = google_healthcare_dataset.dataset.id
  version = "R4"
  enable_update_create          = true
  disable_referential_integrity = true
}

resource "google_storage_bucket" "bucket" {
    name          = "tf_test_example_bucket_name%{random_suffix}"
    location      = "us-central1"
    uniform_bucket_level_access = true
}

resource "google_storage_bucket_object" "mapping_file" {
  name    = "mapping.wstl"
  content = " "
  bucket  = google_storage_bucket.bucket.name
}

resource "google_storage_bucket_iam_member" "hsa" {
    bucket = google_storage_bucket.bucket.name
    role   = "roles/storage.objectUser"
    member = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-healthcare.iam.gserviceaccount.com"
}
`, context)
}

func TestAccHealthcarePipelineJob_healthcarePipelineJobMappingReconDestExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckHealthcarePipelineJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccHealthcarePipelineJob_healthcarePipelineJobMappingReconDestExample(context),
			},
			{
				ResourceName:            "google_healthcare_pipeline_job.example-mapping-pipeline",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dataset", "labels", "location", "self_link", "terraform_labels"},
			},
		},
	})
}

func testAccHealthcarePipelineJob_healthcarePipelineJobMappingReconDestExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}

resource "google_healthcare_pipeline_job" "recon" {
  name  = "tf_test_example_recon_pipeline_job%{random_suffix}"
  location = "us-central1"
  dataset = google_healthcare_dataset.dataset.id
  disable_lineage = true
  reconciliation_pipeline_job {
    merge_config {
      description = "sample description for reconciliation rules"
      whistle_config_source {
        uri = "gs://${google_storage_bucket.bucket.name}/${google_storage_bucket_object.merge_file.name}"
        import_uri_prefix = "gs://${google_storage_bucket.bucket.name}"
      }
    }
    matching_uri_prefix = "gs://${google_storage_bucket.bucket.name}"
    fhir_store_destination = "${google_healthcare_dataset.dataset.id}/fhirStores/${google_healthcare_fhir_store.dest_fhirstore.name}"
  }
}

resource "google_healthcare_pipeline_job" "example-mapping-pipeline" {
  depends_on = [google_healthcare_pipeline_job.recon]
  name  = "tf_test_example_mapping_pipeline_job%{random_suffix}"
  location = "us-central1"
  dataset = google_healthcare_dataset.dataset.id
  disable_lineage = true
  labels = {
    example_label_key = "example_label_value"
  }
  mapping_pipeline_job {
    mapping_config {
      whistle_config_source {
        uri = "gs://${google_storage_bucket.bucket.name}/${google_storage_bucket_object.mapping_file.name}"
        import_uri_prefix = "gs://${google_storage_bucket.bucket.name}"
      }
      description = "example description for mapping configuration"
    }
    fhir_streaming_source {
      fhir_store = "${google_healthcare_dataset.dataset.id}/fhirStores/${google_healthcare_fhir_store.source_fhirstore.name}"
      description = "example description for streaming fhirstore"
    }
    reconciliation_destination = true
  }
}

resource "google_healthcare_dataset" "dataset" {
  name     = "tf_test_example_dataset%{random_suffix}"
  location = "us-central1"
}

resource "google_healthcare_fhir_store" "source_fhirstore" {
  name    = "tf_test_source_fhir_store%{random_suffix}"
  dataset = google_healthcare_dataset.dataset.id
  version = "R4"
  enable_update_create          = true
  disable_referential_integrity = true
}

resource "google_healthcare_fhir_store" "dest_fhirstore" {
  name    = "tf_test_dest_fhir_store%{random_suffix}"
  dataset = google_healthcare_dataset.dataset.id
  version = "R4"
  enable_update_create          = true
  disable_referential_integrity = true
}

resource "google_storage_bucket" "bucket" {
    name          = "tf_test_example_bucket_name%{random_suffix}"
    location      = "us-central1"
    uniform_bucket_level_access = true
}

resource "google_storage_bucket_object" "mapping_file" {
  name    = "mapping.wstl"
  content = " "
  bucket  = google_storage_bucket.bucket.name
}

resource "google_storage_bucket_object" "merge_file" {
  name    = "merge.wstl"
  content = " "
  bucket  = google_storage_bucket.bucket.name
}

resource "google_storage_bucket_iam_member" "hsa" {
    bucket = google_storage_bucket.bucket.name
    role   = "roles/storage.objectUser"
    member = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-healthcare.iam.gserviceaccount.com"
}
`, context)
}

func testAccCheckHealthcarePipelineJobDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_healthcare_pipeline_job" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{HealthcareBasePath}}{{dataset}}/pipelineJobs/{{name}}")
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
				return fmt.Errorf("HealthcarePipelineJob still exists at %s", url)
			}
		}

		return nil
	}
}
