---
subcategory: "Cloud Spanner"
description: |-
  A datasource to retrieve the IAM policy state for a Spanner database.
---


# `google_spanner_database_iam_policy`
Retrieves the current IAM policy data for a Spanner database.

## example

```hcl
data "google_spanner_database_iam_policy" "foo" {
  project     = google_spanner_database.database.project
  database    = google_spanner_database.database.name
  instance    = google_spanner_database.database.instance
}
```

## Argument Reference

The following arguments are supported:

* `database` - (Required) The name of the Spanner database.

* `instance` - (Required) The name of the Spanner instance the database belongs to.

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

The attributes are exported:

* `etag` - (Computed) The etag of the IAM policy.

* `policy_data` - (Computed) The policy data
