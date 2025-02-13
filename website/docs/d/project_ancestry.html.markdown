---
subcategory: "Cloud Platform"
description: |-
  Retrieve the ancestors for a project.
---

# google_project_ancestry

Retrieve the ancestors for a project.
See the [REST API](https://cloud.google.com/resource-manager/reference/rest/v1/projects/getAncestry) for more details.

## Example Usage

```hcl
data "google_project_ancestry" "example" {
  project_id = "example-project"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Optional) The ID of the project. If it is not provided, the provider project is used.

## Attributes Reference

The following attributes are exported:

* `ancestors` - A list of the project's ancestors. Structure is [defined below](#nested_ancestors).

<a name="nested_ancestors"></a>The `ancestors` block supports:

* `id` - If it's a project, the `project_id` is exported, else the numeric folder id or organization id.
* `type` - One of `"project"`, `"folder"` or `"organization"`.

---

* `org_id` - The optional user-assigned display name of the project.
* `parent_id` - The parent's id.
* `parent_type` - One of `"folder"` or `"organization"`.
