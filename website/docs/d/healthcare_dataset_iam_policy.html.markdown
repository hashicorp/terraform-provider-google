---
subcategory: "Cloud Healthcare"
description: |-
  A datasource to retrieve the IAM policy state for a Google Cloud Healthcare dataset.
---


# `google_healthcare_dataset_iam_policy`
Retrieves the current IAM policy data for a Google Cloud Healthcare dataset.

## example

```hcl
data "google_healthcare_dataset_iam_policy" "foo" {
  dataset_id  = google_healthcare_dataset.dataset.id
}
```

## Argument Reference

The following arguments are supported:

* `dataset_id` - (Required) The dataset ID, in the form
    `{project_id}/{location_name}/{dataset_name}` or
    `{location_name}/{dataset_name}`. In the second form, the provider's
    project setting will be used as a fallback.

## Attributes Reference

The attributes are exported:

* `etag` - (Computed) The etag of the IAM policy.

* `policy_data` - (Computed) The policy data
