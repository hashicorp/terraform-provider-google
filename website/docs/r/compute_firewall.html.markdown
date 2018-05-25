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

- - -

* `allow` - (Required) Can be specified multiple times for each allow
    rule. Each allow block supports fields documented below.
    
* `deny` - (Optional) Can be specified multiple times for each deny
    rule. Each deny block supports fields documented below. Can be specified
    instead of allow.

* `description` - (Optional) Textual description field.

* `disabled` - (Optional) Denotes whether the firewall rule is disabled, i.e not applied to the network it is associated with.
    When set to true, the firewall rule is not enforced and the network behaves as if it did not exist.

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

* `priority` - (Optional) The priority for this firewall. Ranges from 0-65535, inclusive. Defaults to 1000. Firewall
    resources with lower priority values have higher precedence (e.g. a firewall resource with a priority value of 0
    takes effect over all other firewall rules with a non-zero priority).

* `source_ranges` - (Optional) A list of source CIDR ranges that this
   firewall applies to. Can't be used for `EGRESS`.

* `source_tags` - (Optional) A list of source tags for this firewall. Can't be used for `EGRESS`.

* `target_tags` - (Optional) A list of target tags for this firewall.

* `direction` - (Optional) Direction of traffic to which this firewall applies;
    One of `INGRESS` or `EGRESS`. Defaults to `INGRESS`.

* `destination_ranges` - (Optional) A list of destination CIDR ranges that this
   firewall applies to. Can't be used for `INGRESS`.

* `source_service_accounts` - (Optional) A list of service accounts such that
    the firewall will apply only to traffic originating from an instance with a service account in this list.  Note that as of May 2018,
    this list can contain only one item, due to a change in the way that these firewall rules are handled.  Source service accounts
    cannot be used to control traffic to an instance's external IP address because service accounts are associated with an instance, not
    an IP address. `source_ranges` can be set at the same time as `source_service_accounts`. If both are set, the firewall will apply to
    traffic that has source IP address within `source_ranges` OR the source IP belongs to an instance with service account listed in
    `source_service_accounts`. The connection does not need to match both properties for the firewall to apply. `source_service_accounts`
    cannot be used at the same time as `source_tags` or `target_tags`.

* `target_service_accounts` - (Optional) A list of service accounts indicating
    sets of instances located in the network that may make network connections as specified in `allow`. `target_service_accounts` cannot
    be used at the same time as `source_tags` or `target_tags`. If neither `target_service_accounts` nor `target_tags` are specified, the
    firewall rule applies to all instances on the specified network.  Note that as of May 2018, this list can contain only one item, due
    to a change in the way that these firewall rules are handled.

The `allow` block supports:

* `protocol` - (Required) The name of the protocol to allow. This value can either be one of the following well
    known protocol strings (tcp, udp, icmp, esp, ah, sctp), or the IP protocol number, or `all`.

* `ports` - (Optional) List of ports and/or port ranges to allow. This can
    only be specified if the protocol is TCP or UDP.

The `deny` block supports:

* `protocol` - (Required) The name of the protocol to deny. This value can either be one of the following well
    known protocol strings (tcp, udp, icmp, esp, ah, sctp), or the IP protocol number, or `all`.

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
