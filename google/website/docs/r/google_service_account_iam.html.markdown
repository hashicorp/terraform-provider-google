---
subcategory: "Cloud Platform"
description: |-
 Collection of resources to manage IAM policy for a service account.
---

# IAM policy for Service Account

When managing IAM roles, you can treat a service account either as a resource or as an identity. This resource is to add iam policy bindings to a service account resource, such as allowing the members to run operations as or modify the service account. To configure permissions for a service account on other GCP resources, use the [google_project_iam](google_project_iam.html) set of resources.

Three different resources help you manage your IAM policy for a service account. Each of these resources serves a different use case:

* `google_service_account_iam_policy`: Authoritative. Sets the IAM policy for the service account and replaces any existing policy already attached.
* `google_service_account_iam_binding`: Authoritative for a given role. Updates the IAM policy to grant a role to a list of members. Other roles within the IAM policy for the service account are preserved.
* `google_service_account_iam_member`: Non-authoritative. Updates the IAM policy to grant a role to a new member. Other members for the role for the service account are preserved.

~> **Note:** `google_service_account_iam_policy` **cannot** be used in conjunction with `google_service_account_iam_binding` and `google_service_account_iam_member` or they will fight over what your policy should be.

~> **Note:** `google_service_account_iam_binding` resources **can be** used in conjunction with `google_service_account_iam_member` resources **only if** they do not grant privilege to the same role.

## google\_service\_account\_iam\_policy

```hcl
data "google_iam_policy" "admin" {
  binding {
    role = "roles/iam.serviceAccountUser"

    members = [
      "user:jane@example.com",
    ]
  }
}

resource "google_service_account" "sa" {
  account_id   = "my-service-account"
  display_name = "A service account that only Jane can interact with"
}

resource "google_service_account_iam_policy" "admin-account-iam" {
  service_account_id = google_service_account.sa.name
  policy_data        = data.google_iam_policy.admin.policy_data
}
```

## google\_service\_account\_iam\_binding

```hcl
resource "google_service_account" "sa" {
  account_id   = "my-service-account"
  display_name = "A service account that only Jane can use"
}

resource "google_service_account_iam_binding" "admin-account-iam" {
  service_account_id = google_service_account.sa.name
  role               = "roles/iam.serviceAccountUser"

  members = [
    "user:jane@example.com",
  ]
}
```

With IAM Conditions:

```hcl
resource "google_service_account" "sa" {
  account_id   = "my-service-account"
  display_name = "A service account that only Jane can use"
}

resource "google_service_account_iam_binding" "admin-account-iam" {
  service_account_id = google_service_account.sa.name
  role               = "roles/iam.serviceAccountUser"

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

## google\_service\_account\_iam\_member

```hcl
data "google_compute_default_service_account" "default" {
}

resource "google_service_account" "sa" {
  account_id   = "my-service-account"
  display_name = "A service account that Jane can use"
}

resource "google_service_account_iam_member" "admin-account-iam" {
  service_account_id = google_service_account.sa.name
  role               = "roles/iam.serviceAccountUser"
  member             = "user:jane@example.com"
}

# Allow SA service account use the default GCE account
resource "google_service_account_iam_member" "gce-default-account-iam" {
  service_account_id = data.google_compute_default_service_account.default.name
  role               = "roles/iam.serviceAccountUser"
  member             = "serviceAccount:${google_service_account.sa.email}"
}
```

With IAM Conditions:

```hcl
resource "google_service_account" "sa" {
  account_id   = "my-service-account"
  display_name = "A service account that Jane can use"
}

resource "google_service_account_iam_member" "admin-account-iam" {
  service_account_id = google_service_account.sa.name
  role               = "roles/iam.serviceAccountUser"
  member             = "user:jane@example.com"

  condition {
    title       = "expires_after_2019_12_31"
    description = "Expiring at midnight of 2019-12-31"
    expression  = "request.time < timestamp(\"2020-01-01T00:00:00Z\")"
  }
}
```

## Argument Reference

The following arguments are supported:

* `service_account_id` - (Required) The fully-qualified name of the service account to apply policy to.

* `member/members` - (Required) Identities that will be granted the privilege in `role`.
  Each entry can have one of the following values:
  * **allUsers**: A special identifier that represents anyone who is on the internet; with or without a Google account.
  * **allAuthenticatedUsers**: A special identifier that represents anyone who is authenticated with a Google account or a service account.
  * **user:{emailid}**: An email address that represents a specific Google account. For example, alice@gmail.com or joe@example.com.
  * **serviceAccount:{emailid}**: An email address that represents a service account. For example, my-other-app@appspot.gserviceaccount.com.
  * **group:{emailid}**: An email address that represents a Google group. For example, admins@example.com.
  * **domain:{domain}**: A G Suite domain (primary, instead of alias) name that represents all the users of that domain. For example, google.com or example.com.

* `role` - (Required) The role that should be applied. Only one
    `google_service_account_iam_binding` can be used per role. Note that custom roles must be of the format
    `[projects|organizations]/{parent-name}/roles/{role-name}`.

* `policy_data` - (Required only by `google_service_account_iam_policy`) The policy data generated by
  a `google_iam_policy` data source.

* `condition` - (Optional) An [IAM Condition](https://cloud.google.com/iam/docs/conditions-overview) for a given binding.
  Structure is [documented below](#nested_condition).

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

* `etag` - (Computed) The etag of the service account IAM policy.

## Import

-> **Custom Roles**: If you're importing a IAM resource with a custom role, make sure to use the
full name of the custom role, e.g. `[projects/my-project|organizations/my-org]/roles/my-custom-role`.

### Importing IAM members

IAM member imports use space-delimited identifiers that contains the `service_account_id`, `role`, and `member`. For example:

* `"projects/{{project_id}}/serviceAccounts/{{service_account_email}} roles/editor user:foo@example.com"`

An [`import` block](https://developer.hashicorp.com/terraform/language/import) (Terraform v1.5.0 and later) can be used to import IAM members:

```tf
import {
  id = "projects/{{project_id}}/serviceAccounts/{{service_account_email}} roles/editor user:foo@example.com"
  to = google_service_account_iam_member.default
}
```

The [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import) can also be used:

```
$ terraform import google_service_account_iam_member.default "projects/{{project_id}}/serviceAccounts/{{service_account_email}} roles/editor user:foo@example.com"
```

### Importing IAM bindings

IAM binding imports use space-delimited identifiers that contains the `service_account_id` and `role`. For example:

* `"projects/{{project_id}}/serviceAccounts/{{service_account_email}} roles/editor"`

An [`import` block](https://developer.hashicorp.com/terraform/language/import) (Terraform v1.5.0 and later) can be used to import IAM bindings:

```tf
import {
  id = "projects/{{project_id}}/serviceAccounts/{{service_account_email}} roles/editor"
  to = google_service_account_iam_binding.default
}
```

The [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import) can also be used:

```
$ terraform import google_service_account_iam_binding.default "projects/{{project_id}}/serviceAccounts/{{service_account_email}} roles/editor"
```

### Importing IAM policies

IAM policy imports use the identifier of the Service Account resource in question. For example:

* `"projects/{{project_id}}/serviceAccounts/{{service_account_email}}"`

An [`import` block](https://developer.hashicorp.com/terraform/language/import) (Terraform v1.5.0 and later) can be used to import IAM policies:

```tf
import {
  id = "projects/{{project_id}}/serviceAccounts/{{service_account_email}}"
  to = google_service_account_iam_policy.default
}
```

The [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import) can also be used:

```
$ terraform import google_service_account_iam_policy.default projects/{{project_id}}/serviceAccounts/{{service_account_email}}
```

### Importing with conditions:

Here are examples of importing IAM memberships and bindings that include conditions:

```
$ terraform import google_service_account_iam_binding.admin-account-iam "projects/{your-project-id}/serviceAccounts/{your-service-account-email} roles/iam.serviceAccountUser expires_after_2019_12_31"

$ terraform import google_service_account_iam_member.admin-account-iam "projects/{your-project-id}/serviceAccounts/{your-service-account-email} roles/iam.serviceAccountUser user:foo@example.com expires_after_2019_12_31"
```