---
subcategory: "Cloud Bigtable"
description: |-
  A datasource to retrieve the IAM policy state for a Bigtable Table.
---


# `google_bigtable_table_iam_policy`
Retrieves the current IAM policy data for a Bigtable Table.

## example

```hcl
data "google_bigtable_table_iam_policy" "policy" {
  instance    = google_bigtable_instance.instance.name
  table       = google_bigtable_table.table.name
}
```

## Argument Reference

The following arguments are supported:

* `instance` - (Required) The name or relative resource id of the instance that owns the table.

* `table` - (Required) The name or relative resource id of the table to manage IAM policies for.

* `project` - (Optional) The project in which the table belongs. If it
    is not provided, Terraform will use the provider default.

## Attributes Reference

The attributes are exported:

* `etag` - (Computed) The etag of the IAM policy.

* `policy_data` - (Computed) The policy data
