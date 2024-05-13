---
subcategory: "ContainerAzure"
description: |-
  Provides lists of available Kubernetes versions for the Container Azure resources.
---

# google_container_azure_versions

Provides access to available Kubernetes versions in a location for a given project.

## Example Usage

```hcl
data "google_container_azure_versions" "central1b" {
  location       = "us-west1"
  project        = "my-project"
}


output "first_available_version" {
  value = data.google_container_azure_versions.versions.valid_versions[0]
}
```

## Argument Reference

The following arguments are supported:

* `location` (Optional) - The location to list versions for.

* `project` (Optional) - ID of the project to list available cluster versions for. Should match the project the cluster will be deployed to.
  Defaults to the project that the provider is authenticated with.

## Attributes Reference

The following attributes are exported:

* `valid_versions` - A list of versions available for use with this project and location.
* `supported_regions` - A list of Azure regions that are available for use with this project and GCP location.
