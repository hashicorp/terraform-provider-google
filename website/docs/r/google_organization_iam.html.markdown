---
subcategory: "Cloud Platform"
description: |-
 Collection of resources to manage IAM policy for a organization.
---

# IAM policy for organizations

Four different resources help you manage your IAM policy for a organization. Each of these resources serves a different use case:

* `google_organization_iam_policy`: Authoritative. Sets the IAM policy for the organization and replaces any existing policy already attached.
* `google_organization_iam_binding`: Authoritative for a given role. Updates the IAM policy to grant a role to a list of members. Other roles within the IAM policy for the organization are preserved.
* `google_organization_iam_member`: Non-authoritative. Updates the IAM policy to grant a role to a new member. Other members for the role for the organization are preserved.
* `google_organization_iam_audit_config`: Authoritative for a given service. Updates the IAM policy to enable audit logging for the given service.


~> **Note:** `google_organization_iam_policy` **cannot** be used in conjunction with `google_organization_iam_binding`, `google_organization_iam_member`, or `google_organization_iam_audit_config` or they will fight over what your policy should be.

~> **Note:** `google_organization_iam_binding` resources **can be** used in conjunction with `google_organization_iam_member` resources **only if** they do not grant privilege to the same role.

## google_organization_iam_policy

!> **Warning:** New organizations have several default policies which will,
   without extreme caution, be **overwritten** by use of this resource.
   The safest alternative is to use multiple `google_organization_iam_binding`
   resources. This resource makes it easy to remove your own access to
   an organization, which will require a call to Google Support to have
   fixed, and can take multiple days to resolve.
   <br /><br />
   In general, this resource should only be used with organizations
   fully managed by Terraform.If you do use this resource,
   the best way to be sure that you are not making dangerous changes is to start
   by **importing** your existing policy, and examining the diff very closely.

```hcl
resource "google_organization_iam_policy" "organization" {
  org_id      = "1234567890"
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

With IAM Conditions:

```hcl
resource "google_organization_iam_policy" "organization" {
  org_id      = "1234567890"
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

## google_organization_iam_binding

~> **Note:** If `role` is set to `roles/owner` and you don't specify a user or service account you have access to in `members`, you can lock yourself out of your organization.

```hcl
resource "google_organization_iam_binding" "organization" {
  org_id  = "1234567890"
  role    = "roles/editor"

  members = [
    "user:jane@example.com",
  ]
}
```

With IAM Conditions:

```hcl
resource "google_organization_iam_binding" "organization" {
  org_id  = "1234567890"
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

## google_organization_iam_member

```hcl
resource "google_organization_iam_member" "organization" {
  org_id  = "1234567890"
  role    = "roles/editor"
  member  = "user:jane@example.com"
}
```

With IAM Conditions:

```hcl
resource "google_organization_iam_member" "organization" {
  org_id  = "1234567890"
  role    = "roles/editor"
  member  = "user:jane@example.com"

  condition {
    title       = "expires_after_2019_12_31"
    description = "Expiring at midnight of 2019-12-31"
    expression  = "request.time < timestamp(\"2020-01-01T00:00:00Z\")"
  }
}
```

## google_organization_iam_audit_config

```hcl
resource "google_organization_iam_audit_config" "organization" {
  org_id = "1234567890"
  service = "allServices"
  audit_log_config {
    log_type = "ADMIN_READ"
  }
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

* `member/members` - (Required except for google_organization_iam_audit_config) Identities that will be granted the privilege in `role`.
  Each entry can have one of the following values:
  * **user:{emailid}**: An email address that represents a specific Google account. For example, alice@gmail.com or joe@example.com.
  * **serviceAccount:{emailid}**: An email address that represents a service account. For example, my-other-app@appspot.gserviceaccount.com.
  * **group:{emailid}**: An email address that represents a Google group. For example, admins@example.com.
  * **domain:{domain}**: A G Suite domain (primary, instead of alias) name that represents all the users of that domain. For example, google.com or example.com.

* `role` - (Required except for google_organization_iam_audit_config) The role that should be applied. Only one
    `google_organization_iam_binding` can be used per role. Note that custom roles must be of the format
    `organizations/{{org_id}}/roles/{{role_id}}`.

* `policy_data` - (Required only by `google_organization_iam_policy`) The `google_iam_policy` data source that represents
    the IAM policy that will be applied to the organization. The policy will be
    merged with any existing policy applied to the organization.

    Changing this updates the policy.

    Deleting this removes all policies from the organization, locking out users without
    organization-level access.

* `org_id` - (Required) The organization id of the target organization.

* `service` - (Required only by google_organization_iam_audit_config) Service which will be enabled for audit logging.  The special value `allServices` covers all services.  Note that if there are google_organization_iam_audit_config resources covering both `allServices` and a specific service then the union of the two AuditConfigs is used for that service: the `log_types` specified in each `audit_log_config` are enabled, and the `exempted_members` in each `audit_log_config` are exempted.

* `audit_log_config` - (Required only by google_organization_iam_audit_config) The configuration for logging of each type of permission.  This can be specified multiple times.  Structure is [documented below](#nested_audit_log_config).

* `condition` - (Optional) An [IAM Condition](https://cloud.google.com/iam/docs/conditions-overview) for a given binding.
  Structure is [documented below](#nested_condition).

---

<a name="nested_audit_log_config"></a>The `audit_log_config` block supports:

* `log_type` - (Required) Permission type for which logging is to be configured.  Must be one of `DATA_READ`, `DATA_WRITE`, or `ADMIN_READ`.

* `exempted_members` - (Optional) Identities that do not cause logging for this type of permission.  The format is the same as that for `members`.

<a name="nested_condition"></a>The `condition` block supports:

* `expression` - (Required) Textual representation of an expression in Common Expression Language syntax.

* `title` - (Required) A title for the expression, i.e. a short string describing its purpose.

* `description` - (Optional) An optional description of the expression. This is a longer text which describes the expression, e.g. when hovered over it in a UI.

~> **Warning:** Terraform considers the `role` and condition contents (`title`+`description`+`expression`) as the
  identifier for the binding. This means that if any part of the condition is changed out-of-band, Terraform will
  consider it to be an entirely different resource and will treat it as such.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `etag` - (Computed) The etag of the organization's IAM policy.


## Import

-> **Custom Roles**: If you're importing a IAM resource with a custom role, make sure to use the
 full name of the custom role, e.g. `organizations/{{org_id}}/roles/{{role_id}}`.

-> **Conditional IAM Bindings**: If you're importing a IAM binding with a condition block, make sure
 to include the title of condition, e.g. `terraform import google_organization_iam_binding.my_organization "your-org-id roles/{{role_id}} condition-title"`

### Importing IAM members

IAM member imports use space-delimited identifiers that contain the resource's `org_id`, `role`, and `member` e.g.

* `"{{org_id}} roles/viewer user:foo@example.com"`

An [`import` block](https://developer.hashicorp.com/terraform/language/import) (Terraform v1.5.0 and later) can be used to import IAM members:

```tf
import {
  id = "{{org_id}} roles/viewer user:foo@example.com"
  to = google_organization_iam_member.default
}
```

The [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import) can also be used:

```
$ terraform import google_organization_iam_member.default "{{org_id}} roles/viewer user:foo@example.com"
```

### Importing IAM bindings

IAM binding imports use space-delimited identifiers that contain the `org_id` and role, e.g.

* `"{{org_id}} roles/viewer"`

An [`import` block](https://developer.hashicorp.com/terraform/language/import) (Terraform v1.5.0 and later) can be used to import IAM bindings:

```tf
import {
  id = "{{org_id}} roles/viewer"
  to = google_organization_iam_binding.default
}
```

The [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import) can also be used:

```
terraform import google_organization_iam_binding.default "{{org_id}} roles/viewer"
```

### Importing IAM policies

IAM policy imports use the identifier of the Organization only. For example:

* `"{{org_id}}"`

An [`import` block](https://developer.hashicorp.com/terraform/language/import) (Terraform v1.5.0 and later) can be used to import IAM policies:

```tf
import {
  id = "{{org_id}}"
  to = google_organization_iam_policy.default
}
```

The [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import) can also be used:

```
$ terraform import google_organization_iam_policy.default "{{org_id}}"
```


### Importing Audit Configs

An audit config can be imported into a `google_organization_iam_audit_config` resource using the resource's `org_id` and the `service`, e.g:

* `"{{org_id}} foo.googleapis.com"`


An [`import` block](https://developer.hashicorp.com/terraform/language/import) (Terraform v1.5.0 and later) can be used to import audit configs:

```tf
import {
  id = "{{org_id}} foo.googleapis.com"
  to = google_organization_iam_audit_config.default
}
```

The [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import) can also be used:

```
terraform import google_organization_iam_audit_config.default "{{org_id}} foo.googleapis.com"
```