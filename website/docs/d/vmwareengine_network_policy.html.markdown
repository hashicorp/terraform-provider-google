---
subcategory: "Cloud VMware Engine"
description: |-
  Get information about a network policy.
---

# google_vmwareengine_network_policy

Use this data source to get details about a network policy resource.

To get more information about network policy, see:
* [API documentation](https://cloud.google.com/vmware-engine/docs/reference/rest/v1/projects.locations.networkPolicies)

## Example Usage

```hcl
data "google_vmwareengine_network_policy" "my_network_policy" {
  name     = "my-network-policy"
  location = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the resource.
* `location` - (Required) Location of the resource.

## Attributes Reference

See [google_vmwareengine_network_policy](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/vmwareengine_network_policy#attributes-reference) resource for details of all the available attributes.