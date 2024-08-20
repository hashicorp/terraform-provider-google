// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package securitycenter_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccSecurityCenterFolderNotificationConfig_basic(t *testing.T) {
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
				Config: testAccSecurityCenterFolderNotificationConfig_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_scc_folder_notification_config.default", "config_id", configID),
				),
			},
			{
				Config:            testAccSecurityCenterFolderNotificationConfig_basic(context),
				ResourceName:      "google_scc_folder_notification_config.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSecurityCenterFolderNotificationConfig_update(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_scc_folder_notification_config.default", "config_id", configID),
				),
			},
			{
				Config:            testAccSecurityCenterFolderNotificationConfig_update(context),
				ResourceName:      "google_scc_folder_notification_config.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSecurityCenterFolderNotificationConfig_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_folder" "folder" {
  parent       = "organizations/%{org_id}"
  display_name = "tf-test-folder-name%{random_suffix}"
}

resource "time_sleep" "wait_1_minute" {
	depends_on = [google_folder.folder]

	create_duration = "3m"
}

resource "google_pubsub_topic" "scc_folder_notification_config" {
  name = "tf-test-topic-%{random_suffix}"
}

resource "google_scc_folder_notification_config" "default" {
  config_id    = "tf-test-config-%{random_suffix}"
  folder 	   = google_folder.folder.folder_id
  description  = "A test folder notification config"
  pubsub_topic = google_pubsub_topic.scc_folder_notification_config.id

  streaming_config {
    filter = "severity = \"HIGH\""
  }

  depends_on = [time_sleep.wait_1_minute]
}
`, context)
}

func testAccSecurityCenterFolderNotificationConfig_update(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_folder" "folder" {
  parent       = "organizations/%{org_id}"
  display_name = "tf-test-folder-name%{random_suffix}"
}

resource "google_pubsub_topic" "scc_folder_notification_config" {
  name = "tf-test-topic-%{random_suffix}"
}

resource "google_scc_folder_notification_config" "default" {
  config_id    = "tf-test-config-%{random_suffix}"
  folder 	   = google_folder.folder.folder_id
  description  = "An updated test folder notification config"
  pubsub_topic = google_pubsub_topic.scc_folder_notification_config.id

  streaming_config {
    filter = "severity = \"CRITICAL\""
  }
}
`, context)
}
