---
layout: "google"
page_title: "Google: google_kms_key_ring_iam_member"
sidebar_current: "docs-google-kms-key-ring-iam-member"
description: |-
 Allows management of a single member for a single binding on the IAM policy for a Google Cloud KMS key ring.
---

# google\_kms\_key\_ring\_iam\_member

Allows creation and management of a single member for a single binding within
the IAM policy for an existing Google Cloud KMS key ring.

~> **Note:** This resource _must not_ be used in conjunction with
   `google_kms_key_ring_iam_policy` or they will fight over what your policy
   should be. Similarly, roles controlled by `google_kms_key_ring_iam_binding`
   should not be assigned to using `google_kms_key_ring_iam_member`.

## Example Usage

```hcl
resource "google_kms_key_ring_iam_member" "key_ring" {
  key_ring_id = "your-key-ring-id"
  role        = "roles/editor"
  member      = "user:jane@example.com"
}
```

## Argument Reference

The following arguments are supported:

* `member` - (Required) The user that the role should apply to.

* `role` - (Required) The role that should be applied.

* `key_ring_id` - (Required) The key ring ID, in the form
    `{project_id}/{location_name}/{key_ring_name}` or
    `{location_name}/{key_ring_name}`. In the second form, the provider's
    project setting will be used as a fallback.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `etag` - (Computed) The etag of the project's IAM policy.
