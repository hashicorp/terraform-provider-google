---
subcategory: "Cloud Platform"
layout: "google"
page_title: "Google: google_projects"
sidebar_current: "docs-google-datasource-projects"
description: |-
  Retrieve a set of projects based on a filter.
---

# google\_projects

Retrieve information about a set of projects based on a filter. See the
[REST API](https://cloud.google.com/resource-manager/reference/rest/v1/projects/list)
for more details.

## Example Usage - searching for projects about to be deleted in an org

```hcl
data "google_projects" "my-org-projects" {
  filter = "parent.id:012345678910 lifecycleState:DELETE_REQUESTED"
}

data "google_project" "deletion-candidate" {
  project_id = data.google_projects.my-org-projects.projects[0].project_id
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) A string filter as defined in the [REST API](https://cloud.google.com/resource-manager/reference/rest/v1/projects/list#query-parameters).


## Attributes Reference

The following attributes are exported:

* `projects` - A list of projects matching the provided filter. Structure is defined below.

The `projects` block supports:

* `project_id` - The project id of the project.

