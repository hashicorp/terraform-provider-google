---
layout: "google"
page_title: "Google: google_kms_key_ring_iam_binding"
sidebar_current: "docs-google-kms-key-ring-iam-binding"
description: |-
 Allows management of a single binding with an IAM policy for a Google Cloud KMS key ring
---

# google\_kms\_key\_ring\_iam\_binding

Allows creation and management of a single binding within IAM policy for
an existing Google Cloud KMS key ring.

## Example Usage

```hcl
resource "google_kms_key_ring_binding" "key_ring" {
  key_ring_id = "your-key-ring-id"
  role        = "roles/editor"

  members = [
    "user:jane@example.com",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `members` - (Required) A list of users that the role should apply to.

* `role` - (Required) The role that should be applied. Only one
    `google_kms_key_ring_iam_binding` can be used per role.

* `key_ring_id` - (Required) The key ring ID, in the form
    `{project_id}/{location_name}/{key_ring_name}` or
    `{location_name}/{key_ring_name}`. In the second form, the provider's
    project setting will be used as a fallback.
## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `etag` - (Computed) The etag of the key ring's IAM policy.

