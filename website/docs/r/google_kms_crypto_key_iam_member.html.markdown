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
resource "google_kms_key_ring" "keyring" {
  name     = "keyring-example"
  location = "global"
}

resource "google_kms_crypto_key" "key" {
  name            = "crypto-key-example"
  key_ring        = google_kms_key_ring.keyring.id
  rotation_period = "100000s"

  lifecycle {
    prevent_destroy = true
  }
}

resource "google_kms_crypto_key_iam_member" "crypto_key" {
  crypto_key_id = "google_kms_crypto_key.key.id"
  role          = "roles/cloudkms.cryptoKeyEncrypter"
  member        = "user:alice@gmail.com"
}
```

With IAM Conditions ([beta](https://terraform.io/docs/providers/google/provider_versions.html)):
```hcl
resource "google_kms_key_ring" "keyring" {
  name     = "keyring-example"
  location = "global"
}

resource "google_kms_crypto_key" "key" {
  name            = "crypto-key-example"
  key_ring        = google_kms_key_ring.keyring.id
  rotation_period = "100000s"

  lifecycle {
    prevent_destroy = true
  }
}

resource "google_kms_crypto_key_iam_member" "crypto_key" {
  crypto_key_id = "google_kms_crypto_key.key.id"
  role          = "roles/cloudkms.cryptoKeyEncrypter"
  member        = "user:alice@gmail.com"

  condition {
    title       = "expires_after_2019_12_31"
    description = "Expiring at midnight of 2019-12-31"
    expression  = "request.time < timestamp(\"2020-01-01T00:00:00Z\")"
  }
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

* `condition` - (Optional, [Beta](https://terraform.io/docs/providers/google/provider_versions.html)) An [IAM Condition](https://cloud.google.com/iam/docs/conditions-overview) for a given binding.
  Structure is documented below.

---

The `condition` block supports:

* `expression` - (Required) Textual representation of an expression in Common Expression Language syntax.

* `title` - (Required) A title for the expression, i.e. a short string describing its purpose.

* `description` - (Optional) An optional description of the expression. This is a longer text which describes the expression, e.g. when hovered over it in a UI.

~> **Warning:** Terraform considers the `role` and condition contents (`title`+`description`+`expression`) as the
  identifier for the binding. This means that if any part of the condition is changed out-of-band, Terraform will
  consider it to be an entirely different resource and will treat it as such.


## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `etag` - (Computed) The etag of the project's IAM policy.

## Import

IAM member imports use space-delimited identifiers; the resource in question, the role, and the account.  This member resource can be imported using the `crypto_key_id`, role, and member identity e.g.

```
$ terraform import google_kms_crypto_key_iam_member.member "your-project-id/location-name/key-ring-name/key-name roles/viewer user:foo@example.com"
```

-> If you're importing a resource with beta features, make sure to include `-provider=google-beta`
as an argument so that Terraform uses the correct provider to import your resource.
