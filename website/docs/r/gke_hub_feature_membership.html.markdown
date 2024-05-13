---
subcategory: "GKEHub"
description: |-
  Contains information about a GKEHub Feature Memberships.
---

# google_gkehub_feature_membership

Contains information about a GKEHub Feature Memberships. Feature Memberships configure GKEHub Features that apply to specific memberships rather than the project as a whole. The google_gke_hub is the Fleet API.

## Example Usage - Config Management

```hcl
resource "google_container_cluster" "cluster" {
  name               = "my-cluster"
  location           = "us-central1-a"
  initial_node_count = 1
}

resource "google_gke_hub_membership" "membership" {
  membership_id = "my-membership"
  endpoint {
    gke_cluster {
      resource_link = "//container.googleapis.com/${google_container_cluster.cluster.id}"
    }
  }
}

resource "google_gke_hub_feature" "feature" {
  name = "configmanagement"
  location = "global"

  labels = {
    foo = "bar"
  }
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
}
```
## Example Usage - Config Management with OCI

```hcl
resource "google_container_cluster" "cluster" {
  name               = "my-cluster"
  location           = "us-central1-a"
  initial_node_count = 1
}

resource "google_gke_hub_membership" "membership" {
  membership_id = "my-membership"
  endpoint {
    gke_cluster {
      resource_link = "//container.googleapis.com/${google_container_cluster.cluster.id}"
    }
  }
}

resource "google_gke_hub_feature" "feature" {
  name = "configmanagement"
  location = "global"

  labels = {
    foo = "bar"
  }
}

resource "google_gke_hub_feature_membership" "feature_member" {
  location = "global"
  feature = google_gke_hub_feature.feature.name
  membership = google_gke_hub_membership.membership.membership_id
  configmanagement {
    version = "1.15.1"
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
}
```

## Example Usage - Service Mesh

```hcl
resource "google_container_cluster" "cluster" {
  name               = "my-cluster"
  location           = "us-central1-a"
  initial_node_count = 1
}

resource "google_gke_hub_membership" "membership" {
  membership_id = "my-membership"
  endpoint {
    gke_cluster {
      resource_link = "//container.googleapis.com/${google_container_cluster.cluster.id}"
    }
  }
}

resource "google_gke_hub_feature" "feature" {
  name = "servicemesh"
  location = "global"
}

resource "google_gke_hub_feature_membership" "feature_member" {
  location = "global"
  feature = google_gke_hub_feature.feature.name
  membership = google_gke_hub_membership.membership.membership_id
  mesh {
    management = "MANAGEMENT_AUTOMATIC"
  }
}
```

## Example Usage - Config Management with Regional Membership

```hcl
resource "google_container_cluster" "cluster" {
  name               = "my-cluster"
  location           = "us-central1-a"
  initial_node_count = 1
}

resource "google_gke_hub_membership" "membership" {
  membership_id = "my-membership"
  location      = "us-central1"
  endpoint {
    gke_cluster {
      resource_link = "//container.googleapis.com/${google_container_cluster.cluster.id}"
    }
  }
}

resource "google_gke_hub_feature" "feature" {
  name = "configmanagement"
  location = "global"

  labels = {
    foo = "bar"
  }
}

resource "google_gke_hub_feature_membership" "feature_member" {
  location = "global"
  feature = google_gke_hub_feature.feature.name
  membership = google_gke_hub_membership.membership.membership_id
  membership_location = google_gke_hub_membership.membership.location
  configmanagement {
    version = "1.6.2"
    config_sync {
      git {
        sync_repo = "https://github.com/hashicorp/terraform"
      }
    }
  }
}
```

## Example Usage - Policy Controller with minimal configuration

```hcl
resource "google_container_cluster" "cluster" {
  name               = "my-cluster"
  location           = "us-central1-a"
  initial_node_count = 1
}

resource "google_gke_hub_membership" "membership" {
  membership_id = "my-membership"
  endpoint {
    gke_cluster {
      resource_link = "//container.googleapis.com/${google_container_cluster.cluster.id}"
    }
  }
}

resource "google_gke_hub_feature" "feature" {
  name = "policycontroller"
  location = "global"
}

resource "google_gke_hub_feature_membership" "feature_member" {
  location = "global"
  feature = google_gke_hub_feature.feature.name
  membership = google_gke_hub_membership.membership.membership_id
  policycontroller {
    policy_controller_hub_config {
      install_spec = "INSTALL_SPEC_ENABLED"
    }
  }
}
```

## Example Usage - Policy Controller with custom configurations

```hcl
resource "google_container_cluster" "cluster" {
  name               = "my-cluster"
  location           = "us-central1-a"
  initial_node_count = 1
}

resource "google_gke_hub_membership" "membership" {
  membership_id = "my-membership"
  endpoint {
    gke_cluster {
      resource_link = "//container.googleapis.com/${google_container_cluster.cluster.id}"
    }
  }
}

resource "google_gke_hub_feature" "feature" {
  name = "policycontroller"
  location = "global"
}

resource "google_gke_hub_feature_membership" "feature_member" {
  location = "global"
  feature = google_gke_hub_feature.feature.name
  membership = google_gke_hub_membership.membership.membership_id
  policycontroller {
    policy_controller_hub_config {
      install_spec = "INSTALL_SPEC_SUSPENDED"
      policy_content {
        template_library {
          installation = "NOT_INSTALLED"
        }
      }
      constraint_violation_limit = 50
      audit_interval_seconds = 120
      referential_rules_enabled = true
      log_denies_enabled = true
      mutation_enabled = true
    }
    version = "1.17.0"
  }
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

* `policycontroller` -
  (Optional)
  Policy Controller-specific spec. Structure is [documented below](#nested_policycontroller).
  
* `feature` -
  (Optional)
  The name of the feature
  
* `location` -
  (Optional)
  The location of the feature
  
* `membership` -
  (Optional)
  The name of the membership  

* `membership_location` -
  (Optional)
  The location of the membership, for example, "us-central1". Default is "global".
  
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
  Specifies the backends Policy Controller should export metrics to. For example, to specify metrics should be exported to Cloud Monitoring and Prometheus, specify backends: ["cloudmonitoring", "prometheus"]. Default: ["cloudmonitoring", "prometheus"]    

<a name="nested_mesh"></a>The `mesh` block supports:

* `management` -
  (Optional)
  Whether to automatically manage Service Mesh. Can either be `MANAGEMENT_AUTOMATIC` or `MANAGEMENT_MANUAL`.

<a name="nested_policycontroller"></a>The `policycontroller` block supports:

* `version` -
  (Optional)
  Version of Policy Controller to install. Defaults to the latest version.

* `policy_controller_hub_config` -
  Policy Controller configuration for the cluster. Structure is [documented below](#nested_policy_controller_hub_config).

<a name="nested_policy_controller_hub_config"></a>The `policy_controller_hub_config` block supports:

* `install_spec` -
  (Optional)
  Configures the mode of the Policy Controller installation. Must be one of `INSTALL_SPEC_NOT_INSTALLED`, `INSTALL_SPEC_ENABLED`, `INSTALL_SPEC_SUSPENDED` or `INSTALL_SPEC_DETACHED`.

* `exemptable_namespaces` -
  (Optional)
  The set of namespaces that are excluded from Policy Controller checks. Namespaces do not need to currently exist on the cluster.

* `referential_rules_enabled` -
  (Optional)
  Enables the ability to use Constraint Templates that reference to objects other than the object currently being evaluated.

* `log_denies_enabled` -
  (Optional)
  Logs all denies and dry run failures.

* `mutation_enabled` -
  (Optional)
  Enables mutation in policy controller. If true, mutation CRDs, webhook, and controller deployment will be deployed to the cluster.

* `monitoring` -
  (Optional)
  Specifies the backends Policy Controller should export metrics to. Structure is [documented below](#nested_monitoring).

* `audit_interval_seconds` -
  (Optional)
  Sets the interval for Policy Controller Audit Scans (in seconds). When set to 0, this disables audit functionality altogether.

* `constraint_violation_limit` -
  (Optional)
  The maximum number of audit violations to be stored in a constraint. If not set, the  default of 20 will be used.

  * `deployment_configs` -
  (Optional)
  Map of deployment configs to deployments ("admission", "audit", "mutation").

* `policy_content` -
  (Optional)
  Specifies the desired policy content on the cluster. Structure is [documented below](#nested_policy_content).

<a name="nested_monitoring"></a>The `monitoring` block supports:

* `backends`
  (Optional)
  Specifies the list of backends Policy Controller will export to. Must be one of `CLOUD_MONITORING` or `PROMETHEUS`. Defaults to [`CLOUD_MONITORING`, `PROMETHEUS`]. Specifying an empty value `[]` disables metrics export.

<a name="nested_deployment_configs"></a>The `deployment_configs` block supports:
    
* `component_name` -
  (Required)
  The name of the component. One of `admission` `audit` or `mutation`
    
* `container_resources` -
  (Optional)
  Container resource requirements.
    
* `pod_affinity` -
  (Optional)
  Pod affinity configuration. Possible values: AFFINITY_UNSPECIFIED, NO_AFFINITY, ANTI_AFFINITY
    
* `pod_tolerations` -
  (Optional)
  Pod tolerations of node taints.
    
* `replica_count` -
  (Optional)
  Pod replica count.
    
<a name="nested_container_resources"></a>The `container_resources` block supports:
    
* `limits` -
  (Optional)
  Limits describes the maximum amount of compute resources allowed for use by the running container.
    
* `requests` -
  (Optional)
  Requests describes the amount of compute resources reserved for the container by the kube-scheduler.
    
<a name="nested_limits"></a>The `limits` block supports:
    
* `cpu` -
  (Optional)
  CPU requirement expressed in Kubernetes resource units.
    
* `memory` -
  (Optional)
  Memory requirement expressed in Kubernetes resource units.
    
<a name="nested_requests"></a>The `requests` block supports:
    
* `cpu` -
  (Optional)
  CPU requirement expressed in Kubernetes resource units.
    
* `memory` -
  (Optional)
  Memory requirement expressed in Kubernetes resource units.
    
<a name="nested_pod_tolerations"></a>The `pod_tolerations` block supports:
    
* `effect` -
  (Optional)
  Matches a taint effect.
    
* `key` -
  (Optional)
  Matches a taint key (not necessarily unique).
    
* `operator` -
  (Optional)
  Matches a taint operator.
    
* `value` -
  (Optional)
  Matches a taint value.

<a name="nested_policy_content"></a>The `policy_content` block supports:

* `bundles` -
  (Optional)
  map of bundle name to BundleInstallSpec. The bundle name maps to the `bundleName` key in the `policycontroller.gke.io/constraintData` annotation on a constraint.

* `template_library`
  (Optional)
  Configures the installation of the Template Library. Structure is [documented below](#nested_template_library).

<a name="nested_bundles"></a>The `bundles` block supports:
    
* `bundle_name` -
  (Required)
  The name of the bundle.
    
* `exempted_namespaces` -
  (Optional)
  The set of namespaces to be exempted from the bundle.

<a name="nested_template_library"></a>The `template_library` block supports:

* `installation`
  (Optional)
  Configures the manner in which the template library is installed on the cluster. Must be one of `ALL`, `NOT_INSTALLED` or `INSTALLATION_UNSPECIFIED`. Defaults to `ALL`.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/locations/{{location}}/features/{{feature}}/membershipId/{{membership}}`

## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options: configuration options:

- `create` - Default is 20 minutes.
- `update` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import

FeatureMembership can be imported using any of these accepted formats:

* `projects/{{project}}/locations/{{location}}/features/{{feature}}/membershipId/{{membership}}`
* `{{project}}/{{location}}/{{feature}}/{{membership}}`
* `{{location}}/{{feature}}/{{membership}}`

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import FeatureMembership using one of the formats above. For example:

```tf
import {
  id = "projects/{{project}}/locations/{{location}}/features/{{feature}}/membershipId/{{membership}}"
  to = google_gke_hub_feature_membership.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), FeatureMembership can be imported using one of the formats above. For example:

```
$ terraform import google_gke_hub_feature_membership.default projects/{{project}}/locations/{{location}}/features/{{feature}}/membershipId/{{membership}}
$ terraform import google_gke_hub_feature_membership.default {{project}}/{{location}}/{{feature}}/{{membership}}
$ terraform import google_gke_hub_feature_membership.default {{location}}/{{feature}}/{{membership}}
```
