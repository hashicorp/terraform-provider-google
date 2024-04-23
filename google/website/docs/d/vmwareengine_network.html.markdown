---
subcategory: "Cloud VMware Engine"
description: |-
  Get information about a VMwareEngine network.
---

# google\_vmwareengine\_network

Use this data source to get details about a VMwareEngine network resource.

To get more information about VMwareEngine Network, see:
* [API documentation](https://cloud.google.com/vmware-engine/docs/reference/rest/v1/projects.locations.vmwareEngineNetworks)

## Example Usage

```hcl
data "google_vmwareengine_network" "my_nw" {
  name     = "us-central1-default"
  location = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the resource.
* `location` - (Required) Location of the resource.

- - -

* `project` - (Optional) The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

## Attributes Reference

See [google_vmwareengine_network](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/vmwareengine_network#attributes-reference) resource for details of all the available attributes.
