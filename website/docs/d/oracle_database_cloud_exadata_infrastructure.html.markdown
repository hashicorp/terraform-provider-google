---
subcategory: "Oracle Database"
description: |-
  Get information about an ExadataInfrastructure.
---

# google_oracle_database_cloud_exadata_infrastructure

Get information about an ExadataInfrastructure.

## Example Usage

```hcl
data "google_oracle_database_cloud_exadata_infrastructure" "my-instance"{
  location = "us-east4"
  cloud_exadata_infrastructure_id = "exadata-id"
}
```

## Argument Reference

The following arguments are supported:

* `cloud_exadata_infrastructure_id` - (Required) The ID of the ExadataInfrastructure.

* `location` - (Required) The location of the resource.

- - -
* `project` - (Optional) The project to which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

See [google_oracle_database_cloud_exadata_infrastructure](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/google_oracle_database_cloud_exadata_infrastructure#argument-reference) resource for details of the available attributes.