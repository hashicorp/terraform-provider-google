---
subcategory: "Cloud Platform"
layout: "google"
page_title: "Google: google_folder_iam_binding"
sidebar_current: "docs-google-folder-iam-binding"
description: |-
 Allows management of a single binding with an IAM policy for a Google Cloud Platform folder.
---

# google\_folder\_iam\_binding

Allows creation and management of a single binding within IAM policy for
an existing Google Cloud Platform folder.

~> **Note:** This resource _must not_ be used in conjunction with
   `google_folder_iam_policy` or they will fight over what your policy
   should be.

~> **Note:** On create, this resource will overwrite members of any existing roles.
    Use `terraform import` and inspect the `terraform plan` output to ensure
    your existing members are preserved.

## Example Usage

```hcl
resource "google_folder" "department1" {
  display_name = "Department 1"
  parent       = "organizations/1234567"
}

resource "google_folder_iam_binding" "admin" {
  folder = google_folder.department1.name
  role   = "roles/editor"

  members = [
    "user:alice@gmail.com",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `folder` - (Required) The resource name of the folder the policy is attached to. Its format is folders/{folder_id}.

* `members` (Required) - An array of identities that will be granted the privilege in the `role`.
  Each entry can have one of the following values:
  * **user:{emailid}**: An email address that is associated with a specific Google account. For example, alice@gmail.com.
  * **serviceAccount:{emailid}**: An email address that represents a service account. For example, my-other-app@appspot.gserviceaccount.com.
  * **group:{emailid}**: An email address that represents a Google group. For example, admins@example.com.
  * **domain:{domain}**: A G Suite domain (primary, instead of alias) name that represents all the users of that domain. For example, google.com or example.com.
  * For more details on format and restrictions see https://cloud.google.com/billing/reference/rest/v1/Policy#Binding

* `role` - (Required) The role that should be applied. Only one
    `google_folder_iam_binding` can be used per role. Note that custom roles must be of the format
    `[projects|organizations]/{parent-name}/roles/{role-name}`.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `etag` - (Computed) The etag of the folder's IAM policy.

## Import

IAM binding imports use space-delimited identifiers; first the resource in question and then the role.  These bindings can be imported using the `folder` and role, e.g.

```
$ terraform import google_folder_iam_binding.viewer "folder-name roles/viewer"
```
