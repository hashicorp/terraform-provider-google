---
subcategory: "Compute Engine"
description: |-
  List networks in a Google Cloud project.
---

# google_compute_networks

List all networks in a specified Google Cloud project.

## Example Usage

```tf
data "google_compute_networks" "my-networks" {
  project = "my-cloud-project"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Optional) The name of the project.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `id` - an identifier for the resource with format projects/{{project}}/global/networks

* `networks` - The list of networks in the specified project.

* `project` - The project name being queried.

* `self_link` - The URI of the resource.
