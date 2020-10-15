---
subcategory: "Cloud Spanner"
layout: "google"
page_title: "Google: google_compute_global_forwarding_rule"
sidebar_current: "docs-google-datasource-spanner-instance"
description: |-
  Get a spanner instance from Google Cloud
---

# google\_compute\_global_\forwarding\_rule

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
See [google_spanner_instance](https://www.terraform.io/docs/providers/google/r/spanner_instance.html) resource for details of all the available attributes.