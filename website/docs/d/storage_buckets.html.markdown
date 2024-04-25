---
subcategory: "Cloud Storage"
description: |-
  Retrieve information about a set of GCS buckets in a project.
---


# google\_storage\_buckets

Gets a list of existing GCS buckets.
See [the official documentation](https://cloud.google.com/storage/docs/introduction)
and [API](https://cloud.google.com/storage/docs/json_api/v1/buckets/list).

## Example Usage

Example GCS buckets.

```hcl
data "google_storage_buckets" "example" {
  project = "example-project"
}
```

## Argument Reference

The following arguments are supported:

* `prefix` - (Optional) Filter results to buckets whose names begin with this prefix.
* `project` - (Optional) The ID of the project. If it is not provided, the provider project is used.


## Attributes Reference

The following attributes are exported:

* `buckets` - A list of all retrieved GCS buckets. Structure is [defined below](#nested_buckets).

<a name="nested_buckets"></a>The `buckets` block supports:

* `labels` - User-provided bucket labels, in key/value pairs.
* `location` - The location of the bucket. 
* `name` - The name of the bucket.
* `self_link` - A url reference to the bucket.
* `storage_class` - The [StorageClass](https://cloud.google.com/storage/docs/storage-classes) of the bucket.
