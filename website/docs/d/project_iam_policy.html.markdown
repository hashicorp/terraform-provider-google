---
subcategory: "Cloud Platform"
description: |-
  A datasource to retrieve the IAM policy state for a project.
---


# `google_project_iam_policy`
Retrieves the current IAM policy data for a project.

## example

```hcl
data "google_project_iam_policy" "policy" {
  project  = "myproject"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required) The project id of the target project. This is not
inferred from the provider.

## Attributes Reference

The attributes are exported:

* `etag` - (Computed) The etag of the IAM policy.

* `policy_data` - (Computed) The policy data
