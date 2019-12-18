---
subcategory: "Kubernetes (Container) Engine"
layout: "google"
page_title: "Google: google_container_cluster"
sidebar_current: "docs-google-datasource-container-cluster"
description: |-
  Get info about a Google Kubernetes Engine cluster.
---

# google\_container\_cluster

Get info about a GKE cluster from its name and location.

## Example Usage

```tf
data "google_container_cluster" "my_cluster" {
  name     = "my-cluster"
  location = "us-east1-a"
}

output "cluster_username" {
  value = data.google_container_cluster.my_cluster.master_auth[0].username
}

output "cluster_password" {
  value = data.google_container_cluster.my_cluster.master_auth[0].password
}

output "endpoint" {
  value = data.google_container_cluster.my_cluster.endpoint
}

output "instance_group_urls" {
  value = data.google_container_cluster.my_cluster.instance_group_urls
}

output "node_config" {
  value = data.google_container_cluster.my_cluster.node_config
}

output "node_pools" {
  value = data.google_container_cluster.my_cluster.node_pool
}
```

## Argument Reference

The following arguments are supported:

* `name` (Required) - The name of the cluster.

* `location` (Optional) - The location (zone or region) this cluster has been
created in. One of `location`, `region`, `zone`, or a provider-level `zone` must
be specified.

* `zone` (Optional) - The zone this cluster has been created in. Deprecated in
favour of `location`.

* `region` (Optional) - The region this cluster has been created in. Deprecated
in favour of `location`.

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

See [google_container_cluster](https://www.terraform.io/docs/providers/google/r/container_cluster.html) resource for details of the available attributes.
