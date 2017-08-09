---
layout: "google"
page_title: "Google: google_sql_database"
sidebar_current: "docs-google-sql-database-x"
description: |-
  Creates a new SQL database in Google Cloud SQL.
---

# google\_sql\_database

Creates a new Google SQL Database on a Google SQL Database Instance. For more information, see the [official documentation](https://cloud.google.com/sql/), or the [JSON API](https://cloud.google.com/sql/docs/admin-api/v1beta4/databases).

## Example Usage

Example creating a SQL Database.

```hcl
resource "google_sql_database_instance" "master" {
  name = "master-instance"

  settings {
    tier = "D0"
  }
}

resource "google_sql_database" "users" {
  name      = "users-db"
  instance  = "${google_sql_database_instance.master.name}"
  charset   = "latin1"
  collation = "latin1_swedish_ci"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the database.

* `instance` - (Required) The name of containing instance.

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

* `charset` - (Optional) The charset value. See MySQL's [Supported Character
    Sets and
    Collations](https://dev.mysql.com/doc/refman/5.7/en/charset-charsets.html)
    and PostgreSQL's [Character Set
    Support](https://www.postgresql.org/docs/9.6/static/multibyte.html)
    for more details and supported values. Note that Cloud SQL's beta
    offering for PostgreSQL databases currently only supports the charset value
    `UTF8`.

* `collation` - (Optional) The collation value. See MySQL's [Supported Character
    Sets and
    Collations](https://dev.mysql.com/doc/refman/5.7/en/charset-charsets.html)
    and PostgreSQL's [Collation
    Support](https://www.postgresql.org/docs/9.6/static/collation.html) for
    more details and supported values. Note that Cloud SQL's beta
    offering for PostgreSQL databases currently only supports the collation
    value `en_US.UTF8`.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `self_link` - The URI of the created resource.

## Import

SQL databases can be imported using the `instance` and `name`, e.g.

```
$ terraform import google_sql_database.database master-instance:users-db
```
