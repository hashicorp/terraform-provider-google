---
layout: "google"
page_title: "Google: google_compute_shared_vpc"
sidebar_current: "docs-google-compute-shared-vpc"
description: |-
 Allows setting up Shared VPC in a Google Cloud Platform project.
---

# google\_compute\_shared\_vpc

Allows setting up Shared VPC in a Google Cloud Platform project. For more information see
[the official documentation](https://cloud.google.com/compute/docs/shared-vpc)
and
[API](https://cloud.google.com/compute/docs/reference/latest/projects).

~> **NOTE on Shared VPCs and Shared VPC Service Project Associations:** Terraform provides
both a standalone [Shared VPC Service Project Association](compute_shared_vpc_service_project_association.html)
resource (an association between a Shared VPC host project and a single `service_project`) and a Shared VPC resource
with a `service_projects` attribute. Do not use the same service project ID in both a Shared VPC resource and a
Shared VPC Service Project Association resource. Doing so will cause a conflict of associations and will overwrite the association.

## Example Usage

```hcl
resource "google_compute_shared_vpc" "vpc" {
  host_project     = "your-project-id"
  service_projects = ["service-project-1", "service-project-2"]
}
```

## Argument Reference

The following arguments are supported:

* `host_project` - (Required) The host project ID.

- - -

* `service_projects` - (Optional) List of IDs of service projects to enable as Shared VPC resources for this host.
