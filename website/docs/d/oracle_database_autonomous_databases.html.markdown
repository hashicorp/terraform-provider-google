---
subcategory: "Oracle Database"
description: |-
  List all AutonomousDatabases.
---

# google_oracle_database_autonomous_databases

List all AutonomousDatabases.

For more information see the
[API](https://cloud.google.com/oracle/database/docs/reference/rest/v1/projects.locations.autonomousDatabases).

## Example Usage

```hcl
data "google_oracle_database_autonomous_databases" "my-adbs"{
  location = "us-east4"
}
```

## Argument Reference

The following arguments are supported:

* `location` - (Required) The location of the resource.

- - -
* `project` - (Optional) The project to which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

The following attributes are exported:

* `AutonomousDatabases` - A list of AutonomousDatabases.

See [google_oracle_database_autonomous_database](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/google_oracle_database_autonomous_database#argument-reference) resource for details of the available attributes.
