---
subcategory: "Compute Engine"
page_title: "Google: google_compute_backend_service"
description: |-
  Get information about a Backend Service.
---

# google\_compute\_backend\_service

Provide access to a Backend Service's attribute. For more information
see [the official documentation](https://cloud.google.com/compute/docs/load-balancing/http/backend-service)
and the [API](https://cloud.google.com/compute/docs/reference/latest/backendServices).

## Example Usage

```tf
data "google_compute_backend_service" "baz" {
  name = "foobar"
}

resource "google_compute_backend_service" "default" {
  name          = "backend-service"
  health_checks = [tolist(data.google_compute_backend_service.baz.health_checks)[0]]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Backend Service.

- - -

* `project` - (Optional) The project in which the resource belongs. If it is not provided, the provider project is used.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `connection_draining_timeout_sec` - Time for which instance will be drained (not accept new connections, but still work to finish started ones).

* `description` - Textual description for the Backend Service.

* `enable_cdn` - Whether or not Cloud CDN is enabled on the Backend Service.

* `fingerprint` - The fingerprint of the Backend Service.

* `id` - an identifier for the resource with format `projects/{{project}}/global/backendServices/{{name}}`

* `port_name` - The name of a service that has been added to an instance group in this backend.

* `protocol` - The protocol for incoming requests.

* `self_link` - The URI of the Backend Service.

* `session_affinity` - The Backend Service session stickiness configuration.

* `timeout_sec` - The number of seconds to wait for a backend to respond to a request before considering the request failed.

* `backend` - The set of backends that serve this Backend Service.

* `health_checks` - The set of HTTP/HTTPS health checks used by the Backend Service.

* `generated_id` - The unique identifier for the resource. This identifier is defined by the server.
