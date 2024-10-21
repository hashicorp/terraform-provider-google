---
subcategory: "Oracle Database"
description: |-
  Get information about an AutonomousDatabase.
---

# google_oracle_database_autonomous_database

Get information about an AutonomousDatabase.

For more information see the
[API](https://cloud.google.com/oracle/database/docs/reference/rest/v1/projects.locations.autonomousDatabases).

## Example Usage

```hcl
data "google_oracle_database_autonomous_database" "my-instance"{
  location = "us-east4"
  autonomous_database_id = "autonomous_database_id"
}
```

## Argument Reference

The following arguments are supported:

* `autonomous_database_id` - (Required) The ID of the AutonomousDatabase.

* `location` - (Required) The location of the resource.

- - -
* `project` - (Optional) The project to which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

See [google_oracle_database_autonomous_database](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/oracle_database_autonomous_database#argument-reference) resource for details of the available attributes.