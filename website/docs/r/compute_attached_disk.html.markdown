---
subcategory: "Compute Engine"
description: |-
  Resource that allows attaching existing persistent disks to compute instances.
---

# google_compute_attached_disk

Persistent disks can be attached to a compute instance using [the `attached_disk`
section within the compute instance configuration](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_instance#attached_disk).
However there may be situations where managing the attached disks via the compute
instance config isn't preferable or possible, such as attaching dynamic
numbers of disks using the `count` variable.


To get more information about attaching disks, see:

* [API documentation](https://cloud.google.com/compute/docs/reference/rest/v1/instances/attachDisk)
* [Resource: google_compute_disk](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_disk)
* How-to Guides
    * [Adding a persistent disk](https://cloud.google.com/compute/docs/disks/add-persistent-disk)

**Note:** When using `google_compute_attached_disk` you **must** use `lifecycle.ignore_changes = ["attached_disk"]` on the `google_compute_instance` resource that has the disks attached. Otherwise the two resources will fight for control of the attached disk block.

## Example Usage
```hcl
resource "google_compute_attached_disk" "default" {
  disk     = google_compute_disk.default.id
  instance = google_compute_instance.default.id
}

resource "google_compute_instance" "default" {
  name         = "attached-disk-instance"
  machine_type = "e2-medium"
  zone         = "us-west1-a"

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
    }
  }

  network_interface {
    network = "default"
  }

  lifecycle {
    ignore_changes = [attached_disk]
  }
}
```

## Argument Reference

The following arguments are supported:


* `instance` -
  (Required)
  `name` or `self_link` of the compute instance that the disk will be attached to.
  If the `self_link` is provided then `zone` and `project` are extracted from the
  self link. If only the name is used then `zone` and `project` must be defined
  as properties on the resource or provider.

* `disk` -
  (Required)
  `name` or `self_link` of the disk that will be attached.


- - -

* `project` -
  (Optional)
  The project that the referenced compute instance is a part of. If `instance` is referenced by its
  `self_link` the project defined in the link will take precedence.

* `zone` -
  (Optional)
  The zone that the referenced compute instance is located within. If `instance` is referenced by its
  `self_link` the zone defined in the link will take precedence.

* `device_name` -
  (Optional)
  Specifies a unique device name of your choice that is
	reflected into the /dev/disk/by-id/google-* tree of a Linux operating
	system running within the instance. This name can be used to
	reference the device for mounting, resizing, and so on, from within
	the instance.

	If not specified, the server chooses a default device name to apply
	to this disk, in the form persistent-disks-x, where x is a number
	assigned by Google Compute Engine.

* `mode` -
  (Optional)
  The mode in which to attach this disk, either READ_WRITE or
	READ_ONLY. If not specified, the default is to attach the disk in
	READ_WRITE mode.

	Possible values:
	  "READ_ONLY"
	  "READ_WRITE"

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/zones/{{zone}}/disks/{{disk.name}}`

## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options: configuration options:

- `create` - Default is 5 minutes.
- `delete` - Default is 5 minutes.

## Import

Attached Disk can be imported the following ways:

* `projects/{{project}}/zones/{{zone}}/instances/{{instance.name}}/{{disk.name}}`
* `{{project}}/{{zone}}/{{instance.name}}/{{disk.name}}`

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Attached Disk using one of the formats above. For example:

```tf
import {
  id = "projects/{{project}}/zones/{{zone}}/instances/{{instance.name}}/{{disk.name}}"
  to = google_compute_attached_disk.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), Attached Disk can be imported using one of the formats above. For example:



```
$ terraform import google_compute_attached_disk.default projects/{{project}}/zones/{{zone}}/instances/{{instance.name}}/{{disk.name}}
$ terraform import google_compute_attached_disk.default {{project}}/{{zone}}/{{instance.name}}/{{disk.name}}
```
