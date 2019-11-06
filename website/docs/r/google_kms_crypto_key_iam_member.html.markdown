---
subcategory: "Cloud KMS"
layout: "google"
page_title: "Google: google_kms_crypto_key_iam_member"
sidebar_current: "docs-google-kms-crypto-key-iam-member"
description: |-
 Allows management of a single member for a single binding on the IAM policy for a Google Cloud KMS crypto key.
---

# google\_kms\_crypto\_key\_iam\_member

Allows creation and management of a single member for a single binding within
the IAM policy for an existing Google Cloud KMS crypto key.

~> **Note:** This resource _must not_ be used in conjunction with
   `google_kms_crypto_key_iam_policy` or they will fight over what your policy
   should be. Similarly, roles controlled by `google_kms_crypto_key_iam_binding`
   should not be assigned to using `google_kms_crypto_key_iam_member`.

## Example Usage

```hcl
resource "google_kms_crypto_key_iam_member" "crypto_key" {
  crypto_key_id = "your-crypto-key-id"
  role          = "roles/editor"
  member        = "user:alice@gmail.com"
}
```

## Argument Reference

The following arguments are supported:

* `member` - (Required) The user that the role should apply to. For more details on format and restrictions see https://cloud.google.com/billing/reference/rest/v1/Policy#Binding

* `role` - (Required) The role that should be applied. Note that custom roles must be of the format
    `[projects|organizations]/{parent-name}/roles/{role-name}`.

* `crypto_key_id` - (Required) The key ring ID, in the form
    `{project_id}/{location_name}/{key_ring_name}/{crypto_key_name}` or
    `{location_name}/{key_ring_name}/{crypto_key_name}`. In the second form,
    the provider's project setting will be used as a fallback.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `etag` - (Computed) The etag of the project's IAM policy.

## Import

IAM member imports use space-delimited identifiers; the resource in question, the role, and the account.  This member resource can be imported using the `crypto_key_id`, role, and member identity e.g.

```
$ terraform import google_kms_crypto_key_iam_member.member "your-project-id/location-name/key-ring-name/key-name roles/viewer user:foo@example.com"
```
