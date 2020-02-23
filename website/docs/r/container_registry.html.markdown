---
subcategory: "Container Registry"
layout: "google"
page_title: "Google: google_container_registry"
sidebar_current: "docs-google-container-registry"
description: |-
  Ensures the GCS bucket backing Google Container Registry exists.
---

# google_container_registry

Ensures that the Google Cloud Storage bucket that backs Google Container Registry exists. Creating this resource will create the backing bucket if it does not exist, or do nothing if the bucket already exists. Destroying this resource does *NOT* destroy the backing bucket. For more information see [the official documentation](https://cloud.google.com/container-registry/docs/overview)

This resource can be used to ensure that the GCS bucket exists prior to assigning permissions. For more information see the [access control page](https://cloud.google.com/container-registry/docs/access-control) for GCR.


## Example Usage

```hcl
resource "google_container_registry" "registry" {
  project  = "my-project"
  location = "EU"
}
```

The `id` field of the `google_container_registry` is the identifier of the storage bucket that backs GCR and can be used to assign permissions to the bucket.

```hcl
resource "google_container_registry" "registry" {
  project  = "my-project"
  location = "EU"
}

resource "google_storage_bucket_iam_member" "viewer" {
  bucket = google_container_registry.registry.id
  role = "roles/storage.objectViewer"
  member = "user:jane@example.com"
}
```

## Argument Reference

The following arguments are supported:

* `location` - (Optional) The location of the registry. One of `ASIA`, `EU`, `US` or not specified. See [the official documentation](https://cloud.google.com/container-registry/docs/pushing-and-pulling#pushing_an_image_to_a_registry) for more information on registry locations.

* `project` - (Optional) The ID of the project in which the resource belongs. If it is not provided, the provider project is used.


## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `bucket_self_link` - The URI of the created resource.

* `id` - The name of the bucket that supports the Container Registry. In the form of `artifacts.{project}.appspot.com` or `{location}.artifacts.{project}.appspot.com` if location is specified.
