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
