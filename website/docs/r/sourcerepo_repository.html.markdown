---
layout: "google"
page_title: "Google: google_sourcerepo_repository"
sidebar_current: "docs-google-sourcerepo_repository"
description: |-
  Manages repositories within Google Cloud Source Repositories.
---

# google\_sourcerepo\_repository

For more information, see [the official
documentation](https://cloud.google.com/compute/docs/source-repositories) and
[API](https://cloud.google.com/source-repositories/docs/reference/rest/v1/projects.repos)

## Example Usage

This example is the common case of creating a repository within Google Cloud Source Repositories:

```hcl
resource "google_sourcerepo_repository" "frontend" {
  name = "frontend"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the repository that will be created.

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

The following attribute is exported:

* `size` - The size of the repository.
