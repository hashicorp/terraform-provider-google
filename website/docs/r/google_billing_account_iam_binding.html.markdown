---
subcategory: "Cloud Platform"
layout: "google"
page_title: "Google: google_billing_account_iam_binding"
sidebar_current: "docs-google-billing-account-iam-binding"
description: |-
 Allows management of a single binding with an IAM policy for a Google Cloud Platform Billing Account.
---

# google\_billing\_account\_iam\_binding

Allows creation and management of a single binding within IAM policy for
an existing Google Cloud Platform Billing Account.

~> **Note:** This resource __must not__ be used in conjunction with
   `google_billing_account_iam_member` for the __same role__ or they will fight over
   what your policy should be.

~> **Note:** On create, this resource will overwrite members of any existing roles.
    Use `terraform import` and inspect the `terraform plan` output to ensure
    your existing members are preserved.

## Example Usage

```hcl
resource "google_billing_account_iam_binding" "binding" {
  billing_account_id = "00AA00-000AAA-00AA0A"
  role               = "roles/billing.viewer"

  members = [
    "user:alice@gmail.com",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `billing_account_id` - (Required) The billing account id.

* `role` - (Required) The role that should be applied.

* `members` - (Required) A list of users that the role should apply to. For more details on format and restrictions see https://cloud.google.com/billing/reference/rest/v1/Policy#Binding

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `etag` - (Computed) The etag of the billing account's IAM policy.

## Import

IAM binding imports use space-delimited identifiers; first the resource in question and then the role.  These bindings can be imported using the `billing_account_id` and role, e.g.

```
$ terraform import google_billing_account_iam_binding.binding "your-billing-account-id roles/viewer"
```
