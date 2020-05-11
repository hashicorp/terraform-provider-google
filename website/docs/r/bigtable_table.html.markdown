---
subcategory: "Cloud Bigtable"
layout: "google"
page_title: "Google: google_bigtable_table"
sidebar_current: "docs-google-bigtable-table"
description: |-
  Creates a Google Cloud Bigtable table inside an instance.
---

# google_bigtable_table

Creates a Google Cloud Bigtable table inside an instance. For more information see
[the official documentation](https://cloud.google.com/bigtable/) and
[API](https://cloud.google.com/bigtable/docs/go/reference).

-> **Note**: It is strongly recommended to set `lifecycle { prevent_destroy = true }`
on tables in order to prevent accidental data loss. See
[Terraform docs](https://www.terraform.io/docs/configuration/resources.html#prevent_destroy)
for more information on lifecycle parameters.


## Example Usage

```hcl
resource "google_bigtable_instance" "instance" {
  name = "tf-instance"

  cluster {
    cluster_id   = "tf-instance-cluster"
    zone         = "us-central1-b"
    num_nodes    = 3
    storage_type = "HDD"
  }

  lifecycle {
    prevent_destroy = true
  }
}

resource "google_bigtable_table" "table" {
  name          = "tf-table"
  instance_name = google_bigtable_instance.instance.name
  split_keys    = ["a", "b", "c"]

  lifecycle {
    prevent_destroy = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the table.

* `instance_name` - (Required) The name of the Bigtable instance.

* `split_keys` - (Optional) A list of predefined keys to split the table on.
!> **Warning:** Modifying the `split_keys` of an existing table will cause Terraform
to delete/recreate the entire `google_bigtable_table` resource.

* `column_family` - (Optional) A group of columns within a table which share a common configuration. This can be specified multiple times. Structure is documented below.

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

-----

`column_family` supports the following arguments:

* `family` - (Optional) The name of the column family.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `id` - an identifier for the resource with format `projects/{{project}}/instances/{{instance_name}}/tables/{{name}}`


## Import

Bigtable Tables can be imported using any of these accepted formats:

```
$ terraform import google_bigtable_table.default projects/{{project}}/instances/{{instance_name}}/tables/{{name}}
$ terraform import google_bigtable_table.default {{project}}/{{instance_name}}/{{name}}
$ terraform import google_bigtable_table.default {{instance_name}}/{{name}}
```

The following fields can't be read and will show diffs if set in config when imported:

- `split_keys`
