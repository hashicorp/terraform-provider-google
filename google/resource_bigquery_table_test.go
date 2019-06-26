package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccBigQueryTable_Basic(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigQueryTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTable(datasetID, tableID),
			},
			{
				ResourceName:      "google_bigquery_table.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBigQueryTableUpdated(datasetID, tableID),
			},
			{
				ResourceName:      "google_bigquery_table.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigQueryTable_View(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigQueryTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableWithView(datasetID, tableID),
			},
			{
				ResourceName:      "google_bigquery_table.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigQueryTable_updateView(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigQueryTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableWithView(datasetID, tableID),
			},
			{
				ResourceName:      "google_bigquery_table.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBigQueryTableWithNewSqlView(datasetID, tableID),
			},
			{
				ResourceName:      "google_bigquery_table.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckBigQueryTableDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_bigquery_table" {
			continue
		}

		config := testAccProvider.Meta().(*Config)
		_, err := config.clientBigQuery.Tables.Get(config.Project, rs.Primary.Attributes["dataset_id"], rs.Primary.Attributes["table_id"]).Do()
		if err == nil {
			return fmt.Errorf("Table still present")
		}
	}

	return nil
}

func testAccBigQueryTable(datasetID, tableID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}

resource "google_bigquery_table" "test" {
  table_id   = "%s"
  dataset_id = "${google_bigquery_dataset.test.dataset_id}"

  time_partitioning {
    type = "DAY"
    field = "ts"
    require_partition_filter = true
  }

  schema = <<EOH
[
  {
    "name": "ts",
    "type": "TIMESTAMP"
  },
  {
    "name": "city",
    "type": "RECORD",
    "fields": [
      {
        "name": "id",
        "type": "INTEGER"
      },
      {
        "name": "coord",
        "type": "RECORD",
        "fields": [
          {
            "name": "lon",
            "type": "FLOAT"
          }
        ]
      }
    ]
  }
]
EOH
}`, datasetID, tableID)
}

func testAccBigQueryTableWithView(datasetID, tableID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}

resource "google_bigquery_table" "test" {
  table_id   = "%s"
  dataset_id = "${google_bigquery_dataset.test.dataset_id}"

  time_partitioning {
    type = "DAY"
  }

  view {
  	query = "SELECT state FROM [lookerdata:cdc.project_tycho_reports]"
  	use_legacy_sql = true
  }
}`, datasetID, tableID)
}

func testAccBigQueryTableWithNewSqlView(datasetID, tableID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}

resource "google_bigquery_table" "test" {
  table_id   = "%s"
  dataset_id = "${google_bigquery_dataset.test.dataset_id}"

  time_partitioning {
    type = "DAY"
  }

  view {
  	query = "%s"
  	use_legacy_sql = false
  }
}`, datasetID, tableID, "SELECT state FROM `lookerdata.cdc.project_tycho_reports`")
}

func testAccBigQueryTableUpdated(datasetID, tableID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}

resource "google_bigquery_table" "test" {
  table_id   = "%s"
  dataset_id = "${google_bigquery_dataset.test.dataset_id}"

  time_partitioning {
    type = "DAY"
  }

  schema = <<EOH
[
  {
    "name": "city",
    "type": "RECORD",
    "fields": [
      {
        "name": "id",
        "type": "INTEGER"
      },
      {
        "name": "coord",
        "type": "RECORD",
        "fields": [
          {
            "name": "lon",
            "type": "FLOAT"
          },
          {
            "name": "lat",
            "type": "FLOAT"
          }
        ]
      }
    ]
  },
  {
    "name": "country",
    "type": "RECORD",
    "fields": [
      {
        "name": "id",
        "type": "INTEGER"
      },
      {
        "name": "name",
        "type": "STRING"
      }
    ]
  }
]
EOH
}`, datasetID, tableID)
}
