---
subcategory: "Oracle Database"
description: |-
  Get information about a CloudVmCluster.
---

# google_oracle_database_cloud_vm_cluster

Get information about a CloudVmCluster.

For more information see the
[API](https://cloud.google.com/oracle/database/docs/reference/rest/v1/projects.locations.cloudVmClusters).

## Example Usage

```hcl
data "google_oracle_database_cloud_vm_cluster" "my-vmcluster"{
  location = "us-east4"
  cloud_vm_cluster_id = "vmcluster-id"
}
```

## Argument Reference

The following arguments are supported:

* `cloud_vm_cluster_id` - (Required) The ID of the VM Cluster.

* `location` - (Required) The location of the resource.

- - -
* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

See [google_oracle_database_cloud_vm_cluster](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/google_oracle_database_cloud_vm_cluster#argument-reference) resource for details of the available attributes.