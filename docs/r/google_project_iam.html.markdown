---
layout: "google"
page_title: "Google: google_project_iam"
sidebar_current: "docs-google-project-iam-x"
description: |-
 Collection of resources to manage IAM policy for a project.
---

# IAM policy for projects

Three different resources help you manage your IAM policy for a project. Each of these resources serves a different use case:

* `google_project_iam_policy`: Authoritative. Sets the IAM policy for the project and replaces any existing policy already attached.
* `google_project_iam_binding`: Authoritative for a given role. Updates the IAM policy to grant a role to a list of members. Other roles within the IAM policy for the project are preserved.
* `google_project_iam_member`: Non-authoritative. Updates the IAM policy to grant a role to a new member. Other members for the role for the project are preserved.

~> **Note:** `google_project_iam_policy` **cannot** be used in conjunction with `google_project_iam_binding` and `google_project_iam_member` or they will fight over what your policy should be.

~> **Note:** `google_project_iam_binding` resources **can be** used in conjunction with `google_project_iam_member` resources **only if** they do not grant privilege to the same role.

## google\_project\_iam\_policy

~> **Be careful!** You can accidentally lock yourself out of your project
   using this resource. Proceed with caution.

```hcl
resource "google_project_iam_policy" "project" {
  project     = "your-project-id"
  policy_data = "${data.google_iam_policy.admin.policy_data}"
}

data "google_iam_policy" "admin" {
  binding {
    role = "roles/editor"

    members = [
      "user:jane@example.com",
    ]
  }
}
```

## google\_project\_iam\_binding

~> **Note:** If `role` is set to `roles/owner` and you don't specify a user or service account you have access to in `members`, you can lock yourself out of your project.

```hcl
resource "google_project_iam_binding" "project" {
  project = "your-project-id"
  role    = "roles/editor"

  members = [
    "user:jane@example.com",
  ]
}
```

## google\_project\_iam\_member

```hcl
resource "google_project_iam_member" "project" {
  project = "your-project-id"
  role    = "roles/editor"
  member  = "user:jane@example.com"
}
```

## Argument Reference

The following arguments are supported:

* `member/members` - (Required) Identities that will be granted the privilege in `role`.
  Each entry can have one of the following values:
  * **user:{emailid}**: An email address that represents a specific Google account. For example, alice@gmail.com or joe@example.com.
  * **serviceAccount:{emailid}**: An email address that represents a service account. For example, my-other-app@appspot.gserviceaccount.com.
  * **group:{emailid}**: An email address that represents a Google group. For example, admins@example.com.
  * **domain:{domain}**: A G Suite domain (primary, instead of alias) name that represents all the users of that domain. For example, google.com or example.com.

* `role` - (Required) The role that should be applied. Only one
    `google_project_iam_binding` can be used per role. Note that custom roles must be of the format
    `[projects|organizations]/{parent-name}/roles/{role-name}`.

* `policy_data` - (Required only by `google_project_iam_policy`) The `google_iam_policy` data source that represents
    the IAM policy that will be applied to the project. The policy will be
    merged with any existing policy applied to the project.

    Changing this updates the policy.

    Deleting this removes the policy, but leaves the original project policy
    intact. If there are overlapping `binding` entries between the original
    project policy and the data source policy, they will be removed.

* `project` - (Optional) The project ID. If not specified, uses the
    ID of the project configured with the provider.

* `authoritative` - (DEPRECATED) (Optional, only for `google_project_iam_policy`)
    A boolean value indicating if this policy
    should overwrite any existing IAM policy on the project. When set to true,
    **any policies not in your config file will be removed**. This can **lock
    you out** of your project until an Organization Administrator grants you
    access again, so please exercise caution. If this argument is `true` and you
    want to delete the resource, you must set the `disable_project` argument to
    `true`, acknowledging that the project will be inaccessible to anyone but the
    Organization Admins, as it will no longer have an IAM policy. Rather than using
    this, you should use `google_project_iam_binding` and
    `google_project_iam_member`.

* `disable_project` - (DEPRECATED) (Optional, only for `google_project_iam_policy`)
    A boolean value that must be set to `true`
    if you want to delete a `google_project_iam_policy` that is authoritative.
    
## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `etag` - (Computed) The etag of the project's IAM policy.

* `restore_policy` - (DEPRECATED) (Computed, only for `google_project_iam_policy`)
    The IAM policy that will be restored when a
    non-authoritative policy resource is deleted.

## Import

IAM resources can be imported using the `project_id`, role, and account.

```
$ terraform import google_project_iam_policy.my_project your-project-id

$ terraform import google_project_iam_binding.my_project "your-project-id roles/viewer"

$ terraform import google_project_iam_member.my_project "your-project-id roles/viewer foo@example.com"
```
