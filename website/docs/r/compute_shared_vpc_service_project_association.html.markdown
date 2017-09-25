---
layout: "google"
page_title: "Google: google_compute_shared_vpc_service_project_association"
sidebar_current: "docs-google-compute-shared-vpc-service-project-association"
description: |-
 Allows associating a service project with a Shared VPC host project.
---

# google\_compute\_shared\_vpc\_service\_project\_association

Allows associating a service project with a Shared VPC host project. For more information see
[the official documentation](https://cloud.google.com/compute/docs/shared-vpc)
and
[API](https://cloud.google.com/compute/docs/reference/latest/projects).

~> **NOTE on Shared VPCs and Shared VPC Service Project Associations:** Terraform provides
both a standalone Shared VPC Service Project Association resource (an association between a Shared VPC host project
and a single `service_project`) and a [Shared VPC](compute_shared_vpc.html) resource with a `service_projects`
attribute. Do not use the same service project ID in both a Shared VPC resource and a Shared VPC Service
Project Association resource. Doing so will cause a conflict of associations and will overwrite the association.

## Example Usage

```hcl
resource "google_compute_shared_vpc_service_project_association" "sp" {
  host_project    = "host-project-id"
  service_project = "service-project-id"
}
```

## Argument Reference

The following arguments are supported:

* `host_project` - (Required) The Shared VPC host project ID.

* `service_project` - (Required) The ID of the service project to enable as a Shared VPC resource for this host project.
