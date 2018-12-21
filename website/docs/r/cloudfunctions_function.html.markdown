---
layout: "google"
page_title: "Google: google_cloudfunctions_function"
sidebar_current: "docs-google-cloudfunctions-function"
description: |-
  Creates a new Cloud Function.
---

# google\_cloudfunctions\_function

Creates a new Cloud Function. For more information see
[the official documentation](https://cloud.google.com/functions/docs/)
and
[API](https://cloud.google.com/functions/docs/apis).

## Example Usage

```hcl
resource "google_storage_bucket" "bucket" {
  name = "test-bucket"
}

resource "google_storage_bucket_object" "archive" {
  name   = "index.zip"
  bucket = "${google_storage_bucket.bucket.name}"
  source = "./path/to/zip/file/which/contains/code"
}

resource "google_cloudfunctions_function" "function" {
  name                  = "function-test"
  description           = "My function"
  available_memory_mb   = 128
  source_archive_bucket = "${google_storage_bucket.bucket.name}"
  source_archive_object = "${google_storage_bucket_object.archive.name}"
  trigger_http          = true
  timeout               = 60
  entry_point           = "helloGET"
  labels {
    my-label = "my-label-value"
  }
  environment_variables {
    MY_ENV_VAR = "my-env-var-value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A user-defined name of the function. Function names must be unique globally.

* `source_archive_bucket` - (Required) The GCS bucket containing the zip archive which contains the function.

* `source_archive_object` - (Required) The source archive object (file) in archive bucket.

- - -

* `description` - (Optional) Description of the function.

* `available_memory_mb` - (Optional) Memory (in MB), available to the function. Default value is 256MB. Allowed values are: 128MB, 256MB, 512MB, 1024MB, and 2048MB.

* `timeout` - (Optional) Timeout (in seconds) for the function. Default value is 60 seconds. Cannot be more than 540 seconds.

* `entry_point` - (Optional) Name of a JavaScript function that will be executed when the Google Cloud Function is triggered.

* `event_trigger` - (Optional) A source that fires events in response to a condition in another service. Structure is documented below. Cannot be used with `trigger_http`.

* `trigger_http` - (Optional) Boolean variable. Any HTTP request (of a supported type) to the endpoint will trigger function execution. Supported HTTP request types are: POST, PUT, GET, DELETE, and OPTIONS. Endpoint is returned as `https_trigger_url`. Cannot be used with `trigger_bucket` and `trigger_topic`.

* `labels` - (Optional) A set of key/value label pairs to assign to the function.

* `runtime` - (Optional) The runtime in which the function is going to run. If empty, defaults to `"nodejs6"`.

* `environment_variables` - (Optional) A set of key/value environment variable pairs to assign to the function.

The `event_trigger` block supports:

* `event_type` - (Required) The type of event to observe. For example: `"providers/cloud.storage/eventTypes/object.change"`
    and `"providers/cloud.pubsub/eventTypes/topic.publish"`. See the documentation on [calling Cloud Functions](https://cloud.google.com/functions/docs/calling/)
    for a full reference. Only Cloud Storage and Cloud Pub/Sub triggers are supported at this time.

* `resource` - (Required) Required. The name of the resource from which to observe events, for example, `"myBucket"`   

* `failure_policy` - (Optional) Specifies policy for failed executions. Structure is documented below.

The `failure_policy` block supports:

* `retry` - (Required) Whether the function should be retried on failure. Defaults to `false`.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `https_trigger_url` - URL which triggers function execution. Returned only if `trigger_http` is used.

* `project` - Project of the function. If it is not provided, the provider project is used.

* `region` - Region of function. Currently can be only "us-central1". If it is not provided, the provider region is used.

## Import

Functions can be imported using the `name`, e.g.

```
$ terraform import google_cloudfunctions_function.default function-test
```
