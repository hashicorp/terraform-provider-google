---
layout: "google"
page_title: "Google: google_storage_bucket_iam_policy"
sidebar_current: "docs-google--storage-bucket-iam-policy"
description: |-
 Allows management of an IAM policy for a Google Storage Bucket.
---

# google\_storage\_bucket\_iam\_policy

Allows creation and management of an IAM policy for an existing Google Storage Bucket.

## Example Usage

```hcl
resource "google_storage_bucket" "bucket" {
  display_name = "%s"
}

resource "google_storage_bucket_iam_policy" "bucket" {
  bucket     = "${google_storage_bucket.bucket.name}"
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

* `bucket` - (Required) The Bucket ID.

* `policy_data` - (Required) The `google_iam_policy` data source that represents
    the IAM policy that will be applied to the project. The policy will be
    merged with any existing policy applied to the project.

    Changing this updates the policy.

    Deleting this removes the policy, but leaves the original bucket policy
    intact. If there are overlapping `binding` entries between the original
    bucket policy and the data source policy, they will be removed.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `etag` - (Computed) The etag of the bucket's IAM policy.

