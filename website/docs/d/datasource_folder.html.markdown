---
layout: "google"
page_title: "Google: google_active_folder"
sidebar_current: "docs-google-datasource-folder"
description: |-
  Get a folder within GCP.
---

# google\_active\_folder

Get a folder within GCP by `display_name` and `parent`.

## Example Usage

```tf
resource "google_folder" "new-folder" {
  display_name = "new-folder"
  parent = "folders/some-folder-id"
}

data "google_active_folder" "new-folder" {
  display_name = "${google_folder.new-folder.display_name}"
  parent = "${google_folder.new-folder.parent}"
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) The folder's display name.

* `parent` - (Required) The resource name of the parent Folder or Organization.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `name` - The resource name of the Folder. This uniquely identifies the folder.
