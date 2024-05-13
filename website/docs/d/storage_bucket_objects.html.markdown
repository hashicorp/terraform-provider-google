---
subcategory: "Cloud Storage"
description: |-
  Retrieve information about a set of GCS bucket objects in a GCS bucket.
---


# google_storage_bucket_objects

Gets existing objects inside an existing bucket in Google Cloud Storage service (GCS).
See [the official documentation](https://cloud.google.com/storage/docs/key-terms#objects)
and [API](https://cloud.google.com/storage/docs/json_api/v1/objects/list).

## Example Usage

Example files stored within a bucket.

```hcl
data "google_storage_bucket_objects" "files" {
  bucket = "file-store"
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required) The name of the containing bucket.
* `match_glob` - (Optional) A glob pattern used to filter results (for example, `foo*bar`).
* `prefix` - (Optional) Filter results to include only objects whose names begin with this prefix.


## Attributes Reference

The following attributes are exported:

* `bucket_objects` - A list of retrieved objects contained in the provided GCS bucket. Structure is [defined below](#nested_bucket_objects).

<a name="nested_bucket_objects"></a>The `bucket_objects` block supports:

* `content_type` - [Content-Type](https://tools.ietf.org/html/rfc7231#section-3.1.1.5) of the object data.
* `media_link` - A url reference to download this object.
* `name` - The name of the object.
* `self_link` - A url reference to this object.
* `storage_class` - The [StorageClass](https://cloud.google.com/storage/docs/storage-classes) of the bucket object.
