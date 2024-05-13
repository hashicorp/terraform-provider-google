---
subcategory: "Cloud VMware Engine"
description: |-
  Get info about a private cloud subnet.
---

# google_vmwareengine_subnet

Use this data source to get details about a subnet. Management subnets support only read operations and should be configured through this data source. User defined subnets can be configured using the resource as well as the datasource.

To get more information about private cloud subnet, see:
* [API documentation](https://cloud.google.com/vmware-engine/docs/reference/rest/v1/projects.locations.privateClouds.subnets)

## Example Usage

```hcl
data "google_vmwareengine_subnet" "my_subnet" {
  name     = "service-1"
  parent   = "project/my-project/locations/us-west1-a/privateClouds/my-cloud"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the resource. 
UserDefined subnets are named in the format of "service-n", where n ranges from 1 to 5. 
Management subnets have arbitary names including "vmotion", "vsan", "system-management" etc. More details about subnet names can be found on the cloud console.
* `parent` - (Required) The resource name of the private cloud that this subnet belongs.

## Attributes Reference

See [google_vmwareengine_subnet](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/vmwareengine_subnet#attributes-reference) resource for details of all the available attributes.