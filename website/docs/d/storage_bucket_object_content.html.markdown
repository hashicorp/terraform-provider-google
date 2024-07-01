---
subcategory: "Cloud Storage"
description: |-
  Get content of a Google Cloud Storage bucket object.
---


# google_storage_bucket_object_content

Gets an existing object content inside an existing bucket in Google Cloud Storage service (GCS).
See [the official documentation](https://cloud.google.com/storage/docs/key-terms#objects)
and
[API](https://cloud.google.com/storage/docs/json_api/v1/objects).

~> **Warning:** The object content will be saved in the state, and visiable to everyone who has access to the state file.

## Example Usage

Example file object  stored within a folder.

```hcl
data "google_storage_bucket_object_content" "key" {
  name   = "encryptedkey"
  bucket = "keystore"
}

output "encrypted" {
  value = data.google_storage_bucket_object_content.key.content
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required) The name of the containing bucket.

* `name` - (Required) The name of the object.

## Attributes Reference

The following attributes are exported:

* `content` - (Computed) The content of the object.
