---
subcategory: "Oracle Database"
description: |-
  List all ExadataInfrastructures.
---

# google_oracle_database_cloud_exadata_infrastructures

List all ExadataInfrastructures.

For more information see the
[API](https://cloud.google.com/oracle/database/docs/reference/rest/v1/projects.locations.cloudExadataInfrastructures).

## Example Usage

```hcl
data "google_oracle_database_cloud_exadata_infrastructures" "my_exadatas"{
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

* `CloudExadataInfrastructures` - A list of ExadataInfrastructures.

See [google_oracle_database_cloud_exadata_infrastructure](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/oracle_database_cloud_exadata_infrastructure#argument-reference) resource for details of the available attributes.
