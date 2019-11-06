---
subcategory: "Cloud Platform"
layout: "google"
page_title: "Google: google_organization_iam_binding"
sidebar_current: "docs-google-organization-iam-binding"
description: |-
 Allows management of a single binding with an IAM policy for a Google Cloud Platform Organization.
---

# google\_organization\_iam\_binding

Allows creation and management of a single binding within IAM policy for
an existing Google Cloud Platform Organization.

~> **Note:** This resource __must not__ be used in conjunction with
   `google_organization_iam_member` for the __same role__ or they will fight over
   what your policy should be.

~> **Note:** On create, this resource will overwrite members of any existing roles.
    Use `terraform import` and inspect the `terraform plan` output to ensure
    your existing members are preserved.

## Example Usage

```hcl
resource "google_organization_iam_binding" "binding" {
  org_id = "123456789"
  role    = "roles/browser"

  members = [
    "user:alice@gmail.com",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `org_id` - (Required) The numeric ID of the organization in which you want to create a custom role.

* `role` - (Required) The role that should be applied. Only one
    `google_organization_iam_binding` can be used per role. Note that custom roles must be of the format
    `[projects|organizations]/{parent-name}/roles/{role-name}`.

* `members` - (Required) A list of users that the role should apply to. For more details on format and restrictions see https://cloud.google.com/billing/reference/rest/v1/Policy#Binding

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `etag` - (Computed) The etag of the organization's IAM policy.

## Import

IAM binding imports use space-delimited identifiers; first the resource in question and then the role.  These bindings can be imported using the `org_id` and role, e.g.

```
$ terraform import google_organization_iam_binding.my_org "your-org-id roles/viewer"
```
