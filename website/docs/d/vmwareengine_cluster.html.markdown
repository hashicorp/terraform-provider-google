---
subcategory: "Cloud VMware Engine"
description: |-
  Get info about a private cloud cluster.
---

# google\_vmwareengine\_cluster

Use this data source to get details about a cluster resource.

~> **Warning:** This data source is in beta, and should be used with the terraform-provider-google-beta provider.
See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html) for more details on beta resources.

To get more information about private cloud cluster, see:
* [API documentation](https://cloud.google.com/vmware-engine/docs/reference/rest/v1/projects.locations.privateClouds.clusters)

## Example Usage

```hcl
data "google_vmwareengine_cluster" "my_cluster" {
  provider = google-beta
  name     = "my-cluster"
  parent   = "project/locations/us-west1-a/privateClouds/my-cloud"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the resource.
* `parent` - (Required) The resource name of the private cloud that this cluster belongs.

## Attributes Reference

See [google_vmwareengine_cluster](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/vmwareengine_cluster#attributes-reference) resource for details of all the available attributes.