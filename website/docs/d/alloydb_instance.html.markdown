---
subcategory: "AlloyDB"
description: |-
  Fetches the details of available instance.
---

# google_alloydb_instance

Use this data source to get information about the available instance. For more details refer the [API docs](https://cloud.google.com/alloydb/docs/reference/rest/v1/projects.locations.clusters.instances).

## Example Usage


```hcl
data "google_alloydb_instance" "qa" {
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` -
  (Required)
  The ID of the alloydb cluster that the instance belongs to.
  'alloydb_cluster_id'

* `instance_id` -
  (Required)
  The ID of the alloydb instance.
  'alloydb_instance_id'

* `project` - 
  (optional) 
  The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

* `location` -
  (optional)
  The canonical id of the location.If it is not provided, the provider project is used. For example: us-east1.

## Attributes Reference

See [google_alloydb_instance](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/alloydb_instance) resource for details of all the available attributes.
