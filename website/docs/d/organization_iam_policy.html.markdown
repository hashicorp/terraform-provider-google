---
subcategory: "Cloud Platform"
description: |-
  A datasource to retrieve the IAM policy state for a organization.
---


# `google_organization_iam_policy`
Retrieves the current IAM policy data for a organization.

## example

```hcl
data "google_organization_iam_policy" "policy" {
  org_id  = "123456789"
}
```

## Argument Reference

The following arguments are supported:

* `org_id` - (Required) The organization id of the target organization.

## Attributes Reference

The attributes are exported:

* `etag` - (Computed) The etag of the IAM policy.

* `policy_data` - (Computed) The policy data
