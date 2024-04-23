---
subcategory: "Cloud Platform"
description: |-
  A datasource to retrieve the IAM policy state for a service account.
---


# `google_service_account_iam_policy`
Retrieves the current IAM policy data for a service account.

## example

```hcl
data "google_service_account_iam_policy" "foo" {
  service_account_id = google_service_account.test_account.name
}
```

## Argument Reference

The following arguments are supported:

* `service_account_id` - (Required) The fully-qualified name of the service account to apply policy to.

## Attributes Reference

The attributes are exported:

* `etag` - (Computed) The etag of the IAM policy.

* `policy_data` - (Computed) The policy data
