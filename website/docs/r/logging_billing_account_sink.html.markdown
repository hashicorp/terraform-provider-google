---
layout: "google"
page_title: "Google: google_logging_billing-account_sink"
sidebar_current: "docs-google-logging-billing-account-sink"
description: |-
  Manages a billing account logging sink.
---

# google\_logging\_billing\_account\_sink

Manages a billing account logging sink. For more information see
[the official documentation](https://cloud.google.com/logging/docs/) and
[Exporting Logs in the API](https://cloud.google.com/logging/docs/api/tasks/exporting-logs).

Note that you must have the "Logs Configuration Writer" IAM role (`roles/logging.configWriter`)
granted to the credentials used with terraform.

## Example Usage

```hcl
resource "google_logging_billing_account_sink" "my-sink" {
    name = "my-sink"
    billing_account = "ABCDEF-012345-GHIJKL"

    # Can export to pubsub, cloud storage, or bigtable
    destination = "storage.googleapis.com/${google_storage_bucket.log-bucket.name}"
}

resource "google_storage_bucket" "log-bucket" {
    name     = "billing-logging-bucket"
}

resource "google_project_iam_binding" "log-writer" {
    role = "roles/storage.objectCreator"

    members = [
        "${google_logging_billing_account_sink.my-sink.writer_identity}",
    ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the logging sink.

* `billing_account` - (Required) The billing account exported to the sink.

* `destination` - (Required) The destination of the sink (or, in other words, where logs are written to). Can be a
    Cloud Storage bucket, a PubSub topic, or a BigQuery dataset. Examples:
```
"storage.googleapis.com/[GCS_BUCKET]"
"bigquery.googleapis.com/projects/[PROJECT_ID]/datasets/[DATASET]"
"pubsub.googleapis.com/projects/[PROJECT_ID]/topics/[TOPIC_ID]"
```
    The writer associated with the sink must have access to write to the above resource.

* `filter` - (Optional) The filter to apply when exporting logs. Only log entries that match the filter are exported.
    See (Advanced Log Filters)[https://cloud.google.com/logging/docs/view/advanced_filters] for information on how to
    write a filter.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `writer_identity` - The identity associated with this sink. This identity must be granted write access to the
    configured `destination`.
