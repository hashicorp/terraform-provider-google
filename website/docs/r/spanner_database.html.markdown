---
layout: "google"
page_title: "Google: google_spanner_database"
sidebar_current: "docs-google-spanner-database"
description: |-
  Creates a Google Spanner Database within a Spanner Instance.
---

# google\_spanner\_instance

Creates a Google Spanner Database within a Spanner Instance. For more information, see the [official documentation](https://cloud.google.com/spanner/), or the [JSON API](https://cloud.google.com/spanner/docs/reference/rest/v1/projects.instances.databases).

## Example Usage

Example creating a Spanner database.

```hcl
resource "google_spanner_instance" "main" {
  config       = "regional-europe-west1"
  display_name = "main-instance"
}

resource "google_spanner_database" "db" {
  instance  = "${google_spanner_instance.main.name}"
  name      = "main-instance"
  ddl       =  [
    "CREATE TABLE t1 (t1 INT64 NOT NULL,) PRIMARY KEY(t1)",
    "CREATE TABLE t2 (t2 INT64 NOT NULL,) PRIMARY KEY(t2)"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `instance` - (Required) The name of the instance that will serve the new database.

* `name` - (Required) The name of the database.

- - -

* `project` - (Optional) The project in which to look for the `instance` specified. If it
    is not provided, the provider project is used.

* `ddl` - (Optional) An optional list of DDL statements to run inside the newly created
   database. Statements can create tables, indexes, etc. These statements execute atomically
   with the creation of the database: if there is an error in any statement, the database
   is not created.


## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `state` - The current state of the database.

## Import

Databases can be imported via their `instance` and `name` values, and optionally
the `project` in which the instance is defined (Often used when the project is different
to that defined in the provider). The format is thus either `{instanceName}/{dbName}` or
`{projectId}/{instanceName}/{dbName}`. e.g.

```
$ terraform import google_spanner_database.db1 instance456/db789

$ terraform import google_spanner_database.db1 project123/instance456/db789

```
