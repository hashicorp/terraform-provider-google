---
layout: "google"
page_title: "Google: google_compute_attached_disk"
sidebar_current: "docs-google-compute-attached-disk"
description: |-
  Resource that allows attaching existing persistent disks to compute instances.
---

# google\_compute\_attached\_disk

Persistent disks can be attached to a compute instance using [the `attach_disk`
section within the compute instance configuration](https://www.terraform.io/docs/providers/google/r/compute_instance.html#attached_disk).
However there may be situations where managing the attached disks via the compute
instance config isn't preferable or possible. For example: attaching dynamic
numbers of disks using the `count` variable.


To get more information about attaching disks, see:

* [API documentation](https://cloud.google.com/compute/docs/reference/rest/v1/instances/attachDisk)
* [Resource: google_compute_disk](https://www.terraform.io/docs/providers/google/r/compute_disk.html)
* How-to Guides
    * [Adding a persistent disk](https://cloud.google.com/compute/docs/disks/add-persistent-disk)


## Example Usage
```hcl
resource "google_compute_attached_disk" "default" {
  disk = "${google_compute_disk.default.self_link}"
  instance = "${google_compute_instance.default.self_link}"
}
```

## Argument Reference

The following arguments are supported:


* `instance` -
  (Required)
  `name` or `self_link` of the compute instance that the disk will be attached to.
  If the `self_link` is provided then `zone` and `project` are extracted from the
  self link. If only the name is used then `zone` and `project` must be defined
  at a global level.

* `disk` -
  (Required)
  `name` or `self_link` of the disk that will be attached.



## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:


* `project` -
  The project that the referenced compute instance is a part of.

* `zone` -
  The zone that the referenced compute instance is located within.

## Timeouts

This resource provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 5 minutes.
- `delete` - Default is 5 minutes.

## Import

Attached Disk is composed of a compute instance and disk and it's `id` is the
names of both concatenated together with a colon `:`. It can be imported the
following ways:

```
$ terraform import google_compute_disk.default projects/{{project}}/zones/{{zone}}/disks/{{instance.name}}:{{disk.name}}
$ terraform import google_compute_disk.default {{project}}/{{zone}}/{{instance.name}}:{{disk.name}}
```
