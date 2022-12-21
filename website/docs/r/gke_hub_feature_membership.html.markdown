---
subcategory: "GKEHub"
page_title: "Google: google_gke_hub_feature_membership"
description: |-
  Contains information about a GKEHub Feature Memberships.
---

# google\_gkehub\_feature\_membership

Contains information about a GKEHub Feature Memberships. Feature Memberships configure GKEHub Features that apply to specific memberships rather than the project as a whole. The google_gke_hub is the Fleet API.

~> **Warning:** This resource is in beta, and should be used with the terraform-provider-google-beta provider.
See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html) for more details on beta resources.


## Example Usage - Config Management

```hcl
resource "google_container_cluster" "cluster" {
  name               = "my-cluster"
  location           = "us-central1-a"
  initial_node_count = 1
  provider = google-beta
}

resource "google_gke_hub_membership" "membership" {
  membership_id = "my-membership"
  endpoint {
    gke_cluster {
      resource_link = "//container.googleapis.com/${google_container_cluster.cluster.id}"
    }
  }
  provider = google-beta
}

resource "google_gke_hub_feature" "feature" {
  name = "configmanagement"
  location = "global"

  labels = {
    foo = "bar"
  }
  provider = google-beta
}

resource "google_gke_hub_feature_membership" "feature_member" {
  location = "global"
  feature = google_gke_hub_feature.feature.name
  membership = google_gke_hub_membership.membership.membership_id
  configmanagement {
    version = "1.6.2"
    config_sync {
      git {
        sync_repo = "https://github.com/hashicorp/terraform"
      }
    }
  }
  provider = google-beta
}
```
## Example Usage - Config Management with OCI

```hcl
resource "google_container_cluster" "cluster" {
  name               = "my-cluster"
  location           = "us-central1-a"
  initial_node_count = 1
  provider = google-beta
}

resource "google_gke_hub_membership" "membership" {
  membership_id = "my-membership"
  endpoint {
    gke_cluster {
      resource_link = "//container.googleapis.com/${google_container_cluster.cluster.id}"
    }
  }
  provider = google-beta
}

resource "google_gke_hub_feature" "feature" {
  name = "configmanagement"
  location = "global"

  labels = {
    foo = "bar"
  }
  provider = google-beta
}

resource "google_gke_hub_feature_membership" "feature_member" {
  location = "global"
  feature = google_gke_hub_feature.feature.name
  membership = google_gke_hub_membership.membership.membership_id
  configmanagement {
    version = "1.12.0"
    config_sync {
      oci {
        sync_repo = "us-central1-docker.pkg.dev/sample-project/config-repo/config-sync-gke:latest"
        policy_dir = "config-connector"
        sync_wait_secs = "20"
        secret_type = "gcpserviceaccount"
        gcp_service_account_email = "sa@project-id.iam.gserviceaccount.com"
      }
    }
  }
  provider = google-beta
}
```

## Example Usage - Multi Cluster Service Discovery

```hcl
resource "google_gke_hub_feature" "feature" {
  name = "multiclusterservicediscovery"
  location = "global"
  labels = {
    foo = "bar"
  }
  provider = google-beta
}
```

## Example Usage - Service Mesh

```hcl
resource "google_container_cluster" "cluster" {
  name               = "my-cluster"
  location           = "us-central1-a"
  initial_node_count = 1
  provider = google-beta
}

resource "google_gke_hub_membership" "membership" {
  membership_id = "my-membership"
  endpoint {
    gke_cluster {
      resource_link = "//container.googleapis.com/${google_container_cluster.cluster.id}"
    }
  }
  provider = google-beta
}

resource "google_gke_hub_feature" "feature" {
  name = "servicemesh"
  location = "global"

  provider = google-beta
}

resource "google_gke_hub_feature_membership" "feature_member" {
  location = "global"
  feature = google_gke_hub_feature.feature.name
  membership = google_gke_hub_membership.membership.membership_id
  mesh {
    management = "MANAGEMENT_AUTOMATIC"
  }
  provider = google-beta
}
```

## Argument Reference

The following arguments are supported:

- - -

* `configmanagement` -
  (Optional)
  Config Management-specific spec. Structure is [documented below](#nested_configmanagement).

* `mesh` -
  (Optional)
  Service mesh specific spec. Structure is [documented below](#nested_mesh).
  
* `feature` -
  (Optional)
  The name of the feature
  
* `location` -
  (Optional)
  The location of the feature
  
* `membership` -
  (Optional)
  The name of the membership
  
* `project` -
  (Optional)
  The project of the feature
  


<a name="nested_configmanagement"></a>The `configmanagement` block supports:
    
* `binauthz` -
  (Optional)
  Binauthz configuration for the cluster. Structure is [documented below](#nested_binauthz).
    
* `config_sync` -
  (Optional)
  Config Sync configuration for the cluster. Structure is [documented below](#nested_config_sync).
    
* `hierarchy_controller` -
  (Optional)
  Hierarchy Controller configuration for the cluster. Structure is [documented below](#nested_hierarchy_controller).
    
* `policy_controller` -
  (Optional)
  Policy Controller configuration for the cluster. Structure is [documented below](#nested_policy_controller).
    
* `version` -
  (Optional)
  Version of ACM installed.
    
<a name="nested_binauthz"></a>The `binauthz` block supports:
    
* `enabled` -
  (Optional)
  Whether binauthz is enabled in this cluster.
    
<a name="nested_config_sync"></a>The `config_sync` block supports:
    
* `git` -
  (Optional) Structure is [documented below](#nested_git).

* `oci` -
  (Optional) Supported from ACM versions 1.12.0 onwards. Structure is [documented below](#nested_oci).
  
  Use either `git` or `oci` config option.

* `prevent_drift` -
  (Optional)
  Supported from ACM versions 1.10.0 onwards. Set to true to enable the Config Sync admission webhook to prevent drifts. If set to "false", disables the Config Sync admission webhook and does not prevent drifts.
    
* `source_format` -
  (Optional)
  Specifies whether the Config Sync Repo is in "hierarchical" or "unstructured" mode.
    
<a name="nested_git"></a>The `git` block supports:
    
* `gcp_service_account_email` -
  (Optional)
  The GCP Service Account Email used for auth when secretType is gcpServiceAccount.

* `https_proxy` -
  (Optional)
  URL for the HTTPS proxy to be used when communicating with the Git repo.
    
* `policy_dir` -
  (Optional)
  The path within the Git repository that represents the top level of the repo to sync. Default: the root directory of the repository.
    
* `secret_type` -
  (Optional)
  Type of secret configured for access to the Git repo.
    
* `sync_branch` -
  (Optional)
  The branch of the repository to sync from. Default: master.
    
* `sync_repo` -
  (Optional)
  The URL of the Git repository to use as the source of truth.
    
* `sync_rev` -
  (Optional)
  Git revision (tag or hash) to check out. Default HEAD.
    
* `sync_wait_secs` -
  (Optional)
  Period in seconds between consecutive syncs. Default: 15.

<a name="nested_oci"></a>The `oci` block supports:
    
* `gcp_service_account_email` -
  (Optional)
  The GCP Service Account Email used for auth when secret_type is gcpserviceaccount.
    
* `policy_dir` -
  (Optional)
  The absolute path of the directory that contains the local resources. Default: the root directory of the image.
    
* `secret_type` -
  (Optional)
  Type of secret configured for access to the OCI Image. Must be one of gcenode, gcpserviceaccount or none.
    
* `sync_repo` -
  (Optional)
  The OCI image repository URL for the package to sync from. e.g. LOCATION-docker.pkg.dev/PROJECT_ID/REPOSITORY_NAME/PACKAGE_NAME.
    
* `sync_wait_secs` -
  (Optional)
  Period in seconds(int64 format) between consecutive syncs. Default: 15.
    
<a name="nested_hierarchy_controller"></a>The `hierarchy_controller` block supports:
    
* `enable_hierarchical_resource_quota` -
  (Optional)
  Whether hierarchical resource quota is enabled in this cluster.
    
* `enable_pod_tree_labels` -
  (Optional)
  Whether pod tree labels are enabled in this cluster.
    
* `enabled` -
  (Optional)
  Whether Hierarchy Controller is enabled in this cluster.
    
<a name="nested_policy_controller"></a>The `policy_controller` block supports:
    
* `audit_interval_seconds` -
  (Optional)
  Sets the interval for Policy Controller Audit Scans (in seconds). When set to 0, this disables audit functionality altogether.
    
* `enabled` -
  (Optional)
  Enables the installation of Policy Controller. If false, the rest of PolicyController fields take no effect.
    
* `exemptable_namespaces` -
  (Optional)
  The set of namespaces that are excluded from Policy Controller checks. Namespaces do not need to currently exist on the cluster.
    
* `log_denies_enabled` -
  (Optional)
  Logs all denies and dry run failures.
    
* `referential_rules_enabled` -
  (Optional)
  Enables the ability to use Constraint Templates that reference to objects other than the object currently being evaluated.
    
* `template_library_installed` -
  (Optional)
  Installs the default template library along with Policy Controller.

* `mutation_enabled` -
  (Optional)
  Enables mutation in policy controller. If true, mutation CRDs, webhook, and controller deployment will be deployed to the cluster.

* `monitoring` -
  (Optional)
  Specifies the backends Policy Controller should export metrics to. For example, to specify metrics should be exported to Cloud Monitoring and Prometheus, specify backends: [\"cloudmonitoring\", \"prometheus\"]. Default: [\"cloudmonitoring\", \"prometheus\"]    

<a name="nested_mesh"></a>The `mesh` block supports:

* `management` -
  (Optional)
  Whether to automatically manage Service Mesh. Can either be `MANAGEMENT_AUTOMATIC` or `MANAGEMENT_MANUAL`.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/locations/{{location}}/features/{{feature}}/membershipId/{{membership}}`

## Timeouts

This resource provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 20 minutes.
- `update` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import

FeatureMembership can be imported using any of these accepted formats:

```
$ terraform import google_gke_hub_feature_membership.default projects/{{project}}/locations/{{location}}/features/{{feature}}/membershipId/{{membership}}
$ terraform import google_gke_hub_feature_membership.default {{project}}/{{location}}/{{feature}}/{{membership}}
$ terraform import google_gke_hub_feature_membership.default {{location}}/{{feature}}/{{membership}}
```
