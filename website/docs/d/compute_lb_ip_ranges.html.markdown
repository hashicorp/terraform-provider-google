---
subcategory: "Cloud Platform"
layout: "google"
page_title: "Google: google_compute_lb_ip_ranges"
sidebar_current: "docs-google-datasource-compute-lb-ip-ranges"
description: |-
  Get information about the IP ranges used when health-checking load balancers.
---

# google_compute_lb_ip_ranges

Use this data source to access IP ranges in your firewall rules.

https://cloud.google.com/compute/docs/load-balancing/health-checks#health_check_source_ips_and_firewall_rules

## Example Usage

```tf
data "google_compute_lb_ip_ranges" "ranges" {
}

resource "google_compute_firewall" "lb" {
  name    = "lb-firewall"
  network = google_compute_network.main.name

  allow {
    protocol = "tcp"
    ports    = ["80"]
  }

  source_ranges = data.google_compute_lb_ip_ranges.ranges.network
  target_tags = [
    "InstanceBehindLoadBalancer",
  ]
}
```

## Argument Reference

There are no arguments available for this data source.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `network` - The IP ranges used for health checks when **Network load balancing** is used

* `http_ssl_tcp_internal` - The IP ranges used for health checks when **HTTP(S), SSL proxy, TCP proxy, and Internal load balancing** is used
