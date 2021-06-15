---
subcategory: "Kubernetes (Container) Engine"
layout: "google"
page_title: "Google: google_container_cluster"
sidebar_current: "docs-google-container-cluster"
description: |-
  Creates a Google Kubernetes Engine (GKE) cluster.
---

# google\_container\_cluster

-> Visit the [Provision a GKE Cluster (Google Cloud)](https://learn.hashicorp.com/tutorials/terraform/gke?in=terraform/kubernetes&utm_source=WEBSITE&utm_medium=WEB_IO&utm_offer=ARTICLE_PAGE&utm_content=DOCS) Learn tutorial to learn how to provision and interact
with a GKE cluster.

-> See the [Using GKE with Terraform](/docs/providers/google/guides/using_gke_with_terraform.html)
guide for more information about using GKE with Terraform.

Manages a Google Kubernetes Engine (GKE) cluster. For more information see
[the official documentation](https://cloud.google.com/container-engine/docs/clusters)
and [the API reference](https://cloud.google.com/kubernetes-engine/docs/reference/rest/v1beta1/projects.locations.clusters).

~> **Note:** All arguments and attributes, including basic auth username and
passwords as well as certificate outputs will be stored in the raw state as
plaintext. [Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Example Usage - with a separately managed node pool (recommended)

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
  location   = "us-central1"
  cluster    = google_container_cluster.primary.name
  node_count = 1

  node_config {
    preemptible  = true
    machine_type = "e2-medium"

    # Google recommends custom service accounts that have cloud-platform scope and permissions granted via IAM Roles.
    service_account = google_service_account.default.email
    oauth_scopes    = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]
  }
}
```

## Example Usage - with the default node pool

```hcl
resource "google_service_account" "default" {
  account_id   = "service-account-id"
  display_name = "Service Account"
}

resource "google_container_cluster" "primary" {
  name               = "marcellus-wallace"
  location           = "us-central1-a"
  initial_node_count = 3
  node_config {
    # Google recommends custom service accounts that have cloud-platform scope and permissions granted via IAM Roles.
    service_account = google_service_account.default.email
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]
    labels = {
      foo = "bar"
    }
    tags = ["foo", "bar"]
  }
  timeouts {
    create = "30m"
    update = "40m"
  }
}
```

## Argument Reference

* `name` - (Required) The name of the cluster, unique within the project and
location.

- - -

* `location` - (Optional) The location (region or zone) in which the cluster
master will be created, as well as the default node location. If you specify a
zone (such as `us-central1-a`), the cluster will be a zonal cluster with a
single cluster master. If you specify a region (such as `us-west1`), the
cluster will be a regional cluster with multiple masters spread across zones in
the region, and with default node locations in those zones as well

* `node_locations` - (Optional) The list of zones in which the cluster's nodes
are located. Nodes must be in the region of their regional cluster or in the
same region as their cluster's zone for zonal clusters. If this is specified for
a zonal cluster, omit the cluster's zone.

-> A "multi-zonal" cluster is a zonal cluster with at least one additional zone
defined; in a multi-zonal cluster, the cluster master is only present in a
single zone while nodes are present in each of the primary zone and the node
locations. In contrast, in a regional cluster, cluster master nodes are present
in multiple zones in the region. For that reason, regional clusters should be
preferred.

* `addons_config` - (Optional) The configuration for addons supported by GKE.
    Structure is documented below.

* `cluster_ipv4_cidr` - (Optional) The IP address range of the Kubernetes pods
in this cluster in CIDR notation (e.g. `10.96.0.0/14`). Leave blank to have one
automatically chosen or specify a `/14` block in `10.0.0.0/8`. This field will
only work for routes-based clusters, where `ip_allocation_policy` is not defined.

* `cluster_autoscaling` - (Optional)
Per-cluster configuration of Node Auto-Provisioning with Cluster Autoscaler to
automatically adjust the size of the cluster and create/delete node pools based
on the current needs of the cluster's workload. See the
[guide to using Node Auto-Provisioning](https://cloud.google.com/kubernetes-engine/docs/how-to/node-auto-provisioning)
for more details. Structure is documented below.

* `database_encryption` - (Optional)
    Structure is documented below.

* `description` - (Optional) Description of the cluster.

* `default_max_pods_per_node` - (Optional) The default maximum number of pods
per node in this cluster. This doesn't work on "routes-based" clusters, clusters
that don't have IP Aliasing enabled. See the [official documentation](https://cloud.google.com/kubernetes-engine/docs/how-to/flexible-pod-cidr)
for more information.

* `enable_binary_authorization` - (Optional) Enable Binary Authorization for this cluster.
    If enabled, all container images will be validated by Google Binary Authorization.

* `enable_kubernetes_alpha` - (Optional) Whether to enable Kubernetes Alpha features for
    this cluster. Note that when this option is enabled, the cluster cannot be upgraded
    and will be automatically deleted after 30 days.

* `enable_tpu` - (Optional) Whether to enable Cloud TPU resources in this cluster.
    See the [official documentation](https://cloud.google.com/tpu/docs/kubernetes-engine-setup).

* `enable_legacy_abac` - (Optional) Whether the ABAC authorizer is enabled for this cluster.
    When enabled, identities in the system, including service accounts, nodes, and controllers,
    will have statically granted permissions beyond those provided by the RBAC configuration or IAM.
    Defaults to `false`

* `enable_shielded_nodes` - (Optional) Enable Shielded Nodes features on all nodes in this cluster.  Defaults to `false`.

* `enable_autopilot` - (Optional) Enable Autopilot for this cluster. Defaults to `false`.
    Note that when this option is enabled, certain features of Standard GKE are not available.
    See the [official documentation](https://cloud.google.com/kubernetes-engine/docs/concepts/autopilot-overview#comparison)
    for available features.

* `initial_node_count` - (Optional) The number of nodes to create in this
cluster's default node pool. In regional or multi-zonal clusters, this is the
number of nodes per zone. Must be set if `node_pool` is not set. If you're using
`google_container_node_pool` objects with no default node pool, you'll need to
set this to a value of at least `1`, alongside setting
`remove_default_node_pool` to `true`.

* `ip_allocation_policy` - (Optional) Configuration of cluster IP allocation for
VPC-native clusters. Adding this block enables [IP aliasing](https://cloud.google.com/kubernetes-engine/docs/how-to/ip-aliases),
making the cluster VPC-native instead of routes-based. Structure is documented
below.

* `networking_mode` - (Optional, [Beta]) Determines whether alias IPs or routes will be used for pod IPs in the cluster.
Options are `VPC_NATIVE` or `ROUTES`. `VPC_NATIVE` enables [IP aliasing](https://cloud.google.com/kubernetes-engine/docs/how-to/ip-aliases),
and requires the `ip_allocation_policy` block to be defined. By default when this field is unspecified, GKE will create a `ROUTES`-based cluster.

* `logging_service` - (Optional) The logging service that the cluster should
    write logs to. Available options include `logging.googleapis.com`(Legacy Stackdriver),
    `logging.googleapis.com/kubernetes`(Stackdriver Kubernetes Engine Logging), and `none`. Defaults to `logging.googleapis.com/kubernetes`

* `maintenance_policy` - (Optional) The maintenance policy to use for the cluster. Structure is
    documented below.

* `master_auth` - (Optional) The authentication information for accessing the
Kubernetes master. Some values in this block are only returned by the API if
your service account has permission to get credentials for your GKE cluster. If
you see an unexpected diff removing a username/password or unsetting your client
cert, ensure you have the `container.clusters.getCredentials` permission.
Structure is documented below. This has been deprecated as of GKE 1.19.

* `master_authorized_networks_config` - (Optional) The desired configuration options
    for master authorized networks. Omit the nested `cidr_blocks` attribute to disallow
    external access (except the cluster node IPs, which GKE automatically whitelists).

* `min_master_version` - (Optional) The minimum version of the master. GKE
    will auto-update the master to new versions, so this does not guarantee the
    current master version--use the read-only `master_version` field to obtain that.
    If unset, the cluster's version will be set by GKE to the version of the most recent
    official release (which is not necessarily the latest version).  Most users will find
    the `google_container_engine_versions` data source useful - it indicates which versions
    are available, and can be use to approximate fuzzy versions in a
    Terraform-compatible way. If you intend to specify versions manually,
    [the docs](https://cloud.google.com/kubernetes-engine/versioning-and-upgrades#specifying_cluster_version)
    describe the various acceptable formats for this field.

-> If you are using the `google_container_engine_versions` datasource with a regional cluster, ensure that you have provided a `location`
to the datasource. A region can have a different set of supported versions than its corresponding zones, and not all zones in a
region are guaranteed to support the same version.

* `monitoring_service` - (Optional) The monitoring service that the cluster
    should write metrics to.
    Automatically send metrics from pods in the cluster to the Google Cloud Monitoring API.
    VM metrics will be collected by Google Compute Engine regardless of this setting
    Available options include
    `monitoring.googleapis.com`(Legacy Stackdriver), `monitoring.googleapis.com/kubernetes`(Stackdriver Kubernetes Engine Monitoring), and `none`.
    Defaults to `monitoring.googleapis.com/kubernetes`

* `network` - (Optional) The name or self_link of the Google Compute Engine
    network to which the cluster is connected. For Shared VPC, set this to the self link of the
    shared network.

* `network_policy` - (Optional) Configuration options for the
    [NetworkPolicy](https://kubernetes.io/docs/concepts/services-networking/networkpolicies/)
    feature. Structure is documented below.

* `node_config` -  (Optional) Parameters used in creating the default node pool.
    Generally, this field should not be used at the same time as a
    `google_container_node_pool` or a `node_pool` block; this configuration
    manages the default node pool, which isn't recommended to be used with
    Terraform. Structure is documented below.

* `node_pool` - (Optional) List of node pools associated with this cluster.
    See [google_container_node_pool](container_node_pool.html) for schema.
    **Warning:** node pools defined inside a cluster can't be changed (or added/removed) after
    cluster creation without deleting and recreating the entire cluster. Unless you absolutely need the ability
    to say "these are the _only_ node pools associated with this cluster", use the
    [google_container_node_pool](container_node_pool.html) resource instead of this property.

* `node_version` - (Optional) The Kubernetes version on the nodes. Must either be unset
    or set to the same value as `min_master_version` on create. Defaults to the default
    version set by GKE which is not necessarily the latest version. This only affects
    nodes in the default node pool. While a fuzzy version can be specified, it's
    recommended that you specify explicit versions as Terraform will see spurious diffs
    when fuzzy versions are used. See the `google_container_engine_versions` data source's
    `version_prefix` field to approximate fuzzy versions in a Terraform-compatible way.
    To update nodes in other node pools, use the `version` attribute on the node pool.

* `notification_config` - (Optional, [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html)) Configuration for the [cluster upgrade notifications](https://cloud.google.com/kubernetes-engine/docs/how-to/cluster-upgrade-notifications) feature. Structure is documented below.

* `pod_security_policy_config` - (Optional, [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html)) Configuration for the
    [PodSecurityPolicy](https://cloud.google.com/kubernetes-engine/docs/how-to/pod-security-policies) feature.
    Structure is documented below.

* `authenticator_groups_config` - (Optional) Configuration for the
    [Google Groups for GKE](https://cloud.google.com/kubernetes-engine/docs/how-to/role-based-access-control#groups-setup-gsuite) feature.
    Structure is documented below.

* `private_cluster_config` - (Optional) Configuration for [private clusters](https://cloud.google.com/kubernetes-engine/docs/how-to/private-clusters),
clusters with private nodes. Structure is documented below.

* `cluster_telemetry` - (Optional, [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html)) Configuration for
   [ClusterTelemetry](https://cloud.google.com/monitoring/kubernetes-engine/installing#controlling_the_collection_of_application_logs) feature,
   Structure is documented below.

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

* `release_channel` - (Optional)
Configuration options for the [Release channel](https://cloud.google.com/kubernetes-engine/docs/concepts/release-channels)
feature, which provide more control over automatic upgrades of your GKE clusters.
When updating this field, GKE imposes specific version requirements. See
[Migrating between release channels](https://cloud.google.com/kubernetes-engine/docs/concepts/release-channels#migrating_between_release_channels)
for more details; the `google_container_engine_versions` datasource can provide
the default version for a channel. Note that removing the `release_channel`
field from your config will cause Terraform to stop managing your cluster's
release channel, but will not unenroll it. Instead, use the `"UNSPECIFIED"`
channel. Structure is documented below.

* `remove_default_node_pool` - (Optional) If `true`, deletes the default node
    pool upon cluster creation. If you're using `google_container_node_pool`
    resources with no default node pool, this should be set to `true`, alongside
    setting `initial_node_count` to at least `1`.

* `resource_labels` - (Optional) The GCE resource labels (a map of key/value pairs) to be applied to the cluster.

* `resource_usage_export_config` - (Optional) Configuration for the
    [ResourceUsageExportConfig](https://cloud.google.com/kubernetes-engine/docs/how-to/cluster-usage-metering) feature.
    Structure is documented below.

* `subnetwork` - (Optional) The name or self_link of the Google Compute Engine
subnetwork in which the cluster's instances are launched.

* `vertical_pod_autoscaling` - (Optional, [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html))
    Vertical Pod Autoscaling automatically adjusts the resources of pods controlled by it.
    Structure is documented below.

* `workload_identity_config` - (Optional)
    Workload Identity allows Kubernetes service accounts to act as a user-managed
    [Google IAM Service Account](https://cloud.google.com/iam/docs/service-accounts#user-managed_service_accounts).
    Structure is documented below.

* `enable_intranode_visibility` - (Optional)
    Whether Intra-node visibility is enabled for this cluster. This makes same node pod to pod traffic visible for VPC network.

* `enable_l4_ilb_subsetting` - (Optional, [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html))
    Whether L4ILB Subsetting is enabled for this cluster.

* `private_ipv6_google_access` - (Optional)
    The desired state of IPv6 connectivity to Google Services. By default, no private IPv6 access to or from Google Services (all access will be via IPv4).

* `datapath_provider` - (Optional)
    The desired datapath provider for this cluster. By default, uses the IPTables-based kube-proxy implementation.

* `default_snat_status` - (Optional)
  [GKE SNAT](https://cloud.google.com/kubernetes-engine/docs/how-to/ip-masquerade-agent#how_ipmasq_works) DefaultSnatStatus contains the desired state of whether default sNAT should be disabled on the cluster, [API doc](https://cloud.google.com/kubernetes-engine/docs/reference/rest/v1beta1/projects.locations.clusters#networkconfig).

The `default_snat_status` block supports

*  `disabled` - (Required) Whether the cluster disables default in-node sNAT rules. In-node sNAT rules will be disabled when defaultSnatStatus is disabled.When disabled is set to false, default IP masquerade rules will be applied to the nodes to prevent sNAT on cluster internal traffic

The `cluster_telemetry` block supports
* `type` - Telemetry integration for the cluster. Supported values (`ENABLED, DISABLED, SYSTEM_ONLY`);
   `SYSTEM_ONLY` (Only system components are monitored and logged) is only available in GKE versions 1.15 and later.

The `addons_config` block supports:

* `horizontal_pod_autoscaling` - (Optional) The status of the Horizontal Pod Autoscaling
    addon, which increases or decreases the number of replica pods a replication controller
    has based on the resource usage of the existing pods.
    It is enabled by default;
    set `disabled = true` to disable.

* `http_load_balancing` - (Optional) The status of the HTTP (L7) load balancing
    controller addon, which makes it easy to set up HTTP load balancers for services in a
    cluster. It is enabled by default; set `disabled = true` to disable.

* `network_policy_config` - (Optional) Whether we should enable the network policy addon
    for the master.  This must be enabled in order to enable network policy for the nodes.
    To enable this, you must also define a [`network_policy`](#network_policy) block,
    otherwise nothing will happen.
    It can only be disabled if the nodes already do not have network policies enabled.
    Defaults to disabled; set `disabled = false` to enable.

* `cloudrun_config` - (Optional). Structure is documented below.

* `istio_config` - (Optional, [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html)).
    Structure is documented below.

* `dns_cache_config` - (Optional, [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html)).
    The status of the NodeLocal DNSCache addon. It is disabled by default.
    Set `enabled = true` to enable.

    **Enabling/Disabling NodeLocal DNSCache in an existing cluster is a disruptive operation.
    All cluster nodes running GKE 1.15 and higher are recreated.**

* `gce_persistent_disk_csi_driver_config` - (Optional, [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html)).
    Whether this cluster should enable the Google Compute Engine Persistent Disk Container Storage Interface (CSI) Driver. Defaults to disabled; set `enabled = true` to enable.

* `kalm_config` - (Optional, [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html)).
    Configuration for the KALM addon, which manages the lifecycle of k8s. It is disabled by default; Set `enabled = true` to enable.

*  `config_connector_config` -  (Optional, [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html)).
    The status of the ConfigConnector addon. It is disabled by default; Set `enabled = true` to enable.

This example `addons_config` disables two addons:

```hcl
addons_config {
  http_load_balancing {
    disabled = true
  }

  horizontal_pod_autoscaling {
    disabled = true
  }
}
```

The `database_encryption` block supports:

* `state` - (Required) `ENCRYPTED` or `DECRYPTED`

* `key_name` - (Required) the key to use to encrypt/decrypt secrets.  See the [DatabaseEncryption definition](https://cloud.google.com/kubernetes-engine/docs/reference/rest/v1beta1/projects.locations.clusters#Cluster.DatabaseEncryption) for more information.

The `cloudrun_config` block supports:

* `disabled` - (Optional) The status of the CloudRun addon. It is disabled by default. Set `disabled=false` to enable.

* `load_balancer_type` - (Optional) The load balancer type of CloudRun ingress service. It is external load balancer by default.
    Set `load_balancer_type=LOAD_BALANCER_TYPE_INTERNAL` to configure it as internal load balancer.

The `istio_config` block supports:

* `disabled` - (Optional) The status of the Istio addon, which makes it easy to set up Istio for services in a
    cluster. It is disabled by default. Set `disabled = false` to enable.

* `auth` - (Optional) The authentication type between services in Istio. Available options include `AUTH_MUTUAL_TLS`.

The `cluster_autoscaling` block supports:

* `enabled` - (Required) Whether node auto-provisioning is enabled. Resource
limits for `cpu` and `memory` must be defined to enable node auto-provisioning.

* `resource_limits` - (Optional) Global constraints for machine resources in the
cluster. Configuring the `cpu` and `memory` types is required if node
auto-provisioning is enabled. These limits will apply to node pool autoscaling
in addition to node auto-provisioning. Structure is documented below.

* `auto_provisioning_defaults` - (Optional) Contains defaults for a node pool created by NAP.
Structure is documented below.

* `autoscaling_profile` - (Optional, [Beta](https://terraform.io/docs/providers/google/provider_versions.html)) Configuration
options for the [Autoscaling profile](https://cloud.google.com/kubernetes-engine/docs/concepts/cluster-autoscaler#autoscaling_profiles)
feature, which lets you choose whether the cluster autoscaler should optimize for resource utilization or resource availability
when deciding to remove nodes from a cluster. Can be `BALANCED` or `OPTIMIZE_UTILIZATION`. Defaults to `BALANCED`.

The `resource_limits` block supports:

* `resource_type` - (Required) The type of the resource. For example, `cpu` and
`memory`.  See the [guide to using Node Auto-Provisioning](https://cloud.google.com/kubernetes-engine/docs/how-to/node-auto-provisioning)
for a list of types.

* `minimum` - (Optional) Minimum amount of the resource in the cluster.

* `maximum` - (Optional) Maximum amount of the resource in the cluster.

The `auto_provisioning_defaults` block supports:

* `min_cpu_platform` - (Optional, [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html))
Minimum CPU platform to be used for NAP created node pools. The instance may be scheduled on the
specified or newer CPU platform. Applicable values are the friendly names of CPU platforms, such
as "Intel Haswell" or "Intel Sandy Bridge".

* `oauth_scopes` - (Optional) Scopes that are used by NAP when creating node pools. Use the "https://www.googleapis.com/auth/cloud-platform" scope to grant access to all APIs. It is recommended that you set `service_account` to a non-default service account and grant IAM roles to that service account for only the resources that it needs.

-> `monitoring.write` is always enabled regardless of user input.  `monitoring` and `logging.write` may also be enabled depending on the values for `monitoring_service` and `logging_service`.

* `service_account` - (Optional) The Google Cloud Platform Service Account to be used by the node VMs.

The `authenticator_groups_config` block supports:

* `security_group` - (Required) The name of the RBAC security group for use with Google security groups in Kubernetes RBAC. Group name must be in format `gke-security-groups@yourdomain.com`.

The `maintenance_policy` block supports:
* `daily_maintenance_window` - (Optional) structure documented below.
* `recurring_window` - (Optional) structure documented below
* `maintenance_exclusion` - (Optional) structure documented below

In beta, one or the other of `recurring_window` and `daily_maintenance_window` is required if a `maintenance_policy` block is supplied.

* `daily_maintenance_window` - Time window specified for daily maintenance operations.
    Specify `start_time` in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format "HH:MM”,
    where HH : \[00-23\] and MM : \[00-59\] GMT. For example:

Examples:
```hcl
maintenance_policy {
  daily_maintenance_window {
    start_time = "03:00"
  }
}
```

* `recurring_window` - Time window for recurring maintenance operations.

Specify `start_time` and `end_time` in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) "Zulu" date format.  The start time's date is
the initial date that the window starts, and the end time is used for calculating duration.  Specify `recurrence` in
[RFC5545](https://tools.ietf.org/html/rfc5545#section-3.8.5.3) RRULE format, to specify when this recurs.
Note that GKE may accept other formats, but will return values in UTC, causing a permanent diff.

Examples:
```
maintenance_policy {
  recurring_window {
    start_time = "2019-08-01T02:00:00Z"
    end_time = "2019-08-01T06:00:00Z"
    recurrence = "FREQ=DAILY"
  }
}
```

```
maintenance_policy {
  recurring_window {
    start_time = "2019-01-01T09:00:00Z"
    end_time = "2019-01-01T17:00:00Z"
    recurrence = "FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR"
  }
}
```

* `maintenance_exclusion` - Exceptions to maintenance window. Non-emergency maintenance should not occur in these windows. A cluster can have up to three maintenance exclusions at a time [Maintenance Window and Exclusions](https://cloud.google.com/kubernetes-engine/docs/concepts/maintenance-windows-and-exclusions)

Specify `start_time` and `end_time` in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) "Zulu" date format.  The start time's date is
the initial date that the window starts, and the end time is used for calculating duration.Specify `recurrence` in
[RFC5545](https://tools.ietf.org/html/rfc5545#section-3.8.5.3) RRULE format, to specify when this recurs.
Note that GKE may accept other formats, but will return values in UTC, causing a permanent diff.

Examples:

```
maintenance_policy {
  recurring_window {
    start_time = "2019-01-01T00:00:00Z"
    end_time = "2019-01-02T00:00:00Z"
    recurrence = "FREQ=DAILY"
  }
  maintenance_exclusion{
    exclusion_name = "batch job"
    start_time = "2019-01-01T00:00:00Z"
    end_time = "2019-01-02T00:00:00Z"
  }
  maintenance_exclusion{
    exclusion_name = "holiday data load"
    start_time = "2019-05-01T00:00:00Z"
    end_time = "2019-05-02T00:00:00Z"
  }
}
```

The `ip_allocation_policy` block supports:

* `cluster_secondary_range_name` - (Optional) The name of the existing secondary
range in the cluster's subnetwork to use for pod IP addresses. Alternatively,
`cluster_ipv4_cidr_block` can be used to automatically create a GKE-managed one.

* `services_secondary_range_name` - (Optional) The name of the existing
secondary range in the cluster's subnetwork to use for service `ClusterIP`s.
Alternatively, `services_ipv4_cidr_block` can be used to automatically create a
GKE-managed one.

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

The `master_auth` block supports:

* `password` - (Optional) The password to use for HTTP basic authentication when accessing
    the Kubernetes master endpoint. This has been deprecated as of GKE 1.19.

* `username` - (Optional) The username to use for HTTP basic authentication when accessing
    the Kubernetes master endpoint. If not present basic auth will be disabled. This has been deprecated as of GKE 1.19.

* `client_certificate_config` - (Optional) Whether client certificate authorization is enabled for this cluster.  For example:

```hcl
master_auth {
  client_certificate_config {
    issue_client_certificate = false
  }
}
```

If this block is provided and both `username` and `password` are empty, basic authentication will be disabled.
This block also contains several computed attributes, documented below. If this block is not provided, GKE will generate a password for you with the username `admin`.

The `master_authorized_networks_config` block supports:

* `cidr_blocks` - (Optional) External networks that can access the
    Kubernetes cluster master through HTTPS.

The `master_authorized_networks_config.cidr_blocks` block supports:

* `cidr_block` - (Optional) External network that can access Kubernetes master through HTTPS.
    Must be specified in CIDR notation.

* `display_name` - (Optional) Field for users to identify CIDR blocks.

The `network_policy` block supports:

* `provider` - (Optional) The selected network policy provider. Defaults to PROVIDER_UNSPECIFIED.

* `enabled` - (Required) Whether network policy is enabled on the cluster.

The `node_config` block supports:

* `disk_size_gb` - (Optional) Size of the disk attached to each node, specified
    in GB. The smallest allowed disk size is 10GB. Defaults to 100GB.

* `disk_type` - (Optional) Type of the disk attached to each node
    (e.g. 'pd-standard', 'pd-balanced' or 'pd-ssd'). If unspecified, the default disk type is 'pd-standard'

* `ephemeral_storage_config` - (Optional, [Beta]) Parameters for the ephemeral storage filesystem. If unspecified, ephemeral storage is backed by the boot disk. Structure is documented below.

```hcl
ephemeral_storage_config {
  local_ssd_count = 2
}
```

* `guest_accelerator` - (Optional) List of the type and count of accelerator cards attached to the instance.
    Structure documented below.
    To support removal of guest_accelerators in Terraform 0.12 this field is an
    [Attribute as Block](/docs/configuration/attr-as-blocks.html)

* `image_type` - (Optional) The image type to use for this node. Note that changing the image type
    will delete and recreate all nodes in the node pool.

* `labels` - (Optional) The Kubernetes labels (key/value pairs) to be applied to each node. The kubernetes.io/ and k8s.io/ prefixes are
    reserved by Kubernetes Core components and cannot be specified.

* `local_ssd_count` - (Optional) The amount of local SSD disks that will be
    attached to each cluster node. Defaults to 0.

* `machine_type` - (Optional) The name of a Google Compute Engine machine type.
    Defaults to `e2-medium`. To create a custom machine type, value should be set as specified
    [here](https://cloud.google.com/compute/docs/reference/latest/instances#machineType).

* `metadata` - (Optional) The metadata key/value pairs assigned to instances in
    the cluster. From GKE `1.12` onwards, `disable-legacy-endpoints` is set to
    `true` by the API; if `metadata` is set but that default value is not
    included, Terraform will attempt to unset the value. To avoid this, set the
    value in your config.

* `min_cpu_platform` - (Optional) Minimum CPU platform to be used by this instance.
    The instance may be scheduled on the specified or newer CPU platform. Applicable
    values are the friendly names of CPU platforms, such as `Intel Haswell`. See the
    [official documentation](https://cloud.google.com/compute/docs/instances/specify-min-cpu-platform)
    for more information.

* `oauth_scopes` - (Optional) The set of Google API scopes to be made available
    on all of the node VMs under the "default" service account.
    Use the "https://www.googleapis.com/auth/cloud-platform" scope to grant access to all APIs. It is recommended that you set `service_account` to a non-default service account and grant IAM roles to that service account for only the resources that it needs.

    See the [official documentation](https://cloud.google.com/kubernetes-engine/docs/how-to/access-scopes) for information on migrating off of legacy access scopes.

* `preemptible` - (Optional) A boolean that represents whether or not the underlying node VMs
    are preemptible. See the [official documentation](https://cloud.google.com/container-engine/docs/preemptible-vm)
    for more information. Defaults to false.

* `sandbox_config` - (Optional, [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html)) [GKE Sandbox](https://cloud.google.com/kubernetes-engine/docs/how-to/sandbox-pods) configuration. When enabling this feature you must specify `image_type = "COS_CONTAINERD"` and `node_version = "1.12.7-gke.17"` or later to use it.
    Structure is documented below.

* `boot_disk_kms_key` - (Optional, [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html)) The Customer Managed Encryption Key used to encrypt the boot disk attached to each node in the node pool. This should be of the form projects/[KEY_PROJECT_ID]/locations/[LOCATION]/keyRings/[RING_NAME]/cryptoKeys/[KEY_NAME]. For more information about protecting resources with Cloud KMS Keys please see: https://cloud.google.com/compute/docs/disks/customer-managed-encryption

* `service_account` - (Optional) The service account to be used by the Node VMs.
    If not specified, the "default" service account is used.

* `shielded_instance_config` - (Optional) Shielded Instance options. Structure is documented below.

* `tags` - (Optional) The list of instance tags applied to all nodes. Tags are used to identify
    valid sources or targets for network firewalls.

* `taint` - (Optional) A list of [Kubernetes taints](https://kubernetes.io/docs/concepts/configuration/taint-and-toleration/)
to apply to nodes. GKE's API can only set this field on cluster creation.
However, GKE will add taints to your nodes if you enable certain features such
as GPUs. If this field is set, any diffs on this field will cause Terraform to
recreate the underlying resource. Taint values can be updated safely in
Kubernetes (eg. through `kubectl`), and it's recommended that you do not use
this field to manage taints. If you do, `lifecycle.ignore_changes` is
recommended. Structure is documented below.

* `workload_metadata_config` - (Optional) Metadata configuration to expose to workloads on the node pool.
    Structure is documented below.

* `kubelet_config` - (Optional, [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html))
Kubelet configuration, currently supported attributes can be found [here](https://cloud.google.com/sdk/gcloud/reference/beta/container/node-pools/create#--system-config-from-file).
Structure is documented below.

```
kubelet_config {
  cpu_manager_policy   = "static"
  cpu_cfs_quota        = true
  cpu_cfs_quota_period = "100us"
}
```

* `linux_node_config` - (Optional, [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html))
Linux node configuration, currently supported attributes can be found [here](https://cloud.google.com/sdk/gcloud/reference/beta/container/node-pools/create#--system-config-from-file).
Note that validations happen all server side. All attributes are optional.
Structure is documented below.

```hcl
linux_node_config {
  sysctls = {
    "net.core.netdev_max_backlog" = "10000"
    "net.core.rmem_max"           = "10000"
  }
}
```

The `ephemeral_storage_config` block supports:

* `local_ssd_count` (Required) - Number of local SSDs to use to back ephemeral storage. Uses NVMe interfaces. Each local SSD is 375 GB in size. If zero, it means to disable using local SSDs as ephemeral storage.

The `guest_accelerator` block supports:

* `type` (Required) - The accelerator type resource to expose to this instance. E.g. `nvidia-tesla-k80`.

* `count` (Required) - The number of the guest accelerator cards exposed to this instance.

The `workload_identity_config` block supports:

* `identity_namespace` (Required) - Currently, the only supported identity namespace is the project's default.

```hcl
workload_identity_config {
  identity_namespace = "${data.google_project.project.project_id}.svc.id.goog"
}
```

The `notification_config` block supports:

* `pubsub` (Required) - The pubsub config for the cluster's upgrade notifications.

The `pubsub` block supports:

* `enabled` (Required) - Whether or not the notification config is enabled

* `topic` (Optional) - The pubsub topic to push upgrade notifications to. Must be in the same project as the cluster. Must be in the format: `projects/{project}/topics/{topic}`.

```hcl
notification_config {
  pubsub {
    enabled = true
    topic = google_pubsub_topic.notifications.id
  }
}
```

The `pod_security_policy_config` block supports:

* `enabled` (Required) - Enable the PodSecurityPolicy controller for this cluster.
    If enabled, pods must be valid under a PodSecurityPolicy to be created.

The `private_cluster_config` block supports:

* `enable_private_nodes` (Optional) - Enables the private cluster feature,
creating a private endpoint on the cluster. In a private cluster, nodes only
have RFC 1918 private addresses and communicate with the master's private
endpoint via private networking.

* `enable_private_endpoint` (Optional) - When `true`, the cluster's private
endpoint is used as the cluster endpoint and access through the public endpoint
is disabled. When `false`, either endpoint can be used. This field only applies
to private clusters, when `enable_private_nodes` is `true`.

* `master_ipv4_cidr_block` (Optional) - The IP range in CIDR notation to use for
the hosted master network. This range will be used for assigning private IP
addresses to the cluster master(s) and the ILB VIP. This range must not overlap
with any other ranges in use within the cluster's network, and it must be a /28
subnet. See [Private Cluster Limitations](https://cloud.google.com/kubernetes-engine/docs/how-to/private-clusters#limitations)
for more details. This field only applies to private clusters, when
`enable_private_nodes` is `true`.

* `master_global_access_config` (Optional) - Controls cluster master global
access settings. If unset, Terraform will no longer manage this field and will
not modify the previously-set value. Structure is documented below.

In addition, the `private_cluster_config` allows access to the following read-only fields:

* `peering_name` - The name of the peering between this cluster and the Google owned VPC.

* `private_endpoint` - The internal IP address of this cluster's master endpoint.

* `public_endpoint` - The external IP address of this cluster's master endpoint.

!> The Google provider is unable to validate certain configurations of
`private_cluster_config` when `enable_private_nodes` is `false`. It's
recommended that you omit the block entirely if the field is not set to `true`.

The `private_cluster_config.master_global_access_config` block supports:

* `enabled` (Optional) - Whether the cluster master is accessible globally or
not.

The `sandbox_config` block supports:

* `sandbox_type` (Required) Which sandbox to use for pods in the node pool.
    Accepted values are:

    * `"gvisor"`: Pods run within a gVisor sandbox.

The `release_channel` block supports:

* `channel` - (Required) The selected release channel.
    Accepted values are:
    * UNSPECIFIED: Not set.
    * RAPID: Weekly upgrade cadence; Early testers and developers who requires new features.
    * REGULAR: Multiple per month upgrade cadence; Production users who need features not yet offered in the Stable channel.
    * STABLE: Every few months upgrade cadence; Production users who need stability above all else, and for whom frequent upgrades are too risky.

The `resource_usage_export_config` block supports:

* `enable_network_egress_metering` (Optional) - Whether to enable network egress metering for this cluster. If enabled, a daemonset will be created
    in the cluster to meter network egress traffic.

* `enable_resource_consumption_metering` (Optional) - Whether to enable resource
consumption metering on this cluster. When enabled, a table will be created in
the resource export BigQuery dataset to store resource consumption data. The
resulting table can be joined with the resource usage table or with BigQuery
billing export. Defaults to `true`.

* `bigquery_destination` (Required) - Parameters for using BigQuery as the destination of resource usage export.

* `bigquery_destination.dataset_id` (Required) - The ID of a BigQuery Dataset. For Example:

```hcl
resource_usage_export_config {
  enable_network_egress_metering = false
  enable_resource_consumption_metering = true

  bigquery_destination {
    dataset_id = "cluster_resource_usage"
  }
}
```

The `shielded_instance_config` block supports:

* `enable_secure_boot` (Optional) - Defines if the instance has Secure Boot enabled.

Secure Boot helps ensure that the system only runs authentic software by verifying the digital signature of all boot components, and halting the boot process if signature verification fails.  Defaults to `false`.

* `enable_integrity_monitoring` (Optional) - Defines if the instance has integrity monitoring enabled.

Enables monitoring and attestation of the boot integrity of the instance. The attestation is performed against the integrity policy baseline. This baseline is initially derived from the implicitly trusted boot image when the instance is created.  Defaults to `true`.

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
    * GKE_METADATA_SERVER: Enables [workload identity](https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity) on the node.

The `kubelet_config` block supports:

* `cpu_manager_policy` - (Required) The CPU management policy on the node. See
[K8S CPU Management Policies](https://kubernetes.io/docs/tasks/administer-cluster/cpu-management-policies/).
One of `"none"` or `"static"`. Defaults to `none` when `kubelet_config` is unset.

* `cpu_cfs_quota` - (Optional) If true, enables CPU CFS quota enforcement for
containers that specify CPU limits.

* `cpu_cfs_quota_period` - (Optional) The CPU CFS quota period value. Specified
as a sequence of decimal numbers, each with optional fraction and a unit suffix,
such as `"300ms"`. Valid time units are "ns", "us" (or "µs"), "ms", "s", "m",
"h". The value must be a positive duration.

-> Note: At the time of writing (2020/08/18) the GKE API rejects the `none`
value and accepts an invalid `default` value instead. While this remains true,
not specifying the `kubelet_config` block should be the equivalent of specifying
`none`.

The `linux_node_config` block supports:

* `sysctls` - (Required)  The Linux kernel parameters to be applied to the nodes
and all pods running on the nodes. Specified as a map from the key, such as
`net.core.wmem_max`, to a string value.

The `vertical_pod_autoscaling` block supports:

* `enabled` (Required) - Enables vertical pod autoscaling

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `id` - an identifier for the resource with format `projects/{{project}}/locations/{{zone}}/clusters/{{name}}`

* `self_link` - The server-defined URL for the resource.

* `endpoint` - The IP address of this cluster's Kubernetes master.

* `instance_group_urls` - List of instance group URLs which have been assigned
    to the cluster.

* `label_fingerprint` - The fingerprint of the set of labels for this cluster.

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

* `tpu_ipv4_cidr_block` - The IP address range of the Cloud TPUs in this cluster, in
    [CIDR](http://en.wikipedia.org/wiki/Classless_Inter-Domain_Routing)
    notation (e.g. `1.2.3.4/29`).

* `services_ipv4_cidr` - The IP address range of the Kubernetes services in this
  cluster, in [CIDR](http://en.wikipedia.org/wiki/Classless_Inter-Domain_Routing)
  notation (e.g. `1.2.3.4/29`). Service addresses are typically put in the last
  `/16` from the container CIDR.

## Timeouts

This resource provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 40 minutes.
- `read`   - Default is 40 minutes.
- `update` - Default is 60 minutes.
- `delete` - Default is 40 minutes.

## Import

GKE clusters can be imported using the `project` , `location`, and `name`. If the project is omitted, the default
provider value will be used. Examples:

```
$ terraform import google_container_cluster.mycluster projects/my-gcp-project/locations/us-east1-a/clusters/my-cluster

$ terraform import google_container_cluster.mycluster my-gcp-project/us-east1-a/my-cluster

$ terraform import google_container_cluster.mycluster us-east1-a/my-cluster
```

~> **Note:** This resource has several fields that control Terraform-specific behavior and aren't present in the API. If they are set in config and you import a cluster, Terraform may need to perform an update immediately after import. Most of these updates should be no-ops but some may modify your cluster if the imported state differs.

For example, the following fields will show diffs if set in config:

- `min_master_version`
- `remove_default_node_pool`

## User Project Overrides

This resource supports [User Project Overrides](https://www.terraform.io/docs/providers/google/guides/provider_reference.html#user_project_override).
