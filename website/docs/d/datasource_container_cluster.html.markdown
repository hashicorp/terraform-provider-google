---
layout: "google"
page_title: "Google: google_container_cluster"
sidebar_current: "docs-google-datasource-container-cluster"
description: |-
  Get info about a Google Kubernetes cluster.
---

# google\_container\_cluster

Get info about a cluster within GCE from its name and zone.

## Example Usage

```tf
data "google_container_cluster" "my_cluster" {
  name   = "my-cluster"
  zone   = "us-east1-a"
}

output "cluster_username" {
  value = "${data.google_container_cluster.my_cluster.master_auth.0.username}"
}

output "cluster_password" {
  value = "${data.google_container_cluster.my_cluster.master_auth.0.password}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - The name of the cluster.

* `zone` - The zones this cluster has been created in.

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `additional_zones` - The list of additional Google Compute Engine
    locations in which the cluster's nodes are located.

* `cluster_ipv4_cidr` - The IP address range of the container pods in
    this cluster.

* `endpoint` - The IP address of this cluster's Kubernetes master.

* `instance_group_urls` - List of instance group URLs which have been assigned
    to the cluster.

* `ip_cidr_range` - The IP address range that machines in this
    network are assigned to, represented as a CIDR block.

* `master_auth.0.client_certificate` - Base64 encoded public certificate
    used by clients to authenticate to the cluster endpoint.

* `master_auth.0.client_key` - Base64 encoded private key used by clients
    to authenticate to the cluster endpoint.

* `master_auth.0.cluster_ca_certificate` - Base64 encoded public certificate
    that is the root of trust for the cluster.

* `master_auth.0.password` - The password to use for HTTP basic
    authentication when accessin the Kubernetes master endpoint.

* `master_auth.0.username` - The username to use for HTTP basic
    authentication when accessin the Kubernetes master endpoint.

* `master_version` - The current version of the master in the cluster.

* `network` - The name or self_link of the Google Compute Engine
    network to which the cluster is connected.

* `node_version` - The Kubernetes version on the nodes.

* `subnetwork` - The name of the Google Compute Engine subnetwork in
    which the cluster's instances are launched.