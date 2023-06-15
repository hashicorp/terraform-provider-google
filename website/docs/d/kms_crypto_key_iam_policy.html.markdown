---
subcategory: "Cloud Key Management Service"
description: |-
  A datasource to retrieve the IAM policy state for a Google Cloud KMS crypto key.
---


# `google_kms_crypto_key_iam_policy`
Retrieves the current IAM policy data for a Google Cloud KMS crypto key.

## example

```hcl
data "google_kms_crypto_key_iam_policy" "foo" {
  crypto_key_id = google_kms_crypto_key.crypto_key.id
}
```

## Argument Reference

The following arguments are supported:

* `crypto_key_id` - (Required) The crypto key ID, in the form

## Attributes Reference

The attributes are exported:

* `etag` - (Computed) The etag of the IAM policy.

* `policy_data` - (Computed) The policy data
