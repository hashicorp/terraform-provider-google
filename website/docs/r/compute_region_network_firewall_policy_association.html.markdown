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
page_title: "Google: google_compute_region_network_firewall_policy_association"
description: |-
  The Compute NetworkFirewallPolicyAssociation resource
---

# google_compute_region_network_firewall_policy_association

The Compute NetworkFirewallPolicyAssociation resource

## Example Usage - regional
```hcl
resource "google_compute_region_network_firewall_policy" "basic_regional_network_firewall_policy" {
  name = "policy"
  project = "my-project-name"
  description = "Sample global network firewall policy"
  region = "us-west1"
}

resource "google_compute_network" "basic_network" {
  name = "network"
}

resource "google_compute_region_network_firewall_policy_association" "primary" {
  name = "association"
  attachment_target = google_compute_network.basic_network.id
  firewall_policy =  google_compute_region_network_firewall_policy.basic_regional_network_firewall_policy.name
  project =  "my-project-name"
  region = "us-west1"
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

* `region` -
  (Optional)
  The location of this resource.
  
* `project` -
  (Optional)
  The project for the resource
  


## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/regions/{{region}}/firewallPolicies/{{firewall_policy}}/associations/{{name}}`

* `short_name` -
  The short name of the firewall policy of the association.
  
## Timeouts

This resource provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import

NetworkFirewallPolicyAssociation can be imported using any of these accepted formats:

```
$ terraform import google_compute_region_network_firewall_policy_association.default projects/{{project}}/regions/{{region}}/firewallPolicies/{{firewall_policy}}/associations/{{name}}
$ terraform import google_compute_region_network_firewall_policy_association.default {{project}}/{{region}}/{{firewall_policy}}/{{name}}
```



