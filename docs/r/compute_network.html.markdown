---
layout: "google"
page_title: "Google: google_compute_network"
sidebar_current: "docs-google-compute-network-x"
description: |-
  Manages a network within GCE.
---

# google\_compute\_network

Manages a network within GCE. For more information see
[the official documentation](https://cloud.google.com/compute/docs/vpc)
and
[API](https://cloud.google.com/compute/docs/reference/latest/networks).

## Example Usage

```hcl
resource "google_compute_network" "default" {
  name                    = "foobar"
  auto_create_subnetworks = "true"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A unique name for the resource, required by GCE.
    Changing this forces a new resource to be created.

- - -

* `auto_create_subnetworks` - (Optional) If set to true, this network will be
    created in auto subnet mode, and Google will create a subnet for each region
    automatically. If set to false, a custom subnetted network will be created that
    can support `google_compute_subnetwork` resources. Defaults to true.

* `ipv4_range` - (Optional) If set to a CIDR block, uses the legacy VPC API with the
  specified range. This API is deprecated. If set, `auto_create_subnetworks` must be
  explicitly set to false.

* `routing_mode` - (Optional) Sets the network-wide routing mode for Cloud Routers
  to use. Accepted values are `"GLOBAL"` or `"REGIONAL"`. Defaults to `"REGIONAL"`.
  Refer to the [Cloud Router documentation](https://cloud.google.com/router/docs/concepts/overview#dynamic-routing-mode)
  for more details.

* `description` - (Optional) A brief description of this resource.

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `gateway_ipv4` - The IPv4 address of the gateway.

* `name` - The unique name of the network.

* `self_link` - The URI of the created resource.


## Import

Networks can be imported using the `name`, e.g.

```
$ terraform import google_compute_network.default foobar
```
