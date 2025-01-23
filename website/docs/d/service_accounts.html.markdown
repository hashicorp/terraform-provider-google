---
subcategory: "Cloud Platform"
description: |-
  Get the service accounts from a project.
---


# google_service_accounts

Gets a list of all service accounts from a project.
See [the official documentation](https://cloud.google.com/iam/docs/service-account-overview)
and [API](https://cloud.google.com/iam/docs/reference/rest/v1/projects.serviceAccounts).

## Example Usage

Get all service accounts from a project

```hcl
data "google_service_accounts" "example" {
  project = "example-project"
}
```

Get all service accounts that are prefixed with `"foo"`

```hcl
data "google_service_accounts" "foo" {
  prefix  = "foo"
}
```

Get all service accounts that contain `"bar"`

```hcl
data "google_service_accounts" "bar" {
  regex   = ".*bar.*"
}
```

Get all service accounts that are prefixed with `"foo"` and contain `"bar"`

```hcl
data "google_service_accounts" "foo_bar" {
  prefix  = "foo"
  regex   = ".*bar.*"
}
```

## Argument Reference

The following arguments are supported:

* `prefix` - (Optional) A prefix for filtering. It's applied with the `account_id`.

* `project` - (Optional) The ID of the project. If it is not provided, the provider project is used.

* `regex` - (Optional) A regular expression for filtering. It's applied with the `email`. Further information about the syntax can be found [here](https://github.com/google/re2/wiki/Syntax).

## Attributes Reference

The following attributes are exported:

* `accounts` - A list of all retrieved service accounts. Structure is [defined below](#nested_accounts).

<a name="nested_accounts"></a>The `accounts` block supports:

* `account_id` - The Google service account ID (the part before the `@` sign in the `email`)

* `disabled` - Whether a service account is disabled or not.

* `display_name` - The display name for the service account.

* `email` - The e-mail address of the service account. This value
    should be referenced from any `google_iam_policy` data sources
    that would grant the service account privileges.

* `member` - The Identity of the service account in the form `serviceAccount:{email}`. This value is often used to refer to the service account in order to grant IAM permissions.

* `name` - The fully-qualified name of the service account.

* `unique_id` - The unique id of the service account.
