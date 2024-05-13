---
subcategory: "Artifact Registry"
description: |-
  Get information about a Google Artifact Registry Repository.
---

# google_artifact_registry_repository

Get information about a Google Artifact Registry Repository. For more information see
the [official documentation](https://cloud.google.com/artifact-registry/docs/)
and [API](https://cloud.google.com/artifact-registry/docs/apis).

## Example Usage

```hcl
data "google_artifact_registry_repository" "my-repo" {
  location      = "us-central1"
  repository_id = "my-repository"
}
```

## Argument Reference

The following arguments are supported:

* `repository_id` - (Required) The last part of the repository name.

* `location` - (Required) The location of the artifact registry repository. eg us-central1

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

See [google_artifact_registry_repository](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/artifact_registry_repository#argument-reference) resource for details of the available attributes.
