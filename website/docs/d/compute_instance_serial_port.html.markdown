---
subcategory: "Compute Engine"
layout: "google"
page_title: "Google: google_compute_instance_serial_port"
sidebar_current: "docs-google-datasource-compute-instance-serial-port"
description: |-
  Get the serial port output from a Compute Instance.
---

# google\_compute\_instance\_serial\_port

Get the serial port output from a Compute Instance. For more information see
the official [API](https://cloud.google.com/compute/docs/instances/viewing-serial-port-output) documentation.

## Example Usage

```hcl
data "google_compute_instance_serial_port" "serial" {
  instance = "my-instance"
  zone = "us-central1-a"
  port = 1
}

output "serial_out" {
  value = data.google_compute_instance_serial_port.serial.contents
}
```

Using the serial port output to generate a windows password, derived from the [official guide](https://cloud.google.com/compute/docs/instances/windows/automate-pw-generation):

```hcl
resource "google_compute_instance" "windows" {
  name         = "windows-instance"
  machine_type = "n1-standard-1"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = "gce-uefi-images/windows-2019"
    }
  }

  network_interface {
    network = "default"

    access_config {
      // Ephemeral IP
    }
  }

  metadata = {
    serial-port-logging-enable = "TRUE"
    // Derived from https://cloud.google.com/compute/docs/instances/windows/automate-pw-generation
    windows-keys = jsonencode(
      {
        email    = "example.user@example.com"
        expireOn = "2020-04-14T01:37:19Z"
        exponent = "AQAB"
        modulus  = "wgsquN4IBNPqIUnu+h/5Za1kujb2YRhX1vCQVQAkBwnWigcCqOBVfRa5JoZfx6KIvEXjWqa77jPvlsxM4WPqnDIM2qiK36up3SKkYwFjff6F2ni/ry8vrwXCX3sGZ1hbIHlK0O012HpA3ISeEswVZmX2X67naOvJXfY5v0hGPWqCADao+xVxrmxsZD4IWnKl1UaZzI5lhAzr8fw6utHwx1EZ/MSgsEki6tujcZfN+GUDRnmJGQSnPTXmsf7Q4DKreTZk49cuyB3prV91S0x3DYjCUpSXrkVy1Ha5XicGD/q+ystuFsJnrrhbNXJbpSjM6sjo/aduAkZJl4FmOt0R7Q=="
        userName = "example-user"
      }
    )
  }

  service_account {
    scopes = ["userinfo-email", "compute-ro", "storage-ro"]
  }
}

data "google_compute_instance_serial_port" "serial" {
  instance = google_compute_instance.windows.name
  zone     = google_compute_instance.windows.zone
  port     = 4
}

output "serial_out" {
  value = data.google_compute_instance_serial_port.serial.contents
}
```

## Argument Reference

The following arguments are supported:

* `instance` - (Required) The name of the Compute Instance to read output from.

* `port` - (Required) The number of the serial port to read output from. Possible values are 1-4.

- - -

* `project` - (Optional) The project in which the Compute Instance exists. If it
    is not provided, the provider project is used.

* `zone` - (Optional) The zone in which the Compute Instance exists.
    If it is not provided, the provider zone is used.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `contents` - The output of the serial port. Serial port output is available only when the VM instance is running, and logs are limited to the most recent 1 MB of output per port.
