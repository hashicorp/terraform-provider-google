---
subcategory: "Cloud Platform"
description: |-
  Get information about a Google Cloud Organization IAM Custom Role.
---

# google_organization_iam_custom_role

Get information about a Google Cloud Organization IAM Custom Role. Note that you must have the `roles/iam.organizationRoleViewer` role (or equivalent permissions) at the organization level to use this datasource.

```hcl
data "google_organization_iam_custom_role" "example" {
  org_id  = "1234567890"
  role_id = "your-role-id"
}

resource "google_project_iam_member" "project" {
  project = "your-project-id"
  role    = data.google_organization_iam_custom_role.example.name
  member  = "user:jane@example.com"
}
```

## Argument Reference

The following arguments are supported:

* `org_id` - (Required) The numeric ID of the organization in which you want to create a custom role.

* `role_id` - (Required) The role id that has been used for this role.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

See [google_organization_iam_custom_role](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/google_organization_iam_custom_role) resource for details of the available attributes.

