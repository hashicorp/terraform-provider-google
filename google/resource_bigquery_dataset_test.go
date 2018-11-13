package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccBigQueryDataset_basic(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigQueryDatasetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryDataset(datasetID),
			},
			{
				ResourceName:      "google_bigquery_dataset.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBigQueryDatasetUpdated(datasetID),
			},
			{
				ResourceName:      "google_bigquery_dataset.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigQueryDataset_access(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_access_%s", acctest.RandString(10))
	otherDatasetID := fmt.Sprintf("tf_test_other_%s", acctest.RandString(10))
	otherTableID := fmt.Sprintf("tf_test_other_%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigQueryDatasetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryDatasetWithOneAccess(datasetID),
			},
			{
				ResourceName:      "google_bigquery_dataset.access_test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBigQueryDatasetWithTwoAccess(datasetID),
			},
			{
				ResourceName:      "google_bigquery_dataset.access_test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBigQueryDatasetWithOneAccess(datasetID),
			},
			{
				ResourceName:      "google_bigquery_dataset.access_test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBigQueryDatasetWithViewAccess(datasetID, otherDatasetID, otherTableID),
			},
			{
				ResourceName:      "google_bigquery_dataset.access_test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckBigQueryDatasetDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_bigquery_dataset" {
			continue
		}

		_, err := config.clientBigQuery.Datasets.Get(config.Project, rs.Primary.Attributes["dataset_id"]).Do()
		if err == nil {
			return fmt.Errorf("Dataset still exists")
		}
	}

	return nil
}

func testAccBigQueryDataset(datasetID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id                  = "%s"
  friendly_name               = "foo"
  description                 = "This is a foo description"
  location                    = "EU"
  default_partition_expiration_ms = 3600000
  default_table_expiration_ms = 3600000

  labels {
    env                         = "foo"
    default_table_expiration_ms = 3600000
  }
}`, datasetID)
}

func testAccBigQueryDatasetUpdated(datasetID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id                  = "%s"
  friendly_name               = "bar"
  description                 = "This is a bar description"
  location                    = "EU"
  default_partition_expiration_ms = 7200000
  default_table_expiration_ms = 7200000

  labels {
    env                         = "bar"
    default_table_expiration_ms = 7200000
  }
}`, datasetID)
}

func testAccBigQueryDatasetWithOneAccess(datasetID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "access_test" {
  dataset_id = "%s"

  access {
    role          = "OWNER"
    user_by_email = "Joe@example.com"
  }

  labels {
    env                         = "foo"
    default_table_expiration_ms = 3600000
  }
}`, datasetID)
}

func testAccBigQueryDatasetWithTwoAccess(datasetID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "access_test" {
  dataset_id = "%s"

  access {
    role          = "OWNER"
    user_by_email = "Joe@example.com"
  }
  access {
    role   = "READER"
    domain = "example.com"
  }

  labels {
    env                         = "foo"
    default_table_expiration_ms = 3600000
  }
}`, datasetID)
}

func testAccBigQueryDatasetWithViewAccess(datasetID, otherDatasetID, otherTableID string) string {
	// Note that we have to add a non-view access to prevent BQ from creating 4 default
	// access entries.
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "other_dataset" {
  dataset_id = "%s"
}

resource "google_bigquery_table" "table_with_view" {
  table_id   = "%s"
  dataset_id = "${google_bigquery_dataset.other_dataset.dataset_id}"

  time_partitioning {
    type = "DAY"
  }

  view {
    query = "SELECT state FROM [lookerdata:cdc.project_tycho_reports]"
    use_legacy_sql = true
  }
}

resource "google_bigquery_dataset" "access_test" {
  dataset_id = "%s"

  access {
    role          = "OWNER"
    user_by_email = "Joe@example.com"
  }
  access {
    view {
      project_id = "${google_bigquery_dataset.other_dataset.project}"
      dataset_id = "${google_bigquery_dataset.other_dataset.dataset_id}"
      table_id   = "${google_bigquery_table.table_with_view.table_id}"
    }
  }

  labels {
    env                         = "foo"
    default_table_expiration_ms = 3600000
  }
}`, otherDatasetID, otherTableID, datasetID)
}
