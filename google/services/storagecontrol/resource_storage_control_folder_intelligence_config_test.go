// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package storagecontrol_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccStorageControlFolderIntelligenceConfig_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccStorageControlFolderIntelligenceConfig_basic(context),
			},
			{
				ResourceName:            "google_storage_control_folder_intelligence_config.folder_intelligence_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name"},
			},
			{
				Config: testAccStorageControlFolderIntelligenceConfig_update_with_filter(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_storage_control_folder_intelligence_config.folder_intelligence_config", "filter.0.excluded_cloud_storage_buckets.0.bucket_id_regexes.0", "random-test-*"),
					resource.TestCheckResourceAttr(
						"google_storage_control_folder_intelligence_config.folder_intelligence_config", "filter.0.excluded_cloud_storage_buckets.0.bucket_id_regexes.1", "random-test2-*"),
					resource.TestCheckResourceAttr(
						"google_storage_control_folder_intelligence_config.folder_intelligence_config", "filter.0.included_cloud_storage_locations.0.locations.0", "us-east-1"),
					resource.TestCheckResourceAttr(
						"google_storage_control_folder_intelligence_config.folder_intelligence_config", "filter.0.included_cloud_storage_locations.0.locations.1", "us-east-2"),
				),
			},
			{
				ResourceName:            "google_storage_control_folder_intelligence_config.folder_intelligence_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name"},
			},
			{
				Config: testAccStorageControlFolderIntelligenceConfig_update_with_empty_filter_fields(context),
			},
			{
				ResourceName:            "google_storage_control_folder_intelligence_config.folder_intelligence_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name"},
			},
			{
				Config: testAccStorageControlFolderIntelligenceConfig_update_with_filter2(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_storage_control_folder_intelligence_config.folder_intelligence_config", "filter.0.included_cloud_storage_buckets.0.bucket_id_regexes.0", "random-test-*"),
					resource.TestCheckResourceAttr(
						"google_storage_control_folder_intelligence_config.folder_intelligence_config", "filter.0.included_cloud_storage_buckets.0.bucket_id_regexes.1", "random-test2-*"),
					resource.TestCheckResourceAttr(
						"google_storage_control_folder_intelligence_config.folder_intelligence_config", "filter.0.excluded_cloud_storage_locations.0.locations.0", "us-east-1"),
					resource.TestCheckResourceAttr(
						"google_storage_control_folder_intelligence_config.folder_intelligence_config", "filter.0.excluded_cloud_storage_locations.0.locations.1", "us-east-2"),
				),
			},
			{
				ResourceName:            "google_storage_control_folder_intelligence_config.folder_intelligence_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name"},
			},
			{
				Config: testAccStorageControlFolderIntelligenceConfig_update_with_empty_filter_fields2(context),
			},
			{
				ResourceName:            "google_storage_control_folder_intelligence_config.folder_intelligence_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name"},
			},
			{
				Config: testAccStorageControlFolderIntelligenceConfig_update_mode_disable(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_storage_control_folder_intelligence_config.folder_intelligence_config", "edition_config", "DISABLED"),
				),
			},
			{
				ResourceName:            "google_storage_control_folder_intelligence_config.folder_intelligence_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name"},
			},
			{
				Config: testAccStorageControlFolderIntelligenceConfig_update_mode_inherit(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_storage_control_folder_intelligence_config.folder_intelligence_config", "edition_config", "INHERIT"),
				),
			},
			{
				ResourceName:            "google_storage_control_folder_intelligence_config.folder_intelligence_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name"},
			},
		},
	})
}

func testAccStorageControlFolderIntelligenceConfig_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_folder" "folder" {
  parent       = "organizations/%{org_id}"
  display_name = "tf-test-folder-name%{random_suffix}"
  deletion_protection=false
}

resource "time_sleep" "wait_120_seconds" {
  depends_on = [google_folder.folder]
  create_duration = "120s"
}

resource "google_storage_control_folder_intelligence_config" "folder_intelligence_config" {
  name = google_folder.folder.folder_id
  edition_config = "STANDARD"
  depends_on = [time_sleep.wait_120_seconds]
}
`, context)
}

func testAccStorageControlFolderIntelligenceConfig_update_with_filter(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_folder" "folder" {
  parent       = "organizations/%{org_id}"
  display_name = "tf-test-folder-name%{random_suffix}"
  deletion_protection=false
}

resource "time_sleep" "wait_120_seconds" {
  depends_on = [google_folder.folder]
  create_duration = "120s"
}

resource "google_storage_control_folder_intelligence_config" "folder_intelligence_config" {
  name = google_folder.folder.folder_id
  edition_config = "STANDARD"
  filter {
    excluded_cloud_storage_buckets{
      bucket_id_regexes = ["random-test-*", "random-test2-*"]
    }
    included_cloud_storage_locations{
      locations = ["us-east-1", "us-east-2"]
    }
  }
  depends_on = [time_sleep.wait_120_seconds]
}
`, context)
}

func testAccStorageControlFolderIntelligenceConfig_update_with_empty_filter_fields(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_folder" "folder" {
  parent       = "organizations/%{org_id}"
  display_name = "tf-test-folder-name%{random_suffix}"
  deletion_protection=false
}

resource "time_sleep" "wait_120_seconds" {
  depends_on = [google_folder.folder]
  create_duration = "120s"
}

resource "google_storage_control_folder_intelligence_config" "folder_intelligence_config" {
  name = google_folder.folder.folder_id
  edition_config = "STANDARD"
  filter {
    excluded_cloud_storage_buckets{
      bucket_id_regexes = []
    }
    included_cloud_storage_locations{
      locations = []
    }
  }
  depends_on = [time_sleep.wait_120_seconds]
}
`, context)
}

func testAccStorageControlFolderIntelligenceConfig_update_with_filter2(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_folder" "folder" {
  parent       = "organizations/%{org_id}"
  display_name = "tf-test-folder-name%{random_suffix}"
  deletion_protection=false
}

resource "time_sleep" "wait_120_seconds" {
  depends_on = [google_folder.folder]
  create_duration = "120s"
}

resource "google_storage_control_folder_intelligence_config" "folder_intelligence_config" {
  name = google_folder.folder.folder_id
  edition_config = "STANDARD"
  filter {
    included_cloud_storage_buckets{
      bucket_id_regexes = ["random-test-*", "random-test2-*"]
    }
    excluded_cloud_storage_locations{
      locations = ["us-east-1", "us-east-2"]
    }
  }
  depends_on = [time_sleep.wait_120_seconds]
}
`, context)
}

func testAccStorageControlFolderIntelligenceConfig_update_with_empty_filter_fields2(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_folder" "folder" {
  parent       = "organizations/%{org_id}"
  display_name = "tf-test-folder-name%{random_suffix}"
	deletion_protection=false
}

resource "time_sleep" "wait_120_seconds" {
  depends_on = [google_folder.folder]
  create_duration = "120s"
}

resource "google_storage_control_folder_intelligence_config" "folder_intelligence_config" {
  name = google_folder.folder.folder_id
  edition_config = "STANDARD"
  filter {
    included_cloud_storage_buckets{
      bucket_id_regexes = []
    }
    excluded_cloud_storage_locations{
      locations = []
    }
  }
	depends_on = [time_sleep.wait_120_seconds]
}
`, context)
}

func testAccStorageControlFolderIntelligenceConfig_update_mode_disable(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_folder" "folder" {
  parent       = "organizations/%{org_id}"
  display_name = "tf-test-folder-name%{random_suffix}"
  deletion_protection=false
}

resource "time_sleep" "wait_120_seconds" {
  depends_on = [google_folder.folder]
  create_duration = "120s"
}

resource "google_storage_control_folder_intelligence_config" "folder_intelligence_config" {
  name = google_folder.folder.folder_id
  edition_config = "DISABLED"
  depends_on = [time_sleep.wait_120_seconds]
}
`, context)
}

func testAccStorageControlFolderIntelligenceConfig_update_mode_inherit(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_folder" "folder" {
  parent       = "organizations/%{org_id}"
  display_name = "tf-test-folder-name%{random_suffix}"
  deletion_protection=false
}

resource "time_sleep" "wait_120_seconds" {
  depends_on = [google_folder.folder]
  create_duration = "120s"
}

resource "google_storage_control_folder_intelligence_config" "folder_intelligence_config" {
  name = google_folder.folder.folder_id
  edition_config = "INHERIT"
  depends_on = [time_sleep.wait_120_seconds]
}
`, context)
}
