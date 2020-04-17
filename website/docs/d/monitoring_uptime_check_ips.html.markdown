---
subcategory: "Stackdriver Monitoring"
layout: "google"
page_title: "Google: google_monitoring_uptime_check_ips"
sidebar_current: "docs-google-datasource-google-monitoring-uptime-check-ips"
description: |-
  Returns the list of IP addresses Stackdriver Monitoring uses for uptime checking.
---

# google\_monitoring\_uptime\_check\_ips

Returns the list of IP addresses that checkers run from. For more information see
the [official documentation](https://cloud.google.com/monitoring/uptime-checks#get-ips).

## Example Usage

```hcl
data "google_monitoring_uptime_check_ips" "ips" {
}

output "ip_list" {
  value = data.google_monitoring_uptime_check_ips.ips.uptime_check_ips
}
```

## Attributes Reference

The following computed attributes are exported:

* `uptime_check_ips` - A list of uptime check IPs used by Stackdriver Monitoring. Each `uptime_check_ip` contains:
  * `region` - A broad region category in which the IP address is located.
  * `location` - A more specific location within the region that typically encodes a particular city/town/metro
  (and its containing state/province or country) within the broader umbrella region category.
  * `ip_address` - The IP address from which the Uptime check originates. This is a fully specified IP address
  (not an IP address range). Most IP addresses, as of this publication, are in IPv4 format; however, one should not
  rely on the IP addresses being in IPv4 format indefinitely, and should support interpreting this field in either
  IPv4 or IPv6 format.
  
