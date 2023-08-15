---
page_title: "Terraform Google Provider 5.0.0 Upgrade Guide"
description: |-
  Terraform Google Provider 5.0.0 Upgrade Guide
---

# Terraform Google Provider 5.0.0 Upgrade Guide

The `5.0.0` release of the Google provider for Terraform is a major version and
includes some changes that you will need to consider when upgrading. This guide
is intended to help with that process and focuses only on the changes necessary
to upgrade from the final `4.X` series release to `5.0.0`.

Most of the changes outlined in this guide have been previously marked as
deprecated in the Terraform `plan`/`apply` output throughout previous provider
releases, up to and including the final `4.X` series release. These changes,
such as deprecation notices, can always be found in the CHANGELOG of the
affected providers. [google](https://github.com/hashicorp/terraform-provider-google/blob/main/CHANGELOG.md)
[google-beta](https://github.com/hashicorp/terraform-provider-google-beta/blob/main/CHANGELOG.md)

## I accidentally upgraded to 5.0.0, how do I downgrade to `4.X`?

If you've inadvertently upgraded to `5.0.0`, first see the
[Provider Version Configuration Guide](#provider-version-configuration) to lock
your provider version; if you've constrained the provider to a lower version
such as shown in the previous version example in that guide, Terraform will pull
in a `4.X` series release on `terraform init`.

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

-> Before upgrading to version 5.0.0, it is recommended to upgrade to the most
recent `4.X` series release of the provider, make the changes noted in this guide,
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
      version = "~> 4.70.0"
    }
  }
}
```

An updated configuration:

```hcl
terraform {
  required_providers {
    google = {
      version = "~> 5.0.0"
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

## Resource: `google_bigquery_table`

### At most one of `view`, `materialized_view`, and `schema` can be set.

The provider will now enforce at plan time that at most one of these fields be set.

## Resource: `google_firebaserules_release`

### Changing `ruleset_name` now triggers replacement

In 4.X.X, changing the `ruleset_name` in `google_firebaserules_release` updates the `Release` in place, which prevents the old `Ruleset` referred to by `ruleset_name` from being destroyed. A workaround is to use a `replace_triggered_by` lifecycle field on the `google_firebaserules_release`. In version 5.0.0, changing `ruleset_name` will trigger a replacement, which allows the `Ruleset` to be deleted. The `replace_triggered_by` workaround becomes unnecessary.

#### Old Config

```hcl
resource "google_firebaserules_release" "primary" {
  name         = "cloud.firestore"
  ruleset_name = "projects/my-project-name/rulesets/${google_firebaserules_ruleset.firestore.name}"
  project      = "my-project-name"

  lifecycle {
    replace_triggered_by = [
      google_firebaserules_ruleset.firestore
    ]
  }
}

resource "google_firebaserules_ruleset" "firestore" {
  source {
    files {
      content = "service cloud.firestore {match /databases/{database}/documents { match /{document=**} { allow read, write: if false; } } }"
      name    = "firestore.rules"
    }
  }

  project = "my-project-name"
}
```

#### New Config

```hcl
resource "google_firebaserules_release" "primary" {
  name         = "cloud.firestore"
  ruleset_name = "projects/my-project-name/rulesets/${google_firebaserules_ruleset.firestore.name}"
  project      = "my-project-name"
}

resource "google_firebaserules_ruleset" "firestore" {
  source {
    files {
      content = "service cloud.firestore {match /databases/{database}/documents { match /{document=**} { allow read, write: if false; } } }"
      name    = "firestore.rules"
    }
  }

  project = "my-project-name"
}
```

## Resource: `google_firebase_web_app`

### `deletion_policy` now defaults to `DELETE`

Previously, `google_firebase_web_app` deletions default to `ABANDON`, which means to only stop tracking the WebApp in Terraform. The actual app is not deleted from the Firebase project. If you are relying on this behavior, set `deletion_policy` to `ABANDON` explicitly in the new version.

## Resource: `google_cloud_run_v2_job`

### `startup_probe` and `liveness_probe` are now removed

These two unsupported fields were introduced incorrectly. They are now removed.

## Resource: `google_cloud_run_v2_service`

### `liveness_probe.tcp_socket` is now removed

This unsupported field was introduced incorrectly. It is now removed.

## Resource: `google_compute_router_nat`

### `enable_endpoint_independent_mapping` now defaults to API's default value which is `FALSE`

Previously, the default value of `enable_endpoint_independent_mapping` was `TRUE`. Now,
it will use the default value from the API which is `FALSE`. If you want to
enable endpoint independent mapping, then explicity set the value of
`enable_endpoint_independent_mapping` field to `TRUE`.

