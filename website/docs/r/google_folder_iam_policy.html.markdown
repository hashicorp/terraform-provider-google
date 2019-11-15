---
subcategory: "Cloud Platform"
layout: "google"
page_title: "Google: google_folder_iam_policy"
sidebar_current: "docs-google-folders-iam-policy"
description: |-
 Allows management of the IAM policy for a Google Cloud Platform folders.
---

# google\_folder\_iam\_policy

Allows creation and management of the IAM policy for an existing Google Cloud
Platform folder.

## Example Usage

```hcl
resource "google_folder_iam_policy" "folder_admin_policy" {
  folder      = google_folder.department1.name
  policy_data = data.google_iam_policy.admin.policy_data
}

resource "google_folder" "department1" {
  display_name = "Department 1"
  parent       = "organizations/1234567"
}

data "google_iam_policy" "admin" {
  binding {
    role = "roles/editor"

    members = [
      "user:jane@example.com",
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `folder` - (Required) The resource name of the folder the policy is attached to. Its format is folders/{folder_id}.

* `policy_data` - (Required) The `google_iam_policy` data source that represents
    the IAM policy that will be applied to the folder. This policy overrides any existing
    policy applied to the folder.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `etag` - (Computed) The etag of the folder's IAM policy. `etag` is used for optimistic concurrency control as a way to help prevent simultaneous updates of a policy from overwriting each other.

## Import

A policy can be imported using the `folder`, e.g.

```
$ terraform import google_folder_iam_policy.my-folder-policy {{folder_id}}
```

