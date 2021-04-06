package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBigQueryJob_withLocation(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
		"location":      "asia-northeast1",
	}

	// Need to construct the import ID manually since the state ID will not contain the location
	importID := fmt.Sprintf("projects/%s/jobs/tf_test_job_query%s/location/%s", getTestProjectFromEnv(), context["random_suffix"], context["location"])

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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
	return Nprintf(`
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
