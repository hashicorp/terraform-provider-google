---
subcategory: "Cloud Run (v2 API)"
description: |-
  Get information about a Google Cloud Run v2 Worker Pool.
---

# google_cloud_run_v2_worker_pool

Get information about a Google Cloud Run v2 Worker Pool. For more information see
the [official documentation](https://cloud.google.com/run/docs/)
and [API](https://cloud.google.com/run/docs/apis).

## Example Usage

```hcl
data "google_cloud_run_v2_worker_pool" "my_worker_pool" {
  name     = "my-worker-pool"
  location = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Cloud Run v2 Worker Pool.

* `location` - (Required) The location of the instance. eg us-central1

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

See [google_cloud_run_v2_worker_pool](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/cloud_run_v2_worker_pool#argument-reference) resource for details of the available attributes.
