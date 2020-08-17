---
subcategory: "Cloud SQL"
layout: "google"
page_title: "Google: google_sql_user"
sidebar_current: "docs-google-sql-user"
description: |-
  Creates a new SQL user in Google Cloud SQL.
---

# google\_sql\_user

Creates a new Google SQL User on a Google SQL User Instance. For more information, see the [official documentation](https://cloud.google.com/sql/), or the [JSON API](https://cloud.google.com/sql/docs/admin-api/v1beta4/users).

~> **Note:** All arguments including the username and password will be stored in the raw state as plain-text.
[Read more about sensitive data in state](/docs/state/sensitive-data.html). Passwords will not be retrieved when running
"terraform import".

## Example Usage

Example creating a SQL User.

```hcl
resource "random_id" "db_name_suffix" {
  byte_length = 4
}

resource "google_sql_database_instance" "master" {
  name = "master-instance-${random_id.db_name_suffix.hex}"

  settings {
    tier = "db-f1-micro"
  }
}

resource "google_sql_user" "users" {
  name     = "me"
  instance = google_sql_database_instance.master.name
  host     = "me.com"
  password = "changeme"
}
```

## Argument Reference

The following arguments are supported:

* `instance` - (Required) The name of the Cloud SQL instance. Changing this
    forces a new resource to be created.

* `name` - (Required) The name of the user. Changing this forces a new resource
    to be created.

* `password` - (Optional) The password for the user. Can be updated. For Postgres
    instances this is a Required field.

- - -

* `host` - (Optional) The host the user can connect from. This is only supported
    for MySQL instances. Don't set this field for PostgreSQL instances.
    Can be an IP address. Changing this forces a new resource to be created.

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

Only the arguments listed above are exposed as attributes.

## Timeouts

This resource provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 10 minutes.
- `update` - Default is 10 minutes.
- `delete` - Default is 10 minutes.

## Import

SQL users for MySQL databases can be imported using the `project`, `instance`, `host` and `name`, e.g.

```
$ terraform import google_sql_user.users my-project/master-instance/my-domain.com/me
```

SQL users for PostgreSQL databases can be imported using the `project`, `instance` and `name`, e.g.

```
$ terraform import google_sql_user.users my-project/master-instance/me
```
