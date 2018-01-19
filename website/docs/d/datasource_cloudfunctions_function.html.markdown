---
layout: "google"
page_title: "Google: google_cloudfunctions_function"
sidebar_current: "docs-google-datasource-cloudfunctions-function"
description: |-
  Get information about a Google Cloud Function.
---

# google\_cloudfunctions\_function

Get information about a Google Cloud Function. For more information see
the [official documentation](https://cloud.google.com/functions/docs/)
and [API](https://cloud.google.com/functions/docs/apis).

## Example Usage

```hcl
data "google_cloudfunctions_function" "my-function" {
  name = "function"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of a Cloud Function.

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

* `region` - (Optional) The region in which the resource belongs. If it
    is not provided, the provider region is used.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `name` - The name of the Cloud Function.
* `source_archive_bucket` - The GCS bucket containing the zip archive which contains the function.
* `source_archive_object` - The source archive object (file) in archive bucket.
* `description` - Description of the function.
* `available_memory_mb` - Available memory (in MB) to the function.
* `timeout` - Function execution timeout (in seconds).
* `entry_point` - Name of a JavaScript function that will be executed when the Google Cloud Function is triggered.
* `trigger_http` - If function is triggered by HTTP, this boolean is set.
* `trigger_bucket` - If function is triggered by bucket, bucket name is set here.
* `trigger_topic` - If function is triggered by Pub/Sub topic, name of topic is set here.
* `https_trigger_url` - If function is triggered by HTTP, trigger URL is set here.
* `labels` - A map of labels applied to this function.
