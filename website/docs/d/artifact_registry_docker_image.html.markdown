---
subcategory: "Artifact Registry"
description: |-
  Get information about a Docker Image within a Google Artifact Registry Repository.
---

# google\_artifact\_registry\_docker\_image

This data source fetches information from a provided Artifact Registry repository, including the fully qualified name and URI for an image, based on a the latest version of image name and optional digest or tag.

~> **Note**
Requires one of the following OAuth scopes: `https://www.googleapis.com/auth/cloud-platform` or `https://www.googleapis.com/auth/cloud-platform.read-only`.

## Example Usage

```hcl
resource "google_artifact_registry_repository" "my_repo" {
  location      = "us-west1"
  repository_id = "my-repository"
  format        = "DOCKER"
}

data "google_artifact_registry_docker_image" "my_image" {
  repository = google_artifact_registry_repository.my_repo.id
  image      = "my-image"
  tag        = "my-tag"
}

resource "google_cloud_run_v2_service" "default" {
 # ...
 
  template {
    containers {
      image = data.google_artifact_registry_docker_image.my_image.self_link
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `location` - (Required) The location of the artifact registry.

* `repository_id` - (Required) The last part of the repository name. to fetch from.

* `image_name` - (Required) The image name to fetch. If no digest or tag is provided, then the latest modified image will be used.

* `project` - (Optional) The project ID in which the resource belongs. If it is not provided, the provider project is used.

## Attributes Reference

The following computed attributes are exported:

* `name` - The fully qualified name of the fetched image.  This name has the form: `projects/{{project}}/locations/{{location}}/repository/{{repository_id}}/dockerImages/{{docker_image}}`. For example, 
```
projects/test-project/locations/us-west4/repositories/test-repo/dockerImages/nginx@sha256:e9954c1fc875017be1c3e36eca16be2d9e9bccc4bf072163515467d6a823c7cf
```

* `self_link` - The URI to access the image.  For example, 
```
us-west4-docker.pkg.dev/test-project/test-repo/nginx@sha256:e9954c1fc875017be1c3e36eca16be2d9e9bccc4bf072163515467d6a823c7cf
```

* `tags` - A list of all tags associated with the image.

* `image_size_bytes` - Calculated size of the image in bytes.

* `media_type` - Media type of this image, e.g. `application/vnd.docker.distribution.manifest.v2+json`. 

* `upload_time` - The time, as a RFC 3339 string, the image was uploaded. For example, `2014-10-02T15:01:23.045123456Z`.

* `build_time` - The time, as a RFC 3339 string, this image was built. 

* `update_time` - The time, as a RFC 3339 string, this image was updated.
