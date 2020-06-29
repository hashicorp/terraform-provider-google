---
subcategory: "Cloud Platform"
layout: "google"
page_title: "Google: google_billing_account"
sidebar_current: "docs-google-datasource-billing-account"
description: |-
  Get information about a Google Billing Account.
---

# google\_billing\_account

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

~> **NOTE:** One of `billing_account` or `display_name` must be specified.

## Attributes Reference

The following additional attributes are exported:

* `id` - The billing account ID.
* `name` - The resource name of the billing account in the form `billingAccounts/{billing_account_id}`.
* `project_ids` - The IDs of any projects associated with the billing account.
