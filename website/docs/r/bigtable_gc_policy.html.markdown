---
subcategory: "Cloud Bigtable"
page_title: "Google: google_bigtable_gc_policy"
description: |-
  Creates a Google Cloud Bigtable GC Policy inside a family.
---

# google_bigtable_gc_policy

Creates a Google Cloud Bigtable GC Policy inside a family. For more information see
[the official documentation](https://cloud.google.com/bigtable/) and
[API](https://cloud.google.com/bigtable/docs/go/reference).

-> **Warning**: We don't recommend having multiple GC policies for the same column
family as it may result in unexpected behavior.

-> **Note**: GC policies associated with a replicated table cannot be destroyed directly.
Destroying a GC policy is translated into never perform garbage collection, this is
considered relaxing from pure age-based or version-based GC policy, hence not allowed.
The workaround is unreplicating the instance first by updating the instance to have one
cluster.

## Example Usage

```hcl
resource "google_bigtable_instance" "instance" {
  name = "tf-instance"
  cluster {
    cluster_id   = "tf-instance-cluster"
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
  deletion_policy = "ABANDON"


  gc_rules = <<EOF
  {
    "rules": [
      {
        "max_age": "168h"
      }
    ]
  }
  EOF
}
```

Multiple conditions is also supported. `UNION` when any of its sub-policies apply (OR). `INTERSECTION` when all its sub-policies apply (AND)

```hcl
resource "google_bigtable_gc_policy" "policy" {
  instance_name = google_bigtable_instance.instance.name
  table         = google_bigtable_table.table.name
  column_family = "name"
  deletion_policy = "ABANDON"

  gc_rules = <<EOF
  {
    "mode": "union",
    "rules": [
      {
        "max_age": "168h"
      },
      {
        "max_version": 10
      }
    ]
  }
  EOF
}
```

An example of more complex GC policy:
```hcl
resource "google_bigtable_instance" "instance" {
  name = "instance_name"

  cluster {
    cluster_id = "cid"
    zone       = "us-central1-b"
  }

  instance_type = "DEVELOPMENT"
  deletion_protection = false
}

resource "google_bigtable_table" "table" {
  name          = "your-table"
  instance_name = google_bigtable_instance.instance.id

  column_family {
    family = "cf1"
  }
}

resource "google_bigtable_gc_policy" "policy" {
  instance_name = google_bigtable_instance.instance.id
  table         = google_bigtable_table.table.name
  column_family = "cf1"
  deletion_policy = "ABANDON"

  gc_rules = <<EOF
  {
    "mode": "union",
    "rules": [
      {
        "max_age": "10h"
      },
      {
        "mode": "intersection",
        "rules": [
          {
            "max_age": "2h"
          },
          {
            "max_version": 2
          }
        ]
      }
    ]
  }
  EOF
}
```
This is equivalent to running the following `cbt` command:
```
cbt setgcpolicy your-table cf1 "(maxage=2d and maxversions=2) or maxage=10h"
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

* `gc_rules` - (Optional) Serialized JSON object to represent a more complex GC policy. Conflicts with `mode`, `max_age` and `max_version`. Conflicts with `mode`, `max_age` and `max_version`.

* `deletion_policy` - (Optional) The deletion policy for the GC policy.
    Setting ABANDON allows the resource to be abandoned rather than deleted. This is useful for GC policy as it cannot be deleted in a replicated instance.

    Possible values are: `ABANDON`.

-----

`max_age` supports the following arguments:

* `days` - (Optional, Deprecated in favor of duration) Number of days before applying GC policy.

* `duration` - (Optional) Duration before applying GC policy (ex. "8h"). This is required when `days` isn't set

-----

`max_version` supports the following arguments:

* `number` - (Required) Number of version before applying the GC policy.

-----
`gc_rules` include 2 fields:
- `mode`: optional, either `intersection` or `union`.
- `rules`: an array of GC policy rule, can be specified as JSON object: `{"max_age": "16h"}` or `{"max_version": 2}`
- If `mode` is not specified, `rules` can only contains one GC policy. If `mode` is specified, `rules` must have at least 2 policies.

## Attributes Reference

Only the arguments listed above are exposed as attributes.

## Import

This resource does not support import.
