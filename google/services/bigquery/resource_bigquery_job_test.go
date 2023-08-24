// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package bigquery_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccBigQueryJob_withLocation(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"location":      "asia-northeast1",
	}

	// Need to construct the import ID manually since the state ID will not contain the location
	importID := fmt.Sprintf("projects/%s/jobs/tf_test_job_query%s/location/%s", envvar.GetTestProjectFromEnv(), context["random_suffix"], context["location"])

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryJob_withLocation(context),
			},
			{
				ResourceName:            "google_bigquery_job.job",
				ImportStateId:           importID,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "status.0.state"},
			},
		},
	})
}

func testAccBigQueryJob_withLocation(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_table" "foo" {
  deletion_protection = false
  dataset_id = google_bigquery_dataset.bar.dataset_id
  table_id   = "tf_test_job_query%{random_suffix}_table"
}

resource "google_bigquery_dataset" "bar" {
  dataset_id                  = "tf_test_job_query%{random_suffix}_dataset"
  friendly_name               = "test"
  description                 = "This is a test description"
  location                    = "%{location}"
}

resource "google_bigquery_job" "job" {
  job_id     = "tf_test_job_query%{random_suffix}"

  labels = {
    "example-label" ="example-value"
  }

  query {
    query = "SELECT state FROM [lookerdata:cdc.project_tycho_reports]"

    destination_table {
      project_id = google_bigquery_table.foo.project
      dataset_id = google_bigquery_table.foo.dataset_id
      table_id   = google_bigquery_table.foo.table_id
    }

    allow_large_results = true
    flatten_results = true

    script_options {
      key_result_statement = "LAST"
    }
  }

  location = "%{location}"
}
`, context)
}
