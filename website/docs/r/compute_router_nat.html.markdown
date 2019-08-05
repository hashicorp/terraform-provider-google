---
layout: "google"
page_title: "Google: google_compute_router_nat"
sidebar_current: "docs-google-compute-router-nat"
description: |-
  Manages a Cloud NAT.
---

# google\_compute\_router\_nat

Manages a Cloud NAT. For more information see
[the official documentation](https://cloud.google.com/nat/docs/overview)
and
[API](https://cloud.google.com/compute/docs/reference/rest/beta/routers).

## Example Usage

A simple NAT configuration: enable NAT for all Subnetworks associated with
the Network associated with the given Router.

```hcl
resource "google_compute_network" "default" {
  name = "my-network"
}

resource "google_compute_subnetwork" "default" {
  name          = "my-subnet"
  network       = google_compute_network.default.self_link
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
}

resource "google_compute_router" "router" {
  name    = "router"
  region  = google_compute_subnetwork.default.region
  network = google_compute_network.default.self_link
  bgp {
    asn = 64514
  }
}

resource "google_compute_router_nat" "simple-nat" {
  name                               = "nat-1"
  router                             = google_compute_router.router.name
  region                             = "us-central1"
  nat_ip_allocate_option             = "AUTO_ONLY"
  source_subnetwork_ip_ranges_to_nat = "ALL_SUBNETWORKS_ALL_IP_RANGES"
}
```

A production-like configuration: enable NAT for one Subnetwork and use a list of
static external IP addresses.

```hcl
resource "google_compute_network" "default" {
  name = "my-network"
}

resource "google_compute_subnetwork" "default" {
  name          = "my-subnet"
  network       = google_compute_network.default.self_link
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
}

resource "google_compute_router" "router" {
  name    = "router"
  region  = google_compute_subnetwork.default.region
  network = google_compute_network.default.self_link
  bgp {
    asn = 64514
  }
}

resource "google_compute_address" "address" {
  count  = 2
  name   = "nat-external-address-${count.index}"
  region = "us-central1"
}

resource "google_compute_router_nat" "advanced-nat" {
  name                               = "nat-1"
  router                             = google_compute_router.router.name
  region                             = "us-central1"
  nat_ip_allocate_option             = "MANUAL_ONLY"
  nat_ips                            = google_compute_address.address[*].self_link
  source_subnetwork_ip_ranges_to_nat = "LIST_OF_SUBNETWORKS"
  subnetwork {
    name                    = google_compute_subnetwork.default.self_link
    source_ip_ranges_to_nat = ["ALL_IP_RANGES"]
  }
  log_config {
    filter = "TRANSLATIONS_ONLY"
    enable = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A unique name for Cloud NAT, required by GCE. Changing
    this forces a new NAT to be created.

* `router` - (Required) The name of the router in which this NAT will be configured.
    Changing this forces a new NAT to be created.

* `nat_ip_allocate_option` - (Required) How external IPs should be allocated for
    this NAT. Valid values are `AUTO_ONLY` or `MANUAL_ONLY`. Changing this forces
    a new NAT to be created.

* `source_subnetwork_ip_ranges_to_nat` - (Required) How NAT should be configured
    per Subnetwork. Valid values include: `ALL_SUBNETWORKS_ALL_IP_RANGES`,
    `ALL_SUBNETWORKS_ALL_PRIMARY_IP_RANGES`, `LIST_OF_SUBNETWORKS`. Changing
    this forces a new NAT to be created.

- - -

* `nat_ips` - (Optional) List of `self_link`s of external IPs. Only valid if
    `nat_ip_allocate_option` is set to `MANUAL_ONLY`. Changing this forces a
    new NAT to be created.

* `subnetwork` - (Optional) One or more subnetwork NAT configurations. Only used
    if `source_subnetwork_ip_ranges_to_nat` is set to `LIST_OF_SUBNETWORKS`. See
    the section below for details on configuration.

* `min_ports_per_vm` - (Optional) Minimum number of ports allocated to a VM
    from this NAT config. If not set, a default number of ports is allocated to a VM.
    Changing this forces a new NAT to be created.

* `udp_idle_timeout_sec` - (Optional) Timeout (in seconds) for UDP connections.
    Defaults to 30s if not set. Changing this forces a new NAT to be created.

* `icmp_idle_timeout_sec` - (Optional) Timeout (in seconds) for ICMP connections.
    Defaults to 30s if not set. Changing this forces a new NAT to be created.

* `tcp_established_idle_timeout_sec` - (Optional) Timeout (in seconds) for TCP
    established connections. Defaults to 1200s if not set. Changing this forces
    a new NAT to be created.

* `tcp_transitory_idle_timeout_sec` - (Optional) Timeout (in seconds) for TCP
    transitory connections. Defaults to 30s if not set. Changing this forces a
    new NAT to be created.

* `project` - (Optional) The ID of the project in which this NAT's router belongs. If it
    is not provided, the provider project is used. Changing this forces a new NAT to be created.

* `region` - (Optional) The region this NAT's router sits in. If not specified,
    the project region will be used. Changing this forces a new NAT to be
    created.

The `subnetwork` block supports:

* `name` - (Required) The `self_link` of the subnetwork to NAT.

* `source_ip_ranges_to_nat` - (Required) List of options for which source IPs in the subnetwork
    should have NAT enabled. Supported values include: `ALL_IP_RANGES`,
    `LIST_OF_SECONDARY_IP_RANGES`, `PRIMARY_IP_RANGE`

* `secondary_ip_range_names` - (Optional) List of the secondary ranges of the subnetwork
    that are allowed to use NAT. This can be populated only if
    `LIST_OF_SECONDARY_IP_RANGES` is one of the values in `source_ip_ranges_to_nat`.

The `log_config` block supports:

* `filter` - (Required) Specifies the desired filtering of logs on this NAT.
    Valid values include: `ALL`, `ERRORS_ONLY`, `TRANSLATIONS_ONLY`

* `enable` - (Required) Whether to export logs.

## Import

Router NATs can be imported using any of these accepted formats:

```
$ terraform import google_compute_router_nat.default {{project}}/{{region}}/{{router}}/{{name}}
$ terraform import google_compute_router_nat.default {{region}}/{{router}}/{{name}}
$ terraform import google_compute_router_nat.default {{router}}/{{name}}
```

-> If you're importing a resource with beta features, make sure to include `-provider=google-beta`
as an argument so that Terraform uses the correct provider to import your resource.
