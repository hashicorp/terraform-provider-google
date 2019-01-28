---
layout: "google"
page_title: "Google: google_project_service"
sidebar_current: "docs-google-project-service-x"
description: |-
 Allows management of a single API service for a Google Cloud Platform project.
---

# google\_project\_service

Allows management of a single API service for an existing Google Cloud Platform project. 

For a list of services available, visit the
[API library page](https://console.cloud.google.com/apis/library) or run `gcloud services list`.

~> **Note:** This resource _must not_ be used in conjunction with
   `google_project_services` or they will fight over which services should be enabled.

## Example Usage

```hcl
resource "google_project_service" "project" {
  project = "your-project-id"
  service = "iam.googleapis.com"

  disable_dependent_services = true
}
```

## Argument Reference

The following arguments are supported:

* `service` - (Required) The service to enable.

* `project` - (Optional) The project ID. If not provided, the provider project is used.

* `disable_dependent_services` - (Optional) If `true`, services that are enabled and which depend on this service should also be disabled when this service is destroyed.
If `false` or unset, an error will be generated if any enabled services depend on this service when destroying it.

* `disable_on_destroy` - (Optional) If true, disable the service when the terraform resource is destroyed.  Defaults to true.  May be useful in the event that a project is long-lived but the infrastructure running in that project changes frequently.

## Import

Project services can be imported using the `project_id` and `service`, e.g.

```
$ terraform import google_project_service.my_project your-project-id/iam.googleapis.com
```
