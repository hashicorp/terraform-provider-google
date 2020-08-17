---
subcategory: "Compute Engine"
layout: "google"
page_title: "Google: google_compute_shared_vpc_service_project"
sidebar_current: "docs-google-compute-shared-vpc-service-project"
description: |-
 Enables the Google Compute Engine Shared VPC feature for a project, assigning it as a service project.
---

# google_compute_shared_vpc_service_project

Enables the Google Compute Engine
[Shared VPC](https://cloud.google.com/compute/docs/shared-vpc)
feature for a project, assigning it as a Shared VPC service project associated
with a given host project.

For more information, see,
[the Project API documentation](https://cloud.google.com/compute/docs/reference/latest/projects),
where the Shared VPC feature is referred to by its former name "XPN".

## Example Usage

```hcl
resource "google_compute_shared_vpc_service_project" "service1" {
  host_project    = "host-project-id"
  service_project = "service-project-id-1"
}
```

For a complete Shared VPC example with both host and service projects, see
[`google_compute_shared_vpc_host_project`](/docs/providers/google/r/compute_shared_vpc_host_project.html).

## Argument Reference

The following arguments are expected:

* `host_project` - (Required) The ID of a host project to associate.

* `service_project` - (Required) The ID of the project that will serve as a Shared VPC service project.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `{{host_project}}/{{service_project}}`

## Timeouts

This resource provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 4 minutes.
- `delete` - Default is 4 minutes.

## Import

Google Compute Engine Shared VPC service project feature can be imported using the `host_project` and `service_project`, e.g.

```
$ terraform import google_compute_shared_vpc_service_project.service1 host-project-id/service-project-id-1
```
