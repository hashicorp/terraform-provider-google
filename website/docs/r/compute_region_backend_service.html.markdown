---
layout: "google"
page_title: "Google: google_compute_region_backend_service"
sidebar_current: "docs-google-compute-region-backend-service"
description: |-
  Creates a Region Backend Service resource for Google Compute Engine.
---

# google\_compute\_region\_backend\_service

A Region Backend Service defines a regionally-scoped group of virtual machines that will serve traffic for load balancing.
For more information see [the official documentation](https://cloud.google.com/compute/docs/load-balancing/internal/)
and [API](https://cloud.google.com/compute/docs/reference/latest/regionBackendServices).

~> **Note**: Region backend services can only be used when using internal load balancing. For external load balancing, use
  [`google_compute_backend_service`](compute_backend_service.html) instead.

## Example Usage

```hcl
resource "google_compute_region_backend_service" "foobar" {
  name             = "blablah"
  description      = "Hello World 1234"
  protocol         = "TCP"
  timeout_sec      = 10
  session_affinity = "CLIENT_IP"

  backend {
    group = "${google_compute_region_instance_group_manager.foo.instance_group}"
  }

  health_checks = ["${google_compute_health_check.default.self_link}"]
}

resource "google_compute_region_instance_group_manager" "foo" {
  name               = "terraform-test"
  instance_template  = "${google_compute_instance_template.foobar.self_link}"
  base_instance_name = "foobar"
  region             = "us-central1"
  target_size        = 1
}

resource "google_compute_instance_template" "foobar" {
  name         = "terraform-test"
  machine_type = "n1-standard-1"

  network_interface {
    network = "default"
  }

  disk {
    source_image = "debian-cloud/debian-9"
    auto_delete  = true
    boot         = true
  }
}

resource "google_compute_health_check" "default" {
  name               = "test"
  check_interval_sec = 1
  timeout_sec        = 1

  tcp_health_check {
    port = "80"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the backend service.

* `health_checks` - (Required) Specifies a list of health checks
    for checking the health of the backend service. Currently at most
    one health check can be specified, and a health check is required.

- - -

* `backend` - (Optional) The list of backends that serve this BackendService.
    Structure is documented below.

* `description` - (Optional) The textual description for the backend service.

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

* `protocol` - (Optional) The protocol for incoming requests. Defaults to
    `TCP`.

* `session_affinity` - (Optional) How to distribute load. Options are `NONE` (no
    affinity), `CLIENT_IP`, `CLIENT_IP_PROTO`, or `CLIENT_IP_PORT_PROTO`.
    Defaults to `NONE`.

* `region` - (Optional) The Region in which the created address should reside.
    If it is not provided, the provider region is used.

* `timeout_sec` - (Optional) The number of secs to wait for a backend to respond
    to a request before considering the request failed. Defaults to `30`.

* `connection_draining_timeout_sec` - (Optional) Time for which instance will be drained
(not accept new connections, but still work to finish started ones). Defaults to `0`.

The `backend` block supports:

* `group` - (Required) The name or URI of a Compute Engine instance group
    (`google_compute_region_instance_group_manager.xyz.instance_group`) that can
    receive traffic. Instance groups must contain at least one instance.

* `description` - (Optional) Textual description for the backend.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `fingerprint` - The fingerprint of the backend service.

* `self_link` - The URI of the created resource.
