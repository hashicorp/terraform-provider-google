---
# ----------------------------------------------------------------------------
#
#     ***     AUTO GENERATED CODE    ***    Type: DCL     ***
#
# ----------------------------------------------------------------------------
#
#     This file is managed by Magic Modules (https:#github.com/GoogleCloudPlatform/magic-modules)
#     and is based on the DCL (https:#github.com/GoogleCloudPlatform/declarative-resource-client-library).
#     Changes will need to be made to the DCL or Magic Modules instead of here.
#
#     We are not currently able to accept contributions to this file. If changes
#     are required, please file an issue at https:#github.com/hashicorp/terraform-provider-google/issues/new/choose
#
# ----------------------------------------------------------------------------
subcategory: "Compute Engine"
layout: "google"
page_title: "Google: google_compute_firewall_policy_rule"
sidebar_current: "docs-google-compute-firewall-policy-rule"
description: |-
  Specific rules to add to a hierarchical firewall policy
---

# google\_compute\_firewall\_policy\_rule

Hierarchical firewall policy rules let you create and enforce a consistent firewall policy across your organization. Rules can explicitly allow or deny connections or delegate evaluation to lower level policies.

For more information see the [official documentation](https://cloud.google.com/vpc/docs/using-firewall-policies#create-rules)

## Example Usage

```hcl
resource "google_compute_firewall_policy" "default" {
  parent      = "organizations/12345"
  short_name  = "my-policy"
  description = "Example Resource"
}

resource "google_compute_firewall_policy_rule" "default" {
  firewall_policy = google_compute_firewall_policy.default.id
  description = "Example Resource"
  priority = 9000
  enable_logging = true
  action = "allow"
  direction = "EGRESS"
  disabled = false
  match {
    layer4_configs {
      ip_protocol = "tcp"
      ports = [80, 8080]
    }
    dest_ip_ranges = ["11.100.0.1/32"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `action` -
  (Required)
  The Action to perform when the client connection triggers the rule. Can currently be either "allow" or "deny()" where valid values for status are 403, 404, and 502.
  
* `direction` -
  (Required)
  The direction in which this rule applies. Possible values: INGRESS, EGRESS
  
* `firewall_policy` -
  (Required)
  The firewall policy of the resource.
  
* `match` -
  (Required)
  A match condition that incoming traffic is evaluated against. If it evaluates to true, the corresponding 'action' is enforced. Structure is [documented below](#nested_match).
  
* `priority` -
  (Required)
  An integer indicating the priority of a rule in the list. The priority must be a positive value between 0 and 2147483647. Rules are evaluated from highest to lowest priority where 0 is the highest priority and 2147483647 is the lowest prority.
  


<a name="nested_match"></a>The `match` block supports:
    
* `dest_ip_ranges` -
  (Optional)
  CIDR IP address range. Maximum number of destination CIDR IP ranges allowed is 256.
    
* `layer4_configs` -
  (Required)
  Pairs of IP protocols and ports that the rule should match. Structure is [documented below](#nested_layer4_configs).
    
* `src_ip_ranges` -
  (Optional)
  CIDR IP address range. Maximum number of source CIDR IP ranges allowed is 256.
    
<a name="nested_layer4_configs"></a>The `layer4_configs` block supports:
    
* `ip_protocol` -
  (Required)
  The IP protocol to which this rule applies. The protocol type is required when creating a firewall rule. This value can either be one of the following well known protocol strings (`tcp`, `udp`, `icmp`, `esp`, `ah`, `ipip`, `sctp`), or the IP protocol number.
    
* `ports` -
  (Optional)
  An optional list of ports to which this rule applies. This field is only applicable for UDP or TCP protocol. Each entry must be either an integer or a range. If not specified, this rule applies to connections through any port. Example inputs include: ``.
    
- - -

* `description` -
  (Optional)
  An optional description for this resource.
  
* `disabled` -
  (Optional)
  Denotes whether the firewall policy rule is disabled. When set to true, the firewall policy rule is not enforced and traffic behaves as if it did not exist. If this is unspecified, the firewall policy rule will be enabled.
  
* `enable_logging` -
  (Optional)
  Denotes whether to enable logging for a particular rule. If logging is enabled, logs will be exported to the configured export destination in Stackdriver. Logs may be exported to BigQuery or Pub/Sub. Note: you cannot enable logging on "goto_next" rules.
  
* `target_resources` -
  (Optional)
  A list of network resource URLs to which this rule applies. This field allows you to control which network's VMs get this rule. If this field is left blank, all VMs within the organization will receive the rule.
  
* `target_service_accounts` -
  (Optional)
  A list of service accounts indicating the sets of instances that are applied with this rule.
  


## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `locations/global/firewallPolicies/{{firewall_policy}}/rules/{{priority}}`

* `kind` -
  Type of the resource. Always `compute#firewallPolicyRule` for firewall policy rules
  
* `rule_tuple_count` -
  Calculation of the complexity of a single firewall policy rule.
  
## Timeouts

This resource provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 20 minutes.
- `update` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import

FirewallPolicyRule can be imported using any of these accepted formats:

```
$ terraform import google_compute_firewall_policy_rule.default locations/global/firewallPolicies/{{firewall_policy}}/rules/{{priority}}
$ terraform import google_compute_firewall_policy_rule.default {{firewall_policy}}/{{priority}}
```



