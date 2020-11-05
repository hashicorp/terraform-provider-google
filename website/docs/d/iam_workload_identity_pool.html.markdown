---
subcategory: "Cloud IAM"
layout: "google"
page_title: "Google: google_iam_workload_identity_pool"
sidebar_current: "docs-google-datasource-iam-workload-identity-pool"
description: |-
  Get a IAM workload identity pool from Google Cloud
---

# google\_iam\_workload_\identity\_pool

Get a IAM workload identity pool from Google Cloud by its id.

~> **Warning:** This resource is in beta, and should be used with the terraform-provider-google-beta provider.
See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html) for more details on beta resources.

## Example Usage

```tf
data "google_iam_workload_identity_pool" "foo" {
  workload_identity_pool_id = "foo-pool"
}
```

## Argument Reference

The following arguments are supported:

* `workload_identity_pool_id` - (Required) The id of the pool which is the
    final component of the resource name.

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference
See [google_iam_workload_identity_pool](https://www.terraform.io/docs/providers/google/r/iam_workload_identity_pool.html) resource for details of all the available attributes.
