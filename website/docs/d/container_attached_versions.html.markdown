---
subcategory: "ContainerAttached"
description: |-
  Provides lists of available platform versions for the Container Attached resources.
---

# google_container_attached_versions

Provides access to available platform versions in a location for a given project.

## Example Usage

```hcl
data "google_container_attached_versions" "uswest" {
  location       = "us-west1"
  project        = "my-project"
}


output "first_available_version" {
  value = data.google_container_attached_versions.versions.valid_versions[0]
}
```

## Argument Reference

The following arguments are supported:

* `location` (Optional) - The location to list versions for.

* `project` (Optional) - ID of the project to list available platform versions for. Should match the project the cluster will be deployed to.
  Defaults to the project that the provider is authenticated with.

## Attributes Reference

The following attributes are exported:

* `valid_versions` - A list of versions available for use with this project and location.
