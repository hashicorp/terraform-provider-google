---
layout: "google"
page_title: "Terraform Google Provider 3.0.0 Upgrade Guide"
sidebar_current: "docs-google-provider-version-3-upgrade"
description: |-
  Terraform Google Provider 3.0.0 Upgrade Guide
---

# Terraform Google Provider 3.0.0 Upgrade Guide

The `3.0.0` release of the Google provider for Terraform is a major version and
includes some changes that you will need to consider when upgrading. This guide
is intended to help with that process and focuses only on the changes necessary
to upgrade from the final `2.X` series release to `3.0.0`.

Most of the changes outlined in this guide have been previously marked as
deprecated in the Terraform `plan`/`apply` output throughout previous provider
releases, up to and including the final `2.X` series release. These changes,
such as deprecation notices, can always be found in the CHANGELOG of the
affected providers. [google](https://github.com/terraform-providers/terraform-provider-google/blob/master/CHANGELOG.md)
[google-beta](https://github.com/terraform-providers/terraform-provider-google-beta/blob/master/CHANGELOG.md)

## I accidentally upgraded to 3.0.0, how do I downgrade to `2.X`?

If you've inadvertently upgraded to `3.0.0`, first see the
[Provider Version Configuration Guide](#provider-version-configuration) to lock
your provider version; if you've constrained the provider to a lower version
such as shown in the previous version example in that guide, Terraform will pull
in a `2.X` series release on `terraform init`.

If you've only ran `terraform init` or `terraform plan`, your state will not
have been modified and downgrading your provider is sufficient.

If you've ran `terraform refresh` or `terraform apply`, Terraform may have made
state changes in the meantime.

* If you're using a local state, or a remote state backend that does not support
versioning, `terraform refresh` with a downgraded provider is likely sufficient
to revert your state. The Google provider generally refreshes most state
information from the API, and the properties necessary to do so have been left
unchanged.

* If you're using a remote state backend that supports versioning such as
[Google Cloud Storage](https://www.terraform.io/docs/backends/types/gcs.html),
you can revert the Terraform state file to a previous version. If you do
so and Terraform had created resources as part of a `terraform apply` in the
meantime, you'll need to either delete them by hand or `terraform import` them
so Terraform knows to manage them.

## Upgrade Topics

<!-- TOC depthFrom:2 depthTo:2 -->

- [Provider Version Configuration](#provider-version-configuration)
- [Data Source: `google_container_engine_versions`](#data-source-google_container_engine_versions)
- [Resource: `google_app_engine_application`](#resource-google_app_engine_application)
- [Resource: `google_cloudfunctions_function`](#resource-google_cloudfunctions_function)
- [Resource: `google_cloudiot_registry`](#resource-google_cloudiot_registry)
- [Resource: `google_composer_environment`](#resource-google_composer_environment)
- [Resource: `google_compute_forwarding_rule`](#resource-google_compute_forwarding_rule)
- [Resource: `google_compute_instance`](#resource-google_compute_instance)
- [Resource: `google_compute_instance_template`](#resource-google_compute_instance_template)
- [Resource: `google_compute_network`](#resource-google_compute_network)
- [Resource: `google_compute_network_peering`](#resource-google_compute_network_peering)
- [Resource: `google_compute_region_instance_group_manager`](#resource-google_compute_region_instance_group_manager)
- [Resource: `google_compute_router_peer`](#resource-google_compute_router_peer)
- [Resource: `google_compute_snapshot`](#resource-google_compute_snapshot)
- [Resource: `google_container_cluster`](#resource-google_container_cluster)
- [Resource: `google_container_node_pool`](#resource-google_container_node_pool)
- [Resource: `google_dataproc_cluster`](#resource-google_dataproc_cluster)
- [Resource: `google_dataproc_job`](#resource-google_dataproc_job)
- [Resource: `google_dns_managed_zone`](#resource-google_dns_managed_zone)
- [Resource: `google_monitoring_alert_policy`](#resource-google_monitoring_alert_policy)
- [Resource: `google_monitoring_uptime_check_config`](#resource-google_monitoring_uptime_check_config)
- [Resource: `google_organization_policy`](#resource-google_organization_policy)
- [Resource: `google_project_services`](#resource-google_project_services)
- [Resource: `google_sql_database_instance`](#resource-google_sql_database_instance)
- [Resource: `google_storage_bucket`](#resource-google_storage_bucket)
- [Resource: `google_storage_transfer_job`](#resource-google_storage_transfer_job)

<!-- /TOC -->

## Provider Version Configuration

-> Before upgrading to version 3.0.0, it is recommended to upgrade to the most
recent `2.X` series release of the provider and ensure that your environment
successfully runs [`terraform plan`](https://www.terraform.io/docs/commands/plan.html)
without unexpected changes or deprecation notices.

It is recommended to use [version constraints](https://www.terraform.io/docs/configuration/providers.html#provider-versions)
when configuring Terraform providers. If you are following that recommendation,
update the version constraints in your Terraform configuration and run
[`terraform init`](https://www.terraform.io/docs/commands/init.html) to download
the new version.

If you aren't using version constraints, you can use `terraform init -upgrade`
in order to upgrade your provider to the latest released version.

For example, given this previous configuration:

```hcl
provider "google" {
  # ... other configuration ...

  version = "~> 2.17.0"
}
```

An updated configuration:

```hcl
provider "google" {
  # ... other configuration ...

  version = "~> 3.0.0"
}
```

## Data Source: `google_container_engine_versions`

### `region` and `zone` are now removed

Use `location` instead.

## Resource: `google_container_cluster`

### Automatic subnetwork creation for VPC-native clusters removed

Automatic creation of subnetworks in GKE has been removed. Now, users of
VPC-native clusters will always need to provide a `google_compute_subnetwork`
resource to use `ip_allocation_policy`. Routes-based clusters are unaffected.

Representing resources managed by another source in Terraform is painful, and
leads to confusing patterns that often involve unnecessarily recreating user
resources. A number of fields in GKE are dedicated to a feature that allows
users to create a GKE-managed subnetwork.

This is a great fit for an imperative tool like `gcloud`, but it's not required
for Terraform. With Terraform, it's relatively easy to specify a subnetwork in
config alongside the cluster. Not only does that allow configuring subnetwork
features like flow logging, it's more explicit, allows the subnetwork to be used
by other resources, and the subnetwork persists through cluster deletion.

Particularly, Shared VPC was incompatible with `create_subnetwork`, and
`node_ipv4_cidr` was easy to confuse with
`ip_allocation_policy.node_ipv4_cidr_block`.

#### Detailed changes:

* `ip_allocation_policy.node_ipv4_cidr_block` removed (This controls the primary range of the created subnetwork)
* `ip_allocation_policy.create_subnetwork`, `ip_allocation_policy.subnetwork_name` removed
* `ip_allocation_policy.use_ip_aliases` removed
  * Enablement is now based on `ip_allocation_policy` being defined instead
* Conflict added between `node_ipv4_cidr`, `ip_allocation_policy`

#### Upgrade instructions

1. Remove the removed fields from `google_container_cluster`
1. Add a `google_compute_subnetwork` to your config, import it using `terraform import`
1. Reference the subnetwork using the `subnetwork` field on your `google_container_cluster`

-> Subnetworks originally created as part of `create_subnetwork` will be deleted
alongside the cluster. If there are other users of the subnetwork, deletion of
the cluster will fail. After the original resources are deleted,
`terraform apply` will recreate the same subnetwork except that it won't be
managed by a GKE cluster and other resources can use it safely.

#### Old Config

```hcl
resource "google_compute_network" "container_network" {
  name                    = "container-network"
  auto_create_subnetworks = false
}

resource "google_container_cluster" "primary" {
  name       = "my-cluster"
  location   = "us-central1"
  network    = "${google_compute_network.container_network.name}"

  initial_node_count = 1

  ip_allocation_policy {
    use_ip_aliases           = true
    create_subnetwork        = true
    cluster_ipv4_cidr_block  = "10.0.0.0/16"
    services_ipv4_cidr_block = "10.1.0.0/16"
    node_ipv4_cidr_block     = "10.2.0.0/16"
  }
}
```

#### New Config

```hcl
resource "google_compute_network" "container_network" {
  name                    = "container-network"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "container_subnetwork" {
  name          = "container-subnetwork"
  description   = "auto-created subnetwork for cluster \"my-cluster\""
  ip_cidr_range = "10.2.0.0/16"
  region        = "us-central1"
  network       = "${google_compute_network.container_network.self_link}"
}

resource "google_container_cluster" "primary" {
  name       = "my-cluster"
  location   = "us-central1"
  network    = "${google_compute_network.container_network.name}"
  subnetwork = "${google_compute_subnetwork.container_subnetwork.name}"

  initial_node_count = 1

  ip_allocation_policy {
    use_ip_aliases           = true
    cluster_ipv4_cidr_block  = "10.0.0.0/16"
    services_ipv4_cidr_block = "10.1.0.0/16"
  }
}
```

### `logging_service` and `monitoring_service` defaults changed

GKE Stackdriver Monitoring (the GKE-specific Stackdriver experience) is now
enabled at cluster creation by default, similar to the default in GKE `1.14`
through other tools.

## Resource: `google_app_engine_application`

### `split_health_checks` is now required on block `google_app_engine_application.feature_settings`

In an attempt to avoid allowing empty blocks in config files, `split_health_checks` is now
required on the `google_app_engine_application.feature_settings` block.

### `taint` field is now authoritative when set

The `taint` field inside of `node_config` blocks on `google_container_cluster`
and `google_container_node_pool` will no longer ignore GPU-related values when
set.

Previously, the field ignored upstream taints when unset and ignored unset GPU
taints when other taints were set. Now it will ignore upstream taints when set
and act authoritatively when set, requiring all taints (including Kubernetes and
GKE-managed ones) to be defined in config.

Additionally, an empty taint can now be specified with `taint = []`. As a result
of this change, the JSON/state representation of the field has changed,
introducing an incompatibility for users who specify config in JSON instead of
HCL or who use `dynamic` blocks. See more details in the [Attributes as Blocks](https://www.terraform.io/docs/configuration/attr-as-blocks.html)
documentation.

## Resource: `google_cloudfunctions_function`

### The `runtime` option `nodejs6` has been deprecated

`nodejs6` has been deprecated and is no longer the default value for `runtime`.
`runtime` is now required.

## Resource: `google_cloudiot_registry`

### `event_notification_config` is now removed

`event_notification_config` has been removed in favor of
`event_notification_configs` (plural). Please switch to using the plural field.

### `public_key_certificate` is now required on block `google_cloudiot_registry.credentials`

In an attempt to avoid allowing empty blocks in config files, `public_key_certificate` is now
required on the `google_cloudiot_registry.credentials` block.

## Resource: `google_composer_environment`

### `use_ip_aliases` is now required on block `google_composer_environment.ip_allocation_policy`

Previously the default value of `use_ip_aliases` was `true`. In an attempt to avoid allowing empty blocks
in config files, `use_ip_aliases` is now required on the `google_composer_environment.ip_allocation_policy` block.

### `enable_private_endpoint` is now required on block `google_composer_environment.private_environment_config`

Previously the default value of `enable_private_endpoint` was `true`. In an attempt to avoid allowing empty blocks
in config files, `enable_private_endpoint` is now required on the `google_composer_environment.private_environment_config` block.

## Resource: `google_compute_forwarding_rule`

### `ip_version` is now removed

`ip_version` is not used for regional forwarding rules.

## Resource: `google_compute_instance`

### `interface` is now required on block `google_compute_instance.scratch_disk`

Previously the default value of `interface` was `SCSI`. In an attempt to avoid allowing empty blocks
in config files, `interface` is now required on the `google_compute_instance.scratch_disk` block.

## Resource: `google_compute_instance_template`

### `kms_key_self_link` is now required on block `google_compute_instance_template.disk_encryption_key`

In an attempt to avoid allowing empty blocks in config files, `kms_key_self_link` is now
required on the `google_compute_instance_template.disk_encryption_key` block.

## Resource: `google_compute_network`

### `ipv4_range` is now removed

Legacy Networks are removed and you will no longer be able to create them
using this field from Feb 1, 2020 onwards.

## Resource: `google_compute_network_peering`

### `auto_create_routes` is now removed

`auto_create_routes` has been removed because it's redundant and not
user-configurable.

## Resource: `google_compute_region_instance_group_manager`

### `update_strategy` no longer has any effect and is removed

With `rolling_update_policy` removed, `update_strategy` has no effect anymore.
Before updating, remove it from your config.

## Resource: `google_compute_router_peer`

### `range` is now required on block `google_compute_router_peer.advertised_ip_ranges`

In an attempt to avoid allowing empty blocks in config files, `range` is now
required on the `google_compute_router_peer.advertised_ip_ranges` block.

## Resource: `google_compute_snapshot`

### `raw_key` is now required on block `google_compute_snapshot.source_disk_encryption_key`

In an attempt to avoid allowing empty blocks in config files, `raw_key` is now
required on the `google_compute_snapshot.source_disk_encryption_key` block.

## Resource: `google_container_cluster`


### `addons_config.kubernetes_dashboard` is now removed

The `kubernetes_dashboard` addon is deprecated for clusters on GKE and
will soon be removed. It is recommended to use alternative GCP Console
dashboards.

### `cidr_blocks` is now required on block `google_container_cluster.master_authorized_networks_config`

In an attempt to avoid allowing empty blocks in config files, `cidr_blocks` is now
required on the `google_container_cluster.master_authorized_networks_config` block.

### The `disabled` field is now required on the `addons_config` blocks for
`http_load_balancing`, `horizontal_pod_autoscaling`, `istio_config`,
`cloudrun_config` and `network_policy_config`.

In an attempt to avoid allowing empty blocks in config files, `disabled` is now
required on the different `google_container_cluster.addons_config` blocks.

### `enabled` is now required on block `google_container_cluster.vertical_pod_autoscaling`

In an attempt to avoid allowing empty blocks in config files, `enabled` is now
required on the `google_container_cluster.vertical_pod_autoscaling` block.

### `enabled` is now required on block `google_container_cluster.network_policy`

Previously the default value of `enabled` was `false`. In an attempt to avoid allowing empty blocks
in config files, `enabled` is now required on the `google_container_cluster.network_policy` block.

### `enable_private_endpoint` is now required on block `google_container_cluster.private_cluster_config`

In an attempt to avoid allowing empty blocks in config files, `enable_private_endpoint` is now
required on the `google_container_cluster.private_cluster_config` block.

### `logging_service` and `monitoring_service` defaults changed

GKE Stackdriver Monitoring (the GKE-specific Stackdriver experience) is now
enabled at cluster creation by default, similar to the default in GKE `1.14`
through other tools.

Terraform will now detect changes out of band when the field(s) are not defined
in config, attempting to return them to their new defaults, and will be clear
about what values will be set when creating a cluster.

`terraform plan` will report changes upon upgrading if the field was previously
unset. Applying this change will enable the new Stackdriver service without
recreating clusters. Users who wish to use another value should record their
intended value in config; the old default values can be added to a
`google_container_cluster` resource config block to preserve them.

#### Old Defaults

```hcl
logging_service    = "logging.googleapis.com"
monitoring_service = "monitoring.googleapis.com"
```

#### New Defaults

```hcl
logging_service    = "logging.googleapis.com/kubernetes"
monitoring_service = "monitoring.googleapis.com/kubernetes"
```

### `use_ip_aliases` is now required on block `google_container_cluster.ip_allocation_policy`

Previously the default value of `use_ip_aliases` was `true`. In an attempt to avoid allowing empty blocks
in config files, `use_ip_aliases` is now required on the `google_container_cluster.ip_allocation_policy` block.

### `zone`, `region` and `additional_zones` are now removed

`zone` and `region` have been removed in favor of `location` and
`additional_zones` has been removed in favor of `node_locations`

## Resource: `google_container_node_pool`

### `zone` and `region` are now removed

`zone` and `region` have been removed in favor of `location`

## Resource: `google_dataproc_cluster`

### `policy_uri` is now required on `google_dataproc_cluster.autoscaling_config` block.

In an attempt to avoid allowing empty blocks in config files, `policy_uri` is now
required on the `google_dataproc_cluster.autoscaling_config` block.

## Resource: `google_dataproc_job`

### `driver_log_levels` is now required on `logging_config` blocks for
`google_dataproc_job.pyspark_config`, `google_dataproc_job.hadoop_config`,
`google_dataproc_job.spark_config`, `google_dataproc_job.pig_config`, and
`google_dataproc_job.sparksql_config`.

In an attempt to avoid allowing empty blocks in config files, `driver_log_levels` is now
required on the different `google_dataproc_job` config blocks.

### `max_failures_per_hour` is now required on block `google_dataproc_job.scheduling`

In an attempt to avoid allowing empty blocks in config files, `max_failures_per_hour` is now
required on the `google_dataproc_job.scheduling` block.

## Resource: `google_dns_managed_zone`

### `networks` is now required on block `google_dns_managed_zone.private_visibility_config`

In an attempt to avoid allowing empty blocks in config files, `networks` is now
required on the `google_dns_managed_zone.private_visibility_config` block.

### `network_url` is now required on block `google_dns_managed_zone.private_visibility_config.networks`

In an attempt to avoid allowing empty blocks in config files, `network_url` is now
required on the `google_dns_managed_zone.private_visibility_config.networks` block.

## Resource: `google_monitoring_alert_policy`

### `labels` is now removed

`labels` is removed as it was never used. See `user_labels` for the correct field.

## Resource: `google_monitoring_uptime_check_config`

### `content` is now required on block `google_monitoring_uptime_check_config.content_matchers`

In an attempt to avoid allowing empty blocks in config files, `content` is now
required on the `google_monitoring_uptime_check_config.content_matchers` block.

### `is_internal` and `internal_checker` are now removed

`is_internal` and `internal_checker` never worked, and are now removed.

## Resource: `google_organization_policy`

### `inherit_from_parent` is now required on block `google_organization_policy.list_policy`

In an attempt to avoid allowing empty blocks in config files, `inherit_from_parent` is now
required on the `google_organization_policy.list_policy` block.

## Resource: `google_project_services`

### `google_project_services` has been removed from the provider

The `google_project_services` resource was authoritative over the list of GCP
services enabled on a project, so that services not explicitly set would be
removed by Terraform.

However, this was dangerous to use in practice. Services have dependencies that
are automatically enabled alongside them and GCP will add dependencies to
services out of band, enabling them. If a user ran Terraform after this,
Terraform would disable the service- and implicitly disable any service that
relied on it.

The `google_project_service` resource is a much better match for most users'
intent, managing a single service at a time. Setting several
`google_project_service` resources is an assertion that "these services are set
on this project", while `google_project_services` was an assertion that "**only**
these services are set on this project".

Users should migrate to using `google_project_service` resources, or using the
[`"terraform-google-modules/project-factory/google//modules/project_services"`](https://registry.terraform.io/modules/terraform-google-modules/project-factory/google/3.3.0/submodules/project_services)
module for a similar interface to `google_project_services`.

-> Prior to `2.13.0`, each `google_project_service` sent separate API enablement
requests. From `2.13.0` onwards, those requests are batched. It's recommended
that you upgrade to `2.13.0+` before migrating if you encounter quota issues
when you migrate off `google_project_services`.

#### Old Config

```hcl
resource "google_project_services" "project" {
  project            = "your-project-id"
  services           = ["iam.googleapis.com", "cloudresourcemanager.googleapis.com"]
  disable_on_destroy = false
}
```

#### New Config (module)

```hcl
module "project_services" {
  source  = "terraform-google-modules/project-factory/google//modules/project_services"
  version = "3.3.0"

  project_id    = "your-project-id"
  activate_apis =  [
    "iam.googleapis.com",
    "cloudresourcemanager.googleapis.com",
  ]

  disable_services_on_destroy = false
  disable_dependent_services  = false
}
```

#### New Config (google_project_service)

```hcl
resource "google_project_service" "project_iam" {
  project = "your-project-id"
  service = "iam.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "project_cloudresourcemanager" {
  project = "your-project-id"
  service = "cloudresourcemanager.googleapis.com"
  disable_on_destroy = false
}
```

## Resource: `google_sql_database_instance`

### `dump_file_path`, `username` and `password` are now required on block `google_sql_database_instance.replica_configuration`

In an attempt to avoid allowing empty blocks in config files, `dump_file_path`, `username` and `password` are now
required on the `google_sql_database_instance.replica_configuration` block.

### `name` and `value` are now required on block `google_sql_database_instance.settings.database_flags`

In an attempt to avoid allowing empty blocks in config files, `name` and `value` are now
required on the `google_sql_database_instance.settings.database_flags` block.

### `value` is now required on block `google_sql_database_instance.settings.ip_configuration.authorized_networks`

In an attempt to avoid allowing empty blocks in config files, `value` is now
required on the `google_sql_database_instance.settings.ip_configuration.authorized_networks` block.

### `zone` is now required on block `google_sql_database_instance.settings.location_preference`

In an attempt to avoid allowing empty blocks in config files, `zone` is now
required on the `google_sql_database_instance.settings.location_preference` block.

## Resource: `google_storage_bucket`

### `enabled` is now required on block `google_storage_bucket.versioning`

Previously the default value of `enabled` was `false`. In an attempt to avoid allowing empty blocks
in config files, `enabled` is now required on the `google_storage_bucket.versioning` block.

### `is_live` is now removed

Please use `with_state` instead, as `is_live` is now removed.

## Resource: `google_storage_transfer_job`

### `overwrite_objects_already_existing_in_sink` is now required on block `google_storage_transfer_job.transfer_options`

In an attempt to avoid allowing empty blocks in config files, `overwrite_objects_already_existing_in_sink` is now
required on the `google_storage_transfer_job.transfer_options` block.