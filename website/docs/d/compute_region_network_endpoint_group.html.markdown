---
subcategory: "Compute Engine"
page_title: "Google: google_compute_region_network_endpoint_group"
description: |-
  Retrieve Region Network Endpoint Group's details.
---

# google\_compute\_region\_network\_endpoint\_group

Use this data source to access a Region Network Endpoint Group's attributes.

The RNEG may be found by providing either a `self_link`, or a `name` and a `region`.

## Example Usage

```hcl
data "google_compute_region_network_endpoint_group" "rneg1" {
  name = "k8s1-abcdef01-myns-mysvc-8080-4b6bac43"
  region = "us-central1"
}

data "google_compute_region_network_endpoint_group" "rneg2" {
  self_link = "https://www.googleapis.com/compute/v1/projects/myproject/regions/us-central1/networkEndpointGroups/k8s1-abcdef01-myns-mysvc-8080-4b6bac43"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Optional) The ID of the project to list versions in. If it is not provided, the provider project is used.

* `name` - (Optional) The Network Endpoint Group name. Provide either this or a `self_link`.

* `region` - (Optional) A reference to the region where the Serverless REGs Reside. Provide either this or a `self_link`.

* `self_link` - (Optional) The Network Endpoint Group self\_link.

## Attributes Reference

In addition the arguments listed above, the following attributes are exported:
* `id` - an identifier for the resource with format projects/{{project}}/regions/{{region}}/networkEndpointGroups/{{name}}
* `network` - The network to which all network endpoints in the RNEG belong.
* `subnetwork` - subnetwork to which all network endpoints in the RNEG belong.
* `description` - The RNEG description.
* `network_endpoint_type` - Type of network endpoints in this network endpoint group.
* `psc_target_service` - The target service url used to set up private service connection to a Google API or a PSC Producer Service Attachment.
* `default_port` - The RNEG default port.
* `size` - Number of network endpoints in the network endpoint group.
