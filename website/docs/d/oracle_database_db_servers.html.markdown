---
subcategory: "Oracle Database"
description: |-
  List all DbServers of a Cloud ExdataInfrastructure.
---

# google_oracle_database_db_servers

List all DbServers of a Cloud Exdata Infrastructure.

## Example Usage

```hcl
data "google_oracle_database_db_servers" "my_db_servers"{
	location = "us-east4"
	cloud_exadata_infrastructure = "exadata-id"
}
```

## Argument Reference

The following arguments are supported:

* `cloud_exadata_infrastructure` - (Required) The Exadata Infrastructure id.

* `location` - (Required) The location of resource.

- - -
* `project` - (Optional) The project to which the resource belongs. If it
    is not provided, the provider project is used.

* `db_servers` - (Output) List of dbServers

<a name="nested_properties"></a> The `db_servers` block supports:

* `display_name` - User friendly name for the resource.

* `properties` - Various properties of the databse server.

<a name="nested_properties"></a> The `properties` block supports:

* `ocid` - The OCID of database server.

* `ocpu_count` - The OCPU count per database.

* `max_ocpu_count` - The total number of CPU cores available.

* `memory_size_gb` - The allocated memory in gigabytes on the database server.

* `max_memory_size_gb` - The total memory available in gigabytes.

* `db_node_storage_size_gb` - The local storage per VM.

* `max_db_node_storage_size_gb` - The total local node storage available in GBs.

* `vm_count` - The VM count per database.

* `state` - The current state of the database server.
<a name="nested_states"></a>Allowed values for `state` are:<br>
`STATE_UNSPECIFIED` - Default unspecified value.<br>
`CREATING` - Indicates that the resource is being created.<br>
`AVAILABLE` - Indicates that the resource is available.<br>
`UNAVAILABLE` - Indicates that the resource is unavailable.<br>
`DELETING` - Indicates that the resource is being deleted.<br>
`DELETED` - Indicates that the resource has been deleted.<br>

* `db_node_ids` - The OCID of database nodes associated with the database server.
