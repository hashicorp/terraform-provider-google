---
subcategory: "Compute Engine"
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

-> Both network must create a peering with each other for the peering
to be functional.

~> Subnets IP ranges across peered VPC networks cannot overlap.

## Example Usage

```hcl
resource "google_compute_network_peering" "peering1" {
  name         = "peering1"
  network      = google_compute_network.default.self_link
  peer_network = google_compute_network.other.self_link
}

resource "google_compute_network_peering" "peering2" {
  name         = "peering2"
  network      = google_compute_network.other.self_link
  peer_network = google_compute_network.default.self_link
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

* `network` - (Required) The primary network of the peering.

* `peer_network` - (Required) The peer network in the peering. The peer network
may belong to a different project.

* `export_custom_routes` - (Optional)
Whether to export the custom routes to the peer network. Defaults to `false`.

* `import_custom_routes` - (Optional)
Whether to export the custom routes from the peer network. Defaults to `false`.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `state` - State for the peering, either `ACTIVE` or `INACTIVE`. The peering is
`ACTIVE` when there's a matching configuration in the peer network.

* `state_details` - Details about the current state of the peering.

## Timeouts

This resource provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 4 minutes.
- `delete` - Default is 4 minutes.

## Import

VPC network peerings can be imported using the name and project of the primary network the peering exists in and the name of the network peering

```
$ terraform import google_compute_network_peering.peering_network project-name/network-name/peering-name
```
