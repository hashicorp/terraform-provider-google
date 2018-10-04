---
layout: "google"
page_title: "Google: google_iam_policy"
sidebar_current: "docs-google-datasource-iam-policy"
description: |-
  Generates an IAM policy that can be referenced by other resources, applying
  the policy to them.
---

# google\_iam\_policy

Generates an IAM policy document that may be referenced by and applied to
other Google Cloud Platform resources, such as the `google_project` resource.

```
data "google_iam_policy" "admin" {
  binding {
    role = "roles/compute.instanceAdmin"

    members = [
      "serviceAccount:your-custom-sa@your-project.iam.gserviceaccount.com",
    ]
  }

  binding {
    role = "roles/storage.objectViewer"

    members = [
      "user:jane@example.com",
    ]
  }
}
```

This data source is used to define IAM policies to apply to other resources.
Currently, defining a policy through a datasource and referencing that policy
from another resource is the only way to apply an IAM policy to a resource.

**Note:** Several restrictions apply when setting IAM policies through this API.
See the [setIamPolicy docs](https://cloud.google.com/resource-manager/reference/rest/v1/projects/setIamPolicy)
for a list of these restrictions.

## Argument Reference

The following arguments are supported:

* `binding` (Required) - A nested configuration block (described below)
  defining a binding to be included in the policy document. Multiple
  `binding` arguments are supported.

Each document configuration must have one or more `binding` blocks, which
each accept the following arguments:

* `role` (Required) - The role/permission that will be granted to the members.
  See the [IAM Roles](https://cloud.google.com/compute/docs/access/iam) documentation for a complete list of roles.
  Note that custom roles must be of the format `[projects|organizations]/{parent-name}/roles/{role-name}`.

* `members` (Required) - An array of identites that will be granted the privilege in the `role`.
  Each entry can have one of the following values:
  * **allUsers**: A special identifier that represents anyone who is on the internet; with or without a Google account. It **can't** be used with the `google_project` resource.
  * **allAuthenticatedUsers**: A special identifier that represents anyone who is authenticated with a Google account or a service account. It **can't** be used with the `google_project` resource.
  * **user:{emailid}**: An email address that represents a specific Google account. For example, alice@gmail.com or joe@example.com.
  * **serviceAccount:{emailid}**: An email address that represents a service account. For example, my-other-app@appspot.gserviceaccount.com.
  * **group:{emailid}**: An email address that represents a Google group. For example, admins@example.com.
  * **domain:{domain}**: A G Suite domain (primary, instead of alias) name that represents all the users of that domain. For example, google.com or example.com.

## Attributes Reference

The following attribute is exported:

* `policy_data` - The above bindings serialized in a format suitable for
  referencing from a resource that supports IAM.
