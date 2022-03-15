---
subcategory: "GKEHub"
layout: "google"
page_title: "Google: google_gke_hub_feature"
sidebar_current: "docs-google-gkehub-feature"
description: |-
  Contains information about a GKEHub Feature.
---

# google\_gkehub\_feature

Contains information about a GKEHub Feature. The google_gke_hub is the Fleet API.

* [API documentation](https://cloud.google.com/anthos/multicluster-management/reference/rest/v1beta/projects.locations.features)

~> **Warning:** This resource is in beta, and should be used with the terraform-provider-google-beta provider.
See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html) for more details on beta resources.


## Example Usage - Multi Cluster Ingress

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
  description = "Membership"
  provider = google-beta
}

resource "google_gke_hub_feature" "feature" {
  name = "multiclusteringress"
  location = "global"
  spec {
    multiclusteringress {
      config_membership = google_gke_hub_membership.membership.id
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

## Argument Reference

The following arguments are supported:

* `location` -
  (Required)
  The location for the resource

- - -

* `labels` -
  (Optional)
  GCP labels for this Feature.
  
* `name` -
  (Optional)
  The full, unique name of this Feature resource
  
* `project` -
  (Optional)
  The project for the resource
  
* `spec` -
  (Optional)
  Optional. Hub-wide Feature configuration. If this Feature does not support any Hub-wide configuration, this field may be unused.
  

The `spec` block supports:
    
* `multiclusteringress` -
  (Optional)
  Multicluster Ingress-specific spec.
    The `multiclusteringress` block supports:
    
* `config_membership` -
  (Optional)
  Fully-qualified Membership name which hosts the MultiClusterIngress CRD. Example: `projects/foo-proj/locations/global/memberships/bar`
    
## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/locations/{{location}}/features/{{name}}`

* `create_time` -
  Output only. When the Feature resource was created.
  
* `delete_time` -
  Output only. When the Feature resource was deleted.
  
* `update_time` -
  Output only. When the Feature resource was last updated.
  
## Timeouts

This resource provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 20 minutes.
- `update` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import

Feature can be imported using any of these accepted formats:

```
$ terraform import google_gke_hub_feature.default projects/{{project}}/locations/{{location}}/features/{{name}}
$ terraform import google_gke_hub_feature.default {{project}}/{{location}}/{{name}}
$ terraform import google_gke_hub_feature.default {{location}}/{{name}}
```



