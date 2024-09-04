// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package securitycenterv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccSecurityCenterV2ProjectNotificationConfig_updateStreamingConfigFilter(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"location":      "global",
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityCenterV2ProjectNotificationConfig_sccV2ProjectNotificationConfigBasicExample(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_scc_v2_project_notification_config.custom_notification_config", "id"),
				),
			},
			{
				ResourceName:            "google_scc_v2_project_notification_config.custom_notification_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project", "location", "config_id"},
			},
			{
				Config: testAccSecurityCenterV2ProjectNotificationConfig_updateStreamingConfigFilter(context),
			},
			{
				ResourceName:            "google_scc_v2_project_notification_config.custom_notification_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project", "location", "config_id"},
			},
			{
				Config: testAccSecurityCenterV2ProjectNotificationConfig_emptyStreamingConfigFilter(context),
			},
			{
				ResourceName:            "google_scc_v2_project_notification_config.custom_notification_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project", "location", "config_id"},
			},
		},
	})
}

func testAccSecurityCenterV2ProjectNotificationConfig_updateStreamingConfigFilter(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_pubsub_topic" "scc_v2_project_notification" {
  name = "tf-test-my-topic%{random_suffix}"
}

resource "google_scc_v2_project_notification_config" "custom_notification_config" {
  config_id    = "tf-test-my-config%{random_suffix}"
  project      = "%{project}"
  description  = "My custom Cloud Security Command Center Finding Notification Configuration"
  pubsub_topic =  google_pubsub_topic.scc_v2_project_notification.id
  location     = "global"

  streaming_config {
    filter = "category = \"OPEN_FIREWALL\""
  }
}
`, context)
}

func testAccSecurityCenterV2ProjectNotificationConfig_emptyStreamingConfigFilter(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_pubsub_topic" "scc_v2_project_notification" {
  name = "tf-test-my-topic%{random_suffix}"
}

resource "google_scc_v2_project_notification_config" "custom_notification_config" {
  config_id    = "tf-test-my-config%{random_suffix}"
  project      = "%{project}"
  description  = "My custom Cloud Security Command Center Finding Notification Configuration"
  pubsub_topic =  google_pubsub_topic.scc_v2_project_notification.id
  location     = "global"

  streaming_config {
    filter = ""
  }
}
`, context)
}
