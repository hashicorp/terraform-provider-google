// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package securitycenter_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccSecurityCenterProjectBigQueryExportConfig_basic(t *testing.T) {
	t.Parallel()

	randomSuffix := acctest.RandString(t, 10)
	datasetID := "tf_test_" + randomSuffix
	orgID := envvar.GetTestOrgFromEnv(t)

	context := map[string]interface{}{
		"org_id":              orgID,
		"random_suffix":       randomSuffix,
		"dataset_id":          datasetID,
		"big_query_export_id": "tf-test-export-" + randomSuffix,
		"project":             envvar.GetTestProjectFromEnv(),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
			"time":   {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityCenterProjectBigQueryExportConfig_basic(context),
			},
			{
				ResourceName:            "google_scc_project_scc_big_query_export.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"update_time", "project"},
			},
			{
				Config: testAccSecurityCenterProjectBigQueryExportConfig_update(context),
			},
			{
				ResourceName:            "google_scc_project_scc_big_query_export.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"update_time", "project"},
			},
		},
	})
}

func testAccSecurityCenterProjectBigQueryExportConfig_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_bigquery_dataset" "default" {
  dataset_id                  = "%{dataset_id}"
  friendly_name               = "test"
  description                 = "This is a test description"
  location                    = "US"
  default_table_expiration_ms = 3600000
  default_partition_expiration_ms = null
  delete_contents_on_destroy  = true

  labels = {
    env = "default"
  }

  lifecycle {
	ignore_changes = [default_partition_expiration_ms]
  }
}

resource "time_sleep" "wait_x_minutes" {
	depends_on = [google_bigquery_dataset.default]
	create_duration = "6m"
	# need to wait for destruction due to 
	# 'still in use' error from api 
	destroy_duration = "1m"
}

resource "google_scc_project_scc_big_query_export" "default" {
  big_query_export_id    = "%{big_query_export_id}"
  project      = "%{project}"
  dataset      = google_bigquery_dataset.default.id
  description  = "Cloud Security Command Center Findings Big Query Export Config"
  filter       = "state=\"ACTIVE\" AND NOT mute=\"MUTED\""

  depends_on = [time_sleep.wait_x_minutes]
}

resource "time_sleep" "wait_for_cleanup" {
	create_duration = "6m"
	depends_on = [google_scc_project_scc_big_query_export.default]
}

`, context)
}

func testAccSecurityCenterProjectBigQueryExportConfig_update(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_bigquery_dataset" "default" {
  dataset_id                  = "%{dataset_id}"
  friendly_name               = "test"
  description                 = "This is a test description"
  location                    = "US"
  default_table_expiration_ms = 3600000
  default_partition_expiration_ms = null
  delete_contents_on_destroy  = true

  labels = {
    env = "default"
  }

  lifecycle {
	ignore_changes = [default_partition_expiration_ms]
  }
}

resource "time_sleep" "wait_x_minutes" {
	depends_on = [google_bigquery_dataset.default]
	create_duration = "6m"
	# need to wait for destruction due to
	# 'still in use' error from api
	destroy_duration = "1m"
}

resource "google_scc_project_scc_big_query_export" "default" {
  big_query_export_id    = "%{big_query_export_id}"
  project      = "%{project}"
  dataset      = google_bigquery_dataset.default.id
  description  = "SCC Findings Big Query Export Update"
  filter       = "state=\"ACTIVE\" AND NOT mute=\"MUTED\""

  depends_on = [time_sleep.wait_x_minutes]
}

resource "time_sleep" "wait_for_cleanup" {
	create_duration = "6m"
	depends_on = [google_scc_project_scc_big_query_export.default]
}

`, context)
}
