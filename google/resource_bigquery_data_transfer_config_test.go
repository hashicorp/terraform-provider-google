package google

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// The service account TF uses needs the permission granted in the configs
// but it will get deleted by parallel tests, so they need to be run serially.
func TestAccBigqueryDataTransferConfig(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		"basic":           testAccBigqueryDataTransferConfig_scheduledQuery_basic,
		"update":          testAccBigqueryDataTransferConfig_scheduledQuery_update,
		"service_account": testAccBigqueryDataTransferConfig_scheduledQuery_with_service_account,
		"no_destintation": testAccBigqueryDataTransferConfig_scheduledQuery_no_destination,
		"booleanParam":    testAccBigqueryDataTransferConfig_copy_booleanParam,
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
	// Uses time.Now
	skipIfVcr(t)
	random_suffix := randString(t, 10)
	now := time.Now().UTC()
	start_time := now.Add(1 * time.Hour).Format(time.RFC3339)
	end_time := now.AddDate(0, 1, 0).Format(time.RFC3339)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigqueryDataTransferConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryDataTransferConfig_scheduledQuery(random_suffix, random_suffix, "third", start_time, end_time, "y"),
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
	// Uses time.Now
	skipIfVcr(t)
	random_suffix := randString(t, 10)
	now := time.Now().UTC()
	first_start_time := now.Add(1 * time.Hour).Format(time.RFC3339)
	first_end_time := now.AddDate(0, 1, 0).Format(time.RFC3339)
	second_start_time := now.Add(2 * time.Hour).Format(time.RFC3339)
	second_end_time := now.AddDate(0, 2, 0).Format(time.RFC3339)
	random_suffix2 := randString(t, 10)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigqueryDataTransferConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryDataTransferConfig_scheduledQuery(random_suffix, random_suffix, "first", first_start_time, first_end_time, "y"),
			},
			{
				Config: testAccBigqueryDataTransferConfig_scheduledQuery(random_suffix, random_suffix, "second", second_start_time, second_end_time, "z"),
			},
			{
				ResourceName:            "google_bigquery_data_transfer_config.query_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
			{
				Config: testAccBigqueryDataTransferConfig_scheduledQuery(random_suffix, random_suffix2, "second", second_start_time, second_end_time, "z"),
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

func testAccBigqueryDataTransferConfig_scheduledQuery_no_destination(t *testing.T) {
	// Uses time.Now
	skipIfVcr(t)
	random_suffix := randString(t, 10)
	now := time.Now().UTC()
	start_time := now.Add(1 * time.Hour).Format(time.RFC3339)
	end_time := now.AddDate(0, 1, 0).Format(time.RFC3339)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigqueryDataTransferConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryDataTransferConfig_scheduledQueryNoDestination(random_suffix, "third", start_time, end_time, "y"),
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

func testAccBigqueryDataTransferConfig_scheduledQuery_with_service_account(t *testing.T) {
	random_suffix := randString(t, 10)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigqueryDataTransferConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryDataTransferConfig_scheduledQuery_service_account(random_suffix),
			},
			{
				ResourceName:            "google_bigquery_data_transfer_config.query_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "service_account_name"},
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

			_, err = sendRequest(config, "GET", "", url, config.userAgent, nil)
			if err == nil {
				return fmt.Errorf("BigqueryDataTransferConfig still exists at %s", url)
			}
		}

		return nil
	}
}

func testAccBigqueryDataTransferConfig_scheduledQuery(random_suffix, random_suffix2, schedule, start_time, end_time, letter string) string {
	return fmt.Sprintf(`
data "google_project" "project" {}

resource "google_project_iam_member" "permissions" {
  project = data.google_project.project.project_id

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

resource "google_pubsub_topic" "my_topic" {
  name = "tf-test-my-topic-%s"
}

resource "google_bigquery_data_transfer_config" "query_config" {
  depends_on = [google_project_iam_member.permissions]

  display_name           = "my-query-%s"
  location               = "asia-northeast1"
  data_source_id         = "scheduled_query"
  schedule               = "%s sunday of quarter 00:00"
  schedule_options {
    disable_auto_scheduling = false
    start_time              = "%s"
    end_time                = "%s"
  }
  destination_dataset_id = google_bigquery_dataset.my_dataset.dataset_id
  notification_pubsub_topic = google_pubsub_topic.my_topic.id
  email_preferences {
    enable_failure_email = true
  }
  params = {
    destination_table_name_template = "my_table"
    write_disposition               = "WRITE_APPEND"
    query                           = "SELECT name FROM tabl WHERE x = '%s'"
  }
}
`, random_suffix, random_suffix, random_suffix2, schedule, start_time, end_time, letter)
}

func testAccBigqueryDataTransferConfig_scheduledQuery_service_account(random_suffix string) string {
	return fmt.Sprintf(`
data "google_project" "project" {}

resource "google_service_account" "bqwriter" {
  account_id = "bqwriter%s"
}

resource "google_project_iam_member" "data_editor" {
  project = data.google_project.project.project_id

  role   = "roles/bigquery.dataEditor"
  member = "serviceAccount:${google_service_account.bqwriter.email}"
}

resource "google_bigquery_dataset" "my_dataset" {
  dataset_id    = "my_dataset%s"
  friendly_name = "foo"
  description   = "bar"
  location      = "asia-northeast1"
}

resource "google_bigquery_data_transfer_config" "query_config" {
  depends_on = [google_project_iam_member.data_editor]

  display_name           = "my-query-%s"
  location               = "asia-northeast1"
  data_source_id         = "scheduled_query"
  schedule               = "every day 00:00"
  destination_dataset_id = google_bigquery_dataset.my_dataset.dataset_id
  service_account_name   = google_service_account.bqwriter.email
  params = {
    destination_table_name_template = "my_table"
    write_disposition               = "WRITE_APPEND"
    query                           = "SELECT 1 AS a"
  }
}
`, random_suffix, random_suffix, random_suffix)
}

func testAccBigqueryDataTransferConfig_scheduledQueryNoDestination(random_suffix, schedule, start_time, end_time, letter string) string {
	return fmt.Sprintf(`
data "google_project" "project" {}

resource "google_project_iam_member" "permissions" {
  project = data.google_project.project.project_id
  role   = "roles/iam.serviceAccountShortTermTokenMinter"
  member = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-bigquerydatatransfer.iam.gserviceaccount.com"
}

resource "google_pubsub_topic" "my_topic" {
  name = "tf-test-my-topic-%s"
}

resource "google_bigquery_data_transfer_config" "query_config" {
  depends_on = [google_project_iam_member.permissions]

  display_name           = "my-query-%s"
  location               = "asia-northeast1"
  data_source_id         = "scheduled_query"
  schedule               = "%s sunday of quarter 00:00"
  schedule_options {
    disable_auto_scheduling = false
    start_time              = "%s"
    end_time                = "%s"
  }
  notification_pubsub_topic = google_pubsub_topic.my_topic.id
  email_preferences {
    enable_failure_email = true
  }
  params = {
    destination_table_name_template = "my_table"
    write_disposition               = "WRITE_APPEND"
    query                           = "SELECT name FROM tabl WHERE x = '%s'"
  }
}
`, random_suffix, random_suffix, schedule, start_time, end_time, letter)
}

func testAccBigqueryDataTransferConfig_booleanParam(random_suffix string) string {
	return fmt.Sprintf(`
data "google_project" "project" {}

resource "google_project_iam_member" "permissions" {
  project = data.google_project.project.project_id
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
