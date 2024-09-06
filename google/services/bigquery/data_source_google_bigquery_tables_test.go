// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package bigquery_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleBigqueryTables_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleBigqueryTables_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_bigquery_tables.example", "tables.#", "1"),
					resource.TestCheckResourceAttr("data.google_bigquery_tables.example", "tables.0.table_id", fmt.Sprintf("tf_test_table_%s", context["random_suffix"])),
					resource.TestCheckResourceAttr("data.google_bigquery_tables.example", "tables.0.labels.%", "1"),
					resource.TestCheckResourceAttr("data.google_bigquery_tables.example", "tables.0.labels.goog-terraform-provisioned", "true"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleBigqueryTables_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
  
  resource "google_bigquery_dataset" "test" {
    dataset_id                  = "tf_test_ds_%{random_suffix}"
    friendly_name               = "testing"
    description                 = "This is a test description"
    location                    = "US"
    default_table_expiration_ms = 3600000
  }

  resource "google_bigquery_table" "test" {
    dataset_id        = google_bigquery_dataset.test.dataset_id
    table_id          = "tf_test_table_%{random_suffix}"
    deletion_protection = false
    schema     = <<EOF
    [
      {
        "name": "name",
        "type": "STRING",
        "mode": "NULLABLE"
      }
    ]
    EOF
  }

  data "google_bigquery_tables" "example" {
    dataset_id = google_bigquery_table.test.dataset_id
  }
`, context)
}
