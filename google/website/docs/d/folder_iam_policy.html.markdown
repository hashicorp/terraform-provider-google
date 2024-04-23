---
subcategory: "Cloud Platform"
description: |-
  A datasource to retrieve the IAM policy state for a folder.
---


# `google_folder_iam_policy`
Retrieves the current IAM policy data for a folder.

## example

```hcl
data "google_folder_iam_policy" "test" {
  folder      = google_folder.permissiontest.name
}
```

## Argument Reference

The following arguments are supported:

* `folder` - (Required) The resource name of the folder the policy is attached to. Its format is folders/{folder_id}.

## Attributes Reference

The attributes are exported:

* `etag` - (Computed) The etag of the IAM policy.

* `policy_data` - (Computed) The policy data
