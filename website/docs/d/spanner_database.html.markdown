---
subcategory: "Cloud Spanner"
description: |-
  Get a spanner database from Google Cloud
---

# google_spanner_database

Get a spanner database from Google Cloud by its name and instance name.

## Example Usage

```tf
data "google_spanner_database" "foo" {
  name     = "foo"
  instance = google_spanner_instance.instance.name
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the spanner database.

* `instance` - (Required) The name of the database's spanner instance.

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference
See [google_spanner_database](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/spanner_database) resource for details of all the available attributes.

**Note** `ddl` is a field where reads are ignored, and thus will show up with a null value.
