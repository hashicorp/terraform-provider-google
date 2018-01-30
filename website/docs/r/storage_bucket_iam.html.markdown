---
layout: "google"
page_title: "Google: google_storage_bucket_iam"
sidebar_current: "docs-google-storage-bucket-iam"
description: |-
 Collection of resources to manage IAM policy for a Google storage bucket.
---

# IAM policy for Google storage bucket

Two different resources help you manage your IAM policy for storage bucket. Each of these resources serves a different use case:

* `google_storage_bucket_iam_binding`: Authoritative for a given role. Updates the IAM policy to grant a role to a list of members. Other roles within the IAM policy for the storage bucket are preserved.
* `google_storage_bucket_iam_member`: Non-authoritative. Updates the IAM policy to grant a role to a new member. Other members for the role for the storage bucket are preserved.

~> **Note:** `google_storage_bucket_iam_binding` resources **can be** used in conjunction with `google_storage_bucket_iam_member` resources **only if** they do not grant privilege to the same role.

## google\_storage\_bucket\_iam\_binding

```hcl
resource "google_storage_bucket_iam_binding" "binding" {
  bucket = "your-bucket-name"
  role        = "roles/storage.objectViewer"

  members = [
    "user:jane@example.com",
  ]
}
```

## google\_storage\_bucket\_iam\_member

```hcl
resource "google_storage_bucket_iam_member" "member" {
  bucket = "your-bucket-name"
  role        = "roles/storage.objectViewer"
  member      = "user:jane@example.com"
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required) The name of the bucket it applies to.

* `member/members` - (Required) Identities that will be granted the privilege in `role`.
  Each entry can have one of the following values:
  * **allUsers**: A special identifier that represents anyone who is on the internet; with or without a Google account.
  * **allAuthenticatedUsers**: A special identifier that represents anyone who is authenticated with a Google account or a service account.
  * **user:{emailid}**: An email address that represents a specific Google account. For example, alice@gmail.com or joe@example.com.
  * **serviceAccount:{emailid}**: An email address that represents a service account. For example, my-other-app@appspot.gserviceaccount.com.
  * **group:{emailid}**: An email address that represents a Google group. For example, admins@example.com.
  * **domain:{domain}**: A Google Apps domain name that represents all the users of that domain. For example, google.com or example.com.

* `role` - (Required) The role that should be applied. Note that custom roles must be of the format
    `[projects|organizations]/{parent-name}/roles/{role-name}`.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `etag` - (Computed) The etag of the storage bucket's IAM policy.