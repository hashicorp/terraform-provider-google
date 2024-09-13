// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package securitycenter_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccSecurityCenterFolderBigQueryExportConfig_update(t *testing.T) {
	t.Parallel()

	randomSuffix := acctest.RandString(t, 10)
	dataset_id := "tf_test_" + randomSuffix
	dataset_id2 := dataset_id + "2"
	orgID := envvar.GetTestOrgFromEnv(t)

	context := map[string]interface{}{
		"org_id":              orgID,
		"random_suffix":       randomSuffix,
		"dataset_id":          dataset_id,
		"dataset_id2":         dataset_id2,
		"big_query_export_id": "tf-test-export-" + randomSuffix,
		"folder_name":         "tf-test-folder-name-" + randomSuffix,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityCenterFolderBigQueryExportConfig_basic(context),
			},
			{
				ResourceName:            "google_scc_folder_scc_big_query_export.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"update_time"},
			},
			{
				Config: testAccSecurityCenterFolderBigQueryExportConfig_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_scc_folder_scc_big_query_export.default", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_scc_folder_scc_big_query_export.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"update_time"},
			},
		},
	})
}

func testAccSecurityCenterFolderBigQueryExportConfig_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_folder" "folder" {
  parent       = "organizations/%{org_id}"
  display_name = "%{folder_name}"
  deletion_protection = false
}
resource "google_bigquery_dataset" "default" {
  dataset_id                  = "%{dataset_id}"
  friendly_name               = "test"
  description                 = "This is a test description"
  location                    = "US"
  default_table_expiration_ms = 3600000
  default_partition_expiration_ms = null
  labels = {
    env = "default"
  }
  lifecycle {
	ignore_changes = [default_partition_expiration_ms]
  }
}
resource "time_sleep" "wait_1_minute" {
	depends_on = [google_bigquery_dataset.default]
	create_duration = "3m"
}
resource "google_scc_folder_scc_big_query_export" "default" {
  big_query_export_id    = "%{big_query_export_id}"
  folder 	   = google_folder.folder.folder_id
  dataset      = google_bigquery_dataset.default.id
  description  = "Cloud Security Command Center Findings Big Query Export Config"
  filter       = "state=\"ACTIVE\" AND NOT mute=\"MUTED\""

  depends_on = [time_sleep.wait_1_minute]
}

`, context)
}

func testAccSecurityCenterFolderBigQueryExportConfig_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_folder" "folder" {
  parent       = "organizations/%{org_id}"
  display_name = "%{folder_name}"
  deletion_protection = false
}
resource "google_bigquery_dataset" "default" {
  dataset_id                  = "%{dataset_id2}"
  friendly_name               = "test"
  description                 = "This is a test description"
  location                    = "US"
  default_table_expiration_ms = 3600000
  default_partition_expiration_ms = null
  labels = {
    env = "default"
  }
  lifecycle {
	ignore_changes = [default_partition_expiration_ms]
  }
}
resource "google_scc_folder_scc_big_query_export" "default" {
  big_query_export_id    = "%{big_query_export_id}"
  folder 	   = google_folder.folder.folder_id
  dataset      = google_bigquery_dataset.default.id
  description  = "SCC Findings Big Query Export Update"
  filter       = ""
}

`, context)
}
