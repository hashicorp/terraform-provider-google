---
layout: "google"
page_title: "Google: google_project_iam_binding"
sidebar_current: "docs-google-project-iam-binding"
description: |-
 Allows management of a single binding with an IAM policy for a Google Cloud Platform project.
---

# google\_project\_iam\_binding

Allows creation and management of a single binding within IAM policy for
an existing Google Cloud Platform project.

~> **Note:** This resource _must not_ be used in conjunction with
   `google_project_iam_policy` or they will fight over what your policy
   should be.

## Example Usage

```hcl
resource "google_project_iam_binding" "project" {
  project = "your-project-id"
  role    = "roles/editor"

  members = [
    "user:jane@example.com",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `members` (Required) - An array of identites that will be granted the privilege in the `role`.
  Each entry can have one of the following values:
  * **user:{emailid}**: An email address that represents a specific Google account. For example, alice@gmail.com or joe@example.com.
  * **serviceAccount:{emailid}**: An email address that represents a service account. For example, my-other-app@appspot.gserviceaccount.com.
  * **group:{emailid}**: An email address that represents a Google group. For example, admins@example.com.
  * **domain:{domain}**: A Google Apps domain name that represents all the users of that domain. For example, google.com or example.com.

* `role` - (Required) The role that should be applied. Only one
    `google_project_iam_binding` can be used per role. Note that custom roles must be of the format
    `[projects|organizations]/{parent-name}/roles/{role-name}`.

* `project` - (Optional) The project ID. If not specified, uses the
    ID of the project configured with the provider.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `etag` - (Computed) The etag of the project's IAM policy.

## Import

IAM binding imports use space-delimited identifiers; first the resource in question and then the role.  These bindings can be imported using the `project_id` and role, e.g.

```
$ terraform import google_project_iam_binding.my_project "your-project-id roles/viewer"
```
