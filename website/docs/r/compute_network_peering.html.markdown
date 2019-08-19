---
layout: "google"
page_title: "Google: google_compute_network_peering"
sidebar_current: "docs-google-compute-network-peering"
description: |-
  Manages a network peering within GCE.
---

# google\_compute\_network\_peering

Manages a network peering within GCE. For more information see
[the official documentation](https://cloud.google.com/compute/docs/vpc/vpc-peering)
and
[API](https://cloud.google.com/compute/docs/reference/latest/networks).

~> **Note:** Both network must create a peering with each other for the peering to be functional.

~> **Note:** Subnets IP ranges across peered VPC networks cannot overlap.

## Example Usage

```hcl
resource "google_compute_network_peering" "peering1" {
  name = "peering1"
  network = "${google_compute_network.default.self_link}"
  peer_network = "${google_compute_network.other.self_link}"
}

resource "google_compute_network_peering" "peering2" {
  name = "peering2"
  network = "${google_compute_network.other.self_link}"
  peer_network = "${google_compute_network.default.self_link}"
}

resource "google_compute_network" "default" {
  name                    = "foobar"
  auto_create_subnetworks = "false"
}

resource "google_compute_network" "other" {
  name                    = "other"
  auto_create_subnetworks = "false"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the peering.

* `network` - (Required) Resource link of the network to add a peering to.

* `peer_network` - (Required) Resource link of the peer network.

* `auto_create_routes` - (Optional) If set to `true`, the routes between the two networks will
  be created and managed automatically. Defaults to `true`.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `state` - State for the peering.

* `state_details` - Details about the current state of the peering.

## Import
VPC Peering Networks can be imported using the name of the network the peering exists in and the name of the peering network

```
$ terraform import google_compute_network_peering.peering_network network-name/peering-network-name
```