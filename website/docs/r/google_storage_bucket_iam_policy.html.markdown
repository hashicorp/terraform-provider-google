---
layout: "google"
page_title: "Google: google_storage_bucket_iam_policy"
sidebar_current: "docs-google-storage-bucket-iam-policy"
description: |-
 Allows management of the IAM policy for a Google Cloud Storage buckets.
---

# google\_storage\_bucket\_iam\_policy

Allows creation and management of the IAM policy for an existing Google Cloud
Storage buckets.

## Example Usage

```hcl
resource "google_storage_bucket_iam_policy" "test" {
  bucket = "${google_storage_bucket.bucket.name}"
  policy_data = "${data.google_iam_policy.test.policy_data}"
}

resource "google_storage_bucket" "bucket" {
  display_name = "my_bucket"
}

data "google_iam_policy" "test" {
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

* `bucket` - (Required) The resource name of the bucket the policy is attached to. Its format is {bucket_id}.

* `policy_data` - (Required) The `google_iam_policy` data source that represents
    the IAM policy that will be applied to the bucket. This policy overrides any existing
    policy applied to the bucket.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `etag` - (Computed) The etag of the bucket's IAM policy. `etag` is used for optimistic concurrency control as a way to help prevent simultaneous updates of a policy from overwriting each other. 
