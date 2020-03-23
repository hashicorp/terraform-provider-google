---
subcategory: "Cloud Storage"
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

In order to enable notifications, a special Google Cloud Storage service account unique to the project
must have the IAM permission "projects.topics.publish" for a Cloud Pub/Sub topic in the project. To get the service
account's email address, use the `google_storage_project_service_account` datasource's `email_address` value, and see below
for an example of enabling notifications by granting the correct IAM permission. See
[the notifications documentation](https://cloud.google.com/storage/docs/gsutil/commands/notification) for more details.

>**NOTE**: This resource can affect your storage IAM policy. If you are using this in the same config as your storage IAM policy resources, consider
making this resource dependent on those IAM resources via `depends_on`. This will safeguard against errors due to IAM race conditions.

## Example Usage

```hcl
resource "google_storage_notification" "notification" {
  bucket         = google_storage_bucket.bucket.name
  payload_format = "JSON_API_V1"
  topic          = google_pubsub_topic.topic.id
  event_types    = ["OBJECT_FINALIZE", "OBJECT_METADATA_UPDATE"]
  custom_attributes = {
    new-attribute = "new-attribute-value"
  }
  depends_on = [google_pubsub_topic_iam_binding.binding]
}

// Enable notifications by giving the correct IAM permission to the unique service account.

data "google_storage_project_service_account" "gcs_account" {
}

resource "google_pubsub_topic_iam_binding" "binding" {
  topic   = google_pubsub_topic.topic.id
  role    = "roles/pubsub.publisher"
  members = ["serviceAccount:${data.google_storage_project_service_account.gcs_account.email_address}"]
}

// End enabling notifications

resource "google_storage_bucket" "bucket" {
  name = "default_bucket"
}

resource "google_pubsub_topic" "topic" {
  name = "default_topic"
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required) The name of the bucket.

* `payload_format` - (Required) The desired content of the Payload. One of `"JSON_API_V1"` or `"NONE"`.

* `topic` - (Required) The Cloud PubSub topic to which this subscription publishes. Expects either the 
    topic name, assumed to belong to the default GCP provider project, or the project-level name, 
    i.e. `projects/my-gcp-project/topics/my-topic` or `my-topic`. If the project is not set in the provider,
    you will need to use the project-level name.
    
- - -

* `custom_attributes` - (Optional)  A set of key/value attribute pairs to attach to each Cloud PubSub message published for this notification subscription

* `event_types` - (Optional) List of event type filters for this notification config. If not specified, Cloud Storage will send notifications for all event types. The valid types are: `"OBJECT_FINALIZE"`, `"OBJECT_METADATA_UPDATE"`, `"OBJECT_DELETE"`, `"OBJECT_ARCHIVE"`

* `object_name_prefix` - (Optional) Specifies a prefix path filter for this notification config. Cloud Storage will only send notifications for objects in this bucket whose names begin with the specified prefix.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `notification_id` - The ID of the created notification.

* `self_link` - The URI of the created resource.

## Import

Storage notifications can be imported using the notification `id` in the format `<bucket_name>/notificationConfigs/<id>` e.g.

```
$ terraform import google_storage_notification.notification default_bucket/notificationConfigs/102
```
