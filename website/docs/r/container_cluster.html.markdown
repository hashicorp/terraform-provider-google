---
layout: "google"
page_title: "Google: google_container_cluster"
sidebar_current: "docs-google-container-cluster"
description: |-
  Creates a Google Kubernetes Engine (GKE) cluster.
---

# google\_container\_cluster

Creates a Google Kubernetes Engine (GKE) cluster. For more information see
[the official documentation](https://cloud.google.com/container-engine/docs/clusters)
and
[API](https://cloud.google.com/container-engine/reference/rest/v1/projects.zones.clusters).

~> **Note:** All arguments including the username and password will be stored in the raw state as plain-text.
[Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Example usage

```hcl
resource "google_container_cluster" "primary" {
  name               = "marcellus-wallace"
  zone               = "us-central1-a"
  initial_node_count = 3

  additional_zones = [
    "us-central1-b",
    "us-central1-c",
  ]

  master_auth {
    username = "mr.yoda"
    password = "adoy.rm"
  }

  node_config {
    oauth_scopes = [
      "https://www.googleapis.com/auth/compute",
      "https://www.googleapis.com/auth/devstorage.read_only",
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
    ]

    labels {
      foo = "bar"
    }

    tags = ["foo", "bar"]
  }
}

# The following outputs allow authentication and connectivity to the GKE Cluster.
output "client_certificate" {
  value = "${google_container_cluster.primary.master_auth.0.client_certificate}"
}

output "client_key" {
  value = "${google_container_cluster.primary.master_auth.0.client_key}"
}

output "cluster_ca_certificate" {
  value = "${google_container_cluster.primary.master_auth.0.cluster_ca_certificate}"
}
```

## Argument Reference

* `name` - (Required) The name of the cluster, unique within the project and
    zone.

* `zone` - (Required) The zone that the master and the number of nodes specified
    in `initial_node_count` should be created in.

- - -

* `additional_zones` - (Optional) The list of additional Google Compute Engine
    locations in which the cluster's nodes should be located. If additional zones are
    configured, the number of nodes specified in `initial_node_count` is created in
    all specified zones.

* `addons_config` - (Optional) The configuration for addons supported by GKE.
    Structure is documented below.

* `cluster_ipv4_cidr` - (Optional) The IP address range of the kubernetes pods in
    this cluster. Default is an automatically assigned CIDR.

* `description` - (Optional) Description of the cluster.

* `enable_kubernetes_alpha` - (Optional) Whether to enable Kubernetes Alpha features for
    this cluster. Note that when this option is enabled, the cluster cannot be upgraded
    and will be automatically deleted after 30 days.

* `enable_legacy_abac` - (Optional) Whether the ABAC authorizer is enabled for this cluster.
    When enabled, identities in the system, including service accounts, nodes, and controllers,
    will have statically granted permissions beyond those provided by the RBAC configuration or IAM.

* `initial_node_count` - (Optional) The number of nodes to create in this
    cluster (not including the Kubernetes master). Must be set if `node_pool` is not set.

* `ip_allocation_policy` - (Optional) Configuration for cluster IP allocation. As of now, only pre-allocated subnetworks (custom type with secondary ranges) are supported.

* `logging_service` - (Optional) The logging service that the cluster should
    write logs to. Available options include `logging.googleapis.com` and
    `none`. Defaults to `logging.googleapis.com`

* `maintenance_policy` - (Optional) The maintenance policy to use for the cluster. Structure is
    documented below.

* `master_auth` - (Optional) The authentication information for accessing the
    Kubernetes master. Structure is documented below.

* `master_authorized_networks_config` - (Optional) The desired configuration options
    for master authorized networks. Omit the nested `cidr_blocks` attribute to disallow
    external access (except the cluster node IPs, which GKE automatically whitelists).

* `min_master_version` - (Optional) The minimum version of the master. GKE
    will auto-update the master to new versions, so this does not guarantee the
    current master version--use the read-only `master_version` field to obtain that.
    If unset, the cluster's version will be set by GKE to the version of the most recent
    official release (which is not necessarily the latest version).

* `monitoring_service` - (Optional) The monitoring service that the cluster
    should write metrics to. Available options include
    `monitoring.googleapis.com` and `none`. Defaults to
    `monitoring.googleapis.com`

* `network` - (Optional) The name or self_link of the Google Compute Engine
    network to which the cluster is connected.

* `network_policy` - (Optional) Configuration options for the
    [NetworkPolicy](https://kubernetes.io/docs/concepts/services-networking/networkpolicies/)
    feature. Structure is documented below.

* `node_config` -  (Optional) Parameters used in creating the cluster's nodes.
    Structure is documented below.

* `node_pool` - (Optional) List of node pools associated with this cluster.
    See [google_container_node_pool](container_node_pool.html) for schema.

* `node_version` - (Optional) The Kubernetes version on the nodes. Must either be unset
    or set to the same value as `min_master_version` on create. Defaults to the default
    version set by GKE which is not necessarily the latest version.

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

* `subnetwork` - (Optional) The name of the Google Compute Engine subnetwork in
    which the cluster's instances are launched.

The `addons_config` block supports:

* `horizontal_pod_autoscaling` - (Optional) The status of the Horizontal Pod Autoscaling
    addon, which increases or decreases the number of replica pods a replication controller
    has based on the resource usage of the existing pods. It is enabled by default;
    set `disabled = true` to disable.

* `http_load_balancing` - (Optional) The status of the HTTP (L7) load balancing
    controller addon, which makes it easy to set up HTTP load balancers for services in a
    cluster. It is enabled by default; set `disabled = true` to disable.

* `kubernetes_dashboard` - (Optional) The status of the Kubernetes Dashboard
    add-on, which controls whether the Kubernetes Dashboard is enabled for this cluster.
    It is enabled by default; set `disabled = true` to disable.

This example `addons_config` disables two addons:

```
addons_config {
  http_load_balancing {
    disabled = true
  }
  horizontal_pod_autoscaling {
    disabled = true
  }
}
```

The `maintenance_policy` block supports:

* `daily_maintenance_window` - (Required) Time window specified for daily maintenance operations.
    Specify `start_time` in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format "HH:MM”,
    where HH : \[00-23\] and MM : \[00-59\] GMT. For example:

```
maintenance_policy {
  daily_maintenance_window {
    start_time = "03:00"
  }
}
```

The `ip_allocation_policy` block supports:

* `cluster_secondary_range_name` - (Optional) The name of the secondary range to be
    used as for the cluster CIDR block. The secondary range will be used for pod IP
    addresses. This must be an existing secondary range associated with the cluster
    subnetwork.

* `services_secondary_range_name` - (Optional) The name of the secondary range to be
    used as for the services CIDR block.  The secondary range will be used for service
    ClusterIPs. This must be an existing secondary range associated with the cluster
    subnetwork.

The `master_auth` block supports:

* `password` - (Required) The password to use for HTTP basic authentication when accessing
    the Kubernetes master endpoint

* `username` - (Required) The username to use for HTTP basic authentication when accessing
    the Kubernetes master endpoint

If this block is provided and both `username` and `password` are empty, basic authentication will be disabled.

The `master_authorized_networks_config` block supports:

* `cidr_blocks` - (Optional) Defines up to 10 external networks that can access
    Kubernetes master through HTTPS.

The `master_authorized_networks_config.cidr_blocks` block supports:

* `cidr_block` - (Optional) External network that can access Kubernetes master through HTTPS.
    Must be specified in CIDR notation.

* `display_name` - (Optional) Field for users to identify CIDR blocks.

The `network_policy` block supports:

* `provider` - (Optional) The selected network policy provider. Defaults to PROVIDER_UNSPECIFIED.

* `enabled` - (Optional) Whether network policy is enabled on the cluster. Defaults to false.

The `node_config` block supports:

* `disk_size_gb` - (Optional) Size of the disk attached to each node, specified
    in GB. The smallest allowed disk size is 10GB. Defaults to 100GB.

* `image_type` - (Optional) The image type to use for this node.

* `labels` - (Optional) The Kubernetes labels (key/value pairs) to be applied to each node.

* `local_ssd_count` - (Optional) The amount of local SSD disks that will be
    attached to each cluster node. Defaults to 0.

* `machine_type` - (Optional) The name of a Google Compute Engine machine type.
    Defaults to `n1-standard-1`.

* `metadata` - (Optional) The metadata key/value pairs assigned to instances in
    the cluster.

* `min_cpu_platform` - (Optional) Minimum CPU platform to be used by this instance.
    The instance may be scheduled on the specified or newer CPU platform. Applicable
    values are the friendly names of CPU platforms, such as `Intel Haswell`. See the
    [official documentation](https://cloud.google.com/compute/docs/instances/specify-min-cpu-platform)
    for more information.

* `oauth_scopes` - (Optional) The set of Google API scopes to be made available
    on all of the node VMs under the "default" service account. These can be
    either FQDNs, or scope aliases. The following scopes are necessary to ensure
    the correct functioning of the cluster:

  * `compute-rw` (`https://www.googleapis.com/auth/compute`)
  * `storage-ro` (`https://www.googleapis.com/auth/devstorage.read_only`)
  * `logging-write` (`https://www.googleapis.com/auth/logging.write`),
    if `logging_service` points to Google
  * `monitoring` (`https://www.googleapis.com/auth/monitoring`),
    if `monitoring_service` points to Google

* `preemptible` - (Optional) A boolean that represents whether or not the underlying node VMs
    are preemptible. See the [official documentation](https://cloud.google.com/container-engine/docs/preemptible-vm)
    for more information. Defaults to false.

* `service_account` - (Optional) The service account to be used by the Node VMs.
    If not specified, the "default" service account is used.

* `tags` - (Optional) The list of instance tags applied to all nodes. Tags are used to identify
    valid sources or targets for network firewalls.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `endpoint` - The IP address of this cluster's Kubernetes master.

* `instance_group_urls` - List of instance group URLs which have been assigned
    to the cluster.

* `maintenance_policy.0.daily_maintenance_window.0.duration` - Duration of the time window, automatically chosen to be
    smallest possible in the given scenario.
    Duration will be in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format "PTnHnMnS".

* `master_auth.0.client_certificate` - Base64 encoded public certificate
    used by clients to authenticate to the cluster endpoint.

* `master_auth.0.client_key` - Base64 encoded private key used by clients
    to authenticate to the cluster endpoint.

* `master_auth.0.cluster_ca_certificate` - Base64 encoded public certificate
    that is the root of trust for the cluster.

* `master_version` - The current version of the master in the cluster. This may
    be different than the `min_master_version` set in the config if the master
    has been updated by GKE.

<a id="timeouts"></a>
## Timeouts

`google_container_cluster` provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - (Default `30 minutes`) Used for clusters
- `update` - (Default `10 minutes`) Used for updates to clusters
- `delete` - (Default `10 minutes`) Used for destroying clusters.

## Import

GKE clusters can be imported using the `zone`, and `name`, e.g.

```
$ terraform import google_container_cluster.mycluster us-east1-a/my-cluster
```
