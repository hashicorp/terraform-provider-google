---
layout: "google"
page_title: "Terraform Google Provider 3.0.0 Upgrade Guide"
sidebar_current: "docs-google-provider-version-3-upgrade"
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
- [Resource: `google_container_cluster`](#resource-google_container_cluster)
- [Resource: `google_project_service`](#resource-google_project_service)
- [Resource: `google_project_services`](#resource-google_project_services)
- [Resource: `google_pubsub_subscription`](#resource-google_pubsub_subscription)

<!-- /TOC -->

## Provider Version Configuration

-> Before upgrading to version 3.0.0, it is recommended to upgrade to the most
recent `2.X` series release of the provider and ensure that your environment
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
  network    = "${google_compute_network.container_network.name}"

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
  network       = "${google_compute_network.container_network.self_link}"
}

resource "google_container_cluster" "primary" {
  name       = "my-cluster"
  location   = "us-central1"
  network    = "${google_compute_network.container_network.name}"
  subnetwork = "${google_compute_subnetwork.container_subnetwork.name}"

  initial_node_count = 1

  ip_allocation_policy {
    use_ip_aliases           = true
    cluster_ipv4_cidr_block  = "10.0.0.0/16"
    services_ipv4_cidr_block = "10.1.0.0/16"
  }
}
```

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
requests. From `2.13.0` onwards, those requests are batched. It's recommended
that you upgrade to `2.13.0+` before migrating if you encounter quota issues
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
resource "google_project_service" "project_iam" {
  project = "your-project-id"
  service = "iam.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "project_cloudresourcemanager" {
  project = "your-project-id"
  service = "cloudresourcemanager.googleapis.com"
  disable_on_destroy = false
}
```

## Resource: `google_pubsub_subscription`

### `name` must now be a short name

`name` previously could have been specified by a long name (e.g. `projects/my-project/subscriptions/my-subscription`)
or a shortname (e.g. `my-subscription`). `name` now must be the shortname.
