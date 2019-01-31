---
layout: "google"
page_title: "Google: google_container_node_pool"
sidebar_current: "docs-google-container-node-pool"
description: |-
  Manages a GKE NodePool resource.
---

# google\_container\_node\_pool

Manages a node pool in a Google Kubernetes Engine (GKE) cluster separately from
the cluster control plane. For more information see [the official documentation](https://cloud.google.com/container-engine/docs/node-pools)
and [the API reference](https://cloud.google.com/container-engine/reference/rest/v1/projects.zones.clusters.nodePools).

### Example Usage - using a separately managed node pool (recommended)

```hcl
resource "google_container_cluster" "primary" {
  name   = "my-gke-cluster"
  region = "us-central1"
  
  # We can't create a cluster with no node pool defined, but we want to only use
  # separately managed node pools. So we create the smallest possible default
  # node pool and immediately delete it.
  remove_default_node_pool = true
  initial_node_count = 1
}

resource "google_container_node_pool" "primary_preemptible_nodes" {
  name       = "my-node-pool"
  region     = "us-central1"
  cluster    = "${google_container_cluster.primary.name}"
  node_count = 1

  node_config {
    preemptible  = true
    machine_type = "n1-standard-1"

    oauth_scopes = [
      "https://www.googleapis.com/auth/compute",
      "https://www.googleapis.com/auth/devstorage.read_only",
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
    ]
  }
}
```

### Example Usage - 2 node pools, 1 separately managed + the default node pool

```hcl
resource "google_container_node_pool" "np" {
  name       = "my-node-pool"
  zone       = "us-central1-a"
  cluster    = "${google_container_cluster.primary.name}"
  node_count = 3
  
  timeouts {
    create = "30m"
    update = "20m"
  }
}

resource "google_container_cluster" "primary" {
  name               = "marcellus-wallace"
  zone               = "us-central1-a"
  initial_node_count = 3

  additional_zones = [
    "us-central1-c",
  ]

  # Setting an empty username and password explicitly disables basic auth
  master_auth {
    username = ""
    password = ""
  }

  node_config {
    oauth_scopes = [
      "https://www.googleapis.com/auth/compute",
      "https://www.googleapis.com/auth/devstorage.read_only",
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
    ]

    guest_accelerator {
      type  = "nvidia-tesla-k80"
      count = 1
    }
  }
}

```

## Argument Reference

* `cluster` - (Required) The cluster to create the node pool for.  Cluster must be present in `zone` provided for zonal clusters.

- - -

* `zone` - (Optional) The zone in which the cluster resides.

* `region` - (Optional) The region in which the cluster resides (for regional clusters).

-> Note: You must be provide `region` for regional clusters and `zone` for zonal clusters

- - -

* `autoscaling` - (Optional) Configuration required by cluster autoscaler to adjust
    the size of the node pool to the current cluster usage. Structure is documented below.

* `initial_node_count` - (Optional) The initial node count for the pool. Changing this will force
    recreation of the resource.

* `management` - (Optional) Node management configuration, wherein auto-repair and
    auto-upgrade is configured. Structure is documented below.

* `max_pods_per_node` - (Optional, [Beta](https://terraform.io/docs/providers/google/provider_versions.html)) The maximum number of pods per node in this node pool.
    Note that this does not work on node pools which are "route-based" - that is, node
    pools belonging to clusters that do not have IP Aliasing enabled.

* `name` - (Optional) The name of the node pool. If left blank, Terraform will
    auto-generate a unique name.

* `node_config` - (Optional) The node configuration of the pool. See
    [google_container_cluster](container_cluster.html) for schema.

* `node_count` - (Optional) The number of nodes per instance group. This field can be used to
    update the number of nodes per instance group but should not be used alongside `autoscaling`.

* `project` - (Optional) The ID of the project in which to create the node pool. If blank,
    the provider-configured project will be used.

* `version` - (Optional) The Kubernetes version for the nodes in this pool. Note that if this field
    and `auto_upgrade` are both specified, they will fight each other for what the node version should
    be, so setting both is highly discouraged.

The `autoscaling` block supports:

* `min_node_count` - (Required) Minimum number of nodes in the NodePool. Must be >=1 and
    <= `max_node_count`.

* `max_node_count` - (Required) Maximum number of nodes in the NodePool. Must be >= min_node_count.

The `management` block supports:

* `auto_repair` - (Optional) Whether the nodes will be automatically repaired.

* `auto_upgrade` - (Optional) Whether the nodes will be automatically upgraded.

<a id="timeouts"></a>
## Timeouts

`google_container_node_pool` provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - (Default `30 minutes`) Used for adding node pools
- `update` - (Default `10 minutes`) Used for updates to node pools
- `delete` - (Default `10 minutes`) Used for removing node pools.

## Import

Node pools can be imported using the `project`, `zone`, `cluster` and `name`. If
the project is omitted, the default provider value will be used. Examples:

```
$ terraform import google_container_node_pool.mainpool my-gcp-project/us-east1-a/my-cluster/main-pool

$ terraform import google_container_node_pool.mainpool us-east1-a/my-cluster/main-pool
```
