---
subcategory: "Cloud Platform"
layout: "google"
page_title: "Google: google_iam_role"
sidebar_current: "docs-google-datasource-iam-role"
description: |-
  Get information about a Google IAM Role.
---

# google\_iam\_role

Use this data source to get information about a Google IAM Role.

```hcl
data "google_iam_role" "roleinfo" {
  name = "roles/compute.viewer"
}

output "the_role_permissions" {
  value = data.google_iam_role.roleinfo.included_permissions
}
```

## Argument Reference

The following arguments are supported:

* `name` (Required) - The name of the Role to lookup in the form `roles/{ROLE_NAME}`, `organizations/{ORGANIZATION_ID}/roles/{ROLE_NAME}` or `projects/{PROJECT_ID}/roles/{ROLE_NAME}`

## Attributes Reference

The following attributes are exported:

* `title` - is a friendly title for the role, such as "Role Viewer"
* `included_permissions` - specifies the list of one or more permissions to include in the custom role, such as - `iam.roles.get`
* `stage` -  indicates the stage of a role in the launch lifecycle, such as `GA`, `BETA` or `ALPHA`.
