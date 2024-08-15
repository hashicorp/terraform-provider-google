// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package biglake_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccBiglakeTable_biglakeTable_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBiglakeTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBiglakeTable_biglakeTableExample(context),
			},
			{
				ResourceName:            "google_biglake_table.table",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "database"},
			},
			{
				Config: testAccBiglakeTable_biglakeTable_update(context),
			},
			{
				ResourceName:            "google_biglake_table.table",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "database"},
			},
		},
	})
}

func testAccBiglakeTable_biglakeTable_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_biglake_catalog" "catalog" {
	name = "tf_test_my_catalog%{random_suffix}"
	location = "US"
}
resource "google_storage_bucket" "bucket" {
	name                        = "tf_test_my_bucket%{random_suffix}"
	location                    = "US"
	force_destroy               = true
	uniform_bucket_level_access = true
}
resource "google_storage_bucket_object" "metadata_folder" {
	name    = "metadata/"
	content = " "
	bucket  = google_storage_bucket.bucket.name
}
resource "google_storage_bucket_object" "data_folder" {
	name    = "data/"
	content = " "
	bucket  = google_storage_bucket.bucket.name
}
resource "google_biglake_database" "database" {
	name = "tf_test_my_database%{random_suffix}"
	catalog = google_biglake_catalog.catalog.id
	type = "HIVE"
	hive_options {
		location_uri = "gs://${google_storage_bucket.bucket.name}/${google_storage_bucket_object.metadata_folder.name}"
		parameters = {
			"owner" = "Alex"
		}
	}
}
resource "google_biglake_table" "table" {
    name = "tf_test_my_table%{random_suffix}"
    database = google_biglake_database.database.id
    type = "HIVE"
    hive_options {
		table_type = "EXTERNAL_TABLE"
		storage_descriptor {
		  location_uri = "gs://${google_storage_bucket.bucket.name}/${google_storage_bucket_object.data_folder.name}/data"
		  input_format = "org.apache.hadoop.mapred.SequenceFileInputFormat2"
		  output_format =  "org.apache.hadoop.hive.ql.io.HiveSequenceFileOutputFormat2"
		}
		# Some Example Parameters.
		parameters = {
		  # Bump the version.
		  "spark.sql.create.version" = "3.1.7"
		  "spark.sql.sources.schema.numParts" = "1"
		  # Update the time.
		  "transient_lastDdlTime" = "1680895000"
		  "spark.sql.partitionProvider" = "catalog"
		  # Change The Name
		  "owner" = "Dana"
		  "spark.sql.sources.schema.part.0" = "{\"type\":\"struct\",\"fields\":[{\"name\":\"id\",\"type\":\"integer\",\"nullable\":true,\"metadata\":{}},{\"name\":\"name\",\"type\":\"string\",\"nullable\":true,\"metadata\":{}},{\"name\":\"age\",\"type\":\"integer\",\"nullable\":true,\"metadata\":{}}]}"
		  "spark.sql.sources.provider": "iceberg"
		  "provider" = "iceberg"
		}
	}
}
`, context)
}
