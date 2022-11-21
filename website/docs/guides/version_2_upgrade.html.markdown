---
page_title: "Terraform Google Provider 2.0.0 Upgrade Guide"
description: |-
  Terraform Google Provider 2.0.0 Upgrade Guide
---

# Terraform Google Provider 2.0.0 Upgrade Guide

Version `2.0.0` of the Google provider for Terraform is a major release and
includes some changes that you will need to consider when upgrading. This guide
is intended to help with that process and focuses only on the changes necessary
to upgrade from version `1.20.0` to `2.0.0`.

-> The "Google provider" refers to both `google` and `google-beta`; each will
have released `2.0.0` at around the same time, and this guide is for both
variants of the Google provider. See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html)
for details if you're new to using `google-beta`.

Most of the changes outlined in this guide have been previously marked as
deprecated in the Terraform `plan`/`apply` output throughout previous provider
releases, up to and including 1.20.0. These changes, such as deprecation notices,
can always be found in the [CHANGELOG](https://github.com/hashicorp/terraform-provider-google/blob/main/CHANGELOG.md).

## Why version 2.0.0?

We introduced version `2.0.0` of the Google provider in order to split the
provider into 2 distinct variants; `google`, the provider for the generally
available (GA) GCP APIs, and `google-beta`, the provider for Beta GCP APIs.

In addition, we made small breaking changes across the provider to enable import
for older resources, enable some new use cases, align field naming / formats
with your expectations based on other GCP tooling, and to facilitate generating
more resources with [Magic Modules](https://github.com/GoogleCloudPlatform/magic-modules).

While you should see some small changes in your configurations as a result of
these changes, we don't expect you'll need to make any major refactorings. As we
develop the provider, we hope to continue to use Magic Modules to provide a
consistent experience across the provider including features like configurable
timeouts, import, and more.

## I accidentally upgraded to 2.0.0, how do I downgrade to `1.X`?

If you've inadvertently upgraded to `2.0.0`, first see the
[Provider Version Configuration Guide](#provider-version-configuration) to lock
your provider version; if you've constrained the provider to a lower version
such as shown in the previous version example in that guide, Terraform will pull
in a `1.X` series release on `terraform init`.

If you've only ran `terraform init` or `terraform plan`, your state will not
have been modified and downgrading your provider is sufficient.

If you've ran `terraform refresh` or `terraform apply`, Terraform may have made
state changes in the meantime.

* If you're using a *local* state, `terraform refresh` with a downgraded
provider is likely sufficient to revert your state. The Google provider
generally refreshes most state information from the API, and the properties
necessary to do so have been left unchanged.

* If you're using a *remote* state backend

  * That does not support versioning, see the local state instructions above

  * That supports *versioning* such as [Google Cloud Storage](https://www.terraform.io/docs/backends/types/gcs.html)
you can revert the Terraform state file to a previous version by hand. If you do
so and Terraform created resources as part of a `terraform apply`, you'll need
to either `terraform import` them or delete them by hand.
  

## Upgrade Topics

<!-- TOC depthFrom:2 depthTo:2 -->

- [Provider Version Configuration](#provider-version-configuration)
- [`google-beta` provider](#google-beta-provider)
- [Data Sources](#data-sources)
- [Resource: `google_bigquery_dataset`](#resource-google_bigquery_dataset)
- [Resource: `google_bigtable_instance`](#resource-google_bigtable_instance)
- [Resource: `google_binary_authorization_attestor`](#resource-google_binary_authorization_attestor)
- [Resource: `google_binary_authorization_policy`](#resource-google_binary_authorization_policy)
- [Resource: `google_cloudbuild_trigger`](#resource-google_cloudbuild_trigger)
- [Resource: `google_cloudfunctions_function`](#resource-google_cloudfunctions_function)
- [Resource: `google_compute_backend_service`](#resource-google_compute_backend_service)
- [Resource: `google_compute_disk`](#resource-google_compute_disk)
- [Resource: `google_compute_global_forwarding_rule`](#resource-google_compute_global_forwarding_rule)
- [Resource: `google_compute_image`](#resource-google_compute_image)
- [Resource: `google_compute_instance`](#resource-google_compute_instance)
- [Resource: `google_compute_instance_from_template`](#resource-google_compute_instance_from_template)
- [Resource: `google_compute_instance_group_manager`](#resource-google_compute_instance_group_manager)
- [Resource: `google_compute_instance_template`](#resource-google_compute_instance_template)
- [Resource: `google_compute_project_metadata`](#resource-google_compute_project_metadata)
- [Resource: `google_compute_region_instance_group_manager`](#resource-google_compute_region_instance_group_manager)
- [Resource: `google_compute_snapshot`](#resource-google_compute_snapshot)
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
- [Resource: `google_project_custom_role`](#resource-google_project_custom_role)
- [Resource: `google_project_iam_policy`](#resource-google_project_iam_policy)
- [Resource: `google_service_account`](#resource-google_service_account)
- [Resource: `google_sql_database_instance`](#resource-google_sql_database_instance)
- [Resource: `google_storage_default_object_acl`](#resource-google_storage_default_object_acl)
- [Resource: `google_storage_object_acl`](#resource-google_storage_object_acl)
- [Resource: `google_*_iam_binding`](#google_*_iam_binding)

<!-- /TOC -->

## Provider Version Configuration

-> Before upgrading to version 2.0.0, it is recommended to upgrade to the most
recent version of the provider (1.20.0) and ensure that your environment
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

  version = "~> 1.20.0"
}
```

An updated configuration:

```hcl
provider "google" {
  # ... other configuration ...

  version = "~> 2.0.0"
}
```

## The `google-beta` provider

The `google-beta` variant of the Google provider is now necessary to be able to
configure products and features that are in beta. The `google-beta` provider
enables full import support of beta features and gives users who wish to use
only the most stable APIs and features more confidence that they are doing so
by continuing to use the `google` provider, which now exclusively uses generally
available (GA) products and features.

Beta GCP features have no deprecation policy and no SLA, but are otherwise considered to be feature-complete
with only minor outstanding issues after their Alpha period. Beta is when GCP
features are publicly announced, and is when they generally become publicly
available. For more information see [the official documentation on GCP launch stages](https://cloud.google.com/terms/launch-stages).

Because the API for beta features can change before their GA launch, there may
be breaking changes in the `google-beta` provider in minor release versions.
These changes will be announced in the [`google-beta` CHANGELOG](https://github.com/hashicorp/terraform-provider-google-beta/blob/main/CHANGELOG.md).

To have resources at different API versions, set up provider blocks for each version:

```hcl
provider "google" {
  project     = "my-project-id"
  region      = "us-central1"
}

provider "google-beta" {
  project     = "my-project-id"
  region      = "us-central1"
}
```

In each resource, explicitly state which provider that resource should be used
with:

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

See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html)
for more details on how to use `google-beta`.

## Data Sources

See the `Resource` sections in this document for properties that may have been
removed.

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

### `cluster_family` has a diff

If you see

```
-/+ google_bigtable_table.my_table (new resource required)
      id:                              "foo" => <computed> (forces new resource)
      column_family.#:                 "1" => "0" (forces new resource)
      column_family.123456789.family:  "my-family" => ""
      instance_name:                   "bar" => "bar"
      name:                            "foo" => "foo"
      project:                         "my-project" => <computed>
```

Add an appropriate `column_family` block to your config, eg:

```diff
+ column_family {
+   family = "my-family"
+ }
```


## Resource: `google_binary_authorization_attestor`

### binary authorization resources have been removed from the GA provider

Use the [`google-beta` provider](#google-beta-provider) to use these resources.

## Resource: `google_binary_authorization_policy`

### binary authorization resources have been removed from the GA provider

Use the [`google-beta` provider](#google-beta-provider) to use these resources.

## Resource: `google_cloudbuild_trigger`

### `build.step.args` is now a list instead of space separated strings.

Example updated configuration:

```hcl
resource "google_cloudbuild_trigger" "build_trigger" {
  trigger_template {
    branch_name = "main-updated"
    repo_name   = "some-repo-updated"
  }

  build {
    images = ["gcr.io/$PROJECT_ID/$REPO_NAME:$SHORT_SHA"]
    tags   = ["team-a", "service-b", "updated"]

    step {
      name = "gcr.io/cloud-builders/gsutil"
      args = ["cp", "gs://mybucket/remotefile.zip", "localfile-updated.zip"]
    }

    step {
      name = "gcr.io/cloud-builders/go"
      args = ["build", "my_package_updated"]
    }

    step {
      name = "gcr.io/cloud-builders/docker"
      args = ["build", "-t", "gcr.io/$PROJECT_ID/$REPO_NAME:$SHORT_SHA", "-f", "Dockerfile", "."]
    }
    step {
      name = "gcr.io/$PROJECT_ID/$REPO_NAME:$SHORT_SHA"
      args = ["test"]
    }
  }
}
```

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
  name     = "example-bucket"
}

resource "google_storage_bucket_object" "archive" {
  name   = "index.zip"
  bucket = "${google_storage_bucket.bucket.name}"
  source = "path/to/source.zip"
}
```

See the documentation at
[`google_cloudfunctions_function`](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/cloudfunctions_function)
for more details.

## Resource: `google_compute_backend_service`

### `custom_request_headers` has been removed from the GA provider

Use the [`google-beta` provider](#google-beta-provider) to set this field.

### `iap` may cause spurious updates

Due to technical limitations around how Terraform can diff fields, you may see a
spurious update where the client secret in your config replaces an incorrect
value that was recorded in state, the SHA256 hash of the secret's value.

You may also encounter the same behaviour on import.

## Resource: `google_compute_disk`

### `disk_encryption_key_raw` and `disk_encryption_key_sha256` have been removed.

Use the `disk_encryption_key` block instead:

```hcl
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
  name  = "example-disk"
  image = "${data.google_compute_image.my_image.self_link}"
  size  = 50
  type  = "pd-ssd"
  zone  = "us-central1-a"
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

### `network_interface.*.address` has been removed

Use `network_interface.*.network_ip` instead.

## Resource: `google_compute_instance_group_manager`

### `version`, `auto_healing_policies`, `rolling_update_policy` have been removed from the GA provider

Use the [`google-beta` provider](#google-beta-provider) to use these fields.
`rolling_update_policy` has been renamed to `update_policy` in `google-beta`.

## Resource: `google_compute_instance_template`

### `network_interface.*.address` has been removed

Use `network_interface.*.network_ip` instead.

## Resource: `google_compute_project_metadata`

### `metadata` is now authoritative

Terraform will remove values not explicitly set in this field. Any `metadata` values
that were added outside of Terraform should be added to the config.

## Resource: `google_compute_region_instance_group_manager`

### `version`, `auto_healing_policies`, `rolling_update_policy` have been removed from the GA provider

Use the [`google-beta` provider](#google-beta-provider) to use these fields.
`rolling_update_policy` has been renamed to `update_policy` in `google-beta`.

### `update_strategy` no longer has any effect and is removed

With `rolling_update_policy` removed, `update_strategy` has no effect anymore.
Before updating, remove it from your config.

## Resource: `google_compute_snapshot`

### `snapshot_encryption_key_raw` and `snapshot_encryption_key_sha256` have been removed.

Use the `snapshot_encryption_key` block instead:

```hcl
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "my_disk" {
  name  = "my-disk"
  image = "${data.google_compute_image.my_image.self_link}"
  size  = 10
  type  = "pd-ssd"
  zone  = "us-central1-a"
}

resource "google_compute_snapshot" "my_snapshot" {
  name        = "my-snapshot"
  source_disk = "${google_compute_disk.my_disk.name}"
  zone        = "us-central1-a"
  snapshot_encryption_key {
    raw_key = "SGVsbG8gZnJvbSBHb29nbGUgQ2xvdWQgUGxhdGZvcm0="
  }
}
```

### `source_disk_encryption_key_raw` and `source_disk_encryption_key_sha256` have been removed.

Use the `source_disk_encryption_key` block instead:

```hcl
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}
resource "google_compute_disk" "my_disk" {
  name  = "my-disk"
  image = "${data.google_compute_image.my_image.self_link}"
  size  = 10
  type  = "pd-ssd"
  zone  = "us-central1-a"
  disk_encryption_key {
    raw_key = "SGVsbG8gZnJvbSBHb29nbGUgQ2xvdWQgUGxhdGZvcm0="
  }
}
resource "google_compute_snapshot" "my_snapshot" {
  name        = "my-snapshot"
  source_disk = "${google_compute_disk.my_disk.name}"
  zone        = "us-central1-a"
  source_disk_encryption_key {
    raw_key = "SGVsbG8gZnJvbSBHb29nbGUgQ2xvdWQgUGxhdGZvcm0="
  }
}

```

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

### `enable_binary_authorization`, `enable_tpu`, `pod_security_policy_config`, `node_config.taints`, `node_config.workload_metadata_config` have been removed from the GA provider

Use the [`google-beta` provider](#google-beta-provider) to set these fields.

### `private_cluster`, `master_ipv4_cidr_block` are removed.

Use `private_cluster_config` and `private_cluster_config.master_ipv4_cidr_block` instead.

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
  name       = "${random_id.np.dec}"
  zone       = "us-central1-a"
  cluster    = "${google_container_cluster.example.name}"
  node_count = 1

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

For more details, see [terraform-provider-google#1054](https://github.com/hashicorp/terraform-provider-google/issues/1054).

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
[`google_app_engine_application` resource](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/app_engine_application) instead.

To avoid errors trying to recreate the resource, import it into your state first by running:

```
terraform import google_app_engine_application.app your-project-id
```


## Resource: `google_project_custom_role`

### `deleted` field is now an output-only attribute

Use `terraform destroy`, or remove the resource from your config instead.

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
[service account IAM resources](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/google_service_account_iam) instead.

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

## Resource: `google_*_iam_binding`

### Create is now authoritative

Every `iam_binding` resource will overwrite the existing member list for a given
role on Create. Running `terraform plan` for the first time will not show members
that have been added via other tools. *To ensure existing `members` are preserved
use `terraform import` instead of creating the resource.*

Previous versions of `google_*_iam_binding` resources would merge the existing
members of a role with the members defined in the terraform config. If there was
a difference between the members defined in the config and the existing members
defined for an existing role it would show a diff if `terraform plan` was run
immediately after create had succeeded.

Affected resources:
* `google_billing_account_iam_binding`
* `google_folder_iam_binding`
* `google_kms_key_ring_iam_binding`
* `google_kms_crypto_key_iam_binding`
* `google_spanner_instance_iam_binding`
* `google_spanner_database_iam_binding`
* `google_organization_iam_binding`
* `google_project_iam_binding`
* `google_pubsub_topic_iam_binding`
* `google_pubsub_subscription_iam_binding`
* `google_service_account_iam_binding`
