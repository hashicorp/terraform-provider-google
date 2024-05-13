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
  Creates a hierarchical firewall policy
---

# google_compute_firewall_policy

Hierarchical firewall policy rules let you create and enforce a consistent firewall policy across your organization. Rules can explicitly allow or deny connections or delegate evaluation to lower level policies. Policies can be created within organizations or folders.

This resource should be generally be used with `google_compute_firewall_policy_association` and `google_compute_firewall_policy_rule`

For more information see the [official documentation](https://cloud.google.com/vpc/docs/firewall-policies)

## Example Usage

```hcl
resource "google_compute_firewall_policy" "default" {
  parent      = "organizations/12345"
  short_name  = "my-policy"
  description = "Example Resource"
}
```

## Argument Reference

The following arguments are supported:

* `parent` -
  (Required)
  The parent of the firewall policy.
  
* `short_name` -
  (Required)
  User-provided name of the Organization firewall policy. The name should be unique in the organization in which the firewall policy is created. The name must be 1-63 characters long, and comply with RFC1035. Specifically, the name must be 1-63 characters long and match the regular expression [a-z]([-a-z0-9]*[a-z0-9])? which means the first character must be a lowercase letter, and all following characters must be a dash, lowercase letter, or digit, except the last character, which cannot be a dash.
  


- - -

* `description` -
  (Optional)
  An optional description of this resource. Provide this property when you create the resource.
  


## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `locations/global/firewallPolicies/{{name}}`

* `creation_timestamp` -
  Creation timestamp in RFC3339 text format.
  
* `fingerprint` -
  Fingerprint of the resource. This field is used internally during updates of this resource.
  
* `id` -
  The unique identifier for the resource. This identifier is defined by the server.
  
* `name` -
  Name of the resource. It is a numeric ID allocated by GCP which uniquely identifies the Firewall Policy.
  
* `rule_tuple_count` -
  Total count of all firewall policy rule tuples. A firewall policy can not exceed a set number of tuples.
  
* `self_link` -
  Server-defined URL for the resource.
  
* `self_link_with_id` -
  Server-defined URL for this resource with the resource id.
  
## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options: configuration options:

- `create` - Default is 20 minutes.
- `update` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import

FirewallPolicy can be imported using any of these accepted formats:

```
$ terraform import google_compute_firewall_policy.default locations/global/firewallPolicies/{{name}}
$ terraform import google_compute_firewall_policy.default {{name}}
```



