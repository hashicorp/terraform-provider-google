---
layout: "google"
page_title: "Google: google_compute_http_health_check"
sidebar_current: "docs-google-compute-http-health-check"
description: |-
  Manages an HTTP Health Check within GCE.
---

# google\_compute\_http\_health\_check

Manages an HTTP health check within GCE. This is used to monitor instances
behind load balancers. Timeouts or HTTP errors cause the instance to be
removed from the pool. For more information, see [the official
documentation](https://cloud.google.com/compute/docs/load-balancing/health-checks#legacy_health_checks)
and
[API](https://cloud.google.com/compute/docs/reference/latest/httpHealthChecks).

~> **Note:** HttpHealthChecks/HTTPSHealthChecks are legacy health checks. The newer [google_compute_health_check](/docs/providers/google/r/compute_health_check.html) should be preferred for all uses except except [Network Load Balancers](https://cloud.google.com/compute/docs/load-balancing/network/) which still require the legacy HttpHealthChecks.

## Example Usage

```hcl
resource "google_compute_http_health_check" "default" {
  name         = "authentication-health-check"
  request_path = "/health_check"

  timeout_sec        = 1
  check_interval_sec = 1
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A unique name for the resource, required by GCE.
    Changing this forces a new resource to be created.

- - -

* `check_interval_sec` - (Optional) The number of seconds between each poll of
    the instance instance (default 5).

* `description` - (Optional) Textual description field.

* `healthy_threshold` - (Optional) Consecutive successes required (default 2).

* `host` - (Optional) HTTP host header field (default instance's public ip).

* `port` - (Optional) TCP port to connect to (default 80).

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

* `request_path` - (Optional) URL path to query (default /).

* `timeout_sec` - (Optional) The number of seconds to wait before declaring
    failure (default 5).

* `unhealthy_threshold` - (Optional) Consecutive failures required (default 2).


## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `creation_timestamp` - Creation timestamp in RFC3339 text format.

* `self_link` - The URI of the created resource.

## Timeouts

This resource provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 4 minutes.
- `update` - Default is 4 minutes.
- `delete` - Default is 4 minutes.

## Import

HTTP health checks can be imported using the `name`, e.g.

```
$ terraform import google_compute_http_health_check.default authentication-health-check
```
