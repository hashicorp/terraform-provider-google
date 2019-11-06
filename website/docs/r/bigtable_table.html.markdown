---
subcategory: "Bigtable"
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


## Example Usage

```hcl
resource "google_bigtable_instance" "instance" {
  name         = "tf-instance"
  cluster_id   = "tf-instance-cluster"
  zone         = "us-central1-b"
  num_nodes    = 3
  storage_type = "HDD"
}

resource "google_bigtable_table" "table" {
  name          = "tf-table"
  instance_name = "${google_bigtable_instance.instance.name}"
  split_keys    = ["a", "b", "c"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the table.

* `instance_name` - (Required) The name of the Bigtable instance.

* `split_keys` - (Optional) A list of predefined keys to split the table on.

* `column_family` - (Optional) A group of columns within a table which share a common configuration. This can be specified multiple times. Structure is documented below.

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

-----

`column_family` supports the following arguments:

* `family` - (Optional) The name of the column family.

## Attributes Reference

Only the arguments listed above are exposed as attributes.
