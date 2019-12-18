---
subcategory: "Container Registry"
layout: "google"
page_title: "Google: google_container_registry_repository"
sidebar_current: "docs-google-datasource-container-repo"
description: |-
  Get URLs for a given project's container registry repository.
---

# google\_container\_registry\_repository

This data source fetches the project name, and provides the appropriate URLs to use for container registry for this project.

The URLs are computed entirely offline - as long as the project exists, they will be valid, but this data source does not contact Google Container Registry (GCR) at any point.

## Example Usage

```hcl
data "google_container_registry_repository" "foo" {
}

output "gcr_location" {
  value = data.google_container_registry_repository.foo.repository_url
}
```

## Argument Reference
* `project`: (Optional) The project ID that this repository is attached to.  If not provided, provider project will be used instead.
* `region`: (Optional) The GCR region to use.  As of this writing, one of `asia`, `eu`, and `us`.  See [the documentation](https://cloud.google.com/container-registry/docs/pushing-and-pulling) for additional information.

## Attributes Reference
In addition to the arguments listed above, this data source exports:

* `repository_url`: The URL at which the repository can be accessed.
