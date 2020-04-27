---
layout: "google"
page_title: "Terraform Google Provider 3.0.0 Upgrade Guide"
sidebar_current: "docs-google-provider-guides-version-3-upgrade"
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

## What is `3.0.0-beta.1`?

With `3.0.0`, we introduced a prerelease window for our major provider releases.
`3.0.0-beta.1` contains all of the changes in `3.0.0`, and allows you to test it
prior to the full upgrade. Currently `3.0.0` is not expected to contain new
features not available in `3.0.0-beta.1`, only bugfixes for issues we're made
aware of before `3.0.0`'s release. Using `3.0.0-beta.1` in production is not
recommended.

```hcl
provider "google" {
  # ... other configuration ...

  version = "~> 3.0.0-beta.1"
}
```

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
- [Provider](#provider)
- [ID Format Changes](#id-format-changes)
- [Data Source: `google_container_engine_versions`](#data-source-google_container_engine_versions)
- [Resource: `google_access_context_manager_access_level`](#resource-google_access_context_manager_access_level)
- [Resource: `google_access_context_manager_service_perimeter`](#resource-google_access_context_manager_service_perimeter)
- [Resource: `google_app_engine_application`](#resource-google_app_engine_application)
- [Resource: `google_app_engine_domain_mapping`](#resource-google_app_engine_domain_mapping)
- [Resource: `google_app_engine_standard_app_version`](#resource-google_app_engine_standard_app_version)
- [Resource: `google_bigquery_table`](#resource-google_bigquery_table)
- [Resource: `google_bigtable_app_profile`](#resource-google_bigtable_app_profile)
- [Resource: `google_binary_authorization_policy`](#resource-google_binary_authorization_policy)
- [Resource: `google_cloudbuild_trigger`](#resource-google_cloudbuild_trigger)
- [Resource: `google_cloudfunctions_function`](#resource-google_cloudfunctions_function)
- [Resource: `google_cloudiot_registry`](#resource-google_cloudiot_registry)
- [Resource: `google_cloudscheduler_job`](#resource-google_cloudscheduler_job)
- [Resource: `google_cloud_run_service`](#resource-google_cloud_run_service)
- [Resource: `google_composer_environment`](#resource-google_composer_environment)
- [Resource: `google_compute_backend_bucket`](#resource-google_compute_backend_bucket)
- [Resource: `google_compute_backend_service`](#resource-google_compute_backend_service)
- [Resource: `google_compute_firewall`](#resource-google_compute_firewall)
- [Resource: `google_compute_forwarding_rule`](#resource-google_compute_forwarding_rule)
- [Resource: `google_compute_global_forwarding_rule`](#resource-google_compute_global_forwarding_rule)
- [Resource: `google_compute_health_check`](#resource-google_compute_health_check)
- [Resource: `google_compute_image`](#resource-google_compute_image)
- [Resource: `google_compute_instance`](#resource-google_compute_instance)
- [Resource: `google_compute_instance_group_manager`](#resource-google_compute_instance_group_manager)
- [Resource: `google_compute_instance_template`](#resource-google_compute_instance_template)
- [Resource: `google_compute_network`](#resource-google_compute_network)
- [Resource: `google_compute_network_peering`](#resource-google_compute_network_peering)
- [Resource: `google_compute_node_template`](#resource-google_compute_node_template)
- [Resource: `google_compute_region_backend_service`](#resource-google_compute_region_backend_service)
- [Resource: `google_compute_region_health_check`](#resource-google_compute_region_health_check)
- [Resource: `google_compute_region_instance_group_manager`](#resource-google_compute_instance_group_manager)
- [Resource: `google_compute_resource_policy`](#resource-google_compute_resource_policy)
- [Resource: `google_compute_route`](#resource-google_compute_route)
- [Resource: `google_compute_router`](#resource-google_compute_router)
- [Resource: `google_compute_router_peer`](#resource-google_compute_router_peer)
- [Resource: `google_compute_snapshot`](#resource-google_compute_snapshot)
- [Resource: `google_compute_subnetwork`](#resource-google_compute_subnetwork)
- [Resource: `google_container_cluster`](#resource-google_container_cluster)
- [Resource: `google_container_node_pool`](#resource-google_container_node_pool)
- [Resource: `google_dataproc_autoscaling_policy`](#resource-google_dataproc_autoscaling_policy)
- [Resource: `google_dataproc_cluster`](#resource-google_dataproc_cluster)
- [Resource: `google_dataproc_job`](#resource-google_dataproc_job)
- [Resource: `google_dns_managed_zone`](#resource-google_dns_managed_zone)
- [Resource: `google_dns_policy`](#resource-google_dns_policy)
- [Resource: `google_folder_organization_policy`](#resource-google_folder_organization_policy)
- [Resource: `google_healthcare_hl7_v2_store`](#resource-google_healthcare_hl7_v2_store)
- [Resource: `google_logging_metric`](#resource-google_logging_metric)
- [Resource: `google_mlengine_model`](#resource-google_mlengine_model)
- [Resource: `google_monitoring_alert_policy`](#resource-google_monitoring_alert_policy)
- [Resource: `google_monitoring_uptime_check_config`](#resource-google_monitoring_uptime_check_config)
- [Resource: `google_organization_policy`](#resource-google_organization_policy)
- [Resource: `google_project_iam_audit_config`](#resource-google_project_iam_audit_config)
- [Resource: `google_project_organization_policy`](#resource-google_project_organization_policy)
- [Resource: `google_project_service`](#resource-google_project_service)
- [Resource: `google_project_services`](#resource-google_project_services)
- [Resource: `google_pubsub_subscription`](#resource-google_pubsub_subscription)
- [Resource: `google_security_scanner_scan_config`](#resource-google_security_scanner_scan_config)
- [Resource: `google_service_account_key`](#resource-google_service_account_key)
- [Resource: `google_sql_database_instance`](#resource-google_sql_database_instance)
- [Resource: `google_storage_bucket`](#resource-google_storage_bucket)
- [Resource: `google_storage_transfer_job`](#resource-google_storage_transfer_job)
- [Resource: `google_tpu_node`](#resource-google_tpu_node)

<!-- /TOC -->

## Provider Version Configuration

-> Before upgrading to version 3.0.0, it is recommended to upgrade to the most
recent `2.X` series release of the provider, make the changes noted in this guide,
and ensure that your environment successfully runs
[`terraform plan`](https://www.terraform.io/docs/commands/plan.html)
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

## Provider

### Terraform 0.11 no longer supported

Support for Terraform 0.11 has been deprecated, and Terraform 0.12 or higher is
required to `terraform init` the provider. See [the blog post](https://www.hashicorp.com/blog/deprecating-terraform-0-11-support-in-terraform-providers/)
for more information. It is recommended that you upgrade to Terraform 0.12 before
upgrading to version 3.0.0 of the provider.

### `userinfo.email` added to default scopes

`userinfo.email` has been added to the default set of OAuth scopes in the
provider. This provides the Terraform user specified by `credentials`' (generally
a service account) email address to GCP APIs in addition to an obfuscated user
id; particularly, it makes the email of the Terraform user available for some
Kubernetes and IAP use cases.

If this was previously defined explicitly, the definition can now be removed.

#### Old Config

```hcl
provider "google" {
  scopes = [
    "https://www.googleapis.com/auth/compute",
    "https://www.googleapis.com/auth/cloud-platform",
    "https://www.googleapis.com/auth/ndev.clouddns.readwrite",
    "https://www.googleapis.com/auth/devstorage.full_control",
    "https://www.googleapis.com/auth/userinfo.email",
  ]
}
```

#### New Config

```hcl
provider "google" {}
```

## ID Format Changes

ID formats on many resources have changed. ID formats have standardized on being similar to the `self_link` of
a resource. Users who depended on particular ID formats in previous versions may be impacted.

## Data Source: `google_container_engine_versions`

### `region` and `zone` are now removed

Use `location` instead.

## Resource: `google_access_context_manager_access_level`

### `os_type` is now required on block `google_access_context_manager_access_level.basic.conditions.device_policy.os_constraints`

In an attempt to avoid allowing empty blocks in config files, `os_type` is now
required on the `basic.conditions.device_policy.os_constraints` block.

## Resource: `google_access_context_manager_service_perimeter`

### At least one of `resources`, `access_levels`, or `restricted_services` is now required on `google_accesscontextmanager_service_perimeter.status`

In an attempt to avoid allowing empty blocks in config files, at least one of `resources`, `access_levels`,
or `restricted_services` is now required on the `status` block.

## Resource: `google_app_engine_application`

### `split_health_checks` is now required on block `google_app_engine_application.feature_settings`

In an attempt to avoid allowing empty blocks in config files, `split_health_checks` is now
required on the `feature_settings` block.

## Resource: `google_app_engine_domain_mapping`

### `ssl_management_type` is now required on `google_app_engine_domain_mapping.ssl_settings`

In an attempt to avoid allowing empty blocks in config files, `ssl_management_type` is now
required on the `ssl_settings` block.

## Resource: `google_app_engine_standard_app_version`

### At least one of `zip` or `files` is now required on `google_app_engine_standard_app_version.deployment`

In an attempt to avoid allowing empty blocks in config files, at least one of `zip` or `files`
is now required on the `deployment` block.

### `shell` is now required on `google_app_engine_standard_app_version.entrypoint`

In an attempt to avoid allowing empty blocks in config files, `shell` is now
required on the `entrypoint` block.

### `script_path` is now required on `google_app_engine_standard_app_version.handlers.script`

In an attempt to avoid allowing empty blocks in config files, `script_path` is now
required on the `handlers.script` block.

### `source_url` is now required on `google_app_engine_standard_app_version.deployment.files` and `google_app_engine_standard_app_version.deployment.zip`

In an attempt to avoid allowing empty blocks in config files, `shell` is now
required on the `deployment.files` and `deployment.zip` blocks.

## Resource: `google_bigquery_table`

### At least one of `range` or `skip_leading_rows` is now required on `external_data_configuration.google_sheets_options`

In an attempt to avoid allowing empty blocks in config files, at least one
of `range` or `skip_leading_rows` is now required on the
`external_data_configuration.google_sheets_options` block.

## Resource: `google_bigtable_app_profile`

### Exactly one of `single_cluster_routing` or `multi_cluster_routing_use_any` is now required on `google_bigtable_app_profile`

In attempt to be more consistent with the API, exactly one of `single_cluster_routing` or
`multi_cluster_routing_use_any` is now required on `google_bigtable_app_profile`.

### `cluster_id` is now required on `google_bigtable_app_profile.single_cluster_routing`

In an attempt to avoid allowing empty blocks in config files, `cluster_id` is now
required on the `single_cluster_routing` block.

## Resource: `google_binary_authorization_policy`

### `name_pattern` is now required on `google_binary_authorization_policy.admission_whitelist_patterns`

In an attempt to avoid allowing empty blocks in config files, `name_pattern` is now
required on the `admission_whitelist_patterns` block.

### `evaluation_mode` and `enforcement_mode` are now required on `google_binary_authorization_policy.cluster_admission_rules`

In an attempt to avoid allowing empty blocks in config files, `evaluation_mode` and `enforcement_mode` are now
required on the `cluster_admission_rules` block.

## Resource: `google_cloudbuild_trigger`

### Exactly one of `filename` or `build` is now required on `google_cloudbuild_trigger`

In attempt to be more consistent with the API, exactly one of `filename` or `build` is now
required on `google_cloudbuild_trigger`.

### Exactly one of `branch_name`, `tag_name` or `commit_sha` is now required on `google_cloudbuild_trigger.trigger_template`

In an attempt to avoid allowing empty blocks in config files, exactly one
of `branch_name`, `tag_name` or `commit_sha` is now required on the
`trigger_template` block.

### Exactly one of `pull_request` or `push` is now required on `google_cloudbuild_trigger.github`

In an attempt to avoid allowing empty blocks in config files, exactly one
of `pull_request` or `push` is now required on the `github` block.

### Exactly one of `branch` or `tag_name` is now required on `google_cloudbuild_trigger.github.push`

In an attempt to avoid allowing empty blocks in config files, exactly one
of `branch` or `tag_name` is now required on the `github.push` block.

### `steps` is now required on `google_cloudbuild_trigger.build`.

In an attempt to avoid allowing empty blocks in config files, `steps` is now
required on the `build` block.

### `name` is now required on `google_cloudbuild_trigger.build.steps`

In an attempt to avoid allowing empty blocks in config files, `name` is now
required on the `build.steps` block.

### `name` and `path` are now required on `google_cloudbuild_trigger.build.steps.volumes`

In an attempt to avoid allowing empty blocks in config files, `name` and `path` are now
required on the `build.volumes` block.

## Resource: `google_cloudfunctions_function`

### The `runtime` option `nodejs6` has been deprecated

`nodejs6` has been deprecated and is no longer the default value for `runtime`.
`runtime` is now required.

## Resource: `google_cloudiot_registry`

### Replace singular event notification config field with plural `event_notification_configs`

Use the plural field `event_notification_configs` instead of
`event_notification_config`, which has now been removed.
Since the Cloud IoT API now accept multiple event notification configs for a
registry, the singular field no longer exists on the API resource and has been
removed from Terraform to prevent conflicts.


#### Old Config

```hcl
resource "google_cloudiot_registry" "myregistry" {
  name = "%s"

  event_notification_config {
    pubsub_topic_name = google_pubsub_topic.event-topic.id
  }
}

```

#### New Config

```hcl
resource "google_cloudiot_registry" "myregistry" {
  name = "%s"

  event_notification_configs {
    pubsub_topic_name = google_pubsub_topic.event-topic.id
  }
}
```

### `public_key_certificate` is now required on block `google_cloudiot_registry.credentials`

In an attempt to avoid allowing empty blocks in config files, `public_key_certificate` is now
required on the `credentials` block.

## Resource: `google_cloud_run_service`

Google Cloud Run Service is being released at v1 and there are breaking schema changes that have arisen from changing the underlying API. These breaking changes only affect the Beta version of the resource as it was not previously available in the GA provider.

To support partial rollouts of different revisions, the `spec` block is now nested under `template` and a second `metadata` block has been added alongside `spec`. Now users can make a change and, using a named revision, they can control the rollout of that revision with a higher granularity.

#### Old Config

```hcl
resource "google_cloud_run_service" "default" {
  spec {
    containers {
      image = "gcr.io/cloudrun/hello"
      args  = ["arrg2", "pirate"]
    }
    container_concurrency = 10
  }
}
```

#### New Config

```hcl
resource "google_cloud_run_service" "default" {
  template {
    spec {
      containers {
        image = "gcr.io/cloudrun/hello"
        args  = ["arrg2", "pirate"]
      }
      container_concurrency = 10
    }

    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale"      = "1000"
        "run.googleapis.com/client-name"        = "cloud-console"
      }
      name = "revision-name"
    }
  }
}
```

## Resource: `google_cloudscheduler_job`

### Exactly one of `pubsub_target`, `http_target` or `app_engine_http_target` is required on `google_cloudscheduler_job`

In attempt to be more consistent with the API, exactly one of `pubsub_target`, `http_target`
or `app_engine_http_target` is now required on `google_cloudscheduler_job`.

### `service_account_email` is now required on `google_cloudscheduler_job.http_target.oauth_token` and `google_cloudscheduler_job.http_target.oidc_token`.

In an attempt to avoid allowing empty blocks in config files, `service_account_email` is now
required on the `http_target.oauth_token` and `http_target.oidc_token` blocks.

### At least one of `retry_count`, `max_retry_duration`, `min_backoff_duration`, `max_backoff_duration`, or `max_doublings` is now required on `google_cloud_scheduler_job.retry_config`

In an attempt to avoid allowing empty blocks in config files, at least one of `retry_count`,
`max_retry_duration`, `min_backoff_duration`, `max_backoff_duration`, or `max_doublings` is
now required on the `retry_config` block.

### At least one of `service`, `version`, or `instance` is now required on `google_cloud_scheduler_job.app_engine_http_target.app_engine_routing`

In an attempt to avoid allowing empty blocks in config files, at least one of `service`,
`version`, or `instance` is now required on the `app_engine_http_target.app_engine_routing` block.

## Resource: `google_composer_environment`

### At least one of `airflow_config_overrides`, `pypi_packages`, `env_variables`, `image_version`, or `python_version` is now required on `google_composer_environment.config.software_config`

In an attempt to avoid allowing empty blocks in config files, at least one of `airflow_config_overrides`,
`pypi_packages`, `env_variables`, `image_version`, or `python_version` is now required on the
`config.software_config` block.

### `use_ip_aliases` is now required on block `google_composer_environment.ip_allocation_policy`

Previously the default value of `use_ip_aliases` was `true`. In an attempt to avoid allowing empty blocks
in config files, `use_ip_aliases` is now required on the `ip_allocation_policy` block.

### At least one of `enable_private_endpoint` or `master_ipv4_cidr_block` is now required on `google_composer_environment.config.private_environment_config`

In an attempt to avoid allowing empty blocks in config files, at least one of `enable_private_endpoint` or `master_ipv4_cidr_block` is now required on the
`config.private_environment_config` block.

### At least one of `node_count`, `node_config`, `software_config` or `private_environment_config` required on `google_composer_environment.config`

In an attempt to avoid allowing empty blocks in config files, at least one of `node_count`, `node_config`, `software_config` or `private_environment_config` is now required on the `config` block.

## Resource: `google_compute_backend_bucket`

### `signed_url_cache_max_age_sec` is now required on `google_compute_backend_bucket.autoscaling_policy.cdn_policy`

Previously the default value of `signed_url_cache_max_age_sec` was `3600`. In an attempt to avoid allowing empty
blocks in config files, `signed_url_cache_max_age_sec` is now required on the
`autoscaling_policy.cdn_policy` block.

## Resource: `google_compute_backend_service`

### At least one of `connect_timeout`, `max_requests_per_connection`, `max_connections`, `max_pending_requests`, `max_requests`,  or `max_retries` is now required on `google_compute_backend_service.circuit_breakers`

In an attempt to avoid allowing empty blocks in config files, at least one of `connect_timeout`,
`max_requests_per_connection`, `max_connections`, `max_pending_requests`, `max_requests`,
or `max_retries` is now required on the `circuit_breakers` block.

###  At least one of `ttl`, `name`, or `path` is now required on `google_compute_backend_service.consistent_hash.http_cookie`

In an attempt to avoid allowing empty blocks in config files, at least one of `ttl`, `name`, or `path`
is now required on the `consistent_hash.http_cookie` block.

### At least one of `http_cookie`, `http_header_name`, or `minimum_ring_size` is now required on `google_compute_backend_service.consistent_hash`

In an attempt to avoid allowing empty blocks in config files, at least one of `http_cookie`,
`http_header_name`, or `minimum_ring_size` is now required on the `consistent_hash` block.

### At least one of `cache_key_policy` or `signed_url_cache_max_age_sec` is now required on `google_compute_backend_service.cdn_policy`

In an attempt to avoid allowing empty blocks in config files, at least one of `cache_key_policy` or
`signed_url_cache_max_age_sec` is now required on the `cdn_policy` block.

### At least one of `include_host`, `include_protocol`, `include_query_string`, `query_string_blacklist`, or `query_string_whitelist` is now required on `google_compute_backend_service.cdn_policy.cache_key_policy`

In an attempt to avoid allowing empty blocks in config files, at least one of `include_host`,
`include_protocol`, `include_query_string`, `query_string_blacklist`, or `query_string_whitelist`
is now required on the `cdn_policy.cache_key_policy` block.

### At least one of `base_ejection_time`, `consecutive_errors`, `consecutive_gateway_failure`, `enforcing_consecutive_errors`, `enforcing_consecutive_gateway_failure`, `enforcing_success_rate`, `interval`, `max_ejection_percent`, `success_rate_minimum_hosts`, `success_rate_request_volume`, or `success_rate_stdev_factor` is now required on `google_compute_backend_service.outlier_detection`

In an attempt to avoid allowing empty blocks in config files, at least one of `base_ejection_time`,
`consecutive_errors`, `consecutive_gateway_failure`, `enforcing_consecutive_errors`,
`enforcing_consecutive_gateway_failure`, `enforcing_success_rate`, `interval`, `max_ejection_percent`,
`success_rate_minimum_hosts`, `success_rate_request_volume`, or `success_rate_stdev_factor`
is now required on the `outlier_detection` block.

### At least one of `enable` or `sample_rate` is now required on `google_compute_backend_service.log_config`

In an attempt to avoid allowing empty blocks in config files, at least one of `enable` or `sample_rate`
is now required on the `log_config` block.

## Resource: `google_compute_firewall`

### Exactly one of `allow` or `deny` is required on `google_compute_firewall`

In attempt to be more consistent with the API, exactly one of `allowed` or `denied`
is now required on `google_compute_firewall`.

## Resource: `google_compute_forwarding_rule`

### `ip_version` is now removed

`ip_version` is not used for regional forwarding rules.

### `ip_address` is now strictly validated to enforce literal IP address format

Previously documentation suggested Terraform could use the same range of valid
IP Address formats for `ip_address` as accepted by the API (e.g. named addresses
or URLs to GCP Address resources). However, the server returns only literal IP
addresses and thus caused diffs on re-apply (i.e. a permadiff). We amended
documenation to say Terraform only accepts literal IP addresses.

This is now strictly validated. While this shouldn't have a large breaking
impact as users would have already run into permadiff issues on re-apply,
there might be validation errors for existing configs. The solution is be to
replace other address formats with the IP address, either manually or by
interpolating values from a `google_compute_address` resource.

#### Old Config (that would have permadiff)

```hcl
resource "google_compute_address" "my-addr" {
  name = "my-addr"
}

resource "google_compute_forwarding_rule" "frule" {
  name = "my-forwarding-rule"

  address = google_compute_address.my-addr.self_link
}
```

#### New Config

```hcl
resource "google_compute_address" "my-addr" {
  name = "my-addr"
}

resource "google_compute_forwarding_rule" "frule" {
  name = "my-forwarding-rule"

  address = google_compute_address.my-addr.address
}
```

## Resource: `google_compute_global_forwarding_rule`

### `ip_address` is now validated to enforce literal IP address format

See [`google_compute_forwarding_rule`](#resource-google_compute_forwarding_rule).

## Resource: `google_compute_health_check`

### Exactly one of `http_health_check`, `https_health_check`, `http2_health_check`, `tcp_health_check` or `ssl_health_check` is required on `google_compute_health_check`

In attempt to be more consistent with the API, exactly one of `http_health_check`, `https_health_check`,
`http2_health_check`, `tcp_health_check` or `ssl_health_check` is now required on
`google_compute_health_check`.

### At least one of `host`, `request_path`, `response`, `port`, `port_name`, `proxy_header`, or `port_specification` is now required on `google_compute_health_check.http_health_check`, `google_compute_health_check.https_health_check` and `google_compute_health_check.http2_health_check`

In an attempt to avoid allowing empty blocks in config files, at least one of `host`, `request_path`, `response`,
`port`, `port_name`, `proxy_header`, or `port_specification` is now required on the
`http_health_check`, `https_health_check` and `http2_health_check` blocks.

### At least one of `request`, `response`, `port`, `port_name`, `proxy_header`, or `port_specification` is now required on `google_compute_health_check.ssl_health_check` and `google_compute_health_check.tcp_health_check`

In an attempt to avoid allowing empty blocks in config files, at least one of `request`, `response`, `port`, `port_name`,
`proxy_header`, or `port_specification` is now required on the `ssl_health_check` and `tcp_health_check` blocks.

## Resource: `google_compute_image`

### `type` is now required on `google_compute_image.guest_os_features`

In an attempt to avoid allowing empty blocks in config files, `type` is now required on the
`guest_os_features` block.

## Resource: `google_compute_instance`

### `interface` is now required on block `google_compute_instance.scratch_disk`

Previously the default value of `interface` was `SCSI`. In an attempt to avoid allowing empty blocks
in config files, `interface` is now required on the `scratch_disk` block.

### At least one of `auto_delete`, `device_name`, `disk_encryption_key_raw`, `kms_key_self_link`, `initialize_params`, `mode` or `source` is now required on `google_compute_instance.boot_disk`

In an attempt to avoid allowing empty blocks in config files, at least one of `auto_delete`, `device_name`,
`disk_encryption_key_raw`, `kms_key_self_link`, `initialize_params`, `mode` or `source` is now required on the
`boot_disk` block.

### At least one of `size`, `type`, `image`, or `labels` is now required on `google_compute_instance.boot_disk.initialize_params`

In an attempt to avoid allowing empty blocks in config files, at least one of `size`, `type`, `image`, or `labels`
is now required on the `initialize_params` block.

### At least one of `enable_secure_boot`, `enable_vtpm`, or `enable_integrity_monitoring` is now required on `google_compute_instance.shielded_instance_config`

In an attempt to avoid allowing empty blocks in config files, at least one of `enable_secure_boot`, `enable_vtpm`,
or `enable_integrity_monitoring` is now required on the `shielded_instance_config` block.

### At least one of `on_host_maintenance`, `automatic_restart`, `preemptible`, or `node_affinities` is now required on `google_compute_instance.scheduling`

In an attempt to avoid allowing empty blocks in config files, at least one of `on_host_maintenance`, `automatic_restart`,
`preemptible`, or `node_affinities` is now required on the `scheduling` block.

## Resource: `google_compute_instance_group_manager`

The following changes apply to both `google_compute_instance_group_manager` and `google_compute_region_instance_group_manager`.

### `instance_template` has been replaced by `version.instance_template`

Instance group managers should be using `version` blocks to reference which
instance template to use for provisioning. To upgrade use a single `version`
block with `instance_template` in your config and by default all traffic will be
directed to that version.

### Old Config

```hcl
resource "google_compute_instance_group_manager" "my_igm" {
    name               = "my-igm"
    zone               = "us-central1-c"
    base_instance_name = "igm"

    instance_template = google_compute_instance_template.my_tmpl.self_link
}
```

### New Config

```hcl
resource "google_compute_instance_group_manager" "my_igm" {
    name               = "my-igm"
    zone               = "us-central1-c"
    base_instance_name = "igm"

    version {
        name = "prod"
        instance_template = google_compute_instance_template.my_tmpl.self_link
    }
}
```

### `update_strategy` has been replaced by `update_policy`

To allow much greater control over the updates happening to instance groups
`update_strategy` has been replaced by `update_policy`. The functionality controlled by `update_strategy` is now controlled by a combination of `update_policy.type` and `update_policy.minimal_action`. `update_strategy = NONE` can be achieved with `type = OPPORTUNISTIC`. The previous values of `RESTART` and `REPLACE` were both `PROACTIVE` types implicitly previously but can now be controlled explicitly.

For more details see the
[official guide](https://cloud.google.com/compute/docs/instance-groups/rolling-out-updates-to-managed-instance-groups).

### Old Config

```hcl
resource "google_compute_instance_group_manager" "my_igm" {
    name               = "my-igm"
    zone               = "us-central1-c"
    base_instance_name = "igm"

    instance_template = "${google_compute_instance_template.my_tmpl.self_link}"

    update_strategy   = "NONE"
}
```

### New Config

```hcl
resource "google_compute_instance_group_manager" "my_igm" {
    name               = "my-igm"
    zone               = "us-central1-c"
    base_instance_name = "igm"

    version {
        name = "prod"
        instance_template = "${google_compute_instance_template.my_tmpl.self_link}"
    }

    update_policy {
      minimal_action = "RESTART"
      type           = "OPPORTUNISTIC"
    }
}
```

## Resource: `google_compute_instance_template`

### At least one of `enable_secure_boot`, `enable_vtpm`, or `enable_integrity_monitoring` is now required on `google_compute_instance_template.shielded_instance_config`

In an attempt to avoid allowing empty blocks in config files, at least one of `enable_secure_boot`, `enable_vtpm`, or
`enable_integrity_monitoring` is now required on the `shielded_instance_config` block.

### At least one of `on_host_maintenance`, `automatic_restart`, `preemptible`, or `node_affinities` is now required on `google_compute_instance_template.scheduling`

In an attempt to avoid allowing empty blocks in config files, at least one of `on_host_maintenance`, `automatic_restart`,
`preemptible`, or `node_affinities` is now required on the `scheduling` block.

### Disks with invalid scratch disk configurations are now rejected

The instance template API allows specifying invalid configurations in some cases,
and an error is only returned when attempting to provision them. Terraform will
now report that some configs that previously appeared valid at plan time are
now invalid.

A disk with `type` `"SCRATCH"` must have `disk_type` `"local-ssd"` and a size of 375GB. For example,
the following is valid:

```hcl
disk {
    auto_delete  = true
    type         = "SCRATCH"
    disk_type    = "local-ssd"
    disk_size_gb = 375
}
```

These configs would have been accepted by Terraform previously, but will now
fail:

```hcl
disk {
    source_image = "https://www.googleapis.com/compute/v1/projects/gce-uefi-images/global/images/centos-7-v20190729"
    auto_delete  = true
    type         = "SCRATCH"
}
```

```hcl
disk {
    source_image = "https://www.googleapis.com/compute/v1/projects/gce-uefi-images/global/images/centos-7-v20190729"
    auto_delete  = true
    disk_type    = "local-ssd"
}
```

```hcl
disk {
    auto_delete  = true
    type         = "SCRATCH"
    disk_type    = "local-ssd"
    disk_size_gb = 300
}
```

### `kms_key_self_link` is now required on block `google_compute_instance_template.disk_encryption_key`

In an attempt to avoid allowing empty blocks in config files, `kms_key_self_link` is now
required on the `disk_encryption_key` block.

## Resource: `google_compute_network`

### `ipv4_range` is now removed

Legacy Networks are removed and you will no longer be able to create them
using this field from Feb 1, 2020 onwards.

## Resource: `google_compute_network_peering`

### `auto_create_routes` is now removed

`auto_create_routes` has been removed because it's redundant and not
user-configurable.

## Resource: `google_compute_node_template`

###  At least one of `cpus` or `memory` is now required on `google_compute_node_template.node_type_flexibility`

In an attempt to avoid allowing empty blocks in config files, at least one of `cpus` or `memory`
is now required on the `node_type_flexibility` block.

## Resource: `google_compute_region_backend_service`

### At least one of `connect_timeout`, `max_requests_per_connection`, `max_connections`, `max_pending_requests`, `max_requests`,  or `max_retries` is now required on `google_compute_region_backend_service.circuit_breakers`

In an attempt to avoid allowing empty blocks in config files, at least one of `connect_timeout`,
`max_requests_per_connection`, `max_connections`, `max_pending_requests`, `max_requests`,
or `max_retries` is now required on the `circuit_breakers` block.

###  At least one of `ttl`, `name`, or `path` is now required on `google_compute_region_backend_service.consistent_hash.http_cookie`

In an attempt to avoid allowing empty blocks in config files, at least one of `ttl`, `name`, or `path`
is now required on the `consistent_hash.http_cookie` block.

### At least one of `http_cookie`, `http_header_name`, or `minimum_ring_size` is now required on `google_compute_region_backend_service.consistent_hash`

In an attempt to avoid allowing empty blocks in config files, at least one of `http_cookie`,
`http_header_name`, or `minimum_ring_size` is now required on the `consistent_hash` block.

### At least one of `disable_connection_drain_on_failover`, `drop_traffic_if_unhealthy`, or `failover_ratio` is now required on `google_compute_region_backend_service.failover_policy`

In an attempt to avoid allowing empty blocks in config files, at least one of `disable_connection_drain_on_failover`,
`drop_traffic_if_unhealthy`, or `failover_ratio` is now required on the `failover_policy` block.

### At least one of `base_ejection_time`, `consecutive_errors`, `consecutive_gateway_failure`, `enforcing_consecutive_errors`, `enforcing_consecutive_gateway_failure`, `enforcing_success_rate`, `interval`, `max_ejection_percent`, `success_rate_minimum_hosts`, `success_rate_request_volume`, or `success_rate_stdev_factor` is now required on `google_compute_region_backend_service.outlier_detection`

In an attempt to avoid allowing empty blocks in config files, at least one of `base_ejection_time`,
`consecutive_errors`, `consecutive_gateway_failure`, `enforcing_consecutive_errors`,
`enforcing_consecutive_gateway_failure`, `enforcing_success_rate`, `interval`, `max_ejection_percent`,
`success_rate_minimum_hosts`, `success_rate_request_volume`, or `success_rate_stdev_factor`
is now required on the `outlier_detection` block.

### At least one of `enable` or `sample_rate` is now required on `google_compute_region_backend_service.log_config`

In an attempt to avoid allowing empty blocks in config files, at least one of `enable` or `sample_rate`
is now required on the `log_config` block.

## Resource: `google_compute_region_health_check`

### Exactly one of `http_health_check`, `https_health_check`, `http2_health_check`, `tcp_health_check` or `ssl_health_check` is required on `google_compute_health_check`

In attempt to be more consistent with the API, exactly one of `http_health_check`, `https_health_check`,
`http2_health_check`, `tcp_health_check` or `ssl_health_check` is now required on the
`google_compute_region_health_check`.

### At least one of `host`, `request_path`, `response`, `port`, `port_name`, `proxy_header`, or `port_specification` is now required on `google_compute_region_health_check.http_health_check`, `google_compute_region_health_check.https_health_check` and `google_compute_region_health_check.http2_health_check`

In an attempt to avoid allowing empty blocks in config files, at least one of `host`, `request_path`, `response`,
`port`, `port_name`, `proxy_header`, or `port_specification` is now required on the
`http_health_check`, `https_health_check` and `http2_health_check` blocks.

### At least one of `request`, `response`, `port`, `port_name`, `proxy_header`, or `port_specification` is now required on `google_compute_region_health_check.ssl_health_check` and `google_compute_region_health_check.tcp_health_check`

In an attempt to avoid allowing empty blocks in config files, at least one of `request`, `response`, `port`, `port_name`,
`proxy_header`, or `port_specification` is now required on the `ssl_health_check` and `tcp_health_check` blocks.

## Resource: `google_compute_resource_policy`

### Exactly one of `hourly_schedule`, `daily_schedule` or `weekly_schedule` is now required on `google_compute_resource_policy.snapshot_schedule_policy.schedule`

In an attempt to avoid allowing empty blocks in config files, exactly one
of `hourly_schedule`, `daily_schedule` or `weekly_schedule` is now required
on the `snapshot_schedule_policy.schedule` block.

### At least one of `labels`, `storage_locations`, or `guest_flush` is now required on `google_compute_resource_policy.snapshot_schedule_policy.snapshot_properties`

In an attempt to avoid allowing empty blocks in config files, at least one of
`labels`, `storage_locations`, or `guest_flush` is now required on the
`snapshot_schedule_policy.snapshot_properties` block.

## Resource: `google_compute_route`

### Exactly one of `next_hop_gateway`, `next_hop_instance`, `next_hop_ip`, `next_hop_vpn_tunnel` or `next_hop_ilb` is required on `google_compute_route`

In attempt to be more consistent with the API, exactly one of `next_hop_gateway`, `next_hop_instance`,
`next_hop_ip`, `next_hop_vpn_tunnel` or `next_hop_ilb` is now required on the
`google_compute_route`.

## Resource: `google_compute_router`

### `range` is now required on `google_compute_router.bgp.advertised_ip_ranges`

In an attempt to avoid allowing empty blocks in config files, `range` is now
required on the `bgp.advertised_ip_ranges` block.

## Resource: `google_compute_router_peer`

### `range` is now required on block `google_compute_router_peer.advertised_ip_ranges`

In an attempt to avoid allowing empty blocks in config files, `range` is now
required on the `advertised_ip_ranges` block.

## Resource: `google_compute_snapshot`

### `raw_key` is now required on block `google_compute_snapshot.source_disk_encryption_key`

In an attempt to avoid allowing empty blocks in config files, `raw_key` is now
required on the `source_disk_encryption_key` block.

## Resource: `google_compute_subnetwork`

### `enable_flow_logs` is now removed

`enable_flow_logs` has been removed and should be replaced by the `log_config` block with configurations
for flow logging. Enablement of flow logs is now controlled by whether `log_config` is defined or not instead
of by the `enable_flow_logs` variable. Users with `enable_flow_logs = false` only need to remove the field.

### At least one of `aggregation_interval`, `flow_sampling`, or `metadata` is now required on `google_compute_subnetwork.log_config`

In an attempt to avoid allowing empty blocks in config files, at least one of
`aggregation_interval`, `flow_sampling`, or `metadata` is now required on the
`log_config` block.


### Old Config

```hcl
resource "google_compute_subnetwork" "subnet-with-logging" {
  name          = "log-test-subnetwork"
  ip_cidr_range = "10.2.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.custom-test.self_link

  enable_flow_logs = true
}
```


### New Config

```hcl
resource "google_compute_subnetwork" "subnet-with-logging" {
  name          = "log-test-subnetwork"
  ip_cidr_range = "10.2.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.custom-test.self_link

  log_config {
    aggregation_interval = "INTERVAL_10_MIN"
    flow_sampling        = 0.5
    metadata             = "INCLUDE_ALL_METADATA"
  }
}
```


## Resource: `google_container_cluster`

### `ip_allocation_policy` will catch out-of-band changes, `use_ip_aliases` removed

-> This change and "Automatic subnetwork creation for VPC-native clusters
removed" are related; see the other entry for more details.

In `2.X`, `ip_allocation_policy` wouldn't cause a diff if it was undefined in
config but was set on the cluster itself. Additionally, it could be defined with
`use_ip_aliases` set to `false`. However, this made it difficult to reason about
whether a cluster was routes-based or VPC-native.

With `3.0.0`, Terraform will detect drift on the block. The configuration has also
been simplified. Terraform creates a VPC-native cluster when
`ip_allocation_policy` is defined (`use_ip_aliases` is implicitly set to true
and is no longer configurable). When the block is undefined, Terraform creates a
routes-based cluster.

Other than removing the `use_ip_aliases` field, most users of VPC-native clusters
won't be affected. `terraform plan` will show a diff if a config doesn't contain
`ip_allocation_policy` but the underlying cluster does. Routes-based cluster
users may need to remove `ip_allocation_policy` if `use_ip_aliases` had been set
to `false`.

#### Old Config

```hcl
resource "google_container_cluster" "primary" {
  name       = "my-cluster"
  location   = "us-central1"

  initial_node_count = 1

  ip_allocation_policy {
    use_ip_aliases = false
  }
}
```

#### New Config

```hcl
resource "google_container_cluster" "primary" {
  name       = "my-cluster"
  location   = "us-central1"

  initial_node_count = 1
}
```


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
* `ip_allocation_policy` will catch drift when not in config
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
  network    = google_compute_network.container_network.name

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
  network       = google_compute_network.container_network.self_link
}

resource "google_container_cluster" "primary" {
  name       = "my-cluster"
  location   = "us-central1"
  network    = google_compute_network.container_network.name
  subnetwork = google_compute_subnetwork.container_subnetwork.name

  initial_node_count = 1

  ip_allocation_policy {
    cluster_ipv4_cidr_block  = "10.0.0.0/16"
    services_ipv4_cidr_block = "10.1.0.0/16"
  }
}
```

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

### `addons_config.kubernetes_dashboard` is now removed

The `kubernetes_dashboard` addon is deprecated for clusters on GKE and
will soon be removed. It is recommended to use alternative GCP Console
dashboards.

### `channel` is now required on `google_container_cluster.release_channel`

In an attempt to avoid allowing empty blocks in config files, `channel` is now
required on the `release_channel` block.

### The `disabled` field is now required on the `addons_config` blocks for `http_load_balancing`, `horizontal_pod_autoscaling`, `istio_config`, `cloudrun_config` and `network_policy_config`.

In an attempt to avoid allowing empty blocks in config files, `disabled` is now
required on the different `google_container_cluster.addons_config` blocks.

### Exactly one of `daily_maintenance_window` or `recurring_window` is now required on `google_container_cluster.maintenance_policy`

In an attempt to avoid allowing empty blocks in config files, exactly one of `daily_maintenance_window` or `recurring_window` is now required on the
`maintenance_policy` block.

### At least one of `http_load_balancing`, `horizontal_pod_autoscaling` , `network_policy_config`, `cloudrun_config`, or `istio_config` is now required on `google_container_cluster.addons_config`

In an attempt to avoid allowing empty blocks in config files, at least one of `http_load_balancing`,
`horizontal_pod_autoscaling` , `network_policy_config`, `cloudrun_config`, or `istio_config` is now required on the
`addons_config` block.

### At least one of `username`, `password` or `client_certificate_config` is now required on `google_container_cluster.master_auth`

In an attempt to avoid allowing empty blocks in config files, at least one of `username`, `password`
or `client_certificate_config` is now required on the `master_auth` block.

### `enabled` is now required on block `google_container_cluster.vertical_pod_autoscaling`

In an attempt to avoid allowing empty blocks in config files, `enabled` is now
required on the `vertical_pod_autoscaling` block.

### `enabled` is now required on block `google_container_cluster.network_policy`

Previously the default value of `enabled` was `false`. In an attempt to avoid allowing empty blocks
in config files, `enabled` is now required on the `network_policy` block.

### `enable_private_endpoint` is now required on block `google_container_cluster.private_cluster_config`

In an attempt to avoid allowing empty blocks in config files, `enable_private_endpoint` is now
required on the `private_cluster_config` block.

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

### `zone`, `region` and `additional_zones` are now removed

`zone` and `region` have been removed in favor of `location` and
`additional_zones` has been removed in favor of `node_locations`

## Resource: `google_container_node_pool`

### `zone` and `region` are now removed

`zone` and `region` have been removed in favor of `location`

## Resource: `google_dataproc_autoscaling_policy`

### At least one of `min_instances`, `max_instances`, or `weight` is now required on `google_dataproc_autoscaling_policy.secondary_worker_config`

In an attempt to avoid allowing empty blocks in config files, at least one of `min_instances`,
`max_instances`, or `weight` is now required on the `secondary_worker_config`
block.

## Resource: `google_dataproc_cluster`

### At least one of `staging_bucket`, `gce_cluster_config`, `master_config`, `worker_config`, `preemptible_worker_config`, `software_config`, `initialization_action` or `encryption_config` is now required on `google_dataproc_cluster.cluster_config`

In an attempt to avoid allowing empty blocks in config files, at least one of `staging_bucket`,
`gce_cluster_config`, `master_config`, `worker_config`, `preemptible_worker_config`, `software_config`,
`initialization_action` or `encryption_config` is now required on the
`cluster_config` block.

### At least one of `image_version`, `override_properties` or `optional_components` is now required on `google_dataproc_cluster.cluster_config.software_config`

In an attempt to avoid allowing empty blocks in config files, at least one of `image_version`,
`override_properties` or `optional_components` is now required on the
`cluster_config.software_config` block.

### At least one of `num_instances` or `disk_config` is now required on `google_dataproc_cluster.cluster_config.preemptible_worker_config`

In an attempt to avoid allowing empty blocks in config files, at least one of `num_instances`
or `disk_config` is now required on the `cluster_config.preemptible_worker_config` block.

### At least one of `zone`, `network`, `subnetwork`, `tags`, `service_account`, `service_account_scopes`, `internal_ip_only` or `metadata` is now required on `google_dataproc_cluster.cluster_config.gce_cluster_config`

In an attempt to avoid allowing empty blocks in config files, at least one of `zone`, `network`, `subnetwork`,
`tags`, `service_account`, `service_account_scopes`, `internal_ip_only` or `metadata` is now required on the
`gce_cluster_config` block.

### At least one of `num_instances`, `image_uri`, `machine_type`, `min_cpu_platform`, `disk_config`, or `accelerators` is now required on `google_dataproc_cluster.cluster_config.master_config` and `google_dataproc_cluster.cluster_config.worker_config`

In an attempt to avoid allowing empty blocks in config files, at least one of `num_instances`, `image_uri`,
`machine_type`, `min_cpu_platform`, `disk_config`, or `accelerators` is now required on the
`cluster_config.master_config` and `cluster_config.worker_config` blocks.

### At least one of `num_local_ssds`, `boot_disk_size_gb` or `boot_disk_type` is now required on `google_dataproc_cluster.cluster_config.preemptible_worker_config.disk_config`, `google_dataproc_cluster.cluster_config.master_config.disk_config` and `google_dataproc_cluster.cluster_config.worker_config.disk_config`

In an attempt to avoid allowing empty blocks in config files, at least one of `num_local_ssds`, `boot_disk_size_gb`
or `boot_disk_type` is now required on the `cluster_config.preemptible_worker_config.disk_config`,
`cluster_config.master_config.disk_config` and `cluster_config.worker_config.disk_config` blocks.


### `policy_uri` is now required on `google_dataproc_cluster.autoscaling_config` block.

In an attempt to avoid allowing empty blocks in config files, `policy_uri` is now
required on the `autoscaling_config` block.

## Resource: `google_dataproc_job`

### At least one of `query_file_uri` or `query_list` is now required on `hive_config`, `pig_config`, and `sparksql_config`

In an attempt to avoid allowing empty blocks in config files, at least one of
`query_file_uri` or `query_list` is now required on the `hive_config`, `pig_config`, and
`sparksql_config` blocks.

### At least one of `main_class` or `main_jar_file_uri` is now required on `google_dataproc_job.spark_config` and `google_dataproc_job.hadoop_config`

In an attempt to avoid allowing empty blocks in config files, at least one of
`main_class` or `main_jar_file_uri` is now required on the `spark_config`
and `hadoop_config` blocks.

### `driver_log_levels` is now required on `logging_config` blocks for `pyspark_config`, `hadoop_config`, `spark_config`, `pig_config`, and `sparksql_config`.

In an attempt to avoid allowing empty blocks in config files, `driver_log_levels` is now
required on `pyspark_config`, `hadoop_config`, `spark_config`, `pig_config`, and
`sparksql_config` blocks.

### `max_failures_per_hour` is now required on block `google_dataproc_job.scheduling`

In an attempt to avoid allowing empty blocks in config files, `max_failures_per_hour` is now
required on the `scheduling` block.

## Resource: `google_dns_managed_zone`

### At least one of `kind`, `non_existence`, `state`,  or `default_key_specs` is now required on `google_dns_managed_zone.dnssec_config`

In an attempt to avoid allowing empty blocks in config files, at least one of
`kind`, `non_existence`, `state`,  or `default_key_specs` is now required on the
`dnssec_config` block.

### `target_network` is now required on block `google_dns_managed_zone.peering_config`

In an attempt to avoid allowing empty blocks in config files, `target_network` is now
required on the `peering_config` block.

### `network_url` is now required on block `google_dns_managed_zone.peering_config.target_network`

In an attempt to avoid allowing empty blocks in config files, `network_url` is now
required on the `peering_config.target_network` block.

### `target_name_servers` is now required on block `google_dns_managed_zone.forwarding_config`

In an attempt to avoid allowing empty blocks in config files, `target_name_servers` is now
required on the `forwarding_config` block.

### `ipv4_address` is now required on block `google_dns_managed_zone.forwarding_config.target_name_servers`

In an attempt to avoid allowing empty blocks in config files, `ipv4_address` is now
required on the `forwarding_config.target_name_servers` block.

### `target_name_servers` is now required on block `google_dns_managed_zone.forwarding_config`

In an attempt to avoid allowing empty blocks in config files, `target_name_servers` is now
required on the `forwarding_config` block.

### `networks` is now required on block `google_dns_managed_zone.private_visibility_config`

In an attempt to avoid allowing empty blocks in config files, `networks` is now
required on the `private_visibility_config` block.

### `network_url` is now required on block `google_dns_managed_zone.private_visibility_config.networks`

In an attempt to avoid allowing empty blocks in config files, `network_url` is now
required on the `private_visibility_config.networks` block.

## Resource: `google_dns_policy`

### `network_url` is now required on block `google_dns_policy.networks`

In an attempt to avoid allowing empty blocks in config files, `network_url` is now
required on the `networks` block.

### `target_name_servers` is now required on block `google_dns_policy.alternative_name_server_config`

In an attempt to avoid allowing empty blocks in config files, `target_name_servers` is now
required on the `alternative_name_server_config` block.

### `ipv4_address` is now required on block `google_dns_policy.alternative_name_server_config.target_name_servers`

In an attempt to avoid allowing empty blocks in config files, `ipv4_address` is now
required on the `alternative_name_server_config.target_name_servers` block.

## Resource: `google_folder_organization_policy`

### Exactly one of `allow` or `deny` is now required on `google_folder_organization_policy.list_policy`

In an attempt to avoid allowing empty blocks in config files, exactly one of `allow` or `deny` is now
required on the `list_policy` block.

### Exactly one of `all` or `values` is now required on `google_folder_organization_policy.list_policy.allow` and `google_folder_organization_policy.list_policy.deny`

In an attempt to avoid allowing empty blocks in config files, exactly one of `all` or `values` is now
required on the `list_policy.allow` and `list_policy.deny` blocks.

## Resource: `google_healthcare_hl7_v2_store`

### At least one of `allow_null_header ` or `segment_terminator` is now required on `google_healthcare_hl7_v2_store.parser_config`

In an attempt to avoid allowing empty blocks in config files, at least one of `allow_null_header `
or `segment_terminator` is now required on the `parser_config` block.

## Resource: `google_logging_metric`

### At least one of `linear_buckets`, `exponential_buckets` or `explicit_buckets` is now required on `google_logging_metric.bucket_options`

In an attempt to avoid allowing empty blocks in config files, at least one of `linear_buckets`,
`exponential_buckets` or `explicit_buckets` is now required on the `bucket_options` block.

### At least one of `num_finite_buckets`, `width` or `offset` is now required on `google_logging_metric.bucket_options.linear_buckets`

In an attempt to avoid allowing empty blocks in config files, at least one of `num_finite_buckets`,
`width` or `offset` is now required on the `bucket_options.linear_buckets` block.

### At least one of `num_finite_buckets`, `growth_factor` or `scale` is now required on `google_logging_metric.bucket_options.exponential_buckets`

In an attempt to avoid allowing empty blocks in config files, at least one of `num_finite_buckets`,
`growth_factor` or `scale` is now required on the `bucket_options.exponential_buckets` block.

### `bounds` is now required on `google_logging_metric.bucket_options.explicit_buckets`

In an attempt to avoid allowing empty blocks in config files, `bounds` is now required on the
`bucket_options.explicit_buckets` block.

## Resource: `google_mlengine_model`

### `name` is now required on `google_mlengine_model.default_version`

In an attempt to avoid allowing empty blocks in config files, `name` is now required on the
`default_version` block.

## Resource: `google_monitoring_alert_policy`

### `labels` is now removed

`labels` is removed as it was never used. See `user_labels` for the correct field.

### At least one of `content` or `mime_type` is now required on `google_monitoring_alert_policy.documentation`

In an attempt to avoid allowing empty blocks in config files, at least one of `content` or `mime_type`
is now required on the `documentation` block.

## Resource: `google_monitoring_uptime_check_config`

### Exactly one of `resource_group` or `monitored_resource` is now required on `google_monitoring_uptime_check_config`

In attempt to be more consistent with the API, exactly one of `resource_group` or `monitored_resource` is now required
on `google_monitoring_uptime_check_config`.

### Exactly one of `http_check` or `tcp_check` is now required on `google_monitoring_uptime_check_config`

In attempt to be more consistent with the API, exactly one of `http_check` or `tcp_check` is now required
on `google_monitoring_uptime_check_config`.

### At least one of `auth_info`, `port`, `headers`, `path`, `use_ssl`, or `mask_headers` is now required on `google_monitoring_uptime_check_config.http_check`

In an attempt to avoid allowing empty blocks in config files, at least one of `auth_info`,
`port`, `headers`, `path`, `use_ssl`, or `mask_headers` is now required on the `http_check` block.

### At least one of `resource_type` or `group_id` is now required on `google_monitoring_uptime_check_config.resource_group`

In an attempt to avoid allowing empty blocks in config files, at least one of `resource_type` or `group_id`
is now required on the `resource_group` block.

### `content` is now required on block `google_monitoring_uptime_check_config.content_matchers`

In an attempt to avoid allowing empty blocks in config files, `content` is now
required on the `content_matchers` block.

### `username` and `password` are now required on block `google_monitoring_uptime_check_config.http_check.auth_info`

In an attempt to avoid allowing empty blocks in config files, `username` and `password` are now
required on the `http_check.auth_info` block.

### `is_internal` and `internal_checker` are now removed

`is_internal` and `internal_checker` never worked, and are now removed.

## Resource: `google_organization_policy`

### Exactly one of `allow` or `deny` is now required on `google_organization_policy.list_policy`

In an attempt to avoid allowing empty blocks in config files, exactly one of `allow` or `deny` is now
required on the `list_policy` block.

### Exactly one of `all` or `values` is now required on `google_organization_policy.list_policy.allow` and `google_organization_policy.list_policy.deny`

In an attempt to avoid allowing empty blocks in config files, exactly one of `all` or `values` is now
required on the `list_policy.allow` and `list_policy.deny` blocks.

## Resource: `google_project_iam_audit_config`

### Audit configs are now authoritative on create

Audit configs are now authoritative on create, rather than merging with existing configs on create.
Writing an audit config resource will now overwrite any existing audit configs on the given project.

## Resource: `google_project_organization_policy`

### Exactly one of `allow` or `deny` is now required on `google_project_organization_policy.list_policy`

In an attempt to avoid allowing empty blocks in config files, exactly one of `allow` or `deny` is now
required on the `list_policy` block.

### Exactly one of `all` or `values` is now required on `google_project_organization_policy.list_policy.allow` and `google_project_organization_policy.list_policy.deny`

In an attempt to avoid allowing empty blocks in config files, exactly one of `all` or `values` is now
required on the `list_policy.allow` and `list_policy.deny` blocks.

## Resource: `google_project_service`

### `bigquery-json.googleapis.com` service can no longer be specified

`bigquery-json.googleapis.com` is being renamed to `bigquery.googleapis.com` in
the upstream API. As a result, `bigquery-json.googleapis.com` has been
disallowed. Instead, please use `bigquery.googleapis.com`. The provider will
automatically convert between them as the upstream API migration continues.

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
requests. From `2.13.0` onwards, those requests are batched on write, and from `2.20.0` onwards,
batched on read. It's recommended that you upgrade to `2.13.0+` before migrating if you
encounter write quota issues or `2.20.0+` before migrating if you encounter read quota issues
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
resource "google_project_service" "service" {
  for_each = toset([
    "iam.googleapis.com",
    "cloudresourcemanager.googleapis.com",
  ])

  service = each.key

  project = "your-project-id"
  disable_on_destroy = false
}
```

## Resource: `google_pubsub_subscription`

### `name` must now be a short name

`name` previously could have been specified by a long name (e.g. `projects/my-project/subscriptions/my-subscription`)
or a shortname (e.g. `my-subscription`). `name` now must be the shortname.

### `ttl` is now required on `google_pubsub_subscription.expiration_policy`

Previously, an empty `expiration_policy` block would allow the resource to never expire. In an attempt to avoid
allowing empty blocks in config files, `ttl` is now required on the `expiration_policy` block.  `ttl` should be set
to `""` for the resource to never expire.

## Resource: `google_security_scanner_scan_config`

### At least one of `google_account` or `custom_account` is now required on `google_security_scanner_scan_config.authentication`

In an attempt to avoid allowing empty blocks in config files, at least one of `google_account` or
`custom_account` is now required on the `authentication` block.

## Resource: `google_service_account_key`

### `pgp_key`, `private_key_fingerprint`, and `private_key_encrypted` are now removed

`google_service_account_key` previously supported encrypting the private key with
a supplied PGP key. This is [no longer supported](https://www.terraform.io/docs/extend/best-practices/sensitive-state.html#don-39-t-encrypt-state)
and has been removed as functionality. State should instead be treated as sensitive,
and ideally encrypted using a remote state backend.

This will require re-provisioning your service account key, unfortunately. There
is no known alternative at this time.

## Resource: `google_sql_database_instance`

### At least one of `ca_certificate`, `client_certificate`, `client_key`, `connect_retry_interval`, `dump_file_path`, `failover_target`, `master_heartbeat_period`, `password`, `ssl_cipher`, `username`, or `verify_server_certificate` is now required on `google_sql_database_instance.settings.replica_configuration`

In an attempt to avoid allowing empty blocks in config files, at least one of `ca_certificate`, `client_certificate`, `client_key`, `connect_retry_interval`,
`dump_file_path`, `failover_target`, `master_heartbeat_period`, `password`, `ssl_cipher`, `username`, or `verify_server_certificate` is now required on the
`settings.replica_configuration` block.

### At least one of `cert`, `common_name`, `create_time`, `expiration_time`, or `sha1_fingerprint` is now required on `google_sql_database_instance.settings.server_ca_cert`

In an attempt to avoid allowing empty blocks in config files, at least one of `cert`, `common_name`, `create_time`, `expiration_time`, or `sha1_fingerprint` is now required on the `settings.server_ca_cert` block.

### At least one of `day`, `hour`, or `update_track` is now required on `google_sql_database_instance.settings.maintenance_window`

In an attempt to avoid allowing empty blocks in config files, at least one of `day`, `hour`,
or `update_track` is now required on the `settings.maintenance_window` block.

### At least one of `binary_log_enabled`, `enabled`, `start_time`, or `location` is now required on `google_sql_database_instance.settings.backup_configuration`

In an attempt to avoid allowing empty blocks in config files, at least one of `binary_log_enabled`, `enabled`, `start_time`, or `location` is now required on the
`settings.backup_configuration` block.

### At least one of `authorized_networks`, `ipv4_enabled`, `require_ssl`, or `private_network` is now required on `google_sql_database_instance.settings.ip_configuration`

In an attempt to avoid allowing empty blocks in config files, at least one of `authorized_networks`, `ipv4_enabled`,
`require_ssl`, and `private_network` is now required on the `settings.ip_configuration` block.

### `name` and `value` are now required on block `google_sql_database_instance.settings.database_flags`

In an attempt to avoid allowing empty blocks in config files, `name` and `value` are now required on the `settings.database_flags` block.

### `value` is now required on block `google_sql_database_instance.settings.ip_configuration.authorized_networks`

In an attempt to avoid allowing empty blocks in config files, `value` is now required on the `settings.ip_configuration.authorized_networks` block.

### `zone` is now required on block `google_sql_database_instance.settings.location_preference`

In an attempt to avoid allowing empty blocks in config files, `zone` is now
required on the `settings.location_preference` block.

## Resource: `google_storage_bucket`

### `enabled` is now required on block `google_storage_bucket.versioning`

Previously the default value of `enabled` was `false`. In an attempt to avoid allowing empty blocks
in config files, `enabled` is now required on the `versioning` block.

### At least one of `main_page_suffix` or `not_found_page` is now required on `google_storage_bucket.website`

In an attempt to avoid allowing empty blocks in config files, at least one of `main_page_suffix` or
`not_found_page` is now required on the `website` block.

### At least one of `min_time_elapsed_since_last_modification`, `max_time_elapsed_since_last_modification`, `include_prefixes`, or `exclude_prefixes` is now required on `google_storage_transfer_job.transfer_spec.object_conditions`

In an attempt to avoid allowing empty blocks in config files, at least one of `min_time_elapsed_since_last_modification`,
`max_time_elapsed_since_last_modification`, `include_prefixes`, or `exclude_prefixes` is now required on the `transfer_spec.object_conditions` block.

### `is_live` is now removed

Please use `with_state` instead, as `is_live` is now removed.

## Resource: `google_storage_transfer_job`

### At least one of `overwrite_objects_already_existing_in_sink`, `delete_objects_unique_in_sink`, or `delete_objects_from_source_after_transfer` is now required on `google_storage_transfer_job.transfer_spec.transfer_options`

In an attempt to avoid allowing empty blocks in config files, at least one of `overwrite_objects_already_existing_in_sink`,
`delete_objects_unique_in_sink`, or `delete_objects_from_source_after_transfer` is now required on the
`transfer_spec.transfer_options` block.

### At least one of `gcs_data_source`, `aws_s3_data_source`, or `http_data_source` is now required on `google_storage_transfer_job.transfer_spec`

In an attempt to avoid allowing empty blocks in config files, at least one of `gcs_data_source`, `aws_s3_data_source`,
or `http_data_source` is now required on the `transfer_spec` block.

## Resource: `google_tpu_node`

### `preemptible` is now required on block `google_tpu_node.scheduling_config`

In an attempt to avoid allowing empty blocks in config files, `preemptible` is now
required on the `scheduling_config` block.
