---
layout: "google"
page_title: "Google: google_storage_object_iam_policy"
sidebar_current: "docs-google--storage-object-iam-policy"
description: |-
 Allows management of an IAM policy for a Google Storage Bucket Object.
---

# google\_storage\-_object\_iam\_policy

Allows creation and management of an IAM policy for an existing Google Storage Bucket Object.

## Example Usage

```hcl
resource "google_storage_bucket" "bucket" {
  display_name = "%s"
}

resource "google_storage_object" "object" {
  display_name = "%s"
	bucket = "${google_storage_bucket.bucket.name}"
}

resource "google_storage_object_iam_policy" "object" {
  object     = "${google_storage_object.object.name}"
  policy_data = "${data.google_iam_policy.admin.policy_data}"
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

* `object` - (Required) The Object ID.

* `policy_data` - (Required) The `google_iam_policy` data source that represents
    the IAM policy that will be applied to the project. The policy will be
    merged with any existing policy applied to the project.

    Changing this updates the policy.

    Deleting this removes the policy, but leaves the original object policy
    intact. If there are overlapping `binding` entries between the original
    object policy and the data source policy, they will be removed.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `etag` - (Computed) The etag of the object's IAM policy.

