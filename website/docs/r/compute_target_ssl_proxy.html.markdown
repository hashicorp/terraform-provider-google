---
layout: "google"
page_title: "Google: google_compute_target_ssl_proxy"
sidebar_current: "docs-google-compute-target-ssl-proxy"
description: |-
  Creates a Target SSL Proxy resource in GCE.
---

# google\_compute\_target\_ssl\_proxy

Creates a target SSL proxy resource in GCE. For more information see
[the official
documentation](https://cloud.google.com/compute/docs/load-balancing/ssl-ssl/) and
[API](https://cloud.google.com/compute/docs/reference/latest/targetSslProxies).


## Example Usage

```hcl
resource "google_compute_target_ssl_proxy" "default" {
  name = "test"
  backend_service = "${google_compute_backend_service.default.self_link}"
  ssl_certificates = ["${google_compute_ssl_certificate.default.self_link}"]
}

resource "google_compute_ssl_certificate" "default" {
  name = "default-cert"
  private_key = "${file("path/to/test.key")}"
  certificate = "${file("path/to/test.crt")}"
}

resource "google_compute_backend_service" "default" {
  name = "default-backend"
  protocol    = "SSL"
  health_checks = ["${google_compute_health_check.default.self_link}"]
}

resource "google_compute_health_check" "default" {
  name = "default-health-check"
  check_interval_sec = 1
  timeout_sec = 1
  tcp_health_check {
    port = "443"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A unique name for the resource, required by GCE. Changing
    this forces a new resource to be created.

* `backend_service` - (Required) The URL of a Backend Service resource to receive the matched traffic.

* `ssl_certificates` - (Required) The URLs of the SSL Certificate resources that
    authenticate connections between users and load balancing.

- - -

* `proxy_header` - (Optional) Type of proxy header to append before sending
    data to the backend, either NONE or PROXY_V1 (default NONE).

* `description` - (Optional) A description of this resource. Changing this
    forces a new resource to be created.

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `proxy_id` - A unique ID assigned by GCE.

* `self_link` - The URI of the created resource.

## Import

SSL proxy can be imported using the `name`, e.g.

```
$ terraform import google_compute_target_ssl_proxy.default test
```
