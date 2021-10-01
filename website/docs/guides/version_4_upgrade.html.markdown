---
layout: "google"
page_title: "Terraform Google Provider 4.0.0 Upgrade Guide"
sidebar_current: "docs-google-provider-guides-version-4-upgrade"
description: |-
  Terraform Google Provider 4.0.0 Upgrade Guide
---

<!-- TOC depthFrom:2 depthTo:2 -->

- [Terraform Google Provider 4.0.0 Upgrade Guide](#terraform-google-provider-400-upgrade-guide)
  - [I accidentally upgraded to 4.0.0, how do I downgrade to `3.X`?](#i-accidentally-upgraded-to-400-how-do-i-downgrade-to-3x)
  - [Provider Version Configuration](#provider-version-configuration)
  - [Provider](#provider)
    - [Provider-level change example](#provider-level-change-example)
  - [Datasource: `google_product_resource`](#datasource-google_product_resource)
    - [Datasource-level change example](#datasource-level-change-example)
  - [Resource: `google_product_resource`](#resource-google_product_resource)
    - [Resource-level change example](#resource-level-change-example)

<!-- /TOC -->

# Terraform Google Provider 4.0.0 Upgrade Guide

The `4.0.0` release of the Google provider for Terraform is a major version and
includes some changes that you will need to consider when upgrading. This guide
is intended to help with that process and focuses only on the changes necessary
to upgrade from the final `3.X` series release to `4.0.0`.

Most of the changes outlined in this guide have been previously marked as
deprecated in the Terraform `plan`/`apply` output throughout previous provider
releases, up to and including the final `3.X` series release. These changes,
such as deprecation notices, can always be found in the CHANGELOG of the
affected providers. [google](https://github.com/hashicorp/terraform-provider-google/blob/master/CHANGELOG.md)
[google-beta](https://github.com/hashicorp/terraform-provider-google-beta/blob/master/CHANGELOG.md)

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
      version = "~> 3.87.0"
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


## Datasource: `google_product_resource`

### Datasource-level change example

Description of the change and how users should adjust their configuration (if needed).

## Resource: `google_product_resource`

### Resource-level change example

Description of the change and how users should adjust their configuration (if needed).
