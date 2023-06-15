---
subcategory: "Cloud Spanner"
description: |-
  A datasource to retrieve the IAM policy state for a Spanner instance.
---


# `google_spanner_instance_iam_policy`
Retrieves the current IAM policy data for a Spanner instance.

## example

```hcl
data "google_spanner_instance_iam_policy" "foo" {
  project     = google_spanner_instance.instance.project
  instance    = google_spanner_instance.instance.name
}
```

## Argument Reference

The following arguments are supported:

* `instance` - (Required) The name of the instance.

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

The attributes are exported:

* `etag` - (Computed) The etag of the IAM policy.

* `policy_data` - (Computed) The policy data
