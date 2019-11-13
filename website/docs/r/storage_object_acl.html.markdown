---
subcategory: "Cloud Storage"
layout: "google"
page_title: "Google: google_storage_object_acl"
sidebar_current: "docs-google-storage-object-acl"
description: |-
  Creates a new object ACL in Google Cloud Storage.
---

# google\_storage\_object\_acl

Authoritatively manages the access control list (ACL) for an object in a Google
Cloud Storage (GCS) bucket. Removing a `google_storage_object_acl` sets the
acl to the `private` [predefined ACL](https://cloud.google.com/storage/docs/access-control#predefined-acl).

For more information see
[the official documentation](https://cloud.google.com/storage/docs/access-control/lists) 
and 
[API](https://cloud.google.com/storage/docs/json_api/v1/objectAccessControls).

-> Want fine-grained control over object ACLs? Use `google_storage_object_access_control` to control individual
role entity pairs.

## Example Usage

Create an object ACL with one owner and one reader.

```hcl
resource "google_storage_bucket" "image-store" {
  name     = "image-store-bucket"
  location = "EU"
}

resource "google_storage_bucket_object" "image" {
  name   = "image1"
  bucket = google_storage_bucket.image-store.name
  source = "image1.jpg"
}

resource "google_storage_object_acl" "image-store-acl" {
  bucket = google_storage_bucket.image-store.name
  object = google_storage_bucket_object.image.output_name

  role_entity = [
    "OWNER:user-my.email@gmail.com",
    "READER:group-mygroup",
  ]
}
```

## Argument Reference

* `bucket` - (Required) The name of the bucket the object is stored in.

* `object` - (Required) The name of the object to apply the acl to.

- - -

* `predefined_acl` - (Optional) The "canned" [predefined ACL](https://cloud.google.com/storage/docs/access-control#predefined-acl) to apply. Must be set if `role_entity` is not.

* `role_entity` - (Optional) List of role/entity pairs in the form `ROLE:entity`. See [GCS Object ACL documentation](https://cloud.google.com/storage/docs/json_api/v1/objectAccessControls) for more details.
Must be set if `predefined_acl` is not.

-> The object's creator will always have `OWNER` permissions for their object, and any attempt to modify that permission would return an error. Instead, Terraform automatically
adds that role/entity pair to your `terraform plan` results when it is omitted in your config; `terraform plan` will show the correct final state at every point except for at
`Create` time, where the object role/entity pair is omitted if not explicitly set.


## Attributes Reference

Only the arguments listed above are exposed as attributes.
