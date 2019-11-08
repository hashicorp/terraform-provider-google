---
subcategory: "Compute Engine"
layout: "google"
page_title: "Google: google_compute_project_metadata_item"
sidebar_current: "docs-google-compute-project-metadata-item"
description: |-
  Manages a single key/value pair on common instance metadata
---

# google\_compute\_project\_metadata\_item

Manages a single key/value pair on metadata common to all instances for
a project in GCE. Using `google_compute_project_metadata_item` lets you
manage a single key/value setting in Terraform rather than the entire
project metadata map.

## Example Usage

```hcl
resource "google_compute_project_metadata_item" "default" {
  key   = "my_metadata"
  value = "my_value"
}
```

## Argument Reference

The following arguments are supported:

* `key` - (Required) The metadata key to set.

* `value` - (Required) The value to set for the given metadata key.

- - -

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

Only the arguments listed above are exposed as attributes.

## Import

Project metadata items can be imported using the `key`, e.g.

```
$ terraform import google_compute_project_metadata_item.default my_metadata
```

## Timeouts

This resource provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 5 minutes.
- `update` - Default is 5 minutes.
- `delete` - Default is 5 minutes.
