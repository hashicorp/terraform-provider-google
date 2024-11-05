// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package bigquery_test

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccBigqueryDatasetIamMember_afterDatasetCreation(t *testing.T) {
	t.Parallel()

	projectID := envvar.GetTestProjectFromEnv()
	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	authDatasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	saID := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	expected_auth := map[string]interface{}{
		"dataset": map[string]interface{}{
			"dataset": map[string]interface{}{
				"projectId": projectID,
				"datasetId": authDatasetID,
			},
			"targetTypes": []interface{}{"VIEWS"},
		},
	}

	expected_sa := map[string]interface{}{
		"role":        "roles/viewer",
		"userByEmail": fmt.Sprintf("%s@%s.iam.gserviceaccount.com", saID, envvar.GetTestProjectFromEnv()),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryDatasetIamMember_afterDatasetAccessCreation(projectID, datasetID, authDatasetID, saID),
				Check:  testAccCheckBigQueryDatasetIamMemberPresent(t, "google_bigquery_dataset.dataset", expected_auth),
			},
			{
				// For iam_member to be non-authoritative, we want authorized datasets to be present after destroy,
				// but the iam resources have to be deleted
				Config: testAccBigqueryDatasetIamMember_destroy(datasetID),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckBigQueryDatasetIamMemberPresent(t, "google_bigquery_dataset.dataset", expected_auth),
					testAccCheckBigQueryDatasetIamMemberAbsent(t, "google_bigquery_dataset.dataset", expected_sa),
				),
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
				Check:  testAccCheckBigQueryDatasetIamMemberPresent(t, "google_bigquery_dataset.dataset", expected),
			},
			{
				// Destroy step instead of CheckDestroy so we can check the access is removed without deleting the dataset
				Config: testAccBigqueryDatasetIamMember_destroy(datasetID),
				Check:  testAccCheckBigQueryDatasetIamMemberAbsent(t, "google_bigquery_dataset.dataset", expected),
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
				Check:  testAccCheckBigQueryDatasetIamMemberPresent(t, "google_bigquery_dataset.dataset", expected),
			},
			{
				// For iam_member to be non-authoritative, we want access block to be present after destroy
				Config: testAccBigqueryDatasetIamMember_destroy(datasetID),
				Check:  testAccCheckBigQueryDatasetIamMemberAbsent(t, "google_bigquery_dataset.dataset", expected),
			},
		},
	})
}

func TestAccBigqueryDatasetIamMember_withDeletedServiceAccount(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	serviceAccountID := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	serviceAccountEmail := fmt.Sprintf("%s@%s.iam.gserviceaccount.com", serviceAccountID, envvar.GetTestProjectFromEnv())

	expected := map[string]interface{}{
		"role":        "roles/viewer",
		"userByEmail": serviceAccountEmail,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryDatasetIamMember_withServiceAccount(datasetID, serviceAccountID),
				Check:  testAccCheckBigQueryDatasetIamMemberPresent(t, "google_bigquery_dataset.dataset", expected),
			},
			{
				// No change from Step 1 except for hard-coding the service account email to remove the explicit resource reference in preparation for the service accoutn deletion.
				Config: testAccBigqueryDatasetIamMember_withServiceAccountEmail(datasetID, serviceAccountID, serviceAccountEmail),
				Check:  testAccCheckBigQueryDatasetIamMemberPresent(t, "google_bigquery_dataset.dataset", expected),
			},
			{
				// Delete the service account but keep the rest the same.
				// The service account will show up with a "deleted:serviceAccount:" prefix and a "?uid=" suffix when getting dataset access from now on.
				// The plan should show an attempt to recreate the IAM member resource because the dataset access no longer contains the original account email so there is a diff.
				Config:             testAccBigqueryDatasetIamMember_withServiceAccountDeleted(datasetID, serviceAccountEmail),
				ExpectNonEmptyPlan: true,
			},
			{
				// Apply the same plan as the previous step.
				// It should fail because the specified account has now been deleted, so the account email is invalid.
				Config:      testAccBigqueryDatasetIamMember_withServiceAccountDeleted(datasetID, serviceAccountEmail),
				ExpectError: regexp.MustCompile(fmt.Sprintf(".+Service account %s does not exist.+", serviceAccountEmail)),
			},
		},
	})
}

func testAccCheckBigQueryDatasetIamMemberPresent(t *testing.T, n string, expected map[string]interface{}) resource.TestCheckFunc {
	return testAccCheckBigQueryDatasetIamMember(t, n, expected, true)
}

func testAccCheckBigQueryDatasetIamMemberAbsent(t *testing.T, n string, expected map[string]interface{}) resource.TestCheckFunc {
	return testAccCheckBigQueryDatasetIamMember(t, n, expected, false)
}

func testAccCheckBigQueryDatasetIamMember(t *testing.T, n string, expected map[string]interface{}, expectPresent bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		config := acctest.GoogleProviderConfig(t)
		url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{BigQueryBasePath}}projects/{{project}}/datasets/{{dataset_id}}")
		if err != nil {
			return err
		}

		ds, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			RawURL:    url,
			UserAgent: config.UserAgent,
		})
		if err != nil {
			return err
		}
		access := ds["access"].([]interface{})
		for _, a := range access {
			if reflect.DeepEqual(a, expected) {
				if !expectPresent {
					return fmt.Errorf("Found access %+v, expected not present", expected)
				}
				return nil
			}
		}
		if expectPresent {
			return fmt.Errorf("Did not find access %+v, expected present", expected)
		}
		return nil
	}
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

func testAccBigqueryDatasetIamMember_withServiceAccount(datasetID, serviceAccountID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "dataset" {
  dataset_id = "%s"
}

resource "google_service_account" "sa" {
  account_id = "%s"
}

resource "google_bigquery_dataset_iam_member" "access" {
  dataset_id = google_bigquery_dataset.dataset.dataset_id
  role       = "roles/viewer"
  member     = "serviceAccount:${google_service_account.sa.email}"
}
`, datasetID, serviceAccountID)
}

func testAccBigqueryDatasetIamMember_withServiceAccountEmail(datasetID, serviceAccountID, serviceAccountEmail string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "dataset" {
  dataset_id = "%s"
}

resource "google_service_account" "sa" {
  account_id = "%s"
}

resource "google_bigquery_dataset_iam_member" "access" {
  dataset_id = google_bigquery_dataset.dataset.dataset_id
  role       = "roles/viewer"
  member     = "serviceAccount:%s"
}
`, datasetID, serviceAccountID, serviceAccountEmail)
}

func testAccBigqueryDatasetIamMember_withServiceAccountDeleted(datasetID, serviceAccountEmail string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "dataset" {
  dataset_id = "%s"
}

resource "google_bigquery_dataset_iam_member" "access" {
  dataset_id = google_bigquery_dataset.dataset.dataset_id
  role       = "roles/viewer"
  member     = "serviceAccount:%s"
}
`, datasetID, serviceAccountEmail)
}
