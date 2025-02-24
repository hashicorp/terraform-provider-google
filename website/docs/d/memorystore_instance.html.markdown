---
subcategory: "Memorystore"
description: |-
  Fetches the details of available instance.
---

# google_memorystore_instance

Use this data source to get information about the available instance. For more details refer the [API docs](https://cloud.google.com/memorystore/docs/valkey/reference/rest/v1/projects.locations.instances).

## Example Usage


```hcl
data "google_memorystore_instance" "qa" {
}
```

## Argument Reference

The following arguments are supported:


* `instance_id` -
  (Required)
  The ID of the memorystore instance.
  'memorystore_instance_id'

* `project` - 
  (optional) 
  The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

* `location` -
  (optional)
  The canonical id of the location.If it is not provided, the provider project is used. For example: us-east1.

## Attributes Reference

See [google_memorystore_instance](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/memorystore_instance) resource for details of all the available attributes.
