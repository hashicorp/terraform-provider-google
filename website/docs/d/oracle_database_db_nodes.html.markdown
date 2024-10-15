---
subcategory: "Oracle Database"
description: |-
  List all database nodes of a Cloud VmCluster.
---

# google_oracle_database_db_nodes

List all DbNodes of a Cloud VmCluster.

For more information see the
[API](https://cloud.google.com/oracle/database/docs/reference/rest/v1/projects.locations.cloudVmClusters.dbNodes).

## Example Usage

```hcl
data "google_oracle_database_db_nodes" "my_db_nodes"{
	location = "us-east4"
	cloud_vm_cluster = "vmcluster-id"
}
```

## Argument Reference

The following arguments are supported:

* `cloud_vm_cluster` - (Required) The ID of the VM Cluster.

* `location` - (Required) The location of the resource.

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes reference

The following attributes are exported:

* `db_nodes` - List of dbNodes. Structure is [documented below](#nested_dbnodes).

<a name="nested_dbnodes"></a> The `db_nodes` block supports:

* `name` - The name of the database node resource in the following format: projects/{project}/locations/{location}/cloudVmClusters/{cloudVmCluster}/dbNodes/{db_node}

* `properties` - Various properties of the database node. Structure is [documented below](#nested_properties).

<a name="nested_properties"></a> The `properties` block supports:

* `ocid`- OCID of database node.

* `ocpu_count` - OCPU count per database node.

* `memory_size_gb` - The allocated memory in GBs on the database node.

* `db_node_storage_size_gb` - The allocated local node storage in GBs on the database node.

* `db_server_ocid` - The OCID of the Database server associated with the database node.

* `hostname` - The host name for the database node.

* `state` - State of the database node.
<a name="nested_states"></a>Possible values for `state` are:<br>
`PROVISIONING` - Indicates that the resource is being provisioned.<br>
`AVAILABLE` - Indicates that the resource is available.<br>
`UPDATING` - Indicates that the resource is being updated.<br>
`STOPPING` - Indicates that the resource is being stopped.<br>
`STOPPED` - Indicates that the resource is stopped.<br>
`STARTING` - Indicates that the resource is being started.<br>
`TERMINATING` - Indicates that the resource is being terminated.<br>
`TERMINATED` - Indicates that the resource is terminated.<br>
`FAILED` - Indicates that the resource has failed.<br>

* `total_cpu_core_count` - The total number of CPU cores reserved on the database node.