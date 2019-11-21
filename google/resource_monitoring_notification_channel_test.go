package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestMonitoringNotificationChannel_labelsObfuscated(t *testing.T) {
	testCases := map[string]struct {
		serverV  string
		expected bool
	}{
		"":         {"", false},
		"foo":      {"foo", false},
		"value":    {"diffValue", false},
		"charcnt8": {"****diff", false},
		"foobar":   {"***bar", true},
		"SECRET":   {"**CRET", true},
	}

	for stateV, testCase := range testCases {
		result := isMonitoringNotificationChannelLabelsObfuscated(testCase.serverV, stateV)
		if result != testCase.expected {
			t.Errorf("expected state value %q and server value %q to return obfuscated=%t, got %t", stateV, testCase.serverV, testCase.expected, result)
		}
	}
}

func TestAccMonitoringNotificationChannel_update(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitoringNotificationChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringNotificationChannel_update("email", `email_address = "fake_email@blahblah.com"`, "true"),
			},
			{
				ResourceName:      "google_monitoring_notification_channel.update",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMonitoringNotificationChannel_update("sms", `number = "+15555379009"`, "false"),
			},
			{
				ResourceName:      "google_monitoring_notification_channel.update",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMonitoringNotificationChannel_update(channel, labels, enabled string) string {
	return fmt.Sprintf(`
resource "google_monitoring_notification_channel" "update" {
  display_name = "IntTest Notification Channel"
  type         = "%s"
  labels = {
    %s
  }

  enabled = "%s"
}
`, channel, labels, enabled,
	)
}
