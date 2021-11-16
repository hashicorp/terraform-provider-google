---
subcategory: "Cloud Platform"
layout: "google"
page_title: "Google: google_compute_shared_vpc_host_project"
sidebar_current: "docs-google-datasource-compute-shared-vpc-host-project"
description: |-
  Retrieve shared VPC host project id
---

# google\_compute\_shared\_vpc\_host\_project

Use this data source to get ID of a host project serving a Shared VPC with a project.
For more information see
[API](https://cloud.google.com/compute/docs/reference/rest/v1/projects/getXpnHost)

## Example Usage

```hcl
data "google_compute_shared_vpc_host_project" "host" {
}

output "host_project" {
  value = data.google_compute_shared_vpc_host_project.host.host_project
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Optional) The ID of the project in which the resource belongs. If it is not provided, the provider project is used.


## Attributes Reference

The following attributes are exported:

* `host_project` - The ID of a host project that is sharing a VPC with the `project`. 
Empty string if no VPC is being shared.  
