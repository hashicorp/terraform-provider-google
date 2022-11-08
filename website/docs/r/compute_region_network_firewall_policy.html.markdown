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
page_title: "Google: google_compute_region_network_firewall_policy"
description: |-
  The Compute NetworkFirewallPolicy resource
---

# google_compute_region_network_firewall_policy

The Compute NetworkFirewallPolicy resource

## Example Usage - regional
```hcl
resource "google_compute_region_network_firewall_policy" "primary" {
  name = "policy"
  project = "my-project-name"
  description = "Sample regional network firewall policy"
  region = "us-west1"
}


```

## Argument Reference

The following arguments are supported:

* `name` -
  (Required)
  User-provided name of the Network firewall policy. The name should be unique in the project in which the firewall policy is created. The name must be 1-63 characters long, and comply with RFC1035. Specifically, the name must be 1-63 characters long and match the regular expression [a-z]([-a-z0-9]*[a-z0-9])? which means the first character must be a lowercase letter, and all following characters must be a dash, lowercase letter, or digit, except the last character, which cannot be a dash.
  


- - -

* `description` -
  (Optional)
  An optional description of this resource. Provide this property when you create the resource.
  
* `region` -
  (Optional)
  The location of this resource.
  
* `project` -
  (Optional)
  The project for the resource
  


## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/regions/{{region}}/firewallPolicies/{{name}}`

* `creation_timestamp` -
  Creation timestamp in RFC3339 text format.
  
* `fingerprint` -
  Fingerprint of the resource. This field is used internally during updates of this resource.
  
* `region_network_firewall_policy_id` -
  The unique identifier for the resource. This identifier is defined by the server.
  
* `rule_tuple_count` -
  Total count of all firewall policy rule tuples. A firewall policy can not exceed a set number of tuples.
  
* `self_link` -
  Server-defined URL for the resource.
  
* `self_link_with_id` -
  Server-defined URL for this resource with the resource id.
  
## Timeouts

This resource provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 20 minutes.
- `update` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import

NetworkFirewallPolicy can be imported using any of these accepted formats:

```
$ terraform import google_compute_region_network_firewall_policy.default projects/{{project}}/regions/{{region}}/firewallPolicies/{{name}}
$ terraform import google_compute_region_network_firewall_policy.default {{project}}/{{region}}/{{name}}
$ terraform import google_compute_region_network_firewall_policy.default {{region}}/{{name}}
$ terraform import google_compute_region_network_firewall_policy.default {{name}}
```



