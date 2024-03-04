---
page_title: "Terraform provider for Google Cloud 4.0.0 Upgrade Guide"
description: |-
  Terraform provider for Google Cloud 4.0.0 Upgrade Guide
---

# Terraform provider for Google Cloud 4.0.0 Upgrade Guide

The `4.0.0` release of the Terraform provider for Google Cloud is a major version and
includes some changes that you will need to consider when upgrading. This guide
is intended to help with that process and focuses only on the changes necessary
to upgrade from the final `3.X` series release to `4.0.0`.

Most of the changes outlined in this guide have been previously marked as
deprecated in the Terraform `plan`/`apply` output throughout previous provider
releases, up to and including the final `3.X` series release. These changes,
such as deprecation notices, can always be found in the CHANGELOG of the
affected providers. [google](https://github.com/hashicorp/terraform-provider-google/blob/main/CHANGELOG.md)
[google-beta](https://github.com/hashicorp/terraform-provider-google-beta/blob/main/CHANGELOG.md)

## I accidentally upgraded to 4.0.0, how do I downgrade to `3.X`?

If you've inadvertently upgraded to `4.0.0`, first see the
[Provider Version Configuration Guide](#provider-version-configuration) to lock
your provider version; if you've constrained the provider to a lower version
such as shown in the previous version example in that guide, Terraform will pull
in a `3.X` series release on `terraform init`.

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

## Provider Version Configuration

-> Before upgrading to version 4.0.0, it is recommended to upgrade to the most
recent `3.X` series release of the provider, make the changes noted in this guide,
and ensure that your environment successfully runs
[`terraform plan`](https://www.terraform.io/docs/commands/plan.html)
without unexpected changes or deprecation notices.

It is recommended to use [version constraints](https://www.terraform.io/docs/language/providers/requirements.html#requiring-providers)
when configuring Terraform providers. If you are following that recommendation,
update the version constraints in your Terraform configuration and run
[`terraform init`](https://www.terraform.io/docs/commands/init.html) to download
the new version.

If you aren't using version constraints, you can use `terraform init -upgrade`
in order to upgrade your provider to the latest released version.

For example, given this previous configuration:

```hcl
terraform {
  # ... other configuration ...
  required_providers {
    google = {
      version = "~> 3.90.0"
    }
  }
}
```

An updated configuration:

```hcl
terraform {
  # ... other configuration ...
  required_providers {
    google = {
      version = "~> 4.0.0"
    }
  }
}
```

## Provider

### `credentials`, `access_token` precedence has changed

Terraform can draw values for both the `credentials` and `access_token` from the
config directly or from environment variables. 

In earlier versions of the provider, `access_token` values specified through
environment variables took precedence over `credentials` values specified in
config. From `4.0.0` onwards, config takes precedence over environment variables,
and the `access_token` environment variable takes precedence over the
`credential` environment variable.

Service account impersonation is unchanged. Terraform will continue to use
the service account if it is specified through an environment variable, even
if `credentials` or `access_token` are specified in config.

### Redundant default scopes are removed

Several default scopes are removed from the provider:

* "https://www.googleapis.com/auth/compute"
* "https://www.googleapis.com/auth/ndev.clouddns.readwrite"
* "https://www.googleapis.com/auth/devstorage.full_control"
* "https://www.googleapis.com/auth/cloud-identity"

They are redundant with the "https://www.googleapis.com/auth/cloud-platform"
scope per [Access scopes](https://cloud.google.com/compute/docs/access/service-accounts#accesscopesiam).
After this change the following scopes are enabled, in line with `gcloud`'s
[list of scopes](https://cloud.google.com/sdk/gcloud/reference/auth/application-default/login):

* "https://www.googleapis.com/auth/cloud-platform"
* "https://www.googleapis.com/auth/userinfo.email"

This change is believed to have no user impact. If you find that Terraform
behaves incorrectly as a result of this change, please report a [bug](https://github.com/hashicorp/terraform-provider-google/issues/new?assignees=&labels=bug&template=bug.md).

### Runtime Configurator (`runtimeconfig`) resources have been removed from the GA provider

Earlier versions of the provider accidentally included the Runtime Configurator
service at GA. `4.0.0` has corrected that error, and Runtime Configurator is
only available in `google-beta`.

Affected Resources:

    * `google_runtimeconfig_config`
    * `google_runtimeconfig_variable`
    * `google_runtimeconfig_config_iam_policy`
    * `google_runtimeconfig_config_iam_binding`
    * `google_runtimeconfig_config_iam_member`

Affected Datasources:

    * `google_runtimeconfig_config`


If you have a configuration using the `google` provider like the following:

```
resource "google_runtimeconfig_config" "my-runtime-config" {
  name        = "my-service-runtime-config"
  description = "Runtime configuration values for my service"
}
```

Add the `google-beta` provider to your configuration:

```
resource "google_runtimeconfig_config" "my-runtime-config" {
  provider = google-beta

  name        = "my-service-runtime-config"
  description = "Runtime configuration values for my service"
}
```

### Service account scopes no longer accept `trace-append` or `trace-ro`, use `trace` instead

Previously users could specify `trace-append` or `trace-ro` as scopes for a given service account.
However, to better align with [Google documentation](https://cloud.google.com/sdk/gcloud/reference/alpha/compute/instances/set-scopes#--scopes), `trace` will now be the only valid scope, as it's an alias for `trace.append` and
`trace-ro` is no longer a documented option.

## Datasources

## Datasource: `google_kms_key_ring`

### `id` now matches the `google_kms_key_ring` id format

The format has changed to better match the resource's ID format.

Interpolations based on the `id` of the datasource may require updates.

## Resources

## Resource: `google_app_engine_standard_app_version`

### `entrypoint` is now required

This resource would fail to deploy without this field defined. Specify the
`entrypoint` block to fix any issues

## Resource: `google_bigquery_job`

### Exactly one of `query`, `load`, `copy` or `extract` is required
The provider will now enforce at plan time that one of these fields be set.

### At least one of `query.0.script_options.0.statement_timeout_ms`, `query.0.script_options.0.statement_byte_budget`, or `query.0.script_options.0.key_result_statement` is required
The provider will now enforce at plan time that one of these fields be set.

### Exactly one of `extract.0.source_table` or `extract.0.source_model` is required
The provider will now enforce at plan time that one of these fields be set.

## Resource: `google_cloudbuild_trigger`

### Exactly one of `build.0.source.0.repo_source.0.branch_name`, `build.0.source.0.repo_source.0.commit_sha` or `build.0.source.0.repo_source.0.tag_name` is required
The provider will now enforce at plan time that one of these fields be set.

## Resource: `google_compute_autoscaler`

### At least one of `autoscaling_policy.0.scale_down_control.0.max_scaled_down_replicas` or `autoscaling_policy.0.scale_down_control.0.time_window_sec` is required
The provider will now enforce at plan time that one of these fields be set.

### At least one of `autoscaling_policy.0.scale_down_control.0.max_scaled_down_replicas.0.fixed` or `autoscaling_policy.0.scale_down_control.0.max_scaled_down_replicas.0.percent` is required
The provider will now enforce at plan time that one of these fields be set.

### At least one of `autoscaling_policy.0.scale_in_control.0.max_scaled_in_replicas` or `autoscaling_policy.0.scale_in_control.0.time_window_sec` is required
The provider will now enforce at plan time that one of these fields be set.

### At least one of `autoscaling_policy.0.scale_in_control.0.max_scaled_in_replicas.0.fixed` or `autoscaling_policy.0.scale_in_control.0.max_scaled_in_replicas.0.percent` is required
The provider will now enforce at plan time that one of these fields be set.

## Resource: `google_compute_region_autoscaler`

### At least one of `autoscaling_policy.0.scale_down_control.0.max_scaled_down_replicas` or `autoscaling_policy.0.scale_down_control.0.time_window_sec` is required
The provider will now enforce at plan time that one of these fields be set.

### At least one of `autoscaling_policy.0.scale_down_control.0.max_scaled_down_replicas.0.fixed` or `autoscaling_policy.0.scale_down_control.0.max_scaled_down_replicas.0.percent` is required
The provider will now enforce at plan time that one of these fields be set.

### At least one of `autoscaling_policy.0.scale_in_control.0.max_scaled_in_replicas` or `autoscaling_policy.0.scale_in_control.0.time_window_sec` is required
The provider will now enforce at plan time that one of these fields be set.

### At least one of `autoscaling_policy.0.scale_in_control.0.max_scaled_in_replicas.0.fixed` or `autoscaling_policy.0.scale_in_control.0.max_scaled_in_replicas.0.percent` is required
The provider will now enforce at plan time that one of these fields be set.

## Resource: `google_compute_firewall`

### One of `source_tags`, `source_ranges` or `source_service_accounts` are required on INGRESS firewalls

Previously, if all of these fields were left empty, the firewall defaulted to allowing traffic from 0.0.0.0/0, which is a suboptimal default.

### `source_ranges` will track changes when unspecified in a config

In `3.X`, `source_ranges` wouldn't cause a diff if it was undefined in
config but was set on the firewall itself. With 4.0.0 Terraform will now
track changes on the block when it is not specified in a user's config.

## Resource: `google_compute_instance`

### `metadata_startup_script` is no longer set on import

Earlier versions of the provider set the `metadata_startup_script` value on
import, omitting the value of `metadata.startup-script` for historical backwards
compatibility. This was dangerous in practice, as `metadata_startup_script`
would flag an instance for recreation if the values differed rather than for
just an update.

In `4.0.0` the behaviour has been flipped, and `metadata.startup-script` is the
default value that gets written. Users who want `metadata_startup_script` set
on an imported instance will need to modify their state manually. This is more
consistent with our expectations for the field, that a user who manages an
instance **only** through Terraform uses it but that most users should prefer
the `metadata` block.

No action is required for user configs with instances already imported. If you
have a config or module where neither is specified- where `import` will be run,
or an old config that is not reconciled with the API- the value that gets set
will change.

## Resource: `google_compute_instance_group_manager`

### `update_policy.min_ready_sec` is removed from the GA provider
This field was incorrectly included in the GA `google` provider in past releases.
In order to continue to use the feature, add `provider = google-beta` to your
resource definition.

## Resource: `google_compute_region_instance_group_manager`

### `update_policy.min_ready_sec` is removed from the GA provider

This field was incorrectly included in the GA `google` provider in past releases.
In order to continue to use the feature, add `provider = google-beta` to your
resource definition.

## Resource: `google_compute_instance_template`

### `enable_display` is removed from the GA provider

This field was incorrectly included in the GA `google` provider in past releases.
In order to continue to use the feature, add `provider = google-beta` to your
resource definition.

### `advanced_machine_features` will track changes when unspecified in a config

In `3.X`, `advanced_machine_features` wouldn't cause a diff if it was undefined in
config but was set on the instance template itself. With 4.0.0 Terraform will now
track changes on the block when it is not specified in a user's config.

## Resource: `google_compute_url_map`

### At least one of `default_route_action.0.fault_injection_policy.0.delay.0.fixed_delay` or `default_route_action.0.fault_injection_policy.0.delay.0.percentage` is required
The provider will now enforce at plan time that one of these fields be set.

## Resource: `google_container_cluster`

### `enable_shielded_nodes` now defaults to `true`

Previously the provider defaulted `enable_shielded_nodes` to false, despite the API default of `true`.
Unless explicitly configured, users may see a diff changing `enable_shielded_nodes` to `true`.

### `instance_group_urls` is now removed

`instance_group_urls` has been removed in favor of `node_pool.managed_instance_group_urls`

### `master_auth.username` and `master_auth.password` are now removed

`master_auth.username` and `master_auth.password` have been removed. 
Basic authentication was removed for GKE cluster versions >= 1.19. The cluster cannot be created with basic authentication enabled. Instructions for choosing an alternative authentication method can be found at: cloud.google.com/kubernetes-engine/docs/how-to/api-server-authentication.

### `master_auth.client_certificate_config` is now required

With the removal of `master_auth.username` and `master_auth.password`, `master_auth.client_certificate_config` is now
the only configurable field in `master_auth`. If you do not wish to configure `master_auth.client_certificate_config`, 
remove the `master_auth` block from your configuration entirely. You will still be able to reference the outputted fields under `master_auth` without the block defined.

### `node_config.workload_metadata_config.node_metadata` is now removed

Removed in favor of `node_config.workload_metadata_config.mode`.

### `workload_identity_config.0.identity_namespace` is now removed

Removed in favor of `workload_identity_config.0.workload_pool`. Switching your
configuration from one value to the other will trigger a diff at plan time, and
a spurious update.

```diff
resource "google_container_cluster" "cluster" {
  name               = "your-cluster"
  location           = "us-central1-a"
  initial_node_count = 1

  workload_identity_config {
-    identity_namespace = "your-project.svc.id.goog"
+   workload_pool = "your-project.svc.id.goog"
  }
```

### `pod_security_policy_config` is removed from the GA provider

This field was incorrectly included in the GA `google` provider in past releases.
In order to continue to use the feature, add `provider = google-beta` to your
resource definition.

## Resource: `google_compute_snapshot`

### `source_disk_link` is now removed

Removed, as the information available was redundant. You can reconstruct a
compatible value based on `source_disk` and `zone`. With a reference such as the
following:

```
google_compute_snapshot.my_snapshot.source_disk_link
```

Substitute the following:

```
"projects/${google_compute_snapshot.my_snapshot.project}/zones/${google_compute_snapshot.my_snapshot.zone}/disks/${google_compute_snapshot.my_snapshot.source_disk}"
```

## Resource: `google_data_loss_prevention_trigger`

### Exactly one of `inspect_job.0.storage_config.0.cloud_storage_options.0.file_set.0.url` or `inspect_job.0.storage_config.0.cloud_storage_options.0.file_set.0.regex_file_set` is required
The provider will now enforce at plan time that one of these fields be set.

### At least one of `inspect_job.0.storage_config.0.timespan_config.0.start_time` or `inspect_job.0.storage_config.0.timespan_config.0.end_time` is required
The provider will now enforce at plan time that one of these fields be set.

## Resource: `google_os_config_patch_deployment`

### At least one of `patch_config.0.reboot_config`, `patch_config.0.apt`, `patch_config.0.yum`, `patch_config.0.goo` `patch_config.0.zypper`, `patch_config.0.windows_update`, `patch_config.0.pre_step` or `patch_config.0.pre_step` is required
The provider will now enforce at plan time that one of these fields be set.

### At least one of `patch_config.0.apt.0.type`, `patch_config.0.apt.0.excludes` or `patch_config.0.apt.0.exclusive_packages` is required
The provider will now enforce at plan time that one of these fields be set.

### At least one of `patch_config.0.yum.0.security`, `patch_config.0.yum.0.minimal`, `patch_config.0.yum.0.excludes` or `patch_config.0.yum.0.exclusive_packages` is required
The provider will now enforce at plan time that one of these fields be set.

### At least one of `patch_config.0.zypper.0.with_optional`, `patch_config.0.zypper.0.with_update`, `patch_config.0.zypper.0.categories`, `patch_config.0.zypper.0.severities`, `patch_config.0.zypper.0.excludes` or `patch_config.0.zypper.0.exclusive_patches` is required
The provider will now enforce at plan time that one of these fields be set.

### Exactly one of `patch_config.0.windows_update.0.classifications`, `patch_config.0.windows_update.0.excludes` or `patch_config.0.windows_update.0.exclusive_patches` is required
The provider will now enforce at plan time that one of these fields be set.

### At least one of `patch_config.0.pre_step.0.linux_exec_step_config` or `patch_config.0.pre_step.0.windows_exec_step_config` is required
The provider will now enforce at plan time that one of these fields be set.

### At least one of `patch_config.0.post_step.0.linux_exec_step_config` or `patch_config.0.post_step.0.windows_exec_step_config` is required
The provider will now enforce at plan time that one of these fields be set.

## Resource: `google_kms_crypto_key`

### `self_link` is now removed

Removed in favor of `id`.

## Resource: `google_kms_key_ring`

### `self_link` is now removed

Removed in favor of `id`.

## Resource: `google_project`

### `org_id`, `folder_id` now conflict at plan time

Previously, they were only checked for conflicts at apply time. Terraform will
now report an error at plan time.

### `org_id`, `folder_id` are unset when removed from config

Previously, these fields kept their old value in state when they were removed
from config, changing the value on next refresh. Going forward, removing one of
the values or switching values will generate a correct plan that removes the
value.

## Resource: `google_project_iam`

### `project` field is now required

The `project` field is now required for all `google_project_iam_*` resources.
Previously, it was only required for `google_project_iam_policy`. This will make
configuration of the project IAM resources more explicit, given that the project
is the targeted resource.

`terraform plan` will indicate any project IAM resources that had drawn a value
with a provider, and you are able to specify the project explicitly to remove
the proposed diff.

## Resource: `google_project_service`

### `bigquery-json.googleapis.com` is no longer a valid service name

`bigquery-json.googleapis.com` was deprecated in the `3.0.0` release, however, at that point the provider
converted it while the upstream API migration was in progress. Now that the API migration has finished,
the provider will no longer convert the service name. Use `bigquery.googleapis.com` instead.

## Resource: `google_pubsub_subscription`

### `path` is now removed

`path` has been removed in favor of `id` which has an identical value.

## Resource: `google_spanner_instance`

### Exactly one of `num_nodes` or `processing_units` is required

The provider will now enforce that you've set one of these fields at plan time.
Earlier versions of the provider set a default value of `1` for `num_nodes`. If
neither field is present in your config, it's likely you can add `num_nodes = 1`
to resolve this change. If that is incorrect, `terraform plan` should inform you
of the correct value.


For example, for a configuration like the following:

```tf
resource "google_spanner_instance" "default" {
  display_name = "main-instance"
  config       = "regional-europe-west1"
}
```

You would amend it to:

```tf
resource "google_spanner_instance" "default" {
  display_name = "main-instance"
  config       = "regional-europe-west1"
  num_nodes    = 1
}
```

## Resource: `google_sql_database_instance`

### First-generation fields have been removed

Removed fields specific to first-generation SQL instances:
`authorized_gae_applications`, `crash_safe_replication`, `replication_type`

### `database_version` field is now required

The `database_version` field is now required.
Previously, it was an optional field and the default value was `MYSQL_5_6`.
Description of the change and how users should adjust their configuration (if needed).

### Drift detection and defaults enabled on fields

Added drift detection and plan-time defaults to several fields used to configure
second-generation SQL instances. If you see changes flagged by Terraform after
running `terraform plan`, amend your config to resolve them.

The affected fields are:

  * `activation_policy` will now default to `ALWAYS` at plan time, and detect
drift even when unset. Previously, Terraform only detected drift when the field
had been set in config explicitly.

  * `availability_type` will now default to `ZONAL` at plan time, and detect
drift even when unset. Previously, Terraform only detected drift when the field
had been set in config explicitly.

  * `disk_type` will now default to `PD_SSD` at plan time, and detect
drift even when unset. Previously, Terraform only detected drift when the field
had been set in config explicitly.

  * `encryption_key_name` will now detect drift even when unset. Previously,
Terraform only detected drift when the field had been set in config explicitly.

## Resource: `google_storage_bucket`

### `bucket_policy_only` field is now removed

`bucket_policy_only` field is now removed in favor of `uniform_bucket_level_access`.

### `location` field is now required.

Previously, the default value of `location` was `US`. In an attempt to avoid allowing invalid 
conbination of `storageClass` value and default `location` value, `location` field is now required.
