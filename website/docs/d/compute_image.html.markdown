---
subcategory: "Compute Engine"
layout: "google"
page_title: "Google: google_compute_image"
sidebar_current: "docs-google-datasource-compute-image"
description: |-
  Get information about a Google Compute Image.
---

# google\_compute\_image

Get information about a Google Compute Image. Check that your service account has the `compute.imageUser` role if you want to share [custom images](https://cloud.google.com/compute/docs/images/sharing-images-across-projects) from another project. If you want to use [public images][pubimg], do not forget to specify the dedicated project. For more information see
[the official documentation](https://cloud.google.com/compute/docs/images) and its [API](https://cloud.google.com/compute/docs/reference/latest/images).

## Example Usage

```hcl
data "google_compute_image" "my_image" {
  family  = "debian-9"
  project = "debian-cloud"
}

resource "google_compute_instance" "default" {
  # ...

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
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

* `project` - (Optional) The project in which the resource belongs. If it is not
  provided, the provider project is used. If you are using a
  [public base image][pubimg], be sure to specify the correct Image Project.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `self_link` - The URI of the image.
* `name` - The name of the image.
* `family` - The family name of the image.
* `disk_size_gb` - The size of the image when restored onto a persistent disk in gigabytes.
* `archive_size_bytes` - The size of the image tar.gz archive stored in Google Cloud Storage in bytes.
* `image_id` - The unique identifier for the image.
* `image_encryption_key_sha256` - The [RFC 4648 base64](https://tools.ietf.org/html/rfc4648#section-4)
    encoded SHA-256 hash of the [customer-supplied encryption key](https://cloud.google.com/compute/docs/disks/customer-supplied-encryption)
    that protects this image.
* `source_image_id` - The ID value of the image used to create this image.
* `source_disk` - The URL of the source disk used to create this image.
* `source_disk_encryption_key_sha256` - The [RFC 4648 base64](https://tools.ietf.org/html/rfc4648#section-4)
    encoded SHA-256 hash of the [customer-supplied encryption key](https://cloud.google.com/compute/docs/disks/customer-supplied-encryption)
    that protects this image.
* `source_disk_id` - The ID value of the disk used to create this image.
* `creation_timestamp` - The creation timestamp in RFC3339 text format.
* `description` - An optional description of this image.
* `labels` - A map of labels applied to this image.
* `label_fingerprint` - A fingerprint for the labels being applied to this image.
* `licenses` - A list of applicable license URI.
* `status` - The status of the image. Possible values are **FAILED**, **PENDING**, or **READY**.

[pubimg]: https://cloud.google.com/compute/docs/images#os-compute-support "Google Cloud Public Base Images"
