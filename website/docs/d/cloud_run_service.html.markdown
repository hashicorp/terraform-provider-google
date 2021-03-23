---
subcategory: "Cloud Run"
layout: "google"
page_title: "Google: google_cloud_run_service"
sidebar_current: "docs-google-datasource-cloud-run-service"
description: |-
  Get information about a Google Cloud Run Service.
---

# google\_cloud\_run\_service

Get information about a Google Cloud Run Service. For more information see
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

* `name` - (Required) The name of the Cloud Run Service.

* `location` - (Required) The location of the cloud run instance. eg us-central1

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

See [google_cloud_run_service](https://www.terraform.io/docs/providers/google/r/cloud_run_service.html#argument-reference) resource for details of the available attributes.
