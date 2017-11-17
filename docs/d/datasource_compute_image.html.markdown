---
layout: "google"
page_title: "Google: google_compute_image"
sidebar_current: "docs-google-datasource-compute-image"
description: |-
  Get information about a Google Compute Image.
---

# google\_compute\_image

Get information about a Google Compute Image. Check that your service account has the `compute.imageUser` role if you want to share [custom images](https://cloud.google.com/compute/docs/images/sharing-images-across-projects) from another project. If you want to use [public images](https://cloud.google.com/compute/docs/images#os-compute-support), do not forget to specify the dedicated project. For more information see
[the official documentation](https://cloud.google.com/compute/docs/images) and its [API](https://cloud.google.com/compute/docs/reference/latest/images).

## Example Usage

```hcl
data "google_compute_image" "my_image" {
  name = "image-family"
}

resource "google_compute_instance" "default" {
  name         = "test"
  machine_type = "n1-standard-1"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = "${data.google_compute_image.my_image.self_link}"
    }
  }

  network_interface {
    network = "default"

    access_config {
      // Ephemeral IP
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` or `family` - (Required) The name of a specific image or a family.
Exactly one of `name` of `family` must be specified. If `name` is specified, it will fetch
the corresponding image. If `family` is specified, it will returns the latest image
that is part of an image family and is not deprecated.

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `self_link` - The URI of the resource.
* `name` - The image name of the resource in case `family` was specified.
* `family` - The family name of the resource in case `name` was specified.
* `disk_size_gb` - The size of the image when restored onto a persistent disk in gigabytes.
* `archive_size_bytes` - The size of the image tar.gz archive stored in Google Cloud Storage in bytes.
* `image_id` - The unique identifier for the resource.
* `image_encryption_key_sha256` - The [RFC 4648 base64](https://tools.ietf.org/html/rfc4648#section-4)
    encoded SHA-256 hash of the [customer-supplied encryption key](https://cloud.google.com/compute/docs/disks/customer-supplied-encryption)
    that protects this resource.
* `source_image_id` - The ID value of the image used to create this image.
* `source_disk` - The URL of the source disk used to create this image.
* `source_disk_encryption_key_sha256` - The [RFC 4648 base64](https://tools.ietf.org/html/rfc4648#section-4)
    encoded SHA-256 hash of the [customer-supplied encryption key](https://cloud.google.com/compute/docs/disks/customer-supplied-encryption)
    that protects this resource.
* `source_disk_id` - The ID value of the disk used to create this image.
* `creation_timestamp` - The creation timestamp in RFC3339 text format.
* `description` - An optional description of this resource.
* `labels` - A map of labels applied to this image.
* `label_fingerprint` - A fingerprint for the labels being applied to this image.
* `licenses` - A list of applicable license URI.
* `status` - The status of the image. Possible values are **FAILED**, **PENDING**, or **READY**.
