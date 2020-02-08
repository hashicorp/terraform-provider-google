---
subcategory: "Cloud Storage"
layout: "google"
page_title: "Google: google_storage_default_object_acl"
sidebar_current: "docs-google-storage-default-object-acl"
description: |-
  Authoritatively manages the default object ACLs for a Google Cloud Storage bucket
---

# google\_storage\_default\_object\_acl

Authoritatively manages the default object ACLs for a Google Cloud Storage bucket
without managing the bucket itself.

-> Note that for each object, its creator will have the `"OWNER"` role in addition
to the default ACL that has been defined.

For more information see
[the official documentation](https://cloud.google.com/storage/docs/access-control/lists) 
and 
[API](https://cloud.google.com/storage/docs/json_api/v1/defaultObjectAccessControls).

-> Want fine-grained control over default object ACLs? Use `google_storage_default_object_access_control`
to control individual role entity pairs.

## Example Usage

Example creating a default object ACL on a bucket with one owner, and one reader.

```hcl
resource "google_storage_bucket" "image-store" {
  name     = "image-store-bucket"
  location = "EU"
}

resource "google_storage_default_object_acl" "image-store-default-acl" {
  bucket = google_storage_bucket.image-store.name
  role_entity = [
    "OWNER:user-my.email@gmail.com",
    "READER:group-mygroup",
  ]
}
```

## Argument Reference

* `bucket` - (Required) The name of the bucket it applies to.

---

* `role_entity` - (Optional) List of role/entity pairs in the form `ROLE:entity`.
See [GCS Object ACL documentation](https://cloud.google.com/storage/docs/json_api/v1/objectAccessControls) for more details.
Omitting the field is the same as providing an empty list.

## Attributes Reference

Only the arguments listed above are exposed as attributes.
