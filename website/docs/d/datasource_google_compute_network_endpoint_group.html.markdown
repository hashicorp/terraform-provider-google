---
subcategory: "Compute Engine"
layout: "google"
page_title: "Google: google_compute_network_endpoint_group"
sidebar_current: "docs-google-datasource-compute-network-endpoint-group"
description: |-
  Retrieve Network Endpoint Group's details.
---

# google\_compute\_network\_endpoint\_group

Use this data source to access a Network Endpoint Group's attributes.

The NEG may be found by providing either a `self_link`, or a `name` and a `zone`.

## Example Usage

```hcl
data "google_compute_network_endpoint_group" "neg1" {
  name = "k8s1-abcdef01-myns-mysvc-8080-4b6bac43"
  zone = "us-central1-a"
}

data "google_compute_network_endpoint_group" "neg2" {
  self_link = "https://www.googleapis.com/compute/v1/projects/myproject/zones/us-central1-a/networkEndpointGroups/k8s1-abcdef01-myns-mysvc-8080-4b6bac43"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Optional) The ID of the project to list versions in.
    If it is not provided, the provider project is used.

* `name` - (Optional) The Network Endpoint Group name.
    Provide either this or a `self_link`.

* `zone` - (Optional) The Network Endpoint Group availability zone.

* `self_link` - (Optional) The Network Endpoint Group self\_link.

## Attributes Reference

In addition the arguments listed above, the following attributes are exported:

* `network` - The network to which all network endpoints in the NEG belong.
* `subnetwork` - subnetwork to which all network endpoints in the NEG belong.
* `description` - The NEG description.
* `network_endpoint_type` - Type of network endpoints in this network endpoint group.
* `default_port` - The NEG default port.
* `size` - Number of network endpoints in the network endpoint group.
