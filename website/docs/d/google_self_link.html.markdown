---
layout: "google"
page_title: "Google: google_self_link"
sidebar_current: "docs-google-datasource-self-link"
description: |-
  Terraform-native interpolation of a Google Cloud Platform self link
---

# google\_self\_link

This datasource allows Terraform-native interpolation of a Google Cloud Platform project-level self link's `name` and relative/partial URI (`relative_uri`).

## Example Usage - Basic

```hcl
data "google_self_link" "instance" {
  self_link = "https://www.googleapis.com/compute/v1/projects/my-gcp-project/regions/us-central1/instances/my-instance"
}

output "relative_uri" {
  value = "${data.google_self_link.instance.relative_uri}"
}
```

## Example Usage - with Instance

```hcl
data "google_self_link" "instance" {
  self_link = "${google_compute_instance.vm_instance.self_link}"
}

output "name" {
  value = "${data.google_self_link.instance.name}"
}

resource "google_compute_instance" "vm_instance" {
	name                      = "vm-instance"
	machine_type              = "f1-micro"
	zone                      = "us-central1-a"
	allow_stopping_for_update = true

	boot_disk {
		initialize_params{
			image = "${data.google_compute_image.debian_image.self_link}"
		}
	}

	network_interface {
		network = "default"
		access_config {
		}
	}
}

data "google_compute_image" "debian_image" {
	family  = "debian-9"
	project = "debian-cloud"
}
```

## Argument Reference

The following arguments are supported:

* `self_link` (Required) - The self link you are interpolating.

## Attributes Reference

The following attributes are exported:

* `relative_uri` - The self link's relative/partial URI (starting from `projects/`)

* `name` - The name of the resource the self link refers to.
