---
layout: "google"
page_title: "Google: google_sourcerepos_repository"
sidebar_current: "docs-google-sourcerepos_repository"
description: |-
  Manages repositories within Google Cloud Source Repositories.
---

# google\_sourcerepos\_repository

Manages repositories within Google Cloud Source Repositories.

## Example Usage

This example is the common case of creating a repository within Google Cloud Source Repositores:

```hcl
resource "google_sourcerepos_repository" "frontend" {
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

Only the arguments listed above are exposed as attributes.
