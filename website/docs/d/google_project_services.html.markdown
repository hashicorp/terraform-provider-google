---
subcategory: "Google Cloud Platform"
layout: "google"
page_title: "Google: google_project_services"
sidebar_current: "docs-google-datasource-project-services"
description: |-
  Retrieve enabled of API services for a Google Cloud Platform project
---

# google\_project\_services

Use this data source to get details on the enabled project services.

For a list of services available, visit the
[API library page](https://console.cloud.google.com/apis/library) or run `gcloud services list`.

## Example Usage

```hcl
data "google_project_services" "project" {
  project = "your-project-id"
}

output "project_services" {
  value = "${join(",", data.google_project_services.project.services)}"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required) The project ID.


## Attributes Reference

The following attributes are exported:

See [google_project_services](https://www.terraform.io/docs/providers/google/r/google_project_services.html) resource for details of the available attributes.
