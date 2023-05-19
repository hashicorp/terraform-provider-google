---
subcategory: "Cloud Key Management Service"
description: |-
  A datasource to retrieve the IAM policy state for a Google Cloud KMS key ring.
---


# `google_kms_key_ring_iam_policy`
Retrieves the current IAM policy data for a Google Cloud KMS key ring.

## example

```hcl
data "google_kms_key_ring_iam_policy" "test_key_ring_iam_policy" {
  key_ring_id = "{project_id}/{location_name}/{key_ring_name}"
}
```

## Argument Reference

The following arguments are supported:

* `key_ring_id` - (Required) The key ring ID, in the form
    `{project_id}/{location_name}/{key_ring_name}` or
    `{location_name}/{key_ring_name}`. In the second form, the provider's
    project setting will be used as a fallback.

## Attributes Reference

The attributes are exported:

* `etag` - (Computed) The etag of the IAM policy.

* `policy_data` - (Computed) The policy data
