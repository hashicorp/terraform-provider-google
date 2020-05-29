package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccBigQueryTable_Basic(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", randString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigQueryTableDestroyProducer(t),
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
	datasetID := fmt.Sprintf("tf_test_%s", randString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", randString(t, 10))
	kms := BootstrapKMSKey(t)
	cryptoKeyName := kms.CryptoKey.Name

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableKms(cryptoKeyName, datasetID, tableID),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigQueryTable_HivePartitioning(t *testing.T) {
	t.Parallel()
	bucketName := testBucketName(t)
	resourceName := "google_bigquery_table.test"
	datasetID := fmt.Sprintf("tf_test_%s", randString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", randString(t, 10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableHivePartitioning(bucketName, datasetID, tableID),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigQueryTable_RangePartitioning(t *testing.T) {
	t.Parallel()
	resourceName := "google_bigquery_table.test"
	datasetID := fmt.Sprintf("tf_test_%s", randString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableRangePartitioning(datasetID, tableID),
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

	datasetID := fmt.Sprintf("tf_test_%s", randString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigQueryTableDestroyProducer(t),
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

	datasetID := fmt.Sprintf("tf_test_%s", randString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigQueryTableDestroyProducer(t),
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

func TestAccBigQueryExternalDataTable_CSV(t *testing.T) {
	t.Parallel()

	bucketName := testBucketName(t)
	objectName := fmt.Sprintf("tf_test_%s.csv", randString(t, 10))

	datasetID := fmt.Sprintf("tf_test_%s", randString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableFromGCS(datasetID, tableID, bucketName, objectName, TEST_CSV, "CSV", "\\\""),
				Check:  testAccCheckBigQueryExtData(t, "\""),
			},
			{
				Config: testAccBigQueryTableFromGCS(datasetID, tableID, bucketName, objectName, TEST_CSV, "CSV", ""),
				Check:  testAccCheckBigQueryExtData(t, ""),
			},
		},
	})
}

func TestAccBigQueryDataTable_sheet(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableFromSheet(context),
			},
			{
				ResourceName:      "google_bigquery_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckBigQueryExtData(t *testing.T, expectedQuoteChar string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_bigquery_table" {
				continue
			}

			config := googleProviderConfig(t)
			dataset := rs.Primary.Attributes["dataset_id"]
			table := rs.Primary.Attributes["table_id"]
			res, err := config.clientBigQuery.Tables.Get(config.Project, dataset, table).Do()
			if err != nil {
				return err
			}

			if res.Type != "EXTERNAL" {
				return fmt.Errorf("Table \"%s.%s\" is of type \"%s\", expected EXTERNAL.", dataset, table, res.Type)
			}
			edc := res.ExternalDataConfiguration
			cvsOpts := edc.CsvOptions
			if cvsOpts == nil || *cvsOpts.Quote != expectedQuoteChar {
				return fmt.Errorf("Table \"%s.%s\" quote should be '%s' but was '%s'", dataset, table, expectedQuoteChar, *cvsOpts.Quote)
			}
		}
		return nil
	}
}

func testAccCheckBigQueryTableDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_bigquery_table" {
				continue
			}

			config := googleProviderConfig(t)
			_, err := config.clientBigQuery.Tables.Get(config.Project, rs.Primary.Attributes["dataset_id"], rs.Primary.Attributes["table_id"]).Do()
			if err == nil {
				return fmt.Errorf("Table still present")
			}
		}

		return nil
	}
}

func testAccBigQueryTable(datasetID, tableID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
	dataset_id = "%s"
}

resource "google_bigquery_table" "test" {
	table_id   = "%s"
	dataset_id = google_bigquery_dataset.test.dataset_id

	time_partitioning {
		type                     = "DAY"
		field                    = "ts"
		require_partition_filter = true
	}
	clustering = ["some_int", "some_string"]
	schema     = <<EOH
[
	{
		"name": "ts",
		"type": "TIMESTAMP"
	},
	{
		"name": "some_string",
		"type": "STRING"
	},
	{
		"name": "some_int",
		"type": "INTEGER"
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

}
`, datasetID, tableID)
}

func testAccBigQueryTableKms(cryptoKeyName, datasetID, tableID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
		dataset_id = "%s"
}

data "google_bigquery_default_service_account" "acct" {}

resource "google_kms_crypto_key_iam_member" "allow" {
	crypto_key_id = "%s"
	role = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
	member = "serviceAccount:${data.google_bigquery_default_service_account.acct.email}"
	depends_on = ["google_bigquery_dataset.test"]
}

resource "google_bigquery_table" "test" {
	table_id   = "%s"
	dataset_id = "${google_bigquery_dataset.test.dataset_id}"

	time_partitioning {
		type = "DAY"
		field = "ts"
	}

	encryption_configuration {
		kms_key_name = "${google_kms_crypto_key_iam_member.allow.crypto_key_id}"
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
}
`, datasetID, cryptoKeyName, tableID)
}

func testAccBigQueryTableHivePartitioning(bucketName, datasetID, tableID string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "test" {
	name          = "%s"
	force_destroy = true
}

resource "google_storage_bucket_object" "test" {
	name    = "key1=20200330/init.csv"
	content = ";"
	bucket  = google_storage_bucket.test.name
}

resource "google_bigquery_dataset" "test" {
        dataset_id = "%s"
}

resource "google_bigquery_table" "test" {
	table_id   = "%s"
	dataset_id = google_bigquery_dataset.test.dataset_id

	external_data_configuration {
            source_format = "CSV"
            autodetect = true
            source_uris= ["gs://${google_storage_bucket.test.name}/*"]

            hive_partitioning_options {
                mode = "AUTO"
                source_uri_prefix = "gs://${google_storage_bucket.test.name}/"
	    }

        }
	depends_on = ["google_storage_bucket_object.test"]
}
`, bucketName, datasetID, tableID)
}

func testAccBigQueryTableRangePartitioning(datasetID, tableID string) string {
	return fmt.Sprintf(`
	resource "google_bigquery_dataset" "test" {
		dataset_id = "%s"
	}

	resource "google_bigquery_table" "test" {
		table_id   = "%s"
		dataset_id = google_bigquery_dataset.test.dataset_id

		range_partitioning {
			field = "id"
			range {
				start    = 1
				end      = 10000
				interval = 100
			}
		}

		schema = <<EOH
[
	{
		"name": "ts",
		"type": "TIMESTAMP"
	},
	{
		"name": "id",
		"type": "INTEGER"
	}
]
EOH
}
	`, datasetID, tableID)
}

func testAccBigQueryTableWithView(datasetID, tableID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
	dataset_id = "%s"
}

resource "google_bigquery_table" "test" {
	table_id   = "%s"
	dataset_id = google_bigquery_dataset.test.dataset_id

	time_partitioning {
		type = "DAY"
	}

	view {
		query          = "SELECT state FROM [lookerdata:cdc.project_tycho_reports]"
		use_legacy_sql = true
	}
}
`, datasetID, tableID)
}

func testAccBigQueryTableWithNewSqlView(datasetID, tableID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
	dataset_id = "%s"
}

resource "google_bigquery_table" "test" {
	table_id   = "%s"
	dataset_id = google_bigquery_dataset.test.dataset_id

	time_partitioning {
		type = "DAY"
	}

	view {
		query          = "%s"
		use_legacy_sql = false
	}
}
`, datasetID, tableID, "SELECT state FROM `lookerdata.cdc.project_tycho_reports`")
}

func testAccBigQueryTableUpdated(datasetID, tableID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
	dataset_id = "%s"
}

resource "google_bigquery_table" "test" {
	table_id   = "%s"
	dataset_id = google_bigquery_dataset.test.dataset_id

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

}
`, datasetID, tableID)
}

func testAccBigQueryTableFromGCS(datasetID, tableID, bucketName, objectName, content, format, quoteChar string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
	dataset_id = "%s"
}

resource "google_storage_bucket" "test" {
	name          = "%s"
	force_destroy = true
}

resource "google_storage_bucket_object" "test" {
	name    = "%s"
	content = <<EOF
%s
EOF

	bucket = google_storage_bucket.test.name
}

resource "google_bigquery_table" "test" {
	table_id   = "%s"
	dataset_id = google_bigquery_dataset.test.dataset_id
	external_data_configuration {
		autodetect    = true
		source_format = "%s"
		csv_options {
			encoding = "UTF-8"
			quote    = "%s"
		}

		source_uris = [
			"gs://${google_storage_bucket.test.name}/${google_storage_bucket_object.test.name}",
		]
	}
}
`, datasetID, bucketName, objectName, content, tableID, format, quoteChar)
}

func testAccBigQueryTableFromSheet(context map[string]interface{}) string {
	return Nprintf(`
	resource "google_bigquery_table" "table" {
		dataset_id = google_bigquery_dataset.dataset.dataset_id
		table_id   = "tf_test_sheet_%{random_suffix}"

		external_data_configuration {
		  autodetect            = true
		  source_format         = "GOOGLE_SHEETS"
		  ignore_unknown_values = true

		  google_sheets_options {
			skip_leading_rows = 1
		  }

		  source_uris = [
			"https://drive.google.com/open?id=xxxx",
		  ]
		}

		schema = <<EOF
	  [
		{
		  "name": "permalink",
		  "type": "STRING",
		  "mode": "NULLABLE",
		  "description": "The Permalink"
		},
		{
		  "name": "state",
		  "type": "STRING",
		  "mode": "NULLABLE",
		  "description": "State where the head office is located"
		}
	  ]
	  EOF
	  }

	  resource "google_bigquery_dataset" "dataset" {
		dataset_id                  = "tf_test_ds_%{random_suffix}"
		friendly_name               = "test"
		description                 = "This is a test description"
		location                    = "EU"
		default_table_expiration_ms = 3600000

		labels = {
		  env = "default"
		}
	  }
`, context)
}

var TEST_CSV = `lifelock,LifeLock,,web,Tempe,AZ,1-May-07,6850000,USD,b
lifelock,LifeLock,,web,Tempe,AZ,1-Oct-06,6000000,USD,a
lifelock,LifeLock,,web,Tempe,AZ,1-Jan-08,25000000,USD,c
mycityfaces,MyCityFaces,7,web,Scottsdale,AZ,1-Jan-08,50000,USD,seed
flypaper,Flypaper,,web,Phoenix,AZ,1-Feb-08,3000000,USD,a
infusionsoft,Infusionsoft,105,software,Gilbert,AZ,1-Oct-07,9000000,USD,a
gauto,gAuto,4,web,Scottsdale,AZ,1-Jan-08,250000,USD,seed
chosenlist-com,ChosenList.com,5,web,Scottsdale,AZ,1-Oct-06,140000,USD,seed
chosenlist-com,ChosenList.com,5,web,Scottsdale,AZ,25-Jan-08,233750,USD,angel
`
