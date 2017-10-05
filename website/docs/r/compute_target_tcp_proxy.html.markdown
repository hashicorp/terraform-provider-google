---
layout: "google"
page_title: "Google: google_compute_target_tcp_proxy"
sidebar_current: "docs-google-compute-target-tcp-proxy"
description: |-
  Creates a Target TCP Proxy resource in GCE.
---

# google\_compute\_target\_tcp\_proxy

Creates a target TCP proxy resource in GCE. For more information see
[the official
documentation](https://cloud.google.com/compute/docs/load-balancing/tcp-ssl/tcp-proxy) and
[API](https://cloud.google.com/compute/docs/reference/latest/targetTcpProxies).


## Example Usage

```hcl
resource "google_compute_target_tcp_proxy" "default" {
  name = "test"
  description = "test"
  backend_service = "${google_compute_backend_service.default.self_link}"
}

resource "google_compute_backend_service" "default" {
  name        = "default-backend"
  protocol    = "TCP"
  timeout_sec = 10

  health_checks = ["${google_compute_health_check.default.self_link}"]
}

resource "google_compute_health_check" "default" {
  name = "default"
  timeout_sec        = 1
  check_interval_sec = 1

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

- - -

* `proxy_header` - (Optional) Type of proxy header to append before sending
    data to the backend, either NONE or PROXY_V1 (default NONE).

* `description` - (Optional) A description of this resource. Changing this
    forces a new resource to be created.

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `proxy_id` - A unique ID assigned by GCE.

* `self_link` - The URI of the created resource.

## Import

TCP proxy can be imported using the `name`, e.g.

```
$ terraform import google_compute_target_tcp_proxy.default test
```