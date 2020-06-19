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

In addition to the arguments listed above, the following computed attributes are exported:

* `memory_size_gb` -
  Redis memory size in GiB.

* `alternative_location_id` -
  Only applicable to STANDARD_HA tier which protects the instance
  against zonal failures by provisioning it across two zones.
  If provided, it must be a different zone from the one provided in
  [locationId].

* `authorized_network` -
  The full name of the Google Compute Engine network to which the
  instance is connected. If left unspecified, the default network
  will be used.

* `connect_mode` -
  The connection mode of the Redis instance.

* `display_name` -
  An arbitrary and optional user-provided name for the instance.

* `labels` -
  Resource labels to represent user provided metadata.

* `redis_configs` -
  Redis configuration parameters, according to http://redis.io/topics/config.
  Please check Memorystore documentation for the list of supported parameters:
  https://cloud.google.com/memorystore/docs/redis/reference/rest/v1/projects.locations.instances#Instance.FIELDS.redis_configs

* `location_id` -
  The zone where the instance will be provisioned. If not provided,
  the service will choose a zone for the instance. For STANDARD_HA tier,
  instances will be created across two zones for protection against
  zonal failures. If [alternativeLocationId] is also provided, it must
  be different from [locationId].

* `redis_version` -
  The version of Redis software. If not provided, latest supported
  version will be used. Currently, the supported values are:
  - REDIS_4_0 for Redis 4.0 compatibility
  - REDIS_3_2 for Redis 3.2 compatibility

* `reserved_ip_range` -
  The CIDR range of internal addresses that are reserved for this
  instance. If not provided, the service will choose an unused /29
  block, for example, 10.0.0.0/29 or 192.168.0.0/29. Ranges must be
  unique and non-overlapping with existing subnets in an authorized
  network.

* `tier` -
  The service tier of the instance. Must be one of these values:
  - BASIC: standalone instance
  - STANDARD_HA: highly available primary/replica instances

  Default value: `BASIC`
  Possible values are:
  * `BASIC`
  * `STANDARD_HA`

* `host` - Hostname or IP address of the exposed Redis endpoint used by clients
  to connect to the service.

* `port` - The port number of the exposed Redis endpoint.

* `create_time` -
  The time the instance was created in RFC3339 UTC "Zulu" format,
  accurate to nanoseconds.

* `current_location_id` -
  The current zone where the Redis endpoint is placed.
  For Basic Tier instances, this will always be the same as the
  [locationId] provided by the user at creation time. For Standard Tier
  instances, this can be either [locationId] or [alternativeLocationId]
  and can change after a failover event.
