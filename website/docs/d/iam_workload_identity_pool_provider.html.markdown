---
subcategory: "Cloud IAM"
page_title: "Google: google_iam_workload_identity_pool_provider"
description: |-
  Get a IAM workload identity pool provider from Google Cloud
---

# google\_iam\_workload_\identity\_pool\_provider

Get a IAM workload identity provider from Google Cloud by its id.

~> **Warning:** This resource is in beta, and should be used with the terraform-provider-google-beta provider.
See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html) for more details on beta resources.

## Example Usage

```tf
data "google_iam_workload_identity_pool_provider" "foo" {
  workload_identity_pool_id          = "foo-pool"
  workload_identity_pool_provider_id = "bar-provider"
}
```

## Argument Reference

The following arguments are supported:

* `workload_identity_pool_id` - (Required) The id of the pool which is the
    final component of the pool resource name.
* `workload_identity_pool_provider_id` - (Required) The id of the provider which is the
    final component of the resource name.

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference
See [google_iam_workload_identity_pool_provider](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/iam_workload_identity_pool_provider) resource for details of all the available attributes.
