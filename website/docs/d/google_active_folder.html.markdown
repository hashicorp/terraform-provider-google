---
subcategory: "Cloud Platform"
layout: "google"
page_title: "Google: google_active_folder"
sidebar_current: "docs-google-datasource-active-folder"
description: |-
  Get a folder within GCP.
---

# google\_active\_folder

Get an active folder within GCP by `display_name` and `parent`.

## Example Usage

```tf
data "google_active_folder" "department1" {
  display_name = "Department 1"
  parent       = "organizations/1234567"
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) The folder's display name.

* `parent` - (Required) The resource name of the parent Folder or Organization.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `name` - The resource name of the Folder. This uniquely identifies the folder.
