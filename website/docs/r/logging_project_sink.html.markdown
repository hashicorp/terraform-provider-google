---
subcategory: "Cloud (Stackdriver) Logging"
description: |-
  Manages a project-level logging sink.
---

# google\_logging\_project\_sink

Manages a project-level logging sink. For more information see:

* [API documentation](https://cloud.google.com/logging/docs/reference/v2/rest/v2/projects.sinks)
* How-to Guides
    * [Exporting Logs](https://cloud.google.com/logging/docs/export)

~> You can specify exclusions for log sinks created by terraform by using the exclusions field of `google_logging_folder_sink`

~> **Note:** You must have [granted the "Logs Configuration Writer"](https://cloud.google.com/logging/docs/access-control) IAM role (`roles/logging.configWriter`) to the credentials used with terraform.

~> **Note** You must [enable the Cloud Resource Manager API](https://console.cloud.google.com/apis/library/cloudresourcemanager.googleapis.com)

~> **Note:** The `_Default` and `_Required` logging sinks are automatically created for a given project and cannot be deleted. Creating a resource of this type will acquire and update the resource that already exists at the desired location. These sinks cannot be removed so deleting this resource will remove the sink config from your terraform state but will leave the logging sink unchanged. The sinks that are currently automatically created are "_Default" and "_Required".


## Example Usage - Basic Sink

```hcl
resource "google_logging_project_sink" "my-sink" {
  name = "my-pubsub-instance-sink"

  # Can export to pubsub, cloud storage, bigquery, log bucket, or another project
  destination = "pubsub.googleapis.com/projects/my-project/topics/instance-activity"

  # Log all WARN or higher severity messages relating to instances
  filter = "resource.type = gce_instance AND severity >= WARNING"

  # Use a unique writer (creates a unique service account used for writing)
  unique_writer_identity = true
}
```

## Example Usage - Cloud Storage Bucket Destination

A more complete example follows: this creates a compute instance, as well as a log sink that logs all activity to a
cloud storage bucket. Because we are using `unique_writer_identity`, we must grant it access to the bucket.

Note that this grant requires the "Project IAM Admin" IAM role (`roles/resourcemanager.projectIamAdmin`) granted to the
credentials used with Terraform.

```hcl
# Our logged compute instance
resource "google_compute_instance" "my-logged-instance" {
  name         = "my-instance"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
    }
  }

  network_interface {
    network = "default"

    access_config {
    }
  }
}

# A gcs bucket to store logs in
resource "google_storage_bucket" "gcs-bucket" {
  name     = "my-unique-logging-bucket"
  location = "US"
}

# Our sink; this logs all activity related to our "my-logged-instance" instance
resource "google_logging_project_sink" "instance-sink" {
  name        = "my-instance-sink"
  description = "some explanation on what this is"
  destination = "storage.googleapis.com/${google_storage_bucket.gcs-bucket.name}"
  filter      = "resource.type = gce_instance AND resource.labels.instance_id = \"${google_compute_instance.my-logged-instance.instance_id}\""

  unique_writer_identity = true
}

# Because our sink uses a unique_writer, we must grant that writer access to the bucket.
resource "google_project_iam_binding" "gcs-bucket-writer" {
  project = "your-project-id"
  role = "roles/storage.objectCreator"

  members = [
    google_logging_project_sink.instance-sink.writer_identity,
  ]
}
```

## Example Usage - User-managed Service Account 

The following example creates a sink that are configured with user-managed service accounts, by specifying
the `custom_writer_identity` field.

Note that you can only create a sink that uses a user-managed service account when the sink destination
is a log bucket.

```hcl
resource "google_service_account" "custom-sa" {
  project      = "other-project-id"
  account_id   = "gce-log-bucket-sink"
  display_name = "gce-log-bucket-sink"
}

# Create a sink that uses user-managed service account
resource "google_logging_project_sink" "my-sink" {
  name = "other-project-log-bucket-sink"

  # Can export to log bucket in another project
  destination = "logging.googleapis.com/projects/other-project-id/locations/global/buckets/gce-logs"

  # Log all WARN or higher severity messages relating to instances
  filter = "resource.type = gce_instance AND severity >= WARNING"

  unique_writer_identity = true
  
  # Use a user-managed service account
  custom_writer_identity = google_service_account.custom-sa.email
}

# grant writer access to the user-managed service account
resource "google_project_iam_member" "custom-sa-logbucket-binding" {
  project = "destination-project-id"
  role   = "roles/logging.bucketWriter"
  member = "serviceAccount:${google_service_account.custom-sa.email}"
}
```

The above example will create a log sink that route logs to destination GCP project using
an user-managed service account. 

## Example Usage - Sink Exclusions

The following example uses `exclusions` to filter logs that will not be exported. In this example logs are exported to a [log bucket](https://cloud.google.com/logging/docs/buckets) and there are 2 exclusions configured

```hcl
resource "google_logging_project_sink" "log-bucket" {
  name        = "my-logging-sink"
  destination = "logging.googleapis.com/projects/my-project/locations/global/buckets/_Default"

  exclusions {
    name        = "nsexcllusion1"
    description = "Exclude logs from namespace-1 in k8s"
    filter      = "resource.type = k8s_container resource.labels.namespace_name=\"namespace-1\" "
  }

  exclusions {
    name        = "nsexcllusion2"
    description = "Exclude logs from namespace-2 in k8s"
    filter      = "resource.type = k8s_container resource.labels.namespace_name=\"namespace-2\" "
  }

  unique_writer_identity = true
}
```



## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the logging sink. Logging automatically creates two sinks: `_Required` and `_Default`.

* `destination` - (Required) The destination of the sink (or, in other words, where logs are written to). Can be a Cloud Storage bucket, a PubSub topic, a BigQuery dataset, a Cloud Logging bucket, or a Google Cloud project. Examples:

    - `storage.googleapis.com/[GCS_BUCKET]`
    - `bigquery.googleapis.com/projects/[PROJECT_ID]/datasets/[DATASET]`
    - `pubsub.googleapis.com/projects/[PROJECT_ID]/topics/[TOPIC_ID]`
    - `logging.googleapis.com/projects/[PROJECT_ID]/locations/global/buckets/[BUCKET_ID]`
    - `logging.googleapis.com/projects/[PROJECT_ID]`

    The writer associated with the sink must have access to write to the above resource.

* `filter` - (Optional) The filter to apply when exporting logs. Only log entries that match the filter are exported.
    See [Advanced Log Filters](https://cloud.google.com/logging/docs/view/advanced_filters) for information on how to
    write a filter.

* `description` - (Optional) A description of this sink. The maximum length of the description is 8000 characters.

* `disabled` - (Optional) If set to True, then this sink is disabled and it does not export any log entries.

* `project` - (Optional) The ID of the project to create the sink in. If omitted, the project associated with the provider is
    used.

* `unique_writer_identity` - (Optional) Whether or not to create a unique identity associated with this sink. If `false`, then the `writer_identity` used is `serviceAccount:cloud-logs@system.gserviceaccount.com`. If `true` (the default),
    then a unique service account is created and used for this sink. If you wish to publish logs across projects or utilize
    `bigquery_options`, you must set `unique_writer_identity` to true.

* `custom_writer_identity` - (Optional) A user managed service account that will be used to write
    the log entries. The format must be `serviceAccount:some@email`. This field can only be specified if you are
    routing logs to a destination outside this sink's project. If not specified, a Logging service account 
    will automatically be generated.

* `bigquery_options` - (Optional) Options that affect sinks exporting data to BigQuery. Structure [documented below](#nested_bigquery_options).

* `exclusions` - (Optional) Log entries that match any of the exclusion filters will not be exported. If a log entry is matched by both `filter` and one of `exclusions.filter`, it will not be exported.  Can be repeated multiple times for multiple exclusions. Structure is [documented below](#nested_exclusions).

<a name="nested_bigquery_options"></a>The `bigquery_options` block supports:

* `use_partitioned_tables` - (Required) Whether to use [BigQuery's partition tables](https://cloud.google.com/bigquery/docs/partitioned-tables).
    By default, Logging creates dated tables based on the log entries' timestamps, e.g. `syslog_20170523`. With partitioned
    tables the date suffix is no longer present and [special query syntax](https://cloud.google.com/bigquery/docs/querying-partitioned-tables)
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

* `id` - an identifier for the resource with format `projects/{{project}}/sinks/{{name}}`

* `writer_identity` - The identity associated with this sink. This identity must be granted write access to the
    configured `destination`.

## Import

Project-level logging sinks can be imported using their URI, e.g.

* `projects/{{project_id}}/sinks/{{name}}`

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import project-level logging sinks using one of the formats above. For example:

```tf
import {
  id = "projects/{{project_id}}/sinks/{{name}}"
  to = google_logging_project_sink.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), project-level logging sinks can be imported using one of the formats above. For example:

```
$ terraform import google_logging_project_sink.default projects/{{project_id}}/sinks/{{name}}
```
