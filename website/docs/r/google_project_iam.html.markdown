---
subcategory: "Cloud Platform"
layout: "google"
page_title: "Google: google_project_iam"
sidebar_current: "docs-google-project-iam-x"
description: |-
 Collection of resources to manage IAM policy for a project.
---

# IAM policy for projects

Four different resources help you manage your IAM policy for a project. Each of these resources serves a different use case:

* `google_project_iam_policy`: Authoritative. Sets the IAM policy for the project and replaces any existing policy already attached.
* `google_project_iam_binding`: Authoritative for a given role. Updates the IAM policy to grant a role to a list of members. Other roles within the IAM policy for the project are preserved.
* `google_project_iam_member`: Non-authoritative. Updates the IAM policy to grant a role to a new member. Other members for the role for the project are preserved.
* `google_project_iam_audit_config`: Authoritative for a given service. Updates the IAM policy to enable audit logging for the given service.


~> **Note:** `google_project_iam_policy` **cannot** be used in conjunction with `google_project_iam_binding`, `google_project_iam_member`, or `google_project_iam_audit_config` or they will fight over what your policy should be.

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
  policy_data = data.google_iam_policy.admin.policy_data
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

With IAM Conditions ([beta](https://terraform.io/docs/providers/google/provider_versions.html), Whitelist-only):

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

    condition {
      title       = "expires_after_2019_12_31"
      description = "Expiring at midnight of 2019-12-31"
      expression  = "request.time < timestamp(\"2020-01-01T00:00:00Z\")"
    }
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

With IAM Conditions ([beta](https://terraform.io/docs/providers/google/provider_versions.html), Whitelist-only):

```hcl
resource "google_project_iam_binding" "project" {
  project = "your-project-id"
  role    = "roles/editor"

  members = [
    "user:jane@example.com",
  ]

  condition {
    title       = "expires_after_2019_12_31"
    description = "Expiring at midnight of 2019-12-31"
    expression  = "request.time < timestamp(\"2020-01-01T00:00:00Z\")"
  }
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

With IAM Conditions ([beta](https://terraform.io/docs/providers/google/provider_versions.html), Whitelist-only):

```hcl
resource "google_project_iam_member" "project" {
  project = "your-project-id"
  role    = "roles/editor"
  member  = "user:jane@example.com"

  condition {
    title       = "expires_after_2019_12_31"
    description = "Expiring at midnight of 2019-12-31"
    expression  = "request.time < timestamp(\"2020-01-01T00:00:00Z\")"
  }
}
```

## google\_project\_iam\_audit\_config

```hcl
resource "google_project_iam_audit_config" "project" {
  project = "your-project-id"
  service = "allServices"
  audit_log_config {
    log_type = "DATA_READ"
    exempted_members = [
      "user:joebloggs@hashicorp.com",
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `member/members` - (Required except for google\_project\_iam\_audit\_config) Identities that will be granted the privilege in `role`.
  Each entry can have one of the following values:
  * **user:{emailid}**: An email address that represents a specific Google account. For example, alice@gmail.com or joe@example.com.
  * **serviceAccount:{emailid}**: An email address that represents a service account. For example, my-other-app@appspot.gserviceaccount.com.
  * **group:{emailid}**: An email address that represents a Google group. For example, admins@example.com.
  * **domain:{domain}**: A G Suite domain (primary, instead of alias) name that represents all the users of that domain. For example, google.com or example.com.

* `role` - (Required except for google\_project\_iam\_audit\_config) The role that should be applied. Only one
    `google_project_iam_binding` can be used per role. Note that custom roles must be of the format
    `[projects|organizations]/{parent-name}/roles/{role-name}`.

* `policy_data` - (Required only by `google_project_iam_policy`) The `google_iam_policy` data source that represents
    the IAM policy that will be applied to the project. The policy will be
    merged with any existing policy applied to the project.

    Changing this updates the policy.

    Deleting this removes all policies from the project, locking out users without
    organization-level access.

* `project` - (Optional) The project ID. If not specified for `google_project_iam_binding`, `google_project_iam_member`, or `google_project_iam_audit_config`, uses the ID of the project configured with the provider.
Required for `google_project_iam_policy` - you must explicitly set the project, and it
will not be inferred from the provider.

* `service` - (Required only by google\_project\_iam\_audit\_config) Service which will be enabled for audit logging.  The special value `allServices` covers all services.  Note that if there are google\_project\_iam\_audit\_config resources covering both `allServices` and a specific service then the union of the two AuditConfigs is used for that service: the `log_types` specified in each `audit_log_config` are enabled, and the `exempted_members` in each `audit_log_config` are exempted.

* `audit_log_config` - (Required only by google\_project\_iam\_audit\_config) The configuration for logging of each type of permission.  This can be specified multiple times.  Structure is documented below.

* `condition` - (Optional, [Beta](https://terraform.io/docs/providers/google/provider_versions.html)) An [IAM Condition](https://cloud.google.com/iam/docs/conditions-overview) for a given binding. You must be whitelisted for the IAM Conditions private beta in order to use them in Terraform.
  Structure is documented below.

---

The `audit_log_config` block supports:

* `log_type` - (Required) Permission type for which logging is to be configured.  Must be one of `DATA_READ`, `DATA_WRITE`, or `ADMIN_READ`.

* `exempted_members` - (Optional) Identities that do not cause logging for this type of permission.  The format is the same as that for `members`.

The `condition` block supports:

* `expression` - (Required) Textual representation of an expression in Common Expression Language syntax.

* `title` - (Required) A title for the expression, i.e. a short string describing its purpose.

* `description` - (Optional) An optional description of the expression. This is a longer text which describes the expression, e.g. when hovered over it in a UI.

~> **Warning:** Terraform considers the `role` and condition contents (`title`+`description`+`expression`) as the
  identifier for the binding. This means that if any part of the condition is changed out-of-band, Terraform will
  consider it to be an entirely different resource and will treat it as such.

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

IAM audit config imports use the identifier of the resource in question and the service, e.g.

```
terraform import google_project_iam_audit_config.my_project "your-project-id foo.googleapis.com"
```
