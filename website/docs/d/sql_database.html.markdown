---
subcategory: "Cloud SQL"
description: |-
  Get a database in a Cloud SQL database instance.
---

# google_sql_database

Use this data source to get information about a database in a Cloud SQL instance.

## Example Usage


```hcl
data "google_sql_database" "qa" {
  name = "test-sql-database"
  instance = google_sql_database_instance.main.name
}
```

## Argument Reference

The following arguments are supported:

* `name` - (required) The name of the database.

* `instance` - (required) The name of the Cloud SQL database instance in which the database belongs.

* `project` - (optional) The ID of the project in which the instance belongs.

## Attributes Reference
See [google_sql_database](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/sql_database) resource for details of all the available attributes.
