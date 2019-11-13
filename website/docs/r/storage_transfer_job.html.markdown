---
subcategory: "Cloud Storage"
layout: "google"
page_title: "Google: google_storage_transfer_job"
sidebar_current: "docs-google-storage-transfer-job-x"
description: |-
  Creates a new Transfer Job in Google Cloud Storage Transfer.
---

# google\_storage\_transfer\_job

Creates a new Transfer Job in Google Cloud Storage Transfer.

To get more information about Google Cloud Storage Transfer, see:

* [Overview](https://cloud.google.com/storage-transfer/docs/overview)
* [API documentation](https://cloud.google.com/storage-transfer/docs/reference/rest/v1/transferJobs#TransferJob)
* How-to Guides
    * [Configuring Access to Data Sources and Sinks](https://cloud.google.com/storage-transfer/docs/configure-access)

## Example Usage

Example creating a nightly Transfer Job from an AWS S3 Bucket to a GCS bucket.

```hcl
data "google_storage_transfer_project_service_account" "default" {
  project = var.project
}

resource "google_storage_bucket" "s3-backup-bucket" {
  name          = "${var.aws_s3_bucket}-backup"
  storage_class = "NEARLINE"
  project       = var.project
}

resource "google_storage_bucket_iam_member" "s3-backup-bucket" {
  bucket     = google_storage_bucket.s3-backup-bucket.name
  role       = "roles/storage.admin"
  member     = "serviceAccount:${data.google_storage_transfer_project_service_account.default.email}"
  depends_on = [google_storage_bucket.s3-backup-bucket]
}

resource "google_storage_transfer_job" "s3-bucket-nightly-backup" {
  description = "Nightly backup of S3 bucket"
  project     = var.project

  transfer_spec {
    object_conditions {
      max_time_elapsed_since_last_modification = "600s"
      exclude_prefixes = [
        "requests.gz",
      ]
    }
    transfer_options {
      delete_objects_unique_in_sink = false
    }
    aws_s3_data_source {
      bucket_name = var.aws_s3_bucket
      aws_access_key {
        access_key_id     = var.aws_access_key
        secret_access_key = var.aws_secret_key
      }
    }
    gcs_data_sink {
      bucket_name = google_storage_bucket.s3-backup-bucket.name
    }
  }

  schedule {
    schedule_start_date {
      year  = 2018
      month = 10
      day   = 1
    }
    schedule_end_date {
      year  = 2019
      month = 1
      day   = 15
    }
    start_time_of_day {
      hours   = 23
      minutes = 30
      seconds = 0
      nanos   = 0
    }
  }

  depends_on = [google_storage_bucket_iam_member.s3-backup-bucket]
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Required) Unique description to identify the Transfer Job.

* `transfer_spec` - (Required) Transfer specification. Structure documented below.

* `schedule` - (Required) Schedule specification defining when the Transfer Job should be scheduled to start, end and and what time to run. Structure documented below.

- - -

* `project` - (Optional) The project in which the resource belongs. If it
	is not provided, the provider project is used.

* `status` - (Optional) Status of the job. Default: `ENABLED`. **NOTE: The effect of the new job status takes place during a subsequent job run. For example, if you change the job status from ENABLED to DISABLED, and an operation spawned by the transfer is running, the status change would not affect the current operation.**

The `transfer_spec` block supports:

* `gcs_data_sink` - (Required) A Google Cloud Storage data sink. Structure documented below.

* `object_conditions` - (Optional) Only objects that satisfy these object conditions are included in the set of data source and data sink objects. Object conditions based on objects' `last_modification_time` do not exclude objects in a data sink. Structure documented below.

* `transfer_options` - (Optional) Characteristics of how to treat files from datasource and sink during job. If the option `delete_objects_unique_in_sink` is true, object conditions based on objects' `last_modification_time` are ignored and do not exclude objects in a data source or a data sink. Structure documented below.

* `gcs_data_source` - (Optional) A Google Cloud Storage data source. Structure documented below.

* `aws_s3_data_source` - (Optional) An AWS S3 data source. Structure documented below.

* `http_data_source` - (Optional) An HTTP URL data source. Structure documented below.

The `schedule` block supports:

* `schedule_start_date` - (Required) The first day the recurring transfer is scheduled to run. If `schedule_start_date` is in the past, the transfer will run for the first time on the following day. Structure documented below.

* `schedule_end_date` - (Optional) The last day the recurring transfer will be run. If `schedule_end_date` is the same as `schedule_start_date`, the transfer will be executed only once. Structure documented below.

* `start_time_of_day` - (Optional) The time in UTC at which the transfer will be scheduled to start in a day. Transfers may start later than this time. If not specified, recurring and one-time transfers that are scheduled to run today will run immediately; recurring transfers that are scheduled to run on a future date will start at approximately midnight UTC on that date. Note that when configuring a transfer with the Cloud Platform Console, the transfer's start time in a day is specified in your local timezone. Structure documented below.

The `object_conditions` block supports:

* `max_time_elapsed_since_last_modification` - (Optional) A duration in seconds with up to nine fractional digits, terminated by 's'. Example: "3.5s".

* `min_time_elapsed_since_last_modification` - (Optional)
A duration in seconds with up to nine fractional digits, terminated by 's'. Example: "3.5s".

* `include_prefixes` - (Optional) If `include_refixes` is specified, objects that satisfy the object conditions must have names that start with one of the `include_prefixes` and that do not start with any of the `exclude_prefixes`. If `include_prefixes` is not specified, all objects except those that have names starting with one of the `exclude_prefixes` must satisfy the object conditions. See [Requirements](https://cloud.google.com/storage-transfer/docs/reference/rest/v1/TransferSpec#ObjectConditions).

* `exclude_prefixes` - (Optional) `exclude_prefixes` must follow the requirements described for `include_prefixes`. See [Requirements](https://cloud.google.com/storage-transfer/docs/reference/rest/v1/TransferSpec#ObjectConditions).

The `transfer_options` block supports:

* `overwrite_objects_already_existing_in_sink` - (Optional) Whether overwriting objects that already exist in the sink is allowed.

* `delete_objects_unique_in_sink` - (Optional) Whether objects that exist only in the sink should be deleted. Note that this option and
`delete_objects_from_source_after_transfer` are mutually exclusive.

* `delete_objects_from_source_after_transfer` - (Optional) Whether objects should be deleted from the source after they are transferred to the sink. Note that this option and `delete_objects_unique_in_sink` are mutually exclusive.

The `gcs_data_sink` block supports:

* `bucket_name` - (Required) Google Cloud Storage bucket name.

The `gcs_data_source` block supports:

* `bucket_name` - (Required) Google Cloud Storage bucket name.

The `aws_s3_data_source` block supports:

* `bucket_name` - (Required) S3 Bucket name.

* `aws_access_key` - (Required) AWS credentials block.

The `aws_access_key` block supports:

* `access_key_id` - (Required) AWS Key ID.

* `secret_access_key` - (Required) AWS Secret Access Key.

The `http_data_source` block supports:

* `list_url` - (Required) The URL that points to the file that stores the object list entries. This file must allow public access. Currently, only URLs with HTTP and HTTPS schemes are supported.

The `schedule_start_date` and `schedule_end_date` blocks support:

* `year` - (Required) Year of date. Must be from 1 to 9999.

* `month` - (Required) Month of year. Must be from 1 to 12.

* `day` - (Required) Day of month. Must be from 1 to 31 and valid for the year and month.

The `start_time_of_day` blocks support:

* `hours` - (Required) Hours of day in 24 hour format. Should be from 0 to 23

* `minutes` - (Required) Minutes of hour of day. Must be from 0 to 59.

* `seconds` - (Optional) Seconds of minutes of the time. Must normally be from 0 to 59.

* `nanos` - (Required) Fractions of seconds in nanoseconds. Must be from 0 to 999,999,999.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `name` - The name of the Transfer Job.

* `creation_time` - When the Transfer Job was created.

* `last_modification_time` - When the Transfer Job was last modified.

* `deletion_time` - When the Transfer Job was deleted.

## Import

Storage buckets can be imported using the Transfer Job's `project` and `name` without the `transferJob/` prefix, e.g.

```
$ terraform import google_storage_transfer_job.nightly-backup-transfer-job my-project-1asd32/8422144862922355674
```
