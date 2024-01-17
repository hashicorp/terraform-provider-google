---
subcategory: "Cloud Platform"
description: |-
  Retrieve a set of folders based on a parent ID.
---

# google\_folders

Retrieve information about a set of folders based on a parent ID. See the
[REST API](https://cloud.google.com/resource-manager/reference/rest/v3/folders/list)
for more details.

## Example Usage - searching for folders at the root of an org

```hcl
data "google_folders" "my-org-folders" {
  parent_id = "organizations/${var.organization_id}"
}

data "google_folder" "first-folder" {
  folder = data.google_folders.my-org-folders.folders[0].name
}
```

## Argument Reference

The following arguments are supported:

* `parent_id` - (Required) A string parent as defined in the [REST API](https://cloud.google.com/resource-manager/reference/rest/v3/folders/list#query-parameters).


## Attributes Reference

The following attributes are exported:

* `folders` - A list of folders matching the provided filter. Structure is defined below.

The `folders` block supports:

* `name` - The id of the folder
* `parent` - The parent id of the folder
* `display_name` - The display name of the folder
* `state` - The lifecycle state of the folder
* `create_time` - The timestamp of when the folder was created
* `update_time` - The timestamp of when the folder was last modified
* `delete_time` - The timestamp of when the folder was requested to be deleted (if applicable)
* `etag` - Entity tag identifier of the folder

