package google

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccBigQueryDatasetAccess_basic(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(10))
	saID := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	expected := map[string]interface{}{
		"role":        "OWNER",
		"userByEmail": fmt.Sprintf("%s@%s.iam.gserviceaccount.com", saID, getTestProjectFromEnv()),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryDatasetAccess_basic(datasetID, saID),
				Check:  testAccCheckBigQueryDatasetAccessPresent("google_bigquery_dataset.dataset", expected),
			},
			{
				// Destroy step instead of CheckDestroy so we can check the access is removed without deleting the dataset
				Config: testAccBigQueryDatasetAccess_destroy(datasetID, "dataset"),
				Check:  testAccCheckBigQueryDatasetAccessAbsent("google_bigquery_dataset.dataset", expected),
			},
		},
	})
}

func TestAccBigQueryDatasetAccess_view(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(10))
	datasetID2 := fmt.Sprintf("tf_test_%s", acctest.RandString(10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(10))

	expected := map[string]interface{}{
		"view": map[string]interface{}{
			"projectId": getTestProjectFromEnv(),
			"datasetId": datasetID2,
			"tableId":   tableID,
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryDatasetAccess_view(datasetID, datasetID2, tableID),
				Check:  testAccCheckBigQueryDatasetAccessPresent("google_bigquery_dataset.private", expected),
			},
			{
				Config: testAccBigQueryDatasetAccess_destroy(datasetID, "private"),
				Check:  testAccCheckBigQueryDatasetAccessAbsent("google_bigquery_dataset.private", expected),
			},
		},
	})
}

func TestAccBigQueryDatasetAccess_multiple(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(10))

	expected1 := map[string]interface{}{
		"role":   "WRITER",
		"domain": "google.com",
	}

	expected2 := map[string]interface{}{
		"role":         "READER",
		"specialGroup": "projectWriters",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryDatasetAccess_multiple(datasetID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBigQueryDatasetAccessPresent("google_bigquery_dataset.dataset", expected1),
					testAccCheckBigQueryDatasetAccessPresent("google_bigquery_dataset.dataset", expected2),
				),
			},
			{
				// Destroy step instead of CheckDestroy so we can check the access is removed without deleting the dataset
				Config: testAccBigQueryDatasetAccess_destroy(datasetID, "dataset"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBigQueryDatasetAccessAbsent("google_bigquery_dataset.dataset", expected1),
					testAccCheckBigQueryDatasetAccessAbsent("google_bigquery_dataset.dataset", expected2),
				),
			},
		},
	})
}

func testAccCheckBigQueryDatasetAccessPresent(n string, expected map[string]interface{}) resource.TestCheckFunc {
	return testAccCheckBigQueryDatasetAccess(n, expected, true)
}

func testAccCheckBigQueryDatasetAccessAbsent(n string, expected map[string]interface{}) resource.TestCheckFunc {
	return testAccCheckBigQueryDatasetAccess(n, expected, false)
}

func testAccCheckBigQueryDatasetAccess(n string, expected map[string]interface{}, expectPresent bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		config := testAccProvider.Meta().(*Config)
		url, err := replaceVarsForTest(config, rs, "{{BigQueryBasePath}}projects/{{project}}/datasets/{{dataset_id}}")
		if err != nil {
			return err
		}

		ds, err := sendRequest(config, "GET", "", url, nil)
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

func testAccBigQueryDatasetAccess_destroy(datasetID, rs string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "%s" {
  dataset_id = "%s"
}
`, rs, datasetID)
}

func testAccBigQueryDatasetAccess_basic(datasetID, saID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset_access" "access" {
  dataset_id    = google_bigquery_dataset.dataset.dataset_id
  role          = "OWNER"
  user_by_email = google_service_account.bqowner.email
}

resource "google_bigquery_dataset" "dataset" {
  dataset_id = "%s"
}

resource "google_service_account" "bqowner" {
  account_id = "%s"
}
`, datasetID, saID)
}

func testAccBigQueryDatasetAccess_view(datasetID, datasetID2, tableID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset_access" "access" {
  dataset_id    = google_bigquery_dataset.private.dataset_id
  view {
    project_id = google_bigquery_table.public.project
    dataset_id = google_bigquery_dataset.public.dataset_id
    table_id   = google_bigquery_table.public.table_id
  }
}

resource "google_bigquery_dataset" "private" {
  dataset_id = "%s"
}

resource "google_bigquery_dataset" "public" {
  dataset_id = "%s"
}

resource "google_bigquery_table" "public" {
  dataset_id = google_bigquery_dataset.public.dataset_id
  table_id   = "%s"

  view {
    query          = "%s"
    use_legacy_sql = false
  }
}

`, datasetID, datasetID2, tableID, "SELECT state FROM `lookerdata.cdc.project_tycho_reports`")
}

func testAccBigQueryDatasetAccess_multiple(datasetID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset_access" "access" {
  dataset_id = google_bigquery_dataset.dataset.dataset_id
  role       = "WRITER"
  domain     = "google.com"
}

resource "google_bigquery_dataset_access" "access2" {
  dataset_id    = google_bigquery_dataset.dataset.dataset_id
  role          = "READER"
  special_group = "projectWriters"
}

resource "google_bigquery_dataset" "dataset" {
  dataset_id = "%s"
}
`, datasetID)
}
