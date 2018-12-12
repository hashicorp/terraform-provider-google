---
layout: "google"
page_title: "Terraform Google Provider Version 2 Upgrade Guide"
sidebar_current: "docs-google-provider-version-2-upgrade"
description: |-
  Terraform Google Provider Version 2 Upgrade Guide
---

# Terraform Google Provider Version 2 Upgrade Guide

Version 2.0.0 of the Google provider for Terraform is a major release and includes some changes that you will need to consider when upgrading. This guide is intended to help with that process and focuses only on changes from version 1.19.1 to version 2.0.0.

Most of the changes outlined in this guide have been previously marked as deprecated in the Terraform plan/apply output throughout previous provider releases, up to and including 1.19.1. These changes, such as deprecation notices, can always be found in the [Terraform Google Provider CHANGELOG](https://github.com/terraform-providers/terraform-provider-google/blob/master/CHANGELOG.md).

Upgrade topics:

<!-- TOC depthFrom:2 depthTo:2 -->

- [Provider Version Configuration](#provider-version-configuration)
- [`google-beta` provider](#google-beta-provider)
- [Open in Cloud Shell](#open-in-cloud-shell)
- [Data Sources](#data-sources)
- [Resource: `google_bigquery_dataset`](#resource-google_bigquery_dataset)
- [Resource: `google_bigtable_instance`](#resource-google_bigtable_instance)
- [Resource: `google_binary_authorizaton_attestor`](#resource-google_binary_authorization_attestor)
- [Resource: `google_binary_authorizaton_policy`](#resource-google_binary_authorization_policy)
- [Resource: `google_cloudfunctions_function`](#resource-google_cloudfunctions_function)
- [Resource: `google_compute_backend_service`](#resource-google_compute_backend_service)
- [Resource: `google_compute_disk`](#resource-google_compute_disk)
- [Resource: `google_compute_global_forwarding_rule`](#resource-google_compute_global_forwarding_rule)
- [Resource: `google_compute_image`](#resource-google_compute_image)
- [Resource: `google_compute_instance`](#resource-google_compute_instance)
- [Resource: `google_compute_instance_from_template`](#resource-google_compute_instance_from_template)
- [Resource: `google_compute_instance_group_manager`](#resource-google_compute_instance_group_manager)
- [Resource: `google_compute_project_metadata`](#resource-google_compute_project_metadata)
- [Resource: `google_compute_region_instance_group_manager`](#resource-google_compute_region_instance_group_manager)
- [Resource: `google_compute_subnetwork_iam_*`](#resource-google_compute_subnetwork_iam_*)
- [Resource: `google_compute_target_pool`](#resource-google_compute_target_pool)
- [Resource: `google_compute_url_map`](#resource-google_compute_url_map)
- [Resource: `google_container_analysis_note`](#resource-google_container_analysis_note)
- [Resource: `google_container_cluster`](#resource-google_container_cluster)
- [Resource: `google_container_node_pool`](#resource-google_container_node_pool)
- [Resource: `google_dataproc_cluster`](#resource-google_dataproc_cluster)
- [Resource: `google_endpoints_service`](#resource-google_endpoints_service)
- [Resource: `google_filestore_instance`](#resource-google_filestore_instance)
- [Resource: `google_organization_custom_role`](#resource-google_organization_custom_role)
- [Resource: `google_project`](#resource-google_project)
- [Resource: `google_project_iam_policy`](#resource-google_project_iam_policy)
- [Resource: `google_service_account`](#resource-google_service_account)
- [Resource: `google_sql_database_instance`](#resource-google_sql_database_instance)
- [Resource: `google_storage_default_object_acl`](#resource-google_storage_default_object_acl)
- [Resource: `google_storage_object_acl`](#resource-google_storage_object_acl)

<!-- /TOC -->

## Provider Version Configuration

!> **WARNING:** This topic is placeholder documentation until version 2.0.0 is released later this year.

-> Before upgrading to version 2.0.0, it is recommended to upgrade to the most recent version of the provider (1.19.1) and ensure that your environment successfully runs [`terraform plan`](https://www.terraform.io/docs/commands/plan.html) without unexpected changes or deprecation notices.

It is recommended to use [version constraints when configuring Terraform providers](https://www.terraform.io/docs/configuration/providers.html#provider-versions). If you are following that recommendation, update the version constraints in your Terraform configuration and run [`terraform init`](https://www.terraform.io/docs/commands/init.html) to download the new version.

For example, given this previous configuration:

```hcl
provider "google" {
  # ... other configuration ...

  version = "~> 1.19.0"
}
```

An updated configuration:

```hcl
provider "google" {
  # ... other configuration ...

  version = "~> 2.0.0"
}
```

## `google-beta` provider

The `google-beta` provider is now necessary to be able to configure resources with beta features.
This new provider enables full import support of beta features and gives users who
wish to use only the most stable APIs and features more confidence that they are doing so
by continuing to use the `google` provider.

Beta GCP Features have no deprecation policy and no SLA, but are otherwise considered to be feature-complete
with only minor outstanding issues after their Alpha period. Beta is when GCP
features are publicly announced, and is when they generally become publicly
available. For more information see [the official documentation on GCP launch stages](https://cloud.google.com/terms/launch-stages).

Because the API for beta features can change before their GA launch, there may be breaking changes
in the `google-beta` provider in minor release versions. These changes will be announced in the
[Terraform `google-beta` Provider CHANGELOG](https://github.com/terraform-providers/terraform-provider-google-beta/blob/master/CHANGELOG.md).

To have resources at different API versions, set up provider blocks for each version:

```hcl
provider "google" {
  credentials = "${file("account.json")}"
  project     = "my-project-id"
  region      = "us-central1"
}

provider "google-beta" {
  credentials = "${file("account.json")}"
  project     = "my-project-id"
  region      = "us-central1"
}
```

In each resource, state which provider that resource should be used with:

```hcl
resource "google_compute_instance" "ga-instance" {
  provider = "google"

  # ...
}

resource "google_compute_instance" "beta-instance" {
  provider = "google-beta"

  # ...
}
```

See [Provider Versions](https://terraform.io/docs/providers/google/provider_versions.html)
for more details on how to use `google-beta`.

## Open in Cloud Shell

2.0.0 is the first release including Open in Cloud Shell. Examples in the documentation for
Magic Modules resources now have Open in Cloud Shell links in their documentation that open
them in an interactive editor and shell - all without leaving the browser. See the
[blog post announcing the feature](https://www.hashicorp.com/blog/kickstart-terraform-on-gcp-with-google-cloud-shell)
for more details.

## Data Sources

See the `Resource` sections in this document for properties that may have been removed.

## Resource: `google_bigquery_dataset`

### `access` is now a Set

The order of entries in `access` no longer matters. Any configurations that
interpolate based on an item at a specific index will need to be updated, as items
may have been reordered.

## Resource: `google_bigtable_instance`

### `cluster_id`, `zone`, `num_nodes`, and `storage_type` have moved into a `cluster` block

Example previous configuration:

```hcl
resource "google_bigtable_instance" "instance" {
  name         = "tf-instance"
  cluster_id   = "tf-instance-cluster"
  zone         = "us-central1-b"
  num_nodes    = 3
  storage_type = "HDD"
}
```

Example updated configuration:

```hcl
resource "google_bigtable_instance" "instance" {
  name = "tf-instance"
  cluster {
    cluster_id   = "tf-instance-cluster"
    zone         = "us-central1-b"
    num_nodes    = 3
    storage_type = "HDD"
  }
}
```

### `zone` is no longer inferred from the provider

`cluster.zone` is now required, even if the provider block has a zone set.

## Resource: `google_binary_authorization_attestor`

### binary authorization resources have been removed from the GA provider

Use the [`google-beta` provider](#google-beta-provider) to use these resources.

## Resource: `google_binary_authorization_policy`

### binary authorization resources have been removed from the GA provider

Use the [`google-beta` provider](#google-beta-provider) to use these resources.

## Resource: `google_cloudfunctions_function`

### `trigger_bucket`, `trigger_topic`, and `retry_on_failure` have been removed

Use the `event_trigger` block instead.

Example updated configuration:

```hcl
resource "google_cloudfunctions_function" "function" {
  name                  = "example-function"
  available_memory_mb   = 128
  source_archive_bucket = "${google_storage_bucket.bucket.name}"
  source_archive_object = "${google_storage_bucket_object.archive.name}"
  timeout               = 61
  entry_point           = "helloGCS"

  event_trigger {
    event_type = "providers/cloud.storage/eventTypes/object.change"
    resource   = "${google_storage_bucket.bucket.name}"
    failure_policy {
      retry = true
    }
  }
}

resource "google_storage_bucket" "bucket" {
  name = "example-bucket"
}

resource "google_storage_bucket_object" "archive" {
  name   = "index.zip"
  bucket = "${google_storage_bucket.bucket.name}"
  source = "path/to/source.zip"
}
```

See the documentation at
[`google_cloudfunctions_function`](https://www.terraform.io/docs/providers/google/r/cloudfunctions_function.html)
for more details.

## Resource: `google_compute_backend_service`

### `custom_request_headers` has been removed from the GA provider

Use the [`google-beta` provider](#google-beta-provider) to set this field.

## Resource: `google_compute_disk`

### `disk_encryption_key_raw` and `disk_encryption_key_sha256` have been removed.

Use the `disk_encryption_key` block instead:

```hcl
data "google_compute_image" "my_image" {
  family  = "debian-9"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
  name = "example-disk"
  image = "${data.google_compute_image.my_image.self_link}"
  size = 50
  type = "pd-ssd"
  zone = "us-central1-a"
  disk_encryption_key {
    raw_key = "SGVsbG8gZnJvbSBHb29nbGUgQ2xvdWQgUGxhdGZvcm0="
  }
}
```

## Resource: `google_compute_global_forwarding_rule`

### `labels` has been removed from the GA provider

Use the [`google-beta` provider](#google-beta-provider) to set this field.

## Resource: `google_compute_image`

### `create_timeout` has been removed

Use the standard [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts)
block instead.

## Resource: `google_compute_instance`

### `create_timeout` has been removed

Use the standard [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts)
block instead.

### `metadata` is now authoritative

Terraform will remove values not explicitly set in this field. Any `metadata` values
that were added outside of Terraform should be added to the config.

### `network` has been removed

Use `network_interface` instead.

### `network_interface.*.address` has been removed

Use `network_interface.*.network_ip` instead.

## Resource: `google_compute_instance_from_template`

### `metadata` is now authoritative

Terraform will remove values not explicitly set in this field. Any `metadata` values
that were added outside of Terraform should be added to the config.

## Resource: `google_compute_instance_group_manager`

### `version`, `auto_healing_policies`, `rolling_update_policy` have been removed from the GA provider

Use the [`google-beta` provider](#google-beta-provider) to use these fields.
`rolling_update_policy` has been renamed to `update_policy` in `google-beta`.

## Resource: `google_compute_project_metadata`

### `metadata` is now authoritative

Terraform will remove values not explicitly set in this field. Any `metadata` values
that were added outside of Terraform should be added to the config.

## Resource: `google_compute_region_instance_group_manager`

### `version`, `auto_healing_policies`, `rolling_update_policy` have been removed from the GA provider

Use the [`google-beta` provider](#google-beta-provider) to use these fields.
`rolling_update_policy` has been renamed to `update_policy` in `google-beta`.

### `update_strategy` no longer has any effect and is deprecated

With `rolling_update_policy` removed, `update_strategy` has no effect anymore.
Remove it from your config at your convenience.

## Resource: `google_compute_subnetwork_iam_*`

### subnetwork IAM resources have been removed from the GA provider

Use the [`google-beta` provider](#google-beta-provider) to use these resources.

## Resource: `google_compute_target_pool`

### `instances` is now a Set

The order of entries in `instances` no longer matters. Any configurations that
interpolate based on an item at a specific index will need to be updated, as items
may have been reordered.

## Resource: `google_compute_url_map`

### `host_rule`, `path_matcher`, and `test` are now authoritative

Terraform will remove values not explicitly set in these fields. Any `host_rule`, `path_matcher`, or `test`
values that were added outside of Terraform should be added to the config.

## Resource: `google_container_analysis_note`

### container analysis resources have been removed from the GA provider

Use the [`google-beta` provider](#google-beta-provider) to use these resources.

## Resource: `google_container_cluster`

### `enable_binary_authorization`, `enable_tpu`, `pod_security_policy_config`, `private_cluster`, and `master_ipv4_cidr_block`, `node_config.taints`, `node_config.workload_metadata_config` have been removed from the GA provider

Use the [`google-beta` provider](#google-beta-provider) to set these fields.

## Resource: `google_container_node_pool`

### `max_pods_per_node`, `node_config.taints`, `node_config.workload_metadata_config` have been removed from the GA provider

Use the [`google-beta` provider](#google-beta-provider) to set these fields.

### `name_prefix` has been removed

Use the `name` field along with the `random` provider instead.

Sample config:

```hcl
variable "machine_type" {}

resource "google_container_cluster" "example" {
  name               = "example-cluster"
  zone               = "us-central1-a"
  initial_node_count = 1

  remove_default_node_pool = true
}

resource "random_id" "np" {
  byte_length = 11
  prefix      = "example-np-"
  keepers = {
    machine_type = "${var.machine_type}"
  }
}

resource "google_container_node_pool" "example" {
  name               = "${random_id.np.dec}"
  zone               = "us-central1-a"
  cluster            = "${google_container_cluster.example.name}"
  node_count         = 1

  node_config {
    machine_type = "${var.machine_type}"
  }

  lifecycle {
    create_before_destroy = true
  }
}
```

The `keepers` parameter in `random_id` takes a map of values that cause the random id to be regenerated.
By tying it to attributes that might change, it makes sure the random id changes too.

To make sure the node pool keeps its old name, figure out what the suffix was by running `terraform show`:

```
google_container_node_pool.example:
  ...
  name = example-np-20180329213336514500000001
```

Determine the base64 encoding of that value by running [this script](https://play.golang.org/p/9KrkDoxRTOw).
Then, import that suffix as the value of `random_id`:

```
terraform import random_id.np example-np-,ELFZ1rbrAThoeQE
```

For more details, see [terraform-provider-google#1054](https://github.com/terraform-providers/terraform-provider-google/issues/1054).

## Resource: `google_endpoints_service`

### `protoc_output` has been removed

Use `protoc_output_base64` instead.

Example previous configuration:

```hcl
resource "google_endpoints_service" "grpc_service" {
  service_name  = "api-name.endpoints.project-id.cloud.goog"
  grpc_config   = "${file("service_spec.yml")}"
  protoc_output = "${file("compiled_descriptor_file.pb")}"
```

Example updated configuration:

```hcl
resource "google_endpoints_service" "grpc_service" {
  service_name         = "api-name.endpoints.project-id.cloud.goog"
  grpc_config          = "${file("service_spec.yml")}"
  protoc_output_base64 = "${base64encode(file("compiled_descriptor_file.pb"))}"
}
```

## Resource: `google_dataproc_cluster`

### `cluster_config.0.delete_autogen_bucket` has been removed

Autogenerated buckets are shared by all clusters in the same region, so deleting
this bucket could adversely harm other dataproc clusters. If you need a bucket
that can be deleted, please create a new one and set the `staging_bucket` field.

### `cluster_config.0.gce_cluster_config.0.tags` is now a Set

The order of entries in `tags` no longer matters. Any configurations that
interpolate based on an item at a specific index will need to be updated, as items
may have been reordered.

## Resource: `google_filestore_instance`

### filestore resources have been removed from the GA provider

Use the [`google-beta` provider](#google-beta-provider) to use these resources.


## Resource: `google_organization_custom_role`

### `deleted` field is now an output-only attribute

Use `terraform destroy`, or remove the resource from your config instead.

## Resource: `google_project`

### `app_engine` has been removed

Use the
[`google_app_engine_application` resource](https://www.terraform.io/docs/providers/google/r/app_engine_application.html) instead.

To avoid errors trying to recreate the resource, import it into your state first by running:

```
terraform import google_app_engine_application.app your-project-id
```

## Resource: `google_project_iam_policy`

### `policy_data` is now authoritative

Terraform will remove values not explicitly set in this field. Any `policy_data`
values that were added outside of Terraform should be added to the config.

### `authoritative`, `restore_policy`, and `disable_project` have been removed

Remove these fields from your config. Ensure that `policy_data` contains all
policy values that exist on the project.

This resource is very dangerous. Consider using `google_project_iam_binding` or
`google_project_iam_member` instead.

## Resource: `google_service_account`

### `policy_data` has been removed

Use one of the other
[service account IAM resources](https://www.terraform.io/docs/providers/google/r/google_service_account_iam.html) instead.

## Resource: `google_sql_database_instance`

### `settings` is now authoritative

Terraform will remove values not explicitly set in this field. Any settings
values that were added outside of Terraform should be added to the config.

## Resource: `google_storage_default_object_acl`

### `role_entity` is now authoritative

Terraform will remove values not explicitly set in this field. Any `role_entity`
values that were added outside of Terraform should be added to the config.

## Resource: `google_storage_object_acl`

### `role_entity` is now authoritative

Terraform will remove values not explicitly set in this field. Any `role_entity`
values that were added outside of Terraform should be added to the config.
For fine-grained management, use `google_storage_object_access_control`.
