---
subcategory: "Compute Engine"
page_title: "Google: google_compute_router"
description: |-
  Get a Cloud Router within GCE.
---

# google\_compute\_router

Get a router within GCE from its name and VPC.

## Example Usage

```hcl
data "google_compute_router" "my-router" {
  name   = "myrouter-us-east1"
  network = "my-network"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the router.

* `network` - (Required) The VPC network on which this router lives.

* `project` - (Optional) The ID of the project in which the resource
    belongs. If it is not provided, the provider project is used.

* `region` - (Optional) The region this router has been created in. If
    unspecified, this defaults to the region configured in the provider.


## Attributes Reference

See [google_compute_router](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_router) resource for details of the available attributes.
