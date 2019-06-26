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
   using this resource. Deleting a `google_project_iam_policy` removes access
   from anyone without organization-level access to the project. Proceed with caution.
   It's not recommended to use `google_project_iam_policy` with your provider project
   to avoid locking yourself out, and it should generally only be used with projects
   fully managed by Terraform.

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

    Deleting this removes all policies from the project, locking out users without
    organization-level access.

* `project` - (Optional) The project ID. If not specified for `google_project_iam_binding`
or `google_project_iam_member`, uses the ID of the project configured with the provider.
Required for `google_project_iam_policy` - you must explicitly set the project, and it
will not be inferred from the provider.
    
## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `etag` - (Computed) The etag of the project's IAM policy.


## Import

IAM member imports use space-delimited identifiers; the resource in question, the role, and the account.  This member resource can be imported using the `project_id`, role, and member e.g.

```
$ terraform import google_project_iam_member.my_project "your-project-id roles/viewer user:foo@example.com"
```

IAM binding imports use space-delimited identifiers; the resource in question and the role.  This binding resource can be imported using the `project_id` and role, e.g.

```
terraform import google_project_iam_binding.my_project "your-project-id roles/viewer"
```

IAM policy imports use the identifier of the resource in question.  This policy resource can be imported using the `project_id`.

```
$ terraform import google_project_iam_policy.my_project your-project-id
```
