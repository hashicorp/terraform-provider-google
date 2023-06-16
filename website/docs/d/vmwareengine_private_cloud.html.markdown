---
subcategory: "Cloud VMware Engine"
description: |-
  Get information about a private cloud.
---

# google\_vmwareengine\_private_cloud

Use this data source to get details about a private cloud resource.

~> **Warning:** This data source is in beta, and should be used with the terraform-provider-google-beta provider.
See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html) for more details on beta resources.

To get more information about private cloud, see:
* [API documentation](https://cloud.google.com/vmware-engine/docs/reference/rest/v1/projects.locations.privateClouds)

## Example Usage

```hcl
data "google_vmwareengine_private_cloud" "my_pc" {
  provider = google-beta
  name     = "my-pc"
  location = "us-central1-a"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the resource.
* `location` - (Required) Location of the resource.

- - -

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

See [google_vmwareengine_private_cloud](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/vmwareengine_private_cloud#attributes-reference) resource for details of all the available attributes.