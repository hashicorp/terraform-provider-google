---
subcategory: "Compute Engine"
layout: "google"
page_title: "Google: google_compute_forwarding_rule"
sidebar_current: "docs-google-datasource-compute-forwarding-rule"
description: |-
  Get a forwarding rule within GCE.
---

# google\_compute\_forwarding\_rule

Get a forwarding rule within GCE from its name.

## Example Usage

```tf
data "google_compute_forwarding_rule" "my-forwarding-rule" {
  name = "forwarding-rule-us-east1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the forwarding rule.


- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

* `region` - (Optional) The region in which the resource belongs. If it
    is not provided, the project region is used.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `description` - Description of this forwarding rule.

* `network` - Network of this forwarding rule.

* `subnetwork` - Subnetwork of this forwarding rule.

* `ip_address` - IP address of this forwarding rule.

* `ip_protocol` - IP protocol of this forwarding rule.

* `ports` - List of ports to use for internal load balancing, if this forwarding rule has any.

* `port_range` - Port range, if this forwarding rule has one.

* `target` - URL of the target pool, if this forwarding rule has one.

* `backend_service` - Backend service, if this forwarding rule has one.

* `load_balancing_scheme` - Type of load balancing of this forwarding rule.

* `region` - Region of this forwarding rule.

* `self_link` - The URI of the resource.
