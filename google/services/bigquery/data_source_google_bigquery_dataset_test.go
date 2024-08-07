// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package bigquery_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleBigqueryDataset_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryDatasetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleBigqueryDataset_basic(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_bigquery_dataset.bar", "google_bigquery_dataset.foo"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleBigqueryDataset_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`

  resource "google_bigquery_dataset" "foo" {
    dataset_id                  = "tf_test_ds_%{random_suffix}"
    friendly_name               = "testing"
    description                 = "This is a test description"
    location                    = "US"
    default_table_expiration_ms = 3600000
  }

  data "google_bigquery_dataset" "bar" {
    dataset_id    = google_bigquery_dataset.foo.dataset_id
  }
`, context)
}
