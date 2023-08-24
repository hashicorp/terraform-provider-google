// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package bigquerydatatransfer_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/services/bigquerydatatransfer"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestBigqueryDataTransferConfig_resourceBigqueryDTCParamsCustomDiffFuncForceNew(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		before   map[string]interface{}
		after    map[string]interface{}
		forcenew bool
	}{
		"changing_data_path_template": {
			before: map[string]interface{}{
				"data_source_id": "google_cloud_storage",
				"params": map[string]interface{}{
					"data_path_template":              "gs://bq-bucket-temp/*.json",
					"destination_table_name_template": "table-old",
					"file_format":                     "JSON",
					"max_bad_records":                 10,
					"write_disposition":               "APPEND",
				},
			},
			after: map[string]interface{}{
				"data_source_id": "google_cloud_storage",
				"params": map[string]interface{}{
					"data_path_template":              "gs://bq-bucket-temp-new/*.json",
					"destination_table_name_template": "table-old",
					"file_format":                     "JSON",
					"max_bad_records":                 10,
					"write_disposition":               "APPEND",
				},
			},
			forcenew: true,
		},
		"changing_destination_table_name_template": {
			before: map[string]interface{}{
				"data_source_id": "google_cloud_storage",
				"params": map[string]interface{}{
					"data_path_template":              "gs://bq-bucket-temp/*.json",
					"destination_table_name_template": "table-old",
					"file_format":                     "JSON",
					"max_bad_records":                 10,
					"write_disposition":               "APPEND",
				},
			},
			after: map[string]interface{}{
				"data_source_id": "google_cloud_storage",
				"params": map[string]interface{}{
					"data_path_template":              "gs://bq-bucket-temp/*.json",
					"destination_table_name_template": "table-new",
					"file_format":                     "JSON",
					"max_bad_records":                 10,
					"write_disposition":               "APPEND",
				},
			},
			forcenew: true,
		},
		"changing_non_force_new_fields": {
			before: map[string]interface{}{
				"data_source_id": "google_cloud_storage",
				"params": map[string]interface{}{
					"data_path_template":              "gs://bq-bucket-temp/*.json",
					"destination_table_name_template": "table-old",
					"file_format":                     "JSON",
					"max_bad_records":                 10,
					"write_disposition":               "APPEND",
				},
			},
			after: map[string]interface{}{
				"data_source_id": "google_cloud_storage",
				"params": map[string]interface{}{
					"data_path_template":              "gs://bq-bucket-temp/*.json",
					"destination_table_name_template": "table-old",
					"file_format":                     "JSON",
					"max_bad_records":                 1000,
					"write_disposition":               "APPEND",
				},
			},
			forcenew: false,
		},
		"changing_destination_table_name_template_for_different_data_source_id": {
			before: map[string]interface{}{
				"data_source_id": "scheduled_query",
				"params": map[string]interface{}{
					"destination_table_name_template": "table-old",
					"query":                           "SELECT 1 AS a",
					"write_disposition":               "WRITE_APPEND",
				},
			},
			after: map[string]interface{}{
				"data_source_id": "scheduled_query",
				"params": map[string]interface{}{
					"destination_table_name_template": "table-new",
					"query":                           "SELECT 1 AS a",
					"write_disposition":               "WRITE_APPEND",
				},
			},
			forcenew: false,
		},
		"changing_data_path_template_for_different_data_source_id": {
			before: map[string]interface{}{
				"data_source_id": "scheduled_query",
				"params": map[string]interface{}{
					"data_path_template": "gs://bq-bucket/*.json",
					"query":              "SELECT 1 AS a",
					"write_disposition":  "WRITE_APPEND",
				},
			},
			after: map[string]interface{}{
				"data_source_id": "scheduled_query",
				"params": map[string]interface{}{
					"data_path_template": "gs://bq-bucket-new/*.json",
					"query":              "SELECT 1 AS a",
					"write_disposition":  "WRITE_APPEND",
				},
			},
			forcenew: false,
		},
	}

	for tn, tc := range cases {
		d := &tpgresource.ResourceDiffMock{
			Before: map[string]interface{}{
				"params":         tc.before["params"],
				"data_source_id": tc.before["data_source_id"],
			},
			After: map[string]interface{}{
				"params":         tc.after["params"],
				"data_source_id": tc.after["data_source_id"],
			},
		}
		err := bigquerydatatransfer.ParamsCustomizeDiffFunc(d)
		if err != nil {
			t.Errorf("failed, expected no error but received - %s for the condition %s", err, tn)
		}
		if d.IsForceNew != tc.forcenew {
			t.Errorf("ForceNew not setup correctly for the condition-'%s', expected:%v; actual:%v", tn, tc.forcenew, d.IsForceNew)
		}
	}
}

// The service account TF uses needs the permission granted in the configs
// but it will get deleted by parallel tests, so they need to be run serially.
func TestAccBigqueryDataTransferConfig(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		"basic":                  testAccBigqueryDataTransferConfig_scheduledQuery_basic,
		"update":                 testAccBigqueryDataTransferConfig_scheduledQuery_update,
		"service_account":        testAccBigqueryDataTransferConfig_scheduledQuery_with_service_account,
		"no_destintation":        testAccBigqueryDataTransferConfig_scheduledQuery_no_destination,
		"booleanParam":           testAccBigqueryDataTransferConfig_copy_booleanParam,
		"update_params":          testAccBigqueryDataTransferConfig_force_new_update_params,
		"update_service_account": testAccBigqueryDataTransferConfig_scheduledQuery_update_service_account,
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
	acctest.SkipIfVcr(t)
	random_suffix := acctest.RandString(t, 10)
	now := time.Now().UTC()
	start_time := now.Add(1 * time.Hour).Format(time.RFC3339)
	end_time := now.AddDate(0, 1, 0).Format(time.RFC3339)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigqueryDataTransferConfigDestroyProducer(t),
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
	acctest.SkipIfVcr(t)
	random_suffix := acctest.RandString(t, 10)
	now := time.Now().UTC()
	first_start_time := now.Add(1 * time.Hour).Format(time.RFC3339)
	first_end_time := now.AddDate(0, 1, 0).Format(time.RFC3339)
	second_start_time := now.Add(2 * time.Hour).Format(time.RFC3339)
	second_end_time := now.AddDate(0, 2, 0).Format(time.RFC3339)
	random_suffix2 := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigqueryDataTransferConfigDestroyProducer(t),
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
	acctest.SkipIfVcr(t)
	random_suffix := acctest.RandString(t, 10)
	now := time.Now().UTC()
	start_time := now.Add(1 * time.Hour).Format(time.RFC3339)
	end_time := now.AddDate(0, 1, 0).Format(time.RFC3339)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigqueryDataTransferConfigDestroyProducer(t),
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
	random_suffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigqueryDataTransferConfigDestroyProducer(t),
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
	random_suffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigqueryDataTransferConfigDestroyProducer(t),
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

func testAccBigqueryDataTransferConfig_force_new_update_params(t *testing.T) {
	random_suffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigqueryDataTransferConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryDataTransferConfig_update_params_force_new(random_suffix, "old", "old"),
			},
			{
				ResourceName:            "google_bigquery_data_transfer_config.update_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
			{
				Config: testAccBigqueryDataTransferConfig_update_params_force_new(random_suffix, "new", "old"),
			},
			{
				ResourceName:            "google_bigquery_data_transfer_config.update_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
			{
				Config: testAccBigqueryDataTransferConfig_update_params_force_new(random_suffix, "new", "new"),
			},
			{
				ResourceName:            "google_bigquery_data_transfer_config.update_config",
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

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{BigqueryDataTransferBasePath}}{{name}}")
			if err != nil {
				return err
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("BigqueryDataTransferConfig still exists at %s", url)
			}
		}

		return nil
	}
}

func testAccBigqueryDataTransferConfig_scheduledQuery_update_service_account(t *testing.T) {
	random_suffix1 := acctest.RandString(t, 10)
	random_suffix2 := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigqueryDataTransferConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryDataTransferConfig_scheduledQuery_updateServiceAccount(random_suffix1, random_suffix1),
			},
			{
				ResourceName:            "google_bigquery_data_transfer_config.query_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "service_account_name"},
			},
			{
				Config: testAccBigqueryDataTransferConfig_scheduledQuery_updateServiceAccount(random_suffix1, random_suffix2),
				Check:  testAccCheckDataTransferServiceAccountNamePrefix("google_bigquery_data_transfer_config.query_config", random_suffix2),
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

// Check if transfer config service account name starts with given prefix
func testAccCheckDataTransferServiceAccountNamePrefix(resourceName string, prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if !strings.HasPrefix(rs.Primary.Attributes["service_account_name"], "bqwriter"+prefix) {
			return fmt.Errorf("Transfer config service account not updated")
		}

		return nil
	}
}

func testAccBigqueryDataTransferConfig_scheduledQuery(random_suffix, random_suffix2, schedule, start_time, end_time, letter string) string {
	return fmt.Sprintf(`
data "google_project" "project" {}

resource "google_project_iam_member" "permissions" {
  project = data.google_project.project.project_id

  role   = "roles/iam.serviceAccountTokenCreator"
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
  role   = "roles/iam.serviceAccountTokenCreator"
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
  role   = "roles/iam.serviceAccountTokenCreator"
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

func testAccBigqueryDataTransferConfig_update_params_force_new(random_suffix, path, table string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "dataset" {
  dataset_id       = "tf_test_%s"
  friendly_name    = "foo"
  description      = "bar"
  location         = "US"
}

resource "google_bigquery_data_transfer_config" "update_config" {
  display_name           = "tf-test-%s"
  data_source_id         = "google_cloud_storage"
  destination_dataset_id = google_bigquery_dataset.dataset.dataset_id
  location               = google_bigquery_dataset.dataset.location

  params = {
    data_path_template              = "gs://bq-bucket-%s-%s/*.json"
    destination_table_name_template = "the-table-%s-%s"
    file_format                     = "JSON"
    max_bad_records                 = 0
    write_disposition               = "APPEND"
  }
}
`, random_suffix, random_suffix, random_suffix, path, random_suffix, table)
}

func testAccBigqueryDataTransferConfig_scheduledQuery_updateServiceAccount(random_suffix string, service_account string) string {
	return fmt.Sprintf(`
data "google_project" "project" {}

resource "google_service_account" "bqwriter%s" {
  account_id = "bqwriter%s"
}

resource "google_project_iam_member" "data_editor" {
  project = data.google_project.project.project_id

  role   = "roles/bigquery.dataEditor"
  member = "serviceAccount:${google_service_account.bqwriter%s.email}"
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
  schedule               = "every 15 minutes"
  destination_dataset_id = google_bigquery_dataset.my_dataset.dataset_id
  service_account_name   = google_service_account.bqwriter%s.email
  params = {
    destination_table_name_template = "my_table"
    write_disposition               = "WRITE_APPEND"
    query                           = "SELECT 1 AS a"
  }
}
`, service_account, service_account, service_account, random_suffix, random_suffix, service_account)
}
