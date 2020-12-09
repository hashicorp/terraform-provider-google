---
subcategory: "Cloud (Stackdriver) Logging"
layout: "google"
page_title: "Google: google_logging_project_sink"
sidebar_current: "docs-google-logging-project-sink"
description: |-
  Manages a project-level logging sink.
---

# google\_logging\_project\_sink

Manages a project-level logging sink. For more information see
[the official documentation](https://cloud.google.com/logging/docs/),
[Exporting Logs in the API](https://cloud.google.com/logging/docs/api/tasks/exporting-logs)
and
[API](https://cloud.google.com/logging/docs/reference/v2/rest/).

~> **Note:** You must have [granted the "Logs Configuration Writer"](https://cloud.google.com/logging/docs/access-control) IAM role (`roles/logging.configWriter`) to the credentials used with terraform.

~> **Note** You must [enable the Cloud Resource Manager API](https://console.cloud.google.com/apis/library/cloudresourcemanager.googleapis.com)

## Example Usage

```hcl
resource "google_logging_project_sink" "my-sink" {
  name = "my-pubsub-instance-sink"

  # Can export to pubsub, cloud storage, or bigquery
  destination = "pubsub.googleapis.com/projects/my-project/topics/instance-activity"

  # Log all WARN or higher severity messages relating to instances
  filter = "resource.type = gce_instance AND severity >= WARNING"

  # Use a unique writer (creates a unique service account used for writing)
  unique_writer_identity = true
}
```

A more complete example follows: this creates a compute instance, as well as a log sink that logs all activity to a
cloud storage bucket. Because we are using `unique_writer_identity`, we must grant it access to the bucket. Note that
this grant requires the "Project IAM Admin" IAM role (`roles/resourcemanager.projectIamAdmin`) granted to the credentials
used with terraform.

```hcl
# Our logged compute instance
resource "google_compute_instance" "my-logged-instance" {
  name         = "my-instance"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-9"
    }
  }

  network_interface {
    network = "default"

    access_config {
    }
  }
}

# A bucket to store logs in
resource "google_storage_bucket" "log-bucket" {
  name = "my-unique-logging-bucket"
}

# Our sink; this logs all activity related to our "my-logged-instance" instance
resource "google_logging_project_sink" "instance-sink" {
  name        = "my-instance-sink"
  description = "some explaination on what this is"
  destination = "storage.googleapis.com/${google_storage_bucket.log-bucket.name}"
  filter      = "resource.type = gce_instance AND resource.labels.instance_id = \"${google_compute_instance.my-logged-instance.instance_id}\""

  unique_writer_identity = true
}

# Because our sink uses a unique_writer, we must grant that writer access to the bucket.
resource "google_project_iam_binding" "log-writer" {
  role = "roles/storage.objectCreator"

  members = [
    google_logging_project_sink.instance-sink.writer_identity,
  ]
}
```

The following example uses `exclusions` to filter logs that will not be exported. In this example logs are exported to a [log bucket](https://cloud.google.com/logging/docs/buckets) and there are 2 exclusions configured

```hcl
resource "google_logging_project_sink" "log-bucket" {
  name        = "my-logging-sink"
  destination = "logging.googleapis.com/projects/my-project/locations/global/buckets/_Default"

  exclusions {
		name = "nsexcllusion1"
		description = "Exclude logs from namespace-1 in k8s"
		filter = "resource.type = k8s_container resource.labels.namespace_name=\"namespace-1\" "
	}

	exclusions {
		name = "nsexcllusion2"
		description = "Exclude logs from namespace-2 in k8s"
		filter = "resource.type = k8s_container resource.labels.namespace_name=\"namespace-2\" "
	}

  unique_writer_identity = true
```



## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the logging sink.

* `destination` - (Required) The destination of the sink (or, in other words, where logs are written to). Can be a
    Cloud Storage bucket, a PubSub topic, a BigQuery dataset or a Cloud Logging bucket . Examples:
```
"storage.googleapis.com/[GCS_BUCKET]"
"bigquery.googleapis.com/projects/[PROJECT_ID]/datasets/[DATASET]"
"pubsub.googleapis.com/projects/[PROJECT_ID]/topics/[TOPIC_ID]"
"logging.googleapis.com/projects/[PROJECT_ID]]/locations/global/buckets/[BUCKET_ID]"
```
    The writer associated with the sink must have access to write to the above resource.

* `filter` - (Optional) The filter to apply when exporting logs. Only log entries that match the filter are exported.
    See [Advanced Log Filters](https://cloud.google.com/logging/docs/view/advanced_filters) for information on how to
    write a filter.

* `description` - (Optional) A description of this sink. The maximum length of the description is 8000 characters.

* `disabled` - (Optional) If set to True, then this sink is disabled and it does not export any log entries.

* `project` - (Optional) The ID of the project to create the sink in. If omitted, the project associated with the provider is
    used.

* `unique_writer_identity` - (Optional) Whether or not to create a unique identity associated with this sink. If `false`
    (the default), then the `writer_identity` used is `serviceAccount:cloud-logs@system.gserviceaccount.com`. If `true`,
    then a unique service account is created and used for this sink. If you wish to publish logs across projects or utilize
    `bigquery_options`, you must set `unique_writer_identity` to true.

* `bigquery_options` - (Optional) Options that affect sinks exporting data to BigQuery. Structure documented below.

* `exclusions` - (Optional) Log entries that match any of the exclusion filters will not be exported. If a log entry is matched by both filter and one of exclusion_filters it will not be exported.  Can be repeated multiple times for multiple exclusions. Structure is documented below.

The `bigquery_options` block supports:

* `use_partitioned_tables` - (Required) Whether to use [BigQuery's partition tables](https://cloud.google.com/bigquery/docs/partitioned-tables).
    By default, Logging creates dated tables based on the log entries' timestamps, e.g. syslog_20170523. With partitioned
    tables the date suffix is no longer present and [special query syntax](https://cloud.google.com/bigquery/docs/querying-partitioned-tables)
    has to be used instead. In both cases, tables are sharded based on UTC timezone.

The `exclusions` block support:

* `name` - (Required) A client-assigned identifier, such as `load-balancer-exclusion`. Identifiers are limited to 100 characters and can include only letters, digits, underscores, hyphens, and periods. First character has to be alphanumeric.
* `description` - (Optional) A description of this exclusion.
* `filter` - (Required) An advanced logs filter that matches the log entries to be excluded. By using the sample function, you can exclude less than 100% of the matching log entries. See [Advanced Log Filters](https://cloud.google.com/logging/docs/view/advanced_filters) for information on how to
    write a filter.
* `disabled` - (Optional) If set to True, then this exclusion is disabled and it does not exclude any log entries.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `id` - an identifier for the resource with format `projects/{{project}}/sinks/{{name}}`

* `writer_identity` - The identity associated with this sink. This identity must be granted write access to the
    configured `destination`.

## Import

Project-level logging sinks can be imported using their URI, e.g.

```
$ terraform import google_logging_project_sink.my_sink projects/my-project/sinks/my-sink
```
