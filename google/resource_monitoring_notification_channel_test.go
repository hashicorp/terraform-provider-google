package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccMonitoringNotificationChannel_update(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitoringNotificationChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringNotificationChannel_update("email", `email_address = "fake_email@blahblah.com"`),
			},
			{
				ResourceName:      "google_monitoring_notification_channel.update",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMonitoringNotificationChannel_update("sms", `number = "+15555379009"`),
			},
			{
				ResourceName:      "google_monitoring_notification_channel.update",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMonitoringNotificationChannel_update(channel, labels string) string {
	return fmt.Sprintf(`
resource "google_monitoring_notification_channel" "update" {
  display_name = "IntTest Notification Channel"
  type = "%s"
  labels = {
    %s
  }
}
`, channel, labels,
	)
}
