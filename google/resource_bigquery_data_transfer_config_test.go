package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// The service account TF uses needs the permission granted in the configs
// but it will get deleted by parallel tests, so they need to be ran serially.
func TestAccBigqueryDataTransferConfig(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		"basic":        testAccBigqueryDataTransferConfig_scheduledQuery_basic,
		"update":       testAccBigqueryDataTransferConfig_scheduledQuery_update,
		"booleanParam": testAccBigqueryDataTransferConfig_copy_booleanParam,
	}

	for name, tc := range testCases {
		// shadow the tc variable into scope so that when
		// the loop continues, if t.Run hasn't executed tc(t)
		// yet, we don't have a race condition
		// see https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		tc := tc
		t.Run(name, func(t *testing.T) {
			tc(t)
		})
	}
}

func testAccBigqueryDataTransferConfig_scheduledQuery_basic(t *testing.T) {
	random_suffix := randString(t, 10)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigqueryDataTransferConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryDataTransferConfig_scheduledQuery(random_suffix, "third", "y"),
			},
			{
				ResourceName:            "google_bigquery_data_transfer_config.query_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
		},
	})
}

func testAccBigqueryDataTransferConfig_scheduledQuery_update(t *testing.T) {
	random_suffix := randString(t, 10)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigqueryDataTransferConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryDataTransferConfig_scheduledQuery(random_suffix, "first", "y"),
			},
			{
				Config: testAccBigqueryDataTransferConfig_scheduledQuery(random_suffix, "second", "z"),
			},
			{
				ResourceName:            "google_bigquery_data_transfer_config.query_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
		},
	})
}

func testAccBigqueryDataTransferConfig_copy_booleanParam(t *testing.T) {
	random_suffix := randString(t, 10)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigqueryDataTransferConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryDataTransferConfig_booleanParam(random_suffix),
			},
			{
				ResourceName:            "google_bigquery_data_transfer_config.copy_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
		},
	})
}

func testAccCheckBigqueryDataTransferConfigDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_bigquery_data_transfer_config" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := googleProviderConfig(t)

			url, err := replaceVarsForTest(config, rs, "{{BigqueryDataTransferBasePath}}{{name}}")
			if err != nil {
				return err
			}

			_, err = sendRequest(config, "GET", "", url, nil)
			if err == nil {
				return fmt.Errorf("BigqueryDataTransferConfig still exists at %s", url)
			}
		}

		return nil
	}
}

func testAccBigqueryDataTransferConfig_scheduledQuery(random_suffix, schedule, letter string) string {
	return fmt.Sprintf(`
data "google_project" "project" {}

resource "google_project_iam_member" "permissions" {
  role   = "roles/iam.serviceAccountShortTermTokenMinter"
  member = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-bigquerydatatransfer.iam.gserviceaccount.com"
}


resource "google_bigquery_dataset" "my_dataset" {
  depends_on = [google_project_iam_member.permissions]

  dataset_id    = "my_dataset%s"
  friendly_name = "foo"
  description   = "bar"
  location      = "asia-northeast1"
}

resource "google_bigquery_data_transfer_config" "query_config" {
  depends_on = [google_project_iam_member.permissions]

  display_name           = "my-query-%s"
  location               = "asia-northeast1"
  data_source_id         = "scheduled_query"
  schedule               = "%s sunday of quarter 00:00"
  destination_dataset_id = google_bigquery_dataset.my_dataset.dataset_id
  params = {
    destination_table_name_template = "my-table"
    write_disposition               = "WRITE_APPEND"
    query                           = "SELECT name FROM tabl WHERE x = '%s'"
  }
}
`, random_suffix, random_suffix, schedule, letter)
}

func testAccBigqueryDataTransferConfig_booleanParam(random_suffix string) string {
	return fmt.Sprintf(`
data "google_project" "project" {}

resource "google_project_iam_member" "permissions" {
  role   = "roles/iam.serviceAccountShortTermTokenMinter"
  member = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-bigquerydatatransfer.iam.gserviceaccount.com"
}

resource "google_bigquery_dataset" "source_dataset" {
  depends_on = [google_project_iam_member.permissions]

  dataset_id    = "source_%s"
  friendly_name = "foo"
  description   = "bar"
  location      = "asia-northeast1"
}

resource "google_bigquery_dataset" "destination_dataset" {
  depends_on = [google_project_iam_member.permissions]

  dataset_id    = "destination_%s"
  friendly_name = "foo"
  description   = "bar"
  location      = "asia-northeast1"
}

resource "google_bigquery_data_transfer_config" "copy_config" {
  depends_on = [google_project_iam_member.permissions]

  location = "asia-northeast1"

  display_name           = "Copy test %s"
  data_source_id         = "cross_region_copy"
  destination_dataset_id = google_bigquery_dataset.destination_dataset.dataset_id
  params = {
    overwrite_destination_table = "true"
    source_dataset_id           = google_bigquery_dataset.source_dataset.dataset_id
    source_project_id           = data.google_project.project.project_id
  }
}
`, random_suffix, random_suffix, random_suffix)
}
