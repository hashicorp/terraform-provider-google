---
subcategory: "Cloud Platform"
layout: "google"
page_title: "Google: google_folder"
sidebar_current: "docs-google-datasource-folder"
description: |-
  Get information about a Google Cloud Folder.
---

# google\_folder

Use this data source to get information about a Google Cloud Folder.

```hcl
# Get folder by id
data "google_folder" "my_folder_1" {
  folder              = "folders/12345"
  lookup_organization = true
}

# Search by fields
data "google_folder" "my_folder_2" {
  folder = "folders/23456"
}

output "my_folder_1_organization" {
  value = data.google_folder.my_folder_1.organization
}

output "my_folder_2_parent" {
  value = data.google_folder.my_folder_2.parent
}
```

## Argument Reference

The following arguments are supported:

* `folder` (Required) - The name of the Folder in the form `{folder_id}` or `folders/{folder_id}`.
* `lookup_organization` (Optional) - `true` to find the organization that the folder belongs, `false` to avoid the lookup. It searches up the tree. (defaults to `false`)

## Attributes Reference

The following attributes are exported:

* `id` - The Folder ID.
* `name` - The resource name of the Folder in the form `folders/{folder_id}`.
* `parent` - The resource name of the parent Folder or Organization.
* `display_name` - The folder's display name.
* `create_time` - Timestamp when the Organization was created. A timestamp in RFC3339 UTC "Zulu" format, accurate to nanoseconds. Example: "2014-10-02T15:01:23.045123456Z".
* `lifecycle_state` - The Folder's current lifecycle state.
* `organization` - If `lookup_organization` is enable, the resource name of the Organization that the folder belongs.
