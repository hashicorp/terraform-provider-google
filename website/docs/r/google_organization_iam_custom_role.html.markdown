---
layout: "google"
page_title: "Google: google_organization_iam_custom_role"
sidebar_current: "docs-google-organization-iam-custom-role"
description: |-
 Allows management of a customized Cloud IAM organization role.
---

# google\_organization\_iam\_custom\_role

Allows management of a customized Cloud IAM organization role. For more information see
[the official documentation](https://cloud.google.com/iam/docs/understanding-custom-roles)
and
[API](https://cloud.google.com/iam/reference/rest/v1/organizations.roles).

## Example Usage

This snippet creates a customized IAM organization role.

```hcl
resource "google_organization_iam_custom_role" "my-custom-role" {
  role_id     = "myCustomRole"
  org_id      = "123456789"
  title       = "My Custom Role"
  description = "A description"
  permissions = ["iam.roles.list", "iam.roles.create", "iam.roles.delete"]
}
```

## Argument Reference

The following arguments are supported:

* `role_id` - (Required) The role id to use for this role.

* `org_id` - (Required) The numeric ID of the organization in which you want to create a custom role.

* `title` - (Required) A human-readable title for the role.

* `permissions` (Required) The names of the permissions this role grants when bound in an IAM policy. At least one permission must be specified.

* `stage` - (Optional) The current launch stage of the role.
    Defaults to `GA`.
    List of possible stages is [here](https://cloud.google.com/iam/reference/rest/v1/organizations.roles#Role.RoleLaunchStage).

* `description` - (Optional) A human-readable description for the role.

* `deleted` - (Optional) The current deleted state of the role. Defaults to `false`.

## Import

Customized IAM organization role can be imported using their URI, e.g.

```
$ terraform import google_organization_iam_custom_role.my-custom-role organizations/123456789/roles/myCustomRole
```
