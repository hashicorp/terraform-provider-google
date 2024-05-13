---
subcategory: "Compute Engine"
description: |-
  Get information about a Google Compute Router NAT.
---

# google_compute_router_nat

To get more information about Snapshot, see:

* [API documentation](https://cloud.google.com/compute/docs/reference/rest/v1/routers)
* How-to Guides
    * [Official Documentation](https://cloud.google.com/router/docs/)

## Example Usage

```hcl
data "google_compute_router_nat" "foo" {
  name = "my-nat"
  router = "my-router"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the NAT service. The name must be 1-63 characters long and
  comply with RFC1035.

* `router` - (Required)
  The name of the Cloud Router in which this NAT will be configured.

- - -

* `project` - (Optional) The ID of the project in which the resource belongs.
    If it is not provided, the provider project is used.

* `region` - (Optional) Region where the router and NAT reside.

## Attributes Reference

See [google_compute_router_nat](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_router_nat) resource for details of the available attributes.
