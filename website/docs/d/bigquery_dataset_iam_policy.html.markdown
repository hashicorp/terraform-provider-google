---
subcategory: "BigQuery"
description: |-
  A datasource to retrieve the IAM policy state for a BigQuery dataset.
---


# `google_bigquery_dataset_iam_policy`
Retrieves the current IAM policy data for a BigQuery dataset.

## example

```hcl
data "google_bigquery_dataset_iam_policy" "policy" {
  dataset_id  = google_bigquery_dataset.dataset.dataset_id
}
```

## Argument Reference

The following arguments are supported:

* `dataset_id` - (Required) The dataset ID.

* `project` - (Optional) The ID of the project in which the resource belongs.
    If it is not provided, the provider project is used.

## Attributes Reference

The attributes are exported:

* `etag` - (Computed) The etag of the IAM policy.

* `policy_data` - (Computed) The policy data
