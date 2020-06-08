---
subcategory: "Memorystore (Redis)"
layout: "google"
page_title: "Google: google_redis_instance"
sidebar_current: "docs-google-datasource-redis-instance"
description: |-
  Get information about a Google Cloud Redis instance.
---

# google\_redis\_instance

Get information about a Google Cloud Redis instance. For more information see
the [official documentation](https://cloud.google.com/memorystore/docs/redis)
and [API](https://cloud.google.com/memorystore/docs/redis/apis).

## Example Usage

```hcl
data "google_redis_instance" "default" {
  name = "my-redis-instance"
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

The following attributes are exported:

* `host` - Hostname or IP address of the exposed Redis endpoint used by clients
  to connect to the service.

* `port` - The port number of the exposed Redis endpoint.
