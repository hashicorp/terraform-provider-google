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

### Provider-level Labels Rework

Labels and annotations are key-value pairs attached on Google cloud resources. Cloud labels are used for organizing resources, filtering resources, breaking down billing, and so on. Annotations are used to attach metadata to Kubernetes resources.

Not all of Google cloud resources support labels and annotations. Please check the Terraform Google provider resource documentation to figure out if the resource supports the `labels` and `annotations` fields.

#### Provider default labels

Default labels configured on the provider through the new `default_labels` field are now supported. The default labels configured on the provider will be applied to all of the resources with the top level `labels` field or the nested `labels` field inside the top level `metadata` field.

Provider-level default annotations are not supported.

#### Resource labels

The new labels model will be applied to all of the resources with the top level `labels` field or the nested `labels` field inside the top level `metadata` field. Some labels fields are for child resources, so the new model will not be applied to labels fields for child resources.

There are now three label-related fields with the new model:

* The `labels` field will be non-authoritative and only manage the labels defined by the users on the resource through Terraform.
* The output-only `effective_labels` will list all of labels present on the resource in GCP, including the labels configured through Terraform, other clients and services.
* The output-only `terraform_labels` will merge the labels defined by the users on the resource through Terraform and the default labels configured on the provider. If the same label exists on both the resource labels and provider default labels, the label on the resource will override the provider label.

After upgrading to `5.0.0`, and then running `terraform refresh` or `terraform apply`, these three fields should show in the state file of the resources with a self-applying `labels` field.

#### Resource annotations

The new annotations model is similar to the new labels model and will be applied to all of the resources with the top level `annotations` field or the nested `annotations` field inside the top level `metadata` field.

There are now two annotation-related fields with the new model, the `annotations` and the output-only `effective_annotations` fields.

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

### `schema` can only be represented as a JSON array with non-null elements.

The provider will now enforce at plan time that `schema` is a valid JSON array with non-null elements.

## Resource: `google_bigquery_routine`

### `routine_type` is now required.

The provider will now enforce at plan time that `routine_type` be set.

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

## Resource: `google_cloud_run_v2_job`

### `startup_probe` and `liveness_probe` are now removed

These two unsupported fields were introduced incorrectly. They are now removed.

## Resource: `google_cloud_run_v2_service`

### `liveness_probe.tcp_socket` is now removed

This unsupported field was introduced incorrectly. It is now removed.

## Resource: `google_dataplex_datascan`

### `dataQualityResult` and `dataProfileResult` output fields are now removed 

`dataQualityResult` and `dataProfileResult` were output-only fields which listed results for the latest job created under a Datascan. These were problematic fields that are unlikely to be relevant in a Terraform context. Removing them reduces the likelihood of additional parsing errors, and reduces maintenance overhead for the API surface.

## Resource: `google_compute_router_nat`

### `enable_endpoint_independent_mapping` now defaults to API's default value which is `FALSE`

Previously, the default value of `enable_endpoint_independent_mapping` was `TRUE`. Now,
it will use the default value from the API which is `FALSE`. If you want to
enable endpoint independent mapping, then explicity set the value of
`enable_endpoint_independent_mapping` field to `TRUE`.


## Resource: `google_firebase_project_location`

### `google_firebase_project_location` is now removed

In `4.X`, `google_firebase_project_location` would implicitly trigger creation of an App Engine application with a default Cloud Storage bucket and Firestore database, located in the specified `location_id`. In `5.0.0`, these resources should instead be set up explicitly using `google_app_engine_application` `google_firebase_storage_bucket`, and `google_firestore_database`.

For more information on configuring Firebase resources with Terraform, see [Get started with Terraform and Firebase](https://firebase.google.com/docs/projects/terraform/get-started).

#### Upgrade instructions

If you have existing resources created using `google_firebase_project_location`:
1. Remove the `google_firebase_project_location` block
1. Add blocks according to "New config" in this section for any of the following that you need: `google_app_engine_application`, `google_firebase_storage_bucket`, and/or `google_firestore_database`.
1. Import the existing resources corresponding to the blocks added in the previous step:
   `terraform import google_app_engine_application.default <project-id>`
   `terraform import google_firebase_storage_bucket.default-bucket <project-id>/<project-id>.appspot.com`
   `terraform import google_firestore_database.default "<project-id>/(default)"`

#### Old config

```hcl
resource "google_firebase_project_location" "basic" {
    provider = google-beta
    project = google_firebase_project.default.project

    location_id = "us-central"
}
```

#### New config

Assuming you use both the default Storage bucket and Firestore, an equivalent configuration would be:

```hcl
resource "google_app_engine_application" "default" {
  provider      = google-beta
  project       = google_firebase_project.default.project
  location_id   = "us-central"
  database_type = "CLOUD_FIRESTORE"

  depends_on = [
    google_firestore_database.default
  ]
}

resource "google_firebase_storage_bucket" "default-bucket" {
  provider  = google-beta
  project   = google_firebase_project.default.project
  bucket_id = google_app_engine_application.default.default_bucket
}

resource "google_firestore_database" "default" {
  project     = google_firebase_project.default.project
  name        = "(default)"
  location_id = "nam5"
  type        = "FIRESTORE_NATIVE"
}
```

## Resource: `google_firebase_web_app`

### `deletion_policy` now defaults to `DELETE`

Previously, `google_firebase_web_app` deletions default to `ABANDON`, which means to only stop tracking the WebApp in Terraform. The actual app is not deleted from the Firebase project. If you are relying on this behavior, set `deletion_policy` to `ABANDON` explicitly in the new version.
## Resource: `google_compute_autoscaler` (beta)

### `metric.filter` now defaults to `resource.type = gce_instance`

Previously, `metric.filter` doesn't have the defult value and causes a UI error.

## Resource: `google_privateca_certificate`

### `config_values` is now removed

Deprecated in favor of field `x509_description`. It is now removed.

### `pem_certificates` is now removed

Deprecated in favor of field `pem_certificate_chain`. It is now removed.

## Product: `gameservices`

### `gameservices` is now removed

This change involved the following resources: `google_game_services_game_server_cluster`, `google_game_services_game_server_deployment`, `google_game_services_game_server_config`, `google_game_services_realm` and `google_game_services_game_server_deployment_rollout`.

## Resource: `google_sql_database`

### `database_flags` is now a set

Previously, `database_flags` was a list, making it order-dependent. It is now a set.

If you were relying on accessing an individual flag by index (for example, `google_sql_database_instance.instance.settings.0.database_flags.0.name`), then that will now need to by hash (for example, `google_sql_database_instance.instance.settings.0.database_flags.<some-hash>.name`).
