---
subcategory: "Cloud Billing"
description: |-
  A datasource to retrieve the IAM policy state for a Billing Account.
---


# `google_billing_account_iam_policy`
Retrieves the current IAM policy data for a Billing Account.

## example

```hcl
data "google_billing_account_iam_policy" "policy" {
  billing_account_id = "MEEP-MEEP-MEEP-MEEP-MEEP"
}
```

## Argument Reference

The following arguments are supported:

* `billing_account_id` - (Required) The billing account id.

## Attributes Reference

The attributes are exported:

* `etag` - (Computed) The etag of the IAM policy.

* `policy_data` - (Computed) The policy data
