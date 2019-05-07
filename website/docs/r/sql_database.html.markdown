---
layout: "google"
page_title: "Google: google_sql_database"
sidebar_current: "docs-google-sql-database-x"
description: |-
  Creates a new SQL database in Google Cloud SQL.
---

# google\_sql\_database

Creates a new Google SQL Database on a Google SQL Database Instance. For more information, see
the [official documentation](https://cloud.google.com/sql/),
or the [JSON API](https://cloud.google.com/sql/docs/admin-api/v1beta4/databases).

## Example Usage

Example creating a SQL Database.

```hcl
resource "random_id" "db_name_suffix" {
  byte_length = 4
}

resource "google_sql_database_instance" "master" {
  name = "master-instance-${random_id.db_name_suffix.hex}"

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

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

* `charset` - (Optional) The charset value. See MySQL's
    [Supported Character Sets and Collations](https://dev.mysql.com/doc/refman/5.7/en/charset-charsets.html)
    and Postgres' [Character Set Support](https://www.postgresql.org/docs/9.6/static/multibyte.html)
    for more details and supported values. Postgres databases are in beta
    and have limited `charset` support; they only support a value of `UTF8` at creation time.

* `collation` - (Optional) The collation value. See MySQL's
    [Supported Character Sets and Collations](https://dev.mysql.com/doc/refman/5.7/en/charset-charsets.html)
    and Postgres' [Collation Support](https://www.postgresql.org/docs/9.6/static/collation.html)
    for more details and supported values. Postgres databases are in beta
    and have limited `collation` support; they only support a value of `en_US.UTF8` at creation time.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `self_link` - The URI of the created resource.

## Import

SQL databases can be imported using one of any of these accepted formats:

```
$ terraform import google_sql_database.database projects/{{project}}/instances/{{instance}}/databases/{{name}}
$ terraform import google_sql_database.database {{project}}/{{instance}}/{{name}}
$ terraform import google_sql_database.database instances/{{name}}/databases/{{name}}
$ terraform import google_sql_database.database {{instance}}/{{name}}
$ terraform import google_sql_database.database {{name}}

```
