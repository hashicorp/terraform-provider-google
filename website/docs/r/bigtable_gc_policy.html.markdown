---
subcategory: "Cloud Bigtable"
layout: "google"
page_title: "Google: google_bigtable_gc_policy"
sidebar_current: "docs-google-bigtable-gc-policy"
description: |-
  Creates a Google Cloud Bigtable GC Policy inside a family.
---

# google_bigtable_gc_policy

Creates a Google Cloud Bigtable GC Policy inside a family. For more information see
[the official documentation](https://cloud.google.com/bigtable/) and
[API](https://cloud.google.com/bigtable/docs/go/reference).


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
}

resource "google_bigtable_table" "table" {
  name          = "tf-table"
  instance_name = google_bigtable_instance.instance.name

  column_family {
    family = "name"
  }
}

resource "google_bigtable_gc_policy" "policy" {
  instance_name = google_bigtable_instance.instance.name
  table         = google_bigtable_table.table.name
  column_family = "name"

  max_age {
    days = 7
  }
}
```

Multiple conditions is also supported. `UNION` when any of its sub-policies apply (OR). `INTERSECTION` when all its sub-policies apply (AND)

```hcl
resource "google_bigtable_gc_policy" "policy" {
  instance_name = google_bigtable_instance.instance.name
  table         = google_bigtable_table.table.name
  column_family = "name"

  mode = "UNION"

  max_age {
    days = 7
  }

  max_version {
    number = 10
  }
}
```

## Argument Reference

The following arguments are supported:

* `table` - (Required) The name of the table.

* `instance_name` - (Required) The name of the Bigtable instance.

* `column_family` - (Required) The name of the column family.

* `project` - (Optional) The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

* `mode` - (Optional) If multiple policies are set, you should choose between `UNION` OR `INTERSECTION`.

* `max_age` - (Optional) GC policy that applies to all cells older than the given age.

* `max_version` - (Optional) GC policy that applies to all versions of a cell except for the most recent.

-----

`max_age` supports the following arguments:

* `days` - (Required) Number of days before applying GC policy.

-----

`max_version` supports the following arguments:

* `number` - (Required) Number of version before applying the GC policy.

## Attributes Reference

Only the arguments listed above are exposed as attributes.
