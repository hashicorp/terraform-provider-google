// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package monitoring_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGoogleMonitoringNotificationChannel_byDisplayName(t *testing.T) {
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleMonitoringNotificationChannel_byDisplayName(fmt.Sprintf("tf-test-%d", acctest.RandInt(t))),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState(
						"data.google_monitoring_notification_channel.default",
						"google_monitoring_notification_channel.default"),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleMonitoringNotificationChannel_byTypeAndLabel(t *testing.T) {
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleMonitoringNotificationChannel_byTypeAndLabel(fmt.Sprintf("tf-test-%d", acctest.RandInt(t))),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState(
						"data.google_monitoring_notification_channel.default",
						"google_monitoring_notification_channel.default"),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleMonitoringNotificationChannel_UserLabel(t *testing.T) {
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleMonitoringNotificationChannel_byTypeAndUserLabel(fmt.Sprintf("tf-test-%d", acctest.RandInt(t))),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState(
						"data.google_monitoring_notification_channel.default",
						"google_monitoring_notification_channel.default"),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleMonitoringNotificationChannel_byDisplayNameAndType(t *testing.T) {
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleMonitoringNotificationChannel_byDisplayNameAndType(fmt.Sprintf("tf-test-%d", acctest.RandInt(t))),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState(
						"data.google_monitoring_notification_channel.email",
						"google_monitoring_notification_channel.email"),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleMonitoringNotificationChannel_ErrorNoDisplayNameOrType(t *testing.T) {
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceGoogleMonitoringNotificationChannel_NoDisplayNameOrType(),
				ExpectError: regexp.MustCompile("At least one of display_name or type must be provided"),
			},
		},
	})
}

func TestAccDataSourceGoogleMonitoringNotificationChannel_ErrorNotFound(t *testing.T) {
	displayName := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceGoogleMonitoringNotificationChannel_NotFound(displayName),
				ExpectError: regexp.MustCompile(fmt.Sprintf(`No NotificationChannel found using filter: display_name="%s"`, displayName)),
			},
		},
	})
}

func TestAccDataSourceGoogleMonitoringNotificationChannel_ErrorNotUnique(t *testing.T) {
	displayName := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleMonitoringNotificationChannel_NotUnique(displayName),
			},
			{
				Config: testAccDataSourceGoogleMonitoringNotificationChannel_NotUniqueWithData(displayName),
				ExpectError: regexp.MustCompile(fmt.Sprintf(
					`Found more than one 1 NotificationChannel matching specified filter: display_name="%s"`, displayName)),
			},
		},
	})
}

func testAccDataSourceGoogleMonitoringNotificationChannel_byDisplayName(displayName string) string {
	return fmt.Sprintf(`
resource "google_monitoring_notification_channel" "default" {
  display_name = "%s"
  type         = "webhook_tokenauth"

  labels = {
    url = "http://www.google.com"
  }
}

data "google_monitoring_notification_channel" "default" {
  display_name = google_monitoring_notification_channel.default.display_name
}
`, displayName)
}

// Include label so we don't fail on dangling resources
func testAccDataSourceGoogleMonitoringNotificationChannel_byTypeAndLabel(displayName string) string {
	return fmt.Sprintf(`
resource "google_monitoring_notification_channel" "default" {
  display_name = "%s"
  type         = "email"

  labels = {
    email_address = "%s@google.com"
  }
}

data "google_monitoring_notification_channel" "default" {
  type = google_monitoring_notification_channel.default.type
  labels =  google_monitoring_notification_channel.default.labels
}
`, displayName, displayName)
}

func testAccDataSourceGoogleMonitoringNotificationChannel_byTypeAndUserLabel(testName string) string {
	return fmt.Sprintf(`
resource "google_monitoring_notification_channel" "default" {
  display_name = "%s"
  type         = "email"

  labels = {
    email_address = "%s@google.com"
  }

  user_labels = {
    testname = "foo-%s"
  }
}

data "google_monitoring_notification_channel" "default" {
  type = google_monitoring_notification_channel.default.type
  user_labels =  google_monitoring_notification_channel.default.user_labels
}
`, testName, testName, testName)
}

func testAccDataSourceGoogleMonitoringNotificationChannel_byDisplayNameAndType(displayName string) string {
	return fmt.Sprintf(`
resource "google_monitoring_notification_channel" "webhook" {
  display_name = "%s"
  type         = "webhook_tokenauth"

  labels = {
    url = "http://www.google.com"
  }
}

resource "google_monitoring_notification_channel" "email" {
  display_name = "%s"
  type         = "email"

  labels = {
    email_address = "%s@google.com"
  }
}

data "google_monitoring_notification_channel" "email" {
  display_name = google_monitoring_notification_channel.email.display_name
  type = google_monitoring_notification_channel.email.type
}
`, displayName, displayName, displayName)
}

func testAccDataSourceGoogleMonitoringNotificationChannel_NotFound(displayName string) string {
	return fmt.Sprintf(`
data "google_monitoring_notification_channel" "default" {
  display_name = "%s"
}
`, displayName)
}

func testAccDataSourceGoogleMonitoringNotificationChannel_NoDisplayNameOrType() string {
	return `
data "google_monitoring_notification_channel" "default" {
	labels = {
		email = "doesntmatter@google.com'"
	}
    user_labels = {
		foo = "bar"
	}
}
`
}

func testAccDataSourceGoogleMonitoringNotificationChannel_NotUnique(displayName string) string {
	return fmt.Sprintf(`
resource "google_monitoring_notification_channel" "channel-1" {
  display_name = "%[1]s"
  type         = "webhook_tokenauth"

  labels = {
    url = "http://%[1]s.google.com"
  }
}

resource "google_monitoring_notification_channel" "channel-2" {
  display_name = google_monitoring_notification_channel.channel-1.display_name
  type         = "webhook_tokenauth"

  labels = {
    url = "http://%[1]s-copy.google.org"
  }
}
`, displayName)
}

func testAccDataSourceGoogleMonitoringNotificationChannel_NotUniqueWithData(displayName string) string {
	return testAccDataSourceGoogleMonitoringNotificationChannel_NotUnique(displayName) + `

data "google_monitoring_notification_channel" "ds" {
  display_name = google_monitoring_notification_channel.channel-2.display_name
}
`
}
