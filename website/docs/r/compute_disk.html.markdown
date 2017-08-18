---
layout: "google"
page_title: "Google: google_compute_disk"
sidebar_current: "docs-google-compute-disk"
description: |-
  Creates a new persistent disk within GCE, based on another disk.
---

# google\_compute\_disk

Creates a new persistent disk within GCE, based on another disk. For more information see
[the official documentation](https://cloud.google.com/compute/docs/disks/add-persistent-disk)
and
[API](https://cloud.google.com/compute/docs/reference/latest/disks).

~> **Note:** All arguments including the disk encryption key will be stored in the raw state as plain-text.
[Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Example Usage

```hcl
resource "google_compute_disk" "default" {
  name  = "test-disk"
  type  = "pd-ssd"
  zone  = "us-central1-a"
  image = "debian-8-jessie-v20170523"
  labels {
    environment = "dev"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A unique name for the resource, required by GCE.
    Changing this forces a new resource to be created.

* `zone` - (Required) The zone where this disk will be available.

- - -

* `disk_encryption_key_raw` - (Optional) A 256-bit [customer-supplied encryption key]
    (https://cloud.google.com/compute/docs/disks/customer-supplied-encryption),
    encoded in [RFC 4648 base64](https://tools.ietf.org/html/rfc4648#section-4)
    to encrypt this disk.

* `image` - (Optional) The image from which to initialize this disk. This can be
    one of: the image's `self_link`, of a full name and version, e.g.
    `debian-8-jessie-v20170523`

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

* `size` - (Optional) The size of the image in gigabytes. If not specified, it
    will inherit the size of its base image.

* `snapshot` - (Optional) Name of snapshot from which to initialize this disk.

* `type` - (Optional) The GCE disk type.

* `labels` - (Optional) A set of key/value label pairs to assign to the image.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `disk_encryption_key_sha256` - The [RFC 4648 base64]
    (https://tools.ietf.org/html/rfc4648#section-4) encoded SHA-256 hash of the
    [customer-supplied encryption key](https://cloud.google.com/compute/docs/disks/customer-supplied-encryption)
    that protects this resource.

* `self_link` - The URI of the created resource.

* `users` - The Users of the created resource.

* `label_fingerprint` - The fingerprint of the assigned labels.

## Import

Disks can be imported using the `name`, e.g.

```
$ terraform import google_compute_disk.default test-disk
```
