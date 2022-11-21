---
subcategory: "Memorystore (Redis)"
page_title: "Google: google_redis_instance"
description: |-
  Get information about a Google Cloud Redis instance.
---

# google\_redis\_instance

Get info about a Google Cloud Redis instance.

## Example Usage

```tf
data "google_redis_instance" "my_instance" {
  name = "my-redis-instance"
}

output "instance_memory_size_gb" {
  value = data.google_redis_instance.my_instance.memory_size_gb
}

output "instance_connect_mode" {
  value = data.google_redis_instance.my_instance.connect_mode
}

output "instance_authorized_network" {
  value = data.google_redis_instance.my_instance.authorized_network
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of a Redis instance.

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

* `region` - (Optional) The region in which the resource belongs. If it
    is not provided, the provider region is used.

## Attributes Reference

See [google_redis_instance](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/redis_instance) resource for details of the available attributes.
