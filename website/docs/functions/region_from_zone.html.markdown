---
page_title: region_from_zone Function - terraform-provider-google
description: |-
  Returns the region within a provided zone.
---

# Function: region_from_zone

Returns a region name derived from a provided zone.

For more information about using provider-defined functions with Terraform [see the official documentation](https://developer.hashicorp.com/terraform/plugin/framework/functions/concepts).

## Example Usage

### Use with the `google` provider

```terraform
terraform {
  required_providers {
    google = {
      source = "hashicorp/google"
    }
  }
}

resource "google_compute_instance" "default" {
  name         = "my-instance"
  machine_type = "n2-standard-2"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
      labels = {
        my_label = "value"
      }
    }
  }

  network_interface {
    network    = "default"
    subnetwork = google_compute_subnetwork.default.id
    access_config {
      // Ephemeral public IP
    }
  }

  metadata_startup_script = "echo hi > /test.txt"
}

data "google_compute_network" "default" {
  name = "default"
}

resource "google_compute_subnetwork" "default" {
  name          = "my-subnet"
  region        = "us-central1"
  network       = data.google_compute_network.default.id
  ip_cidr_range = "192.168.10.0/24"
}

// The region_from_zone function is used to assert that the VM and subnet are in the same region
check "vm_subnet_compatibility_check" {
  assert {
    condition     = google_compute_subnetwork.default.region == provider::google::region_from_zone(google_compute_instance.default.zone)
    error_message = "Subnet ${google_compute_subnetwork.default.id} and VM ${google_compute_instance.default.id} are not in the same region"
  }
}
```

### Use with the `google-beta` provider

```terraform
terraform {
  required_providers {
    google-beta = {
      source = "hashicorp/google-beta"
    }
  }
}

resource "google_compute_instance" "default" {
  provider     = google-beta
  name         = "my-instance"
  machine_type = "n2-standard-2"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
      labels = {
        my_label = "value"
      }
    }
  }

  network_interface {
    network    = "default"
    subnetwork = google_compute_subnetwork.default.id
    access_config {
      // Ephemeral public IP
    }
  }

  metadata_startup_script = "echo hi > /test.txt"
}

data "google_compute_network" "default" {
  provider   = google-beta
  name       = "default"
}

resource "google_compute_subnetwork" "default" {
  provider      = google-beta
  name          = "my-subnet"
  region        = "us-central1"
  network       = data.google_compute_network.default.id
  ip_cidr_range = "192.168.10.0/24"
}

// The region_from_zone function is used to assert that the VM and subnet are in the same region
check "vm_subnet_compatibility_check" {
  assert {
    condition     = google_compute_subnetwork.default.region == provider::google-beta::region_from_zone(google_compute_instance.default.zone)
    error_message = "Subnet ${google_compute_subnetwork.default.id} and VM ${google_compute_instance.default.id} are not in the same region"
  }
}
```

## Signature

```text
region_from_zone(zone string) string
```

## Arguments

1. `zone` (String) A string of a resource's zone
