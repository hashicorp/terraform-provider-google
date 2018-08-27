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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBigQueryDatasetExists(
						"google_bigquery_dataset.test"),
				),
			},

			{
				Config: testAccBigQueryDatasetUpdated(datasetID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBigQueryDatasetExists(
						"google_bigquery_dataset.test"),
				),
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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBigQueryDatasetExists(
						"google_bigquery_dataset.access_test"),
				),
			},

			{
				Config: testAccBigQueryDatasetWithTwoAccess(datasetID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBigQueryDatasetExists(
						"google_bigquery_dataset.access_test"),
				),
			},

			{
				Config: testAccBigQueryDatasetWithOneAccess(datasetID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBigQueryDatasetExists(
						"google_bigquery_dataset.access_test"),
				),
			},

			{
				Config: testAccBigQueryDatasetWithViewAccess(datasetID, otherDatasetID, otherTableID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBigQueryDatasetExists(
						"google_bigquery_dataset.access_test"),
				),
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

func testAccCheckBigQueryDatasetExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		found, err := config.clientBigQuery.Datasets.Get(config.Project, rs.Primary.Attributes["dataset_id"]).Do()
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("Dataset not found")
		}

		return nil
	}
}

func testAccBigQueryDataset(datasetID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id                  = "%s"
  friendly_name               = "foo"
  description                 = "This is a foo description"
  location                    = "EU"
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
  dataset_id                  = "%s"

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
  dataset_id                  = "%s"

  access {
    role          = "OWNER"
    user_by_email = "Joe@example.com"
  }
  access {
    role	      = "READER"
    domain	      = "example.com"
  }

  labels {
    env                         = "foo"
    default_table_expiration_ms = 3600000
  }
}`, datasetID)
}

func getBigQueryTableWithView(datasetID, tableID string) string {
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
}`, datasetID, tableID)
}

func testAccBigQueryDatasetWithViewAccess(datasetID, otherDatasetID, otherTableID string) string {
	otherTable := getBigQueryTableWithView(otherDatasetID, otherTableID)
	// Note that we have to add a non-view access to prevent BQ from creating 4 default
	// access entries.
	return fmt.Sprintf(`
%s

resource "google_bigquery_dataset" "access_test" {
  dataset_id                  = "%s"
  friendly_name               = "foo"
  description                 = "This is a foo description"
  location                    = "EU"
  default_table_expiration_ms = 3600000

  access = [
    {
	role          = "OWNER"
	user_by_email = "Joe@example.com"
    },
    {
	view = {
		project_id = "${google_bigquery_dataset.other_dataset.project}"
		dataset_id = "${google_bigquery_dataset.other_dataset.dataset_id}"
		table_id   = "${google_bigquery_table.table_with_view.table_id}"
	}
    }
  ]

  labels {
    env                         = "foo"
    default_table_expiration_ms = 3600000
  }
}`, otherTable, datasetID)
}
