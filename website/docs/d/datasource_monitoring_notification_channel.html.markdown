---
subcategory: "Stackdriver Monitoring"
layout: "google"
page_title: "Google: google_monitoring_notification_channel"
sidebar_current: "docs-google-datasource-monitoring-notification-channel"
description: |-
  A NotificationChannel is a medium through which an alert is delivered
  when a policy violation is detected.
---

# google\_monitoring\_notification\_channel

A NotificationChannel is a medium through which an alert is delivered
when a policy violation is detected. Examples of channels include email, SMS,
and third-party messaging applications. Fields containing sensitive information
like authentication tokens or contact info are only partially populated on retrieval.


To get more information about NotificationChannel, see:

* [API documentation](https://cloud.google.com/monitoring/api/ref_v3/rest/v3/projects.notificationChannels)
* How-to Guides
    * [Notification Options](https://cloud.google.com/monitoring/support/notification-options)
    * [Monitoring API Documentation](https://cloud.google.com/monitoring/api/v3/)


## Example Usage - Notification Channel Basic


```hcl
data "google_monitoring_notification_channel" "basic" {
  display_name = "Test Notification Channel"
}

resource "google_monitoring_alert_policy" "alert_policy" {
  display_name = "My Alert Policy"
  notification_channels = [data.google_monitoring_notification_channel.basic.name]
  combiner     = "OR"
  conditions {
    display_name = "test condition"
    condition_threshold {
      filter     = "metric.type=\"compute.googleapis.com/instance/disk/write_bytes_count\" AND resource.type=\"gce_instance\""
      duration   = "60s"
      comparison = "COMPARISON_GT"
      aggregations {
        alignment_period   = "60s"
        per_series_aligner = "ALIGN_RATE"
      }
    }
  }
}

```

## Argument Reference

The arguments of this data source act as filters for querying the available notification channels. The given filters must match exactly one notification channel whose data will be exported as attributes. The following arguments are supported:


* `display_name` -
  (Optional)
    The display name for this notification channel.

* `type` -
  (Optional)
  The type of the notification channel.

~> **NOTE:** One of `display_name` or `type` must be specified.

- - -


* `project` - (Optional) The ID of the project in which the resource belongs.
    If it is not provided, the provider project is used.


## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:


* `name` -
  The full REST resource name for this channel. The syntax is:
  `projects/[PROJECT_ID]/notificationChannels/[CHANNEL_ID]`.

* `verification_status` -
  Indicates whether this channel has been verified or not.

* `labels` -
  Configuration fields that define the channel and its behavior.

* `user_labels` -
  User-supplied key/value data that does not need to conform to the corresponding NotificationChannelDescriptor's schema, unlike the labels field.

* `description` -
  An optional human-readable description of this notification channel.

* `enabled` -
  Whether notifications are forwarded to the described channel.
