package google

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceGoogleMonitoringNotificationChannel_byDisplayName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleMonitoringNotificationChannel_byDisplayName(acctest.RandomWithPrefix("tf-test")),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceState(
						"data.google_monitoring_notification_channel.my",
						"google_monitoring_notification_channel.my"),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleMonitoringNotificationChannel_byType(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleMonitoringNotificationChannel_byType(acctest.RandomWithPrefix("tf-test")),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceState(
						"data.google_monitoring_notification_channel.my",
						"google_monitoring_notification_channel.my"),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleMonitoringNotificationChannel_byDisplayNameAndType(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleMonitoringNotificationChannel_byDisplayNameAndType(acctest.RandomWithPrefix("tf-test")),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceState(
						"data.google_monitoring_notification_channel.my",
						"google_monitoring_notification_channel.myemail"),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleMonitoringNotificationChannel_NotFound(t *testing.T) {
	displayName := acctest.RandomWithPrefix("tf-test")
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceGoogleMonitoringNotificationChannel_NotFound(displayName),
				ExpectError: regexp.MustCompile(fmt.Sprintf("No NotificationChannel found using filter=display_name=\"%s\"", displayName)),
			},
		},
	})
}

func TestAccDataSourceGoogleMonitoringNotificationChannel_NotUnique(t *testing.T) {
	displayName := acctest.RandomWithPrefix("tf-test")
	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleMonitoringNotificationChannel_NotUnique(displayName),
			},
			{
				Config:      testAccDataSourceGoogleMonitoringNotificationChannel_NotUniqueDS(displayName),
				ExpectError: regexp.MustCompile(fmt.Sprintf("More than one matching NotificationChannel found using filter=display_name=\"%s\"", displayName)),
			},
		},
	})
}

func testAccDataSourceGoogleMonitoringNotificationChannel_byDisplayName(displayName string) string {
	return fmt.Sprintf(`
resource "google_monitoring_notification_channel" "my" {
  display_name = "%s"
  type         = "webhook_tokenauth"

  labels = {
    url = "http://www.acme.org"
  }
}

data "google_monitoring_notification_channel" "my" {
  display_name = google_monitoring_notification_channel.my.display_name
}
`, displayName)
}

func testAccDataSourceGoogleMonitoringNotificationChannel_byType(displayName string) string {
	return fmt.Sprintf(`
resource "google_monitoring_notification_channel" "my" {
  display_name = "%s"
  type         = "sms"

  labels = {
    number = "+1555"
  }
}

data "google_monitoring_notification_channel" "my" {
  type = google_monitoring_notification_channel.my.type
}
`, displayName)
}

func testAccDataSourceGoogleMonitoringNotificationChannel_byDisplayNameAndType(displayName string) string {
	return fmt.Sprintf(`
resource "google_monitoring_notification_channel" "mywebhook" {
  display_name = "%s"
  type         = "webhook_tokenauth"

  labels = {
    url = "http://www.acme.org"
  }
}

resource "google_monitoring_notification_channel" "myemail" {
  display_name = google_monitoring_notification_channel.mywebhook.display_name
  type         = "email"

  labels = {
    email_address = "mailme@acme.org"
  }
}

data "google_monitoring_notification_channel" "my" {
  display_name = google_monitoring_notification_channel.myemail.display_name
  type = google_monitoring_notification_channel.myemail.type
}
`, displayName)
}

func testAccDataSourceGoogleMonitoringNotificationChannel_NotFound(displayName string) string {
	return fmt.Sprintf(`
data "google_monitoring_notification_channel" "my" {
  display_name = "%s"
}
`, displayName)
}

func testAccDataSourceGoogleMonitoringNotificationChannel_NotUnique(displayName string) string {
	return fmt.Sprintf(`
resource "google_monitoring_notification_channel" "default" {
  display_name = "%s"
  type         = "webhook_tokenauth"

  labels = {
    url = "http://www.acme1.org"
  }
}

resource "google_monitoring_notification_channel" "default2" {
  display_name = google_monitoring_notification_channel.default.display_name
  type         = "webhook_tokenauth"

  labels = {
    url = "http://www.acme2.org"
  }
}
`, displayName)
}

func testAccDataSourceGoogleMonitoringNotificationChannel_NotUniqueDS(displayName string) string {
	return fmt.Sprintf(`
data "google_monitoring_notification_channel" "my" {
  display_name = "%s"
}
`, displayName)
}
