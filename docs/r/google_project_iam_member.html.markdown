---
layout: "google"
page_title: "Google: google_project_iam_member"
sidebar_current: "docs-google-project-iam-member"
description: |-
 Allows management of a single member for a single binding on the IAM policy for a Google Cloud Platform project.
---

# google\_project\_iam\_member

Allows creation and management of a single member for a single binding within
the IAM policy for an existing Google Cloud Platform project.

~> **Note:** This resource _must not_ be used in conjunction with
   `google_project_iam_policy` or they will fight over what your policy
   should be. Similarly, roles controlled by `google_project_iam_binding`
   should not be assigned to using `google_project_iam_member`.

## Example Usage

```hcl
resource "google_project_iam_member" "project" {
  project = "your-project-id"
  role    = "roles/editor"
  member  = "user:jane@example.com"
}
```

## Argument Reference

The following arguments are supported:

* `member` - (Required) The user that the role should apply to.

* `role` - (Required) The role that should be applied.

* `project` - (Optional) The project ID. If not specified, uses the
    ID of the project configured with the provider.## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `etag` - (Computed) The etag of the project's IAM policy.
