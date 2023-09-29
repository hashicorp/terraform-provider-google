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

Default labels configured on the provider through the new `default_labels` field are now supported. The default labels configured on the provider will be applied to all of the resources with the top level `labels` field or the nested `labels` field inside the top level `metadata` field. This change introduced a regression we were unable to resolve, and labels with an empty value (`""`) should be avoided, as they will be ignored and not included in `terraform_labels`, `effective_labels` and then API requests. Instead, labels with the value `_` or `true` are recommended.

Provider-level default annotations are not supported.

#### Resource labels

Labels and annotations fields on Terraform Google provider were authoritative and Terraform thought it was the only owner of the fields. This model worked well initially, but with the introduction of system labels and other client managed labels, Terraform would conflict with their labels and show a diff. We've reworked the `labels` field to resolve this class of problem.

The new labels model will be applied to all of the resources with the top level `labels` field or the nested `labels` field inside the top level `metadata` field. Some labels fields are for child resources, so the new model will not be applied to labels fields for child resources.

There are now three label-related fields with the new model:

* The `labels` field will be non-authoritative and only manage the labels defined by the users on the resource through Terraform. If a label was added outside of Terraform, it will not be managed by Terraform, unless it is added to the `labels` field in the configuration. The out of band labels will be listed in the `effective_labels` field. The new model introduced a regression we were unable to resolve, and the labels with an empty value (`""`) should be avoided, as they will be ignored and not included in `terraform_labels`, `effective_labels` and then API requests. Instead, labels with the value `_` or `true` are recommended.
* The output-only `terraform_labels` will merge the labels defined by the users on the resource through Terraform and the default labels configured on the provider. If the same label exists on both the resource labels and provider default labels, the label on the resource will override the provider label.
* The output-only `effective_labels` will list all of labels present on the resource in GCP, including the labels configured through Terraform, the system, and other clients.

**Note:** `ignore_changes` can be applied to `labels` field to ignore the changes of the user defined labels. It is not recommended to apply `ignore_changes` to `terraform_labels` or `effective_labels`, as it may unintuitively affect the final API call.

The following changes will be observed after upgrading to `5.0.0`.
* Running `terraform import` on Google provider `5.0.0` or later, these three fields will show in the state file with an empty `labels` and `terraform_labels` value. `effective_labels` will have all of labels present on the resource in GCP. You can update the resource to bring labels defined in your configuration under management by Terraform.
* Running `terraform refresh` on Google provider `5.0.0` or later  with an existing pre-`5.0.0` resource before an `apply`, these three fields will show in the state file. `labels` field will have your current labels, `terraform_labels` will be empty, and `effective_labels` will have all of labels present on the resource in GCP.
* Running `terraform plan` or `terraform apply` on Google provider `5.0.0` or later with an existing pre-`5.0.0` resource before an `apply`, the plan will show an updated adding your current labels to `terraform_labels`. After running `terraform apply`, these three fields will show in the state file. `labels` will have your current labels, `terraform_labels` will have the combination of `labels` and your provider-default labels, and `effective_labels` will have all of labels present on the resource in GCP.
* In the resource `google_cloud_run_domain_mapping`, the system labels `cloud.googleapis.com/location` and `run.googleapis.com/overrideAt` will be removed from `labels` inside `metadata` field in the state file as part of a one-time resource schema upgrade. If any of these label keys are in the configuration, after upgrading to `5.0.0`, the plan will show that these keys will be added. You can safely accept this change, and Terraform will begin to manage them again.
* In the resource `google_cloud_run_service`, the system label `cloud.googleapis.com/location` will be removed from `labels` inside `metadata` field in the state file as part of a one-time resource schema upgrade. If this system label key is in the configuration, after upgrading to `5.0.0`, the plan will show that this key will be added. You can safely accept this change, and Terraform will begin to manage it again.
* In the resource `google_dataflow_flex_template_job`, the system labels with the prefix `goog-dataflow-provided` are removed from `labels` field in the state file as part of a one-time resource schema upgrade. If any of these label keys are in the configuration, after upgrading to `5.0.0`, the plan will show that these keys will be added. You can safely accept this change, and Terraform will begin to manage them again.
* In the resource `google_dataflow_job`, the system labels with the prefix `goog-dataflow-provided` will be removed from `labels` field in the state file as part of a one-time resource schema upgrade. If any of these label keys are in the configuration, after upgrading to `5.0.0`, the plan will show that these keys will be added. You can safely accept this change, and Terraform will begin to manage them again.
* In the resource `google_dataproc_cluster`, the system labels with the prefix `goog-dataproc` will be removed from `labels` field in the state file as part of a one-time resource schema upgrade. If any of these label keys are in the configuration, after upgrading to `5.0.0`, the plan will show that these keys will be added. You can safely accept this change.
* In the resource `google_notebooks_instance`, the system labels with the prefix `goog-caip-notebook` will be removed from `labels` field in the state file as part of a one-time resource schema upgrade. If any of these label keys are in the configuration, after upgrading to `5.0.0`, the plan will show that these keys will be added. You can safely accept this change, and Terraform will begin to manage them again.
* In the resource `google_storage_bucket`, the system labels with the prefix `goog-dataplex` will be removed from `labels` field in the state file as part of a one-time resource schema upgrade. If any of these label keys are in the configuration, after upgrading to `5.0.0`, the plan will show that these keys will be added. You can safely accept this change, and Terraform will begin to manage them again.

#### Data source labels

For most resource-based datasources, all three of `labels`, `effective_labels` and `terraform_labels` will now be present. All of these three fields will have all of labels present on the resource in GCP including the labels configured through Terraform, the system, and other clients, equivalent to `effective_labels` on the resource.

#### Resource annotations

The new annotations model is similar to the new labels model and will be applied to all of the resources with the top level `annotations` field or the nested `annotations` field inside the top level `metadata` field.

There are now two annotation-related fields with the new model, `annotations` and the output-only `effective_annotations`.

**Note:** `ignore_changes` can be applied to `annotations` field to ignore the changes of the user defined annotations. It is not recommended to apply `ignore_changes` to `effective_annotations`, as it may unintuitively affect the final API call.

The following changes will be observed after upgrading to `5.0.0`.
* Running `terraform import` on Google provider `5.0.0` or later, these two fields will show in the state file with an empty `annotations` value. `effective_annotations` will have all of annotations present on the resource in GCP. You can update the resource to bring annotations defined in your configuration under management by Terraform.
* Running `terraform refresh` on Google provider `5.0.0` or later with an existing pre-`5.0.0` resource before an `apply`, these two fields will show in the state file. `annotations` field will have your current annotations and `effective_annotations` will have all of annotations present on the resource in GCP.
* Running `terraform apply` on Google provider `5.0.0` or later with an existing pre-`5.0.0` resource before an `apply`, these two fields will show in the state file. `annotations` field will have your current annotations and `effective_annotations` will have all of annotations present on the resource in GCP.
* In the resource `google_cloud_run_domain_mapping`, the system annotations `serving.knative.dev/creator`, `serving.knative.dev/lastModifier`, `run.googleapis.com/operation-id`, `run.googleapis.com/ingress-status`, `run.googleapis.com/ingress` with the value `all` will be removed from `annotations` inside `metadata` field in the state file as part of a one-time resource schema upgrade. If any of these annotation keys are in the configuration, after upgrading to `5.0.0`, the plan will show that these keys will be added. You can safely accept this change, and Terraform will begin to manage them again.
* In the resource `google_cloud_run_service`, the system annotations `serving.knative.dev/creator`, `serving.knative.dev/lastModifier`, `run.googleapis.com/operation-id`, `run.googleapis.com/ingress-status`, `run.googleapis.com/ingress` with the value `all` will be removed from `annotations` inside `metadata` field in the state file as part of a one-time resource schema upgrade. If any of these annotation keys are in the configuration, after upgrading to `5.0.0`, the plan will show that these keys will be added. You can safely accept this change, and Terraform will begin to manage them again.

#### Data source annotations

For most resource-based datasources, both `annotations` and `effective_annotations` will now be present. Both fields will have all of annotations present on the resource in GCP including the annotations configured through Terraform, the system, and other clients, equivalent to `effective_annotations` on the resource.

#### Example
##### Config
```hcl
provider "google" {
  default_labels = {
    default_key = "default_value"
  }
}

resource "google_dataproc_cluster" "with_labels" {
  name   = "tf-test-dproc-test-1"
  region = "us-central1"

  labels = {
    key1 = "value1"
    key2 = "value2"
  }
}
```
After the configuration is applied, Terraform is managing `key1` and `key2` in the `labels` field. `terraform_labels` field has label `default_key`, `key1` and `key2`. `effective_labels` has label `default_key`, `key1`, `key2` and system labels.

##### Config
```hcl
provider "google" {
  default_labels = {
    default_key = "default_value"
  }
}

resource "google_dataproc_cluster" "with_labels" {
  name   = "tf-test-dproc-test-1"
  region = "us-central1"

  labels = {
    key1 = "value1"
  }
}
```
After the configuration is applied, Terraform stops managing `key2` and is managing `key1` in the `labels` field. `terraform_labels` field has label `default_key` and `key1`. `effective_labels` has label `default_key`, `key1` and system labels.

### Updates to how empty strings are handled in the `provider` block

In 5.0.0+ any empty strings set in the `provider` block will be used and not ignored. Previously any empty strings used as arguments in the `provider` block were ignored and did not contribute to configuration of the provider. 

Users should remove empty string arguments to avoid errors during plan/apply stages.

```tf
provider "google" {
  credentials = "" # this line should be removed
  project = "my-project"
  region = "us-central1"
  zone = "" # this line should be removed
}
```

### Changes to how default `location`, `region` and `zone` values are obtained for resources

Currently, when configuring resources that require a `location`, `region` or `zone` field you have the choice of specifying it in the resource block or allowing default values to be used. Default [region](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/provider_reference#region) or [zone](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/provider_reference#zone) values can be configured in the provider block or by providing values through environment variables.

Changes in 5.0.0 make the way the provider handles `location`/`region`/`zone` values more consistent:

* Resources that have a `location` field will now use the default `region` value preferentially over the default `zone` value set on the provider. This is only relevant to resources where `location` is not provided in the resource block directly.
* Previously, default `region` and `zone` values set as URIs were incompatible with resources that have `location` or `region` arguments. In 5.0.0+ those values will now be valid and won't result in errors during plan/apply stages.


#### When you may need to take action

There is only one change that we anticipate can lead to unexpected diffs in Terraform plans after upgrading to 5.0.0, which is:

> Resources that have a `location` field will now use the default `region` value preferentially over the default `zone` value set on the provider. This is only relevant to resources where `location` is not provided in the resource block directly.

Users will need to check for unexpected `location` changes for resources. If an unexpected change is seen, the solution is to explicitly set the `location` value in that resource's configuration block to match the desired value.

This will only affect users whose configuration contains resource blocks that have missing `location` values and whose [default zone](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/provider_reference#zone) value belongs to a region that's different than the [default region](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/provider_reference#region) value. For example, if you set `us-central1-a` as the default zone and `us-central2` as the default region on the provider you may see plans that contain unexpected diffs to move resources from `us-central1` to `us-central2`.


### Provider default values shown at plan-time

`project`, `region`, and `zone` fields will now display their values during plan-time instead of the placeholder `(known after apply)` value normally displayed for fields without fixed Terraform default values. These values will be taken from either the Terraform resource config file, provider config, or local environment variables, depending on which variables are supplied by the user, matching the existing per-resource functionality for what default values are used in execution of a Terraform plan.

### Resource import formats have improved validation

Throughout the provider there were many resources which erroneously gave false positives to poorly formatted import input if a subset of the provided input was valid to their configured import formats. All GCP resource IDs supplied to "terraform import" must match the documentation specified import formats exactly.

## Datasources

### Datasources now error universally on 404

All non-IAM data sources have been updated to return an error when a target resource URI can not be reached. Previously this was inconsistent between different datasources in whether an empty value was returned to Terraform state upon 404 or if an error was returned, but this has been standardized. Any plans that reference datasources which no longer exist (or do not exist yet) will need to be revised to have these datasources removed from configuration files.

## Datasource: `google_product_datasource`

### Datasource-level change example header

Description of the change and how users should adjust their configuration (if needed).

## Resources

## Resource: `google_product_resource`

### Resource-level change example header

Description of the change and how users should adjust their configuration (if needed).

## Resource: `google_access_context_manager_service_perimeters`

### Fields changed from lists to sets where appropriate

The following fields have been changed from a list to a set:

`google_access_context_manager_service_perimeters.spec.access_levels`
`google_access_context_manager_service_perimeters.spec.resources`
`google_access_context_manager_service_perimeters.spec.restricted_services`
`google_access_context_manager_service_perimeters.spec.vpc_accessible_services.allowed_services`
`google_access_context_manager_service_perimeters.spec.egress_policies.egress_from.identities`
`google_access_context_manager_service_perimeters.spec.egress_policies.egress_to.external_resources`
`google_access_context_manager_service_perimeters.spec.egress_policies.egress_to.resources`
`google_access_context_manager_service_perimeters.spec.ingress_policies.ingress_from.identities`
`google_access_context_manager_service_perimeters.spec.ingress_policies.ingress_to.resources`
`google_access_context_manager_service_perimeters.status.access_levels`
`google_access_context_manager_service_perimeters.status.resources`
`google_access_context_manager_service_perimeters.status.restricted_services`
`google_access_context_manager_service_perimeters.status.egress_policies.egress_from.identities`
`google_access_context_manager_service_perimeters.status.egress_policies.egress_to.external_resources`
`google_access_context_manager_service_perimeters.status.egress_policies.egress_to.resources`
`google_access_context_manager_service_perimeters.status.ingress_policies.ingress_from.identities`
`google_access_context_manager_service_perimeters.status.ingress_policies.ingress_to.resources`

If you were relying on accessing an individual field by index (for example, google_access_context_manager_service_perimeters.spec.access_levels.0, then that will now need to by hash (for example, google_access_context_manager_service_perimeters.spec.access_levels.\<some-hash\>).

## Resource: `google_access_context_manager_service_perimeter`

### Fields changed from lists to sets where appropriate

The following fields have been changed from a list to a set:

`google_access_context_manager_service_perimeter.spec.access_levels`
`google_access_context_manager_service_perimeter.spec.resources`
`google_access_context_manager_service_perimeter.spec.restricted_services`
`google_access_context_manager_service_perimeter.spec.vpc_accessible_services.allowed_services`
`google_access_context_manager_service_perimeter.spec.egress_policies.egress_from.identities`
`google_access_context_manager_service_perimeter.spec.egress_policies.egress_to.external_resources`
`google_access_context_manager_service_perimeter.spec.egress_policies.egress_to.resources`
`google_access_context_manager_service_perimeter.spec.ingress_policies.ingress_from.identities`
`google_access_context_manager_service_perimeter.spec.ingress_policies.ingress_to.resources`
`google_access_context_manager_service_perimeter.status.access_levels`
`google_access_context_manager_service_perimeter.status.resources`
`google_access_context_manager_service_perimeter.status.restricted_services`
`google_access_context_manager_service_perimeter.status.egress_policies.egress_from.identities`
`google_access_context_manager_service_perimeter.status.egress_policies.egress_to.external_resources`
`google_access_context_manager_service_perimeter.status.egress_policies.egress_to.resources`
`google_access_context_manager_service_perimeter.status.ingress_policies.ingress_from.identities`
`google_access_context_manager_service_perimeter.status.ingress_policies.ingress_to.resources`

If you were relying on accessing an individual field by index (for example, google_access_context_manager_service_perimeter.spec.access_levels.0, then that will now need to by hash (for example, google_access_context_manager_service_perimeter.spec.access_levels.\<some-hash\>).


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


## Resource: `google_container_cluster`

### Clusters created in error states are now tainted rather than deleted

GKE clusters that are created but do not become healthy will now be recorded in
state and marked as tainted for cleanup on next apply rather than immediately
deleted.

This behavior was changed to allow users to collect internal logs from the
cluster and/or manually resolve the issues and untaint their failed clusters.

### `node_config.0.taint` and `node_pool.0.node_config.0.taint` field change

The `taint` field has been changed to manage a subset of the taint keys on a node pool
and the `effective_taints` output field has been added to record the complete set of
taints applied to the node pool by GKE.

Previously, the field was authoritative and would require every taint on the node pool
to be recorded, causing friction when users used GPUs or configured sandbox settings,
actions which added taints. After this change, only "Terraform-managed" taints will be
managed by the `taint` field. Other taints, including new taints injected by the
server, will not have drift detected.

Currently, the set of managed taints and their values are immutable in Terraform, and
any changes will cause a recreate to be planned. However, taints can be unmanaged by
simultaneously removing the taint entry from GKE and your Terraform configuration at
the same time.

The set of taints Terraform manages (and their values) will be determined based on
how the cluster or node pool resource was added to your Terraform state file:

* If you created the cluster or node pool with Terraform with Google provider 5.0.0
or later, the set of taints specified during resource creation will be managed.
* If you imported the cluster or node pool with Google provider 5.0.0 or later, no
taints will be managed by Terraform
* If you upgraded from an earlier version, the complete set of taint values applied to the
node pool at the time of your last refresh will be managed by Terraform

Most existing configurations will not be affected with this change as they already specify
the whole set of managed taints, or are already ignoring changes with `lifecycle.ignore_changes`,
preventing a diff.

A limited number of users may see a diff if they are using the `google-beta` provider
and have specified a `sandbox_config` value. If that's the case, you can safely add the
proposed value to configuration (below) or apply `lifecycle.ignore_changes` to the field to resolve.


```diff
+    taint {
+      key    = "sandbox.gke.io/runtime"
+      value  = "gvisor"
+      effect = "NO_SCHEDULE"
+    }
```

### Cluster deletion now prevented by default with `deletion_protection`

The field `deletion_protection` has been added with a default of `true`. This field prevents
Terraform from destroying or recreating the cluster.
**`deletion_protection` does NOT prevent deletion outside of Terraform.**
To destroy a `google_container_cluster`, this field must be explicitly set to `false`.

### `networking_mode` defaults to `VPC_NATIVE` for newly created clusters

New clusters will default to `VPC_NATIVE` which enables [IP aliasing](https://cloud.google.com/kubernetes-engine/docs/how-to/ip-aliases). Previously, `google_container_cluster` would default to using routes as
the networking mode unless `ip_allocation_policy` policy was set. Now, `networking_mode` will
default to `VPC_NATIVE` and `ip_allocation_policy` will be set by the server if unset in
configuration. Existing clusters should not be affected.

#### New Minimal Config for VPC-native cluster

```hcl
resource "google_container_cluster" "primary" {
  name               = "my_cluster"
  location           = "us-central1-a"
  initial_node_count = 1
}
```

#### New Minimal Config for Routes-based cluster

```hcl
resource "google_container_cluster" "primary" {
  name               = "my_cluster"
  location           = "us-central1-a"
  initial_node_count = 1
  networking_mode    = "ROUTES"
}
```

### `enable_binary_authorization` is now removed

`enable_binary_authorization` has been removed in favor of `binary_authorization.enabled`.

### Default value of `network_policy.provider` is now removed

Previously `network_policy.provider` defaulted to "PROVIDER_UNSPECIFIED". It no longer
has a default value.

## Resource: `google_container_node_pool`

### `management.auto_repair` and `management.auto_upgrade` now default to true

Previously both fields defaulted to false if `management` was set as an empty block.
They now default to true.

### `logging_variant` no longer has a provider default value

The provider no longer shows a default value at plan-time, and defers to the 
API server default.

### `node_config.0.taint` field change

The `taint` field has been changed to manage a subset of the taint keys on a node pool
and the `effective_taints` output field has been added to record the complete set of
taints applied to the node pool by GKE.

Previously, the field was authoritative and would require every taint on the node pool
to be recorded, causing friction when users used GPUs or configured sandbox settings,
actions which added taints. After this change, only "Terraform-managed" taints will be
managed by the `taint` field. Other taints, including new taints injected by the
server, will not have drift detected.

Currently, the set of managed taints and their values are immutable in Terraform, and
any changes will cause a recreate to be planned. However, taints can be unmanaged by
simultaneously removing the taint entry from GKE and your Terraform configuration at
the same time.

The set of taints Terraform manages (and their values) will be determined based on
how the cluster or node pool resource was added to your Terraform state file:

* If you created the cluster or node pool with Terraform with Google provider 5.0.0
or later, the set of taints specified during resource creation will be managed.
* If you imported the cluster or node pool with Google provider 5.0.0 or later, no
taints will be managed by Terraform
* If you upgraded from an earlier version, the complete set of taint values applied to the
node pool at the time of your last refresh will be managed by Terraform

Most existing configurations will not be affected with this change as they already specify
the whole set of managed taints, or are already ignoring changes with `lifecycle.ignore_changes`,
preventing a diff.

A limited number of users may see a diff if they are using the `google-beta` provider
and have specified a `sandbox_config` value. If that's the case, you can safely add the
proposed value to configuration (below) or apply `lifecycle.ignore_changes` to the field to resolve.


```diff
+    taint {
+      key    = "sandbox.gke.io/runtime"
+      value  = "gvisor"
+      effect = "NO_SCHEDULE"
+    }
```

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

## Resource: `google_monitoring_dashboard`

### `dashboard_json` suppresses removal diffs more aggressively

To prevent permanent diffs from default values, Terraform will now attempt to suppress diffs where the value is returned in the JSON
string but doesn't exist in the configuration. Consequently, legitmate remove-only diffs will also be suppressed.
For Terraform to detect the diff, JSON key removals must also be accompanied by a non-removal change (trivial or not).

## Resource: `google_monitoring_metric_descriptor`

### Changing `labels` now triggers replacement

Previously, attempting to update `labels` failed and created a permadiff. The `labels` 
field is now immutable without destroying and recreating the resource.

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

## Resource: `google_compute_security_policy`

### `rule.rate_limit_options.encorce_on_key` no longer has default value

Previously, the default value for `rule.rate_limit_options.encorce_on_key` is "ALL", now this field no longer has a default value.

## Resource: `google_logging_metric`

### Additional `bucket_options` subfields are now properly required

When setting the `bucket_options` block, the following fields may be required:

* `num_finite_buckets`, `width`, and `offset` are now required when `bucket_options.linear_buckets` is set.

* `num_finite_buckets`, `growth_factor`, and `scale` are now required when `bucket_options.exponential_buckets` is set.

Previously these fields should have been required but were not, which allowed for invalid `google_logging_metric` configurations.

## Resource: `google_logging_project_sink`

### `unique_writer_identity` now defaults to `TRUE`

Previously, the default value of `unique_writer_identity` was `FALSE`. Now it will be `TRUE`.

This will change the behavior for new sinks created using the default value. Previously, all sinks created using the default value had a `writer_identity` of `serviceAccount:cloud-logs@system.gserviceaccount.com`. Now sinks created using the default value will have a `writer_identity` that differs depending on the parent resource, for example: `serviceAccount:service-<PROJECT_NUMBER>@gcp-sa-logging.iam.gserviceaccount.com` for a project-level sink.

IAM permissions that were manually configured for `cloud-logs@system.gserviceaccount.com` and `iam_bindings` that are hard-coded to use `cloud-logs@system.gserviceaccount.com` will not properly apply permissions to the `writer_identity` of new sinks created using the default value.  **If a sink is missing the proper permissions it will be successfully created but it will fail to export log data.**

Currently there are only two types of log sinks that populate `writer_identity` and can be created with `unique_writer_identity = false`.  Only these types of sinks may be affected:
* Sinks with a Cloud Pub/Sub topic `destination` for which the topic is in the same project as the sink.
* Sinks for a BigQuery dataset `destination` for which the dataset is in the same project as the sink.

To ensure that proper permissions are in place for new sinks created using the default value, check that the related `iam_bindings` are configured and reference the sink's `writer_identity` property.

Here is an example of proper `iam_bindings`:

```hcl
resource "google_logging_project_sink" "gcs-bucket-sink" {
  name        = "my-gcs-bucket-sink"
  description = "Routes all admin activity logs to a GCS bucket"
  destination = "storage.googleapis.com/${google_storage_bucket.log-bucket.name}"
  filter      = "log_id(\"cloudaudit.googleapis.com/activity\")"
  # `unique_writer_identity is explicitly set to true here, but will now default to 'true'.
  unique_writer_identity = true
}

# We must grant proper permissions for the log sink to access the GCS bucket.
resource "google_project_iam_binding" "gcs-bucket-writer" {
  project = "your-project-id"
  role = "roles/storage.objectCreator"

  members = [
    google_logging_project_sink.gcs-bucket-sink.writer_identity,
  ]
}
```

## Resource: `google_cloudfunctions2_function`
### `location` now a required field
Deployment would fail if this field was unspecified. Marked this field as requied to align with implementation. This value cannot be inferred from any provider level config. No change is necessary for upgrade as this field is already needed for any deployments.

## Resource: `google_cloud_run_v2_service`
### transitioned `volumes.cloud_sql_instance.instances` to SET from ARRAY for `google_cloud_run_v2_service`
Previously, `database_flags` was a list, making it order-dependent. It is now a set.

If you were relying on accessing an individual flag by index (for example, `google_sql_database_instance.instance.settings.0.database_flags.0.name`), then that will now need to by hash (for example, `google_sql_database_instance.instance.settings.0.database_flags.<some-hash>.name`).

## Product: `cloudiot`

### resource `google_cloudiot_device` is now removed

### resource `google_cloudiot_registry` is now removed

### resource `google_cloudiot_registry_iam_*` is now removed

### datasource `google_cloudiot_registry_iam_policy` is now removed

## Resource: `google_service_networking_connection`

### `Create` endpoint is used to create the resource

`google_service_networking_connection` now uses the Create endpoint instead of the Patch endpoint during the creation step. Previously, Patch was used as a workaround for an issue that has since been resolved.

### "terraform destroy" now fully deletes the resource instead of abandoning

`google_service_networking_connection` now uses API `deleteConnection` method instead of `removePeering` method during the deletion step. Previously, `removePeering` method was used because `deleteConnection` method was unavailable. In some cases a private connection cannot be deleted immediately after the resource using that connection is deleted, and users may have to delete the private connection after a waiting period.

## Resource: `google_secret_manager_secret`

### `replication.automatic` is now removed

Deprecated in favor of field `replication.auto`. It is now removed.

#### Old Config

```hcl
resource "google_secret_manager_secret" "my-secret" {
  secret_id = "tf-secret"
  
  replication {
    automatic = true
  }
}
```

#### New Config

```hcl
resource "google_secret_manager_secret" "my-secret" {
  secret_id = "tf-secret"
  
  replication {
    auto {}
  }
}
```

## Resource: `google_compute_service_attachment`

### `reconcile_connections` now defaults from API

`reconcile_connections` previously defaults to true. Now it will default from the API.

### Retyped `consumer_accept_lists` to SET from ARRAY

Previously, `consumer_accept_lists` was a list, making it order-dependent. It is now a set.

If you were relying on accessing an individual flag by index (for example, `google_compute_service_attachment.consumer_accept_lists.0.project_id_or_num`), then that will now need to by hash (for example, `google_compute_service_attachment.consumer_accept_lists.<some-hash>.project_id_or_num`).

## Resource: `google_dataflow_flex_template_job`

### Fields that are a part of the [environment block](https://cloud.google.com/dataflow/docs/reference/rest/v1b3/projects.locations.flexTemplates/launch#FlexTemplateRuntimeEnvironment) will be overriden to be sent via their fields even when supplied via parameters.

Several fields within the `google_dataflow_flex_template_job` resource can be supplied through either the `parameters{}` block or a field on the resource object. Support for these fields on the resource object was added in the `4.66.0` release of the Google provider. That version introduced an issue where the values were being double-sent to the API due to being recorded in Terraform state in two places. To resolve this issue, these fields will be deduplicated and sent to the API through the resource object.

Additionally, due to the API returning these fields to the user they will now be considered computed and users will see values twice within their state when configuring these fields' values via the `parameters{}` block.

## Resource: `google_compute_node_group`

### Node groups are now mutable

Due to limitations in previous field configurations, the only field that could be updated previously was `node_template`. It is now possible to adjust the `autoscaling_policy` without recreating the group, nor will any adjustment to the `size` of the nodepool prompt resource recration.

### `size` is now an output only field.

`size` previously served as an alias for `initial_size` on resource creation, and users would be required to recreate the resource if the `size` value ever adjusted due to either direct user update or auto-scaling adjustment outside of Terraform.

It will now mirror its API functionality and serve as an output only field to show how many nodes currently exist within the resource. All existing configurations which used `size` as an input field must be updated for its removal.

### One of `initial_size` or `autoscaling_policy{}` must be configured on resource creation.

These fields will supply the base node-count for a node group and one of them will be required for successful resource creation. Both will be freely updateable or removable on future state changes that do not require recreation.

## Resource: `google_looker_instance`

### `LOOKER_MODELER` has been removed as a platform edition.

Looker Modeler edition is deprecated as a platform edition.

Deprecated in favor of field `pem_certificate_chain`. It is now removed.

## Resource: `google_gkeonprem_bare_metal_admin_cluster`

Delete operation is disabled. The command `terraform destroy` maps to no-op. Users need to delete resource manually. Please refer to the [user guide](https://cloud.google.com/anthos/clusters/docs/bare-metal/latest/how-to/reset-nodes) for the instructions of cluster deletion.
