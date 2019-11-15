---
subcategory: "Cloud Platform"
layout: "google"
page_title: "Google: google_project"
sidebar_current: "docs-google-datasource-project"
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

See [google_project](https://www.terraform.io/docs/providers/google/r/google_project.html) resource for details of the available attributes.

