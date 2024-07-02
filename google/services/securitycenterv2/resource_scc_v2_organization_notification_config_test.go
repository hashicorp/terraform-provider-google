// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package securitycenterv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccSecurityCenterV2OrganizationNotificationConfig_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityCenterV2OrganizationNotificationConfig_basic(context),
			},
			{
				ResourceName:      "google_scc_v2_organization_notification_config.default",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"config_id",
				},
			},
			{
				Config: testAccSecurityCenterV2OrganizationNotificationConfig_update(context),
			},
			{
				ResourceName:      "google_scc_v2_organization_notification_config.default",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"config_id",
				},
			},
		},
	})
}

func testAccSecurityCenterV2OrganizationNotificationConfig_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_pubsub_topic" "scc_v2_organization_notification_config" {
  name = "tf-test-topic-%{random_suffix}"
}

resource "google_scc_v2_organization_notification_config" "default" {
  config_id    = "tf-test-config-%{random_suffix}"
  organization = "%{org_id}"
  location     = "global"
  description  = "A test organization notification config"
  pubsub_topic = google_pubsub_topic.scc_v2_organization_notification_config.id

  streaming_config {
    filter = "severity = \"HIGH\""
  }
}
`, context)
}

func testAccSecurityCenterV2OrganizationNotificationConfig_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_pubsub_topic" "scc_v2_organization_notification_config" {
  name = "tf-test-topic-%{random_suffix}"
}

resource "google_scc_v2_organization_notification_config" "default" {
  config_id    = "tf-test-config-%{random_suffix}"
  organization = "%{org_id}"
  location     = "global"
  description  = "An updated test organization notification config"
  pubsub_topic = google_pubsub_topic.scc_v2_organization_notification_config.id

  streaming_config {
    filter = "severity = \"CRITICAL\""
  }
}
`, context)
}
