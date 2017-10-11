---
layout: "google"
page_title: "Google: google_compute_shared_vpc_host_project"
sidebar_current: "docs-google-compute-shared-vpc-host-project"
description: |-
 Allows enabling and disabling Shared VPC for the host Google Cloud Platform project.
---

# google\_compute\_shared\_vpc\_host\_project

Allows enabling and disabling Shared VPC for the host Google Cloud Platform project. For more information see
[the official documentation](https://cloud.google.com/compute/docs/shared-vpc)
and
[API](https://cloud.google.com/compute/docs/reference/latest/projects).

## Example Usage

```hcl
resource "google_compute_shared_vpc_host_project" "host" {
  project = "your-host-project-id"
}

resource "google_compute_shared_vpc_service_project" "service1" {
  project    = "your-service-project-id-1"
  // The host project must enable shared VPC first
  depends_on = ["google_compute_shared_vpc_host_project.host"]
}

resource "google_compute_shared_vpc_service_project" "service2" {
  project    = "your-service-project-id-2"
  // The host project must enable shared VPC first
  depends_on = ["google_compute_shared_vpc_host_project.host"]
}
```

## Argument Reference

The following arguments are supported:

* `host_project` - (Required) The host project ID.

* `service_project` - (Required) The service project ID.
