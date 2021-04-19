---
subcategory: "Compute Engine"
layout: "google"
page_title: "Google: google_compute_ha_vpn_gateway"
sidebar_current: "docs-google-datasource-compute-ha-vpn-gateway"
description: |-
  Get a HA VPN Gateway within GCE.
---

# google\_compute\_forwarding\_rule

Get a HA VPN Gateway within GCE from its name.

## Example Usage

```tf
data "google_compute_ha_vpn_gateway" "gateway" {
  name = "foobar"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the forwarding rule.


- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

* `region` - (Optional) The region in which the resource belongs. If it
    is not provided, the project region is used.

## Attributes Reference
See [google_compute_ha_vpn_gateway](https://www.terraform.io/docs/providers/google/r/compute_ha_vpn_gateway.html) resource for details of the available attributes.
