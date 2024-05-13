---
subcategory: "Cloud SQL"
page_title: "Google: google_sql_databases"
description: |-
  Get a list of databases in a Cloud SQL database instance.
---

# google_sql_databases

Use this data source to get information about a list of databases in a Cloud SQL instance.
## Example Usage


```hcl
data "google_sql_databases" "qa" {
  instance = google_sql_database_instance.main.name
}
```

## Argument Reference

The following arguments are supported:

* `instance` - (required) The name of the Cloud SQL database instance in which the database belongs.

* `project` - (optional) The ID of the project in which the instance belongs.

-> **Note** This datasource performs client-side sorting to provide consistent ordering of the databases.

## Attributes Reference
See [google_sql_database](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/sql_database) resource for details of all the available attributes.
