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
  Applies a hierarchical firewall policy to a target resource
---

# google_compute_firewall_policy_association

Allows associating hierarchical firewall policies with the target where they are applied. This allows creating policies and rules in a different location than they are applied.

For more information on applying hierarchical firewall policies see the [official documentation](https://cloud.google.com/vpc/docs/firewall-policies#managing_hierarchical_firewall_policy_resources)

## Example Usage

```hcl
resource "google_compute_firewall_policy" "default" {
  parent      = "organizations/12345"
  short_name  = "my-policy"
  description = "Example Resource"
}

resource "google_compute_firewall_policy_association" "default" {
  firewall_policy = google_compute_firewall_policy.default.id
  attachment_target = google_folder.folder.name
  name = "my-association"
}
```


## Argument Reference

The following arguments are supported:

* `attachment_target` -
  (Required)
  The target that the firewall policy is attached to.
  
* `firewall_policy` -
  (Required)
  The firewall policy ID of the association.
  
* `name` -
  (Required)
  The name for an association.
  


- - -



## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `locations/global/firewallPolicies/{{firewall_policy}}/associations/{{name}}`

* `short_name` -
  The short name of the firewall policy of the association.
  
## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options: configuration options:

- `create` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import

FirewallPolicyAssociation can be imported using any of these accepted formats:

* `locations/global/firewallPolicies/{{firewall_policy}}/associations/{{name}}`
* `{{firewall_policy}}/{{name}}`

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import FirewallPolicyAssociation using one of the formats above. For example:

```tf
import {
  id = "locations/global/firewallPolicies/{{firewall_policy}}/associations/{{name}}"
  to = google_compute_firewall_policy_association.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), FirewallPolicyAssociation can be imported using one of the formats above. For example:

```
$ terraform import google_compute_firewall_policy_association.default locations/global/firewallPolicies/{{firewall_policy}}/associations/{{name}}
$ terraform import google_compute_firewall_policy_association.default {{firewall_policy}}/{{name}}
```



