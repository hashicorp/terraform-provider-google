---
layout: "google"
page_title: "Google: google_compute_image"
sidebar_current: "docs-google-compute-image"
description: |-
  Creates a bootable VM image for Google Compute Engine from an existing tarball.
---

# google\_compute\_image

Creates a bootable VM image resource for Google Compute Engine from an existing
tarball. For more information see [the official documentation](https://cloud.google.com/compute/docs/images) and
[API](https://cloud.google.com/compute/docs/reference/latest/images).


## Example Usage

```hcl
resource "google_compute_image" "bootable-image" {
  name = "my-custom-image"

  raw_disk {
    source = "https://storage.googleapis.com/my-bucket/my-disk-image-tarball.tar.gz"
  }

  licenses = [
    "https://www.googleapis.com/compute/v1/projects/vm-options/global/licenses/enable-vmx",
  ]
}

resource "google_compute_instance" "vm" {
  name         = "vm-from-custom-image"
  machine_type = "n1-standard-1"
  zone         = "us-east1-c"

  boot_disk {
    initialize_params {
      image = "${google_compute_image.bootable-image.self_link}"
    }
  }

  network_interface {
    network = "default"
  }
}
```

## Argument Reference

The following arguments are supported: (Note that one of either source_disk or
  raw_disk is required)

* `name` - (Required) A unique name for the resource, required by GCE.
    Changing this forces a new resource to be created.

- - -

* `description` - (Optional) The description of the image to be created

* `family` - (Optional) The name of the image family to which this image belongs.

* `labels` - (Optional) A set of key/value label pairs to assign to the image.

* `source_disk` - (Optional) The URL of a disk that will be used as the source of the
    image. Changing this forces a new resource to be created.

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

* `raw_disk` - (Optional) The raw disk that will be used as the source of the image.
    Changing this forces a new resource to be created. Structure is documented
    below.

* `licenses` - (Optional) A list of license URIs to apply to this image. Changing this
    forces a new resource to be created.

The `raw_disk` block supports:

* `source` - (Required) The full Google Cloud Storage URL where the disk
    image is stored.

* `sha1` - (Optional) SHA1 checksum of the source tarball that will be used
    to verify the source before creating the image.

* `container_type` - (Optional) The format used to encode and transmit the
    block device. TAR is the only supported type and is the default.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `self_link` - The URI of the created resource.

* `label_fingerprint` - The fingerprint of the assigned labels.

## Timeouts

`google_compute_image` provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default `4 minutes`
- `update` - Default `4 minutes`
- `delete` - Default `4 minutes`

## Import

VM image can be imported using the `name`, e.g.

```
$ terraform import google_compute_image.web-image my-custom-image
```
