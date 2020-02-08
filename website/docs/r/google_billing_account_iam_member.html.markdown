---
subcategory: "Cloud Platform"
layout: "google"
page_title: "Google: google_billing_account_iam_member"
sidebar_current: "docs-google-billing-account-iam-member"
description: |-
 Allows management of a single member for a single binding on the IAM policy for a Google Cloud Platform Billing Account.
---

# google\_billing\_account\_iam\_member

Allows creation and management of a single member for a single binding within
the IAM policy for an existing Google Cloud Platform Billing Account.

~> **Note:** This resource __must not__ be used in conjunction with
   `google_billing_account_iam_binding` for the __same role__ or they will fight over
   what your policy should be.

## Example Usage

```hcl
resource "google_billing_account_iam_member" "binding" {
  billing_account_id = "00AA00-000AAA-00AA0A"
  role               = "roles/billing.viewer"
  member             = "user:alice@gmail.com"
}
```

## Argument Reference

The following arguments are supported:

* `billing_account_id` - (Required) The billing account id.

* `role` - (Required) The role that should be applied.

* `member` - (Required) The user that the role should apply to. For more details on format and restrictions see https://cloud.google.com/billing/reference/rest/v1/Policy#Binding

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `etag` - (Computed) The etag of the billing account's IAM policy.

## Import

IAM member imports use space-delimited identifiers; the resource in question, the role, and the account.  This member resource can be imported using the `billing_account_id`, role, and member identity, e.g.

```
$ terraform import google_billing_account_iam_member.binding "your-billing-account-id roles/viewer user:foo@example.com"
```

-> **Custom Roles**: If you're importing a IAM member with a custom role, make sure to use the
 full name of the custom role, e.g. `[projects/my-project|organizations/my-org]/roles/my-custom-role`.
