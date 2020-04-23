---
subcategory: "Cloud (Stackdriver) Logging"
layout: "google"
page_title: "Google: google_logging_organization_sink"
sidebar_current: "docs-google-logging-organization-sink"
description: |-
  Manages a organization-level logging sink.
---

# google\_logging\_organization\_sink

Manages a organization-level logging sink. For more information see
[the official documentation](https://cloud.google.com/logging/docs/) and
[Exporting Logs in the API](https://cloud.google.com/logging/docs/api/tasks/exporting-logs).

Note that you must have the "Logs Configuration Writer" IAM role (`roles/logging.configWriter`)
granted to the credentials used with terraform.

## Example Usage

```hcl
resource "google_logging_organization_sink" "my-sink" {
  name   = "my-sink"
  org_id = "123456789"

  # Can export to pubsub, cloud storage, or bigquery
  destination = "storage.googleapis.com/${google_storage_bucket.log-bucket.name}"

  # Log all WARN or higher severity messages relating to instances
  filter = "resource.type = gce_instance AND severity >= WARN"
}

resource "google_storage_bucket" "log-bucket" {
  name = "organization-logging-bucket"
}

resource "google_project_iam_member" "log-writer" {
  role = "roles/storage.objectCreator"

  member = google_logging_organization_sink.my-sink.writer_identity
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the logging sink.

* `org_id` - (Required) The numeric ID of the organization to be exported to the sink.

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

* `include_children` - (Optional) Whether or not to include children organizations in the sink export. If true, logs
    associated with child projects are also exported; otherwise only logs relating to the provided organization are included.

* `bigquery_options` - (Optional) Options that affect sinks exporting data to BigQuery. Structure documented below.

The `bigquery_options` block supports:

* `use_partitioned_tables` - (Required) Whether to use [BigQuery's partition tables](https://cloud.google.com/bigquery/docs/partitioned-tables).
    By default, Logging creates dated tables based on the log entries' timestamps, e.g. syslog_20170523. With partitioned
    tables the date suffix is no longer present and [special query syntax](https://cloud.google.com/bigquery/docs/querying-partitioned-tables)
    has to be used instead. In both cases, tables are sharded based on UTC timezone.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `writer_identity` - The identity associated with this sink. This identity must be granted write access to the
    configured `destination`.

## Import

Organization-level logging sinks can be imported using this format:

```
$ terraform import google_logging_organization_sink.my_sink organizations/{{organization_id}}/sinks/{{sink_id}}
```
