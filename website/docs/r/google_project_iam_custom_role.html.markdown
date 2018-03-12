---
layout: "google"
page_title: "Google: google_project_iam_custom_role"
sidebar_current: "docs-google-project-iam-custom-role"
description: |-
 Allows management of a customized Cloud IAM project role.
---

# google\_project\_iam\_custom\_role

Allows management of a customized Cloud IAM project role. For more information see
[the official documentation](https://cloud.google.com/iam/docs/understanding-custom-roles)
and
[API](https://cloud.google.com/iam/reference/rest/v1/projects.roles).

## Example Usage

This snippet creates a customized IAM role.

```hcl
resource "google_project_iam_custom_role" "my-custom-role" {
  role_id     = "myCustomRole"
  title       = "My Custom Role"
  description = "A description"
  permissions = ["iam.roles.list", "iam.roles.create", "iam.roles.delete"]
}
```

## Argument Reference

The following arguments are supported:

* `role_id` - (Required) The role id to use for this role.

* `title` - (Required) A human-readable title for the role.

* `permissions` (Required) The names of the permissions this role grants when bound in an IAM policy. At least one permission must be specified.

* `project` - (Optional) The project that the service account will be created in.
    Defaults to the provider project configuration.

* `stage` - (Optional) The current launch stage of the role.
    Defaults to `GA`.
    List of possible stages is [here](https://cloud.google.com/iam/reference/rest/v1/organizations.roles#Role.RoleLaunchStage).

* `description` - (Optional) A human-readable description for the role.

* `deleted` - (Optional) The current deleted state of the role. Defaults to `false`.

## Import

Customized IAM project role can be imported using their URI, e.g.

```
$ terraform import google_project_iam_custom_role.my-custom-role projects/my-project/roles/myCustomRole
```
