---
subcategory: "App Hub"
description: |-
  Application is a functional grouping of Services and Workloads that helps achieve a desired end-to-end business functionality.
---

# google_apphub_application

Application is a functional grouping of Services and Workloads that helps achieve a desired end-to-end business functionality. Services and Workloads are owned by the Application.


## Example Usage


```hcl
data "google_apphub_application" "application" {
  project = "project-id"
  application_id = "application"
  location = "location"
}
```

## Argument Reference

See [google_resource_application](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/apphub_application#argument-reference) resource for details of the available attributes.

