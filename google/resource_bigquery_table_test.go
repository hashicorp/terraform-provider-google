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

func TestAccBigQueryTable_Kms(t *testing.T) {
	t.Parallel()
	resourceName := "google_bigquery_table.test"
	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(10))
	projectId := "terraform-" + acctest.RandString(10)
	projectOrg := getTestOrgFromEnv(t)
	projectBillingAccount := getTestBillingAccountFromEnv(t)
	keyRingName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	cryptoKeyName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigQueryTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableKms(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName, datasetID, tableID),
			},
			{
				ResourceName:      resourceName,
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

func TestAccBigQueryTable_ViewWithLegacySQL(t *testing.T) {
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

func testAccBigQueryTableKms(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName, datasetID, tableID string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
	name            = "%s"
	project_id      = "%s"
	org_id          = "%s"
	billing_account = "%s"
}

resource "google_project_services" "acceptance" {
	project = "${google_project.acceptance.project_id}"

	services = [
	  "cloudkms.googleapis.com",
		"bigquery-json.googleapis.com",
	]
}

resource "google_kms_key_ring" "key_ring" {
	project  = "${google_project_services.acceptance.project}"
	name     = "%s"
	location = "us-central1"
}

resource "google_kms_crypto_key" "crypto_key" {
	name            = "%s"
	key_ring        = "${google_kms_key_ring.key_ring.self_link}"
	rotation_period = "1000000s"
}

resource "google_bigquery_dataset" "test" {
	project = "${google_project.acceptance.project_id}"
  dataset_id = "%s"
	depends_on = ["google_project_services.acceptance"]
}

data "google_bigquery_default_service_account" "acct" {
	project = "${google_project_services.acceptance.project}"
}

resource "google_kms_crypto_key_iam_member" "allow" {
	crypto_key_id = "${google_kms_crypto_key.crypto_key.self_link}"
	role = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
	member = "serviceAccount:${data.google_bigquery_default_service_account.acct.email}"
	depends_on = ["google_bigquery_dataset.test"]
}

resource "google_bigquery_table" "test" {
	project = "${google_project.acceptance.project_id}"
  table_id   = "%s"
  dataset_id = "${google_bigquery_dataset.test.dataset_id}"

  time_partitioning {
    type = "DAY"
    field = "ts"	
  }

	encryption_configuration {
		kms_key_name = "${google_kms_crypto_key.crypto_key.self_link}"
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

  depends_on = ["google_kms_crypto_key_iam_member.allow"]
}
`, projectId, projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName, datasetID, tableID)
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
}`, datasetID, tableID, "SELECT state FROM `lookerdata:cdc.project_tycho_reports`")
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
