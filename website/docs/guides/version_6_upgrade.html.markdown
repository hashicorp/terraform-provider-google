---
page_title: "Terraform provider for Google Cloud 6.0.0 Upgrade Guide"
description: |-
  Terraform provider for Google Cloud 6.0.0 Upgrade Guide
---

# Terraform Google Provider 6.0.0 Upgrade Guide

The `6.0.0` release of the Google provider for Terraform is a major version and
includes some changes that you will need to consider when upgrading. This guide
is intended to help with that process and focuses only on the changes necessary
to upgrade from the final `5.X` series release to `6.0.0`.

Most of the changes outlined in this guide have been previously marked as
deprecated in the Terraform `plan`/`apply` output throughout previous provider
releases, up to and including the final `5.X` series release. These changes,
such as deprecation notices, can always be found in the CHANGELOG of the
affected providers. [google](https://github.com/hashicorp/terraform-provider-google/blob/main/CHANGELOG.md)
[google-beta](https://github.com/hashicorp/terraform-provider-google-beta/blob/main/CHANGELOG.md)

## I accidentally upgraded to 6.0.0, how do I downgrade to `5.X`?

If you've inadvertently upgraded to `6.0.0`, first see the
[Provider Version Configuration Guide](#provider-version-configuration) to lock
your provider version; if you've constrained the provider to a lower version
such as shown in the previous version example in that guide, Terraform will pull
in a `5.X` series release on `terraform init`.

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
[Google Cloud Storage](https://developer.hashicorp.com/terraform/language/settings/backends/gcs),
you can revert the Terraform state file to a previous version. If you do
so and Terraform had created resources as part of a `terraform apply` in the
meantime, you'll need to either delete them by hand or `terraform import` them
so Terraform knows to manage them.

## Provider Version Configuration

-> Before upgrading to version 6.0.0, it is recommended to upgrade to the most
recent `5.X` series release of the provider, make the changes noted in this guide,
and ensure that your environment successfully runs
[`terraform plan`](https://developer.hashicorp.com/terraform/cli/commands/plan)
without unexpected changes or deprecation notices.

It is recommended to use [version constraints](https://developer.hashicorp.com/terraform/language/providers/requirements#requiring-providers)
when configuring Terraform providers. If you are following that recommendation,
update the version constraints in your Terraform configuration and run
[`terraform init`](https://developer.hashicorp.com/terraform/cli/commands/init) to download
the new version.

If you aren't using version constraints, you can use `terraform init -upgrade`
in order to upgrade your provider to the latest released version.

For example, given this previous configuration:

```hcl
terraform {
  required_providers {
    google = {
      version = "~> 5.30.0"
    }
  }
}
```

An updated configuration:

```hcl
terraform {
  required_providers {
    google = {
      version = "~> 6.0.0"
    }
  }
}
```

## Provider

### Provider-level change example header

Description of the change and how users should adjust their configuration (if needed).

## Datasources

## Datasource: `google_product_datasource`

### Datasource-level change example header

Description of the change and how users should adjust their configuration (if needed).

## Resources

## Resource: `google_bigquery_table`

### View creation now validates `schema`

A `view` can no longer be created when `schema` contains required fields

## Resource: `google_bigquery_reservation`

### `multi_region_auxiliary` is now removed

This field is no longer supported by the BigQuery Reservation API.

## Resource: `google_sql_database_instance`

### `settings.ip_configuration.require_ssl` is now removed

Removed in favor of field `settings.ip_configuration.ssl_mode`.

## Resource: `google_pubsub_topic`

### `schema_settings` no longer has a default value

An empty value means the setting should be cleared.

## Resources: `google_container_cluster`, `google_container_node_pool`, and `google_compute_instance`

### `guest_accelerator = []` is no longer valid configuration

To explicitly set an empty list of objects, set `guest_accelerator.count = 0`.

Previously, to explicitly set `guest_accelerator` as an empty list of objects, the specific configuration `guest_accelerator = []` was necessary.
This was to maintain compatability in behavior between Terraform versions 0.11 and 0.12 using a special setting ["attributes as blocks"](https://developer.hashicorp.com/terraform/language/attr-as-blocks).
This special setting causes other breakages so it is now removed, with setting `guest_accelerator.count = 0` available as an alternative form of empty `guest_accelerator` object.

### `guest_accelerator.gpu_driver_installation_config = []` and `guest_accelerator.gpu_sharing_config = []` are no longer valid configuration

These were never intended to be set this way. Removing the fields from configuration should not produce a diff.

## Resource: `google_domain`

### Domain deletion now prevented by default with `deletion_protection`

The field `deletion_protection` has been added with a default value of `true`. This field prevents
Terraform from destroying or recreating the Domain. In 6.0.0, existing domains will have 
`deletion_protection` set to `true` during the next refresh unless otherwise set in configuration.

**`deletion_protection` does NOT prevent deletion outside of Terraform.**

To disable deletion protection, explicitly set this field to `false` in configuration
and then run `terraform apply` to apply the change.

## Resource: `google_cloud_run_v2_job`

### retyped `containers.env` to SET from ARRAY

Previously, `containers.env` was a list, making it order-dependent. It is now a set.

If you were relying on accessing an individual environment variable by index (for example, `google_cloud_run_v2_job.template.containers.0.env.0.name`), then that will now need to by hash (for example, `google_cloud_run_v2_job.template.containers.0.env.<some-hash>.name`).

## Resource: `google_cloud_run_v2_service`

### `liveness_probe` no longer defaults from API

Cloud Run does not provide a default value for liveness probe. Now removing this field
will remove the liveness probe from the Cloud Run service.

### retyped `containers.env` to SET from ARRAY

Previously, `containers.env` was a list, making it order-dependent. It is now a set.

If you were relying on accessing an individual environment variable by index (for example, `google_cloud_run_v2_service.template.containers.0.env.0.name`), then that will now need to by hash (for example, `google_cloud_run_v2_service.template.containers.0.env.<some-hash>.name`).

## Resource: `google_compute_subnetwork`

### `secondary_ip_range = []` is no longer valid configuration

To explicitly set an empty list of objects, use `send_secondary_ip_range_if_empty = true` and completely remove `secondary_ip_range` from config.

Previously, to explicitly set `secondary_ip_range` as an empty list of objects, the specific configuration `secondary_ip_range = []` was necessary.
This was to maintain compatability in behavior between Terraform versions 0.11 and 0.12 using a special setting ["attributes as blocks"](https://developer.hashicorp.com/terraform/language/attr-as-blocks).
This special setting causes other breakages so it is now removed, with `send_secondary_ip_range_if_empty` available instead.

## Resource: `google_compute_backend_service`

## Resource: `google_compute_region_backend_service`

### `iap.enabled` is now required in the `iap` block

To apply the IAP settings to the backend service, `true` needs to be set for `enabled` field.

### `outlier_detection` subfields default values removed

Empty values mean the setting should be cleared.

### `connection_draining_timeout_sec` default value changed

An empty value now means 300.

### `balancing_mode` default value changed

An empty value now means UTILIZATION.

## Resource: `google_redis_cluster`

### `deletion_protection_enabled` field with default value added

Support for the deletionProtectionEnabled field has been added. Redis clusters will now be created with a `deletion_protection_enabled = true` value by default. 
 
## Resource: `google_vpc_access_connector`

### Fields `min_throughput` and `max_throughput` no longer have default values

The fields `min_throughput` and `max_throughput` no longer have default values 
set by the provider. This was necessary to add conflicting field validation, also
described in this guide.

No configuration changes are needed for existing resources as these fields' values
will default to values present in data returned from the API.

### Conflicting field validation added for `min_throughput` and `min_instances`, and `max_throughput` and `max_instances`

The provider will now enforce that `google_vpc_access_connector` resources can only
include one of `min_throughput` and `min_instances` and one of `max_throughput`and 
`max_instances`. Previously if a user included all four fields in a resource block
they would experience a permadiff. This is a result of how `min_instances` and
`max_instances` fields' values take precedence in the API, and how the API calculates
values for `min_throughput` and `max_throughput` that match the number of instances.

Users will need to check their configuration for any `google_vpc_access_connector`
resource blocks that contain both fields in a conflicting pair, and remove one of those fields.
The fields that are removed from the configuration will still have Computed values,
that are derived from the API.

## Removals

### Resource: `google_identity_platform_project_default_config` is now removed

`google_identity_platform_project_default_config` is removed in favor of `google_identity_platform_project_config`

## Resource: `google_storage_bucket`

### `lifecycle_rule.condition.no_age` is now removed

Previously `lifecycle_rule.condition.age` attribute was being set to zero by default and `lifecycle_rule.condition.no_age` was introduced to prevent that.
Now `lifecycle_rule.condition.no_age` is no longer supported and `lifecycle_rule.condition.age` won't be set to zero by default.
Removed in favor of the field `lifecycle_rule.condition.send_age_if_zero` which can be used to set a zero value for the `lifecycle_rule.condition.age` attribute. 

For a seamless update, if your state today uses `no_age=true`, update it to remove `no_age` and set `send_age_if_zero=false`. If you do not use `no_age=true` and desire to continue creating rules with an `age=0` condition, you will need to add `send_age_if_zero=true` to your state to avoid any changes after updating to 6.0.0. 

With the 6.0.0 update, `send_age_if_zero` will be set to `false` by default unless declared explicitly `true`, and `age=0` conditions will be removed from existing buckets next time your `lifecycle_rule.condition` configuration is updated.
