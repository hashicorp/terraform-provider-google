---
subcategory: "Cloud Platform"
description: |-
  Get information about a Google Cloud IAM Custom Role from a project.
---

# google_project_iam_custom_role

Get information about a Google Cloud Project IAM Custom Role. Note that you must have the `roles/iam.roleViewer` role (or equivalent permissions) at the project level to use this datasource.

```hcl
data "google_project_iam_custom_role" "example" {
  project = "your-project-id"
  role_id = "your-role-id"
}

resource "google_project_iam_member" "project" {
  project = "your-project-id"
  role    = data.google_project_iam_custom_role.example.name
  member  = "user:jane@example.com"
}
```

## Argument Reference

The following arguments are supported:

* `role_id` - (Required) The role id that has been used for this role.

* `project` - (Optional) The project were the custom role has been created in. Defaults to the provider project configuration.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

See [google_project_iam_custom_role](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/google_project_iam_custom_role) resource for details of the available attributes.

