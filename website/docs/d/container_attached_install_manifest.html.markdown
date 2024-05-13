---
subcategory: "ContainerAttached"
description: |-
  Generates a YAML manifest for boot-strapping an Attached cluster registration.
---

# google_container_attached_install_manifest

Provides access to available platform versions in a location for a given project.

## Example Usage

```hcl
data "google_container_attached_install_manifest" "manifest" {
	location         = "us-west1"
	project          = "my-project"
	cluster_id       = "test-cluster-1"
	platform_version = "1.25.0-gke.1"
}


output "install_manifest" {
  value = data.google_container_attached_install_manifest.manifest
}
```

## Argument Reference

The following arguments are supported:

* `location` (Optional) - The location to list versions for.

* `project` (Optional) - ID of the project to list available platform versions for. Should match the project the cluster will be deployed to.
  Defaults to the project that the provider is authenticated with.

* `cluster_id` (Required) - The name that will be used when creating the attached cluster resource.

* `platform_version` (Required) - The platform version for the cluster. A list of valid values can be retrieved using the `google_container_attached_versions` data source.

## Attributes Reference

The following attributes are exported:

* `manifest` - A string with the YAML manifest that needs to be applied to the cluster.
