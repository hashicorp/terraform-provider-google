---
subcategory: "Compute Engine"
description: |-
  Get information about a Regional Backend Service.
---

# google_compute_region_backend_service

Get information about a Regional Backend Service. For more information see
[the official documentation](https://cloud.google.com/compute/docs/load-balancing/internal/backend-service) and
[API](https://cloud.google.com/compute/docs/reference/rest/beta/regionBackendServices).

## Example Usage

```hcl
data "google_compute_region_backend_service" "my_backend" {
  name   = "my-backend-service"
  region = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the regional backend service.

* `region` - (Optional) The region where the backend service resides.

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

See [google_compute_region_backend_service](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_region_backend_service) resource for details of the available attributes.
