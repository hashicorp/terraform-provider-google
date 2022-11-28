---
subcategory: "Kubernetes (Container) Engine"
page_title: "Google: google_container_node_pool"
description: |-
  Manages a GKE NodePool resource.
---

# google\_container\_node\_pool

-> See the [Using GKE with Terraform](/docs/providers/google/guides/using_gke_with_terraform.html)
guide for more information about using GKE with Terraform.

Manages a node pool in a Google Kubernetes Engine (GKE) cluster separately from
the cluster control plane. For more information see [the official documentation](https://cloud.google.com/container-engine/docs/node-pools)
and [the API reference](https://cloud.google.com/kubernetes-engine/docs/reference/rest/v1beta1/projects.locations.clusters.nodePools).

### Example Usage - using a separately managed node pool (recommended)

```hcl
resource "google_service_account" "default" {
  account_id   = "service-account-id"
  display_name = "Service Account"
}

resource "google_container_cluster" "primary" {
  name     = "my-gke-cluster"
  location = "us-central1"

  # We can't create a cluster with no node pool defined, but we want to only use
  # separately managed node pools. So we create the smallest possible default
  # node pool and immediately delete it.
  remove_default_node_pool = true
  initial_node_count       = 1
}

resource "google_container_node_pool" "primary_preemptible_nodes" {
  name       = "my-node-pool"
  cluster    = google_container_cluster.primary.id
  node_count = 1

  node_config {
    preemptible  = true
    machine_type = "e2-medium"

    # Google recommends custom service accounts that have cloud-platform scope and permissions granted via IAM Roles.
    service_account = google_service_account.default.email
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]
  }
}
```

### Example Usage - 2 node pools, 1 separately managed + the default node pool

```hcl
resource "google_service_account" "default" {
  account_id   = "service-account-id"
  display_name = "Service Account"
}

resource "google_container_node_pool" "np" {
  name       = "my-node-pool"
  cluster    = google_container_cluster.primary.id
  node_config {
    machine_type = "e2-medium"
    # Google recommends custom service accounts that have cloud-platform scope and permissions granted via IAM Roles.
    service_account = google_service_account.default.email
    oauth_scopes    = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]
  }
  timeouts {
    create = "30m"
    update = "20m"
  }
}

resource "google_container_cluster" "primary" {
  name               = "marcellus-wallace"
  location           = "us-central1-a"
  initial_node_count = 3

  node_locations = [
    "us-central1-c",
  ]

  node_config {
    # Google recommends custom service accounts that have cloud-platform scope and permissions granted via IAM Roles.
    service_account = google_service_account.default.email
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]
    guest_accelerator {
      type  = "nvidia-tesla-k80"
      count = 1
    }
  }
}
```

## Argument Reference

* `cluster` - (Required) The cluster to create the node pool for. Cluster must be present in `location` provided for clusters. May be specified in the format `projects/{{project}}/locations/{{location}}/clusters/{{cluster}}` or as just the name of the cluster.

- - -

* `location` - (Optional) The location (region or zone) of the cluster.

- - -

* `autoscaling` - (Optional) Configuration required by cluster autoscaler to adjust
    the size of the node pool to the current cluster usage. Structure is [documented below](#nested_autoscaling).

* `initial_node_count` - (Optional) The initial number of nodes for the pool. In
    regional or multi-zonal clusters, this is the number of nodes per zone. Changing
    this will force recreation of the resource. WARNING: Resizing your node pool manually
    may change this value in your existing cluster, which will trigger destruction
    and recreation on the next Terraform run (to rectify the discrepancy).  If you don't
    need this value, don't set it.  If you do need it, you can [use a lifecycle block to
    ignore subsequent changes to this field](https://github.com/hashicorp/terraform-provider-google/issues/6901#issuecomment-667369691).

* `management` - (Optional) Node management configuration, wherein auto-repair and
    auto-upgrade is configured. Structure is [documented below](#nested_management).

* `max_pods_per_node` - (Optional) The maximum number of pods per node in this node pool.
    Note that this does not work on node pools which are "route-based" - that is, node
    pools belonging to clusters that do not have IP Aliasing enabled.
    See the [official documentation](https://cloud.google.com/kubernetes-engine/docs/how-to/flexible-pod-cidr)
    for more information.

* `node_locations` - (Optional)
The list of zones in which the node pool's nodes should be located. Nodes must
be in the region of their regional cluster or in the same region as their
cluster's zone for zonal clusters. If unspecified, the cluster-level
`node_locations` will be used.

-> Note: `node_locations` will not revert to the cluster's default set of zones
upon being unset. You must manually reconcile the list of zones with your
cluster.

* `name` - (Optional) The name of the node pool. If left blank, Terraform will
    auto-generate a unique name.

* `name_prefix` - (Optional) Creates a unique name for the node pool beginning
    with the specified prefix. Conflicts with `name`.

* `node_config` - (Optional) Parameters used in creating the node pool. See
    [google_container_cluster](container_cluster.html#nested_node_config) for schema.

* `network_config` - (Optional) The network configuration of the pool. Such as
    configuration for [Adding Pod IP address ranges](https://cloud.google.com/kubernetes-engine/docs/how-to/multi-pod-cidr)) to the node pool. Or enabling private nodes. Structure is
    [documented below](#nested_network_config)

* `node_count` - (Optional) The number of nodes per instance group. This field can be used to
    update the number of nodes per instance group but should not be used alongside `autoscaling`.

* `project` - (Optional) The ID of the project in which to create the node pool. If blank,
    the provider-configured project will be used.

* `upgrade_settings` (Optional) Specify node upgrade settings to change how GKE upgrades nodes.
    The maximum number of nodes upgraded simultaneously is limited to 20. Structure is [documented below](#nested_upgrade_settings).

* `version` - (Optional) The Kubernetes version for the nodes in this pool. Note that if this field
    and `auto_upgrade` are both specified, they will fight each other for what the node version should
    be, so setting both is highly discouraged. While a fuzzy version can be specified, it's
    recommended that you specify explicit versions as Terraform will see spurious diffs
    when fuzzy versions are used. See the `google_container_engine_versions` data source's
    `version_prefix` field to approximate fuzzy versions in a Terraform-compatible way.

* `placement_policy` - (Optional, [Beta](https://terraform.io/docs/providers/google/provider_versions.html)) Specifies a custom placement policy for the
  nodes.

<a name="nested_autoscaling"></a>The `autoscaling` block supports (either total or per zone limits are required):

* `min_node_count` - (Optional) Minimum number of nodes per zone in the NodePool.
    Must be >=0 and <= `max_node_count`. Cannot be used with total limits.

* `max_node_count` - (Optional) Maximum number of nodes per zone in the NodePool.
    Must be >= min_node_count. Cannot be used with total limits.

* `total_min_node_count` - (Optional) Total minimum number of nodes in the NodePool.
    Must be >=0 and <= `total_max_node_count`. Cannot be used with per zone limits.
    Total size limits are supported only in 1.24.1+ clusters.

* `total_max_node_count` - (Optional) Total maximum number of nodes in the NodePool.
    Must be >= total_min_node_count. Cannot be used with per zone limits.
    Total size limits are supported only in 1.24.1+ clusters.

* `location_policy` - (Optional) Location policy specifies the algorithm used when
  scaling-up the node pool. Location policy is supported only in 1.24.1+ clusters.
    * "BALANCED" - Is a best effort policy that aims to balance the sizes of available zones.
    * "ANY" - Instructs the cluster autoscaler to prioritize utilization of unused reservations,
      and reduce preemption risk for Spot VMs.

<a name="nested_management"></a>The `management` block supports:

* `auto_repair` - (Optional) Whether the nodes will be automatically repaired.

* `auto_upgrade` - (Optional) Whether the nodes will be automatically upgraded.

<a name="nested_network_config"></a>The `network_config` block supports:

* `create_pod_range` - (Optional) Whether to create a new range for pod IPs in this node pool. Defaults are provided for `pod_range` and `pod_ipv4_cidr_block` if they are not specified.

* `enable_private_nodes` - (Optional) Whether nodes have internal IP addresses only.

* `pod_ipv4_cidr_block` - (Optional) The IP address range for pod IPs in this node pool. Only applicable if createPodRange is true. Set to blank to have a range chosen with the default size. Set to /netmask (e.g. /14) to have a range chosen with a specific netmask. Set to a CIDR notation (e.g. 10.96.0.0/14) to pick a specific range to use.

* `pod_range` - (Optional) The ID of the secondary range for pod IPs. If `create_pod_range` is true, this ID is used for the new range. If `create_pod_range` is false, uses an existing secondary range with this ID.

<a name="nested_upgrade_settings"></a>The `upgrade_settings` block supports:

* `max_surge` - (Optional) The number of additional nodes that can be added to the node pool during
    an upgrade. Increasing `max_surge` raises the number of nodes that can be upgraded simultaneously.
    Can be set to 0 or greater.

* `max_unavailable` - (Optional) The number of nodes that can be simultaneously unavailable during
    an upgrade. Increasing `max_unavailable` raises the number of nodes that can be upgraded in
    parallel. Can be set to 0 or greater.

`max_surge` and `max_unavailable` must not be negative and at least one of them must be greater than zero.

* `strategy` - (Default `SURGE`) The upgrade stragey to be used for upgrading the nodes.

* `blue_green_settings` - (Optional) The settings to adjust [blue green upgrades](https://cloud.google.com/kubernetes-engine/docs/concepts/node-pool-upgrade-strategies#blue-green-upgrade-strategy).
    Structure is [documented below](#nested_blue_green_settings)

<a name="nested_blue_green_settings"></a>The `blue_green_settings` block supports:

* `standard_rollout_policy` - (Required) Specifies the standard policy settings for blue-green upgrades.
    * `batch_percentage` - (Optional) Percentage of the blue pool nodes to drain in a batch.
    * `batch_node_count` - (Optional) Number of blue nodes to drain in a batch.
    * `batch_soak_duration` - (Optionial) Soak time after each batch gets drained.

* `node_pool_soak_duration` - (Optional) Time needed after draining the entire blue pool.
    After this period, the blue pool will be cleaned up.

<a name="nested_placement_policy"></a>The `placement_policy` block supports:

* `type` - (Required) The type of the policy. Supports a single value: COMPACT.
  Specifying COMPACT placement policy type places node pool's nodes in a closer
  physical proximity in order to reduce network latency between nodes.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `{{project}}/{{location}}/{{cluster}}/{{name}}`

* `instance_group_urls` - The resource URLs of the managed instance groups associated with this node pool.

* `managed_instance_group_urls` - List of instance group URLs which have been assigned to this node pool.

<a id="timeouts"></a>
## Timeouts

`google_container_node_pool` provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - (Default `30 minutes`) Used for adding node pools
- `update` - (Default `30 minutes`) Used for updates to node pools
- `delete` - (Default `30 minutes`) Used for removing node pools.

## Import

Node pools can be imported using the `project`, `location`, `cluster` and `name`. If
the project is omitted, the project value in the provider configuration will be used. Examples:

```
$ terraform import google_container_node_pool.mainpool my-gcp-project/us-east1-a/my-cluster/main-pool

$ terraform import google_container_node_pool.mainpool us-east1/my-cluster/main-pool
```
