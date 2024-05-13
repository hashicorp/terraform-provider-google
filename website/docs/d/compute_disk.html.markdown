---
subcategory: "Compute Engine"
description: |-
  Get information about a Google Compute Persistent disks.
---

# google_compute_disk

Get information about a Google Compute Persistent disks.

[the official documentation](https://cloud.google.com/compute/docs/disks) and its [API](https://cloud.google.com/compute/docs/reference/latest/disks).

## Example Usage

```hcl
data "google_compute_disk" "persistent-boot-disk" {
  name    = "persistent-boot-disk"
  project = "example"
}

resource "google_compute_instance" "default" {
  # ...
    
  boot_disk {
    source = data.google_compute_disk.persistent-boot-disk.self_link
    auto_delete = false
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of a specific disk.

- - -

* `zone` - (Optional) A reference to the zone where the disk resides.

* `project` - (Optional) The ID of the project in which the resource belongs.
    If it is not provided, the provider project is used.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/zones/{{zone}}/disks/{{name}}`

* `label_fingerprint` -
  The fingerprint used for optimistic locking of this resource.  Used
  internally during updates.

* `creation_timestamp` -
  Creation timestamp in RFC3339 text format.

* `last_attach_timestamp` -
  Last attach timestamp in RFC3339 text format.

* `last_detach_timestamp` -
  Last detach timestamp in RFC3339 text format.

* `users` -
  Links to the users of the disk (attached instances) in form:
  project/zones/zone/instances/instance

* `source_image_id` -
  The ID value of the image used to create this disk. This value
  identifies the exact image that was used to create this persistent
  disk. For example, if you created the persistent disk from an image
  that was later deleted and recreated under the same name, the source
  image ID would identify the exact version of the image that was used.

* `source_snapshot_id` -
  The unique ID of the snapshot used to create this disk. This value
  identifies the exact snapshot that was used to create this persistent
  disk. For example, if you created the persistent disk from a snapshot
  that was later deleted and recreated under the same name, the source
  snapshot ID would identify the exact version of the snapshot that was
  used.

* `description` -
  The optional description of this resource.

* `labels` - All of labels (key/value pairs) present on the resource in GCP, including the labels configured through Terraform, other clients and services.

* `size` -
  Size of the persistent disk, specified in GB.

* `physical_block_size_bytes` -
  Physical block size of the persistent disk, in bytes.

* `type` -
  URL of the disk type resource describing which disk type to use to
  create the disk.

* `image` -
  The image from which to initialize this disk.

* `zone` -
  A reference to the zone where the disk resides.

* `source_image_encryption_key` -
  The customer-supplied encryption key of the source image.

* `snapshot` -
  The source snapshot used to create this disk.

* `source_snapshot_encryption_key` -
  (Optional)
  The customer-supplied encryption key of the source snapshot.

* `self_link` - The URI of the created resource.
