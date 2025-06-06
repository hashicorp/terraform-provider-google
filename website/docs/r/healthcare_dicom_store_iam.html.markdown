---
# ----------------------------------------------------------------------------
#
#     ***     AUTO GENERATED CODE    ***    Type: Handwritten     ***
#
# ----------------------------------------------------------------------------
#
#     This code is generated by Magic Modules using the following:
#
#     Source file: https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/third_party/terraform/website/docs/r/healthcare_dicom_store_iam.html.markdown
#
#     DO NOT EDIT this file directly. Any changes made to this file will be
#     overwritten during the next generation cycle.
#
# ----------------------------------------------------------------------------
subcategory: "Cloud Healthcare"
description: |-
 Collection of resources to manage IAM policy for a Google Cloud Healthcare DICOM store.
---

# IAM policy for Google Cloud Healthcare DICOM store

~> **Warning:** These resources are in beta, and should be used with the terraform-provider-google-beta provider.
See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html) for more details on beta resources.

Three different resources help you manage your IAM policy for Healthcare DICOM store. Each of these resources serves a different use case:

* `google_healthcare_dicom_store_iam_policy`: Authoritative. Sets the IAM policy for the DICOM store and replaces any existing policy already attached.
* `google_healthcare_dicom_store_iam_binding`: Authoritative for a given role. Updates the IAM policy to grant a role to a list of members. Other roles within the IAM policy for the DICOM store are preserved.
* `google_healthcare_dicom_store_iam_member`: Non-authoritative. Updates the IAM policy to grant a role to a new member. Other members for the role for the DICOM store are preserved.

~> **Note:** `google_healthcare_dicom_store_iam_policy` **cannot** be used in conjunction with `google_healthcare_dicom_store_iam_binding` and `google_healthcare_dicom_store_iam_member` or they will fight over what your policy should be.

~> **Note:** `google_healthcare_dicom_store_iam_binding` resources **can be** used in conjunction with `google_healthcare_dicom_store_iam_member` resources **only if** they do not grant privilege to the same role.

## google_healthcare_dicom_store_iam_policy

```hcl
data "google_iam_policy" "admin" {
  binding {
    role = "roles/editor"

    members = [
      "user:jane@example.com",
    ]
  }
}

resource "google_healthcare_dicom_store_iam_policy" "dicom_store" {
  dicom_store_id = "your-dicom-store-id"
  policy_data    = data.google_iam_policy.admin.policy_data
}
```

## google_healthcare_dicom_store_iam_binding

```hcl
resource "google_healthcare_dicom_store_iam_binding" "dicom_store" {
  dicom_store_id = "your-dicom-store-id"
  role           = "roles/editor"

  members = [
    "user:jane@example.com",
  ]
}
```

## google_healthcare_dicom_store_iam_member

```hcl
resource "google_healthcare_dicom_store_iam_member" "dicom_store" {
  dicom_store_id = "your-dicom-store-id"
  role           = "roles/editor"
  member         = "user:jane@example.com"
}
```

## Argument Reference

The following arguments are supported:

* `dicom_store_id` - (Required) The DICOM store ID, in the form
    `{project_id}/{location_name}/{dataset_name}/{dicom_store_name}` or
    `{location_name}/{dataset_name}/{dicom_store_name}`. In the second form, the provider's
    project setting will be used as a fallback.

* `member/members` - (Required) Identities that will be granted the privilege in `role`.
  Each entry can have one of the following values:
  * **allUsers**: A special identifier that represents anyone who is on the internet; with or without a Google account.
  * **allAuthenticatedUsers**: A special identifier that represents anyone who is authenticated with a Google account or a service account.
  * **user:{emailid}**: An email address that represents a specific Google account. For example, alice@gmail.com or joe@example.com.
  * **serviceAccount:{emailid}**: An email address that represents a service account. For example, my-other-app@appspot.gserviceaccount.com.
  * **group:{emailid}**: An email address that represents a Google group. For example, admins@example.com.
  * **domain:{domain}**: A G Suite domain (primary, instead of alias) name that represents all the users of that domain. For example, google.com or example.com.

* `role` - (Required) The role that should be applied. Only one
    `google_healthcare_dicom_store_iam_binding` can be used per role. Note that custom roles must be of the format
    `[projects|organizations]/{parent-name}/roles/{role-name}`.

* `policy_data` - (Required only by `google_healthcare_dicom_store_iam_policy`) The policy data generated by
  a `google_iam_policy` data source.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `etag` - (Computed) The etag of the DICOM store's IAM policy.

## Import

-> **Custom Roles** If you're importing a IAM resource with a custom role, make sure to use the
 full name of the custom role, e.g. `[projects/my-project|organizations/my-org]/roles/my-custom-role`.

### Importing IAM members

IAM member imports use space-delimited identifiers that contains the `dicom_store_id`, `role`, and `member`. For example:

* `"{{project_id}}/{{location}}/{{dataset}}/{{dicom_store}} roles/editor jane@example.com"`

An [`import` block](https://developer.hashicorp.com/terraform/language/import) (Terraform v1.5.0 and later) can be used to import IAM members:

```tf
import {
  id = "{{project_id}}/{{location}}/{{dataset}}/{{dicom_store}} roles/editor jane@example.com"
  to = google_healthcare_dicom_store_iam_member.default
}
```

The [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import) can also be used:

```
$ terraform import google_healthcare_dicom_store_iam_member.default "{{project_id}}/{{location}}/{{dataset}}/{{dicom_store}} roles/editor jane@example.com"
```

### Importing IAM bindings

IAM binding imports use space-delimited identifiers that contain the resource's `dicom_store_id` and `role`. For example:

* `"{{project_id}}/{{location}}/{{dataset}}/{{dicom_store}} roles/editor"`

An [`import` block](https://developer.hashicorp.com/terraform/language/import) (Terraform v1.5.0 and later) can be used to import IAM bindings:

```tf
import {
  id = "{{project_id}}/{{location}}/{{dataset}}/{{dicom_store}} roles/editor"
  to = google_healthcare_dicom_store_iam_binding.default
}
```

The [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import) can also be used:

```
$ terraform import google_healthcare_dicom_store_iam_binding.default "{{project_id}}/{{location}}/{{dataset}}/{{dicom_store}} roles/editor"
```

### Importing IAM policies

IAM policy imports use the identifier of the Healthcare DICOM store resource. For example:

* `"{{project_id}}/{{location}}/{{dataset}}/{{dicom_store}}"`

An [`import` block](https://developer.hashicorp.com/terraform/language/import) (Terraform v1.5.0 and later) can be used to import IAM policies:

```tf
import {
  id = "{{project_id}}/{{location}}/{{dataset}}/{{dicom_store}}"
  to = google_healthcare_dicom_store_iam_policy.default
}
```

The [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import) can also be used:

```
$ terraform import google_healthcare_dicom_store_iam_policy.default {{project_id}}/{{location}}/{{dataset}}/{{dicom_store}}
```