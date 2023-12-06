---
subcategory: "Cloud VMware Engine"
description: |-
  Get information about a external address.
---

# google\_vmwareengine\_external_address

Use this data source to get details about a external address resource.

To get more information about external address, see:
* [API documentation](https://cloud.google.com/vmware-engine/docs/reference/rest/v1/projects.locations.privateClouds.externalAddresses)

## Example Usage

```hcl
data "google_vmwareengine_external_address" "my_external_address" {
  name     = "my-external-address"
  parent   = "project/my-project/locations/us-west1-a/privateClouds/my-cloud"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the resource.
* `parent` - (Required) The resource name of the private cloud that this cluster belongs.

## Attributes Reference

See [google_vmwareengine_external_address](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/vmwareengine_external_address#attributes-reference) resource for details of all the available attributes.