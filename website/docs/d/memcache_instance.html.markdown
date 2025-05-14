---
subcategory: "Memcache"
description: |-
  Fetches the details of available instance.
---

# google_memcache_instance

Use this data source to get information about the available instance. For more details refer the [API docs](https://cloud.google.com/memorystore/docs/memcached/reference/rest/v1/projects.locations.instances).

## Example Usage


```hcl
data "google_memcache_instance" "qa" {
}
```

## Argument Reference

The following arguments are supported:


* `name` -
  (Required)
  The ID of the memcache instance.
  'memcache_instance_id'

* `project` - 
  (optional) 
  The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

* `region` -
  (optional)
  The canonical id of the region. If it is not provided, the provider project is used. For example: us-east1.

## Attributes Reference

See [google_memcache_instance](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/memcache_instance) resource for details of all the available attributes.