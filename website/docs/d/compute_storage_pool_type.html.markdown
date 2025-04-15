---
subcategory: "Compute Engine"
description: |-
  Represents a StoragePoolType data source.
---

# google_compute_storage_pool_type

Represents a StoragePoolType data source.
The type of Hyperdisk Storage Pool that you create determines the type of disks that you can create in the storage pool.


To get more information about StoragePoolType, see:

* [API documentation](https://cloud.google.com/compute/docs/reference/rest/v1/storagePoolTypes)
* How-to Guides
    * [Types of Hyperdisk Storage Pools](https://cloud.google.com/compute/docs/disks/storage-pools#sp-types)

## Argument Reference

The following arguments are supported:


* `zone` -
  (Required)
  The name of the zone.

* `storage_pool_type` -
  (Required)
  Name of the storage pool type.


- - -


* `project` - (Optional) The ID of the project in which the resource belongs.
    If it is not provided, the provider project is used.


## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/zones/{{zone}}/storagePoolTypes/{{storagePoolType}}`

* `kind` -
  Type of the resource. Always compute#storagePoolType for storage pool types.

* `id` -
  The unique identifier for the resource. This identifier is defined by the server.

* `creation_timestamp` -
  Creation timestamp in RFC3339 text format.

* `name` -
  Name of the resource.

* `description` -
  An optional description of this resource.

* `deprecated` -
  The deprecation status associated with this storage pool type.
  Structure is [documented below](#nested_deprecated).

* `self_link` -
  Server-defined URL for the resource.

* `self_link_with_id` -
  Server-defined URL for this resource with the resource id.

* `min_pool_provisioned_capacity_gb` -
  Minimum storage pool size in GB.

* `max_pool_provisioned_capacity_gb` -
  Maximum storage pool size in GB.

* `min_pool_provisioned_iops` -
  Minimum provisioned IOPS.

* `max_pool_provisioned_iops` -
  Maximum provisioned IOPS.

* `min_pool_provisioned_throughput` -
  Minimum provisioned throughput.

* `max_pool_provisioned_throughput` -
  Maximum provisioned throughput.

* `supported_disk_types` -
  The list of disk types supported in this storage pool type.


<a name="nested_deprecated"></a>The `deprecated` block contains:

* `state` -
  (Output)
  The deprecation state of this resource. This can be ACTIVE, DEPRECATED, OBSOLETE, or DELETED.
  Operations which communicate the end of life date for an image, can use ACTIVE.
  Operations which create a new resource using a DEPRECATED resource will return successfully,
  but with a warning indicating the deprecated resource and recommending its replacement.
  Operations which use OBSOLETE or DELETED resources will be rejected and result in an error.

* `replacement` -
  (Output)
  The URL of the suggested replacement for a deprecated resource.
  The suggested replacement resource must be the same kind of resource as the deprecated resource.

* `deprecated` -
  (Output)
  An optional RFC3339 timestamp on or after which the state of this resource is intended to change to DEPRECATED.
  This is only informational and the status will not change unless the client explicitly changes it.

* `obsolete` -
  (Output)
  An optional RFC3339 timestamp on or after which the state of this resource is intended to change to OBSOLETE.
  This is only informational and the status will not change unless the client explicitly changes it.

* `deleted` -
  (Output)
  An optional RFC3339 timestamp on or after which the state of this resource is intended to change to DELETED.
  This is only informational and the status will not change unless the client explicitly changes it.

## User Project Overrides

This resource supports [User Project Overrides](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/provider_reference#user_project_override).
