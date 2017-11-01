---
layout: "google"
page_title: "Google: google_project_service"
sidebar_current: "docs-google-project-service"
description: |-
 Allows management of a single API service for a Google Cloud Platform project.
---

# google\_project\_service

Allows management of a single API service for an existing Google Cloud Platform project. 

For a list of services available, visit the
[API library page](https://console.cloud.google.com/apis/library) or run `gcloud service-management list`.

~> **Note:** This resource _must not_ be used in conjunction with
   `google_project_services` or they will fight over which services should be enabled.

## Example Usage

```hcl
resource "google_project_service" "project" {
  project = "your-project-id"
  service = "iam.googleapis.com"
}
```

## Argument Reference

The following arguments are supported:

* `service` - (Required) The service to enable.

* `project` - (Optional) The project ID. If not provided, the provider project is used.
