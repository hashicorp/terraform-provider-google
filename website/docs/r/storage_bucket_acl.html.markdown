---
subcategory: "Cloud Storage"
page_title: "Google: google_storage_bucket_acl"
description: |-
  Creates a new bucket ACL in Google Cloud Storage.
---

# google\_storage\_bucket\_acl

Authoritatively manages a bucket's ACLs in Google cloud storage service (GCS). For more information see
[the official documentation](https://cloud.google.com/storage/docs/access-control/lists)
and
[API](https://cloud.google.com/storage/docs/json_api/v1/bucketAccessControls).

Bucket ACLs can be managed non authoritatively using the [`storage_bucket_access_control`](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/storage_bucket_access_control) resource. Do not use these two resources in conjunction to manage the same bucket.

Permissions can be granted either by ACLs or Cloud IAM policies. In general, permissions granted by Cloud IAM policies do not appear in ACLs, and permissions granted by ACLs do not appear in Cloud IAM policies. The only exception is for ACLs applied directly on a bucket and certain bucket-level Cloud IAM policies, as described in [Cloud IAM relation to ACLs](https://cloud.google.com/storage/docs/access-control/iam#acls).

**NOTE** This resource will not remove the `project-owners-<project_id>` entity from the `OWNER` role.

## Example Usage

Example creating an ACL on a bucket with one owner, and one reader.

```hcl
resource "google_storage_bucket" "image-store" {
  name     = "image-store-bucket"
  location = "EU"
}

resource "google_storage_bucket_acl" "image-store-acl" {
  bucket = google_storage_bucket.image-store.name

  role_entity = [
    "OWNER:user-my.email@gmail.com",
    "READER:group-mygroup",
  ]
}
```

## Argument Reference

* `bucket` - (Required) The name of the bucket it applies to.

- - -

* `predefined_acl` - (Optional) The [canned GCS ACL](https://cloud.google.com/storage/docs/access-control/lists#predefined-acl) to apply. Must be set if `role_entity` is not.

* `role_entity` - (Optional) List of role/entity pairs in the form `ROLE:entity`. See [GCS Bucket ACL documentation](https://cloud.google.com/storage/docs/json_api/v1/bucketAccessControls)  for more details. Must be set if `predefined_acl` is not.

* `default_acl` - (Optional) Configure this ACL to be the default ACL.

## Attributes Reference

Only the arguments listed above are exposed as attributes.

## Import

This resource does not support import.
