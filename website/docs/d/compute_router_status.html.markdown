---
subcategory: "Compute Engine"
page_title: "Google: google_compute_router_status"
description: |-
  Get a Cloud Router's Status.
---

# google\_compute\_router\_status

Get a Cloud Router's status within GCE from its name and region. This data source exposes the
routes learned by a Cloud Router via BGP peers.

For more information see [the official documentation](https://cloud.google.com/network-connectivity/docs/router/how-to/viewing-router-details)
and
[API](https://cloud.google.com/compute/docs/reference/rest/v1/routers/getRouterStatus).

## Example Usage

```hcl
data "google_compute_router_status" "my-router" {
  name   = "myrouter"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the router.

* `project` - (Optional) The ID of the project in which the resource
    belongs. If it is not provided, the provider project is used.

* `region` - (Optional) The region this router has been created in. If
    unspecified, this defaults to the region configured in the provider.


## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `network` - The network name or resource link to the parent
    network of this subnetwork.

* `best_routes` - List of best `compute#routes` configurations for this router's network. See [google_compute_route](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_route) resource for available attributes.

* `best_routes_for_router` - List of best `compute#routes` for this specific router. See [google_compute_route](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_route) resource for available attributes.
