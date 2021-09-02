---
subcategory: "Cloud SQL"
layout: "google"
page_title: "Google: google_sql_database_instance"
sidebar_current: "docs-google-datasource-sql-database-instance"
description: |-
  Get a SQL database instance in Google Cloud SQL.
---

# google\_sql\_database\_instance

Use this data source to get information about a Cloud SQL instance.

## Example Usage


```hcl
data "google_sql_database_instance" "qa" {
  name = "test-sql-instance"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (required) The name of the instance.

* `project` - (optional) The ID of the project in which the resource belongs.

## Attributes Reference
See [google_sql_database_instance](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/sql_database_instance) resource for details of all the available attributes.
