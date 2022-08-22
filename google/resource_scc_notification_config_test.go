package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSecurityCenterNotificationConfig_updateStreamingConfigFilter(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        getTestOrgFromEnv(t),
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecurityCenterNotificationConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityCenterNotificationConfig_sccNotificationConfigBasicExample(context),
			},
			{
				ResourceName:            "google_scc_notification_config.custom_notification_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"organization", "config_id"},
			},
			{
				Config: testAccSecurityCenterNotificationConfig_updateStreamingConfigFilter(context),
			},
			{
				ResourceName:            "google_scc_notification_config.custom_notification_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"organization", "config_id"},
			},
		},
	})
}

func testAccSecurityCenterNotificationConfig_updateStreamingConfigFilter(context map[string]interface{}) string {
	return Nprintf(`
resource "google_pubsub_topic" "scc_notification" {
  name = "tf-test-my-topic%{random_suffix}"
}

resource "google_scc_notification_config" "custom_notification_config" {
  config_id    = "tf-test-my-config%{random_suffix}"
  organization = "%{org_id}"
  description  = "My custom Cloud Security Command Center Finding Notification Configuration"
  pubsub_topic =  google_pubsub_topic.scc_notification.id

  streaming_config {
    filter = "category = \"OPEN_FIREWALL\""
  }
}
`, context)
}
