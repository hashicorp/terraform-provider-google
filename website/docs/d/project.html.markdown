---
subcategory: "Cloud Platform"
page_title: "Google: google_project"
description: |-
  Retrieve project details
---

# google\_project

Use this data source to get project details.
For more information see
[API](https://cloud.google.com/resource-manager/reference/rest/v1/projects#Project)

## Example Usage

```hcl
data "google_project" "project" {
}

output "project_number" {
  value = data.google_project.project.number
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Optional) The project ID. If it is not provided, the provider project is used.


## Attributes Reference

The following attributes are exported:

* `number` - The numeric identifier of the project.

See [google_project](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/google_project) resource for details of the available attributes.

