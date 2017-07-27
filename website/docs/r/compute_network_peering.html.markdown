---
layout: "google"
page_title: "Google: google_compute_network_peering"
sidebar_current: "docs-google-compute-network-peering"
description: |-
  Manages a network peering within GCE.
---

# google\_compute\_network\_peering

Manages a network peering within GCE.

## Example Usage

```hcl
resource "google_compute_network" "default" {
  name                    = "foobar"
  auto_create_subnetworks = "false"
}

resource "google_compute_network" "other" {
  name                    = "other"
  auto_create_subnetworks = "false"
}

// Both network must create a peering with each other for the peering
// to be functional.
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

// Subnets IP ranges across peered VPC networks cannot overlap.
resource "google_compute_subnetwork" "network1-subnet1" {
  name = "network1-sub1"
  ip_cidr_range = "10.128.0.0/20"  
  network = "${google_compute_network.network1.self_link}"
  region = "us-east1"
}

resource "google_compute_subnetwork" "network2-subnet1" {
  name = "network1-sub2"
  ip_cidr_range = "10.132.0.0/20"  
  network = "${google_compute_network.network2.self_link}"
  region = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the peering.

* `network` - (Required) Resource link of the network to add a peering to.

* `peer_network` - (Required) Resource link of the peer network.

* `auto_create_routes` - (Optional) If set to true, the routes between the two networks will
  be created and managed automatically. Defaults to true.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `state` - (Computed) State for the peering.

* `state_details` - (Computed) Details about the current state of the peering.
