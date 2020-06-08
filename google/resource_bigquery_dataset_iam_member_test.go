package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccBigqueryDatasetIamMember_basic(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", randString(t, 10))
	saID := fmt.Sprintf("tf-test-%s", randString(t, 10))

	expected := map[string]interface{}{
		"role":        "roles/viewer",
		"userByEmail": fmt.Sprintf("%s@%s.iam.gserviceaccount.com", saID, getTestProjectFromEnv()),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryDatasetIamMember_basic(datasetID, saID),
				Check:  testAccCheckBigQueryDatasetAccessPresent(t, "google_bigquery_dataset.dataset", expected),
			},
			{
				// Destroy step instead of CheckDestroy so we can check the access is removed without deleting the dataset
				Config: testAccBigqueryDatasetIamMember_destroy(datasetID, "dataset"),
				Check:  testAccCheckBigQueryDatasetAccessAbsent(t, "google_bigquery_dataset.dataset", expected),
			},
		},
	})
}

func testAccBigqueryDatasetIamMember_destroy(datasetID, rs string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "%s" {
  dataset_id = "%s"
}
`, rs, datasetID)
}

func testAccBigqueryDatasetIamMember_basic(datasetID, saID string) string {
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
