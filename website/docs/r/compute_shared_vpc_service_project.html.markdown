---
layout: "google"
page_title: "Google: google_compute_shared_vpc_service_project"
sidebar_current: "docs-google-compute-shared-vpc-service-project"
description: |-
 Allows enabling and disabling Shared VPC for a service Google Cloud Platform project.
---

# google\_compute\_shared\_vpc\_service\_project

Allows enabling and disabling Shared VPC for a service Google Cloud Platform project. For more information see
[the official documentation](https://cloud.google.com/compute/docs/shared-vpc)
and
[API](https://cloud.google.com/compute/docs/reference/latest/projects).

## Example Usage

```hcl
resource "google_compute_shared_vpc_host_project" "host" {
  project     = "your-project-id"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required) The host project ID.
