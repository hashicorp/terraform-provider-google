---
subcategory: "Cloud SQL"
description: |-
  Creates a new SQL user in Google Cloud SQL.
---

# google\_sql\_user

Creates a new Google SQL User on a Google SQL User Instance. For more information, see the [official documentation](https://cloud.google.com/sql/), or the [JSON API](https://cloud.google.com/sql/docs/admin-api/v1beta4/users).

~> **Note:** All arguments including the username and password will be stored in the raw state as plain-text.
[Read more about sensitive data in state](https://www.terraform.io/language/state/sensitive-data). Passwords will not be retrieved when running
"terraform import".

## Example Usage

Example creating a SQL User.

```hcl
resource "random_id" "db_name_suffix" {
  byte_length = 4
}

resource "google_sql_database_instance" "main" {
  name             = "main-instance-${random_id.db_name_suffix.hex}"
  database_version = "MYSQL_5_7"

  settings {
    tier = "db-f1-micro"
  }
}

resource "google_sql_user" "users" {
  name     = "me"
  instance = google_sql_database_instance.main.name
  host     = "me.com"
  password = "changeme"
}
```

Example using [Cloud SQL IAM database authentication](https://cloud.google.com/sql/docs/mysql/authentication).

```hcl
resource "random_id" "db_name_suffix" {
  byte_length = 4
}

resource "google_sql_database_instance" "main" {
  name             = "main-instance-${random_id.db_name_suffix.hex}"
  database_version = "POSTGRES_15"

  settings {
    tier = "db-f1-micro"

    database_flags {
      name  = "cloudsql.iam_authentication"
      value = "on"
    }
  }
}

resource "google_sql_user" "iam_user" {
  name     = "me@example.com"
  instance = google_sql_database_instance.main.name
  type     = "CLOUD_IAM_USER"
}

resource "google_sql_user" "iam_service_account_user" {
  # Note: for Postgres only, GCP requires omitting the ".gserviceaccount.com" suffix
  # from the service account email due to length limits on database usernames.
  name     = trimsuffix(google_service_account.service_account.email, ".gserviceaccount.com")
  instance = google_sql_database_instance.main.name
  type     = "CLOUD_IAM_SERVICE_ACCOUNT"
}
```

## Argument Reference

The following arguments are supported:

* `instance` - (Required) The name of the Cloud SQL instance. Changing this
    forces a new resource to be created.

* `name` - (Required) The name of the user. Changing this forces a new resource
    to be created.

* `password` - (Optional) The password for the user. Can be updated. For Postgres
    instances this is a Required field, unless type is set to either CLOUD_IAM_USER
    or CLOUD_IAM_SERVICE_ACCOUNT. Don't set this field for CLOUD_IAM_USER
    and CLOUD_IAM_SERVICE_ACCOUNT user types for any Cloud SQL instance.

* `type` - (Optional) The user type. It determines the method to authenticate the
    user during login. The default is the database's built-in user type. Flags
    include "BUILT_IN", "CLOUD_IAM_USER", or "CLOUD_IAM_SERVICE_ACCOUNT".

* `deletion_policy` - (Optional) The deletion policy for the user.
    Setting `ABANDON` allows the resource to be abandoned rather than deleted. This is useful
    for Postgres, where users cannot be deleted from the API if they have been granted SQL roles.
    
    Possible values are: `ABANDON`.

- - -

* `host` - (Optional) The host the user can connect from. This is only supported
    for BUILT_IN users in MySQL instances. Don't set this field for PostgreSQL and SQL Server instances.
    Can be an IP address. Changing this forces a new resource to be created.

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

The optional `password_policy` block is only supported by Mysql. The `password_policy` block supports:

* `allowed_failed_attempts` - (Optional) Number of failed attempts allowed before the user get locked.

* `password_expiration_duration` - (Optional) Password expiration duration with one week grace period.

* `enable_failed_attempts_check` - (Optional) If true, the check that will lock user after too many failed login attempts will be enabled.

* `enable_password_verification` - (Optional) If true, the user must specify the current password before changing the password. This flag is supported only for MySQL.

The read only `password_policy.status` subblock supports:

* `locked` - (read only) If true, user does not have login privileges.

* `password_expiration_time` - (read only) Password expiration duration with one week grace period.

## Attributes Reference

Only the arguments listed above are exposed as attributes.

## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options: configuration options:

- `create` - Default is 10 minutes.
- `update` - Default is 10 minutes.
- `delete` - Default is 10 minutes.

## Import

SQL users for MySQL databases can be imported using the `project`, `instance`, `host` and `name`, e.g.

* `{{project_id}}/{{instance}}/{{host}}/{{name}}`

SQL users for PostgreSQL databases can be imported using the `project`, `instance` and `name`, e.g.

* `{{project_id}}/{{instance}}/{{name}}`

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import NAME_HERE using one of the formats above. For example:

```tf
# MySQL database
import {
  id = "{{project_id}}/{{instance}}/{{host}}/{{name}}"
  to = google_sql_user.default
}

# PostgreSQL database
import {
  id = "{{project_id}}/{{instance}}/{{name}}"
  to = google_sql_user.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), NAME_HERE can be imported using one of the formats above. For example:

```
# MySQL database
$ terraform import google_sql_user.default {{project_id}}/{{instance}}/{{host}}/{{name}}

# PostgreSQL database
$ terraform import google_sql_user.default {{project_id}}/{{instance}}/{{name}}
```