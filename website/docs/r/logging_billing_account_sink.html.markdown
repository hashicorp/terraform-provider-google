---
subcategory: "Cloud (Stackdriver) Logging"
description: |-
  Manages a billing account logging sink.
---

# google_logging_billing_account_sink

* [API documentation](https://cloud.google.com/logging/docs/reference/v2/rest/v2/billingAccounts.sinks)
* How-to Guides
    * [Exporting Logs](https://cloud.google.com/logging/docs/export)

~> **Note** You must have the "Logs Configuration Writer" IAM role (`roles/logging.configWriter`)
[granted on the billing account](https://cloud.google.com/billing/reference/rest/v1/billingAccounts/getIamPolicy) to
the credentials used with Terraform. [IAM roles granted on a billing account](https://cloud.google.com/billing/docs/how-to/billing-access) are separate from the
typical IAM roles granted on a project.

## Example Usage

```hcl
resource "google_logging_billing_account_sink" "my-sink" {
  name            = "my-sink"
  description = "some explanation on what this is"
  billing_account = "ABCDEF-012345-GHIJKL"

  # Can export to pubsub, cloud storage, or bigquery
  destination = "storage.googleapis.com/${google_storage_bucket.log-bucket.name}"
}

resource "google_storage_bucket" "log-bucket" {
  name     = "billing-logging-bucket"
  location = "US"
}

resource "google_project_iam_binding" "log-writer" {
  project = "your-project-id"
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
    Cloud Storage bucket, a PubSub topic, a BigQuery dataset or a Cloud Logging bucket. Examples:

    - `storage.googleapis.com/[GCS_BUCKET]`
    - `bigquery.googleapis.com/projects/[PROJECT_ID]/datasets/[DATASET]`
    - `pubsub.googleapis.com/projects/[PROJECT_ID]/topics/[TOPIC_ID]`
    - `logging.googleapis.com/projects/[PROJECT_ID]]/locations/global/buckets/[BUCKET_ID]`

    The writer associated with the sink must have access to write to the above resource.

* `filter` - (Optional) The filter to apply when exporting logs. Only log entries that match the filter are exported.
    See [Advanced Log Filters](https://cloud.google.com/logging/docs/view/advanced_filters) for information on how to
    write a filter.

* `description` - (Optional) A description of this sink. The maximum length of the description is 8000 characters.

* `disabled` - (Optional) If set to True, then this sink is disabled and it does not export any log entries.

* `bigquery_options` - (Optional) Options that affect sinks exporting data to BigQuery. Structure [documented below](#nested_bigquery_options).

* `exclusions` - (Optional) Log entries that match any of the exclusion filters will not be exported. If a log entry is matched by both `filter` and one of `exclusions.filter`, it will not be exported.  Can be repeated multiple times for multiple exclusions. Structure is [documented below](#nested_exclusions).

<a name="nested_bigquery_options"></a>The `bigquery_options` block supports:

* `use_partitioned_tables` - (Required) Whether to use [BigQuery's partition tables](https://cloud.google.com/bigquery/docs/partitioned-tables).
    By default, Logging creates dated tables based on the log entries' timestamps, e.g. syslog_20170523. With partitioned
    tables, the date suffix is no longer present and [special query syntax](https://cloud.google.com/bigquery/docs/querying-partitioned-tables)
    has to be used instead. In both cases, tables are sharded based on UTC timezone.

<a name="nested_exclusions"></a>The `exclusions` block supports:

* `name` - (Required) A client-assigned identifier, such as `load-balancer-exclusion`. Identifiers are limited to 100 characters and can include only letters, digits, underscores, hyphens, and periods. First character has to be alphanumeric.
* `description` - (Optional) A description of this exclusion.
* `filter` - (Required) An advanced logs filter that matches the log entries to be excluded. By using the sample function, you can exclude less than 100% of the matching log entries. See [Advanced Log Filters](https://cloud.google.com/logging/docs/view/advanced_filters) for information on how to
    write a filter.
* `disabled` - (Optional) If set to True, then this exclusion is disabled and it does not exclude any log entries.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `id` - an identifier for the resource with format `billingAccounts/{{billing_account_id}}/sinks/{{sink_id}}`

* `writer_identity` - The identity associated with this sink. This identity must be granted write access to the
    configured `destination`.

## Import

Billing account logging sinks can be imported using this format:

* `billingAccounts/{{billing_account_id}}/sinks/{{sink_id}}`

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import billing account logging sinks using one of the formats above. For example:

```tf
import {
  id = "billingAccounts/{{billing_account_id}}/sinks/{{sink_id}}"
  to = google_logging_billing_account_sink.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), billing account logging sinks can be imported using one of the formats above. For example:

```
$ terraform import google_logging_billing_account_sink.default billingAccounts/{{billing_account_id}}/sinks/{{sink_id}}
```
