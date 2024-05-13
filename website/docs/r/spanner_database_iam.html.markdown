---
subcategory: "Cloud Spanner"
description: |-
 Collection of resources to manage IAM policy for a Spanner database.
---

# IAM policy for Spanner Databases

Three different resources help you manage your IAM policy for a Spanner database. Each of these resources serves a different use case:

* `google_spanner_database_iam_policy`: Authoritative. Sets the IAM policy for the database and replaces any existing policy already attached.

~> **Warning:** It's entirely possibly to lock yourself out of your database using `google_spanner_database_iam_policy`. Any permissions granted by default will be removed unless you include them in your config.

* `google_spanner_database_iam_binding`: Authoritative for a given role. Updates the IAM policy to grant a role to a list of members. Other roles within the IAM policy for the database are preserved.
* `google_spanner_database_iam_member`: Non-authoritative. Updates the IAM policy to grant a role to a new member. Other members for the role for the database are preserved.

~> **Note:** `google_spanner_database_iam_policy` **cannot** be used in conjunction with `google_spanner_database_iam_binding` and `google_spanner_database_iam_member` or they will fight over what your policy should be.

~> **Note:** `google_spanner_database_iam_binding` resources **can be** used in conjunction with `google_spanner_database_iam_member` resources **only if** they do not grant privilege to the same role.

## google_spanner_database_iam_policy

```hcl
data "google_iam_policy" "admin" {
  binding {
    role = "roles/editor"

    members = [
      "user:jane@example.com",
    ]
  }
}

resource "google_spanner_database_iam_policy" "database" {
  instance    = "your-instance-name"
  database    = "your-database-name"
  policy_data = data.google_iam_policy.admin.policy_data
}
```

With IAM Conditions:

```hcl
data "google_iam_policy" "admin" {
  binding {
    role = "roles/editor"

    members = [
      "user:jane@example.com",
    ]
    
    condition {
      title       = "My Role"
      description = "Grant permissions on my_role"
      expression  = "(resource.type == \"spanner.googleapis.com/DatabaseRole\" && (resource.name.endsWith(\"/myrole\")))"
    }
  }
}

resource "google_spanner_database_iam_policy" "database" {
  instance    = "your-instance-name"
  database    = "your-database-name"
  policy_data = data.google_iam_policy.admin.policy_data
}
```

## google_spanner_database_iam_binding

```hcl
resource "google_spanner_database_iam_binding" "database" {
  instance = "your-instance-name"
  database = "your-database-name"
  role     = "roles/compute.networkUser"

  members = [
    "user:jane@example.com",
  ]
}
```

With IAM Conditions:

```hcl
resource "google_spanner_database_iam_binding" "database" {
  instance = "your-instance-name"
  database = "your-database-name"
  role     = "roles/compute.networkUser"

  members = [
    "user:jane@example.com",
  ]
  
  condition {
    title       = "My Role"
    description = "Grant permissions on my_role"
    expression  = "(resource.type == \"spanner.googleapis.com/DatabaseRole\" && (resource.name.endsWith(\"/myrole\")))"
  }
}
```

## google_spanner_database_iam_member

```hcl
resource "google_spanner_database_iam_member" "database" {
  instance = "your-instance-name"
  database = "your-database-name"
  role     = "roles/compute.networkUser"
  member   = "user:jane@example.com"
}
```

With IAM Conditions:

```hcl
resource "google_spanner_database_iam_member" "database" {
  instance = "your-instance-name"
  database = "your-database-name"
  role     = "roles/compute.networkUser"
  member   = "user:jane@example.com"
  
  condition {
    title       = "My Role"
    description = "Grant permissions on my_role"
    expression  = "(resource.type == \"spanner.googleapis.com/DatabaseRole\" && (resource.name.endsWith(\"/myrole\")))"
  }
}
```

## Argument Reference

The following arguments are supported:

* `database` - (Required) The name of the Spanner database.

* `instance` - (Required) The name of the Spanner instance the database belongs to.

* `member/members` - (Required) Identities that will be granted the privilege in `role`.
  Each entry can have one of the following values:
  * **allUsers**: A special identifier that represents anyone who is on the internet; with or without a Google account.
  * **allAuthenticatedUsers**: A special identifier that represents anyone who is authenticated with a Google account or a service account.
  * **user:{emailid}**: An email address that represents a specific Google account. For example, alice@gmail.com or joe@example.com.
  * **serviceAccount:{emailid}**: An email address that represents a service account. For example, my-other-app@appspot.gserviceaccount.com.
  * **group:{emailid}**: An email address that represents a Google group. For example, admins@example.com.
  * **domain:{domain}**: A G Suite domain (primary, instead of alias) name that represents all the users of that domain. For example, google.com or example.com.

* `role` - (Required) The role that should be applied. Only one
    `google_spanner_database_iam_binding` can be used per role. Note that custom roles must be of the format
    `[projects|organizations]/{parent-name}/roles/{role-name}`.

* `policy_data` - (Required only by `google_spanner_database_iam_policy`) The policy data generated by
  a `google_iam_policy` data source.

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

* `condition` - (Optional) An [IAM Condition](https://cloud.google.com/iam/docs/conditions-overview) for a given binding.
  Structure is [documented below](#nested_condition).

---

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

* `etag` - (Computed) The etag of the database's IAM policy.

## Import

-> **Custom Roles:** If you're importing a IAM resource with a custom role, make sure to use the
 full name of the custom role, e.g. `[projects/my-project|organizations/my-org]/roles/my-custom-role`.

For all import syntaxes, the "resource in question" can take any of the following forms:

* {{project}}/{{instance}}/{{database}}
* {{instance}}/{{database}} (project is taken from provider project)

### Importing IAM members

IAM member imports use space-delimited identifiers that contains the `database`, `role`, and `member`. For example:

* `"{{project}}/{{instance}}/{{database}} roles/viewer user:foo@example.com"`

An [`import` block](https://developer.hashicorp.com/terraform/language/import) (Terraform v1.5.0 and later) can be used to import IAM members:

```tf
import {
  id = "{{project}}/{{instance}}/{{database}} roles/viewer user:foo@example.com"
  to = google_spanner_database_iam_member.default
}
```

The [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import) can also be used:

```
$ terraform import google_spanner_database_iam_member.default "{{project}}/{{instance}}/{{database}} roles/viewer user:foo@example.com"
```

### Importing IAM bindings

IAM binding imports use space-delimited identifiers that contain the resource's `database` and `role`. For example:

* `"{{project}}/{{instance}}/{{database}} roles/viewer"`

An [`import` block](https://developer.hashicorp.com/terraform/language/import) (Terraform v1.5.0 and later) can be used to import IAM bindings:

```tf
import {
  id = "{{project}}/{{instance}}/{{database}} roles/viewer"
  to = google_spanner_database_iam_binding.default
}
```

The [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import) can also be used:

```
$ terraform import google_spanner_database_iam_binding.default "{{project}}/{{instance}}/{{database}} roles/viewer"
```

### Importing IAM policies

IAM policy imports use the identifier of the Spanner Database resource in question. For example:

* `{{project}}/{{instance}}/{{database}}`


An [`import` block](https://developer.hashicorp.com/terraform/language/import) (Terraform v1.5.0 and later) can be used to import IAM policies:

```tf
import {
  id = {{project}}/{{instance}}/{{database}}
  to = google_spanner_database_iam_policy.default
}
```

The [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import) can also be used:

```
$ terraform import google_spanner_database_iam_policy.default {{project}}/{{instance}}/{{database}}
```