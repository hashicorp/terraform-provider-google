---
subcategory: "Cloud Bigtable"
description: |-
  A datasource to retrieve the IAM policy state for a Bigtable instance.
---


# `google_bigtable_instance_iam_policy`
Retrieves the current IAM policy data for a Bigtable instance.

## example

```hcl
data "google_bigtable_instance_iam_policy" "policy" {
  instance    = google_bigtable_instance.instance.name
}
```

## Argument Reference

The following arguments are supported:

* `instance` - (Required) The name or relative resource id of the instance to manage IAM policies for.


* `project` - (Optional) The project in which the instance belongs. If it
    is not provided, Terraform will use the provider default.

## Attributes Reference

The attributes are exported:

* `etag` - (Computed) The etag of the IAM policy.

* `policy_data` - (Computed) The policy data
