---
subcategory: "Compute Engine"
layout: "google"
page_title: "Google: google_compute_backend_bucket"
sidebar_current: "docs-google-datasource-compute-backend-bucket"
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

See [google_compute_backend_bucket](https://www.terraform.io/docs/providers/google/r/compute_backend_bucket.html) resource for details of the available attributes.
