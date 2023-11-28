---
subcategory: "Cloud VMware Engine"
description: |-
  Get information about a network peering.
---

# google\_vmwareengine\_network_peering

Use this data source to get details about a network peering resource.

To get more information about network peering, see:
* [API documentation](https://cloud.google.com/vmware-engine/docs/reference/rest/v1/projects.locations.networkPeerings)

## Example Usage

```hcl
data "google_vmwareengine_network_peering" "my_network_peering" {
  name     = "my-network-peering"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the resource.

## Attributes Reference

See [google_vmwareengine_network_peering](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/vmwareengine_network_peering#attributes-reference) resource for details of all the available attributes.