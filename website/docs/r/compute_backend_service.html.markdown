---
layout: "google"
page_title: "Google: google_compute_backend_service"
sidebar_current: "docs-google-compute-backend-service"
description: |-
  Creates a Backend Service resource for Google Compute Engine.
---

# google\_compute\_backend\_service

A Backend Service defines a group of virtual machines that will serve traffic for load balancing. For more information
see [the official documentation](https://cloud.google.com/compute/docs/load-balancing/http/backend-service)
and the [API](https://cloud.google.com/compute/docs/reference/latest/backendServices).

For internal load balancing, use a [google_compute_region_backend_service](/docs/providers/google/r/compute_region_backend_service.html).

## Example Usage

```hcl
resource "google_compute_backend_service" "website" {
  name        = "my-backend"
  description = "Our company website"
  port_name   = "http"
  protocol    = "HTTP"
  timeout_sec = 10
  enable_cdn  = false

  backend {
    group = "${google_compute_instance_group_manager.webservers.instance_group}"
  }

  health_checks = ["${google_compute_http_health_check.default.self_link}"]
}

resource "google_compute_instance_group_manager" "webservers" {
  name               = "my-webservers"
  instance_template  = "${google_compute_instance_template.webserver.self_link}"
  base_instance_name = "webserver"
  zone               = "us-central1-f"
  target_size        = 1
}

resource "google_compute_instance_template" "webserver" {
  name         = "standard-webserver"
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

resource "google_compute_http_health_check" "default" {
  name               = "test"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the backend service.

* `health_checks` - (Required) Specifies a list of HTTP/HTTPS health checks
    for checking the health of the backend service. Currently at most one health
    check can be specified, and a health check is required.

- - -

* `backend` - (Optional) The list of backends that serve this BackendService. Structure is documented below.

* `iap` - (Optional) Specification for the Identity-Aware proxy. Disabled if not specified. Structure is documented below.

* `cdn_policy` - (Optional) Cloud CDN configuration for this BackendService. Structure is documented below.

* `connection_draining_timeout_sec` - (Optional) Time for which instance will be drained (not accept new connections,
but still work to finish started ones). Defaults to `300`.

* `custom_request_headers` - (Optional, [Beta](https://terraform.io/docs/providers/google/provider_versions.html)) Headers that the
    HTTP/S load balancer should add to proxied requests. See [guide](https://cloud.google.com/compute/docs/load-balancing/http/backend-service#user-defined-request-headers) for details.

* `description` - (Optional) The textual description for the backend service.

* `enable_cdn` - (Optional) Whether or not to enable the Cloud CDN on the backend service.

* `port_name` - (Optional) The name of a service that has been added to an
    instance group in this backend. See [related docs](https://cloud.google.com/compute/docs/instance-groups/#specifying_service_endpoints) for details. Defaults to http.

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

* `protocol` - (Optional) The protocol for incoming requests. Defaults to
    `HTTP`.

* `security_policy` - (Optional) Name or URI of a
    [security policy](https://cloud.google.com/armor/docs/security-policy-concepts) to add to the backend service.

* `session_affinity` - (Optional) How to distribute load. Options are `NONE` (no
    affinity), `CLIENT_IP` (hash of the source/dest addresses / ports), and
    `GENERATED_COOKIE` (distribute load using a generated session cookie).

* `timeout_sec` - (Optional) The number of secs to wait for a backend to respond
    to a request before considering the request failed. Defaults to `30`.

The `backend` block supports:

* `group` - (Required) The name or URI of a Compute Engine instance group
    (`google_compute_instance_group_manager.xyz.instance_group`) that can
    receive traffic.

* `balancing_mode` - (Optional) Defines the strategy for balancing load.
    Defaults to `UTILIZATION`

* `capacity_scaler` - (Optional) A float in the range [0, 1.0] that scales the
    maximum parameters for the group (e.g., max rate). A value of 0.0 will cause
    no requests to be sent to the group (i.e., it adds the group in a drained
    state). The default is 1.0.

* `description` - (Optional) Textual description for the backend.

* `max_rate` - (Optional) Maximum requests per second (RPS) that the group can
    handle.

* `max_rate_per_instance` - (Optional) The maximum per-instance requests per
    second (RPS).

* `max_connections` - (Optional) The max number of simultaneous connections for the
    group. Can be used with either CONNECTION or UTILIZATION balancing
    modes. For CONNECTION mode, either maxConnections or
    maxConnectionsPerInstance must be set.

* `max_connections_per_instance` - (Optional) The max number of simultaneous connections
    that a single backend instance can handle. This is used to calculate
    the capacity of the group. Can be used in either CONNECTION or
    UTILIZATION balancing modes. For CONNECTION mode, either
    maxConnections or maxConnectionsPerInstance must be set.

* `max_utilization` - (Optional) The target CPU utilization for the group as a
    float in the range [0.0, 1.0]. This flag can only be provided when the
    balancing mode is `UTILIZATION`. Defaults to `0.8`.

The `cdn_policy` block supports:

* `cache_key_policy` - (Optional) The CacheKeyPolicy for this CdnPolicy.
    Structure is documented below.

The `cache_key_policy` block supports:

* `include_host` - (Optional) If true, requests to different hosts will be cached separately.

* `include_protocol` - (Optional) If true, http and https requests will be cached separately.

* `include_query_string` - (Optional) If true, include query string parameters in the cache key
    according to `query_string_whitelist` and `query_string_blacklist`. If neither is set, the entire
    query string will be included. If false, the query string will be excluded from the cache key entirely.

* `query_string_blacklist` - (Optional) Names of query string parameters to exclude in cache keys.
    All other parameters will be included. Either specify `query_string_whitelist` or
    `query_string_blacklist`, not both. '&' and '=' will be percent encoded and not treated as delimiters.

* `query_string_whitelist` - (Optional) Names of query string parameters to include in cache keys.
    All other parameters will be excluded. Either specify `query_string_whitelist` or
    `query_string_blacklist`, not both. '&' and '=' will be percent encoded and not treated as delimiters.

The `iap` block supports:

* `oauth2_client_id` - (Required) The client ID for use with OAuth 2.0.

* `oauth2_client_secret` - (Required) The client secret for use with OAuth 2.0.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `fingerprint` - The fingerprint of the backend service.

* `self_link` - The URI of the created resource.

## Import

Backend services can be imported using the `name`, e.g.

```
$ terraform import google_compute_backend_service.website my-backend
```
