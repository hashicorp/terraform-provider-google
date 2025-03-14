---
subcategory: "Compute Engine"
description: |-
  Get information about Google Compute Images.
---

# google_compute_images

Get information about Google Compute Images. Check that your service account has the `compute.imageUser` role if you want to share [custom images](https://cloud.google.com/compute/docs/images/sharing-images-across-projects) from another project. If you want to use [public images][pubimg], do not forget to specify the dedicated project. For more information see
[the official documentation](https://cloud.google.com/compute/docs/images) and its [API](https://cloud.google.com/compute/docs/reference/latest/images).

## Example Usage

```hcl
data "google_compute_images" "debian" {
  filter = "name eq my-image.*"
}

resource "google_compute_instance" "default" {
  name         = "test"
  machine_type = "f1-micro"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = data.google_compute_images.debian.images[0].self_link
    }
  }

  network_interface {
    network = google_compute_network.default.name
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` -Filter for the images to be returned by the data source. Syntax can be found [here](https://cloud.google.com/compute/docs/reference/rest/v1/images/list) in the filter section.

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
* `source_image_id` - The ID value of the image used to create this image.
* `source_disk` - The URL of the source disk used to create this image.
* `source_disk_id` - The ID value of the disk used to create this image.
* `creation_timestamp` - The creation timestamp in RFC3339 text format.
* `description` - An optional description of this image.
* `labels` - All of labels (key/value pairs) present on the resource in GCP, including the labels configured through Terraform, other clients and services.

[pubimg]: https://cloud.google.com/compute/docs/images#os-compute-support "Google Cloud Public Base Images"
