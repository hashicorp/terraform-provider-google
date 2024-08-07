---
subcategory: "Cloud (Stackdriver) Logging"
description: |-
  Manages a folder-level logging bucket config.
---

# google_logging_folder_bucket_config

Manages a folder-level logging bucket config. For more information see
[the official logging documentation](https://cloud.google.com/logging/docs/) and
[Storing Logs](https://cloud.google.com/logging/docs/storage).

~> **Note:** Logging buckets are automatically created for a given folder, project, organization, billingAccount and cannot be deleted. Creating a resource of this type will acquire and update the resource that already exists at the desired location. These buckets cannot be removed so deleting this resource will remove the bucket config from your terraform state but will leave the logging bucket unchanged. The buckets that are currently automatically created are "_Default" and "_Required".

## Example Usage

```hcl
resource "google_folder" "default" {
  display_name = "some-folder-name"
  parent       = "organizations/123456789"
}

resource "google_logging_folder_bucket_config" "basic" {
  folder         = google_folder.default.name
  location       = "global"
  retention_days = 30
  bucket_id      = "_Default"
  
  index_configs {
    field_path = "jsonPayload.request.status"
    type       = "INDEX_TYPE_STRING"
  }
}
```

## Argument Reference

The following arguments are supported:

* `folder` - (Required) The parent resource that contains the logging bucket.

* `location` - (Required) The location of the bucket.

* `bucket_id` - (Required) The name of the logging bucket. Logging automatically creates two log buckets: `_Required` and `_Default`.

* `description` - (Optional) Describes this bucket.

* `retention_days` - (Optional) Logs will be retained by default for this amount of time, after which they will automatically be deleted. The minimum retention period is 1 day. If this value is set to zero at bucket creation time, the default time of 30 days will be used. Bucket retention can not be increased on buckets outside of projects.

* `index_configs` - (Optional) A list of indexed fields and related configuration data. Structure is [documented below](#nested_index_configs).

<a name="nested_index_configs"></a>The `index_configs` block supports:

* `field_path` - The LogEntry field path to index.
  Note that some paths are automatically indexed, and other paths are not eligible for indexing. See [indexing documentation](https://cloud.google.com/logging/docs/analyze/custom-index) for details.

* `type` - The type of data in this index. Allowed types include `INDEX_TYPE_UNSPECIFIED`, `INDEX_TYPE_STRING` and `INDEX_TYPE_INTEGER`.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `id` - an identifier for the resource with format `folders/{{folder}}/locations/{{location}}/buckets/{{bucket_id}}`

* `name` -  The resource name of the bucket. For example: "folders/my-folder-id/locations/my-location/buckets/my-bucket-id"

* `lifecycle_state` -  The bucket's lifecycle such as active or deleted. See [LifecycleState](https://cloud.google.com/logging/docs/reference/v2/rest/v2/billingAccounts.buckets#LogBucket.LifecycleState).

## Import

This resource can be imported using the following format:

* `folders/{{folder}}/locations/{{location}}/buckets/{{bucket_id}}`

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import this resource using one of the formats above. For example:

```tf
import {
  id = "folders/{{folder}}/locations/{{location}}/buckets/{{bucket_id}}"
  to = google_logging_folder_bucket_config.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), this resource can be imported using one of the formats above. For example:

```
$ terraform import google_logging_folder_bucket_config.default folders/{{folder}}/locations/{{location}}/buckets/{{bucket_id}}
```
