---
subcategory: "Cloud Platform"
description: |-
 Verify the API service for the Google Cloud Platform project to see if it is enabled or not.
---

# google_project_service

Verify the API service for the Google Cloud Platform project to see if it is enabled or not.

For a list of services available, visit the [API library page](https://console.cloud.google.com/apis/library)
or run `gcloud services list --available`.

This datasource requires the [Service Usage API](https://console.cloud.google.com/apis/library/serviceusage.googleapis.com)
to use.


To get more information about `google_project_service`, see:

* [API documentation](https://cloud.google.com/service-usage/docs/reference/rest/v1/services)
* How-to Guides
    * [Enabling and Disabling Services](https://cloud.google.com/service-usage/docs/enable-disable)

## Example Usage

```hcl
data "google_project_service" "my-project-service" {
  service = "my-project-service"
}
```

## Argument Reference

The following arguments are supported:

* `service` - (Required) The name of the Google Platform project service.

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

See [google_project_service](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/google_project_service#argument-reference) resource for details of the available attributes.