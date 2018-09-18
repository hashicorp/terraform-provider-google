---
layout: "google"
page_title: "Google: google_netblock_ip_ranges"
sidebar_current: "docs-google-datasource-netblock-ip-ranges"
description: |-
  Use this data source to get the IP ranges from the sender policy framework (SPF) record of \_cloud-netblocks.googleusercontent.com
---

# google_netblock_ip_ranges

Use this data source to get the IP ranges from the sender policy framework (SPF) record of \_cloud-netblocks.googleusercontent

https://cloud.google.com/compute/docs/faq#where_can_i_find_product_name_short_ip_ranges

## Example Usage

```tf
data "google_netblock_ip_ranges" "netblock" {}

output "cidr_blocks" {
  value = "${data.google_netblock_ip_ranges.netblock.cidr_blocks}"
}

output "cidr_blocks_ipv4" {
  value = "${data.google_netblock_ip_ranges.netblock.cidr_blocks_ipv4}"
}

output "cidr_blocks_ipv6" {
  value = "${data.google_netblock_ip_ranges.netblock.cidr_blocks_ipv6}"
}
```

## Attributes Reference

* `cidr_blocks` - Retrieve list of all CIDR blocks.

* `cidr_blocks_ipv4` - Retrieve list of the IP4 CIDR blocks

* `cidr_blocks_ipv6` - Retrieve list of the IP6 CIDR blocks.
