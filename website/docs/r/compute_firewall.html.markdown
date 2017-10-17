---
layout: "google"
page_title: "Google: google_compute_firewall"
sidebar_current: "docs-google-compute-firewall"
description: |-
  Manages a firewall resource within GCE.
---

# google\_compute\_firewall

Manages a firewall resource within GCE. For more information see
[the official documentation](https://cloud.google.com/compute/docs/vpc/firewalls)
and
[API](https://cloud.google.com/compute/docs/reference/latest/firewalls).

## Example Usage

```hcl
resource "google_compute_firewall" "default" {
  name    = "test-firewall"
  network = "${google_compute_network.other.name}"

  allow {
    protocol = "icmp"
  }

  allow {
    protocol = "tcp"
    ports    = ["80", "8080", "1000-2000"]
  }

  source_tags = ["web"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A unique name for the resource, required by GCE.
    Changing this forces a new resource to be created.

* `network` - (Required) The name or self_link of the network to attach this firewall to.

* `allow` - (Required) Can be specified multiple times for each allow
    rule. Each allow block supports fields documented below.

- - -

* `description` - (Optional) Textual description field.

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

* `priority` - (Optional) The priority for this firewall. Ranges from 0-65535, inclusive. Defaults to 1000. Firewall
    resources with lower priority values have higher precedence (e.g. a firewall resource with a priority value of 0
    takes effect over all other firewall rules with a non-zero priority).

* `source_ranges` - (Optional) A list of source CIDR ranges that this
   firewall applies to. Can't be used for `EGRESS`.

* `source_tags` - (Optional) A list of source tags for this firewall. Can't be used for `EGRESS`.

* `target_tags` - (Optional) A list of target tags for this firewall.

- - -

* `deny` - (Optional, [Beta](/docs/providers/google/index.html#beta-features)) Can be specified multiple times for each deny
    rule. Each deny block supports fields documented below. Can be specified
    instead of allow.

* `direction` - (Optional, [Beta](/docs/providers/google/index.html#beta-features)) Direction of traffic to which this firewall applies;
    One of `INGRESS` or `EGRESS`. Defaults to `INGRESS`.

* `destination_ranges` - (Optional, [Beta](/docs/providers/google/index.html#beta-features)) A list of destination CIDR ranges that this
   firewall applies to. Can't be used for `INGRESS`.

The `allow` block supports:

* `protocol` - (Required) The name of the protocol to allow.

* `ports` - (Optional) List of ports and/or port ranges to allow. This can
    only be specified if the protocol is TCP or UDP.

The `deny` block supports:

* `protocol` - (Required) The name of the protocol to allow.

* `ports` - (Optional) List of ports and/or port ranges to allow. This can
    only be specified if the protocol is TCP or UDP.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `self_link` - The URI of the created resource.


## Import

Firewalls can be imported using the `name`, e.g.

```
$ terraform import google_compute_firewall.default test-firewall
```
