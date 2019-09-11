---
layout: "google"
page_title: "Google: google_project_services"
sidebar_current: "docs-google-project-services"
description: |-
 Allows management of API services for a Google Cloud Platform project.
---

# google\_project\_services

Allows management of enabled API services for an existing Google Cloud
Platform project. Services in an existing project that are not defined
in the config will be removed, except for those services that are well
known to be required by a service defined in the config.

For a list of services available, visit the
[API library page](https://console.cloud.google.com/apis/library) or run `gcloud services list`.

~> **Note:** This resource attempts to be the authoritative source on *all* enabled APIs, which often
	leads to conflicts when certain actions enable other APIs. If you do not need to ensure that
	*exclusively* a particular set of APIs are enabled, you should most likely use the
	[google_project_service](google_project_service.html) resource, one resource per API.

## Example Usage

```hcl
resource "google_project_services" "project" {
  project = "your-project-id"
  services   = ["iam.googleapis.com", "cloudresourcemanager.googleapis.com"]
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required) The project ID.
    Changing this forces Terraform to attempt to disable all previously managed
    API services in the previous project.

* `services` - (Required) The list of services that are enabled. Supports
    update.

* `disable_on_destroy` - (Optional) Whether or not to disable APIs on project
    when destroyed. Defaults to true. **Note**: When `disable_on_destroy` is
    true and the project is changed, Terraform will force disable API services
    managed by Terraform for the previous project.

* `ignore_auto_enabled` - (Optional) Whether or not to ignore APIs which are
    known to be auto-enabled by GCP when enabling the configured APIs.
    Defaults to true. **Note**: The Google Cloud provider does not know
    about all service dependencies. Also, service dependencies can change
    over time. If needed, or if you want to do so for documentation
    purposes, you can always add the auto-enabled APIs to your config, too.

## Import

Project services can be imported using the `project_id`, e.g.

```
$ terraform import google_project_services.my_project your-project-id
```
