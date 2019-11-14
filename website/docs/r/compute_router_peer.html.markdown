---
subcategory: "Compute Engine"
layout: "google"
page_title: "Google: google_compute_router_peer"
sidebar_current: "docs-google-compute-router-peer"
description: |-
  Manages a Cloud Router BGP peer.
---

# google\_compute\_router\_peer

Manages a Cloud Router BGP peer. For more information see
[the official documentation](https://cloud.google.com/compute/docs/cloudrouter)
and
[API](https://cloud.google.com/compute/docs/reference/latest/routers).

## Example Usage

```hcl
resource "google_compute_router_peer" "foobar" {
  name                      = "peer-1"
  router                    = "router-1"
  region                    = "us-central1"
  peer_ip_address           = "169.254.1.2"
  peer_asn                  = 65513
  advertised_route_priority = 100
  interface                 = "interface-1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A unique name for BGP peer, required by GCE. Changing
    this forces a new peer to be created.

* `router` - (Required) The name of the router in which this BGP peer will be configured.
    Changing this forces a new peer to be created.

* `interface` - (Required) The name of the interface the BGP peer is associated with.
    Changing this forces a new peer to be created.

* `peer_ip_address` - (Required) IP address of the BGP interface outside Google Cloud.
    Changing this forces a new peer to be created.

* `peer_asn` - (Required) Peer BGP Autonomous System Number (ASN).
    Changing this forces a new peer to be created.

- - -

* `advertised_route_priority` - (Optional) The priority of routes advertised to this BGP peer.
    Changing this forces a new peer to be created.

* `advertise_mode` - (Optional) User-specified flag to indicate which mode to use for advertisement.
    Options include `DEFAULT` or `CUSTOM`.

* `advertised_groups` - (Optional) User-specified list of prefix groups to advertise in custom mode,
    which can take one of the following options:

    `ALL_SUBNETS`: Advertises all available subnets, including peer VPC subnets.  
    `ALL_VPC_SUBNETS`: Advertises the router's own VPC subnets.  
    `ALL_PEER_VPC_SUBNETS`: Advertises peer subnets of the router's VPC network.

    Note that this field can only be populated if `advertise_mode` is `CUSTOM` and overrides the list
    defined for the router (in the "bgp" message). These groups are advertised in addition to any
    specified prefixes. Leave this field blank to advertise no custom groups.

* `advertised_ip_ranges` - (Optional) User-specified list of individual IP ranges to advertise in
    custom mode. This field can only be populated if `advertise_mode` is `CUSTOM` and overrides
    the list defined for the router (in the "bgp" message). These IP ranges are advertised in
    addition to any specified groups. Leave this field blank to advertise no custom IP ranges.

* `project` - (Optional) The ID of the project in which this peer's router belongs. If it
    is not provided, the provider project is used. Changing this forces a new peer to be created.

* `region` - (Optional) The region this peer's router sits in. If not specified,
    the project region will be used. Changing this forces a new peer to be
    created.


The `advertised_ip_ranges` block supports:

* `description` -
  (Optional) User-specified description for the IP range.

* `range` -
  (Required) The IP range to advertise. The value must be a CIDR-formatted string.


## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `ip_address` - IP address of the interface inside Google Cloud Platform.

## Import

Router BGP peers can be imported using the `region`, `router`, and `name`, e.g.

```
$ terraform import google_compute_router_peer.foobar us-central1/router-1/peer-1
```
