---
subcategory: "Cloud Platform"
layout: "google"
page_title: "Google: google_organization_iam_policy"
sidebar_current: "docs-google-organization-iam-policy"
description: |-
 Allows management of the entire IAM policy for a Google Cloud Platform Organization.
---

# google\_organization\_iam\_policy

Allows management of the entire IAM policy for an existing Google Cloud Platform Organization.

~> **Warning:** New organizations have several default policies which will,
   without extreme caution, be **overwritten** by use of this resource.
   The safest alternative is to use multiple `google_organization_iam_binding`
   resources.  It is easy to use this resource to remove your own access to
   an organization, which will require a call to Google Support to have
   fixed, and can take multiple days to resolve.  If you do use this resource,
   the best way to be sure that you are not making dangerous changes is to start
   by importing your existing policy, and examining the diff very closely.

~> **Note:** This resource __must not__ be used in conjunction with
   `google_organization_iam_member` or `google_organization_iam_binding`
   or they will fight over what your policy should be.

## Example Usage

```hcl
resource "google_organization_iam_policy" "policy" {
  org_id      = "123456789"
  policy_data = data.google_iam_policy.admin.policy_data
}

data "google_iam_policy" "admin" {
  binding {
    role = "roles/editor"

    members = [
      "user:jane@example.com",
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `org_id` - (Required) The numeric ID of the organization in which you want to create a custom role.

* `policy_data` - (Required) The `google_iam_policy` data source that represents
    the IAM policy that will be applied to the organization. This policy overrides any existing
    policy applied to the organization.

## Import

```
$ terraform import google_organization_iam_policy.my_org your-org-id
```
