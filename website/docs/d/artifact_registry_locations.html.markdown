---
subcategory: "Artifact Registry"
description: |-
  Get Artifact Registry locations available for a project.
---

# google_artifact_registry_locations

Get Artifact Registry locations available for a project. 

To get more information about Artifact Registry, see:

* [API documentation](https://cloud.google.com/artifact-registry/docs/reference/rest/v1/projects.locations/list)
* How-to Guides
    * [Official Documentation](https://cloud.google.com/artifact-registry/docs/overview)
    
## Example Usage

```hcl
data "google_artifact_registry_locations" "available" {
}
```

## Example Usage: Multi-regional Artifact Registry deployment

```hcl
data "google_artifact_registry_locations" "available" {
}

resource "google_artifact_registry_repository" "repo_one" {
  location = data.google_artifact_registry_locations.available.locations[0]
  repository_id = "repo-one"
  format        = "apt"
}

resource "google_artifact_registry_repository" "repo_two" {
  location = data.google_artifact_registry_locations.available.locations[1]
  repository_id = "repo-two"
  format        = "apt"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Optional) The project to list versions for. If it
    is not provided, the provider project is used.

## Attributes Reference

The following attributes are exported:

* `locations` - The list of Artifact Registry locations available for the given project.
