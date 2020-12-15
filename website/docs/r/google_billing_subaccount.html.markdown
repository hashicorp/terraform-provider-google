---
subcategory: "Cloud Platform"
layout: "google"
page_title: "Google: google_billing_subaccount"
sidebar_current: "docs-google-billing-subaccount"
description: |-
 Allows management of a Google Cloud Billing Subaccount.
---

# google\_billing\_subaccount

Allows creation and management of a Google Cloud Billing Subaccount.

!> **WARNING:** Deleting this Terraform resource will not delete or close the billing subaccount.

```hcl
resource "google_billing_subaccount" "subaccount" {
    display_name = "My Billing Account"
    master_billing_account = "012345-567890-ABCDEF"
}
```

## Argument Reference

* `display_name` (Required) - The display name of the billing account.

* `master_billing_account` (Required) - The name of the master billing account that the subaccount
  will be created under in the form `{billing_account_id}` or `billingAccounts/{billing_account_id}`.

* `deletion_policy` (Optional) - If set to "RENAME_ON_DESTROY" the billing account display_name
  will be changed to "Terraform Destroyed" along with a timestamp.  If set to "" this will not occur.
  Default is "".

## Attributes Reference

The following additional attributes are exported:

* `open` - `true` if the billing account is open, `false` if the billing account is closed.

* `name` - The resource name of the billing account in the form `billingAccounts/{billing_account_id}`.

* `billing_account_id` - The billing account id.

## Import

Billing Subaccounts can be imported using any of these accepted formats:

```
$ terraform import google_billing_subaccount.default billingAccounts/{billing_account_id}
```
