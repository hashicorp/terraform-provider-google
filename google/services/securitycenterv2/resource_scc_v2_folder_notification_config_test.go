// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package securitycenterv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccSecurityCenterV2FolderNotificationConfig_basic(t *testing.T) {
	t.Parallel()

	orgID := envvar.GetTestOrgFromEnv(t)
	randomSuffix := acctest.RandString(t, 10)

	context := map[string]interface{}{
		"org_id":        orgID,
		"random_suffix": randomSuffix,
	}

	configID := fmt.Sprintf("tf-test-config-%s", randomSuffix)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
			"time":   {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityCenterV2FolderNotificationConfig_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_scc_v2_folder_notification_config.default", "config_id", configID),
				),
			},
			{
				Config:            testAccSecurityCenterV2FolderNotificationConfig_basic(context),
				ResourceName:      "google_scc_v2_folder_notification_config.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSecurityCenterV2FolderNotificationConfig_update(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_scc_v2_folder_notification_config.default", "config_id", configID),
				),
			},
			{
				Config:            testAccSecurityCenterV2FolderNotificationConfig_update(context),
				ResourceName:      "google_scc_v2_folder_notification_config.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSecurityCenterV2FolderNotificationConfig_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_folder" "folder" {
  parent       = "organizations/%{org_id}"
  display_name = "tf-test-folder-name%{random_suffix}"
  deletion_protection = false
}

resource "time_sleep" "wait_1_minute" {
	depends_on = [google_folder.folder]

	create_duration = "3m"
}

resource "google_pubsub_topic" "scc_v2_folder_notification_config" {
  name = "tf-test-topic-%{random_suffix}"
}

resource "google_scc_v2_folder_notification_config" "default" {
  config_id    = "tf-test-config-%{random_suffix}"
  folder 	   = google_folder.folder.folder_id
  location     = "global"
  description  = "A test folder notification config"
  pubsub_topic = google_pubsub_topic.scc_v2_folder_notification_config.id

  streaming_config {
    filter = "severity = \"HIGH\""
  }

  depends_on = [time_sleep.wait_1_minute]
}
`, context)
}

func testAccSecurityCenterV2FolderNotificationConfig_update(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_folder" "folder" {
  parent       = "organizations/%{org_id}"
  display_name = "tf-test-folder-name%{random_suffix}"
  deletion_protection = false
}

resource "google_pubsub_topic" "scc_v2_folder_notification_config" {
  name = "tf-test-topic-%{random_suffix}"
}

resource "google_scc_v2_folder_notification_config" "default" {
  config_id    = "tf-test-config-%{random_suffix}"
  folder 	   = google_folder.folder.folder_id
  location     = "global"
  description  = "An updated test folder notification config"
  pubsub_topic = google_pubsub_topic.scc_v2_folder_notification_config.id

  streaming_config {
    filter = "severity = \"CRITICAL\""
  }
}
`, context)
}
