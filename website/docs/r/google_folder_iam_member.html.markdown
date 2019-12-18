---
subcategory: "Cloud Platform"
layout: "google"
page_title: "Google: google_folder_iam_member"
sidebar_current: "docs-google-folder-iam-member"
description: |-
 Allows management of a single member for a single binding on the IAM policy for a Google Cloud Platform folder.
---

# google\_folder\_iam\_member

Allows creation and management of a single member for a single binding within
the IAM policy for an existing Google Cloud Platform folder.

~> **Note:** This resource _must not_ be used in conjunction with
   `google_folder_iam_policy` or they will fight over what your policy
   should be. Similarly, roles controlled by `google_folder_iam_binding`
   should not be assigned to using `google_folder_iam_member`.

## Example Usage

```hcl
resource "google_folder" "department1" {
  display_name = "Department 1"
  parent       = "organizations/1234567"
}

resource "google_folder_iam_member" "admin" {
  folder = google_folder.department1.name
  role   = "roles/editor"
  member = "user:alice@gmail.com"
}
```

## Argument Reference

The following arguments are supported:

* `folder` - (Required) The resource name of the folder the policy is attached to. Its format is folders/{folder_id}.

* `member` - (Required) The identity that will be granted the privilege in the `role`. For more details on format and restrictions see https://cloud.google.com/billing/reference/rest/v1/Policy#Binding
  This field can have one of the following values:
  * **user:{emailid}**: An email address that represents a specific Google account. For example, alice@gmail.com or joe@example.com.
  * **serviceAccount:{emailid}**: An email address that represents a service account. For example, my-other-app@appspot.gserviceaccount.com.
  * **group:{emailid}**: An email address that represents a Google group. For example, admins@example.com.
  * **domain:{domain}**: A G Suite domain (primary, instead of alias) name that represents all the users of that domain. For example, google.com or example.com.

* `role` - (Required) The role that should be applied. Note that custom roles must be of the format
    `[projects|organizations]/{parent-name}/roles/{role-name}`.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `etag` - (Computed) The etag of the folder's IAM policy.

## Import

IAM member imports use space-delimited identifiers; the resource in question, the role, and the account.  This member resource can be imported using the `folder`, role, and member identity e.g.

```
$ terraform import google_folder_iam_member.my_project "folder-name roles/viewer user:foo@example.com"
```
