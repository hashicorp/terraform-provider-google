---
subcategory: "Compute Engine"
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
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options: configuration options:

- `create` - Default is 4 minutes.
- `delete` - Default is 4 minutes.

## Import

Google Compute Engine Shared VPC host project feature can be imported using `project`, e.g.

* `{{project_id}}`

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Google Compute Engine Shared VPC host projects using one of the formats above. For example:

```tf
import {
  id = "{{project_id}}"
  to = google_compute_shared_vpc_host_project.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), Google Compute Engine Shared VPC host projects can be imported using one of the formats above. For example:


```
$ terraform import google_compute_shared_vpc_host_project.default {{project_id}}
```
