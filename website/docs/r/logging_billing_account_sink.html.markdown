---
subcategory: "Cloud (Stackdriver) Logging"
layout: "google"
page_title: "Google: google_logging_billing_account_sink"
sidebar_current: "docs-google-logging-billing-account-sink"
description: |-
  Manages a billing account logging sink.
---

# google\_logging\_billing\_account\_sink

Manages a billing account logging sink. For more information see
[the official documentation](https://cloud.google.com/logging/docs/) and
[Exporting Logs in the API](https://cloud.google.com/logging/docs/api/tasks/exporting-logs).

~> **Note** You must have the "Logs Configuration Writer" IAM role (`roles/logging.configWriter`)
[granted on the billing account](https://cloud.google.com/billing/reference/rest/v1/billingAccounts/getIamPolicy) to
the credentials used with Terraform. [IAM roles granted on a billing account](https://cloud.google.com/billing/docs/how-to/billing-access) are separate from the
typical IAM roles granted on a project.

## Example Usage

```hcl
resource "google_logging_billing_account_sink" "my-sink" {
  name            = "my-sink"
  billing_account = "ABCDEF-012345-GHIJKL"

  # Can export to pubsub, cloud storage, or bigquery
  destination = "storage.googleapis.com/${google_storage_bucket.log-bucket.name}"
}

resource "google_storage_bucket" "log-bucket" {
  name = "billing-logging-bucket"
}

resource "google_project_iam_binding" "log-writer" {
  role = "roles/storage.objectCreator"

  members = [
    google_logging_billing_account_sink.my-sink.writer_identity,
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
    See [Advanced Log Filters](https://cloud.google.com/logging/docs/view/advanced_filters) for information on how to
    write a filter.

* `bigquery_options` - (Optional) Options that affect sinks exporting data to BigQuery. Structure documented below.

The `bigquery_options` block supports:

* `use_partitioned_tables` - (Required) Whether to use [BigQuery's partition tables](https://cloud.google.com/bigquery/docs/partitioned-tables).
    By default, Logging creates dated tables based on the log entries' timestamps, e.g. syslog_20170523. With partitioned
    tables the date suffix is no longer present and [special query syntax](https://cloud.google.com/bigquery/docs/querying-partitioned-tables)
    has to be used instead. In both cases, tables are sharded based on UTC timezone.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `id` - an identifier for the resource with format `billingAccounts/{{billing_account_id}}/sinks/{{sink_id}}`

* `writer_identity` - The identity associated with this sink. This identity must be granted write access to the
    configured `destination`.

## Import

Billing account logging sinks can be imported using this format:

```
$ terraform import google_logging_billing_account_sink.my_sink billingAccounts/{{billing_account_id}}/sinks/{{sink_id}}
```
