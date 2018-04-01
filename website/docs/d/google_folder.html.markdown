---
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
  folder = "folders/12345"
  lookup_organization = true
}

# Search by fields
data "google_folder" "my_folder_2" {
  parent = "organizations/23456"
  lifecycle_state = "ACTIVE"
  display_name = "test folder"
}

output "my_folder_1_organization" {
  value = "${data.google_folder.my_folder_1.organization}"
}

output "my_folder_2_name" {
  value = "${data.google_folder.my_folder_2.name}"
}

```

## Argument Reference

The arguments of this data source act as filters for querying the available Folders.
The given filters must match exactly one Folders whose data will be exported as attributes.
The following arguments are supported:

* `folder` (Optional) - The name of the Folder in the form `{folder_id}` or `folders/{folder_id}`.

* `name` (Optional) - The resource name of the Folder in the form `folders/{folder_id}`. Example: "folders/343562343".
* `parent` (Optional) - The resource name of the parent Folder or Organization. Example: "organizations/{organization_id}" or "folders/{folder_id}".
* `display_name` (Optional) - The folder's display name. Example: "test folder".
* `lifecycle_state` (Optional) - The Folder's current lifecycle state. Example: "ACTIVE".

* `lookup_organization` (Optional) - `true` to find the organization that the folder belongs, `false` to avoid the lookup. It searches up the tree.

~> **NOTE:** One of `folder` or {`name` or `parent` or `display_name` or `lifecycle_state`} must be specified, the `folder` is to get by id and the other group is to search using Query String which support wildcards.

## Attributes Reference

The following additional attributes are exported:

* `id` - The Folder ID.
* `name` - The resource name of the Folder in the form `folders/{organization_id}`.
* `parent` - The resource name of the parent Folder or Organization.
* `display_name` - The folder's display name.
* `create_time` - Timestamp when the Organization was created. A timestamp in RFC3339 UTC "Zulu" format, accurate to nanoseconds. Example: "2014-10-02T15:01:23.045123456Z".
* `lifecycle_state` - The Folder's current lifecycle state.
* `organization` - If `lookup_organization` is enable, the resource name of the Organization that the folder belongs.
