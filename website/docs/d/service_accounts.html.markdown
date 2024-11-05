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

Example service accounts.

```hcl
data "google_service_accounts" "example" {
  project = "example-project"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Optional) The ID of the project. If it is not provided, the provider project is used.


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
