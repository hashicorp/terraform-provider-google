---
subcategory: "Cloud SQL"
description: |-
  Get a list of SQL database instances in a project in Google Cloud SQL.
---

# google_sql_database_instances

Use this data source to get information about a list of Cloud SQL instances in a project. You can also apply some filters over this list to get a more filtered list of Cloud SQL instances.

## Example Usage


```hcl
data "google_sql_database_instances" "qa" {
  project = "test-project"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (optional) The ID of the project in which the resources belong. If it is not provided, the provider project is used.

* `database_version` - (optional) To filter out the Cloud SQL instances which are of the specified database version.

* `region` - (optional) To filter out the Cloud SQL instances which are located in the specified region.

* `zone` - (optional) To filter out the Cloud SQL instances which are located in the specified zone. This zone refers to the Compute Engine zone that the instance is currently serving from.

* `tier` - (optional) To filter out the Cloud SQL instances based on the tier(or machine type) of the database instances.

* `state` - (optional) To filter out the Cloud SQL instances based on the current serving state of the database instance. Supported values include `SQL_INSTANCE_STATE_UNSPECIFIED`, `RUNNABLE`, `SUSPENDED`, `PENDING_DELETE`, `PENDING_CREATE`, `MAINTENANCE`, `FAILED`.

## Attributes Reference
See [google_sql_database_instance](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/sql_database_instance) resource for details of all the available attributes.
