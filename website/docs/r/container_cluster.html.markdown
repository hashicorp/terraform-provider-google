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

- - -

* `zone` - (Optional) The zone that the master and the number of nodes specified
    in `initial_node_count` should be created in. Only one of `zone` and `region`
    may be set. If neither zone nor region are set, the provider zone is used.

* `region` (Optional)
    The region to create the cluster in, for
    [Regional Clusters](https://cloud.google.com/kubernetes-engine/docs/concepts/multi-zone-and-regional-clusters#regional).
    In a Regional Cluster, the number of nodes specified in `initial_node_count` is
    created in three zones of the region (this can be changed by setting `additional_zones`).

* `additional_zones` - (Optional) The list of additional Google Compute Engine
    locations in which the cluster's nodes should be located. If additional zones are
    configured, the number of nodes specified in `initial_node_count` is created in
    all specified zones.

* `addons_config` - (Optional) The configuration for addons supported by GKE.
    Structure is documented below.

* `cluster_ipv4_cidr` - (Optional) The IP address range of the kubernetes pods in
    this cluster. Default is an automatically assigned CIDR.

* `cluster_autoscaling` - (Optional, [Beta](https://terraform.io/docs/providers/google/provider_versions.html))
    Configuration for cluster autoscaling (also called autoprovisioning), as described in
    [the docs](https://cloud.google.com/kubernetes-engine/docs/how-to/node-auto-provisioning).
    Structure is documented below.

* `description` - (Optional) Description of the cluster.

* `enable_binary_authorization` - (Optional, [Beta](https://terraform.io/docs/providers/google/provider_versions.html)) Enable Binary Authorization for this cluster.
    If enabled, all container images will be validated by Google Binary Authorization.

* `enable_kubernetes_alpha` - (Optional) Whether to enable Kubernetes Alpha features for
    this cluster. Note that when this option is enabled, the cluster cannot be upgraded
    and will be automatically deleted after 30 days.

* `enable_tpu` - (Optional, [Beta](https://terraform.io/docs/providers/google/provider_versions.html)) Whether to enable Cloud TPU resources in this cluster.
    See the [official documentation](https://cloud.google.com/tpu/docs/kubernetes-engine-setup).

* `enable_legacy_abac` - (Optional) Whether the ABAC authorizer is enabled for this cluster.
    When enabled, identities in the system, including service accounts, nodes, and controllers,
    will have statically granted permissions beyond those provided by the RBAC configuration or IAM.
    Defaults to `false`

* `initial_node_count` - (Optional) The number of nodes to create in this
    cluster (not including the Kubernetes master). Must be set if `node_pool` is not set.

* `ip_allocation_policy` - (Optional) Configuration for cluster IP allocation. As of now, only pre-allocated subnetworks (custom type with secondary ranges) are supported.
    This will activate IP aliases. See the [official documentation](https://cloud.google.com/kubernetes-engine/docs/how-to/ip-aliases)
    Structure is documented below.

* `logging_service` - (Optional) The logging service that the cluster should
    write logs to. Available options include `logging.googleapis.com`,
    `logging.googleapis.com/kubernetes` (beta), and `none`. Defaults to `logging.googleapis.com`

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
    official release (which is not necessarily the latest version).  Most users will find
    the `google_container_engine_versions` data source useful - it indicates which versions
    are available.  If you intend to specify versions manually, [the docs](https://cloud.google.com/kubernetes-engine/versioning-and-upgrades#specifying_cluster_version)
    describe the various acceptable formats for this field.

-> If you are using the `google_container_engine_versions` datasource with a regional cluster, ensure that you have provided a `region`
to the datasource. A `region` can have a different set of supported versions than its corresponding `zone`s, and not all `zone`s in a
`region` are guaranteed to support the same version.

* `monitoring_service` - (Optional) The monitoring service that the cluster
    should write metrics to.
    Automatically send metrics from pods in the cluster to the Google Cloud Monitoring API.
    VM metrics will be collected by Google Compute Engine regardless of this setting
    Available options include
    `monitoring.googleapis.com`, `monitoring.googleapis.com/kubernetes` (beta) and `none`.
    Defaults to `monitoring.googleapis.com`

* `network` - (Optional) The name or self_link of the Google Compute Engine
    network to which the cluster is connected. For Shared VPC, set this to the self link of the
    shared network.

* `network_policy` - (Optional) Configuration options for the
    [NetworkPolicy](https://kubernetes.io/docs/concepts/services-networking/networkpolicies/)
    feature. Structure is documented below.

* `node_config` -  (Optional) Parameters used in creating the cluster's nodes.
    Structure is documented below.

* `node_pool` - (Optional) List of node pools associated with this cluster.
    See [google_container_node_pool](container_node_pool.html) for schema.
    **Warning:** node pools defined inside a cluster can't be changed (or added/removed) after
    cluster creation without deleting and recreating the entire cluster. Unless you absolutely need the ability
    to say "these are the _only_ node pools associated with this cluster", use the
    [google_container_node_pool](container_node_pool.html) resource instead of this property.

* `node_version` - (Optional) The Kubernetes version on the nodes. Must either be unset
    or set to the same value as `min_master_version` on create. Defaults to the default
    version set by GKE which is not necessarily the latest version.

* `pod_security_policy_config` - (Optional, [Beta](https://terraform.io/docs/providers/google/provider_versions.html)) Configuration for the
    [PodSecurityPolicy](https://cloud.google.com/kubernetes-engine/docs/how-to/pod-security-policies) feature.
    Structure is documented below.

* `private_cluster_config` - (Optional) A set of options for creating
    a private cluster. Structure is documented below.

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

* `remove_default_node_pool` - (Optional) If true, deletes the default node pool upon cluster creation.

* `resource_labels` - (Optional) The GCE resource labels (a map of key/value pairs) to be applied to the cluster.

* `subnetwork` - (Optional) The name or self_link of the Google Compute Engine subnetwork in
    which the cluster's instances are launched.

The `addons_config` block supports:

* `horizontal_pod_autoscaling` - (Optional) The status of the Horizontal Pod Autoscaling
    addon, which increases or decreases the number of replica pods a replication controller
    has based on the resource usage of the existing pods.
    It ensures that a Heapster pod is running in the cluster, which is also used by the Cloud Monitoring service.
    It is enabled by default;
    set `disabled = true` to disable.
* `http_load_balancing` - (Optional) The status of the HTTP (L7) load balancing
    controller addon, which makes it easy to set up HTTP load balancers for services in a
    cluster. It is enabled by default; set `disabled = true` to disable.
* `kubernetes_dashboard` - (Optional) The status of the Kubernetes Dashboard
    add-on, which controls whether the Kubernetes Dashboard is enabled for this cluster.
    It is enabled by default; set `disabled = true` to disable.
* `network_policy_config` - (Optional) Whether we should enable the network policy addon
    for the master.  This must be enabled in order to enable network policy for the nodes.
    It can only be disabled if the nodes already do not have network policies enabled.
    Set `disabled = true` to disable.

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

The `cluster_autoscaling` block supports:
* `enabled` - (Required) Whether cluster autoscaling (also called autoprovisioning) is
    enabled.  To set this to true, make sure your config meets the rest of the
    requirements.  Notably, you'll need `min_master_version` of at least `1.11.2`.
* `resource_limits` - (Optional) A list of limits on the autoprovisioning.
    See [the docs](https://cloud.google.com/kubernetes-engine/docs/how-to/node-auto-provisioning)
    for an explanation of what options are available.  If enabling autoprovisioning, make
    sure to set at least `cpu` and `memory`.  Structure is documented below.

The `resource_limits` block supports:
* `resource_type` - (Required) See [the docs](https://cloud.google.com/kubernetes-engine/docs/how-to/node-auto-provisioning)
    for a list of permitted types - `cpu`, `memory`, and others.
* `minimum` - (Optional) The minimum value for the resource type specified.
* `maximum` - (Optional) The maximum value for the resource type specified.

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

* `cluster_ipv4_cidr_block` - (Optional) The IP address range for the cluster pod IPs.
    Set to blank to have a range chosen with the default size. Set to /netmask (e.g. /14)
    to have a range chosen with a specific netmask. Set to a CIDR notation (e.g. 10.96.0.0/14)
    from the RFC-1918 private networks (e.g. 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16) to
    pick a specific range to use.

* `services_ipv4_cidr_block` - (Optional) The IP address range of the services IPs in this cluster.
    Set to blank to have a range chosen with the default size. Set to /netmask (e.g. /14)
    to have a range chosen with a specific netmask. Set to a CIDR notation (e.g. 10.96.0.0/14)
    from the RFC-1918 private networks (e.g. 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16) to
    pick a specific range to use.

* `create_subnetwork`- (Optional) Whether a new subnetwork will be created automatically for the cluster.

* `subnetwork_name` - (Optional) A custom subnetwork name to be used if create_subnetwork is true.
    If this field is empty, then an automatic name will be chosen for the new subnetwork.

The `master_auth` block supports:

* `password` - (Optional) The password to use for HTTP basic authentication when accessing
    the Kubernetes master endpoint.

* `username` - (Optional) The username to use for HTTP basic authentication when accessing
    the Kubernetes master endpoint. If not present basic auth will be disabled.

* `client_certificate_config` - (Optional) Whether client certificate authorization is enabled for this cluster.  For example:

```
master_auth {
  client_certificate_config {
    issue_client_certificate = false
  }
}
```

If this block is provided and both `username` and `password` are empty, basic authentication will be disabled.
This block also contains several computed attributes, documented below. If this block is not provided, GKE will generate a password for you with the username `admin`.

The `master_authorized_networks_config` block supports:

* `cidr_blocks` - (Optional) Defines up to 20 external networks that can access
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

* `disk_type` - (Optional) Type of the disk attached to each node
    (e.g. 'pd-standard' or 'pd-ssd'). If unspecified, the default disk type is 'pd-standard'

* `guest_accelerator` - (Optional) List of the type and count of accelerator cards attached to the instance.
    Structure documented below.

* `image_type` - (Optional) The image type to use for this node. Note that changing the image type
    will delete and recreate all nodes in the node pool.

* `labels` - (Optional) The Kubernetes labels (key/value pairs) to be applied to each node.

* `local_ssd_count` - (Optional) The amount of local SSD disks that will be
    attached to each cluster node. Defaults to 0.

* `machine_type` - (Optional) The name of a Google Compute Engine machine type.
    Defaults to `n1-standard-1`. To create a custom machine type, value should be set as specified
    [here](https://cloud.google.com/compute/docs/reference/latest/instances#machineType).

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
    In order to use the configured `oauth_scopes` for logging and monitoring, the service account being used needs the
    [roles/logging.logWriter](https://cloud.google.com/iam/docs/understanding-roles#stackdriver_logging_roles) and
    [roles/monitoring.metricWriter](https://cloud.google.com/iam/docs/understanding-roles#stackdriver_monitoring_roles) roles.

     -> Projects that enable the [Cloud Compute Engine API](https://cloud.google.com/compute/) with Terraform may need these roles added manually to the service account. Projects that enable the API in the Cloud Console should have them added automatically.

* `tags` - (Optional) The list of instance tags applied to all nodes. Tags are used to identify
    valid sources or targets for network firewalls.

* `taint` - (Optional, [Beta](https://terraform.io/docs/providers/google/provider_versions.html)) List of
    [kubernetes taints](https://kubernetes.io/docs/concepts/configuration/taint-and-toleration/)
    to apply to each node. Structure is documented below.

* `workload_metadata_config` - (Optional, [Beta](https://terraform.io/docs/providers/google/provider_versions.html)) Metadata configuration to expose to workloads on the node pool.
    Structure is documented below.

The `guest_accelerator` block supports:

* `type` (Required) - The accelerator type resource to expose to this instance. E.g. `nvidia-tesla-k80`.

* `count` (Required) - The number of the guest accelerator cards exposed to this instance.

The `pod_security_policy_config` block supports:

* `enabled` (Required) - Enable the PodSecurityPolicy controller for this cluster.
    If enabled, pods must be valid under a PodSecurityPolicy to be created.

The `private_cluster_config` block supports:

* `enable_private_endpoint` (Optional) - Whether the master's internal IP address is used as the cluster endpoint.

* `enable_private_nodes` (Optional) - Whether nodes have internal IP addresses only. If enabled, all nodes are given only RFC 1918 private
    addresses and communicate with the master via private networking.

* `master_ipv4_cidr_block` (Optional) - The IP range in CIDR notation to use for the hosted master network. This range will be used for
    assigning internal IP addresses to the master or set of masters, as well as the ILB VIP. This range must not overlap with any other ranges
    in use within the cluster's network.

In addition, the `private_cluster_config` allows access to the following read-only fields:

* `private_endpoint` - The internal IP address of this cluster's master endpoint.

* `public_endpoint` - The external IP address of this cluster's master endpoint.

The `taint` block supports:

* `key` (Required) Key for taint.

* `value` (Required) Value for taint.

* `effect` (Required) Effect for taint. Accepted values are `NO_SCHEDULE`, `PREFER_NO_SCHEDULE`, and `NO_EXECUTE`.

The `workload_metadata_config` block supports:

* `node_metadata` (Required) How to expose the node metadata to the workload running on the node.
    Accepted values are:
    * UNSPECIFIED: Not Set
    * SECURE: Prevent workloads not in hostNetwork from accessing certain VM metadata, specifically kube-env, which contains Kubelet credentials, and the instance identity token. See [Metadata Concealment](https://cloud.google.com/kubernetes-engine/docs/how-to/metadata-proxy) documentation.
    * EXPOSE: Expose all VM metadata to pods.

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

* `tpu_ipv4_cidr_block` - ([Beta](https://terraform.io/docs/providers/google/provider_versions.html)) The IP address range of the Cloud TPUs in this cluster, in
    [CIDR](http://en.wikipedia.org/wiki/Classless_Inter-Domain_Routing)
    notation (e.g. `1.2.3.4/29`).

<a id="timeouts"></a>
## Timeouts

`google_container_cluster` provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - (Default `30 minutes`) Used for clusters
- `update` - (Default `10 minutes`) Used for updates to clusters
- `delete` - (Default `10 minutes`) Used for destroying clusters.

## Import

GKE clusters can be imported using the `project` , `zone` or `region`, and `name`. If the project is omitted, the default
provider value will be used. Examples:

```
$ terraform import google_container_cluster.mycluster my-gcp-project/us-east1-a/my-cluster

$ terraform import google_container_cluster.mycluster us-east1-a/my-cluster
```
