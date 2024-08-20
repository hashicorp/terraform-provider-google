// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package bigquery_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccBigqueryDatasetIamMember_afterDatasetCreation(t *testing.T) {
	t.Parallel()

	projectID := envvar.GetTestProjectFromEnv()
	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	authDatasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	saID := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	expected := map[string]interface{}{
		"dataset": map[string]interface{}{
			"dataset": map[string]interface{}{
				"projectId": projectID,
				"datasetId": authDatasetID,
			},
			"targetTypes": []interface{}{"VIEWS"},
		},
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryDatasetIamMember_afterDatasetAccessCreation(projectID, datasetID, authDatasetID, saID),
				Check:  testAccCheckBigQueryDatasetAccessPresent(t, "google_bigquery_dataset.dataset", expected),
			},
			{
				// For iam_member to be non-authoritative, we want access block to be present after destroy
				Config: testAccBigqueryDatasetIamMember_destroy(datasetID),
				Check:  testAccCheckBigQueryDatasetAccessPresent(t, "google_bigquery_dataset.dataset", expected),
			},
		},
	})
}

func TestAccBigqueryDatasetIamMember_serviceAccount(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	saID := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	expected := map[string]interface{}{
		"role":        "roles/viewer",
		"userByEmail": fmt.Sprintf("%s@%s.iam.gserviceaccount.com", saID, envvar.GetTestProjectFromEnv()),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryDatasetIamMember_serviceAccount(datasetID, saID),
				Check:  testAccCheckBigQueryDatasetAccessPresent(t, "google_bigquery_dataset.dataset", expected),
			},
			{
				// Destroy step instead of CheckDestroy so we can check the access is removed without deleting the dataset
				Config: testAccBigqueryDatasetIamMember_destroy(datasetID),
				Check:  testAccCheckBigQueryDatasetAccessAbsent(t, "google_bigquery_dataset.dataset", expected),
			},
		},
	})
}

func TestAccBigqueryDatasetIamMember_iamMember(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	wifIDs := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	expected := map[string]interface{}{
		"role":      "roles/viewer",
		"iamMember": fmt.Sprintf("principal://iam.googleapis.com/projects/%s/locations/global/workloadIdentityPools/%s/subject/test", envvar.GetTestProjectNumberFromEnv(), wifIDs),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryDatasetIamMember_iamMember(datasetID, wifIDs),
				Check:  testAccCheckBigQueryDatasetAccessPresent(t, "google_bigquery_dataset.dataset", expected),
			},
			{
				// For iam_member to be non-authoritative, we want access block to be present after destroy
				Config: testAccBigqueryDatasetIamMember_destroy(datasetID),
				Check:  testAccCheckBigQueryDatasetAccessPresent(t, "google_bigquery_dataset.dataset", expected),
			},
		},
	})
}

func testAccBigqueryDatasetIamMember_destroy(datasetID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "dataset" {
  dataset_id = "%s"
}
`, datasetID)
}

func testAccBigqueryDatasetIamMember_serviceAccount(datasetID, saID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset_iam_member" "access" {
  dataset_id    = google_bigquery_dataset.dataset.dataset_id
  role          = "roles/viewer"
  member        = "serviceAccount:${google_service_account.bqviewer.email}"
}

resource "google_bigquery_dataset" "dataset" {
  dataset_id = "%s"
}

resource "google_service_account" "bqviewer" {
  account_id = "%s"
}
`, datasetID, saID)
}

func testAccBigqueryDatasetIamMember_afterDatasetAccessCreation(projectID, datasetID, authDatasetID, saID string) string {
	return fmt.Sprintf(`

resource "google_bigquery_dataset" "auth_dataset" {
	dataset_id = "%s"
}

resource "google_bigquery_dataset" "dataset" {
  dataset_id = "%s"
  access {
    dataset {
      dataset {
        project_id = "%s"
        dataset_id = google_bigquery_dataset.auth_dataset.dataset_id
      }
      target_types = ["VIEWS"]
    }
  }
  lifecycle {
    ignore_changes = [access]
  }
}

resource "google_service_account" "bqviewer" {
  account_id = "%s"
}

resource "google_bigquery_dataset_iam_member" "access" {
  dataset_id    = google_bigquery_dataset.dataset.dataset_id
  role          = "roles/viewer"
  member        = "serviceAccount:${google_service_account.bqviewer.email}"
}
`, authDatasetID, datasetID, projectID, saID)
}

func testAccBigqueryDatasetIamMember_iamMember(datasetID, wifIDs string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset_iam_member" "access" {
  dataset_id = google_bigquery_dataset.dataset.dataset_id
  role       = "roles/viewer"
  member     = "iamMember:principal://iam.googleapis.com/${google_iam_workload_identity_pool.wif_pool.name}/subject/test"
}

resource "google_bigquery_dataset" "dataset" {
  dataset_id = "%s"
}

resource "google_iam_workload_identity_pool" "wif_pool" {
  workload_identity_pool_id = "%s"
}

resource "google_iam_workload_identity_pool_provider" "wif_provider" {
  workload_identity_pool_id          = google_iam_workload_identity_pool.wif_pool.workload_identity_pool_id
  workload_identity_pool_provider_id = "%s"
  attribute_mapping                  = {
    "google.subject" = "assertion.sub"
  }
  oidc {
    issuer_uri = "https://issuer-uri.com"
  }
}
`, datasetID, wifIDs, wifIDs)
}
