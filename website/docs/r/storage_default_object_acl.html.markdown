---
layout: "google"
page_title: "Google: google_storage_default_object_acl"
sidebar_current: "docs-google-storage-default-object-acl"
description: |-
  Creates a new default object ACL in Google Cloud Storage.
---

# google\_storage\_default\_object\_acl

Creates a new default object ACL in Google Cloud Storage service (GCS). For more information see
[the official documentation](https://cloud.google.com/storage/docs/access-control/lists) 
and 
[API](https://cloud.google.com/storage/docs/json_api/v1/defaultObjectAccessControls).

## Example Usage

Example creating a default object ACL on a bucket with one owner, and one reader.

```hcl
resource "google_storage_bucket" "image-store" {
  name     = "image-store-bucket"
  location = "EU"
}

resource "google_storage_default_object_acl" "image-store-default-acl" {
  bucket = "${google_storage_bucket.image-store.name}"
  role_entity = [
    "OWNER:user-my.email@gmail.com",
    "READER:group-mygroup",
  ]
}
```

## Argument Reference

* `bucket` - (Required) The name of the bucket it applies to.

* `role_entity` - (Required) List of role/entity pairs in the form `ROLE:entity`. See [GCS Object ACL documentation](https://cloud.google.com/storage/docs/json_api/v1/objectAccessControls) for more details.

## Attributes Reference

Only the arguments listed above are exposed as attributes.
