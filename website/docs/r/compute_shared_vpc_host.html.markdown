---
layout: "google"
page_title: "Google: google_compute_shared_vpc_host"
sidebar_current: "docs-google-compute-shared-vpc-host"
description: |-
 Allows setting a Google Cloud Platform project to be a Shared VPC Host.
---

# google\_compute\_shared\_vpc\_host

Allows setting a Google Cloud Platform project to be a Shared VPC Host. For more information see
[the official documentation](https://cloud.google.com/compute/docs/shared-vpc)
and
[API](https://cloud.google.com/compute/docs/reference/latest/projects/enableXpnHost).

## Example Usage

```hcl
resource "google_shared_vpc_host" "host" {
  project = "your-project-id"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required) The project ID.
