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
page_title: "Google: google_compute_region_network_firewall_policy_rule"
description: |-
  The Compute NetworkFirewallPolicyRule resource
---

# google_compute_region_network_firewall_policy_rule

The Compute NetworkFirewallPolicyRule resource

## Example Usage - regional
```hcl
resource "google_compute_region_network_firewall_policy" "basic_regional_network_firewall_policy" {
  name = "policy"
  project = "my-project-name"
  description = "Sample regional network firewall policy"
  region = "us-west1"
}

resource "google_compute_region_network_firewall_policy_rule" "primary" {
 firewall_policy = google_compute_region_network_firewall_policy.basic_regional_network_firewall_policy.name
 action = "allow"
 direction = "INGRESS"
 priority = 1000
 rule_name = "test-rule"
 description = "This is a simple rule description"
match {
 src_secure_tags {
 name = "tagValues/${google_tags_tag_value.basic_value.name}"
 }
 src_ip_ranges = ["10.100.0.1/32"]
layer4_configs {
ip_protocol = "all"
 }
 }
 target_service_accounts = ["emailAddress:my@service-account.com"]
 region = "us-west1"
 enable_logging = true
 disabled = false
}

resource "google_compute_network" "basic_network" {
  name = "network"
}
resource "google_tags_tag_key" "basic_key" {
  parent = "organizations/123456789"
  short_name = "tagkey"
  purpose = "GCE_FIREWALL"
  purpose_data = {
  network= "my-project-name/${google_compute_network.basic_network.name}"
  }
  description = "For keyname resources."
}


resource "google_tags_tag_value" "basic_value" {
    parent = "tagKeys/${google_tags_tag_key.basic_key.name}"
    short_name = "tagvalue"
    description = "For valuename resources."
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
  A match condition that incoming traffic is evaluated against. If it evaluates to true, the corresponding 'action' is enforced.
  
* `priority` -
  (Required)
  An integer indicating the priority of a rule in the list. The priority must be a positive value between 0 and 2147483647. Rules are evaluated from highest to lowest priority where 0 is the highest priority and 2147483647 is the lowest prority.
  


The `match` block supports:
    
* `dest_ip_ranges` -
  (Optional)
  CIDR IP address range. Maximum number of destination CIDR IP ranges allowed is 5000.
    
* `layer4_configs` -
  (Required)
  Pairs of IP protocols and ports that the rule should match.
    
* `src_ip_ranges` -
  (Optional)
  CIDR IP address range. Maximum number of source CIDR IP ranges allowed is 5000.
    
* `src_secure_tags` -
  (Optional)
  List of secure tag values, which should be matched at the source of the traffic. For INGRESS rule, if all the <code>srcSecureTag</code> are INEFFECTIVE, and there is no <code>srcIpRange</code>, this rule will be ignored. Maximum number of source tag values allowed is 256.
    
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
  
* `region` -
  (Optional)
  The location of this resource.
  
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

* `id` - an identifier for the resource with format `projects/{{project}}/regions/{{region}}/firewallPolicies/{{firewall_policy}}/{{priority}}`

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

NetworkFirewallPolicyRule can be imported using any of these accepted formats:

```
$ terraform import google_compute_region_network_firewall_policy_rule.default projects/{{project}}/regions/{{region}}/firewallPolicies/{{firewall_policy}}/{{priority}}
$ terraform import google_compute_region_network_firewall_policy_rule.default {{project}}/{{region}}/{{firewall_policy}}/{{priority}}
$ terraform import google_compute_region_network_firewall_policy_rule.default {{region}}/{{firewall_policy}}/{{priority}}
$ terraform import google_compute_region_network_firewall_policy_rule.default {{firewall_policy}}/{{priority}}
```



