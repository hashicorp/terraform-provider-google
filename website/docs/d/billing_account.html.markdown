---
subcategory: "Cloud Billing"
description: |-
  Get information about a Google Billing Account.
---

# google_billing_account

Use this data source to get information about a Google Billing Account.

```hcl
data "google_billing_account" "acct" {
  display_name = "My Billing Account"
  open         = true
}

resource "google_project" "my_project" {
  name       = "My Project"
  project_id = "your-project-id"
  org_id     = "1234567"

  billing_account = data.google_billing_account.acct.id
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available billing accounts.
The given filters must match exactly one billing account whose data will be exported as attributes.
The following arguments are supported:

* `billing_account` (Optional) - The name of the billing account in the form `{billing_account_id}` or `billingAccounts/{billing_account_id}`.
* `display_name` (Optional) - The display name of the billing account.
* `open` (Optional) - `true` if the billing account is open, `false` if the billing account is closed.
* `lookup_projects` (Optional) - `true` if projects associated with the billing account should be read, `false` if this step
should be skipped. Setting `false` may be useful if the user permissions do not allow listing projects. Defaults to `true`.

~> **NOTE:** One of `billing_account` or `display_name` must be specified.

## Attributes Reference

The following additional attributes are exported:

* `id` - The billing account ID.
* `name` - The resource name of the billing account in the form `billingAccounts/{billing_account_id}`.
* `project_ids` - The IDs of any projects associated with the billing account. `lookup_projects` must not be false
for this to be populated.
