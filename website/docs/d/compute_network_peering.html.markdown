---
subcategory: "Compute Engine"
description: |-
  Get information of a specified compute network peering.
---

# google_compute_network_peering

Get information of a specified compute network peering. For more information see
[the official documentation](https://cloud.google.com/compute/docs/vpc/vpc-peering)
and
[API](https://cloud.google.com/compute/docs/reference/latest/networks).

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
data "google_compute_network_peering" "peering1_ds" {
  name       = google_compute_network_peering.peering1.name
  network    = google_compute_network_peering.peering1.network
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the peering.

* `network` - (Required) The primary network of the peering.

## Attributes Reference

See [google_compute_network_peering](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_network_peering#argument-reference) resource for details of the available attributes.

## Timeouts

This datasource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options: configuration options:

- `read` - Default is 4 minutes.
