---
subcategory: "Cloud Run"
description: |-
  Get information about a Google Cloud Run v2 Job.
---

# google_cloud_run_v2_job

Get information about a Google Cloud Run v2 Job. For more information see
the [official documentation](https://cloud.google.com/run/docs/)
and [API](https://cloud.google.com/run/docs/apis).

## Example Usage

```hcl
data "google_cloud_run_v2_job" "my_job" {
  name     = "my-job"
  location = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Cloud Run v2 Job.

* `location` - (Required) The location of the instance. eg us-central1

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

See [google_cloud_run_v2_job](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/cloud_run_v2_job#argument-reference) resource for details of the available attributes.
