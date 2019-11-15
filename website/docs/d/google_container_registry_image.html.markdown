---
subcategory: "Container Registry"
layout: "google"
page_title: "Google: google_container_registry_image"
sidebar_current: "docs-google-datasource-container-image"
description: |-
  Get URLs for a given project's container registry image.
---

# google\_container\_registry\_image

This data source fetches the project name, and provides the appropriate URLs to use for container registry for this project.

The URLs are computed entirely offline - as long as the project exists, they will be valid, but this data source does not contact Google Container Registry (GCR) at any point.

## Example Usage

```hcl
data "google_container_registry_image" "debian" {
  name = "debian"
}

output "gcr_location" {
  value = data.google_container_registry_image.debian.image_url
}
```

## Argument Reference
* `name`: (Required) The image name.
* `project`: (Optional) The project ID that this image is attached to.  If not provider, provider project will be used instead.
* `region`: (Optional) The GCR region to use.  As of this writing, one of `asia`, `eu`, and `us`.  See [the documentation](https://cloud.google.com/container-registry/docs/pushing-and-pulling) for additional information.
* `tag`: (Optional) The tag to fetch, if any.
* `digest`: (Optional) The image digest to fetch, if any.

## Attributes Reference
In addition to the arguments listed above, this data source exports:
* `image_url`: The URL at which the image can be accessed.
