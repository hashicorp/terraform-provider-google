---
layout: "google"
page_title: "Google: google_kms_crypto_key_iam_binding"
sidebar_current: "docs-google-kms-crypto-key-iam-binding"
description: |-
 Allows management of a single binding with an IAM policy for a Google Cloud KMS crypto key
---

# google\_kms\_crypto\_key\_iam\_binding

Allows creation and management of a single binding within IAM policy for
an existing Google Cloud KMS crypto key.

## Example Usage

```hcl
resource "google_kms_crypto_key_iam_binding" "crypto_key" {
  crypto_key_id = "your-crypto-key-id"
  role          = "roles/editor"

  members = [
    "user:jane@example.com",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `members` - (Required) A list of users that the role should apply to.

* `role` - (Required) The role that should be applied. Only one
    `google_kms_crypto_key_iam_binding` can be used per role. Note that custom roles must be of the format
    `[projects|organizations]/{parent-name}/roles/{role-name}`.

* `crypto_key_id` - (Required) The crypto key ID, in the form
    `{project_id}/{location_name}/{key_ring_name}/{crypto_key_name}` or
    `{location_name}/{key_ring_name}/{crypto_key_name}`.
    In the second form, the provider's project setting will be used as a fallback.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `etag` - (Computed) The etag of the crypto key's IAM policy.

## Import

IAM binding imports use space-delimited identifiers; first the resource in question and then the role.  These bindings can be imported using the `crypto_key_id` and role, e.g.

```
$ terraform import google_kms_crypto_key_iam_binding.my_binding "your-project-id/location-name/key-name roles/viewer"
```
