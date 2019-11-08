---
subcategory: "Cloud Storage"
layout: "google"
page_title: "Google: google_storage_bucket_iam"
sidebar_current: "docs-google-storage-bucket-iam"
description: |-
 Collection of resources to manage IAM policy for a Google storage bucket.
---

# IAM policy for Google storage bucket

Three different resources help you manage your IAM policy for storage bucket. Each of these resources serves a different use case:

* `google_storage_bucket_iam_binding`: Authoritative for a given role. Updates the IAM policy to grant a role to a list of members. Other roles within the IAM policy for the storage bucket are preserved.
* `google_storage_bucket_iam_member`: Non-authoritative. Updates the IAM policy to grant a role to a new member. Other members for the role for the storage bucket are preserved.
* `google_storage_bucket_iam_policy`: Setting a policy removes all other permissions on the bucket, and if done incorrectly, there's a real chance you will lock yourself out of the bucket. If possible for your use case, using multiple google_storage_bucket_iam_binding resources will be much safer. See the usage example on how to work with policy correctly.


~> **Note:** `google_storage_bucket_iam_binding` resources **can be** used in conjunction with `google_storage_bucket_iam_member` resources **only if** they do not grant privilege to the same role.

## google\_storage\_bucket\_iam\_binding

```hcl
resource "google_storage_bucket_iam_binding" "binding" {
  bucket = "your-bucket-name"
  role   = "roles/storage.objectViewer"

  members = [
    "user:jane@example.com",
  ]
}
```

## google\_storage\_bucket\_iam\_member

```hcl
resource "google_storage_bucket_iam_member" "member" {
  bucket = "your-bucket-name"
  role   = "roles/storage.objectViewer"
  member = "user:jane@example.com"
}
```

## google\_storage\_bucket\_iam\_policy

When applying a policy that does not include the roles listed below, you lose the default permissions which google adds to your bucket:
* `roles/storage.legacyBucketOwner`
* `roles/storage.legacyBucketReader`

If this happens only an entity with `roles/storage.admin` privileges can repair this bucket's policies. It is recommended to include the above roles in policies to get the same behaviour as with the other two options.

```hcl
data "google_iam_policy" "foo-policy" {
  binding {
    role = "roles/your-role"

    members = ["group:yourgroup@example.com"]
  }
}

resource "google_storage_bucket_iam_policy" "member" {
  bucket      = "your-bucket-name"
  policy_data = data.google_iam_policy.foo-policy.policy_data
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
  * **domain:{domain}**: A G Suite domain (primary, instead of alias) name that represents all the users of that domain. For example, google.com or example.com.

* `role` - (Required) The role that should be applied. Note that custom roles must be of the format
    `[projects|organizations]/{parent-name}/roles/{role-name}`.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `etag` - (Computed) The etag of the storage bucket's IAM policy.


## Import

For `google_storage_bucket_iam_policy`:

IAM member imports use space-delimited identifiers - generally the resource in question, the role, and the member identity (i.e. `serviceAccount: my-sa@my-project.iam.gserviceaccount.com` or `user:foo@example.com`). Policies, bindings, and members can be respectively imported as follows:

```
$ terraform import google_storage_bucket_iam_policy.policy "my-bucket user:foo@example.com"

$ terraform import google_storage_bucket_iam_binding.binding "my-bucket roles/my-role "

$ terraform import google_storage_bucket_iam_member.member "my-bucket roles/my-role user:foo@example.com"
```
