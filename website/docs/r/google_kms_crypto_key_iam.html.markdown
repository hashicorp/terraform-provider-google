---
subcategory: "Cloud Key Management Service"
description: |-
 Collection of resources to manage IAM policy for a Google Cloud KMS crypto key.
---

# IAM policy for Google Cloud KMS crypto key

Three different resources help you manage your IAM policy for KMS crypto key. Each of these resources serves a different use case:

* `google_kms_crypto_key_iam_policy`: Authoritative. Sets the IAM policy for the crypto key and replaces any existing policy already attached.
* `google_kms_crypto_key_iam_binding`: Authoritative for a given role. Updates the IAM policy to grant a role to a list of members. Other roles within the IAM policy for the crypto key are preserved.
* `google_kms_crypto_key_iam_member`: Non-authoritative. Updates the IAM policy to grant a role to a new member. Other members for the role for the crypto key are preserved.

~> **Note:** `google_kms_crypto_key_iam_policy` **cannot** be used in conjunction with `google_kms_crypto_key_iam_binding` and `google_kms_crypto_key_iam_member` or they will fight over what your policy should be.

~> **Note:** `google_kms_crypto_key_iam_binding` resources **can be** used in conjunction with `google_kms_crypto_key_iam_member` resources **only if** they do not grant privilege to the same role.

# google\_kms\_crypto\_key\_iam\_policy

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

data "google_iam_policy" "admin" {
  binding {
    role = "roles/cloudkms.cryptoKeyEncrypter"

    members = [
      "user:jane@example.com",
    ]
  }
}

resource "google_kms_crypto_key_iam_policy" "crypto_key" {
  crypto_key_id = google_kms_crypto_key.key.id
  policy_data = data.google_iam_policy.admin.policy_data
}
```

With IAM Conditions ([beta](https://terraform.io/docs/providers/google/provider_versions.html)):

```hcl
data "google_iam_policy" "admin" {
  binding {
    role = "roles/cloudkms.cryptoKeyEncrypter"

    members = [
      "user:jane@example.com",
    ]

    condition {
      title       = "expires_after_2019_12_31"
      description = "Expiring at midnight of 2019-12-31"
      expression  = "request.time < timestamp(\"2020-01-01T00:00:00Z\")"
    }
  }
}
```

# google\_kms\_crypto\_key\_iam\_binding

```hcl
resource "google_kms_crypto_key_iam_binding" "crypto_key" {
  crypto_key_id = google_kms_crypto_key.key.id
  role          = "roles/cloudkms.cryptoKeyEncrypter"

  members = [
    "user:jane@example.com",
  ]
}
```

With IAM Conditions ([beta](https://terraform.io/docs/providers/google/provider_versions.html)):

```hcl
resource "google_kms_crypto_key_iam_binding" "crypto_key" {
  crypto_key_id = google_kms_crypto_key.key.id
  role          = "roles/cloudkms.cryptoKeyEncrypter"

  members = [
    "user:jane@example.com",
  ]

  condition {
    title       = "expires_after_2019_12_31"
    description = "Expiring at midnight of 2019-12-31"
    expression  = "request.time < timestamp(\"2020-01-01T00:00:00Z\")"
  }
}
```

# google\_kms\_crypto\_key\_iam\_member

```hcl
resource "google_kms_crypto_key_iam_member" "crypto_key" {
  crypto_key_id = google_kms_crypto_key.key.id
  role          = "roles/cloudkms.cryptoKeyEncrypter"
  member        = "user:jane@example.com"
}
```

With IAM Conditions ([beta](https://terraform.io/docs/providers/google/provider_versions.html)):

```hcl
resource "google_kms_crypto_key_iam_member" "crypto_key" {
  crypto_key_id = google_kms_crypto_key.key.id
  role          = "roles/cloudkms.cryptoKeyEncrypter"
  member        = "user:jane@example.com"

  condition {
    title       = "expires_after_2019_12_31"
    description = "Expiring at midnight of 2019-12-31"
    expression  = "request.time < timestamp(\"2020-01-01T00:00:00Z\")"
  }
}
```

## Argument Reference

The following arguments are supported:

* `crypto_key_id` - (Required) The crypto key ID, in the form
    `{project_id}/{location_name}/{key_ring_name}/{crypto_key_name}` or
    `{location_name}/{key_ring_name}/{crypto_key_name}`. In the second form,
    the provider's project setting will be used as a fallback.

* `member/members` - (Required) Identities that will be granted the privilege in `role`.
  Each entry can have one of the following values:
  * **allUsers**: A special identifier that represents anyone who is on the internet; with or without a Google account.
  * **allAuthenticatedUsers**: A special identifier that represents anyone who is authenticated with a Google account or a service account.
  * **user:{emailid}**: An email address that represents a specific Google account. For example, jane@example.com or joe@example.com.
  * **serviceAccount:{emailid}**: An email address that represents a service account. For example, my-other-app@appspot.gserviceaccount.com.
  * **group:{emailid}**: An email address that represents a Google group. For example, admins@example.com.
  * **domain:{domain}**: A G Suite domain (primary, instead of alias) name that represents all the users of that domain. For example, google.com or example.com.

* `role` - (Required) The role that should be applied. Note that custom roles must be of the format
    `[projects|organizations]/{parent-name}/roles/{role-name}`.

* `policy_data` - (Required only by `google_kms_crypto_key_iam_policy`) The policy data generated by
  a `google_iam_policy` data source.

* `condition` - (Optional, [Beta](https://terraform.io/docs/providers/google/provider_versions.html)) An [IAM Condition](https://cloud.google.com/iam/docs/conditions-overview) for a given binding.
  Structure is [documented below](#nested_condition).

---

<a name="nested_condition"></a>The `condition` block supports:

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

### Importing IAM members

IAM member imports use space-delimited identifiers that contain the resource's `crypto_key_id`, `role`, and member identity e.g.

* `"{{project_id}}/{{location}}/{{key_ring_name}}/{{crypto_key_name}} roles/viewer user:foo@example.com"`

An [`import` block](https://developer.hashicorp.com/terraform/language/import) (Terraform v1.5.0 and later) can be used to import IAM members:

```tf
import {
  id = "{{project_id}}/{{location}}/{{key_ring_name}}/{{crypto_key_name}} roles/viewer user:foo@example.com"
  to = google_kms_crypto_key_iam_member.default
}
```

The [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import) can also be used:

```
$ terraform import google_kms_crypto_key_iam_member.default "{{project_id}}/{{location}}/{{key_ring_name}}/{{crypto_key_name}} roles/viewer user:foo@example.com"
```

### Importing IAM bindings

IAM binding imports use space-delimited identifiers that contain the `crypto_key_id` and role, e.g.

* `"{{project_id}}/{{location}}/{{key_ring_name}}/{{crypto_key_name}} roles/viewer"`

An [`import` block](https://developer.hashicorp.com/terraform/language/import) (Terraform v1.5.0 and later) can be used to import IAM bindings:

```tf
import {
  id = "{{project_id}}/{{location}}/{{key_ring_name}}/{{crypto_key_name}} roles/viewer"
  to = google_kms_crypto_key_iam_binding.default
}
```

The [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import) can also be used:

```
$ terraform import google_kms_crypto_key_iam_binding.default "{{project_id}}/{{location}}/{{key_ring_name}}/{{crypto_key_name}} roles/viewer"
```

### Importing IAM policies

IAM policy imports use the identifier of the KMS crypto key only. For example:

* `{{project_id}}/{{location}}/{{key_ring_name}}/{{crypto_key_name}}`

An [`import` block](https://developer.hashicorp.com/terraform/language/import) (Terraform v1.5.0 and later) can be used to import IAM policies:

```tf
import {
  id = "{{project_id}}/{{location}}/{{key_ring_name}}/{{crypto_key_name}}"
  to = google_kms_crypto_key_iam_policy.default
}
```

The [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import) can also be used:

```
$ terraform import google_kms_crypto_key_iam_policy.default {{project_id}}/{{location}}/{{key_ring_name}}/{{crypto_key_name}}
```