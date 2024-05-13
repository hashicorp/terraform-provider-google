---
subcategory: "Cloud Platform"
description: |-
 Collection of resources to manage IAM policy for a folder.
---

# IAM policy for folders

Four different resources help you manage your IAM policy for a folder. Each of these resources serves a different use case:

* `google_folder_iam_policy`: Authoritative. Sets the IAM policy for the folder and replaces any existing policy already attached.
* `google_folder_iam_binding`: Authoritative for a given role. Updates the IAM policy to grant a role to a list of members. Other roles within the IAM policy for the folder are preserved.
* `google_folder_iam_member`: Non-authoritative. Updates the IAM policy to grant a role to a new member. Other members for the role for the folder are preserved.
* `google_folder_iam_audit_config`: Authoritative for a given service. Updates the IAM policy to enable audit logging for the given service.


~> **Note:** `google_folder_iam_policy` **cannot** be used in conjunction with `google_folder_iam_binding`, `google_folder_iam_member`, or `google_folder_iam_audit_config` or they will fight over what your policy should be.

~> **Note:** `google_folder_iam_binding` resources **can be** used in conjunction with `google_folder_iam_member` resources **only if** they do not grant privilege to the same role.

~> **Note:** The underlying API method `projects.setIamPolicy` has constraints which are documented [here](https://cloud.google.com/resource-manager/reference/rest/v1/projects/setIamPolicy). In addition to these constraints, 
   IAM Conditions cannot be used with Basic Roles such as Owner. Violating these constraints will result in the API returning a 400 error code so please review these if you encounter errors with this resource.

## google_folder_iam_policy

!> **Be careful!** You can accidentally lock yourself out of your folder
   using this resource. Deleting a `google_folder_iam_policy` removes access
   from anyone without permissions on its parent folder/organization. Proceed with caution.
   It's not recommended to use `google_folder_iam_policy` with your provider folder
   to avoid locking yourself out, and it should generally only be used with folders
   fully managed by Terraform. If you do use this resource, it is recommended to **import** the policy before
   applying the change.

```hcl
resource "google_folder_iam_policy" "folder" {
  folder      = "folders/1234567"
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
resource "google_folder_iam_policy" "folder" {
  folder      = "folders/1234567"
  policy_data = "${data.google_iam_policy.admin.policy_data}"
}

data "google_iam_policy" "admin" {
  binding {
    role = "roles/compute.admin"

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

## google_folder_iam_binding

```hcl
resource "google_folder_iam_binding" "folder" {
  folder  = "folders/1234567"
  role    = "roles/editor"

  members = [
    "user:jane@example.com",
  ]
}
```

With IAM Conditions:

```hcl
resource "google_folder_iam_binding" "folder" {
  folder  = "folders/1234567"
  role    = "roles/container.admin"

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

## google_folder_iam_member

```hcl
resource "google_folder_iam_member" "folder" {
  folder  = "folders/1234567"
  role    = "roles/editor"
  member  = "user:jane@example.com"
}
```

With IAM Conditions:

```hcl
resource "google_folder_iam_member" "folder" {
  folder  = "folders/1234567"
  role    = "roles/firebase.admin"
  member  = "user:jane@example.com"

  condition {
    title       = "expires_after_2019_12_31"
    description = "Expiring at midnight of 2019-12-31"
    expression  = "request.time < timestamp(\"2020-01-01T00:00:00Z\")"
  }
}
```

## google_folder_iam_audit_config

```hcl
resource "google_folder_iam_audit_config" "folder" {
  folder  = "folders/1234567"
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

* `member/members` - (Required except for google_folder_iam_audit_config) Identities that will be granted the privilege in `role`.
  Each entry can have one of the following values:
  * **user:{emailid}**: An email address that represents a specific Google account. For example, alice@gmail.com or joe@example.com.
  * **serviceAccount:{emailid}**: An email address that represents a service account. For example, my-other-app@appspot.gserviceaccount.com.
  * **group:{emailid}**: An email address that represents a Google group. For example, admins@example.com.
  * **domain:{domain}**: A G Suite domain (primary, instead of alias) name that represents all the users of that domain. For example, google.com or example.com.

* `role` - (Required except for google_folder_iam_audit_config) The role that should be applied. Only one
    `google_folder_iam_binding` can be used per role. Note that custom roles must be of the format
    `organizations/{{org_id}}/roles/{{role_id}}`.

* `policy_data` - (Required only by `google_folder_iam_policy`) The `google_iam_policy` data source that represents
    the IAM policy that will be applied to the folder. The policy will be
    merged with any existing policy applied to the folder.

    Changing this updates the policy.

    Deleting this removes all policies from the folder, locking out users without
    folder-level access.

* `folder` - (Required) The resource name of the folder the policy is attached to. Its format is folders/{folder_id}.

* `service` - (Required only by google_folder_iam_audit_config) Service which will be enabled for audit logging.  The special value `allServices` covers all services.  Note that if there are google_folder_iam_audit_config resources covering both `allServices` and a specific service then the union of the two AuditConfigs is used for that service: the `log_types` specified in each `audit_log_config` are enabled, and the `exempted_members` in each `audit_log_config` are exempted.

* `audit_log_config` - (Required only by google_folder_iam_audit_config) The configuration for logging of each type of permission.  This can be specified multiple times.  Structure is [documented below](#nested_audit_log_config).

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

* `etag` - (Computed) The etag of the folder's IAM policy.


## Import

-> **Custom Roles**: If you're importing a IAM resource with a custom role, make sure to use the
 full name of the custom role, e.g. `organizations/{{org_id}}/roles/{{role_id}}`.
 
-> **Conditional IAM Bindings**: If you're importing a IAM binding with a condition block, make sure
 to include the title of condition, e.g. `terraform import google_folder_iam_binding.my_folder "folder roles/{{role_id}} condition-title"`
 

### Importing IAM members

IAM member imports use space-delimited identifiers that contain the resource's `folder_id`, `role`, and `member` e.g.

* `"folders/{{folder_id}} roles/viewer user:foo@example.com"`

An [`import` block](https://developer.hashicorp.com/terraform/language/import) (Terraform v1.5.0 and later) can be used to import IAM members:

```tf
import {
  id = "folders/{{folder_id}} roles/viewer user:foo@example.com"
  to = google_folder_iam_member.default
}
```

The [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import) can also be used:

```
$ terraform import google_folder_iam_member.default "folder roles/viewer user:foo@example.com"
```

### Importing IAM bindings

IAM binding imports use space-delimited identifiers that contain the resource's `folder_id` and `role`, e.g.

* `"folders/{{folder_id}} roles/viewer"`

An [`import` block](https://developer.hashicorp.com/terraform/language/import) (Terraform v1.5.0 and later) can be used to import IAM bindings:

```tf
import {
  id = "folders/{{folder_id}} roles/viewer"
  to = google_folder_iam_binding.default
}
```

The [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import) can also be used:

```
$ terraform import google_folder_iam_binding.default "folder roles/viewer"
```

### Importing IAM policies

IAM policy imports use the identifier of the Folder only. For example:

* `folders/{{folder_id}}`

An [`import` block](https://developer.hashicorp.com/terraform/language/import) (Terraform v1.5.0 and later) can be used to import IAM policies:

```tf
import {
  id = "folders/{{folder_id}}"
  to = google_folder_iam_policy.default
}
```

The [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import) can also be used:

```
$ terraform import google_folder_iam_policy.default folders/{{folder_id}}
```


### Importing Audit Configs

An audit config can be imported into a `google_folder_iam_audit_config` resource using the resource's `folder_id` and the `service`, e.g:

* `"folder/{{folder_id}} foo.googleapis.com"`


An [`import` block](https://developer.hashicorp.com/terraform/language/import) (Terraform v1.5.0 and later) can be used to import audit configs:

```tf
import {
  id = "folder/{{folder_id}} foo.googleapis.com"
  to = google_folder_iam_audit_config.default
}
```

The [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import) can also be used:

```
terraform import google_folder_iam_audit_config.default "folder/{{folder_id}} foo.googleapis.com"
```