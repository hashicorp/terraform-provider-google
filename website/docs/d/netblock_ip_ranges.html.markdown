---
subcategory: "Cloud Platform"
layout: "google"
page_title: "Google: google_netblock_ip_ranges"
sidebar_current: "docs-google-datasource-netblock-ip-ranges"
description: |-
  Use this data source to get the IP addresses from different special IP ranges on Google Cloud Platform.
---

# google_netblock_ip_ranges

Use this data source to get the IP addresses from different special IP ranges on Google Cloud Platform.

## Example Usage - Cloud Ranges

```tf
data "google_netblock_ip_ranges" "netblock" {
}

output "cidr_blocks" {
  value = data.google_netblock_ip_ranges.netblock.cidr_blocks
}

output "cidr_blocks_ipv4" {
  value = data.google_netblock_ip_ranges.netblock.cidr_blocks_ipv4
}

output "cidr_blocks_ipv6" {
  value = data.google_netblock_ip_ranges.netblock.cidr_blocks_ipv6
}
```

## Example Usage - Allow Health Checks

```tf
data "google_netblock_ip_ranges" "legacy-hcs" {
  range_type = "legacy-health-checkers"
}

resource "google_compute_firewall" "allow-hcs" {
  name    = "allow-hcs"
  network = google_compute_network.default.name

  allow {
    protocol = "tcp"
    ports    = ["80"]
  }

  source_ranges = data.google_netblock_ip_ranges.legacy-hcs.cidr_blocks_ipv4
}

resource "google_compute_network" "default" {
  name = "test-network"
}
```

## Argument Reference

The following arguments are supported:

* `range_type` (Optional) - The type of range for which to provide results.

  Defaults to `cloud-netblocks`. The following `range_type`s are supported:

  * `cloud-netblocks` - Corresponds to the IP addresses used for resources on Google Cloud Platform. [More details.](https://cloud.google.com/compute/docs/faq#where_can_i_find_product_name_short_ip_ranges)

  * `google-netblocks` - Corresponds to IP addresses used for Google services. [More details.](https://cloud.google.com/compute/docs/faq#where_can_i_find_product_name_short_ip_ranges)

  * `restricted-googleapis` - Corresponds to the IP addresses used for Private Google Access only for services that support VPC Service Controls API access. [More details.](https://cloud.google.com/vpc/docs/private-access-options#domain-vips)

  * `private-googleapis` - Corresponds to the IP addresses used for Private Google Access for services that do not support VPC Service Controls. [More details.](https://cloud.google.com/vpc/docs/private-access-options#domain-vips)

  * `dns-forwarders` - Corresponds to the IP addresses used to originate Cloud DNS outbound forwarding. [More details.](https://cloud.google.com/dns/zones/#creating-forwarding-zones)

  * `iap-forwarders` - Corresponds to the IP addresses used for Cloud IAP for TCP forwarding. [More details.](https://cloud.google.com/iap/docs/using-tcp-forwarding)

  * `health-checkers` - Corresponds to the IP addresses used for health checking in Cloud Load Balancing. [More details.](https://cloud.google.com/load-balancing/docs/health-checks)

  * `legacy-health-checkers` - Corresponds to the IP addresses used for legacy style health checkers (used by Network Load Balancing). [ More details.](https://cloud.google.com/load-balancing/docs/health-checks)


## Attributes Reference

* `cidr_blocks` - Retrieve list of all CIDR blocks.

* `cidr_blocks_ipv4` - Retrieve list of the IPv4 CIDR blocks

* `cidr_blocks_ipv6` - Retrieve list of the IPv6 CIDR blocks, if available.
