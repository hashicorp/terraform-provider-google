---
subcategory: "Cloud Run"
layout: "google"
page_title: "Google: google_cloud_run_service"
sidebar_current: "docs-google-cloud-run-service"
description: |-
  Get information about a Google Cloud Run Service.
---

# google\_cloud\_run\_service

Get information about a Google Cloud Run. For more information see
the [official documentation](https://cloud.google.com/run/docs/)
and [API](https://cloud.google.com/run/docs/apis).

## Example Usage

```hcl
data "google_cloud_run_service" "run-service" {
  name = "my-service"
  location = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of a Cloud Function.
* `location` - (Required) The name of a Cloud Function.

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

* `region` - (Optional) The region in which the resource belongs. If it
    is not provided, the provider region is used.

## Attributes Reference

TBD
