---
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
data "google_project" "project" {}

output "project_number" {
  value = "${data.google_project.project.project_number}"
} 
```

## Argument Reference

There are no arguments available for this data source.


## Attributes Reference

The following attributes are exported:

* `name` - The user-assigned display name of the project

* `project_number` - The number uniquely identifying the project

