---
subcategory: "Compute Engine"
description: |-
  Get subnetworks within GCE.
---

# google\_compute\_subnetworks

Get subnetworks within GCE.
See [the official documentation](https://cloud.google.com/vpc/docs/subnets)
and [API](https://cloud.google.com/compute/docs/reference/rest/v1/subnetworks/list).

## Example Usage

```hcl
data "google_compute_subnetworks" "my-subnetworks" {
  filter  = "ipCidrRange eq 192.168.178.0/24"
  project = "my-project"
  region  = "us-east1"
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) - A string filter as defined in the [REST API](https://cloud.google.com/compute/docs/reference/rest/v1/subnetworks/list#query-parameters).

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

* `region` - (Optional) The region this subnetwork has been created in. If
    unspecified, this defaults to the region configured in the provider.

## Attributes Reference

* `subnetworks` - A list of all retrieved GCE subnetworks. Structure is [defined below](#nested_subnetworks).

<a name="nested_subnetworks"></a>The `subnetworks` block supports:

* `description` - Description of the subnetwork.
* `ip_cidr_range` - The IP address range represented as a CIDR block.
* `name` - The name of the subnetwork.
* `network` - The self link of the parent network.
* `network_name` - The name of the parent network computed from `network` attribute.
* `private_ip_google_access` - Whether the VMs in the subnet can access Google services without assigned external IP addresses.
* `self_link` - The self link of the subnetwork.
