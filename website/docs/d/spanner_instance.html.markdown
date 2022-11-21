---
subcategory: "Cloud Spanner"
page_title: "Google: google_spanner_instance"
description: |-
  Get a spanner instance from Google Cloud
---

# google\_spanner\_instance

Get a spanner instance from Google Cloud by its name.

## Example Usage

```tf
data "google_spanner_instance" "foo" {
  name = "bar"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the spanner instance.

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference
See [google_spanner_instance](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/spanner_instance) resource for details of all the available attributes.
