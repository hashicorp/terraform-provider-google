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

* `filter` - (Required) A string filter as defined in the [REST API](https://cloud.google.com/resource-manager/reference/rest/v1/projects/list#query-parameters).


## Attributes Reference

The following attributes are exported:

* `projects` - A list of projects matching the provided filter. Structure is [defined below](#nested_projects).

<a name="nested_projects"></a>The `projects` block supports:

* `project_id` - The project id of the project.
* `number` - The numeric identifier of the project.
* `name` - The optional user-assigned display name of the project.
* `labels` - A set of key/value label pairs assigned on a project.
* `lifecycle_state` - The Project lifecycle state.
* `create_time` - Creation time in RFC3339 UTC "Zulu" format.
* `parent` - An optional reference to a parent resource.

