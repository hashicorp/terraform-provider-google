---
subcategory: "Compute Engine"
layout: "google"
page_title: "Google: google_compute_shared_vpc_host_project"
sidebar_current: "docs-google-compute-shared-vpc-host-project"
description: |-
 Enables the Google Compute Engine Shared VPC feature for a project, assigning it as a host project.
---

# google_compute_shared_vpc_host_project

Enables the Google Compute Engine
[Shared VPC](https://cloud.google.com/compute/docs/shared-vpc)
feature for a project, assigning it as a Shared VPC host project.

For more information, see,
[the Project API documentation](https://cloud.google.com/compute/docs/reference/latest/projects),
where the Shared VPC feature is referred to by its former name "XPN".

## Example Usage

```hcl
# A host project provides network resources to associated service projects.
resource "google_compute_shared_vpc_host_project" "host" {
  project = "host-project-id"
}

# A service project gains access to network resources provided by its
# associated host project.
resource "google_compute_shared_vpc_service_project" "service1" {
  host_project    = google_compute_shared_vpc_host_project.host.project
  service_project = "service-project-id-1"
}

resource "google_compute_shared_vpc_service_project" "service2" {
  host_project    = google_compute_shared_vpc_host_project.host.project
  service_project = "service-project-id-2"
}
```

## Argument Reference

The following arguments are expected:

* `project` - (Required) The ID of the project that will serve as a Shared VPC host project

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `{{project}}`

## Timeouts

This resource provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 4 minutes.
- `delete` - Default is 4 minutes.

## Import

Google Compute Engine Shared VPC host project feature can be imported using the `project`, e.g.

```
$ terraform import google_compute_shared_vpc_host_project.host host-project-id
```
