---
subcategory: "Compute Engine"
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

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `{{key}}``

## Import

Project metadata items can be imported using the `key`, e.g.

* `{{key}}`

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import project metadata items using one of the formats above. For example:

```tf
import {
  id = "{{key}}"
  to = google_compute_project_metadata_item.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), project metadata items can be imported using one of the formats above. For example:

```
$ terraform import google_compute_project_metadata_item.default {{key}}
```

## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options: configuration options:

- `create` - Default is 5 minutes.
- `update` - Default is 5 minutes.
- `delete` - Default is 5 minutes.
