---
subcategory: "Cloud Composer"
layout: "google"
page_title: "Google: google_composer_image_versions"
sidebar_current: "docs-google-datasource-composer-image-versions"
description: |-
  Provides available Cloud Composer versions.
---

# google\_composer\_image\_versions

Provides access to available Cloud Composer versions in a region for a given project.

## Example Usage

```hcl
data "google_composer_image_versions" "all" {
}

resource "google_composer_environment" "test" {
  name   = "test-env"
  region = "us-central1"
  config {
    software_config {
      image_version = data.google_composer_image_versions.all.image_versions[0].image_version_id
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Optional) The ID of the project to list versions in.
    If it is not provided, the provider project is used.

* `region` - (Optional) The location to list versions in.
    If it is not provider, the provider region is used.

## Attributes Reference

The following attributes are exported:

* `image_versions` - A list of composer image versions available in the given project and location. Each `image_version` contains:
  * `image_version_id` - The string identifier of the image version, in the form: "composer-x.y.z-airflow-a.b(.c)"
  * `supported_python_versions` - Supported python versions for this image version
