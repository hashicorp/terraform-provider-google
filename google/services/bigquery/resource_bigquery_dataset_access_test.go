// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package bigquery_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccBigQueryDatasetAccess_basic(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	saID := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	expected := map[string]interface{}{
		"role":        "OWNER",
		"userByEmail": fmt.Sprintf("%s@%s.iam.gserviceaccount.com", saID, envvar.GetTestProjectFromEnv()),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryDatasetAccess_basic(datasetID, saID),
				Check:  testAccCheckBigQueryDatasetAccessPresent(t, "google_bigquery_dataset.dataset", expected),
			},
			{
				// Destroy step instead of CheckDestroy so we can check the access is removed without deleting the dataset
				Config: testAccBigQueryDatasetAccess_destroy(datasetID, "dataset"),
				Check:  testAccCheckBigQueryDatasetAccessAbsent(t, "google_bigquery_dataset.dataset", expected),
			},
		},
	})
}

func TestAccBigQueryDatasetAccess_view(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	datasetID2 := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	expected := map[string]interface{}{
		"view": map[string]interface{}{
			"projectId": envvar.GetTestProjectFromEnv(),
			"datasetId": datasetID2,
			"tableId":   tableID,
		},
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryDatasetAccess_view(datasetID, datasetID2, tableID),
				Check:  testAccCheckBigQueryDatasetAccessPresent(t, "google_bigquery_dataset.private", expected),
			},
			{
				Config: testAccBigQueryDatasetAccess_destroy(datasetID, "private"),
				Check:  testAccCheckBigQueryDatasetAccessAbsent(t, "google_bigquery_dataset.private", expected),
			},
		},
	})
}

func TestAccBigQueryDatasetAccess_authorizedDataset(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	datasetID2 := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	expected := map[string]interface{}{
		"dataset": map[string]interface{}{
			"dataset": map[string]interface{}{
				"projectId": envvar.GetTestProjectFromEnv(),
				"datasetId": datasetID2,
			},
			"targetTypes": []interface{}{"VIEWS"},
		},
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryDatasetAccess_authorizedDataset(datasetID, datasetID2),
				Check:  testAccCheckBigQueryDatasetAccessPresent(t, "google_bigquery_dataset.private", expected),
			},
			{
				Config: testAccBigQueryDatasetAccess_destroy(datasetID, "private"),
				Check:  testAccCheckBigQueryDatasetAccessAbsent(t, "google_bigquery_dataset.private", expected),
			},
		},
	})
}

func TestAccBigQueryDatasetAccess_authorizedRoutine(t *testing.T) {
	// Multiple fine-grained resources
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"public_dataset":  fmt.Sprintf("tf_test_public_dataset_%s", acctest.RandString(t, 10)),
		"public_routine":  fmt.Sprintf("tf_test_public_routine_%s", acctest.RandString(t, 10)),
		"private_dataset": fmt.Sprintf("tf_test_private_dataset_%s", acctest.RandString(t, 10)),
	}

	expected := map[string]interface{}{
		"routine": map[string]interface{}{
			"projectId": envvar.GetTestProjectFromEnv(),
			"datasetId": context["public_dataset"],
			"routineId": context["public_routine"],
		},
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryDatasetAccess_authorizedRoutine(context),
				Check:  testAccCheckBigQueryDatasetAccessPresent(t, "google_bigquery_dataset.private", expected),
			},
			{
				// Destroy step instead of CheckDestroy so we can check the access is removed without deleting the dataset
				Config: testAccBigQueryDatasetAccess_destroy(context["private_dataset"].(string), "private"),
				Check:  testAccCheckBigQueryDatasetAccessAbsent(t, "google_bigquery_dataset.private", expected),
			},
		},
	})
}

func TestAccBigQueryDatasetAccess_multiple(t *testing.T) {
	// Multiple fine-grained resources
	acctest.SkipIfVcr(t)
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	expected1 := map[string]interface{}{
		"role":   "WRITER",
		"domain": "google.com",
	}

	expected2 := map[string]interface{}{
		"role":         "READER",
		"specialGroup": "projectWriters",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryDatasetAccess_multiple(datasetID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBigQueryDatasetAccessPresent(t, "google_bigquery_dataset.dataset", expected1),
					testAccCheckBigQueryDatasetAccessPresent(t, "google_bigquery_dataset.dataset", expected2),
				),
			},
			{
				// Destroy step instead of CheckDestroy so we can check the access is removed without deleting the dataset
				Config: testAccBigQueryDatasetAccess_destroy(datasetID, "dataset"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBigQueryDatasetAccessAbsent(t, "google_bigquery_dataset.dataset", expected1),
					testAccCheckBigQueryDatasetAccessAbsent(t, "google_bigquery_dataset.dataset", expected2),
				),
			},
		},
	})
}

func TestAccBigQueryDatasetAccess_predefinedRole(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	expected1 := map[string]interface{}{
		"role":   "WRITER",
		"domain": "google.com",
	}

	expected2 := map[string]interface{}{
		"role":   "READER",
		"domain": "google.com",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryDatasetAccess_predefinedRole("roles/bigquery.dataEditor", datasetID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBigQueryDatasetAccessPresent(t, "google_bigquery_dataset.dataset", expected1),
				),
			},
			{
				// Update role
				Config: testAccBigQueryDatasetAccess_predefinedRole("roles/bigquery.dataViewer", datasetID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBigQueryDatasetAccessPresent(t, "google_bigquery_dataset.dataset", expected2),
				),
			},
			{
				// Destroy step instead of CheckDestroy so we can check the access is removed without deleting the dataset
				Config: testAccBigQueryDatasetAccess_destroy(datasetID, "dataset"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBigQueryDatasetAccessAbsent(t, "google_bigquery_dataset.dataset", expected1),
				),
			},
		},
	})
}

func TestAccBigQueryDatasetAccess_iamMember(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	sinkName := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryDatasetAccess_iamMember(datasetID, sinkName),
			},
		},
	})
}

func TestAccBigQueryDatasetAccess_allUsers(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryDatasetAccess_allUsers(datasetID),
			},
			{
				Config: testAccBigQueryDatasetAccess_allAuthenticatedUsers(datasetID),
			},
		},
	})
}

func TestAccBigQueryDatasetAccess_allAuthenticatedUsers(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryDatasetAccess_allAuthenticatedUsers(datasetID),
			},
		},
	})
}

func testAccCheckBigQueryDatasetAccessPresent(t *testing.T, n string, expected map[string]interface{}) resource.TestCheckFunc {
	return testAccCheckBigQueryDatasetAccess(t, n, expected, true)
}

func testAccCheckBigQueryDatasetAccessAbsent(t *testing.T, n string, expected map[string]interface{}) resource.TestCheckFunc {
	return testAccCheckBigQueryDatasetAccess(t, n, expected, false)
}

func testAccCheckBigQueryDatasetAccess(t *testing.T, n string, expected map[string]interface{}, expectPresent bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		config := acctest.GoogleProviderConfig(t)
		url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{BigQueryBasePath}}projects/{{project}}/datasets/{{dataset_id}}")
		if err != nil {
			return err
		}

		ds, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			RawURL:    url,
			UserAgent: config.UserAgent,
		})
		if err != nil {
			return err
		}
		access := ds["access"].([]interface{})
		for _, a := range access {
			if reflect.DeepEqual(a, expected) {
				if !expectPresent {
					return fmt.Errorf("Found access %+v, expected not present", expected)
				}
				return nil
			}
		}
		if expectPresent {
			return fmt.Errorf("Did not find access %+v, expected present", expected)
		}
		return nil
	}
}

func testAccBigQueryDatasetAccess_destroy(datasetID, rs string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "%s" {
  dataset_id = "%s"
}
`, rs, datasetID)
}

func testAccBigQueryDatasetAccess_basic(datasetID, saID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset_access" "access" {
  dataset_id    = google_bigquery_dataset.dataset.dataset_id
  role          = "OWNER"
  user_by_email = google_service_account.bqowner.email
}

resource "google_bigquery_dataset" "dataset" {
  dataset_id = "%s"
}

resource "google_service_account" "bqowner" {
  account_id = "%s"
}
`, datasetID, saID)
}

func testAccBigQueryDatasetAccess_view(datasetID, datasetID2, tableID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset_access" "access" {
  dataset_id    = google_bigquery_dataset.private.dataset_id
  view {
    project_id = google_bigquery_table.public.project
    dataset_id = google_bigquery_dataset.public.dataset_id
    table_id   = google_bigquery_table.public.table_id
  }
}

resource "google_bigquery_dataset" "private" {
  dataset_id = "%s"
}

resource "google_bigquery_dataset" "public" {
  dataset_id = "%s"
}

resource "google_bigquery_table" "public" {
  deletion_protection = false
  dataset_id = google_bigquery_dataset.public.dataset_id
  table_id   = "%s"

  view {
    query          = "%s"
    use_legacy_sql = false
  }
}

`, datasetID, datasetID2, tableID, "SELECT state FROM `lookerdata.cdc.project_tycho_reports`")
}

func testAccBigQueryDatasetAccess_authorizedDataset(datasetID, datasetID2 string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset_access" "access" {
  dataset_id    = google_bigquery_dataset.private.dataset_id
  dataset {
    dataset{
      project_id = google_bigquery_dataset.public.project
      dataset_id = google_bigquery_dataset.public.dataset_id
    }
    target_types = ["VIEWS"]
  }
}

resource "google_bigquery_dataset" "private" {
  dataset_id = "%s"
}

resource "google_bigquery_dataset" "public" {
  dataset_id = "%s"
}
`, datasetID, datasetID2)
}

func testAccBigQueryDatasetAccess_authorizedRoutine(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_dataset" "public" {
  dataset_id  = "%{public_dataset}"
  description = "This dataset is public"
}

resource "google_bigquery_routine" "public" {
  dataset_id      = google_bigquery_dataset.public.dataset_id
  routine_id      = "%{public_routine}"
  routine_type    = "TABLE_VALUED_FUNCTION"
  language        = "SQL"
  definition_body = <<-EOS
    SELECT 1 + value AS value
  EOS
  arguments {
    name          = "value"
    argument_kind = "FIXED_TYPE"
    data_type     = jsonencode({ "typeKind" = "INT64" })
  }
  return_table_type = jsonencode({ "columns" = [
    { "name" = "value", "type" = { "typeKind" = "INT64" } },
  ] })
}

resource "google_bigquery_dataset" "private" {
  dataset_id  = "%{private_dataset}"
}

resource "google_bigquery_dataset_access" "authorized_routine" {
  dataset_id = google_bigquery_dataset.private.dataset_id
  routine {
    project_id = google_bigquery_routine.public.project
    dataset_id = google_bigquery_routine.public.dataset_id
    routine_id = google_bigquery_routine.public.routine_id
  }
}
`, context)
}

func testAccBigQueryDatasetAccess_multiple(datasetID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset_access" "access" {
  dataset_id = google_bigquery_dataset.dataset.dataset_id
  role       = "WRITER"
  domain     = "google.com"
}

resource "google_bigquery_dataset_access" "access2" {
  dataset_id    = google_bigquery_dataset.dataset.dataset_id
  role          = "READER"
  special_group = "projectWriters"
}

resource "google_bigquery_dataset" "dataset" {
  dataset_id = "%s"
}
`, datasetID)
}

func testAccBigQueryDatasetAccess_predefinedRole(role, datasetID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset_access" "access" {
  dataset_id = google_bigquery_dataset.dataset.dataset_id
  role       = "%s"
  domain     = "google.com"
}

resource "google_bigquery_dataset" "dataset" {
  dataset_id = "%s"
}
`, role, datasetID)
}

func testAccBigQueryDatasetAccess_iamMember(datasetID, sinkName string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset_access" "dns_query_sink" {
  dataset_id = google_bigquery_dataset.dataset.dataset_id
  role = "roles/bigquery.dataEditor"
  iam_member = google_logging_project_sink.logging_sink.writer_identity
}

resource "google_bigquery_dataset" "dataset" {
  dataset_id    = "%s"
}

resource "google_logging_project_sink" "logging_sink" {
  name = "%s_logging_project_sink"

  destination = "bigquery.googleapis.com/${google_bigquery_dataset.dataset.id}"

  filter = "resource.type=\"dns_query\""

  unique_writer_identity = true
}
`, datasetID, sinkName)
}

func testAccBigQueryDatasetAccess_allUsers(datasetID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset_access" "dns_query_sink" {
  dataset_id = google_bigquery_dataset.dataset.dataset_id
  role = "roles/bigquery.dataEditor"
  iam_member = "allUsers"
}

resource "google_bigquery_dataset" "dataset" {
  dataset_id    = "%s"
}
`, datasetID)
}

func testAccBigQueryDatasetAccess_allAuthenticatedUsers(datasetID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset_access" "dns_query_sink" {
  dataset_id = google_bigquery_dataset.dataset.dataset_id
  role = "roles/bigquery.dataEditor"
  iam_member = "allAuthenticatedUsers"
}

resource "google_bigquery_dataset" "dataset" {
  dataset_id    = "%s"
}
`, datasetID)
}
