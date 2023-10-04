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
description: |-
  The Compute NetworkFirewallPolicyRule resource
---

# google_compute_network_firewall_policy_rule

The Compute NetworkFirewallPolicyRule resource

## Example Usage - global
```hcl
resource "google_network_security_address_group" "basic_global_networksecurity_address_group" {
  name        = "policy"
  parent      = "projects/my-project-name"
  description = "Sample global networksecurity_address_group"
  location    = "global"
  items       = ["208.80.154.224/32"]
  type        = "IPV4"
  capacity    = 100
}

resource "google_compute_network_firewall_policy" "basic_network_firewall_policy" {
  name        = "policy"
  description = "Sample global network firewall policy"
  project     = "my-project-name"
}

resource "google_compute_network_firewall_policy_rule" "primary" {
  action                  = "allow"
  description             = "This is a simple rule description"
  direction               = "INGRESS"
  disabled                = false
  enable_logging          = true
  firewall_policy         = google_compute_network_firewall_policy.basic_network_firewall_policy.name
  priority                = 1000
  rule_name               = "test-rule"
  target_service_accounts = ["my@service-account.com"]

  match {
    src_ip_ranges = ["10.100.0.1/32"]
    src_fqdns = ["google.com"]
    src_region_codes = ["US"]
    src_threat_intelligences = ["iplist-known-malicious-ips"]

    src_secure_tags {
      name = "tagValues/${google_tags_tag_value.basic_value.name}"
    }

    layer4_configs {
      ip_protocol = "all"
    }
    
    src_address_groups = [google_network_security_address_group.basic_global_networksecurity_address_group.id]
  }
}

resource "google_compute_network" "basic_network" {
  name = "network"
}

resource "google_tags_tag_key" "basic_key" {
  description = "For keyname resources."
  parent      = "organizations/123456789"
  purpose     = "GCE_FIREWALL"
  short_name  = "tagkey"
  purpose_data = {
    network = "my-project-name/${google_compute_network.basic_network.name}"
  }
}

resource "google_tags_tag_value" "basic_value" {
  description = "For valuename resources."
  parent      = "tagKeys/${google_tags_tag_key.basic_key.name}"
  short_name  = "tagvalue"
}

```

## Argument Reference

The following arguments are supported:

* `action` -
  (Required)
  The Action to perform when the client connection triggers the rule. Valid actions are "allow", "deny" and "goto_next".
  
* `direction` -
  (Required)
  The direction in which this rule applies. Possible values: INGRESS, EGRESS
  
* `firewall_policy` -
  (Required)
  The firewall policy of the resource.
  
* `match` -
  (Required)
  A match condition that incoming traffic is evaluated against. If it evaluates to true, the corresponding 'action' is enforced.
  
* `priority` -
  (Required)
  An integer indicating the priority of a rule in the list. The priority must be a positive value between 0 and 2147483647. Rules are evaluated from highest to lowest priority where 0 is the highest priority and 2147483647 is the lowest prority.
  


The `match` block supports:
    
* `dest_address_groups` -
  (Optional)
  Address groups which should be matched against the traffic destination. Maximum number of destination address groups is 10. Destination address groups is only supported in Egress rules.
    
* `dest_fqdns` -
  (Optional)
  Domain names that will be used to match against the resolved domain name of destination of traffic. Can only be specified if DIRECTION is egress.
    
* `dest_ip_ranges` -
  (Optional)
  CIDR IP address range. Maximum number of destination CIDR IP ranges allowed is 5000.
    
* `dest_region_codes` -
  (Optional)
  The Unicode country codes whose IP addresses will be used to match against the source of traffic. Can only be specified if DIRECTION is egress.
    
* `dest_threat_intelligences` -
  (Optional)
  Name of the Google Cloud Threat Intelligence list.
    
* `layer4_configs` -
  (Required)
  Pairs of IP protocols and ports that the rule should match.
    
* `src_address_groups` -
  (Optional)
  Address groups which should be matched against the traffic source. Maximum number of source address groups is 10. Source address groups is only supported in Ingress rules.
    
* `src_fqdns` -
  (Optional)
  Domain names that will be used to match against the resolved domain name of source of traffic. Can only be specified if DIRECTION is ingress.
    
* `src_ip_ranges` -
  (Optional)
  CIDR IP address range. Maximum number of source CIDR IP ranges allowed is 5000.
    
* `src_region_codes` -
  (Optional)
  The Unicode country codes whose IP addresses will be used to match against the source of traffic. Can only be specified if DIRECTION is ingress.
    
* `src_secure_tags` -
  (Optional)
  List of secure tag values, which should be matched at the source of the traffic. For INGRESS rule, if all the <code>srcSecureTag</code> are INEFFECTIVE, and there is no <code>srcIpRange</code>, this rule will be ignored. Maximum number of source tag values allowed is 256.
    
* `src_threat_intelligences` -
  (Optional)
  Name of the Google Cloud Threat Intelligence list.
    
The `layer4_configs` block supports:
    
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
  
* `project` -
  (Optional)
  The project for the resource
  
* `rule_name` -
  (Optional)
  An optional name for the rule. This field is not a unique identifier and can be updated.
  
* `target_secure_tags` -
  (Optional)
  A list of secure tags that controls which instances the firewall rule applies to. If <code>targetSecureTag</code> are specified, then the firewall rule applies only to instances in the VPC network that have one of those EFFECTIVE secure tags, if all the target_secure_tag are in INEFFECTIVE state, then this rule will be ignored. <code>targetSecureTag</code> may not be set at the same time as <code>targetServiceAccounts</code>. If neither <code>targetServiceAccounts</code> nor <code>targetSecureTag</code> are specified, the firewall rule applies to all instances on the specified network. Maximum number of target label tags allowed is 256.
  
* `target_service_accounts` -
  (Optional)
  A list of service accounts indicating the sets of instances that are applied with this rule.
  


The `src_secure_tags` block supports:
    
* `name` -
  (Required)
  Name of the secure tag, created with TagManager's TagValue API. @pattern tagValues/[0-9]+
    
* `state` -
  [Output Only] State of the secure tag, either `EFFECTIVE` or `INEFFECTIVE`. A secure tag is `INEFFECTIVE` when it is deleted or its network is deleted.
    
The `target_secure_tags` block supports:
    
* `name` -
  (Required)
  Name of the secure tag, created with TagManager's TagValue API. @pattern tagValues/[0-9]+
    
* `state` -
  [Output Only] State of the secure tag, either `EFFECTIVE` or `INEFFECTIVE`. A secure tag is `INEFFECTIVE` when it is deleted or its network is deleted.
    
## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/global/firewallPolicies/{{firewall_policy}}/rules/{{priority}}`

* `kind` -
  Type of the resource. Always `compute#firewallPolicyRule` for firewall policy rules
  
* `rule_tuple_count` -
  Calculation of the complexity of a single firewall policy rule.
  
## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options:

- `create` - Default is 20 minutes.
- `update` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import

NetworkFirewallPolicyRule can be imported using any of these accepted formats:
* `projects/{{project}}/global/firewallPolicies/{{firewall_policy}}/rules/{{priority}}`
* `{{project}}/{{firewall_policy}}/{{priority}}`
* `{{firewall_policy}}/{{priority}}`

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import NetworkFirewallPolicyRule using one of the formats above. For example:


```tf
import {
  id = "projects/{{project}}/global/firewallPolicies/{{firewall_policy}}/rules/{{priority}}"
  to = google_compute_network_firewall_policy_rule.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), NetworkFirewallPolicyRule can be imported using one of the formats above. For example:

```
$ terraform import google_compute_network_firewall_policy_rule.default projects/{{project}}/global/firewallPolicies/{{firewall_policy}}/rules/{{priority}}
$ terraform import google_compute_network_firewall_policy_rule.default {{project}}/{{firewall_policy}}/{{priority}}
$ terraform import google_compute_network_firewall_policy_rule.default {{firewall_policy}}/{{priority}}
```



