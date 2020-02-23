---
subcategory: "Cloud Platform"
layout: "google"
page_title: "Google: google_billing_account_iam_policy"
sidebar_current: "docs-google-billing-account-iam-policy"
description: |-
 Allows management of the entire IAM policy for a Google Cloud Platform Billing Account.
---

# google\_billing\_account\_iam\_policy

Allows management of the entire IAM policy for an existing Google Cloud Platform Billing Account.

~> **Warning:** Billing accounts have a default user that can be **overwritten**
by use of this resource. The safest alternative is to use multiple `google_billing_account_iam_binding`
   resources. If you do use this resource, the best way to be sure that you are
   not making dangerous changes is to start by importing your existing policy,
   and examining the diff very closely.

~> **Note:** This resource __must not__ be used in conjunction with
   `google_billing_account_iam_member` or `google_billing_account_iam_binding`
   or they will fight over what your policy should be.

## Example Usage

```hcl
resource "google_billing_account_iam_policy" "policy" {
  billing_account_id = "00AA00-000AAA-00AA0A"
  policy_data        = data.google_iam_policy.admin.policy_data
}

data "google_iam_policy" "admin" {
  binding {
    role = "roles/billing.viewer"

    members = [
      "user:jane@example.com",
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `billing_account_id` - (Required) The billing account id.

* `policy_data` - (Required) The `google_iam_policy` data source that represents
    the IAM policy that will be applied to the billing account. This policy overrides any existing
    policy applied to the billing account.

## Import

```
$ terraform import google_billing_account_iam_policy.policy billing-account-id
```
