---
layout: "google"
page_title: "Google: google_storage_notification"
sidebar_current: "docs-google-storage-notification"
description: |-
  Creates a new notification configuration on a specified bucket.
---

# google\_storage\_notification

Creates a new notification configuration on a specified bucket, establishing a flow of event notifications from GCS to a Cloud Pub/Sub topic.
 For more information see 
[the official documentation](https://cloud.google.com/storage/docs/pubsub-notifications) 
and 
[API](https://cloud.google.com/storage/docs/json_api/v1/notifications).

## Example Usage

```hcl
resource "google_storage_bucket" "bucket" {
	name = "default_bucket"
}
		
resource "google_pubsub_topic" "topic" {
	name = "default_topic"
}

// In order to enable notifications,
// a GCS service account unique to each project
// must have the IAM permission "projects.topics.publish" to a Cloud Pub/Sub topic from this project
// The only reference to this requirement can be found here:
// https://cloud.google.com/storage/docs/gsutil/commands/notification
// The GCS service account has the format of <project-id>@gs-project-accounts.iam.gserviceaccount.com
// API for retrieving it https://cloud.google.com/storage/docs/json_api/v1/projects/serviceAccount/get

resource "google_pubsub_topic_iam_binding" "binding" {
	topic       = "${google_pubsub_topic.topic.name}"
	role        = "roles/pubsub.publisher"
		  
	members     = ["serviceAccount:my-project-id@gs-project-accounts.iam.gserviceaccount.com"]
}

resource "google_storage_notification" "notification" {
	bucket            = "${google_storage_bucket.bucket.name}"
	payload_format    = "JSON_API_V1"
	topic             = "${google_pubsub_topic.topic.id}"
	event_types       = ["%s","%s"]
	custom_attributes {
		new-attribute = "new-attribute-value"
	}
	depends_on        = ["google_pubsub_topic_iam_binding.binding"]
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required) The name of the bucket.

* `payload_format` - (Required) The desired content of the Payload. One of `"JSON_API_V1"` or `"NONE"`.

* `topic` - (Required) The Cloud PubSub topic to which this subscription publishes.

- - -

* `custom_attributes` - (Optional)  A set of key/value attribute pairs to attach to each Cloud PubSub message published for this notification subscription

* `event_types` - (Optional) List of event type filters for this notification config. If not specified, Cloud Storage will send notifications for all event types. The valid types are: `"OBJECT_FINALIZE"`, `"OBJECT_METADATA_UPDATE"`, `"OBJECT_DELETE"`, `"OBJECT_ARCHIVE"`

* `object_name_prefix` - (Optional) Specifies a prefix path filter for this notification config. Cloud Storage will only send notifications for objects in this bucket whose names begin with the specified prefix.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `self_link` - The URI of the created resource.

## Import

Storage notifications can be imported using the notification `id` in the format `<bucket_name>/notificationConfigs/<id>` e.g.

```
$ terraform import google_storage_notification.notification default_bucket/notificationConfigs/102
```
