---
subcategory: "Cloud KMS"
layout: "google"
page_title: "Google: google_kms_crypto_key_iam_binding"
sidebar_current: "docs-google-kms-crypto-key-iam-binding"
description: |-
 Allows management of a single binding with an IAM policy for a Google Cloud KMS crypto key
---

# google\_kms\_crypto\_key\_iam\_binding

Allows creation and management of a single binding within IAM policy for
an existing Google Cloud KMS crypto key.

~> **Note:** On create, this resource will overwrite members of any existing roles.
    Use `terraform import` and inspect the `terraform plan` output to ensure
    your existing members are preserved.

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

resource "google_kms_crypto_key_iam_binding" "crypto_key" {
  crypto_key_id = "google_kms_crypto_key.key.id"
  role          = "roles/cloudkms.cryptoKeyEncrypter"

  members = [
    "user:alice@gmail.com",
  ]
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

resource "google_kms_crypto_key_iam_binding" "crypto_key" {
  crypto_key_id = "google_kms_crypto_key.key.id"
  role          = "roles/cloudkms.cryptoKeyEncrypter"

  members = [
    "user:alice@gmail.com",
  ]

  condition {
    title       = "expires_after_2019_12_31"
    description = "Expiring at midnight of 2019-12-31"
    expression  = "request.time < timestamp(\"2020-01-01T00:00:00Z\")"
  }
}
```

## Argument Reference

The following arguments are supported:

* `members` - (Required) A list of users that the role should apply to. For more details on format and restrictions see https://cloud.google.com/billing/reference/rest/v1/Policy#Binding

* `role` - (Required) The role that should be applied. Only one
    `google_kms_crypto_key_iam_binding` can be used per role. Note that custom roles must be of the format
    `[projects|organizations]/{parent-name}/roles/{role-name}`.

* `crypto_key_id` - (Required) The crypto key ID, in the form
    `{project_id}/{location_name}/{key_ring_name}/{crypto_key_name}` or
    `{location_name}/{key_ring_name}/{crypto_key_name}`.
    In the second form, the provider's project setting will be used as a fallback.

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

* `etag` - (Computed) The etag of the crypto key's IAM policy.

## Import

IAM binding imports use space-delimited identifiers; first the resource in question and then the role.  These bindings can be imported using the `crypto_key_id` and role, e.g.

```
$ terraform import google_kms_crypto_key_iam_binding.crypto_key "my-gcp-project/us-central1/my-key-ring/my-crypto-key roles/editor"
```

-> If you're importing a resource with beta features, make sure to include `-provider=google-beta`
as an argument so that Terraform uses the correct provider to import your resource.
