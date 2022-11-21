---
subcategory: "Compute Engine"
page_title: "Google: google_compute_backend_bucket"
description: |-
  Get information about a BackendBucket.
---

# google\_compute\_backend\_bucket

Get information about a BackendBucket.

## Example Usage

```tf
data "google_compute_backend_bucket" "my-backend-bucket" {
  name = "my-backend"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the resource.

- - -

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

See [google_compute_backend_bucket](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_backend_bucket) resource for details of the available attributes.
