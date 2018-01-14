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
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A unique name for the resource, required by GCE.
    Changing this forces a new resource to be created.

* `source_archive_bucket` - (Required) The source archive bucket. Changing this forces a new resource to be created.

* `source_archive_object` - (Required) The source archive object (file) in archive bucket. Changing this forces a new resource to be created.

- - -

* `description` - (Optional) Description, which you would like to add to your function.

* `available_memory_mb` - (Optional) Memory (in MB), available to you function. Default value is 256MB. Allowed values are: 128MB, 256MB, 512MB, 1024MB, and 2048MB.

* `timeout` - (Optional) Timeout (in seconds) for the function. Default value is 60 seconds. Cannot be more than 540 seconds.

* `entry_point` - (Optional) Name of a JavaScript function that will be executed when the Google Cloud Function is triggered.

* `trigger_http` - (Optional) Boolean variable. Any HTTP request (of a supported type) to the endpoint will trigger function execution. Supported HTTP request types are: POST, PUT, GET, DELETE, and OPTIONS. Endpoint would be returned as `https_trigger_url`. Blocks usage of `trigger_bucket` and `trigger_topic`.

* `trigger_bucket` - (Optional) Google Cloud Storage bucket name. Every change in files in this bucket will trigger function execution. Blocks usage of `trigger_http` and `trigger_topic`.

* `trigger_topic` - (Optional) Name of Pub/Sub topic. Every message published in this topic will trigger function execution with message contents passed as input data. Blocks usage of `trigger_http` and `trigger_bucket`.

* `labels` - (Optional) A set of key/value label pairs to assign to the function.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `https_trigger_url` - URL, which could trigger function execution. Returned only if `trigger_http` is used.

* `project` - GCE project of the function. If it is not provided, the provider region is used.

* `region` - GCE region of function. Currently can be only "us-central1". If it is not provided, the provider region is used.

## Import

Functions can be imported using the `name`, e.g.

```
$ terraform import google_cloudfunctions_function.default test-function
```
