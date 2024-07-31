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

## Resource: `google_product_resource`

### Resource-level change example header

Description of the change and how users should adjust their configuration (if needed).

## Resource: `google_sql_database_instance`

### `settings.ip_configuration.require_ssl` is now removed

Removed in favor of field `settings.ip_configuration.ssl_mode`.

## Resource: `google_pubsub_topic`

### `schema_settings` no longer has a default value

An empty value means the setting should be cleared.

## Resource: `google_cloud_run_v2_service`

### `liveness_probe` no longer defaults from API

Cloud Run does not provide a default value for liveness probe. Now removing this field
will remove the liveness probe from the Cloud Run service.

## Resource: `google_storage_bucket`

### `lifecycle_rule.condition.no_age` is now removed

Previously `lifecycle_rule.condition.age` attirbute was being set zero value by default and `lifecycle_rule.condition.no_age` was introduced to prevent that.
Now `lifecycle_rule.condition.no_age` is no longer supported and `lifecycle_rule.condition.age` won't set a zero value by default.
Removed in favor of the field `lifecycle_rule.condition.send_age_if_zero` which can be used to set zero value for `lifecycle_rule.condition.age` attribute. 

For a seamless update, if your state today uses `no_age=true`, update it to remove `no_age` and set `send_age_if_zero=false`. If you do not use `no_age=true`, you will need to add `send_age_if_zero=true` to your state to avoid any changes after updating to 6.0.0. 
