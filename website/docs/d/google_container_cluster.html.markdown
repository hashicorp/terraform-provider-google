---
layout: "google"
page_title: "Google: google_container_cluster"
sidebar_current: "docs-google-datasource-container-cluster"
description: |-
  Get info about a Google Kubernetes cluster.
---

# google\_container\_cluster

Get info about a cluster within GKE from its name and zone.

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

output "endpoint" {
  value = "${data.google_container_cluster.my_cluster.endpoint}"
}

output "instance_group_urls" {
  value = "${data.google_container_cluster.my_cluster.instance_group_urls}"
}

output "node_config" {
  value = "${data.google_container_cluster.my_cluster.node_config}"
}

output "node_pools" {
  value = "${data.google_container_cluster.my_cluster.node_pool}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - The name of the cluster.

* `zone` or `region` - The zone or region this cluster has been created in.

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

See [google_container_cluster](https://www.terraform.io/docs/providers/google/r/container_cluster.html) resource for details of the available attributes.