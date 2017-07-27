---
layout: "google"
page_title: "Google: google_project_iam_binding"
sidebar_current: "docs-google-project-iam-binding"
description: |-
 Allows management of a single binding with an IAM policy for a Google Cloud Platform project.
---

# google\_project\_iam\_binding

Allows creation and management of a single binding within IAM policy for
an existing Google Cloud Platform project.

~> **Note:** This resource _must not_ be used in conjunction with
   `google_project_iam_policy` or they will fight over what your policy
   should be.

## Example Usage

```hcl
resource "google_project_iam_binding" "project" {
  project = "your-project-id"
  role    = "roles/editor"

  members = [
    "user:jane@example.com",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `members` - (Required) A list of users that the role should apply to.

* `role` - (Required) The role that should be applied. Only one
    `google_project_iam_binding` can be used per role.

* `project` - (Optional) The project ID. If not specified, uses the
    ID of the project configured with the provider.## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `etag` - (Computed) The etag of the project's IAM policy.

